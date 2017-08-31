package loader

import (
	//"archive/zip"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/alex-shch/hlc-2017/proto"
	"github.com/alex-shch/hlc-2017/store"
)

var (
	usersCount     = 0
	visitsCount    = 0
	locationsCount = 0

	maxPlaceLen = 0
)

func LoadOptions(store *store.Store, optFilePath string) {
	file, err := os.Open(optFilePath)
	if err != nil {
		log.Printf("error open %s: %s\n", optFilePath, err)
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	buf, _, err := reader.ReadLine()

	if err != nil {
		log.Printf("error read %s err: %s\n", optFilePath, err)
	} else {
		store.SetDataTime(buf)
	}

}

func Load(store *store.Store, dataPath string) error {
	files, err := ioutil.ReadDir(dataPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		fileName := file.Name()
		pos := strings.IndexByte(fileName, '_')
		if pos == -1 {
			fmt.Errorf("Invalid file name %q", fileName)
			continue
		}

		err = loadFile(store, dataPath+"/"+fileName, fileName[0:pos])
		if err != nil {
			return err
		}

	}

	store.UpdateLinks()

	log.Printf("phase1, users: %d, visits: %d, locations: %d, max place len: %d",
		usersCount,
		visitsCount,
		locationsCount,
		maxPlaceLen,
	)

	return nil
}

func loadFile(store *store.Store, fileName, entity string) error {
	rc, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer rc.Close()

	// TODO сначала грузить пользователей и места, потом посещения
	// второй вариант, сначала прогрузить все, потомо настроить связи на пользователей и места
	// вроде второй вариант выглядит предпочтительнее

	// TODO к сохраняемым объектам сохранять и исходный json, пересчитывать в отдельной рутине, если объект обновился

	switch entity {
	case "users":
		if ch, err := splitter(rc, "users"); err != nil {
			return err
		} else if err := loadUsers(&store.Users, ch); err != nil {
			return err
		}
	case "locations":
		if ch, err := splitter(rc, "locations"); err != nil {
			return err
		} else if err := loadLocations(&store.Locations, ch); err != nil {
			return err
		}
	case "visits":
		if ch, err := splitter(rc, "visits"); err != nil {
			return err
		} else if err := loadVisits(&store.Visits, ch); err != nil {
			return err
		}
	default:
		return fmt.Errorf("Invalid file name %q", fileName)
	}

	return nil
}

func splitter(reader io.Reader, checkName string) (<-chan []byte, error) {
	dec := json.NewDecoder(reader)

	// '{'
	if t, err := dec.Token(); err != nil {
		return nil, err
	} else if fmt.Sprintf("%v", t) != "{" {
		return nil, fmt.Errorf(`invalid file data, expected "{", found "%v"`, t)
	}

	// "users" | "visits" | "locations"
	if t, err := dec.Token(); err != nil {
		return nil, err
	} else if t != checkName {
		return nil, fmt.Errorf("invalid file data, expected %q, found %q", checkName, t)
	}

	// '['
	if t, err := dec.Token(); err != nil {
		return nil, err
	} else if fmt.Sprintf("%v", t) != "[" {
		return nil, fmt.Errorf(`invalid file data, expected "[", found "%v"`, t)
	}

	ch := make(chan []byte)
	go func() {
		defer close(ch)
		for dec.More() { // прочитает до закрывающегося тега массива
			m := make(json.RawMessage, 0)
			if err := dec.Decode(&m); err != nil {
				panic(err)
			} else {
				ch <- m
			}
		}
	}()

	return ch, nil
}

var user1006found bool

func loadUsers(users *store.Users, stream <-chan []byte) error {
	user := &proto.User{}

	for buf := range stream {
		if err := user.UnmarshalJSON(buf); err != nil {
			return fmt.Errorf("Failed parse %q: %s", buf, err)
		}

		//user.Json = buf
		users.Add(user)

		usersCount++
	}

	return nil
}

func loadLocations(locations *store.Locations, stream <-chan []byte) error {
	location := &proto.Location{}

	for buf := range stream {
		if err := location.UnmarshalJSON(buf); err != nil {
			return fmt.Errorf("Failed parse %q: %s", buf, err)
		}

		if l := len(location.Place.Val); l > maxPlaceLen {
			maxPlaceLen = l
		}

		//location.Json = buf
		locations.Add(location)

		locationsCount++
	}
	return nil
}

func loadVisits(visits *store.Visits, stream <-chan []byte) error {
	visit := &proto.Visit{}

	for buf := range stream {
		if err := visit.UnmarshalJSON(buf); err != nil {
			return fmt.Errorf("Failed parse %q: %s", buf, err)
		}

		//visit.Json = buf
		visits.Add(visit)

		visitsCount++
	}
	return nil
}
