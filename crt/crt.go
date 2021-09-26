package crt

import (
	"fmt"
	"math/big"
)

// One is a big.Int equal to 1.
var One = big.NewInt(1)

// Pair contains a remainder A and a divisor N.
type Pair struct {
	Remainder, Divisor *big.Int
}

// Do does the CRT. It's adapted from
// https://rosettacode.org/wiki/Chinese_remainder_theorem#Go. It fails
// if the divisors are not pairwise coprime.
func Do(pairs []Pair) (*big.Int, error) {
	if len(pairs) == 0 {
		return nil, fmt.Errorf("no pairs provided")
	}

	var product big.Int // The product of all divisors.

	product.Set(pairs[0].Divisor)
	for _, p := range pairs[1:] {
		product.Mul(&product, p.Divisor)
	}

	var x, q, s, z big.Int

	for _, p := range pairs {
		q.Div(&product, p.Divisor)    // q = product / divisor
		z.GCD(nil, &s, p.Divisor, &q) // z = gcd(divisor, q), then s = z / q

		if z.Cmp(One) != 0 { // if z == 1
			return nil, fmt.Errorf("divisor %d violates pairwise coprime requirement", p.Divisor)
		}

		s.Mul(&s, &q)          // s *= q
		s.Mul(&s, p.Remainder) // s *= remainder
		x.Add(&x, &s)          // x += s
	}

	return x.Mod(&x, &product), nil // x %= product
}
