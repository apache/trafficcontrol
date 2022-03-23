package v4

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/assert"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/deliveryservice"
	client "github.com/apache/trafficcontrol/traffic_ops/v4-client"
)

func TestDeliveryServicesUrlSigKeys(t *testing.T) {
	if !includeSystemTests {
		t.Skip()
	}
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Users, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, ServerCapabilities, ServiceCategories, DeliveryServices}, func() {

		methodTests := utils.V4TestCase{
			"GET": {
				// "deliveryservices/xmlId/%s/sslkeys" // t.Run("Verify SSL key generation on DS creation", VerifySSLKeysOnDsCreationTest)
				// Verifies fields are not empty strings // cert can be decoded
				"OK when VALID request": {
					ClientSession: TOSession, RequestOpts: client.RequestOptions{QueryParameters: url.Values{"xmlID": {"ds1"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1)),
				},
				// CreateTestDeliveryServicesURLSignatureKeys
				// /xmlId/%s/urlkeys

				///xmlId/%s/urlkeys // GetTestDeliveryServicesURLSignatureKeys
				// Verify no error

				//%s/urisignkeys  // CreateTestDeliveryServicesURISigningKeys
				// Verify they match with was created

				//%s/urisignkeys  // GetTestDeliveryServicesURISigningKeys
				// Verify no error
			},
			"POST": {
				// Blank CDN "sslkeytransfer" // NEW DS type=HTTP CDN="sslkeytransfer"
				// Generate keys // deliveryservices/sslkeys/generate
				// get keys for oldcdn: %s/name/%s/sslkeys"
				// get keys for new cdn: sslkeytransfer1

				// CreateTestDeliveryServicesURLSignatureKeys
				// FIRST POST
				///deliveryservices/xmlId/%s/urlkeys/generate
				// GET AND VERIFY they were created
				// Create new keys again and check that they are different

				//CreateTestDeliveryServicesURISigningKeys // deliveryservices/%s/urisignkeys
				// Perform another create using different keyset
				// verify the get req matches new keyset not old one

				// DeliveryServiceSSLKeys /sslkeys/add
			},
			"PUT": {
				"BAD REQUEST when updating DS ROUTING NAME when DS has SSL KEYS": {},
				"BAD REQUEST when updating DS CDN when DS has SSL KEYS":          {},
			},
			"DELETE": {
				///xmlId/%s/urlkeys // DeleteTestDeliveryServicesURLSignatureKeys

				///%s/urisignkeys // DeleteTestDeliveryServicesURISigningKeys
				// Expect certain keys to not be empty?

				//t.Errorf("expected 1 ssl keys for CDN %v, got %d instead", cdn.Name, len(newCdnKeys))
				// Expect on CDN key one deleting a DS that belonged to that cdn
				// PreREq: create cdn, 2 dses attached, get keys, delete ds get cdn keys // expect len=1

			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					switch method {
					case "GET":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.GetDeliveryServiceSSLKeys(testCase.RequestOpts.QueryParameters.Get("xmlID"), testCase.RequestOpts)
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
		t.Run("Verify SSL key generation on DS creation", VerifySSLKeysOnDsCreationTest)
		t.Run("Update CDN for a Delivery Service with SSL keys", SSLDeliveryServiceCDNUpdateTest)
		t.Run("Create URL Signature keys for a Delivery Service", CreateTestDeliveryServicesURLSignatureKeys)
		t.Run("Retrieve URL Signature keys for a Delivery Service", GetTestDeliveryServicesURLSignatureKeys)
		t.Run("Delete URL Signature keys for a Delivery Service", DeleteTestDeliveryServicesURLSignatureKeys)
		t.Run("Create URI Signing Keys for a Delivery Service", CreateTestDeliveryServicesURISigningKeys)
		t.Run("Retrieve URI Signing keys for a Delivery Service", GetTestDeliveryServicesURISigningKeys)
		t.Run("Delete URI Signing keys for a Delivery Service", DeleteTestDeliveryServicesURISigningKeys)
		t.Run("Delete old CDN SSL keys", DeleteCDNOldSSLKeys)
		t.Run("Create and retrieve SSL keys for a Delivery Service", DeliveryServiceSSLKeys)
	})
}

func validateSSLKeyExpectedFields(expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		dsSSLKey := resp.(tc.DeliveryServiceSSLKeys)
		for field, expected := range expectedResp {
			switch field {
			case "Key":
				assert.Equal(t, expected, dsSSLKey.Certificate.Key, "Expected Certificate Key to be %v, but got %v", expected, dsSSLKey.Certificate.Key)
			case "Crt":
				assert.Equal(t, expected, dsSSLKey.Certificate.Crt, "Expected Certificate Crt to be %v, but got %v", expected, dsSSLKey.Certificate.Crt)
			case "CSR":
				assert.Equal(t, expected, dsSSLKey.Certificate.CSR, "Expected Certificate CSR to be %v, but got %v", expected, dsSSLKey.Certificate.CSR)
			default:
				t.Errorf("Expected field: %v, does not exist in response", field)
			}
		}
	}
}

func validateCertDecoding() utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		dsSSLKey := resp.(tc.DeliveryServiceSSLKeys)
		err := deliveryservice.Base64DecodeCertificate(&dsSSLKey.Certificate)
		assert.RequireNoError(t, err, "Couldn't decode certificate: %v", err)
	}
}

