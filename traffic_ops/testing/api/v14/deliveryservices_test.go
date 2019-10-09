package v14

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

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	toclient "github.com/apache/trafficcontrol/traffic_ops/client"
)

func TestDeliveryServices(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Users, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, DeliveryServices}, func() {
		UpdateTestDeliveryServices(t)
		UpdateNullableTestDeliveryServices(t)
		GetTestDeliveryServices(t)
		DeliveryServiceMinorVersionsTest(t)
		DeliveryServiceTenancyTest(t)
	})
}

func CreateTestDeliveryServices(t *testing.T) {
	pl := tc.Parameter{
		ConfigFile: "remap.config",
		Name:       "location",
		Value:      "/remap/config/location/parameter/",
	}
	_, _, err := TOSession.CreateParameter(pl)
	if err != nil {
		t.Errorf("cannot create parameter: %v\n", err)
	}
	for _, ds := range testData.DeliveryServices {
		_, err = TOSession.CreateDeliveryService(&ds)
		if err != nil {
			t.Errorf("could not CREATE delivery service '%s': %v\n", ds.XMLID, err)
		}
	}
}

func GetTestDeliveryServices(t *testing.T) {
	actualDSes, _, err := TOSession.GetDeliveryServices()
	if err != nil {
		t.Errorf("cannot GET DeliveryServices: %v - %v\n", err, actualDSes)
	}
	actualDSMap := map[string]tc.DeliveryService{}
	for _, ds := range actualDSes {
		actualDSMap[ds.XMLID] = ds
	}
	cnt := 0
	for _, ds := range testData.DeliveryServices {
		if _, ok := actualDSMap[ds.XMLID]; !ok {
			t.Errorf("GET DeliveryService missing: %v\n", ds.XMLID)
		}
		// exactly one ds should have exactly 3 query params. the rest should have none
		if c := len(ds.ConsistentHashQueryParams); c > 0 {
			if c != 3 {
				t.Errorf("deliveryservice %s has %d query params; expected %d or %d", ds.XMLID, c, 3, 0)
			}
			cnt++
		}
	}
	if cnt > 2 {
		t.Errorf("exactly 2 deliveryservices should have more than one query param; found %d", cnt)
	}
}

func UpdateTestDeliveryServices(t *testing.T) {
	firstDS := testData.DeliveryServices[0]

	dses, _, err := TOSession.GetDeliveryServices()
	if err != nil {
		t.Errorf("cannot GET Delivery Services: %v\n", err)
	}

	remoteDS := tc.DeliveryService{}
	found := false
	for _, ds := range dses {
		if ds.XMLID == firstDS.XMLID {
			found = true
			remoteDS = ds
			break
		}
	}
	if !found {
		t.Errorf("GET Delivery Services missing: %v\n", firstDS.XMLID)
	}

	updatedLongDesc := "something different"
	updatedMaxDNSAnswers := 164598
	updatedMaxOriginConnections := 100
	remoteDS.LongDesc = updatedLongDesc
	remoteDS.MaxDNSAnswers = updatedMaxDNSAnswers
	remoteDS.MaxOriginConnections = updatedMaxOriginConnections
	remoteDS.MatchList = nil // verify that this field is optional in a PUT request, doesn't cause nil dereference panic

	if updateResp, err := TOSession.UpdateDeliveryService(strconv.Itoa(remoteDS.ID), &remoteDS); err != nil {
		t.Errorf("cannot UPDATE DeliveryService by ID: %v - %v\n", err, updateResp)
	}

	// Retrieve the server to check rack and interfaceName values were updated
	resp, _, err := TOSession.GetDeliveryService(strconv.Itoa(remoteDS.ID))
	if err != nil {
		t.Errorf("cannot GET Delivery Service by ID: %v - %v\n", remoteDS.XMLID, err)
	}
	if resp == nil {
		t.Errorf("cannot GET Delivery Service by ID: %v - nil\n", remoteDS.XMLID)
	}

	if resp.LongDesc != updatedLongDesc || resp.MaxDNSAnswers != updatedMaxDNSAnswers || resp.MaxOriginConnections != updatedMaxOriginConnections {
		t.Errorf("results do not match actual: %s, expected: %s\n", resp.LongDesc, updatedLongDesc)
		t.Errorf("results do not match actual: %v, expected: %v\n", resp.MaxDNSAnswers, updatedMaxDNSAnswers)
		t.Errorf("results do not match actual: %v, expected: %v\n", resp.MaxOriginConnections, updatedMaxOriginConnections)
	}
}

