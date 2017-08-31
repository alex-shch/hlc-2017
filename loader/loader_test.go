package loader

import (
	"strings"
	"testing"

	"github.com/alex-shch/travels/store"
)

func TestSplitter(t *testing.T) {
	vals := []string{
		`{"Name": "Ed", "Text": "Knock knock."}`,
		`{"Name": "Sam", "Text": "Who's there?"}`,
		`{"Name": "Ed", "Text": "Go fmt."}`,
		`{"Name": "Sam", "Text": "Go fmt who?"}`,
		`{"Name": "Ed", "Text": "Go fmt yourself!"}`,
	}

	jsonStream := `
		{
			"users": [` + strings.Join(vals, ",\r\n") + `], "xxx": 123
		}
	`

	ch, err := splitter(strings.NewReader(jsonStream), "users")
	if err != nil {
		t.Fatal(err)
	}

	i := 0
	for buf := range ch {
		if expected := vals[i]; string(buf) != expected {
			t.Errorf("Ln %d, %q != %q", i, buf, expected)
		}
		i++
	}
	if expected := len(vals); i != expected {
		t.Errorf("%d != %d", i, expected)
	}
}

func TestLoader(t *testing.T) {
	storage := store.New()

	if err := Load(storage, "../testdata/data.zip"); err != nil {
		t.Error(err)
	}

	if u, err := storage.GetUser(24); err != nil {
		t.Error(err)
	} else {
		if expected := int64(76291200); u.BirthDate != expected {
			t.Errorf("%d != %d", u.BirthDate, expected)
		}
		if expected := store.Gender('f'); u.Gender != expected {
			t.Errorf("%v != %v", u.Gender, expected)
		}
	}

	if l, err := storage.GetLocation(36); err != nil {
		t.Error(err)
	} else if expected := "Шри-Ланка"; l.Country != expected {
		t.Errorf("%q != %q", l.Country, expected)
	}

	if v, err := storage.GetVisit(17); err != nil {
		t.Error(err)
	} else {
		if expected := 105; v.LocationId != expected {
			t.Errorf("%d != %d", v.LocationId, expected)
		}
		if expected := 19; v.UserId != expected {
			t.Errorf("%d != %d", v.UserId, expected)
		}
	}
}