func createBlankCDN(cdnName string, t *testing.T) tc.CDN {
	_, _, err := TOSession.CreateCDN(tc.CDN{
		DNSSECEnabled: false,
		DomainName:    cdnName + ".ai",
		Name:          cdnName,
	}, client.RequestOptions{})
	if err != nil {
		t.Fatal("unable to create cdn: " + err.Error())
	}

	originalKeys, _, err := TOSession.GetCDNSSLKeys(cdnName, client.RequestOptions{})
	if err != nil {
		t.Fatalf("unable to get keys on cdn %v: %v - alerts: %+v", cdnName, err, originalKeys.Alerts)
	}

	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("name", cdnName)
	cdns, _, err := TOSession.GetCDNs(opts)
	if err != nil {
		t.Fatalf("unable to get cdn %v: %v - alerts: %+v", cdnName, err, cdns.Alerts)
	}
	if len(cdns.Response) < 1 {
		t.Fatal("expected more than 0 cdns")
	}
	keys, _, err := TOSession.GetCDNSSLKeys(cdnName, client.RequestOptions{})
	if err != nil {
		t.Fatalf("unable to get keys on cdn %v: %v", cdnName, err)
	}
	if len(keys.Response) != len(originalKeys.Response) {
		t.Fatalf("expected %v ssl keys on cdn %v, got %v", len(originalKeys.Response), cdnName, len(keys.Response))
	}
	return cdns.Response[0]
}

func cleanUp(t *testing.T, ds tc.DeliveryServiceV4, oldCDNID int, newCDNID int, sslKeyVersions []string) {
	if ds.ID == nil || ds.XMLID == nil {
		t.Error("Cannot clean up Delivery Service with nil ID and/or XMLID")
		return
	}
	xmlid := *ds.XMLID
	id := *ds.ID

	opts := client.NewRequestOptions()
	for _, version := range sslKeyVersions {
		opts.QueryParameters.Set("version", version)
		resp, _, err := TOSession.DeleteDeliveryServiceSSLKeys(xmlid, opts)
		if err != nil {
			t.Errorf("Unexpected error deleting Delivery Service SSL Keys: %v - alerts: %+v", err, resp.Alerts)
		}
	}
	resp, _, err := TOSession.DeleteDeliveryService(id, client.RequestOptions{})
	if err != nil {
		t.Errorf("Unexpected error deleting Delivery Service '%s' (#%d) during cleanup: %v - alerts: %+v", xmlid, id, err, resp.Alerts)
	}
	if oldCDNID != -1 {
		resp2, _, err := TOSession.DeleteCDN(oldCDNID, client.RequestOptions{})
		if err != nil {
			t.Errorf("Unexpected error deleting CDN (#%d) during cleanup: %v - alerts: %+v", oldCDNID, err, resp2.Alerts)
		}
	}
	if newCDNID != -1 {
		resp2, _, err := TOSession.DeleteCDN(newCDNID, client.RequestOptions{})
		if err != nil {
			t.Errorf("Unexpected error deleting CDN (#%d) during cleanup: %v - alerts: %+v", newCDNID, err, resp2.Alerts)
		}
	}
}

