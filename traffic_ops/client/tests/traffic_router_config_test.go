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

	"github.com/jheitz200/traffic_control/traffic_ops/client"
	"github.com/jheitz200/traffic_control/traffic_ops/client/fixtures"
)

func TestTRConfig(t *testing.T) {
	resp := fixtures.TrafficRouterConfig()
	server := validServer(resp)
	defer server.Close()

	var httpClient http.Client
	to := client.Session{
		URL:       server.URL,
		UserAgent: &httpClient,
	}

	Context(t, "Given the need to test a successful Traffic Ops request for TR Config")

	tr, err := to.TrafficRouterConfigMap("title-vi")
	if err != nil {
		Error(t, "Should be able to make a request to Traffic Ops")
	} else {
		Success(t, "Should be able to make a request to Traffic Ops")
	}

	Context(t, "Given the need to test a successful Traffic Ops request for TR Config - Traffic Server")

	ts := tr.TrafficServer
	if len(ts) != 2 {
		Error(t, "Should get back \"2\" Traffic Servers, got: %d", len(ts))
	} else {
		Success(t, "Should get back \"2\" Traffic Servers")
	}

	hashID := "tr-chi-05"
	if val, ok := ts[hashID]; ok {
		Success(t, "Should get back map entry for \"%s\"", hashID)

		if val.IP != "10.10.10.10" {
			Error(t, "Should get back \"10.10.10.10\" for \"IP\", got: %s", val.IP)
		} else {
			Success(t, "Should get back \"10.10.10.10\" for \"IP\"")
		}

	} else {
		Error(t, "Should get back map entry for \"%s\"", hashID)
	}

	hashID = "edge-test-01"
	if val, ok := ts[hashID]; ok {
		Success(t, "Should get back map entry for for \"%s\"", hashID)

		if val.Type != "EDGE" {
			Error(t, "Should get back \"EDGE\" for \"%s\" \"Type\", got: %s", hashID, ts[hashID].Type)
		} else {
			Success(t, "Should get back \"EDGE\" for \"%s\" \"Type\"", hashID)
		}

	} else {
		Error(t, "Should get back map entry for \"%s\"", hashID)
	}

	Context(t, "Given the need to test a successful Traffic Ops request for TR Config - Traffic Monitor")

	tm := tr.TrafficMonitor
	if len(tm) != 1 {
		Error(t, "Should get back \"1\" Traffic Servers, got: %d", len(tm))
	} else {
		Success(t, "Should get back \"1\" Traffic Servers")
	}

	hostname := "traffic-monitor-01"
	if val, ok := tm[hostname]; ok {
		Success(t, "Should get back map entry for \"%s\"", hostname)

		if val.Profile != "tr-123" {
			Error(t, "Should get back \"tr-123\" for \"Profile\", got: %s", val.Profile)
		} else {
			Success(t, "Should get back \"tr-123\" for \"Profile\"")
		}

	} else {
		Error(t, "Should get back map entry for \"%s\"", hostname)
	}

	Context(t, "Given the need to test a successful Traffic Ops request for TR Config - Traffic Router")

	r := tr.TrafficRouter
	if len(r) != 1 {
		Error(t, "Should get back \"1\" Traffic Servers, got: %d", len(r))
	} else {
		Success(t, "Should get back \"1\" Traffic Servers")
	}

	fqdn := "tr-01@ga.atlanta.kabletown.net"
	if val, ok := r[fqdn]; ok {
		Success(t, "Should get back map entry for \"%s\"", fqdn)

		if val.Location != "tr-chicago" {
			Error(t, "Should get back \"tr-chicago\" for \"Location\", got: %s", val.Location)
		} else {
			Success(t, "Should get back \"tr-chicago\" for \"Location\"")
		}

	} else {
		Error(t, "Should get back map entry for \"%s\"", fqdn)
	}

	Context(t, "Given the need to test a successful Traffic Ops request for TR Config - CacheGroups")

	c := tr.CacheGroup
	if len(c) != 2 {
		Error(t, "Should get back \"2\" Traffic Servers, got: %d", len(c))
	} else {
		Success(t, "Should get back \"2\" Traffic Servers")
	}

	name := "philadelphia"
	if val, ok := c[name]; ok {
		Success(t, "Should get back map entry for \"%s\"", name)

		if val.Coordinates.Latitude != 99 {
			Error(t, "Should get back \"99\" for \"Coordinates.Latitude\", got: %v", val.Coordinates.Latitude)
		} else {
			Success(t, "Should get back \"99\" for \"Coordinates.Latitude\"")
		}

	} else {
		Error(t, "Should get back map entry for \"%s\"", name)
	}

	name = "tr-chicago"
	if val, ok := c[name]; ok {
		Success(t, "Should get back map entry for \"%s\"", name)

		if val.Coordinates.Longitude != 9 {
			Error(t, "Should get back \"9\" for \"Coordinates.Longitude\", got: %v", val.Coordinates.Longitude)
		} else {
			Success(t, "Should get back \"9\" for \"Coordinates.Longitude\"")
		}

	} else {
		Error(t, "Should get back map entry for \"%s\"", name)
	}

	Context(t, "Given the need to test a successful Traffic Ops request for TR Config - Delivery Services")

	ds := tr.DeliveryService
	if len(ds) != 1 {
		Error(t, "Should get back \"1\" TR Delivery Service, got: %d", len(ds))
	} else {
		Success(t, "Should get back \"1\" TR Delivery Service")
	}

	xmlID := "ds-06"
	if val, ok := ds[xmlID]; ok {
		Success(t, "Should get back map entry for \"%s\"", xmlID)

		if val.TTL != 3600 {
			Error(t, "Should get back \"3600\" for \"TTL\", got: %d", val.TTL)
		} else {
			Success(t, "Should get back \"3600\" for \"TTL\"")
		}

	} else {
		Error(t, "Should get back map entry for \"%s\"", xmlID)
	}

	Context(t, "Given the need to test a successful Traffic Ops request for TR Config - Config")

	conf := tr.Config
	if _, ok := conf["peers.polling.interval"]; ok {
		if conf["peers.polling.interval"] != float64(1000) {
			Error(t, "Should get back \"1000\" for map entry for \"peers.polling.interval\", got: \"%v\"", conf["peers.polling.interval"])
		} else {
			Success(t, "Should get back \"1000\" for map entry for \"peers.polling.interval\"")
		}
	}

	Context(t, "Given the need to test a successful Traffic Ops request for TR Config - Stats")

	stats := tr.Stat
	if _, ok := stats["cdnName"]; ok {
		if stats["cdnName"] != "test-cdn" {
			Error(t, "Should get back \"test-cdn\" for map entry for \"cdnName\", got: \"%s\"", stats["cdnName"])
		} else {
			Success(t, "Should get back \"test-cdn\" for map entry for \"cdnName\"")
		}
	}
}
