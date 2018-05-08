package remap

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
	"strings"
	"time"

	"github.com/apache/incubator-trafficcontrol/grove/cacheobj"
	"github.com/apache/incubator-trafficcontrol/grove/remapdata"
	"github.com/apache/incubator-trafficcontrol/grove/web"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
)

// ValidHTTPCodes provides fast lookup whether a HTTP response code is valid per RFC7234§3
var ValidHTTPCodes = map[int]struct{}{
	200: {},
	201: {},
	202: {},
	203: {},
	204: {},
	205: {},
	206: {},
	207: {},
	208: {},
	226: {},

	300: {},
	301: {},
	302: {},
	303: {},
	304: {},
	305: {},
	306: {},
	307: {},
	308: {},

	400: {},
	401: {},
	402: {},
	403: {},
	404: {},
	405: {},
	406: {},
	407: {},
	408: {},
	409: {},
	410: {},
	411: {},
	412: {},
	413: {},
	414: {},
	415: {},
	416: {},
	417: {},
	418: {},
	421: {},
	422: {},
	423: {},
	424: {},
	428: {},
	429: {},
	431: {},
	451: {},

	500: {},
	501: {},
	502: {},
	503: {},
	504: {},
	505: {},
	506: {},
	507: {},
	508: {},
	510: {},
	511: {},
}

// cacheableResponseCodes provides fast lookup whether a HTTP response code is cacheable by default, per RFC7234§3
var defaultCacheableResponseCodes = map[int]struct{}{
	200: {},
	203: {},
	204: {},
	206: {},
	300: {},
	301: {},
	404: {},
	405: {},
	410: {},
	414: {},
	501: {},
}

// codeUnderstood returns whether the given response code is understood by this cache. Required by RFC7234§3
func codeUnderstood(code int) bool {
	_, ok := ValidHTTPCodes[code]
	return ok
}

// CanCache returns whether an object can be cached per RFC 7234, based on the request headers, response headers, and response code. If strictRFC is false, this ignores request headers denying cacheability such as `no-cache`, in order to protect origins.
// TODO add options to ignore/violate request cache-control (to protect origins)
func CanCache(reqMethod string, reqHeaders http.Header, respCode int, respHeaders http.Header, strictRFC bool) bool {
	log.Debugf("CanCache start\n")
	if reqMethod != http.MethodGet {
		return false // for now, we only support GET as a cacheable method.
	}
	reqCacheControl := web.ParseCacheControl(reqHeaders)
	respCacheControl := web.ParseCacheControl(respHeaders)
	log.Debugf("CanCache reqCacheControl %+v respCacheControl %+v\n", reqCacheControl, respCacheControl)
	return canStoreResponse(respCode, respHeaders, reqCacheControl, respCacheControl, strictRFC) && canStoreAuthenticated(reqCacheControl, respCacheControl)
}

// CanReuseStored checks the constraints in RFC7234§4
func CanReuseStored(reqHeaders http.Header, respHeaders http.Header, reqCacheControl web.CacheControl, respCacheControl web.CacheControl, respReqHeaders http.Header, respReqTime time.Time, respRespTime time.Time, strictRFC bool) remapdata.Reuse {
	// TODO: remove allowed_stale, check in cache manager after revalidate fails? (since RFC7234§4.2.4 prohibits serving stale response unless disconnected).

	if !selectedHeadersMatch(reqHeaders, respReqHeaders, strictRFC) {
		log.Debugf("CanReuseStored false - selected headers don't match\n") // debug
		return remapdata.ReuseCannot
	}

	if !fresh(respHeaders, respCacheControl, respReqTime, respRespTime) {
		allowedStale := allowedStale(respHeaders, reqCacheControl, respCacheControl, respReqTime, respRespTime, strictRFC)
		log.Debugf("CanReuseStored not fresh, allowed stale: %v\n", allowedStale) // debug
		return allowedStale
	}

	if hasPragmaNoCache(reqHeaders) && strictRFC {
		log.Debugf("CanReuseStored MustRevalidate - has pragma no-cache\n")
		return remapdata.ReuseMustRevalidate
	}

	if _, ok := reqCacheControl["no-cache"]; ok && strictRFC {
		log.Debugf("CanReuseStored false - request has cache-control no-cache\n")
		return remapdata.ReuseCannot
	}

	if _, ok := respCacheControl["no-cache"]; ok {
		log.Debugf("CanReuseStored false - response has cache-control no-cache\n")
		return remapdata.ReuseCannot
	}

	if strictRFC && !inMinFresh(respHeaders, reqCacheControl, respCacheControl, respReqTime, respRespTime) {
		return remapdata.ReuseMustRevalidate
	}

	log.Debugf("CanReuseStored true (respCacheControl %+v)\n", respCacheControl)
	return remapdata.ReuseCan
}

