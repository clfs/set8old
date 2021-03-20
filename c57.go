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
// It's unoptimized, so try to use a small bound (< 1 million).
func PrimeFactorsLessThan(n, bound *big.Int) []*big.Int {
	if n.Sign() < 1 || bound.Sign() < 1 {
		return nil
	}

	var (
		res []*big.Int
		tmp big.Int
	)

	for f := big.NewInt(1); f.Cmp(bound) < 0; f.Add(f, big1) {
		// if n divisible by f AND f is prime, then save it
		if tmp.Mod(n, f).Sign() == 0 && f.ProbablyPrime(1) {
			res = append(res, new(big.Int).Set(f))
		}
	}

	return res
}

// SubgroupConfinementAttack recovers Bob's secret key in a DHKE scheme, by way
// of the Pohlig-Hellman algorithm for discrete logarithms. The result is stored
// in dst.
func SubgroupConfinementAttack(bob *C57Bob, p, g, q, dst *big.Int) error {
	var (
		// Invalid public key to force subgroups.
		h big.Int
		// We'll need this for the CRT step.
		crtPairs []*crt.Pair
	)

	dst.Div(dst.Sub(p, big1), q) // j, or (p - 1) // q
	jFactors := PrimeFactorsLessThan(dst, big65536)
	for _, n := range jFactors {
		// Pick an element h of order f.
		for {
			rnd, err := rand.Int(rand.Reader, dst.Sub(p, big1)) // [0, p-1)
			if err != nil {
				return err
			}
			rnd.Add(rnd, big1) // [1, p)

			// Try to pick h != 1.
			h.Exp(rnd, dst.Div(dst.Sub(p, big1), n), p)
			if h.Cmp(big1) != 0 {
				break
			}
		}

		// Query Bob.
		msg, tag, err := bob.Query(&h)
		if err != nil {
			return nil
		}

		// Brute-force Bob's secret key mod f.
		for a := big.NewInt(0); a.Cmp(n) < 0; a.Add(a, big1) {
			guess, err := HMACSHA256(dst.Exp(&h, a, p).Bytes(), msg)
			if err != nil {
				return err
			}
			// No need for hmac.Equal since we're the attacker.
			if bytes.Equal(guess, tag) {
				crtPairs = append(crtPairs, &crt.Pair{
					A: new(big.Int).Set(a),
					N: new(big.Int).Set(n),
				})
				break
			}
		}
	}

	err := crt.Do(crtPairs, dst)
	if err != nil {
		// CRT can fail if there are no pairs, or if the divisors aren't
		// pairwise coprime. Neither of these apply here.
		return fmt.Errorf("this should never happen")
	}
	return nil
}
