package model

const (
	ReqAdd = iota
	ReqAvg
	ReqRandom
	ReqSpellCheck
	ReqSearch
)

const ReqDataSize = 1024

type ClientRequest struct {
	ID          uint
	RequestType int
	Data        [ReqDataSize]byte
	Size        int
}
