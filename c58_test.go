package set8

import "testing"

func TestPollardsKangaroo(t *testing.T) {
	t.Parallel()
	cases := map[string]struct {
		// Pollard's Kangaroo inputs
		p, g, a, b, y string
		// Constants for the Pollard mapper
		c, k string
		// Desired result
		want string
	}{
		"small": {
			p:    "11470374874925275658116663507232161402086650258453896274534991676898999262641581519101074740642369848233294239851519212341844337347119899874391456329785623",
			g:    "622952335333961296978159266084741085889881358738459939978290179936063635566740258555167783009058567397963466103140082647486611657350811560630587013183357",
			a:    "0",
			b:    "1048576", // 2^20
			y:    "7760073848032689505395005705677365876654629189298052775754597607446617558600394076764814236081991643094239886772481052254010323780165093955236429914607119",
			c:    "4",
			k:    "10",
			want: "705485",
		},
		"large": {
			p:    "11470374874925275658116663507232161402086650258453896274534991676898999262641581519101074740642369848233294239851519212341844337347119899874391456329785623",
			g:    "622952335333961296978159266084741085889881358738459939978290179936063635566740258555167783009058567397963466103140082647486611657350811560630587013183357",
			a:    "0",
			b:    "1099511627776", // 2^40
			y:    "9388897478013399550694114614498790691034187453089355259602614074132918843899833277397448144245883225611726912025846772975325932794909655215329941809013733",
			c:    "4",
			k:    "10",
			want: "1",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			var (
				p    = HelperBigIntFromString(t, tc.p)
				g    = HelperBigIntFromString(t, tc.g)
				a    = HelperBigIntFromString(t, tc.a)
				b    = HelperBigIntFromString(t, tc.b)
				y    = HelperBigIntFromString(t, tc.y)
				c    = HelperBigIntFromString(t, tc.c)
				k    = HelperBigIntFromString(t, tc.k)
				want = HelperBigIntFromString(t, tc.want)
			)
			pm, err := NewPollardMapper(k, c, p)
			if err != nil {
				t.Errorf("failed to create pollard mapper: %v", err)
			}
			got, err := PollardsKangaroo(p, g, a, b, y, pm)
			if err != nil {
				t.Errorf("failed to find index of y: %v", err)
			}
			if got.Cmp(want) != 0 {
				t.Errorf("got %v, want %v", got, want)
			}
		})
	}
}
