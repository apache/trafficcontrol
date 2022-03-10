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
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-rfc"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/assert"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
	client "github.com/apache/trafficcontrol/traffic_ops/v4-client"
)

func TestDeliveryServices(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Users, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, ServerCapabilities, ServiceCategories, DeliveryServices, DeliveryServicesRequiredCapabilities, DeliveryServiceServerAssignments}, func() {

		tomorrow := time.Now().AddDate(0, 0, 1).Format(time.RFC1123)
		currentTime := time.Now().UTC().Add(-5 * time.Second)
		currentTimeRFC := currentTime.Format(time.RFC1123)

		tenant4UserSession := utils.CreateV4Session(t, Config.TrafficOps.URL, "tenant4user", "pa$$word", Config.Default.Session.TimeoutInSecs)

		methodTests := utils.V4TestCase{
			"GET": {
				"NOT MODIFIED when NO CHANGES made": {
					ClientSession: TOSession, RequestOpts: client.RequestOptions{Header: http.Header{rfc.IfModifiedSince: {tomorrow}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"OK when VALID request": {
					ClientSession: TOSession, Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				//GetAccessibleToTest
				// ADDITIONAL TESTS: GET USING ROOT TENANT ID: 1 -> should match length of testdata
				// GET USING NEW TENANT BELONGING TO NO DSs -> length = 0 // tenant1 = child of root -> len = testdata - 1
				"OK when VALID ACCESSIBLETO parameter": {
					ClientSession: TOSession, RequestOpts: client.RequestOptions{QueryParameters: url.Values{"accessibleTo": {"1"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1)),
				},
				"OK when ACTIVE=TRUE": {
					ClientSession: TOSession, RequestOpts: client.RequestOptions{QueryParameters: url.Values{"active": {"true"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateDSExpectedFields(map[string]interface{}{"Active": true})),
				},
				"OK when ACTIVE=FALSE": {
					ClientSession: TOSession, RequestOpts: client.RequestOptions{QueryParameters: url.Values{"active": {"false"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateDSExpectedFields(map[string]interface{}{"Active": false})),
				},
				"OK when VALID CDN parameter": {
					ClientSession: TOSession, RequestOpts: client.RequestOptions{QueryParameters: url.Values{"cdn": {"cdn1"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateDSExpectedFields(map[string]interface{}{"CDNName": "cdn1"})),
				},
				"OK when VALID LOGSENABLED parameter": {
					ClientSession: TOSession, RequestOpts: client.RequestOptions{QueryParameters: url.Values{"logsEnabled": {"false"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateDSExpectedFields(map[string]interface{}{"LogsEnabled": false})),
				},
				"OK when VALID PROFILE parameter": {
					ClientSession: TOSession, RequestOpts: client.RequestOptions{QueryParameters: url.Values{"profile": {"ATS_EDGE_TIER_CACHE"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateDSExpectedFields(map[string]interface{}{"ProfileName": "ATS_EDGE_TIER_CACHE"})),
				},
				"OK when VALID SERVICECATEGORY parameter": {
					ClientSession: TOSession, RequestOpts: client.RequestOptions{QueryParameters: url.Values{"serviceCategory": {"serviceCategory1"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateDSExpectedFields(map[string]interface{}{"ServiceCategory": "serviceCategory1"})),
				},
				"OK when VALID TENANT parameter": {
					ClientSession: TOSession, RequestOpts: client.RequestOptions{QueryParameters: url.Values{"tenant": {"tenant1"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateDSExpectedFields(map[string]interface{}{"Tenant": "tenant1"})),
				},
				"OK when VALID TOPOLOGY parameter": {
					ClientSession: TOSession, RequestOpts: client.RequestOptions{QueryParameters: url.Values{"topology": {"mso-topology"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateDSExpectedFields(map[string]interface{}{"Topology": "mso-topology"})),
				},
				"OK when VALID TYPE parameter": {
					ClientSession: TOSession, RequestOpts: client.RequestOptions{QueryParameters: url.Values{"type": {"HTTP"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateDSExpectedFields(map[string]interface{}{"Type": tc.DSTypeHTTP})),
				},
				"OK when VALID XMLID parameter": {
					ClientSession: TOSession, RequestOpts: client.RequestOptions{QueryParameters: url.Values{"xmlId": {"ds1"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1),
						validateDSExpectedFields(map[string]interface{}{"XMLID": "ds1"})),
				},
				"EMPTY RESPONSE when INVALID ACCESSIBLETO parameter": {
					ClientSession: TOSession, RequestOpts: client.RequestOptions{QueryParameters: url.Values{"accessibleTo": {"10000"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
				"EMPTY RESPONSE when INVALID CDN parameter": {
					ClientSession: TOSession, RequestOpts: client.RequestOptions{QueryParameters: url.Values{"cdn": {"10000"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
				"EMPTY RESPONSE when INVALID PROFILE parameter": {
					ClientSession: TOSession, RequestOpts: client.RequestOptions{QueryParameters: url.Values{"profile": {"10000"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
				"EMPTY RESPONSE when INVALID TENANT parameter": {
					ClientSession: TOSession, RequestOpts: client.RequestOptions{QueryParameters: url.Values{"tenant": {"10000"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
				"EMPTY RESPONSE when INVALID TYPE parameter": {
					ClientSession: TOSession, RequestOpts: client.RequestOptions{QueryParameters: url.Values{"type": {"10000"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
				"EMPTY RESPONSE when INVALID XMLID parameter": {
					ClientSession: TOSession, RequestOpts: client.RequestOptions{QueryParameters: url.Values{"xmlId": {"invalid_xml_id"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
				"FIRST RESULT when LIMIT=1": {
					ClientSession: TOSession, RequestOpts: client.RequestOptions{QueryParameters: url.Values{"orderby": {"id"}, "limit": {"1"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validatePagination("limit")),
				},
				"SECOND RESULT when LIMIT=1 OFFSET=1": {
					ClientSession: TOSession, RequestOpts: client.RequestOptions{QueryParameters: url.Values{"orderby": {"id"}, "limit": {"1"}, "offset": {"1"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validatePagination("offset")),
				},
				"SECOND RESULT when LIMIT=1 PAGE=2": {
					ClientSession: TOSession, RequestOpts: client.RequestOptions{QueryParameters: url.Values{"orderby": {"id"}, "limit": {"1"}, "page": {"2"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validatePagination("page")),
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
				"VALID when SORTORDER param is DESC": {
					ClientSession: TOSession, RequestOpts: client.RequestOptions{QueryParameters: url.Values{"sortOrder": {"desc"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateDescSort()),
				},
				"EMPTY RESPONSE when TENANT attempts reading DS OUTSIDE TENANCY": {
					ClientSession: tenant4UserSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"xmlId": {"ds3"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
			},
			"POST": {
				"CREATED when VALID request WITH GEO LIMIT COUNTRIES": {
					ClientSession: TOSession,
					RequestBody: generateDeliveryService(t, map[string]interface{}{
						"geoLimit":          2,
						"geoLimitCountries": []string{"US", "CA"},
						"xmlId":             "geolimit-test",
					}),
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusCreated), utils.ResponseHasLength(1),
						validateDSExpectedFields(map[string]interface{}{"GeoLimitCountries": tc.GeoLimitCountriesType{"US", "CA"}})),
				},
				"BAD REQUEST when using LONG DESCRIPTION 2 and 3 fields": {
					ClientSession: TOSession,
					RequestBody: generateDeliveryService(t, map[string]interface{}{
						"longDesc1": "long desc 1",
						"longDesc2": "long desc 2",
						"xmlId":     "ld1-ld2-test",
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when XMLID left EMPTY": {
					ClientSession: TOSession,
					RequestBody: generateDeliveryService(t, map[string]interface{}{
						"xmlId": "",
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when XMLID is NIL": {
					ClientSession: TOSession,
					RequestBody: generateDeliveryService(t, map[string]interface{}{
						"xmlId": nil,
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when TOPOLOGY DOESNT EXIST": {
					ClientSession: TOSession,
					RequestBody: generateDeliveryService(t, map[string]interface{}{
						"topology": "topology-doesnt-exist",
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when creating STEERING DS with TLS VERSIONS": {
					ClientSession: TOSession,
					RequestBody: generateDeliveryService(t, map[string]interface{}{
						"tlsVersions": []string{"1.1"},
						"typeId":      GetTypeId(t, "STEERING"),
						"xmlId":       "test-TLS-creation-steering",
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"OK when creating HTTP DS with TLS VERSIONS": {
					ClientSession: TOSession,
					RequestBody: generateDeliveryService(t, map[string]interface{}{
						"tlsVersions": []string{"1.1"},
						"xmlId":       "test-TLS-creation-http",
					}),
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusCreated), utils.ResponseHasLength(1)),
				},
				"BAD REQUEST when creating DS with TENANCY NOT THE SAME AS CURRENT TENANT": {
					ClientSession: tenant4UserSession,
					RequestBody: generateDeliveryService(t, map[string]interface{}{
						"tenantId": GetTenantId(t, "tenant3"),
						"xmlId":    "test-tenancy",
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden), utils.ResponseHasLength(0)),
				},
			},
			"PUT": {
				"BAD REQUEST when using LONG DESCRIPTION 2 and 3 fields": {
					EndpointId: GetDeliveryServiceId(t, "ds1"), ClientSession: TOSession,
					RequestBody: generateDeliveryService(t, map[string]interface{}{
						"longDesc1": "long desc 1",
						"longDesc2": "long desc 2",
						"xmlId":     "ds1",
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				//"OK when VALID request": {
				//	EndpointId: GetDeliveryServiceId(t, "ds1"), ClientSession: TOSession,
				//	RequestBody: map[string]interface{}{
				//		"active": false,
				//		"cdnName": "cdn1",
				//		"ccrDnsTtl": 3600,
				//		"checkPath": "",
				//		"consistentHashQueryParams": [],
				//		"deepCachingType": "NEVER",
				//		"displayName": "newds1displayname",
				//		"dnsBypassCname": null,
				//		"dnsBypassIp": "",
				//		"dnsBypassIp6": "",
				//		"dnsBypassTtl": 30,
				//		"dscp": 41,
				//		"edgeHeaderRewrite": "edgeHeader1\nedgeHeader2",
				//		"exampleURLs": [
				//		"http://ccr.ds1.example.net",
				//		"https://ccr.ds1.example.net"
				//	],
				//		"fqPacingRate": 0,
				//		"geoLimit": 1,
				//		"geoLimitCountries": "",
				//		"geoLimitRedirectURL": null,
				//		"geoProvider": 0,
				//		"globalMaxMbps": 0,
				//		"globalMaxTps": 0,
				//		"httpBypassFqdn": "",
				//		"infoUrl": "TBD",
				//		"initialDispersion": 2,
				//		"ipv6RoutingEnabled": false,
				//		"logsEnabled": true,
				//		"longDesc": "something different",
				//		"longDesc1": "ds1",
				//		"longDesc2": "ds1",
				//		"matchList": [
				//	{
				//		"pattern": ".*\\.ds1\\..*",
				//		"setNumber": 0,
				//		"type": "HOST_REGEXP"
				//	}
				//	],
				//		"maxDnsAnswers": 164598,
				//		"midHeaderRewrite": "midHeader1\nmidHeader2",
				//		"missLat": 42.881944,
				//		"missLong": -88.627778,
				//		"multiSiteOrigin": true,
				//		"orgServerFqdn": "http://origin.update.example.net",
				//		"originShield": null,
				//		"profileDescription": null,
				//		"profileName": "ATS_EDGE_TIER_CACHE",
				//		"protocol": 2,
				//		"qstringIgnore": 0,
				//		"rangeRequestHandling": 0,
				//		"regexRemap": "rr1\nrr2",
				//		"regionalGeoBlocking": true,
				//		"remapText": "@plugin=tslua.so @pparam=/opt/trafficserver/etc/trafficserver/remapPlugin1.lua",
				//		"routingName": "ccr-ds1",
				//		"signed": false,
				//		"signingAlgorithm": "url_sig",
				//		"sslKeyVersion": 2,
				//		"tenant": "tenant1",
				//		"tenantName": "tenant1",
				//		"type": "HTTP",
				//		"xmlId": "ds1",
				//		"anonymousBlockingEnabled": true,
				//		"maxOriginConnections": 100
				//		"maxRequestHeaderBytes": 131080
				//	},
				//	Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				//},
				//UpdateValidateORGServerCacheGroup
				//"Assign an Origin not in a Cache Group used by a Delivery Service's Topology to that Delivery Service": {},
				"BAD REQUEST when INVALID REMAP TEXT": {
					EndpointId: GetDeliveryServiceId(t, "ds1"), ClientSession: TOSession,
					RequestBody: generateDeliveryService(t, map[string]interface{}{
						"remapText": "@plugin=tslua.so @pparam=/opt/trafficserver/etc/trafficserver/remapPlugin1.lua\nline2",
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when MISSING SLICE PLUGIN SIZE": {
					EndpointId: GetDeliveryServiceId(t, "ds1"), ClientSession: TOSession,
					RequestBody: generateDeliveryService(t, map[string]interface{}{
						"rangeRequestHandling": 3,
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when SLICE PLUGIN SIZE SET with INVALID RANGE REQUEST SETTING": {
					EndpointId: GetDeliveryServiceId(t, "ds1"), ClientSession: TOSession,
					RequestBody: generateDeliveryService(t, map[string]interface{}{
						"rangeRequestHandling": 1,
						"rangeSliceBlockSize":  262144,
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when SLICE PLUGIN SIZE TOO SMALL": {
					EndpointId: GetDeliveryServiceId(t, "ds1"), ClientSession: TOSession,
					RequestBody: generateDeliveryService(t, map[string]interface{}{
						"rangeRequestHandling": 3,
						"rangeSliceBlockSize":  0,
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when SLICE PLUGIN SIZE TOO LARGE": {
					EndpointId: GetDeliveryServiceId(t, "ds1"), ClientSession: TOSession,
					RequestBody: generateDeliveryService(t, map[string]interface{}{
						"rangeRequestHandling": 3,
						"rangeSliceBlockSize":  40000000,
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when ADDING TOPOLOGY to CLIENT STEERING DS": {
					EndpointId: GetDeliveryServiceId(t, "ds-client-steering"), ClientSession: TOSession,
					RequestBody: generateDeliveryService(t, map[string]interface{}{
						"topology": "mso-topology",
						"xmlId":    "ds-client-steering",
						"typeId":   GetTypeId(t, "CLIENT_STEERING"),
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when TOPOLOGY DOESNT EXIST": {
					EndpointId: GetDeliveryServiceId(t, "ds1"), ClientSession: TOSession,
					RequestBody: generateDeliveryService(t, map[string]interface{}{
						"topology": "",
						"xmlId":    "ds1",
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when ADDING TOPOLOGY to DS with DS REQUIRED CAPABILITY": {
					EndpointId: GetDeliveryServiceId(t, "ds1"), ClientSession: TOSession,
					RequestBody: generateDeliveryService(t, map[string]interface{}{
						"topology": "top-for-ds-req",
						"xmlId":    "ds1",
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"OK when REMOVING TOPOLOGY": {
					EndpointId: GetDeliveryServiceId(t, "ds-based-top-with-no-mids"), ClientSession: TOSession,
					RequestBody: generateDeliveryService(t, map[string]interface{}{
						"topology": nil,
						"xmlId":    "ds-based-top-with-no-mids",
					}),
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				//t.Fatalf("expected 400-level error assigning Topology %s to Delivery Service %s because Cache Group %s has no Servers in it in CDN %d, no error received", *dsTopology, xmlID, cacheGroupName, *ds.CDNID)
				// "top-ds-in-cdn2"
				//"BAD REQUEST when ASSIGNING TOPOLOGY when CG has NO SERVERS": {
				//	EndpointId: GetDeliveryServiceId(t, "top-ds-in-cdn2"), ClientSession: TOSession,
				//	RequestBody: map[string]interface{}{
				//		"topology": "top-cg-no-servers",
				//		"xmlId":    "top-ds-in-cdn2",
				//	},
				//	Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				//},
				"OK when DS with TOPOLOGY updates HEADER REWRITE FIELDS": {
					EndpointId: GetDeliveryServiceId(t, "ds-top"), ClientSession: TOSession,
					RequestBody: generateDeliveryService(t, map[string]interface{}{
						"firstHeaderRewrite": "foo",
						"innerHeaderRewrite": "bar",
						"lastHeaderRewrite":  "baz",
						"topology":           "mso-topology",
						"xmlId":              "ds-top",
					}),
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"BAD REQUEST when DS with NO TOPOLOGY updates HEADER REWRITE FIELDS": {
					EndpointId: GetDeliveryServiceId(t, "ds1"), ClientSession: TOSession,
					RequestBody: generateDeliveryService(t, map[string]interface{}{
						"firstHeaderRewrite": "foo",
						"innerHeaderRewrite": "bar",
						"lastHeaderRewrite":  "baz",
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when DS with TOPOLOGY updates LEGACY HEADER REWRITE FIELDS": {
					EndpointId: GetDeliveryServiceId(t, "ds-top"), ClientSession: TOSession,
					RequestBody: generateDeliveryService(t, map[string]interface{}{
						"edgeHeaderRewrite": "foo",
						"midHeaderRewrite":  "bar",
						"topology":          "mso-topology",
						"xmlId":             "ds-top",
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"OK when DS with NO TOPOLOGY updates LEGACY HEADER REWRITE FIELDS": {
					EndpointId: GetDeliveryServiceId(t, "ds1"), ClientSession: TOSession,
					RequestBody: generateDeliveryService(t, map[string]interface{}{
						"profileId":         GetProfileId(t, "ATS_EDGE_TIER_CACHE"),
						"edgeHeaderRewrite": "foo",
						"midHeaderRewrite":  "bar",
						"xmlId":             "ds1",
					}),
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"OK when UPDATING MINOR VERSION FIELDS": {
					EndpointId: GetDeliveryServiceId(t, "ds-test-minor-versions"), ClientSession: TOSession,
					RequestBody: generateDeliveryService(t, map[string]interface{}{
						"consistentHashQueryParams": []string{"d", "e", "f"},
						"consistentHashRegex":       "foo",
						"deepCachingType":           "NEVER",
						"fqPacingRate":              41,
						"maxOriginConnections":      500,
						"routingName":               "cdn",
						"signingAlgorithm":          "uri_signing",
						"tenantId":                  GetTenantId(t, "tenant1"),
						"trRequestHeaders":          "X-ooF\nX-raB",
						"trResponseHeaders":         "Access-Control-Max-Age: 600\nContent-Type: text/html; charset=utf-8",
						"xmlId":                     "ds-test-minor-versions",
					}),
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateDSExpectedFields(map[string]interface{}{"ConsistentHashQueryParams": []string{"d", "e", "f"},
							"ConsistentHashRegex": "foo", "DeepCachingType": tc.DeepCachingTypeNever, "FQPacingRate": 41, "MaxOriginConnections": 500,
							"SigningAlgorithm": "uri_signing", "Tenant": "tenant1", "TRRequestHeaders": "X-ooF\nX-raB",
							"TRResponseHeaders": "Access-Control-Max-Age: 600\nContent-Type: text/html; charset=utf-8",
						})),
				},
				"BAD REQUEST when INVALID COUNTRY CODE": {
					EndpointId: GetDeliveryServiceId(t, "ds1"), ClientSession: TOSession,
					RequestBody: generateDeliveryService(t, map[string]interface{}{
						"geoLimit":          2,
						"geoLimitCountries": []string{"US", "CA", "12"},
						"xmlId":             "invalid-geolimit-test",
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when CHANGING TOPOLOGY of DS with ORG SERVERS ASSIGNED": {
					EndpointId: GetDeliveryServiceId(t, "ds-top"), ClientSession: TOSession,
					RequestBody: generateDeliveryService(t, map[string]interface{}{
						"topology": "another-topology",
						"xmlId":    "ds-top",
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when UPDATING DS OUTSIDE TENANCY": {
					EndpointId: GetDeliveryServiceId(t, "ds3"), ClientSession: tenant4UserSession,
					RequestBody:  generateDeliveryService(t, map[string]interface{}{"xmlId": "ds3"}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
				"PRECONDITION FAILED when updating with IMS & IUS Headers": {
					EndpointId: GetDeliveryServiceId(t, "ds1"), ClientSession: TOSession,
					RequestOpts: client.RequestOptions{
						Header: http.Header{
							rfc.IfModifiedSince: {currentTimeRFC}, rfc.IfUnmodifiedSince: {currentTimeRFC},
						},
					},
					RequestBody:  generateDeliveryService(t, map[string]interface{}{"xmlId": "ds1"}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
				"PRECONDITION FAILED when updating with IFMATCH ETAG Header": {
					EndpointId: GetDeliveryServiceId(t, "ds1"), ClientSession: TOSession,
					RequestBody:  generateDeliveryService(t, map[string]interface{}{"xmlId": "ds1"}),
					RequestOpts:  client.RequestOptions{Header: http.Header{rfc.IfMatch: {rfc.ETag(currentTime)}}},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
			},
			"DELETE": {
				"BAD REQUEST when DELETING DS OUTSIDE TENANCY": {
					EndpointId: GetDeliveryServiceId(t, "ds3"), ClientSession: tenant4UserSession,
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
			},
			"GET AFTER CHANGES": {
				"OK when CHANGES made": {
					ClientSession: TOSession,
					RequestOpts: client.RequestOptions{
						Header: http.Header{
							rfc.IfModifiedSince: {currentTimeRFC}, rfc.IfUnmodifiedSince: {currentTimeRFC},
						},
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
			},
			//"DELIVERY SERVICES CAPACITY": {
			//	// capDS, _, err := TOSession.GetDeliveryServiceCapacity(*ds.ID, client.RequestOptions{})
			//	"Basic GET request for /deliveryservices/{{ID}}/capacity": {},
			//},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					ds := tc.DeliveryServiceV4{}

					if val, ok := testCase.RequestOpts.QueryParameters["accessibleTo"]; ok {
						if _, err := strconv.Atoi(val[0]); err != nil {
							testCase.RequestOpts.QueryParameters.Set("accessibleTo", strconv.Itoa(GetTenantId(t, val[0])))
						}
					}
					if val, ok := testCase.RequestOpts.QueryParameters["cdn"]; ok {
						if _, err := strconv.Atoi(val[0]); err != nil {
							testCase.RequestOpts.QueryParameters.Set("cdn", strconv.Itoa(GetCDNId(t, val[0])))
						}
					}
					if val, ok := testCase.RequestOpts.QueryParameters["profile"]; ok {
						if _, err := strconv.Atoi(val[0]); err != nil {
							testCase.RequestOpts.QueryParameters.Set("profile", strconv.Itoa(GetProfileId(t, val[0])))
						}
					}
					if val, ok := testCase.RequestOpts.QueryParameters["type"]; ok {
						if _, err := strconv.Atoi(val[0]); err != nil {
							testCase.RequestOpts.QueryParameters.Set("type", strconv.Itoa(GetTypeId(t, val[0])))
						}
					}
					if val, ok := testCase.RequestOpts.QueryParameters["tenant"]; ok {
						if _, err := strconv.Atoi(val[0]); err != nil {
							testCase.RequestOpts.QueryParameters.Set("tenant", strconv.Itoa(GetTenantId(t, val[0])))
						}
					}

					if testCase.RequestBody != nil {
						dat, err := json.Marshal(testCase.RequestBody)
						assert.NoError(t, err, "Error occurred when marshalling request body: %v", err)
						err = json.Unmarshal(dat, &ds)
						assert.NoError(t, err, "Error occurred when unmarshalling request body: %v", err)
					}

					switch method {
					case "GET", "GET AFTER CHANGES":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.GetDeliveryServices(testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							fmt.Println(*ds.Tenant)
							resp, reqInf, err := testCase.ClientSession.CreateDeliveryService(ds, testCase.RequestOpts)
							fmt.Println(resp)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "PUT":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.UpdateDeliveryService(testCase.EndpointId(), ds, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.DeleteDeliveryService(testCase.EndpointId(), testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, resp.Alerts, err)
							}
						})
					}
				}
			})
		}
	})
}

func validateDSExpectedFields(expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		dsResp := resp.([]tc.DeliveryServiceV40)
		for field, expected := range expectedResp {
			for _, ds := range dsResp {
				switch field {
				case "Active":
					assert.Equal(t, expected, *ds.Active, "Expected active to be %v, but got %v", expected, *ds.Active)
				case "DeepCachingType":
					assert.Equal(t, expected, *ds.DeepCachingType, "Expected deepCachingType to be %v, but got %v", expected, *ds.DeepCachingType)
				case "CDNName":
					assert.Equal(t, expected, *ds.CDNName, "Expected CDNName to be %v, but got %v", expected, *ds.CDNName)
				case "ConsistentHashRegex":
					assert.Equal(t, expected, *ds.ConsistentHashRegex, "Expected ConsistentHashRegex to be %v, but got %v", expected, *ds.ConsistentHashRegex)
				case "ConsistentHashQueryParams":
					assert.Exactly(t, expected, ds.ConsistentHashQueryParams, "Expected ConsistentHashQueryParams to be %v, but got %v", expected, ds.ConsistentHashQueryParams)
				case "FQPacingRate":
					assert.Equal(t, expected, *ds.FQPacingRate, "Expected FQPacingRate to be %v, but got %v", expected, *ds.FQPacingRate)
				case "GeoLimitCountries":
					assert.Exactly(t, expected, ds.GeoLimitCountries, "Expected GeoLimitCountries to be %v, but got &v", expected, ds.GeoLimitCountries)
				case "LogsEnabled":
					assert.Equal(t, expected, *ds.LogsEnabled, "Expected LogsEnabled to be %v, but got %v", expected, *ds.LogsEnabled)
				case "MaxOriginConnections":
					assert.Equal(t, expected, *ds.MaxOriginConnections, "Expected MaxOriginConnections to be %v, but got %v", expected, *ds.MaxOriginConnections)
				case "ProfileName":
					assert.Equal(t, expected, *ds.ProfileName, "Expected ProfileName to be %v, but got %v", expected, *ds.ProfileName)
				case "ServiceCategory":
					assert.Equal(t, expected, *ds.ServiceCategory, "Expected ServiceCategory to be %v, but got %v", expected, *ds.ServiceCategory)
				case "SigningAlgorithm":
					assert.Equal(t, expected, *ds.SigningAlgorithm, "Expected SigningAlgorithm to be %v, but got %v", expected, *ds.SigningAlgorithm)
				case "Tenant":
					assert.Equal(t, expected, *ds.Tenant, "Expected Topology to be %v, but got %v", expected, *ds.Tenant)
				case "Topology":
					assert.Equal(t, expected, *ds.Topology, "Expected Topology to be %v, but got %v", expected, *ds.Topology)
				case "TRRequestHeaders":
					assert.Equal(t, expected, *ds.TRRequestHeaders, "Expected TRRequestHeaders to be %v, but got %v", expected, *ds.TRRequestHeaders)
				case "TRResponseHeaders":
					assert.Equal(t, expected, *ds.TRResponseHeaders, "Expected TRResponseHeaders to be %v, but got %v", expected, *ds.TRResponseHeaders)
				case "Type":
					assert.Equal(t, expected, *ds.Type, "Expected Type to be %v, but got %v", expected, *ds.Type)
				case "XMLID":
					assert.Equal(t, expected, *ds.XMLID, "Expected XMLID to be %v, but got %v", expected, *ds.XMLID)
				default:
					t.Errorf("Expected field: %v, does not exist in response", field)
				}
			}
		}
	}
}

func validatePagination(paginationParam string) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		paginationResp := resp.([]tc.DeliveryServiceV40)

		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("orderby", "id")
		respBase, _, err := TOSession.GetDeliveryServices(opts)
		assert.RequireNoError(t, err, "Cannot get Delivery Services: %v - alerts: %+v", err, respBase.Alerts)

		ds := respBase.Response
		assert.RequireGreaterOrEqual(t, len(ds), 3, "Need at least 3 Delivery Services in Traffic Ops to test pagination support, found: %d", len(ds))
		switch paginationParam {
		case "limit:":
			assert.Exactly(t, ds[:1], paginationResp, "expected GET Delivery Services with limit = 1 to return first result")
		case "offset":
			assert.Exactly(t, ds[1:2], paginationResp, "expected GET Delivery Services with limit = 1, offset = 1 to return second result")
		case "page":
			assert.Exactly(t, ds[1:2], paginationResp, "expected GET Delivery Services with limit = 1, page = 2 to return second result")
		}
	}
}

func validateDescSort() utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, alerts tc.Alerts, _ error) {
		dsResp := resp.([]tc.DeliveryServiceV40)
		var sortedList []string
		assert.RequireGreaterOrEqual(t, len(dsResp), 2, "Need at least 2 XMLIDs in Traffic Ops to test desc sort, found: %d", len(dsResp))

		for _, ds := range dsResp {
			sortedList = append(sortedList, *ds.XMLID)
		}

		res := sort.SliceIsSorted(sortedList, func(p, q int) bool {
			return sortedList[p] > sortedList[q]
		})
		assert.Equal(t, res, true, "List is not sorted by their XMLIDs: %v", sortedList)
	}
}

func GetCDNId(t *testing.T, cdnName string) int {
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("name", cdnName)
	resp, _, err := TOSession.GetCDNs(opts)

	assert.RequireNoError(t, err, "Get CDNs Request failed with error: %v", err)
	assert.RequireEqual(t, 1, len(resp.Response), "Expected response object length 1, but got %d", len(resp.Response))
	assert.RequireNotNil(t, &resp.Response[0].ID, "Expected id to not be nil")

	return resp.Response[0].ID
}

func GetDeliveryServiceId(t *testing.T, xmlId string) func() int {
	return func() int {
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("xmlId", xmlId)

		resp, _, err := TOSession.GetDeliveryServices(opts)
		assert.RequireNoError(t, err, "Get Delivery Service Request failed with error: %v", err)
		assert.RequireEqual(t, len(resp.Response), 1, "Expected response object length 1, but got %d", len(resp.Response))
		assert.RequireNotNil(t, resp.Response[0].ID, "Expected id to not be nil")

		return *resp.Response[0].ID
	}
}

func GetProfileId(t *testing.T, profileName string) int {
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("name", profileName)
	resp, _, err := TOSession.GetProfiles(opts)

	assert.RequireNoError(t, err, "Get Profiles Request failed with error: %v", err)
	assert.RequireEqual(t, 1, len(resp.Response), "Expected response object length 1, but got %d", len(resp.Response))
	assert.RequireNotNil(t, &resp.Response[0].ID, "Expected id to not be nil")

	return resp.Response[0].ID
}

func GetTenantId(t *testing.T, tenantName string) int {
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("name", tenantName)
	resp, _, err := TOSession.GetTenants(opts)

	assert.RequireNoError(t, err, "Get Tenants Request failed with error: %v", err)
	assert.RequireEqual(t, 1, len(resp.Response), "Expected response object length 1, but got %d", len(resp.Response))
	assert.RequireNotNil(t, &resp.Response[0].ID, "Expected id to not be nil")

	return resp.Response[0].ID
}

func generateDeliveryService(t *testing.T, requestDS map[string]interface{}) map[string]interface{} {
	// map for the most basic HTTP Delivery Service a user can create
	genericHTTPDS := map[string]interface{}{
		"active":               true,
		"cdnName":              "cdn1",
		"cdnId":                GetCDNId(t, "cdn1"),
		"displayName":          "test ds",
		"dscp":                 0,
		"geoLimit":             0,
		"geoProvider":          0,
		"initialDispersion":    1,
		"ipv6RoutingEnabled":   false,
		"logsEnabled":          false,
		"missLat":              0.0,
		"missLong":             0.0,
		"multiSiteOrigin":      false,
		"orgServerFqdn":        "http://ds.test",
		"protocol":             0,
		"qstringIgnore":        0,
		"rangeRequestHandling": 0,
		"regionalGeoBlocking":  false,
		"routingName":          "ccr-ds1",
		"tenant":               "tenant1",
		"type":                 tc.DSTypeHTTP,
		"typeId":               GetTypeId(t, "HTTP"),
		"xmlId":                "da1",
	}
	for k, v := range requestDS {
		genericHTTPDS[k] = v
	}
	return genericHTTPDS
}

func CreateTestDeliveryServices(t *testing.T) {
	for _, ds := range testData.DeliveryServices {
		ds = ds.RemoveLD1AndLD2()
		if ds.XMLID == nil {
			t.Error("Found a Delivery Service in testing data with null or undefined XMLID")
			continue
		}
		resp, _, err := TOSession.CreateDeliveryService(ds, client.RequestOptions{})
		assert.NoError(t, err, "Could not create Delivery Service '%s': %v - alerts: %+v", *ds.XMLID, err, resp.Alerts)
	}
}

func DeleteTestDeliveryServices(t *testing.T) {
	dses, _, err := TOSession.GetDeliveryServices(client.RequestOptions{})
	assert.NoError(t, err, "Cannot get Delivery Services: %v - alerts: %+v", err, dses.Alerts)

	for _, ds := range dses.Response {
		delResp, _, err := TOSession.DeleteDeliveryService(*ds.ID, client.RequestOptions{})
		assert.NoError(t, err, "Could not delete Delivery Service: %v - alerts: %+v", err, delResp.Alerts)
		// Retrieve Delivery Service to see if it got deleted
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("id", strconv.Itoa(*ds.ID))
		getDS, _, err := TOSession.GetDeliveryServices(opts)
		assert.NoError(t, err, "Error deleting Delivery Service for '%s' : %v - alerts: %+v", *ds.XMLID, err, getDS.Alerts)
		assert.Equal(t, 0, len(getDS.Response), "Expected Delivery Service '%s' to be deleted", *ds.XMLID)
	}
}
