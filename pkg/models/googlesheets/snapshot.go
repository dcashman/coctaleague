package googlesheets

import (
	"time"

	"dcashman.net/coctaleague/pkg/models"
)

type GoogleSheetsSnapshot struct {
	startingFunds int
	lineupInfo    LineupInfo
	teams         []Team
	players       map[models.PlayerType][]Player
	hotseat       string
	times         map[string]time.Duration
}

func (g *GoogleSheetsSnapshot) StartingFunds() int {
	return g.startingFunds
}

func (g *GoogleSheetsSnapshot) Teams() []models.Team {
	s := make([]models.Team, len(g.teams))
	for i := range g.teams {
		s[i] = &g.teams[i]
	}

	return s
}

func (g *GoogleSheetsSnapshot) LineupInfo() models.LineupInfo {
	return &g.lineupInfo
}

func (g *GoogleSheetsSnapshot) Players() map[models.PlayerType][]models.Player {
	retMap := make(map[models.PlayerType][]models.Player)
	for k, v := range g.players {
		s := make([]models.Player, len(v))
		for i := range v {
			s[i] = &v[i]
		}
		retMap[k] = s
	}

	return retMap
}

func (g *GoogleSheetsSnapshot) Hotseat() string {
	return g.hotseat
}

func (g *GoogleSheetsSnapshot) Times() map[string]time.Duration {
	return g.times
}
