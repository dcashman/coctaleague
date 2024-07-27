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
	return 0
}

func (g *GoogleSheetsSnapshot) Teams() []models.Team {
	return nil
}

func (g *GoogleSheetsSnapshot) LineupInfo() models.LineupInfo {
	return nil
}

func (g *GoogleSheetsSnapshot) Players() map[models.PlayerType][]models.Player {
	return nil
}

func (g *GoogleSheetsSnapshot) Hotseat() string {
	return g.hotseat
}

func (g *GoogleSheetsSnapshot) Times() map[string]time.Duration {
	return g.times
}
