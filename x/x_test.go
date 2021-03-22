package x

import (
	"math/big"
	"testing"
)

func HelperBigIntFromString(tb testing.TB, s string) *big.Int {
	tb.Helper()
	n, ok := new(big.Int).SetString(s, 10)
	if !ok {
		tb.Errorf("cannot convert to big.Int: %s", s)
	}
	return n
}

func BenchmarkExp_InPlace(b *testing.B) {
	g := HelperBigIntFromString(b, "622952335333961296978159266084741085889881358738459939978290179936063635566740258555167783009058567397963466103140082647486611657350811560630587013183357")
	y := HelperBigIntFromString(b, "7760073848032689505395005705677365876654629189298052775754597607446617558600394076764814236081991643094239886772481052254010323780165093955236429914607119")
	p := HelperBigIntFromString(b, "11470374874925275658116663507232161402086650258453896274534991676898999262641581519101074740642369848233294239851519212341844337347119899874391456329785623")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.Exp(g, y, p)
	}
}

func BenchmarkExp_NotInPlace(b *testing.B) {
	g := HelperBigIntFromString(b, "622952335333961296978159266084741085889881358738459939978290179936063635566740258555167783009058567397963466103140082647486611657350811560630587013183357")
	y := HelperBigIntFromString(b, "7760073848032689505395005705677365876654629189298052775754597607446617558600394076764814236081991643094239886772481052254010323780165093955236429914607119")
	p := HelperBigIntFromString(b, "11470374874925275658116663507232161402086650258453896274534991676898999262641581519101074740642369848233294239851519212341844337347119899874391456329785623")
	var t big.Int

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		t.Exp(g, y, p)
	}
}