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

type PlayerType int

// Keep in sync w/below
const (
	QB PlayerType = iota
	RB
	WR
	TE
	K
	D
	numPlayerTypes // Always keep at the end
)

var AllPlayerTypes []PlayerType

func init() {
	for pt := QB; pt < numPlayerTypes; pt++ {
		AllPlayerTypes = append(AllPlayerTypes, pt)
	}
}

// Type representing players in our table.
type Player struct {
	Name           string
	Organization   string
	Type           PlayerType
	PredictedValue int
	Bid            *Bid
}

type Bid struct {
	Player *Player
	Bidder *Team
	Amount int
}

type DraftStore interface {
	PlaceBid(bid Bid) error

	ParseDraft() (DraftSnapshot, error)
}

type LineupInfo interface {
	PlayerSlots() int
	StarterSlots() int
	PositionSlots() map[PlayerType]struct {
		Min int
		Max int
	}
}

type DraftSnapshot interface {
	StartingFunds() int
	Teams() []*Team
	LineupInfo() LineupInfo
	Players() map[PlayerType][]*Player
}

// Convenience function that relies on the currently-true assumption that each team has a unique Name.
func TeamFromName(d DraftSnapshot, name string) (*Team, error) {
	for _, t := range d.Teams() {
		if t.Name == name {
			return t, nil
		}
	}
	return nil, ErrNoRecord
}
