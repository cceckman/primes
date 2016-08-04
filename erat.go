package primes

import (
	"math"
)

// SimpleErat gives a simple / naive Sieve of Eratosthenes solver.
// It's very very simple- it even starts at 2.
type simpleErat struct{}

func (p *simpleErat) PrimesUpTo(n int, out chan<- int) {
	// No primes less than or equal to 1.
	if n <= 1 {
		close(out)
	}

	composite := make([]bool, n+1)
	for i := 2; i <= n; i++ {
		if composite[i] {
			continue
		}
		// it's still prime!
		// Write out
		out <- i
		// And sieve
		for k := i * 2; k < len(composite); k += i {
			composite[k] = true
		}
	}
	close(out)
}

func (p *simpleErat) IsPrime(n int) bool {
	if n <= 1 {
		return false
	}

	c := make(chan int)

	// It's no more expensive to compute all primes up to n
	// with a sieve of Eratosthenes, vs. just computing
	// whether n is prime.
	go p.PrimesUpTo(n, c)
	for p := range c {
		if p == n {
			return true
		}
	}
	return false
}

// erat2 is a Sieve of Eratosthenes with some opmtimizations,
// namely:
// - Discount multiples of 2 altogether
type erat2 struct{}

func erat2_num(i int) int { return (i * 2) + 1 }
func erat2_idx(n int) int { return (n - 1) / 2 }

func (p *erat2) PrimesUpTo(n int, out chan<- int) {
	if n <= 1 {
		close(out)
		return
	}
	if n == 2 {
		out <- 2
		close(out)
		return
	}
	out <- 2

	// In an odds-only slice,
	// index i refers to the number (i*2)+1;
	// number n is at index (n-1) / 2
	// TODO: Upper limit can be optimized to be something to do with sqrt(n).
	// prime = not-composite, until proven otherwise.
	composite := make([]bool, erat2_idx(n)+1)

	// Start at index 1 == number 3
	for i := 1; i < len(composite); i++ {
		if composite[i] { // non-default; has been explicitly set to be composite.
			continue
		}
		// Found a prime; record it...
		out <- erat2_num(i)

		// Perform sieve:
		// run through multiples of i, marking as composite.
		k := 3
		for {
			// The multiples of N are at
			// j := idx(num(i) * K)
			// for all K > 2.
			// But we cut out all even multiples by hopping to k+=2,
			// TODO optimize this instruction. This is readable, but mathier is probably faster.
			j := erat2_idx(erat2_num(i) * k)
			if j >= len(composite) {
				break
			}
			composite[j] = true
			k += 2
		} // end sieve
	}
	close(out)
}

func (p *erat2) IsPrime(n int) bool {
	if n <= 1 {
		return false
	}
	if n%2 == 0 {
		return false
	}

	c := make(chan int)

	// It's no more expensive to compute all primes up to n
	// with a sieve of Eratosthenes, vs. just computing
	// whether n is prime.
	go p.PrimesUpTo(n, c)
	for p := range c {
		if p == n {
			return true
		}
	}
	return false
}

// erat3 is like erat2, but also:
// Only looks for primes up to sqrt(n)
type erat3 struct{}

func erat3_num(i int) int { return (i * 2) + 1 }
func erat3_idx(n int) int { return (n - 1) / 2 }

func (p *erat3) PrimesUpTo(n int, out chan<- int) {
	if n <= 1 {
		close(out)
		return
	}
	if n == 2 {
		out <- 2
		close(out)
		return
	}
	out <- 2

	// In an odds-only slice,
	// index i refers to the number (i*2)+1;
	// number n is at index (n-1) / 2
	// prime = not-composite, until proven otherwise.
	composite := make([]bool, erat3_idx(n)+1)

	// Only need to look for primes "less than or equal to" sqrt(n)
	// before assuming all remaining (un-sieved) ones are prime
	sqrt := int(math.Ceil(math.Sqrt(float64(n))))

	// Start at index 1 == number 3
	for i := 1; i < len(composite); i++ {
		if composite[i] { // non-default; has been explicitly set to be composite.
			continue
		}
		// Found a prime; record it...
		out <- erat3_num(i)

		if erat3_num(i) > sqrt {
			// Skip sieving; we've covered all the primes already.
			continue
		}

		// run through multiples of i, marking as composite.
		k := 3
		for {
			// The multiples of N are at
			// j := idx(num(i) * K)
			// for all K > 2.
			// But we cut out all even multiples by hopping to k+=2,
			// TODO optimize this instruction. This is readable, but mathier is probably faster.
			j := erat3_idx(erat3_num(i) * k)
			if j >= len(composite) {
				break
			}
			composite[j] = true
			k += 2
		} // end sieve
	}
	close(out)
}

func (p *erat3) IsPrime(n int) bool {
	if n <= 1 {
		return false
	}
	if n%2 == 0 {
		return false
	}

	c := make(chan int)

	// It's no more expensive to compute all primes up to n
	// with a sieve of Eratosthenes, vs. just computing
	// whether n is prime.
	go p.PrimesUpTo(n, c)
	for p := range c {
		if p == n {
			return true
		}
	}
	return false
}

// erat4 is like erat3, but also:
// Use the actual values rather than indices for math (hopefully, reduce num ops)
type erat4 struct{}

func erat4_idx(n int) int { return (n - 1) / 2 }

func (p *erat4) PrimesUpTo(n int, out chan<- int) {
	if n <= 1 {
		close(out)
		return
	}
	if n == 2 {
		out <- 2
		close(out)
		return
	}
	out <- 2

	// In an odds-only slice,
	// index i refers to the number (i*2)+1;
	// number n is at index (n-1) / 2
	// prime = not-composite, until proven otherwise.
	composite := make([]bool, erat4_idx(n)+1)

	// Only need to look for primes "less than or equal to" sqrt(n)
	// before assuming all remaining (un-sieved) ones are prime
	sqrt := int(math.Ceil(math.Sqrt(float64(n))))

	// Start at index 1 == number 3
	for i := 1; i < n; i += 2 {
		if composite[(i - 1) / 2] { // non-default; has been explicitly set to be composite.
			continue
		}
		// Found a prime; record it...
		out <- i

		if i > sqrt {
			// Skip sieving; we've covered all the primes already.
			continue
		}

		// run through odd multiples of i, marking as composite.
		// Start with i * 3; add 2i each time.
		for j := i * 3; j <= n; j += (i + i) {
			composite[(j - 1) / 2] = true
		} // end sieve
	}
	close(out)
}

func (p *erat4) IsPrime(n int) bool {
	if n <= 1 {
		return false
	}
	if n%2 == 0 {
		return false
	}

	c := make(chan int)

	// It's no more expensive to compute all primes up to n
	// with a sieve of Eratosthenes, vs. just computing
	// whether n is prime.
	go p.PrimesUpTo(n, c)
	for p := range c {
		if p == n {
			return true
		}
	}
	return false
}


