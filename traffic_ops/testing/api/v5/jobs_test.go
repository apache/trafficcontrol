package v5

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

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
	client "github.com/apache/trafficcontrol/v8/traffic_ops/v5-client"
)

func TestJobs(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, ServiceCategories, ServerCapabilities, ServerServerCapabilities, DeliveryServices, Jobs}, func() {

		currentTime := time.Now()
		pastTime := currentTime.AddDate(0, 0, -1)
		futureTime := currentTime.AddDate(0, 0, 1)
		startTime := currentTime.UTC().Add(time.Minute)
		pastTimeRFC := pastTime.Format(time.RFC3339)
		futureTimeRFC := futureTime.Format(time.RFC3339)

		methodTests := utils.V5TestCase{
			"GET": {
				"OK when VALID request": {
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1)),
				},
				"OK when VALID ASSETURL parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"assetUrl": {"http://origin.example.net/older"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1),
						validateInvalidationJobsFields(map[string]interface{}{"AssetURL": "http://origin.example.net/older"})),
				},
				"OK when VALID CREATEDBY parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"createdBy": {"admin"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateInvalidationJobsFields(map[string]interface{}{"CreatedBy": "admin"})),
				},
				"OK when VALID DELIVERYSERVICE parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"deliveryService": {"ds2"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1),
						validateInvalidationJobsFields(map[string]interface{}{"DeliveryService": "ds2"})),
				},
				"OK when VALID DSID parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"dsId": {strconv.Itoa(GetDeliveryServiceId(t, "ds2")())}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1),
						validateInvalidationJobsFields(map[string]interface{}{"DeliveryService": "ds2"})),
				},
				"OK when VALID ID parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"id": {strconv.Itoa(GetJobID(t, "http://origin.example.net/oldest")())}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1),
						validateInvalidationJobsFields(map[string]interface{}{"ID": GetJobID(t, "http://origin.example.net/oldest")()})),
				},
				"OK when VALID INVALIDATIONTYPE parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"invalidationType": {"REFRESH"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateInvalidationJobsFields(map[string]interface{}{"InvalidationType": "REFRESH"})),
				},
				"OK when VALID MAXREVALDURATIONDAYS parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"maxRevalDurationDays": {""}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateMaxRevalDurationDays()),
				},
				"OK when VALID USERID parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"userId": {strconv.Itoa(GetUserID(t, "admin")())}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateInvalidationJobsFields(map[string]interface{}{"CreatedBy": "admin"})),
				},
				"OK when VALID CDN parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"cdn": {"cdn2"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateInvalidationJobsFields(map[string]interface{}{"DeliveryService": "ds-forked-topology"})),
				},
				"EMPTY RESPONSE when INVALID ASSETURL parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"assetUrl": {"doesntexist"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
				"EMPTY RESPONSE when INVALID CREATEDBY parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"createdBy": {"doesntexist"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
				"EMPTY RESPONSE when INVALID DELIVERYSERVICE parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"deliveryService": {"doesntexist"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
				"EMPTY RESPONSE when INVALID DSID parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"dsId": {"1111111111"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
				"EMPTY RESPONSE when INVALID ID parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"id": {"11111111111"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
				"EMPTY RESPONSE when INVALID INVALIDATIONTYPE parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"invalidationType": {"DOESNT EXIST"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
				"EMPTY RESPONSE when INVALID USERID parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"userId": {"1111111111"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
			},
			"POST": {
				"OK when STARTTIME is a FUTURE DATE": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"deliveryService":  "ds1",
						"regex":            "/.*",
						"startTime":        startTime.AddDate(0, 0, 1),
						"ttlHours":         36,
						"invalidationType": "REFRESH",
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"OK when STARTTIME is RFC FORMAT AND a FUTURE DATE": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"deliveryService":  "ds1",
						"regex":            "/.*",
						"startTime":        futureTimeRFC,
						"ttlHours":         36,
						"invalidationType": "REFRESH",
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"BAD REQUEST when TTLHours value GREATER than MAXREVALDURATIONDAYS": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"deliveryService":  "ds1",
						"regex":            "/.*",
						"startTime":        startTime,
						"ttlHours":         9999,
						"invalidationType": "REFRESH",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"OK when ALREADY EXISTS": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"deliveryService":  "ds1",
						"regex":            "/.*",
						"startTime":        startTime,
						"ttlHours":         72,
						"invalidationType": "REFRESH",
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.HasAlertLevel(tc.WarnLevel.String())),
				},
				"NOT FOUND when DELIVERYSERVICE DOESNT EXIST": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"deliveryService":  "doesntExist",
						"regex":            "/.*",
						"startTime":        startTime,
						"ttlHours":         36,
						"invalidationType": "REFRESH",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
				"BAD REQUEST when MISSING DELIVERYSERVICE": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"regex":            "/.*",
						"startTime":        startTime,
						"ttlHours":         36,
						"invalidationType": "REFRESH",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when BLANK DELIVERYSERVICE": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"deliveryService":  "",
						"regex":            "/.*",
						"startTime":        startTime,
						"ttlHours":         36,
						"invalidationType": "REFRESH",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when STARTTIME is a PAST DATE": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"deliveryService":  "ds1",
						"regex":            "/.*",
						"startTime":        startTime.AddDate(0, 0, -1),
						"ttlHours":         36,
						"invalidationType": "REFRESH",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when STARTTIME is RFC FORMAT AND is a PAST DATE": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"deliveryService":  "ds1",
						"regex":            "/.*",
						"startTime":        pastTimeRFC,
						"ttlHours":         36,
						"invalidationType": "REFRESH",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when MISSING STARTTIME": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"deliveryService":  "ds1",
						"regex":            "/.*",
						"ttlHours":         36,
						"invalidationType": "REFRESH",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when MISSING REGEX": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"deliveryService":  "ds1",
						"startTime":        startTime,
						"ttlHours":         36,
						"invalidationType": "REFRESH",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when EMPTY REGEX": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"deliveryService":  "ds1",
						"regex":            "",
						"startTime":        startTime,
						"ttlHours":         36,
						"invalidationType": "REFRESH",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when MISSING TTL": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"deliveryService":  "ds1",
						"regex":            "/.*",
						"startTime":        startTime,
						"invalidationType": "REFRESH",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when TTL is ZERO": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"deliveryService":  "ds1",
						"regex":            "/.*",
						"startTime":        startTime,
						"ttlHours":         0,
						"invalidationType": "REFRESH",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"WARNING ALERT when JOB COLLISION": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"deliveryService":  "ds1",
						"regex":            "/foo",
						"startTime":        startTime.Add(time.Hour),
						"ttlHours":         36,
						"invalidationType": "REFRESH",
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.HasAlertLevel(tc.WarnLevel.String())),
				},
			},
			"PUT": {
				"OK when STARTTIME is a FUTURE DATE": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"id":               GetJobID(t, "http://origin.example.net/.*")(),
						"assetUrl":         "http://origin.example.net/.*",
						"createdBy":        "admin",
						"deliveryService":  "ds1",
						"regex":            "/.*",
						"startTime":        startTime.AddDate(0, 0, 1),
						"ttlHours":         36,
						"invalidationType": "REFRESH",
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"NOT FOUND when INVALID ID": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"assetUrl":         "http://origin.example.net/.*",
						"createdBy":        "admin",
						"deliveryService":  "ds1",
						"id":               111111111,
						"regex":            "/old",
						"startTime":        startTime.AddDate(0, 0, 3),
						"ttlHours":         2160,
						"invalidationType": "REFRESH",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
				"BAD REQUEST when STARTTIME NOT within 2 DAYS FROM NOW": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"id":               GetJobID(t, "http://origin.example.net/.*")(),
						"assetUrl":         "http://origin.example.net/.*",
						"createdBy":        "admin",
						"deliveryService":  "ds1",
						"regex":            "/old",
						"startTime":        startTime.AddDate(0, 0, 3),
						"ttlHours":         2160,
						"invalidationType": "REFRESH",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"CONFLICT when DIFFERENT DELIVERY SERVICE": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"id":               GetJobID(t, "http://origin.example.net/.*")(),
						"assetUrl":         "http://origin.example.net/.*",
						"createdBy":        "admin",
						"deliveryService":  "ds3",
						"regex":            "/old",
						"startTime":        startTime,
						"ttlHours":         2160,
						"invalidationType": "REFRESH",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusConflict)),
				},
				"CONFLICT when INVALID DELIVERY SERVICE": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"id":               GetJobID(t, "http://origin.example.net/.*")(),
						"assetUrl":         "http://origin.example.net/.*",
						"createdBy":        "admin",
						"deliveryService":  "doesntexist",
						"regex":            "/old",
						"startTime":        startTime,
						"ttlHours":         2160,
						"invalidationType": "REFRESH",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusConflict)),
				},
				"BAD REQUEST when BLANK DELIVERY SERVICE": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"id":               GetJobID(t, "http://origin.example.net/.*")(),
						"assetUrl":         "http://origin.example.net/.*",
						"createdBy":        "admin",
						"deliveryService":  "",
						"regex":            "/old",
						"startTime":        startTime,
						"ttlHours":         2160,
						"invalidationType": "REFRESH",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when INVALID ASSETURL": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"id":               GetJobID(t, "http://origin.example.net/.*")(),
						"assetUrl":         "http://google.com",
						"createdBy":        "admin",
						"deliveryService":  "ds1",
						"regex":            "/old",
						"startTime":        startTime,
						"ttlHours":         2160,
						"invalidationType": "REFRESH",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when BLANK ASSETURL": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"id":               GetJobID(t, "http://origin.example.net/.*")(),
						"assetUrl":         "",
						"createdBy":        "admin",
						"deliveryService":  "ds1",
						"regex":            "/old",
						"startTime":        startTime,
						"ttlHours":         2160,
						"invalidationType": "REFRESH",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when BLANK CREATEDBY": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"id":               GetJobID(t, "http://origin.example.net/.*")(),
						"assetUrl":         "http://origin.example.net/.*",
						"deliveryService":  "ds1",
						"regex":            "/old",
						"startTime":        startTime,
						"ttlHours":         2160,
						"invalidationType": "REFRESH",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"CONFLICT when DIFFERENT CREATEDBY": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"id":               GetJobID(t, "http://origin.example.net/.*")(),
						"assetUrl":         "http://origin.example.net/.*",
						"createdBy":        "operator",
						"deliveryService":  "ds1",
						"regex":            "/old",
						"startTime":        startTime,
						"ttlHours":         2160,
						"invalidationType": "REFRESH",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusConflict)),
				},
				"BAD REQUEST when BLANK INVALIDATION TYPE": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"id":               GetJobID(t, "http://origin.example.net/.*")(),
						"assetUrl":         "http://origin.example.net/.*",
						"createdBy":        "operator",
						"deliveryService":  "ds1",
						"regex":            "/old",
						"startTime":        startTime,
						"ttlHours":         2160,
						"invalidationType": "",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when STARTTIME is a PAST DATE": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"id":               GetJobID(t, "http://origin.example.net/.*")(),
						"assetUrl":         "http://origin.example.net/.*",
						"createdBy":        "admin",
						"deliveryService":  "ds1",
						"regex":            "/.*",
						"startTime":        startTime.AddDate(0, 0, -1),
						"ttlHours":         36,
						"invalidationType": "REFRESH",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when STARTTIME is RFC FORMAT AND is a PAST DATE": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"id":               GetJobID(t, "http://origin.example.net/.*")(),
						"assetUrl":         "http://origin.example.net/.*",
						"createdBy":        "admin",
						"deliveryService":  "ds1",
						"regex":            "/.*",
						"startTime":        pastTimeRFC,
						"ttlHours":         36,
						"invalidationType": "REFRESH",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"WARNING ALERT when JOB COLLISION": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"id":               GetJobID(t, "http://origin.example.net/foo")(),
						"assetUrl":         "http://origin.example.net/foo",
						"createdBy":        "admin",
						"deliveryService":  "ds1",
						"regex":            "/foo",
						"startTime":        startTime.Add(time.Hour),
						"ttlHours":         36,
						"invalidationType": "REFETCH",
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.HasAlertLevel(tc.WarnLevel.String())),
				},
			},
			"DELETE": {
				"NOT FOUND when JOB DOESNT EXIST": {
					EndpointID:    func() int { return 1111111111 },
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					job := tc.InvalidationJobCreateV4{}
					jobUpdate := tc.InvalidationJobV4{}

					if testCase.RequestBody != nil {
						dat, err := json.Marshal(testCase.RequestBody)
						assert.NoError(t, err, "Error occurred when marshalling request body: %v", err)
						if method == "POST" {
							err = json.Unmarshal(dat, &job)
							assert.NoError(t, err, "Error occurred when unmarshalling request body: %v", err)
						} else if method == "PUT" {
							err = json.Unmarshal(dat, &jobUpdate)
							assert.NoError(t, err, "Error occurred when unmarshalling request body: %v", err)
						}
					}

					switch method {
					case "GET":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.GetInvalidationJobs(testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.CreateInvalidationJob(job, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					case "PUT":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.UpdateInvalidationJob(jobUpdate, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.DeleteInvalidationJob(uint64(testCase.EndpointID()), testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					}
				}
			})
		}
		t.Run("POST/BAD REQUEST when REFETCH PARAMETER NOT ENABLED", func(t *testing.T) { CreateRefetchJobParameterFail(t) })
	})
}

