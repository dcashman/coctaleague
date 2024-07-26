package googlesheets

import (
	"fmt"
	"log"

	"google.golang.org/api/sheets/v4"

	"dcashman.net/coctaleague/pkg/models"
)

type GoogleSheetsDb struct {
	bounds  string // Read range
	id      string // Sheet id
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
	resp, err := g.service.Spreadsheets.Values.Get(g.id, g.bounds).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
		return nil, err
	}

	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
	} else {
		// Get number of teams

		// Get
		fmt.Printf("first row: %s\n", resp.Values[4])
		for i, row := range resp.Values {
			fmt.Printf("%d, %s\n", i, row[0])
		}
	}
	return nil, nil
}

func (g *GoogleSheetsDb) PlaceBid(models.Bid) error {
	// Use the underlying sheet to place a bid, returning an error if it couldn't be placed.
	return nil
}
