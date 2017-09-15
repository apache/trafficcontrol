package cache

import (
	"net/http"
	"time"

	"github.com/apache/incubator-trafficcontrol/grove/web"
)

type CacheObj struct {
	Body             []byte
	ReqHeaders       http.Header
	RespHeaders      http.Header
	RespCacheControl web.CacheControl
	Code             int
	ReqTime          time.Time // this is our client's time when the object was requested
	ReqRespTime      time.Time // this is our client's time when the object was received
	RespRespTime     time.Time // this is the origin server's Date time when the object was sent
	Size             uint64
}

// ComputeSize computes the size of the given CacheObj. This computation is expensive, as the headers must be iterated over. Thus, the size should be computed once and stored, not computed on-the-fly for every new request for the cached object.
func (c CacheObj) ComputeSize() uint64 {
	// TODO include headers size
	return uint64(len(c.Body))
}

func New(reqHeader http.Header, bytes []byte, code int, respHeader http.Header, reqTime time.Time, reqRespTime time.Time, respRespTime time.Time) *CacheObj {
	obj := &CacheObj{
		Body:             bytes,
		ReqHeaders:       reqHeader,
		RespHeaders:      respHeader,
		RespCacheControl: web.ParseCacheControl(respHeader),
		Code:             code,
		ReqTime:          reqTime,
		ReqRespTime:      reqRespTime,
		RespRespTime:     respRespTime,
	}
	// copyHeader(reqHeader, &obj.reqHeaders)
	// copyHeader(respHeader, &obj.respHeaders)
	obj.Size = obj.ComputeSize()
	return obj
}
