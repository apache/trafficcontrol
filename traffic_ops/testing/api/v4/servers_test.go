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
	"time"

	"github.com/apache/trafficcontrol/lib/go-rfc"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/assert"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
	client "github.com/apache/trafficcontrol/traffic_ops/v4-client"
)

func TestServers(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Users, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, ServiceCategories, DeliveryServices}, func() {

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
				// else if resp.Summary.Count != 1 {
				"OK when VALID HOSTNAME parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"hostName": {"atlanta-edge-01"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1)),
				},
				// Ds ASsignments as prereqs // validate length
				"OK when VALID DSID parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"dsId": {"atlanta-edge-01"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1)),
				},
				"EMPTY RESPONSE when INVALID DSID parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"dsId": {"999999"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
				"FIRST RESULT when LIMIT=1": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"id"}, "limit": {"1"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validatePagination("limit")),
				},
				"SECOND RESULT when LIMIT=1 OFFSET=1": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"id"}, "limit": {"1"}, "offset": {"1"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validatePagination("offset")),
				},
				"SECOND RESULT when LIMIT=1 PAGE=2": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"id"}, "limit": {"1"}, "page": {"2"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validatePagination("page")),
				},
				"BAD REQUEST when INVALID LIMIT parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"limit": {"-2"}}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when INVALID OFFSET parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"limit": {"1"}, "offset": {"0"}}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when INVALID PAGE parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"limit": {"1"}, "page": {"0"}}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
			},
			"POST": {
				"BAD REQUEST when BLANK PROFILENAMES": {
					ClientSession: TOSession,
					RequestBody:   generateServer(t, map[string]interface{}{"profileNames": []string{""}}),
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusBadRequest)),
				},
			},
			"PUT": {
				// XMPP SHOULD NEVER CHANGE ?? VERIFY THIS
				"OK when VALID request": {
					EndpointId:    GetServerId(t, ""),
					ClientSession: TOSession,
					RequestBody: map[string]interface{}{
						"cdnId":        GetCDNId(t, "cdn1"),
						"cachegroupId": GetCacheGroupId(t, "cachegroup1")(),
						"domainName":   "updateddomainname",
						"hostName":     "atl-edge-01",
						"httpsPort":    8080,
						"interfaces": []map[string]interface{}{{
							"ipAddresses": []map[string]interface{}{{
								"address":        "127.0.0.1",
								"serviceAddress": true,
							}},
							"name": "eth0",
						}},
						"interfaceName":  "bond1",
						"interfaceMtu":   uint64(1280),
						"physLocationId": GetPhysLocationId(t, ""),
						"profileNames":   []string{""},
						"rack":           "RR 119.03",
						"statusId":       GetStatusId(t, "REPORTED"),
						"tcpPort":        8080,
						"typeId":         GetTypeId(t, "EDGE"),
						"updPending":     true,
						"xmppId":         "change-it",
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateServerFields(map[string]interface{}{"UpdPending": true, "TCPPort": 8080, "HTTPSPort": 8080,
							"DomainName": "updateddomainname", "XMPPID": "", "HostName": "atl-edge-01", "Rack": "RR 119.03",
							"InterfaceName": "bond1", "MTU": uint64(1280)})),
				},
				// Cannot update server type when assigned to a delivery service // prereq: assignment
				"BAD REQUEST when UPDATING TYPE of SERVER ASSIGNED to DS":                         {},
				"NO CHANGE to STATUSLASTUPDATED when STATUS is UNCHANGED":                         {},
				"STATUSLASTUPDATED UPDATED when STATUS CHANGES":                                   {},
				"BAD REQUEST when UPDATING CDN when LAST SERVER IN CACHEGROUP IN TOPOLOGY":        {},
				"BAD REQUEST when UPDATING CACHEGROUP when LAST SERVER IN CACHEGROUP IN TOPOLOGY": {},
				"BAD REQUEST when IPADDRESS EXISTS with SAME PROFILE": {
					EndpointId:    GetServerId(t, ""),
					ClientSession: TOSession,
					RequestBody:   map[string]interface{}{},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when BLANK HOSTNAME": {
					EndpointId:    GetServerId(t, "atlanta-edge-01"),
					ClientSession: TOSession,
					RequestBody:   generateServer(t, map[string]interface{}{"hostName": ""}),
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when BLANK DOMAINNAME": {
					EndpointId:    GetServerId(t, "atlanta-edge-01"),
					ClientSession: TOSession,
					RequestBody:   generateServer(t, map[string]interface{}{"domainName": ""}),
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"PRECONDITION FAILED when updating with IMS & IUS Headers": {
					EndpointId:    GetServerId(t, "atlanta-edge-01"),
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfUnmodifiedSince: {currentTimeRFC}}},
					RequestBody:   generateServer(t, map[string]interface{}{"hostName": "atlanta-edge-01"}),
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
				"PRECONDITION FAILED when updating with IFMATCH ETAG Header": {
					EndpointId:    GetServerId(t, "atlanta-edge-01"),
					ClientSession: TOSession,
					RequestBody:   generateServer(t, map[string]interface{}{"hostName": "atlanta-edge-01"}),
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfMatch: {rfc.ETag(currentTime)}}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
			},
			"DELETE": {
				"BAD REQUEST when LAST SERVER in CACHE GROUP": {
					EndpointId:    GetServerId(t, "midInTopologyMidCg01"),
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
			},
			"GET AFTER CHANGES": {
				"OK when CHANGES made": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfModifiedSince: {currentTimeRFC}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
			},
			"SERVER DETAILS GET": {
				// validate interface routerportName and routerName resp.Response[0].ServerInterfaces[0].RouterHostName
				"OK when VALID HOSTNAME parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"hostName": {"atlanta-edge-01"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1)),
				},
			},
		}

		GetTestServersQueryParameters(t)
	})
}

