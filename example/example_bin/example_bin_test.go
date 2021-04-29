package example_bin

import (
	"testing"
)

func Test(t *testing.T) {
	res := getName("Red", "Stone")
	if res != "Red Stone Team" {
		t.Error("Error", res)
	}
}

// This is the XorShift random generator
type XorShiftRng struct {
	x, y, z, w uint32
}

// The code is from [here](https://docs.rs/rand/0.5.0/src/rand/prng/xorshift.rs.html#29-34)
func (rng *XorShiftRng) next_u32() uint32 {
	x := rng.x
	t := x ^ (x << 11)
	rng.x = rng.y
	rng.y = rng.z
	rng.z = rng.w
	w := rng.w
	rng.w = w ^ (w >> 19) ^ (t ^ (t >> 8))
	return rng.w
}

func BenchmarkSomething(b *testing.B) {
	// you may ask, why I implemented it here?
	// I dunno, just for fun :)
	rng := XorShiftRng{
		x: 0xbad5eed,
		y: 0xbad5eed,
		z: 0xbad5eed,
		w: 0xbad5eed,
	}
	for n := 0; n < b.N; n++ {
		_ = rng.next_u32()
	}
}
