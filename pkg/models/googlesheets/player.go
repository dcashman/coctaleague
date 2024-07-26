package googlesheets

import "dcashman.net/coctaleague/pkg/models"

type GoogleSheetsSnapshot struct {
}

func (*GoogleSheetsSnapshot) StartingFunds() int {
	return 0
}

func (*GoogleSheetsSnapshot) Teams() []*models.Team {
	return nil
}

func (*GoogleSheetsSnapshot) LineupInfo() models.LineupInfo {
	return nil
}

func (*GoogleSheetsSnapshot) Players() map[models.PlayerType][]*models.Player {
	return nil
}
