/*
   Copyright 2015 Comcast Cable Communications Management, LLC

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

package test

import (
	"net/http"
	"testing"

	"github.com/Comcast/test_helper"
	"github.com/jheitz200/traffic_control/traffic_ops/client"
	"github.com/jheitz200/traffic_control/traffic_ops/client/fixtures"
)

func TestTRConfig(t *testing.T) {
	resp := fixtures.TrafficRouterConfig()
	server := test.ValidHTTPServer(resp)
	defer server.Close()

	var httpClient http.Client
	to := client.Session{
		URL:       server.URL,
		UserAgent: &httpClient,
	}

	test.Context(t, "Given the need to test a successful Traffic Ops request for TR Config")

	tr, err := to.TrafficRouterConfigMap("title-vi")
	if err != nil {
		test.Error(t, "Should be able to make a request to Traffic Ops")
	} else {
		test.Success(t, "Should be able to make a request to Traffic Ops")
	}

	test.Context(t, "Given the need to test a successful Traffic Ops request for TR Config - Traffic Server")

	ts := tr.TrafficServer
	if len(ts) != 2 {
		test.Error(t, "Should get back \"2\" Traffic Servers, got: %d", len(ts))
	} else {
		test.Success(t, "Should get back \"2\" Traffic Servers")
	}

	hashID := "tr-chi-05"
	if val, ok := ts[hashID]; ok {
		test.Success(t, "Should get back map entry for \"%s\"", hashID)

		if val.IP != "10.10.10.10" {
			test.Error(t, "Should get back \"10.10.10.10\" for \"IP\", got: %s", val.IP)
		} else {
			test.Success(t, "Should get back \"10.10.10.10\" for \"IP\"")
		}

	} else {
		test.Error(t, "Should get back map entry for \"%s\"", hashID)
	}

	hashID = "edge-test-01"
	if val, ok := ts[hashID]; ok {
		test.Success(t, "Should get back map entry for for \"%s\"", hashID)

		if val.Type != "EDGE" {
			test.Error(t, "Should get back \"EDGE\" for \"%s\" \"Type\", got: %s", hashID, ts[hashID].Type)
		} else {
			test.Success(t, "Should get back \"EDGE\" for \"%s\" \"Type\"", hashID)
		}

	} else {
		test.Error(t, "Should get back map entry for \"%s\"", hashID)
	}

	test.Context(t, "Given the need to test a successful Traffic Ops request for TR Config - Traffic Monitor")

	tm := tr.TrafficMonitor
	if len(tm) != 1 {
		test.Error(t, "Should get back \"1\" Traffic Servers, got: %d", len(tm))
	} else {
		test.Success(t, "Should get back \"1\" Traffic Servers")
	}

	hostname := "traffic-monitor-01"
	if val, ok := tm[hostname]; ok {
		test.Success(t, "Should get back map entry for \"%s\"", hostname)

		if val.Profile != "tr-123" {
			test.Error(t, "Should get back \"tr-123\" for \"Profile\", got: %s", val.Profile)
		} else {
			test.Success(t, "Should get back \"tr-123\" for \"Profile\"")
		}

	} else {
		test.Error(t, "Should get back map entry for \"%s\"", hostname)
	}

	test.Context(t, "Given the need to test a successful Traffic Ops request for TR Config - Traffic Router")

	r := tr.TrafficRouter
	if len(r) != 1 {
		test.Error(t, "Should get back \"1\" Traffic Servers, got: %d", len(r))
	} else {
		test.Success(t, "Should get back \"1\" Traffic Servers")
	}

	fqdn := "tr-01@ga.atlanta.kabletown.net"
	if val, ok := r[fqdn]; ok {
		test.Success(t, "Should get back map entry for \"%s\"", fqdn)

		if val.Location != "tr-chicago" {
			test.Error(t, "Should get back \"tr-chicago\" for \"Location\", got: %s", val.Location)
		} else {
			test.Success(t, "Should get back \"tr-chicago\" for \"Location\"")
		}

	} else {
		test.Error(t, "Should get back map entry for \"%s\"", fqdn)
	}

	test.Context(t, "Given the need to test a successful Traffic Ops request for TR Config - CacheGroups")

	c := tr.CacheGroup
	if len(c) != 2 {
		test.Error(t, "Should get back \"2\" Traffic Servers, got: %d", len(c))
	} else {
		test.Success(t, "Should get back \"2\" Traffic Servers")
	}

	name := "philadelphia"
	if val, ok := c[name]; ok {
		test.Success(t, "Should get back map entry for \"%s\"", name)

		if val.Coordinates.Latitude != 99 {
			test.Error(t, "Should get back \"99\" for \"Coordinates.Latitude\", got: %v", val.Coordinates.Latitude)
		} else {
			test.Success(t, "Should get back \"99\" for \"Coordinates.Latitude\"")
		}

	} else {
		test.Error(t, "Should get back map entry for \"%s\"", name)
	}

	name = "tr-chicago"
	if val, ok := c[name]; ok {
		test.Success(t, "Should get back map entry for \"%s\"", name)

		if val.Coordinates.Longitude != 9 {
			test.Error(t, "Should get back \"9\" for \"Coordinates.Longitude\", got: %v", val.Coordinates.Longitude)
		} else {
			test.Success(t, "Should get back \"9\" for \"Coordinates.Longitude\"")
		}

	} else {
		test.Error(t, "Should get back map entry for \"%s\"", name)
	}

	test.Context(t, "Given the need to test a successful Traffic Ops request for TR Config - Delivery Services")

	ds := tr.DeliveryService
	if len(ds) != 1 {
		test.Error(t, "Should get back \"1\" TR Delivery Service, got: %d", len(ds))
	} else {
		test.Success(t, "Should get back \"1\" TR Delivery Service")
	}

	xmlID := "ds-06"
	if val, ok := ds[xmlID]; ok {
		test.Success(t, "Should get back map entry for \"%s\"", xmlID)

		if val.TTL != 3600 {
			test.Error(t, "Should get back \"3600\" for \"TTL\", got: %d", val.TTL)
		} else {
			test.Success(t, "Should get back \"3600\" for \"TTL\"")
		}

	} else {
		test.Error(t, "Should get back map entry for \"%s\"", xmlID)
	}

	test.Context(t, "Given the need to test a successful Traffic Ops request for TR Config - Config")

	conf := tr.Config
	if _, ok := conf["peers.polling.interval"]; ok {
		if conf["peers.polling.interval"] != float64(1000) {
			test.Error(t, "Should get back \"1000\" for map entry for \"peers.polling.interval\", got: \"%v\"", conf["peers.polling.interval"])
		} else {
			test.Success(t, "Should get back \"1000\" for map entry for \"peers.polling.interval\"")
		}
	}

	test.Context(t, "Given the need to test a successful Traffic Ops request for TR Config - Stats")

	stats := tr.Stat
	if _, ok := stats["cdnName"]; ok {
		if stats["cdnName"] != "test-cdn" {
			test.Error(t, "Should get back \"test-cdn\" for map entry for \"cdnName\", got: \"%s\"", stats["cdnName"])
		} else {
			test.Success(t, "Should get back \"test-cdn\" for map entry for \"cdnName\"")
		}
	}
}
