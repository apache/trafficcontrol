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

	"github.com/apache/trafficcontrol/v6/lib/go-tc"
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
		got, _, err := TOSession.GetOSVersions()
		if err != nil {
			t.Fatalf("unexpected error from authenticated GetOSVersions(): %v", err)
		}

		t.Logf("GetOSVersions() response: %#v", got)

		if lenGot, lenExp := len(got), len(expected); lenGot != lenExp {
			t.Fatalf("incorrect map length: got %d map entries, expected %d", lenGot, lenExp)
		}
		for k, expectedVal := range expected {
			if gotVal := got[k]; gotVal != expectedVal {
				t.Fatalf("incorrect map entry for key %q: got %q, expected %q", k, gotVal, expectedVal)
			}
		}
	})

	// Ensure request with an un-authenticated client returns an error.
	t.Run("un-authenticated", func(t *testing.T) {
		_, _, err := NoAuthTOSession.GetOSVersions()
		if err == nil {
			t.Fatalf("expected error from unauthenticated GetOSVersions(), got: %v", err)
		}
		t.Logf("unauthenticated GetOSVersions() error (expected): %v", err)
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
		if _, _, err := TOSession.CreateParameter(p); err != nil {
			t.Fatalf("could not CREATE parameter: %v\n", err)
		}
		// Cleanup DB entry
		defer func() {
			resp, _, err := TOSession.GetParameterByNameAndConfigFileAndValue(p.Name, p.ConfigFile, p.Value)
			if err != nil {
				t.Fatalf("cannot GET Parameter by name: %v - %v\n", p.Name, err)
			}
			if len(resp) != 1 {
				t.Fatalf("unexpected response length %d", len(resp))
			}

			if delResp, _, err := TOSession.DeleteParameterByID(resp[0].ID); err != nil {
				t.Fatalf("cannot DELETE Parameter by name: %v - %v\n", err, delResp)
			}
		}()

		_, _, err := TOSession.GetOSVersions()
		if err == nil {
			t.Fatalf("expected error from GetOSVersions() after adding invalid Parameter DB entry, got: %v", err)
		}
		t.Logf("got expected error from GetOSVersions() after adding Parameter DB entry with config directory %q: %v", p.Value, err)
	})
}
