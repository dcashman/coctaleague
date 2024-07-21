package testutil

import (
	"dcashman.net/coctaleague/pkg/models"
)

var (
	TestQB01 = NewTestPlayer("a", "A", models.QB, 60)

	TestQB02 = NewTestPlayer("b", "B", models.QB, 54)

	TestQB03 = NewTestPlayer("c", "C", models.QB, 40)

	QB_PLAYERS = []models.Player{TestQB01, TestQB02, TestQB03}

	TestRB01 = NewTestPlayer("d", "A", models.RB, 54)

	TestRB02 = NewTestPlayer("e", "B", models.RB, 52)

	TestRB03 = NewTestPlayer("f", "C", models.RB, 50)

	TestRB04 = NewTestPlayer("g", "A", models.RB, 30)

	RB_PLAYERS = []models.Player{TestRB01, TestRB02, TestRB03, TestRB04}

	TestWR01 = NewTestPlayer("h", "A", models.WR, 40)

	TestWR02 = NewTestPlayer("i", "B", models.WR, 36)

	TestWR03 = NewTestPlayer("j", "C", models.WR, 35)

	TestWR04 = NewTestPlayer("k", "B", models.WR, 24)

	WR_PLAYERS = []models.Player{TestWR01, TestWR02, TestWR03, TestWR04}

	TestTE01 = NewTestPlayer("l", "A", models.TE, 32)

	TestTE02 = NewTestPlayer("m", "B", models.TE, 28)

	TestTE03 = NewTestPlayer("n", "C", models.TE, 20)

	TestTE04 = NewTestPlayer("o", "C", models.TE, 12)

	TE_PLAYERS = []models.Player{TestTE01, TestTE02, TestTE03, TestTE04}

	TestD01 = NewTestPlayer("p", "B", models.D, 32)

	TestD02 = NewTestPlayer("q", "C", models.D, 30)

	TestD03 = NewTestPlayer("r", "B", models.D, 29)

	D_PLAYERS = []models.Player{TestD01, TestD02, TestD03}

	TestK01 = NewTestPlayer("s", "B", models.K, 36)

	TestK02 = NewTestPlayer("t", "C", models.K, 35)

	TestK03 = NewTestPlayer("u", "A", models.K, 31)

	K_PLAYERS = []models.Player{TestK01, TestK02, TestK03}

	TestPlayerMap = map[models.PlayerType][]models.Player{
		models.QB: QB_PLAYERS,
		models.RB: RB_PLAYERS,
		models.WR: WR_PLAYERS,
		models.TE: TE_PLAYERS,
		models.D:  D_PLAYERS,
		models.K:  K_PLAYERS,
	}
)

type testPlayer struct {
	name  string // Player's name
	org   string // Team player is a part of in the 'real world'
	pt    models.PlayerType
	value int        // Value we expect this player to produce. potentially used for bids.
	bid   models.Bid // Current 'winning bid' for the player.
}

func (p *testPlayer) Name() string {
	return p.name
}

func (p *testPlayer) Organization() string {
	return p.org
}

func (p *testPlayer) Type() models.PlayerType {
	return p.pt
}

func (p *testPlayer) PredictedValue() int {
	return p.value
}

func (p *testPlayer) Bid() models.Bid {
	return p.bid
}

func (p *testPlayer) UpdateBid(b models.Bid) error {
	// TODO PLACE BID for player HERE?
	p.bid = b
	return nil
}

func NewTestPlayer(name string, org string, pt models.PlayerType, pv int) models.Player {
	return &testPlayer{
		name:  name,
		org:   org,
		pt:    pt,
		value: pv,
		bid:   models.Bid{},
	}
}
