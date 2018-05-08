package web

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
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	"math"
	"strconv"
)

type Hdr struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ModHdrs struct {
	Set  []Hdr    `json:"set"`
	Drop []string `json:"drop"`
}

// Any returns whether any header modifications exist
func (mh *ModHdrs) Any() bool {
	return len(mh.Set) > 0 || len(mh.Drop) > 0
}

// Mod drops and sets the headers in h according to its rules.
func (mh *ModHdrs) Mod(h http.Header) {
	if len(h) == 0 { // this happens on a dial tcp timeout
		log.Debugf("modifyheaders: Header is  a nil map")
		return
	}
	for _, hdr := range mh.Drop {
		log.Debugf("modifyheaders: Dropping header %s\n", hdr)
		h.Del(hdr)
	}
	for _, hdr := range mh.Set {
		log.Debugf("modifyheaders: Setting header %s: %s \n", hdr.Name, hdr.Value)
		h.Set(hdr.Name, hdr.Value)
	}
}

func CopyHeaderTo(source http.Header, dest *http.Header) {
	for n, v := range source {
		for _, vv := range v {
			dest.Add(n, vv)
		}
	}
}

func CopyHeader(source http.Header) http.Header {
	dest := http.Header{}
	for n, v := range source {
		for _, vv := range v {
			dest.Add(n, vv)
		}
	}
	return dest
}

// GetClientIPPort returns the client IP address of the given request, and the port. It returns the first x-forwarded-for IP if any, else the RemoteAddr.
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

func GetIP(r *http.Request) (net.IP, error) {
	clientIPStr, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return nil, errors.New("malformed client address '" + r.RemoteAddr + "'")
	}
	clientIP := net.ParseIP(clientIPStr)
	if clientIP == nil {
		return nil, errors.New("malformed client IP address '" + clientIPStr + "'")
	}
	return clientIP, nil
}

// TryFlush calls Flush on w if it's an http.Flusher. If it isn't, it returns without error.
func TryFlush(w http.ResponseWriter) {
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
}

