package grove

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// TODO add logging

type Cache interface {
	AddSize(key string, value interface{}, size uint64) bool
	Get(key string) (interface{}, bool)
	Remove(key string)
	RemoveOldest()
	Size() uint64
}

type cacheHandler struct {
	cache          Cache
	remapper       HTTPRequestRemapper
	getter         Getter
	ruleThrottlers map[string]Throttler // doesn't need threadsafe keys, because it's never added to or deleted after creation. TODO fix for hot rule reloading
	scheme         string
	strictRFC      bool
	stats          Stats
	conns          *ConnMap
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
func NewCacheHandler(cache Cache, remapper HTTPRequestRemapper, ruleLimit uint64, stats Stats, scheme string, conns *ConnMap, strictRFC bool) http.Handler {
	return &cacheHandler{
		cache:          cache,
		remapper:       remapper,
		getter:         NewGetter(),
		ruleThrottlers: makeRuleThrottlers(remapper, ruleLimit),
		strictRFC:      strictRFC,
		scheme:         scheme,
		stats:          stats,
		conns:          conns,
		// keyThrottlers:     NewThrottlers(keyLimit),
		// nocacheThrottlers: NewThrottlers(nocacheLimit),
	}
}

func makeRuleThrottlers(remapper HTTPRequestRemapper, limit uint64) map[string]Throttler {
	remapRules := remapper.Rules()
	ruleThrottlers := make(map[string]Throttler, len(remapRules))
	for _, rule := range remapRules {
		ruleLimit := uint64(rule.ConcurrentRuleRequests)
		if rule.ConcurrentRuleRequests == 0 {
			ruleLimit = limit
		}
		ruleThrottlers[rule.Name] = NewThrottler(ruleLimit)
	}
	return ruleThrottlers
}

// NewCacheHandlerFunc creates and returns an http.HandleFunc, which may be pipelined with other http.HandleFuncs via `http.HandleFunc`. This is a convenience wrapper around the `http.Handler` object obtainable via `New`. If you prefer objects, use `NewCacheHandler`.
func NewCacheHandlerFunc(cache Cache, remapper HTTPRequestRemapper, ruleLimit uint64, stats Stats, scheme string, conns *ConnMap, strictRFC bool) http.HandlerFunc {
	handler := NewCacheHandler(cache, remapper, ruleLimit, stats, scheme, conns, strictRFC)
	return func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	}
}

func (h *cacheHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.TryServe(w, r)
}

// TryServe attempts to serve the given request, as a caching reverse proxy.
// Serving acts as a state machine.
func (h *cacheHandler) TryServe(w http.ResponseWriter, r *http.Request) {
	// inBytes := getBytes(r)
	reqTime := time.Now()

	// copy request header, because it's not guaranteed valid after actually issuing the request
	reqHeader := http.Header{}
	copyHeader(r.Header, &reqHeader)

	// TODO fix host header
	remappedReq, remapName, cacheKey, allowed, ok, err := h.remapper.Remap(r, h.scheme)
	if err != nil {
		fmt.Printf("DEBUG request error: %v\n", err)
		h.serveReqErr(w)
		return
	}
	if !ok {
		fmt.Printf("DEBUG rule not found for %v\n", r.RequestURI)
		h.serveRuleNotFound(w)
		return
	}
	if !allowed {
		fmt.Printf("DEBUG IP %v not allowed\n", r.RemoteAddr)
		h.serveNotAllowed(w)
		return
	}

	getAndCache := func() *CacheObj {
		get := func() *CacheObj {
			// TODO figure out why respReqTime isn't used
			respCode, respHeader, respBody, _, respRespTime, err := h.request(remappedReq)
			if err != nil {
				fmt.Printf("DEBUG origin err for %v rule %v err %v\n", cacheKey, remapName, err)
				code := http.StatusInternalServerError
				body := []byte(http.StatusText(code))
				return NewCacheObj(reqHeader, body, code, respHeader, reqTime, respRespTime)
			}
			obj := NewCacheObj(reqHeader, respBody, respCode, respHeader, reqTime, respRespTime)
			if CanCache(reqHeader, respCode, respHeader, h.strictRFC) {
				fmt.Printf("h.cache.AddSize %v\n", cacheKey)
				h.cache.AddSize(cacheKey, obj, obj.size) // TODO store pointer?
			}
			return obj
		}

		c := (*CacheObj)(nil)
		ruleThrottler, ok := h.ruleThrottlers[remapName]
		if !ok {
			fmt.Printf("ERROR rule %v returned, but not in ruleThrottlers map. Requesting with no rule (origin) limit!\n", remapName)
			ruleThrottler = NewNoThrottler()
		}
		ruleThrottler.Throttle(func() {
			c = get()
		})
		return c
	}

	reqCacheControl := ParseCacheControl(reqHeader)
	fmt.Printf("DEBUG TryServe got Cache-Control %+v\n", reqCacheControl)
	// return true for Revalidate, and issue revalidate requests separately.
	canReuse := func(cacheObj *CacheObj) bool {
		canReuse := CanReuseStored(reqHeader, cacheObj.respHeaders, reqCacheControl, cacheObj.respCacheControl, cacheObj.reqHeaders, cacheObj.reqTime, cacheObj.respTime, h.strictRFC)
		switch canReuse {
		case ReuseCan:
			return true
		case ReuseCannot:
			return false
		case ReuseMustRevalidate:
			return true
		default:
			fmt.Printf("Error: CanReuseStored returned unknown %v\n", canReuse)
			return false
		}
	}

	iCacheObj, ok := h.cache.Get(cacheKey)
	if !ok {
		fmt.Printf("DEBUG cacheHandler.ServeHTTP: '%v' not in cache\n", cacheKey)
		cacheObj := h.getter.Get(cacheKey, getAndCache, canReuse)
		h.respond(w, cacheObj.code, cacheObj.respHeaders, cacheObj.body, h.stats, h.conns, r.RemoteAddr, remapName)
		return
	}

	cacheObj, ok := iCacheObj.(*CacheObj)
	if !ok {
		// should never happen
		fmt.Printf("Error: cache key '%v' value '%v' type '%T' expected *CacheObj\n", cacheKey, iCacheObj, iCacheObj)
		cacheObj = h.getter.Get(cacheKey, getAndCache, canReuse)
		// TODO check for ReuseMustRevalidate
		h.respond(w, cacheObj.code, cacheObj.respHeaders, cacheObj.body, h.stats, h.conns, r.RemoteAddr, remapName)
		return
	}

	reqHeaders := r.Header

	canReuseStored := CanReuseStored(reqHeaders, cacheObj.respHeaders, reqCacheControl, cacheObj.respCacheControl, cacheObj.reqHeaders, cacheObj.reqTime, cacheObj.respTime, h.strictRFC)

	switch canReuseStored {
	case ReuseCan:
		fmt.Printf("DEBUG cacheHandler.ServeHTTP: '%v' cache hit!\n", cacheKey)
	case ReuseCannot:
		fmt.Printf("DEBUG cacheHandler.ServeHTTP: '%v' can't reuse\n", cacheKey)
		cacheObj = h.getter.Get(cacheKey, getAndCache, canReuse)
	case ReuseMustRevalidate:
		fmt.Printf("DEBUG cacheHandler.ServeHTTP: '%v' must revalidate\n", cacheKey)
		// TODO implement revalidate
		cacheObj = h.getter.Get(cacheKey, getAndCache, canReuse)
	}
	h.respond(w, cacheObj.code, cacheObj.respHeaders, cacheObj.body, h.stats, h.conns, r.RemoteAddr, remapName)
}

