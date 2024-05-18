package googlesheets

import (
	"google.golang.org/api/sheets/v4"

	"dcashman.net/coctaleague/pkg/models"
)

type Player struct {
	p    models.Player
	cell string
}

type Team struct {
	t    models.Team
	cell string
}

type GoogleSheetsDb struct {
	bounds  string
	id      string
	service *sheets.Service
	title   string
}

func NewGoogleSheetsDb(bounds string, id string, service *sheets.Service, title string) *GoogleSheetsDb {
	return &GoogleSheetsDb{
		bounds:  bounds,
		id:      id,
		service: service,
		title:   title,
	}
}

func (g *GoogleSheetsDb) ParseDraft() (models.DraftSnapshot, error) {
	// Use the underlying sheet to populate a SheetDraft type, which implements the DraftSnapshot interface.
	return models.DraftSnapshot{}, nil
}

func (g *GoogleSheetsDb) PlaceBid(models.Bid) error {
	// Use the underlying sheet to place a bid, returning an error if it couldn't be placed.
	return nil
}
