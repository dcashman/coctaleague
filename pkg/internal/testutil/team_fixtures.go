package testutil

import (
	"dcashman.net/coctaleague/pkg/models"
)

var (
	TestTeam01 = NewEmptyTestTeam("a", DEFAULT_STARTING_FUNDS)

	TestTeam02 = NewEmptyTestTeam("b", DEFAULT_STARTING_FUNDS)

	TestTeams = []*models.Team{&TestTeam01, &TestTeam02}
)

func NewEmptyTestTeam(name string, funds int) models.Team {
	return models.Team{Name: name, Funds: funds, Roster: make(map[models.PlayerType]map[*models.Player]bool)}
}

func NewTestTeam(name string, funds int, players []*models.Player) {
	team := NewEmptyTestTeam(name, funds)
	for _, p := range players {
		team.AddPlayer(p, TestLineupInfo, p.Bid.Amount+1)
	}
}
