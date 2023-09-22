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

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	toclient "github.com/apache/trafficcontrol/v8/traffic_ops/v5-client"
)

func CreateTestDeliveryServicesRequiredCapabilities(t *testing.T, cl *toclient.Session, td TrafficControl) {
	// Assign all required capability to delivery services listed in `tc-fixtures.json`.
	for _, dsrc := range td.DeliveryServicesRequiredCapabilities {
		dsId := GetDeliveryServiceId(t, cl, *dsrc.XMLID)()
		dsrc = tc.DeliveryServicesRequiredCapability{
			DeliveryServiceID:  &dsId,
			RequiredCapability: dsrc.RequiredCapability,
		}
		resp, _, err := cl.CreateDeliveryServicesRequiredCapability(dsrc, toclient.RequestOptions{})
		assert.NoError(t, err, "Unexpected error creating a Delivery Service/Required Capability relationship: %v - alerts: %+v", err, resp.Alerts)
	}
}

func DeleteTestDeliveryServicesRequiredCapabilities(t *testing.T, cl *toclient.Session) {
	// Get Required Capabilities to delete them
	dsrcs, _, err := cl.GetDeliveryServicesRequiredCapabilities(toclient.RequestOptions{})
	assert.NoError(t, err, "Error getting Delivery Service/Required Capability relationships: %v - alerts: %+v", err, dsrcs.Alerts)

	for _, dsrc := range dsrcs.Response {
		alerts, _, err := cl.DeleteDeliveryServicesRequiredCapability(*dsrc.DeliveryServiceID, *dsrc.RequiredCapability, toclient.RequestOptions{})
		assert.NoError(t, err, "Error deleting a relationship between a Delivery Service and a Capability: %v - alerts: %+v", err, alerts.Alerts)
	}
}
