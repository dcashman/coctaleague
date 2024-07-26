package googlesheets

import "dcashman.net/coctaleague/pkg/models"

type Team struct {
}

func (*Team) StartingFunds() int {
	return 0
}

func (*Team) Teams() []*models.Team {
	return nil
}

func (*Team) LineupInfo() models.LineupInfo {
	return nil
}

func (*Team) Players() map[models.PlayerType][]*models.Player {
	return nil
}
