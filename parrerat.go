package primes

import (
	"os"
	"runtime"
	"reflect"
)

// ParrErat is a prallelized, segmented Eratosthenes sieve.
// It must be instantiated with
// - How big the segments should be, in # of ints. "related to memory" is a good size.
// - How many workers to run. "Number of cores" is a good size.
type parrErat struct{
	segmentSize uintptr // in # of ints
	workers int
}

func NewParrErat() Primer {
	pageSize := uintptr(os.Getpagesize())
	intSize := reflect.TypeOf(int(0)).Size()
	// How many ints we can store in a single page
	segmentSize := pageSize / intSize

	return &parrErat{
		segmentSize: segmentSize,
		workers: runtime.GOMAXPROCS(-1),
	}
}

func m2_num(i int) int { return (i * 2) + 1 }
func m2_idx(n int) int { return (n - 1) / 2 }

func (p *parrErat) PrimesUpTo(n int, out chan<- int) {
	if n <= 1 {
		close(out)
		return
	}

	// Set up memory regions.
	// We need n+1 locations to store numbers up to n,
	// because we use n as the index.
	nSegments := uintptr(n+1) / p.segmentSize
	if uintptr(n+1) / p.segmentSize != 0 {
		nSegments += 1
	}

	// TODO more efficient allocation here, using a larger block?
	composite := make([][]bool, nSegments)
	for s := 0; s < len(composite); s++ {
		composite[s] = make([]bool, p.segmentSize)
	}

	for i := 2; i < n; i++ {
	}

	close(out)
}

func (p *parrErat) IsPrime(n int) bool {
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