// CanReuse is a helper wrapping CanReuseStored, returning a boolean rather than an enum, for when it's known whether MustRevalidate can be used.
func CanReuse(reqHeader http.Header, reqCacheControl web.CacheControl, cacheObj *cacheobj.CacheObj, strictRFC bool, revalidateCanReuse bool) bool {
	canReuse := CanReuseStored(reqHeader, cacheObj.RespHeaders, reqCacheControl, cacheObj.RespCacheControl, cacheObj.ReqHeaders, cacheObj.ReqRespTime, cacheObj.RespRespTime, strictRFC)
	return canReuse == remapdata.ReuseCan || (canReuse == remapdata.ReuseMustRevalidate && revalidateCanReuse)
}

// canStoreAuthenticated checks the constraints in RFC7234§3.2
// TODO: ensure RFC7234§3.2 requirements that max-age=0, must-revlaidate, s-maxage=0 are revalidated
func canStoreAuthenticated(reqCacheControl, respCacheControl web.CacheControl) bool {
	if _, ok := reqCacheControl["authorization"]; !ok {
		return true
	}
	if _, ok := respCacheControl["must-revalidate"]; ok {
		return true
	}
	if _, ok := respCacheControl["public"]; ok {
		return true
	}
	if _, ok := respCacheControl["s-maxage"]; ok {
		return true
	}
	log.Debugf("CanStoreAuthenticated false: has authorization, and no must-revalidate/public/s-maxage\n")
	return false
}

// CanStoreResponse checks the constraints in RFC7234
func canStoreResponse(
	respCode int,
	respHeaders http.Header,
	reqCacheControl web.CacheControl,
	respCacheControl web.CacheControl,
	strictRFC bool,
) bool {
	if _, ok := reqCacheControl["no-store"]; strictRFC && ok {
		log.Debugf("CanStoreResponse false: request has no-store\n")
		return false
	}
	if _, ok := respCacheControl["no-store"]; ok {
		log.Debugf("CanStoreResponse false: response has no-store\n") // RFC7234§5.2.2.3
		return false
	}
	if _, ok := respCacheControl["no-cache"]; ok {
		log.Debugf("CanStoreResponse false: response has no-cache\n") // RFC7234§5.2.2.2
		return false
	}
	if _, ok := respCacheControl["private"]; ok {
		log.Debugf("CanStoreResponse false: has private\n")
		return false
	}
	if _, ok := respCacheControl["authorization"]; ok {
		log.Debugf("CanStoreResponse false: has authorization\n")
		return false
	}
	if !cacheControlAllows(respCode, respHeaders, respCacheControl) {
		log.Debugf("CanStoreResponse false: CacheControlAllows false\n")
		return false
	}
	log.Debugf("CanStoreResponse true\n")
	return true
}

func cacheControlAllows(
	respCode int,
	respHeaders http.Header,
	respCacheControl web.CacheControl,
) bool {
	if _, ok := respHeaders["Expires"]; ok {
		return true
	}
	if _, ok := respCacheControl["max-age"]; ok {
		return true
	}
	if _, ok := respCacheControl["s-maxage"]; ok {
		return true
	}
	if extensionAllows() {
		return true
	}
	if codeDefaultCacheable(respCode) {
		return true
	}
	log.Debugf("CacheControlAllows false: no expires, no max-age, no s-max-age, no extension allows, code not default cacheable\n")
	return false
}

// extensionAllows returns whether a cache-control extension allows the response to be cached, per RFC7234§3 and RFC7234§5.2.3.
func extensionAllows() bool {
	// This MUST return false unless a specific Cache Control cache-extension token exists for an extension which allows. Which is to say, returning true here without a cache-extension token is in strict violation of RFC7234.
	// In practice, all returning true does is override whether a response code is default-cacheable. If we wanted to do that, it would be better to make codeDefaultCacheable take a strictRFC parameter.
	return false
}

func codeDefaultCacheable(code int) bool {
	_, ok := defaultCacheableResponseCodes[code]
	return ok
}

// Fresh checks the constraints in RFC7234§4 via RFC7234§4.2
func fresh(
	respHeaders http.Header,
	respCacheControl web.CacheControl,
	respReqTime time.Time,
	respRespTime time.Time,
) bool {
	freshnessLifetime := web.GetFreshnessLifetime(respHeaders, respCacheControl)
	currentAge := web.GetCurrentAge(respHeaders, respReqTime, respRespTime)
	log.Debugf("Fresh: freshnesslifetime %v currentAge %v\n", freshnessLifetime, currentAge)
	fresh := freshnessLifetime > currentAge
	return fresh
}

// inMinFresh returns whether the given response is within the `min-fresh` request directive. If no `min-fresh` directive exists in the request, `true` is returned.
func inMinFresh(respHeaders http.Header, reqCacheControl web.CacheControl, respCacheControl web.CacheControl, respReqTime time.Time, respRespTime time.Time) bool {
	minFresh, ok := web.GetHTTPDeltaSecondsCacheControl(reqCacheControl, "min-fresh")
	if !ok {
		return true // no min-fresh => within min-fresh
	}
	freshnessLifetime := web.GetFreshnessLifetime(respHeaders, respCacheControl)
	currentAge := web.GetCurrentAge(respHeaders, respReqTime, respRespTime)
	inMinFresh := minFresh < (freshnessLifetime - currentAge)
	log.Debugf("inMinFresh minFresh %v freshnessLifetime %v currentAge %v => %v < (%v - %v) = %v\n", minFresh, freshnessLifetime, currentAge, minFresh, freshnessLifetime, currentAge, inMinFresh)
	return inMinFresh
}

