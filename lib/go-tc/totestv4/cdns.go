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

func GetCDNID(t *testing.T, cl *toclient.Session, cdnName string) func() int {
	return func() int {
		opts := toclient.NewRequestOptions()
		opts.QueryParameters.Set("name", cdnName)
		cdnsResp, _, err := cl.GetCDNs(opts)
		assert.RequireNoError(t, err, "Get CDNs Request failed with error:", err)
		assert.RequireEqual(t, 1, len(cdnsResp.Response), "Expected response object length 1, but got %d", len(cdnsResp.Response))
		assert.RequireNotNil(t, cdnsResp.Response[0].ID, "Expected id to not be nil")
		return cdnsResp.Response[0].ID
	}
}

func CreateTestCDNs(t *testing.T, cl *toclient.Session, dat TrafficControl) {
	for _, cdn := range dat.CDNs {
		resp, _, err := cl.CreateCDN(cdn, toclient.RequestOptions{})
		assert.NoError(t, err, "Could not create CDN: %v - alerts: %+v", err, resp.Alerts)
	}
}

func DeleteTestCDNs(t *testing.T, cl *toclient.Session) {
	resp, _, err := cl.GetCDNs(toclient.RequestOptions{})
	assert.NoError(t, err, "Cannot get CDNs: %v - alerts: %+v", err, resp.Alerts)
	for _, cdn := range resp.Response {
		delResp, _, err := cl.DeleteCDN(cdn.ID, toclient.RequestOptions{})
		assert.NoError(t, err, "Cannot delete CDN '%s' (#%d): %v - alerts: %+v", cdn.Name, cdn.ID, err, delResp.Alerts)

		// Retrieve the CDN to see if it got deleted
		opts := toclient.NewRequestOptions()
		opts.QueryParameters.Set("id", strconv.Itoa(cdn.ID))
		cdns, _, err := cl.GetCDNs(opts)
		assert.NoError(t, err, "Error deleting CDN '%s': %v - alerts: %+v", cdn.Name, err, cdns.Alerts)
		assert.Equal(t, 0, len(cdns.Response), "Expected CDN '%s' to be deleted", cdn.Name)
	}
}
