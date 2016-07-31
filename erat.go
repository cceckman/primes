package primes

var (
	_ Primer = &simpleErat{}
)

// SimpleErat gives a simple / naive Sieve of Eratosthenes solver.
// It's very very simple- it even starts at 2.
func SimpleErat() Primer {
	return &simpleErat{}
}

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
	go p.PrimesUpTo(n,c)
	for p := range c {
		if p == n {
			return true
		}
	}
	return false
}

func Erat2() Primer {
	return &erat2{}
}

// erat is a Sieve of Eratosthenes with some opmtimizations,
// most notably only looking at odd numbers.
type erat2 struct{}

func (p *erat2) PrimesUpTo(n int, out chan<- int) {
	if n <= 1 {
		close(out)
		return
	}

	// In an odds-only slice,
	// index i refers to the number (i*2)+1;
	// number n is at index (n-1) / 2
	num := func(i int) int { return (i * 2) + 1 }
	idx := func(n int) int { return (n - 1) / 2 }

	// TODO: Upper limit can be optimized to be something to do with sqrt(n).
	// prime = not-composite, until proven otherwise.
	composite := make([]bool, idx(n) + 1)

	// Start at index 1 == number 3
	for i := 1; i < len(composite); i++ {
		if composite[i] { // non-default; has been explicitly set to be composite.
			continue
		}
		// Found a prime; record it...
		out <- num(i)

		// Perform sieve:
		// run through multiples of i, marking as composite.
		k := 3
		for {
			// The multiples of N are at
			// j := idx(num(i) * K)
			// for all K > 2.
			// But we cut out all even multiples by hopping to k+=2,
			// TODO optimize this instruction. This is readable, but mathier is probably faster.
			j := idx(num(i) * k)
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
	if n % 2 == 0 {
		return false
	}

	c := make(chan int)

	// It's no more expensive to compute all primes up to n
	// with a sieve of Eratosthenes, vs. just computing
	// whether n is prime.
	go p.PrimesUpTo(n,c)
	for p := range c {
		if p == n {
			return true
		}
	}
	return false
}
