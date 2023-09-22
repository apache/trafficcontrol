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
	"sort"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
	client "github.com/apache/trafficcontrol/v8/traffic_ops/v4-client"
)

func TestStaticDNSEntries(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, ServiceCategories, DeliveryServices, StaticDNSEntries}, func() {

		currentTime := time.Now().UTC().Add(-15 * time.Second)
		currentTimeRFC := currentTime.Format(time.RFC1123)
		tomorrow := currentTime.AddDate(0, 0, 1).Format(time.RFC1123)

		methodTests := utils.TestCase[client.Session, client.RequestOptions, tc.StaticDNSEntry]{
			"GET": {
				"NOT MODIFIED when NO CHANGES made": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfModifiedSince: {tomorrow}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"OK when VALID request": {
					ClientSession: TOSession,
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateStaticDNSEntriesSort()),
				},
				"OK when VALID HOST parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"host": {"host1"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1),
						validateStaticDNSEntriesFields(map[string]interface{}{"Host": "host1"})),
				},
			},
			"PUT": {
				"OK when VALID request": {
					EndpointID:    GetStaticDNSEntryID(t, "host2"),
					ClientSession: TOSession,
					RequestBody: tc.StaticDNSEntry{
						Address:         "192.168.0.2",
						CacheGroupName:  "cachegroup2",
						DeliveryService: "ds2",
						Host:            "host2",
						Type:            "A_RECORD",
						TTL:             10,
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateStaticDNSEntriesUpdateCreateFields("host2", map[string]interface{}{"Address": "192.168.0.2"})),
				},
				"BAD REQUEST when INVALID IPV4 ADDRESS for A_RECORD": {
					EndpointID:    GetStaticDNSEntryID(t, "host2"),
					ClientSession: TOSession,
					RequestBody: tc.StaticDNSEntry{
						Address:         "test.testdomain.net.",
						CacheGroupName:  "cachegroup2",
						DeliveryService: "ds2",
						Host:            "host2",
						Type:            "A_RECORD",
						TTL:             10,
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when INVALID DNS for CNAME_RECORD": {
					EndpointID:    GetStaticDNSEntryID(t, "host1"),
					ClientSession: TOSession,
					RequestBody: tc.StaticDNSEntry{
						Address:         "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
						CacheGroupName:  "cachegroup1",
						DeliveryService: "ds1",
						Host:            "host1",
						Type:            "CNAME_RECORD",
						TTL:             0,
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when MISSING TRAILING PERIOD for CNAME_RECORD": {
					EndpointID:    GetStaticDNSEntryID(t, "host1"),
					ClientSession: TOSession,
					RequestBody: tc.StaticDNSEntry{
						Address:         "cdn.test.com",
						CacheGroupName:  "cachegroup1",
						DeliveryService: "ds1",
						Host:            "host1",
						Type:            "CNAME_RECORD",
						TTL:             0,
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when INVALID IPV6 ADDRESS for AAAA_RECORD": {
					EndpointID:    GetStaticDNSEntryID(t, "host3"),
					ClientSession: TOSession,
					RequestBody: tc.StaticDNSEntry{
						Address:         "192.168.0.1",
						CacheGroupName:  "cachegroup2",
						DeliveryService: "ds1",
						Host:            "host3",
						TTL:             10,
						Type:            "AAAA_RECORD",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"PRECONDITION FAILED when updating with IMS & IUS Headers": {
					EndpointID:    GetStaticDNSEntryID(t, "host3"),
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfUnmodifiedSince: {currentTimeRFC}}},
					RequestBody: tc.StaticDNSEntry{
						Address:         "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
						CacheGroupName:  "cachegroup2",
						DeliveryService: "ds1",
						Host:            "host3",
						TTL:             10,
						Type:            "AAAA_RECORD",
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
				"PRECONDITION FAILED when updating with IFMATCH ETAG Header": {
					EndpointID:    GetStaticDNSEntryID(t, "host3"),
					ClientSession: TOSession,
					RequestBody: tc.StaticDNSEntry{
						Address:         "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
						CacheGroupName:  "cachegroup2",
						DeliveryService: "ds1",
						Host:            "host3",
						TTL:             10,
						Type:            "AAAA_RECORD",
					},
					RequestOpts:  client.RequestOptions{Header: http.Header{rfc.IfMatch: {rfc.ETag(currentTime)}}},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					switch method {
					case "GET":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.GetStaticDNSEntries(testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.CreateStaticDNSEntry(testCase.RequestBody, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					case "PUT":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.UpdateStaticDNSEntry(testCase.EndpointID(), testCase.RequestBody, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.DeleteStaticDNSEntry(testCase.EndpointID(), testCase.RequestOpts)
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

func validateStaticDNSEntriesFields(expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Static DNS Entries response to not be nil.")
		staticDNSEntriesResp := resp.([]tc.StaticDNSEntry)
		for field, expected := range expectedResp {
			for _, staticDNSEntry := range staticDNSEntriesResp {
				switch field {
				case "Address":
					assert.Equal(t, expected, staticDNSEntry.Address, "Expected Address to be %v, but got %s", expected, staticDNSEntry.Address)
				case "Host":
					assert.Equal(t, expected, staticDNSEntry.Host, "Expected Host to be %v, but got %s", expected, staticDNSEntry.Host)
				default:
					t.Errorf("Expected field: %v, does not exist in response", field)
				}
			}
		}
	}
}

func validateStaticDNSEntriesUpdateCreateFields(host string, expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("host", host)
		staticDNSEntries, _, err := TOSession.GetStaticDNSEntries(opts)
		assert.RequireNoError(t, err, "Error getting Static DNS Entries: %v - alerts: %+v", err, staticDNSEntries.Alerts)
		assert.RequireEqual(t, 1, len(staticDNSEntries.Response), "Expected one Static DNS Entry returned Got: %d", len(staticDNSEntries.Response))
		validateStaticDNSEntriesFields(expectedResp)(t, toclientlib.ReqInf{}, staticDNSEntries.Response, tc.Alerts{}, nil)
	}
}

func validateStaticDNSEntriesSort() utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, alerts tc.Alerts, _ error) {
		assert.RequireNotNil(t, resp, "Expected Static DNS Entries response to not be nil.")
		var staticDNSEntryHosts []string
		staticDNSEntryResp := resp.([]tc.StaticDNSEntry)
		for _, staticDNSEntry := range staticDNSEntryResp {
			staticDNSEntryHosts = append(staticDNSEntryHosts, staticDNSEntry.Host)
		}
		assert.Equal(t, true, sort.StringsAreSorted(staticDNSEntryHosts), "List is not sorted by their hosts: %v", staticDNSEntryHosts)
	}
}

func GetStaticDNSEntryID(t *testing.T, host string) func() int {
	return func() int {
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("host", host)
		staticDNSEntries, _, err := TOSession.GetStaticDNSEntries(opts)
		assert.RequireNoError(t, err, "Get Static DNS Entries Request failed with error:", err)
		assert.RequireEqual(t, 1, len(staticDNSEntries.Response), "Expected response object length 1, but got %d", len(staticDNSEntries.Response))
		return staticDNSEntries.Response[0].ID
	}
}
