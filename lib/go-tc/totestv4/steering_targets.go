package totestv4

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

	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	toclient "github.com/apache/trafficcontrol/v8/traffic_ops/v4-client"
)

func CreateTestSteeringTargets(t *testing.T, cl *toclient.Session, td TrafficControl) {
	for _, st := range td.SteeringTargets {
		st.TypeID = util.IntPtr(GetTypeId(t, cl, *st.Type))
		st.DeliveryServiceID = util.UInt64Ptr(uint64(GetDeliveryServiceId(t, cl, string(*st.DeliveryService))()))
		st.TargetID = util.UInt64Ptr(uint64(GetDeliveryServiceId(t, cl, string(*st.Target))()))
		resp, _, err := cl.CreateSteeringTarget(st, toclient.RequestOptions{})
		assert.RequireNoError(t, err, "Creating steering target: %v - alerts: %+v", err, resp.Alerts)
	}
}

func DeleteTestSteeringTargets(t *testing.T, cl *toclient.Session, td TrafficControl) {
	dsIDs := []uint64{}
	for _, st := range td.SteeringTargets {
		opts := toclient.NewRequestOptions()
		opts.QueryParameters.Set("xmlId", string(*st.DeliveryService))
		respDS, _, err := cl.GetDeliveryServices(opts)
		assert.RequireNoError(t, err, "Deleting steering target: getting ds: %v - alerts: %+v", err, respDS.Alerts)
		assert.RequireEqual(t, 1, len(respDS.Response), "Deleting steering target: getting ds: expected 1 delivery service")
		assert.RequireNotNil(t, respDS.Response[0].ID, "Deleting steering target: getting ds: nil ID returned")

		dsID := uint64(*respDS.Response[0].ID)
		st.DeliveryServiceID = &dsID
		dsIDs = append(dsIDs, dsID)

		opts.QueryParameters.Set("xmlId", string(*st.Target))
		respTarget, _, err := cl.GetDeliveryServices(opts)
		assert.RequireNoError(t, err, "Deleting steering target: getting target ds: %v - alerts: %+v", err, respTarget.Alerts)
		assert.RequireEqual(t, 1, len(respTarget.Response), "Deleting steering target: getting target ds: expected 1 delivery service")
		assert.RequireNotNil(t, respTarget.Response[0].ID, "Deleting steering target: getting target ds: not found")

		targetID := uint64(*respTarget.Response[0].ID)
		st.TargetID = &targetID

		resp, _, err := cl.DeleteSteeringTarget(int(*st.DeliveryServiceID), int(*st.TargetID), toclient.RequestOptions{})
		assert.NoError(t, err, "Deleting steering target: deleting: %v - alerts: %+v", err, resp.Alerts)
	}

	for _, dsID := range dsIDs {
		sts, _, err := cl.GetSteeringTargets(int(dsID), toclient.RequestOptions{})
		assert.NoError(t, err, "deleting steering targets: getting steering target: %v - alerts: %+v", err, sts.Alerts)
		assert.Equal(t, 0, len(sts.Response), "Deleting steering targets: after delete, getting steering target: expected 0 actual %d", len(sts.Response))
	}
}
