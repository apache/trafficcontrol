package plugin

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
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/v8/grove/web"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-rfc"

	"code.cloudfoundry.org/bytefmt"
)

func init() {
	AddPlugin(10000, Funcs{onRequest: cacheinspect})
}

// CacheStatsEndpoint is our reserved path
const CacheStatsEndpoint = "/_cacheinspect"

func writeHTMLPageHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`
<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<title>Grove Cache Inspector</title>
</head>
<body>
<pre>
`))

}

func writeHTMLPageFooter(w http.ResponseWriter) {
	w.Write([]byte(`
</pre>
</body>
</html>
`))
}

func cacheinspect(icfg interface{}, d OnRequestData) bool {
	if !strings.HasPrefix(d.R.URL.Path, CacheStatsEndpoint) {
		log.Debugf("plugin onrequest http_cacheinspect returning, not in path '" + d.R.URL.Path + "'\n")
		return false
	}

	log.Debugf("plugin onrequest http_cacheinspect calling\n")

	reqTime := time.Now()
	w := d.W
	req := d.R

	// TODO access log? Stats byte count?
	ip, err := web.GetIP(req)
	if err != nil {
		code := http.StatusInternalServerError
		w.WriteHeader(code)
		w.Write([]byte(http.StatusText(code)))
		log.Errorln("statHandler ServeHTTP failed to get IP: " + ip.String())
		return true
	}
	if !d.StatRules.Allowed(ip) {
		code := http.StatusForbidden
		w.WriteHeader(code)
		w.Write([]byte(http.StatusText(code)))
		log.Debugln("statHandler.ServeHTTP IP " + ip.String() + " FORBIDDEN") // TODO event?
		return true
	}

	respCode := http.StatusOK
	w.WriteHeader(respCode)
	writeHTMLPageHeader(w)
	qstringOptions := req.URL.Query()

	// The default cache = "", which is the default memcache
	cacheToDisplay := ""
	showSelectCache := false
	if cNameArr, cNamePresent := qstringOptions["cache"]; cNamePresent {
		showSelectCache = true
		cacheToDisplay = cNameArr[0]
	}
	if keyArr, showKey := qstringOptions["key"]; showKey {
		hLine := fmt.Sprintf("Key: %s cache: \"%s\"\n\n", keyArr[0], cacheToDisplay)
		w.Write([]byte(hLine))
		if cacheObject, ok := d.Stats.CachePeek(keyArr[0], cacheToDisplay); ok {
			for k, v := range cacheObject.ReqHeaders {
				w.Write([]byte(fmt.Sprintf("  > %s: %s\n", k, strings.Join(v, ","))))
			}
			w.Write([]byte("\n"))
			for k, v := range cacheObject.RespHeaders {
				w.Write([]byte(fmt.Sprintf("  < %s: %s\n", k, strings.Join(v, ","))))
			}
			w.Write([]byte("\n"))
			w.Write([]byte(fmt.Sprintf("  Code:                         %d\n", cacheObject.Code)))
			w.Write([]byte(fmt.Sprintf("  OriginCode:                   %d\n", cacheObject.OriginCode)))
			w.Write([]byte(fmt.Sprintf("  ProxyURL:                     %s\n", cacheObject.ProxyURL)))
			w.Write([]byte(fmt.Sprintf("  ReqTime:                      %v\n", cacheObject.ReqTime)))
			w.Write([]byte(fmt.Sprintf("  ReqRespTime:                  %v\n", cacheObject.ReqRespTime)))
			w.Write([]byte(fmt.Sprintf("  RespRespTime:                 %v\n", cacheObject.RespRespTime)))
			w.Write([]byte(fmt.Sprintf("  LastModified:                 %v\n", cacheObject.LastModified)))
			w.Write([]byte(fmt.Sprintf("  HitCount:                     %v\n", cacheObject.HitCount)))
		} else {
			w.Write([]byte("Not Found"))
		}
	} else {
		searchArr, doSearch := qstringOptions["search"]
		cacheNames := d.Stats.CacheNames()
		sort.Strings(cacheNames)
		w.Write([]byte(fmt.Sprintf("Jump to:")))
		for _, cName := range cacheNames {
			w.Write([]byte(fmt.Sprintf("<a href=#%s>%s</a>  ", cName, cName)))
		}
		w.Write([]byte(fmt.Sprintf("\n")))
		for _, cName := range cacheNames {
			if showSelectCache && cName != cacheToDisplay {
				continue
			}
			w.Write([]byte(fmt.Sprintf("<a name=%s></a>", cName)))
			w.Write([]byte(fmt.Sprintf("\n\n<b>*** Cache \"%s\" ***</b>\n", cName)))
			keys := d.Stats.CacheKeys(cName)
			size, _ := d.Stats.CacheSizeByName(cName)
			capacity, _ := d.Stats.CacheCapacityByName(cName)
			w.Write([]byte(fmt.Sprintf("\n  * Size of in use cache:      %s \n", bytefmt.ByteSize(size))))
			w.Write([]byte(fmt.Sprintf("  * Cache capacity:            %s \n", bytefmt.ByteSize(capacity))))
			w.Write([]byte(fmt.Sprintf("  * Number of elements in LRU: %d\n", len(keys))))
			// tail is how much from the top of the LRU to display, top of the LRU is most recently used. head is the other side.
			head := 100
			tail := 100
			tailStr, ok := qstringOptions["tail"]
			if ok {
				tail, err = strconv.Atoi(tailStr[0])
				if err != nil {
					log.Errorf("Error converting tail value to int: %v", err)
				}
			}
			headStr, ok := qstringOptions["head"]
			if ok {
				head, err = strconv.Atoi(headStr[0])
				if err != nil {
					log.Errorf("Error converting head value to int: %v", err)
					head = 100
				}
			}

			w.Write([]byte(fmt.Sprintf("  * Objects in cache sorted by Least Recently Used on top, ")))
			if doSearch {
				w.Write([]byte(fmt.Sprintf("showing only strings matching %s:\n", searchArr[0])))
			} else {
				w.Write([]byte(fmt.Sprintf("showing only first %d and last %d:\n\n", head, tail)))
			}

			w.Write([]byte(fmt.Sprintf("<b>            #    Code      Size                   Age              FreshFor    HitCount      Key</b>\n")))
			for i, key := range keys {
				if (doSearch && !strings.Contains(key, searchArr[0])) || !doSearch && (i >= tail && i < len(keys)-head) {
					continue
				}

				cacheObject, _ := d.Stats.CachePeek(key, cName)
				age := time.Now().Sub(cacheObject.ReqRespTime)
				freshFor := rfc.FreshFor(cacheObject.RespHeaders, cacheObject.RespCacheControl, cacheObject.ReqRespTime, cacheObject.RespRespTime)
				w.Write([]byte(fmt.Sprintf("     %8d%8d%10s%22v%22v%12d      <a href=\"http://%s%s?key=%s&cache=%s\">%s</a>\n",
					i, cacheObject.Code, bytefmt.ByteSize(cacheObject.Size), age, freshFor, cacheObject.HitCount, req.Host, CacheStatsEndpoint, url.QueryEscape(key), cName, key)))
			}

		}
	}

	writeHTMLPageFooter(w)
	clientIP, _ := web.GetClientIPPort(req)
	now := time.Now()
	// TODO add eventId?
	log.EventRaw(atsEventLogStr(now, clientIP, d.Hostname, req.Host, d.Port, "-", d.Scheme, req.URL.String(), req.Method, req.Proto, respCode, now.Sub(reqTime), 0, 0, 0, true, true, getCacheHitStr(true, false), "-", "-", req.UserAgent(), req.Header.Get("X-Money-Trace"), 1))

	return true
}
