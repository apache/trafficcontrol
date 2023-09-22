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

func CreateTestCacheGroups(t *testing.T, cl *toclient.Session, dat TrafficControl) {
	for _, cg := range dat.CacheGroups {

		resp, _, err := cl.CreateCacheGroup(cg, toclient.RequestOptions{})
		if err != nil {
			t.Errorf("could not create Cache Group: %v - alerts: %+v", err, resp.Alerts)
			continue
		}

		// Testing 'join' fields during create
		if cg.ParentName != nil && resp.Response.ParentName == nil {
			t.Error("Parent cachegroup is null in response when it should have a value")
		}
		if cg.SecondaryParentName != nil && resp.Response.SecondaryParentName == nil {
			t.Error("Secondary parent cachegroup is null in response when it should have a value")
		}
		if cg.Type != nil && resp.Response.Type == nil {
			t.Error("Type is null in response when it should have a value")
		}
		assert.NotNil(t, resp.Response.LocalizationMethods, "Localization methods are null")
		assert.NotNil(t, resp.Response.Fallbacks, "Fallbacks are null")
	}
}

func DeleteTestCacheGroups(t *testing.T, cl *toclient.Session, dat TrafficControl) {
	var parentlessCacheGroups []tc.CacheGroupNullableV5
	opts := toclient.NewRequestOptions()

	// delete the edge caches.
	for _, cg := range dat.CacheGroups {
		if cg.Name == nil {
			t.Error("Found a Cache Group with null or undefined name")
			continue
		}

		// Retrieve the CacheGroup by name so we can get the id for Deletion
		opts.QueryParameters.Set("name", *cg.Name)
		resp, _, err := cl.GetCacheGroups(opts)
		assert.NoError(t, err, "Cannot GET CacheGroup by name '%s': %v - alerts: %+v", *cg.Name, err, resp.Alerts)

		if len(resp.Response) < 1 {
			t.Errorf("Could not find test data Cache Group '%s' in Traffic Ops", *cg.Name)
			continue
		}
		cg = resp.Response[0]

		// Cachegroups that are parents (usually mids but sometimes edges)
		// need to be deleted only after the children cachegroups are deleted.
		if cg.ParentCachegroupID == nil && cg.SecondaryParentCachegroupID == nil {
			parentlessCacheGroups = append(parentlessCacheGroups, cg)
			continue
		}

		if cg.ID == nil {
			t.Error("Traffic Ops returned a Cache Group with null or undefined ID")
			continue
		}

		alerts, _, err := cl.DeleteCacheGroup(*cg.ID, toclient.RequestOptions{})
		assert.NoError(t, err, "Cannot delete Cache Group: %v - alerts: %+v", err, alerts)

		// Retrieve the CacheGroup to see if it got deleted
		opts.QueryParameters.Set("name", *cg.Name)
		cgs, _, err := cl.GetCacheGroups(opts)
		assert.NoError(t, err, "Error deleting Cache Group by name: %v - alerts: %+v", err, cgs.Alerts)
		assert.Equal(t, 0, len(cgs.Response), "Expected CacheGroup name: %s to be deleted", *cg.Name)
	}

	opts = toclient.NewRequestOptions()
	// now delete the parentless cachegroups
	for _, cg := range parentlessCacheGroups {
		// nil check for cg.Name occurs prior to insertion into parentlessCacheGroups
		opts.QueryParameters.Set("name", *cg.Name)
		// Retrieve the CacheGroup by name so we can get the id for Deletion
		resp, _, err := cl.GetCacheGroups(opts)
		assert.NoError(t, err, "Cannot get Cache Group by name '%s': %v - alerts: %+v", *cg.Name, err, resp.Alerts)

		if len(resp.Response) < 1 {
			t.Errorf("Cache Group '%s' somehow stopped existing since the last time we ask Traffic Ops about it", *cg.Name)
			continue
		}

		respCG := resp.Response[0]
		if respCG.ID == nil {
			t.Errorf("Traffic Ops returned Cache Group '%s' with null or undefined ID", *cg.Name)
			continue
		}
		delResp, _, err := cl.DeleteCacheGroup(*respCG.ID, toclient.RequestOptions{})
		assert.NoError(t, err, "Cannot delete Cache Group '%s': %v - alerts: %+v", *respCG.Name, err, delResp.Alerts)

		// Retrieve the CacheGroup to see if it got deleted
		opts.QueryParameters.Set("name", *cg.Name)
		cgs, _, err := cl.GetCacheGroups(opts)
		assert.NoError(t, err, "Error attempting to fetch Cache Group '%s' after deletion: %v - alerts: %+v", *cg.Name, err, cgs.Alerts)
		assert.Equal(t, 0, len(cgs.Response), "Expected Cache Group '%s' to be deleted", *cg.Name)
	}
}

func GetCacheGroupId(t *testing.T, cl *toclient.Session, cacheGroupName string) func() int {
	return func() int {
		opts := toclient.NewRequestOptions()
		opts.QueryParameters.Set("name", cacheGroupName)

		resp, _, err := cl.GetCacheGroups(opts)
		assert.RequireNoError(t, err, "Get Cache Groups Request failed with error: %v", err)
		assert.RequireEqual(t, len(resp.Response), 1, "Expected response object length 1, but got %d", len(resp.Response))
		assert.RequireNotNil(t, resp.Response[0].ID, "Expected id to not be nil")

		return *resp.Response[0].ID
	}
}
