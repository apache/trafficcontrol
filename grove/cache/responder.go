package cache

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

	"github.com/apache/trafficcontrol/v8/grove/cachedata"
	"github.com/apache/trafficcontrol/v8/grove/plugin"
	"github.com/apache/trafficcontrol/v8/grove/stat"
	"github.com/apache/trafficcontrol/v8/grove/web"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
)

// Responder is an object encapsulating the cache's response to the client. It holds all the data necessary to respond, log the response, and add the stats.
type Responder struct {
	W             http.ResponseWriter
	RequestID     uint64
	PluginCfg     map[string]interface{}
	Plugins       plugin.Plugins
	PluginContext map[string]*interface{}
	Stats         stat.Stats
	F             RespondFunc
	ResponseCode  *int
	cachedata.ParentRespData
	cachedata.SrvrData
	cachedata.ReqData
}

func DefaultParentRespData() cachedata.ParentRespData {
	return cachedata.ParentRespData{
		Reuse:               rfc.ReuseCannot,
		OriginCode:          0,
		OriginReqSuccess:    false,
		OriginConnectFailed: false,
		OriginBytes:         0,
		ProxyStr:            "-",
	}
}

func DefaultRespCode() *int {
	c := http.StatusBadRequest
	return &c
}

type RespondFunc func() (uint64, error)

// NewResponder creates a Responder, which defaults to a generic error response.
func NewResponder(w http.ResponseWriter, pluginCfg map[string]interface{}, pluginContext map[string]*interface{}, srvrData cachedata.SrvrData, reqData cachedata.ReqData, plugins plugin.Plugins, stats stat.Stats, reqID uint64) *Responder {
	responder := &Responder{
		W:              w,
		RequestID:      reqID,
		PluginCfg:      pluginCfg,
		Plugins:        plugins,
		PluginContext:  pluginContext,
		Stats:          stats,
		ResponseCode:   DefaultRespCode(),
		ParentRespData: DefaultParentRespData(),
		SrvrData:       srvrData,
		ReqData:        reqData,
	}
	responder.F = func() (uint64, error) { return web.ServeErr(w, *responder.ResponseCode) }
	return responder
}

// SetResponse is a helper which sets the RespondFunc of r to `web.Respond` with the given code, headers, body, and connectionClose. Note it takes a pointer to the headers and body, which may be modified after calling this but before the Do() sends the response.
func (r *Responder) SetResponse(code *int, hdrs *http.Header, body *[]byte, connectionClose bool) {
	r.ResponseCode = code
	r.F = func() (uint64, error) {
		if r.Req.Method == http.MethodHead {
			*body = nil
		}
		return web.Respond(r.W, *code, *hdrs, *body, connectionClose)
	}
}

// Do responds to the client, according to the data in r, with the given code, headers, and body. It additionally writes to the event log, and adds statistics about this request. This should always be called for the final response to a client, in order to properly log, stat, and other final operations.
// For cache misses, reuse should be ReuseCannot.
// For parent connect failures, originCode should be 0.
func (r *Responder) Do() {
	// TODO move plugins.BeforeRespond here? How do we distinguish between success, and know to set headers? r.OriginReqSuccess?
	bytesSent, err := r.F()
	if err != nil {
		log.Errorf("%s %s %s %v : responding: %v", r.Req.RemoteAddr, r.Req.Method, r.Req.RequestURI, r.ResponseCode, err.Error())
	}
	web.TryFlush(r.W) // TODO remove? Let plugins do it, if they need to?

	respSuccess := err != nil
	respData := cachedata.RespData{RespCode: *r.ResponseCode, BytesWritten: bytesSent, RespSuccess: respSuccess, CacheHit: isCacheHit(r.Reuse, r.OriginCode)}
	arData := plugin.AfterRespondData{W: r.W, Stats: r.Stats, ReqData: r.ReqData, SrvrData: r.SrvrData, ParentRespData: r.ParentRespData, RespData: respData, RequestID: r.RequestID}
	r.Plugins.OnAfterRespond(r.PluginCfg, r.PluginContext, arData)
}

func isCacheHit(reuse rfc.Reuse, originCode int) bool {
	// TODO move to web? remap?
	return reuse == rfc.ReuseCan || ((reuse == rfc.ReuseMustRevalidate || reuse == rfc.ReuseMustRevalidateCanStale) && originCode == http.StatusNotModified)
}
