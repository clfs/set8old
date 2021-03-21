package set8

import "math/big"

// Avoid using these in tests or benchmarks.
var (
	big1     = big.NewInt(1)
	big2     = big.NewInt(2)
	big65536 = big.NewInt(65536) // 2^16
)
