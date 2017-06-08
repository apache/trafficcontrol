package grove

import (
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/common/log"
)

type Reuse int

const (
	ReuseCan = iota
	ReuseCannot
	ReuseMustRevalidate
)

// validHttpCodes provides fast lookup whether a HTTP response code is valid per RFC7234§3
var validHttpCodes = map[int]struct{}{
	200: struct{}{},
	201: struct{}{},
	202: struct{}{},
	203: struct{}{},
	204: struct{}{},
	205: struct{}{},
	206: struct{}{},
	207: struct{}{},
	208: struct{}{},
	226: struct{}{},

	300: struct{}{},
	301: struct{}{},
	302: struct{}{},
	303: struct{}{},
	304: struct{}{},
	305: struct{}{},
	306: struct{}{},
	307: struct{}{},
	308: struct{}{},

	400: struct{}{},
	401: struct{}{},
	402: struct{}{},
	403: struct{}{},
	404: struct{}{},
	405: struct{}{},
	406: struct{}{},
	407: struct{}{},
	408: struct{}{},
	409: struct{}{},
	410: struct{}{},
	411: struct{}{},
	412: struct{}{},
	413: struct{}{},
	414: struct{}{},
	415: struct{}{},
	416: struct{}{},
	417: struct{}{},
	418: struct{}{},
	421: struct{}{},
	422: struct{}{},
	423: struct{}{},
	424: struct{}{},
	428: struct{}{},
	429: struct{}{},
	431: struct{}{},
	451: struct{}{},

	500: struct{}{},
	501: struct{}{},
	502: struct{}{},
	503: struct{}{},
	504: struct{}{},
	505: struct{}{},
	506: struct{}{},
	507: struct{}{},
	508: struct{}{},
	510: struct{}{},
	511: struct{}{},
}

// cacheableResponseCodes provides fast lookup whether a HTTP response code is cacheable by default, per RFC7234§3
var defaultCacheableResponseCodes = map[int]struct{}{
	200: struct{}{},
	203: struct{}{},
	204: struct{}{},
	206: struct{}{},
	300: struct{}{},
	301: struct{}{},
	404: struct{}{},
	405: struct{}{},
	410: struct{}{},
	414: struct{}{},
	501: struct{}{},
}

// CodeUnderstood returns whether the given response code is understood by this cache. Required by RFC7234§3
func CodeUnderstood(code int) bool {
	_, ok := validHttpCodes[code]
	return ok
}

// TODO add options to ignore/violate request cache-control (to protect origins)
// CanCache returns whether an object can be cached per RFC 7234, based on the request headers, response headers, and response code. If strictRFC is false, this ignores request headers denying cacheability such as `no-cache`, in order to protect origins.
func CanCache(reqHeaders http.Header, respCode int, respHeaders http.Header, strictRFC bool) bool {
	log.Debugf("CanCache start\n")
	reqCacheControl := ParseCacheControl(reqHeaders)
	respCacheControl := ParseCacheControl(respHeaders)
	log.Debugf("CanCache reqCacheControl %+v respCacheControl %+v\n", reqCacheControl, respCacheControl)
	return CanStoreResponse(respCode, respHeaders, reqCacheControl, respCacheControl, strictRFC) && CanStoreAuthenticated(reqCacheControl, respCacheControl)
}

