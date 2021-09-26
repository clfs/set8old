package crt

import (
	"math/big"
	"testing"
)

func BenchmarkDo(b *testing.B) {
	pairs := []Pair{
		{big.NewInt(2), big.NewInt(3)},
		{big.NewInt(3), big.NewInt(4)},
		{big.NewInt(1), big.NewInt(5)},
	}
	for i := 0; i < b.N; i++ {
		_, _ = Do(pairs)
	}
}

func TestDo(t *testing.T) {
	t.Parallel()
	cases := map[string]struct {
		pairs []Pair
		want  *big.Int
	}{
		"one pair": {
			pairs: []Pair{
				{big.NewInt(3), big.NewInt(5)},
			},
			want: big.NewInt(3),
		},
		"two pairs": {
			pairs: []Pair{
				{big.NewInt(1), big.NewInt(5)},
				{big.NewInt(3), big.NewInt(7)},
			},
			want: big.NewInt(31),
		},
		"three pairs": {
			pairs: []Pair{
				{big.NewInt(2), big.NewInt(3)},
				{big.NewInt(3), big.NewInt(4)},
				{big.NewInt(1), big.NewInt(5)},
			},
			want: big.NewInt(11),
		},
	}
	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got, err := Do(tc.pairs)
			if err != nil {
				t.Errorf("Do failed: %v", err)
			}
			if got.Cmp(tc.want) != 0 {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}
