package utils

import (
	"strings"
	"testing"
)

var path = "visits/34366/avg/"

/*
func split(dst *[]string, src string) {
	idx := 0
	for {
		pos := strings.IndexByte(src[idx:], '/')
		if pos == -1 {
			break
		} else {
			*dst = append(*dst, src[idx:idx+pos])
			idx += pos + 1
		}
	}
	if idx < len(src) {
		*dst = append(*dst, src[idx:])
	}
}
*/

func TestSplit(t *testing.T) {
	var dst []string
	Split(&dst, path)

	expected := []string{"visits", "34366", "avg"}
	if x, y := len(expected), len(dst); x != y {
		t.Fatalf("%d != %d", x, y)
	}

	for i, s := range expected {
		if val := dst[i]; val != s {
			t.Errorf("%q != %q", val, s)
		}
	}
}

func BenchmarkSplitNative(b *testing.B) {
	var resp []string

	for i := 0; i < b.N; i++ {
		resp = strings.Split(path, "/")
	}

	_ = resp // skip compilation error
}

func BenchmarkSplitSearch(b *testing.B) {
	src := append([]byte(nil), '/')
	src = append(src, path...)

	dst := make([]string, 0, 16)

	for i := 0; i < b.N; i++ {
		dst = dst[:0]
		//Split(&dst, path)
		Split(&dst, string(path[1:])) // именно так используется при разборе параметров
	}
}
