package models

import (
	"fmt"
	"time"
)

// Type representing teams in our table.
type Team struct {
	ID     int
	Name   string
	Funds  int
	Roster map[PlayerType]map[*Player]bool
}

// TODO: This does not currently change the underlying data store, but in the future it should, as such it is currently useful only for
// initialization and tests.
func (t *Team) AddPlayer(p *Player, li LineupInfo, bid int) error {
	if bid > t.MaxBidValue(li) {
		return fmt.Errorf("team %s has insufficient funds (%d with %d players registered) to make bid (%d)", t.Name, t.Funds, len(t.Roster), bid)
	}
	if p.Bid.Amount >= bid {
		return fmt.Errorf("bid of %d is inadequate to purchase player with existing bid of %d ", bid, p.Bid.Amount)
	}
	if t.RosterSize() == li.PlayerSlots() {
		return fmt.Errorf("team %s already has a full roster of %d players", t.Name, li.PlayerSlots())
	}

	// Linked list logic, remove from old team if exists.
	// TODO: This is not currently necessary, since we don't link to backing data store, but eventually it may be.
	if p.Bid != nil {

		// Must remove player from previous team before we insert our new bid, since it uses the existing bid to credit the
		// previous owner.
		p.Bid.Bidder.RmPlayer(p)
	}

	// Finally, add the player to our team, first by updating the 'winning bid on the player', then by adding it to the roster.
	p.Bid = &Bid{
		ID:        0,
		Submitted: time.Now(),
		Player:    p,
		Bidder:    t,
		Amount:    bid,
	}
	t.Funds -= p.Bid.Amount
	t.Roster[p.Type][p] = true
	return nil
}

func (t *Team) RmPlayer(p *Player) error {
	if _, ok := t.Roster[p.Type][p]; ok {
		delete(t.Roster[p.Type], p)

		// Must have this team's bid available before removal.
		t.Funds += p.Bid.Amount
	}
	return nil
}

func (t *Team) RosterSize() int {
	var count int
	for _, v := range t.Roster {
		count += len(v)
	}
	return count
}

// Available funds for a team may contains some already spoken-for amounts, since there is a minimum
// number of players required for each team. For example, if a team has drafted 9 of 16 players, then
// bids must be made on the remaining 7 players, meaning 7 funds of the available (unspent) must be
// discounted.
func (t *Team) MaxBidValue(li LineupInfo) int {
	return t.Funds - (li.PlayerSlots() - len(t.Roster))
}
