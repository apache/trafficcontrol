package towrap

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
 */

import (
	"testing"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

func TestMonitorConfigValid(t *testing.T) {
	mc := (*tc.LegacyTrafficMonitorConfigMap)(nil)
	if MonitorConfigValid(mc) == nil {
		t.Errorf("MonitorCopnfigValid(nil) expected: error, actual: nil")
	}
	mc = &tc.LegacyTrafficMonitorConfigMap{}
	if MonitorConfigValid(mc) == nil {
		t.Errorf("MonitorConfigValid({}) expected: error, actual: nil")
	}

	validMC := &tc.LegacyTrafficMonitorConfigMap{
		TrafficServer:   map[string]tc.LegacyTrafficServer{"a": {}},
		CacheGroup:      map[string]tc.TMCacheGroup{"a": {}},
		TrafficMonitor:  map[string]tc.TrafficMonitor{"a": {}},
		DeliveryService: map[string]tc.TMDeliveryService{"a": {}},
		Profile:         map[string]tc.TMProfile{"a": {}},
		Config: map[string]interface{}{
			"peers.polling.interval":  42.0,
			"health.polling.interval": 24.0,
		},
	}
	if err := MonitorConfigValid(validMC); err != nil {
		t.Errorf("MonitorConfigValid(%++v) expected: nil, actual: %+v", validMC, err)
	}
}
