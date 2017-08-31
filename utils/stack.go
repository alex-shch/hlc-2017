package utils

import (
	"sync/atomic"
	"unsafe"

	"github.com/alex-shch/hlc-2017/store"
)

//
// bytes atomic stack
//
type BytesAtomicStackNode struct {
	Buf  []byte
	next unsafe.Pointer
}

type BytesAtomicStack struct {
	data []BytesAtomicStackNode

	head unsafe.Pointer

	lock uint32
}

func NewBytesAtomicStack(bufSize, stackCapacity int) *BytesAtomicStack {
	q := &BytesAtomicStack{
		data: make([]BytesAtomicStackNode, stackCapacity),
	}
	for i := 0; i < stackCapacity; i++ {
		q.data[i].Buf = make([]byte, 0, bufSize)
		if i > 0 {
			q.data[i].next = unsafe.Pointer(&q.data[i-1])
		}
	}
	q.head = unsafe.Pointer(&q.data[stackCapacity-1])

	return q
}

func (q *BytesAtomicStack) Pop() *BytesAtomicStackNode {
	for {
		head := q.head
		if head != nil && head == q.head {
			next := unsafe.Pointer((*BytesAtomicStackNode)(head).next)
			if atomic.CompareAndSwapPointer(&q.head, head, next) {
				n := (*BytesAtomicStackNode)(head)
				n.Buf = n.Buf[:0]
				return n
			}
		}
	}
}

func (q *BytesAtomicStack) Push(node *BytesAtomicStackNode) {
	for {
		head := q.head
		node.next = head
		if head == q.head {
			if atomic.CompareAndSwapPointer(&q.head, head, unsafe.Pointer(node)) {
				return
			}
		}
	}
}

//
// strings atomic stack
//
/*
type StringsAtomicStackNode struct {
	Buf  []string
	next unsafe.Pointer
}

type StringsAtomicStack struct {
	data []StringsAtomicStackNode

	head unsafe.Pointer

	lock uint32
}

func NewStringsAtomicStack(bufSize, stackCapacity int) *StringsAtomicStack {
	q := &StringsAtomicStack{
		data: make([]StringsAtomicStackNode, stackCapacity),
	}
	for i := 0; i < stackCapacity; i++ {
		q.data[i].Buf = make([]string, 0, bufSize)
		if i > 0 {
			q.data[i].next = unsafe.Pointer(&q.data[i-1])
		}
	}
	q.head = unsafe.Pointer(&q.data[stackCapacity-1])

	return q
}

func (q *StringsAtomicStack) Pop() *StringsAtomicStackNode {
	type NodePtr *StringsAtomicStackNode

	for {
		head := q.head
		if head != nil && head == q.head {
			next := unsafe.Pointer(NodePtr(head).next)
			if atomic.CompareAndSwapPointer(&q.head, head, next) {
				n := NodePtr(head)
				n.Buf = n.Buf[:0]
				return n
			}
		}
	}
}

func (q *StringsAtomicStack) Push(node *StringsAtomicStackNode) {
	for {
		head := q.head
		node.next = head
		if head == q.head {
			if atomic.CompareAndSwapPointer(&q.head, head, unsafe.Pointer(node)) {
				return
			}
		}
	}
}
*/

//
// visit array atomic stack
//
type VisitsAtomicStackNode struct {
	Array store.VisitArray
	next  unsafe.Pointer
}

type VisitsAtomicStack struct {
	data []VisitsAtomicStackNode

	head unsafe.Pointer

	lock uint32
}

func NewVisitsAtomicStack(bufSize, stackCapacity int) *VisitsAtomicStack {
	q := &VisitsAtomicStack{
		data: make([]VisitsAtomicStackNode, stackCapacity),
	}
	for i := 0; i < stackCapacity; i++ {
		q.data[i].Array = make(store.VisitArray, 0, bufSize)
		if i > 0 {
			q.data[i].next = unsafe.Pointer(&q.data[i-1])
		}
	}
	q.head = unsafe.Pointer(&q.data[stackCapacity-1])

	return q
}

func (q *VisitsAtomicStack) Pop() *VisitsAtomicStackNode {
	for {
		head := q.head
		if head != nil && head == q.head {
			next := unsafe.Pointer((*VisitsAtomicStackNode)(head).next)
			if atomic.CompareAndSwapPointer(&q.head, head, next) {
				n := (*VisitsAtomicStackNode)(head)
				n.Array = n.Array[:0]
				return n
			}
		}
	}
}

func (q *VisitsAtomicStack) Push(node *VisitsAtomicStackNode) {
	for {
		head := q.head
		node.next = head
		if head == q.head {
			if atomic.CompareAndSwapPointer(&q.head, head, unsafe.Pointer(node)) {
				return
			}
		}
	}
}
