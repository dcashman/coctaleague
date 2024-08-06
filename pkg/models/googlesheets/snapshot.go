package googlesheets

import (
	"log"
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
	for i, _ := range g.teams {
		s[i] = &g.teams[i]
	}
	log.Printf("Name of team in inter: %s vs. raw: %s\n", s[0].Name(), g.teams[0].name)

	return s
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
