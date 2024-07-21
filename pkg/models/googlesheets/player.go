package googlesheets

import "dcashman.net/coctaleague/pkg/models"

type Player struct {
	name  string // Player's name
	org   string // Team player is a part of in the 'real world'
	pt    models.PlayerType
	value int        // Value we expect this player to produce. potentially used for bids.
	bid   models.Bid // Current 'winning bid' for the player.
	cell  string
}

func (p *Player) Name() string {
	return p.name
}

func (p *Player) Organization() string {
	return p.org
}

func (p *Player) Type() models.PlayerType {
	return p.pt
}

func (p *Player) PredictedValue() int {
	return p.value
}

func (p *Player) Bid() models.Bid {
	return p.bid
}

func (p *Player) UpdateBid(b models.Bid) error {
	// TODO PLACE BID for player HERE?
	p.bid = b
	return nil
}
