package httpsrvr

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 *
 */

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/v8/experimental/traffic_router_golang/availableservers"
	"github.com/apache/trafficcontrol/v8/experimental/traffic_router_golang/cgsrch"
	"github.com/apache/trafficcontrol/v8/experimental/traffic_router_golang/coveragezone"
	"github.com/apache/trafficcontrol/v8/experimental/traffic_router_golang/crconfigregex"
	"github.com/apache/trafficcontrol/v8/experimental/traffic_router_golang/nextcache"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
)

// TODO config
// var DefaultPos = tc.CRConfigLatitudeLongitude{Lat: 39.578968, Lon: -104.934333}
var DefaultPos = tc.CRConfigLatitudeLongitude{Lat: 39.579244, Lon: -104.934282}

// TODO config
const UseXForwardedFor = true

func getHandler(
	regexes crconfigregex.Ths,
	availSrvrs availableservers.AvailableServers,
	cgSrchThs cgsrch.Ths,
	nextCacherThs nextcache.Ths,
	cz coveragezone.CoverageZone,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// host := r.Header.Get("Host")

		// TODO parse subdomains more efficiently
		fqdnParts := strings.Split(r.Host, ".")
		if len(fqdnParts) < 3 {
			fmt.Println("EVENT request '" + r.Host + "' doesn't have enough parts (must be 'subsubdomain.subdomain.domain'), returning 404")
			w.WriteHeader(http.StatusNotFound)
			return
		}

		subsubdomain := fqdnParts[0]
		subdomain := fqdnParts[1]
		domain := strings.Join(fqdnParts[2:len(fqdnParts)-1], ".")

		fmt.Println("DEBUG request '" + r.Host + "' split ssd '" + subsubdomain + "' sd '" + subdomain + "' d '" + domain + "'")

		dsRegexes := (*crconfigregex.Regexes)(regexes.Get())

		dsName, ok := dsRegexes.DeliveryService(domain, subdomain, subsubdomain)
		if !ok {
			fmt.Println("EVENT request '" + r.Host + "' has no match, returning 404")
			w.WriteHeader(http.StatusNotFound)
			return
		}

		fmt.Println("EVENT " + r.RemoteAddr + " request '" + r.Host + "' matched " + string(dsName))

		ipStr := r.Header.Get("X-Forwarded-For")
		if ipStr == "" {
			err := error(nil)
			ipStr, _, err = net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				fmt.Println("ERROR request from" + r.RemoteAddr + " failed to parse: " + err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		fmt.Println("DEBUG " + r.RemoteAddr + " IP '" + ipStr + "'")

		ip := net.ParseIP(ipStr)
		if ip == nil {
			fmt.Println("ERROR request from" + r.RemoteAddr + " IP failed to parse.")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		pos, ok := cz.Get(ip)
		if !ok {
			pos = DefaultPos
			log.Warnln("request from" + r.RemoteAddr + " IP " + ip.String() + " not found, using default")
		}
		log.Infof("LATLON: Request from"+r.RemoteAddr+" IP "+ip.String()+" got %+v\n", pos)

		cgSrch := cgSrchThs.Get()
		cgDat, ok := cgSrch.Nearest(pos.Lat, pos.Lon)
		if !ok {
			fmt.Println("ERROR request from" + r.RemoteAddr + " has no nearest cachegroup (should only happen if there are no cachegroups)")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		cg := tc.CacheGroupName(cgDat.Obj)

		srvrs, err := availSrvrs.Get(dsName, cg)
		if err != nil {
			fmt.Println("EVENT request '" + r.Host + "' with cg '" + string(cg) + "' ds '" + string(dsName) + "' failed to get available servers, returning 404: " + err.Error())
			w.WriteHeader(http.StatusNotFound)
			return
		}

		fmt.Printf("DEBUG GOT AVAILABLE SERVERS %+v\n", srvrs)

		if len(srvrs) == 0 {
			fmt.Println("EVENT request '" + r.Host + "' with cg '" + string(cg) + "' ds '" + string(dsName) + "' no available servers, returning 500")
			w.WriteHeader(http.StatusInternalServerError) // TODO better code?
			return
		}

		nextCacher := nextCacherThs.Get()
		nextSrvrI, ok := nextCacher.NextCache(dsName)
		if !ok {
			// should never happen
			fmt.Println("ERROR request '" + r.Host + "' with cg '" + string(cg) + "' ds '" + string(dsName) + "' not found in Nextcacher, returning 500")
			w.WriteHeader(http.StatusInternalServerError) // TODO better code?
			return
		}

		srvr := srvrs[nextSrvrI%uint64(len(srvrs))]

		newURL := string(srvr) + "." + subdomain + "." + domain + r.URL.Path
		if r.URL.RawQuery != "" {
			newURL += "?" + r.URL.RawQuery
		}

		w.Header().Add(rfc.Location, newURL)
		w.WriteHeader(http.StatusFound)
	}
}

func Start(
	regexes crconfigregex.Ths,
	availableServers availableservers.AvailableServers,
	cgSrch cgsrch.Ths,
	nextCacher nextcache.Ths,
	cz coveragezone.CoverageZone,
	port uint,
) *http.Server {
	srvr := http.Server{}
	srvr.Addr = ":" + strconv.Itoa(int(port))
	srvr.Handler = getHandler(regexes, availableServers, cgSrch, nextCacher, cz)
	go func() {
		err := srvr.ListenAndServe()
		if err != nil {
			fmt.Println("Serving: " + err.Error())
		}
	}()
	return &srvr
}
