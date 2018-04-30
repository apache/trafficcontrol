package cache

import (
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/apache/incubator-trafficcontrol/grove/cacheobj"
	"github.com/apache/incubator-trafficcontrol/grove/icache"
	"github.com/apache/incubator-trafficcontrol/grove/remap"
	"github.com/apache/incubator-trafficcontrol/grove/thread"
	"github.com/apache/incubator-trafficcontrol/grove/web"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
)

const CodeConnectFailure = http.StatusBadGateway

type Retrier struct {
	H                 *Handler
	ReqHdr            http.Header
	ReqTime           time.Time
	ReqCacheControl   web.CacheControl
	RemappingProducer *remap.RemappingProducer
	ReqID             uint64
}

func NewRetrier(h *Handler, reqHdr http.Header, reqTime time.Time, reqCacheControl web.CacheControl, remappingProducer *remap.RemappingProducer, reqID uint64) *Retrier {
	return &Retrier{
		H:                 h,
		ReqHdr:            reqHdr,
		ReqCacheControl:   reqCacheControl,
		RemappingProducer: remappingProducer,
		ReqID:             reqID,
	}
}

// Get takes the HTTP request and the cached object if there is one, and makes a new request, retrying according to its RemappingProducer. If no cached object exists, pass a nil obj. Returns the cache object (which may be a failure), the parent URL used, and any error.
func (r *Retrier) Get(req *http.Request, obj *cacheobj.CacheObj) (*cacheobj.CacheObj, string, error) {
	retryGetFunc := func(remapping remap.Remapping, retryFailures bool, obj *cacheobj.CacheObj) *cacheobj.CacheObj {
		// return true for Revalidate, and issue revalidate requests separately.
		canReuse := func(cacheObj *cacheobj.CacheObj) bool {
			return remap.CanReuse(r.ReqHdr, r.ReqCacheControl, cacheObj, r.H.strictRFC, true)
		}
		getAndCache := func() *cacheobj.CacheObj {
			return GetAndCache(remapping.Request, remapping.ProxyURL, remapping.CacheKey, remapping.Name, remapping.Request.Header, r.ReqTime, r.H.strictRFC, remapping.Cache, r.H.ruleThrottlers[remapping.Name], obj, remapping.Timeout, retryFailures, remapping.RetryNum, remapping.RetryCodes, remapping.Transport, r.ReqID)
		}
		gotObj, getReqID := r.H.getter.Get(remapping.CacheKey, getAndCache, canReuse, r.ReqID)

		req := remapping.Request
		log.Debugf("Retrier.Get Y URI %v %v %v remapping.CacheKey %v rule %v parent %v code %v headers %+v len(body) %v getterid %v (reqid %v)\n", req.URL.Scheme, req.URL.Host, req.URL.EscapedPath(), remapping.CacheKey, remapping.Name, remapping.ProxyURL, gotObj.Code, gotObj.RespHeaders, len(gotObj.Body), getReqID, r.ReqID)

		return gotObj
	}

	return retryingGet(retryGetFunc, req, r.RemappingProducer, obj)
}

// retryingGet takes a function, and retries failures up to the RemappingProducer RetryNum limit. On failure, it creates a new remapping. The func f should use `remapping` to make its request. If it hits failures up to the limit, it returns the last received cacheobj.CacheObj
// Returns the cache object (which may be a failure), the last parent URL used, and any error.
// TODO refactor to not close variables - it's awkward and confusing.
func retryingGet(getCacheObj func(remapping remap.Remapping, retryFailures bool, obj *cacheobj.CacheObj) *cacheobj.CacheObj, request *http.Request, remappingProducer *remap.RemappingProducer, cachedObj *cacheobj.CacheObj) (*cacheobj.CacheObj, string, error) {
	obj := (*cacheobj.CacheObj)(nil)
	for {
		remapping, retryAllowed, err := remappingProducer.GetNext(request)
		parentURL := remapping.Request.URL.Host
		if err == remap.ErrNoMoreRetries {
			if obj == nil {
				return nil, parentURL, errors.New("remapping producer allows no requests") // should never happen
			}
			return obj, parentURL, nil
		} else if err != nil {
			return nil, parentURL, err
		}
		obj = getCacheObj(remapping, retryAllowed, cachedObj)
		if !isFailure(obj, remapping.RetryCodes) {
			return obj, parentURL, nil
		}
	}
}

func isFailure(o *cacheobj.CacheObj, retryCodes map[int]struct{}) bool {
	_, failureCode := retryCodes[o.Code]
	return failureCode || o.Code == CodeConnectFailure
}

