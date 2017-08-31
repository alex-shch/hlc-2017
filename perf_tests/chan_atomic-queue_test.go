package test

import (
	"log"
	"sync/atomic"
	"testing"
	"time"
	"unsafe"
)

//
// _ChanQueue
//
type _ChanQueue struct {
	bufs chan []byte
}

func _NewChanQueue(size int) *_ChanQueue {
	q := &_ChanQueue{
		bufs: make(chan []byte, size),
	}

	for i := 0; i < size; i++ {
		q.bufs <- make([]byte, 0, 512) // TODO передать размер параметром
	}

	return q
}

func (self *_ChanQueue) Pop() []byte {
	select {
	case buf := <-self.bufs:
		return buf[:0]
	default:
		log.Println("queue is empty")
		return nil
	}
}

func (self *_ChanQueue) Push(buf []byte) {
	select {
	case self.bufs <- buf:
	default:
		log.Println("queue is full")
	}
}

//
// _AtomicStack
//
type _AtomicStackNode struct {
	Buf  []byte
	next unsafe.Pointer
}

type _AtomicStack struct {
	data []_AtomicStackNode

	head unsafe.Pointer

	lock uint32
}

func _NewAtomicStack(size int) *_AtomicStack {
	q := &_AtomicStack{
		data: make([]_AtomicStackNode, size),
	}
	for i := 0; i < size; i++ {
		q.data[i].Buf = make([]byte, 0, 512) // TODO передать размер параметром
		if i > 0 {
			q.data[i].next = unsafe.Pointer(&q.data[i-1])
		}
	}
	q.head = unsafe.Pointer(&q.data[size-1])

	return q
}

func (q *_AtomicStack) Pop() *_AtomicStackNode {
	type _AtomicStackNodePtr *_AtomicStackNode

	for {
		head := q.head
		if head != nil && head == q.head {
			next := unsafe.Pointer(_AtomicStackNodePtr(head).next)
			if atomic.CompareAndSwapPointer(&q.head, head, next) {
				return _AtomicStackNodePtr(head)
			}
		}
	}
}

func (q *_AtomicStack) Push(node *_AtomicStackNode) {
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

func TestAtomicQueue(t *testing.T) {
	const size = 4
	q := _NewAtomicStack(size)

	{
		n := q.Pop()
		n.Buf = append(n.Buf, "123"...)
		q.Push(n)
	}

	{
		n := q.Pop()

		if expected, val := "123", string(n.Buf); val != expected {
			t.Errorf("%s != %s", val, expected)
		}

		q.Push(n)
	}

	{
		ns := [4]*_AtomicStackNode{}
		for i := 0; i < size; i++ {
			ns[i] = q.Pop()
		}

		for i := 0; i < size; i++ {
			q.Push(ns[i])
		}
	}

	//
	//
	{
		// достали все
		ns := [4]*_AtomicStackNode{}
		for i := 0; i < size; i++ {
			ns[i] = q.Pop()
		}

		b := false

		done := make(chan bool)
		go func() {
			// ждем, пока появится свободная
			n := q.Pop()
			b = true
			close(done)
			q.Push(n)
		}()

		time.Sleep(time.Millisecond * 500)

		if b {
			t.Error("true != false")
		}

		// возвращаем
		for i := 0; i < size; i++ {
			q.Push(ns[i])
		}

		<-done
	}
}

func BenchmarkChanQueueS(b *testing.B) {
	q := _NewChanQueue(4)

	for i := 0; i < b.N; i++ {
		buf := q.Pop()
		q.Push(buf)
	}
}

func BenchmarkChanQueueM(b *testing.B) {
	q := _NewChanQueue(4)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			buf := q.Pop()
			q.Push(buf)
		}
	})
}

func BenchmarkAtomicStackS(b *testing.B) {
	q := _NewAtomicStack(4)

	for i := 0; i < b.N; i++ {
		n := q.Pop()
		q.Push(n)
	}
}

func BenchmarkAtomicStackM(b *testing.B) {
	q := _NewAtomicStack(4)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			n := q.Pop()
			q.Push(n)
		}
	})
}
