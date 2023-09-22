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

func CreateTestDeliveryServices(t *testing.T, cl *toclient.Session, td TrafficControl) {
	for _, ds := range td.DeliveryServices {
		if ds.XMLID == "" {
			t.Error("Found a Delivery Service in testing data with null or undefined XMLID")
			continue
		}
		resp, _, err := cl.CreateDeliveryService(ds, toclient.RequestOptions{})
		assert.NoError(t, err, "Could not create Delivery Service '%s': %v - alerts: %+v", ds.XMLID, err, resp.Alerts)
	}
}

func DeleteTestDeliveryServices(t *testing.T, cl *toclient.Session) {
	dses, _, err := cl.GetDeliveryServices(toclient.RequestOptions{})
	assert.NoError(t, err, "Cannot get Delivery Services: %v - alerts: %+v", err, dses.Alerts)

	for _, ds := range dses.Response {
		delResp, _, err := cl.DeleteDeliveryService(*ds.ID, toclient.RequestOptions{})
		assert.NoError(t, err, "Could not delete Delivery Service: %v - alerts: %+v", err, delResp.Alerts)
		// Retrieve Delivery Service to see if it got deleted
		opts := toclient.NewRequestOptions()
		opts.QueryParameters.Set("id", strconv.Itoa(*ds.ID))
		getDS, _, err := cl.GetDeliveryServices(opts)
		assert.NoError(t, err, "Error deleting Delivery Service for '%s' : %v - alerts: %+v", ds.XMLID, err, getDS.Alerts)
		assert.Equal(t, 0, len(getDS.Response), "Expected Delivery Service '%s' to be deleted", ds.XMLID)
	}
}

func GetDeliveryServiceId(t *testing.T, cl *toclient.Session, xmlId string) func() int {
	return func() int {
		opts := toclient.NewRequestOptions()
		opts.QueryParameters.Set("xmlId", xmlId)

		resp, _, err := cl.GetDeliveryServices(opts)
		assert.RequireNoError(t, err, "Get Delivery Service Request failed with error: %v", err)
		assert.RequireEqual(t, 1, len(resp.Response), "Expected delivery service response object length 1, but got %d", len(resp.Response))
		assert.RequireNotNil(t, resp.Response[0].ID, "Expected id to not be nil")

		return *resp.Response[0].ID
	}
}

func setInactive(t *testing.T, cl *toclient.Session, dsID int) {
	opts := toclient.NewRequestOptions()
	opts.QueryParameters.Set("id", strconv.Itoa(dsID))
	resp, _, err := cl.GetDeliveryServices(opts)
	assert.RequireNoError(t, err, "Failed to fetch details for Delivery Service #%d: %v - alerts: %+v", dsID, err, resp.Alerts)
	assert.RequireEqual(t, len(resp.Response), 1, "Expected exactly one Delivery Service to exist with ID %d, found: %d", dsID, len(resp.Response))

	ds := resp.Response[0]
	if ds.Active == "" {
		t.Errorf("Deliver Service #%d had null or undefined 'active'", dsID)
		ds.Active = tc.DSActiveStateInactive
	}
	if ds.Active == tc.DSActiveStateActive || ds.Active == tc.DSActiveStatePrimed {
		ds.Active = tc.DSActiveStateInactive
		_, _, err = cl.UpdateDeliveryService(dsID, ds, toclient.RequestOptions{})
		assert.NoError(t, err, "Failed to set Delivery Service #%d to inactive: %v", dsID, err)
	}
}
