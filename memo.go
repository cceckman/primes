//
// memo.go
// Copyright (C) 2016 cceckman <charles@cceckman.com>
//
// Distributed under terms of the MIT license.
//

package primes

import(
	"math"
	"sync"
	"sync/atomic"
)

var (
	_ Primer = NewMemoizingPrimer()
)

// MemoizingPrimer is threadsafe, remembered Primer.
type MemoizingPrimer struct {
	// max is the largest number checked for primacy
	max int64
	// listing is the list of all primes found below max.
	listed []int
	lock sync.RWMutex
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
	// Quick answer...
	if n % 2 == 0 {
		return false
	}

	// Otherwise, see if we're up to that value yet.
	max := atomic.LoadInt64(&p.max)

	// Need to compute up to and including n; see if it shows up in the result.
	if max <= int64(n) {
		c := make(chan int)
		p.PrimesUpTo(n+1, c)
		for i := range c {
			if i == n {
				backgroundFlush(c)
				return true
			}
			if i > n {
				backgroundFlush(c)
				return false
			}
		}
	}

	// Have already computed n; do a binary search.
	p.lock.RLock()
	defer p.lock.RUnlock()

	lower, upper := 0, len(p.listed) - 1
	pivot := (upper - lower) / 2 + lower
	for {
		pivot = (upper - lower) / 2 + lower

		if p.listed[pivot] == n {
			return true
		}

		if lower >= upper {
			break
		}

		if p.listed[pivot] > n {
			lower = pivot + 1
		} else if p.listed[pivot] < n {
			upper = pivot - 1
		}
	}
	return false
}

// PrimesUpTo streams all the primes up to n, and closes 'out' when complete.
// It's non-blocking.
func (p *MemoizingPrimer) PrimesUpTo(n int, out chan<- int) {
	go func() {
		if int64(n) >= atomic.LoadInt64(&p.max) {
			// *May* need to compute more, so head into the write-locked section.
			p.computeUpTo(n)
		}
		// We have now asserted we're caught up, so take the just-read lock.

		p.lock.RLock()
		defer p.lock.RUnlock()
		for _, v := range p.listed {
			out <- v
		}
		close(out)
	}()
}

// computeUpTo is a blocking call that returns once p has computed primes up to n.
func (p *MemoizingPrimer) computeUpTo(n int) {
	p.lock.Lock()
	defer p.lock.Unlock()

	if int64(n) < p.max {
		// Already computed (we didn't have the RWlock last time we checked.) Return immediately.
		return
	}

	// In an odds-only slice,
	// index i refers to the number (i*2)+1 + p.max;
	// number n is at index (n-1 - p.max) / 2
	// prime = not-composite, until proven otherwise.
	composite := make([]bool, int64(n) / 2 + 1 - p.max)

	// Only need to look for primes "less than or equal to" sqrt(n)
	// before assuming all remaining (un-sieved) ones are prime
	sqrt := int64(math.Ceil(math.Sqrt(float64(n))))

	// First, mark composites from known primes already.
	for _, pr := range p.listed {
		prime := int64(pr)
		if prime == 2 {
			continue
		}
		// Skip those that would already have been covered by p.max;
		// so, start from the first multiple of prime that is greater than p.max
		factor := (p.max / prime) + 1

		for i := prime * factor; i < sqrt; i += (prime * 2) {
			composite[(i - 1 - p.max) / 2] = true
		}
	}

	// Now, walk up, checking / marking composites along the way.
	for i := p.max + 1; i < int64(n); i += 2 {
		if composite[(i - 1) / 2] { // non-default; has been explicitly set to be composite.
			continue
		}
		// Found a prime; record it.
		p.listed = append(p.listed, int(i))

		if i > sqrt {
			// Skip sieving; we've covered all the primes already.
			continue
		}

		// run through odd multiples of i, marking as composite.
		// Start with i * i; lower multiples of i will have already been marked as multiples
		// of another, smaller prime. Add 2i each time to ignore the even multiples.
		for j := i * i; j <= int64(n); j += (i + i) {
			composite[(j - 1 - p.max) / 2] = true
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
