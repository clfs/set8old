package crt

import (
	"math/big"
	"testing"
)

func TestDo(t *testing.T) {
	t.Parallel()
	cases := map[string]struct {
		pairs []*Pair
		want  *big.Int
	}{
		"one pair": {
			pairs: []*Pair{
				{A: big.NewInt(3), N: big.NewInt(5)},
			},
			want: big.NewInt(3),
		},
		"two pairs": {
			pairs: []*Pair{
				{A: big.NewInt(1), N: big.NewInt(5)},
				{A: big.NewInt(3), N: big.NewInt(7)},
			},
			want: big.NewInt(31),
		},
		"three pairs": {
			pairs: []*Pair{
				{A: big.NewInt(2), N: big.NewInt(3)},
				{A: big.NewInt(3), N: big.NewInt(4)},
				{A: big.NewInt(1), N: big.NewInt(5)},
			},
			want: big.NewInt(11),
		},
	}
	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := new(big.Int)
			err := Do(tc.pairs, got)
			if err != nil {
				t.Errorf("Do failed: %v", err)
			}
			if got.Cmp(tc.want) != 0 {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}
