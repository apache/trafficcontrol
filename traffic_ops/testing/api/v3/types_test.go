package v3

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
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
	tc "github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
	client "github.com/apache/trafficcontrol/v8/traffic_ops/v4-client"
)

func TestTypes(t *testing.T) {
	WithObjs(t, []TCObj{Parameters, Types}, func() {

		currentTime := time.Now().UTC().Add(-15 * time.Second)
		currentTimeRFC := currentTime.Format(time.RFC1123)
		tomorrow := currentTime.AddDate(0, 0, 1).Format(time.RFC1123)

		methodTests := utils.V3TestCaseT[tc.Type]{
			"GET": {
				"NOT MODIFIED when NO CHANGES made": {
					ClientSession:  TOSession,
					RequestHeaders: http.Header{rfc.IfModifiedSince: {tomorrow}},
					Expectations:   utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"OK when VALID request": {
					ClientSession: TOSession,
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateTypeSort()),
				},
				"OK when VALID NAME parameter": {
					ClientSession: TOSession,
					RequestParams: url.Values{"name": {"ORG"}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1),
						validateTypeFields(map[string]interface{}{"Name": "ORG"})),
				},
				"OK when CHANGES made": {
					ClientSession:  TOSession,
					RequestHeaders: http.Header{rfc.IfModifiedSince: {currentTimeRFC}},
					Expectations:   utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
			},
			"POST": {
				"BAD REQUEST when useInTable NOT server": {
					ClientSession: TOSession,
					RequestBody: tc.Type{
						Description: "Host header regular expression-Test",
						Name:        "TEST_1",
						UseInTable:  "regex",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"OK when VALID request when useInTable=server": {
					ClientSession: TOSession,
					RequestBody: tc.Type{
						Description: "Host header regular expression-Test",
						Name:        "TEST_4",
						UseInTable:  "server",
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateTypeUpdateCreateFields("TEST_4", map[string]interface{}{"Name": "TEST_4"})),
				},
			},
			"PUT": {
				"BAD REQUEST when useInTable NOT server": {
					EndpointID:    GetTypeID(t, "ACTIVE_DIRECTORY"),
					ClientSession: TOSession,
					RequestBody: tc.Type{
						Description: "Active Directory User",
						Name:        "TEST_3",
						UseInTable:  "cachegroup",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"OK when VALID request when useInTable=server": {
					EndpointID:    GetTypeID(t, "RIAK"),
					ClientSession: TOSession,
					RequestBody: tc.Type{
						Description: "riak type",
						Name:        "TEST_5",
						UseInTable:  "server",
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateTypeUpdateCreateFields("TEST_5", map[string]interface{}{"Name": "TEST_5"})),
				},
			},
			"DELETE": {
				"OK when VALID request": {
					EndpointID:    GetTypeID(t, "INFLUXDB"),
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {

					params := make(map[string]string)
					if testCase.RequestParams != nil {
						for k, v := range testCase.RequestParams {
							params[k] = v[0]
						}
					}

					switch method {
					case "GET":
						t.Run(name, func(t *testing.T) {
							var resp []tc.Type
							var reqInf toclientlib.ReqInf
							var err error
							if len(params) != 0 {
								resp, reqInf, err = testCase.ClientSession.GetTypeByNameWithHdr(params["name"], testCase.RequestHeaders)
							} else {
								resp, reqInf, err = testCase.ClientSession.GetTypesWithHdr(testCase.RequestHeaders)
							}
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp, tc.Alerts{}, err)
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.CreateType(testCase.RequestBody)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					case "PUT":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.UpdateTypeByIDWithHdr(testCase.EndpointID(), testCase.RequestBody, testCase.RequestHeaders)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.DeleteTypeByID(testCase.EndpointID())
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
		typ, _, err := TOSession.GetTypeByNameWithHdr(name, nil)
		assert.RequireNoError(t, err, "Error getting Types: %v", err)
		assert.RequireEqual(t, 1, len(typ), "Expected one Type returned, Got: %d", len(typ))
		validateTypeFields(expectedResp)(t, toclientlib.ReqInf{}, typ, tc.Alerts{}, nil)
	}
}

func GetTypeID(t *testing.T, typeName string) func() int {
	return func() int {
		resp, _, err := TOSession.GetTypeByNameWithHdr(typeName, nil)
		assert.RequireNoError(t, err, "Get Types Request failed with error: %v", err)
		assert.RequireEqual(t, 1, len(resp), "Expected response object length 1, but got %d", len(resp))

		return resp[0].ID
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

	for _, typ := range testData.Types {
		if typ.UseInTable != "server" {
			err = execSQL(db, fmt.Sprintf(dbQueryTemplate, typ.Name, typ.Description, typ.UseInTable))
			assert.RequireNoError(t, err, "could not create Type using database operations: %v", err)
		} else {
			alerts, _, err := TOSession.CreateType(typ)
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

	types, _, err := TOSession.GetTypesWithHdr(nil)
	assert.NoError(t, err, "Cannot get Types: %v: %+v", err)

	for _, typ := range types {
		if typ.Name == "CHECK_EXTENSION_BOOL" || typ.Name == "CHECK_EXTENSION_NUM" || typ.Name == "CHECK_EXTENSION_OPEN_SLOT" {
			continue
		}

		if typ.UseInTable != "server" {
			err := execSQL(db, fmt.Sprintf(dbDeleteTemplate, typ.Name))
			assert.RequireNoError(t, err, "cannot delete Type using database operations: %v", err)
		} else {
			delResp, _, err := TOSession.DeleteTypeByID(typ.ID)
			assert.RequireNoError(t, err, "cannot delete Type using the API: %v - alerts: %+v", err, delResp.Alerts)
		}

		// Retrieve the Type by name to see if it was deleted.
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("name", typ.Name)
		types, _, err := TOSession.GetTypesWithHdr(nil, typ.Name)
		assert.NoError(t, err, "error fetching Types filtered by presumably deleted name: %v", err)
		assert.Equal(t, 0, len(types), "expected Type '%s' to be deleted", typ.Name)
	}
}
