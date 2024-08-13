package googlesheets

import (
	"testing"

	"dcashman.net/coctaleague/pkg/models"
)

func Test_indicesToCellStr(t *testing.T) {
	type input struct {
		row int
		col int
	}
	tests := []struct {
		name string
		in   input
		want string
	}{
		{name: "cellid_A1", in: input{0, 0}, want: "A1"},
		{name: "cellid_B3", in: input{2, 1}, want: "B3"},
		{name: "cellid_AA53", in: input{52, 26}, want: "AA53"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := indicesToCellStr(tt.in.row, tt.in.col); got != tt.want {
				t.Errorf("indicesToCellStr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cellStrToIndices(t *testing.T) {
	type output struct {
		row int
		col int
	}
	tests := []struct {
		name string
		in   string
		want output
	}{
		{name: "cellid_A1", want: output{0, 0}, in: "A1"},
		{name: "cellid_B3", want: output{2, 1}, in: "B3"},
		{name: "cellid_AA53", want: output{52, 26}, in: "AA53"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if row, col, err := cellStrToIndices(tt.in); err != nil || row != tt.want.row || col != tt.want.col {
				t.Errorf("cellStrToIndices() = %d, %d, %v, want %d, %d, nil", row, col, err, tt.want.row, tt.want.col)
			}
		})
	}
}

func Test_bidTocell(t *testing.T) {
	tests := []struct {
		name string
		in   models.Bid
		want string
	}{
		{
			name: "bidToCell",
			in: models.Bid{
				Amount: 4,
				Bidder: &Team{
					cell: "A6",
				},
				Player: &Player{
					cell: "F4",
				},
			},
			want: "K4",
		},
		{
			name: "bidToCell_D",
			in: models.Bid{
				Amount: 3,
				Bidder: &Team{
					cell: "A6",
				},
				Player: &Player{
					pt:   models.D,
					cell: "DG7",
				},
			},
			want: "DK7",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := bidToCell(tt.in); err != nil || got != tt.want {
				t.Errorf("bidToCell() = %s, want %s, err: %v", got, tt.want, err)
			}
		})
	}
}
