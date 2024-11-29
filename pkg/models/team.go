package models

// Type representing teams in our table.
type Team interface {
	Name() string
	Funds() int
	// TODO: stop using Player as a map key (using interfaces as map keys is prone to error).
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

func rosterSize(t Team) int {
	sum := 0
	r := t.Roster()
	for _, pt := range AllPlayerTypes {
		sum += len(r[pt])
	}
	return sum
}
// Available funds for a team may contains some already spoken-for amounts, since there is a minimum
// number of players required for each team. For example, if a team has drafted 9 of 16 players, then
// bids must be made on the remaining 7 players, meaning 7 funds of the available (unspent) must be
// discounted.
func MaxBidValue(t Team, li LineupInfo) int {
	current := rosterSize(t)
	capacity := li.PlayerSlots()

	if current == capacity {
		return 0
	}
	if current > capacity {
		panic("team already has more players than capacity allows")
	}

	slack := t.Funds() - (capacity - current)

	// Add 1 to account for the min allocated to each player.
	return slack + 1
}
