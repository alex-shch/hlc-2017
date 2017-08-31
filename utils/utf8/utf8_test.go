package utf8

import (
	"testing"
)

func TestUnquote(t *testing.T) {
	src := `"\u0412\u0438\u043a\u0442\u043e\u0440"`
	dst := make([]byte, 0, 3*len(src)/2)
	if err := Unquote(&dst, src); err != nil {
		t.Error(err)
	}
	if expected, val := "Виктор", string(dst); val != expected {
		t.Errorf("%q != %q", val, expected)
	}
}
