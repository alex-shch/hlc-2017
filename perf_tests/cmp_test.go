package test

import (
	"bytes"
	"math/rand"
	"testing"
)

var stringDict = []string{
	"users", "visits", "locations",
	"users", "visits", "locations",
	"users", "visits", "locations",
	"users", "visits", "locations",
	"users", "visits", "locations",
	"use23456", "vis", "loca",
}
var dictLen = len(stringDict)

var byteDict = func() [][]byte {
	a := make([][]byte, 0, dictLen)
	for i := 0; i < dictLen; i++ {
		a = append(a, []byte(stringDict[i]))
	}
	return a
}()

func BenchmarkStringCmp(b *testing.B) {
	_true := 0
	_false := 0

	for i := 0; i < b.N; i++ {
		s1 := stringDict[rand.Intn(dictLen)]
		s2 := stringDict[rand.Intn(dictLen)]

		if s1 == s2 {
			_true++
		} else {
			_false++
		}
	}

	b.Logf("true: %d, false: %d", _true, _false)
}

func BenchmarkBytesCmp(b *testing.B) {
	_true := 0
	_false := 0

	for i := 0; i < b.N; i++ {
		b1 := byteDict[rand.Intn(dictLen)]
		b2 := byteDict[rand.Intn(dictLen)]

		if bytes.Equal(b1, b2) {
			_true++
		} else {
			_false++
		}
	}

	b.Logf("true: %d, false: %d", _true, _false)
}

func BenchmarkBytesStringCmp(b *testing.B) {
	_true := 0
	_false := 0

	for i := 0; i < b.N; i++ {
		b1 := byteDict[rand.Intn(dictLen)]
		b2 := byteDict[rand.Intn(dictLen)]

		if string(b1) == string(b2) {
			_true++
		} else {
			_false++
		}
	}

	b.Logf("true: %d, false: %d", _true, _false)
}
