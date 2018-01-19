package web

import (
	"net"
	"net/http"
	"strings"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
)

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

type Hdr struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// ModHdrs drops and sets headers in h according to drop and set.
func ModHdrs(h *http.Header, drop []string, set []Hdr) {
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

// TryFlush calls Flush on w if it's an http.Flusher. If it isn't, it returns without error.
func TryFlush(w http.ResponseWriter) {
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
}
