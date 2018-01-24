package cache

import (
	"net/http"
	"time"

	"github.com/apache/incubator-trafficcontrol/grove/remap"
	"github.com/apache/incubator-trafficcontrol/grove/web"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
)

// Responder is an object encapsulating the cache's response to the client. It holds all the data necessary to respond, log the response, and add the stats.
type Responder struct {
	W                         http.ResponseWriter
	R                         *http.Request
	StatLog                   *StatLogger
	F                         RespondFunc
	Reuse                     remap.Reuse
	OriginCode                int
	ResponseCode              int
	SuccessfullyGotFromOrigin bool
	OriginConnectFailed       bool
	OriginBytes               uint64
	ProxyStr                  string
}

type RespondFunc func() (uint64, error)

// NewResponder creates a Responder, which defaults to a generic error response.
func NewResponder(w http.ResponseWriter, r *http.Request, statLog *StatLogger) *Responder {
	responder := &Responder{
		W:                         w,
		R:                         r,
		StatLog:                   statLog,
		Reuse:                     remap.ReuseCannot,
		OriginCode:                0,
		ResponseCode:              http.StatusBadRequest,
		SuccessfullyGotFromOrigin: false,
		OriginConnectFailed:       false,
		OriginBytes:               0,
		ProxyStr:                  "-",
	}
	responder.F = func() (uint64, error) { return web.ServeErr(w, responder.ResponseCode) }
	return responder
}

// SetResponse is a helper which sets the RespondFunc of r to `web.Respond` with the given code, headers, body, and connectionClose
func (r *Responder) SetResponse(code int, hdrs http.Header, body []byte, connectionClose bool) {
	r.F = func() (uint64, error) { return web.Respond(r.W, code, hdrs, body, connectionClose) }
}

// Do responds to the client, according to the data in r, with the given code, headers, and body. It additionally writes to the event log, and adds statistics about this request. This should always be called for the final response to a client, in order to properly log, stat, and other final operations.
// For cache misses, reuse should be ReuseCannot.
// For parent connect failures, originCode should be 0.
func (r *Responder) Do() {
	bytesSent, err := r.F()
	if err != nil {
		log.Errorln(time.Now().Format(time.RFC3339Nano) + " " + r.R.RemoteAddr + " " + r.R.Method + " " + r.R.RequestURI + ": responding: " + err.Error())
	}
	web.TryFlush(r.W)

	successfullyRespondedToClient := err != nil

	r.StatLog.Log(
		r.ResponseCode,
		bytesSent,
		successfullyRespondedToClient,
		r.SuccessfullyGotFromOrigin,
		isCacheHit(r.Reuse, r.OriginCode),
		r.OriginConnectFailed,
		r.OriginCode,
		r.OriginBytes,
		r.ProxyStr)
}

func isCacheHit(reuse remap.Reuse, originCode int) bool {
	return reuse == remap.ReuseCan || ((reuse == remap.ReuseMustRevalidate || reuse == remap.ReuseMustRevalidateCanStale) && (originCode > 299 && originCode < 400))
}
