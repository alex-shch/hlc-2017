package store

import (
	"strconv"
	//"sync/atomic"

	"github.com/alex-shch/hlc-2017/proto"
)

const blockCount = 32 // должно быть степенью двойки

//easyjson:json
type Visit struct {
	Id         int `json:"id"`
	LocationId int `json:"location"`
	UserId     int `json:"user"`
	VisitedAt  int `json:"visited_at"`
	Mark       int `json:"mark"`

	//Json      []byte
	//json      [512]byte
	//Json2List []byte
	//json2List [128]byte
	location *Location
	user     *User
	lock     uint32
}

func (self *Visit) Location() *Location {
	return self.location
}

func (self *Visit) User() *User {
	return self.user
}

type VisitPool []Visit

/*
func (self *VisitPool) Alloc(visit *Visit) *Visit {
	index := len(*self)
	*self = append(*self, *visit) // Проверить через atomic или использовать mutex
	return &(*self)[index]
}
*/

/*
type VisitArrayNode struct {
	Visit     *Visit
	VisitedAt int64
}
type VisitArray []VisitArrayNode
*/
type VisitArray []*Visit

func (a VisitArray) Len() int           { return len(a) }
func (a VisitArray) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a VisitArray) Less(i, j int) bool { return a[i].VisitedAt < a[j].VisitedAt }

func (a VisitArray) Sort() {
	left := 0
	right := len(a) - 1

	for left <= right {
		for i := left; i < right; i++ {
			j := i + 1
			if a[i].VisitedAt > a[j].VisitedAt {
				a[i], a[j] = a[j], a[i]
			}
		}
		right -= 1

		for i := right; i > left; i-- {
			j := i - 1
			if a[i].VisitedAt < a[j].VisitedAt {
				a[i], a[j] = a[j], a[i]
			}
		}
		left++
	}
}

/*
func (a VisitArray) Less(i, j int) bool {
	ai := a[i]
	aj := a[j]
	if ai != nil && aj != nil {
		return ai.VisitedAt < aj.VisitedAt
	}
	if ai == nil && aj == nil {
		return false
	}

	// null last
	if ai == nil {
		return false
	}
	return true
}
*/

// для сущности Visit
func (self *Visit) ToJson(buf []byte) []byte {
	buf = append(buf, `{"id":`...)
	buf = strconv.AppendInt(buf, int64(self.Id), 10)
	buf = append(buf, `,"location":`...)
	buf = strconv.AppendInt(buf, int64(self.LocationId), 10)
	buf = append(buf, `,"user":`...)
	buf = strconv.AppendInt(buf, int64(self.UserId), 10)
	buf = append(buf, `,"visited_at":`...)
	buf = strconv.AppendInt(buf, int64(self.VisitedAt), 10)
	buf = append(buf, `,"mark":`...)
	buf = strconv.AppendInt(buf, int64(self.Mark), 10)
	buf = append(buf, '}')
	return buf
}

/*
// для сущности Visit
func (self *Visit) JsonBuf() []byte {
	if self.Json == nil {
		waitForSet(&self.lock)
		if self.Json == nil {
			buf := self.json[:0]
			buf = append(buf, `{"id":`...)
			buf = strconv.AppendInt(buf, int64(self.Id), 10)
			buf = append(buf, `,"location":`...)
			buf = strconv.AppendInt(buf, int64(self.LocationId), 10)
			buf = append(buf, `,"user":`...)
			buf = strconv.AppendInt(buf, int64(self.UserId), 10)
			buf = append(buf, `,"visited_at":`...)
			buf = strconv.AppendInt(buf, int64(self.VisitedAt), 10)
			buf = append(buf, `,"mark":`...)
			buf = strconv.AppendInt(buf, int64(self.Mark), 10)
			buf = append(buf, '}')
			self.Json = buf
		}
		atomic.StoreUint32(&self.lock, 0)
	}
	return self.Json
}
*/

// для списка Visit

func (self *Visit) ToJsonList(buf []byte) []byte {
	buf = append(buf, `{"mark":`...)
	buf = strconv.AppendInt(buf, int64(self.Mark), 10)
	buf = append(buf, `,"visited_at":`...)
	buf = strconv.AppendInt(buf, int64(self.VisitedAt), 10)
	buf = append(buf, `,"place":"`...)
	buf = append(buf, self.location.Place...)
	buf = append(buf, `"}`...)
	return buf
}

/*
func (self *Visit) ToJsonList(dstBuf []byte) []byte {
	if self.Json2List == nil {
		waitForSet(&self.lock)

		if self.Json2List == nil {
			buf := self.json2List[:0]
			buf = append(buf, `{"mark":`...)
			buf = strconv.AppendInt(buf, int64(self.Mark), 10)
			buf = append(buf, `,"visited_at":`...)
			buf = strconv.AppendInt(buf, int64(self.VisitedAt), 10)
			buf = append(buf, `,"place":"`...)
			buf = append(buf, self.location.Place...)
			buf = append(buf, `"}`...)
			self.Json2List = buf
		}

		atomic.StoreUint32(&self.lock, 0)
	}

	return append(dstBuf, self.Json2List...)
}
*/

type VisitsData struct {
	Pool VisitPool
	Ids  map[int]*Visit

	lock uint32
}

type Visits struct {
	//data [8]*VisitsData
	data    []Visit
	dataLen int
	added   map[int]*Visit
}

func (self *Visits) Get(id int) (*Visit, bool) {
	/*
		idx := id & 0x07 // 8 штук
		data := self.data[idx]
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

func (self *Visits) Add(visitData *proto.Visit) *Visit {
	/*
		idx := visit.Id & 0x07 // 8 штук
		data := self.data[idx]
		waitForSet(&data.lock)
		ptr := data.Pool.Alloc(visit)
		data.Ids[visit.Id] = ptr
		atomic.StoreUint32(&data.lock, 0)
	*/

	var ptr *Visit
	id := visitData.Id.Val
	if id < self.dataLen {
		ptr = &self.data[id]
	} else {
		ptr = &Visit{}
		self.added[id] = ptr
	}

	ptr.Id = id
	ptr.LocationId = visitData.Location.Val
	ptr.UserId = visitData.User.Val
	ptr.Mark = visitData.Mark.Val
	ptr.VisitedAt = int(visitData.VisitedAt.Val)

	return ptr

	/*
		visitedAt := visit.VisitedAt
		indexId := (visit.Id & (blockCount - 1)) // индекс от 0 до blockCount, должно быть степенью двойки
		indexes := self.Indexes[indexId].data

		waitForSet(&self.Indexes[indexId].lock)
		go func() {
			defer atomic.StoreUint32(&self.Indexes[indexId].lock, 0)
			lb := sort.Search(len(indexes), func(i int) bool { return indexes[i].VisitedAt >= visitedAt })
			if lb < len(indexes) {
				indexes = append(indexes, nil)
				copy(indexes[lb+1:], indexes[lb:])
				indexes[lb] = visit
				self.Indexes[indexId].data = indexes
			} else {
				self.Indexes[indexId].data = append(indexes, visit)
			}
		}()

		return ptr
	*/
}

/*
func (self *Visits) Count() int {
	count := 0
	for _, data := range self.data {
		count += len(data.Pool)
	}
	return count
}
*/
