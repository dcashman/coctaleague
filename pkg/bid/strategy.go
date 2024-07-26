package bid

import (
	"math"

	"dcashman.net/coctaleague/pkg/models"
)

type Style int
type ValueBasis int
type Preemptive int

const (
	// Place only the minimum bet on each bid required, with the bid determined by the points/price ratio
	Value Style = iota

	// Place only the minimum bid, with the bid determined by lowest absolute cost
	Minimum

	// Determine the best possible team with current funds and buy all players for cheapest amount
	FullTeam

	// Use predicted points from parse source, e.g. ESPN, for a player as a measure of the player's expected 'value'.
	// This should be considered the 'default', since the value is present in the original parsed db.
	Predicted ValueBasis = iota

	// Determine bid based on historical costs for this position, e.g. #1 RB typically goes for XX, #2 goes for...
	// The 'value' present for a position then should be the difference between the historical cost and the current rate.
	// Conceptually, if a player noramal goes for $X but is available for $(X - N), then its 'value' $N.
	Historical

	// NOT YET SUPPORTED: Allow users to input their own values
	Custom

	// Only spend the minimum required on the extra required players.
	OnePointMin Preemptive = iota

	// A slightly more aggressive strategy, instead of trying to win the race for the 'top 1 point player' and risk making a
	// bid at the bottom of a long run of unvalued players, try to win the '2 point' cutoff. This may be easier because
	// historically players haven't been fought over the 2pt threshold as much, so probability of gobbling this up is likely
	// higher, since not competing with other players trying to guess where the '1 point line' starts, and worst-case this
	// still ends up above the entire run of 1pt players.  The downside, of course, is dedicating double the points of potentially
	// disposable players.
	TwoPointMin
)

type Strategy struct {
	Style      Style
	Value      ValueBasis
	Preemptive Preemptive
}

type PositionDistribution struct {
	Start int
	Bench int
}

type TeamComposition map[models.PlayerType]PositionDistribution

// Current bidder returns multiple teams in the event that multiple teams have the same remaining funds.
func CurrentBidder(snapshot models.DraftSnapshot) []*models.Team {
	maxAvailable := struct {
		currMax int
		teams   []*models.Team
	}{0, nil}
	for _, t := range snapshot.Teams() {

		available := t.MaxBidValue(snapshot.LineupInfo())
		if available > maxAvailable.currMax {
			maxAvailable.currMax = available
			maxAvailable.teams = []*models.Team{t}
		} else if available == maxAvailable.currMax {
			maxAvailable.teams = append(maxAvailable.teams, t)
		}
	}
	return maxAvailable.teams
}

func RecommendBids(snapshot models.DraftSnapshot, team *models.Team, strategy Strategy) []models.Bid {
	// Get the preemptive bids out of the way first, since they don't affect our effective maximum bid
	// amount, and if using non-1-pt preemptive bids, may change whether or not we are required to make
	// a bid.
	bids := preemptiveBids(snapshot, team, strategy)

	// Calculate our next bid if required, for a 'real value' player.
	// TODO: implement bid calculation logic.

	return bids
}

// There is no 'right answer' in terms of the balance of starters vs. bench players and which positions.
// We choose the an approximately balanecd option where each starter has, if possible, a sub.
// TODO: add test for this function
// TODO: consider new type for team comp
func DesiredTeamComposition(snapshot models.DraftSnapshot, strategy Strategy) TeamComposition {
	lineupInfo := snapshot.LineupInfo()
	numBenchPlayers := lineupInfo.PlayerSlots() - lineupInfo.StarterSlots()

	// Determine how many starters we can choose positions for after the minimum has been allocated. While doing
	// this, also record the gap sizes for each position, which represents which positions have the most potential
	// in terms of adding players.
	numFlexStarters := lineupInfo.StarterSlots()
	gaps := make(map[models.PlayerType]int)
	roster := make(TeamComposition)

	// Go through each playerType
	for _, pt := range models.AllPlayerTypes {
		min := lineupInfo.PositionSlots()[pt].Min

		roster[pt] = PositionDistribution{
			Start: min,
			Bench: 0,
		}
		gaps[pt] = lineupInfo.PositionSlots()[pt].Max - min
		numFlexStarters -= min
	}

	// While we still have flexible starters, assign one to the position with the greatest avaliable number of starters.
	// The intuition here is that we will try to have more reserves for these positions to give us greater flexibility.
	for numFlexStarters > 0 {
		// Get the position with the highest remaining values
		maxGap := 0
		var maxPT models.PlayerType

		// Iterate through all positions using deterministic order. Eventually we should write a clear preference, but for
		// now just rely on our declared preferences, which happen to align with the order of positions.
		for _, pt := range models.AllPlayerTypes {
			if gaps[pt] > maxGap {
				maxGap = gaps[pt]
				maxPT = pt
			}
		}
		// Add a starter
		pd := roster[maxPT]
		pd.Start += 1
		roster[maxPT] = pd

		// Take away a gap count
		gaps[maxPT] = gaps[maxPT] - 1

		numFlexStarters -= 1
	}

	// Now do the same with the bench: assign bench players according to which positions have the most starters currently
	// without backups.
	// TODO: enable us to give priority to certain positions in the event of tie.  For example: after the first round, RBs and WRs may both have a gap
	// of 1 player, so which should get the next starter?
	for numBenchPlayers > 0 {
		maxGap := math.MinInt // We could find ourselves in a situation in which every position has more bench players than starters
		var maxPT models.PlayerType
		for k, v := range roster {

			// We never want backups for Defense or Kickers
			if k == models.D || k == models.K {
				continue
			}

			gap := v.Start - v.Bench
			if gap > maxGap {
				maxPT = k
				maxGap = gap
			}
		}
		pd := roster[maxPT]
		pd.Bench += 1
		roster[maxPT] = pd
		numBenchPlayers -= 1
	}

	return roster
}
