package bid

import (
	"testing"

	"dcashman.net/coctaleague/pkg/internal/testutil"
	"dcashman.net/coctaleague/pkg/models"
)

var (
	TestStrategy = Strategy{
		Style:      Value,
		Value:      Predicted,
		Preemptive: TwoPointMin,
	}
	TestTeamComposition = TeamComposition{
		models.QB: PositionDistribution{Start: 1, Bench: 1},
		models.RB: PositionDistribution{Start: 2, Bench: 2},
		models.WR: PositionDistribution{Start: 3, Bench: 3},
		models.TE: PositionDistribution{Start: 1, Bench: 1},
		models.K:  PositionDistribution{Start: 1, Bench: 0},
		models.D:  PositionDistribution{Start: 1, Bench: 0},
	}
)

type DesiredTeamCompositionInput struct {
	ds       models.DraftSnapshot
	strategy Strategy
}

func DesiredTeamCompositionEqual(a TeamComposition, b TeamComposition) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		/*if k == models.WR {
			fmt.Printf("WTFAFKLJSDF key value: %v\n", k)
			return false
		}*/
		if v != b[k] {
			return false
		}
	}
	return true
}

func TestDesiredTeamComposition(t *testing.T) {
	var tests = []struct {
		input DesiredTeamCompositionInput
		want  TeamComposition
	}{
		{
			input: DesiredTeamCompositionInput{
				ds:       testutil.NewEmptyDraftSnapshot(100, testutil.TestLineupInfo),
				strategy: TestStrategy,
			},
			want: TestTeamComposition,
		},
	}

	for _, test := range tests {
		res := DesiredTeamComposition(test.input.ds, test.input.strategy)

		if !DesiredTeamCompositionEqual(test.want, res) {
			t.Errorf("Expected result: %v but got %v", test.want, res)
		}
	}
}
