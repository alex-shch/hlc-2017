package store

import (
	"sync/atomic"
)

func waitForSet(lock *uint32) {
	for {
		if *lock == 0 {
			if atomic.CompareAndSwapUint32(lock, 0, 1) {
				break
			}
		}
	}
}
