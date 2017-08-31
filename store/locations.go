package store

import (
	"strconv"
	"sync/atomic"

	"github.com/alex-shch/hlc-2017/proto"
)

//easyjson:json
type Location struct {
	Id int `json:"id"`
	//Place    string `json:"place"`
	place    [21]byte
	Place    []byte `json:"place"`
	country  [50]byte
	Country  []byte `json:"country"`
	city     [50]byte
	City     []byte `json:"city"`
	Distance int    `json:"distance"`

	//Json   []byte
	//json   [512]byte
	visits_data [100]*Visit
	visits      []*Visit
	lock        uint32 // блокировка на добавление visit или генерацию json
}

func (self *Location) Visits() []*Visit {
	return self.visits
}

func (self *Location) addVisit(visit *Visit) {
	waitForSet(&self.lock)
	self.visits = append(self.visits, visit)

	atomic.StoreUint32(&self.lock, 0)
}

func (self *Location) delVisit(visit *Visit) {
	waitForSet(&self.lock)

	// как вариант, просто писать nil не перемещая память
	// при создании списка nil-значения пропускать
	for i, v := range self.visits {
		if v == visit {
			self.visits[i] = nil
			// TODO посчитать, количество удалений и вставок
			break
		}
	}

	atomic.StoreUint32(&self.lock, 0)
}

func (self *Location) ToJson(buf []byte) []byte {
	buf = append(buf, `{"id":`...)
	buf = strconv.AppendInt(buf, int64(self.Id), 10)
	buf = append(buf, `,"place":"`...)
	buf = append(buf, self.Place...)
	buf = append(buf, `","country":"`...)
	buf = append(buf, self.Country...)
	buf = append(buf, `","city":"`...)
	buf = append(buf, self.City...)
	buf = append(buf, `","distance":`...)
	buf = strconv.AppendInt(buf, int64(self.Distance), 10)
	buf = append(buf, '}')
	return buf
}

/*
func (self *Location) JsonBuf() []byte {
	if self.Json == nil {
		waitForSet(&self.lock)
		if self.Json == nil {
			buf := self.json[:0]
			buf = append(buf, `{"id":`...)
			buf = strconv.AppendInt(buf, int64(self.Id), 10)
			buf = append(buf, `,"place":"`...)
			buf = append(buf, self.Place...)
			buf = append(buf, `","country":"`...)
			buf = append(buf, self.Country...)
			buf = append(buf, `","city":"`...)
			buf = append(buf, self.City...)
			buf = append(buf, `","distance":`...)
			buf = strconv.AppendInt(buf, int64(self.Distance), 10)
			buf = append(buf, '}')
			self.Json = buf
		}
		atomic.StoreUint32(&self.lock, 0)
	}
	return self.Json
}
*/

type LocationPool []Location

/*
func (self *LocationPool) Alloc(location *Location) *Location {
	index := len(*self)
	*self = append(*self, *location) // Проверить через atomic или использовать mutex
	ptr := &(*self)[index]
	ptr.visits = make([]*Visit, 0, 2048) // TODO выделять буфер из пула
	return ptr
}
*/

type LocationsData struct {
	Pool LocationPool
	Ids  map[int]*Location

	lock uint32
}

type Locations struct {
	//data [8]*LocationsData
	data    []Location
	dataLen int
	added   map[int]*Location
}

func (self *Locations) Get(id int) (*Location, bool) {
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

func (self *Locations) Add(location *proto.Location) *Location {
	// TODO по ТЗ все приходящие id будут уникальны
	//if _, found := self.Ids[location.Id]; found {
	//	return nil, fmt.Errorf("Location %d already exists", location.Id)
	//}

	/*
		idx := location.Id & 0x07 // 8 штук
		data := self.data[idx]
		waitForSet(&data.lock)
		ptr := data.Pool.Alloc(location)
		data.Ids[location.Id] = ptr
		atomic.StoreUint32(&data.lock, 0)
		return ptr
	*/

	id := location.Id.Val
	var ptr *Location
	if id < self.dataLen {
		ptr = &self.data[id]
	} else {
		ptr = &Location{}
		self.added[id] = ptr
	}

	ptr.Id = id
	ptr.City = append(ptr.city[:0], location.City.Val...)
	ptr.Country = append(ptr.country[:0], location.Country.Val...)
	ptr.Place = append(ptr.place[:0], location.Place.Val...)
	ptr.Distance = location.Distance.Val

	ptr.visits = ptr.visits_data[:0]

	return ptr
}
