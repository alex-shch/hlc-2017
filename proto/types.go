package proto

import (
	"fmt"
	"strconv"

	"github.com/alex-shch/hlc-2017/utils/utf8"
)

var nullValueError error = fmt.Errorf("Null value")

type Int struct {
	Val   int
	IsSet bool
}

func (self *Int) UnmarshalJSON(data []byte) error {
	if val, err := strconv.ParseInt(string(data), 10, 64); err != nil {
		return err
	} else {
		self.Val = int(val)
		self.IsSet = true
	}
	return nil
}

type Int64 struct {
	Val   int64
	IsSet bool
}

func (self *Int64) UnmarshalJSON(data []byte) error {
	if val, err := strconv.ParseInt(string(data), 10, 64); err != nil {
		return err
	} else {
		self.Val = val
		self.IsSet = true
	}
	return nil
}

type String struct {
	Val   string
	IsSet bool
}

func (self *String) UnmarshalJSON(data []byte) error {
	self.IsSet = true
	storage := [256]byte{}
	buf := storage[:0]
	err := utf8.Unquote(&buf, data)
	if err == nil {
		self.Val = string(buf)
	}
	return err
}

type String100 struct {
	data  [100]byte
	Val   []byte
	IsSet bool
}

func (self *String100) UnmarshalJSON(data []byte) error {
	self.IsSet = true
	self.Val = self.data[:0]
	return utf8.Unquote(&self.Val, data)
}
