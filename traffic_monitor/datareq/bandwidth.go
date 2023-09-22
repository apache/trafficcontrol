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

package datareq

import (
	"fmt"

	"github.com/apache/trafficcontrol/v8/traffic_monitor/ds"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/threadsafe"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/todata"
)

func srvAPIBandwidthKbps(toData todata.TODataThreadsafe, lastStats threadsafe.LastStats) []byte {
	kbpsStats := lastStats.Get()
	sum := float64(0.0)
	for _, data := range kbpsStats.Caches {
		sum += data.Bytes.PerSec / ds.BytesPerKilobit
	}
	return []byte(fmt.Sprintf("%f", sum))
}
