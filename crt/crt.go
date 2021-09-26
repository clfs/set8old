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

// Do does the CRT. It's adapted from
// https://rosettacode.org/wiki/Chinese_remainder_theorem#Go. It fails
// if the divisors are not pairwise coprime.
func Do(pairs []Pair) (*big.Int, error) {
	if len(pairs) == 0 {
		return nil, fmt.Errorf("no pairs provided")
	}

	var product big.Int // The product of all divisors.

	product.Set(pairs[0].N)
	for _, p := range pairs[1:] {
		product.Mul(&product, p.N)
	}

	var x, q, s, z big.Int

	for _, p := range pairs {
		q.Div(&product, p.N)    // q = product / p.N
		z.GCD(nil, &s, p.N, &q) // z = gcd(p.N, q), then s = z / b

		if z.Cmp(One) != 0 { // if z == 1
			return nil, fmt.Errorf("divisor %d violates pairwise coprime requirement", p.N)
		}

		s.Mul(&s, &q)  // s *= q
		s.Mul(&s, p.A) // s *= p.A
		x.Add(&x, &s)  // x += s
	}

	return x.Mod(&x, &product), nil // x %= product
}
