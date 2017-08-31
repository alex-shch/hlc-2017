package test

import (
	//"bytes"
	"fmt"
	"strconv"
	"testing"
)

// implements io.Writer
type bufWriter struct {
	buf []byte
}

func (self *bufWriter) Write(p []byte) (n int, err error) {
	self.buf = append(self.buf, p...)
	return len(p), nil
}

var (
	mark      = 123414
	visitedAt = int64(34523423234)
	place     = "reghfbsvcasfa"

	size = 256
	buf  = make([]byte, 0, size)
	w    = &bufWriter{make([]byte, 0, size)}
)

func printfJson() {
	fmt.Fprintf(
		w,
		`{"mark":%d, "visited_at":%d, "place":%q}`,
		mark,
		visitedAt,
		place,
	)
}

func appendJson() {
	buf = append(buf, `{"mark":`...)
	buf = strconv.AppendInt(buf, int64(mark), 10)
	buf = append(buf, `, "visited_at":`...)
	buf = strconv.AppendInt(buf, visitedAt, 10)
	buf = append(buf, `, "place":"`...)
	buf = append(buf, place...)
	buf = append(buf, `"}`...)
}

func TestPrintfAppendEqual(t *testing.T) {
	printfJson()
	s1 := string(w.buf)
	t.Log(s1)
	appendJson()
	s2 := string(buf)
	t.Log(s2)

	if s1 != s2 {
		t.Errorf("%q != %q", s1, s2)
	}
}

func BenchmarkFprintf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		printfJson()
		w.buf = w.buf[:0]
	}
}

func BenchmarkBytesBuf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		appendJson()
		buf = buf[:0]
	}
}
