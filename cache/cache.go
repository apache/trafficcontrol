package cache

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"sync/atomic"
	"time"
	"unsafe"

	cacheobj "github.com/apache/incubator-trafficcontrol/grove/cacheobj"
	"github.com/apache/incubator-trafficcontrol/grove/thread"
	"github.com/apache/incubator-trafficcontrol/grove/web"

	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/common/log"
)

// TODO add logging

type Cache interface {
	AddSize(key string, value interface{}, size uint64) bool
	Get(key string) (interface{}, bool)
	Remove(key string)
	RemoveOldest()
	Size() uint64
}

type CacheHandlerPointer struct {
	realHandler *unsafe.Pointer
}

func NewCacheHandlerPointer(realHandler *CacheHandler) *CacheHandlerPointer {
	p := (unsafe.Pointer)(realHandler)
	return &CacheHandlerPointer{realHandler: &p}
}

func (h *CacheHandlerPointer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	realHandler := (*CacheHandler)(atomic.LoadPointer(h.realHandler))
	realHandler.TryServe(w, r)
}

func (h *CacheHandlerPointer) Set(newHandler *CacheHandler) {
	p := (unsafe.Pointer)(newHandler)
	atomic.StorePointer(h.realHandler, p)
}

type CacheHandler struct {
	cache           Cache
	remapper        HTTPRequestRemapper
	getter          thread.Getter
	ruleThrottlers  map[string]thread.Throttler // doesn't need threadsafe keys, because it's never added to or deleted after creation. TODO fix for hot rule reloading
	scheme          string
	strictRFC       bool
	stats           Stats
	conns           *web.ConnMap
	connectionClose bool
	transport       *http.Transport
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
// ruleLimit uint64, keyLimit uint64, nocacheLimit uint64
//
// The connectionClose parameter determines whether to send a `Connection: close` header. This is primarily designed for maintenance, to drain the cache of incoming requestors. This overrides rule-specific `connection-close: false` configuration, under the assumption that draining a cache is a temporary maintenance operation, and if connectionClose is true on the service and false on some rules, those rules' configuration is probably a permament setting whereas the operator probably wants to drain all connections if the global setting is true. If it's necessary to leave connection close false on some rules, set all other rules' connectionClose to true and leave the global connectionClose unset.
func NewCacheHandler(
	cache Cache,
	remapper HTTPRequestRemapper,
	ruleLimit uint64,
	stats Stats,
	scheme string,
	conns *web.ConnMap,
	strictRFC bool,
	connectionClose bool,
	reqTimeout time.Duration,
	reqKeepAlive time.Duration,
	reqMaxIdleConns int,
	reqIdleConnTimeout time.Duration,
) *CacheHandler {
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

	return &CacheHandler{
		cache:           cache,
		remapper:        remapper,
		getter:          thread.NewGetter(),
		ruleThrottlers:  makeRuleThrottlers(remapper, ruleLimit),
		strictRFC:       strictRFC,
		scheme:          scheme,
		stats:           stats,
		conns:           conns,
		connectionClose: connectionClose,
		transport:       transport,
		// keyThrottlers:     NewThrottlers(keyLimit),
		// nocacheThrottlers: NewThrottlers(nocacheLimit),
	}
}

func makeRuleThrottlers(remapper HTTPRequestRemapper, limit uint64) map[string]thread.Throttler {
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

// NewCacheHandlerFunc creates and returns an http.HandleFunc, which may be pipelined with other http.HandleFuncs via `http.HandleFunc`. This is a convenience wrapper around the `http.Handler` object obtainable via `New`. If you prefer objects, use `NewCacheHandler`.
func NewCacheHandlerFunc(
	cache Cache,
	remapper HTTPRequestRemapper,
	ruleLimit uint64,
	stats Stats,
	scheme string,
	conns *web.ConnMap,
	strictRFC bool,
	connectionClose bool,
	reqTimeout time.Duration,
	reqKeepAlive time.Duration,
	reqMaxIdleConns int,
	reqIdleConnTimeout time.Duration,
) http.HandlerFunc {
	handler := NewCacheHandler(
		cache,
		remapper,
		ruleLimit,
		stats,
		scheme,
		conns,
		strictRFC,
		connectionClose,
		reqTimeout,
		reqKeepAlive,
		reqMaxIdleConns,
		reqIdleConnTimeout,
	)
	return func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	}
}

func (h *CacheHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.TryServe(w, r)
}

const CodeConnectFailure = http.StatusBadGateway

func isFailure(o *cacheobj.CacheObj, retryCodes map[int]struct{}) bool {
	_, failureCode := retryCodes[o.Code]
	return failureCode || o.Code == CodeConnectFailure
}

// RetryingGet takes a function, and retries failures up to the RemappingProducer RetryNum limit. On failure, it creates a new remapping. The func f should use `remapping` to make its request. If it hits failures up to the limit, it returns the last received cacheobj.CacheObj
// TODO refactor to not close variables - it's awkward and confusing.
func RetryingGet(getCacheObj func(remapping Remapping, retryFailures bool, obj *cacheobj.CacheObj) *cacheobj.CacheObj, request *http.Request, remappingProducer *RemappingProducer, cachedObj *cacheobj.CacheObj) (*cacheobj.CacheObj, error) {
	obj := (*cacheobj.CacheObj)(nil)
	for {
		remapping, retryAllowed, err := remappingProducer.GetNext(request)
		if err == ErrNoMoreRetries {
			if obj == nil {
				return nil, fmt.Errorf("remapping producer allows no requests") // should never happen
			}
			return obj, nil
		} else if err != nil {
			return nil, err
		}
		obj = getCacheObj(remapping, retryAllowed, cachedObj)
		if !isFailure(obj, remapping.RetryCodes) {
			return obj, nil
		}
	}
}

// GetAndCache makes a client request for the given `http.Request` and caches it if `CanCache`.
// THe `ruleThrottler` may be nil, in which case the request will be unthrottled.
func GetAndCache(
	req *http.Request,
	proxyURL *url.URL,
	cacheKey string,
	remapName string,
	reqHeader http.Header,
	reqTime time.Time,
	strictRFC bool,
	cache Cache,
	ruleThrottler thread.Throttler,
	revalidateObj *cacheobj.CacheObj,
	timeout time.Duration,
	cacheFailure bool,
	retryNum int,
	retryCodes map[int]struct{},
	transport *http.Transport,
) *cacheobj.CacheObj {
	// TODO this is awkward, with 'revalidateObj' indicating whether the request is a Revalidate. Should Getting and Caching be split up? How?
	get := func() *cacheobj.CacheObj {
		// TODO figure out why respReqTime isn't used by rules
		log.Debugf("GetAndCache calling request %v %v %v %v %v\n", req.Method, req.URL.Scheme, req.URL.Host, req.URL.EscapedPath(), req.Header)
		// TODO Verify overriding the passed reqTime is the right thing to do
		respCode, respHeader, respBody, reqTime, reqRespTime, err := request(transport, req, proxyURL)
		if err != nil {
			log.Errorf("Parent error for URI %v %v %v cacheKey %v rule %v parent %v error %v\n", req.URL.Scheme, req.URL.Host, req.URL.EscapedPath(), cacheKey, remapName, proxyURL.String(), err)
			code := CodeConnectFailure
			body := []byte(http.StatusText(code))
			return cacheobj.New(reqHeader, body, code, respHeader, reqTime, reqRespTime, reqRespTime)
		}
		if _, ok := retryCodes[respCode]; ok && !cacheFailure {
			return cacheobj.New(reqHeader, respBody, respCode, respHeader, reqTime, reqRespTime, reqRespTime)
		}

		log.Debugf("GetAndCache request returned %v headers %+v\n", respCode, respHeader)
		respRespTime, ok := GetHTTPDate(respHeader, "Date")
		if !ok {
			log.Errorf("request %v returned no Date header - RFC Violation!\n", req.RequestURI)
			respRespTime = reqRespTime // if no Date was returned using the client response time simulates latency 0
		}

		obj := (*cacheobj.CacheObj)(nil)
		// TODO This means if we can't cache the object, we return nil. Verify this is ok
		if !CanCache(reqHeader, respCode, respHeader, strictRFC) {
			return cacheobj.New(reqHeader, respBody, respCode, respHeader, reqTime, reqRespTime, reqRespTime)
		}
		log.Debugf("h.cache.AddSize %v\n", cacheKey)
		log.Debugf("GetAndCache respCode %v\n", respCode)
		if revalidateObj == nil || respCode < 300 || respCode > 399 {
			log.Debugf("GetAndCache new %v\n", cacheKey)
			obj = cacheobj.New(reqHeader, respBody, respCode, respHeader, reqTime, reqRespTime, respRespTime)
		} else {
			log.Debugf("GetAndCache revalidating %v\n", cacheKey)
			// must copy, because this cache object may be concurrently read by other goroutines
			newRespHeader := http.Header{}
			web.CopyHeader(revalidateObj.RespHeaders, &newRespHeader)
			newRespHeader.Set("Date", respHeader.Get("Date"))
			obj = &cacheobj.CacheObj{
				Body:             revalidateObj.Body,
				ReqHeaders:       revalidateObj.ReqHeaders,
				RespHeaders:      newRespHeader,
				RespCacheControl: revalidateObj.RespCacheControl,
				Code:             revalidateObj.Code,
				ReqTime:          reqTime,
				ReqRespTime:      reqRespTime,
				RespRespTime:     respRespTime,
				Size:             revalidateObj.Size,
			}
		}
		cache.AddSize(cacheKey, obj, obj.Size) // TODO store pointer?
		return obj
	}

	c := (*cacheobj.CacheObj)(nil)
	if ruleThrottler == nil {
		log.Errorf("rule %v not in ruleThrottlers map. Requesting with no origin limit!\n", remapName)
		ruleThrottler = thread.NewNoThrottler()
	}
	ruleThrottler.Throttle(func() { c = get() })
	return c
}

func CanReuse(reqHeader http.Header, reqCacheControl web.CacheControl, cacheObj *cacheobj.CacheObj, strictRFC bool, revalidateCanReuse bool) bool {
	canReuse := CanReuseStored(reqHeader, cacheObj.RespHeaders, reqCacheControl, cacheObj.RespCacheControl, cacheObj.ReqHeaders, cacheObj.ReqRespTime, cacheObj.RespRespTime, strictRFC)
	return canReuse == ReuseCan || (canReuse == ReuseMustRevalidate && revalidateCanReuse)
}

// TryServe attempts to serve the given request, as a caching reverse proxy.
// Serving acts as a state machine.
func (h *CacheHandler) TryServe(w http.ResponseWriter, r *http.Request) {
	log.EventRaw(time.Now().Format(time.RFC3339Nano) + " " + r.RemoteAddr + " " + r.Method + " " + r.RequestURI + "\n")

	h.stats.IncConnections()
	defer h.stats.DecConnections()

	// inBytes := getBytes(r)
	reqTime := time.Now()

	// copy request header, because it's not guaranteed valid after actually issuing the request
	reqHeader := http.Header{}
	web.CopyHeader(r.Header, &reqHeader)

	// ok = 'rule found'
	// remappedReq, remapName, cacheKey, allowed, ruleConnectionClose, ok, err := h.remapper.Remap(r, h.scheme, 0) // TODO handle failures

	remappingProducer, err := h.remapper.RemappingProducer(r, h.scheme)
	if err != nil {
		if err == ErrRuleNotFound {
			log.Debugf("rule not found for %v\n", r.RequestURI)
			h.serveRuleNotFound(w)
		} else if err == ErrIPNotAllowed {
			log.Debugf("IP %v not allowed\n", r.RemoteAddr)
			h.serveNotAllowed(w)
		} else {
			log.Debugf("request error: %v\n", err)
			h.serveReqErr(w)
		}
		return
	}

	reqCacheControl := web.ParseCacheControl(reqHeader)
	log.Debugf("TryServe got Cache-Control %+v\n", reqCacheControl)

	connectionClose := h.connectionClose || remappingProducer.ConnectionClose()
	cacheKey := remappingProducer.CacheKey()
	reqFQDN := r.Host

	retryGetFunc := func(remapping Remapping, retryFailures bool, obj *cacheobj.CacheObj) *cacheobj.CacheObj {
		// return true for Revalidate, and issue revalidate requests separately.
		canReuse := func(cacheObj *cacheobj.CacheObj) bool {
			return CanReuse(reqHeader, reqCacheControl, cacheObj, h.strictRFC, true)
		}

		getAndCache := func() *cacheobj.CacheObj {
			return GetAndCache(remapping.Request, remapping.ProxyURL, remapping.CacheKey, remapping.Name, remapping.Request.Header, reqTime, h.strictRFC, h.cache, h.ruleThrottlers[remapping.Name], obj, remapping.Timeout, retryFailures, remapping.RetryNum, remapping.RetryCodes, h.transport)
		}

		return h.getter.Get(cacheKey, getAndCache, canReuse)
	}

	retryingGet := func(r *http.Request, obj *cacheobj.CacheObj) (*cacheobj.CacheObj, error) {
		return RetryingGet(retryGetFunc, r, remappingProducer, obj)
	}

	iCacheObj, ok := h.cache.Get(cacheKey)

	if !ok {
		log.Debugf("cacheHandler.ServeHTTP: '%v' not in cache\n", cacheKey)

		// func RetryingGet(getcacheobj.CacheObj func(remapping Remapping, retryFailures bool) *cacheobj.CacheObj, request *http.Request, remappingProducer *RemappingProducer) (*cacheobj.CacheObj, error) {

		cacheObj, err := retryingGet(r, nil)
		if err != nil {
			log.Errorf("retrying get error (in uncached): %v\n", err)
			h.serveReqErr(w)
			return
		}

		h.respond(w, cacheObj.Code, cacheObj.RespHeaders, cacheObj.Body, h.stats, h.conns, r.RemoteAddr, reqFQDN, connectionClose)
		return
	}

	cacheObj, ok := iCacheObj.(*cacheobj.CacheObj)
	if !ok {
		// should never happen
		log.Errorf("cache key '%v' value '%v' type '%T' expected *cacheobj.CacheObj\n", cacheKey, iCacheObj, iCacheObj)
		cacheObj, err = retryingGet(r, nil)
		if err != nil {
			log.Errorf("retrying get error (in unexpected cacheobj): %v\n", err)
			h.serveReqErr(w)
			return
		}

		// TODO check for ReuseMustRevalidate
		h.respond(w, cacheObj.Code, cacheObj.RespHeaders, cacheObj.Body, h.stats, h.conns, r.RemoteAddr, reqFQDN, connectionClose)
		return
	}

	reqHeaders := r.Header

	canReuseStored := CanReuseStored(reqHeaders, cacheObj.RespHeaders, reqCacheControl, cacheObj.RespCacheControl, cacheObj.ReqHeaders, cacheObj.ReqRespTime, cacheObj.RespRespTime, h.strictRFC)

	switch canReuseStored {
	case ReuseCan:
		log.Debugf("cacheHandler.ServeHTTP: '%v' cache hit!\n", cacheKey)
	case ReuseCannot:
		log.Debugf("cacheHandler.ServeHTTP: '%v' can't reuse\n", cacheKey)
		cacheObj, err = retryingGet(r, nil)
		if err != nil {
			log.Errorf("retrying get error (in reuse-cannot): %v\n", err)
			h.serveReqErr(w)
			return
		}
	case ReuseMustRevalidate:
		log.Debugf("cacheHandler.ServeHTTP: '%v' must revalidate\n", cacheKey)
		// r := remapping.Request
		// TODO verify setting the existing request header here works
		r.Header.Set("If-Modified-Since", cacheObj.RespRespTime.Format(time.RFC1123))
		cacheObj, err = retryingGet(r, cacheObj)
		if err != nil {
			log.Errorf("retrying get error: %v\n", err)
			h.serveReqErr(w)
			return
		}
	case ReuseMustRevalidateCanStale:
		log.Debugf("cacheHandler.ServeHTTP: '%v' must revalidate (but allowed stale)\n", cacheKey)
		// r := remapping.Request
		// TODO verify setting the existing request header here works
		r.Header.Set("If-Modified-Since", cacheObj.RespRespTime.Format(time.RFC1123))
		oldCacheObj := cacheObj
		cacheObj, err = retryingGet(r, cacheObj)
		if err != nil {
			log.Errorf("retrying get error - serving stale as allowed: %v\n", err)
			cacheObj = oldCacheObj
		}
	}
	log.Debugf("cacheHandler.ServeHTTP: '%v' responding with %v\n", cacheKey, cacheObj.Code)
	h.respond(w, cacheObj.Code, cacheObj.RespHeaders, cacheObj.Body, h.stats, h.conns, r.RemoteAddr, reqFQDN, connectionClose)
}

// serveRuleNotFound writes the appropriate response to the client, via given writer, for when no remap rule was found for a request.
func (h *CacheHandler) serveRuleNotFound(w http.ResponseWriter) {
	code := http.StatusNotFound
	w.WriteHeader(code)
	w.Write([]byte(http.StatusText(code)))
}

// serveNotAllowed writes the appropriate response to the client, via given writer, for when the client's IP is not allowed for the requested rule.
func (h *CacheHandler) serveNotAllowed(w http.ResponseWriter) {
	code := http.StatusForbidden
	w.WriteHeader(code)
	w.Write([]byte(http.StatusText(code)))
}

// serveReqErr writes the appropriate response to the client, via given writer, for a generic request error.
func (h *CacheHandler) serveReqErr(w http.ResponseWriter) {
	code := http.StatusBadRequest
	w.WriteHeader(code)
	w.Write([]byte(http.StatusText(code)))
}

// func (h *cacheHandler) ServeCacheMiss(w http.ResponseWriter, r *http.Request, reqTime time.Time, key string) {
// 	fmt.Printf("DEBUG cacheHandler.ServeCacheMiss\n")
// 	reqHeader := http.Header{}
// 	copyHeader(r.Header, &reqHeader) // copy before ServeHTTP invalidates the request

// 	noCache := false // TODO fix
// 	h.ThrottleRequest(remapName, key, noCache, func() {
// 		respCode, respHeader, respBody, respReqTime, respRespTime, err = request(remappedReq)
// 	})
// 	if err != nil {
// 		fmt.Printf("DEBUG origin err for %v rule %v err %v\n", key, remapName, err)
// 		h.serveOriginErr(w)
// 		return
// 	}

// 	respHeader.Add("Requested-Host", remappedReq.Host)
// 	go h.respond(w, respCode, respHeader, respBody)

// 	h.TryCache(key, reqHeader, respBody, respCode, respHeader, respReqTime, respRespTime)
// }

// func (h *cacheHandler) ServeCacheHit(w http.ResponseWriter, r *http.Request, cacheObj cacheobj.CacheObj) {
// 	fmt.Printf("DEBUG cacheHandler.ServeCacheHit\n")
// 	h.respond(w, cacheObj.code, cacheObj.respHeaders, cacheObj.body)
// }

// func (h *cacheHandler) ServeCacheRevalidate(w http.ResponseWriter, r *http.Request, cacheObj cacheobj.CacheObj) {
// 	fmt.Printf("DEBUG cacheHandler.ServeCacheRevalidate\n")
// 	// TODO implement
// 	h.respond(w, cacheObj.code, cacheObj.respHeaders, cacheObj.body)
// 	// h.parent.ServeHTTP(w, r)
// }

// // TryCache determines if it can validly cache the given response per RFC 7234. If so, it caches it in this handler's cache.
// func (h *cacheHandler) TryCache(key string, reqHeader http.Header, bytes []byte, code int, respHeader http.Header, reqTime time.Time, respTime time.Time) {
// 	canCache := CanCache(reqHeader, code, respHeader)
// 	fmt.Printf("TryCache canCache '%v': %v\n", key, canCache)
// 	if canCache {
// 		obj := Newcacheobj.CacheObj(reqHeader, bytes, code, respHeader, reqTime, respTime)
// 		h.cache.AddSize(key, obj, obj.size)
// 	}
// }

// request makes the given request and returns its response code, headers, body, the request time, response time, and any error.
func request(transport *http.Transport, r *http.Request, proxyURL *url.URL) (int, http.Header, []byte, time.Time, time.Time, error) {
	log.Debugf("request requesting %v headers %v\n", r.RequestURI, r.Header)
	rr := r

	if proxyURL != nil && proxyURL.Host != "" {
		transport.Proxy = http.ProxyURL(proxyURL)
	}
	reqTime := time.Now()
	resp, err := transport.RoundTrip(rr)
	respTime := time.Now()
	if err != nil {
		return 0, nil, nil, reqTime, respTime, fmt.Errorf("request error: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	// TODO determine if respTime should go here

	if err != nil {
		return 0, nil, nil, reqTime, respTime, fmt.Errorf("reading response body: %v", err)
	}

	return resp.StatusCode, resp.Header, body, reqTime, respTime, nil
}

// respond writes the given code, header, and body to the ResponseWriter.
func (h *CacheHandler) respond(w http.ResponseWriter, code int, header http.Header, body []byte, stats Stats, conns *web.ConnMap, remoteAddr string, reqFQDN string, connectionClose bool) {
	dH := w.Header()
	web.CopyHeader(header, &dH)
	if connectionClose {
		dH.Add("Connection", "close")
	}
	w.WriteHeader(code)
	w.Write(body)

	remapRuleStats, ok := stats.Remap().Stats(reqFQDN)
	if !ok {
		log.Errorf("Remap rule %v not in Stats\n", reqFQDN)
		return
	}

	conn, ok := conns.Pop(remoteAddr)
	if !ok {
		log.Errorf("RemoteAddr %v not in Conns\n", remoteAddr)
		return
	}

	interceptConn, ok := conn.(*web.InterceptConn)
	if !ok {
		log.Errorf("Could not get Conn info: Conn is not an InterceptConn: %T\n", conn)
		return
	}

	if wFlusher, ok := w.(http.Flusher); !ok {
		log.Errorf("ResponseWriter is not a Flusher, could not flush written bytes, stat out_bytes will be inaccurate!\n")
	} else {
		wFlusher.Flush()
	}

	bytesRead := interceptConn.BytesRead()
	bytesWritten := interceptConn.BytesWritten()

	// bytesRead, bytesWritten := getConnInfoAndDestroyWriter(w, stats, remapRuleName)
	remapRuleStats.AddInBytes(uint64(bytesRead))
	remapRuleStats.AddOutBytes(uint64(bytesWritten))
	switch {
	case code < 200:
		log.Errorf("responded with invalid code %v\n", code)
	case code < 300:
		remapRuleStats.AddStatus2xx(1)
	case code < 400:
		remapRuleStats.AddStatus3xx(1)
	case code < 500:
		remapRuleStats.AddStatus4xx(1)
	case code < 600:
		remapRuleStats.AddStatus5xx(1)
	default:
		log.Errorf("responded with invalid code %v\n", code)
	}
}
