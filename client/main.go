package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"thread-pool/model"
)

const (
	requestsPerClient = 100000
	maxBatchSize      = (requestsPerClient / 10) * 2
)

var (
	defaultServer = "http://localhost:8080"
	r             = rand.New(rand.NewSource(time.Now().Unix()))
)

var (
	numClients = flag.Int("num_clients", 1, "numbers of clients to send 100,000 requests")
)

func submitRequests(url string) {
	var (
		req   *model.ClientRequest
		reqID uint
		wg    sync.WaitGroup
	)

	msgLeft := requestsPerClient
	for 0 < msgLeft {
		batch := r.Intn(maxBatchSize)
		if batch > msgLeft {
			batch = msgLeft
		}
		msgLeft -= batch

		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < batch; i++ {
				req = new(model.ClientRequest)
				reqID++
				req.ID = reqID
				req.Size = r.Intn(model.ReqDataSize)

				buf := new(bytes.Buffer)
				err := json.NewEncoder(buf).Encode(req)
				if err != nil {
					panic(err)
				}

				resp, err := http.Post(url, "text/json", buf)
				if err != nil {
					panic(err)
				}
				resp.Body.Close()
			}
		}()
	}
	wg.Wait()
}

func main() {
	var wg sync.WaitGroup
	flag.Parse()

	wg.Add(*numClients)
	for i := 0; i < *numClients; i++ {
		go func() {
			defer wg.Done()
			submitRequests(defaultServer)
		}()
	}
	wg.Wait()
}
