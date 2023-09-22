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
	"strconv"
	"testing"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	toclient "github.com/apache/trafficcontrol/v8/traffic_ops/v5-client"
)

// this resets the IDs of things attached to a DS, which needs to be done
// because the WithObjs flow destroys and recreates those object IDs
// non-deterministically with each test - BUT, the client method permanently
// alters the DSR structures by adding these referential IDs. Older clients
// got away with it by not making 'DeliveryService' a pointer, but to add
// original/requested fields you need to sometimes allow each to be nil, so
// this is a problem that needs to be solved at some point.
// A better solution _might_ be to reload all the test fixtures every time
// to wipe any and all referential modifications made to any test data, but
// for now that's overkill.
func resetDS(ds *tc.DeliveryServiceV5) {
	if ds == nil {
		return
	}
	ds.CDNID = 0
	ds.ID = nil
	ds.ProfileID = nil
	ds.TenantID = 0
	ds.TypeID = 0
}

func CreateTestDeliveryServiceRequests(t *testing.T, cl *toclient.Session, td TrafficControl) {
	for _, dsr := range td.DeliveryServiceRequests {
		resetDS(dsr.Original)
		resetDS(dsr.Requested)
		respDSR, _, err := cl.CreateDeliveryServiceRequest(dsr, toclient.RequestOptions{})
		assert.NoError(t, err, "Could not create Delivery Service Requests: %v - alerts: %+v", err, respDSR.Alerts)
	}
}

func DeleteTestDeliveryServiceRequests(t *testing.T, cl *toclient.Session) {
	resp, _, err := cl.GetDeliveryServiceRequests(toclient.RequestOptions{})
	assert.NoError(t, err, "Cannot get Delivery Service Requests: %v - alerts: %+v", err, resp.Alerts)
	for _, request := range resp.Response {
		alert, _, err := cl.DeleteDeliveryServiceRequest(*request.ID, toclient.RequestOptions{})
		assert.NoError(t, err, "Cannot delete Delivery Service Request #%d: %v - alerts: %+v", request.ID, err, alert.Alerts)

		// Retrieve the DeliveryServiceRequest to see if it got deleted
		opts := toclient.NewRequestOptions()
		opts.QueryParameters.Set("id", strconv.Itoa(*request.ID))
		dsr, _, err := cl.GetDeliveryServiceRequests(opts)
		assert.NoError(t, err, "Unexpected error fetching Delivery Service Request #%d after deletion: %v - alerts: %+v", *request.ID, err, dsr.Alerts)
		assert.Equal(t, len(dsr.Response), 0, "Expected Delivery Service Request #%d to be deleted, but it was found in Traffic Ops", *request.ID)
	}
}
