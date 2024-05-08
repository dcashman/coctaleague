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
	Creator User
	Year    int
	Funds   int
}

// Type representing teams in our table.
type Team struct {
	ID                  int
	Name                string
	Owner               User
	Season              Season
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
	Season              Season
	SpreadsheetPosition int
	espnId              int
	espnPredictedPoints int
	espnActualPoints    int
}

type Bid struct {
	ID        int
	Submitted time.Time
	Player    Player
	Bidder    Team
	Amount    int
}

type Draft struct {
	// Teams -> Player mapping

	// Player weight and value

	// maybe sort players in order of value, or just position, or cost

}
