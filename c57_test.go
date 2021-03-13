package set8

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrimeFactorsLessThan(t *testing.T) {
	t.Parallel()
	cases := map[string]struct {
		n, bound *big.Int
		want     []*big.Int
	}{
		"zero n": {
			n:     big.NewInt(0),
			bound: big.NewInt(100),
			want:  []*big.Int{},
		},
		"zero bound": {
			n:     big.NewInt(100),
			bound: big.NewInt(0),
			want:  []*big.Int{},
		},
		"zero n and zero bound": {
			n:     big.NewInt(0),
			bound: big.NewInt(0),
			want:  []*big.Int{},
		},
		"negative n": {
			n:     big.NewInt(-1),
			bound: big.NewInt(100),
			want:  []*big.Int{},
		},
		"negative bound": {
			n:     big.NewInt(100),
			bound: big.NewInt(-1),
			want:  []*big.Int{},
		},
		"negative n and negative bound": {
			n:     big.NewInt(-1),
			bound: big.NewInt(-1),
			want:  []*big.Int{},
		},
		"prime n": {
			n:     big.NewInt(7),
			bound: big.NewInt(7),
			want:  []*big.Int{},
		},
		"one factor": {
			n:     big.NewInt(10),
			bound: big.NewInt(3),
			want:  []*big.Int{big.NewInt(2)},
		},
		"multiple factors": {
			n:     big.NewInt(10),
			bound: big.NewInt(10),
			want:  []*big.Int{big.NewInt(2), big.NewInt(5)},
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := PrimeFactorsLessThan(tc.n, tc.bound)
			assert.ElementsMatch(t, got, tc.want)
		})
	}
}

func TestPrimeFactorsLessThan_Large(t *testing.T) {
	t.Parallel()
	// This test uses the actual j parameter.
	n, ok := new(big.Int).SetString("30477252323177606811760882179058908038824640750610513771646768011063128035873508507547741559514324673960576895059570", 10)
	if !ok {
		t.Error("failed to create n")
	}
	bound := big.NewInt(65536) // 2^16
	want := []*big.Int{
		big.NewInt(2),
		big.NewInt(3),
		big.NewInt(5),
		big.NewInt(109),
		big.NewInt(7963),
		big.NewInt(8539),
		big.NewInt(20641),
		big.NewInt(38833),
		big.NewInt(39341),
		big.NewInt(46337),
		big.NewInt(51977),
		big.NewInt(54319),
		big.NewInt(57529),
	}
	got := PrimeFactorsLessThan(n, bound)
	assert.ElementsMatch(t, got, want)
}

func TestSubgroupConfinementAttack(t *testing.T) {
	t.Parallel()

	p, ok := new(big.Int).SetString("7199773997391911030609999317773941274322764333428698921736339643928346453700085358802973900485592910475480089726140708102474957429903531369589969318716771", 10)
	if !ok {
		t.Error("failed to create p")
	}
	g, ok := new(big.Int).SetString("4565356397095740655436854503483826832136106141639563487732438195343690437606117828318042418238184896212352329118608100083187535033402010599512641674644143", 10)
	if !ok {
		t.Error("failed to create g")
	}
	q, ok := new(big.Int).SetString("236234353446506858198510045061214171961", 10)
	if !ok {
		t.Error("failed to create q")
	}
	client, err := NewC57BobClient(p, g, q)
	if err != nil {
		t.Errorf("failed to create Bob client: %v", err)
	}
	key, err := SubgroupConfinementAttack(p, g, q, client)
	if err != nil {
		t.Errorf("failed to find Bob key: %v", err)
	}
	if key.Cmp(client.key) != 0 {
		t.Errorf("incorrect key: got %v, want %v", key, client.key)
	}
}
