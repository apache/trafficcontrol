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
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/assert"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
	client "github.com/apache/trafficcontrol/traffic_ops/v4-client"
)

// All prerequisite Snapshots are associated to this cdn
var cdn = "cdn1"

func TestSnapshot(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, ServiceCategories, DeliveryServices, DeliveryServiceServerAssignments}, func() {

		readOnlyUserSession := utils.CreateV4Session(t, Config.TrafficOps.URL, "readonlyuser", "pa$$word", Config.Default.Session.TimeoutInSecs)

		methodTests := utils.V4TestCase{
			"PUT": {
				"OK when VALID CDN parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"cdn": {"cdn1"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"OK when VALID CDNID parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"cdnID": {strconv.Itoa(GetCDNID(t, "cdn1")())}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"NOT FOUND when NON-EXISTENT CDN": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"cdn": {"cdn-invalid"}}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
				"NOT FOUND when NON-EXISTENT CDNID": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"cdnID": {"999999"}}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusNotFound)),
				},
				"FORBIDDEN when READ-ONLY user": {
					ClientSession: readOnlyUserSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"cdn": {"cdn1"}}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					switch method {
					case "PUT":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.SnapshotCRConfig(testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					}
				}
			})
		}
	})
}

func TestCDNNameSnapshot(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, ServiceCategories, DeliveryServices, Snapshot}, func() {

		methodTests := utils.V4TestCase{
			"GET": {
				"ANY-MAP DELIVERY SERVICE NOT IN CRCONFIG": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"cdnID": {strconv.Itoa(GetCDNID(t, cdnName)())}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateDeliveryServiceNotInResponse("anymap-ds")),
				},
				"TMPATH is NIL and TMHOST is CORRECT": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"cdnID": {strconv.Itoa(GetCDNID(t, cdnName)())}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateCRConfigFields(map[string]interface{}{"TMHost": "crconfig.tm.url.test.invalid"})),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					var cdn string
					switch method {
					case "GET":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.GetCRConfig(cdn, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					}
				}
			})
		}
	})
}

func TestCDNNameSnapshotNew(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, ServiceCategories, DeliveryServices, Snapshot}, func() {

		// Prerequiste: Delete Parameter and update snapshot

		methodTests := utils.V4TestCase{
			"GET": {
				"SNAPSHOT UPDATE CAPTURED CORRECTLY": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"cdn": {cdn}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateCRConfigFields(map[string]interface{}{"TMHost": ""})),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					var cdn string
					switch method {
					case "GET":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.GetCRConfigNew(cdn, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					}
				}
			})
		}
	})
}

func validateCRConfigFields(expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Snapshot response to not be nil.")
		crconfig := resp.(tc.CRConfig)
		assert.Equal(t, nil, crconfig.Stats.TMPath, "Expected no TMPath in APIv4, but it was: %s", *crconfig.Stats.TMPath)
		for field, expected := range expectedResp {
			switch field {
			case "TMHost":
				assert.RequireNotNil(t, crconfig.Stats.TMHost, "Expected Stats TM Host to not be nil.")
				assert.Equal(t, expected, *crconfig.Stats.TMHost, "Expected Stats TM Host to be %v, but got %s", expected, *crconfig.Stats.TMHost)
			default:
				t.Errorf("Expected field: %v, does not exist in response", field)
			}
		}
	}
}

func validateDeliveryServiceNotInResponse(deliveryService string) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected response to not be nil.")
		crconfig := resp.(tc.CRConfig)
		for ds := range crconfig.DeliveryServices {
			assert.NotEqual(t, ds, deliveryService, "Found unexpected delivery service: %s in CRConfig Delivery Services.", deliveryService)
		}
		for server := range crconfig.ContentServers {
			for ds := range crconfig.ContentServers[server].DeliveryServices {
				assert.NotEqual(t, ds, deliveryService, "Found unexpected delivery service: %s in CRConfig Content Servers Delivery Services.", deliveryService)
			}
		}
	}
}

func CreateSnapshot(t *testing.T) {
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("cdnID", strconv.Itoa(GetCDNID(t, cdn)()))
	resp, _, err := TOSession.SnapshotCRConfig(opts)
	assert.RequireNoError(t, err, "Could not create Snapshot: %v - alerts: %+v", err, resp.Alerts)
}

