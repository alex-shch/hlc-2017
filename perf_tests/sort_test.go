package test

import (
	"math/rand"
	"sort"
	"testing"
)

func shakerSort(a []int) {
	left := 0
	right := len(a) - 1

	for left <= right {
		for i := left; i < right; i++ {
			j := i + 1
			if a[i] > a[j] {
				a[i], a[j] = a[j], a[i]
			}
		}
		right -= 1

		for i := right; i > left; i-- {
			j := i - 1
			if a[i] < a[j] {
				a[i], a[j] = a[j], a[i]
			}
		}
		left++
	}
}

func TestShakerSort(t *testing.T) {
	a := [10]int{}
	for i := 0; i < len(a); i++ {
		a[i] = rand.Intn(20)
	}

	b := make(sort.IntSlice, len(a))
	copy(b, a[:])

	t.Logf("\n%v\n%v\n", a[:], b)

	shakerSort(a[:])
	sort.Sort(b)
	t.Logf("\n%v\n%v\n", a[:], b)

	for i := 0; i < len(a); i++ {
		if expected, val := b[i], a[i]; val != expected {
			t.Fatalf("idx:%d, %d != %d", i, val, expected)
		}
	}
}

func benchmarkQSort(size int, b *testing.B) {
	a := make(sort.IntSlice, size)

	for i := 0; i < b.N; i++ {
		//b.StopTimer()
		for i := 0; i < len(a); i++ {
			a[i] = rand.Intn(20)
		}
		//b.StartTimer()

		sort.Sort(a)
	}
}

func benchmarkShakerSort(size int, b *testing.B) {
	a := make([]int, size)

	for i := 0; i < b.N; i++ {
		//b.StopTimer()
		for i := 0; i < len(a); i++ {
			a[i] = rand.Intn(20)
		}
		//b.StartTimer()

		shakerSort(a)
	}
}

func BenchmarkSort(b *testing.B) {
	b.Run("qsort 3", func(b *testing.B) { benchmarkQSort(3, b) })
	b.Run("qsort 4", func(b *testing.B) { benchmarkQSort(4, b) })
	b.Run("qsort 5", func(b *testing.B) { benchmarkQSort(5, b) })
	b.Run("qsort 7", func(b *testing.B) { benchmarkQSort(7, b) })
	b.Run("qsort 10", func(b *testing.B) { benchmarkQSort(10, b) })
	b.Run("qsort 15", func(b *testing.B) { benchmarkQSort(15, b) })
	b.Run("qsort 20", func(b *testing.B) { benchmarkQSort(20, b) })
	b.Run("qsort 25", func(b *testing.B) { benchmarkQSort(25, b) })
	b.Run("qsort 30", func(b *testing.B) { benchmarkQSort(30, b) })
	b.Run("qsort 35", func(b *testing.B) { benchmarkQSort(35, b) })
	b.Run("qsort 40", func(b *testing.B) { benchmarkQSort(40, b) })
	b.Run("qsort 50", func(b *testing.B) { benchmarkQSort(50, b) })

	b.Run("shakersort 3", func(b *testing.B) { benchmarkShakerSort(3, b) })
	b.Run("shakersort 4", func(b *testing.B) { benchmarkShakerSort(4, b) })
	b.Run("shakersort 5", func(b *testing.B) { benchmarkShakerSort(5, b) })
	b.Run("shakersort 7", func(b *testing.B) { benchmarkShakerSort(7, b) })
	b.Run("shakersort 10", func(b *testing.B) { benchmarkShakerSort(10, b) })
	b.Run("shakersort 15", func(b *testing.B) { benchmarkShakerSort(15, b) })
	b.Run("shakersort 20", func(b *testing.B) { benchmarkShakerSort(20, b) })
	b.Run("shakersort 25", func(b *testing.B) { benchmarkShakerSort(25, b) })
	b.Run("shakersort 30", func(b *testing.B) { benchmarkShakerSort(30, b) })
	b.Run("shakersort 35", func(b *testing.B) { benchmarkShakerSort(35, b) })
	b.Run("shakersort 40", func(b *testing.B) { benchmarkShakerSort(40, b) })
	b.Run("shakersort 50", func(b *testing.B) { benchmarkShakerSort(50, b) })
}
