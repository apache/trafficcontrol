package cacheobj

/*
   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

import (
	"net/http"
	"time"

	"github.com/apache/trafficcontrol/grove/web"
)

type CacheObj struct {
	Body             []byte
	ReqHeaders       http.Header
	RespHeaders      http.Header
	RespCacheControl web.CacheControl
	Code             int
	OriginCode       int
	ProxyURL         string
	ReqTime          time.Time // our client's time when the object was requested
	ReqRespTime      time.Time // our client's time when the object was received
	RespRespTime     time.Time // the origin server's Date time when the object was sent
	LastModified     time.Time // the origin LastModified if it exists, or Date if it doesn't
	Size             uint64
	HitCount         uint64 // the number of times this object was hit
}

// ComputeSize computes the size of the given CacheObj. This computation is expensive, as the headers must be iterated over. Thus, the size should be computed once and stored, not computed on-the-fly for every new request for the cached object.
func (c CacheObj) ComputeSize() uint64 {
	// TODO include headers size
	return uint64(len(c.Body))
}

func New(reqHeader http.Header, bytes []byte, code int, originCode int, proxyURL string, respHeader http.Header, reqTime time.Time, reqRespTime time.Time, respRespTime time.Time, lastModified time.Time) *CacheObj {
	obj := &CacheObj{
		Body:             bytes,
		ReqHeaders:       reqHeader,
		RespHeaders:      respHeader,
		RespCacheControl: web.ParseCacheControl(respHeader),
		Code:             code,
		OriginCode:       originCode,
		ProxyURL:         proxyURL,
		ReqTime:          reqTime,
		ReqRespTime:      reqRespTime,
		RespRespTime:     respRespTime,
		LastModified:     lastModified,
		HitCount:         1,
	}
	// copyHeader(reqHeader, &obj.reqHeaders)
	// copyHeader(respHeader, &obj.respHeaders)
	obj.Size = obj.ComputeSize()
	return obj
}
