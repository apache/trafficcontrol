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
	"reflect"
	"strconv"
	"strings"
	"testing"

	tc "github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	client "github.com/apache/trafficcontrol/v8/traffic_ops/v4-client"
)

func TestCDNsDNSSEC(t *testing.T) {
	if !includeSystemTests {
		t.Skip()
	}
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Users, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, ServerCapabilities, ServiceCategories, DeliveryServices}, func() {
		t.Run("GENERATE DNSSEC KEYS", func(t *testing.T) { GenerateDNSSECKeys(t) })
		t.Run("REFRESH DNSSEC KEYS", func(t *testing.T) { RefreshDNSSECKeys(t) }) // NOTE: testing refresh last (while no keys exist) because it's asynchronous and might affect other tests
	})
}

func RefreshDNSSECKeys(t *testing.T) {
	resp, reqInf, err := TOSession.RefreshDNSSECKeys(client.RequestOptions{})
	assert.NoError(t, err, "Unable to refresh DNSSEC keys: %v - alerts: %+v", err, resp.Alerts)
	assert.Equal(t, reqInf.StatusCode, http.StatusAccepted, "Refreshing DNSSEC keys - Expected: status code %d, Actual: %d", http.StatusAccepted, reqInf.StatusCode)

	loc := reqInf.RespHeaders.Get("Location")
	if loc == "" {
		t.Fatalf("Refreshing DNSSEC keys - Expected: non-empty 'Location' response header, Actual: empty")
	}
	locSplit := strings.Split(loc, "/")
	assert.RequireGreaterOrEqual(t, len(locSplit), 5, "Expected 'Location' response header to split into at least 5 parts, Got: %v", len(locSplit))
	asyncID, err := strconv.Atoi(locSplit[4])
	assert.RequireNoError(t, err, "Parsing async_status ID from 'Location' response header - Expected: no error, Actual: %v", err)

	status, _, err := TOSession.GetAsyncStatus(asyncID, client.RequestOptions{})
	assert.NoError(t, err, "Getting async status id %d - Expected: no error, Actual: %v", asyncID, err)
	assert.NotNil(t, status.Response.Message, "Getting async status for DNSSEC refresh job - Expected: non-nil message, Actual: nil")
}

