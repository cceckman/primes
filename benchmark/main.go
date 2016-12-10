package main

import (
	"flag"
	"fmt"
	"github.com/cceckman/primes"
	"github.com/cceckman/bencher"
	"os"
)

var (
	maxLevel = flag.Int("max_level", 5, "How far to run benchmarks: up to 1..1 with this many zeros in the middle.")
	help = flag.Bool("help", false, "Display a usage message.")
)

type printInt int
func (p printInt) String() string {
	return fmt.Sprintf("%d", int(p))
}

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

	levels := levelGen(*maxLevel)

	closures := make(bencher.Cases)

	// benchmark PrimesUpTo, to get a set of primes to use
	for name, primer := range primes.Implementations {
		for _, level := range levels {
			// Copy locally.
			primer := primer
			level := level

			caseName := fmt.Sprintf("%s: %s(%d)", "PrimesUpTo", name, level)

			closure := func() fmt.Stringer {
				c := make(chan int)
				go primer.PrimesUpTo(level, c)
				x := 0
				for x = range c { }
				return printInt(x)
			}

			closures[caseName] = bencher.Runnable(closure)
		}
	}

	bencher.AutoBenchmark(closures)
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
