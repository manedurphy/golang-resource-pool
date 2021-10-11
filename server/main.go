package main

import (
	"flag"
	"net/http"
	"time"
)

var (
	strategy = flag.String("strategy", "local", "strategy for memory allocation of request body")
)

func main() {
	flag.Parse()
	srv := newTCPServer(":8080")

	go func() {
		if err := srv.Start(); err != nil {
			if err != http.ErrServerClosed {
				panic(err)
			}
		}
	}()

	time.Sleep(50 * time.Second)
	srv.Stop()
}
