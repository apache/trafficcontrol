/*

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package v4

import (
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-rfc"
	tc "github.com/apache/trafficcontrol/lib/go-tc"
)

func TestParameters(t *testing.T) {

	//toReqTimeout := time.Second * time.Duration(Config.Default.Session.TimeoutInSecs)
	//SwitchSession(toReqTimeout, Config.TrafficOps.URL, Config.TrafficOps.Users.Admin, Config.TrafficOps.UserPassword, Config.TrafficOps.Users.Portal, Config.TrafficOps.UserPassword)

	WithObjs(t, []TCObj{Parameters}, func() {
		GetTestParametersIMS(t)
		currentTime := time.Now().UTC().Add(-5 * time.Second)
		time := currentTime.Format(time.RFC1123)
		var header http.Header
		header = make(map[string][]string)
		header.Set(rfc.IfModifiedSince, time)
		header.Set(rfc.IfUnmodifiedSince, time)
		UpdateTestParameters(t)
		UpdateTestParametersWithHeaders(t, header)
		GetTestParameters(t)
		GetTestParametersIMSAfterChange(t, header)
		header = make(map[string][]string)
		etag := rfc.ETag(currentTime)
		header.Set(rfc.IfMatch, etag)
		UpdateTestParametersWithHeaders(t, header)
		GetTestPaginationSupportParameters(t)
		GetTestParametersByConfigfile(t)
		GetTestParametersByValue(t)
		GetTestParametersByName(t)
		GetParametersByInvalidId(t)
		GetParametersByInvalidName(t)
		GetParametersByInvalidConfigfile(t)
		GetParametersByInvalidValue(t)
		CreateTestParametersAlreadyExist(t)
		CreateTestParametersMissingName(t)
		CreateTestParametersMissingconfigFile(t)
		CreateMultipleTestParameters(t)
		DeleteTestParametersByInvalidId(t)
		UpdateParametersInvalidValue(t)
		UpdateParametersInvalidName(t)
		UpdateParametersInvalidConfigFile(t)
	})
}

func UpdateTestParametersWithHeaders(t *testing.T, header http.Header) {
	if len(testData.Parameters) > 0 {
		firstParameter := testData.Parameters[0]
		// Retrieve the Parameter by name so we can get the id for the Update
		resp, _, err := TOSession.GetParametersByProfileName(firstParameter.Name, header)
		if err != nil {
			t.Errorf("cannot GET Parameter by name: %v - %v", firstParameter.Name, err)
		}
		if len(resp) > 0 {
			remoteParameter := resp[0]
			expectedParameterValue := "UPDATED"
			remoteParameter.Value = expectedParameterValue
			_, reqInf, err := TOSession.UpdateParameter(remoteParameter.ID, remoteParameter, header)
			if err == nil {
				t.Errorf("Expected error about precondition failed, but got none")
			}
			if reqInf.StatusCode != http.StatusPreconditionFailed {
				t.Errorf("Expected status code 412, got %v", reqInf.StatusCode)
			}
		}
	}
}

func GetTestParametersIMSAfterChange(t *testing.T, header http.Header) {
	params := url.Values{}
	for _, pl := range testData.Parameters {
		params.Set("name", pl.Name)
		_, reqInf, err := TOSession.GetParameters(params, header)
		if err != nil {
			t.Errorf("Expected no error, but got %v", err.Error())
		}
		if reqInf.StatusCode != http.StatusOK {
			t.Errorf("Expected 200 status code, got %v", reqInf.StatusCode)
		}
	}
	currentTime := time.Now().UTC()
	currentTime = currentTime.Add(1 * time.Second)
	timeStr := currentTime.Format(time.RFC1123)
	header.Set(rfc.IfModifiedSince, timeStr)
	for _, pl := range testData.Parameters {
		params.Set("name", pl.Name)
		_, reqInf, err := TOSession.GetParameters(params, header)
		if err != nil {
			t.Errorf("Expected no error, but got %v", err.Error())
		}
		if reqInf.StatusCode != http.StatusNotModified {
			t.Errorf("Expected 304 status code, got %v", reqInf.StatusCode)
		}
	}
}

func CreateTestParameters(t *testing.T) {

	for _, pl := range testData.Parameters {
		resp, _, err := TOSession.CreateParameter(pl)
		t.Log("Response: ", resp)
		if err != nil {
			t.Errorf("could not CREATE parameters: %v", err)
		}
	}
}

func CreateMultipleTestParameters(t *testing.T) {
	DeleteTestParameters(t)

	pls := []tc.Parameter{}
	for _, pl := range testData.Parameters {
		pls = append(pls, pl)
	}
	resp, _, err := TOSession.CreateMultipleParameters(pls)
	t.Log("Response: ", resp)
	if err != nil {
		t.Errorf("could not CREATE parameters: %v", err)
	}

}

func UpdateTestParameters(t *testing.T) {

	firstParameter := testData.Parameters[0]
	// Retrieve the Parameter by name so we can get the id for the Update
	params := url.Values{}
	params.Set("name", firstParameter.Name)
	resp, _, err := TOSession.GetParameters(params, nil)
	if err != nil {
		t.Errorf("cannot GET Parameter by name: %v - %v", firstParameter.Name, err)
	}
	remoteParameter := resp[0]
	old := resp[0]
	expectedParameterValue := "UPDATED"
	expectedParameterConfigFile := "updatedConfigFile"
	expectedParameterName := "updateName"
	remoteParameter.Value = expectedParameterValue
	remoteParameter.ConfigFile = expectedParameterConfigFile
	remoteParameter.Name = expectedParameterName
	var alert tc.Alerts
	alert, _, err = TOSession.UpdateParameter(remoteParameter.ID, remoteParameter, nil)
	if err != nil {
		t.Errorf("cannot UPDATE Parameter by id: %v - %v", err, alert)
	}

	// Retrieve the Parameter to check Parameter name got updated
	resp, _, err = TOSession.GetParameterByID(remoteParameter.ID, nil)
	if err != nil {
		t.Errorf("cannot GET Parameter by name: %v - %v", firstParameter.Name, err)
	}
	respParameter := resp[0]
	if respParameter.Value != expectedParameterValue {
		t.Errorf("results do not match actual: %s, expected: %s", respParameter.Value, expectedParameterValue)
	}
		//update back to old value for safe deletion in pre-requisite

	alert, _, err = TOSession.UpdateParameter(remoteParameter.ID, old, nil)
	if err != nil {
		t.Errorf("cannot UPDATE Parameter by id: %v - %v", err, alert)
	}
}

func UpdateParametersInvalidValue(t *testing.T){
	firstParameter := testData.Parameters[0]
	// Retrieve the Parameter by name so we can get the id for the Update
	params := url.Values{}
	params.Set("name", firstParameter.Name)
	resp, _, err := TOSession.GetParameters(params, nil)
	if err != nil {
		t.Errorf("cannot GET Parameter by name: %v - %v", firstParameter.Name, err)
	}
	remoteParameter := resp[0]
	remoteParameter.Value = ""

	var alert tc.Alerts
	alert, reqInf, err := TOSession.UpdateParameter(remoteParameter.ID, remoteParameter, nil)
	if err != nil {
		t.Errorf("cannot UPDATE Parameter by id: %v - %v", err, alert)
	}
	if reqInf.StatusCode != http.StatusOK {
		t.Errorf("Expected 200 status code, got %v", reqInf.StatusCode)
	}
}

func UpdateParametersInvalidName(t *testing.T){
	firstParameter := testData.Parameters[0]
	// Retrieve the Parameter by name so we can get the id for the Update
	params := url.Values{}
	params.Set("name", firstParameter.Name)
	resp, _, err := TOSession.GetParameters(params, nil)
	if err != nil {
		t.Errorf("cannot GET Parameter by name: %v - %v", firstParameter.Name, err)
	}
	remoteParameter := resp[0]
	remoteParameter.Name = ""

	var alert tc.Alerts
	alert, reqInf, err := TOSession.UpdateParameter(remoteParameter.ID, remoteParameter, nil)
	if err == nil {
		t.Errorf("Invalid name has been updated by ID: %v - %v", err, alert)
	}
	if reqInf.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected 400 status code, got %v", reqInf.StatusCode)
	}
}

func UpdateParametersInvalidConfigFile(t *testing.T){
	firstParameter := testData.Parameters[0]
	// Retrieve the Parameter by name so we can get the id for the Update
	params := url.Values{}
	params.Set("name", firstParameter.Name)
	resp, _, err := TOSession.GetParameters(params, nil)
	if err != nil {
		t.Errorf("cannot GET Parameter by name: %v - %v", firstParameter.Name, err)
	}
	remoteParameter := resp[0]
	remoteParameter.ConfigFile = ""

	var alert tc.Alerts
	alert, reqInf, err := TOSession.UpdateParameter(remoteParameter.ID, remoteParameter, nil)
	if err == nil {
		t.Errorf("Invalid Config File has been updated by ID: %v - %v", err, alert)
	}
	if reqInf.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected 400 status code, got %v", reqInf.StatusCode)
	}
}

func GetTestParametersIMS(t *testing.T) {
	var header http.Header
	header = make(map[string][]string)
	futureTime := time.Now().AddDate(0, 0, 1)
	time := futureTime.Format(time.RFC1123)
	header.Set(rfc.IfModifiedSince, time)
	params := url.Values{}
	for _, pl := range testData.Parameters {
		params.Set("name", pl.Name)
		_, reqInf, err := TOSession.GetParameters(params, header)
		if err != nil {
			t.Errorf("Expected no error, but got %v", err.Error())
		}
		if reqInf.StatusCode != http.StatusNotModified {
			t.Errorf("Expected 304 status code, got %v", reqInf.StatusCode)
		}
	}
}

func GetTestParameters(t *testing.T) {
	params := url.Values{}
	for _, pl := range testData.Parameters {
		params.Set("name", pl.Name)
		resp, _, err := TOSession.GetParameters(params, nil)
		if err != nil {
			t.Errorf("cannot GET Parameter by name: %v - %v", err, resp)
		}
	}
}

func DeleteTestParametersParallel(t *testing.T) {

	var wg sync.WaitGroup
	for _, pl := range testData.Parameters {

		wg.Add(1)
		go func(p tc.Parameter) {
			defer wg.Done()
			DeleteTestParameter(t, p)
		}(pl)

	}
	wg.Wait()
}

func DeleteTestParameters(t *testing.T) {

	for _, pl := range testData.Parameters {
		DeleteTestParameter(t, pl)
	}
}

func DeleteTestParameter(t *testing.T, pl tc.Parameter) {

	// Retrieve the Parameter by name so we can get the id for the Update
	params := url.Values{}
	params.Set("name", pl.Name)
	params.Set("configFile", pl.ConfigFile)
	resp, _, err := TOSession.GetParameters(params, nil)
	if err != nil {
		t.Errorf("cannot GET Parameter by name: %v - %v", pl.Name, err)
	}

	if len(resp) == 0 {
		// TODO This fails for the ProfileParameters test; determine a way to check this, even for ProfileParameters
		// t.Errorf("DeleteTestParameter got no params for %+v %+v", pl.Name, pl.ConfigFile)
	} else if len(resp) > 1 {
		// TODO figure out why this happens, and be more precise about deleting things where created.
		// t.Errorf("DeleteTestParameter params for %+v %+v expected 1, actual %+v", pl.Name, pl.ConfigFile, len(resp))
	}

	for _, respParameter := range resp {
		delResp, _, err := TOSession.DeleteParameter(respParameter.ID)
		if err != nil {
			t.Errorf("cannot DELETE Parameter by name: %v - %v", err, delResp)
		}

		// Retrieve the Parameter to see if it got deleted
		pls, _, err := TOSession.GetParameterByID(pl.ID, nil)
		if err != nil {
			t.Errorf("error deleting Parameter name: %s", err.Error())
		}
		if len(pls) > 0 {
			t.Errorf("expected Parameter Name: %s and ConfigFile: %s to be deleted", pl.Name, pl.ConfigFile)
		}
	}
}

func DeleteTestParametersByInvalidId(t *testing.T) {
	delResp, _, err := TOSession.DeleteParameter(10000)
	if err == nil {
		t.Errorf("cannot DELETE Parameters by Invalid ID: %v - %v", err, delResp)
	}
}

func GetTestPaginationSupportParameters(t *testing.T) {
	qparams := url.Values{}
	qparams.Set("orderby", "id")
	parameters, _, err := TOSession.GetParameters(qparams, nil)
	if err != nil {
		t.Errorf("cannot GET Parameters: %v", err)
	}

	if len(parameters) > 0 {
		qparams = url.Values{}
		qparams.Set("orderby", "id")
		qparams.Set("limit", "1")
		parametersWithLimit, _, err := TOSession.GetParameters(qparams, nil)
		if err == nil {
			if !reflect.DeepEqual(parameters[:1], parametersWithLimit) {
				t.Error("expected GET Parameters with limit = 1 to return first result")
			}
		} else {
			t.Error("Error in getting parameter by limit")
		}
		if len(parameters) > 1 {
			qparams = url.Values{}
			qparams.Set("orderby", "id")
			qparams.Set("limit", "1")
			qparams.Set("offset", "1")
			parametersWithOffset, _, err := TOSession.GetParameters(qparams, nil)
			if err == nil {
				if !reflect.DeepEqual(parameters[1:2], parametersWithOffset) {
					t.Error("expected GET Parameters with limit = 1, offset = 1 to return second result")
				}
			} else {
				t.Error("Error in getting parameter by limit and offset")
			}

			qparams = url.Values{}
			qparams.Set("orderby", "id")
			qparams.Set("limit", "1")
			qparams.Set("page", "2")
			parametersWithPage, _, err := TOSession.GetParameters(qparams, nil)
			if err == nil {
				if !reflect.DeepEqual(parameters[1:2], parametersWithPage) {
					t.Error("expected GET Parameters with limit = 1, page = 2 to return second result")
				}
			} else {
				t.Error("Error in getting parameters by limit and page")
			}
		} else {
			t.Errorf("only one parameters found, so offset functionality can't test")
		}
	} else {
		t.Errorf("No parameters found to check pagination")
	}

	qparams = url.Values{}
	qparams.Set("limit", "-2")
	_, _, err = TOSession.GetParameters(qparams, nil)
	if err == nil {
		t.Error("expected GET Parameters to return an error when limit is not bigger than -1")
	} else if !strings.Contains(err.Error(), "must be bigger than -1") {
		t.Errorf("expected GET Parameters to return an error for limit is not bigger than -1, actual error: " + err.Error())
	}

	qparams = url.Values{}
	qparams.Set("limit", "1")
	qparams.Set("offset", "0")
	_, _, err = TOSession.GetParameters(qparams, nil)
	if err == nil {
		t.Error("expected GET Parameters to return an error when offset is not a positive integer")
	} else if !strings.Contains(err.Error(), "must be a positive integer") {
		t.Errorf("expected GET Parameters to return an error for offset is not a positive integer, actual error: " + err.Error())
	}

	qparams = url.Values{}
	qparams.Set("limit", "1")
	qparams.Set("page", "0")
	_, _, err = TOSession.GetParameters(qparams, nil)
	if err == nil {
		t.Error("expected GET Parameters to return an error when page is not a positive integer")
	} else if !strings.Contains(err.Error(), "must be a positive integer") {
		t.Errorf("expected GET Parameters to return an error for page is not a positive integer, actual error: " + err.Error())
	}
}

func GetTestParametersByConfigfile(t *testing.T) {
	for _, parameters := range testData.Parameters {
		resp, reqInf, err := TOSession.GetParametersByConfigFile(parameters.ConfigFile, nil)
		if err != nil {
			t.Errorf("cannot GET Parameter by Config File: %v - %v", parameters.ConfigFile, err)
		}
		if reqInf.StatusCode != http.StatusOK {
			t.Errorf("Expected 200 status code, got %v", reqInf.StatusCode)
		}
		if len(resp) <= 0 {
			t.Errorf("No data available for Get Parameters by Config file")
		}
	}
}

func GetTestParametersByValue(t *testing.T) {
	for _, parameters := range testData.Parameters {
		_, reqInf, err := TOSession.GetParametersByValue(parameters.Value, nil)
		if err != nil {
			t.Errorf("cannot GET Parameter by Value: %v - %v", parameters.Value, err)
		}
		if reqInf.StatusCode != http.StatusOK {
			t.Errorf("Expected 200 status code, got %v", reqInf.StatusCode)
		}
	}
}

func GetTestParametersByName(t *testing.T) {
	for _, parameters := range testData.Parameters {
		resp, reqInf, err := TOSession.GetParametersByName(parameters.Name, nil)
		if err != nil {
			t.Errorf("cannot GET Parameter by Name: %v - %v", parameters.Name, err)
		}
		if reqInf.StatusCode != http.StatusOK {
			t.Errorf("Expected 200 status code, got %v", reqInf.StatusCode)
		}
		if len(resp) <= 0 {
			t.Errorf("No data available for Get Parameters by Name")
		}
	}
}

func GetParametersByInvalidId(t *testing.T) {
	resp, _, err := TOSession.GetParameterByID(10000, nil)
	if err != nil {
		t.Errorf("Getting Parameters by Invalid ID %v", err)
	}
	if len(resp) >= 1 {
		t.Errorf("Invalid ID shouldn't have any response %v Error %v", resp, err)
	}
}

func GetParametersByInvalidName(t *testing.T) {
	resp, _, err := TOSession.GetParametersByName("abcd", nil)
	if err != nil {
		t.Errorf("Getting Parameters by Invalid Name %v", err)
	}
	if len(resp) >= 1 {
		t.Errorf("Invalid name shouldn't have any response %v Error %v", resp, err)
	}
}

func GetParametersByInvalidConfigfile(t *testing.T) {
	resp, _, err := TOSession.GetParametersByConfigFile("abcd", nil)
	if err != nil {
		t.Errorf("Getting Parameters by Invalid ConfigFile %v", err)
	}
	if len(resp) >= 1 {
		t.Errorf("Invalid config file shouldn't have any response %v Error %v", resp, err)
	}
}

func GetParametersByInvalidValue(t *testing.T) {
	resp, _, err := TOSession.GetParametersByValue("abcd", nil)
	if err != nil {
		t.Errorf("Getting Parameters by Invalid value %v", err)
	}
	if len(resp) >= 1 {
		t.Errorf("Invalid value shouldn't have any response %v Error %v", resp, err)
	}
}

func CreateTestParametersAlreadyExist(t *testing.T) {
	resp, _, _ := TOSession.GetParameters(nil, nil)
	_, reqInf, _ := TOSession.CreateParameter(resp[0])
	if reqInf.StatusCode != http.StatusBadRequest{
		t.Errorf("Expected 400 status code, got %v", reqInf.StatusCode)
	}
}

func CreateTestParametersMissingName(t *testing.T) {
	firstParameter := testData.Parameters[0]
	firstParameter.Name = ""
	_, reqInf, _ := TOSession.CreateParameter(firstParameter)
	if reqInf.StatusCode != http.StatusBadRequest{
		t.Errorf("Expected 400 status code, got %v", reqInf.StatusCode)
	}
}

func CreateTestParametersMissingconfigFile(t *testing.T) {
	firstParameter := testData.Parameters[0]
	firstParameter.ConfigFile = ""
	_, reqInf, _ := TOSession.CreateParameter(firstParameter)
	if reqInf.StatusCode != http.StatusBadRequest{
		t.Errorf("Expected 400 status code, got %v", reqInf.StatusCode)
	}
}
