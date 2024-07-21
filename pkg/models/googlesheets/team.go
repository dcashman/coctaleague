package googlesheets

import (
	"fmt"

	"dcashman.net/coctaleague/pkg/models"
)

type Team struct {
	name   string
	funds  int
	roster map[models.PlayerType]map[models.Player]bool
	cell   string
}

func (t *Team) Name() string {
	return t.name
}

func (t *Team) Funds() int {
	return t.funds
}

func (t *Team) Roster() map[models.PlayerType]map[models.Player]bool {
	return t.roster
}

func (t *Team) Players() map[models.PlayerType][]*models.Player {
	return nil
}

// TODO: This does not currently change the underlying data store, but in the future it should, as such it is currently useful only for
// initialization and tests.
func (t *Team) AddPlayer(p models.Player, li models.LineupInfo, bid int) error {
	if bid > models.MaxBidValue(t, li) {
		return fmt.Errorf("team %s has insufficient funds (%d with %d players registered) to make bid (%d)", t.Name, t.Funds, len(t.Roster()), bid)
	}
	if p.Bid().Amount >= bid {
		return fmt.Errorf("bid of %d is inadequate to purchase player with existing bid of %d ", bid, p.Bid().Amount)
	}
	if models.RosterSize(t) == li.PlayerSlots() {
		return fmt.Errorf("team %s already has a full roster of %d players", t.Name, li.PlayerSlots())
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

func (t *Team) RmPlayer(p models.Player) error {
	if _, ok := t.Roster()[p.Type()][p]; ok {
		delete(t.Roster()[p.Type()], p)

		// Must have this team's bid available before removal.
		t.funds += p.Bid().Amount
	}
	return nil
}
