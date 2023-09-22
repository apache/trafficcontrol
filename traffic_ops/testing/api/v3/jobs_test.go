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
	"encoding/json"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
)

func TestJobs(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, DeliveryServices, Jobs}, func() {

		currentTime := time.Now()
		startTime := currentTime.UTC().Add(time.Minute)

		methodTests := utils.V3TestCase{
			"GET": {
				"OK when VALID request": {
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1)),
				},
				"OK when VALID DELIVERYSERVICE parameter": {
					ClientSession: TOSession,
					RequestParams: url.Values{"deliveryService": {"ds2"}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1),
						validateInvalidationJobsFields(map[string]interface{}{"DeliveryService": "ds2"})),
				},
			},
			"POST": {
				"BAD REQUEST when TTLHours value GREATER than MAXREVALDURATIONDAYS": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"deliveryService": "ds1",
						"regex":           "/.*",
						"startTime":       startTime,
						"ttl":             9999,
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"WARNING ALERT when JOB COLLISION": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"deliveryService": "ds1",
						"regex":           "/foo",
						"startTime":       startTime.Add(time.Hour),
						"ttl":             2160,
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.HasAlertLevel(tc.WarnLevel.String())),
				},
			},
			"PUT": {
				"WARNING ALERT when JOB COLLISION": {
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"id":              GetJobID(t, "ds1", "admin")(),
						"assetUrl":        "http://origin.example.net/foo",
						"createdBy":       "admin",
						"deliveryService": "ds1",
						"regex":           "/foo",
						"keyword":         "PURGE",
						"parameters":      "TTL:2h",
						"startTime":       startTime.Add(time.Hour),
						"ttl":             2160,
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.HasAlertLevel(tc.WarnLevel.String())),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					job := tc.InvalidationJobInput{}
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
					case "GET":
						t.Run(name, func(t *testing.T) {
							var ds interface{}
							if val, ok := testCase.RequestParams["deliveryService"]; ok {
								ds = val[0]
								resp, reqInf, err := testCase.ClientSession.GetInvalidationJobsWithHdr(&ds, nil, testCase.RequestHeaders)
								for _, check := range testCase.Expectations {
									check(t, reqInf, resp, tc.Alerts{}, err)
								}
							} else {
								resp, reqInf, err := testCase.ClientSession.GetInvalidationJobsWithHdr(nil, nil, testCase.RequestHeaders)
								for _, check := range testCase.Expectations {
									check(t, reqInf, resp, tc.Alerts{}, err)
								}
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.CreateInvalidationJob(job)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					case "PUT":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.UpdateInvalidationJob(jobUpdate)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.DeleteInvalidationJob(uint64(testCase.EndpointID()))
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

func validateInvalidationJobsFields(expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Invalidation Jobs response to not be nil.")
		jobResp := resp.([]tc.InvalidationJob)
		for field, expected := range expectedResp {
			for _, job := range jobResp {
				switch field {
				case "AssetURL":
					assert.RequireNotNil(t, job.AssetURL, "Expected AssetURL to not be nil.")
					assert.Equal(t, expected, *job.AssetURL, "Expected AssetURL to be %v, but got %s", expected, *job.AssetURL)
				case "CreatedBy":
					assert.RequireNotNil(t, job.CreatedBy, "Expected CreatedBy to not be nil.")
					assert.Equal(t, expected, *job.CreatedBy, "Expected CreatedBy to be %v, but got %s", expected, *job.CreatedBy)
				case "DeliveryService":
					assert.RequireNotNil(t, job.DeliveryService, "Expected DeliveryService to not be nil.")
					assert.Equal(t, expected, *job.DeliveryService, "Expected DeliveryService to be %v, but got %s", expected, *job.DeliveryService)
				case "ID":
					assert.RequireNotNil(t, job.ID, "Expected ID to not be nil.")
					assert.Equal(t, uint64(expected.(int)), *job.ID, "Expected ID to be %v, but got %s", expected, *job.ID)
				case "Keyword":
					assert.RequireNotNil(t, job.Keyword, "Expected Keyword to not be nil.")
					assert.Equal(t, expected, *job.Keyword, "Expected Keyword to be %v, but got %s", expected, *job.Keyword)
				case "Parameters":
					assert.RequireNotNil(t, job.Parameters, "Expected Parameters to not be nil.")
					assert.Equal(t, expected, *job.Parameters, "Expected Parameters to be %v, but got %s", expected, *job.Parameters)
				case "StartTime":
					assert.RequireNotNil(t, job.StartTime, "Expected StartTime to not be nil.")
					assert.Equal(t, true, job.StartTime.Round(time.Minute).Equal(expected.(time.Time).Round(time.Minute)), "Expected StartTime to be %v, but got %s", expected, job.StartTime)
				default:
					t.Errorf("Expected field: %v, does not exist in response", field)
				}
			}
		}
	}
}

func GetJobID(t *testing.T, ds interface{}, user interface{}) func() int {
	return func() int {
		t.Helper()
		jobs, _, err := TOSession.GetInvalidationJobsWithHdr(&ds, &user, nil)
		assert.RequireNoError(t, err, "Get Jobs Request failed with error:", err)
		assert.RequireGreaterOrEqual(t, len(jobs), 1, "Expected at least 1 response object, but got %d", len(jobs))
		assert.RequireNotNil(t, jobs[0].ID, "Expected Job ID to not be nil.")
		return int(*jobs[0].ID)
	}
}

func CreateTestJobs(t *testing.T) {
	for _, job := range testData.InvalidationJobs {
		job.StartTime = &tc.Time{
			Time:  time.Now().Add(time.Minute).UTC(),
			Valid: true,
		}
		resp, _, err := TOSession.CreateInvalidationJob(job)
		assert.RequireNoError(t, err, "Could not create job: %v - alerts: %+v", err, resp.Alerts)
	}
}

func DeleteTestJobs(t *testing.T) {
	jobs, _, err := TOSession.GetInvalidationJobsWithHdr(nil, nil, nil)
	assert.NoError(t, err, "Cannot get Jobs: %v - alerts: %+v", err)

	for _, job := range jobs {
		assert.RequireNotNil(t, job.ID, "Expected JOB ID to not be nil.")
		alerts, _, err := TOSession.DeleteInvalidationJob(*job.ID)
		assert.NoError(t, err, "Unexpected error deleting Job with ID: (#%d): %v - alerts: %+v", *job.ID, err, alerts.Alerts)
	}
	// Retrieve the Jobs to see if they got deleted
	getJobs, _, err := TOSession.GetInvalidationJobsWithHdr(nil, nil, nil)
	assert.NoError(t, err, "Error getting Jobs after deletion: %v", err)
	assert.Equal(t, 0, len(getJobs), "Expected Jobs to be deleted, but %d jobs were found in Traffic Ops", len(getJobs))
}
