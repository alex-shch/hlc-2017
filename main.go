package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	//"runtime/pprof"
	//"runtime/debug"

	"github.com/alex-shch/hlc-2017/http"
	"github.com/alex-shch/hlc-2017/loader"
	"github.com/alex-shch/hlc-2017/store"
)

// //go :generate go run vendor/github.com/mailru/easyjson/easyjson/main.go -all proto/proto.go
//go:generate go run vendor/github.com/mailru/easyjson/easyjson/main.go proto/
//go:generate go run vendor/github.com/mailru/easyjson/easyjson/main.go store/

// после генерации заменить в файле proto_easyjson.go для типов User, Visit и Location
//		if in.IsNull() {
//			in.Skip()
//			in.WantComma()
//			continue
//		}
// на
//		if in.IsNull() {
//			in.AddError(nullValueError)
//			return
//		}
//

/*
В общем с нашей стороны проверка выглядит как "взять текущую дату, поменять ей год и посмотреть кто оказался справа а кто слева"
И воткнут специальный костыль, чтобы вы не столкнулись с таймзонами и прочей веселухой. Отражу это в тз, многие спрашивают
*/

func main() {
	port := 8080
	data := "testdata"
	opt := "testdata/options.txt"

	if os.Getenv("DOCKER_RUN") == "yes" {
		port = 80
		data = "./data"
		opt = "/tmp/data/options.txt"
	}

	addr := fmt.Sprintf(":%d", port)

	log.Println("start")
	store := store.New()
	log.Println("init store completed")

	loader.LoadOptions(store, opt)
	log.Printf("ammo date: %v (%s)", store.DataTime(), store.DataTimeSrc())

	if err := loader.Load(store, data); err != nil {
		panic(err)
	}

	httpHandler := http.HttpHandler{
		Store: store,
		Addr:  addr,
	}

	//go http.Preheat(port)

	runtime.GC()
	log.Println("gc")

	//debug.SetGCPercent(-1)

	log.Fatal(httpHandler.Run())
}