func UpdateNullableTestDeliveryServices(t *testing.T) {
	firstDS := testData.DeliveryServices[0]

	dses, _, err := TOSession.GetDeliveryServicesNullable()
	if err != nil {
		t.Fatalf("cannot GET Delivery Services: %v\n", err)
	}

	remoteDS := tc.DeliveryServiceNullable{}
	found := false
	for _, ds := range dses {
		if ds.XMLID == nil || ds.ID == nil {
			continue
		}
		if *ds.XMLID == firstDS.XMLID {
			found = true
			remoteDS = ds
			break
		}
	}
	if !found {
		t.Fatalf("GET Delivery Services missing: %v\n", firstDS.XMLID)
	}

	updatedLongDesc := "something else different"
	updatedMaxDNSAnswers := 164599
	remoteDS.LongDesc = &updatedLongDesc
	remoteDS.MaxDNSAnswers = &updatedMaxDNSAnswers

	if updateResp, err := TOSession.UpdateDeliveryServiceNullable(strconv.Itoa(*remoteDS.ID), &remoteDS); err != nil {
		t.Fatalf("cannot UPDATE DeliveryService by ID: %v - %v\n", err, updateResp)
	}

	// Retrieve the server to check rack and interfaceName values were updated
	resp, _, err := TOSession.GetDeliveryServiceNullable(strconv.Itoa(*remoteDS.ID))
	if err != nil {
		t.Fatalf("cannot GET Delivery Service by ID: %v - %v\n", remoteDS.XMLID, err)
	}
	if resp == nil {
		t.Fatalf("cannot GET Delivery Service by ID: %v - nil\n", remoteDS.XMLID)
	}

	if resp.LongDesc == nil || resp.MaxDNSAnswers == nil {
		t.Errorf("results do not match actual: %v, expected: %s\n", resp.LongDesc, updatedLongDesc)
		t.Fatalf("results do not match actual: %v, expected: %d\n", resp.MaxDNSAnswers, updatedMaxDNSAnswers)
	}

	if *resp.LongDesc != updatedLongDesc || *resp.MaxDNSAnswers != updatedMaxDNSAnswers {
		t.Errorf("results do not match actual: %s, expected: %s\n", *resp.LongDesc, updatedLongDesc)
		t.Fatalf("results do not match actual: %d, expected: %d\n", *resp.MaxDNSAnswers, updatedMaxDNSAnswers)
	}
}

func DeleteTestDeliveryServices(t *testing.T) {
	dses, _, err := TOSession.GetDeliveryServices()
	if err != nil {
		t.Errorf("cannot GET deliveryservices: %v\n", err)
	}
	for _, testDS := range testData.DeliveryServices {
		ds := tc.DeliveryService{}
		found := false
		for _, realDS := range dses {
			if realDS.XMLID == testDS.XMLID {
				ds = realDS
				found = true
				break
			}
		}
		if !found {
			t.Errorf("DeliveryService not found in Traffic Ops: %v\n", ds.XMLID)
		}

		delResp, err := TOSession.DeleteDeliveryService(strconv.Itoa(ds.ID))
		if err != nil {
			t.Errorf("cannot DELETE DeliveryService by ID: %v - %v\n", err, delResp)
		}

		// Retrieve the Server to see if it got deleted
		foundDS, err := TOSession.DeliveryService(strconv.Itoa(ds.ID))
		if err == nil && foundDS != nil {
			t.Errorf("expected Delivery Service: %s to be deleted\n", ds.XMLID)
		}
	}

	// clean up parameter created in CreateTestDeliveryServices()
	params, _, err := TOSession.GetParameterByNameAndConfigFile("location", "remap.config")
	for _, param := range params {
		deleted, _, err := TOSession.DeleteParameterByID(param.ID)
		if err != nil {
			t.Errorf("cannot DELETE parameter by ID (%d): %v - %v\n", param.ID, err, deleted)
		}
	}
}

