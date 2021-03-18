package set8

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/clfs/set8/crt"
)

// C57Bob represents Bob for challenge 57.
type C57Bob struct {
	p, g *big.Int
	key  *big.Int
	msg  []byte
}

// NewC57Bob returns a new C57Bob.
func NewC57Bob(p, g, q *big.Int) (*C57Bob, error) {
	key, err := rand.Int(rand.Reader, q)
	if err != nil {
		return nil, err
	}
	return &C57Bob{
		p:   p,
		g:   g,
		key: key,
		msg: []byte("crazy flamboyant for the rap enjoyment"),
	}, nil
}

// Query accepts a public key without validating it.
// Bob computes a shared secret and returns a message
// and its MAC.
func (c *C57Bob) Query(h *big.Int) ([]byte, []byte, error) {
	sharedKey := new(big.Int).Exp(h, c.key, c.p)
	tag, err := HMACSHA256(sharedKey.Bytes(), c.msg)
	if err != nil {
		return nil, nil, err
	}
	return c.msg, tag, nil
}

// HMACSHA256 computes the HMAC-SHA256 of a message under a key.
// If you're signing lots of messages with the same key, don't
// use this - it'll be inefficient.
func HMACSHA256(key, msg []byte) ([]byte, error) {
	mac := hmac.New(sha256.New, key)
	_, err := mac.Write(msg)
	if err != nil {
		return nil, err
	}
	return mac.Sum(nil), nil
}

// PrimeFactorsLessThan returns all prime factors of n that are less
// than bound. It returns nil if either n or bound is non-positive.
// It's fully unoptimized, so try to use a small bound (< 1 million).
func PrimeFactorsLessThan(n, bound *big.Int) []*big.Int {
	if n.Sign() < 1 || bound.Sign() < 1 {
		return nil
	}

	var res []*big.Int

	one := big.NewInt(1)
	var tmp big.Int

	// for f in [1, bound)
	for f := big.NewInt(1); f.Cmp(bound) < 0; f.Add(f, one) {
		// if n divisible by f AND f is prime
		if tmp.Mod(n, f).Sign() == 0 && f.ProbablyPrime(1) {
			// save a copy of f
			res = append(res, new(big.Int).Set(f))
		}
	}

	return res
}

// SubgroupConfinementAttack recovers Bob's secret key in a DHKE scheme, by way
// of the Pohlig-Hellman algorithm for discrete logarithms.
func SubgroupConfinementAttack(p, g, q *big.Int, bob *C57Bob) (*big.Int, error) {
	// Define some constants in advance.
	one := big.NewInt(1)

	// Compute j and its factors.
	j := new(big.Int).Sub(p, one) // p - 1
	j.Div(j, q)                   // (p - 1) // q
	jFactors := PrimeFactorsLessThan(j, big.NewInt(65536))

	// Set up the CRT pairs for the eventual Chinese Remainder Theorem.
	var crtPairs []*crt.Pair

	// Avoid allocating excess bigints in the for loop.
	h := new(big.Int)
	tmp := new(big.Int)

	// It's possible to stop early if the product of the remainders is greater
	// than q, but I'm too lazy to write that check. Instead, every factor is
	// used.
	for _, f := range jFactors {
		// Find an element h of order f.
		for {
			// Generate a random integer in [1, p).
			rnd, err := rand.Int(rand.Reader, tmp.Sub(p, one)) // [0, p-1)
			if err != nil {
				return nil, err
			}
			rnd.Add(one, rnd) // [0, p-1) becomes [1, p)

			// Compute h.
			tmp.Div(tmp.Sub(p, one), f) // (p - 1) // f
			h.Exp(rnd, tmp, p)

			// Retry if h is 1.
			if h.Cmp(one) != 0 {
				break
			}
		}

		// Query Bob with h.
		msg, tag, err := bob.Query(h)
		if err != nil {
			return nil, err
		}

		// Brute-force Bob's secret key mod f.
		for i := big.NewInt(0); i.Cmp(f) < 0; i.Add(i, one) {
			// Make a guess.
			guess, err := HMACSHA256(tmp.Exp(h, i, p).Bytes(), msg)
			if err != nil {
				return nil, err
			}
			// If the guess was correct, obtain a CRT pair and move on
			// to the next factor of j.
			if bytes.Equal(guess, tag) {
				crtPairs = append(crtPairs, &crt.Pair{
					A: new(big.Int).Set(i),
					N: new(big.Int).Set(f),
				})
				break
			}
		}
	}

	err := crt.Do(crtPairs, tmp)
	if err != nil {
		// CRT can fail if the divisors aren't pairwise coprime, but that should
		// never happen, since all the divisors we chose were prime.
		return nil, fmt.Errorf("this should never happen")
	}
	return tmp, nil
}
