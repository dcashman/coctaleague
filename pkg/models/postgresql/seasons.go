package postgresql

import (
	"database/sql"
	"errors"

	"dcashman.net/coctaleague/pkg/models"
)

type SeasonModel struct {
	DB *sql.DB
}

// We'll use the Insert method to add a new record to the users table.
func (m *SeasonModel) Insert(year, funds int) error {

	stmt := `INSERT INTO users (year, funds) VALUES($1, $2)`

	// Use the Exec() method to insert the season details into the seasons table
	_, err := m.DB.Exec(stmt, year, funds)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return models.ErrDuplicateEmail
		default:
			return err
		}
	}

	return nil
}

func (m *SeasonModel) Get(id int) (*models.Season, error) {
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
