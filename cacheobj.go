package grove

import (
	"net/http"
	"time"
)

type CacheObj struct {
	body             []byte
	reqHeaders       http.Header
	respHeaders      http.Header
	respCacheControl CacheControl
	code             int
	reqTime          time.Time // this is our client's time when the object was requested
	reqRespTime      time.Time // this is our client's time when the object was received
	respRespTime     time.Time // this is the origin server's Date time when the object was sent
	size             uint64
}

// ComputeSize computes the size of the given CacheObj. This computation is expensive, as the headers must be iterated over. Thus, the size should be computed once and stored, not computed on-the-fly for every new request for the cached object.
func (c CacheObj) ComputeSize() uint64 {
	// TODO include headers size
	return uint64(len(c.body))
}

func NewCacheObj(reqHeader http.Header, bytes []byte, code int, respHeader http.Header, reqTime time.Time, reqRespTime time.Time, respRespTime time.Time) *CacheObj {
	obj := &CacheObj{
		body:             bytes,
		reqHeaders:       reqHeader,
		respHeaders:      respHeader,
		respCacheControl: ParseCacheControl(respHeader),
		code:             code,
		reqTime:          reqTime,
		reqRespTime:      reqRespTime,
		respRespTime:     respRespTime,
	}
	// copyHeader(reqHeader, &obj.reqHeaders)
	// copyHeader(respHeader, &obj.respHeaders)
	obj.size = obj.ComputeSize()
	return obj
}