func validateInvalidationJobsFields(expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Invalidation Jobs response to not be nil.")
		jobResp := resp.([]tc.InvalidationJobV4)
		for field, expected := range expectedResp {
			for _, job := range jobResp {
				switch field {
				case "AssetURL":
					assert.Equal(t, expected, job.AssetURL, "Expected AssetURL to be %v, but got %s", expected, job.AssetURL)
				case "CreatedBy":
					assert.Equal(t, expected, job.CreatedBy, "Expected CreatedBy to be %v, but got %s", expected, job.CreatedBy)
				case "DeliveryService":
					assert.Equal(t, expected, job.DeliveryService, "Expected DeliveryService to be %v, but got %s", expected, job.DeliveryService)
				case "ID":
					assert.Equal(t, uint64(expected.(int)), job.ID, "Expected ID to be %v, but got %s", expected, job.ID)
				case "InvalidationType":
					assert.Equal(t, expected, job.InvalidationType, "Expected InvalidationType to be %v, but got %s", expected, job.InvalidationType)
				case "StartTime":
					assert.Equal(t, true, job.StartTime.Round(time.Minute).Equal(expected.(time.Time).Round(time.Minute)), "Expected StartTime to be %v, but got %s", expected, job.StartTime)
				case "TTLHours":
					assert.Equal(t, expected, job.TTLHours, "Expected TTLHours to be %v, but got %s", expected, job.TTLHours)
				default:
					t.Errorf("Expected field: %v, does not exist in response", field)
				}
			}
		}
	}
}

