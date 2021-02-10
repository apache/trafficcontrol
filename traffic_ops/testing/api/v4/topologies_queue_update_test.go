package v4

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
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

type topologiesQueueUpdateTestCase struct {
	Description string
	tc.TopologiesQueueUpdateRequest
}

func TestTopologiesQueueUpdate(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Users, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, DeliveryServices}, func() {
		const topologyName = "mso-topology"
		cdnID, dsID := getCDNIDAndDSID(t)
		InvalidCDNIDIsRejected(t, topologyName)
		InvalidActionIsRejected(t, topologyName, cdnID)
		NonexistentTopologyIsRejected(t, cdnID)
		UpdatesAreQueued(t, topologyName, cdnID, dsID)
	})
}

func getCDNIDAndDSID(t *testing.T) (int64, int) {
	xmlID := "ds-top"
	params := url.Values{}
	params.Set("xmlId", xmlID)
	dses, _, err := TOSession.GetDeliveryServices(nil, params)
	if err != nil {
		t.Fatalf("unable to get deliveryservice %s: %s", xmlID, err)
	}
	if len(dses) < 1 {
		t.Fatalf("deliveryservice with xmlId %s not found!", xmlID)
	}
	ds := dses[0]
	return int64(*ds.CDNID), *ds.ID
}

func InvalidCDNIDIsRejected(t *testing.T, topologyName tc.TopologyName) {
	testCase := topologiesQueueUpdateTestCase{
		Description:                  "invalid CDN ID",
		TopologiesQueueUpdateRequest: tc.TopologiesQueueUpdateRequest{Action: "queue", CDNID: -1},
	}
	_, reqInf, _ := TOSession.TopologiesQueueUpdate(topologyName, testCase.TopologiesQueueUpdateRequest)
	if reqInf.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status code %d for request with %s, got status code %d", http.StatusBadRequest, testCase.Description, reqInf.StatusCode)
	}
}

func InvalidActionIsRejected(t *testing.T, topologyName tc.TopologyName, cdnID int64) {
	testCase := topologiesQueueUpdateTestCase{
		Description:                  "invalid update action",
		TopologiesQueueUpdateRequest: tc.TopologiesQueueUpdateRequest{Action: "requeue", CDNID: cdnID},
	}
	_, reqInf, _ := TOSession.TopologiesQueueUpdate(topologyName, testCase.TopologiesQueueUpdateRequest)
	if reqInf.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status code %d for request with %s, got status code %d", http.StatusBadRequest, testCase.Description, reqInf.StatusCode)
	}
}

func NonexistentTopologyIsRejected(t *testing.T, cdnID int64) {
	const topologyName = "nonexistent"
	testCase := topologiesQueueUpdateTestCase{
		Description:                  "nonexistent topology",
		TopologiesQueueUpdateRequest: tc.TopologiesQueueUpdateRequest{Action: "queue", CDNID: cdnID},
	}
	_, reqInf, _ := TOSession.TopologiesQueueUpdate(topologyName, testCase.TopologiesQueueUpdateRequest)
	if reqInf.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status code %d for request with %s, got status code %d", http.StatusBadRequest, testCase.Description, reqInf.StatusCode)
	}
}

func UpdatesAreQueued(t *testing.T, topologyName tc.TopologyName, cdnID int64, dsID int) {
	const action = "queue"
	testCase := topologiesQueueUpdateTestCase{
		Description:                  "invalid update action",
		TopologiesQueueUpdateRequest: tc.TopologiesQueueUpdateRequest{Action: action, CDNID: cdnID},
	}
	resp, _, err := TOSession.TopologiesQueueUpdate(topologyName, testCase.TopologiesQueueUpdateRequest)
	if err != nil {
		t.Fatalf("received error queueing server updates on topology %s: %s", topologyName, err)
	}
	if resp.Action != action {
		t.Fatalf("expected action %s, received action %s", action, resp.Action)
	}
	if resp.CDNID != cdnID {
		t.Fatalf("expected CDN ID %d, received CDN ID %d", cdnID, resp.CDNID)
	}
	if topologyName != resp.Topology {
		t.Fatalf("expected topology %s, received topology %s", topologyName, resp.Topology)
	}
	params := url.Values{}
	dsIDString := strconv.Itoa(dsID)
	params.Set("dsId", dsIDString)
	serversResponse, _, err := TOSession.GetServers(params, nil)
	if err != nil {
		t.Fatalf("getting servers for delivery service with id %s: %s", dsIDString, err)
	}
	servers := serversResponse.Response
	for _, server := range servers {
		if *server.CDNID != int(cdnID) {
			continue
		}
		if !*server.UpdPending {
			t.Fatalf("expected UpdPending = %t for server with hostname %s, got UpdPending = %t", true, *server.HostName, *server.UpdPending)
		}
	}
}
