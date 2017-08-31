package http

import (
	"bytes"
	//"fmt"
	"strconv"
)

const (
	ENTITY_USER     = 1
	ENTITY_VISIT    = 2
	ENTITY_LOCATION = 3
)

type EntityCtx struct {
	Entity  int
	Id      int
	SubTask bool
}

func parseRequest(uri []byte) (ctx EntityCtx, err error) {
	ctx.Entity = -1
	ctx.Id = -1
	ctx.SubTask = false

	pos := bytes.IndexByte(uri, '?')
	var src []byte
	if pos == -1 {
		src = uri
	} else {
		src = uri[:pos]
	}

	// parse entity
	idx := 1
	pos = bytes.IndexByte(src[idx:], '/')
	if pos == -1 {
		err = parseRequestError
		return
	} else {
		val := string(src[1 : idx+pos])

		//fmt.Println(val)

		if val == "users" {
			ctx.Entity = ENTITY_USER
		} else if val == "visits" {
			ctx.Entity = ENTITY_VISIT
		} else if val == "locations" {
			ctx.Entity = ENTITY_LOCATION
		} else {
			err = parseRequestError
			return
		}
		idx += pos + 1
	}

	// parse id
	{
		pos = bytes.IndexByte(src[idx:], '/')
		var idVal int64
		var val string
		if pos == -1 {
			val = string(src[idx:])
		} else {
			val = string(src[idx : idx+pos])
		}

		if val != "new" {
			idVal, err = strconv.ParseInt(val, 10, 32)
			if err != nil {
				return
			}
			ctx.Id = int(idVal)
		}

		if pos == -1 {
			return
		}

		idx += pos + 1
	}

	// parse additional command
	{
		pos = bytes.IndexByte(src[idx:], '/')
		var val string
		if pos == -1 {
			val = string(src[idx:])
		} else {
			val = string(src[idx : idx+pos])
		}

		if ctx.Entity == ENTITY_LOCATION && ctx.Id != -1 && val == "avg" {
			ctx.SubTask = true
			return
		} else if ctx.Entity == ENTITY_USER && ctx.Id != -1 && val == "visits" {
			ctx.SubTask = true
			return
		}
	}

	err = parseRequestError
	return
}
