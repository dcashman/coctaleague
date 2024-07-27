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

// Used for web application, which doesn't really work yet.
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
type Player interface {
	Name() string         // Player's name
	Organization() string // Team player is a part of in the 'real world'
	Type() PlayerType
	PredictedValue() int // Value we expect this player to produce. potentially used for bids.
	Bid() Bid            // Current 'winning bid' for the player.
	UpdateBid(b Bid) error
}

type Bid struct {
	Player Player
	Bidder Team
	Amount int
}

func NewBid(p Player, t Team, a int) Bid {
	return Bid{
		Player: p,
		Bidder: t,
		Amount: a,
	}
}

type DraftStore interface {
	PlaceBid(bid Bid) error

	ParseDraft() (DraftSnapshot, error)
	WriteShotclock(d time.Duration, td time.Duration, h map[string]time.Duration) error
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
	Teams() []Team
	LineupInfo() LineupInfo
	Players() map[PlayerType][]Player
	Hotseat() string
	Times() map[string]time.Duration
}

// Convenience function that relies on the currently-true assumption that each team has a unique Name.
func TeamFromName(d DraftSnapshot, name string) (Team, error) {
	for _, t := range d.Teams() {
		if t.Name() == name {
			return t, nil
		}
	}
	return nil, ErrNoRecord
}
