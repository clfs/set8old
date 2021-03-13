package set8

import (
	"fmt"
	"math/big"
)

type CRTPair struct {
	Remainder, Divisor *big.Int
}

// CRT is adapted from https://rosettacode.org/wiki/Chinese_remainder_theorem#Go
func CRT(pairs []*CRTPair) (*big.Int, error) {
	if len(pairs) == 0 {
		return nil, fmt.Errorf("no pairs provided")
	}

	one := big.NewInt(1)

	p := new(big.Int).Set(pairs[0].Divisor)
	for _, pair := range pairs[1:] {
		p.Mul(p, pair.Divisor)
	}

	var x, q, s, z big.Int
	for _, pair := range pairs {
		q.Div(p, pair.Divisor)
		z.GCD(nil, &s, pair.Divisor, &q)
		if z.Cmp(one) != 0 {
			return nil, fmt.Errorf("not pairwise coprime with divisor %v", pair.Divisor)
		}
		x.Add(&x, s.Mul(pair.Remainder, s.Mul(&s, &q)))
	}

	return x.Mod(&x, p), nil
}
