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
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
)

func TestDeliveryServicesKeys(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Users, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, ServerCapabilities, ServiceCategories, DeliveryServices}, func() {
		if !includeSystemTests {
			t.Skip()
		}
		SSLDeliveryServiceCDNUpdateTest(t)
		GetTestDeliveryServicesURLSigKeys(t)
	})
}

func createBlankCDN(cdnName string, t *testing.T) tc.CDN {
	_, _, err := TOSession.CreateCDN(tc.CDN{
		DNSSECEnabled: false,
		DomainName:    cdnName + ".ai",
		Name:          cdnName,
	})
	assert.RequireNoError(t, err, "Expected no error when creating cdn: %v", err)

	originalKeys, _, err := TOSession.GetCDNSSLKeysWithHdr(cdnName, nil)
	assert.RequireNoError(t, err, "Expected no error when getting cdn ssl keys: %v", err)

	cdns, _, err := TOSession.GetCDNByNameWithHdr(cdnName, nil)
	assert.RequireNoError(t, err, "Unable to get cdn: %v", err)
	assert.RequireGreaterOrEqual(t, len(cdns), 1, "Expected more than 0 cdns")

	keys, _, err := TOSession.GetCDNSSLKeysWithHdr(cdnName, nil)
	assert.RequireNoError(t, err, "Expected no error when getting cdn ssl keys: %v", err)
	assert.RequireEqual(t, len(keys), len(originalKeys), "Expected %v ssl keys on cdn %v, got %v", len(originalKeys), cdnName, len(keys))

	return cdns[0]
}

func cleanUp(t *testing.T, ds tc.DeliveryServiceNullableV30, oldCDNID int, newCDNID int, sslKeyVersions []string) {
	_, _, err := TOSession.DeleteDeliveryServiceSSLKeysByID(*ds.XMLID)
	assert.NoError(t, err, "Expected no error when cleaning up delivery service ssl keys.")

	params := url.Values{}
	for _, version := range sslKeyVersions {
		params.Set("version", version)
		_, _, err := TOSession.DeleteDeliveryServiceSSLKeysByVersion(*ds.XMLID, params)
		assert.NoError(t, err, "Expected no error when cleaning up delivery service ssl keys by versions.")
	}
	_, err = TOSession.DeleteDeliveryService(strconv.Itoa(*ds.ID))
	assert.NoError(t, err, "Expected no error when cleaning up delivery services.")

	_, _, err = TOSession.DeleteCDNByID(oldCDNID)
	assert.NoError(t, err, "Expected no error when cleaning up cdns.")

	_, _, err = TOSession.DeleteCDNByID(newCDNID)
	assert.NoError(t, err, "Expected no error when cleaning up cdns.")
}

