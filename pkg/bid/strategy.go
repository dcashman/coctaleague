package bid

import (
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

// Current bidder returns multiple teams in the event that multiple teams have the same remaining funds.
func CurrentBidder(snapshot models.DraftSnapshot) []*models.Team {
	maxAvailable := struct {
		currMax int
		teams   []*models.Team
	}{0, nil}
	for _, t := range snapshot.Teams {
		// Available funds for a team may contains some already spoken-for amounts, since there is a minimum
		// number of players required for each team. For example, if a team has drafted 9 of 16 players, then
		// bids must be made on the remaining 7 players, meaning 7 funds of the available (unspent) must be
		// discounted.
		available := t.Funds - (snapshot.LineupInfo.PlayerSlots() - len(t.Roster))
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
	return nil
}
