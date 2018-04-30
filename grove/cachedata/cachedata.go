package cachedata

// cachedata exists as a package to avoid import cycles

import (
	"net/http"
	"time"

	"github.com/apache/incubator-trafficcontrol/grove/remapdata"
	"github.com/apache/incubator-trafficcontrol/grove/web"
)

// ParentResponseData contains data about the parent/origin response.
type ParentRespData struct {
	Reuse            remapdata.Reuse
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
