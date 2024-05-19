package bid

import (
	"dcashman.net/coctaleague/pkg/models"
)

type Strategy int

const (
	// Place only the minimum bet on each bid required, with the bid determined by the points/price ratio
	Value Strategy = iota

	// Place only the minimum bid, with the bid determined by lowest absolute cost
	Minimum

	// Determine the best possible team with current funds and buy all players for cheapest amount
	FullTeam

	// Determine bid based on historical costs for this position, e.g. #1 RB typically goes for XX, #2 goes for...
	Historical
)

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