// TODO add warning generation funcs

// AllowedStale checks the constraints in RFC7234§4 via RFC7234§4.2.4
func allowedStale(respHeaders http.Header, reqCacheControl web.CacheControl, respCacheControl web.CacheControl, respReqTime time.Time, respRespTime time.Time, strictRFC bool) remapdata.Reuse {
	// TODO return remapdata.ReuseMustRevalidate where permitted
	_, reqHasMaxAge := reqCacheControl["max-age"]
	_, reqHasMaxStale := reqCacheControl["max-stale"]
	_, respHasMustReval := respCacheControl["must-revalidate"]
	_, respHasProxyReval := respCacheControl["proxy-revalidate"]
	log.Debugf("AllowedStale: reqHasMaxAge %v reqHasMaxStale %v strictRFC %v\n", reqHasMaxAge, reqHasMaxStale, strictRFC)
	if respHasMustReval || respHasProxyReval {
		log.Debugf("AllowedStale: returning mustreval - must-revalidate\n")
		return remapdata.ReuseMustRevalidate
	}
	if strictRFC && reqHasMaxAge && !reqHasMaxStale {
		log.Debugf("AllowedStale: returning can - strictRFC & reqHasMaxAge & !reqHasMaxStale\n")
		return remapdata.ReuseMustRevalidateCanStale
	}
	if _, ok := respCacheControl["no-cache"]; ok {
		log.Debugf("AllowedStale: returning reusecannot - no-cache\n")
		return remapdata.ReuseCannot // TODO verify RFC doesn't allow Revalidate here
	}
	if _, ok := respCacheControl["no-store"]; ok {
		log.Debugf("AllowedStale: returning reusecannot - no-store\n")
		return remapdata.ReuseCannot // TODO verify RFC doesn't allow revalidate here
	}
	if !inMaxStale(respHeaders, respCacheControl, respReqTime, respRespTime) {
		log.Debugf("AllowedStale: returning mustreval - not in max stale\n")
		return remapdata.ReuseMustRevalidate // TODO verify RFC allows
	}
	log.Debugf("AllowedStale: returning can - all preconditions passed\n")
	return remapdata.ReuseMustRevalidateCanStale
}

// InMaxStale returns whether the given response is within the `max-stale` request directive. If no `max-stale` directive exists in the request, `true` is returned.
func inMaxStale(respHeaders http.Header, respCacheControl web.CacheControl, respReqTime time.Time, respRespTime time.Time) bool {
	maxStale, ok := web.GetHTTPDeltaSecondsCacheControl(respCacheControl, "max-stale")
	if !ok {
		// maxStale = 5 // debug
		return true // no max-stale => within max-stale
	}
	freshnessLifetime := web.GetFreshnessLifetime(respHeaders, respCacheControl)
	currentAge := web.GetCurrentAge(respHeaders, respReqTime, respRespTime)
	log.Errorf("DEBUGR InMaxStale maxStale %v freshnessLifetime %v currentAge %v => %v > (%v, %v)\n", maxStale, freshnessLifetime, currentAge, maxStale, currentAge, freshnessLifetime) // DEBUG
	inMaxStale := maxStale > (currentAge - freshnessLifetime)
	return inMaxStale
}

// SelectedHeadersMatch checks the constraints in RFC7234§4.1
// TODO: change caching to key on URL+headers, so multiple requests for the same URL with different vary headers can be cached?
func selectedHeadersMatch(reqHeaders http.Header, respReqHeaders http.Header, strictRFC bool) bool {
	varyHeaders, ok := reqHeaders["vary"]
	if !strictRFC && !ok {
		return true
	}
	if len(varyHeaders) == 0 {
		return true
	}
	varyHeader := varyHeaders[0]

	if varyHeader == "*" {
		return false
	}
	varyHeader = strings.ToLower(varyHeader)
	varyHeaderHeaders := strings.Split(varyHeader, ",")
	for _, header := range varyHeaderHeaders {
		if _, ok := respReqHeaders[header]; !ok {
			return false
		}
	}
	return true
}

// HasPragmaNoCache returns whether the given headers have a `pragma: no-cache` which is to be considered per HTTP/1.1. This specifically returns false if `cache-control` exists, even if `pragma: no-cache` exists, per RFC7234§5.4
func hasPragmaNoCache(reqHeaders http.Header) bool {
	if _, ok := reqHeaders["Cache-Control"]; ok {
		return false
	}
	pragmas, ok := reqHeaders["pragma"]
	if !ok {
		return false
	}
	if len(pragmas) == 0 {
		return false
	}
	pragma := pragmas[0]

	if strings.HasPrefix(pragma, "no-cache") { // RFC7234§5.4 specifically requires no-cache be the first pragma
		return true
	}
	return false
}