func SSLDeliveryServiceCDNUpdateTest(t *testing.T) {
	cdnNameOld := "sslkeytransfer"
	oldCdn := createBlankCDN(cdnNameOld, t)
	cdnNameNew := "sslkeytransfer1"
	newCdn := createBlankCDN(cdnNameNew, t)

	types, _, err := TOSession.GetTypeByNameWithHdr("HTTP", nil)
	assert.RequireNoError(t, err, "Unable to get types: %v", err)
	assert.RequireGreaterOrEqual(t, len(types), 1, "Expected at least one type.")

	customDS := tc.DeliveryServiceNullableV30{}
	customDS.Active = util.BoolPtr(true)
	customDS.CDNID = util.IntPtr(oldCdn.ID)
	customDS.DSCP = util.IntPtr(0)
	customDS.DisplayName = util.StrPtr("displayName")
	customDS.RoutingName = util.StrPtr("routingName")
	customDS.GeoLimit = util.IntPtr(0)
	customDS.GeoProvider = util.IntPtr(0)
	customDS.IPV6RoutingEnabled = util.BoolPtr(false)
	customDS.InitialDispersion = util.IntPtr(1)
	customDS.LogsEnabled = util.BoolPtr(true)
	customDS.MissLat = util.FloatPtr(0)
	customDS.MissLong = util.FloatPtr(0)
	customDS.MultiSiteOrigin = util.BoolPtr(false)
	customDS.OrgServerFQDN = util.StrPtr("https://test.com")
	customDS.Protocol = util.IntPtr(2)
	customDS.QStringIgnore = util.IntPtr(0)
	customDS.RangeRequestHandling = util.IntPtr(0)
	customDS.RegionalGeoBlocking = util.BoolPtr(false)
	customDS.TenantID = util.IntPtr(1)
	customDS.TypeID = util.IntPtr(types[0].ID)
	customDS.XMLID = util.StrPtr("dsID")
	customDS.MaxRequestHeaderBytes = nil

	ds, _, err := TOSession.CreateDeliveryServiceV30(customDS)
	assert.RequireNoError(t, err, "Unable to create delivery service: %v", err)

	ds.CDNName = &oldCdn.Name

	defer cleanUp(t, ds, oldCdn.ID, newCdn.ID, []string{"1"})

	_, _, err = TOSession.GenerateSSLKeysForDS(*ds.XMLID, *ds.CDNName, tc.SSLKeyRequestFields{
		BusinessUnit: util.StrPtr("BU"),
		City:         util.StrPtr("CI"),
		Organization: util.StrPtr("OR"),
		HostName:     util.StrPtr("*.test.com"),
		Country:      util.StrPtr("CO"),
		State:        util.StrPtr("ST"),
	})
	assert.RequireNoError(t, err, "Unable to generate sslkeys for cdn %v: %v", oldCdn.Name, err)

	tries := 0
	var oldCDNKeys []tc.CDNSSLKeys
	for tries < 5 {
		time.Sleep(time.Second)
		oldCDNKeys, _, err = TOSession.GetCDNSSLKeysWithHdr(oldCdn.Name, nil)
		if err == nil && len(oldCDNKeys) > 0 {
			break
		}
		tries++
	}
	assert.RequireNoError(t, err, "Unable to get cdn %v keys: %v", oldCdn.Name, err)
	assert.RequireGreaterOrEqual(t, len(oldCDNKeys), 1, "Expected at least one key.")

	newCDNKeys, _, err := TOSession.GetCDNSSLKeysWithHdr(newCdn.Name, nil)
	assert.RequireNoError(t, err, "Unable to get cdn %v keys: %v", newCdn.Name, err)

	ds.RoutingName = util.StrPtr("anothername")
	_, _, err = TOSession.UpdateDeliveryServiceV30WithHdr(*ds.ID, ds, nil)
	assert.Error(t, err, "Should not be able to update delivery service (routing name) as it has ssl keys")

	ds.RoutingName = util.StrPtr("routingName")

	ds.CDNID = &newCdn.ID
	ds.CDNName = &newCdn.Name
	_, _, err = TOSession.UpdateDeliveryServiceV30WithHdr(*ds.ID, ds, nil)
	assert.Error(t, err, "Should not be able to update delivery service (cdn) as it has ssl keys")

	// Check new CDN still has an ssl key
	keys, _, err := TOSession.GetCDNSSLKeysWithHdr(newCdn.Name, nil)
	assert.RequireNoError(t, err, "Unable to get cdn %v keys: %v", newCdn.Name, err)
	assert.Equal(t, len(newCDNKeys), len(keys), "Expected %v keys, got %v", len(newCDNKeys), len(keys))

	// Check old CDN does not have ssl key
	keys, _, err = TOSession.GetCDNSSLKeysWithHdr(oldCdn.Name, nil)
	assert.RequireNoError(t, err, "Unable to get cdn %v keys: %v", oldCdn.Name, err)
	assert.Equal(t, len(oldCDNKeys), len(keys), "Expected %v key, got %v", len(oldCDNKeys), len(keys))
}

func GetTestDeliveryServicesURLSigKeys(t *testing.T) {
	if len(testData.DeliveryServices) == 0 {
		t.Fatal("couldn't get the xml ID of test DS")
	}
	firstDS := testData.DeliveryServices[0]
	if firstDS.XMLID == nil {
		t.Fatal("couldn't get the xml ID of test DS")
	}

	_, _, err := TOSession.GetDeliveryServiceURLSigKeysWithHdr(*firstDS.XMLID, nil)
	if err != nil {
		t.Error("failed to get url sig keys: " + err.Error())
	}
}
