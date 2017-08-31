package http

import (
	"fmt"
	"log"
	"sort"
	"strconv"

	//"github.com/alex-shch/hlc-2017/prof"
	"github.com/alex-shch/hlc-2017/proto"
	"github.com/alex-shch/hlc-2017/store"
	"github.com/alex-shch/hlc-2017/utils"

	"github.com/valyala/fasthttp"
	"github.com/valyala/tcplisten"
)

var emptyJsonObj = "{}"
var emptyJsonAvg = `{"avg": 0.0}`
var emptyJsonVisits = `{"visits": []}`

var parseRequestError = fmt.Errorf("Invalid request")

type HttpHandler struct {
	Store *store.Store
	Addr  string
}

// TODO разбить на 16 корзин
var jsonBufs = utils.NewBytesAtomicStack(1024*1024, 16)
var avgJsonBufs = utils.NewBytesAtomicStack(128, 16)
var visitsBufs = utils.NewVisitsAtomicStack(1024, 16)

var requestCounter int

type NullLogger struct{}

func (NullLogger) Printf(format string, args ...interface{}) {}

func (self *HttpHandler) Run() error {
	cfg := tcplisten.Config{
		FastOpen:    true,
		DeferAccept: true,
		Backlog:     1024,
	}

	ln, err := cfg.NewListener("tcp4", self.Addr)
	if err != nil {
		log.Fatalf("error in reuseport listener: %s", err)
	}

	httpServer := &fasthttp.Server{
		Handler: self.HandleFastHTTP,
	}

	httpServer.DisableHeaderNamesNormalizing = true
	httpServer.Logger = NullLogger{}
	//httpServer.ReduceMemoryUsage = true

	if err = httpServer.Serve(ln); err != nil {
		log.Fatalf("error in fasthttp Server: %s", err)
	}

	//log.Fatal(fasthttp.ListenAndServe(self.Addr, self.HandleFastHTTP))

	return nil
}

func (self *HttpHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	entityCtx, err := parseRequest(ctx.Path())
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}
	self.handleFastHTTP(ctx, entityCtx)

	/*
		log.Println(requestCounter)

		if requestCounter < 5 {
			requestCounter++
		} else if requestCounter == 5 {
			requestCounter++
			prof.Prof.Stop()
		}
	*/
}

func (self *HttpHandler) handleFastHTTP(ctx *fasthttp.RequestCtx, entityCtx EntityCtx) {
	if ctx.Request.Header.IsGet() {
		//ctx.Response.Header.SetBytesKV([]byte("Connection"), []byte("keep-alive"))
		if entityCtx.SubTask {
			if entityCtx.Entity == ENTITY_USER {
				self.getUserVisits(ctx, entityCtx.Id)
			} else {
				self.getLocationAvg(ctx, entityCtx.Id)
			}
			return
		}

		switch entityCtx.Entity {
		case ENTITY_USER:
			self.getUser(ctx, entityCtx.Id)
		case ENTITY_VISIT:
			self.getVisit(ctx, entityCtx.Id)
		case ENTITY_LOCATION:
			self.getLocation(ctx, entityCtx.Id)
		}
	} else if ctx.Request.Header.IsPost() {
		/*
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("url: %s, body: %s\n", ctx.URI(), ctx.PostBody())
				}
			}()
		*/

		statusCode := fasthttp.StatusNotFound
		switch entityCtx.Entity {
		case ENTITY_USER:
			if entityCtx.Id == -1 {
				statusCode = self.addUser(ctx)
			} else {
				statusCode = self.updateUser(ctx, entityCtx.Id)
			}
		case ENTITY_VISIT:
			if entityCtx.Id == -1 {
				statusCode = self.addVisit(ctx)
			} else {
				statusCode = self.updateVisit(ctx, entityCtx.Id)
			}
		case ENTITY_LOCATION:
			if entityCtx.Id == -1 {
				statusCode = self.addLocation(ctx)
			} else {
				statusCode = self.updateLocation(ctx, entityCtx.Id)
			}
		}
		ctx.SetStatusCode(statusCode)
		//ctx.Response.SetBody(emptyJsonObj)
		//ctx.Response.BodyWriter().Write(emptyJsonObj)
		ctx.Response.AppendBodyString(emptyJsonObj)
		ctx.SetConnectionClose()
	} else {
		ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
	}
}

func (self *HttpHandler) getUser(ctx *fasthttp.RequestCtx, id int) {
	user, err := self.Store.GetUser(id)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	node := jsonBufs.Pop()
	buf := user.ToJson(node.Buf)
	ctx.Response.AppendBodyString(string(buf))
	jsonBufs.Push(node)

	/*
		if user.Json == nil {
			//node := visitsJsonBufs.Pop()
			//user.ToJson(&node.Buf)

			//ctx.Response.SetBody(buf)
			//ctx.Response.BodyWriter().Write(buf)
			//ctx.Response.AppendBodyString(string(node.Buf))
			ctx.Response.AppendBodyString(string(user.JsonBuf()))

			//visitsJsonBufs.Push(node)
		} else {
			//ctx.Response.BodyWriter().Write(user.Json)
			ctx.Response.AppendBodyString(string(user.Json))
		}
	*/
}

