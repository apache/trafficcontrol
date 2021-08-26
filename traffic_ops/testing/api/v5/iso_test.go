package v5

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

	"github.com/apache/trafficcontrol/lib/go-tc"
	client "github.com/apache/trafficcontrol/traffic_ops/v5-client"
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
	t.Run("authenticated", func(t *testing.T) {
		got, _, err := TOSession.GetOSVersions(client.RequestOptions{})
		if err != nil {
			t.Fatalf("unexpected error from authenticated GetOSVersions(): %v - alerts: %+v", err, got.Alerts)
		}

		if lenGot, lenExp := len(got.Response), len(expected); lenGot != lenExp {
			t.Fatalf("incorrect map length: got %d map entries, expected %d", lenGot, lenExp)
		}
		for k, expectedVal := range expected {
			if gotVal := got.Response[k]; gotVal != expectedVal {
				t.Fatalf("incorrect map entry for key %q: got %q, expected %q", k, gotVal, expectedVal)
			}
		}
	})

	// Ensure request with an un-authenticated client returns an error.
	t.Run("un-authenticated", func(t *testing.T) {
		_, _, err := NoAuthTOSession.GetOSVersions(client.RequestOptions{})
		if err == nil {
			t.Fatal("expected error from unauthenticated GetOSVersions(), got: <nil>")
		}
	})

	// Update database with a Parameter entry. This should cause the endpoint
	// to use the Parameter's value as the configuration file's directory. In this
	// case, an intentionally missing/invalid directory is provided.
	// Ensure authenticated request client returns an error.
	// NOTE: This does not assume this test and TO are using the same filesystem, but
	// does make the reasonable assumption that `/DOES/NOT/EXIST/osversions.json` will not exist
	// on the TO host.
	t.Run("parameter-invalid", func(t *testing.T) {
		p := tc.Parameter{
			ConfigFile: "mkisofs",
			Name:       "kickstart.files.location",
			Value:      "/DOES/NOT/EXIST",
		}
		if alerts, _, err := TOSession.CreateParameter(p, client.RequestOptions{}); err != nil {
			t.Fatalf("could not create Parameter: %v - alerts: %+v", err, alerts.Alerts)
		}
		// Cleanup DB entry
		defer func() {
			opts := client.NewRequestOptions()
			opts.QueryParameters.Set("name", p.Name)
			opts.QueryParameters.Set("configFile", p.ConfigFile)
			opts.QueryParameters.Set("value", p.Value)
			resp, _, err := TOSession.GetParameters(opts)
			if err != nil {
				t.Fatalf("cannot GET Parameter by name '%s', configFile '%s' and value '%s': %v - alerts: %+v", p.Name, p.ConfigFile, p.Value, err, resp.Alerts)
			}
			if len(resp.Response) != 1 {
				t.Fatalf("unexpected response length %d", len(resp.Response))
			}

			if delResp, _, err := TOSession.DeleteParameter(resp.Response[0].ID, client.RequestOptions{}); err != nil {
				t.Fatalf("cannot delete Parameter #%d: %v - alerts: %+v", resp.Response[0].ID, err, delResp.Alerts)
			}
		}()

		_, _, err := TOSession.GetOSVersions(client.RequestOptions{})
		if err == nil {
			t.Fatal("expected error from GetOSVersions() after adding invalid Parameter DB entry, got: <nil>")
		}
	})
}
