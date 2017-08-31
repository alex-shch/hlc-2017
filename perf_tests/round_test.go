package test

import (
	"math/rand"
	"strconv"
	"testing"
)

var roundBuf = make([]byte, 0, 64)

func round32_5(val float64) []byte {
	buf := strconv.AppendFloat(roundBuf[:0], val, 'f', 5, 32)
	return buf
}

func round64_5(val float64) []byte {
	buf := strconv.AppendFloat(roundBuf[:0], val, 'f', 5, 64)
	return buf
}

func round32_1(val float64) []byte {
	buf := strconv.AppendFloat(roundBuf[:0], val, 'f', -1, 32)
	return buf
}

func round64_1(val float64) []byte {
	buf := strconv.AppendFloat(roundBuf[:0], val, 'f', -1, 64)
	return buf
}

func TestCompareRoundX32X64(t *testing.T) {
	b32 := make([]byte, 0, 64)
	b64 := make([]byte, 0, 64)

	for i := 0; i < 1000; i++ {
		val := rand.Float64()
		val = float64(int64(val*100000+0.5)) / 100000
		b32 = strconv.AppendFloat(b32[:0], val, 'f', -1, 64)
		b64 = strconv.AppendFloat(b64[:0], val, 'f', -1, 64)
		s32 := string(b32)
		s64 := string(b64)
		if s32 != s64 {
			t.Errorf("%s != %s", s32, s64)
		}
	}
}

func TestRound32(t *testing.T) {
	if expected, val := "1.00000", round32_5(1.0000011); string(val) != expected {
		t.Errorf("%s != %s", val, expected)
	}

	if expected, val := "1.00001", round32_5(1.0000053); string(val) != expected {
		t.Errorf("%s != %s", val, expected)
	}
}

func BenchmarkRound(b *testing.B) {

}

func BenchmarkRound32x5(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		val := rand.Float64()
		b.StartTimer()

		round32_5(val)
	}
}

func BenchmarkRound64x5(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		val := rand.Float64()
		b.StartTimer()

		round64_5(val)
	}
}

func BenchmarkRound32x1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		val := rand.Float64()
		val = float64(int64(val*100000+0.5)) / 100000
		b.StartTimer()

		round32_1(val)
	}
}

func BenchmarkRound64x1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		val := rand.Float64()
		val = float64(int64(val*100000+0.5)) / 100000
		b.StartTimer()

		round64_1(val)
	}
}
