package cache

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/apache/incubator-trafficcontrol/grove/cacheobj"
	"github.com/apache/incubator-trafficcontrol/grove/thread"
	"github.com/apache/incubator-trafficcontrol/grove/web"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"

	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

// TODO add logging

type Cache interface {
	AddSize(key string, value interface{}, size uint64) bool
	Get(key string) (interface{}, bool)
	Remove(key string)
	RemoveOldest()
	Size() uint64
}

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
	cache           Cache
	remapper        HTTPRequestRemapper
	getter          thread.Getter
	ruleThrottlers  map[string]thread.Throttler // doesn't need threadsafe keys, because it's never added to or deleted after creation. TODO fix for hot rule reloading
	scheme          string
	port            string
	hostname        string
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
func NewHandler(
	cache Cache,
	remapper HTTPRequestRemapper,
	ruleLimit uint64,
	stats Stats,
	scheme string,
	port string,
	conns *web.ConnMap,
	strictRFC bool,
	connectionClose bool,
	reqTimeout time.Duration,
	reqKeepAlive time.Duration,
	reqMaxIdleConns int,
	reqIdleConnTimeout time.Duration,
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
		cache:           cache,
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

const CodeConnectFailure = http.StatusBadGateway
const NSPerSec = 1000000000

// NewHandlerFunc creates and returns an http.HandleFunc, which may be pipelined with other http.HandleFuncs via `http.HandleFunc`. This is a convenience wrapper around the `http.Handler` object obtainable via `New`. If you prefer objects, use `NewCacheHandler`.
func NewHandlerFunc(
	cache Cache,
	remapper HTTPRequestRemapper,
	ruleLimit uint64,
	stats Stats,
	scheme string,
	port string,
	conns *web.ConnMap,
	strictRFC bool,
	connectionClose bool,
	reqTimeout time.Duration,
	reqKeepAlive time.Duration,
	reqMaxIdleConns int,
	reqIdleConnTimeout time.Duration,
) http.HandlerFunc {
	handler := NewHandler(
		cache,
		remapper,
		ruleLimit,
		stats,
		scheme,
		port,
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

func setDSCP(conn *web.InterceptConn, dscp int) error {
	if dscp == 0 {
		return nil
	}
	if conn == nil {
		return errors.New("Conn is nil")
	}
	realConn := conn.Real()
	if realConn == nil {
		return errors.New("real Conn is nil")
	}
	ipv4Err := ipv4.NewConn(realConn).SetTOS(dscp)
	ipv6Err := ipv6.NewConn(realConn).SetTrafficClass(dscp)
	if ipv4Err != nil || ipv6Err != nil {
		errStr := ""
		if ipv4Err != nil {
			errStr = "setting IPv4 TOS: " + ipv4Err.Error()
		}
		if ipv6Err != nil {
			if ipv4Err != nil {
				errStr += "; "
			}
			errStr += "setting IPv6 TrafficClass: " + ipv6Err.Error()
		}
		return errors.New(errStr)
	}
	return nil
}

// modHdrs drops and sets headers in h according to the input drop and set lists
func modHdrs(h *http.Header, drop []string, set []Hdr) {
	if h == nil || len(*h) == 0 { // this happens on a dial tcp timeout
		log.Debugf("modHdrs: Header is  a nil map")
		return
	}
	for _, hdr := range drop {
		log.Debugf("modHdrs: Dropping header %s\n", hdr)
		h.Del(hdr)
	}
	for _, hdr := range set {
		log.Debugf("modHdrs: Setting header %s: %s \n", hdr.Name, hdr.Value)
		h.Set(hdr.Name, hdr.Value)
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.stats.IncConnections()
	defer h.stats.DecConnections()

	conn := (*web.InterceptConn)(nil)
	if realConn, ok := h.conns.Pop(r.RemoteAddr); !ok {
		log.Errorf("RemoteAddr '%v' not in Conns\n", r.RemoteAddr)
	} else {
		if conn, ok = realConn.(*web.InterceptConn); !ok {
			log.Errorf("Could not get Conn info: Conn is not an InterceptConn: %T\n", realConn)
		}
	}

	remappingProducer, err := h.remapper.RemappingProducer(r, h.scheme)

	if err == nil { // if we failed to get a remapping, there's no DSCP to set.
		if err := setDSCP(conn, remappingProducer.DSCP()); err != nil {
			log.Errorln(time.Now().Format(time.RFC3339Nano) + " " + r.RemoteAddr + " " + r.Method + " " + r.RequestURI + ": could not set DSCP: " + err.Error())
		}
	}

	reqTime := time.Now()
	reqHeader := web.CopyHeader(r.Header) // copy request header, because it's not guaranteed valid after actually issuing the request
	moneyTraceHdr := reqHeader.Get("X-Money-Trace")
	clientIP, _ := GetClientIPPort(r)
	statLog := NewStatLogger(w, conn, h, r, moneyTraceHdr, clientIP, reqTime, remappingProducer)

	if err != nil {
		code := 0
		if err == ErrRuleNotFound {
			log.Debugf("rule not found for %v\n", r.RequestURI)
			code = http.StatusNotFound
		} else if err == ErrIPNotAllowed {
			log.Debugf("IP %v not allowed\n", r.RemoteAddr)
			code = http.StatusForbidden
		} else {
			log.Debugf("request error: %v\n", err)
			code = http.StatusBadRequest
		}
		bytesWritten := uint64(0)
		code, bytesWritten, err = serveErr(w, code)
		tryFlush(w)
		statLog.Log(code, bytesWritten, err == nil, false, isCacheHit(ReuseCannot, 0), true, 0, 0, "-")
		return
	}

	reqCacheControl := web.ParseCacheControl(reqHeader)
	log.Debugf("Serve got Cache-Control %+v\n", reqCacheControl)

	connectionClose := h.connectionClose || remappingProducer.ConnectionClose()
	cacheKey := remappingProducer.CacheKey()
	retrier := NewRetrier(h, reqHeader, reqTime, reqCacheControl, cacheKey, remappingProducer)
	iCacheObj, ok := h.cache.Get(cacheKey)

	if !ok {
		log.Debugf("cache.Handler.ServeHTTP: '%v' not in cache\n", cacheKey)
		modHdrs(&r.Header, remappingProducer.rule.ToOriginHeaders.Drop, remappingProducer.rule.ToOriginHeaders.Set)
		cacheObj, err := retrier.Get(r, nil)
		if err != nil {
			log.Errorf("retrying get error (in uncached): %v\n", err)

			code, bytesWritten, err := serveReqErr(w)
			if err != nil {
				log.Errorln(time.Now().Format(time.RFC3339Nano) + " " + r.RemoteAddr + " " + r.Method + " " + r.RequestURI + ": responding: " + err.Error())
			}
			tryFlush(w)
			statLog.Log(code, bytesWritten, err == nil, false, isCacheHit(ReuseCannot, 0), true, 0, 0, cacheObj.ProxyURL)
			return
		}
		modHdrs(&cacheObj.RespHeaders, remappingProducer.rule.ToClientHeaders.Drop, remappingProducer.rule.ToClientHeaders.Set)
		bytesWritten, err := h.respond(w, cacheObj.Code, cacheObj.RespHeaders, cacheObj.Body, connectionClose)
		if err != nil {
			log.Errorln(time.Now().Format(time.RFC3339Nano) + " " + r.RemoteAddr + " " + r.Method + " " + r.RequestURI + ": responding: " + err.Error())
		}
		tryFlush(w)
		statLog.Log(cacheObj.Code, bytesWritten, true, err == nil, isCacheHit(ReuseCannot, cacheObj.OriginCode), false, 0, 0, cacheObj.ProxyURL)
		return
	}

	cacheObj, ok := iCacheObj.(*cacheobj.CacheObj)
	if !ok {
		// should never happen
		log.Errorf("cache key '%v' value '%v' type '%T' expected *cacheobj.CacheObj\n", cacheKey, iCacheObj, iCacheObj)
		cacheObj, err = retrier.Get(r, nil)
		if err != nil {
			log.Errorf("retrying get error (in unexpected cacheobj): %v\n", err)
			code, bytesWritten, err := serveReqErr(w)
			tryFlush(w)
			statLog.Log(code, bytesWritten, err == nil, false, isCacheHit(ReuseCannot, 0), false, 0, 0, "-")
			return
		}
		bytesWritten, err := h.respond(w, cacheObj.Code, cacheObj.RespHeaders, cacheObj.Body, connectionClose)
		if err != nil {
			log.Errorln(time.Now().Format(time.RFC3339Nano) + " " + r.RemoteAddr + " " + r.Method + " " + r.RequestURI + ": responding: " + err.Error())
		}
		tryFlush(w)
		statLog.Log(cacheObj.Code, bytesWritten, err == nil, true, isCacheHit(ReuseCannot, cacheObj.OriginCode), false, cacheObj.OriginCode, cacheObj.Size, cacheObj.ProxyURL)
		return
	}

	reqHeaders := r.Header
	canReuseStored := CanReuseStored(reqHeaders, cacheObj.RespHeaders, reqCacheControl, cacheObj.RespCacheControl, cacheObj.ReqHeaders, cacheObj.ReqRespTime, cacheObj.RespRespTime, h.strictRFC)

	switch canReuseStored {
	case ReuseCan:
		log.Debugf("cache.Handler.ServeHTTP: '%v' cache hit!\n", cacheKey)
	case ReuseCannot:
		log.Debugf("cache.Handler.ServeHTTP: '%v' can't reuse\n", cacheKey)
		cacheObj, err = retrier.Get(r, nil)
		if err != nil {
			log.Errorf("retrying get error (in reuse-cannot): %v\n", err)
			code, bytesWritten, err := serveReqErr(w)
			if err != nil {
				log.Errorln(time.Now().Format(time.RFC3339Nano) + " " + r.RemoteAddr + " " + r.Method + " " + r.RequestURI + ": responding: " + err.Error())
			}
			tryFlush(w)
			statLog.Log(code, bytesWritten, err == nil, false, isCacheHit(ReuseCannot, 0), false, 0, 0, "-")
			return
		}
	case ReuseMustRevalidate:
		log.Debugf("cache.Handler.ServeHTTP: '%v' must revalidate\n", cacheKey)
		r.Header.Set("If-Modified-Since", cacheObj.RespRespTime.Format(time.RFC1123))
		cacheObj, err = retrier.Get(r, cacheObj)
		if err != nil {
			log.Errorf("retrying get error: %v\n", err)
			code, bytesWritten, err := serveReqErr(w)
			if err != nil {
				log.Errorln(time.Now().Format(time.RFC3339Nano) + " " + r.RemoteAddr + " " + r.Method + " " + r.RequestURI + ": responding: " + err.Error())
			}
			tryFlush(w)
			statLog.Log(code, bytesWritten, err == nil, false, isCacheHit(ReuseCannot, code), false, 0, 0, "-")
			return
		}

	case ReuseMustRevalidateCanStale:
		log.Debugf("cache.Handler.ServeHTTP: '%v' must revalidate (but allowed stale)\n", cacheKey)
		r.Header.Set("If-Modified-Since", cacheObj.RespRespTime.Format(time.RFC1123))
		oldCacheObj := cacheObj
		cacheObj, err = retrier.Get(r, cacheObj)
		if err != nil {
			log.Errorf("retrying get error - serving stale as allowed: %v\n", err)
			cacheObj = oldCacheObj
		}
	}
	log.Debugf("cache.Handler.ServeHTTP: '%v' responding with %v\n", cacheKey, cacheObj.Code)
	modHdrs(&cacheObj.RespHeaders, remappingProducer.rule.ToClientHeaders.Drop, remappingProducer.rule.ToClientHeaders.Set)
	bytesSent, err := h.respond(w, cacheObj.Code, cacheObj.RespHeaders, cacheObj.Body, connectionClose)
	if err != nil {
		log.Errorln(time.Now().Format(time.RFC3339Nano) + " " + r.RemoteAddr + " " + r.Method + " " + r.RequestURI + ": responding: " + err.Error())
	}
	tryFlush(w)
	statLog.Log(cacheObj.Code, bytesSent, err == nil, true, isCacheHit(canReuseStored, cacheObj.OriginCode), false, cacheObj.OriginCode, cacheObj.Size, cacheObj.ProxyURL)
}

func tryFlush(w http.ResponseWriter) {
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
}

// GetClientIPPort returns the client IP address of the given request, and the port. It returns the first x-forwarded-for IP if any, else the RemoteAddr
func GetClientIPPort(r *http.Request) (string, string) {
	xForwardedFor := r.Header.Get("X-FORWARDED-FOR")
	ips := strings.Split(xForwardedFor, ",")
	ip, port, err := net.SplitHostPort(r.RemoteAddr)
	if len(ips) < 1 || ips[0] == "" {
		if err != nil {
			return r.RemoteAddr, port // TODO log?
		}
		return ip, port
	}
	return strings.TrimSpace(ips[0]), port
}

func atsEventLogStr(
	timestamp time.Time, // (prefix)
	clientIP string, // chi
	selfHostname string,
	reqHost string, // phn
	reqPort string, // php
	originHost string, // shn
	scheme string, // url
	url string, // url
	method string, // cqhm
	protocol string, // cqhv
	respCode int, // pssc
	timeToServe time.Duration, // ttms
	bytesSent uint64, // b
	originStatus int, // sssc
	originBytes uint64, // sscl
	successfullyRespondedToClient bool, // cfsc
	successfullyGotFromOrigin bool, // pfsc
	cacheHit string, // crc
	proxyUsed string, // phr
	thisProxyName string, // pqsn
	clientUserAgent string, // client user agent
	xmt string, // moneytrace header
) string {
	unixNano := timestamp.UnixNano()
	unixSec := unixNano / NSPerSec
	unixFrac := 1 / (unixNano % NSPerSec)
	unixFracStr := strconv.FormatInt(unixFrac, 10)
	if len(unixFracStr) > 3 {
		unixFracStr = unixFracStr[:3]
	}
	cfsc := "FIN"
	if !successfullyRespondedToClient {
		cfsc = "INTR"
	}
	pfsc := "FIN"
	if !successfullyGotFromOrigin {
		pfsc = "INTR"
	}

	// TODO escape quotes within useragent, moneytrace
	clientUserAgent = `"` + clientUserAgent + `"`
	if xmt == "" {
		xmt = `"-"`
	} else {
		xmt = `"` + xmt + `"`
	}

	return strconv.FormatInt(unixSec, 10) + "." + unixFracStr + " chi=" + clientIP + " phn=" + selfHostname + " php=" + reqPort + " shn=" + originHost + " url=" + scheme + "://" + reqHost + url + " cqhn=" + method + " cqhv=" + protocol + " pssc=" + strconv.FormatInt(int64(respCode), 10) + " ttms=" + strconv.FormatInt(int64(timeToServe/time.Millisecond), 10) + " b=" + strconv.FormatInt(int64(bytesSent), 10) + " sssc=" + strconv.FormatInt(int64(originStatus), 10) + " sscl=" + strconv.FormatInt(int64(originBytes), 10) + " cfsc=" + cfsc + " pfsc=" + pfsc + " crc=" + cacheHit + " phr=" + proxyUsed + " psqn=" + thisProxyName + " uas=" + clientUserAgent + " xmt=" + xmt + "\n"
}

func isCacheHit(reuse Reuse, originCode int) bool {
	return reuse == ReuseCan || ((reuse == ReuseMustRevalidate || reuse == ReuseMustRevalidateCanStale) && (originCode > 299 && originCode < 400))
}

// serveRuleNotFound writes the appropriate response to the client, via given writer, for when no remap rule was found for a request.
func serveRuleNotFound(w http.ResponseWriter) (int, uint64, error) {
	return serveErr(w, http.StatusNotFound)
}

// serveNotAllowed writes the appropriate response to the client, via given writer, for when the client's IP is not allowed for the requested rule.
func serveNotAllowed(w http.ResponseWriter) (int, uint64, error) {
	return serveErr(w, http.StatusForbidden)
}

// serveReqErr writes the appropriate response to the client, via given writer, for a generic request error. Returns the code sent, the body bytes written, and any write error.
func serveReqErr(w http.ResponseWriter) (int, uint64, error) {
	return serveErr(w, http.StatusBadRequest)
}

func serveErr(w http.ResponseWriter, code int) (int, uint64, error) {
	w.WriteHeader(code)
	bytesWritten, err := w.Write([]byte(http.StatusText(code)))
	return code, uint64(bytesWritten), err
}

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
func (h *Handler) respond(w http.ResponseWriter, code int, header http.Header, body []byte, connectionClose bool) (uint64, error) {
	dH := w.Header()
	web.CopyHeaderTo(header, &dH)
	if connectionClose {
		dH.Add("Connection", "close")
	}
	w.WriteHeader(code)
	bytesWritten, err := w.Write(body) // get the less-accurate body bytes written, in case we can't get the more accurate intercepted data

	// bytesWritten = int(WriteStats(stats, w, conn, reqFQDN, remoteAddr, code, uint64(bytesWritten))) // TODO write err to stats?
	return uint64(bytesWritten), err
}

// WriteStats writes to the remapRuleStats, and returns the bytes written to the connection
func WriteStats(stats Stats, w http.ResponseWriter, conn *web.InterceptConn, reqFQDN string, remoteAddr string, code int, bytesWritten uint64, cacheHit bool) uint64 {
	remapRuleStats, ok := stats.Remap().Stats(reqFQDN)
	if !ok {
		log.Errorf("Remap rule %v not in Stats\n", reqFQDN)
		return bytesWritten
	}

	if wFlusher, ok := w.(http.Flusher); !ok {
		log.Errorf("ResponseWriter is not a Flusher, could not flush written bytes, stat out_bytes will be inaccurate!\n")
	} else {
		wFlusher.Flush()
	}

	bytesRead := 0 // TODO get somehow? Count body? Sum header?
	if conn != nil {
		bytesRead = conn.BytesRead()
		bytesWritten = uint64(conn.BytesWritten()) // get the more accurate interceptConn bytes written, if we can
		// Don't log - the Handler has already logged the failure to get the conn
	}

	// bytesRead, bytesWritten := getConnInfoAndDestroyWriter(w, stats, remapRuleName)
	remapRuleStats.AddInBytes(uint64(bytesRead))
	remapRuleStats.AddOutBytes(uint64(bytesWritten))

	if cacheHit {
		stats.AddCacheHit()
		remapRuleStats.AddCacheHit()
	} else {
		stats.AddCacheMiss()
		remapRuleStats.AddCacheMiss()
	}

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
	return bytesWritten
}
