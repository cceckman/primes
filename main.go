package main

import(
	"flag"
	"fmt"
	"testing"

	"github.com/cceckman/primes"
)

var(
	maxLevel = flag.Int("max_level", 5, "How far to run benchmarks: up to 1..1 with this many zeros in the middle.")
	primes []int
)

// Benchmark the performance of different prime algorithms.
// Output as TSV
func main() {
	evaluators := map[string]primes.Primer{
		"SimpleErat": primes.SimpleErat(),
	}
	levels := levelGen(*maxLevel)

	// Output header
	fmt.Print("Function\tAlgorithm\tParameter\tResults\n")

	// First: bench PrimesUpTo, to get a set of primes to use
	for name, primer := range evaluators {
		for l, level := range levels {
			// create a benchmark that closes over level and primer
			// Not running in parallel; OK to use those variables
			b := func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					c := make(chan int)
					go primer.PrimesUpTo(level, c)
					for _ = range c {}
				}
			}
			// Run the benchmark
			result := testing.Benchmark(b)
			// And write it out
			formatOutput("PrimesUpTo", name, level, result.String())
		}
	}

	// Then: bench IsPrime
}

func formatOutput(function, primer string, parameter int, result string) string {
	buf := bytes.NewBuffer()
	fmt.Fprintf(
		buf, "%s\t%s\t%d\n%s",
		function, primer, parameter, result
	)

	return buf.String()
}

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
