package main

import (
	"fmt"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/experimental/tocookie"
	"log" // TODO change to traffic_monitor_golang/common/log
	"net/http"
	"time"
)

func wrapAuth(h RegexHandlerFunc, noAuth bool, secret string) RegexHandlerFunc {
	if noAuth {
		return h
	}
	return func(w http.ResponseWriter, r *http.Request, p ParamMap) {
		handleUnauthorized := func(reason string) {
			log.Printf("%v %v %v sent 401 - %v\n", time.Now(), r.RemoteAddr, r.URL.Path, reason)
			status := http.StatusUnauthorized
			w.WriteHeader(status)
			fmt.Fprintf(w, http.StatusText(status))
		}

		cookie, err := r.Cookie(tocookie.Name)
		if err != nil {
			handleUnauthorized("error getting cookie: " + err.Error())
			return
		}

		if cookie == nil {
			handleUnauthorized("no auth cookie")
			return
		}

		oldCookie, err := tocookie.Parse(secret, cookie.Value)
		if err != nil {
			handleUnauthorized("cookie error: " + err.Error())
			return
		}

		newCookieVal := tocookie.Refresh(oldCookie, secret)
		http.SetCookie(w, &http.Cookie{Name: tocookie.Name, Value: newCookieVal})
		h(w, r, p)
	}
}

func wrapLogTime(h RegexHandlerFunc) RegexHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, p ParamMap) {
		start := time.Now()
		defer func() {
			now := time.Now()
			log.Printf("%v %v served %v in %v\n", now, r.RemoteAddr, r.URL.Path, now.Sub(start))
		}()
		h(w, r, p)
	}
}
