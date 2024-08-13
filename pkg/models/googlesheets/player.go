package googlesheets

import (
	"fmt"

	"dcashman.net/coctaleague/pkg/models"
	"google.golang.org/api/sheets/v4"
)

type Player struct {
	name  string // Player's name
	org   string // Team player is a part of in the 'real world'
	pt    models.PlayerType
	value int        // Value we expect this player to produce. potentially used for bids.
	bid   models.Bid // Current 'winning bid' for the player.
	cell  string
}

func (p *Player) Name() string {
	return p.name
}

func (p *Player) Organization() string {
	return p.org
}

func (p *Player) Type() models.PlayerType {
	return p.pt
}

func (p *Player) PredictedValue() int {
	return p.value
}

func (p *Player) Bid() models.Bid {
	return p.bid
}

func (p *Player) UpdateBid(b models.Bid) error {
	p.bid = b
	return nil
}

func NewPlayer(name string, org string, pt models.PlayerType, value int, cell string) Player {
	return Player{
		name:  name,
		org:   org,
		pt:    pt,
		value: value,
		cell:  cell,
	}
}

// Parsing functions
func parsePlayersPos(vr *sheets.ValueRange, teams []Team, row int, col int, pt models.PlayerType) ([]Player, error) {
	v := vr.Values
	var players []Player

	if !(row < len(v) && col < len(v[row])) {
		// We should have at least the first value for a position defined.
		return nil, fmt.Errorf("Could not parse players for position %v because sheets data is empty at cell %s.", pt, indicesToCellStr(row, col))
	}
	pn := interfaceToString(v[row][col])
	for pn != "" {
		// Get basic player info
		valCol := col + PLAYERS_VAL_OFFSET
		bidsStart := col + PLAYERS_BIDS_OFFSET
		var po string
		if pt == models.D {
			// Defense doesn't have an associated organizaiton because the player already represents the organization as a whole. Adjust
			// accordingly.
			po = pn

			// This also means that the value and bids columns are 1 closer than desired
			valCol--
			bidsStart--
		} else {
			po = interfaceToString(v[row][col+PLAYERS_ORG_OFFSET])
		}

		pv, err := interfaceToInt(v[row][valCol])
		if err != nil {
			return nil, fmt.Errorf("Value read for player %s, in cell: %s not an integer. %s\n", pn, indicesToCellStr(row, valCol), err.Error())
		}
		pc := indicesToCellStr(row, col)
		p := NewPlayer(pn, po, pt, pv, pc)

		// Get winning bid, if one exists
		maxBid := 0
		var maxTeam Team
		for i := 0; i < len(teams); i++ {
			bc := bidsStart + i

			bs := interfaceToString(v[row][bc])
			if bs != "" {
				bv, err := interfaceToInt(v[row][bc])
				if err != nil {
					return nil, fmt.Errorf("Bid read for player %s, in cell: %s not an integer. %s\n", pn, indicesToCellStr(row, bc), err.Error())
				}
				if bv > maxBid {
					maxBid = bv
					maxTeam = teams[i]
				}
			}
		}
		if maxBid > 0 {
			// Add player to team and players list with winning bid
			p.UpdateBid(models.NewBid(&p, &maxTeam, maxBid))
			_, ok := maxTeam.roster[pt]
			if !ok {
				maxTeam.roster[pt] = make(map[models.Player]bool)
			}
			maxTeam.roster[pt][&p] = true
		}

		// Player is ready, add to list
		players = append(players, p)

		// Move on to next player, if there is one
		row++

		if row < len(v) && col < len(v[row]) {
			pn = interfaceToString(v[row][col])
		} else {
			// Sheets has no more data in this dimension, must be empty
			pn = ""
		}
	}
	return players, nil
}

func parsePlayers(vr *sheets.ValueRange, teams []Team) (map[models.PlayerType][]Player, error) {
	var players = make(map[models.PlayerType][]Player)

	// For the players, we will injest by position.  Current version has the player position in the
	// same order as the declaration in models, but that could change in a future version.

	row, col, err := cellStrToIndices(PLAYERS_CELL)
	if err != nil {
		return nil, fmt.Errorf("unable to get starting indices for player list at starting cell %s in spreasheet. %s\n", PLAYERS_CELL, err.Error())
	}

	posOffset := len(teams) + PLAYERS_PADDING_COLS + 3 // + 3 accounts for the player name, org and value columns
	for i, pt := range models.AllPlayerTypes {
		players[pt], err = parsePlayersPos(vr, teams, row, col+(i*posOffset), pt)
		if err != nil {
			return nil, fmt.Errorf("Unable to parse players for position %v in spreasheet. %s\n", pt, err.Error())
		}
	}

	return players, nil
}
