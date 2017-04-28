package grove

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// TODO add logging

type Cache interface {
	AddSize(key, value interface{}, size uint64) bool
	Get(key interface{}) (interface{}, bool)
	Remove(key interface{})
	RemoveOldest()
	Size() uint64
}

type cacheHandler struct {
	cache    Cache
	remapper HTTPRequestRemapper
	getter   Getter
	// ruleThrottlers    map[string]Throttler // doesn't need threadsafe keys, because it's never added to or deleted after creation.
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
func NewCacheHandler(cache Cache, remapper HTTPRequestRemapper) http.Handler {
	return &cacheHandler{
		cache:    cache,
		remapper: remapper,
		getter:   NewGetter(),
		// ruleThrottlers:    makeRuleThrottlers(remapper, ruleLimit),
		// keyThrottlers:     NewThrottlers(keyLimit),
		// nocacheThrottlers: NewThrottlers(nocacheLimit),
	}
}

// func makeRuleThrottlers(remapper HTTPRequestRemapper, limit uint64) map[string]Throttler {
// 	remapRules := remapper.Rules()
// 	ruleThrottlers := make(map[string]Throttler, len(remapRules))
// 	for _, rule := range remapRules {
// 		ruleThrottlers[rule] = NewThrottler(limit)
// 	}
// 	return ruleThrottlers
// }

// NewHandlerFunc creates and returns an http.HandleFunc, which may be pipelined with other http.HandleFuncs via `http.HandleFunc`. This is a convenience wrapper around the `http.Handler` object obtainable via `New`. If you prefer objects, use Java. I mean, `New`.
func NewCacheHandlerFunc(cache Cache, remapper HTTPRequestRemapper) http.HandlerFunc {
	handler := NewCacheHandler(cache, remapper)
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
	reqTime := time.Now()

	// copy request header, because it's not guaranteed valid after actually issuing the request
	reqHeader := http.Header{}
	copyHeader(r.Header, &reqHeader)

	key := buildKey(r) // key MUST be built before calling ServeHTTP, because the http.Request is not guaranteed valid after.

	remappedReq, remapName, ruleFound := h.remapper.Remap(r)
	if !ruleFound {
		fmt.Printf("DEBUG rule not found for %v\n", key)
		h.serveRuleNotFound(w)
		return
	}

	getAndCache := func() *CacheObj {
		// TODO figure out why respReqTime isn't used
		respCode, respHeader, respBody, _, respRespTime, err := h.request(remappedReq)
		if err != nil {
			fmt.Printf("DEBUG origin err for %v rule %v err %v\n", key, remapName, err)
			code := http.StatusInternalServerError
			body := []byte(http.StatusText(code))
			return NewCacheObj(reqHeader, body, code, respHeader, reqTime, respRespTime)
		}
		obj := NewCacheObj(reqHeader, respBody, respCode, respHeader, reqTime, respRespTime)
		if CanCache(reqHeader, respCode, respHeader) {
			fmt.Printf("h.cache.AddSize %v\n", key)
			h.cache.AddSize(key, obj, obj.size) // TODO store pointer?
		}
		return obj
	}

	reqCacheControl := ParseCacheControl(reqHeader)
	fmt.Printf("DEBUG TryServe got Cache-Control %+v\n", reqCacheControl)
	// return true for Revalidate, and issue revalidate requests separately.
	canReuse := func(cacheObj *CacheObj) bool {
		canReuse := CanReuseStored(reqHeader, cacheObj.respHeaders, reqCacheControl, cacheObj.respCacheControl, cacheObj.reqHeaders, cacheObj.reqTime, cacheObj.respTime)
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

	iCacheObj, ok := h.cache.Get(key)
	if !ok {
		fmt.Printf("DEBUG cacheHandler.ServeHTTP: '%v' not in cache\n", key)
		cacheObj := h.getter.Get(key, getAndCache, canReuse)
		h.respond(w, cacheObj.code, cacheObj.respHeaders, cacheObj.body)
		return
	}

	cacheObj, ok := iCacheObj.(*CacheObj)
	if !ok {
		// should never happen
		fmt.Printf("Error: cache key '%v' value '%v' type '%T' expected *CacheObj\n", key, iCacheObj, iCacheObj)
		cacheObj = h.getter.Get(key, getAndCache, canReuse)
		// TODO check for ReuseMustRevalidate
		h.respond(w, cacheObj.code, cacheObj.respHeaders, cacheObj.body)
		return
	}

	reqHeaders := r.Header

	canReuseStored := CanReuseStored(reqHeaders, cacheObj.respHeaders, reqCacheControl, cacheObj.respCacheControl, cacheObj.reqHeaders, cacheObj.reqTime, cacheObj.respTime)

	switch canReuseStored {
	case ReuseCan:
		fmt.Printf("DEBUG cacheHandler.ServeHTTP: '%v' cache hit!\n", key)
	case ReuseCannot:
		fmt.Printf("DEBUG cacheHandler.ServeHTTP: '%v' can't reuse\n", key)
		cacheObj = h.getter.Get(key, getAndCache, canReuse)
	case ReuseMustRevalidate:
		fmt.Printf("DEBUG cacheHandler.ServeHTTP: '%v' must revalidate\n", key)
		// TODO implement revalidate
		cacheObj = h.getter.Get(key, getAndCache, canReuse)
	}
	h.respond(w, cacheObj.code, cacheObj.respHeaders, cacheObj.body)
}

// serveRuleNotFound writes the appropriate response to the client, via given writer, for when no remap rule was found for a request.
func (h *cacheHandler) serveRuleNotFound(w http.ResponseWriter) {
	code := http.StatusNotFound
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

// TODO make configurable (method, uri, query params, etc)
func buildKey(r *http.Request) string {
	uri := fmt.Sprintf("%s://%s%s", getScheme(r), r.Host, r.RequestURI)
	key := fmt.Sprintf("%s:%s", r.Method, uri)
	return key
}

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

func (h *cacheHandler) respond(w http.ResponseWriter, code int, header http.Header, body []byte) {
	dH := w.Header()
	copyHeader(header, &dH)
	w.WriteHeader(code)
	w.Write(body)
}

// func (h *cacheHandler) ThrottleRequest(remapRuleName string, remapKey string, noCache bool, request func()) {
// 	h.ruleThrottlers[remapRuleName].Throttle(func() {
// 		var throttlers Throttlers
// 		if noCache {
// 			throttlers = h.nocacheThrottlers
// 		} else {
// 			throttlers = h.keyThrottlers
// 		}
// 		throttlers.Throttle(remapKey, request)
// 	})
// }
