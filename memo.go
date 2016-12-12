//
// memo.go
// Copyright (C) 2016 cceckman <charles@cceckman.com>
//
// Distributed under terms of the MIT license.
//

package primes

import(
	"math"
)

var (
	_ Primer = NewMemoizingPrimer()
)

// MemoizingPrimer is primer that stores found primes.
// It is not (yet) threadsafe.
type MemoizingPrimer struct {
	// max is the largest number checked for primacy
	max int64
	// listing is the list of all primes found below max.
	listed []int
}

func NewMemoizingPrimer() *MemoizingPrimer {
	p := &MemoizingPrimer{
		max: 10,
		listed: []int{2, 3, 5, 7},
	}
	return p
}

// IsPrime returns whether or not n is prime. It blocks until it can determine a result.
func (p *MemoizingPrimer) IsPrime(n int) bool {
	// Quick answers.
	if n <= 1 {
		return false
	}
	if n == 2 {
		return true
	}
	if n%2 == 0 {
		return false
	}

	// Iterate through the primes and try to find it.
	c := make(chan int)
	p.PrimesUpTo(n, c)
	for i := range c {
		if i == n {
			// Get rid of the rest of the values in the background.
			backgroundFlush(c)
			return true
		}
	}

	// Better way: do a binary search in the list.
	// Log rather than linear.

	return false
}

// PrimesUpTo streams all the primes up to n, and closes 'out' when complete.
// It's non-blocking.
func (p *MemoizingPrimer) PrimesUpTo(n int, out chan<- int) {
	p.computeUpTo(n)
	// We have now asserted we're caught up.

	go func() {
		for _, v := range p.listed {
			if v <= n {
				out <- v
			}
		}
		close(out)
	}()
}

// computeUpTo is a blocking call that returns once p has computed primes up to n.
func (p *MemoizingPrimer) computeUpTo(n int) {
	if n <= int(p.max) {
		return
	}

	// In an odds-only slice,
	// index i refers to the number (i*2)+1;
	// number n is at index (n-1) / 2
	// (We could use less memory by subtracting out p.max.)
	// prime = not-composite, until proven otherwise.
	composite := make([]bool, (n / 2) + 1)
	composite[0] = true  // 1 is not prime

	// First, mark composites from already-known primes.
	// Skip the first prime (2), since it isn't even in our array.
	for _, prime := range p.listed[1:] {
		// We could optimize this by starting at the first *odd* multiple of prime greater than or
		// equal to p.max. Instead, just start at prime*prime. (All smaller multiples of prime have
		// another prime factor, that is smaller than prime.)
		for i := prime * prime; i <= n; i += (prime * 2) {
			composite[(i - 1) / 2] = true
		}
	}

	// Only need to look for primes "less than or equal to" sqrt(n)
	// before assuming all remaining (un-sieved) ones are prime
	sqrt := int(math.Ceil(math.Sqrt(float64(n))))
	// Now, walk up, checking / marking composites along the way.
	// Start with p.max or p.max-1, whichever is odd.
	oddMax := int(p.max) - 1 + (int(p.max) % 2)
	for i := oddMax; i <= n; i += 2 {
		if composite[(i - 1) / 2] { // non-default; has been explicitly set to be composite.
			continue
		}
		// Found a prime; record it.
		p.listed = append(p.listed, int(i))

		if i > sqrt {
			// Skip sieving; we've covered all the primes already.
			continue
		}

		// Found a new prime.
		// run through odd multiples of i, marking as composite.
		// Start with i * i; lower multiples of i will have already been marked as multiples
		// of another, smaller prime. Add 2i each time to ignore the even multiples.
		for j := i * i; j <= n; j += (i + i) {
			composite[(j - 1) / 2] = true
		} // end sieve
	}
	// Finally, update max.
	p.max = int64(n)
}


// backgroundFlush starts a thread that flushes c.
func backgroundFlush(c <-chan int) {
	go func() {
		for _ = range c {}
	}()
}