// getCustomDS returns a DS that is guaranteed to have non-nil:
//
//    Active
//    CDNID
//    DSCP
//    DisplayName
//    RoutingName
//    GeoLimit
//    GeoProvider
//    IPV6RoutingEnabled
//    InitialDispersion
//    LogsEnabled
//    MissLat
//    MissLong
//    MultiSiteOrigin
//    OrgServerFQDN
//    Protocol
//    QStringIgnore
//    RangeRequestHandling
//    RegionalGeoBlocking
//    TenantID
//    TypeID
//    XMLID
//
// BUT, will ALWAYS have nil MaxRequestHeaderBytes.
// Note that the Tenant is hard-coded to #1.
func getCustomDS(cdnID, typeID int, displayName, routingName, orgFQDN, dsID string) tc.DeliveryServiceV4 {
	customDS := tc.DeliveryServiceV4{}
	customDS.Active = util.BoolPtr(true)
	customDS.CDNID = util.IntPtr(cdnID)
	customDS.DSCP = util.IntPtr(0)
	customDS.DisplayName = util.StrPtr(displayName)
	customDS.RoutingName = util.StrPtr(routingName)
	customDS.GeoLimit = util.IntPtr(0)
	customDS.GeoProvider = util.IntPtr(0)
	customDS.IPV6RoutingEnabled = util.BoolPtr(false)
	customDS.InitialDispersion = util.IntPtr(1)
	customDS.LogsEnabled = util.BoolPtr(true)
	customDS.MissLat = util.FloatPtr(0)
	customDS.MissLong = util.FloatPtr(0)
	customDS.MultiSiteOrigin = util.BoolPtr(false)
	customDS.OrgServerFQDN = util.StrPtr(orgFQDN)
	customDS.Protocol = util.IntPtr(2)
	customDS.QStringIgnore = util.IntPtr(0)
	customDS.RangeRequestHandling = util.IntPtr(0)
	customDS.RegionalGeoBlocking = util.BoolPtr(false)
	customDS.TenantID = util.IntPtr(1)
	customDS.TypeID = util.IntPtr(typeID)
	customDS.XMLID = util.StrPtr(dsID)
	customDS.MaxRequestHeaderBytes = nil
	return customDS
}

func DeleteCDNOldSSLKeys(t *testing.T) {
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

	// First DS creation
	customDS := getCustomDS(cdn.ID, types.Response[0].ID, "displayName", "routingName", "https://test.com", "dsID")

	resp, _, err := TOSession.CreateDeliveryService(customDS, client.RequestOptions{})
	if err != nil {
		t.Errorf("Unexpected error creating a Delivery Service: %v - alerts: %+v", err, resp.Alerts)
	}
	if len(resp.Response) != 1 {
		t.Fatalf("Expected Delivery Service creation to return exactly one Delivery Service, got: %d", len(resp.Response))
	}
	ds := resp.Response[0]
	if ds.XMLID == nil {
		t.Fatal("Traffic Ops returned a representation for a Delivery Service with null or undefined XMLID")
	}

	ds.CDNName = &cdn.Name
	sslKeyRequestFields := tc.SSLKeyRequestFields{
		BusinessUnit: util.StrPtr("BU"),
		City:         util.StrPtr("CI"),
		Organization: util.StrPtr("OR"),
		HostName:     util.StrPtr("*.test.com"),
		Country:      util.StrPtr("CO"),
		State:        util.StrPtr("ST"),
	}
	genResp, _, err := TOSession.GenerateSSLKeysForDS(*ds.XMLID, *ds.CDNName, sslKeyRequestFields, client.RequestOptions{})
	if err != nil {
		t.Fatalf("Unexpected error generaing SSL Keys for Delivery Service '%s': %v - alerts: %+v", *ds.XMLID, err, genResp.Alerts)
	}
	defer cleanUp(t, ds, cdn.ID, -1, []string{"1"})

	// Second DS creation
	customDS2 := getCustomDS(cdn.ID, types.Response[0].ID, "displayName2", "routingName2", "https://test2.com", "dsID2")

	resp, _, err = TOSession.CreateDeliveryService(customDS2, client.RequestOptions{})
	if err != nil {
		t.Errorf("Unexpected error creating a Delivery Service: %v - alerts: %+v", err, resp.Alerts)
	}
	if len(resp.Response) != 1 {
		t.Fatalf("Expected Delivery Service creation to return exactly one Delivery Service, got: %d", len(resp.Response))
	}
	ds2 := resp.Response[0]
	if ds2.XMLID == nil || ds2.ID == nil {
		t.Fatal("Traffic Ops returned a representation for a Delivery Service with null or undefined XMLID and/or ID")
	}

	ds2.CDNName = &cdn.Name
	sslKeyRequestFields.HostName = util.StrPtr("*.test2.com")
	genResp, _, err = TOSession.GenerateSSLKeysForDS(*ds2.XMLID, *ds2.CDNName, sslKeyRequestFields, client.RequestOptions{})
	if err != nil {
		t.Fatalf("Unexpected error generaing SSL Keys for Delivery Service '%s': %v - alerts: %+v", *ds2.XMLID, err, genResp.Alerts)
	}

	var cdnKeys []tc.CDNSSLKeys
	for tries := 0; tries < 5; tries++ {
		time.Sleep(time.Second)
		var sslKeysResp tc.CDNSSLKeysResponse
		sslKeysResp, _, err = TOSession.GetCDNSSLKeys(cdn.Name, client.RequestOptions{})
		if err != nil {
			continue
		}
		cdnKeys = sslKeysResp.Response
		if len(cdnKeys) != 0 {
			break
		}
	}

	if err != nil {
		t.Fatalf("unable to get CDN %v SSL keys: %v", cdn.Name, err)
	}
	if len(cdnKeys) != 2 {
		t.Errorf("expected two ssl keys for CDN %v, got %d instead", cdn.Name, len(cdnKeys))
	}

	delResp, _, err := TOSession.DeleteDeliveryService(*ds2.ID, client.RequestOptions{})
	if err != nil {
		t.Errorf("Unexpected error deleting Delivery Service #%d: %v - alerts: %+v", *ds2.ID, err, delResp.Alerts)
	}

	opts = client.NewRequestOptions()
	opts.QueryParameters.Set("cdnID", strconv.Itoa(cdn.ID))
	snapResp, _, err := TOSession.SnapshotCRConfig(opts)
	if err != nil {
		t.Fatalf("Failed to take Snapshot of CDN #%d: %v - alerts: %+v", cdn.ID, err, snapResp.Alerts)
	}
	var newCdnKeys []tc.CDNSSLKeys
	for tries := 0; tries < 5; tries++ {
		time.Sleep(time.Second)
		var sslKeysResp tc.CDNSSLKeysResponse
		sslKeysResp, _, err = TOSession.GetCDNSSLKeys(cdn.Name, client.RequestOptions{})
		newCdnKeys = sslKeysResp.Response
		if err == nil && len(newCdnKeys) == 1 {
			break
		}
	}

	if err != nil {
		t.Fatalf("unable to get CDN %v SSL keys: %v", cdn.Name, err)
	}
	if len(newCdnKeys) != 1 {
		t.Errorf("expected 1 ssl keys for CDN %v, got %d instead", cdn.Name, len(newCdnKeys))
	}
}

