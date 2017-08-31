package utils

import (
	"strings"
)

func Split(dst *[]string, src string) {
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
