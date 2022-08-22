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
	"strings"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/assert"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
	client "github.com/apache/trafficcontrol/traffic_ops/v4-client"
)

func TestJobs(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, ServiceCategories, DeliveryServices, Jobs}, func() {

		startTime := time.Now().Add(time.Minute).UTC()

		methodTests := utils.V4TestCase{
			"GET": {
				"OK when VALID request": {
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1)),
				},
				"OK when VALID ASSETURL parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"assetUrl": {""}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1),
						validateInvalidationJobsFields(map[string]interface{}{"AssetURL": ""})),
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
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"id": {strconv.Itoa(GetJobID(t, "")())}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1),
						validateInvalidationJobsFields(map[string]interface{}{"ID": GetJobID(t, "")()})),
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
						validateInvalidationJobsFields(map[string]interface{}{"MaxRevalDurationDays": ""})),
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
						validateInvalidationJobsFields(map[string]interface{}{"DeliveryServcie": "ds-forked-topology"})),
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
				"OK when STARTTIME is UNIX FORMAT AND a FUTURE DATE": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"deliveryService":  "ds1",
						"regex":            "/.*",
						"startTime":        startTime.AddDate(0, 0, 1).Unix(),
						"ttlHours":         36,
						"invalidationType": "REFRESH",
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"OK when STARTTIME is NON-STANDARD FORMAT AND a FUTURE DATE": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"deliveryService":  "ds1",
						"regex":            "/.*",
						"startTime":        startTime.AddDate(0, 0, 1).Format("2020-03-11 14:12:20-06"),
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
						"ttlHours":         9999999999,
						"invalidationType": "REFRESH",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when ALREADY EXISTS": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"deliveryService":  "ds1",
						"regex":            "/.*",
						"startTime":        startTime,
						"ttlHours":         72,
						"invalidationType": "REFRESH",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
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
				"BAD REQUEST when STARTTIME is NON-STANDARD FORMAT AND a PAST DATE": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"deliveryService":  "ds1",
						"regex":            "/.*",
						"startTime":        startTime.AddDate(0, 0, -1).Format("2020-03-11 14:12:20-06"),
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
						"startTime":        startTime.AddDate(0, 0, -1).Format(time.RFC1123),
						"ttlHours":         36,
						"invalidationType": "REFRESH",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when STARTTIME is UNIX FORMAT AND is a PAST DATE": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"deliveryService":  "ds1",
						"regex":            "/.*",
						"startTime":        startTime.AddDate(0, 0, -1).Unix,
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
			},
			"PUT": {
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
				"OK when STARTTIME is UNIX FORMAT AND a FUTURE DATE": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"deliveryService":  "ds1",
						"regex":            "/.*",
						"startTime":        startTime.AddDate(0, 0, 1).Unix(),
						"ttlHours":         36,
						"invalidationType": "REFRESH",
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"OK when STARTTIME is NON-STANDARD FORMAT AND a FUTURE DATE": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"deliveryService":  "ds1",
						"regex":            "/.*",
						"startTime":        startTime.AddDate(0, 0, 1).Format("2020-03-11 14:12:20-06"),
						"ttlHours":         36,
						"invalidationType": "REFRESH",
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"BAD REQUEST when INVALID ID": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"assetUrl":         "",
						"createdBy":        "admin",
						"deliveryService":  "ds1",
						"id":               111111111,
						"regex":            "/old",
						"startTime":        startTime.AddDate(0, 0, 3),
						"ttlHours":         2160,
						"invalidationType": "REFRESH",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when STARTTIME NOT within 2 DAYS FROM NOW": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"assetUrl":         "",
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
						"assetUrl":         "",
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
						"assetUrl":         "",
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
						"assetUrl":         "",
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
						"assetUrl":         "http://google.com",
						"createdBy":        "admin",
						"deliveryService":  "",
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
						"assetUrl":         "",
						"createdBy":        "admin",
						"deliveryService":  "",
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
						"assetUrl":         "",
						"deliveryService":  "",
						"regex":            "/old",
						"startTime":        startTime,
						"ttlHours":         2160,
						"invalidationType": "REFRESH",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when DIFFERENT CREATEDBY": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"assetUrl":         "",
						"createdBy":        "operator",
						"deliveryService":  "",
						"regex":            "/old",
						"startTime":        startTime,
						"ttlHours":         2160,
						"invalidationType": "REFRESH",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when BLANK INVALIDATION TYPE": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"assetUrl":         "",
						"createdBy":        "operator",
						"deliveryService":  "",
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
						"deliveryService":  "ds1",
						"regex":            "/.*",
						"startTime":        startTime.AddDate(0, 0, -1),
						"ttlHours":         36,
						"invalidationType": "REFRESH",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when STARTTIME is NON-STANDARD FORMAT AND a PAST DATE": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"deliveryService":  "ds1",
						"regex":            "/.*",
						"startTime":        startTime.AddDate(0, 0, -1).Format("2020-03-11 14:12:20-06"),
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
						"startTime":        startTime.AddDate(0, 0, -1).Format(time.RFC1123),
						"ttlHours":         36,
						"invalidationType": "REFRESH",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when STARTTIME is UNIX FORMAT AND is a PAST DATE": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"deliveryService":  "ds1",
						"regex":            "/.*",
						"startTime":        startTime.AddDate(0, 0, -1).Unix,
						"ttlHours":         36,
						"invalidationType": "REFRESH",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
			},
			"DELETE": {
				"NOT FOUND when JOB DOESNT EXIST": {
					EndpointId:    GetJobID(t, ""),
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					job := tc.InvalidationJobCreateV4{}
					jobUpdate := tc.InvalidationJob{}

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
					case "GET", "GET AFTER CHANGES":
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
							alerts, reqInf, err := testCase.ClientSession.UpdateInvalidationJob(job, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.DeleteInvalidationJob(uint64(testCase.EndpointId()), testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					}
				}
			})
		}

		JobCollisionWarningTest(t)
		CreateRefetchJobParameterFail(t)
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
					//			if !strings.HasSuffix(toJob.AssetURL, testJob.Regex) {
					//				continue
					//			}
					assert.Equal(t, expected, job.AssetURL, "Expected AssetURL to be %v, but got %s", expected, job.AssetURL)
				case "CreatedBy":
					assert.Equal(t, expected, job.CreatedBy, "Expected CreatedBy to be %v, but got %s", expected, job.CreatedBy)
				case "DeliveryService":
					assert.Equal(t, expected, job.DeliveryService, "Expected DeliveryService to be %v, but got %s", expected, job.DeliveryService)
				case "ID":
					assert.Equal(t, expected, job.ID, "Expected ID to be %v, but got %s", expected, job.ID)
				case "InvalidationType":
					assert.Equal(t, expected, job.InvalidationType, "Expected InvalidationType to be %v, but got %s", expected, job.InvalidationType)
				case "StartTime":
					/*
								toJobTime := toJob.StartTime.Round(time.Minute)
						testJobTime := testJob.StartTime.Round(time.Minute)
						if !toJobTime.Equal(testJobTime) {
							t.Errorf("test job ds %v regex %s start time expected '%+v' actual '%+v'", testJob.DeliveryService, testJob.Regex, testJobTime, toJobTime)
							continue
						}
							if time.Since(j.StartTime) > time.Duration(maxRevalDurationDays)*24*time.Hour {
							t.Errorf("GET /jobs by maxRevalDurationDays returned job that is older than %d days: {%s, %s, %v}", maxRevalDurationDays, j.DeliveryService, j.AssetURL, j.StartTime)
						}
					*/
					assert.Equal(t, expected, job.StartTime, "Expected StartTime to be %v, but got %s", expected, job.StartTime)
				case "TTLHours":
					assert.Equal(t, expected, job.TTLHours, "Expected TTLHours to be %v, but got %s", expected, job.TTLHours)
				default:
					t.Errorf("Expected field: %v, does not exist in response", field)
				}
			}
		}
	}
}