func DeliveryServiceSSLKeys(t *testing.T) {
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

	customDS := getCustomDS(cdn.ID, types.Response[0].ID, "displayName", "routingName", "https://test.com", "dsID")

	resp, _, err := TOSession.CreateDeliveryService(customDS, client.RequestOptions{})
	if err != nil {
		t.Errorf("Unexpected error creating a Delivery Service: %v - alerts: %+v", err, resp.Alerts)
	}
	if len(resp.Response) != 1 {
		t.Fatalf("Expected Delivery Service creation to return exactly one Delivery Service, got: %d", len(resp.Response))
	}
	ds := resp.Response[0]
	if ds.XMLID == nil {
		t.Fatal("Traffic Ops returned a representation for a Delivery Service with null or undefined XMLID")
	}

	ds.CDNName = &cdn.Name
	genResp, _, err := TOSession.GenerateSSLKeysForDS(*ds.XMLID, *ds.CDNName, tc.SSLKeyRequestFields{
		BusinessUnit: util.StrPtr("BU"),
		City:         util.StrPtr("CI"),
		Organization: util.StrPtr("OR"),
		HostName:     util.StrPtr("*.test2.com"),
		Country:      util.StrPtr("CO"),
		State:        util.StrPtr("ST"),
	}, client.RequestOptions{})
	if err != nil {
		t.Fatalf("Unexpected error generating SSL Keys for Delivery Service '%s': %v - alerts: %+v", *ds.XMLID, err, genResp.Alerts)
	}
	defer cleanUp(t, ds, cdn.ID, -1, []string{"1"})

	dsSSLKey := new(tc.DeliveryServiceSSLKeys)
	for tries := 0; tries < 5; tries++ {
		time.Sleep(time.Second)
		var sslKeysResp tc.DeliveryServiceSSLKeysResponse
		sslKeysResp, _, err = TOSession.GetDeliveryServiceSSLKeys(*ds.XMLID, client.RequestOptions{})
		*dsSSLKey = sslKeysResp.Response
		if err == nil && dsSSLKey != nil {
			break
		}
	}

	if err != nil || dsSSLKey == nil {
		t.Fatalf("unable to get DS %s SSL key: %v", *ds.XMLID, err)
	}
	if dsSSLKey.Certificate.Key == "" {
		t.Errorf("expected a valid key but got nothing")
	}
	if dsSSLKey.Certificate.Crt == "" {
		t.Errorf("expected a valid certificate, but got nothing")
	}
	if dsSSLKey.Certificate.CSR == "" {
		t.Errorf("expected a valid CSR, but got nothing")
	}

	err = deliveryservice.Base64DecodeCertificate(&dsSSLKey.Certificate)
	if err != nil {
		t.Fatalf("couldn't decode certificate: %v", err)
	}

	dsSSLKeyReq := tc.DeliveryServiceSSLKeysReq{
		AuthType:        &dsSSLKey.AuthType,
		CDN:             &dsSSLKey.CDN,
		DeliveryService: &dsSSLKey.DeliveryService,
		BusinessUnit:    &dsSSLKey.BusinessUnit,
		City:            &dsSSLKey.City,
		Organization:    &dsSSLKey.Organization,
		HostName:        &dsSSLKey.Hostname,
		Country:         &dsSSLKey.Country,
		State:           &dsSSLKey.State,
		Key:             &dsSSLKey.Key,
		Version:         &dsSSLKey.Version,
		Certificate:     &dsSSLKey.Certificate,
	}
	addSSLKeysResp, _, err := TOSession.AddSSLKeysForDS(tc.DeliveryServiceAddSSLKeysReq{DeliveryServiceSSLKeysReq: dsSSLKeyReq}, client.RequestOptions{})
	if err != nil {
		t.Errorf("Unexpected error adding SSL keys for Delivery Service '%s': %v - alerts: %+v", dsSSLKey.DeliveryService, err, addSSLKeysResp.Alerts)
	}

	dsSSLKey = new(tc.DeliveryServiceSSLKeys)
	for tries := 0; tries < 5; tries++ {
		time.Sleep(time.Second)
		var sslKeysResp tc.DeliveryServiceSSLKeysResponse
		sslKeysResp, _, err = TOSession.GetDeliveryServiceSSLKeys(*ds.XMLID, client.RequestOptions{})
		*dsSSLKey = sslKeysResp.Response
		if err == nil && dsSSLKey != nil {
			break
		}
	}

	if err != nil || dsSSLKey == nil {
		t.Fatalf("unable to get DS %s SSL key: %v", *ds.XMLID, err)
	}
	if dsSSLKey.Certificate.Key == "" {
		t.Errorf("expected a valid key but got nothing")
	}
	if dsSSLKey.Certificate.Crt == "" {
		t.Errorf("expected a valid certificate, but got nothing")
	}
	if dsSSLKey.Certificate.CSR == "" {
		t.Errorf("expected a valid CSR, but got nothing")
	}
}

