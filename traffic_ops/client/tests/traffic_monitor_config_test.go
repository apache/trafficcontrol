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

func TestTMConfig(t *testing.T) {
	resp := fixtures.TrafficMonitorConfig()
	server := test.ValidHTTPServer(resp)
	defer server.Close()

	var httpClient http.Client
	to := client.Session{
		URL:       server.URL,
		UserAgent: &httpClient,
	}

	test.Context(t, "Given the need to test a successful Traffic Ops request for TM Config")

	tm, err := to.TrafficMonitorConfigMap("test-cdn")
	if err != nil {
		test.Error(t, "Should be able to make a request to Traffic Ops")
	} else {
		test.Success(t, "Should be able to make a request to Traffic Ops")
	}

	test.Context(t, "Given the need to test a successful Traffic Ops request for TM Config - Traffic Server")

	ts := tm.TrafficServer
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

	m := tm.TrafficMonitor
	if len(m) != 1 {
		test.Error(t, "Should get back \"1\" Traffic Servers, got: %d", len(m))
	} else {
		test.Success(t, "Should get back \"1\" Traffic Servers")
	}

	hostName := "traffic-monitor-01"
	if val, ok := m[hostName]; ok {
		test.Success(t, "Should get back map entry for \"%s\"", hostName)

		if val.Status != "ONLINE" {
			test.Error(t, "Should get back \"ONLINE\" for \"Status\", got: %s", val.Status)
		} else {
			test.Success(t, "Should get back \"ONLINE\" for \"Status\"")
		}

	} else {
		test.Error(t, "Should get back map entry for \"%s\"", hostName)
	}

	test.Context(t, "Given the need to test a successful Traffic Ops request for TM Config - CacheGroups")

	c := tm.CacheGroup
	if len(c) != 2 {
		test.Error(t, "Should get back \"2\" Traffic Servers, got: %d", len(c))
	} else {
		test.Success(t, "Should get back \"2\" Traffic Servers")
	}

	name := "philadelphia"
	if val, ok := c[name]; ok {
		test.Success(t, "Should get back map entry for \"%s\"", name)

		if val.Coordinates.Latitude != 55 {
			test.Error(t, "Should get back \"55\" for \"Coordinates.Latitude\", got: %s", val.Coordinates.Latitude)
		} else {
			test.Success(t, "Should get back \"55\" for \"Coordinates.Latitude\"")
		}

	} else {
		test.Error(t, "Should get back map entry for \"%s\"", name)
	}

	name = "tr-chicago"
	if val, ok := c[name]; ok {
		test.Success(t, "Should get back map entry for \"%s\"", name)

		if val.Coordinates.Longitude != 9 {
			test.Error(t, "Should get back \"9\" for \"Coordinates.Longitude\", got: %s", val.Coordinates.Longitude)
		} else {
			test.Success(t, "Should get back \"9\" for \"Coordinates.Longitude\"")
		}

	} else {
		test.Error(t, "Should get back map entry for \"%s\"", name)
	}

	test.Context(t, "Given the need to test a successful Traffic Ops request for TM Config - Delivery Services")

	ds := tm.DeliveryService
	if len(ds) != 1 {
		test.Error(t, "Should get back \"1\" TM Delivery Service, got: %d", len(ds))
	} else {
		test.Success(t, "Should get back \"1\" TM Delivery Service")
	}

	xmlID := "ds-05"
	if val, ok := ds[xmlID]; ok {
		test.Success(t, "Should get back map entry for \"%s\"", xmlID)

		if val.Status != "REPORTED" {
			test.Error(t, "Should get back \"REPORTED\" for \"Status\", got: %s", val.Status)
		} else {
			test.Success(t, "Should get back \"REPORTED\" for \"Status\"")
		}

	} else {
		test.Error(t, "Should get back map entry for \"%s\"", xmlID)
	}

	test.Context(t, "Given the need to test a successful Traffic Ops request for TM Config - TM Profiles")

	p := tm.Profile
	if len(p) != 1 {
		test.Error(t, "Should get back \"1\" TM Profie, got: %d", len(p))
	} else {
		test.Success(t, "Should get back \"1\" TM Profile")
	}

	name = "tm-123"
	if val, ok := p[name]; ok {
		test.Success(t, "Should get back map entry for \"%s\"", name)

		if val.Parameters.HealthConnectionTimeout != 2000 {
			test.Error(t, "Should get back \"2000\" for \"Parameters.HealthConnectionTimeout\", got: %v", val.Parameters.HealthConnectionTimeout)
		} else {
			test.Success(t, "Should get back \"2000\" for \"Parameters.HealthConnectionTimeout\"")
		}

	} else {
		test.Error(t, "Should get back map entry for \"%s\"", name)
	}

	test.Context(t, "Given the need to test a successful Traffic Ops request for TM Config - Config")

	conf := tm.Config
	if _, ok := conf["peers.polling.interval"]; ok {
		if conf["peers.polling.interval"] != float64(1000) {
			test.Error(t, "Should get back \"1000\" for map entry for \"peers.polling.interval\", got: \"%v\"", conf["peers.polling.interval"])
		} else {
			test.Success(t, "Should get back \"1000\" for map entry for \"peers.polling.interval\"")
		}
	}
}
