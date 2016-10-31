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

package test

import (
	"net/http"
	"testing"

	"github.com/apache/incubator-trafficcontrol/traffic_ops/client"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/client/fixtures"
	"github.com/jheitz200/test_helper"
)

func TestTRConfig(t *testing.T) {
	resp := fixtures.TrafficRouterConfig()
	server := testHelper.ValidHTTPServer(resp)
	defer server.Close()

	var httpClient http.Client
	to := client.Session{
		URL:       server.URL,
		UserAgent: &httpClient,
	}

	testHelper.Context(t, "Given the need to test a successful Traffic Ops request for TR Config")

	tr, err := to.TrafficRouterConfigMap("title-vi")
	if err != nil {
		testHelper.Error(t, "Should be able to make a request to Traffic Ops")
	} else {
		testHelper.Success(t, "Should be able to make a request to Traffic Ops")
	}

	testHelper.Context(t, "Given the need to test a successful Traffic Ops request for TR Config - Traffic Server")

	ts := tr.TrafficServer
	if len(ts) != 2 {
		testHelper.Error(t, "Should get back \"2\" Traffic Servers, got: %d", len(ts))
	} else {
		testHelper.Success(t, "Should get back \"2\" Traffic Servers")
	}

	hashID := "tr-chi-05"
	if val, ok := ts[hashID]; ok {
		testHelper.Success(t, "Should get back map entry for \"%s\"", hashID)

		if val.IP != "10.10.10.10" {
			testHelper.Error(t, "Should get back \"10.10.10.10\" for \"IP\", got: %s", val.IP)
		} else {
			testHelper.Success(t, "Should get back \"10.10.10.10\" for \"IP\"")
		}

	} else {
		testHelper.Error(t, "Should get back map entry for \"%s\"", hashID)
	}

	hashID = "edge-test-01"
	if val, ok := ts[hashID]; ok {
		testHelper.Success(t, "Should get back map entry for for \"%s\"", hashID)

		if val.Type != "EDGE" {
			testHelper.Error(t, "Should get back \"EDGE\" for \"%s\" \"Type\", got: %s", hashID, ts[hashID].Type)
		} else {
			testHelper.Success(t, "Should get back \"EDGE\" for \"%s\" \"Type\"", hashID)
		}

	} else {
		testHelper.Error(t, "Should get back map entry for \"%s\"", hashID)
	}

	testHelper.Context(t, "Given the need to test a successful Traffic Ops request for TR Config - Traffic Monitor")

	tm := tr.TrafficMonitor
	if len(tm) != 1 {
		testHelper.Error(t, "Should get back \"1\" Traffic Servers, got: %d", len(tm))
	} else {
		testHelper.Success(t, "Should get back \"1\" Traffic Servers")
	}

	hostname := "traffic-monitor-01"
	if val, ok := tm[hostname]; ok {
		testHelper.Success(t, "Should get back map entry for \"%s\"", hostname)

		if val.Profile != "tr-123" {
			testHelper.Error(t, "Should get back \"tr-123\" for \"Profile\", got: %s", val.Profile)
		} else {
			testHelper.Success(t, "Should get back \"tr-123\" for \"Profile\"")
		}

	} else {
		testHelper.Error(t, "Should get back map entry for \"%s\"", hostname)
	}

	testHelper.Context(t, "Given the need to test a successful Traffic Ops request for TR Config - Traffic Router")

	r := tr.TrafficRouter
	if len(r) != 1 {
		testHelper.Error(t, "Should get back \"1\" Traffic Servers, got: %d", len(r))
	} else {
		testHelper.Success(t, "Should get back \"1\" Traffic Servers")
	}

	fqdn := "tr-01@ga.atlanta.kabletown.net"
	if val, ok := r[fqdn]; ok {
		testHelper.Success(t, "Should get back map entry for \"%s\"", fqdn)

		if val.Location != "tr-chicago" {
			testHelper.Error(t, "Should get back \"tr-chicago\" for \"Location\", got: %s", val.Location)
		} else {
			testHelper.Success(t, "Should get back \"tr-chicago\" for \"Location\"")
		}

	} else {
		testHelper.Error(t, "Should get back map entry for \"%s\"", fqdn)
	}

	testHelper.Context(t, "Given the need to test a successful Traffic Ops request for TR Config - CacheGroups")

	c := tr.CacheGroup
	if len(c) != 2 {
		testHelper.Error(t, "Should get back \"2\" Traffic Servers, got: %d", len(c))
	} else {
		testHelper.Success(t, "Should get back \"2\" Traffic Servers")
	}

	name := "philadelphia"
	if val, ok := c[name]; ok {
		testHelper.Success(t, "Should get back map entry for \"%s\"", name)

		if val.Coordinates.Latitude != 99 {
			testHelper.Error(t, "Should get back \"99\" for \"Coordinates.Latitude\", got: %v", val.Coordinates.Latitude)
		} else {
			testHelper.Success(t, "Should get back \"99\" for \"Coordinates.Latitude\"")
		}

	} else {
		testHelper.Error(t, "Should get back map entry for \"%s\"", name)
	}

	name = "tr-chicago"
	if val, ok := c[name]; ok {
		testHelper.Success(t, "Should get back map entry for \"%s\"", name)

		if val.Coordinates.Longitude != 9 {
			testHelper.Error(t, "Should get back \"9\" for \"Coordinates.Longitude\", got: %v", val.Coordinates.Longitude)
		} else {
			testHelper.Success(t, "Should get back \"9\" for \"Coordinates.Longitude\"")
		}

	} else {
		testHelper.Error(t, "Should get back map entry for \"%s\"", name)
	}

	testHelper.Context(t, "Given the need to test a successful Traffic Ops request for TR Config - Delivery Services")

	ds := tr.DeliveryService
	if len(ds) != 1 {
		testHelper.Error(t, "Should get back \"1\" TR Delivery Service, got: %d", len(ds))
	} else {
		testHelper.Success(t, "Should get back \"1\" TR Delivery Service")
	}

	xmlID := "ds-06"
	if val, ok := ds[xmlID]; ok {
		testHelper.Success(t, "Should get back map entry for \"%s\"", xmlID)

		if val.TTL != 3600 {
			testHelper.Error(t, "Should get back \"3600\" for \"TTL\", got: %d", val.TTL)
		} else {
			testHelper.Success(t, "Should get back \"3600\" for \"TTL\"")
		}

	} else {
		testHelper.Error(t, "Should get back map entry for \"%s\"", xmlID)
	}

	testHelper.Context(t, "Given the need to test a successful Traffic Ops request for TR Config - Config")

	conf := tr.Config
	if _, ok := conf["peers.polling.interval"]; ok {
		if conf["peers.polling.interval"] != float64(1000) {
			testHelper.Error(t, "Should get back \"1000\" for map entry for \"peers.polling.interval\", got: \"%v\"", conf["peers.polling.interval"])
		} else {
			testHelper.Success(t, "Should get back \"1000\" for map entry for \"peers.polling.interval\"")
		}
	}

	testHelper.Context(t, "Given the need to test a successful Traffic Ops request for TR Config - Stats")

	stats := tr.Stat
	if _, ok := stats["cdnName"]; ok {
		if stats["cdnName"] != "test-cdn" {
			testHelper.Error(t, "Should get back \"test-cdn\" for map entry for \"cdnName\", got: \"%s\"", stats["cdnName"])
		} else {
			testHelper.Success(t, "Should get back \"test-cdn\" for map entry for \"cdnName\"")
		}
	}
}
