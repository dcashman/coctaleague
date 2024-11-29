package bid

import "dcashman.net/coctaleague/pkg/models"

// Type to enable us to enumerate all player selections for a type and treat them as one element for
// consideration by our knapsack problem. This allows us to return the appropriate players for any
// given weight without having to worry about how to mix and match the different players.
type playerGroup struct {
	val     int
	weight  int
	players []models.Player
	pt      models.PlayerType
}

func newPlayerGroup(players []models.Player, pt models.PlayerType) *playerGroup {
	return &playerGroup{players: players, pt: pt}
}

func (pg *playerGroup) Value() int {
	// Lazy init
	if pg.val == 0 {
		var sum int
		for _, p := range pg.players {
			sum += p.PredictedValue()
		}
		// Prevent us from having to calculate this again
		pg.val = sum
	}
	return pg.val
}

func (pg *playerGroup) Weight() int {
	// Lazy init
	if pg.weight == 0 {
		var sum int
		for _, p := range pg.players {
			// Bid will always cost at least 1 pt to make
			sum += 1

			// If the player already has an owner with a winning bid, add that to our cost, since we have to
			// overcome it.
			sum += p.Bid().Amount
		}
		pg.weight = sum
	}
	return pg.weight
}

// Helper function which creates a new player group from an old one, using the same underlying interface refs.
// This is added due to the lazy-init feature which assumes that the contents of a playergroup don't change after
// creation.
func (pg *playerGroup) AddPlayer(p models.Player) *playerGroup {
	newPlayers := append(pg.players, p)
	return newPlayerGroup(newPlayers, pg.pt)
}

func (pg *playerGroup) ToBids(t models.Team) []models.Bid {
	var bids []models.Bid
	for _, p := range pg.players {
		bids = append(bids, models.Bid{Bidder: t, Player: p, Amount: p.Bid().Amount + 1})
	}
	return bids
}
