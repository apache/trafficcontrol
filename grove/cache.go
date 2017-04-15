package grove

import (
	"fmt"
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
	cache  Cache
	parent http.Handler
}

// NewHandler returns an http.Handler objectn, which may be pipelined with other http.Handlers via `http.ListenAndServe`. If you prefer pipelining functions, use `GetHandlerFunc`.
func NewCacheHandler(parent http.Handler, cache Cache) http.Handler {
	return &cacheHandler{
		cache:  cache,
		parent: parent,
	}
}

// NewHandlerFunc creates and returns an http.HandleFunc, which may be pipelined with other http.HandleFuncs via `http.HandleFunc`. This is a convenience wrapper around the `http.Handler` object obtainable via `New`. If you prefer objects, use Java. I mean, `New`.
func NewCacheHandlerFunc(parent http.HandlerFunc, cache Cache) http.HandlerFunc {
	handler := NewCacheHandler(parent, cache)
	return func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	}
}

func (h *cacheHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqTime := time.Now()
	key := buildKey(r) // key MUST be built before calling ServeHTTP, because the http.Request is not guaranteed valid after.
	iCacheObj, ok := h.cache.Get(key)
	if !ok {
		fmt.Printf("DEBUG cacheHandler.ServeHTTP: '%v' not in cache\n", key)
		h.ServeCacheMiss(w, r, reqTime, key)
		return
	}
	cacheObj, ok := iCacheObj.(CacheObj)
	if !ok {
		// should never happen
		fmt.Printf("Error: cache key '%v' value '%v' type '%T' expected CacheObj\n", key, iCacheObj, iCacheObj)
		h.ServeCacheMiss(w, r, reqTime, key)
		return
	}

	reqHeaders := r.Header
	reqCacheControl := ParseCacheControl(reqHeaders)

	canReuseStored := CanReuseStored(reqHeaders, cacheObj.respHeaders, reqCacheControl, cacheObj.respCacheControl, cacheObj.reqHeaders, cacheObj.reqTime, cacheObj.respTime)

	switch canReuseStored {
	case ReuseCan:
		fmt.Printf("DEBUG cacheHandler.ServeHTTP: '%v' cache hit!\n", key)
		h.ServeCacheHit(w, r, cacheObj)
	case ReuseCannot:
		fmt.Printf("DEBUG cacheHandler.ServeHTTP: '%v' can't reuse\n", key)
		h.ServeCacheMiss(w, r, reqTime, key)
	case ReuseMustRevalidate:
		fmt.Printf("DEBUG cacheHandler.ServeHTTP: '%v' must revalidate\n", key)
		h.ServeCacheRevalidate(w, r, cacheObj)
	}
}

func (h *cacheHandler) ServeCacheMiss(w http.ResponseWriter, r *http.Request, reqTime time.Time, key string) {
	fmt.Printf("DEBUG cacheHandler.ServeCacheMiss\n")
	teeWriter := NewHTTPResponseWriterTee(w)
	reqHeader := http.Header{}
	copyHeader(r.Header, &reqHeader) // copy before ServeHTTP invalidates the request

	h.parent.ServeHTTP(teeWriter, r)
	respTime := time.Now() // TODO get response time as soon as it's returned. This is used to estimate latency per RFC 7234, it needs to be _immediately_ after the origin responds.
	h.TryCache(key, reqHeader, teeWriter.Bytes, teeWriter.Code, teeWriter.WrittenHeader, reqTime, respTime)
}

func (h *cacheHandler) ServeCacheHit(w http.ResponseWriter, r *http.Request, cacheObj CacheObj) {
	fmt.Printf("DEBUG cacheHandler.ServeCacheHit\n")
	h.parent.ServeHTTP(w, r)
}

func (h *cacheHandler) ServeCacheRevalidate(w http.ResponseWriter, r *http.Request, cacheObj CacheObj) {
	fmt.Printf("DEBUG cacheHandler.ServeCacheRevalidate\n")
	h.parent.ServeHTTP(w, r)
}

// TODO make configurable (method, uri, query params, etc)
func buildKey(r *http.Request) string {
	uri := fmt.Sprintf("%s://%s%s", getScheme(r), r.Host, r.RequestURI)
	key := fmt.Sprintf("%s:%s", r.Method, uri)
	return key
}

// TryCache determines if it can validly cache the given response per RFC 7234. If so, it caches it in this handler's cache.
func (h *cacheHandler) TryCache(key string, reqHeader http.Header, bytes []byte, code int, respHeader http.Header, reqTime time.Time, respTime time.Time) {
	canCache := CanCache(reqHeader, code, respHeader)
	fmt.Printf("TryCache canCache '%v': %v\n", key, canCache)
	if canCache {
		obj := NewCacheObj(reqHeader, bytes, code, respHeader, reqTime, respTime)
		h.cache.AddSize(key, obj, obj.size)
	}
}
