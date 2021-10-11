package rp

import (
	"fmt"
	"sync"
	"sync/atomic"
	"thread-pool/model"
)

const CR_POOL_SIZE = 10

var pool chan *model.ClientRequest
var total, allocated uint64

func init() {
	pool = make(chan *model.ClientRequest, CR_POOL_SIZE)
}

func Alloc() *model.ClientRequest {
	IncrementTotal()
	select {
	case cr := <-pool:
		return cr
	default:
		IncrementAllocated()
		return new(model.ClientRequest)
	}
}

func Release(cr *model.ClientRequest) {
	select {
	case pool <- cr:
	default:
	}
}

func GetNewClientRequest() model.ClientRequest {
	atomic.AddUint64(&total, 1)
	return model.ClientRequest{}
}

func GetNewClientRequestPtr() *model.ClientRequest {
	IncrementTotal()
	IncrementAllocated()
	return new(model.ClientRequest)
}

func GetSyncPool() *sync.Pool {
	return &sync.Pool{
		New: func() interface{} {
			IncrementAllocated()
			return new(model.ClientRequest)
		},
	}
}

func IncrementTotal() {
	atomic.AddUint64(&total, 1)
}

func IncrementAllocated() {
	atomic.AddUint64(&allocated, 1)
}

func PrintStats(strategy string) {
	var reused uint64
	if strategy == "local" {
		reused = 0
	} else {
		reused = total - allocated
	}
	fmt.Printf(`Stats
Total: %d
Reused: %d
Allocated: %d
`, total, reused, allocated)
}