func VerifySSLKeysOnDsCreationTest(t *testing.T) {
	for _, ds := range testData.DeliveryServices {
		if !(*ds.Protocol == tc.DSProtocolHTTPS || *ds.Protocol == tc.DSProtocolHTTPAndHTTPS || *ds.Protocol == tc.DSProtocolHTTPToHTTPS) {
			continue
		}
		var err error
		dsSSLKey := new(tc.DeliveryServiceSSLKeys)
		for tries := 0; tries < 5; tries++ {
			time.Sleep(time.Second)
			var sslKeysResp tc.DeliveryServiceSSLKeysResponse
			sslKeysResp, _, err = TOSession.GetDeliveryServiceSSLKeys(*ds.XMLID, client.RequestOptions{})
			*dsSSLKey = sslKeysResp.Response
			if err == nil && dsSSLKey != nil {
				break
			}
		}

		if err != nil || dsSSLKey == nil {
			t.Fatalf("unable to get DS %s SSL key: %v", *ds.XMLID, err)
		}
		if dsSSLKey.Certificate.Key == "" {
			t.Errorf("expected a valid key but got nothing")
		}
		if dsSSLKey.Certificate.Crt == "" {
			t.Errorf("expected a valid certificate, but got nothing")
		}
		if dsSSLKey.Certificate.CSR == "" {
			t.Errorf("expected a valid CSR, but got nothing")
		}

		err = deliveryservice.Base64DecodeCertificate(&dsSSLKey.Certificate)
		if err != nil {
			t.Fatalf("couldn't decode certificate: %v", err)
		}
	}
}

