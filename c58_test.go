package set8

import (
	"math/big"
	"testing"
)

var testCasesPollardsKangaroo = map[string]struct {
	p, g, a, b, y string // Pollard's Kangaroo inputs
	c, k          string // Constants for the Pollard mapper
	want          string // Desired result
}{
	"small": {
		p:    "11470374874925275658116663507232161402086650258453896274534991676898999262641581519101074740642369848233294239851519212341844337347119899874391456329785623",
		g:    "622952335333961296978159266084741085889881358738459939978290179936063635566740258555167783009058567397963466103140082647486611657350811560630587013183357",
		a:    "0",
		b:    "1048576", // 2^20
		y:    "7760073848032689505395005705677365876654629189298052775754597607446617558600394076764814236081991643094239886772481052254010323780165093955236429914607119",
		c:    "3",
		k:    "11",
		want: "705485",
	},
	// "large": {
	// 	p:    "11470374874925275658116663507232161402086650258453896274534991676898999262641581519101074740642369848233294239851519212341844337347119899874391456329785623",
	// 	g:    "622952335333961296978159266084741085889881358738459939978290179936063635566740258555167783009058567397963466103140082647486611657350811560630587013183357",
	// 	a:    "0",
	// 	b:    "1099511627776", // 2^40
	// 	y:    "9388897478013399550694114614498790691034187453089355259602614074132918843899833277397448144245883225611726912025846772975325932794909655215329941809013733",
	// 	c:    "3",
	// 	k:    "16",
	// 	want: "359579674340",
	// },
}

func BenchmarkPollardsKangaroo(b *testing.B) {
	var (
		p  = bigInt(b, testCasesPollardsKangaroo["small"].p)
		g  = bigInt(b, testCasesPollardsKangaroo["small"].g)
		a  = bigInt(b, testCasesPollardsKangaroo["small"].a)
		bb = bigInt(b, testCasesPollardsKangaroo["small"].b)
		y  = bigInt(b, testCasesPollardsKangaroo["small"].y)
		c  = bigInt(b, testCasesPollardsKangaroo["small"].c)
		k  = bigInt(b, testCasesPollardsKangaroo["small"].k)
		pm = NewPollardMapper(k, c, p)
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = PollardsKangaroo(pm, p, g, a, bb, y)
	}
}

func BenchmarkPollardMapper_F(b *testing.B) {
	var (
		k = bigInt(b, testCasesPollardsKangaroo["small"].k)
		c = bigInt(b, testCasesPollardsKangaroo["small"].c)
		p = bigInt(b, testCasesPollardsKangaroo["small"].p)
		y = bigInt(b, testCasesPollardsKangaroo["small"].y)
		m = NewPollardMapper(k, c, p)
		t big.Int
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.F(&t, y)
	}
}

func TestPollardsKangaroo(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.SkipNow()
	}
	for name, tc := range testCasesPollardsKangaroo {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			var (
				p    = bigInt(t, tc.p)
				g    = bigInt(t, tc.g)
				a    = bigInt(t, tc.a)
				b    = bigInt(t, tc.b)
				y    = bigInt(t, tc.y)
				c    = bigInt(t, tc.c)
				k    = bigInt(t, tc.k)
				want = bigInt(t, tc.want)
				pm   = NewPollardMapper(k, c, p)
			)

			// Can we compute the index of y?
			got, err := PollardsKangaroo(pm, p, g, a, b, y)
			if err != nil {
				t.Fatalf("failed to find index of y: %v", err)
			}

			// Was the index of y correct?
			if new(big.Int).Exp(g, got, p).Cmp(y) != 0 {
				t.Errorf("incorrect index: %v", got)
			}

			// Double-check against the saved value.
			if got.Cmp(want) != 0 {
				t.Errorf("got %v, want %v", got, want)
			}
		})
	}
}
