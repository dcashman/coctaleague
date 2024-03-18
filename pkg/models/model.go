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
	ID      int
	Name    string
	Creator int
	Year    int
	Funds   int
}

// Type representing teams in our table.
type Team struct {
	ID                  int
	Name                string
	Owner               int
	Season              int
	SpreadsheetPosition int
}

type PlayerType int

const (
	QB PlayerType = iota
	RB
	WR
	TE
	K
	D
)

// Type representing players in our table.
type Player struct {
	ID                  int
	Name                string
	Organization        string
	Type                PlayerType
	Season              int
	SpreadsheetPosition int
	espnId              int
	espnPredictedPoints int
	espnActualPoints    int
}

type Bid struct {
	ID        int
	Submitted time.Time
	Player    int
	Bidder    int
	Amount    int
}
