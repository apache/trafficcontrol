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

	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	toclient "github.com/apache/trafficcontrol/v8/traffic_ops/v5-client"
)

func CreateTestStaticDNSEntries(t *testing.T, cl *toclient.Session, td TrafficControl) {
	for _, staticDNSEntry := range td.StaticDNSEntries {
		resp, _, err := cl.CreateStaticDNSEntry(staticDNSEntry, toclient.RequestOptions{})
		assert.RequireNoError(t, err, "Could not create Static DNS Entry: %v - alerts: %+v", err, resp.Alerts)
	}
}

func DeleteTestStaticDNSEntries(t *testing.T, cl *toclient.Session) {
	staticDNSEntries, _, err := cl.GetStaticDNSEntries(toclient.RequestOptions{})
	assert.NoError(t, err, "Cannot get Static DNS Entries: %v - alerts: %+v", err, staticDNSEntries.Alerts)

	for _, staticDNSEntry := range staticDNSEntries.Response {
		alerts, _, err := cl.DeleteStaticDNSEntry(*staticDNSEntry.ID, toclient.RequestOptions{})
		assert.NoError(t, err, "Unexpected error deleting Static DNS Entry '%s' (#%d): %v - alerts: %+v", staticDNSEntry.Host, staticDNSEntry.ID, err, alerts.Alerts)
		// Retrieve the Static DNS Entry to see if it got deleted
		opts := toclient.NewRequestOptions()
		opts.QueryParameters.Set("host", *staticDNSEntry.Host)
		getStaticDNSEntry, _, err := cl.GetStaticDNSEntries(opts)
		assert.NoError(t, err, "Error getting Static DNS Entry '%s' after deletion: %v - alerts: %+v", staticDNSEntry.Host, err, getStaticDNSEntry.Alerts)
		assert.Equal(t, 0, len(getStaticDNSEntry.Response), "Expected Static DNS Entry '%s' to be deleted, but it was found in Traffic Ops", staticDNSEntry.Host)
	}
}
