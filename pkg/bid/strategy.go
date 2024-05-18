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

func CurrentBidder(snapshot models.DraftSnapshot) *models.Team {
	return nil
}

func RecommendBids(snapshot models.DraftSnapshot, team *models.Team, strategy Strategy) []models.Bid {
	return nil
}
