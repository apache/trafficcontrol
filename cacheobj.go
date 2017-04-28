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
	reqTime          time.Time
	respTime         time.Time
	size             uint64
}

// ComputeSize computes the size of the given CacheObj. This computation is expensive, as the headers must be iterated over. Thus, the size should be computed once and stored, not computed on-the-fly for every new request for the cached object.
func (c CacheObj) ComputeSize() uint64 {
	// TODO include headers size
	return uint64(len(c.body))
}

func NewCacheObj(reqHeader http.Header, bytes []byte, code int, respHeader http.Header, reqTime time.Time, respTime time.Time) *CacheObj {
	obj := &CacheObj{
		body:             bytes,
		reqHeaders:       http.Header{},
		respHeaders:      http.Header{},
		respCacheControl: ParseCacheControl(respHeader),
		code:             code,
		reqTime:          reqTime,
		respTime:         respTime,
		size:             0,
	}
	// copyHeader(reqHeader, &obj.reqHeaders)
	// copyHeader(respHeader, &obj.respHeaders)
	obj.size = obj.ComputeSize()
	return obj
}
