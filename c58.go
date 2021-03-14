package set8

import (
	"fmt"
	"math/big"
)

type PollardMapper struct {
	// Constants
	k, c, p *big.Int
	// Temporary storage for 2
	two *big.Int
}

func NewPollardMapper(k, c, p *big.Int) (*PollardMapper, error) {
	if k.Sign() != 1 || c.Sign() != 1 || p.Sign() != 1 {
		return nil, fmt.Errorf("k, c, p must be positive: %v, %v, %v", k, c, p)
	}
	return &PollardMapper{k: k, c: c, p: p, two: big.NewInt(2)}, nil
}

func (p PollardMapper) F(y, dst *big.Int) {
	dst.Exp(p.two, dst.Mod(y, p.k), p.p) // dst = 2^(y mod k) mod p
}

func (p PollardMapper) N() *big.Int {
	// No need to go hard on pre-allocating, since this isn't a hot path.
	tmp := big.NewInt(0)
	one := big.NewInt(1)

	// Will hold all possible outputs of F.
	seen := make([]*big.Int, 0)
	// for i in [0, k)
	for i := big.NewInt(0); i.Cmp(p.k) < 0; i.Add(i, one) {
		// seen[i] = F(i)
		p.F(i, tmp)
		seen = append(seen, new(big.Int).Set(tmp))
	}

	// Compute N.
	tmp.SetInt64(0)
	for _, s := range seen {
		tmp.Add(tmp, s)
	}
	count := big.NewInt(int64(len(seen)))    // len(seen) is never 0
	return tmp.Mul(tmp.Div(tmp, count), p.c) // mean(seen) * c
}

func PollardsKangaroo(p, g, a, b, y *big.Int, pm *PollardMapper) (*big.Int, error) {
	var (
		// Tame kangaroo
		xT, yT big.Int
		// Wild kangaroo
		xW, yW big.Int
		// Temporary storage
		tmp = big.NewInt(0)
		// Constants
		one = big.NewInt(1)
	)

	// TODO set n correctly.
	n := pm.N()

	xT.SetInt64(0)  // xT = 0
	yT.Exp(g, b, p) // yT = g^b
	xW.SetInt64(0)  // xW = 0
	yW.Set(y)       // yW = y

	// for i in [0, n)
	for i := big.NewInt(0); i.Cmp(n) < 0; i.Add(i, one) {
		pm.F(&yT, tmp)
		xT.Add(&xT, tmp)                // xT = xT + f(yT)
		yT.Mul(&yT, tmp.Exp(g, tmp, p)) // yT = yT * g^f(yT)
	}

	// while xW < b - a + xT
	for xW.Cmp(tmp.Add(tmp.Sub(b, a), &xT)) < 0 {
		pm.F(&yW, tmp)
		xW.Add(&xW, tmp)                // xW = xW + f(yW)
		yW.Mul(&yW, tmp.Exp(g, tmp, p)) // yW = yW * g^f(yW)

		if yW.Cmp(&yT) == 0 {
			return tmp.Sub(tmp.Add(b, &xT), &xW), nil // b + xT - xW
		}
	}

	return nil, fmt.Errorf("no index found")
}
