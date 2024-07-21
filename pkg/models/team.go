package models

// Type representing teams in our table.
type Team interface {
	Name() string
	Funds() int
	Roster() map[PlayerType]map[Player]bool
	AddPlayer(p Player, li LineupInfo, bid int) error
	RmPlayer(p Player) error
}

func RosterSize(t Team) int {
	var count int
	for _, v := range t.Roster() {
		count += len(v)
	}
	return count
}

// Available funds for a team may contains some already spoken-for amounts, since there is a minimum
// number of players required for each team. For example, if a team has drafted 9 of 16 players, then
// bids must be made on the remaining 7 players, meaning 7 funds of the available (unspent) must be
// discounted.
func MaxBidValue(t Team, li LineupInfo) int {
	return t.Funds() - (li.PlayerSlots() - len(t.Roster()))
}
