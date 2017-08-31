package proto

import (
	"testing"
)

func TestFullData(t *testing.T) {
	js := []byte(`{
        "id": 1,
        "email": "robosen@icloud.com",
        "first_name": "Данила",
        "last_name": "Стамленский",
        "gender": "m",
        "birth_date": 345081600
    }`)

	u := &User{}
	if err := u.UnmarshalJSON(js); err != nil {
		t.Error(err)
	}

	if expected, val := 1, u.Id.Val; val != expected {
		t.Errorf("%v != %v", val, expected)
	}

	if expected, val := "robosen@icloud.com", u.Email.Val; val != expected {
		t.Errorf("%v != %v", val, expected)
	}

	if expected, val := int64(345081600), u.BirthDate.Val; val != expected {
		t.Errorf("%d != %d", val, expected)
	}
}

func TestPartialData(t *testing.T) {
	js := []byte(`{
        "first_name": "\u0412\u0438\u043a\u0442\u043e\u0440",
        "last_name": ""
    }`)

	u := &User{}
	if err := u.UnmarshalJSON(js); err != nil {
		t.Error(err)
	}

	if expected, val := false, u.Email.IsSet; val != expected {
		t.Errorf("%v != %v", val, expected)
	}

	if expected, val := false, u.Email.IsSet; val != expected {
		t.Errorf("%v != %v", val, expected)
	}

	if expected, val := true, u.LastName.IsSet; val != expected {
		t.Errorf("%v != %v", val, expected)
	} else {
		if expected, val := "", u.LastName.Val; val != expected {
			t.Errorf("%q != %q", val, expected)
		}
	}

	if expected, val := true, u.FirstName.IsSet; val != expected {
		t.Errorf("%v != %v", val, expected)
	} else {
		if expected, val := "Виктор", u.FirstName.Val; val != expected {
			t.Errorf("%q != %q", val, expected)
		}
	}
}

func TestNullData(t *testing.T) {
	u := &User{}

	if err := u.UnmarshalJSON([]byte(`{"last_name": null}`)); err == nil {
		t.Error("need null error")
	}

	if err := u.UnmarshalJSON([]byte(`{"gender": null}`)); err == nil {
		t.Error("need null error")
	}

	if err := u.UnmarshalJSON([]byte(`{"birth_date": null}`)); err == nil {
		t.Error("need null error")
	}
}