func SSLDeliveryServiceCDNUpdateTest(t *testing.T) {
	cdnNameOld := "sslkeytransfer"
	oldCdn := createBlankCDN(cdnNameOld, t)
	cdnNameNew := "sslkeytransfer1"
	newCdn := createBlankCDN(cdnNameNew, t)

	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("name", "HTTP")
	types, _, err := TOSession.GetTypes(opts)
	if err != nil {
		t.Fatalf("unable to get Types: %v - alerts: %+v", err, types.Alerts)
	}
	if len(types.Response) < 1 {
		t.Fatal("expected at least one type")
	}

	customDS := getCustomDS(oldCdn.ID, types.Response[0].ID, "displayName", "routingName", "https://test.com", "dsID")

	resp, _, err := TOSession.CreateDeliveryService(customDS, client.RequestOptions{})
	if err != nil {
		t.Errorf("Unexpected error creating a custom Delivery Service: %v - alerts: %+v", err, resp.Alerts)
	}
	if len(resp.Response) != 1 {
		t.Fatalf("Expected Delivery Service creation to create exactly one Delivery Service, Traffic Ops indicates %d were created", len(resp.Response))
	}
	ds := resp.Response[0]
	if ds.ID == nil || ds.XMLID == nil {
		t.Fatal("Traffic Ops created a Delivery Service with null or undefined XMLID and/or ID")
	}
	ds.CDNName = &oldCdn.Name

	defer cleanUp(t, ds, oldCdn.ID, newCdn.ID, []string{"1"})

	_, _, err = TOSession.GenerateSSLKeysForDS(*ds.XMLID, *ds.CDNName, tc.SSLKeyRequestFields{
		BusinessUnit: util.StrPtr("BU"),
		City:         util.StrPtr("CI"),
		Organization: util.StrPtr("OR"),
		HostName:     util.StrPtr("*.test.com"),
		Country:      util.StrPtr("CO"),
		State:        util.StrPtr("ST"),
	}, client.RequestOptions{})
	if err != nil {
		t.Fatalf("unable to generate sslkeys for cdn %v: %v", oldCdn.Name, err)
	}

	var oldCDNKeys []tc.CDNSSLKeys
	for tries := 0; tries < 5; tries++ {
		time.Sleep(time.Second)
		resp, _, err := TOSession.GetCDNSSLKeys(oldCdn.Name, client.RequestOptions{})
		oldCDNKeys = resp.Response
		if err == nil && len(oldCDNKeys) > 0 {
			break
		}
	}
	if err != nil {
		t.Fatalf("unable to get cdn %v keys: %v", oldCdn.Name, err)
	}
	if len(oldCDNKeys) < 1 {
		t.Fatal("expected at least 1 key")
	}

	newCDNKeys, _, err := TOSession.GetCDNSSLKeys(newCdn.Name, client.RequestOptions{})
	if err != nil {
		t.Fatalf("unable to get cdn %v keys: %v", newCdn.Name, err)
	}

	ds.RoutingName = util.StrPtr("anothername")
	_, _, err = TOSession.UpdateDeliveryService(*ds.ID, ds, client.RequestOptions{})
	if err == nil {
		t.Fatal("should not be able to update delivery service (routing name) as it has ssl keys")
	}
	ds.RoutingName = util.StrPtr("routingName")

	ds.CDNID = &newCdn.ID
	ds.CDNName = &newCdn.Name
	_, _, err = TOSession.UpdateDeliveryService(*ds.ID, ds, client.RequestOptions{})
	if err == nil {
		t.Fatal("should not be able to update delivery service (cdn) as it has ssl keys")
	}

	// Check new CDN still has an ssl key
	keys, _, err := TOSession.GetCDNSSLKeys(newCdn.Name, client.RequestOptions{})
	if err != nil {
		t.Fatalf("unable to get cdn %v keys: %v - alerts: %+v", newCdn.Name, err, keys.Alerts)
	}
	if len(keys.Response) != len(newCDNKeys.Response) {
		t.Fatalf("expected %v keys, got %v", len(newCDNKeys.Response), len(keys.Response))
	}

	// Check old CDN does not have ssl key
	keys, _, err = TOSession.GetCDNSSLKeys(oldCdn.Name, client.RequestOptions{})
	if err != nil {
		t.Fatalf("unable to get cdn %v keys: %v - %+v", oldCdn.Name, err, keys.Alerts)
	}
	if len(keys.Response) != len(oldCDNKeys) {
		t.Fatalf("expected %v key, got %v", len(oldCDNKeys), len(keys.Response))
	}
}

func GetTestDeliveryServicesURLSignatureKeys(t *testing.T) {
	if len(testData.DeliveryServices) == 0 {
		t.Fatal("couldn't get the xml ID of test DS")
	}
	firstDS := testData.DeliveryServices[0]
	if firstDS.XMLID == nil {
		t.Fatal("Found a Delivery Service in testing data with a null or undefined XMLID")
	}

	_, _, err := TOSession.GetDeliveryServiceURLSignatureKeys(*firstDS.XMLID, client.RequestOptions{})
	if err != nil {
		t.Errorf("failed to get url sig keys: %v", err)
	}
}

