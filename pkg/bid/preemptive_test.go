package bid

import (
	"testing"

	"dcashman.net/coctaleague/pkg/internal/testutil"
	"dcashman.net/coctaleague/pkg/models"
)

type minBidQuantityInput struct {
	ds models.DraftSnapshot
	t  models.Team
	pt models.PlayerType
	s  Strategy
}

func defaultMinBidQualityInput(pt models.PlayerType) minBidQuantityInput {
	return minBidQuantityInput{
		ds: testutil.NewEmptyDraftSnapshot(100, testutil.TestLineupInfo),
		t:  testutil.TestTeam01,
		pt: pt,
		s:  TestStrategy,
	}
}

func TestMinBidQuantity(t *testing.T) {
	var tests = []struct {
		input minBidQuantityInput
		want  int
	}{
		{
			input: defaultMinBidQualityInput(models.QB),
			want:  1,
		},
		{
			input: defaultMinBidQualityInput(models.RB),
			want:  2,
		},
		{
			input: defaultMinBidQualityInput(models.WR),
			want:  3,
		},
		{
			input: defaultMinBidQualityInput(models.TE),
			want:  1,
		},
		{
			input: defaultMinBidQualityInput(models.D),
			want:  1,
		},
		{
			input: defaultMinBidQualityInput(models.K),
			want:  1,
		},
	}

	for _, test := range tests {
		res := minBidQuantity(test.input.ds, test.input.t, test.input.pt, test.input.s)

		if res != test.want {
			t.Errorf("Expected result: %d but got %d for pos: %v", test.want, res, test.input.pt)
		}
	}
}

func TestMinBidAmount(t *testing.T) {
	ts2 := TestStrategy
	ts2.Preemptive = TwoPointMin
	if minBidAmount(ts2) != 2 {
		t.Error("Expected 2 points ")
	}

	ts1 := TestStrategy
	ts1.Preemptive = OnePointMin
	if minBidAmount(ts1) != 1 {
		t.Error("Expected 1 points ")
	}
}

func TestPreemptiveBids(t *testing.T) {
	res := preemptiveBids(testutil.TestDraftSnapshot, testutil.TestTeam01, TestStrategy)
	if len(res) != 9 {
		t.Error("Expected a full 9 bids")
	}
	bidDist := make(map[models.PlayerType][]models.Bid)
	for _, b := range res {
		bidDist[b.Player.Type()] = append(bidDist[b.Player.Type()], b)
	}
	if len(bidDist[models.QB]) != 1 {
		t.Error("Expected min bid for QB backup")
	}
	if bidDist[models.QB][0].Player != testutil.TestQB01 {
		t.Error("Expected QB min bid to be for most valuable available plaeyr")
	}
	if len(bidDist[models.RB]) != 2 {
		t.Error("Expected 2 min bid for RB backups")
	}
	if bidDist[models.RB][0].Player != testutil.TestRB01 {
		t.Error("Expected RB min bid to be for most valuable available plaeyr")
	}
	if len(bidDist[models.WR]) != 3 {
		t.Error("Expected 3 min bids for WR backups")
	}
	if bidDist[models.WR][0].Player != testutil.TestWR01 {
		t.Error("Expected WR min bid to be for most valuable available plaeyr")
	}
	if len(bidDist[models.TE]) != 1 {
		t.Error("Expected min bid for TE backup")
	}
	if bidDist[models.TE][0].Player != testutil.TestTE01 {
		t.Error("Expected TE min bid to be for most valuable available plaeyr")
	}
	if len(bidDist[models.D]) != 1 {
		t.Error("Expected min bid for D")
	}
	if bidDist[models.D][0].Player != testutil.TestD01 {
		t.Error("Expected D min bid to be for most valuable available plaeyr")
	}
	if len(bidDist[models.K]) != 1 {
		t.Error("Expected min bid for K")
	}
	if bidDist[models.K][0].Player != testutil.TestK01 {
		t.Error("Expected K min bid to be for most valuable available plaeyr")
	}

}
