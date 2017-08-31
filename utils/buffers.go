package utils

/*
import (
	"log"

	"github.com/alex-shch/travels/store"
)

//
// byte buffer pool
//
type ByteBufferPool struct {
	bufs   chan []byte
	bufCap int
	name   string
}

func NewByteBufferPool(bufCapacity, poolSize, poolCapacity int, name string) *ByteBufferPool {
	pool := &ByteBufferPool{
		bufs:   make(chan []byte, poolCapacity),
		bufCap: bufCapacity,
	}
	for i := 0; i < poolSize; i++ {
		pool.bufs <- make([]byte, 0, bufCapacity)
	}
	return pool
}

func (self *ByteBufferPool) Pop() []byte {
	select {
	case buf := <-self.bufs:
		return buf[:0]
	default:
		log.Println("buf ", self.name, " is empty")
		return make([]byte, 0, self.bufCap)
	}
}

func (self *ByteBufferPool) Push(buf []byte) {
	select {
	case self.bufs <- buf:
	default:
		log.Println("buf ", self.name, " is full")
	}
}

//
// String buffer pool
//
type StringBufferPool struct {
	bufs   chan []string
	bufCap int
	name   string
}

func NewStringBufferPool(bufCapacity, poolSize, poolCapacity int, name string) *StringBufferPool {
	pool := &StringBufferPool{
		bufs:   make(chan []string, poolCapacity),
		bufCap: bufCapacity,
	}
	for i := 0; i < poolSize; i++ {
		pool.bufs <- make([]string, 0, bufCapacity)
	}
	return pool
}

func (self *StringBufferPool) Pop() []string {
	select {
	case buf := <-self.bufs:
		return buf[:0]
	default:
		log.Println("buf ", self.name, " is empty")
		return make([]string, 0, self.bufCap)
	}
}

func (self *StringBufferPool) Push(buf []string) {
	select {
	case self.bufs <- buf:
	default:
		log.Println("buf ", self.name, " is full")
	}
}

//
// VisitArray pool
//
type VisitArrayPool struct {
	bufs   chan store.VisitArray
	bufCap int
	name   string
}

func NewVisitArrayPool(bufCapacity, poolSize, poolCapacity int, name string) *VisitArrayPool {
	pool := &VisitArrayPool{
		bufs:   make(chan store.VisitArray, poolCapacity),
		bufCap: bufCapacity,
	}
	for i := 0; i < poolSize; i++ {
		pool.bufs <- make(store.VisitArray, 0, bufCapacity)
	}
	return pool
}

func (self *VisitArrayPool) Pop() store.VisitArray {
	select {
	case buf := <-self.bufs:
		return buf[:0]
	default:
		log.Println("buf ", self.name, " is empty")
		return make(store.VisitArray, 0, self.bufCap)
	}
}

func (self *VisitArrayPool) Push(buf store.VisitArray) {
	select {
	case self.bufs <- buf:
	default:
		log.Println("buf ", self.name, " is full")
	}
}
*/
