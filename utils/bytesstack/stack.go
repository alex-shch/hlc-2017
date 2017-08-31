package bytesstack

import (
	"sync/atomic"
	"unsafe"
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
