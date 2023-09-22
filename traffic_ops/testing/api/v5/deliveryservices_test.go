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

	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
	client "github.com/apache/trafficcontrol/v8/traffic_ops/v5-client"
)

func TestDeliveryServices(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Users, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, ServerCapabilities, ServerServerCapabilities, ServiceCategories, DeliveryServices, DeliveryServiceServerAssignments}, func() {

		currentTime := time.Now().UTC().Add(-15 * time.Second)
		currentTimeRFC := currentTime.Format(time.RFC1123)
		tomorrow := currentTime.AddDate(0, 0, 1).Format(time.RFC1123)

		tenant1UserSession := utils.CreateV5Session(t, Config.TrafficOps.URL, "tenant1user", "pa$$word", Config.Default.Session.TimeoutInSecs)
		tenant2UserSession := utils.CreateV5Session(t, Config.TrafficOps.URL, "tenant2user", "pa$$word", Config.Default.Session.TimeoutInSecs)
		tenant3UserSession := utils.CreateV5Session(t, Config.TrafficOps.URL, "tenant3user", "pa$$word", Config.Default.Session.TimeoutInSecs)
		tenant4UserSession := utils.CreateV5Session(t, Config.TrafficOps.URL, "tenant4user", "pa$$word", Config.Default.Session.TimeoutInSecs)

		methodTests := utils.V5TestCase{
			"GET": {
				"NOT MODIFIED when NO CHANGES made": {
					ClientSession: TOSession, RequestOpts: client.RequestOptions{Header: http.Header{rfc.IfModifiedSince: {tomorrow}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"OK when VALID request": {
					ClientSession: TOSession, Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						utils.ResponseLengthGreaterOrEqual(1)),
				},
				"OK when VALID ACCESSIBLETO parameter": {
					ClientSession: TOSession, RequestOpts: client.RequestOptions{QueryParameters: url.Values{"accessibleTo": {"1"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1)),
				},
				"OK when Active=ACTIVE": {
					ClientSession: TOSession, RequestOpts: client.RequestOptions{QueryParameters: url.Values{"active": {string(tc.DSActiveStateActive)}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateDSExpectedFields(map[string]interface{}{"Active": tc.DSActiveStateActive}, true)),
				},
				"OK when Active=PRIMED": {
					ClientSession: TOSession, RequestOpts: client.RequestOptions{QueryParameters: url.Values{"active": {string(tc.DSActiveStatePrimed)}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateDSExpectedFields(map[string]interface{}{"Active": tc.DSActiveStatePrimed}, true)),
				},
				"OK when Active=INACTIVE": {
					ClientSession: TOSession, RequestOpts: client.RequestOptions{QueryParameters: url.Values{"active": {string(tc.DSActiveStateInactive)}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateDSExpectedFields(map[string]interface{}{"Active": tc.DSActiveStateInactive}, true)),
				},
				"OK when VALID CDN parameter": {
					ClientSession: TOSession, RequestOpts: client.RequestOptions{QueryParameters: url.Values{"cdn": {"cdn1"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateDSExpectedFields(map[string]interface{}{"CDNName": "cdn1"}, true)),
				},
				"OK when VALID LOGSENABLED parameter": {
					ClientSession: TOSession, RequestOpts: client.RequestOptions{QueryParameters: url.Values{"logsEnabled": {"false"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateDSExpectedFields(map[string]interface{}{"LogsEnabled": false}, true)),
				},
				"OK when VALID PROFILE parameter": {
					ClientSession: TOSession, RequestOpts: client.RequestOptions{QueryParameters: url.Values{"profile": {"ATS_EDGE_TIER_CACHE"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateDSExpectedFields(map[string]interface{}{"ProfileName": "ATS_EDGE_TIER_CACHE"}, true)),
				},
				"OK when VALID SERVICECATEGORY parameter": {
					ClientSession: TOSession, RequestOpts: client.RequestOptions{QueryParameters: url.Values{"serviceCategory": {"serviceCategory1"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateDSExpectedFields(map[string]interface{}{"ServiceCategory": "serviceCategory1"}, true)),
				},
				"OK when VALID TENANT parameter": {
					ClientSession: TOSession, RequestOpts: client.RequestOptions{QueryParameters: url.Values{"tenant": {"tenant1"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateDSExpectedFields(map[string]interface{}{"Tenant": "tenant1"}, true)),
				},
				"OK when VALID TOPOLOGY parameter": {
					ClientSession: TOSession, RequestOpts: client.RequestOptions{QueryParameters: url.Values{"topology": {"mso-topology"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateDSExpectedFields(map[string]interface{}{"Topology": "mso-topology"}, true)),
				},
				"OK when VALID TYPE parameter": {
					ClientSession: TOSession, RequestOpts: client.RequestOptions{QueryParameters: url.Values{"type": {"HTTP"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateDSExpectedFields(map[string]interface{}{"Type": string(tc.DSTypeHTTP)}, true)),
				},
				"OK when VALID XMLID parameter": {
					ClientSession: TOSession, RequestOpts: client.RequestOptions{QueryParameters: url.Values{"xmlId": {"ds1"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1),
						validateDSExpectedFields(map[string]interface{}{"XMLID": "ds1"}, true)),
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
				"OK when PARENT TENANT reads DS of INACTIVE CHILD TENANT": {
					ClientSession: tenant1UserSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"xmlId": {"ds2"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1)),
				},
				"EMPTY RESPONSE when DS BELONGS to TENANT but PARENT TENANT is INACTIVE": {
					ClientSession: tenant3UserSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"xmlId": {"ds3"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
				"EMPTY RESPONSE when INACTIVE TENANT reads DS of SAME TENANCY": {
					ClientSession: tenant2UserSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"xmlId": {"ds2"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
				"EMPTY RESPONSE when TENANT reads DS OUTSIDE TENANCY": {
					ClientSession: tenant4UserSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"xmlId": {"ds3"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
				"EMPTY RESPONSE when CHILD TENANT reads DS of PARENT TENANT": {
					ClientSession: tenant3UserSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"xmlId": {"ds2"}}},
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
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusCreated),
						validateDSExpectedFields(map[string]interface{}{"GeoLimitCountries": []string{"US", "CA"}}, false)),
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
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusCreated)),
				},
				"BAD REQUEST when creating DS with TENANCY NOT THE SAME AS CURRENT TENANT": {
					ClientSession: tenant4UserSession,
					RequestBody: generateDeliveryService(t, map[string]interface{}{
						"tenantId": GetTenantID(t, "tenant3")(),
						"xmlId":    "test-tenancy",
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
				"BAD REQUEST when creating DS with INVALID ACTIVE STATE": {
					ClientSession: TOSession,
					RequestBody: generateDeliveryService(
						t,
						map[string]interface{}{
							"active": "this is a totally invalid active value",
						},
					),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
			},
			"PUT": {
				"BAD REQUEST when INVALID ACTIVE STATE": {
					ClientSession: TOSession,
					EndpointID:    GetDeliveryServiceId(t, "ds1"),
					RequestBody: generateDeliveryService(
						t,
						map[string]interface{}{
							"active": "this is a totally invalid active value",
						},
					),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"OK when VALID request": {
					EndpointID: GetDeliveryServiceId(t, "ds2"), ClientSession: TOSession,
					RequestBody: generateDeliveryService(t, map[string]interface{}{
						"maxRequestHeaderBytes": 131080,
						"longDesc":              "something different",
						"maxDNSAnswers":         164598,
						"maxOriginConnections":  100,
						"active":                tc.DSActiveStatePrimed,
						"displayName":           "newds2displayname",
						"dscp":                  41,
						"geoLimit":              1,
						"initialDispersion":     2,
						"ipv6RoutingEnabled":    false,
						"logsEnabled":           false,
						"missLat":               42.881944,
						"missLong":              -88.627778,
						"multiSiteOrigin":       true,
						"orgServerFqdn":         "http://origin.example.net",
						"protocol":              2,
						"regional":              true,
						"routingName":           "ccr-ds2",
						"qStringIgnore":         0,
						"regionalGeoBlocking":   true,
						"xmlId":                 "ds2",
					}),
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateDSExpectedFields(map[string]interface{}{"MaxRequestHeaderSize": 131080,
							"LongDesc": "something different", "MaxDNSAnswers": 164598, "MaxOriginConnections": 100,
							"Active": tc.DSActiveStatePrimed, "DisplayName": "newds2displayname", "DSCP": 41, "GeoLimit": 1,
							"InitialDispersion": 2, "IPV6RoutingEnabled": false, "LogsEnabled": false, "MissLat": 42.881944,
							"MissLong": -88.627778, "MultiSiteOrigin": true, "OrgServerFQDN": "http://origin.example.net",
							"Protocol": 2, "QStringIgnore": 0, "RegionalGeoBlocking": true,
						}, false)),
				},
				"BAD REQUEST when INVALID REMAP TEXT": {
					EndpointID: GetDeliveryServiceId(t, "ds1"), ClientSession: TOSession,
					RequestBody: generateDeliveryService(t, map[string]interface{}{
						"remapText": "@plugin=tslua.so @pparam=/opt/trafficserver/etc/trafficserver/remapPlugin1.lua\nline2",
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when MISSING SLICE PLUGIN SIZE": {
					EndpointID: GetDeliveryServiceId(t, "ds1"), ClientSession: TOSession,
					RequestBody: generateDeliveryService(t, map[string]interface{}{
						"rangeRequestHandling": 3,
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when SLICE PLUGIN SIZE SET with INVALID RANGE REQUEST SETTING": {
					EndpointID: GetDeliveryServiceId(t, "ds1"), ClientSession: TOSession,
					RequestBody: generateDeliveryService(t, map[string]interface{}{
						"rangeRequestHandling": 1,
						"rangeSliceBlockSize":  262144,
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when SLICE PLUGIN SIZE TOO SMALL": {
					EndpointID: GetDeliveryServiceId(t, "ds1"), ClientSession: TOSession,
					RequestBody: generateDeliveryService(t, map[string]interface{}{
						"rangeRequestHandling": 3,
						"rangeSliceBlockSize":  0,
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when SLICE PLUGIN SIZE TOO LARGE": {
					EndpointID: GetDeliveryServiceId(t, "ds1"), ClientSession: TOSession,
					RequestBody: generateDeliveryService(t, map[string]interface{}{
						"rangeRequestHandling": 3,
						"rangeSliceBlockSize":  40000000,
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when ADDING TOPOLOGY to CLIENT STEERING DS": {
					EndpointID: GetDeliveryServiceId(t, "ds-client-steering"), ClientSession: TOSession,
					RequestBody: generateDeliveryService(t, map[string]interface{}{
						"topology": "mso-topology",
						"xmlId":    "ds-client-steering",
						"typeId":   GetTypeId(t, "CLIENT_STEERING"),
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when TOPOLOGY DOESNT EXIST": {
					EndpointID: GetDeliveryServiceId(t, "ds1"), ClientSession: TOSession,
					RequestBody: generateDeliveryService(t, map[string]interface{}{
						"topology": "",
						"xmlId":    "ds1",
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when ADDING TOPOLOGY to DS with DS REQUIRED CAPABILITY": {
					EndpointID: GetDeliveryServiceId(t, "ds1"), ClientSession: TOSession,
					RequestBody: generateDeliveryService(t, map[string]interface{}{
						"topology":             "top-for-ds-req",
						"xmlId":                "ds1",
						"requiredCapabilities": []string{"foo"},
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when ADDING TOPOLOGY to DS when NO CACHES in SAME CDN as DS": {
					EndpointID: GetDeliveryServiceId(t, "top-ds-in-cdn2"), ClientSession: TOSession,
					RequestBody: generateDeliveryService(t, map[string]interface{}{
						"cdnId":    GetCDNID(t, "cdn2")(),
						"topology": "top-with-caches-in-cdn1",
						"xmlId":    "top-ds-in-cdn2",
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"OK when REMOVING TOPOLOGY": {
					EndpointID: GetDeliveryServiceId(t, "ds-based-top-with-no-mids"), ClientSession: TOSession,
					RequestBody: generateDeliveryService(t, map[string]interface{}{
						"topology": nil,
						"xmlId":    "ds-based-top-with-no-mids",
					}),
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"OK when DS with TOPOLOGY updates HEADER REWRITE FIELDS": {
					EndpointID: GetDeliveryServiceId(t, "ds-top"), ClientSession: TOSession,
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
					EndpointID: GetDeliveryServiceId(t, "ds1"), ClientSession: TOSession,
					RequestBody: generateDeliveryService(t, map[string]interface{}{
						"firstHeaderRewrite": "foo",
						"innerHeaderRewrite": "bar",
						"lastHeaderRewrite":  "baz",
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when DS with TOPOLOGY updates LEGACY HEADER REWRITE FIELDS": {
					EndpointID: GetDeliveryServiceId(t, "ds-top"), ClientSession: TOSession,
					RequestBody: generateDeliveryService(t, map[string]interface{}{
						"edgeHeaderRewrite": "foo",
						"midHeaderRewrite":  "bar",
						"topology":          "mso-topology",
						"xmlId":             "ds-top",
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"OK when DS with NO TOPOLOGY updates LEGACY HEADER REWRITE FIELDS": {
					EndpointID: GetDeliveryServiceId(t, "ds2"), ClientSession: TOSession,
					RequestBody: generateDeliveryService(t, map[string]interface{}{
						"profileId":         GetProfileID(t, "ATS_EDGE_TIER_CACHE")(),
						"edgeHeaderRewrite": "foo",
						"midHeaderRewrite":  "bar",
						"routingName":       "ccr-ds2",
						"xmlId":             "ds2",
					}),
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"OK when UPDATING MINOR VERSION FIELDS": {
					EndpointID: GetDeliveryServiceId(t, "ds-test-minor-versions"), ClientSession: TOSession,
					RequestBody: generateDeliveryService(t, map[string]interface{}{
						"consistentHashQueryParams": []string{"d", "e", "f"},
						"consistentHashRegex":       "foo",
						"deepCachingType":           "NEVER",
						"fqPacingRate":              41,
						"maxOriginConnections":      500,
						"routingName":               "cdn",
						"signingAlgorithm":          "uri_signing",
						"tenantId":                  GetTenantID(t, "tenant1")(),
						"trRequestHeaders":          "X-ooF\nX-raB",
						"trResponseHeaders":         "Access-Control-Max-Age: 600\nContent-Type: text/html; charset=utf-8",
						"xmlId":                     "ds-test-minor-versions",
					}),
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateDSExpectedFields(map[string]interface{}{"ConsistentHashQueryParams": []string{"d", "e", "f"},
							"ConsistentHashRegex": "foo", "DeepCachingType": tc.DeepCachingTypeNever, "FQPacingRate": 41, "MaxOriginConnections": 500,
							"SigningAlgorithm": "uri_signing", "Tenant": "tenant1", "TRRequestHeaders": "X-ooF\nX-raB",
							"TRResponseHeaders": "Access-Control-Max-Age: 600\nContent-Type: text/html; charset=utf-8",
						}, false)),
				},
				"BAD REQUEST when INVALID COUNTRY CODE": {
					EndpointID: GetDeliveryServiceId(t, "ds1"), ClientSession: TOSession,
					RequestBody: generateDeliveryService(t, map[string]interface{}{
						"geoLimit":          2,
						"geoLimitCountries": []string{"US", "CA", "12"},
						"xmlId":             "invalid-geolimit-test",
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when CHANGING TOPOLOGY of DS with ORG SERVERS ASSIGNED": {
					EndpointID: GetDeliveryServiceId(t, "ds-top"), ClientSession: TOSession,
					RequestBody: generateDeliveryService(t, map[string]interface{}{
						"topology": "another-topology",
						"xmlId":    "ds-top",
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when UPDATING DS OUTSIDE TENANCY": {
					EndpointID: GetDeliveryServiceId(t, "ds3"), ClientSession: tenant4UserSession,
					RequestBody:  generateDeliveryService(t, map[string]interface{}{"xmlId": "ds3"}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
				"PRECONDITION FAILED when updating with IMS & IUS Headers": {
					EndpointID: GetDeliveryServiceId(t, "ds1"), ClientSession: TOSession,
					RequestOpts:  client.RequestOptions{Header: http.Header{rfc.IfUnmodifiedSince: {currentTimeRFC}}},
					RequestBody:  generateDeliveryService(t, map[string]interface{}{"xmlId": "ds1"}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
				"PRECONDITION FAILED when updating with IFMATCH ETAG Header": {
					EndpointID: GetDeliveryServiceId(t, "ds1"), ClientSession: TOSession,
					RequestBody:  generateDeliveryService(t, map[string]interface{}{"xmlId": "ds1"}),
					RequestOpts:  client.RequestOptions{Header: http.Header{rfc.IfMatch: {rfc.ETag(currentTime)}}},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
			},
			"DELETE": {
				"BAD REQUEST when DELETING DS OUTSIDE TENANCY": {
					EndpointID: GetDeliveryServiceId(t, "ds3"), ClientSession: tenant4UserSession,
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
			},
			"GET AFTER CHANGES": {
				"OK when CHANGES made": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfModifiedSince: {currentTimeRFC}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
			},
			"DELIVERY SERVICES CAPACITY": {
				"OK when VALID request": {
					EndpointID: GetDeliveryServiceId(t, "ds1"), ClientSession: TOSession,
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					var ds tc.DeliveryServiceV5

					if val, ok := testCase.RequestOpts.QueryParameters["accessibleTo"]; ok {
						if _, err := strconv.Atoi(val[0]); err != nil {
							testCase.RequestOpts.QueryParameters.Set("accessibleTo", strconv.Itoa(GetTenantID(t, val[0])()))
						}
					}
					if val, ok := testCase.RequestOpts.QueryParameters["cdn"]; ok {
						if _, err := strconv.Atoi(val[0]); err != nil {
							testCase.RequestOpts.QueryParameters.Set("cdn", strconv.Itoa(GetCDNID(t, val[0])()))
						}
					}
					if val, ok := testCase.RequestOpts.QueryParameters["profile"]; ok {
						if _, err := strconv.Atoi(val[0]); err != nil {
							testCase.RequestOpts.QueryParameters.Set("profile", strconv.Itoa(GetProfileID(t, val[0])()))
						}
					}
					if val, ok := testCase.RequestOpts.QueryParameters["type"]; ok {
						if _, err := strconv.Atoi(val[0]); err != nil {
							testCase.RequestOpts.QueryParameters.Set("type", strconv.Itoa(GetTypeId(t, val[0])))
						}
					}
					if val, ok := testCase.RequestOpts.QueryParameters["tenant"]; ok {
						if _, err := strconv.Atoi(val[0]); err != nil {
							testCase.RequestOpts.QueryParameters.Set("tenant", strconv.Itoa(GetTenantID(t, val[0])()))
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
							resp, reqInf, err := testCase.ClientSession.CreateDeliveryService(ds, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "PUT":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.UpdateDeliveryService(testCase.EndpointID(), ds, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.DeleteDeliveryService(testCase.EndpointID(), testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, resp.Alerts, err)
							}
						})
					case "DELIVERY SERVICES CAPACITY":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.GetDeliveryServiceCapacity(testCase.EndpointID(), testCase.RequestOpts)
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

func validateDSExpectedFields(expectedResp map[string]interface{}, multi bool) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		var dsResp []tc.DeliveryServiceV5
		if multi {
			dsResp = resp.([]tc.DeliveryServiceV5)
		} else {
			dsResp = []tc.DeliveryServiceV5{resp.(tc.DeliveryServiceV5)}
		}
		for field, expected := range expectedResp {
			for _, ds := range dsResp {
				switch field {
				case "Active":
					assert.Equal(t, expected, ds.Active, "Expected Active to be %v, but got %v", expected, ds.Active)
				case "DeepCachingType":
					assert.Equal(t, expected, ds.DeepCachingType, "Expected DeepCachingType to be %v, but got %v", expected, ds.DeepCachingType)
				case "CDNName":
					assert.Equal(t, expected, *ds.CDNName, "Expected CDNName to be %v, but got %v", expected, *ds.CDNName)
				case "ConsistentHashRegex":
					assert.Equal(t, expected, *ds.ConsistentHashRegex, "Expected ConsistentHashRegex to be %v, but got %v", expected, *ds.ConsistentHashRegex)
				case "ConsistentHashQueryParams":
					assert.Exactly(t, expected, ds.ConsistentHashQueryParams, "Expected ConsistentHashQueryParams to be %v, but got %v", expected, ds.ConsistentHashQueryParams)
				case "DisplayName":
					assert.Equal(t, expected, ds.DisplayName, "Expected DisplayName to be %v, but got %v", expected, ds.DisplayName)
				case "DSCP":
					assert.Equal(t, expected, ds.DSCP, "Expected DSCP to be %v, but got %v", expected, ds.DSCP)
				case "FQPacingRate":
					assert.Equal(t, expected, *ds.FQPacingRate, "Expected FQPacingRate to be %v, but got %v", expected, *ds.FQPacingRate)
				case "GeoLimit":
					assert.Equal(t, expected, ds.GeoLimit, "Expected GeoLimit to be %v, but got &v", expected, ds.GeoLimit)
				case "GeoLimitCountries":
					assert.Exactly(t, expected, ds.GeoLimitCountries, "Expected GeoLimitCountries to be %v, but got &v", expected, ds.GeoLimitCountries)
				case "InitialDispersion":
					assert.Equal(t, expected, *ds.InitialDispersion, "Expected InitialDispersion to be %v, but got &v", expected, ds.InitialDispersion)
				case "IPV6RoutingEnabled":
					assert.Equal(t, expected, *ds.IPV6RoutingEnabled, "Expected IPV6RoutingEnabled to be %v, but got &v", expected, ds.IPV6RoutingEnabled)
				case "LogsEnabled":
					assert.Equal(t, expected, ds.LogsEnabled, "Expected LogsEnabled to be %v, but got %v", expected, ds.LogsEnabled)
				case "LongDesc":
					assert.Equal(t, expected, ds.LongDesc, "Expected LongDesc to be %v, but got %v", expected, ds.LongDesc)
				case "MaxDNSAnswers":
					assert.Equal(t, expected, *ds.MaxDNSAnswers, "Expected MaxDNSAnswers to be %v, but got %v", expected, *ds.MaxDNSAnswers)
				case "MaxOriginConnections":
					assert.Equal(t, expected, *ds.MaxOriginConnections, "Expected MaxOriginConnections to be %v, but got %v", expected, *ds.MaxOriginConnections)
				case "MaxRequestHeaderSize":
					assert.Equal(t, expected, *ds.MaxRequestHeaderBytes, "Expected MaxRequestHeaderBytes to be %v, but got %v", expected, *ds.MaxRequestHeaderBytes)
				case "MissLat":
					assert.Equal(t, expected, *ds.MissLat, "Expected MissLat to be %v, but got %v", expected, *ds.MissLat)
				case "MissLong":
					assert.Equal(t, expected, *ds.MissLong, "Expected MissLong to be %v, but got %v", expected, *ds.MissLong)
				case "MultiSiteOrigin":
					assert.Equal(t, expected, ds.MultiSiteOrigin, "Expected MultiSiteOrigin to be %v, but got %v", expected, ds.MultiSiteOrigin)
				case "OrgServerFQDN":
					assert.Equal(t, expected, *ds.OrgServerFQDN, "Expected OrgServerFQDN to be %v, but got %v", expected, *ds.OrgServerFQDN)
				case "ProfileName":
					assert.Equal(t, expected, *ds.ProfileName, "Expected ProfileName to be %v, but got %v", expected, *ds.ProfileName)
				case "Protocol":
					assert.Equal(t, expected, *ds.Protocol, "Expected Protocol to be %v, but got %v", expected, *ds.Protocol)
				case "QStringIgnore":
					assert.Equal(t, expected, *ds.QStringIgnore, "Expected QStringIgnore to be %v, but got %v", expected, *ds.QStringIgnore)
				case "RegionalGeoBlocking":
					assert.Equal(t, expected, ds.RegionalGeoBlocking, "Expected QStringIgnore to be %v, but got %v", expected, ds.RegionalGeoBlocking)
				case "ServiceCategory":
					assert.Equal(t, expected, *ds.ServiceCategory, "Expected ServiceCategory to be %v, but got %v", expected, *ds.ServiceCategory)
				case "SigningAlgorithm":
					assert.Equal(t, expected, *ds.SigningAlgorithm, "Expected SigningAlgorithm to be %v, but got %v", expected, *ds.SigningAlgorithm)
				case "Tenant":
					assert.Equal(t, expected, *ds.Tenant, "Expected Tenant to be %v, but got %v", expected, *ds.Tenant)
				case "Topology":
					assert.Equal(t, expected, *ds.Topology, "Expected Topology to be %v, but got %v", expected, *ds.Topology)
				case "TRRequestHeaders":
					assert.Equal(t, expected, *ds.TRRequestHeaders, "Expected TRRequestHeaders to be %v, but got %v", expected, *ds.TRRequestHeaders)
				case "TRResponseHeaders":
					assert.Equal(t, expected, *ds.TRResponseHeaders, "Expected TRResponseHeaders to be %v, but got %v", expected, *ds.TRResponseHeaders)
				case "Type":
					assert.Equal(t, expected, *ds.Type, "Expected Type to be %v, but got %v", expected, *ds.Type)
				case "XMLID":
					assert.Equal(t, expected, ds.XMLID, "Expected XMLID to be %v, but got %v", expected, ds.XMLID)
				default:
					t.Errorf("Expected field: %v, does not exist in response", field)
				}
			}
		}
	}
}

func validatePagination(paginationParam string) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		paginationResp := resp.([]tc.DeliveryServiceV5)

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
		dsDescResp := resp.([]tc.DeliveryServiceV5)
		var descSortedList []string
		var ascSortedList []string
		assert.GreaterOrEqual(t, len(dsDescResp), 2, "Need at least 2 XMLIDs in Traffic Ops to test desc sort, found: %d", len(dsDescResp))
		// Get delivery services in the default ascending order for comparison.
		dsAscResp, _, err := TOSession.GetDeliveryServices(client.RequestOptions{})
		assert.NoError(t, err, "Unexpected error getting Delivery Services with default sort order: %v - alerts: %+v", err, dsAscResp.Alerts)
		assert.GreaterOrEqual(t, len(dsAscResp.Response), 2, "Need at least 2 XMLIDs in Traffic Ops to test sort, found %d", len(dsAscResp.Response))
		// Verify the response match in length, i.e. equal amount of delivery services.
		assert.Equal(t, len(dsAscResp.Response), len(dsDescResp), "Expected descending order response length: %v, to match ascending order response length %v", len(dsAscResp.Response), len(dsDescResp))
		// Insert xmlIDs to the front of a new list, so they are now reversed to be in ascending order.
		for _, ds := range dsDescResp {
			descSortedList = append([]string{ds.XMLID}, descSortedList...)
		}
		// Insert xmlIDs by appending to a new list, so they stay in ascending order.
		for _, ds := range dsAscResp.Response {
			ascSortedList = append(ascSortedList, ds.XMLID)
		}
		assert.Exactly(t, ascSortedList, descSortedList, "Delivery Service responses are not equal after reversal: %v - %v", ascSortedList, descSortedList)
	}
}

func GetDeliveryServiceId(t *testing.T, xmlId string) func() int {
	return func() int {
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("xmlId", xmlId)

		resp, _, err := TOSession.GetDeliveryServices(opts)
		assert.RequireNoError(t, err, "Get Delivery Service Request failed with error: %v", err)
		assert.RequireEqual(t, 1, len(resp.Response), "Expected delivery service response object length 1, but got %d", len(resp.Response))
		assert.RequireNotNil(t, resp.Response[0].ID, "Expected id to not be nil")

		return *resp.Response[0].ID
	}
}

func generateDeliveryService(t *testing.T, requestDS map[string]interface{}) map[string]interface{} {
	// map for the most basic HTTP Delivery Service a user can create
	genericHTTPDS := map[string]interface{}{
		"active":               tc.DSActiveStateActive,
		"cdnName":              "cdn1",
		"cdnId":                GetCDNID(t, "cdn1")(),
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
		"profileName":          "ATS_EDGE_TIER_CACHE",
		"qstringIgnore":        0,
		"rangeRequestHandling": 0,
		"regional":             false,
		"regionalGeoBlocking":  false,
		"routingName":          "ccr-ds1",
		"tenant":               "tenant1",
		"type":                 tc.DSTypeHTTP,
		"typeId":               GetTypeId(t, "HTTP"),
		"xmlId":                "ds1",
	}
	for k, v := range requestDS {
		genericHTTPDS[k] = v
	}
	return genericHTTPDS
}

func CreateTestDeliveryServices(t *testing.T) {
	for _, ds := range testData.DeliveryServices {
		resp, _, err := TOSession.CreateDeliveryService(ds, client.RequestOptions{})
		assert.NoError(t, err, "Could not create Delivery Service '%s': %v", ds.XMLID, err, resp.Alerts)
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
		assert.NoError(t, err, "Error deleting Delivery Service for '%s' : %v", ds.XMLID, err)
		assert.Equal(t, 0, len(getDS.Response), "Expected Delivery Service '%s' to be deleted", ds.XMLID)
	}
}
