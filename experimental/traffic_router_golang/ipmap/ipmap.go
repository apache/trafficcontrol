package ipmap

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
	"errors"
	"net"

	"github.com/apache/trafficcontrol/v8/experimental/traffic_router_golang/coveragezone"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
)

// TODO implement, with Coverage Zone File + Maxmind DB

var ErrNotFound = errors.New("not found")

type LatLon struct {
	Lat float64
	Lon float64
}

// DummyLocations returns dummy locations for IPs.
// TODO remove when real geo lookup is implemented
func DummyLocations() map[string]LatLon {
	return map[string]LatLon{
		"127.0.0.1":   {39.579244, -104.934282},
		"192.168.0.1": {39.579244, -104.934282},
		"::1":         {39.579244, -104.934282},
	}
}

func New() coveragezone.CoverageZone {
	return &ipmap{}
}

type ipmap struct{}

// Get takes an IP and returns the Latitude and Longitude.
func (i *ipmap) Get(ip net.IP) (tc.CRConfigLatitudeLongitude, bool) {
	// TODO: get IP via req.RemoteAddr => net.SplitHostPort => net.ParseIP
	locations := DummyLocations()
	ll, ok := locations[ip.String()]
	if !ok {
		return tc.CRConfigLatitudeLongitude{}, false
	}
	return tc.CRConfigLatitudeLongitude{Lat: ll.Lat, Lon: ll.Lon}, true
}
