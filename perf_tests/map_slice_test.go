package test

import (
	"math/rand"
	"sort"
	"testing"
)

//
// Hash map
//
type Map map[int]*int

func (self Map) add(id int) bool {
	_, found := self[id]
	if found {
		return false
	} else {
		self[id] = nil
		return true
	}
}

func (self Map) get(id int) (*int, bool) {
	val, found := self[id]
	return val, found
}

//
// Lower bound
//
type Node struct {
	id  int
	val *int
}
type Slice []Node

func (self Slice) add(id int) bool {
	size := len(self)
	lb := sort.Search(size, func(i int) bool { return self[i].id >= id })
	if lb < size {
		//fmt.Printf("lower bound for %d at index: %d, value: %d\n", x, lb, a[lb])
		if self[lb].id == id {
			return false
		}

		self = append(self, Node{})
		copy(self[lb+1:], self[lb:])
		ptr := &self[lb]
		ptr.id = id
		ptr.val = nil
		//fmt.Println(a)
	} else {
		//fmt.Printf("lower bound for value %d doesn't exist\n", x)
		self = append(self, Node{id: id, val: nil})
	}
	return true
}

func (self Slice) get(id int) (*int, bool) {
	size := len(self)
	lb := sort.Search(size, func(i int) bool { return self[i].id >= id })
	if lb < size {
		ptr := &self[lb]
		if ptr.id == id {
			return ptr.val, true
		}
	}
	return nil, false
}

//
// Testing
//
func BenchmarkAddToMap(b *testing.B) {
	count := b.N + b.N/2
	store := make(Map, count)
	added := 0
	exists := 0

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		id := rand.Intn(count)
		b.StartTimer()

		if ok := store.add(id); ok {
			added++
		} else {
			exists++
		}
	}

	b.Logf("added: %d, exists: %d", added, exists)
}

func BenchmarkAddToSlice(b *testing.B) {
	count := b.N + b.N/2
	store := make(Slice, count)
	added := 0
	exists := 0

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		id := rand.Intn(count)
		b.StartTimer()

		if ok := store.add(id); ok {
			added++
		} else {
			exists++
		}
	}

	b.Logf("added: %d, exists: %d", added, exists)
}

func BenchmarkGetFromMap(b *testing.B) {
	count := b.N + b.N/2
	store := make(Map, count)
	fails := 0
	exists := 0

	for i := 0; i < b.N; i++ {
		id := rand.Intn(count)
		store[id] = nil
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		id := rand.Intn(count)
		b.StartTimer()

		if _, found := store.get(id); found {
			exists++
		} else {
			fails++
		}
	}

	b.Logf("exists: %d, fails: %d", exists, fails)
}

func BenchmarkGetFromSlice(b *testing.B) {
	count := b.N + b.N/2
	store := make(Slice, count)
	fails := 0
	exists := 0

	for i := 0; i < b.N; i++ {
		id := rand.Intn(count)
		store.add(id)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		id := rand.Intn(count)
		b.StartTimer()

		if _, found := store.get(id); found {
			exists++
		} else {
			fails++
		}
	}

	b.Logf("exists: %d, fails: %d", exists, fails)
}
