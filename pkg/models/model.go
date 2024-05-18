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

// Used for web application.
type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
	Active         bool
}

// Type representing teams in our table.
type Team struct {
	ID     int
	Name   string
	Funds  int
	Roster []*Player
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
	ID             int
	Name           string
	Organization   string
	Type           PlayerType
	espnId         int
	PredictedValue int
	Bid            *Bid
}

type Bid struct {
	ID        int
	Submitted time.Time
	Player    *Player
	Bidder    *Team
	Amount    int
}

type DraftStore interface {
	PlaceBid(bid Bid) error

	ParseDraft() (DraftSnapshot, error)
}

type DraftSnapshot struct {
	Teams   []*Team
	Players map[PlayerType][]*Player
}

func (d DraftSnapshot) TeamFromName(name string) *Team {
	return nil
}
