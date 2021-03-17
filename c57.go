package set8

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/clfs/crt"
)

// C57BobClient is a client representing Bob for challenge 57.
type C57BobClient struct {
	key, p, g, q *big.Int
}

// NewC57BobClient returns a queryable client representing Bob.
func NewC57BobClient(p, g, q *big.Int) (*C57BobClient, error) {
	key, err := rand.Int(rand.Reader, q)
	if err != nil {
		return nil, err
	}
	return &C57BobClient{key: key, p: p, g: g, q: q}, nil
}

// Query accepts a "supposed" public key h without validation.
// Bob computes a shared secret, and sends back a message
// with a MAC.
func (c *C57BobClient) Query(h *big.Int) (*TaggedMessage, error) {
	sharedKey := new(big.Int).Exp(h, c.key, c.p)
	msg := []byte("crazy flamboyant for the rap enjoyment")
	tag, err := HMACSHA256(sharedKey, msg)
	if err != nil {
		return nil, err
	}
	return &TaggedMessage{
		Message: msg,
		Tag:     tag,
	}, nil
}

// HMACSHA256 computes the HMAC-SHA256 of a message under a key.
func HMACSHA256(key *big.Int, msg []byte) ([]byte, error) {
	mac := hmac.New(sha256.New, key.Bytes())
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
func SubgroupConfinementAttack(p, g, q *big.Int, client *C57BobClient) (*big.Int, error) {
	// Define some constants in advance.
	one := big.NewInt(1)

	// Compute j and its factors.
	j := new(big.Int).Sub(p, one) // p - 1
	j.Div(j, q)                   // (p - 1) // q
	jFactors := PrimeFactorsLessThan(j, big.NewInt(65536))

	// Set up the CRT pairs for the eventual Chinese Remainder Theorem.
	var crtPairs []*crt.CRTPair

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
		response, err := client.Query(h)
		if err != nil {
			return nil, err
		}

		// Brute-force Bob's secret key mod f.
		for i := big.NewInt(0); i.Cmp(f) < 0; i.Add(i, one) {
			// Make a guess.
			tag, err := HMACSHA256(tmp.Exp(h, i, p), response.Message)
			if err != nil {
				return nil, err
			}
			// If the guess was correct, obtain a CRT pair and move on
			// to the next factor of j.
			if bytes.Equal(tag, response.Tag) {
				crtPairs = append(crtPairs, &crt.CRTPair{
					Remainder: new(big.Int).Set(i),
					Divisor:   new(big.Int).Set(f),
				})
				break
			}
		}
	}

	bobKey, err := crt.CRT(crtPairs)
	if err != nil {
		// CRT can fail if the divisors aren't pairwise coprime, but that should
		// never happen, since all the divisors we chose were prime.
		return nil, fmt.Errorf("this should never happen")
	}

	return bobKey, nil
}
