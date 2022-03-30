package v4

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
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-rfc"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/assert"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
	client "github.com/apache/trafficcontrol/traffic_ops/v4-client"
)

func TestDeliveryServicesRequiredCapabilities(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Users, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, ServerCapabilities, Topologies, ServiceCategories, DeliveryServices, DeliveryServicesRequiredCapabilities, DeliveryServiceServerAssignments, ServerServerCapabilities}, func() {

		tomorrow := time.Now().AddDate(0, 0, 1).Format(time.RFC1123)
		currentTime := time.Now().UTC().Add(-15 * time.Second)
		currentTimeRFC := currentTime.Format(time.RFC1123)

		methodTests := utils.V4TestCase{
			"GET": {
				"NOT MODIFIED when NO CHANGES made": {
					ClientSession: TOSession, RequestOpts: client.RequestOptions{Header: http.Header{rfc.IfModifiedSince: {tomorrow}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"OK when VALID request": {
					ClientSession: TOSession, Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"OK when VALID DELIVERYSERVICEID parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"deliveryServiceId": {"ds1"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateDSRCExpectedFields(map[string]interface{}{"DeliveryServiceId": "ds1"})),
				},
				"OK when VALID XMLID parameter": {
					ClientSession: TOSession, RequestOpts: client.RequestOptions{QueryParameters: url.Values{"xmlID": {"ds2"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateDSRCExpectedFields(map[string]interface{}{"XMLID": "ds2"})),
				},
				"OK when VALID REQUIREDCAPABILITY parameter": {
					ClientSession: TOSession, RequestOpts: client.RequestOptions{QueryParameters: url.Values{"requiredCapability": {"bar"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateDSRCExpectedFields(map[string]interface{}{"RequiredCapability": "bar"})),
				},
				"FIRST RESULT when LIMIT=1": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"requiredCapability"}, "limit": {"1"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateDSRCPagination("limit")),
				},
				"SECOND RESULT when LIMIT=1 OFFSET=1": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"requiredCapability"}, "limit": {"1"}, "offset": {"1"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateDSRCPagination("offset")),
				},
				"SECOND RESULT when LIMIT=1 PAGE=2": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"requiredCapability"}, "limit": {"1"}, "page": {"2"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateDSRCPagination("page")),
				},
				"BAD REQUEST when INVALID LIMIT parameter": {
					ClientSession: TOSession, RequestOpts: client.RequestOptions{QueryParameters: url.Values{"limit": {"-2"}}},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when INVALID OFFSET parameter": {
					ClientSession: TOSession, RequestOpts: client.RequestOptions{QueryParameters: url.Values{"limit": {"1"}, "offset": {"0"}}},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when INVALID PAGE parameter": {
					ClientSession: TOSession, RequestOpts: client.RequestOptions{QueryParameters: url.Values{"limit": {"1"}, "page": {"0"}}},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
			},
			"POST": {
				"BAD REQUEST when REASSIGNING REQUIRED CAPABILITY to DELIVERY SERVICE": {
					ClientSession: TOSession, RequestBody: map[string]interface{}{
						"deliveryServiceID":  GetDeliveryServiceId(t, "ds1")(),
						"RequiredCapability": "foo",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when SERVERS DONT have CAPABILITY": {
					ClientSession: TOSession, RequestBody: map[string]interface{}{
						"deliveryServiceID":  GetDeliveryServiceId(t, "ds3")(),
						"RequiredCapability": "bar",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when DELIVERY SERVICE HAS TOPOLOGY where SERVERS DONT have CAPABILITY": {
					ClientSession: TOSession, RequestBody: map[string]interface{}{
						"deliveryServiceID":  GetDeliveryServiceId(t, "ds-top-req-cap")(),
						"RequiredCapability": "bar",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when DELIVERY SERVICE ID EMPTY": {
					ClientSession: TOSession, RequestBody: map[string]interface{}{
						"requiredCapability": "bar",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when REQUIRED CAPABILITY EMPTY": {
					ClientSession: TOSession, RequestBody: map[string]interface{}{
						"deliveryServiceID": GetDeliveryServiceId(t, "ds1")(),
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"NOT FOUND when NON-EXISTENT REQUIRED CAPABILITY": {
					ClientSession: TOSession, RequestBody: map[string]interface{}{
						"deliveryServiceID":  GetDeliveryServiceId(t, "ds1")(),
						"requiredCapability": "bogus",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
				"NOT FOUND when NON-EXISTENT DELIVERY SERVICE ID": {
					ClientSession: TOSession, RequestBody: map[string]interface{}{
						"deliveryServiceID":  -1,
						"requiredCapability": "foo",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
				"BAD REQUEST when INVALID DELIVERY SERVICE TYPE": {
					ClientSession: TOSession, RequestBody: map[string]interface{}{
						"deliveryServiceID":  GetDeliveryServiceId(t, "anymap-ds")(),
						"requiredCapability": "foo",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
			},
			"DELETE": {
				"OK when VALID request": {
					EndpointId: GetDeliveryServiceId(t, "msods1"), ClientSession: TOSession,
					RequestBody:  map[string]interface{}{"requiredCapability": "bar"},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"NOT FOUND when NON-EXISTENT DELIVERYSERVICEID parameter": {
					EndpointId: func() int { return -1 }, ClientSession: TOSession,
					RequestBody:  map[string]interface{}{"requiredCapability": "foo"},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
				"NOT FOUND when NON-EXISTENT REQUIREDCAPABILITY parameter": {
					EndpointId: GetDeliveryServiceId(t, "ds1"), ClientSession: TOSession,
					RequestBody:  map[string]interface{}{"requiredCapability": "bogus"},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
			},
			"GET AFTER CHANGES": {
				"OK when CHANGES made": {
					ClientSession: TOSession, RequestOpts: client.RequestOptions{Header: http.Header{rfc.IfModifiedSince: {currentTimeRFC}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					dsrc := tc.DeliveryServicesRequiredCapability{}

					if val, ok := testCase.RequestOpts.QueryParameters["deliveryServiceId"]; ok {
						if _, err := strconv.Atoi(val[0]); err != nil {
							testCase.RequestOpts.QueryParameters.Set("deliveryServiceId", strconv.Itoa(GetDeliveryServiceId(t, val[0])()))
						}
					}

					if testCase.RequestBody != nil {
						dat, err := json.Marshal(testCase.RequestBody)
						assert.NoError(t, err, "Error occurred when marshalling request body: %v", err)
						err = json.Unmarshal(dat, &dsrc)
						assert.NoError(t, err, "Error occurred when unmarshalling request body: %v", err)
					}

					switch method {
					case "GET", "GET AFTER CHANGES":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.GetDeliveryServicesRequiredCapabilities(testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.CreateDeliveryServicesRequiredCapability(dsrc, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, resp, err)
							}
						})
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.DeleteDeliveryServicesRequiredCapability(testCase.EndpointId(), testCase.RequestBody["requiredCapability"].(string), testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					}
				}
			})
		}
	})

}

func validateDSRCExpectedFields(expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		dsrcResp := resp.([]tc.DeliveryServicesRequiredCapability)
		for field, expected := range expectedResp {
			for _, dsrc := range dsrcResp {
				switch field {
				case "DeliveryServiceID":
					assert.Equal(t, expected, *dsrc.DeliveryServiceID, "Expected deliveryServiceId to be %v, but got %v", expected, dsrc.DeliveryServiceID)
				case "XMLID":
					assert.Equal(t, expected, *dsrc.XMLID, "Expected xmlID to be %v, but got %v", expected, dsrc.XMLID)
				case "RequiredCapability":
					assert.Equal(t, expected, *dsrc.RequiredCapability, "Expected requiredCapability to be %v, but got %v", expected, dsrc.RequiredCapability)
				}
			}
		}
	}
}

func validateDSRCPagination(paginationParam string) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		paginationResp := resp.([]tc.DeliveryServicesRequiredCapability)

		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("orderby", "requiredCapability")
		respBase, _, err := TOSession.GetDeliveryServicesRequiredCapabilities(opts)
		assert.RequireNoError(t, err, "Cannot get Delivery Services Required Capabilities: %v - alerts: %+v", err, respBase.Alerts)

		dsrc := respBase.Response
		assert.RequireGreaterOrEqual(t, len(dsrc), 3, "Need at least 3 Delivery Services Required Capabilities in Traffic Ops to test pagination support, found: %d", len(dsrc))
		switch paginationParam {
		case "limit:":
			assert.Exactly(t, dsrc[:1], paginationResp, "Expected GET deliveryservices_required_capabilities with limit = 1 to return first result")
		case "offset":
			assert.Exactly(t, dsrc[1:2], paginationResp, "Expected GET deliveryservices_required_capabilities with limit = 1, offset = 1 to return second result")
		case "page":
			assert.Exactly(t, dsrc[1:2], paginationResp, "Expected GET deliveryservices_required_capabilities with limit = 1, page = 2 to return second result")
		}
	}
}

func OriginAssignTopologyBasedDeliveryServiceWithRequiredCapabilities(t *testing.T) {
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("xmlId", "ds-top-req-cap2")
	resp, _, err := TOSession.GetDeliveryServices(opts)
	if err != nil {
		t.Errorf("getting Delivery Service 'ds-top-req-cap2': %v - alerts: %+v", err, resp.Alerts)
	}
	if len(resp.Response) != 1 {
		t.Fatalf("expected to get only one Delivery Service with XMLID 'ds-top-req-cap2', but got %d", len(resp.Response))
	}
	if resp.Response[0].ID == nil {
		t.Fatal("no ID in the resulting delivery service")
	}
	dsID := *resp.Response[0].ID
	opts.QueryParameters = url.Values{}
	alerts, _, err := TOSession.AssignServersToDeliveryService([]string{"denver-mso-org-01"}, "ds-top-req-cap2", client.RequestOptions{})
	if err != nil {
		t.Errorf("assigning server 'denver-mso-org-01' to Delivery Service 'ds-top-req-cap2': %v - alerts: %+v", err, alerts)
	}
	opts.QueryParameters.Set("dsId", strconv.Itoa(dsID))
	opts.QueryParameters.Set("type", tc.OriginTypeName)
	responseServers, _, err := TOSession.GetServers(opts)
	if err != nil {
		t.Fatalf("getting servers for the 'ds-top-req-cap2' Delivery Service: %v - alerts: %+v", err, responseServers.Alerts)
	}
	if len(responseServers.Response) != 1 {
		t.Fatalf("expected just one ORG server in the response, but got %d", len(responseServers.Response))
	}
	if responseServers.Response[0].HostName == nil {
		t.Fatal("expected a valid host name for the resulting ORG server, but got nothing")
	}
	if *responseServers.Response[0].HostName != "denver-mso-org-01" {
		t.Errorf("expected host name of the resulting ORG server to be %v, but got %v", "denver-mso-org-01", *responseServers.Response[0].HostName)
	}
}

func CreateTestDeliveryServicesRequiredCapabilities(t *testing.T) {
	// Assign all required capability to delivery services listed in `tc-fixtures.json`.
	for _, dsrc := range testData.DeliveryServicesRequiredCapabilities {
		dsId := GetDeliveryServiceId(t, *dsrc.XMLID)()
		dsrc = tc.DeliveryServicesRequiredCapability{
			DeliveryServiceID:  &dsId,
			RequiredCapability: dsrc.RequiredCapability,
		}
		resp, _, err := TOSession.CreateDeliveryServicesRequiredCapability(dsrc, client.RequestOptions{})
		assert.NoError(t, err, "Unexpected error creating a Delivery Service/Required Capability relationship: %v - alerts: %+v", err, resp.Alerts)
	}
}

func DeleteTestDeliveryServicesRequiredCapabilities(t *testing.T) {
	// Get Required Capabilities to delete them
	dsrcs, _, err := TOSession.GetDeliveryServicesRequiredCapabilities(client.RequestOptions{})
	assert.NoError(t, err, "Error getting Delivery Service/Required Capability relationships: %v - alerts: %+v", err, dsrcs.Alerts)

	for _, dsrc := range dsrcs.Response {
		alerts, _, err := TOSession.DeleteDeliveryServicesRequiredCapability(*dsrc.DeliveryServiceID, *dsrc.RequiredCapability, client.RequestOptions{})
		assert.NoError(t, err, "Error deleting a relationship between a Delivery Service and a Capability: %v - alerts: %+v", err, alerts.Alerts)
	}
}

func helperGetDeliveryServiceID(t *testing.T, xmlID *string) *int {
	t.Helper()
	if xmlID == nil {
		t.Error("xml id must not be nil")
		return nil
	}
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("xmlId", *xmlID)
	ds, _, err := TOSession.GetDeliveryServices(opts)
	if err != nil {
		t.Errorf("Unexpected error getting Delivery Services filtered by XMLID '%s': %v - alerts: %+v", *xmlID, err, ds.Alerts)
		return nil
	}
	if len(ds.Response) != 1 {
		t.Errorf("Expected exactly one Delivery Service to have XMLID '%s', found: %d", *xmlID, len(ds.Response))
		return nil
	}
	return ds.Response[0].ID
}

func CreateTestTopologyBasedDeliveryServicesRequiredCapabilities(t *testing.T) {
	for _, td := range testData.TopologyBasedDeliveryServicesRequiredCapabilities {

		c := tc.DeliveryServicesRequiredCapability{
			DeliveryServiceID:  helperGetDeliveryServiceID(t, td.XMLID),
			RequiredCapability: td.RequiredCapability,
		}

		alerts, _, err := TOSession.CreateDeliveryServicesRequiredCapability(c, client.RequestOptions{})
		if err != nil {
			t.Fatalf("cannot create delivery service required capability: %v - %+v", err, alerts.Alerts)
		}
	}

	invalid := tc.DeliveryServicesRequiredCapability{
		DeliveryServiceID:  helperGetDeliveryServiceID(t, util.StrPtr("ds-top-req-cap")),
		RequiredCapability: util.StrPtr("asdf"),
	}
	_, reqInf, err := TOSession.CreateDeliveryServicesRequiredCapability(invalid, client.RequestOptions{})
	if err == nil {
		t.Fatal("when adding delivery service required capability to a delivery service with a topology that " +
			"doesn't have cachegroups with at least one server with the required capabilities - expected: error, actual: nil")
	}
	if reqInf.StatusCode != http.StatusBadRequest {
		t.Fatalf("when adding delivery service required capability to a delivery service with a topology that "+
			"doesn't have cachegroups with at least one server with the required capabilities - expected status code: "+
			"%d, actual: %d", http.StatusBadRequest, reqInf.StatusCode)
	}
}
