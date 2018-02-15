package cache

import (
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/apache/incubator-trafficcontrol/grove/cacheobj"
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
	CacheKey          string
	RemappingProducer *remap.RemappingProducer
}

func NewRetrier(h *Handler, reqHdr http.Header, reqTime time.Time, reqCacheControl web.CacheControl, cacheKey string, remappingProducer *remap.RemappingProducer) *Retrier {
	return &Retrier{
		H:                 h,
		ReqHdr:            reqHdr,
		ReqCacheControl:   reqCacheControl,
		RemappingProducer: remappingProducer,
	}
}

// Get takes the HTTP request and the cached object if there is one, and makes a new request, retrying according to its RemappingProducer. If no cached object exists, pass a nil obj.
func (r *Retrier) Get(req *http.Request, obj *cacheobj.CacheObj) (*cacheobj.CacheObj, error) {
	retryGetFunc := func(remapping remap.Remapping, retryFailures bool, obj *cacheobj.CacheObj) *cacheobj.CacheObj {
		// return true for Revalidate, and issue revalidate requests separately.
		canReuse := func(cacheObj *cacheobj.CacheObj) bool {
			return remap.CanReuse(r.ReqHdr, r.ReqCacheControl, cacheObj, r.H.strictRFC, true)
		}
		getAndCache := func() *cacheobj.CacheObj {
			return GetAndCache(remapping.Request, remapping.ProxyURL, remapping.CacheKey, remapping.Name, remapping.Request.Header, r.ReqTime, r.H.strictRFC, r.H.cache, r.H.ruleThrottlers[remapping.Name], obj, remapping.Timeout, retryFailures, remapping.RetryNum, remapping.RetryCodes, r.H.transport)
		}
		return r.H.getter.Get(r.CacheKey, getAndCache, canReuse)
	}

	return retryingGet(retryGetFunc, req, r.RemappingProducer, obj)
}

// retryingGet takes a function, and retries failures up to the RemappingProducer RetryNum limit. On failure, it creates a new remapping. The func f should use `remapping` to make its request. If it hits failures up to the limit, it returns the last received cacheobj.CacheObj
// TODO refactor to not close variables - it's awkward and confusing.
func retryingGet(getCacheObj func(remapping remap.Remapping, retryFailures bool, obj *cacheobj.CacheObj) *cacheobj.CacheObj, request *http.Request, remappingProducer *remap.RemappingProducer, cachedObj *cacheobj.CacheObj) (*cacheobj.CacheObj, error) {
	obj := (*cacheobj.CacheObj)(nil)
	for {
		remapping, retryAllowed, err := remappingProducer.GetNext(request)
		if err == remap.ErrNoMoreRetries {
			if obj == nil {
				return nil, errors.New("remapping producer allows no requests") // should never happen
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
		proxyURLStr := ""
		if proxyURL != nil {
			proxyURLStr = proxyURL.Host
		}
		if revalidateObj != nil {
			req.Header.Set(ModifiedSinceHdr, revalidateObj.RespRespTime.Format(time.RFC1123))
		} else {
			req.Header.Del(ModifiedSinceHdr)
		}
		respCode, respHeader, respBody, reqTime, reqRespTime, err := web.Request(transport, req, proxyURL)
		if err != nil {
			log.Errorf("Parent error for URI %v %v %v cacheKey %v rule %v parent %v error %v\n", req.URL.Scheme, req.URL.Host, req.URL.EscapedPath(), cacheKey, remapName, proxyURLStr, err)
			code := CodeConnectFailure
			body := []byte(http.StatusText(code))
			return cacheobj.New(reqHeader, body, code, code, proxyURLStr, respHeader, reqTime, reqRespTime, reqRespTime, time.Time{})
		}
		if _, ok := retryCodes[respCode]; ok && !cacheFailure {
			return cacheobj.New(reqHeader, respBody, respCode, respCode, proxyURLStr, respHeader, reqTime, reqRespTime, reqRespTime, time.Time{})
		}

		log.Debugf("GetAndCache request returned %v headers %+v\n", respCode, respHeader)
		respRespTime, ok := web.GetHTTPDate(respHeader, "Date")
		if !ok {
			log.Errorf("request %v returned no Date header - RFC Violation! Using local response timestamp.\n", req.RequestURI)
			respRespTime = reqRespTime // if no Date was returned using the client response time simulates latency 0
		}

		lastModified, ok := web.GetHTTPDate(respHeader, "Last-Modified")
		if !ok {
			lastModified = respRespTime
		}

		obj := (*cacheobj.CacheObj)(nil)
		log.Debugf("h.cache.AddSize %v\n", cacheKey)
		log.Debugf("GetAndCache respCode %v\n", respCode)
		if revalidateObj == nil || respCode != http.StatusNotModified {
			log.Debugf("GetAndCache new %v\n", cacheKey)
			obj = cacheobj.New(reqHeader, respBody, respCode, respCode, proxyURLStr, respHeader, reqTime, reqRespTime, respRespTime, lastModified)
			if !remap.CanCache(reqHeader, respCode, respHeader, strictRFC) {
				return obj // return without caching
			}
		} else {
			log.Debugf("GetAndCache revalidating %v\n", cacheKey)
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