// serveRuleNotFound writes the appropriate response to the client, via given writer, for when no remap rule was found for a request.
func (h *cacheHandler) serveRuleNotFound(w http.ResponseWriter) {
	code := http.StatusNotFound
	w.WriteHeader(code)
	w.Write([]byte(http.StatusText(code)))
}

// serveNotAllowed writes the appropriate response to the client, via given writer, for when the client's IP is not allowed for the requested rule.
func (h *cacheHandler) serveNotAllowed(w http.ResponseWriter) {
	code := http.StatusForbidden
	w.WriteHeader(code)
	w.Write([]byte(http.StatusText(code)))
}

// serveReqErr writes the appropriate response to the client, via given writer, for a generic request error.
func (h *cacheHandler) serveReqErr(w http.ResponseWriter) {
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
// 		respCode, respHeader, respBody, respReqTime, respRespTime, err = h.request(remappedReq)
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

// func (h *cacheHandler) ServeCacheHit(w http.ResponseWriter, r *http.Request, cacheObj CacheObj) {
// 	fmt.Printf("DEBUG cacheHandler.ServeCacheHit\n")
// 	h.respond(w, cacheObj.code, cacheObj.respHeaders, cacheObj.body)
// }

// func (h *cacheHandler) ServeCacheRevalidate(w http.ResponseWriter, r *http.Request, cacheObj CacheObj) {
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
// 		obj := NewCacheObj(reqHeader, bytes, code, respHeader, reqTime, respTime)
// 		h.cache.AddSize(key, obj, obj.size)
// 	}
// }

// request makes the given request and returns its response code, headers, body, the request time, response time, and any error.
func (h *cacheHandler) request(r *http.Request) (int, http.Header, []byte, time.Time, time.Time, error) {
	rr := r
	// Create a client and query the target
	var transport http.Transport
	reqTime := time.Now()
	resp, err := transport.RoundTrip(rr)
	respTime := time.Now()
	if err != nil {
		return 0, nil, nil, reqTime, respTime, fmt.Errorf("request error: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, nil, reqTime, respTime, fmt.Errorf("reading response body: %v", err)
	}

	return resp.StatusCode, resp.Header, body, reqTime, respTime, nil
}

// respond writes the given code, header, and body to the ResponseWriter.
func (h *cacheHandler) respond(w http.ResponseWriter, code int, header http.Header, body []byte, stats Stats, conns *ConnMap, remoteAddr string, remapRuleName string) {
	dH := w.Header()
	copyHeader(header, &dH)
	w.WriteHeader(code)
	w.Write(body)

	remapRuleStats, ok := stats.Remap().Stats(remapRuleName)
	if !ok {
		fmt.Printf("ERROR Remap rule %v not in Stats\n", remapRuleName)
		return
	}

	conn, ok := conns.Pop(remoteAddr)
	if !ok {
		fmt.Printf("ERROR RemoteAddr %v not in Conns\n", remoteAddr)
		return
	}

	interceptConn, ok := conn.(*InterceptConn)
	if !ok {
		fmt.Printf("ERROR Could not get Conn info: Conn is not an InterceptConn: %T\n", conn)
		return
	}

	if wFlusher, ok := w.(http.Flusher); !ok {
		fmt.Printf("ERROR ResponseWriter is not a Flusher, could not flush written bytes, stat out_bytes will be inaccurate!\n")
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
		fmt.Printf("ERROR responded with invalid code %v\n", code)
	case code < 300:
		remapRuleStats.AddStatus2xx(1)
	case code < 400:
		remapRuleStats.AddStatus3xx(1)
	case code < 500:
		remapRuleStats.AddStatus4xx(1)
	case code < 600:
		remapRuleStats.AddStatus5xx(1)
	default:
		fmt.Printf("ERROR responded with invalid code %v\n", code)
	}
}
