package googlesheets

import (
	"fmt"
	"log"
	"sort"
	"time"

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

func sheetsRange(title string, bounds string) string {
	return fmt.Sprintf("%s!%s", title, bounds)
}

func parseTimes(vr *sheets.ValueRange) (map[string]time.Duration, error) {
	// Times are hard-coded to be B46 - B60, and our values range starts at A1
	// We start i at 45 because slices are 0-indexed but gsheets are not.
	// We go row-by-row, assigning the key as the name (column A, index 0) with
	// the value time (column B, index 1)
	res := make(map[string]time.Duration)
	for i := 45; i < 60; i++ {
		tdString := interfaceToString(vr.Values[i][1])
		if tdString == "" {
			tdString = "0s"
		}
		td, err := time.ParseDuration(tdString)
		if err != nil {
			log.Fatalf("Unable to parse shotclock times from sheet: %v", err)
			return nil, err
		}
		res[interfaceToString(vr.Values[i][0])] = td
	}
	return res, nil
}

func interfaceToString(data interface{}) string {
	if str, ok := data.(string); ok {
		return str
	}
	return fmt.Sprintf("%v", data)
}

func (g *GoogleSheetsDb) ParseDraft() (models.DraftSnapshot, error) {
	// Use the underlying sheet to populate a SheetDraft type, which implements the DraftSnapshot interface.
	readRange := sheetsRange(g.title, g.bounds)
	resp, err := g.service.Spreadsheets.Values.Get(g.id, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
		return nil, err
	}

	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
		return nil, fmt.Errorf("couldn't get spreadsheet from google")
	}

	ss := GoogleSheetsSnapshot{}

	ss.hotseat = interfaceToString(resp.Values[0][0])

	ss.times, err = parseTimes(resp)
	if err != nil {
		return nil, err
	}
	// Get teams slice
	// get players for position:

	// Get number of teams

	// Get
	//fmt.Printf("first row: %s\n", resp.Values[4])
	/*for i, row := range resp.Values {
		fmt.Printf("%d, %s\n", i, row[0])
	}*/

	return &ss, nil
}

func (g *GoogleSheetsDb) PlaceBid(models.Bid) error {
	// Use the underlying sheet to place a bid, returning an error if it couldn't be placed.
	return nil
}

func timesToWriteRange(h map[string]time.Duration) *sheets.ValueRange {
	// First let's sort them so that we always get the same order
	var keys []string
	var values [][]interface{}
	for key := range h {
		keys = append(keys, key)
	}
	// Sort the keys
	sort.Strings(keys)

	for _, k := range keys {
		values = append(values, []interface{}{k, h[k].String()})
	}
	return &sheets.ValueRange{
		Values: values,
	}
}

// Special-cased hard-coding for now (during draf on draft-day)
func (g *GoogleSheetsDb) WriteShotclock(d time.Duration, td time.Duration, h map[string]time.Duration) error {
	writeRange := sheetsRange(g.title, "B1:C1")
	valueRange := &sheets.ValueRange{
		Values: [][]interface{}{
			{d.String(), td.String()},
		},
	}
	_, err := g.service.Spreadsheets.Values.Update(g.id, writeRange, valueRange).ValueInputOption("RAW").Do()
	if err != nil {
		log.Fatalf("Unable to write data to sheet: %v", err)
		return err
	}

	// Write the histogram stuff
	if len(valueRange.Values) > 14 {
		log.Fatalf("Too many entries in our times to write section: %d\n", len(valueRange.Values))
		return fmt.Errorf("Received histogram range with extra values. %v", valueRange.Values)
	}
	valueRange = timesToWriteRange(h)
	writeRange = sheetsRange(g.title, "A46:B60")
	_, err = g.service.Spreadsheets.Values.Update(g.id, writeRange, valueRange).ValueInputOption("RAW").Do()
	if err != nil {
		log.Fatalf("Unable to write data to sheet: %v", err)
		return err
	}
	return nil
}
