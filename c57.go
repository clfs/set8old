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

// C57Bob represents Bob in challenge 57.
type C57Bob struct {
	p, g *big.Int // Group constants
	key  *big.Int // Bob's secret key
	msg  []byte   // Bob's static message
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

// Query accepts a public key without validating it. Bob computes a shared
// secret, then replies with a message and MAC.
func (c *C57Bob) Query(h *big.Int) ([]byte, []byte, error) {
	var sharedSecret big.Int
	sharedSecret.Exp(h, c.key, c.p) // shared secret
	tag, err := HMACSHA256(sharedSecret.Bytes(), c.msg)
	if err != nil {
		return nil, nil, err
	}
	return c.msg, tag, nil
}

// HMACSHA256 returns the HMAC-SHA256 of a message under a key. Don't use this
// if you're signing lots of messages with the same key - it's inefficient.
func HMACSHA256(key, msg []byte) ([]byte, error) {
	mac := hmac.New(sha256.New, key)
	_, err := mac.Write(msg)
	if err != nil {
		return nil, err
	}
	return mac.Sum(nil), nil
}

// PrimeFactorsLessThan returns all prime factors of n less than bound. It
// returns nil if n or bound are non-positive. It's also very unoptimized, so
// use a small bound.
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
func SubgroupConfinementAttack(bob *C57Bob, p, g, q *big.Int) (*big.Int, error) {
	var (
		h        big.Int // Invalid public key to force subgroups.
		dst      big.Int
		crtPairs []crt.Pair // We'll need this for the CRT step.
	)

	dst.Div(dst.Sub(p, big1), q) // j, or (p - 1) // q
	jFactors := PrimeFactorsLessThan(&dst, big65536)
	for _, n := range jFactors {
		// Pick an element h of order f.
		for {
			rnd, err := rand.Int(rand.Reader, dst.Sub(p, big1)) // [0, p-1)
			if err != nil {
				return nil, err
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
			return nil, err
		}

		// Brute-force Bob's secret key mod f.
		for a := big.NewInt(0); a.Cmp(n) < 0; a.Add(a, big1) {
			guess, err := HMACSHA256(dst.Exp(&h, a, p).Bytes(), msg)
			if err != nil {
				return nil, err
			}
			// No need for hmac.Equal since we're the attacker.
			if bytes.Equal(guess, tag) {
				crtPairs = append(crtPairs, crt.Pair{
					Remainder: new(big.Int).Set(a),
					Divisor:   new(big.Int).Set(n),
				})
				break
			}
		}
	}

	res, err := crt.Do(crtPairs)
	if err != nil {
		// CRT fails if there are no pairs, or if the divisors aren't pairwise
		// coprime. Neither of these apply here.
		return nil, fmt.Errorf("this should never happen")
	}

	return res, nil
}
