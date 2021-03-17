package set8

import "math/big"

// One is a big.Int equal to 1.
var One = big.NewInt(1)

// TaggedMessage represents a message and tag pair, both as byte slices.
type TaggedMessage struct {
	Message, Tag []byte
}
