// Test the primes package.
package primes

import (
	"sync"
	"reflect"
	"testing"
)

var (
	refPrimes = []int{
		2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97,
		101, 103, 107, 109, 113, 127, 131, 137, 139, 149, 151, 157, 163, 167, 173, 179, 181, 191, 193,
		197, 199, 211, 223, 227, 229, 233, 239, 241, 251, 257, 263, 269, 271, 277, 281, 283, 293, 307,
		311, 313, 317, 331, 337, 347, 349, 353, 359, 367, 373, 379, 383, 389, 397, 401, 409, 419, 421,
		431, 433, 439, 443, 449, 457, 461, 463, 467, 479, 487, 491, 499, 503, 509, 521, 523, 541, 547,
		557, 563, 569, 571, 577, 587, 593, 599, 601, 607, 613, 617, 619, 631, 641, 643, 647, 653, 659,
		661, 673, 677, 683, 691, 701, 709, 719, 727, 733, 739, 743, 751, 757, 761, 769, 773, 787, 797,
		809, 811, 821, 823, 827, 829, 839, 853, 857, 859, 863, 877, 881, 883, 887, 907, 911, 919, 929,
		937, 941, 947, 953, 967, 971, 977, 983, 991, 997, 1009, 1013, 1019, 1021, 1031, 1033, 1039,
		1049, 1051, 1061, 1063, 1069, 1087, 1091, 1093, 1097, 1103, 1109, 1117, 1123, 1129, 1151, 1153,
		1163, 1171, 1181, 1187, 1193, 1201, 1213, 1217, 1223,
	}
)

func TestPrimesUpTo(t *testing.T) {
	max := refPrimes[len(refPrimes)-1]
	for name, p := range Implementations {
		got := PrimesUpTo(max, p)
		if !reflect.DeepEqual(got, refPrimes) {
			t.Errorf("Got incorrect result for Primer %s: got: %v wanted: %v",
				name, got, refPrimes,
			)
		}
	}
}

func TestIsPrime(t *testing.T) {
	max := refPrimes[len(refPrimes)-1]
	min := -10
	pointer := 0 // into refPrimes

	// Probably more efficient to keep the bigger loop on the outside.
	for i := min; i < max; i++ {
		want := i == refPrimes[pointer]
		for name, p := range Implementations {
			// i == refPrimes[pointer] means "is this prime".
			got := p.IsPrime(i)
			if got != want {
				t.Errorf("Got incorrect result for Primer %s on value %v: got: %v wanted: %v",
					name, i, got, want,
				)
			}
		}
		if i == refPrimes[pointer] { // We've passed this prime; move to the next.
			pointer++
		}
	}
}

// TestIsPrimeParallel spawns a new thread for each (i, impl) instance.
// It gives some coverage of parallelism, though no guarantees.
func TestIsPrimeParallel(t *testing.T) {
	max := refPrimes[len(refPrimes)-1]
	min := -10
	pointer := 0 // into refPrimes

	var wg sync.WaitGroup

	// Probably more efficient to keep the bigger loop on the outside.
	for i := min; i < max; i++ {
		want := i == refPrimes[pointer]
		for name, p := range Implementations {
			// Idiomatic override of loop variables.
			name := name
			p := p
			i := i
			wg.Add(1)
			go func() {
				// i == refPrimes[pointer] means "is this prime".
				got := p.IsPrime(i)
				if got != want {
					t.Errorf("Got incorrect result for Primer %s on value %v: got: %v wanted: %v",
						name, i, got, want,
					)
				}
			}()
		}
		if i == refPrimes[pointer] { // We've passed this prime; move to the next.
			pointer++
		}
	}
	wg.Done()
}

func TestRegression(t *testing.T) {
	mem := NewMemoizingPrimer()

	r := mem.IsPrime(1683)
	if r {
		t.Errorf("for %d: got: %v want: %v", 1683, r, false)
	}
	r = mem.IsPrime(1765)
	if r {
		t.Errorf("for %d: got: %v want: %v", 1765, r, false)
	}

	mem = NewMemoizingPrimer()
	r = mem.IsPrime(11)
	if !r {
		t.Errorf("for %d: got: %v want: %v", 11, r, true)
	}
	r = mem.IsPrime(1765)
	if r {
		t.Errorf("for %d: got: %v want: %v", 1765, r, false)
	}

	mem = NewMemoizingPrimer()
	r = mem.IsPrime(12)
	if r {
		t.Errorf("for %d: got: %v want: %v", 12, r, false)
	}
	r = mem.IsPrime(1765)
	if r {
		t.Errorf("for %d: got: %v want: %v", 1765, r, false)
	}
	// What we find from the above: MemoizingPrimer is wrong when it
	// previously has an odd max. Also, note that when using only
	// IsPrime, p.max will always be odd- PrimesUpTo doesn't get invoked
	// from IsPrime if n is even.

	mem = NewMemoizingPrimer()
	c := make(chan int)
	mem.PrimesUpTo(11, c)
	for _ = range c {}
	r = mem.IsPrime(1765)
	if r {
		t.Errorf("for %d: got: %v want: %v", 1765, r, false)
	}

	mem = NewMemoizingPrimer()
	c = make(chan int)
	mem.PrimesUpTo(12, c)
	for _ = range c {}
	r = mem.IsPrime(1765)
	if r {
		t.Errorf("for %d: got: %v want: %v", 1765, r, false)
	}
	// And what we find from the above... both of them fail, so it's not simply a question of p.max
	// even/odd beforehand.
	mem = NewMemoizingPrimer()
	r = mem.IsPrime(1765)
	if r {
		t.Errorf("for %d: got: %v want: %v", 1765, r, false)
	}
	// But: the case above succeeds!

}