func DeliveryServiceMinorVersionsTest(t *testing.T) {
	testDS := testData.DeliveryServices[4]
	if testDS.XMLID != "ds-test-minor-versions" {
		t.Errorf("expected XMLID: ds-test-minor-versions, actual: %s\n", testDS.XMLID)
	}

	dses, _, err := TOSession.GetDeliveryServicesNullable()
	if err != nil {
		t.Errorf("cannot GET DeliveryServices: %v - %v\n", err, dses)
	}
	ds := tc.DeliveryServiceNullable{}
	for _, d := range dses {
		if *d.XMLID == testDS.XMLID {
			ds = d
			break
		}
	}
	// GET latest, verify expected values for 1.3 and 1.4 fields
	if ds.DeepCachingType == nil {
		t.Errorf("expected DeepCachingType: %s, actual: nil\n", testDS.DeepCachingType.String())
	} else if *ds.DeepCachingType != testDS.DeepCachingType {
		t.Errorf("expected DeepCachingType: %s, actual: %s\n", testDS.DeepCachingType.String(), ds.DeepCachingType.String())
	}
	if ds.FQPacingRate == nil {
		t.Errorf("expected FQPacingRate: %d, actual: nil\n", testDS.FQPacingRate)
	} else if *ds.FQPacingRate != testDS.FQPacingRate {
		t.Errorf("expected FQPacingRate: %d, actual: %d\n", testDS.FQPacingRate, *ds.FQPacingRate)
	}
	if ds.SigningAlgorithm == nil {
		t.Errorf("expected SigningAlgorithm: %s, actual: nil\n", testDS.SigningAlgorithm)
	} else if *ds.SigningAlgorithm != testDS.SigningAlgorithm {
		t.Errorf("expected SigningAlgorithm: %s, actual: %s\n", testDS.SigningAlgorithm, *ds.SigningAlgorithm)
	}
	if ds.Tenant == nil {
		t.Errorf("expected Tenant: %s, actual: nil\n", testDS.Tenant)
	} else if *ds.Tenant != testDS.Tenant {
		t.Errorf("expected Tenant: %s, actual: %s\n", testDS.Tenant, *ds.Tenant)
	}
	if ds.TRRequestHeaders == nil {
		t.Errorf("expected TRRequestHeaders: %s, actual: nil\n", testDS.TRRequestHeaders)
	} else if *ds.TRRequestHeaders != testDS.TRRequestHeaders {
		t.Errorf("expected TRRequestHeaders: %s, actual: %s\n", testDS.TRRequestHeaders, *ds.TRRequestHeaders)
	}
	if ds.TRResponseHeaders == nil {
		t.Errorf("expected TRResponseHeaders: %s, actual: nil\n", testDS.TRResponseHeaders)
	} else if *ds.TRResponseHeaders != testDS.TRResponseHeaders {
		t.Errorf("expected TRResponseHeaders: %s, actual: %s\n", testDS.TRResponseHeaders, *ds.TRResponseHeaders)
	}
	if ds.ConsistentHashRegex == nil {
		t.Errorf("expected ConsistentHashRegex: %s, actual: nil\n", testDS.ConsistentHashRegex)
	} else if *ds.ConsistentHashRegex != testDS.ConsistentHashRegex {
		t.Errorf("expected ConsistentHashRegex: %s, actual: %s\n", testDS.ConsistentHashRegex, *ds.ConsistentHashRegex)
	}
	if ds.ConsistentHashQueryParams == nil {
		t.Errorf("expected ConsistentHashQueryParams: %v, actual: nil\n", testDS.ConsistentHashQueryParams)
	} else if !reflect.DeepEqual(ds.ConsistentHashQueryParams, testDS.ConsistentHashQueryParams) {
		t.Errorf("expected ConsistentHashQueryParams: %v, actual: %v\n", testDS.ConsistentHashQueryParams, ds.ConsistentHashQueryParams)
	}
	if ds.MaxOriginConnections == nil {
		t.Errorf("expected MaxOriginConnections: %d, actual: nil\n", testDS.MaxOriginConnections)
	} else if *ds.MaxOriginConnections != testDS.MaxOriginConnections {
		t.Errorf("expected MaxOriginConnections: %d, actual: %d\n", testDS.MaxOriginConnections, *ds.MaxOriginConnections)
	}

	// GET 1.1, verify 1.3 and 1.4 fields are nil
	data := tc.DeliveryServicesNullableResponse{}
	if err = makeV11Request(http.MethodGet, "deliveryservices/"+strconv.Itoa(*ds.ID), nil, &data); err != nil {
		t.Errorf("cannot GET 1.1 deliveryservice: %s\n", err.Error())
	}
	respDS := data.Response[0]
	if !dsV13FieldsAreNil(respDS) || !dsV14FieldsAreNil(respDS) {
		t.Errorf("expected 1.3 and 1.4 values to be nil, actual: non-nil")
	}

	// GET 1.3, verify 1.3 fields are non-nil and 1.4 fields are nil
	data = tc.DeliveryServicesNullableResponse{}
	if err = makeV13Request(http.MethodGet, "deliveryservices/"+strconv.Itoa(*ds.ID), nil, &data); err != nil {
		t.Errorf("cannot GET 1.3 deliveryservice: %s\n", err.Error())
	}
	respDS = data.Response[0]
	if dsV13FieldsAreNil(respDS) {
		t.Errorf("expected 1.3 values to be non-nil, actual: nil\n")
	}
	if !dsV14FieldsAreNil(respDS) {
		t.Errorf("expected 1.4 values to be nil, actual: non-nil")
	}
	if _, err = TOSession.DeleteDeliveryService(strconv.Itoa(*ds.ID)); err != nil {
		t.Errorf("cannot DELETE deliveryservice: %s\n", err.Error())
	}

	ds.ID = nil
	dsBody, err := json.Marshal(ds)
	if err != nil {
		t.Errorf("cannot POST deliveryservice, failed to marshal JSON: %s\n", err.Error())
	}
	dsV11Body, err := json.Marshal(ds.DeliveryServiceNullableV11)
	if err != nil {
		t.Errorf("cannot POST deliveryservice, failed to marshal JSON: %s\n", err.Error())
	}

	// POST 1.3 w/ 1.4 data, verify 1.4 fields were ignored
	postDSResp := tc.CreateDeliveryServiceNullableResponse{}
	if err = makeV13Request(http.MethodPost, "deliveryservices", bytes.NewBuffer(dsBody), &postDSResp); err != nil {
		t.Errorf("cannot POST 1.3 deliveryservice, failed to make request: %s\n", err.Error())
	}
	if !dsV14FieldsAreNil(postDSResp.Response[0]) {
		t.Errorf("POST 1.3 expected 1.4 values to be nil, actual: non-nil")
	}
	respID := postDSResp.Response[0].ID
	getDS, _, err := TOSession.GetDeliveryServiceNullable(strconv.Itoa(*respID))
	if err != nil {
		t.Errorf("cannot GET deliveryservice: %s\n", err.Error())
	}
	if !dsV14FieldsAreNilOrDefault(*getDS) {
		t.Errorf("POST 1.3 expected 1.4 values to be nil/default, actual: non-nil/default")
	}
	if _, err = TOSession.DeleteDeliveryService(strconv.Itoa(*respID)); err != nil {
		t.Errorf("cannot DELETE deliveryservice: %s\n", err.Error())
	}

	// POST 1.1 w/ 1.4 data, verify 1.3 and 1.4 fields were ignored
	postDSResp = tc.CreateDeliveryServiceNullableResponse{}
	if err = makeV11Request(http.MethodPost, "deliveryservices", bytes.NewBuffer(dsBody), &postDSResp); err != nil {
		t.Errorf("cannot POST 1.1 deliveryservice, failed to make request: %s\n", err.Error())
	}
	if !dsV13FieldsAreNil(postDSResp.Response[0]) || !dsV14FieldsAreNil(postDSResp.Response[0]) {
		t.Errorf("POST 1.1 expected 1.3 and 1.4 values to be nil, actual: non-nil %++v\n", postDSResp.Response[0])
	}
	respID = postDSResp.Response[0].ID
	getDS, _, err = TOSession.GetDeliveryServiceNullable(strconv.Itoa(*respID))
	if err != nil {
		t.Errorf("cannot GET deliveryservice: %s\n", err.Error())
	}
	if !dsV13FieldsAreNilOrDefault(*getDS) || !dsV14FieldsAreNilOrDefault(*getDS) {
		t.Errorf("POST 1.1 expected 1.3 and 1.4 values to be nil/default, actual: non-nil/default %++v\n", *getDS)
	}

	// PUT 1.4 w/ 1.4 data, then verify that a PUT 1.1 with 1.1 data preserves the existing 1.3 and 1.4 data
	if _, err = TOSession.UpdateDeliveryServiceNullable(strconv.Itoa(*respID), &ds); err != nil {
		t.Errorf("cannot PUT deliveryservice: %s\n", err.Error())
	}
	putDSResp := tc.UpdateDeliveryServiceNullableResponse{}
	if err = makeV11Request(http.MethodPut, "deliveryservices/"+strconv.Itoa(*respID), bytes.NewBuffer(dsV11Body), &putDSResp); err != nil {
		t.Errorf("cannot PUT 1.1 deliveryservice, failed to make request: %s\n", err.Error())
	}
	if !dsV13FieldsAreNil(putDSResp.Response[0]) || !dsV14FieldsAreNil(putDSResp.Response[0]) {
		t.Errorf("PUT 1.1 expected 1.3 and 1.4 values to be nil, actual: non-nil %++v\n", putDSResp.Response[0])
	}
	getDS, _, err = TOSession.GetDeliveryServiceNullable(strconv.Itoa(*respID))
	if err != nil {
		t.Errorf("cannot GET deliveryservice: %s\n", err.Error())
	}
	if getDS.FQPacingRate == nil {
		t.Errorf("expected FQPacingRate: %d, actual: nil\n", testDS.FQPacingRate)
	} else if *getDS.FQPacingRate != testDS.FQPacingRate {
		t.Errorf("expected FQPacingRate: %d, actual: %d\n", testDS.FQPacingRate, *getDS.FQPacingRate)
	}
	if getDS.MaxOriginConnections == nil {
		t.Errorf("expected MaxOriginConnections: %d, actual: nil\n", testDS.MaxOriginConnections)
	} else if *getDS.MaxOriginConnections != testDS.MaxOriginConnections {
		t.Errorf("expected MaxOriginConnections: %d, actual: %d\n", testDS.MaxOriginConnections, *getDS.MaxOriginConnections)
	}

	// PUT 1.3 w/ 1.1 data, verify that 1.4 fields were preserved
	putDSResp = tc.UpdateDeliveryServiceNullableResponse{}
	if err = makeV13Request(http.MethodPut, "deliveryservices/"+strconv.Itoa(*respID), bytes.NewBuffer(dsV11Body), &putDSResp); err != nil {
		t.Errorf("cannot PUT 1.3 deliveryservice, failed to make request: %s\n", err.Error())
	}
	if !dsV14FieldsAreNil(putDSResp.Response[0]) {
		t.Errorf("PUT 1.3 expected 1.4 values to be nil, actual: non-nil %++v\n", putDSResp.Response[0])
	}
	getDS, _, err = TOSession.GetDeliveryServiceNullable(strconv.Itoa(*respID))
	if err != nil {
		t.Errorf("cannot GET deliveryservice: %s\n", err.Error())
	}
	if getDS.MaxOriginConnections == nil {
		t.Errorf("expected MaxOriginConnections: %d, actual: nil\n", testDS.MaxOriginConnections)
	} else if *getDS.MaxOriginConnections != testDS.MaxOriginConnections {
		t.Errorf("expected MaxOriginConnections: %d, actual: %d\n", testDS.MaxOriginConnections, *getDS.MaxOriginConnections)
	}

	// DELETE+POST 1.1 again, so that 1.3 and 1.4 fields are back to nil/default
	if _, err = TOSession.DeleteDeliveryService(strconv.Itoa(*respID)); err != nil {
		t.Errorf("cannot DELETE deliveryservice: %s\n", err.Error())
	}
	postDSResp = tc.CreateDeliveryServiceNullableResponse{}
	if err = makeV11Request(http.MethodPost, "deliveryservices", bytes.NewBuffer(dsV11Body), &postDSResp); err != nil {
		t.Errorf("cannot POST 1.1 deliveryservice, failed to make request: %s\n", err.Error())
	}
	respID = postDSResp.Response[0].ID

	// PUT 1.1 w/ 1.4 data - make sure 1.3 and 1.4 fields were ignored
	putDSResp = tc.UpdateDeliveryServiceNullableResponse{}
	if err = makeV11Request(http.MethodPut, "deliveryservices/"+strconv.Itoa(*respID), bytes.NewBuffer(dsBody), &putDSResp); err != nil {
		t.Errorf("cannot PUT 1.1 deliveryservice, failed to make request: %s\n", err.Error())
	}
	if !dsV13FieldsAreNil(putDSResp.Response[0]) || !dsV14FieldsAreNil(putDSResp.Response[0]) {
		t.Errorf("PUT 1.1 expected 1.3 and 1.4 values to be nil, actual: non-nil %++v\n", putDSResp.Response[0])
	}
	respID = putDSResp.Response[0].ID
	getDS, _, err = TOSession.GetDeliveryServiceNullable(strconv.Itoa(*respID))
	if err != nil {
		t.Errorf("cannot GET deliveryservice: %s\n", err.Error())
	}
	if !dsV13FieldsAreNilOrDefault(*getDS) || !dsV14FieldsAreNilOrDefault(*getDS) {
		t.Errorf("PUT 1.1 expected 1.3 and 1.4 values to be nil/default, actual: non-nil/default %++v\n", *getDS)
	}

	// PUT 1.3 w/ 1.4 data, make sure 1.4 fields were ignored
	putDSResp = tc.UpdateDeliveryServiceNullableResponse{}
	if err = makeV13Request(http.MethodPut, "deliveryservices/"+strconv.Itoa(*respID), bytes.NewBuffer(dsBody), &putDSResp); err != nil {
		t.Errorf("cannot PUT 1.1 deliveryservice, failed to make request: %s\n", err.Error())
	}
	if !dsV14FieldsAreNil(putDSResp.Response[0]) {
		t.Errorf("PUT 1.3 expected 1.4 values to be nil, actual: non-nil\n")
	}
	respID = putDSResp.Response[0].ID
	getDS, _, err = TOSession.GetDeliveryServiceNullable(strconv.Itoa(*respID))
	if err != nil {
		t.Errorf("cannot GET deliveryservice: %s\n", err.Error())
	}
	if !dsV14FieldsAreNilOrDefault(*getDS) {
		t.Errorf("PUT 1.3 expected 1.4 values to be nil/default, actual: non-nil/default\n")
	}
}

