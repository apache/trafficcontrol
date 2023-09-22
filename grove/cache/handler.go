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
	"os"
	"strconv"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/apache/trafficcontrol/v8/grove/cachedata"
	"github.com/apache/trafficcontrol/v8/grove/plugin"

	"github.com/apache/trafficcontrol/v8/grove/remap"
	"github.com/apache/trafficcontrol/v8/grove/stat"
	"github.com/apache/trafficcontrol/v8/grove/thread"
	"github.com/apache/trafficcontrol/v8/grove/web"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
)

type HandlerPointer struct {
	realHandler *unsafe.Pointer
}

func NewHandlerPointer(realHandler *Handler) *HandlerPointer {
	p := (unsafe.Pointer)(realHandler)
	return &HandlerPointer{realHandler: &p}
}

func (h *HandlerPointer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	realHandler := (*Handler)(atomic.LoadPointer(h.realHandler))
	realHandler.ServeHTTP(w, r)
}

func (h *HandlerPointer) Set(newHandler *Handler) {
	p := (unsafe.Pointer)(newHandler)
	atomic.StorePointer(h.realHandler, p)
}

type Handler struct {
	remapper        remap.HTTPRequestRemapper
	getter          thread.Getter
	ruleThrottlers  map[string]thread.Throttler // doesn't need threadsafe keys, because it's never added to or deleted after creation. TODO fix for hot rule reloading
	scheme          string
	port            string
	hostname        string
	strictRFC       bool
	stats           stat.Stats
	conns           *web.ConnMap
	connectionClose bool
	plugins         plugin.Plugins
	pluginContext   map[string]*interface{}
	httpConns       *web.ConnMap
	httpsConns      *web.ConnMap
	interfaceName   string
	requestID       uint64 // Atomic - DO NOT access or modify without atomic operations
	// keyThrottlers     Throttlers
	// nocacheThrottlers Throttlers
}

// func (h *cacheHandler) checkoutKeyThrottler(k string) Throttler {
// 	keyThrottlersM.Lock()
// 	defer keyThrottlersM.Unlock()
// 	if t, ok := keyThrottlers[k]; !ok {
// 		keyThrottlers[k] = NewThrottler
// 	}
// 	return keyThrottlers[k]
// }

// NewHandler returns an http.Handler object, which may be pipelined with other http.Handlers via `http.ListenAndServe`. If you prefer pipelining functions, use `GetHandlerFunc`.
//
// This needs rate-limited in 3 ways.
// 1. ruleLimit - Simultaneous requests to the origin (remap rule) should be configurably limited. For example, "only allow 1000 simultaneous requests to the origin
// 2. keyLimit - Simultaneous requests, on cache miss, for the same key (Method+Path+Qstring), should be configurably limited. For example, "Only allow 10 simultaneous requests per unique URL on cache miss. Additional requestors must wait until others complete. Once another requestor completes, all waitors for the same URL are signalled to use the cache, or proceed to the third uncacheable limiter"
// 3. nocacheLimit - If simultaneous requestors exceed the URL limiter, and some request for the same key gets a result which is uncacheable, waitors for the same URL may then proceed at a third configurable limit for uncacheable requests.
//
// Note these only apply to cache misses. Cache hits are not limited in any way, the origin is not hit and the cache value is immediately returned to the client.
//
// This prevents a large number of uncacheable requests for the same URL from timing out because they're required to proceed serially from the low simultaneous-requests-per-URL limit, while at the same time only hitting the origin with a very low limit for many simultaneous cacheable requests.
//
// Example: Origin limit is 10,000, key limit is 1, the uncacheable limit is 1,000.
// Then, 2,000 requests come in for the same URL, simultaneously. They are all within the Origin limit, so they are all allowed to proceed to the key limiter. Then, the first request is allowed to make an actual request to the origin, while the other 1,999 wait at the key limiter.
//
// The connectionClose parameter determines whether to send a `Connection: close` header. This is primarily designed for maintenance, to drain the cache of incoming requestors. This overrides rule-specific `connection-close: false` configuration, under the assumption that draining a cache is a temporary maintenance operation, and if connectionClose is true on the service and false on some rules, those rules' configuration is probably a permament setting whereas the operator probably wants to drain all connections if the global setting is true. If it's necessary to leave connection close false on some rules, set all other rules' connectionClose to true and leave the global connectionClose unset.
func NewHandler(
	remapper remap.HTTPRequestRemapper,
	ruleLimit uint64,
	stats stat.Stats,
	scheme string,
	port string,
	conns *web.ConnMap,
	strictRFC bool,
	connectionClose bool,
	plugins plugin.Plugins,
	pluginContext map[string]*interface{},
	httpConns *web.ConnMap,
	httpsConns *web.ConnMap,
	interfaceName string,
) *Handler {
	hostname, err := os.Hostname()
	if err != nil {
		log.Errorf("getting  hostname: %v\n", err)
	}

	return &Handler{
		remapper:        remapper,
		getter:          thread.NewGetter(),
		ruleThrottlers:  makeRuleThrottlers(remapper, ruleLimit),
		strictRFC:       strictRFC,
		scheme:          scheme,
		port:            port,
		hostname:        hostname,
		stats:           stats,
		conns:           conns,
		connectionClose: connectionClose,
		plugins:         plugins,
		pluginContext:   pluginContext,
		httpConns:       httpConns,
		httpsConns:      httpsConns,
		interfaceName:   interfaceName,
		// keyThrottlers:     NewThrottlers(keyLimit),
		// nocacheThrottlers: NewThrottlers(nocacheLimit),
	}
}

