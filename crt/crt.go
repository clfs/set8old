package crt

import (
	"fmt"
	"math/big"
)

// One is a big.Int equal to 1.
var One = big.NewInt(1)

// Pair contains a remainder A and a divisor N.
type Pair struct {
	A, N *big.Int
}

// Do performs the CRT and stores the result in dst. It's adapted from
// https://rosettacode.org/wiki/Chinese_remainder_theorem#Go.
func Do(pairs []*Pair, dst *big.Int) error {
	if len(pairs) == 0 {
		return fmt.Errorf("TODO no pairs provided")
	}

	dst.Set(pairs[0].N)
	for _, p := range pairs[1:] {
		dst.Mul(dst, p.N)
	}

	var x, q, s, z big.Int

	for _, p := range pairs {
		q.Div(dst, p.N)
		z.GCD(nil, &s, p.N, &q)
		if z.Cmp(One) != 0 {
			return fmt.Errorf("TODO not pairwise coprime with divisor %v", p.N)
		}
		x.Add(&x, s.Mul(p.A, s.Mul(&s, &q)))
	}

	dst.Mod(&x, dst)
	return nil
}
