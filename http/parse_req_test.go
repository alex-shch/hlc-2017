package http

import (
	"testing"
)

func TestParseRequest(t *testing.T) {
	type Answer struct {
		ctx EntityCtx
		err error
	}

	tests := map[string]Answer{
		"/users/1":         Answer{EntityCtx{ENTITY_USER, 1, false}, nil},
		"/users/new":       Answer{EntityCtx{ENTITY_USER, -1, false}, nil},
		"/users/1/visits":  Answer{EntityCtx{ENTITY_USER, 1, true}, nil},
		"/locations/1":     Answer{EntityCtx{ENTITY_LOCATION, 1, false}, nil},
		"/locations/new":   Answer{EntityCtx{ENTITY_LOCATION, -1, false}, nil},
		"/locations/1/avg": Answer{EntityCtx{ENTITY_LOCATION, 1, true}, nil},
		"/visits/1":        Answer{EntityCtx{ENTITY_VISIT, 1, false}, nil},
		"/visits/new":      Answer{EntityCtx{ENTITY_VISIT, -1, false}, nil},
	}

	for req, expected := range tests {
		a, err := parseRequest([]byte(req))

		if a != expected.ctx || err != expected.err {
			t.Errorf("%q -> %v != %v, %s", req, a, expected.ctx, err)
		}
	}
}
