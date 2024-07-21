package testutil

import (
	"fmt"

	"dcashman.net/coctaleague/pkg/models"
)

var (
	TestTeam01 = NewEmptyTestTeam("a", DEFAULT_STARTING_FUNDS)

	TestTeam02 = NewEmptyTestTeam("b", DEFAULT_STARTING_FUNDS)

	TestTeams = []models.Team{TestTeam01, TestTeam02}
)

func NewEmptyTestTeam(name string, funds int) models.Team {
	return &testTeam{
		name:   name,
		funds:  funds,
		roster: make(map[models.PlayerType]map[models.Player]bool),
	}
}

func NewTestTeam(name string, funds int, players []models.Player) {
	team := NewEmptyTestTeam(name, funds)
	for _, p := range players {
		team.AddPlayer(p, TestLineupInfo, p.Bid().Amount+1)
	}
}

type testTeam struct {
	name   string
	funds  int
	roster map[models.PlayerType]map[models.Player]bool
}

func (t *testTeam) Name() string {
	return t.name
}

func (t *testTeam) Funds() int {
	return t.funds
}

func (t *testTeam) Roster() map[models.PlayerType]map[models.Player]bool {
	return t.roster
}

func (*testTeam) Players() map[models.PlayerType][]*models.Player {
	return nil
}

// TODO: this is shared code between test and googlesheets implementations.  This is a bad code smell.  Remove.
func (t *testTeam) AddPlayer(p models.Player, li models.LineupInfo, bid int) error {
	if bid > models.MaxBidValue(t, li) {
		return fmt.Errorf("team %s has insufficient funds (%d with %d players registered) to make bid (%d)", t.name, t.funds, len(t.Roster()), bid)
	}
	if p.Bid().Amount >= bid {
		return fmt.Errorf("bid of %d is inadequate to purchase player with existing bid of %d ", bid, p.Bid().Amount)
	}
	if models.RosterSize(t) == li.PlayerSlots() {
		return fmt.Errorf("team %s already has a full roster of %d players", t.name, li.PlayerSlots())
	}

	// Linked list logic, remove from old team if exists.
	// TODO: This is not currently necessary, since we don't link to backing data store, but eventually it may be.
	if p.Bid().Bidder != nil {

		// Must remove player from previous team before we insert our new bid, since it uses the existing bid to credit the
		// previous owner.
		p.Bid().Bidder.RmPlayer(p)
	}

	// Finally, add the player to our team, first by updating the 'winning bid on the player', then by adding it to the roster.
	p.UpdateBid(models.NewBid(p, t, bid))
	t.funds -= p.Bid().Amount
	t.roster[p.Type()][p] = true
	return nil
}

func (t *testTeam) RmPlayer(p models.Player) error {
	if _, ok := t.Roster()[p.Type()][p]; ok {
		delete(t.Roster()[p.Type()], p)

		// Must have this team's bid available before removal.
		t.funds += p.Bid().Amount
	}
	return nil
}
