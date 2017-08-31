package test

import (
	"sync"
	"sync/atomic"
	"testing"
)

var x uint32

//go:noinline
func checkAtomic() {
	atomic.CompareAndSwapUint32(&x, 0, 1)
	atomic.StoreUint32(&x, 0)
}

var m sync.Mutex

//go:noinline
func checkMutex() {
	m.Lock()
	m.Unlock()
}

func BenchmarkMutex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		checkMutex()
	}
}

func BenchmarkAtomic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		checkAtomic()
	}
}
