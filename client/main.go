package main

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"thread-pool/model"
)

var (
	defaultServer = ":8080"
)

const (
	requestsPerClient = 100000
	maxBatchSize      = (requestsPerClient / 10) * 2
)

var (
	s = rand.NewSource(time.Now().Unix())
	r = rand.New(s)
)

func submitRequests(url string) {
	var (
		req   *model.ClientRequest
		reqID uint
	)

	msgLeft := requestsPerClient
	for 0 < msgLeft {
		batch := r.Intn(maxBatchSize)
		if batch > msgLeft {
			batch = msgLeft
		}
		msgLeft -= batch

		for i := 0; i < batch; i++ {
			req = new(model.ClientRequest)
			reqID++
			req.ID = reqID
			req.Size = r.Intn(model.ReqDataSize)

			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(req)
			if err != nil {
				panic(err)
			}

			resp, err := http.Post(url, "text/json", &buf)
			if err != nil {
				panic(err)
			}
			resp.Body.Close()
		}
	}
}

func main() {
	url := "http://localhost"
	submitRequests(url + defaultServer)
}
