package v3

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
)

func TestGetOSVersions(t *testing.T) {
	if Config.NoISO {
		t.Skip("No ISO generation available")
	}
	// Default value per ./traffic_ops/install/data/json/osversions.json file.
	// This should be the data returned in the CiaB environment.
	expected := map[string]string{
		"CentOS 7.2": "centos72",
	}

	// Ensure request with an authenticated client returns expected data.
	t.Run("OK when AUTHENTICATED", func(t *testing.T) {
		got, _, err := TOSession.GetOSVersions()
		assert.NoError(t, err, "Unexpected error from authenticated GetOSVersions(): %v", err)
		assert.RequireEqual(t, len(expected), len(got), "Incorrect map length: got %d map entries, expected %d", len(expected), len(got))

		for k, expectedVal := range expected {
			assert.RequireEqual(t, expectedVal, got[k], "Incorrect map entry for key %q: got %q, expected %q", k, got[k], expectedVal)
		}
	})

	// Ensure request with an un-authenticated client returns an error.
	t.Run("ERROR when UNAUTHENTICATED", func(t *testing.T) {
		_, _, err := NoAuthTOSession.GetOSVersions()
		assert.Error(t, err, "Expected error from unauthenticated GetOSVersions(), got: <nil>")
	})

	// Update database with a Parameter entry. This should cause the endpoint
	// to use the Parameter's value as the configuration file's directory. In this
	// case, an intentionally missing/invalid directory is provided.
	// Ensure authenticated request client returns an error.
	// NOTE: This does not assume this test and TO are using the same filesystem, but
	// does make the reasonable assumption that `/DOES/NOT/EXIST/osversions.json` will not exist
	// on the TO host.
	t.Run("ERROR when INVALID PARAMETER", func(t *testing.T) {
		p := tc.Parameter{
			ConfigFile: "mkisofs",
			Name:       "kickstart.files.location",
			Value:      "/DOES/NOT/EXIST",
		}
		alerts, _, err := TOSession.CreateParameter(p)
		assert.RequireNoError(t, err, "Could not create Parameter: %v - alerts: %+v", err, alerts.Alerts)

		// Cleanup DB entry
		defer func() {
			resp, _, err := TOSession.GetParameterByNameAndConfigFileAndValueWithHdr(p.Name, p.ConfigFile, p.Value, nil)
			assert.RequireNoError(t, err, "Cannot GET Parameter by name '%s', configFile '%s' and value '%s': %v", p.Name, p.ConfigFile, p.Value, err)
			assert.RequireEqual(t, 1, len(resp), "Unexpected response length %d", len(resp))

			delResp, _, err := TOSession.DeleteParameterByID(resp[0].ID)
			assert.RequireNoError(t, err, "Cannot delete Parameter #%d: %v - alerts: %+v", resp[0].ID, err, delResp.Alerts)

		}()

		_, _, err = TOSession.GetOSVersions()
		assert.Error(t, err, "Expected error from GetOSVersions() after adding invalid Parameter DB entry, got: <nil>")
	})
}