func (self *HttpHandler) getVisit(ctx *fasthttp.RequestCtx, id int) {
	visit, err := self.Store.GetVisit(id)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	node := jsonBufs.Pop()
	buf := visit.ToJson(node.Buf)
	ctx.Response.AppendBodyString(string(buf))
	jsonBufs.Push(node)

	/*
		if visit.Json == nil {
			//node := visitsJsonBufs.Pop()
			//visit.ToJson(&node.Buf)

			//ctx.Response.SetBody(buf)
			//ctx.Response.BodyWriter().Write(buf)
			//ctx.Response.AppendBodyString(string(node.Buf))
			ctx.Response.AppendBodyString(string(visit.JsonBuf()))

			//visitsJsonBufs.Push(node)
		} else {
			//ctx.Response.BodyWriter().Write(visit.Json)
			ctx.Response.AppendBodyString(string(visit.Json))
		}
	*/
}

func (self *HttpHandler) getLocation(ctx *fasthttp.RequestCtx, id int) {
	location, err := self.Store.GetLocation(id)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)

	node := jsonBufs.Pop()
	buf := location.ToJson(node.Buf)
	ctx.Response.AppendBodyString(string(buf))
	jsonBufs.Push(node)
}

func (self *HttpHandler) addUser(ctx *fasthttp.RequestCtx) int {
	user := &proto.User{}

	if err := user.UnmarshalJSON(ctx.Request.Body()); err != nil {
		return fasthttp.StatusBadRequest
	}

	// TODO проверить, что все поля передали

	err := self.Store.AddUser(user)
	if err != nil {
		return fasthttp.StatusBadRequest
	}

	return fasthttp.StatusOK
}

func (self *HttpHandler) addVisit(ctx *fasthttp.RequestCtx) int {
	visit := &proto.Visit{}

	if err := visit.UnmarshalJSON(ctx.Request.Body()); err != nil {
		return fasthttp.StatusBadRequest
	}

	// TODO проверить, что все поля передали

	err := self.Store.AddVisit(visit)
	if err != nil {
		return fasthttp.StatusBadRequest
	}

	return fasthttp.StatusOK
}

func (self *HttpHandler) addLocation(ctx *fasthttp.RequestCtx) int {
	location := &proto.Location{}

	if err := location.UnmarshalJSON(ctx.Request.Body()); err != nil {
		return fasthttp.StatusBadRequest
	}

	// TODO проверить, что все поля передали

	err := self.Store.AddLocation(location)
	if err != nil {
		return fasthttp.StatusBadRequest
	}

	return fasthttp.StatusOK
}

func (self *HttpHandler) updateUser(ctx *fasthttp.RequestCtx, id int) int {
	user := &proto.User{}

	if err := user.UnmarshalJSON(ctx.Request.Body()); err != nil {
		return fasthttp.StatusBadRequest
	}

	// возможно,проверка все же нужна
	if user.Gender.IsSet {
		if val := string(user.Gender.Val); val != "f" && val != "m" {
			return fasthttp.StatusBadRequest
		}
	}

	if user.Id.IsSet {
		//log.Printf("id: %d in update user request", user.Id.Val)
		return fasthttp.StatusBadRequest
	}

	err := self.Store.UpdateUser(id, user)
	if err != nil {
		return fasthttp.StatusNotFound
	}

	return fasthttp.StatusOK
}

func (self *HttpHandler) updateVisit(ctx *fasthttp.RequestCtx, id int) int {
	visit := &proto.Visit{}

	if err := visit.UnmarshalJSON(ctx.Request.Body()); err != nil {
		return fasthttp.StatusBadRequest
	}

	if visit.Id.IsSet {
		//log.Printf("id: %d in update visit request", visit.Id.Val)
		return fasthttp.StatusBadRequest
	}

	err := self.Store.UpdateVisit(id, visit)
	if err != nil {
		if err == store.NotFoundError {
			return fasthttp.StatusNotFound
		}
		return fasthttp.StatusBadRequest
	}

	return fasthttp.StatusOK
}

func (self *HttpHandler) updateLocation(ctx *fasthttp.RequestCtx, id int) int {
	location := &proto.Location{}

	if err := location.UnmarshalJSON(ctx.Request.Body()); err != nil {
		return fasthttp.StatusBadRequest
	}

	if location.Id.IsSet {
		//log.Printf("id: %d in update location request", location.Id.Val)
		return fasthttp.StatusBadRequest
	}

	err := self.Store.UpdateLocation(id, location)
	if err != nil {
		return fasthttp.StatusNotFound
	}

	return fasthttp.StatusOK
}

