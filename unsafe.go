package primes

import (
	"math"
	"sort"
)

// DB is a thread-unsafe memoized list of prime numbers.
type DB struct {
	primes []int
}

func New() *DB {
	return &DB{
		primes: []int{2, 3, 5, 7},
	}
}

// Iterator is a thread-unsafe handle to iterate over a list of primes.
type Iterator struct {
	parent *DB
	index int
}

// Next returns the next prime number and advances the iterator.
func (i *Iterator) Next() int {
	if i.index == len(i.parent.primes) {
		// Need to grow the list.
		// Arbitrarily choose max*2 as the factor.
		max := i.parent.primes[i.index-1]
		i.parent.computeBeyond(2 * max)
	}

	i.index +=1
	return i.parent.primes[i.index-1]
}

// Iterator returns a new Iterator backed by this DB.
func (p *DB) Iterator() Iterator {
	return Iterator{
		parent: p,
	}
}

// IsPrime returns whether or not n is prime. It blocks until it can determine a result.
func (p *DB) IsPrime(n int) bool {
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

	// Ensure that we have enough in the list.
	p.computeBeyond(n)

	i := sort.SearchInts(p.primes, n)
	return i < len(p.primes) && p.primes[i] == n
}

// computeBeyond is a thread-unsave blocking call that returns once p has computed a prime greater than m.
func (p *DB) computeBeyond(m int) {
	// Do we know at least one prime beyond the requested number?
	max := p.primes[len(p.primes)-1]
	if m < max {
		return
	}


	// In order for callers to guarantee progress, we only want to
	// return if we've added primes to our list.
	for n, initLen := m, len(p.primes); len(p.primes) == initLen; n = n*2 + 1 {

		// In an odds-only slice,
		// index i refers to the number (i*2)+1;
		// number n is at index (n-1) / 2
		// (We could use less memory by subtracting out p.max.)
		// prime = not-composite, until proven otherwise.
		composite := make([]bool, (n/2)+1)
		composite[0] = true // 1 is not prime

		// First, mark composites from already-known primes.
		// Skip the first prime (2), since it isn't even in our array.
		for _, prime := range p.primes[1:] {
			// We could optimize this by starting at the first *odd* multiple of prime greater than or
			// equal to p.max. Instead, just start at prime*prime. (All smaller multiples of prime have
			// another prime factor, that is smaller than prime.)
			for i := prime * prime; i <= n; i += (prime * 2) {
				composite[(i-1)/2] = true
			}
		}

		// Only need to look for primes "less than or equal to" sqrt(n)
		// before assuming all remaining (un-sieved) ones are prime
		sqrt := int(math.Ceil(math.Sqrt(float64(n))))
		// Now, walk up, checking / marking composites along the way.
		// Start with p.max+1 or p.max+2, whichever is odd.
		oddMax := max + 1 + (max % 2)
		for i := oddMax; i <= n; i += 2 {
			if composite[(i-1)/2] { // non-default; has been explicitly set to be composite.
				continue
			}
			// Found a prime; record it.
			p.primes = append(p.primes, int(i))

			if i > sqrt {
				// Skip sieving; we've covered all the primes already.
				continue
			}

			// Found a new prime.
			// run through odd multiples of i, marking as composite.
			// Start with i * i; lower multiples of i will have already been marked as multiples
			// of another, smaller prime. Add 2i each time to ignore the even multiples.
			for j := i * i; j <= n; j += (i + i) {
				composite[(j-1)/2] = true
			} // end sieve
		}
	}
}
