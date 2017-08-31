package http

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
)

func Preheat(port int) {
	time.Sleep(time.Second)
	preheat(port)

	runtime.GC()
	log.Println("gc")
}

func preheat(port int) {
	log.Println("preheating")

	dst := make([]byte, 0, 128)

	size := 512 * 1024
	sendBuf := make([]byte, size)
	for i := 0; i < size; i++ {
		sendBuf[i] = ' '
	}
	sendBuf[0] = '{'
	sendBuf[size-1] = '}'

	args := fasthttp.Args{}
	args.AppendBytes(sendBuf)

	entities := []string{"users", "visits", "locations"}

	wg := sync.WaitGroup{}

	shot := func(id int) {
		for _, entity := range entities {
			url := fmt.Sprintf("http://0.0.0.0:%d/%s/%d", port, entity, id)

			//statusCode, body, err := fasthttp.GetTimeout(dst, url, time.Second)
			_, _, err := fasthttp.GetTimeout(dst, url, time.Second)
			if err != nil {
				log.Printf("%s %s", url, err)
			} else {
				//log.Printf("%s %d, body len: %d", url, statusCode, len(body))
			}

			//statusCode, body, err = fasthttp.Post(dst, url, &args)
			_, _, err = fasthttp.Post(dst, url, &args)
			if err != nil {
				log.Printf("%s %s", url, err)
			} else {
				//log.Printf("%s %d, body len: %d", url, statusCode, len(body))
			}
		}
		wg.Done()
	}

	count := 10000
	wg.Add(count)
	for i := 0; i < count; i++ {
		id := rand.Intn(5000000)
		go shot(id)
	}
	wg.Wait()

	log.Println("preheating completed")
}