func validateServerFields(expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		serverResp := resp.([]tc.ServerV40)
		for field, expected := range expectedResp {
			for _, server := range serverResp {
				switch field {
				case "CacheGroup":
					assert.Equal(t, expected, *server.Cachegroup, "Expected Cachegroup to be %v, but got %v", expected, *server.Cachegroup)
				case "DomainName":
					assert.Equal(t, expected, *server.DomainName, "Expected DomainName to be %v, but got %v", expected, *server.DomainName)
				case "HostName":
					assert.Equal(t, expected, *server.HostName, "Expected HostName to be %v, but got %v", expected, *server.HostName)
				case "HTTPSPort":
					assert.Equal(t, expected, *server.HTTPSPort, "Expected HTTPSPort to be %v, but got %v", expected, *server.HTTPSPort)
				case "InterfaceName":
					assert.Equal(t, expected, server.Interfaces[0].Name, "Expected InterfaceName to be %v, but got %v", expected, server.Interfaces[0].Name)
				case "MTU":
					assert.Equal(t, expected, *server.Interfaces[0].MTU, "Expected MTU to be %v, but got %v", expected, *server.Interfaces[0].MTU)
				case "PhysLocation":
					assert.Equal(t, expected, *server.PhysLocation, "Expected PhysLocation to be %v, but got %v", expected, *server.PhysLocation)
				case "Rack":
					assert.Equal(t, expected, *server.Rack, "Expected Rack to be %v, but got %v", expected, *server.Rack)
				case "TCPPort":
					assert.Equal(t, expected, *server.TCPPort, "Expected TCPPort to be %v, but got %v", expected, *server.TCPPort)
				case "Type":
					assert.Equal(t, expected, server.Type, "Expected Type to be %v, but got %v", expected, server.Type)
				case "UpdPending":
					assert.Equal(t, expected, server.UpdPending, "Expected UpdPending to be %v, but got %v", expected, *server.UpdPending)
				case "XMPPID":
					assert.Equal(t, expected, *server.XMPPID, "Expected XMPPID to be %v, but got %v", expected, *server.XMPPID)
				default:
					t.Errorf("Expected field: %v, does not exist in response", field)
				}
			}
		}
	}
}

func GetServerId(t *testing.T, hostName string) func() int {
	return func() int {
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("hostName", hostName)
		serversResp, _, err := TOSession.GetServers(opts)
		assert.NoError(t, err, "Get Servers Request failed with error:", err)
		assert.Equal(t, 1, len(serversResp.Response), "Expected response object length 1, but got %d", len(serversResp.Response))
		assert.NotNil(t, serversResp.Response[0].ID, "Expected id to not be nil")
		return *serversResp.Response[0].ID
	}
}

