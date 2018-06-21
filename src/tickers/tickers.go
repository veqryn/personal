package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"time"
)

func main() {
	// Setup metrics in the default mux (pprof import also goes to the default mux) and start metrics+pprof server
	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			panic(err)
		}
	}()

	go func() {
		for range time.NewTicker(100 * time.Millisecond).C {
			fmt.Println("tick")
		}
	}()

	time.Sleep(100 * time.Second)

	buf := make([]byte, 1<<16)
	runtime.Stack(buf, true)
	fmt.Printf("%s", buf)
}
