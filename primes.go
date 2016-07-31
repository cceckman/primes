// Prime generator interface.
package primes

import (
	"math"
)

var (
	Implementations = map[string]Primer{
		"SimpleErat": &simpleErat{},
		"Erat2":      &erat2{},
		"Erat3":      &erat3{},
	}
)

// TODO use math/big?

// Primer provides functionality around prime numbers.
// It may be implemented with more-or-less efficient algorithms.
type Primer interface {
	// PrimesUpTo sends to `out`, in order, the primes up to (possibly including) `n`.
	// `out` is provided as a parameter so that the caller may set e.g. buffers.
	// `out` is closed once all primes have been returned.
	PrimesUpTo(n int, out chan<- int)
	// IsPrime tests whether `n` is a prime number.
	IsPrime(n int) bool
}

// Alternative, easier to test: collect into an array.
func PrimesUpTo(n int, p Primer) []int {
	// Use pi(x) ~ x / log x to estimate capacity.
	// OK to be only approximate for capacity; append will allocate more if needed.
	est := n / int(math.Log(float64(n)))

	result := make([]int, est, 0)
	c := make(chan int, est)
	// allow buffering within the channel- don't block computation if we can help it.
	// Not likely, since all we're doing is an append, but whatever.

	go p.PrimesUpTo(n, c)
	for prime := range c {
		result = append(result, prime)
	}
	return result
}
