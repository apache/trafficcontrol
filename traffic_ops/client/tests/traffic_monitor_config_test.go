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

func TestTMConfig(t *testing.T) {
	resp := fixtures.TrafficMonitorConfig()
	server := validServer(resp)
	defer server.Close()

	var httpClient http.Client
	to := client.Session{
		URL:       server.URL,
		UserAgent: &httpClient,
	}

	Context(t, "Given the need to test a successful Traffic Ops request for TM Config")

	tm, err := to.TrafficMonitorConfigMap("test-cdn")
	if err != nil {
		Error(t, "Should be able to make a request to Traffic Ops")
	} else {
		Success(t, "Should be able to make a request to Traffic Ops")
	}

	Context(t, "Given the need to test a successful Traffic Ops request for TM Config - Traffic Server")

	ts := tm.TrafficServer
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

	m := tm.TrafficMonitor
	if len(m) != 1 {
		Error(t, "Should get back \"1\" Traffic Servers, got: %d", len(m))
	} else {
		Success(t, "Should get back \"1\" Traffic Servers")
	}

	hostName := "traffic-monitor-01"
	if val, ok := m[hostName]; ok {
		Success(t, "Should get back map entry for \"%s\"", hostName)

		if val.Status != "ONLINE" {
			Error(t, "Should get back \"ONLINE\" for \"Status\", got: %s", val.Status)
		} else {
			Success(t, "Should get back \"ONLINE\" for \"Status\"")
		}

	} else {
		Error(t, "Should get back map entry for \"%s\"", hostName)
	}

	Context(t, "Given the need to test a successful Traffic Ops request for TM Config - CacheGroups")

	c := tm.CacheGroup
	if len(c) != 2 {
		Error(t, "Should get back \"2\" Traffic Servers, got: %d", len(c))
	} else {
		Success(t, "Should get back \"2\" Traffic Servers")
	}

	name := "philadelphia"
	if val, ok := c[name]; ok {
		Success(t, "Should get back map entry for \"%s\"", name)

		if val.Coordinates.Latitude != 55 {
			Error(t, "Should get back \"55\" for \"Coordinates.Latitude\", got: %s", val.Coordinates.Latitude)
		} else {
			Success(t, "Should get back \"55\" for \"Coordinates.Latitude\"")
		}

	} else {
		Error(t, "Should get back map entry for \"%s\"", name)
	}

	name = "tr-chicago"
	if val, ok := c[name]; ok {
		Success(t, "Should get back map entry for \"%s\"", name)

		if val.Coordinates.Longitude != 9 {
			Error(t, "Should get back \"9\" for \"Coordinates.Longitude\", got: %s", val.Coordinates.Longitude)
		} else {
			Success(t, "Should get back \"9\" for \"Coordinates.Longitude\"")
		}

	} else {
		Error(t, "Should get back map entry for \"%s\"", name)
	}

	Context(t, "Given the need to test a successful Traffic Ops request for TM Config - Delivery Services")

	ds := tm.DeliveryService
	if len(ds) != 1 {
		Error(t, "Should get back \"1\" TM Delivery Service, got: %d", len(ds))
	} else {
		Success(t, "Should get back \"1\" TM Delivery Service")
	}

	xmlID := "ds-05"
	if val, ok := ds[xmlID]; ok {
		Success(t, "Should get back map entry for \"%s\"", xmlID)

		if val.Status != "REPORTED" {
			Error(t, "Should get back \"REPORTED\" for \"Status\", got: %s", val.Status)
		} else {
			Success(t, "Should get back \"REPORTED\" for \"Status\"")
		}

	} else {
		Error(t, "Should get back map entry for \"%s\"", xmlID)
	}

	Context(t, "Given the need to test a successful Traffic Ops request for TM Config - TM Profiles")

	p := tm.Profile
	if len(p) != 1 {
		Error(t, "Should get back \"1\" TM Profie, got: %d", len(p))
	} else {
		Success(t, "Should get back \"1\" TM Profile")
	}

	name = "tm-123"
	if val, ok := p[name]; ok {
		Success(t, "Should get back map entry for \"%s\"", name)

		if val.Parameters.HealthConnectionTimeout != 2000 {
			Error(t, "Should get back \"2000\" for \"Parameters.HealthConnectionTimeout\", got: %v", val.Parameters.HealthConnectionTimeout)
		} else {
			Success(t, "Should get back \"2000\" for \"Parameters.HealthConnectionTimeout\"")
		}

	} else {
		Error(t, "Should get back map entry for \"%s\"", name)
	}

	Context(t, "Given the need to test a successful Traffic Ops request for TM Config - Config")

	conf := tm.Config
	if _, ok := conf["peers.polling.interval"]; ok {
		if conf["peers.polling.interval"] != float64(1000) {
			Error(t, "Should get back \"1000\" for map entry for \"peers.polling.interval\", got: \"%v\"", conf["peers.polling.interval"])
		} else {
			Success(t, "Should get back \"1000\" for map entry for \"peers.polling.interval\"")
		}
	}
}