func validateMaxRevalDurationDays() utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Invalidation Jobs response to not be nil.")
		maxRevalDurationDays := 90
		jobResp := resp.([]tc.InvalidationJobV4)
		for _, job := range jobResp {
			if time.Since(job.StartTime) > time.Duration(maxRevalDurationDays)*24*time.Hour {
				t.Errorf("GET /jobs by maxRevalDurationDays returned job that is older than %d days: %v}", maxRevalDurationDays, time.Since(job.StartTime))
			}
		}
	}
}

func GetJobID(t *testing.T, assetUrl string) func() int {
	return func() int {
		t.Helper()
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("assetUrl", assetUrl)
		jobs, _, err := TOSession.GetInvalidationJobs(opts)
		assert.RequireNoError(t, err, "Get Jobs Request failed with error:", err)
		assert.RequireGreaterOrEqual(t, len(jobs.Response), 1, "Expected at least 1 response object, but got %d", len(jobs.Response))
		return int(jobs.Response[0].ID)
	}
}

func CreateTestJobs(t *testing.T) {
	for _, job := range testData.InvalidationJobs {
		job.StartTime = time.Now().Add(time.Minute).UTC()
		resp, _, err := TOSession.CreateInvalidationJob(job, client.RequestOptions{})
		assert.RequireNoError(t, err, "Could not create job: %v - alerts: %+v", err, resp.Alerts)
	}
}