func GenerateDNSSECKeys(t *testing.T) {
	assert.RequireGreaterOrEqual(t, len(testData.CDNs), 1, "Need at least one CDN to test updating CDNs")
	firstCDN := testData.CDNs[0]
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("name", firstCDN.Name)
	cdns, _, err := TOSession.GetCDNs(opts)
	assert.RequireNoError(t, err, "Unexpected error getting CDNs filtered by name '%s': %v - alerts: %+v", firstCDN.Name, err, cdns.Alerts)
	assert.RequireEqual(t, len(cdns.Response), 1, "Expected exactly one CDN named '%s' to exist, found: %d", firstCDN.Name, len(cdns.Response))

	cdn := cdns.Response[0]

	ttl := util.JSONIntStr(60)
	req := tc.CDNDNSSECGenerateReq{
		Key:               util.StrPtr(firstCDN.Name),
		TTL:               &ttl,
		KSKExpirationDays: &ttl,
		ZSKExpirationDays: &ttl,
	}
	resp, _, err := TOSession.GenerateCDNDNSSECKeys(req, client.RequestOptions{})
	assert.RequireNoError(t, err, "Unexpected error generating CDN DNSSEC keys: %v - alerts: %+v", err, resp.Alerts)

	res, _, err := TOSession.GetCDNDNSSECKeys(firstCDN.Name, client.RequestOptions{})
	assert.RequireNoError(t, err, "Unexpected error getting CDN DNSSEC keys: %v - alerts: %+v", err, res.Alerts)

	if _, ok := res.Response[firstCDN.Name]; !ok {
		t.Errorf("getting CDN DNSSEC keys - expected: key %s, actual: missing", firstCDN.Name)
	}
	originalKeys := res.Response

	resp, _, err = TOSession.GenerateCDNDNSSECKeys(req, client.RequestOptions{})
	assert.RequireNoError(t, err, "Unexpected error generating CDN DNSSEC keys: %v - alerts: %+v", err, resp.Alerts)

	res, _, err = TOSession.GetCDNDNSSECKeys(firstCDN.Name, client.RequestOptions{})
	assert.RequireNoError(t, err, "Unexpected error getting CDN DNSSEC keys: %v - alerts: %+v", err, res.Alerts)

	newKeys := res.Response

	if reflect.DeepEqual(originalKeys, newKeys) {
		t.Errorf("Generating CDN DNSSEC keys - expected: original keys to differ from new keys, actual: they are the same")
	}

	kskReq := tc.CDNGenerateKSKReq{
		ExpirationDays: util.Uint64Ptr(30),
	}
	originalKSK := newKeys
	resp, _, err = TOSession.GenerateCDNDNSSECKSK(firstCDN.Name, kskReq, client.RequestOptions{})
	assert.NoError(t, err, "Unexpected error generating DNSSEC KSK: %v - alerts: %+v", err, resp.Alerts)

	res, _, err = TOSession.GetCDNDNSSECKeys(firstCDN.Name, client.RequestOptions{})
	assert.NoError(t, err, "Unexpected error getting CDN DNSSEC keys: %v - alerts: %+v", err, res.Alerts)

	if _, ok := res.Response[firstCDN.Name]; !ok {
		t.Fatalf("getting CDN DNSSEC keys - expected: key %s, actual: missing", firstCDN.Name)
	}
	newKSK := res.Response
	if reflect.DeepEqual(originalKSK[firstCDN.Name].KSK, newKSK[firstCDN.Name].KSK) {
		t.Error("Generating CDN DNSSEC KSK - Expected: KSK to be different, Actual: KSK is the same")
	}
	if !reflect.DeepEqual(originalKSK[firstCDN.Name].ZSK, newKSK[firstCDN.Name].ZSK) {
		t.Error("Generating CDN DNSSEC KSK - Expected: ZSK to be equal, Actual: ZSK is different")
	}

	// ensure that when DNSSEC is enabled on a CDN, creating a new DS will generate DNSSEC keys for that DS:
	if !cdn.DNSSECEnabled {
		cdn.DNSSECEnabled = true
		resp, _, err := TOSession.UpdateCDN(cdn.ID, cdn, client.RequestOptions{})
		assert.NoError(t, err, "Unexpected error updating CDN: %v - alerts: %+v", err, resp.Alerts)

		defer func() {
			cdn.DNSSECEnabled = false
			resp, _, err := TOSession.UpdateCDN(cdn.ID, cdn, client.RequestOptions{})
			assert.NoError(t, err, "Unexpected error updating CDN: %v - alerts: %+v", err, resp.Alerts)
		}()
	}

	opts.QueryParameters.Set("name", "HTTP")
	types, _, err := TOSession.GetTypes(opts)
	assert.RequireNoError(t, err, "Unexpected error getting Types filtered by name 'HTTP': %v - alerts: %+v", err, types.Alerts)
	assert.RequireEqual(t, len(types.Response), 1, "Expected exactly one Type to exist with name 'HTTP', found: %d", len(types.Response))

	dsXMLID := "testdnssecgen"
	customDS := getCustomDS(cdn.ID, types.Response[0].ID, dsXMLID, "cdn", "https://testdnssecgen.example.com", dsXMLID)
	ds, _, err := TOSession.CreateDeliveryService(customDS, client.RequestOptions{})
	assert.NoError(t, err, "Unexpected error creating Delivery Service: %v - alerts: %+v", err, ds.Alerts)
	assert.RequireEqual(t, len(ds.Response), 1, "Expected creating a Delivery Service to create exactly one Delivery Service, Traffic Ops returned: %d", len(ds.Response))
	assert.RequireNotNil(t, ds.Response[0].ID, nil, "Traffic Ops returned a representation for a created Delivery Service with null or undefined ID")

	res, _, err = TOSession.GetCDNDNSSECKeys(firstCDN.Name, client.RequestOptions{})
	assert.RequireNoError(t, err, "Unexpected error getting CDN DNSSEC keys: %v - alerts: %+v", err, res.Alerts)

	if _, ok := res.Response[dsXMLID]; !ok {
		t.Error("after creating a new delivery service for a DNSSEC-enabled CDN - expected: DNSSEC keys to be found for the delivery service, actual: no DNSSEC keys found for the delivery service")
	}
	alerts, _, err := TOSession.DeleteDeliveryService(*ds.Response[0].ID, client.RequestOptions{})
	assert.NoError(t, err, "Unexpected error deleting Delivery Service: %v - alerts: %+v", err, alerts.Alerts)

	delResp, _, err := TOSession.DeleteCDNDNSSECKeys(firstCDN.Name, client.RequestOptions{})
	assert.NoError(t, err, "Unexpected error deleting CDN DNSSEC keys: %v - alerts: %+v", err, delResp.Alerts)
}