func makeRuleThrottlers(remapper remap.HTTPRequestRemapper, limit uint64) map[string]thread.Throttler {
	remapRules := remapper.Rules()
	ruleThrottlers := make(map[string]thread.Throttler, len(remapRules))
	for _, rule := range remapRules {
		ruleLimit := uint64(rule.ConcurrentRuleRequests)
		if rule.ConcurrentRuleRequests == 0 {
			ruleLimit = limit
		}
		ruleThrottlers[rule.Name] = thread.NewThrottler(ruleLimit)
	}
	return ruleThrottlers
}

func copyPluginContext(context map[string]*interface{}) map[string]*interface{} {
	new := make(map[string]*interface{}, len(context))
	for k, v := range context {
		newV := *v
		new[k] = &newV
	}
	return new
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqTime := time.Now()
	reqID := atomic.AddUint64(&h.requestID, 1)
	pluginContext := copyPluginContext(h.pluginContext) // must give each request a copy, because they can modify in parallel
	srvrData := cachedata.SrvrData{Hostname: h.hostname, Port: h.port, Scheme: h.scheme}
	onReqData := plugin.OnRequestData{W: w, R: r, Stats: h.stats, StatRules: h.remapper.StatRules(), HTTPConns: h.httpConns, HTTPSConns: h.httpsConns, InterfaceName: h.interfaceName, SrvrData: srvrData, RequestID: reqID}
	stop := h.plugins.OnRequest(h.remapper.PluginCfg(), pluginContext, onReqData)
	if stop {
		return
	}

	conn := (*web.InterceptConn)(nil)
	if realConn, ok := h.conns.Get(r.RemoteAddr); !ok {
		log.Infof("RemoteAddr '%v' not in Conns (reqid %v)\n", r.RemoteAddr, reqID)
	} else {
		if conn, ok = realConn.(*web.InterceptConn); !ok {
			log.Infof("Could not get Conn info: Conn is not an InterceptConn: %T (reqid %v)\n", realConn, reqID)
		}
	}

	remappingProducer, err := h.remapper.RemappingProducer(r, h.scheme)

	if err == nil { // if we failed to get a remapping, there's no DSCP to set.
		if err := conn.SetDSCP(remappingProducer.DSCP()); err != nil {
			log.Infoln(time.Now().Format(time.RFC3339Nano) + " " + r.RemoteAddr + " " + r.Method + " " + r.RequestURI + ": could not set DSCP: " + err.Error() + " (reqid " + strconv.FormatUint(reqID, 10) + ")")
		}
	}

	reqHeader := web.CopyHeader(r.Header) // copy request header, because it's not guaranteed valid after actually issuing the request
	clientIP, _ := web.GetClientIPPort(r)

	toFQDN := ""
	pluginCfg := map[string]interface{}{}
	if remappingProducer != nil {
		toFQDN = remappingProducer.FirstFQDN()
		pluginCfg = remappingProducer.PluginCfg()
	}

	reqData := cachedata.ReqData{Req: r, Conn: conn, ClientIP: clientIP, ReqTime: reqTime, ToFQDN: toFQDN}
	responder := NewResponder(w, pluginCfg, pluginContext, srvrData, reqData, h.plugins, h.stats, reqID)

	if err != nil {
		switch err {
		case remap.ErrRuleNotFound:
			log.Debugf("rule not found for %v (reqid %v)\n", r.RequestURI, reqID)
			*responder.ResponseCode = http.StatusNotFound
		case remap.ErrIPNotAllowed:
			log.Debugf("IP %v not allowed (reqid %v)\n", r.RemoteAddr, reqID)
			*responder.ResponseCode = http.StatusForbidden
		default:
			log.Debugf("request error: %v (reqid %v)\n", err, reqID)
		}
		responder.OriginConnectFailed = true
		responder.Do()
		return
	}

	reqCacheControl := rfc.ParseCacheControl(reqHeader)
	log.Debugf("Serve got Cache-Control %+v (reqid %v)\n", reqCacheControl, reqID)

	connectionClose := h.connectionClose || remappingProducer.ConnectionClose()

	beforeCacheLookUpData := plugin.BeforeCacheLookUpData{Req: r, DefaultCacheKey: remappingProducer.CacheKey(), CacheKeyOverrideFunc: remappingProducer.OverrideCacheKey}
	h.plugins.OnBeforeCacheLookup(remappingProducer.PluginCfg(), pluginContext, beforeCacheLookUpData)

	cacheKey := remappingProducer.CacheKey()
	retrier := NewRetrier(h, reqHeader, reqTime, reqCacheControl, remappingProducer, reqID)

	cache := remappingProducer.Cache()

	var reqHost *string
	cacheObj, ok := cache.Get(cacheKey)
	if !ok {
		log.Debugf("cache.Handler.ServeHTTP: '%v' not in cache (reqid %v)\n", cacheKey, reqID)
		beforeParentRequestData := plugin.BeforeParentRequestData{Req: r, RemapRule: remappingProducer.Name()}
		h.plugins.OnBeforeParentRequest(remappingProducer.PluginCfg(), pluginContext, beforeParentRequestData)
		cacheObj, reqHost, err = retrier.Get(r, nil)
		if err != nil {
			log.Errorf("retrying get error (in uncached): %v (reqid %v)\n", err, reqID)
			responder.OriginConnectFailed = true
			responder.ProxyStr = cacheObj.ProxyURL
			if reqHost != nil {
				responder.ToFQDN = *reqHost
			}
			responder.Do()
			return
		}

		responder.OriginCode = cacheObj.OriginCode
		// create new pointers, so plugins don't modify the cacheObj
		codePtr, hdrsPtr, bodyPtr := cacheObj.Code, cacheObj.RespHeaders, cacheObj.Body
		responder.SetResponse(&codePtr, &hdrsPtr, &bodyPtr, connectionClose)
		responder.OriginReqSuccess = true
		responder.ProxyStr = cacheObj.ProxyURL
		if reqHost != nil {
			responder.ToFQDN = *reqHost
		}
		beforeRespData := plugin.BeforeRespondData{Req: r, CacheObj: cacheObj, Code: &codePtr, Hdr: &hdrsPtr, Body: &bodyPtr, RemapRule: remappingProducer.Name()}
		h.plugins.OnBeforeRespond(remappingProducer.PluginCfg(), pluginContext, beforeRespData)
		responder.Do()
		return
	}

	reqHeaders := r.Header
	canReuseStored := rfc.CanReuseStored(reqHeaders, cacheObj.RespHeaders, reqCacheControl, cacheObj.RespCacheControl, cacheObj.ReqHeaders, cacheObj.ReqRespTime, cacheObj.RespRespTime, h.strictRFC)

	if canReuseStored != rfc.ReuseCan { // run the BeforeParentRequest hook for revalidations / ReuseCannot
		beforeParentRequestData := plugin.BeforeParentRequestData{Req: r, RemapRule: remappingProducer.Name()}
		h.plugins.OnBeforeParentRequest(remappingProducer.PluginCfg(), pluginContext, beforeParentRequestData)
	}

	switch canReuseStored {
	case rfc.ReuseCan:
		log.Debugf("cache.Handler.ServeHTTP: '%v' cache hit! (reqid %v)\n", cacheKey, reqID)
	case rfc.ReuseCannot:
		log.Debugf("cache.Handler.ServeHTTP: '%v' can't reuse (reqid %v)\n", cacheKey, reqID)
		cacheObj, reqHost, err = retrier.Get(r, nil)
		if err != nil {
			log.Errorf("retrying get error (in reuse-cannot): %v (reqid %v)\n", err, reqID)
			responder.Do()
			return
		}
	case rfc.ReuseMustRevalidate:
		log.Debugf("cache.Handler.ServeHTTP: '%v' must revalidate (reqid %v)\n", cacheKey, reqID)
		cacheObj, reqHost, err = retrier.Get(r, cacheObj)
		if err != nil {
			log.Errorf("retrying get error: %v (reqid %v)\n", err, reqID)
			responder.Do()
			return
		}
	case rfc.ReuseMustRevalidateCanStale:
		log.Debugf("cache.Handler.ServeHTTP: '%v' must revalidate (but allowed stale) (reqid %v)\n", cacheKey, reqID)
		oldCacheObj := cacheObj
		cacheObj, reqHost, err = retrier.Get(r, cacheObj)
		if err != nil {
			log.Errorf("retrying get error - serving stale as allowed: %v (reqid %v)\n", err, reqID)
			cacheObj = oldCacheObj
		}
	}
	log.Debugf("cache.Handler.ServeHTTP: '%v' responding with %v (reqid %v)\n", cacheKey, cacheObj.Code, reqID)

	// create new pointers, so plugins don't modify the cacheObj
	codePtr, hdrsPtr, bodyPtr := cacheObj.Code, cacheObj.RespHeaders, cacheObj.Body
	responder.SetResponse(&codePtr, &hdrsPtr, &bodyPtr, connectionClose)
	responder.OriginReqSuccess = true
	responder.Reuse = canReuseStored
	responder.OriginCode = cacheObj.OriginCode
	responder.OriginBytes = cacheObj.Size
	responder.ProxyStr = cacheObj.ProxyURL
	if reqHost != nil {
		responder.ToFQDN = *reqHost
	}
	beforeRespData := plugin.BeforeRespondData{Req: r, CacheObj: cacheObj, Code: &codePtr, Hdr: &hdrsPtr, Body: &bodyPtr, RemapRule: remappingProducer.Name()}
	h.plugins.OnBeforeRespond(remappingProducer.PluginCfg(), pluginContext, beforeRespData)
	responder.Do()
}
