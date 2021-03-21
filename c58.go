package set8

import (
	"fmt"
	"math/big"
)

type PollardMapper struct {
	k, c, p *big.Int
}

func NewPollardMapper(k, c, p *big.Int) (*PollardMapper, error) {
	if k.Sign() != 1 || c.Sign() != 1 || p.Sign() != 1 {
		return nil, fmt.Errorf("k, c, p must be positive: %v, %v, %v", k, c, p)
	}
	return &PollardMapper{k: k, c: c, p: p}, nil
}

func (p PollardMapper) F(y, dst *big.Int) {
	dst.Exp(big2, dst.Mod(y, p.k), p.p) // dst = 2^(y mod k) mod p
}

func (p PollardMapper) N(dst *big.Int) {
	// Will hold all possible outputs of F.
	seen := make([]*big.Int, 0, p.k.Int64())
	for i := big.NewInt(0); i.Cmp(p.k) < 0; i.Add(i, big1) {
		p.F(i, dst)
		seen = append(seen, new(big.Int).Set(dst))
	}

	// Compute N = mean(seen) * c.
	dst.SetInt64(0)
	for _, s := range seen {
		dst.Add(dst, s)
	}
	dst.Mul(dst.Div(dst, p.k), p.c)
}

func PollardsKangaroo(pm *PollardMapper, p, g, a, b, y, dst *big.Int) error {
	var (
		// Tame kangaroo
		xT, yT big.Int
		// Wild kangaroo
		xW, yW big.Int
		// N value for the mapper
		n big.Int
	)

	// Set N.
	pm.N(&n)

	xT.SetInt64(0)  // xT = 0
	yT.Exp(g, b, p) // yT = g^b
	xW.SetInt64(0)  // xW = 0
	yW.Set(y)       // yW = y

	for i := big.NewInt(0); i.Cmp(&n) < 0; i.Add(i, big1) {
		pm.F(&yT, dst)                              // Compute f(yT) only once per iteration.
		xT.Add(&xT, dst).Mod(&xT, p)                // xT = xT + f(yT)
		yT.Mul(&yT, dst.Exp(g, dst, p)).Mod(&yT, p) // yT = yT * g^f(yT)
	}

	// while xW < b - a + xT
	forBound := dst.Add(dst.Sub(b, a), &xT)
	for xW.Cmp(forBound) < 0 {
		pm.F(&yW, dst)                              // Compute f(yW) only once per iteration.
		xW.Add(&xW, dst).Mod(&xW, p)                // xW = xW + f(yW)
		yW.Mul(&yW, dst.Exp(g, dst, p)).Mod(&yW, p) // yW = yW * g^f(yW)

		// If wild y and tame y collide, success!
		if yW.Cmp(&yT) == 0 {
			dst.Sub(dst.Add(b, &xT), &xW) // b + xT - xW
			return nil
		}
	}

	return fmt.Errorf("no index found")
}
