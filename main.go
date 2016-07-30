package main

import (
	"flag"
	"fmt"
	"github.com/cceckman/primes"
	"os"
	"strings"
	"text/tabwriter"
	"testing"
)

var (
	maxLevel = flag.Int("max_level", 5, "How far to run benchmarks: up to 1..1 with this many zeros in the middle.")
	// human = flag.Bool("human", false, "Human-formatted output.") // TODO implement.
	help = flag.Bool("help", false, "Display a usage message.")
)

// Benchmark the performance of different prime algorithms.
// Output as TSV
func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s: Benchmark performance of prime algorithms.\nUsage:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	if *help {
		flag.Usage()
		os.Exit(-1)
	}
	// TODO: I should really put all the above in a template for main.go.


	w := tabwriter.NewWriter(os.Stdout, 0, 8, 0, '\t', 0)

	evaluators := map[string]primes.Primer{
		"SimpleErat": primes.SimpleErat(),
	}
	levels := levelGen(*maxLevel)

	// Output header
	fmt.Fprint(w, "Function\tAlgorithm\tParameter\tIterations\tTotal Time (s)\tAverage Time (ns)\tAllocs\tBytes\n")

	// First: bench PrimesUpTo, to get a set of primes to use
	for name, primer := range evaluators {
		for _, level := range levels {
			// create a benchmark that closes over level and primer
			// Not running in parallel; OK to use those variables
			b := func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					c := make(chan int)
					go primer.PrimesUpTo(level, c)
					for _ = range c {
					}
				}
			}
			// Run the benchmark
			result := testing.Benchmark(b)
			// And write it out

			fmt.Fprintln(w, strings.Join([]string{
				"PrimesUpTo",
				name,
				fmt.Sprint(level),
				fmt.Sprint(result.N),
				fmt.Sprint(result.T),
				fmt.Sprint(result.NsPerOp()),
				result.MemString(),
			}, "\t"))
		}
	}

	// Then: bench IsPrime

	w.Flush()
}

// Generate levels of stressyness; how many zeros we want.
func levelGen(n int) []int {
	r := make([]int, n)
	r[0] = 101
	for i := range r {
		if i == 0 {
			continue
		}
		r[i] = (r[i-1]-1)*10 + 1
	}
	return r
}
