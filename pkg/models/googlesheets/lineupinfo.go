package googlesheets

import "dcashman.net/coctaleague/pkg/models"

type LineupInfo struct {
	playerSlots   int
	starterSlots  int
	positionSlots map[models.PlayerType]struct {
		Min int
		Max int
	}
}

// Right now we just hard-code the values.  Eventually we can parse the spreadsheet for this, as it
// does change.
var (
	PositionSlots = map[models.PlayerType]struct {
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
)

func NewLineupInfo() LineupInfo {
	return LineupInfo{
		playerSlots:   16,
		starterSlots:  9,
		positionSlots: PositionSlots,
	}
}

func (li *LineupInfo) PlayerSlots() int {
	return li.playerSlots
}

func (li *LineupInfo) StarterSlots() int {
	return li.starterSlots
}

func (li *LineupInfo) PositionSlots() map[models.PlayerType]struct {
	Min int
	Max int
} {
	return li.positionSlots
}
