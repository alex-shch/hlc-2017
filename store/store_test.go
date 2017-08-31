package store

import (
	"fmt"
	"testing"

	"github.com/alex-shch/travels/proto"
)

func TestStoreForUser(t *testing.T) {
	store := Store{
		Users: Users{
			Pool: make(UserPool, 0, 1),
			Ids:  make(map[int]*User, 1),
		},
	}

	uId := 7
	u1 := User{Id: uId, Gender: 'm', BirthDate: 30, FirstName: "Xxx"}

	u2 := proto.User{}
	u2.UnmarshalJSON([]byte(fmt.Sprintf(`{"id": %d, "gender": "f", "birth_date": 25}`, uId)))

	// TODO тестировать на генерацию ошибку

	store.AddUser(&u1)

	if u, err := store.GetUser(uId); err != nil {
		t.Error(err)
	} else {
		if u.BirthDate != u1.BirthDate {
			t.Errorf("%d != %d", u.BirthDate, u1.BirthDate)
		}
	}

	if err := store.UpdateUser(uId, &u2); err != nil {
		t.Error(err)
	}

	if u, err := store.GetUser(uId); err != nil {
		t.Error(err)
	} else {
		if u.BirthDate != u2.BirthDate.Val {
			t.Errorf("%d != %d", u.BirthDate, u1.BirthDate)
		}
	}

	if u, err := store.GetUser(uId); err != nil {
		t.Error(err)
	} else {
		if u.Id != uId { // должно было обновиться
			t.Errorf("%d != %d", u.Id, uId)
		}
		if u.FirstName != u1.FirstName { // должно было остаться старое значение
			t.Errorf("%d != %d", u.FirstName, u1.FirstName)
		}
	}
}

func TestStoreUpdateVisit(t *testing.T) {
	// TODO проверить, что обновляются ссылки на location и user
	t.Skip()
}
