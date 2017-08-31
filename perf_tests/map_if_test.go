package test

import (
	"math/rand"
	"testing"
)

//go:noinline
func checkIf(entity string) int {
	if entity != "users" && entity != "visits" && entity != "locations" {
		return 404
	}
	return 200
}

var _map = map[string]int{
	"users":     200,
	"visits":    200,
	"locations": 200,
}

//go:noinline
func checkMap(entity string) int {
	_, found := _map[entity]
	if found {
		return 200
	}
	return 404
}

var dict = []string{
	"users", "visits", "locations",
	"users", "visits", "locations",
	"users", "visits", "locations",
	"users", "visits", "locations",
	"users", "visits", "locations",
	"use23456", "vis", "loca",
}
var dictLen = len(dict)

func BenchmarkCheckIf(b *testing.B) {
	_200 := 0
	_404 := 0

	for i := 0; i < b.N; i++ {
		s := dict[rand.Intn(dictLen)]

		code := checkIf(s)

		if code == 200 {
			_200++
		} else {
			_404++
		}
	}

	b.Logf("200: %d, 404: %d", _200, _404)
}

func BenchmarkCheckMap(b *testing.B) {
	_200 := 0
	_404 := 0

	for i := 0; i < b.N; i++ {
		s := dict[rand.Intn(dictLen)]

		code := checkMap(s)

		if code == 200 {
			_200++
		} else {
			_404++
		}
	}

	b.Logf("200: %d, 404: %d", _200, _404)
}
