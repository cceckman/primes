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

	composite := make([]bool, n + 1)
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
	composite := make([]bool, n + 1)
	for i := 2; i <= n; i++ {
		if composite[i] {
			continue
		}
		// it's still prime!
		// Sieve...
		for k := i * 2; k < len(composite); k += i {
			composite[k] = true
		}
	}
	return ! composite[n]
}
