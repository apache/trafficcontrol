package totest

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

	"github.com/apache/trafficcontrol/lib/go-util/assert"
	toclient "github.com/apache/trafficcontrol/traffic_ops/v5-client"
)

func CreateTestServerChecks(t *testing.T, cl *toclient.Session, td TrafficControl) {
	for _, servercheck := range td.Serverchecks {
		resp, _, err := cl.InsertServerCheckStatus(servercheck, toclient.RequestOptions{})
		assert.RequireNoError(t, err, "Could not insert Servercheck: %v - alerts: %+v", err, resp.Alerts)
	}
}

// Need to define no-op function as TCObj interface expects a delete function
// There is no delete path for serverchecks
func DeleteTestServerChecks(*testing.T, *toclient.Session) {
	return
}
