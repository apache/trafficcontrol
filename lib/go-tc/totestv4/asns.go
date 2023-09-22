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

func CreateTestASNs(t *testing.T, cl *toclient.Session, dat TrafficControl) {
	for _, asn := range dat.ASNs {
		asn.CachegroupID = GetCacheGroupId(t, cl, asn.Cachegroup)()
		resp, _, err := cl.CreateASN(asn, toclient.RequestOptions{})
		assert.RequireNoError(t, err, "Could not create ASN: %v - alerts: %+v", err, resp)
	}
}

func DeleteTestASNs(t *testing.T, cl *toclient.Session) {
	asns, _, err := cl.GetASNs(toclient.RequestOptions{})
	assert.NoError(t, err, "Error trying to fetch ASNs for deletion: %v - alerts: %+v", err, asns.Alerts)

	for _, asn := range asns.Response {
		alerts, _, err := cl.DeleteASN(asn.ID, toclient.RequestOptions{})
		assert.NoError(t, err, "Cannot delete ASN %d: %v - alerts: %+v", asn.ASN, err, alerts)
		// Retrieve the ASN to see if it got deleted
		opts := toclient.NewRequestOptions()
		opts.QueryParameters.Set("asn", strconv.Itoa(asn.ASN))
		asns, _, err := cl.GetASNs(opts)
		assert.NoError(t, err, "Error trying to fetch ASN after deletion: %v - alerts: %+v", err, asns.Alerts)
		assert.Equal(t, 0, len(asns.Response), "Expected ASN %d to be deleted, but it was found in Traffic Ops", asn.ASN)
	}
}