func DeleteTestJobs(t *testing.T) {
	jobs, _, err := TOSession.GetInvalidationJobs(client.RequestOptions{})
	assert.NoError(t, err, "Cannot get Jobs: %v - alerts: %+v", err, jobs.Alerts)

	for _, job := range jobs.Response {
		alerts, _, err := TOSession.DeleteInvalidationJob(job.ID, client.RequestOptions{})
		assert.NoError(t, err, "Unexpected error deleting Job with ID: (#%d): %v - alerts: %+v", job.ID, err, alerts.Alerts)
		// Retrieve the Job to see if it got deleted
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("id", strconv.Itoa(int(job.ID)))
		getJobs, _, err := TOSession.GetInvalidationJobs(opts)
		assert.NoError(t, err, "Error getting Job with ID: '%d' after deletion: %v - alerts: %+v", job.ID, err, getJobs.Alerts)
		assert.Equal(t, 0, len(getJobs.Response), "Expected Job to be deleted, but it was found in Traffic Ops")
	}
}

func CreateRefetchJobParameterFail(t *testing.T) {
	// Delete the refetch parameter as a prerequisite
	clearRefetchEnabledParameter(t)
	createJob := tc.InvalidationJobCreateV4{
		DeliveryService:  "ds1",
		Regex:            "/.*",
		TTLHours:         72,
		StartTime:        time.Now().Add(time.Hour).UTC(),
		InvalidationType: "REFETCH",
	}
	_, reqInf, err := TOSession.CreateInvalidationJob(createJob, client.RequestOptions{})
	assert.Error(t, err, "Expected error preventing the creation of the Refetch Invalidation Job.")
	assert.Equal(t, http.StatusBadRequest, reqInf.StatusCode, "Expected Status Code: 400, Got: %d", reqInf.StatusCode)
}

func clearRefetchEnabledParameter(t *testing.T) {
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("name", string(tc.RefetchEnabled))
	paramsResp, _, err := TOSession.GetParameters(opts)
	assert.RequireNoError(t, err, "Error retrieving parameters. err: %v \n alerts: %v", err, paramsResp.Alerts)
	assert.RequireEqual(t, 1, len(paramsResp.Response), "Expected one parameter returned from response Got: %d", len(paramsResp.Response))
	assert.RequireEqual(t, string(tc.RefetchEnabled), paramsResp.Response[0].Name, "Expected the RefetchEnabled parameter Got: %s", paramsResp.Response[0].Name)
	alerts, _, err := TOSession.DeleteParameter(paramsResp.Response[0].ID, client.RequestOptions{})
	assert.RequireNoError(t, err, "Expected no error when deleting RefetchEnabled parameter: %v Alerts: %v", err, alerts.Alerts)
}
