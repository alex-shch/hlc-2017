package store

import (
	"testing"
)

func TestUserAdd(t *testing.T) {
	users := Users{
		Pool: make(UserPool, 0, 1),
		Ids:  make(map[int]*User, 1),
	}

	u1 := User{Id: 7, Gender: 'm', BirthDate: 30}
	u2 := User{Id: 8, Gender: 'f', BirthDate: 25}

	users.Add(&u1)

	if u, found := users.Get(u1.Id); !found {
		t.Errorf("user %d not found", u1.Id)
	} else {
		if u.Id != u1.Id {
			t.Errorf("%d != %d", u.Id, u1.Id)
		}
	}

	users.Add(&u2)

	if val := len(users.Pool); val != 2 {
		t.Error("%d != 2", val)
	}

	if val := len(users.Ids); val != 2 {
		t.Error("%d != 2", val)
	}

	if u, found := users.Get(u1.Id); !found {
		t.Errorf("user %d not found", u1.Id)
	} else {
		if u.Id != u1.Id {
			t.Errorf("%d != %d", u.Id, u1.Id)
		}
	}

	if u, found := users.Get(u2.Id); !found {
		t.Errorf("user %d not found", u2.Id)
	} else {
		if u.Id != u2.Id {
			t.Errorf("%d != %d", u.Id, u2.Id)
		}
	}
}

func TestUserAddVisit(t *testing.T) {
	u := User{}

	v1 := &Visit{}
	v2 := &Visit{}
	v3 := &Visit{}

	f := func(v *Visit, done chan<- bool) {
		for i := 0; i < 50; i++ {
			u.addVisit(v)
		}
		done <- true
	}

	done := make(chan bool)
	go f(v1, done)
	go f(v2, done)
	go f(v3, done)

	<-done
	<-done
	<-done

	if val := len(u.visits); val != 150 {
		t.Errorf("%d != %d", val, 150)
	}
}
