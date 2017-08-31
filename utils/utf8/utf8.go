package utf8

import (
	"fmt"
	"unicode/utf8"
)

var ErrSyntax = fmt.Errorf("Utf-8 syntax error")

func Unquote(buf *[]byte, s []byte) error {
	n := len(s)
	if n < 2 {
		return ErrSyntax
	}
	s = s[1 : n-1]

	var runeTmp [utf8.UTFMax]byte
	//buf := make([]byte, 0, 3*len(s)/2) // Try to avoid more allocations.
	k := 6
	for i, size := 0, n-2; i < size; i += k {
		if string(s[i:i+2]) == "\\u" {
			n := unquoteChar(runeTmp[:], s[i+2:])
			*buf = append(*buf, runeTmp[:n]...)
			k = 6
		} else {
			*buf = append(*buf, s[i])
			k = 1
		}
	}
	return nil
}

func unquoteChar(buf []byte, s []byte) int {
	var v rune
	for j := 0; j < 4; j++ {
		x := unhex(s[j])
		v = v<<4 | x
	}
	return utf8.EncodeRune(buf, v)
}

func unhex(b byte) rune {
	c := rune(b) // ? зачем?
	switch {
	case '0' <= c && c <= '9':
		return c - '0'
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10
	case 'A' <= c && c <= 'F':
		return c - 'A' + 10
	}
	return '?'
}
