package store

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"time"

	"github.com/alex-shch/hlc-2017/proto"
)

var NotFoundError = fmt.Errorf("Not found")
var AlreadyExistsError = fmt.Errorf("Already exists")

type Store struct {
	Users     Users
	Visits    Visits
	Locations Locations

	dataTime       time.Time
	dataTimeSource string
}

func New() *Store {
	t := time.Now()

	store := &Store{
		//Users: Users{},
		//Locations: Locations{},
		//Visits:         Visits{},
		dataTime:       time.Date(t.Year(), t.Month(), t.Day()-1, 0, 0, 0, 0, time.UTC),
		dataTimeSource: "undefined",
	}

	store.Locations.added = make(map[int]*Location, 1000)
	store.Locations.dataLen = 800000
	store.Locations.data = make([]Location, store.Locations.dataLen)
	/*
		for i := 0; i < len(store.Locations.data); i++ {
			store.Locations.data[i] = &LocationsData{
				Pool: make(LocationPool, 0, 20000),
				Ids:  make(map[int]*Location, 20000),
			}
		}
	*/

	store.Users.added = make(map[int]*User, 1000)
	store.Users.dataLen = 100200
	store.Users.data = make([]User, store.Users.dataLen)
	/*
		for i := 0; i < len(store.Users.data); i++ {
			store.Users.data[i] = &UsersData{
				Pool: make(UserPool, 0, 30000),
				Ids:  make(map[int]*User, 30000),
			}
		}
	*/

	store.Visits.added = make(map[int]*Visit, 1000)
	store.Visits.dataLen = 10001500
	store.Visits.data = make([]Visit, store.Visits.dataLen)

	/*
		for i := 0; i < len(store.Visits.data); i++ {
			store.Visits.data[i] = &VisitsData{
				Pool: make(VisitPool, 0, 200000),
				Ids:  make(map[int]*Visit, 200000),
			}
		}
	*/

	/*
		size := cap(store.Visits.Pool)/blockCount + 1
		for i, _ := range store.Visits.Indexes {
			store.Visits.Indexes[i] = VisitIndex{
				data: make([]*Visit, 0, size),
			}
		}
	*/

	return store
}

func (self *Store) GetUser(id int) (*User, error) {
	user, found := self.Users.Get(id)

	if !found {
		return nil, NotFoundError
	}

	return user, nil
}

func (self *Store) AddUser(userData *proto.User) error {
	if _, found := self.Users.Get(userData.Id.Val); found {
		return AlreadyExistsError
	} else {
		self.Users.Add(userData)
	}

	return nil
}

func (self *Store) UpdateUser(id int, userData *proto.User) error {
	user, found := self.Users.Get(id)
	if !found {
		return NotFoundError
	}

	if userData.Gender.IsSet {
		user.Gender = Gender(userData.Gender.Val[0])
	}
	if userData.BirthDate.IsSet {
		user.BirthDate = int(userData.BirthDate.Val)
	}
	if userData.Email.IsSet {
		user.Email = append(user.email[:0], userData.Email.Val...)
	}
	if userData.FirstName.IsSet {
		user.FirstName = append(user.firsName[:0], userData.FirstName.Val...)
	}
	if userData.LastName.IsSet {
		user.LastName = append(user.lastName[:0], userData.LastName.Val...)
	}

	//user.Json = nil
	// TODO go user.UpdateJsonSrc()

	return nil
}

func (self *Store) GetLocation(id int) (*Location, error) {
	location, found := self.Locations.Get(id)
	if !found {
		return nil, NotFoundError
	}

	return location, nil
}

func (self *Store) AddLocation(locationData *proto.Location) error {
	if _, found := self.Locations.Get(locationData.Id.Val); found {
		return AlreadyExistsError
	} else {
		self.Locations.Add(locationData)
	}
	return nil
}

func (self *Store) UpdateLocation(id int, locationData *proto.Location) error {
	location, found := self.Locations.Get(id)
	if !found {
		return NotFoundError
	}

	if locationData.Place.IsSet {
		//location.Place = locationData.Place.Val
		location.Place = append(location.place[:0], locationData.Place.Val...)
		/*
			for _, v := range location.visits {
				if v != nil {
					v.Json2List = nil
				}
			}
		*/
	}
	if locationData.Country.IsSet {
		location.Country = append(location.country[:0], locationData.Country.Val...)
	}
	if locationData.City.IsSet {
		location.City = append(location.city[:0], locationData.City.Val...)
	}
	if locationData.Distance.IsSet {
		location.Distance = int(locationData.Distance.Val)
	}

	//location.Json = nil
	// TODO go user.UpdateJsonSrc()s

	return nil
}

func (self *Store) GetVisit(id int) (*Visit, error) {
	visit, found := self.Visits.Get(id)
	if !found {
		return nil, NotFoundError
	}

	return visit, nil
}

