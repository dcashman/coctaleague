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