func (self *HttpHandler) getUserVisits(ctx *fasthttp.RequestCtx, id int) {
	user, err := self.Store.GetUser(id)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	f := _VisitFilter{}

	if err := f.Parse(ctx.QueryArgs()); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	srcVisits := user.Visits()

	//dstVisits := visitsBufs.Pop()
	visitsNode := visitsBufs.Pop()
	dstVisits := visitsNode.Array

	// TODO т.к. уже все отсортировано, можно сразу закидывать json в буфер

	for _, v := range srcVisits {
		//for _, node := range srcVisits {
		if v != nil {
			//if v := node.Visit; v != nil {
			if f.fromDate.isSet && f.fromDate.val >= v.VisitedAt {
				continue
			}
			if f.toDate.isSet && f.toDate.val <= v.VisitedAt {
				continue
			}

			location := v.Location()

			if f.distance.isSet && f.distance.val <= location.Distance {
				continue
			}

			if f.country.isSet && f.country.val != string(location.Country) {
				continue
			}

			dstVisits = append(dstVisits, v)
			//dstVisits = append(dstVisits, node)
		}
	}

	ctx.SetStatusCode(fasthttp.StatusOK)

	if size := len(dstVisits); size == 0 {
		ctx.Response.AppendBodyString(emptyJsonVisits)

	} else {
		if size > 1 {
			// надо сортировать

			if size == 2 {
				if dstVisits[0].VisitedAt > dstVisits[1].VisitedAt {
					dstVisits[0], dstVisits[1] = dstVisits[1], dstVisits[0]
				}
			} else {
				if size < 40 {
					dstVisits.Sort()
				} else {
					sort.Sort(dstVisits)
				}
			}
		}
		//ctx.Response.Header.SetContentLength(-2)

		node := jsonBufs.Pop()
		buf := node.Buf
		buf = append(buf, `{"visits": [`...)
		first := true
		for _, v := range dstVisits {
			//for _, node := range dstVisits {
			if first {
				first = false
			} else {
				buf = append(buf, ',')
			}
			buf = v.ToJsonList(buf)
			//node.Visit.ExportToJsonBuf(&buf)
		}
		buf = append(buf, "]}"...)
		//ctx.SetBody(buf)
		//ctx.Response.BodyWriter().Write(buf)
		ctx.Response.AppendBodyString(string(buf))
		jsonBufs.Push(node)
	}

	visitsBufs.Push(visitsNode)
}

func (self *HttpHandler) getLocationAvg(ctx *fasthttp.RequestCtx, id int) {
	location, err := self.Store.GetLocation(id)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	f := _LocationFilter{}

	if err := f.Parse(ctx.QueryArgs(), self.Store.DataTime()); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	count := 0
	marks := 0

	srcVisits := location.Visits()
	for _, v := range srcVisits {
		if v != nil {
			if f.fromDate.isSet && f.fromDate.val >= v.VisitedAt {
				continue
			}
			if f.toDate.isSet && f.toDate.val <= v.VisitedAt {
				continue
			}

			user := v.User()

			if f.gender.isSet && byte(f.gender.val) != byte(user.Gender) {
				continue
			}

			// чем больше unix-time, тем младше человек

			// считать тех, кто старше => пропускать тех,кто младше
			if f.fromAge.isSet && f.fromAge.val < user.BirthDate {
				continue
			}
			// считать тех, кто младше => пропускать тех,кто старше
			// если др сегодня, тоже пропускаем
			if f.toAge.isSet && user.BirthDate <= f.toAge.val {
				continue
			}

			count++
			marks += v.Mark
		}
	}

	ctx.SetStatusCode(fasthttp.StatusOK)

	if count > 0 {
		node := avgJsonBufs.Pop()
		node.Buf = append(node.Buf, `{"avg": `...)

		var avg float64
		if count > 1 {
			avg = float64(marks) / float64(count)
			avg = float64(int64(avg*100000+0.5)) / 100000
		} else {
			avg = float64(marks)
		}

		node.Buf = strconv.AppendFloat(node.Buf, avg, 'f', -1, 32)

		node.Buf = append(node.Buf, '}')
		//ctx.Response.Header.SetContentLength(len(buf))
		//ctx.Response.SetBody(buf)
		//ctx.Response.BodyWriter().Write(buf)
		ctx.Response.AppendBodyString(string(node.Buf))
		avgJsonBufs.Push(node)
	} else {
		//ctx.Response.Header.SetContentLength(len(empty))
		//ctx.Response.SetBodyString(empty)
		//ctx.Response.BodyWriter().Write(emptyJsonAvg)
		ctx.Response.AppendBodyString(emptyJsonAvg)
	}
}
