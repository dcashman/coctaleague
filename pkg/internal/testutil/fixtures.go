package testutil

import (
	"dcashman.net/coctaleague/pkg/models"
)

const (
	DEFAULT_STARTING_FUNDS = 100
)

var (
	TestPositionSlots = map[models.PlayerType]struct {
		Min int
		Max int
	}{
		models.QB: {Min: 1, Max: 1},
		models.RB: {Min: 1, Max: 2},
		models.WR: {Min: 2, Max: 4},
		models.TE: {Min: 1, Max: 2},
		models.D:  {Min: 1, Max: 1},
		models.K:  {Min: 1, Max: 1},
	}
	TestLineupInfo = testLineupInfo{
		playerSlots:   16,
		starterSlots:  9,
		positionSlots: TestPositionSlots,
	}

	TestDraftSnapshot = testDraftSnapshot{
		startingFunds: DEFAULT_STARTING_FUNDS,
		teams:         TestTeams,
		lineupInfo:    TestLineupInfo,
		players:       TestPlayerMap,
	}
)

type testLineupInfo struct {
	playerSlots   int
	starterSlots  int
	positionSlots map[models.PlayerType]struct {
		Min int
		Max int
	}
}

func (t testLineupInfo) PlayerSlots() int {
	return t.playerSlots
}

func (t testLineupInfo) StarterSlots() int {
	return t.starterSlots
}

func (t testLineupInfo) PositionSlots() map[models.PlayerType]struct {
	Min int
	Max int
} {
	return t.positionSlots
}

func NewEmptyDraftSnapshot(startingFunds int, lineupInfo models.LineupInfo) models.DraftSnapshot {
	return testDraftSnapshot{
		startingFunds: startingFunds,
		lineupInfo:    lineupInfo,
		teams:         []*models.Team{},
		players:       make(map[models.PlayerType][]*models.Player),
	}
}

type testDraftSnapshot struct {
	startingFunds int
	teams         []*models.Team
	lineupInfo    models.LineupInfo
	players       map[models.PlayerType][]*models.Player
}

func (t testDraftSnapshot) StartingFunds() int {
	return t.startingFunds
}

func (t testDraftSnapshot) Teams() []*models.Team {
	return t.teams
}

func (t testDraftSnapshot) LineupInfo() models.LineupInfo {
	return t.lineupInfo
}

func (t testDraftSnapshot) Players() map[models.PlayerType][]*models.Player {
	return t.players
}
