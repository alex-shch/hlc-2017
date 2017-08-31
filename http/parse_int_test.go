package http

import (
	"math/rand"
	"strconv"
	"testing"

	"github.com/valyala/fasthttp"
)

var dict = []string{
	"5673452", "12934878", "235241851",
	"3964242", "85424113", "354623247",
	"8532426", "55774745", "875645476",
	"7456533", "57865760", "054757656",
}
var dictLen = len(dict)

func BenchmarkStdParseInt32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		num := dict[rand.Intn(dictLen)]

		_, err := strconv.ParseInt(num, 10, 32)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkStdParseInt64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		num := dict[rand.Intn(dictLen)]

		_, err := strconv.ParseInt(num, 10, 64)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkFasthttpParseInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		num := dict[rand.Intn(dictLen)]

		_, err := fasthttp.ParseUint([]byte(num))
		if err != nil {
			b.Fatal(err)
		}
	}
}
