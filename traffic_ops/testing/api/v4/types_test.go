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
	"fmt"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/assert"
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
	"net/http"
	"net/url"
	"sort"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-rfc"
	tc "github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/utils"
	client "github.com/apache/trafficcontrol/traffic_ops/v4-client"
)

func TestTypes(t *testing.T) {
	WithObjs(t, []TCObj{Parameters, Types}, func() {

		currentTime := time.Now().UTC().Add(-15 * time.Second)
		currentTimeRFC := currentTime.Format(time.RFC1123)
		tomorrow := currentTime.AddDate(0, 0, 1).Format(time.RFC1123)

		methodTests := utils.V4TestCase{
			"GET": {
				"NOT MODIFIED when NO CHANGES made": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfModifiedSince: {tomorrow}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"OK when VALID request": {
					ClientSession: TOSession,
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateTypeSort()),
				},
				"OK when VALID NAME parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"name": {"ORG"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1),
						validateTypeFields(map[string]interface{}{"Name": "ORG"})),
				},
			},
			"POST": {
				"BAD REQUEST when NOT in server table": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"description": "Host header regular expression-Test",
						"name":        "TEST_1",
						"useInTable":  "regex",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"OK when VALID request in SERVER table": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"description": "Host header regular expression-Test",
						"name":        "TEST_4",
						"useInTable":  "server",
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateTypeUpdateCreateFields("TEST_4", map[string]interface{}{"Name": "TEST_4"})),
				},
			},
			"PUT": {
				"BAD REQUEST when NOT in server table": {
					EndpointId:    GetTypeID(t, "ACTIVE_DIRECTORY"),
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"description": "Active Directory User",
						"name":        "TEST_3",
						"useInTable":  "cachegroup",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"OK when VALID request in SERVER table": {
					EndpointId:    GetTypeID(t, "RIAK"),
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"description": "riak type",
						"name":        "TEST_5",
						"useInTable":  "server",
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateTypeUpdateCreateFields("TEST_5", map[string]interface{}{"Name": "TEST_5"})),
				},
			},
			"DELETE": {
				"OK when VALID request": {
					EndpointId:    GetTypeID(t, "INFLUXDB"),
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
			},
			"GET AFTER CHANGES": {
				"OK when CHANGES made": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfModifiedSince: {currentTimeRFC}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					typ := tc.Type{}

					if testCase.RequestBody != nil {
						dat, err := json.Marshal(testCase.RequestBody)
						assert.NoError(t, err, "Error occurred when marshalling request body: %v", err)
						err = json.Unmarshal(dat, &typ)
						assert.NoError(t, err, "Error occurred when unmarshalling request body: %v", err)
					}

					switch method {
					case "GET", "GET AFTER CHANGES":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.GetTypes(testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.CreateType(typ, client.RequestOptions{})
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					case "PUT":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.UpdateType(testCase.EndpointId(), typ, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.DeleteType(testCase.EndpointId(), testCase.RequestOpts)
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

func validateTypeSort() utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, alerts tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Type response to not be nil.")
		var typeNames []string
		typeResp := resp.([]tc.Type)
		for _, typ := range typeResp {
			typeNames = append(typeNames, typ.Name)
		}
		assert.Equal(t, true, sort.StringsAreSorted(typeNames), "List is not sorted by their names: %v", typeNames)
	}
}

func validateTypeFields(expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Type response to not be nil.")
		typeResp := resp.([]tc.Type)
		for field, expected := range expectedResp {
			for _, typ := range typeResp {
				switch field {
				case "Name":
					assert.Equal(t, expected, typ.Name, "Expected Name to be %v, but got %s", expected, typ.Name)
				case "UseInTable":
					assert.Equal(t, expected, typ.UseInTable, "Expected UseInTable to be %v, but got %s", expected, typ.UseInTable)
				default:
					t.Errorf("Expected field: %v, does not exist in response", field)
				}
			}
		}
	}
}

func validateTypeUpdateCreateFields(name string, expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("name", name)
		typ, _, err := TOSession.GetTypes(opts)
		assert.RequireNoError(t, err, "Error getting Types: %v - alerts: %+v", err, typ.Alerts)
		assert.RequireEqual(t, 1, len(typ.Response), "Expected one Type returned, Got: %d", len(typ.Response))
		validateTypeFields(expectedResp)(t, toclientlib.ReqInf{}, typ.Response, tc.Alerts{}, nil)
	}
}

func GetTypeID(t *testing.T, typeName string) func() int {
	return func() int {
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("name", typeName)
		resp, _, err := TOSession.GetTypes(opts)

		assert.RequireNoError(t, err, "Get Types Request failed with error: %v", err)
		assert.RequireEqual(t, 1, len(resp.Response), "Expected response object length 1, but got %d", len(resp.Response))

		return resp.Response[0].ID
	}
}

func CreateTestTypes(t *testing.T) {
	db, err := OpenConnection()
	assert.RequireNoError(t, err, "cannot open db")

	defer func() {
		err := db.Close()
		assert.NoError(t, err, "unable to close connection to db, error: %v", err)
	}()
	dbQueryTemplate := "INSERT INTO type (name, description, use_in_table) VALUES ('%s', '%s', '%s');"

	//opts := client.NewRequestOptions()
	for _, typ := range testData.Types {
		if typ.UseInTable != "server" {
			err = execSQL(db, fmt.Sprintf(dbQueryTemplate, typ.Name, typ.Description, typ.UseInTable))
			assert.RequireNoError(t, err, "could not create Type using database operations: %v", err)
		} else {
			alerts, _, err := TOSession.CreateType(typ, client.RequestOptions{})
			assert.RequireNoError(t, err, "could not create Type: %v - alerts: %+v", err, alerts.Alerts)
		}
	}
}

func DeleteTestTypes(t *testing.T) {
	db, err := OpenConnection()
	assert.RequireNoError(t, err, "cannot open db")

	defer func() {
		err := db.Close()
		assert.NoError(t, err, "unable to close connection to db, error: %v", err)
	}()
	dbDeleteTemplate := "DELETE FROM type WHERE name='%s';"

	types, _, err := TOSession.GetTypes(client.RequestOptions{})
	assert.NoError(t, err, "Cannot get Types: %v - alerts: %+v", err, types.Alerts)

	for _, typ := range types.Response {
		if typ.Name == "CHECK_EXTENSION_BOOL" || typ.Name == "CHECK_EXTENSION_NUM" || typ.Name == "CHECK_EXTENSION_OPEN_SLOT" {
			continue
		}

		if typ.UseInTable != "server" {
			err := execSQL(db, fmt.Sprintf(dbDeleteTemplate, typ.Name))
			assert.RequireNoError(t, err, "cannot delete Type using database operations: %v", err)
		} else {
			delResp, _, err := TOSession.DeleteType(typ.ID, client.RequestOptions{})
			assert.RequireNoError(t, err, "cannot delete Type using the API: %v - alerts: %+v", err, delResp.Alerts)
		}

		// Retrieve the Type by name so we can get the id for the Update
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("name", typ.Name)
		types, _, err := TOSession.GetTypes(opts)
		assert.NoError(t, err, "error fetching Types filtered by presumably deleted name: %v - alerts: %+v", err, types.Alerts)
		assert.Equal(t, 0, len(types.Response), "expected Type '%s' to be deleted", typ.Name)
	}
}