func CreateTestDeliveryServicesURLSignatureKeys(t *testing.T) {
	if len(testData.DeliveryServices) == 0 {
		t.Fatal("couldn't get the xml ID of test DS")
	}
	firstDS := testData.DeliveryServices[0]

	if firstDS.XMLID == nil {
		t.Fatal("Found a Delivery Service in testing data with a null or undefined XMLID")
	}

	resp, _, err := TOSession.CreateDeliveryServiceURLSignatureKeys(*firstDS.XMLID, client.RequestOptions{})
	if err != nil {
		t.Errorf("Unexpected error creaetting URL signing keys: %v - alerts: %+v", err, resp.Alerts)
	}
	firstKeys, _, err := TOSession.GetDeliveryServiceURLSignatureKeys(*firstDS.XMLID, client.RequestOptions{})
	if err != nil {
		t.Errorf("Unexpected error getting URL signing keys: %v - alerts: %+v", err, firstKeys.Alerts)
	}
	if len(firstKeys.Response) == 0 {
		t.Errorf("failed to create URL signing keys")
	}
	firstKey, ok := firstKeys.Response["key0"]
	if !ok {
		t.Fatal("Expected to find 'key0' in URL signing keys, but didn't")
	}

	// Create new keys again and check that they are different
	resp, _, err = TOSession.CreateDeliveryServiceURLSignatureKeys(*firstDS.XMLID, client.RequestOptions{})
	if err != nil {
		t.Errorf("Unexpected error creating URL signing keys: %v - alerts: %+v", err, resp.Alerts)
	}
	secondKeys, _, err := TOSession.GetDeliveryServiceURLSignatureKeys(*firstDS.XMLID, client.RequestOptions{})
	if err != nil {
		t.Errorf("Unexpected error getting URL signing keys: %v - alerts: %+v", err, secondKeys.Alerts)
	}
	if len(secondKeys.Response) == 0 {
		t.Errorf("failed to create url sig keys")
	}
	secondKey, ok := secondKeys.Response["key0"]
	if !ok {
		t.Fatal("Expected to find 'key0' in URL signing keys, but didn't")
	}

	if secondKey == firstKey {
		t.Errorf("second create did not generate new url sig keys")
	}
}

func DeleteTestDeliveryServicesURLSignatureKeys(t *testing.T) {
	if len(testData.DeliveryServices) == 0 {
		t.Fatal("couldn't get the xml ID of test DS")
	}
	firstDS := testData.DeliveryServices[0]

	if firstDS.XMLID == nil {
		t.Fatal("Found a Delivery Service in testing data with a null or undefined XMLID")
	}

	resp, _, err := TOSession.DeleteDeliveryServiceURLSignatureKeys(*firstDS.XMLID, client.RequestOptions{})
	if err != nil {
		t.Errorf("Unexpected error deletining URL signing keys: %v - alerts: %+v", err, resp.Alerts)
	}

}

func GetTestDeliveryServicesURISigningKeys(t *testing.T) {
	if len(testData.DeliveryServices) == 0 {
		t.Fatal("couldn't get the xml ID of test DS")
	}
	firstDS := testData.DeliveryServices[0]

	if firstDS.XMLID == nil {
		t.Fatal("Found a Delivery Service in testing data with a null or undefined XMLID")
	}

	_, _, err := TOSession.GetDeliveryServiceURISigningKeys(*firstDS.XMLID, client.RequestOptions{})
	if err != nil {
		t.Errorf("Unexpected error getting URI signing keys for Delivery Service '%s': %v", *firstDS.XMLID, err)
	}
}

const (
	keySet1 = `
{
  "Kabletown URI Authority 1": {
    "renewal_kid": "First Key",
    "keys": [
      {
        "alg": "HS256",
        "kid": "First Key",
        "kty": "oct",
        "k": "Kh_RkUMj-fzbD37qBnDf_3e_RvQ3RP9PaSmVEpE24AM"
      }
    ]
  }
}`
	keySet2 = `
{
"Kabletown URI Authority 1": {
    "renewal_kid": "New First Key",
    "keys": [
      {
        "alg": "HS256",
        "kid": "New First Key",
        "kty": "oct",
        "k": "Kh_RkUMj-fzbD37qBnDf_3e_RvQ3RP9PaSmVEpE24AM"
      }
    ]
  }
}`
)