func DeleteSnapshot(t *testing.T) {
	return
}

func TestCDNNameConfigsMonitoring(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, ServiceCategories, DeliveryServices}, func() {

		methodTests := utils.V4TestCase{
			"GET": {
				"OK when VALID CDN parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"cdn": {"cdn1"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					var cdn string
					switch method {
					case "GET":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.GetCRConfigNew(cdn, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					}
				}
			})
		}
	})
}

func MonitoringConfig(t *testing.T) {
	if len(testData.CDNs) < 1 {
		t.Fatalf("no cdn test data")
	}
	const cdnName = "cdn1"
	const profileName = "EDGE1"
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("name", cdnName)
	cdns, _, err := TOSession.GetCDNs(opts)
	if err != nil {
		t.Fatalf("getting CDNs with name '%s': %v - alerts: %+v", cdnName, err, cdns.Alerts)
	}
	if len(cdns.Response) != 1 {
		t.Fatalf("expected exactly 1 CDN named '%s' but found %d CDNs", cdnName, len(cdns.Response))
	}
	opts.QueryParameters.Set("name", profileName)
	profiles, _, err := TOSession.GetProfiles(opts)
	if err != nil {
		t.Fatalf("getting Profiles with name '%s': %v - alerts: %+v", profileName, err, profiles.Alerts)
	}
	if len(profiles.Response) != 1 {
		t.Fatalf("expected exactly 1 Profiles named %s but found %d Profiles", profileName, len(profiles.Response))
	}
	parameters, _, err := TOSession.GetParametersByProfileName(profileName, client.RequestOptions{})
	if err != nil {
		t.Fatalf("getting Parameters by Profile name '%s': %v - alerts: %+v", profileName, err, parameters.Alerts)
	}
	parameterMap := map[string]tc.HealthThreshold{}
	parameterFound := map[string]bool{}
	const thresholdPrefixLength = len(tc.ThresholdPrefix)
	for _, parameter := range parameters.Response {
		if !strings.HasPrefix(parameter.Name, tc.ThresholdPrefix) {
			continue
		}
		parameterName := parameter.Name[thresholdPrefixLength:]
		parameterMap[parameterName], err = tc.StrToThreshold(parameter.Value)
		if err != nil {
			t.Fatalf("converting string '%s' to HealthThreshold: %s", parameter.Value, err.Error())
		}
		parameterFound[parameterName] = false
	}
	const expectedThresholdParameters = 3
	if len(parameterMap) != expectedThresholdParameters {
		t.Fatalf("expected Profile '%s' to contain %d Parameters with names starting with '%s' but %d such Parameters were found", profileName, expectedThresholdParameters, tc.ThresholdPrefix, len(parameterMap))
	}
	tmConfig, _, err := TOSession.GetTrafficMonitorConfig(cdnName, client.RequestOptions{})
	if err != nil {
		t.Fatalf("getting Traffic Monitor Config: %v - alerts: %+v", err, tmConfig.Alerts)
	}
	profileFound := false
	var profile tc.TMProfile
	for _, profile = range tmConfig.Response.Profiles {
		if profile.Name == profileName {
			profileFound = true
			break
		}
	}
	if !profileFound {
		t.Fatalf("Traffic Monitor Config contained no Profile named '%s", profileName)
	}
	for parameterName, value := range profile.Parameters.Thresholds {
		if _, ok := parameterFound[parameterName]; !ok {
			t.Fatalf("unexpected Threshold Parameter name '%s' found in Profile '%s' in Traffic Monitor Config", parameterName, profileName)
		}
		parameterFound[parameterName] = true
		if parameterMap[parameterName].String() != value.String() {
			t.Fatalf("expected '%s' but received '%s' for Threshold Parameter '%s' in Profile '%s' in Traffic Monitor Config", parameterMap[parameterName].String(), value.String(), parameterName, profileName)
		}
	}
	missingParameters := []string{}
	for parameterName, found := range parameterFound {
		if !found {
			missingParameters = append(missingParameters, parameterName)
		}
	}
	if len(missingParameters) != 0 {
		t.Fatalf("Threshold parameters defined for Profile '%s' but missing for Profile '%s' in Traffic Monitor Config: %s", profileName, profileName, strings.Join(missingParameters, ", "))
	}
}
