package main

import (
	"encoding/json"
	"net/http"
	"sync"
	"sync/atomic"
	"thread-pool/model"
	"thread-pool/rp"
	"time"

	"github.com/sirupsen/logrus"
)

type (
	Server interface {
		Start() error
		Stop()
	}

	TCPServer struct {
		numReqs uint64
		port    string
		log     *logrus.Logger
		s       *http.Server
	}
)

var pool *sync.Pool

func newTCPServer(port string) Server {
	if *strategy == "sync_pool" {
		pool = rp.GetSyncPool()
	}

	log := logrus.New()
	log.Formatter = &logrus.JSONFormatter{
		PrettyPrint: true,
	}
	log.Level = logrus.InfoLevel

	srv := &TCPServer{
		port: port,
		log:  log,
		s: &http.Server{
			Addr:         port,
			ReadTimeout:  1 * time.Second,
			WriteTimeout: 500 * time.Millisecond,
		},
	}

	srv.s.Handler = srv
	return srv
}

func (srv *TCPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	srv.log.Debug("received a request")
	atomic.AddUint64(&srv.numReqs, 1)
	defer r.Body.Close()

	switch *strategy {
	case "local":
		req := rp.GetNewClientRequest()
		rp.IncrementTotal()
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			panic(err)
		}
		srv.log.WithFields(logrus.Fields{
			"size":       req.Size,
			"type":       req.RequestType,
			"request_id": req.ID,
		}).Debug("request body decoded")
	case "new":
		req := rp.GetNewClientRequestPtr()
		rp.IncrementTotal()
		rp.IncrementAllocated()
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			panic(err)
		}
		srv.log.WithFields(logrus.Fields{
			"size":       req.Size,
			"type":       req.RequestType,
			"request_id": req.ID,
		}).Debug("request body decoded")
	case "resource_pool":
		req := rp.Alloc()
		defer rp.Release(req)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			panic(err)
		}
		srv.log.WithFields(logrus.Fields{
			"size":       req.Size,
			"type":       req.RequestType,
			"request_id": req.ID,
		}).Debug("request body decoded")
	default:
		req := pool.Get().(*model.ClientRequest)
		defer pool.Put(req)
		rp.IncrementTotal()
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			panic(err)
		}
		srv.log.WithFields(logrus.Fields{
			"size":       req.Size,
			"type":       req.RequestType,
			"request_id": req.ID,
		}).Debug("request body decoded")
	}
}

func (srv *TCPServer) Start() error {
	srv.log.WithFields(logrus.Fields{
		"addr":     srv.s.Addr,
		"strategy": *strategy,
	}).Infoln("starting server...")

	return srv.s.ListenAndServe()
}

func (srv *TCPServer) Stop() {
	srv.log.WithFields(logrus.Fields{
		"addr": srv.s.Addr,
	}).Infoln("stopping server...")
	srv.s.Close()
	srv.log.WithFields(logrus.Fields{
		"request_count": srv.numReqs,
	}).Infoln("server stopped")

	rp.PrintStats(*strategy)
}
