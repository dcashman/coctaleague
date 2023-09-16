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
func (m *SeasonModel) Insert(year, funds int) error {

	stmt := `INSERT INTO seasons (year, funds) VALUES($1, $2)`

	// Use the Exec() method to insert the season details into the seasons table
	_, err := m.DB.Exec(stmt, year, funds)
	if err != nil {
		return err
	}

	return nil
}

func (m *SeasonModel) GetId(id int) (*models.Season, error) {
	s := &models.Season{}

	stmt := `SELECT id, year, funds, FROM seasons WHERE id = $1`
	err := m.DB.QueryRow(stmt, id).Scan(&s.ID, &s.Year, &s.Funds)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

func (m *SeasonModel) Get(year int) (*models.Season, error) {
	s := &models.Season{}

	stmt := `SELECT id, year, funds, FROM seasons WHERE year = $1`
	err := m.DB.QueryRow(stmt, year).Scan(&s.ID, &s.Year, &s.Funds)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}
