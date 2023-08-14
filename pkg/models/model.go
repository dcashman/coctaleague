package models

import (
	"errors"
	"time"
)

var (
	ErrNoRecord           = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("models: duplicate email")
)

// Type representing users in our table.
type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
	Active         bool
}

// Type representing seasons in our table.
type Season struct {
	ID    int
	Year  int
	Funds int
}

// Type representing teams in our table.
type Team struct {
	ID                  int
	Name                string
	Owner               int
	Season              int
	SpreadsheetPosition int
}