func dsV13FieldsAreNilOrDefault(ds tc.DeliveryServiceNullable) bool {
	return (ds.DeepCachingType == nil || *ds.DeepCachingType == tc.DeepCachingTypeNever) &&
		(ds.FQPacingRate == nil || *ds.FQPacingRate == 0) &&
		(ds.TRRequestHeaders == nil || *ds.TRRequestHeaders == "") &&
		(ds.TRResponseHeaders == nil || *ds.TRResponseHeaders == "")
}

func dsV14FieldsAreNilOrDefault(ds tc.DeliveryServiceNullable) bool {
	return (ds.ConsistentHashRegex == nil || *ds.ConsistentHashRegex == "") &&
		(ds.ConsistentHashQueryParams == nil || len(ds.ConsistentHashQueryParams) == 0) &&
		(ds.MaxOriginConnections == nil || *ds.MaxOriginConnections == 0)
}

func dsV13FieldsAreNil(ds tc.DeliveryServiceNullable) bool {
	return ds.DeepCachingType == nil &&
		ds.FQPacingRate == nil &&
		ds.SigningAlgorithm == nil &&
		ds.Tenant == nil &&
		ds.TRRequestHeaders == nil &&
		ds.TRResponseHeaders == nil
}

func dsV14FieldsAreNil(ds tc.DeliveryServiceNullable) bool {
	return ds.ConsistentHashRegex == nil &&
		(ds.ConsistentHashQueryParams == nil || len(ds.ConsistentHashQueryParams) == 0) &&
		ds.MaxOriginConnections == nil
}

