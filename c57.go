package set8

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
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

// Respond accepts a public key without validating it. Bob computes a shared
// secret, then replies with a message and MAC.
func (c *C57Bob) Respond(h *big.Int) ([]byte, []byte, error) {
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

// PrimeFactorsLessThan returns all prime factors of n less than bound. If n or
// bound are non-positive, no factors are returned. It's also very unoptimized,
// so use a small bound.
func PrimeFactorsLessThan(n, bound *big.Int) []*big.Int {
	if n.Sign() < 1 || bound.Sign() < 1 {
		return nil
	}

	var (
		res []*big.Int
		tmp big.Int
	)

	for f := big.NewInt(1); f.Cmp(bound) < 0; f.Add(f, big1) {
		// If n is divisible by f, and f is prime, then save it.
		if tmp.Mod(n, f).Sign() == 0 && f.ProbablyPrime(1) {
			res = append(res, new(big.Int).Set(f))
		}
	}

	return res
}

// RandInt returns a uniform random value in [a, b). It panics if a >= b.
func RandInt(a, b *big.Int) (*big.Int, error) {
	tmp := new(big.Int).Sub(b, a)          // tmp = b - a
	tmp, err := rand.Int(rand.Reader, tmp) // tmp in [0, b-a)
	if err != nil {
		return nil, err
	}
	tmp.Add(tmp, a) // tmp += a
	return tmp, nil // tmp in [a, b)
}

// SubgroupConfinementAttack recovers Bob's secret key in a DHKE scheme, by way
// of the Pohlig-Hellman algorithm for discrete logarithms.
func SubgroupConfinementAttack(bob *C57Bob, p, g, q *big.Int) (*big.Int, error) {
	var pairs []crt.Pair // Remainder/divisor pairs for the CRT step.

	var j big.Int // j = (p - 1) // q
	j.Sub(p, big1)
	j.Div(&j, q)

	for _, r := range PrimeFactorsLessThan(&j, big65536) {
		var h big.Int // Holds an invalid public key.

		// Find an invalid public key that's an element of order r.
		for {
			rnd, err := RandInt(big1, p) // rnd in [1, p)
			if err != nil {
				return nil, err
			}

			h.Exp(rnd, h.Div(h.Sub(p, big1), r), p) // h = rnd^((p-1)/r) mod p

			// Try until h != 1.
			if h.Cmp(big1) != 0 {
				break
			}
		}

		// Query Bob.
		msg, tag, err := bob.Respond(&h)
		if err != nil {
			return nil, err
		}

		// Brute-force Bob's secret key mod f.
		for a := big.NewInt(0); a.Cmp(r) < 0; a.Add(a, big1) {
			// Reuse j to save memory.
			j.Exp(&h, a, p) // j = h^a mod p

			guess, err := HMACSHA256(j.Bytes(), msg)
			if err != nil {
				return nil, err
			}

			// No need for hmac.Equal since we're the attacker.
			if bytes.Equal(guess, tag) {
				pairs = append(pairs, crt.Pair{
					Remainder: new(big.Int).Set(a),
					Divisor:   new(big.Int).Set(r),
				})
				break
			}
		}
	}

	return crt.Do(pairs)
}
