package _integration

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
	"testing"
)

func TestClient(t *testing.T) {
	if actual, err := TMClient.CacheCount(); err != nil {
		t.Errorf("client CacheCount error expected nil, actual %v\n", err)
	} else if actual <= 0 {
		t.Errorf("client CacheCount expected > 0 actual %v\n", actual)
	}

	if actual, err := TMClient.CacheAvailableCount(); err != nil {
		t.Errorf("client CacheAvailableCount error expected nil, actual %v\n", err)
	} else if actual <= 0 {
		t.Errorf("client CacheAvailableCount expected > 0 actual %v\n", actual)
	}

	if actual, err := TMClient.CacheDownCount(); err != nil {
		t.Errorf("client CacheDownCount error expected nil, actual %v\n", err)
	} else if actual < 0 {
		t.Errorf("client CacheDownCount expected >= 0 actual %v\n", actual)
	}

	if actual, err := TMClient.Version(); err != nil {
		t.Errorf("client Version error expected nil, actual %v\n", err)
	} else if actual == "" {
		t.Errorf("client Version expected not empty, actual empty\n")
	}

	if actual, err := TMClient.TrafficOpsURI(); err != nil {
		t.Errorf("client TrafficOpsURI error expected nil, actual %v\n", err)
	} else if actual == "" {
		t.Errorf("client TrafficOpsURI expected not empty, actual empty\n")
	}

	if actual, err := TMClient.BandwidthKBPS(); err != nil {
		t.Errorf("client BandwidthKBPS error expected nil, actual %v\n", err)
	} else if actual < 0 {
		t.Errorf("client BandwidthKBPS expected >=0 , actual %v\n", actual)
	}

	if actual, err := TMClient.BandwidthCapacityKBPS(); err != nil {
		t.Errorf("client BandwidthCapacityKBPS error expected nil, actual %v\n", err)
	} else if actual <= 0 {
		t.Errorf("client BandwidthCapacityKBPS expected >0 , actual %v\n", actual)
	}

	if actual, err := TMClient.CacheStatuses(); err != nil {
		t.Errorf("client CacheStatuses error expected nil, actual %v\n", err)
	} else if len(actual) == 0 {
		t.Errorf("client len(CacheStatuses) expected >0 , actual %v\n", actual)
	}

	if actual, err := TMClient.MonitorConfig(); err != nil {
		t.Errorf("client MonitorConfig error expected nil, actual %v\n", err)
	} else if len(actual.TrafficServer) == 0 {
		t.Errorf("client len(TrafficMonitorConfig.TrafficServers) expected not empty, actual %v\n", actual)
	}

	if actual, err := TMClient.CRConfigHistory(); err != nil {
		t.Errorf("client CRConfigHistory error expected nil, actual %v\n", err)
	} else if len(actual) == 0 {
		t.Errorf("client len(CRConfigHistory) expected !=0, actual %v\n", actual)
	}

	if actual, err := TMClient.EventLog(); err != nil {
		t.Errorf("client EventLog error expected nil, actual %v\n", err)
	} else if len(actual.Events) == 0 {
		t.Errorf("client len(EventLog.Events) expected !=0, actual %v\n", actual)
	}

	if actual, err := TMClient.CacheStats(); err != nil {
		t.Errorf("client CacheStats error expected nil, actual %v\n", err)
	} else if len(actual.Caches) == 0 {
		t.Errorf("client len(CacheStats.Caches) expected !=0, actual %v\n", actual)
	}

	if actual, err := TMClient.CacheStatsNew(); err != nil {
		t.Errorf("client CacheStatsNew error expected nil, actual %v\n", err)
	} else if len(actual.Caches) == 0 {
		t.Errorf("client len(CacheStatsNew.Caches) expected !=0, actual %v\n", actual)
	}

	if actual, err := TMClient.DSStats(); err != nil {
		t.Errorf("client DSStats error expected nil, actual %v\n", err)
	} else if len(actual.DeliveryService) == 0 {
		t.Errorf("client len(DSStats.DeliveryService) expected !=0, actual %v\n", actual)
	}

	if actual, err := TMClient.CRStates(false); err != nil {
		t.Errorf("client CRStates error expected nil, actual %v\n", err)
	} else if len(actual.Caches) == 0 {
		t.Errorf("client len(CRStates.Caches) expected !=0, actual %v\n", actual)
	}

	if actual, err := TMClient.CRConfig(); err != nil {
		t.Errorf("client CRConfig error expected nil, actual %v\n", err)
	} else if len(actual.ContentServers) == 0 {
		t.Errorf("client len(CRConfig.ContentServers) expected !=0, actual %v\n", actual)
	}
}
