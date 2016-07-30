package primes

import(
	"testing"
)

const(
	maxLevel = 5
)

// Generate levels of stressyness; how many zeros we want.
func levelGen(n int) []int {
	r := make([]int, n)
	r[0] = 101
	for i := range r {
		if i == 0 {
			continue
		}
		r[i] = (r[i-1] -1) * 10 + 1
	}
	return r
}

// Several test cases, parameterized on a primer.
func benchmarkPrimesUpTo(b *testing.B, p Primer, n int) {
	for i := 0; i < b.N; i++ {
		c := make(chan int)
		go p.PrimesUpTo(n, c)
		for _ = range c {}
	}
}

func BenchmarkSimpleErat(b *testing.B) {
	t := SimpleErat()
	for _, n := range levelGen(maxLevel) {
		benchmarkPrimesUpTo(b, t, n)
	}
}