// CanStoreAuthenticated checks the constraints in RFC7234§3.2
// TODO: ensure RFC7234§3.2 requirements that max-age=0, must-revlaidate, s-maxage=0 are revalidated
func CanStoreAuthenticated(reqCacheControl, respCacheControl CacheControl) bool {
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
func CanStoreResponse(
	respCode int,
	respHeaders http.Header,
	reqCacheControl CacheControl,
	respCacheControl CacheControl,
	strictRFC bool,
) bool {
	if _, ok := reqCacheControl["no-store"]; !strictRFC && ok {
		log.Debugf("CanStoreResponse false: request has no-store\n")
		return false
	}
	if _, ok := respCacheControl["no-store"]; ok {
		log.Debugf("CanStoreResponse false: response has no-store\n")
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
	if !CacheControlAllows(respCode, respHeaders, respCacheControl) {
		log.Debugf("CanStoreResponse false: CacheControlAllows false\n")
		return false
	}
	log.Debugf("CanStoreResponse true\n")
	return true
}

func CacheControlAllows(
	respCode int,
	respHeaders http.Header,
	respCacheControl CacheControl,
) bool {
	if _, ok := respHeaders["Expires"]; ok {
		return true
	}
	if _, ok := respCacheControl["max-age"]; ok {
		return true
	}
	if _, ok := respCacheControl["s-max-age"]; ok {
		return true
	}
	if ExtensionAllows() {
		return true
	}
	if CodeDefaultCacheable(respCode) {
		return true
	}
	log.Debugf("CacheControlAllows false: no expires, no max-age, no s-max-age, no extension allows, code not default cacheable\n")
	return false
}

// ExtensionAllows returns whether a cache-control extension allows the response to be cached, per RFC7234§3 and RFC7234§5.2.3. Note this currently fulfills the literal wording of the section, but cache-control extensions may override any requirements, in which case the logic of can_cache? outside this function would have to be changed.
func ExtensionAllows() bool {
	return true
}

func CodeDefaultCacheable(code int) bool {
	if _, ok := defaultCacheableResponseCodes[code]; ok {
		return true
	}
	return false
}

// CanReuseStored checks the constraints in RFC7234§4
func CanReuseStored(reqHeaders http.Header, respHeaders http.Header, reqCacheControl CacheControl, respCacheControl CacheControl, respReqHeaders http.Header, respReqTime time.Time, respRespTime time.Time, strictRFC bool) Reuse {
	// TODO: remove allowed_stale, check in cache manager after revalidate fails? (since RFC7234§4.2.4 prohibits serving stale response unless disconnected).

	if !SelectedHeadersMatch(reqHeaders, respReqHeaders, strictRFC) {
		log.Debugf("CanReuseStored false - selected headers don't match\n") // debug
		return ReuseCannot
	}

	if !Fresh(respHeaders, respCacheControl, respReqTime, respRespTime) {
		allowedStale := AllowedStale(respHeaders, reqCacheControl, respCacheControl, respReqTime, respRespTime, strictRFC)
		log.Debugf("CanReuseStored not fresh, allowed stale: %v\n", allowedStale) // debug
		return allowedStale
	}

	if HasPragmaNoCache(reqHeaders) && strictRFC {
		log.Debugf("CanReuseStored MustRevalidate - has pragma no-cache\n")
		return ReuseMustRevalidate
	}
	if _, ok := reqCacheControl["no-cache"]; ok && strictRFC {
		log.Debugf("CanReuseStored false - request has cache-control no-cache\n")
		return ReuseCannot
	}
	if _, ok := respCacheControl["no-cache"]; ok {
		log.Debugf("CanReuseStored false - response has cache-control no-cache\n")
		return ReuseCannot
	}
	log.Debugf("CanReuseStored true (respCacheControl %+v)\n", respCacheControl)
	return ReuseCan
}

// Fresh checks the constraints in RFC7234§4 via RFC7234§4.2
func Fresh(
	respHeaders http.Header,
	respCacheControl CacheControl,
	respReqTime time.Time,
	respRespTime time.Time,
) bool {
	freshnessLifetime := GetFreshnessLifetime(respHeaders, respCacheControl)
	currentAge := GetCurrentAge(respHeaders, respReqTime, respRespTime)
	fresh := freshnessLifetime > currentAge
	return fresh
}

// GetHTTPDate is a helper function which gets an HTTP date from the given map (which is typically a `http.Header` or `CacheControl`. Returns false if the given key doesn't exist in the map, or if the value isn't a valid HTTP Date per RFC2616§3.3.
func GetHTTPDate(headers http.Header, key string) (time.Time, bool) {
	maybeDate, ok := headers[key]
	if !ok {
		return time.Time{}, false
	}
	if len(maybeDate) == 0 {
		return time.Time{}, false
	}
	return ParseHTTPDate(maybeDate[0])
}

// GetHTTPDeltaSeconds is a helper function which gets an HTTP Delta Seconds from the given map (which is typically a `http.Header` or `CacheControl`. Returns false if the given key doesn't exist in the map, or if the value isn't a valid Delta Seconds per RFC2616§3.3.2.
func GetHTTPDeltaSeconds(m map[string][]string, key string) (time.Duration, bool) {
	maybeSeconds, ok := m[key]
	if !ok {
		return 0, false
	}
	if len(maybeSeconds) == 0 {
		return 0, false
	}
	maybeSec := maybeSeconds[0]

	seconds, err := strconv.ParseUint(maybeSec, 10, 64)
	if err != nil {
		return 0, false
	}
	return time.Duration(seconds) * time.Second, true
}

// GetHTTPDeltaSeconds is a helper function which gets an HTTP Delta Seconds from the given map (which is typically a `http.Header` or `CacheControl`. Returns false if the given key doesn't exist in the map, or if the value isn't a valid Delta Seconds per RFC2616§3.3.2.
func GetHTTPDeltaSecondsCacheControl(m map[string]string, key string) (time.Duration, bool) {
	maybeSec, ok := m[key]
	if !ok {
		return 0, false
	}
	seconds, err := strconv.ParseUint(maybeSec, 10, 64)
	if err != nil {
		return 0, false
	}
	return time.Duration(seconds) * time.Second, true
}

// GetFreshnessLifetime calculates the freshness_lifetime per RFC7234§4.2.1
func GetFreshnessLifetime(respHeaders http.Header, respCacheControl CacheControl) time.Duration {
	if s, ok := GetHTTPDeltaSecondsCacheControl(respCacheControl, "s-maxage"); ok {
		return s
	}
	if s, ok := GetHTTPDeltaSecondsCacheControl(respCacheControl, "max-age"); ok {
		return s
	}

	getExpires := func() (time.Duration, bool) {
		expires, ok := GetHTTPDate(respHeaders, "Expires")
		if !ok {
			return 0, false
		}
		date, ok := GetHTTPDate(respHeaders, "Date")
		if !ok {
			return 0, false
		}
		return expires.Sub(date), true
	}
	if s, ok := getExpires(); ok {
		return s
	}
	return 0
}

const Day = time.Hour * time.Duration(24)

// HeuristicFreshness follows the recommendation of RFC7234§4.2.2 and returns the min of 10% of the (Date - Last-Modified) headers and 24 hours, if they exist, and 24 hours if they don't.
// TODO: smarter and configurable heuristics
func HeuristicFreshness(respHeaders http.Header) time.Duration {
	sinceLastModified, ok := SinceLastModified(respHeaders)
	if !ok {
		return Day
	}
	freshness := time.Duration(math.Min(float64(Day), float64(sinceLastModified)))
	return freshness
}

func SinceLastModified(headers http.Header) (time.Duration, bool) {
	lastModified, ok := GetHTTPDate(headers, "last-modified")
	if !ok {
		return 0, false
	}
	date, ok := GetHTTPDate(headers, "date")
	if !ok {
		return 0, false
	}
	return date.Sub(lastModified), true
}

// ParseHTTPDate parses the given RFC7231§7.1.1 HTTP-date
func ParseHTTPDate(d string) (time.Time, bool) {
	if t, err := time.Parse(time.RFC1123, d); err == nil {
		return t, true
	}
	if t, err := time.Parse(time.RFC850, d); err == nil {
		return t, true
	}
	if t, err := time.Parse(time.ANSIC, d); err == nil {
		return t, true
	}
	return time.Time{}, false

}

// AgeValue is used to calculate current_age per RFC7234§4.2.3
func AgeValue(respHeaders http.Header) time.Duration {
	s, ok := GetHTTPDeltaSeconds(respHeaders, "age")
	if !ok {
		return 0
	}
	return s
}

// DateValue is used to calculate current_age per RFC7234§4.2.3. It returns time, or false if the response had no Date header (in violation of HTTP/1.1).
func DateValue(respHeaders http.Header) (time.Time, bool) {
	return GetHTTPDate(respHeaders, "date")
}

func ApparentAge(respHeaders http.Header, respRespTime time.Time) time.Duration {
	dateValue, ok := DateValue(respHeaders)
	if !ok {
		return 0 // TODO log warning?
	}
	rawAge := respRespTime.Sub(dateValue)
	return time.Duration(math.Max(0.0, float64(rawAge)))
}

func ResponseDelay(respReqTime time.Time, respRespTime time.Time) time.Duration {
	return respRespTime.Sub(respReqTime)
}

func CorrectedAgeValue(respHeaders http.Header, respReqTime time.Time, respRespTime time.Time) time.Duration {
	return AgeValue(respHeaders) + ResponseDelay(respReqTime, respRespTime)
}

func CorrectedInitialAge(respHeaders http.Header, respReqTime time.Time, respRespTime time.Time) time.Duration {
	return time.Duration(math.Max(float64(ApparentAge(respHeaders, respRespTime)), float64(CorrectedAgeValue(respHeaders, respReqTime, respRespTime))))
}

func ResidentTime(respRespTime time.Time) time.Duration {
	return time.Now().Sub(respRespTime)
}

func GetCurrentAge(respHeaders http.Header, respReqTime time.Time, respRespTime time.Time) time.Duration {
	return CorrectedInitialAge(respHeaders, respReqTime, respRespTime) + ResidentTime(respRespTime)
}

// TODO add min-fresh check

// TODO add warning generation funcs

// AllowedStale checks the constraints in RFC7234§4 via RFC7234§4.2.4
func AllowedStale(respHeaders http.Header, reqCacheControl CacheControl, respCacheControl CacheControl, respReqTime time.Time, respRespTime time.Time, strictRFC bool) Reuse {
	// TODO return ReuseMustRevalidate where permitted
	_, reqHasMaxAge := reqCacheControl["max-age"]
	_, reqHasMaxStale := reqCacheControl["max-stale"]
	if strictRFC && reqHasMaxAge && !reqHasMaxStale {
		return ReuseCan
	}
	if _, ok := respCacheControl["must-revalidate"]; ok {
		return ReuseMustRevalidate
	}
	if _, ok := respCacheControl["no-cache"]; ok {
		return ReuseCannot // TODO verify RFC doesn't allow Revalidate here
	}
	if _, ok := respCacheControl["no-store"]; ok {
		return ReuseCannot // TODO verify RFC doesn't allow revalidate here
	}
	if !InMaxStale(respHeaders, respCacheControl, respReqTime, respRespTime) {
		return ReuseMustRevalidate // TODO verify RFC allows
	}
	return ReuseCan
}

// InMaxStale returns whether the given response is within the `max-stale` request directive. If no `max-stale` directive exists in the request, `true` is returned.
func InMaxStale(respHeaders http.Header, respCacheControl CacheControl, respReqTime time.Time, respRespTime time.Time) bool {
	maxStale, ok := GetHTTPDeltaSecondsCacheControl(respCacheControl, "max-stale")
	if !ok {
		// maxStale = 5 // debug
		return true // no max-stale => within max-stale
	}
	freshnessLifetime := GetFreshnessLifetime(respHeaders, respCacheControl)
	currentAge := GetCurrentAge(respHeaders, respReqTime, respRespTime)
	log.Errorf("DEBUGR InMaxStale maxStale %v freshnessLifetime %v currentAge %v => %v > (%v, %v)\n", maxStale, freshnessLifetime, currentAge, maxStale, currentAge, freshnessLifetime) // DEBUG
	inMaxStale := maxStale > (currentAge - freshnessLifetime)
	return inMaxStale
}

// SelectedHeadersMatch checks the constraints in RFC7234§4.1
// TODO: change caching to key on URL+headers, so multiple requests for the same URL with different vary headers can be cached?
func SelectedHeadersMatch(reqHeaders http.Header, respReqHeaders http.Header, strictRFC bool) bool {
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
func HasPragmaNoCache(reqHeaders http.Header) bool {
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
