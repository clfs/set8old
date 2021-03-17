package crt

import (
	"math/big"
	"testing"
)

func TestCRT(t *testing.T) {
	t.Parallel()
	cases := map[string]struct {
		pairs []*CRTPair
		want  *big.Int
	}{
		"two pairs": {
			pairs: []*CRTPair{
				{
					Remainder: big.NewInt(1),
					Divisor:   big.NewInt(5),
				},
				{
					Remainder: big.NewInt(3),
					Divisor:   big.NewInt(7),
				},
			},
			want: big.NewInt(31),
		},
		"three pairs": {
			pairs: []*CRTPair{
				{
					Remainder: big.NewInt(2),
					Divisor:   big.NewInt(3),
				},
				{
					Remainder: big.NewInt(3),
					Divisor:   big.NewInt(4),
				},
				{
					Remainder: big.NewInt(1),
					Divisor:   big.NewInt(5),
				},
			},
			want: big.NewInt(11),
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got, err := CRT(tc.pairs)
			if err != nil {
				t.Errorf("CRT failed: %v", err)
			}
			if got.Cmp(tc.want) != 0 {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}