func GetJobID(t *testing.T, name string) func() int {
	return func() int {
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("name", name)
		jobs, _, err := TOSession.GetInvalidationJobs(opts)
		assert.RequireNoError(t, err, "Get Jobs Request failed with error:", err)
		assert.RequireEqual(t, 1, len(jobs.Response), "Expected response object length 1, but got %d", len(jobs.Response))
		return int(jobs.Response[0].ID)
	}
}

func JobCollisionWarningTest(t *testing.T) {
	if len(testData.DeliveryServices) < 1 {
		t.Fatal("Need at least one Delivery Service to test Invalidation Job collisions")
	}
	if testData.DeliveryServices[0].XMLID == nil {
		t.Fatal("Found a Delivery Service in the testing data with null or undefined XMLID")
	}
	xmlID := *testData.DeliveryServices[0].XMLID

	firstJob := tc.InvalidationJobCreateV4{
		DeliveryService:  xmlID,
		Regex:            `/\.*([A-Z]0?)`,
		TTLHours:         16,
		StartTime:        time.Now().Add(time.Hour),
		InvalidationType: tc.REFRESH,
	}

	resp, _, err := TOSession.CreateInvalidationJob(firstJob, client.RequestOptions{})
	if err != nil {
		t.Fatalf("Unexpected error creating a content invalidation Job: %v - alerts: %+v", err, resp.Alerts)
	}

	newJob := tc.InvalidationJobCreateV4{
		DeliveryService:  firstJob.DeliveryService,
		Regex:            firstJob.Regex,
		TTLHours:         firstJob.TTLHours,
		StartTime:        firstJob.StartTime.Add(time.Hour),
		InvalidationType: tc.REFRESH,
	}

	alerts, _, err := TOSession.CreateInvalidationJob(newJob, client.RequestOptions{})
	if err != nil {
		t.Fatalf("expected invalidation job create to succeed: %v - %+v", err, alerts.Alerts)
	}

	if len(alerts.Alerts) < 2 {
		t.Fatalf("expected at least 2 alerts on creation, got %v", len(alerts.Alerts))
	}

	found := false
	for _, alert := range alerts.Alerts {
		if alert.Level == tc.WarnLevel.String() && strings.Contains(alert.Text, firstJob.Regex) {
			found = true
		}
	}
	if !found {
		t.Error("Expected a warning-level error about the regular expression, but couldn't find one")
	}

	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("deliveryService", xmlID)
	jobs, _, err := TOSession.GetInvalidationJobs(opts)
	if err != nil {
		t.Fatalf("unable to get invalidation jobs: %v - alerts: %+v", err, jobs.Alerts)
	}

	var realJob *tc.InvalidationJobV4
	for i, job := range jobs.Response {
		diff := newJob.StartTime.Sub(job.StartTime)
		if job.DeliveryService == xmlID && job.CreatedBy == "admin" && diff < time.Second {
			realJob = &jobs.Response[i]
			break
		}
	}

	if realJob == nil || realJob.ID == 0 {
		t.Fatal("could not find new job")
	}

	time := firstJob.StartTime.Add(time.Hour * 2)
	realJob.StartTime = time
	alerts, _, err = TOSession.UpdateInvalidationJob(*realJob, client.RequestOptions{})
	if err != nil {
		t.Fatalf("expected invalidation job update to succeed: %v - alerts: %+v", err, alerts.Alerts)
	}

	if len(alerts.Alerts) < 2 {
		t.Fatalf("expected at least 2 alerts on update, got %v", len(alerts.Alerts))
	}

	found = false
	for _, alert := range alerts.Alerts {
		if alert.Level == tc.WarnLevel.String() && strings.Contains(alert.Text, firstJob.Regex) {
			found = true
		}
	}
	if !found {
		t.Error("Expected a warning-level error about the regular expression, but couldn't find one")
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
	// Ensure clean slate for parameters
	clearRefetchEnabledParameter(t)

	// Attempt to create Refetch job w/o refetch_enabled
	job := testData.InvalidationJobsRefetch[0]
	createJob := tc.InvalidationJobCreateV4{
		DeliveryService:  job.DeliveryService,
		Regex:            job.Regex,
		TTLHours:         job.TTLHours,
		StartTime:        time.Now().Add(time.Hour).UTC(),
		InvalidationType: job.InvalidationType,
	}

	_, _, err := TOSession.CreateInvalidationJob(createJob, client.RequestOptions{})
	assert.Error(t, err, "Expected error preventing the creation of the Refetch Invalidation Job.")
}

func clearRefetchEnabledParameter(t *testing.T) {
	// Ensure Parameter is not set
	paramsResp, _, err := TOSession.GetParameters(client.RequestOptions{})
	assert.RequireNoError(t, err, "Error retrieving parameters. err: %v \n alerts: %v", err, paramsResp.Alerts)

	for _, param := range paramsResp.Response {
		if param.Name == string(tc.RefetchEnabled) {
			TOSession.DeleteParameter(param.ID, client.RequestOptions{})
		}
	}
}
