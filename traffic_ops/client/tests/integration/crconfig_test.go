package integration

import "testing"

//TestCachegroupResults compares the results of the Cachegroup api and Cachegroup client
func TestGetCrConfig(t *testing.T) {
	//Get a CDN from the to client
	cdn, err := GetCdn()
	if err != nil {
		t.Errorf("TestGetCrConfig -- Could not get CDNs from TO...%v\n", err)
	}

	crConfig, cacheHitStatus, err := to.GetCRConfig(cdn.Name)
	if err != nil {
		t.Errorf("Could not get CrConfig for %s.  Error is...%v\n", cdn.Name, err)
	}

	if cacheHitStatus == "" {
		t.Error("cacheHitStatus is empty...")
	}

	if len(crConfig) == 0 {
		t.Error("Raw CrConfig reponse was 0...")
	}
}
