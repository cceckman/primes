package primes

import (
	"testing"
)

func TestIteratorNext(t *testing.T) {
	db := New()
	it := db.Iterator()

	for i, p := range refPrimes {
		c := it.Next()
		if c != p {
			t.Error("unexpected prime #%d: got: %d want: %d", i, c, p)
		}
	}
}

func TestDB(t *testing.T) {
	db := New()
	max := refPrimes[len(refPrimes)-1]

	j := 0
	for n := 0; n < max; n++ {
		want := refPrimes[j] == n
		if want {
			j++
		}

		got := db.IsPrime(n)
		if got != want {
			t.Errorf("unexpected primacy for %d: got: %v want: %v", got, want)
		}
	}
}