func CreateTestDeliveryServicesURISigningKeys(t *testing.T) {
	if len(testData.DeliveryServices) == 0 {
		t.Fatal("couldn't get the xml ID of test DS")
	}
	firstDS := testData.DeliveryServices[0]
	if firstDS.XMLID == nil {
		t.Fatal("Found a Delivery Service in testing data with a null or undefined XMLID")
	}

	var keyset tc.JWKSMap

	if err := json.Unmarshal([]byte(keySet1), &keyset); err != nil {
		t.Errorf("json.UnMarshal(): expected nil error, actual: %v", err)
	}

	_, _, err := TOSession.CreateDeliveryServiceURISigningKeys(*firstDS.XMLID, keyset, client.RequestOptions{})
	if err != nil {
		t.Error("failed to create uri sig keys: " + err.Error())
	}

	firstKeysBytes, _, err := TOSession.GetDeliveryServiceURISigningKeys(*firstDS.XMLID, client.RequestOptions{})
	if err != nil {
		t.Error("failed to get uri sig keys: " + err.Error())
	}

	firstKeys := tc.JWKSMap{}
	if err := json.Unmarshal(firstKeysBytes, &firstKeys); err != nil {
		t.Errorf("failed to unmarshal uri sig keys")
	}

	kabletownFirstKeys, ok := firstKeys["Kabletown URI Authority 1"]
	if !ok {
		t.Fatal("failed to create uri sig keys: 'Kabletown URI Authority 1' not found in response after creation")
	}
	if kabletownFirstKeys.Len() < 1 {
		t.Fatal("failed to create URI signing keys: 'Kabletown URI Authority 1' had zero keys after creation")
	}

	// Create new keys again and check that they are different
	var keyset2 tc.JWKSMap

	if err := json.Unmarshal([]byte(keySet2), &keyset2); err != nil {
		t.Errorf("json.UnMarshal(): expected nil error, actual: %v", err)
	}

	alerts, _, err := TOSession.CreateDeliveryServiceURISigningKeys(*firstDS.XMLID, keyset2, client.RequestOptions{})
	if err != nil {
		t.Errorf("Unexpected error creating URI Signature Keys for Delivery Service '%s': %v - alerts: %+v", *firstDS.XMLID, err, alerts.Alerts)
	}

	secondKeysBytes, _, err := TOSession.GetDeliveryServiceURISigningKeys(*firstDS.XMLID, client.RequestOptions{})
	if err != nil {
		t.Error("failed to get uri sig keys: " + err.Error())
	}

	secondKeys := tc.JWKSMap{}
	if err := json.Unmarshal(secondKeysBytes, &secondKeys); err != nil {
		t.Errorf("failed to unmarshal uri sig keys")
	}

	kabletownSecondKeys, ok := secondKeys["Kabletown URI Authority 1"]
	if !ok {
		t.Fatal("failed to create uri sig keys: 'Kabletown URI Authority 1' not found in response after creation")
	}
	if kabletownSecondKeys.Len() < 1 {
		t.Fatal("failed to create URI signing keys: 'Kabletown URI Authority 1' had zero keys after creation")
	}

	k1, ok := kabletownFirstKeys.Get(0)
	if !ok {
		t.Errorf(`failed to get key 0 from kabletownFirstKeys`)
	}
	k2, ok := kabletownSecondKeys.Get(0)
	if !ok {
		t.Errorf(`failed to get key 0 from kabletownSecondKeys`)
	}

	if k2.KeyID() == k1.KeyID() {
		t.Errorf("second create did not generate new uri sig keys - key mismatch")
	}
}

func DeleteTestDeliveryServicesURISigningKeys(t *testing.T) {
	if len(testData.DeliveryServices) == 0 {
		t.Fatal("couldn't get the xml ID of test DS")
	}
	firstDS := testData.DeliveryServices[0]
	if firstDS.XMLID == nil {
		t.Fatal("Found a Delivery Service in testing data with a null or undefined XMLID")
	}

	resp, _, err := TOSession.DeleteDeliveryServiceURISigningKeys(*firstDS.XMLID, client.RequestOptions{})
	if err != nil {
		t.Errorf("Unexpected error deleting URI Signing keys for Delivery Service '%s': %v - alerts: %+v", *firstDS.XMLID, err, resp.Alerts)
	}

	emptyBytes, _, err := TOSession.GetDeliveryServiceURISigningKeys(*firstDS.XMLID, client.RequestOptions{})
	if err != nil {
		t.Errorf("Unexpected error getting URI signing keys for Delivery Service '%s': %v", *firstDS.XMLID, err)
	}
	emptyMap := make(map[string]interface{})
	err = json.Unmarshal(emptyBytes, &emptyMap)
	if err != nil {
		t.Errorf("unexpected error unmarshalling empty URI signing keys response: %v", err)
	}
	renewalKid, hasRenewalKid := emptyMap["renewal_kid"]
	keys, hasKeys := emptyMap["keys"]
	if !hasRenewalKid {
		t.Error("getting empty URI signing keys - expected: 'renewal_kid' key, actual: not present")
	}
	if !hasKeys {
		t.Error("getting empty URI signing keys - expected: 'keys' key, actual: not present")
	}
	if renewalKid != nil {
		t.Errorf("getting empty URI signing keys - expected: 'renewal_kid' value to be nil, actual: %+v", renewalKid)
	}
	if keys != nil {
		t.Errorf("getting empty URI signing keys - expected: 'keys' value to be nil, actual: %+v", keys)
	}
}
