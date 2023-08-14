package postgresql

import (
	"database/sql"
	"errors"

	"dcashman.net/coctaleague/pkg/models"
)

type TeamModel struct {
	DB *sql.DB
}

// We'll use the Insert method to add a new record to the teams table.
func (m *TeamModel) Insert(name string, owner, season, spreadsheet_position int) error {

	stmt := `INSERT INTO teams (name, owner, season, spreadsheet_position) VALUES($1, $2, $3, $4)`

	// Use the Exec() method to insert the Team details into the Teams table
	_, err := m.DB.Exec(stmt, name, owner, season, spreadsheet_position)
	if err != nil {
		return err
	}

	return nil
}

func (m *TeamModel) Get(id int) (*models.Team, error) {
	t := &models.Team{}

	stmt := `SELECT id, name, owner, season, spreadsheet_position FROM Teams WHERE id = $1`
	err := m.DB.QueryRow(stmt, id).Scan(&t.ID, &t.Name, &t.Owner, &t.Season, &t.SpreadsheetPosition)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return t, nil
}
