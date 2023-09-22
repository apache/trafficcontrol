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
	"strconv"
	"testing"

	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	toclient "github.com/apache/trafficcontrol/v8/traffic_ops/v4-client"
)

func CreateTestDeliveryServiceRequestComments(t *testing.T, cl *toclient.Session, td TrafficControl) {
	for _, comment := range td.DeliveryServiceRequestComments {
		opts := toclient.NewRequestOptions()
		opts.QueryParameters.Set("xmlId", comment.XMLID)
		resp, _, err := cl.GetDeliveryServiceRequests(opts)
		assert.NoError(t, err, "Cannot get Delivery Service Request by XMLID '%s': %v - alerts: %+v", comment.XMLID, err, resp.Alerts)
		assert.Equal(t, len(resp.Response), 1, "Found %d Delivery Service request by XMLID '%s, expected exactly one", len(resp.Response), comment.XMLID)
		assert.NotNil(t, resp.Response[0].ID, "Got Delivery Service Request with xml_id '%s' that had a null ID", comment.XMLID)

		comment.DeliveryServiceRequestID = *resp.Response[0].ID
		alerts, _, err := cl.CreateDeliveryServiceRequestComment(comment, toclient.RequestOptions{})
		assert.NoError(t, err, "Could not create Delivery Service Request Comment: %v - alerts: %+v", err, alerts.Alerts)
	}
}

func DeleteTestDeliveryServiceRequestComments(t *testing.T, cl *toclient.Session) {
	comments, _, err := cl.GetDeliveryServiceRequestComments(toclient.RequestOptions{})
	assert.NoError(t, err, "Unexpected error getting Delivery Service Request Comments: %v - alerts: %+v", err, comments.Alerts)

	for _, comment := range comments.Response {
		resp, _, err := cl.DeleteDeliveryServiceRequestComment(comment.ID, toclient.RequestOptions{})
		assert.NoError(t, err, "Cannot delete Delivery Service Request Comment #%d: %v - alerts: %+v", comment.ID, err, resp.Alerts)

		// Retrieve the delivery service request comment to see if it got deleted
		opts := toclient.NewRequestOptions()
		opts.QueryParameters.Set("id", strconv.Itoa(comment.ID))
		comments, _, err := cl.GetDeliveryServiceRequestComments(opts)
		assert.NoError(t, err, "Unexpected error fetching Delivery Service Request Comment %d after deletion: %v - alerts: %+v", comment.ID, err, comments.Alerts)
		assert.Equal(t, len(comments.Response), 0, "Expected Delivery Service Request Comment #%d to be deleted, but it was found in Traffic Ops", comment.ID)
	}
}
