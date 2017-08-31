package test

import (
	"sync"
	"testing"
)

var lock sync.Mutex
var rwLock sync.Mutex

//go:noinline
func checkLockMutex() {
	lock.Lock()
	lock.Unlock()
}

//go:noinline
func checkRWLock() {
	rwLock.Lock()
	rwLock.Unlock()
}

func BenchmarkLockMutex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		checkLockMutex()
	}
}

func BenchmarkRWLockMutex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		checkRWLock()
	}
}
