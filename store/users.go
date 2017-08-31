package store

import (
	"strconv"
	"sync/atomic"

	"github.com/alex-shch/hlc-2017/proto"
)

type Gender byte

func (self *Gender) UnmarshalJSON(data []byte) error {
	*self = Gender(data[1]) // [", f | m, "]
	return nil
}

func (self Gender) MarshalJSON() ([]byte, error) {
	return []byte{'"', byte(self), '"'}, nil
}

//easyjson:json
type User struct {
	Id        int    `json:"id"`
	Gender    Gender `json:"gender"`
	BirthDate int    `json:"birth_date"`
	email     [100]byte
	Email     []byte `json:"email"`
	firsName  [50]byte
	FirstName []byte `json:"first_name"`
	lastName  [50]byte
	LastName  []byte `json:"last_name"`

	//Json   []byte
	//json   [512]byte
	visits_data [100]*Visit
	visits      VisitArray
	lock        uint32 // блокировка на добавление visit или генерацию json
}

func (self *User) Visits() VisitArray {
	return self.visits
}

func (self *User) addVisit(visit *Visit) {
	waitForSet(&self.lock)

	self.visits = append(self.visits, visit)

	//sort.Sort(self.visits)

	atomic.StoreUint32(&self.lock, 0)
}

func (self *User) delVisit(visit *Visit) {
	// TODO пока считаем что обновление не придет одновременно с чтением, или надо заменить на (chan RWMutex, 4)
	waitForSet(&self.lock)

	// TODO 1 - сортировать, 2 использовать поиск в отсортированном массиве

	// как вариант, просто писать nil не перемещая память
	// при создании списка nil-значения пропускать
	for i, v := range self.visits {
		//for _, node := range self.visits {
		if v == visit {
			//if node.Visit == visit {
			//node.Visit = nil
			self.visits[i] = nil

			// TODO посчитать, количество удалений и вставок
			break
		}
	}

	atomic.StoreUint32(&self.lock, 0)
}

/*
func (self *User) UpdateVisit(visit *Visit, visitedAt int64) {
		waitForSet(&self.lock)

		// TODO использовать двоичный поиск
		for _, node := range self.visits {
			if node.Visit == visit {
				node.VisitedAt = visitedAt
				break
			}
		}
		//sort.Sort(self.visits)

		atomic.StoreUint32(&self.lock, 0)
}
*/

func (self *User) ToJson(buf []byte) []byte {
	buf = append(buf, `{"id":`...)
	buf = strconv.AppendInt(buf, int64(self.Id), 10)
	buf = append(buf, `,"gender":"`...)
	buf = append(buf, byte(self.Gender))
	buf = append(buf, `","birth_date":`...)
	buf = strconv.AppendInt(buf, int64(self.BirthDate), 10)
	buf = append(buf, `,"email":"`...)
	buf = append(buf, self.Email...)
	buf = append(buf, `","first_name":"`...)
	buf = append(buf, self.FirstName...)
	buf = append(buf, `","last_name":"`...)
	buf = append(buf, self.LastName...)
	buf = append(buf, `"}`...)
	return buf
}

/*
func (self *User) JsonBuf() []byte {
	if self.Json == nil {
		waitForSet(&self.lock)
		if self.Json == nil {
			buf := self.json[:0]
			buf = append(buf, `{"id":`...)
			buf = strconv.AppendInt(buf, int64(self.Id), 10)
			buf = append(buf, `,"gender":"`...)
			buf = append(buf, byte(self.Gender))
			buf = append(buf, `","birth_date":`...)
			buf = strconv.AppendInt(buf, int64(self.BirthDate), 10)
			buf = append(buf, `,"email":"`...)
			buf = append(buf, self.Email...)
			buf = append(buf, `","first_name":"`...)
			buf = append(buf, self.FirstName...)
			buf = append(buf, `","last_name":"`...)
			buf = append(buf, self.LastName...)
			buf = append(buf, `"}`...)
			self.Json = buf
		}
		atomic.StoreUint32(&self.lock, 0)
	}
	return self.Json
}
*/

type UserPool []User

/*
func (self *UserPool) Alloc(user *User) *User {
	index := len(*self)
	// *self = append(*self, *user) // Проверить через atomic или использовать mutex
	*self = append(*self, User{}) // Проверить через atomic или использовать mutex
	ptr := &(*self)[index]
	ptr.visits = make(VisitArray, 0, 1024) // TODO на visits могла быть ленивая загрузка - НЕТ, но проверить!
	return ptr
}
*/

type UsersData struct {
	Pool UserPool
	Ids  map[int]*User

	lock uint32
}

type Users struct {
	//data [8]*UsersData
	data    []User
	dataLen int
	added   map[int]*User
}

func (self *Users) Get(id int) (*User, bool) {
	/*
		idx := id & 0x07 // 8 штук
		data := self.data[idx]
		//ptr, found := self.Ids[id]
		ptr, found := data.Ids[id]
		return ptr, found
	*/

	if id < self.dataLen {
		ptr := &self.data[id]
		return ptr, ptr.Id != 0
	} else {
		val, ok := self.added[id]
		return val, ok
	}
}

func (self *Users) Add(userData *proto.User) *User {
	id := userData.Id.Val

	var ptr *User
	if id < self.dataLen {
		ptr = &self.data[id]
	} else {
		ptr = &User{}
		self.added[id] = ptr
	}
	ptr.Id = id
	ptr.Email = append(ptr.email[:0], userData.Email.Val...)
	ptr.FirstName = append(ptr.firsName[:0], userData.FirstName.Val...)
	ptr.LastName = append(ptr.lastName[:0], userData.LastName.Val...)
	ptr.BirthDate = int(userData.BirthDate.Val)
	ptr.Gender = Gender(userData.Gender.Val[0])

	ptr.visits = ptr.visits_data[:0]

	return ptr

	/*
		idx := user.Id & 0x07 // 8 штук
		data := self.data[idx]
		waitForSet(&data.lock)
		ptr := data.Pool.Alloc(user)
		data.Ids[user.Id] = ptr
		atomic.StoreUint32(&data.lock, 0)
		return ptr
	*/
}
