package test

import (
	"testing"
)

func BenchmarkArrayIndex(b *testing.B) {
	a := make([]int, b.N)
	size := 0

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		a[size] = i
		size++
	}
}

func BenchmarkArrayAppend(b *testing.B) {
	a := make([]int, 0, b.N)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		a = append(a, i)
	}
}