func (self *Store) AddVisit(visitData *proto.Visit) error {

	if _, found := self.Visits.Get(visitData.Id.Val); found {
		return AlreadyExistsError
	}

	// TODO вот тут возможно гонка и в соседней рутине пытается добавиться Location
	location, found := self.Locations.Get(visitData.Location.Val)
	if !found {
		return NotFoundError
	}

	// TODO вот тут возможно гонка и в соседней рутине пытается добавиться User
	user, found := self.Users.Get(visitData.User.Val)
	if !found {
		return NotFoundError
	}

	v := self.Visits.Add(visitData)

	v.location = location
	v.user = user
	user.addVisit(v)
	location.addVisit(v)
	return nil
}

func (self *Store) UpdateVisit(id int, visitData *proto.Visit) error {
	visit, found := self.Visits.Get(id)
	if !found {
		return NotFoundError
	}

	//needUpdateUser := false
	//var newVisitedAt int64

	// если обновляем user и location, уже нужно новое значение
	// P.S. решение себя не оправдало
	if visitData.VisitedAt.IsSet {
		visit.VisitedAt = int(visitData.VisitedAt.Val)
		//needUpdateUser = true
	}

	// если передали locationId
	locId := int(visitData.Location.Val)
	if visitData.Location.IsSet && visit.LocationId != locId {
		// TODO возможна вставка location в параллельном запросе
		location, found := self.Locations.Get(locId)
		if !found {
			return NotFoundError
		}
		visit.LocationId = locId
		// переключение рутины оказалось слишком дорогим, может откладывать запуск до момента, пока не завершился запрос пользователя?
		//go location.addVisit(visit)
		//go visit.location.delVisit(visit)
		location.addVisit(visit)
		visit.location.delVisit(visit)
		visit.location = location
	}

	// если передали userId
	userId := int(visitData.User.Val)
	if visitData.User.IsSet && visit.UserId != userId {
		// TODO возможна вставка user в параллельном запросе
		user, found := self.Users.Get(userId)
		if !found {
			return NotFoundError
		}

		visit.UserId = userId
		// переключение рутины оказалось слишком дорогим
		//go user.addVisit(visit)
		//go visit.user.delVisit(visit)
		user.addVisit(visit)
		visit.user.delVisit(visit)
		visit.user = user
		//needUpdateUser = false // только что обновили
	}

	if visitData.Mark.IsSet {
		visit.Mark = int(visitData.Mark.Val)
	}

	//visit.Json = nil
	//visit.Json2List = nil

	//if needUpdateUser && visit.user != nil {
	//	visit.user.UpdateVisit(visit, newVisitedAt)
	//}

	return nil
}

func (self *Store) GetUserVisits(userId int) (VisitArray, error) {
	user, found := self.Users.Get(userId)
	if !found {
		return nil, NotFoundError
	}

	return user.Visits(), nil
}

func (self *Store) UpdateLinks() {
	userVisitsCount := 0
	locationVisitsCount := 0

	tmpBuf := make([]byte, 512)

	for i, _ := range self.Visits.data {
		v := &self.Visits.data[i]
		if v.Id != 0 {
			if user, found := self.Users.Get(v.UserId); !found {
				panic(fmt.Sprintf("user %d not found", v.UserId))
			} else {
				v.user = user
				self.Visits.data[i] = *v
				user.addVisit(v)
				if count := len(user.visits); userVisitsCount < count {
					userVisitsCount = count
				}
			}

			if location, found := self.Locations.Get(v.LocationId); !found {
				panic(fmt.Sprintf("location %d not found", v.LocationId))
			} else {
				v.location = location
				self.Visits.data[i] = *v
				location.addVisit(v)
				if count := len(location.visits); locationVisitsCount < count {
					locationVisitsCount = count
				}
			}

			tmpBuf = tmpBuf[:0]
			v.ToJsonList(tmpBuf)
		}
	}
	log.Printf("user visits: %d, location visits: %d\n", userVisitsCount, locationVisitsCount)

	// возможно частичное упорядочивание здесь позволит быстрее отсортировать на выходе
	for i, _ := range self.Users.data {
		u := &self.Users.data[i]
		if u.Id != 0 {
			sort.Sort(u.visits)
		}
	}
}

func (self *Store) SetDataTime(src []byte) {
	self.dataTimeSource = string(src)
	timeVal, err := strconv.ParseInt(self.dataTimeSource, 10, 64)
	if err != nil {
		log.Println("err parse %s: %v", src, err)
	}
	t := time.Unix(timeVal, 0)
	self.dataTime = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

func (self *Store) DataTime() time.Time {
	return self.dataTime
}

func (self *Store) DataTimeSrc() string {
	return self.dataTimeSource
}
