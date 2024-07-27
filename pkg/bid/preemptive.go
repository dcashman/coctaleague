package bid

import (
	"sort"

	"dcashman.net/coctaleague/pkg/models"
)

// Check to see if any "basic" bids need to be cast before calculating. These bids are ones
// which we will always want to opportunistically make, such as making sure we have the best
// possible player already selected for any of the positions for which we only want to pay one
// point.
func preemptiveBids(snapshot models.DraftSnapshot, team models.Team, strategy Strategy) []models.Bid {
	bids := []models.Bid{}

	// Iterate through each position type
	for _, p := range models.AllPlayerTypes {
		bids = append(bids, minBids(snapshot, team, p, strategy)...)
	}

	return bids
}

func minBids(snapshot models.DraftSnapshot, team models.Team, position models.PlayerType, strategy Strategy) []models.Bid {
	bids := []models.Bid{}

	//  Determine price to bid
	minValue := minBidAmount(strategy)
	numBids := minBidQuantity(snapshot, team, position, strategy)

	if strategy.Value != Predicted {
		panic("Unsupported value basis for players")
	}

	// Get the highest valued player available for that price
	allPlayers := snapshot.Players()
	players := allPlayers[position]
	sort.Slice(players, func(i, j int) bool {
		// We want to sort by greatest value first, not lowest
		return players[i].PredictedValue() > players[j].PredictedValue()
	})

	i := 0
	for numBids > 0 && i < len(players) {
		if players[i].Bid().Amount < minValue && players[i].Bid().Bidder != team {
			bids = append(bids, models.Bid{Bidder: team, Player: players[i], Amount: players[i].Bid().Amount + 1})
			numBids--
		}
		i++
	}

	return bids
}

// Determine how many of the given position need to receive bids.
func minBidQuantity(snapshot models.DraftSnapshot, team models.Team, position models.PlayerType, strategy Strategy) int {
	var currentMinBids int
	for k := range team.Roster()[position] {
		if k.Bid().Amount <= minBidAmount(strategy) {
			// This player could represent one of our default bids
			currentMinBids++
		}
	}

	// We may have some open 'min bid slots', but we may also have all of them filled, and potentailly some non-min
	// entries going for cheap enough to qualify as min-bid.
	positionComp := DesiredTeamComposition(snapshot, strategy)[position]
	minBidsForPosition := positionComp.Bench

	// Defense and kickers are both considered minbids even for their starters, so add those in too
	if position == models.D || position == models.K {
		minBidsForPosition += positionComp.Start
	}
	playersNeeded := minBidsForPosition - currentMinBids
	if playersNeeded > 0 {
		return playersNeeded
	}
	return 0
}

func minBidAmount(strategy Strategy) int {
	if strategy.Preemptive == TwoPointMin {
		return 2
	} else {
		return 1
	}
}
