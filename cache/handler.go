package cache

import (
	"net"
	"net/http"
	"os"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/apache/incubator-trafficcontrol/grove/cachedata"
	"github.com/apache/incubator-trafficcontrol/grove/plugin"
	"github.com/apache/incubator-trafficcontrol/grove/plugin/beforerespond"
	"github.com/apache/incubator-trafficcontrol/grove/remap"
	"github.com/apache/incubator-trafficcontrol/grove/remapdata"
	"github.com/apache/incubator-trafficcontrol/grove/stat"
	"github.com/apache/incubator-trafficcontrol/grove/thread"
	"github.com/apache/incubator-trafficcontrol/grove/web"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
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
	transport       *http.Transport
	plugins         plugin.Plugins
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
	reqTimeout time.Duration,
	reqKeepAlive time.Duration,
	reqMaxIdleConns int,
	reqIdleConnTimeout time.Duration,
	plugins plugin.Plugins,
) *Handler {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   reqTimeout,
			KeepAlive: reqKeepAlive,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          reqMaxIdleConns,
		IdleConnTimeout:       reqIdleConnTimeout,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	transport.Dial = func(network, address string) (net.Conn, error) {
		d := net.Dialer{DualStack: true, FallbackDelay: time.Millisecond * 50}
		return d.Dial(network, address)
	}

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
		transport:       transport,
		plugins:         plugins,
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

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqTime := time.Now()

	conn := (*web.InterceptConn)(nil)
	if realConn, ok := h.conns.Get(r.RemoteAddr); !ok {
		log.Errorf("RemoteAddr '%v' not in Conns\n", r.RemoteAddr)
	} else {
		if conn, ok = realConn.(*web.InterceptConn); !ok {
			log.Errorf("Could not get Conn info: Conn is not an InterceptConn: %T\n", realConn)
		}
	}

	remappingProducer, err := h.remapper.RemappingProducer(r, h.scheme)

	if err == nil { // if we failed to get a remapping, there's no DSCP to set.
		if err := conn.SetDSCP(remappingProducer.DSCP()); err != nil {
			log.Errorln(time.Now().Format(time.RFC3339Nano) + " " + r.RemoteAddr + " " + r.Method + " " + r.RequestURI + ": could not set DSCP: " + err.Error())
		}
	}

	reqHeader := web.CopyHeader(r.Header) // copy request header, because it's not guaranteed valid after actually issuing the request
	clientIP, _ := web.GetClientIPPort(r)

	toFQDN := ""
	pluginCfg := map[string]interface{}{}
	if remappingProducer != nil {
		toFQDN = remappingProducer.ToFQDN()
		pluginCfg = remappingProducer.PluginCfg()
	}

	reqData := cachedata.ReqData{r, conn, clientIP, reqTime, toFQDN}
	srvrData := cachedata.SrvrData{h.hostname, h.port, h.scheme}

	responder := NewResponder(w, pluginCfg, srvrData, reqData, h.plugins, h.stats)

	if err != nil {
		switch err {
		case remap.ErrRuleNotFound:
			log.Debugf("rule not found for %v\n", r.RequestURI)
			*responder.ResponseCode = http.StatusNotFound
		case remap.ErrIPNotAllowed:
			log.Debugf("IP %v not allowed\n", r.RemoteAddr)
			*responder.ResponseCode = http.StatusForbidden
		default:
			log.Debugf("request error: %v\n", err)
		}
		responder.OriginConnectFailed = true
		responder.Do()
		return
	}

	reqCacheControl := web.ParseCacheControl(reqHeader)
	log.Debugf("Serve got Cache-Control %+v\n", reqCacheControl)

	connectionClose := h.connectionClose || remappingProducer.ConnectionClose()
	cacheKey := remappingProducer.CacheKey()
	retrier := NewRetrier(h, reqHeader, reqTime, reqCacheControl, cacheKey, remappingProducer)

	cache := remappingProducer.Cache()

	cacheObj, ok := cache.Get(cacheKey)
	if !ok {
		log.Debugf("cache.Handler.ServeHTTP: '%v' not in cache\n", cacheKey)
		remappingProducer.ToOriginHdrs().Mod(&r.Header)
		cacheObj, err := retrier.Get(r, nil)
		if err != nil {
			log.Errorf("retrying get error (in uncached): %v\n", err)
			responder.OriginConnectFailed = true
			responder.ProxyStr = cacheObj.ProxyURL
			responder.Do()
			return
		}

		responder.OriginCode = cacheObj.OriginCode
		// create new pointers, so plugins don't modify the cacheObj
		codePtr, hdrsPtr, bodyPtr := cacheObj.Code, cacheObj.RespHeaders, cacheObj.Body
		responder.SetResponse(&codePtr, &hdrsPtr, &bodyPtr, connectionClose)
		responder.OriginReqSuccess = true
		responder.ProxyStr = cacheObj.ProxyURL
		beforeRespData := beforerespond.Data{r, cacheObj, &codePtr, &hdrsPtr, &bodyPtr}
		h.plugins.BeforeRespond.Call(remappingProducer.PluginCfg(), beforeRespData)
		responder.Do()
		return
	}

	reqHeaders := r.Header
	canReuseStored := remap.CanReuseStored(reqHeaders, cacheObj.RespHeaders, reqCacheControl, cacheObj.RespCacheControl, cacheObj.ReqHeaders, cacheObj.ReqRespTime, cacheObj.RespRespTime, h.strictRFC)

	switch canReuseStored {
	case remapdata.ReuseCan:
		log.Debugf("cache.Handler.ServeHTTP: '%v' cache hit!\n", cacheKey)
	case remapdata.ReuseCannot:
		log.Debugf("cache.Handler.ServeHTTP: '%v' can't reuse\n", cacheKey)
		cacheObj, err = retrier.Get(r, nil)
		if err != nil {
			log.Errorf("retrying get error (in reuse-cannot): %v\n", err)
			responder.Do()
			return
		}
	case remapdata.ReuseMustRevalidate:
		log.Debugf("cache.Handler.ServeHTTP: '%v' must revalidate\n", cacheKey)
		cacheObj, err = retrier.Get(r, cacheObj)
		if err != nil {
			log.Errorf("retrying get error: %v\n", err)
			responder.Do()
			return
		}
	case remapdata.ReuseMustRevalidateCanStale:
		log.Debugf("cache.Handler.ServeHTTP: '%v' must revalidate (but allowed stale)\n", cacheKey)
		oldCacheObj := cacheObj
		cacheObj, err = retrier.Get(r, cacheObj)
		if err != nil {
			log.Errorf("retrying get error - serving stale as allowed: %v\n", err)
			cacheObj = oldCacheObj
		}
	}
	log.Debugf("cache.Handler.ServeHTTP: '%v' responding with %v\n", cacheKey, cacheObj.Code)

	// create new pointers, so plugins don't modify the cacheObj
	codePtr, hdrsPtr, bodyPtr := cacheObj.Code, cacheObj.RespHeaders, cacheObj.Body
	responder.SetResponse(&codePtr, &hdrsPtr, &bodyPtr, connectionClose)
	responder.OriginReqSuccess = true
	responder.Reuse = canReuseStored
	responder.OriginCode = cacheObj.OriginCode
	responder.OriginBytes = cacheObj.Size
	responder.ProxyStr = cacheObj.ProxyURL
	beforeRespData := beforerespond.Data{r, cacheObj, &codePtr, &hdrsPtr, &bodyPtr}
	h.plugins.BeforeRespond.Call(remappingProducer.PluginCfg(), beforeRespData)
	responder.Do()
}
