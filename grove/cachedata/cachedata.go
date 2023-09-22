package cachedata

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

// cachedata exists as a package to avoid import cycles

import (
	"net/http"
	"time"

	"github.com/apache/trafficcontrol/v8/grove/web"
	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
)

// ParentResponseData contains data about the parent/origin response.
type ParentRespData struct {
	Reuse            rfc.Reuse
	OriginCode       int
	OriginReqSuccess bool
	// OriginConnectFailed is whether the connection to the origin succeeded. It's possible to get a failure response from an origin, but have the connection succeed.
	OriginConnectFailed bool
	OriginBytes         uint64
	ProxyStr            string
}

// HandlerData contains data generally held by the Handler, and known as soon as the request is received.
type SrvrData struct {
	Hostname string
	Port     string
	Scheme   string
}

type ReqData struct {
	Req      *http.Request
	Conn     *web.InterceptConn
	ClientIP string
	ReqTime  time.Time
	ToFQDN   string
}

type RespData struct {
	RespCode     int
	BytesWritten uint64
	RespSuccess  bool
	CacheHit     bool
}
