package main

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/alex-shch/travels/loader"
	"github.com/alex-shch/travels/proto"
	"github.com/alex-shch/travels/store"
)

var (
	testStore      *store.Store
	usersCount     int
	visitsCount    int
	locationsCount int
)

func initBench() {
	if testStore != nil {
		return
	}

	startTime := time.Now()
	rand.Seed(startTime.Unix())

	testStore = store.New()
	if err := loader.Load(testStore, "testdata/full/data.zip"); err != nil {
		panic(err)
	}

	usersCount = len(testStore.Users.Pool)
	visitsCount = len(testStore.Visits.Pool)
	locationsCount = len(testStore.Locations.Pool)

	fmt.Printf("users count: %d, visits count: %d, locations count: %d\n",
		usersCount, visitsCount, locationsCount,
	)

	fmt.Printf("init ok, load time: %f\n", time.Now().Sub(startTime).Seconds())
}

func BenchmarkGetUser(b *testing.B) {
	initBench()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		id := rand.Intn(usersCount + int(usersCount/2))
		testStore.GetUser(id)
	}
}

func BenchmarkGetVisit(b *testing.B) {
	initBench()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		id := rand.Intn(visitsCount + int(visitsCount/2))
		testStore.GetVisit(id)
	}
}

func BenchmarkGetLocation(b *testing.B) {
	initBench()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		id := rand.Intn(locationsCount + int(locationsCount/2))
		testStore.GetLocation(id)
	}
}

func BenchmarkUpdateVisit(b *testing.B) {
	initBench()
	id := rand.Intn(visitsCount + int(visitsCount/2))

	data := proto.Visit{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if rand.Intn(5) == 0 {
			data.Location.IsSet = true
			data.Location.Val = rand.Intn(locationsCount + int(locationsCount/2))
		}
		if rand.Intn(5) == 0 {
			data.User.IsSet = true
			data.User.Val = rand.Intn(usersCount + int(usersCount/2))
		}

		testStore.UpdateVisit(id, &data)
	}
}

func BenchmarkLocationAvg(b *testing.B) {

}
