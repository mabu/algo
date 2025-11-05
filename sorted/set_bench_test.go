package sorted

import (
	"math/rand/v2"
	"testing"
)

func benchmarkSet() *Set[int] {
	return NewSet[int]()
	// return NewSetFunc(func(a, b int) int { return a - b })
}

const n int = 1e5

var (
	rnd          = rand.New(rand.NewPCG(1, 2))
	permutation1 = rnd.Perm(n)
	permutation2 = rnd.Perm(n)
	random       []int
)

func init() {
	for range n {
		random = append(random, rnd.Int())
	}
}

func BenchmarkInsert(b *testing.B) {
	for b.Loop() {
		s := benchmarkSet()
		for _, v := range random {
			s.Insert(v)
		}
	}
}

func BenchmarkAll(b *testing.B) {
	s := benchmarkSet()
	for _, v := range permutation1 {
		s.Insert(v)
	}
	for b.Loop() {
		for range s.All() {
		}
	}
}

func BenchmarkDelete(b *testing.B) {
	for b.Loop() {
		b.StopTimer()
		s := benchmarkSet()
		for _, v := range permutation1 {
			s.Insert(v)
		}
		b.StartTimer()
		for _, v := range permutation2 {
			s.Delete(v)
		}
	}
}