const ModifiedSinceHdr = "If-Modified-Since"

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
	cache icache.Cache,
	ruleThrottler thread.Throttler,
	revalidateObj *cacheobj.CacheObj,
	timeout time.Duration,
	cacheFailure bool,
	retryNum int,
	retryCodes map[int]struct{},
	transport *http.Transport,
	reqID uint64,
) *cacheobj.CacheObj {
	// TODO this is awkward, with 'revalidateObj' indicating whether the request is a Revalidate. Should Getting and Caching be split up? How?
	get := func() *cacheobj.CacheObj {
		// TODO figure out why respReqTime isn't used by rules
		log.Debugf("GetAndCache calling request %v %v %v %v %v (reqid %v)\n", req.Method, req.URL.Scheme, req.URL.Host, req.URL.EscapedPath(), req.Header, reqID)
		// TODO Verify overriding the passed reqTime is the right thing to do
		proxyURLStr := ""
		if proxyURL != nil {
			proxyURLStr = proxyURL.Host
		}
		// else {
		// 	proxyURLStr = req.URL.Host
 		// }

		if revalidateObj != nil {
			req.Header.Set(ModifiedSinceHdr, revalidateObj.RespRespTime.Format(time.RFC1123))
		} else {
			req.Header.Del(ModifiedSinceHdr)
		}
		respCode, respHeader, respBody, reqTime, reqRespTime, err := web.Request(transport, req)
		log.Debugf("GetAndCache web.Request URI %v %v %v cacheKey %v rule %v parent %v error %v reval %v code %v len(body) %v (reqid %v)\n", req.URL.Scheme, req.URL.Host, req.URL.EscapedPath(), cacheKey, remapName, proxyURLStr, err, revalidateObj != nil, respCode, len(respBody), reqID)

		if err != nil {
			log.Errorf("Parent error for URI %v %v %v cacheKey %v rule %v parent %v error %v (reqid %v)\n", req.URL.Scheme, req.URL.Host, req.URL.EscapedPath(), cacheKey, remapName, proxyURLStr, err, reqID)
			code := CodeConnectFailure
			body := []byte(http.StatusText(code))
			return cacheobj.New(reqHeader, body, code, code, proxyURLStr, respHeader, reqTime, reqRespTime, reqRespTime, time.Time{})
		}
		if _, ok := retryCodes[respCode]; ok && !cacheFailure {
			return cacheobj.New(reqHeader, respBody, respCode, respCode, proxyURLStr, respHeader, reqTime, reqRespTime, reqRespTime, time.Time{})
		}

		log.Debugf("GetAndCache request returned %v headers %+v (reqid %v)\n", respCode, respHeader, reqID)
		respRespTime, ok := web.GetHTTPDate(respHeader, "Date")
		if !ok {
			log.Errorf("request %v returned no Date header - RFC Violation! Using local response timestamp (reqid %v)\n", req.RequestURI, reqID)
			respRespTime = reqRespTime // if no Date was returned using the client response time simulates latency 0
		}

		lastModified, ok := web.GetHTTPDate(respHeader, "Last-Modified")
		if !ok {
			lastModified = respRespTime
		}

		obj := (*cacheobj.CacheObj)(nil)
		log.Debugf("h.cache.Add %v (reqid %v)\n", cacheKey, reqID)
		log.Debugf("GetAndCache respCode %v (reqid %v)\n", respCode, reqID)
		if revalidateObj == nil || respCode != http.StatusNotModified {
			log.Debugf("GetAndCache new %v (reqid %v)\n", cacheKey, reqID)
			obj = cacheobj.New(reqHeader, respBody, respCode, respCode, proxyURLStr, respHeader, reqTime, reqRespTime, respRespTime, lastModified)
			if !remap.CanCache(reqHeader, respCode, respHeader, strictRFC) {
				return obj // return without caching
			}
		} else {
			log.Debugf("GetAndCache revalidating %v len(revalidateObj.Body) %v (reqid %v)\n", cacheKey, len(revalidateObj.Body), reqID)
			// must copy, because this cache object may be concurrently read by other goroutines
			newRespHeader := web.CopyHeader(revalidateObj.RespHeaders)
			newRespHeader.Set("Date", respHeader.Get("Date"))
			obj = &cacheobj.CacheObj{
				Body:             revalidateObj.Body,
				ReqHeaders:       revalidateObj.ReqHeaders,
				RespHeaders:      newRespHeader,
				RespCacheControl: revalidateObj.RespCacheControl,
				Code:             revalidateObj.Code,
				OriginCode:       respCode,
				ProxyURL:         proxyURLStr,
				ReqTime:          reqTime,
				ReqRespTime:      reqRespTime,
				RespRespTime:     respRespTime,
				LastModified:     revalidateObj.LastModified,
				Size:             revalidateObj.Size,
			}
		}
		cache.Add(cacheKey, obj) // TODO store pointer?
		return obj
	}

	c := (*cacheobj.CacheObj)(nil)
	if ruleThrottler == nil {
		log.Errorf("rule %v not in ruleThrottlers map. Requesting with no origin limit! (reqid %v)\n", remapName, reqID)
		ruleThrottler = thread.NewNoThrottler()
	}
	ruleThrottler.Throttle(func() { c = get() })
	return c
}
