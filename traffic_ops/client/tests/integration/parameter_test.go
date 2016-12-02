package integration

import (
	// "encoding/json"
	"encoding/json"
	"fmt"
	"testing"

	traffic_ops "github.com/apache/incubator-trafficcontrol/traffic_ops/client"
)

func TestParameters(t *testing.T) {
	profile, err := GetProfile()
	if err != nil {
		t.Errorf("Could not get a profile, error was: %v\n", err)
	}

	uri := fmt.Sprintf("/api/1.2/parameters/profile/%s.json", profile.Name)
	resp, err := Request(*to, "GET", uri, nil)
	if err != nil {
		t.Errorf("Could not get %s reponse was: %v\n", uri, err)
		t.FailNow()
	}

	defer resp.Body.Close()
	var apiParamRes traffic_ops.ParamResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiParamRes); err != nil {
		t.Errorf("Could not decode parameter json.  Error is: %v\n", err)
		t.FailNow()
	}
	apiParams := apiParamRes.Response

	clientParams, err := to.Parameters(profile.Name)
	if err != nil {
		t.Errorf("Could not get Hardware from client.  Error is: %v\n", err)
		t.FailNow()
	}

	if len(apiParams) != len(clientParams) {
		t.Errorf("Params Response Length -- expected %v, got %v\n", len(apiParams), len(clientParams))
	}

	for _, apiParam := range apiParams {
		match := false
		for _, clientParam := range clientParams {
			if apiParam.Name == clientParam.Name && apiParam.Value == clientParam.Value && apiParam.ConfigFile == clientParam.ConfigFile {
				match = true
			}
		}
		if !match {
			t.Errorf("Did not get a param matching %+v from the client\n", apiParam)
		}
	}
}
