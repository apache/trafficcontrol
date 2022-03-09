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
	"github.com/apache/trafficcontrol/lib/go-util"
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

		tenant4UserSession := createSession(t, "tenant4user", "pa$$word")

		methodTests := map[string]map[string]struct {
			endpointId    func() int
			clientSession *client.Session
			requestOpts   client.RequestOptions
			requestBody   map[string]interface{}
			expectations  []utils.CkReqFunc
		}{
			"GET": {
				"NOT MODIFIED when NO CHANGES made": {
					clientSession: TOSession, requestOpts: client.RequestOptions{Header: http.Header{rfc.IfModifiedSince: {tomorrow}}},
					expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				// Performs GET -> Validates against test data
				"OK when VALID request": {
					clientSession: TOSession, expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				//GetAccessibleToTest
				// ADDITIONAL TESTS: GET USING ROOT TENANT ID: 1 -> should match length of testdata
				// GET USING NEW TENANT BELONGING TO NO DSs -> length = 0 // tenant1 = child of root -> len = testdata - 1
				"OK when VALID ACCESSIBLETO parameter": {
					clientSession: TOSession, requestOpts: client.RequestOptions{QueryParameters: url.Values{"accessibleTo": {"1"}}},
					expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1)),
				},
				"OK when ACTIVE=TRUE": {
					clientSession: TOSession, requestOpts: client.RequestOptions{QueryParameters: url.Values{"active": {"true"}}},
					expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateExpectedFields(map[string]interface{}{"Active": true})),
				},
				"OK when ACTIVE=FALSE": {
					clientSession: TOSession, requestOpts: client.RequestOptions{QueryParameters: url.Values{"active": {"false"}}},
					expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateExpectedFields(map[string]interface{}{"Active": false})),
				},
				"OK when VALID CDN parameter": {
					clientSession: TOSession, requestOpts: client.RequestOptions{QueryParameters: url.Values{"cdn": {"cdn1"}}},
					expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateExpectedFields(map[string]interface{}{"CDNName": "cdn1"})),
				},
				"OK when VALID LOGSENABLED parameter": {
					clientSession: TOSession, requestOpts: client.RequestOptions{QueryParameters: url.Values{"logsEnabled": {"false"}}},
					expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateExpectedFields(map[string]interface{}{"LogsEnabled": false})),
				},
				"OK when VALID PROFILE parameter": {
					clientSession: TOSession, requestOpts: client.RequestOptions{QueryParameters: url.Values{"profile": {"ATS_EDGE_TIER_CACHE"}}},
					expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateExpectedFields(map[string]interface{}{"ProfileName": "ATS_EDGE_TIER_CACHE"})),
				},
				"OK when VALID SERVICECATEGORY parameter": {
					clientSession: TOSession, requestOpts: client.RequestOptions{QueryParameters: url.Values{"serviceCategory": {"serviceCategory1"}}},
					expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateExpectedFields(map[string]interface{}{"ServiceCategory": "serviceCategory1"})),
				},
				"OK when VALID TENANT parameter": {
					clientSession: TOSession, requestOpts: client.RequestOptions{QueryParameters: url.Values{"tenant": {"tenant1"}}},
					expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateExpectedFields(map[string]interface{}{"Tenant": "tenant1"})),
				},
				"OK when VALID TOPOLOGY parameter": {
					clientSession: TOSession, requestOpts: client.RequestOptions{QueryParameters: url.Values{"topology": {"mso-topology"}}},
					expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateExpectedFields(map[string]interface{}{"Topology": "mso-topology"})),
				},
				"OK when VALID TYPE parameter": {
					clientSession: TOSession, requestOpts: client.RequestOptions{QueryParameters: url.Values{"type": {"HTTP"}}},
					expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateExpectedFields(map[string]interface{}{"Type": tc.DSTypeHTTP})),
				},
				"OK when VALID XMLID parameter": {
					clientSession: TOSession, requestOpts: client.RequestOptions{QueryParameters: url.Values{"xmlId": {"ds1"}}},
					expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1),
						validateExpectedFields(map[string]interface{}{"XMLID": "ds1"})),
				},
				"EMPTY RESPONSE when INVALID ACCESSIBLETO parameter": {
					clientSession: TOSession, requestOpts: client.RequestOptions{QueryParameters: url.Values{"accessibleTo": {"10000"}}},
					expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
				"EMPTY RESPONSE when INVALID CDN parameter": {
					clientSession: TOSession, requestOpts: client.RequestOptions{QueryParameters: url.Values{"cdn": {"10000"}}},
					expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
				"EMPTY RESPONSE when INVALID PROFILE parameter": {
					clientSession: TOSession, requestOpts: client.RequestOptions{QueryParameters: url.Values{"profile": {"10000"}}},
					expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
				"EMPTY RESPONSE when INVALID TENANT parameter": {
					clientSession: TOSession, requestOpts: client.RequestOptions{QueryParameters: url.Values{"tenant": {"10000"}}},
					expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
				"EMPTY RESPONSE when INVALID TYPE parameter": {
					clientSession: TOSession, requestOpts: client.RequestOptions{QueryParameters: url.Values{"type": {"10000"}}},
					expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
				"EMPTY RESPONSE when INVALID XMLID parameter": {
					clientSession: TOSession, requestOpts: client.RequestOptions{QueryParameters: url.Values{"xmlId": {"invalid_xml_id"}}},
					expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
				"FIRST RESULT when LIMIT=1": {
					clientSession: TOSession, requestOpts: client.RequestOptions{QueryParameters: url.Values{"orderby": {"id"}, "limit": {"1"}}},
					expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validatePagination("limit")),
				},
				"SECOND RESULT when LIMIT=1 OFFSET=1": {
					clientSession: TOSession, requestOpts: client.RequestOptions{QueryParameters: url.Values{"orderby": {"id"}, "limit": {"1"}, "offset": {"1"}}},
					expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validatePagination("offset")),
				},
				"SECOND RESULT when LIMIT=1 PAGE=2": {
					clientSession: TOSession, requestOpts: client.RequestOptions{QueryParameters: url.Values{"orderby": {"id"}, "limit": {"1"}, "page": {"2"}}},
					expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validatePagination("page")),
				},
				"BAD REQUEST when INVALID LIMIT parameter": {
					clientSession: TOSession, requestOpts: client.RequestOptions{QueryParameters: url.Values{"limit": {"-2"}}},
					expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when INVALID OFFSET parameter": {
					clientSession: TOSession, requestOpts: client.RequestOptions{QueryParameters: url.Values{"limit": {"1"}, "offset": {"0"}}},
					expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when INVALID PAGE parameter": {
					clientSession: TOSession, requestOpts: client.RequestOptions{QueryParameters: url.Values{"limit": {"1"}, "page": {"0"}}},
					expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"VALID when SORTORDER param is DESC": {
					clientSession: TOSession, requestOpts: client.RequestOptions{QueryParameters: url.Values{"sortOrder": {"desc"}}},
					expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateSorted()),
				},
				"EMPTY RESPONSE when TENANT attempts reading DS OUTSIDE TENANCY": {
					clientSession: tenant4UserSession,
					requestOpts:   client.RequestOptions{QueryParameters: url.Values{"xmlId": {"ds3"}}},
					expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
			},
			"POST": {
				"CREATED when VALID request WITH GEO LIMIT COUNTRIES": {
					clientSession: TOSession,
					requestBody: generateDeliveryService(t, map[string]interface{}{
						"geoLimit":          2,
						"geoLimitCountries": []string{"US", "CA"},
						"xmlId":             "geolimit-test",
					}),
					expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusCreated), utils.ResponseHasLength(1),
						validateExpectedFields(map[string]interface{}{"GeoLimitCountries": tc.GeoLimitCountriesType{"US", "CA"}})),
				},
				"BAD REQUEST when using LONG DESCRIPTION 2 and 3 fields": {
					clientSession: TOSession,
					requestBody: generateDeliveryService(t, map[string]interface{}{
						"longDesc1": "long desc 1",
						"longDesc2": "long desc 2",
						"xmlId":     "ld1-ld2-test",
					}),
					expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when XMLID left EMPTY": {
					clientSession: TOSession,
					requestBody: generateDeliveryService(t, map[string]interface{}{
						"xmlId": "",
					}),
					expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when XMLID is NIL": {
					clientSession: TOSession,
					requestBody: generateDeliveryService(t, map[string]interface{}{
						"xmlId": nil,
					}),
					expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when TOPOLOGY DOESNT EXIST": {
					clientSession: TOSession,
					requestBody: generateDeliveryService(t, map[string]interface{}{
						"topology": "topology-doesnt-exist",
					}),
					expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when creating STEERING DS with TLS VERSIONS": {
					clientSession: TOSession,
					requestBody: generateDeliveryService(t, map[string]interface{}{
						"tlsVersions": []string{"1.1"},
						"typeId":      GetTypeId(t, "STEERING"),
						"xmlId":       "test-TLS-creation-steering",
					}),
					expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"OK when creating HTTP DS with TLS VERSIONS": {
					clientSession: TOSession,
					requestBody: generateDeliveryService(t, map[string]interface{}{
						"tlsVersions": []string{"1.1"},
						"xmlId":       "test-TLS-creation-http",
					}),
					expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusCreated), utils.ResponseHasLength(1)),
				},
				"BAD REQUEST when creating DS with TENANCY NOT THE SAME AS CURRENT TENANT": {
					clientSession: tenant4UserSession,
					requestBody: generateDeliveryService(t, map[string]interface{}{
						"tenantId": GetTenantId(t, "tenant3"),
						"xmlId":    "test-tenancy",
					}),
					expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden), utils.ResponseHasLength(0)),
				},
			},
			"PUT": {
				"BAD REQUEST when using LONG DESCRIPTION 2 and 3 fields": {
					endpointId: GetDeliveryServiceId(t, "ds1"), clientSession: TOSession,
					requestBody: generateDeliveryService(t, map[string]interface{}{
						"longDesc1": "long desc 1",
						"longDesc2": "long desc 2",
						"xmlId":     "ds1",
					}),
					expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				//"OK when VALID request": {
				//	endpointId: GetDeliveryServiceId(t, "ds1"), clientSession: TOSession,
				//	requestBody: map[string]interface{}{
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
				//	expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				//},
				//UpdateValidateORGServerCacheGroup
				//"Assign an Origin not in a Cache Group used by a Delivery Service's Topology to that Delivery Service": {},
				"BAD REQUEST when INVALID REMAP TEXT": {
					endpointId: GetDeliveryServiceId(t, "ds1"), clientSession: TOSession,
					requestBody: generateDeliveryService(t, map[string]interface{}{
						"remapText": "@plugin=tslua.so @pparam=/opt/trafficserver/etc/trafficserver/remapPlugin1.lua\nline2",
					}),
					expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when MISSING SLICE PLUGIN SIZE": {
					endpointId: GetDeliveryServiceId(t, "ds1"), clientSession: TOSession,
					requestBody: generateDeliveryService(t, map[string]interface{}{
						"rangeRequestHandling": 3,
					}),
					expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when SLICE PLUGIN SIZE SET with INVALID RANGE REQUEST SETTING": {
					endpointId: GetDeliveryServiceId(t, "ds1"), clientSession: TOSession,
					requestBody: generateDeliveryService(t, map[string]interface{}{
						"rangeRequestHandling": 1,
						"rangeSliceBlockSize":  262144,
					}),
					expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when SLICE PLUGIN SIZE TOO SMALL": {
					endpointId: GetDeliveryServiceId(t, "ds1"), clientSession: TOSession,
					requestBody: generateDeliveryService(t, map[string]interface{}{
						"rangeRequestHandling": 3,
						"rangeSliceBlockSize":  0,
					}),
					expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when SLICE PLUGIN SIZE TOO LARGE": {
					endpointId: GetDeliveryServiceId(t, "ds1"), clientSession: TOSession,
					requestBody: generateDeliveryService(t, map[string]interface{}{
						"rangeRequestHandling": 3,
						"rangeSliceBlockSize":  40000000,
					}),
					expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when ADDING TOPOLOGY to CLIENT STEERING DS": {
					endpointId: GetDeliveryServiceId(t, "ds-client-steering"), clientSession: TOSession,
					requestBody: generateDeliveryService(t, map[string]interface{}{
						"topology": "mso-topology",
						"xmlId":    "ds-client-steering",
						"typeId":   GetTypeId(t, "CLIENT_STEERING"),
					}),
					expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when TOPOLOGY DOESNT EXIST": {
					endpointId: GetDeliveryServiceId(t, "ds1"), clientSession: TOSession,
					requestBody: generateDeliveryService(t, map[string]interface{}{
						"topology": "",
						"xmlId":    "ds1",
					}),
					expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when ADDING TOPOLOGY to DS with DS REQUIRED CAPABILITY": {
					endpointId: GetDeliveryServiceId(t, "ds1"), clientSession: TOSession,
					requestBody: generateDeliveryService(t, map[string]interface{}{
						"topology": "top-for-ds-req",
						"xmlId":    "ds1",
					}),
					expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"OK when REMOVING TOPOLOGY": {
					endpointId: GetDeliveryServiceId(t, "ds-based-top-with-no-mids"), clientSession: TOSession,
					requestBody: generateDeliveryService(t, map[string]interface{}{
						"topology": nil,
						"xmlId":    "ds-based-top-with-no-mids",
					}),
					expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				//t.Fatalf("expected 400-level error assigning Topology %s to Delivery Service %s because Cache Group %s has no Servers in it in CDN %d, no error received", *dsTopology, xmlID, cacheGroupName, *ds.CDNID)
				// "top-ds-in-cdn2"
				//"BAD REQUEST when ASSIGNING TOPOLOGY when CG has NO SERVERS": {
				//	endpointId: GetDeliveryServiceId(t, "top-ds-in-cdn2"), clientSession: TOSession,
				//	requestBody: map[string]interface{}{
				//		"topology": "top-cg-no-servers",
				//		"xmlId":    "top-ds-in-cdn2",
				//	},
				//	expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				//},
				"OK when DS with TOPOLOGY updates HEADER REWRITE FIELDS": {
					endpointId: GetDeliveryServiceId(t, "ds-top"), clientSession: TOSession,
					requestBody: generateDeliveryService(t, map[string]interface{}{
						"firstHeaderRewrite": "foo",
						"innerHeaderRewrite": "bar",
						"lastHeaderRewrite":  "baz",
						"topology":           "mso-topology",
						"xmlId":              "ds-top",
					}),
					expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"BAD REQUEST when DS with NO TOPOLOGY updates HEADER REWRITE FIELDS": {
					endpointId: GetDeliveryServiceId(t, "ds1"), clientSession: TOSession,
					requestBody: generateDeliveryService(t, map[string]interface{}{
						"firstHeaderRewrite": "foo",
						"innerHeaderRewrite": "bar",
						"lastHeaderRewrite":  "baz",
					}),
					expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when DS with TOPOLOGY updates LEGACY HEADER REWRITE FIELDS": {
					endpointId: GetDeliveryServiceId(t, "ds-top"), clientSession: TOSession,
					requestBody: generateDeliveryService(t, map[string]interface{}{
						"edgeHeaderRewrite": "foo",
						"midHeaderRewrite":  "bar",
						"topology":          "mso-topology",
						"xmlId":             "ds-top",
					}),
					expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"OK when DS with NO TOPOLOGY updates LEGACY HEADER REWRITE FIELDS": {
					endpointId: GetDeliveryServiceId(t, "ds1"), clientSession: TOSession,
					requestBody: generateDeliveryService(t, map[string]interface{}{
						"profileId":         GetProfileId(t, "ATS_EDGE_TIER_CACHE"),
						"edgeHeaderRewrite": "foo",
						"midHeaderRewrite":  "bar",
						"xmlId":             "ds1",
					}),
					expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"OK when UPDATING MINOR VERSION FIELDS": {
					endpointId: GetDeliveryServiceId(t, "ds-test-minor-versions"), clientSession: TOSession,
					requestBody: generateDeliveryService(t, map[string]interface{}{
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
					expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateExpectedFields(map[string]interface{}{"ConsistentHashQueryParams": []string{"d", "e", "f"},
							"ConsistentHashRegex": "foo", "DeepCachingType": tc.DeepCachingTypeNever, "FQPacingRate": 41, "MaxOriginConnections": 500,
							"SigningAlgorithm": "uri_signing", "Tenant": "tenant1", "TRRequestHeaders": "X-ooF\nX-raB",
							"TRResponseHeaders": "Access-Control-Max-Age: 600\nContent-Type: text/html; charset=utf-8",
						})),
				},
				"BAD REQUEST when INVALID COUNTRY CODE": {
					endpointId: GetDeliveryServiceId(t, "ds1"), clientSession: TOSession,
					requestBody: generateDeliveryService(t, map[string]interface{}{
						"geoLimit":          2,
						"geoLimitCountries": []string{"US", "CA", "12"},
						"xmlId":             "invalid-geolimit-test",
					}),
					expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when CHANGING TOPOLOGY of DS with ORG SERVERS ASSIGNED": {
					endpointId: GetDeliveryServiceId(t, "ds-top"), clientSession: TOSession,
					requestBody: generateDeliveryService(t, map[string]interface{}{
						"topology": "another-topology",
						"xmlId":    "ds-top",
					}),
					expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when UPDATING DS OUTSIDE TENANCY": {
					endpointId: GetDeliveryServiceId(t, "ds3"), clientSession: tenant4UserSession,
					requestBody:  generateDeliveryService(t, map[string]interface{}{"xmlId": "ds3"}),
					expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
				"PRECONDITION FAILED when updating with IMS & IUS Headers": {
					endpointId: GetDeliveryServiceId(t, "ds1"), clientSession: TOSession,
					requestOpts: client.RequestOptions{
						Header: http.Header{
							rfc.IfModifiedSince: {currentTimeRFC}, rfc.IfUnmodifiedSince: {currentTimeRFC},
						},
					},
					requestBody:  generateDeliveryService(t, map[string]interface{}{"xmlId": "ds1"}),
					expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
				"PRECONDITION FAILED when updating with IFMATCH ETAG Header": {
					endpointId: GetDeliveryServiceId(t, "ds1"), clientSession: TOSession,
					requestBody:  generateDeliveryService(t, map[string]interface{}{"xmlId": "ds1"}),
					requestOpts:  client.RequestOptions{Header: http.Header{rfc.IfMatch: {rfc.ETag(currentTime)}}},
					expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
			},
			"DELETE": {
				"BAD REQUEST when DELETING DS OUTSIDE TENANCY": {
					endpointId: GetDeliveryServiceId(t, "ds3"), clientSession: tenant4UserSession,
					expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
			},
			"GET AFTER CHANGES": {
				"OK when CHANGES made": {
					clientSession: TOSession,
					requestOpts: client.RequestOptions{
						Header: http.Header{
							rfc.IfModifiedSince: {currentTimeRFC}, rfc.IfUnmodifiedSince: {currentTimeRFC},
						},
					},
					expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
			},
			//"CDN LOCK": {
			//	//CUDDeliveryServiceWithLocks
			//	"Create/ Update/ Delete delivery services with locks": {},
			//},
			//"DELIVERY SERVICES CAPACITY": {
			//	// capDS, _, err := TOSession.GetDeliveryServiceCapacity(*ds.ID, client.RequestOptions{})
			//	"Basic GET request for /deliveryservices/{{ID}}/capacity": {},
			//},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					ds := tc.DeliveryServiceV4{}

					if val, ok := testCase.requestOpts.QueryParameters["accessibleTo"]; ok {
						if _, err := strconv.Atoi(val[0]); err != nil {
							testCase.requestOpts.QueryParameters.Set("accessibleTo", strconv.Itoa(GetTenantId(t, val[0])))
						}
					}
					if val, ok := testCase.requestOpts.QueryParameters["cdn"]; ok {
						if _, err := strconv.Atoi(val[0]); err != nil {
							testCase.requestOpts.QueryParameters.Set("cdn", strconv.Itoa(GetCDNId(t, val[0])))
						}
					}
					if val, ok := testCase.requestOpts.QueryParameters["profile"]; ok {
						if _, err := strconv.Atoi(val[0]); err != nil {
							testCase.requestOpts.QueryParameters.Set("profile", strconv.Itoa(GetProfileId(t, val[0])))
						}
					}
					if val, ok := testCase.requestOpts.QueryParameters["type"]; ok {
						if _, err := strconv.Atoi(val[0]); err != nil {
							testCase.requestOpts.QueryParameters.Set("type", strconv.Itoa(GetTypeId(t, val[0])))
						}
					}
					if val, ok := testCase.requestOpts.QueryParameters["tenant"]; ok {
						if _, err := strconv.Atoi(val[0]); err != nil {
							testCase.requestOpts.QueryParameters.Set("tenant", strconv.Itoa(GetTenantId(t, val[0])))
						}
					}

					if testCase.requestBody != nil {
						dat, err := json.Marshal(testCase.requestBody)
						assert.NoError(t, err, "Error occurred when marshalling request body: %v", err)
						err = json.Unmarshal(dat, &ds)
						assert.NoError(t, err, "Error occurred when unmarshalling request body: %v", err)
					}

					switch method {
					case "GET", "GET AFTER CHANGES":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.clientSession.GetDeliveryServices(testCase.requestOpts)
							for _, check := range testCase.expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							fmt.Println(*ds.Tenant)
							resp, reqInf, err := testCase.clientSession.CreateDeliveryService(ds, testCase.requestOpts)
							fmt.Println(resp)
							for _, check := range testCase.expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "PUT":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.clientSession.UpdateDeliveryService(testCase.endpointId(), ds, testCase.requestOpts)
							for _, check := range testCase.expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.clientSession.DeleteDeliveryService(testCase.endpointId(), testCase.requestOpts)
							for _, check := range testCase.expectations {
								check(t, reqInf, nil, resp.Alerts, err)
							}
						})
					case "CDNLOCK":
						t.Run(name, func(t *testing.T) {
							UpdateCachegroupWithLocks(t)
						})
					}
				}
			})
		}
	})
}

func validateExpectedFields(expectedResp map[string]interface{}) utils.CkReqFunc {
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

func validateSorted() utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, alerts tc.Alerts, _ error) {
		dsResp := resp.([]tc.DeliveryServiceV40)
		var sortedList []string
		assert.RequireGreaterOrEqual(t, len(dsResp), 2, "Need at least 2 ASNs in Traffic Ops to test sorted, found: %d", len(dsResp))

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

// createSession creates a session using the passed in username and password.
func createSession(t *testing.T, username string, password string) *client.Session {
	userSession, _, err := client.LoginWithAgent(Config.TrafficOps.URL, username, password, true, "to-api-v4-client-tests", false, toReqTimeout)
	assert.RequireNoError(t, err, "Could not login with user %v: %v", username, err)
	return userSession
}

func CUDDeliveryServiceWithLocks(t *testing.T) {
	// Create a new user with operations level privileges
	user1 := tc.UserV4{
		Username:             "lock_user1",
		RegistrationSent:     new(time.Time),
		LocalPassword:        util.StrPtr("test_pa$$word"),
		ConfirmLocalPassword: util.StrPtr("test_pa$$word"),
		Role:                 "operations",
	}
	user1.Email = util.StrPtr("lockuseremail@domain.com")
	user1.TenantID = 1
	//util.IntPtr(resp.Response[0].ID)
	user1.FullName = util.StrPtr("firstName LastName")
	_, _, err := TOSession.CreateUser(user1, client.RequestOptions{})
	if err != nil {
		t.Fatalf("could not create test user with username: %s", user1.Username)
	}
	defer ForceDeleteTestUsersByUsernames(t, []string{"lock_user1"})

	// Establish a session with the newly created non admin level user
	userSession, _, err := client.LoginWithAgent(Config.TrafficOps.URL, user1.Username, *user1.LocalPassword, true, "to-api-v4-client-tests", false, toReqTimeout)
	if err != nil {
		t.Fatalf("could not login with user lock_user1: %v", err)
	}
	if len(testData.DeliveryServices) == 0 {
		t.Fatalf("no deliveryservices to run the test on, quitting")
	}

	cdn := createBlankCDN("sslkeytransfer", t)
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("name", "HTTP")
	types, _, err := TOSession.GetTypes(opts)
	if err != nil {
		t.Fatalf("unable to get Types: %v - alerts: %+v", err, types.Alerts)
	}
	if len(types.Response) < 1 {
		t.Fatal("expected at least one type")
	}
	customDS := getCustomDS(cdn.ID, types.Response[0].ID, "cdn-locks-test-ds-name", "edge", "https://test-cdn-locks.com", "cdn-locks-test-ds-xml-id")

	// Create a lock for this user
	_, _, err = userSession.CreateCDNLock(tc.CDNLock{
		CDN:     cdn.Name,
		Message: util.StrPtr("test lock"),
		Soft:    util.BoolPtr(false),
	}, client.RequestOptions{})
	if err != nil {
		t.Fatalf("couldn't create cdn lock: %v", err)
	}
	// Try to create a new ds on a CDN that another user has a hard lock on -> this should fail
	_, reqInf, err := TOSession.CreateDeliveryService(customDS, client.RequestOptions{})
	if err == nil {
		t.Error("expected an error while creating a new ds for a CDN for which a hard lock is held by another user, but got nothing")
	}
	if reqInf.StatusCode != http.StatusForbidden {
		t.Errorf("expected a 403 forbidden status while creating a new ds for a CDN for which a hard lock is held by another user, but got %d", reqInf.StatusCode)
	}

	// Try to create a new ds on a CDN that the same user has a hard lock on -> this should succeed
	dsResp, reqInf, err := userSession.CreateDeliveryService(customDS, client.RequestOptions{})
	if err != nil {
		t.Errorf("expected no error while creating a new ds for a CDN for which a hard lock is held by the same user, but got %v", err)
	}
	if len(dsResp.Response) != 1 {
		t.Fatalf("one response expected, but got %d", len(dsResp.Response))
	}
	opts = client.NewRequestOptions()
	opts.QueryParameters.Set("xmlId", *customDS.XMLID)
	deliveryServices, _, err := userSession.GetDeliveryServices(opts)
	if err != nil {
		t.Fatalf("couldn't get ds: %v", err)
	}
	if len(deliveryServices.Response) != 1 {
		t.Fatal("couldn't get exactly one ds in the response, quitting")
	}
	dsID := dsResp.Response[0].ID
	// Try to update a ds on a CDN that another user has a hard lock on -> this should fail
	customDS.LongDesc = util.StrPtr("changed_long_desc")
	_, reqInf, err = TOSession.UpdateDeliveryService(*dsID, customDS, client.RequestOptions{})
	if err == nil {
		t.Error("expected an error while updating a ds for a CDN for which a hard lock is held by another user, but got nothing")
	}
	if reqInf.StatusCode != http.StatusForbidden {
		t.Errorf("expected a 403 forbidden status while updating a ds for a CDN for which a hard lock is held by another user, but got %d", reqInf.StatusCode)
	}
	// Try to update a ds on a CDN that the same user has a hard lock on -> this should succeed
	_, reqInf, err = userSession.UpdateDeliveryService(*dsID, customDS, client.RequestOptions{})
	if err != nil {
		t.Errorf("expected no error while updating a ds for a CDN for which a hard lock is held by the same user, but got %v", err)
	}
	// Try to delete a ds on a CDN that another user has a hard lock on -> this should fail
	_, reqInf, err = TOSession.DeleteDeliveryService(*dsID, client.RequestOptions{})
	if err == nil {
		t.Error("expected an error while deleting a ds for a CDN for which a hard lock is held by another user, but got nothing")
	}
	if reqInf.StatusCode != http.StatusForbidden {
		t.Errorf("expected a 403 forbidden status while deleting a ds for a CDN for which a hard lock is held by another user, but got %d", reqInf.StatusCode)
	}

	// Try to delete a ds on a CDN that the same user has a hard lock on -> this should succeed
	_, reqInf, err = userSession.DeleteDeliveryService(*dsID, client.RequestOptions{})
	if err != nil {
		t.Errorf("expected no error while deleting a ds for a CDN for which a hard lock is held by the same user, but got %v", err)
	}

	// Delete the lock
	_, _, err = userSession.DeleteCDNLocks(client.RequestOptions{QueryParameters: url.Values{"cdn": []string{cdn.Name}}})
	if err != nil {
		t.Errorf("expected no error while deleting other user's lock using admin endpoint, but got %v", err)
	}
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