func GetStatusId(t *testing.T, name string) func() int {
	return func() int {
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("name", name)
		statusResp, _, err := TOSession.GetStatuses(opts)
		assert.NoError(t, err, "Get Statuses Request failed with error:", err)
		assert.Equal(t, 1, len(statusResp.Response), "Expected response object length 1, but got %d", len(statusResp.Response))
		assert.NotNil(t, statusResp.Response[0].ID, "Expected id to not be nil")
		return statusResp.Response[0].ID
	}
}

func GetPhysLocationId(t *testing.T, name string) func() int {
	return func() int {
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("name", name)
		physLocResp, _, err := TOSession.GetPhysLocations(opts)
		assert.NoError(t, err, "Get PhysLocation Request failed with error:", err)
		assert.Equal(t, 1, len(physLocResp.Response), "Expected response object length 1, but got %d", len(physLocResp.Response))
		assert.NotNil(t, physLocResp.Response[0].ID, "Expected id to not be nil")
		return physLocResp.Response[0].ID
	}
}

func generateServer(t *testing.T, requestServer map[string]interface{}) map[string]interface{} {
	// map for the most basic Server a user can create
	genericServer := map[string]interface{}{
		"cdnId":        GetCDNId(t, "cdn1"),
		"cachegroupId": GetCacheGroupId(t, "cachegroup1")(),
		"domainName":   "localhost",
		"hostName":     "testserver",
		"interfaces": []map[string]interface{}{{
			"ipAddresses": []map[string]interface{}{{
				"address":        "127.0.0.1",
				"serviceAddress": true,
			}},
			"name": "eth0",
		}},
		"physLocationId": GetPhysLocationId(t, "Denver")(),
		"profileNames":   []string{"EDGE1"},
		"statusId":       GetStatusId(t, "REPORTED")(),
		"typeId":         GetTypeId(t, "EDGE"),
	}

	for k, v := range requestServer {
		genericServer[k] = v
	}
	return genericServer
}

func CreateTestServers(t *testing.T) {
	for _, server := range testData.Servers {
		resp, _, err := TOSession.CreateServer(server, client.RequestOptions{})
		assert.RequireNoError(t, err, "Could not create server '%s': %v - alerts: %+v", *server.HostName, err, resp.Alerts)
	}
}

