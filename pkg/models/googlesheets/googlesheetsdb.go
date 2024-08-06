package googlesheets

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"time"
	"unicode"

	"google.golang.org/api/sheets/v4"

	"dcashman.net/coctaleague/pkg/models"
)

const (
	STARTING_FUNDS_CELL = "C2"
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

func interfaceToString(data interface{}) string {
	if str, ok := data.(string); ok {
		return str
	}
	return fmt.Sprintf("%v", data)
}

func interfaceToInt(data interface{}) (int, error) {
	s := interfaceToString(data)
	return strconv.Atoi(s)
}

func indicesToCellStr(row int, col int) string {
	// Get the Column
	letter := ""
	for col >= 0 {
		letter = string(rune('A'+col%26)) + letter
		col = col/26 - 1
	}

	return fmt.Sprintf("%s%d", letter, row+1)
}

// cellStrToIndices converts a cell string (e.g., B1) back to row and column indices.
func cellStrToIndices(cell string) (int, int, error) {
	var colStr string
	var rowStr string

	// Separate the column letters and row numbers
	for i, char := range cell {
		if unicode.IsDigit(char) {
			colStr = cell[:i]
			rowStr = cell[i:]
			break
		}
	}

	// Convert column string to index
	col := 0
	for i := 0; i < len(colStr); i++ {
		col = col*26 + int(colStr[i]-'A'+1)
	}
	col-- // Adjust for zero-based indexing

	// Convert row string to index
	row, err := strconv.Atoi(rowStr)
	if err != nil {
		return 0, 0, err
	}
	row-- // Adjust for zero-based indexing

	return row, col, nil
}

func parseStartingFunds(vr *sheets.ValueRange) (int, error) {
	// Current spreadsheet version has funds listed in cell C2
	i, j, err := cellStrToIndices(STARTING_FUNDS_CELL)
	if err != nil {
		return 0, err
	}
	sf, err := interfaceToInt(vr.Values[i][j])
	if err != nil {
		fmt.Errorf("Could not read in starting funds for league from cell %s. %s\n", STARTING_FUNDS_CELL, err.Error())
	}
	return sf, nil
}

func (g *GoogleSheetsDb) ParseDraft(numMembers int) (models.DraftSnapshot, error) {
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

	// Parse the core draft info
	// 0) Get our starting funds and lineup info
	ss.startingFunds, err = parseStartingFunds(resp)
	if err != nil {
		return nil, err
	}

	// Lineup info is hard-coded for now.  Parsing this from the sheet would likely be brittle and not worth
	// the effort.
	ss.lineupInfo = NewLineupInfo()

	// 1) Get the teams
	ss.teams, err = parseTeams(resp, numMembers, ss.startingFunds)
	if err != nil {
		return nil, err
	}

	// 2) Get the players
	ss.players = make(map[models.PlayerType][]Player)

	// extra logic just for shot-clock during draft.
	ss.hotseat = interfaceToString(resp.Values[0][0])
	ss.times, err = parseTimes(resp)
	if err != nil {
		return nil, err
	}

	return &ss, nil
}

func (g *GoogleSheetsDb) PlaceBid(models.Bid) error {
	// Use the underlying sheet to place a bid, returning an error if it couldn't be placed.
	return nil
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