// request makes the given request and returns its response code, headers, body, the request time, response time, and any error.
func Request(transport *http.Transport, r *http.Request) (int, http.Header, []byte, time.Time, time.Time, error) {
	log.Debugf("request requesting %v headers %v\n", r.RequestURI, r.Header)
	rr := r

	reqTime := time.Now()
	resp, err := transport.RoundTrip(rr)
	respTime := time.Now()
	if err != nil {
		return 0, nil, nil, reqTime, respTime, errors.New("request error: " + err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	// TODO determine if respTime should go here

	if err != nil {
		return 0, nil, nil, reqTime, respTime, errors.New("reading response body: " + err.Error())
	}

	return resp.StatusCode, resp.Header, body, reqTime, respTime, nil
}

// Respond writes the given code, header, and body to the ResponseWriter. If connectionClose, a Connection: Close header is also written. Returns the bytes written, and any error.
func Respond(w http.ResponseWriter, code int, header http.Header, body []byte, connectionClose bool) (uint64, error) {
	// TODO move connectionClose to modhdr plugin
	dH := w.Header()
	CopyHeaderTo(header, &dH)
	if connectionClose {
		dH.Add("Connection", "close")
	}
	w.WriteHeader(code)
	bytesWritten, err := w.Write(body) // get the less-accurate body bytes written, in case we can't get the more accurate intercepted data

	// bytesWritten = int(WriteStats(stats, w, conn, reqFQDN, remoteAddr, code, uint64(bytesWritten))) // TODO write err to stats?
	return uint64(bytesWritten), err
}

// ServeReqErr writes the appropriate response to the client, via given writer, for a generic request error. Returns the code sent, the body bytes written, and any write error.
func ServeReqErr(w http.ResponseWriter) (int, uint64, error) {
	code := http.StatusBadRequest
	bytes, err := ServeErr(w, http.StatusBadRequest)
	return code, bytes, err
}

// ServeErr writes the given error code to w, writes the text for that code to the body, and returns the code sent, bytes written, and any write error.
func ServeErr(w http.ResponseWriter, code int) (uint64, error) {
	w.WriteHeader(code)
	bytesWritten, err := w.Write([]byte(http.StatusText(code)))
	return uint64(bytesWritten), err
}

// TryGetBytesWritten attempts to get the real bytes written to the given conn. It takes the bytesWritten as returned by Write(). It forcibly calls Flush() in order to force a write to the conn. Then, it attempts to get the more accurate bytes written to the Conn. If this fails, the given and less accurate bytesWritten is returned. If everything succeeds, the accurate bytes written to the Conn is returned.
func TryGetBytesWritten(w http.ResponseWriter, conn *InterceptConn, bytesWritten uint64) uint64 {
	if wFlusher, ok := w.(http.Flusher); !ok {
		log.Errorln("ResponseWriter is not a Flusher, could not flush written bytes, stat out_bytes will be inaccurate!")
	} else {
		wFlusher.Flush()
	}
	if conn != nil {
		return uint64(conn.BytesWritten()) // get the more accurate interceptConn bytes written, if we can
	}
	return bytesWritten
}

// GetHTTPDate is a helper function which gets an HTTP date from the given map (which is typically a `http.Header` or `CacheControl`. Returns false if the given key doesn't exist in the map, or if the value isn't a valid HTTP Date per RFC2616§3.3.
func GetHTTPDate(headers http.Header, key string) (time.Time, bool) {
	maybeDate := headers.Get(key)
	if maybeDate == "" {
		return time.Time{}, false
	}
	return ParseHTTPDate(maybeDate)
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

// RemapTextKey is the plugin shared data key inserted by grovetccfg for the Remap Line of the Delivery Service in Traffic Control, Traffic Ops.
const RemapTextKey = "remap_text"

const Day = time.Hour * time.Duration(24)

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

// HeuristicFreshness follows the recommendation of RFC7234§4.2.2 and returns the min of 10% of the (Date - Last-Modified) headers and 24 hours, if they exist, and 24 hours if they don't.
// TODO: smarter and configurable heuristics
func HeuristicFreshness(respHeaders http.Header) time.Duration {
	sinceLastModified, ok := sinceLastModified(respHeaders)
	if !ok {
		return Day
	}
	freshness := time.Duration(math.Min(float64(Day), float64(sinceLastModified)))
	return freshness
}

func sinceLastModified(headers http.Header) (time.Duration, bool) {
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
	return HeuristicFreshness(respHeaders)
}

// t6AgeValue is used to calculate current_age per RFC7234§4.2.3
func AgeValue(respHeaders http.Header) time.Duration {
	s, ok := GetHTTPDeltaSeconds(respHeaders, "age")
	if !ok {
		return 0
	}
	return s
}

func GetCurrentAge(respHeaders http.Header, respReqTime time.Time, respRespTime time.Time) time.Duration {
	correctedInitial := CorrectedInitialAge(respHeaders, respReqTime, respRespTime)
	resident := residentTime(respRespTime)
	log.Debugf("getCurrentAge: correctedInitialAge %v residentTime %v\n", correctedInitial, resident)
	return correctedInitial + resident
}

func CorrectedInitialAge(respHeaders http.Header, respReqTime time.Time, respRespTime time.Time) time.Duration {
	return time.Duration(math.Max(float64(ApparentAge(respHeaders, respRespTime)), float64(CorrectedAgeValue(respHeaders, respReqTime, respRespTime))))
}

func CorrectedAgeValue(respHeaders http.Header, respReqTime time.Time, respRespTime time.Time) time.Duration {
	return AgeValue(respHeaders) + responseDelay(respReqTime, respRespTime)
}

func responseDelay(respReqTime time.Time, respRespTime time.Time) time.Duration {
	return respRespTime.Sub(respReqTime)
}

func residentTime(respRespTime time.Time) time.Duration {
	return time.Now().Sub(respRespTime)
}

func ApparentAge(respHeaders http.Header, respRespTime time.Time) time.Duration {
	dateValue, ok := dateValue(respHeaders)
	if !ok {
		return 0 // TODO log warning?
	}
	rawAge := respRespTime.Sub(dateValue)
	return time.Duration(math.Max(0.0, float64(rawAge)))
}

// dateValue is used to calculate current_age per RFC7234§4.2.3. It returns time, or false if the response had no Date header (in violation of HTTP/1.1).
func dateValue(respHeaders http.Header) (time.Time, bool) {
	return GetHTTPDate(respHeaders, "date")
}

// FreshFor checks returns how long this object is still good for
func FreshFor(
	respHeaders http.Header,
	respCacheControl CacheControl,
	respReqTime time.Time,
	respRespTime time.Time,
) time.Duration {
	freshnessLifetime := GetFreshnessLifetime(respHeaders, respCacheControl)
	currentAge := GetCurrentAge(respHeaders, respReqTime, respRespTime)
	log.Debugf("FreshFor: freshnesslifetime %v currentAge %v\n", freshnessLifetime, currentAge)
	//fresh := freshnessLifetime > currentAge
	return freshnessLifetime - currentAge
}