func makeV11Request(method string, path string, body io.Reader, respStruct interface{}) error {
	return makeRequest("1.1", method, path, body, respStruct)
}

func makeV13Request(method string, path string, body io.Reader, respStruct interface{}) error {
	return makeRequest("1.3", method, path, body, respStruct)
}

// TODO: move this helper function into a better location
func makeRequest(version string, method string, path string, body io.Reader, respStruct interface{}) error {
	req, err := http.NewRequest(method, TOSession.URL+"/api/"+version+"/"+path, body)
	if err != nil {
		return fmt.Errorf("failed to create request: %s", err.Error())
	}
	resp, err := TOSession.Client.Do(req)
	if err != nil {
		return fmt.Errorf("running request: %s", err.Error())
	}
	defer resp.Body.Close()
	bts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading body: " + err.Error())
	}
	if err = json.Unmarshal(bts, respStruct); err != nil {
		return fmt.Errorf("unmarshalling body '" + string(bts) + "': " + err.Error())
	}
	return nil
}

func DeliveryServiceTenancyTest(t *testing.T) {
	dses, _, err := TOSession.GetDeliveryServicesNullable()
	if err != nil {
		t.Errorf("cannot GET deliveryservices: %v\n", err)
	}
	tenant3DS := tc.DeliveryServiceNullable{}
	foundTenant3DS := false
	for _, d := range dses {
		if *d.XMLID == "ds3" {
			tenant3DS = d
			foundTenant3DS = true
		}
	}
	if !foundTenant3DS || *tenant3DS.Tenant != "tenant3" {
		t.Error("expected to find deliveryservice 'ds3' with tenant 'tenant3'")
	}

	toReqTimeout := time.Second * time.Duration(Config.Default.Session.TimeoutInSecs)
	tenant4TOClient, _, err := toclient.LoginWithAgent(TOSession.URL, "tenant4user", "pa$$word", true, "to-api-v14-client-tests/tenant4user", true, toReqTimeout)
	if err != nil {
		t.Fatalf("failed to log in with tenant4user: %v", err.Error())
	}

	dsesReadableByTenant4, _, err := tenant4TOClient.GetDeliveryServicesNullable()
	if err != nil {
		t.Error("tenant4user cannot GET deliveryservices")
	}

	// assert that tenant4user cannot read deliveryservices outside of its tenant
	for _, ds := range dsesReadableByTenant4 {
		if *ds.XMLID == "ds3" {
			t.Error("expected tenant4 to be unable to read delivery services from tenant 3")
		}
	}

	// assert that tenant4user cannot update tenant3user's deliveryservice
	if _, err = tenant4TOClient.UpdateDeliveryServiceNullable(string(*tenant3DS.ID), &tenant3DS); err == nil {
		t.Errorf("expected tenant4user to be unable to update tenant3's deliveryservice (%s)", *tenant3DS.XMLID)
	}

	// assert that tenant4user cannot delete tenant3user's deliveryservice
	if _, err = tenant4TOClient.DeleteDeliveryService(string(*tenant3DS.ID)); err == nil {
		t.Errorf("expected tenant4user to be unable to delete tenant3's deliveryservice (%s)", *tenant3DS.XMLID)
	}

	// assert that tenant4user cannot create a deliveryservice outside of its tenant
	tenant3DS.XMLID = util.StrPtr("deliveryservicetenancytest")
	tenant3DS.DisplayName = util.StrPtr("deliveryservicetenancytest")
	if _, err = tenant4TOClient.CreateDeliveryServiceNullable(&tenant3DS); err == nil {
		t.Errorf("expected tenant4user to be unable to create a deliveryservice outside of its tenant")
	}

}
