package postgresql

import (
	"database/sql"
	"errors"

	"dcashman.net/coctaleague/pkg/models"
)

type SeasonModel struct {
	DB *sql.DB
}

// We'll use the Insert method to add a new record to the seasons table.
func (m *SeasonModel) Insert(name string, creator, year, funds int) error {

	stmt := `INSERT INTO seasons (name, creator, year, funds) VALUES($1, $2, $3, $4)`

	// Use the Exec() method to insert the season details into the seasons table
	_, err := m.DB.Exec(stmt, name, creator, year, funds)
	if err != nil {
		return err
	}

	return nil
}

func (m *SeasonModel) GetFromId(id int) (*models.Season, error) {
	s := &models.Season{}

	stmt := `SELECT id, name, creator, year, funds, FROM seasons WHERE id = $1`
	err := m.DB.QueryRow(stmt, id).Scan(&s.ID, &s.Name, &s.Creator, &s.Year, &s.Funds)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

func (m *SeasonModel) Get(name string, owner, year int) (*models.Season, error) {
	s := &models.Season{}

	stmt := `SELECT id, name, creator, year, funds, FROM seasons WHERE name = $1 AND owner = $2 AND year = $3`
	err := m.DB.QueryRow(stmt, name, owner, year).Scan(&s.ID, &s.Name, &s.Creator, &s.Year, &s.Funds)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}
