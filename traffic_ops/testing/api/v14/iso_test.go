package v14

import (
	"testing"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

func TestGetOSVersions(t *testing.T) {
	// Default value per ./traffic_ops/install/data/perl/osversions.cfg file.
	// This is what the CIAB environment uses.
	expected := map[string]string{
		"CentOS 7.2": "centos72",
	}

	// Ensure request with an authenticated client returns expected data.

	got, _, err := TOSession.GetOSVersions()
	if err != nil {
		t.Errorf("unexpected error from authenticated GetOSVersions(): %v", err)
	}
	if lenGot, lenExp := len(got), len(expected); lenGot != lenExp {
		t.Errorf("incorrect map length: got %d map entries, expected %d", lenGot, lenExp)
	}
	for k, expectedVal := range expected {
		if gotVal := got[k]; gotVal != expectedVal {
			t.Errorf("incorrect map entry for key %q: got %q, expected %q", k, gotVal, expectedVal)
		}
	}

	t.Logf("GetOSVersions() response: %#v", got)

	// Ensure request with an un-authenticated client returns an error.

	_, _, err = NoAuthTOSession.GetOSVersions()
	if err == nil {
		t.Errorf("expected error from unauthenticated GetOSVersions(), got: %v", err)
	} else {
		t.Logf("unauthenticated GetOSVersions() error (expected): %v", err)
	}

	// Update database with a Parameter entry. This should cause the endpoint
	// to use the Parameter's value as the configuration file's directory. In this
	// case, an intentionally missing/invalid directory is provided.
	// Ensure authenticated request client returns an error.

	p := tc.Parameter{
		ConfigFile: "mkisofs",
		Name:       "kickstart.files.location",
		Value:      "/DOES/NOT/EXIST",
	}

	if _, _, err = TOSession.CreateParameter(p); err != nil {
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

	got, _, err = TOSession.GetOSVersions()
	if err == nil {
		t.Errorf("expected error from GetOSVersions() after adding invalid Parameter DB entry, got: %v", err)
	} else {
		t.Logf("got expected error from GetOSVersions() after adding Parameter DB entry with config directory: %q", p.Value)
	}
}