func GetTestServersQueryParameters(t *testing.T) {
	dses, _, err := TOSession.GetDeliveryServices(client.RequestOptions{QueryParameters: url.Values{"xmlId": []string{"ds1"}}})
	if err != nil {
		t.Fatalf("Failed to get Delivery Services: %v - alerts: %+v", err, dses.Alerts)
	}
	if len(dses.Response) < 1 {
		t.Fatal("Failed to get at least one Delivery Service")
	}

	ds := dses.Response[0]
	if ds.ID == nil {
		t.Fatal("Traffic Ops returned a representation of a Delivery Service with null or undefined ID")
	}

	AssignTestDeliveryService(t)
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("dsId", strconv.Itoa(*ds.ID))
	servers, _, err := TOSession.GetServers(opts)
	if err != nil {
		t.Fatalf("Failed to get server by Delivery Service ID: %v - alerts: %+v", err, servers.Alerts)
	}
	if len(servers.Response) != 3 {
		t.Fatalf("expected to get 3 servers for Delivery Service: %d, actual: %d", *ds.ID, len(servers.Response))
	}

	dses, _, err = TOSession.GetDeliveryServices(client.RequestOptions{})
	if err != nil {
		t.Fatalf("Failed to get Delivery Services: %v - alerts: %+v", err, dses.Alerts)
	}

	foundTopDs := false
	const (
		topDSXmlID = "ds-top"
		topology   = "mso-topology"
	)
	for _, ds = range dses.Response {
		if ds.XMLID == nil || ds.ID == nil {
			t.Error("Traffic Ops returned a representation of a Delivery Service that had a null or undefined XMLID and/or ID")
			continue
		}
		if *ds.XMLID != topDSXmlID {
			continue
		}
		if ds.Topology == nil || ds.FirstHeaderRewrite == nil || ds.InnerHeaderRewrite == nil || ds.LastHeaderRewrite == nil {
			t.Errorf("Traffic Ops returned a representation of Delivery Service '%s' that had a null or undefined Topology and/or First Header Rewrite text and/or Inner Header Rewrite text and/or Last Header Rewrite text", topDSXmlID)
			continue
		}
		foundTopDs = true
		break
	}
	if !foundTopDs {
		t.Fatalf("unable to find deliveryservice %s", topDSXmlID)
	}

	/* Create a deliveryservice server assignment that should not show up in the
	 * client.GetServers( response because ds-top is topology-based
	 */
	const otherServerHostname = "topology-edge-02"
	serverResponse, _, err := TOSession.GetServers(client.RequestOptions{QueryParameters: url.Values{"hostName": []string{otherServerHostname}}})
	if err != nil {
		t.Fatalf("getting server by Host Name %s: %v - alerts: %+v", otherServerHostname, err, serverResponse.Alerts)
	}
	if len(serverResponse.Response) != 1 {
		t.Fatalf("unable to find server with hostname %s", otherServerHostname)
	}
	otherServer := serverResponse.Response[0]
	if otherServer.ID == nil || otherServer.HostName == nil {
		t.Fatal("Traffic Ops returned a representation of a Server that had a null or undefined ID and/or Host Name")
	}

	dsTopologyField, dsFirstHeaderRewriteField, innerHeaderRewriteField, lastHeaderRewriteField := *ds.Topology, *ds.FirstHeaderRewrite, *ds.InnerHeaderRewrite, *ds.LastHeaderRewrite
	ds.Topology, ds.FirstHeaderRewrite, ds.InnerHeaderRewrite, ds.LastHeaderRewrite = nil, nil, nil, nil
	updResp, _, err := TOSession.UpdateDeliveryService(*ds.ID, ds, client.RequestOptions{})
	if err != nil {
		t.Fatalf("unable to temporary remove topology-related fields from deliveryservice '%s': %v - alerts: %+v", topDSXmlID, err, updResp.Alerts)
	}
	if len(updResp.Response) != 1 {
		t.Fatalf("Expected updating a Delivery Service to update exactly one Delivery Service, but Traffic Ops indicates that %d were updated", len(updResp.Response))
	}
	ds = updResp.Response[0]
	if ds.ID == nil {
		t.Fatal("Traffic Ops returned a representation of a Delivery Service that had null or undefined ID")
	}
	assignResp, _, err := TOSession.CreateDeliveryServiceServers(*ds.ID, []int{*otherServer.ID}, false, client.RequestOptions{})
	if err != nil {
		t.Fatalf("unable to assign server '%s' to Delivery Service '%s': %v - alerts: %+v", *otherServer.HostName, topDSXmlID, err, assignResp.Alerts)
	}
	ds.Topology, ds.FirstHeaderRewrite, ds.InnerHeaderRewrite, ds.LastHeaderRewrite = &dsTopologyField, &dsFirstHeaderRewriteField, &innerHeaderRewriteField, &lastHeaderRewriteField
	updResp, _, err = TOSession.UpdateDeliveryService(*ds.ID, ds, client.RequestOptions{})
	if err != nil {
		t.Fatalf("unable to re-add topology-related fields to deliveryservice %s: %v - alerts: %+v", topDSXmlID, err, updResp.Alerts)
	}

	opts.Header = nil
	opts.QueryParameters.Set("dsId", strconv.Itoa(*ds.ID))
	expectedHostnames := map[string]bool{
		"edge1-cdn1-cg3":                 false,
		"edge2-cdn1-cg3":                 false,
		"atlanta-mid-01":                 false,
		"atlanta-mid-16":                 false,
		"edgeInCachegroup3":              false,
		"midInSecondaryCachegroupInCDN1": false,
	}
	response, _, err := TOSession.GetServers(opts)
	if err != nil {
		t.Fatalf("Failed to get servers by Topology-based Delivery Service ID with xmlId %s: %v - alerts: %+v", topDSXmlID, err, response.Alerts)
	}
	if len(response.Response) == 0 {
		t.Fatalf("Did not find any servers for Topology-based Delivery Service with xmlId %s", topDSXmlID)
	}
	for _, server := range response.Response {
		if server.HostName == nil {
			t.Fatal("Traffic Ops responded with a representation for a server with null or undefined Host Name")
		}
		if _, exists := expectedHostnames[*server.HostName]; !exists {
			t.Fatalf("expected hostnames %v, actual %s", expectedHostnames, *server.HostName)
		}
		expectedHostnames[*server.HostName] = true
	}
	var notInResponse []string
	for hostName, inResponse := range expectedHostnames {
		if !inResponse {
			notInResponse = append(notInResponse, hostName)
		}
	}
	if len(notInResponse) != 0 {
		t.Fatalf("%d servers missing from the response: %s", len(notInResponse), strings.Join(notInResponse, ", "))
	}
	const originHostname = "denver-mso-org-01"
	if resp, _, err := TOSession.AssignServersToDeliveryService([]string{originHostname}, topDSXmlID, client.RequestOptions{}); err != nil {
		t.Fatalf("assigning origin server '%s' to Delivery Service '%s': %v - alerts: %+v", originHostname, topDSXmlID, err, resp.Alerts)
	}
	response, _, err = TOSession.GetServers(opts)
	if err != nil {
		t.Fatalf("Failed to get servers by Topology-based Delivery Service ID with xmlId %s: %v - alerts: %+v", topDSXmlID, err, response.Alerts)
	}
	if len(response.Response) == 0 {
		t.Fatalf("Did not find any servers for Topology-based Delivery Service with xmlId %s", topDSXmlID)
	}
	containsOrigin := false
	for _, server := range response.Response {
		if server.HostName == nil || *server.HostName != originHostname {
			continue
		}
		containsOrigin = true
		break
	}
	if !containsOrigin {
		t.Fatalf("did not find origin server %s when querying servers by dsId after assigning %s to delivery service %s", originHostname, originHostname, topDSXmlID)
	}

	const topDsWithNoMids = "ds-based-top-with-no-mids"
	dses, _, err = TOSession.GetDeliveryServices(client.RequestOptions{QueryParameters: url.Values{"xmlId": []string{topDsWithNoMids}}})
	if err != nil {
		t.Fatalf("Failed to get Delivery Services: %v - alerts: %+v", err, dses.Alerts)
	}
	if len(dses.Response) < 1 {
		t.Fatal("Failed to get at least one Delivery Service")
	}

	ds = dses.Response[0]
	if ds.ID == nil {
		t.Fatal("Got Delivery Service with nil ID")
	}
	opts.QueryParameters.Set("dsId", strconv.Itoa(*ds.ID))

	response, _, err = TOSession.GetServers(opts)
	if err != nil {
		t.Fatalf("Failed to get servers by Topology-based Delivery Service ID with xmlId %s: %s", topDsWithNoMids, err)
	}
	if len(response.Response) == 0 {
		t.Fatalf("Did not find any servers for Topology-based Delivery Service with xmlId %s: %s", topDsWithNoMids, err)
	}
	for _, server := range response.Response {
		if server.HostName == nil {
			t.Fatal("Traffic Ops returned a server with null or undefined Host Name")
		}
		if server.Type == tc.CacheTypeMid.String() {
			t.Fatalf("Expected to find no %s-typed servers when querying servers by the ID for Delivery Service with XMLID %s but found %s-typed server %s", tc.CacheTypeMid, topDsWithNoMids, tc.CacheTypeMid, *server.HostName)
		}
	}

	opts.QueryParameters.Del("dsId")
	opts.QueryParameters.Set("topology", topology)
	expectedHostnames = map[string]bool{
		originHostname:                   false,
		"denver-mso-org-02":              false,
		"edge1-cdn1-cg3":                 false,
		"edge2-cdn1-cg3":                 false,
		"atlanta-mid-01":                 false,
		"atlanta-mid-16":                 false,
		"atlanta-mid-17":                 false,
		"edgeInCachegroup3":              false,
		"midInParentCachegroup":          false,
		"midInSecondaryCachegroup":       false,
		"midInSecondaryCachegroupInCDN1": false,
		"test-mso-org-01":                false,
	}
	response, _, err = TOSession.GetServers(opts)
	if err != nil {
		t.Fatalf("Failed to get servers belonging to Cache Groups in Topology %s: %v - alerts: %+v", topology, err, response.Alerts)
	}
	if len(response.Response) == 0 {
		t.Fatalf("Did not find any servers belonging to Cache Groups in Topology %s:", topology)
	}
	for _, server := range response.Response {
		if server.HostName == nil {
			t.Fatal("Traffic Ops returned a server with null or undefined Host Name")
		}
		if _, exists := expectedHostnames[*server.HostName]; !exists {
			t.Fatalf("expected hostnames %v, actual %s", expectedHostnames, *server.HostName)
		}
		expectedHostnames[*server.HostName] = true
	}
	notInResponse = []string{}
	for hostName, inResponse := range expectedHostnames {
		if !inResponse {
			notInResponse = append(notInResponse, hostName)
		}
	}
	if len(notInResponse) != 0 {
		t.Fatalf("%d servers missing from the response: %s", len(notInResponse), strings.Join(notInResponse, ", "))
	}
	opts.QueryParameters.Del("topology")

	resp, _, err := TOSession.GetServers(client.RequestOptions{})
	if err != nil {
		t.Fatalf("Failed to get servers: %v - alerts: %+v", err, resp.Alerts)
	}

	if len(resp.Response) < 1 {
		t.Fatal("Failed to get at least one server")
	}

	s := resp.Response[0]

	opts.QueryParameters.Set("type", s.Type)
	if resp, _, err := TOSession.GetServers(opts); err != nil {
		t.Errorf("Error getting servers by Type: %v - alerts: %+v", err, resp.Alerts)
	}
	opts.QueryParameters.Del("type")

	if s.CachegroupID == nil {
		t.Error("Found server with no Cache Group ID")
	} else {
		opts.QueryParameters.Add("cachegroup", strconv.Itoa(*s.CachegroupID))
		if resp, _, err := TOSession.GetServers(opts); err != nil {
			t.Errorf("Error getting servers by Cache Group ID: %v - alerts: %+v", err, resp.Alerts)
		}
		opts.QueryParameters.Del("cachegroup")
	}

	if s.Status == nil {
		t.Error("Found server with no status")
	} else {
		opts.QueryParameters.Add("status", *s.Status)
		if resp, _, err := TOSession.GetServers(opts); err != nil {
			t.Errorf("Error getting servers by status: %v - alerts: %+v", err, resp.Alerts)
		}
		opts.QueryParameters.Del("status")
	}

	opts.QueryParameters.Add("name", s.ProfileNames[0])
	pr, _, err := TOSession.GetProfiles(opts)
	if err != nil {
		t.Fatalf("failed to query profile: %v", err)
	}
	if len(pr.Response) != 1 {
		t.Error("Found server with no Profile ID")
	} else {
		profileID := pr.Response[0].ID
		opts.QueryParameters.Add("profileId", strconv.Itoa(profileID))
		if resp, _, err := TOSession.GetServers(opts); err != nil {
			t.Errorf("Error getting servers by Profile ID: %v - alerts: %+v", err, resp.Alerts)
		}
		opts.QueryParameters.Del("profileId")
	}

	cgs, _, err := TOSession.GetCacheGroups(client.RequestOptions{})
	if err != nil {
		t.Fatalf("Failed to get Cache Groups: %v", err)
	}
	if len(cgs.Response) < 1 {
		t.Fatal("Failed to get at least one Cache Group")
	}
	if cgs.Response[0].ID == nil {
		t.Fatal("Cache Group found with no ID")
	}

	opts.QueryParameters.Add("parentCacheGroup", strconv.Itoa(*cgs.Response[0].ID))
	if resp, _, err = TOSession.GetServers(opts); err != nil {
		t.Errorf("Error getting servers by parent Cache Group: %v - alerts: %+v", err, resp.Alerts)
	}
	opts.QueryParameters.Del("parentCacheGroup")
}

func DeleteTestServers(t *testing.T) {
	servers, _, err := TOSession.GetServers(client.RequestOptions{})
	assert.NoError(t, err, "Cannot get Servers: %v - alerts: %+v", err, servers.Alerts)

	for _, server := range servers.Response {
		delResp, _, err := TOSession.DeleteServer(*server.ID, client.RequestOptions{})
		assert.NoError(t, err, "Could not delete Server: %v - alerts: %+v", err, delResp.Alerts)
		// Retrieve Server to see if it got deleted
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("id", strconv.Itoa(*server.ID))
		getServer, _, err := TOSession.GetServers(opts)
		assert.NoError(t, err, "Error deleting Server for '%s' : %v - alerts: %+v", *server.HostName, err, getServer.Alerts)
		assert.Equal(t, 0, len(getServer.Response), "Expected Delivery Service '%s' to be deleted", *server.HostName)
	}
}
