package test

import (
	//"runtime"
	"sync/atomic"
	"testing"
)

type IntNode struct {
	val int64
}

type PtrArray []*IntNode

func (a PtrArray) Len() int           { return len(a) }
func (a PtrArray) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a PtrArray) Less(i, j int) bool { return a[i].val < a[j].val }

type _ArrayStorage struct {
	data PtrArray
	lock uint32
}

func (self *_ArrayStorage) add(val *IntNode) {
	_waitForSet(&self.lock)

	self.data = append(self.data, val)

	atomic.StoreUint32(&self.lock, 0)
}

func _waitForSet(lock *uint32) {
	locked := atomic.CompareAndSwapUint32(lock, 0, 1)

	for !locked {
		// ждем завершения предыдущей вставки
		if *lock == 0 {
			if atomic.CompareAndSwapUint32(lock, 0, 1) {
				return
			}
		}
	}
}

func (self *_ArrayStorage) UpdateGo(val *IntNode) {
	go self.add(val)
}

func (self *_ArrayStorage) UpdateSelf(val *IntNode) {
	self.add(val)
}

func BenchmarkUpdateGo(b *testing.B) {
	storage := &_ArrayStorage{
		data: make(PtrArray, 0, b.N),
	}

	node := &IntNode{
		val: 7,
	}

	for i := 0; i < b.N; i++ {
		storage.UpdateGo(node)
	}
}

func BenchmarkUpdateSelf(b *testing.B) {
	storage := &_ArrayStorage{
		data: make(PtrArray, 0, b.N),
	}

	node := &IntNode{
		val: 7,
	}

	for i := 0; i < b.N; i++ {
		storage.UpdateSelf(node)
	}
}
