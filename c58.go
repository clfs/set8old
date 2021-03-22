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

func PollardsKangaroo(pm *PollardMapper, p, g, a, b, y *big.Int) (*big.Int, error) {
	var (
		// Tame kangaroo
		xT, yT big.Int
		// Wild kangaroo
		xW, yW big.Int
		// N value for the mapper
		n big.Int
		// Temporary values
		t1, t2 big.Int
	)

	// The cache maps cacheKey to Exp(g, cacheKey, p).
	cache := make(map[uint64]*big.Int)
	for i := big.NewInt(0); i.Cmp(pm.k) < 0; i.Add(i, big1) {
		pm.F(i, &t1)
		t2.Exp(g, &t1, p)
		cache[t1.Uint64()] = new(big.Int).Set(&t2)
	}

	pm.N(&n)

	xT.SetInt64(0)  // xT = 0
	yT.Exp(g, b, p) // yT = g^b
	xW.SetInt64(0)  // xW = 0
	yW.Set(y)       // yW = y

	for i := big.NewInt(0); i.Cmp(&n) < 0; i.Add(i, big1) {
		pm.F(&yT, &t1) // Compute f(yT) only once.

		t2.Add(&xT, &t1)
		xT.Mod(&t2, p) // xT = xT + f(yT)

		t2.Set(cache[t1.Uint64()]) // g^f(yT) mod p, cached
		t1.Mul(&yT, &t2)
		yT.Mod(&t1, p) // yT = yT * g^f(yT)
	}

	t1.Sub(b, a)
	forBound := t2.Add(&t1, &xT) // b - a + xT

	for xW.Cmp(forBound) < 0 {
		pm.F(&yW, &t1) // Compute f(yW) only once.

		t2.Add(&xW, &t1)
		xW.Mod(&t2, p) // xW = xW + f(yW)

		t2.Set(cache[t1.Uint64()]) // g^f(yW) mod p, cached
		t1.Mul(&yW, &t2)
		yW.Mod(&t1, p) // yW = yW * g^f(yW)

		// If wild y and tame y collide, success!
		if yW.Cmp(&yT) == 0 {
			t1.Add(b, &xT)
			return t2.Sub(&t1, &xW), nil // b + xT - xW
		}
	}

	return nil, fmt.Errorf("no index found")
}
