package http

import (
	"fmt"
	"strconv"
	"time"

	"github.com/valyala/fasthttp"
)

var InvalidGenderValueError = fmt.Errorf("Invalid gender value")

//
// IntFilter
//
type IntFilter struct {
	val   int
	isSet bool
}

func (self *IntFilter) Set(val string) error {
	if t, err := strconv.ParseInt(val, 10, 64); err != nil {
		return err
	} else {
		self.val = int(t)
		self.isSet = true
	}
	return nil
}

//
// Int64Filter
//
type Int64Filter struct {
	val   int64
	isSet bool
}

func (self *Int64Filter) Set(val string) error {
	if t, err := strconv.ParseInt(val, 10, 64); err != nil {
		return err
	} else {
		self.val = t
		self.isSet = true
	}
	return nil
}

//
// StringFilter
//
type StringFilter struct {
	val   string
	isSet bool
}

func (self *StringFilter) Set(val string) error {
	self.val = string(val)
	self.isSet = true
	return nil
}

//
// GenderFilter
//
type GenderFilter struct {
	val   byte
	isSet bool
}

func (self *GenderFilter) Set(val string) error {
	if l := len(val); l == 1 {
		if g := val[0]; g == 'f' || g == 'm' {
			self.val = g
			self.isSet = true
			return nil
		}
	}
	return InvalidGenderValueError
}

//
// AgeFilter
//
type AgeFilter struct {
	val   int
	isSet bool
}

func (self *AgeFilter) Set(val string, startTime time.Time) error {
	if t, err := strconv.ParseInt(val, 10, 64); err != nil {
		return err
	} else {
		birthDate := startTime.AddDate(-int(t), 0, 0)
		self.val = int(birthDate.Unix())
		self.isSet = true
	}
	return nil
}

//
// VisitFilter
//
type _VisitFilter struct {
	fromDate IntFilter
	toDate   IntFilter
	country  StringFilter
	distance IntFilter
}

func (self *_VisitFilter) Parse(args *fasthttp.Args) error {
	var err error
	args.VisitAll(func(key, val []byte) {
		if err == nil {
			strVal := string(val)
			switch string(key) {
			case "fromDate":
				err = self.fromDate.Set(strVal)
			case "toDate":
				err = self.toDate.Set(strVal)
			case "country":
				err = self.country.Set(strVal)
			case "toDistance":
				err = self.distance.Set(strVal)
			}
		}
	})

	return err
}

//
// LocationFilter
//
type _LocationFilter struct {
	fromDate IntFilter
	toDate   IntFilter
	fromAge  AgeFilter
	toAge    AgeFilter
	gender   GenderFilter
}

func (self *_LocationFilter) Parse(args *fasthttp.Args, startTime time.Time) error {
	var err error
	args.VisitAll(func(key, val []byte) {
		if err == nil {
			strVal := string(val)
			switch string(key) {
			case "fromDate":
				err = self.fromDate.Set(strVal)
			case "toDate":
				err = self.toDate.Set(strVal)
			case "fromAge":
				err = self.fromAge.Set(strVal, startTime)
			case "toAge":
				err = self.toAge.Set(strVal, startTime)
			case "gender":
				err = self.gender.Set(strVal)
			}
		}
	})

	return err
}
