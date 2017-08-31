package test

import (
	"sync/atomic"
	"testing"
)

//go:noinline
func deferOn() {
	var x uint32
	atomic.CompareAndSwapUint32(&x, 0, 1)
	defer atomic.StoreUint32(&x, 0)
}

//go:noinline
func deferOff() {
	var x uint32
	atomic.CompareAndSwapUint32(&x, 0, 1)
	atomic.StoreUint32(&x, 0)
}

func BenchmarkAtomicDefer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		deferOn()
	}
}

func BenchmarkAtomicPlain(b *testing.B) {
	for i := 0; i < b.N; i++ {
		deferOff()
	}
}
