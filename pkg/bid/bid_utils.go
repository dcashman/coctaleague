package bid

import (
	"sort"

	"dcashman.net/coctaleague/pkg/models"
)

// Check to see if any "basic" bids need to be cast before calculating. These bids are ones
// which we will always want to opportunistically make, such as making sure we have the best
// possible player already selected for any of the positions for which we only want to pay one
// point.
func preemptiveBids(snapshot models.DraftSnapshot, team models.Team, strategy Strategy, desiredComp TeamComposition) []models.Bid {
	bids := []models.Bid{}

	// Iterate through each position type
	for _, p := range models.AllPlayerTypes {
		bids = append(bids, minBids(snapshot, team, p, strategy, desiredComp)...)
	}

	return bids
}

// Filter out all players not on the provided team. These players are eligible for future bids.
func unownedPlayers(snapshot models.DraftSnapshot, team models.Team, pt models.PlayerType) []models.Player {
	var uop []models.Player
	for _, p := range snapshot.Players()[pt] {
		if _, found := team.Roster()[pt][p]; !found {
			uop = append(uop, p)
		}
	}
	return uop
}

// n-choose-k combinations: would be nice to not need to do this
func expandPlayerGroups(players []models.Player, size int) []*playerGroup {
	var pgs []*playerGroup
	if size < 1 {
		return pgs
	}

	if size == 1 {
		// This is equivalent to the existing players
		for _, p := range players {
			pa := []models.Player{p}
			pgs = append(pgs, newPlayerGroup(pa, p.Type()))
		}
		return pgs
	}

	// Otherwise, we need to build up our combinations recursively
	for i := 0; i <= len(players)-size; i++ {
		pgsSub := expandPlayerGroups(players[i+1:], size-1)
		for _, pgSub := range pgsSub {
			pg := pgSub.AddPlayer(players[i])
			pgs = append(pgs, pg)
		}
	}
	return pgs

}

func fullBids(snapshot models.DraftSnapshot, team models.Team, strategy Strategy, desiredComp TeamComposition, funds int) []models.Bid {
	bids := []models.Bid{}

	// We need to determine which players we still need after the preemptive bids, so we can determine which
	// to bid on.
	needed := fullBidsNeeded(snapshot, team, strategy, desiredComp)

	// Strip out the players that are already a part of the team from our candidates.
	strippedPlayers := make(map[models.PlayerType][]models.Player)
	for _, pt := range models.AllPlayerTypes {
		strippedPlayers[pt] = unownedPlayers(snapshot, team, pt)
	}

	// These needed slots need to be combined into grouped entries for the knapsack algorithm to work.
	// TODO: consider moving this logic into the knapsack model itself, calculating the tuple that represents each price-point best, rather than
	// pre-calculating all tuples.
	// TODO: parallelize this.
	var groupedPlayers [][]Item
	for _, pt := range models.AllPlayerTypes {
		if needed[pt] == 0 {
			// We don't need any of this category, so don't include it for knapsack calculations.
			continue
		}
		var items []Item

		pgs := expandPlayerGroups(strippedPlayers[pt], needed[pt])
		for _, pg := range pgs {
			// Convert to Item
			items = append(items, pg)
		}
		groupedPlayers = append(groupedPlayers, items)
	}
	pgs := MultiChoiceKnapsack(groupedPlayers, funds)

	for _, pg := range pgs {
		pg, ok := pg.(*playerGroup)
		if ok {
			bids = append(bids, pg.ToBids(team)...)
		} else {
			panic("MultiChoiceKnapsack returned Item which is not an underlying playerGroup")
		}
	}

	return bids
}

func minBids(snapshot models.DraftSnapshot, team models.Team, position models.PlayerType, strategy Strategy, desiredComp TeamComposition) []models.Bid {
	bids := []models.Bid{}

	//  Determine price to bid
	// TODO: BUG: need to adjust if team has less $ available
	minValue := minBidAmount(strategy)
	numBids := minBidQuantity(snapshot, team, position, strategy, desiredComp)

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
func minBidQuantity(snapshot models.DraftSnapshot, team models.Team, position models.PlayerType, strategy Strategy, desiredComp TeamComposition) int {
	var currentMinBids int
	for k := range team.Roster()[position] {
		if k.Bid().Amount <= minBidAmount(strategy) {
			// This player could represent one of our default bids
			currentMinBids++
		}
	}

	// We may have some open 'min bid slots', but we may also have all of them filled, and potentailly some non-min
	// entries going for cheap enough to qualify as min-bid.
	positionComp := desiredComp[position]
	minBidsForPosition := positionComp.Bench

	// Defense and kickers are both considered minbids even for their starters, so add those in too
	if position == models.D || position == models.K {
		minBidsForPosition += positionComp.Start
	}

	//TODO: BUG: need to ensure we have space on the team
	playersNeeded := minBidsForPosition - currentMinBids
	if playersNeeded > 0 {
		return playersNeeded
	}
	return 0
}

// Determine how many of the given position need to receive bids.
func normBidQuantity(snapshot models.DraftSnapshot, team models.Team, position models.PlayerType, strategy Strategy, desiredComp TeamComposition) int {
	// The number of normal bids is equal to totalPlayersNeeded - (minBidsNeeded + currPlayers)
	positionComp := desiredComp[position]
	total := positionComp.Bench + positionComp.Start
	mbNeeded := minBidQuantity(snapshot, team, position, strategy, desiredComp)
	curr := len(team.Roster()[position])

	return total - (mbNeeded + curr)
}

func fullBidsNeeded(snapshot models.DraftSnapshot, team models.Team, strategy Strategy, desiredComp TeamComposition) map[models.PlayerType]int {
	gaps := make(map[models.PlayerType]int)

	for _, pt := range models.AllPlayerTypes {
		gaps[pt] = normBidQuantity(snapshot, team, pt, strategy, desiredComp)
	}
	return gaps
}

func minBidAmount(strategy Strategy) int {
	if strategy.Preemptive == TwoPointMin {
		return 2
	} else {
		return 1
	}
}
