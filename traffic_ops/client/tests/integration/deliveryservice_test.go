package integration

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	traffic_ops "github.com/apache/incubator-trafficcontrol/traffic_ops/client"
)

var (
	testDs           traffic_ops.DeliveryService
	testDsID         string
	existingTestDS   traffic_ops.DeliveryService
	existingTestDSID string
	sslDs            traffic_ops.DeliveryService
)

func init() {
	cdn, err := GetCdn()
	if err != nil {
		fmt.Printf("Deliverservice_test init -- Could not get CDNs from TO...%v\n", err)
		os.Exit(1)
	}

	profile, err := GetProfile()
	if err != nil {
		fmt.Printf("Deliverservice_test init -- Could not get Profiles from TO...%v\n", err)
		os.Exit(1)
	}

	dsType, err := GetType("deliveryservice")
	if err != nil {
		fmt.Printf("Deliverservice_test init -- Could not get Types from TO...%v\n", err)
		os.Exit(1)
	}

	//create DeliveryService object for testing
	testDs.Active = false
	testDs.CCRDNSTTL = 30
	testDs.CDNName = cdn.Name
	testDs.CDNID = cdn.ID
	testDs.CacheURL = "cacheURL"
	testDs.CheckPath = "CheckPath"
	testDs.DNSBypassCname = "DNSBypassCNAME"
	testDs.DNSBypassIP = "10.10.10.10"
	testDs.DNSBypassIP6 = "FF01:0:0:0:0:0:0:FB"
	testDs.DNSBypassTTL = 30
	testDs.DSCP = 0
	testDs.DisplayName = "DisplayName"
	testDs.EdgeHeaderRewrite = "EdgeHeaderRewrite"
	testDs.GeoLimit = 5
	testDs.GeoProvider = 1
	testDs.GlobalMaxMBPS = 15000
	testDs.GlobalMaxTPS = 15000
	testDs.HTTPBypassFQDN = "HTTPBypassFQDN"
	testDs.IPV6RoutingEnabled = true
	testDs.InfoURL = "InfoUrl"
	testDs.InitialDispersion = 5
	testDs.LongDesc = "LongDesc"
	testDs.LongDesc1 = "LongDesc1"
	testDs.LongDesc2 = "LongDesc2"
	testDs.MaxDNSAnswers = 5
	testDs.MidHeaderRewrite = "MidHeaderRewrite"
	testDs.MissLat = 5.555
	testDs.MissLong = -50.5050
	testDs.MultiSiteOrigin = true
	testDs.OrgServerFQDN = "http://OrgServerFQDN"
	testDs.ProfileDesc = profile.Description
	testDs.ProfileName = profile.Name
	testDs.ProfileID = profile.ID
	testDs.Protocol = 1
	testDs.QStringIgnore = 1
	testDs.RangeRequestHandling = 0
	testDs.RegexRemap = "regexRemap"
	testDs.RemapText = "remapText"
	testDs.Signed = false
	testDs.TRResponseHeaders = "TRResponseHeaders"
	testDs.Type = dsType.Name
	testDs.TypeID = dsType.ID
	testDs.XMLID = "Test-DS-" + strconv.FormatInt(time.Now().Unix(), 10)
	testDs.RegionalGeoBlocking = false
	testDs.LogsEnabled = false

	//Create method currently does not support MatchList...
	// testDsMatch1 := new(traffic_ops.DeliveryServiceMatch)
	// testDsMatch1.Pattern = "Pattern1"
	// testDsMatch1.SetNumber = "0"
	// testDsMatch1.Type = "HOST"

	// testDsMatch2 := new(traffic_ops.DeliveryServiceMatch)
	// testDsMatch2.Pattern = "Pattern2"
	// testDsMatch2.SetNumber = "1"
	// testDsMatch2.Type = "HOST"

	// testDs.MatchList = append(testDs.MatchList, *testDsMatch1)
	// testDs.MatchList = append(testDs.MatchList, *testDsMatch2)

}

func TestDeliveryServices(t *testing.T) {
	resp, err := Request(*to, "GET", "/api/1.2/deliveryservices.json", nil)
	if err != nil {
		t.Errorf("Could not get deliveryservices.json reponse was: %v\n", err)
	}

	defer resp.Body.Close()
	var apiDsRes traffic_ops.GetDeliveryServiceResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiDsRes); err != nil {
		t.Errorf("Could not decode Deliveryservice json.  Error is: %v\n", err)
	}
	apiDss := apiDsRes.Response

	clientDss, err := to.DeliveryServices()
	if err != nil {
		t.Errorf("Could not get Deliveryservices from client.  Error is: %v\n", err)
	}

	if len(apiDss) != len(clientDss) {
		t.Errorf("Array lengths from client and API are different...API = %v, Client = %v\n", apiDss, clientDss)
	}

	matchFound := false
	for _, apiDs := range apiDss {
		//set these to use later...this saves time over doing it in the init() method
		if apiDs.Protocol == 0 && existingTestDS.ID == 0 {
			existingTestDS = apiDs
			existingTestDSID = strconv.Itoa(existingTestDS.ID)
		}
		if apiDs.Protocol > 0 && strings.Contains(apiDs.Type, "DNS") && sslDs.ID == 0 {
			sslDs = apiDs
		}

		for _, clientDs := range clientDss {
			if clientDs.XMLID != apiDs.XMLID {
				continue
			}
			matchFound = true
			compareDs(apiDs, clientDs, t)
		}
		if !matchFound {
			t.Errorf("A match for %s from the API was not found in the client results\n", apiDs.XMLID)
		}
	}
}

func TestCreateDs(t *testing.T) {
	//create a DS and validate response
	res, err := to.CreateDeliveryService(&testDs)
	if err != nil {
		t.Error("Failed to create deliveryservice!  Error is: ", err)
	} else {
		testDs.ID = res.Response[0].ID
		testDsID = strconv.Itoa(testDs.ID)
		compareDs(testDs, res.Response[0], t)
	}
}

func TestUpdateDs(t *testing.T) {
	testDs.DisplayName = "New Display Name"
	testDs.LongDesc += "-- Update"
	testDs.LongDesc1 += "-- Update"
	testDs.LongDesc2 += "-- Update"
	testDs.EdgeHeaderRewrite += "-- Update"
	res, err := to.UpdateDeliveryService(testDsID, &testDs)
	if err != nil {
		t.Error("Failed to update deliveryservice!  Error is: ", err)
	} else {
		compareDs(testDs, res.Response[0], t)
	}
}

func TestDeliveryService(t *testing.T) {
	uri := fmt.Sprintf("/api/1.2/deliveryservices/%s.json", testDsID)
	resp, err := Request(*to, "GET", uri, nil)
	if err != nil {
		t.Errorf("Could not get %s reponse was: %v\n", uri, err)
	}

	defer resp.Body.Close()
	var apiDsRes traffic_ops.GetDeliveryServiceResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiDsRes); err != nil {
		t.Errorf("Could not decode Deliveryservice json.  Error is: %v\n", err)
	}

	clientDss, err := to.DeliveryService(testDsID)
	if err != nil {
		t.Errorf("Could not get Deliveryservice from client.  Error is: %v\n", err)
	}

	compareDs(apiDsRes.Response[0], *clientDss, t)
}

//Put this Test after anything using the testDS or testDsID variables
func TestDeleteDeliveryService(t *testing.T) {
	res, err := to.DeleteDeliveryService(testDsID)
	if err != nil {
		t.Errorf("Could not delete Deliveryserivce %s reponse was: %v\n", testDsID, err)
	}
	if res.Alerts[0].Level != "success" {
		t.Errorf("Alert.Level -- Expected \"success\" got %s", res.Alerts[0].Level)
	}
	if res.Alerts[0].Text != "Delivery service was deleted." {
		t.Errorf("Alert.Level -- Expected \"Delivery service was deleted.\" got %s", res.Alerts[0].Text)
	}
}

func TestDeliveryServiceState(t *testing.T) {
	uri := fmt.Sprintf("/api/1.2/deliveryservices/%s/state.json", existingTestDSID)
	resp, err := Request(*to, "GET", uri, nil)
	if err != nil {
		t.Errorf("Could not get %s reponse was: %v\n", uri, err)
	}

	defer resp.Body.Close()
	var apiDsStateRes traffic_ops.DeliveryServiceStateResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiDsStateRes); err != nil {
		t.Errorf("Could not decode DeliveryserviceState reponse.  Error is: %v\n", err)
	}

	apiDsState := apiDsStateRes.Response

	clientDsState, err := to.DeliveryServiceState(existingTestDSID)
	if err != nil {
		t.Errorf("Could not get DS State from client for %s reponse was: %v\n", existingTestDSID, err)
	}

	if apiDsState.Enabled != clientDsState.Enabled {
		t.Errorf("Enabled -- Expected %v got %v for ID %s", apiDsState.Enabled, clientDsState.Enabled, existingTestDSID)
	}
	if apiDsState.Failover.Configured != clientDsState.Failover.Configured {
		t.Errorf("Failover.Configured -- Expected %v got %v", apiDsState.Failover.Configured, clientDsState.Failover.Configured)
	}
	if apiDsState.Failover.Destination.Location != clientDsState.Failover.Destination.Location {
		t.Errorf("Failover.Destination.Location -- Expected %v got %v", apiDsState.Failover.Destination.Location, clientDsState.Failover.Destination.Location)
	}
	if apiDsState.Failover.Destination.Type != clientDsState.Failover.Destination.Type {
		t.Errorf("Failover.Destination.Type -- Expected %v got %v", apiDsState.Failover.Destination.Type, clientDsState.Failover.Destination.Type)
	}
	if apiDsState.Failover.Enabled != clientDsState.Failover.Enabled {
		t.Errorf("res.Failover.Enabled -- Expected %v got %v", apiDsState.Failover.Enabled, clientDsState.Failover.Enabled)
	}
	if len(apiDsState.Failover.Locations) != len(clientDsState.Failover.Locations) {
		t.Errorf("res.Failover.Locations len -- Expected %v got %v", len(apiDsState.Failover.Locations), len(clientDsState.Failover.Locations))
	}
}

func TestDeliveryServiceHealth(t *testing.T) {
	uri := fmt.Sprintf("/api/1.2/deliveryservices/%s/health.json", existingTestDSID)
	resp, err := Request(*to, "GET", uri, nil)
	if err != nil {
		t.Errorf("Could not get %s reponse was: %v\n", uri, err)
	}

	defer resp.Body.Close()
	var apiDsHealthRes traffic_ops.DeliveryServiceHealthResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiDsHealthRes); err != nil {
		t.Errorf("Could not decode DeliveryserviceHealth reponse.  Error is: %v\n", err)
	}

	apiDsHealth := apiDsHealthRes.Response

	clientDsHealth, err := to.DeliveryServiceHealth(existingTestDSID)
	if err != nil {
		t.Errorf("Could not ge Deliveryserivce Health for %s reponse was: %v\n", existingTestDSID, err)
	}

	if apiDsHealth.TotalOnline != clientDsHealth.TotalOnline {
		t.Errorf("TotalOnline -- Expected %v got %v", apiDsHealth.TotalOnline, apiDsHealth.TotalOnline)
	}

	if apiDsHealth.TotalOffline != clientDsHealth.TotalOffline {
		t.Errorf("TotalOffline -- Expected %v got %v", apiDsHealth.TotalOffline, clientDsHealth.TotalOffline)
	}

	if len(apiDsHealth.CacheGroups) != len(clientDsHealth.CacheGroups) {
		t.Errorf("len Cachegroups -- Expected %v got %v", len(apiDsHealth.CacheGroups), len(clientDsHealth.CacheGroups))
	}

	for _, apiCg := range apiDsHealth.CacheGroups {
		match := false
		for _, clientCg := range clientDsHealth.CacheGroups {
			if apiCg.Name != clientCg.Name {
				continue
			}
			match = true
			if apiCg.Offline != clientCg.Offline {
				t.Errorf("Cachegroup.Offline -- Expected %v got %v", apiCg.Offline, clientCg.Offline)
			}
			if apiCg.Online != clientCg.Online {
				t.Errorf("Cachegroup.Online -- Expected %v got %v", apiCg.Online, clientCg.Online)
			}
		}
		if !match {
			t.Errorf("Cachegroup -- No match from client for api cachgroup %v", apiCg.Name)
		}
	}
}

func TestDeliveryServiceCapacity(t *testing.T) {
	uri := fmt.Sprintf("/api/1.2/deliveryservices/%s/capacity.json", existingTestDSID)
	resp, err := Request(*to, "GET", uri, nil)
	if err != nil {
		t.Errorf("Could not get %s reponse was: %v\n", uri, err)
	}

	defer resp.Body.Close()
	var apiDsCapacityRes traffic_ops.DeliveryServiceCapacityResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiDsCapacityRes); err != nil {
		t.Errorf("Could not decode DeliveryserviceCapacity reponse.  Error is: %v\n", err)
	}

	apiDsCapacity := apiDsCapacityRes.Response

	clientDsCapacity, err := to.DeliveryServiceCapacity(existingTestDSID)
	if err != nil {
		t.Errorf("Could not ge Deliveryserivce Capacity for %s reponse was: %v\n", existingTestDSID, err)
	}

	if fmt.Sprintf("%6.5f", apiDsCapacity.AvailablePercent) != fmt.Sprintf("%6.5f", clientDsCapacity.AvailablePercent) {
		t.Errorf("AvailablePercent -- Expected %v got %v", apiDsCapacity.AvailablePercent, clientDsCapacity.AvailablePercent)
	}

	if fmt.Sprintf("%6.5f", apiDsCapacity.MaintenancePercent) != fmt.Sprintf("%6.5f", clientDsCapacity.MaintenancePercent) {
		t.Errorf("MaintenenancePercent -- Expected %v got %v", apiDsCapacity.MaintenancePercent, clientDsCapacity.MaintenancePercent)
	}

	if fmt.Sprintf("%6.5f", apiDsCapacity.UnavailablePercent) != fmt.Sprintf("%6.5f", clientDsCapacity.UnavailablePercent) {
		t.Errorf("UnavailablePercent -- Expected %v got %v", apiDsCapacity.UnavailablePercent, clientDsCapacity.UnavailablePercent)
	}

	if fmt.Sprintf("%6.5f", apiDsCapacity.UtilizedPercent) != fmt.Sprintf("%6.5f", clientDsCapacity.UtilizedPercent) {
		t.Errorf("UtilizedPercent -- Expected %v got %v", apiDsCapacity.UtilizedPercent, clientDsCapacity.UtilizedPercent)
	}

}

func TestDeliveryServiceRouting(t *testing.T) {
	uri := fmt.Sprintf("/api/1.2/deliveryservices/%s/routing.json", existingTestDSID)
	resp, err := Request(*to, "GET", uri, nil)
	if err != nil {
		t.Errorf("Could not get %s reponse was: %v\n", uri, err)
	}

	defer resp.Body.Close()
	var apiDsRoutingRes traffic_ops.DeliveryServiceRoutingResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiDsRoutingRes); err != nil {
		t.Errorf("Could not decode DeliveryserviceRouting reponse.  Error is: %v\n", err)
	}

	apiDsRouting := apiDsRoutingRes.Response

	clientDsRouting, err := to.DeliveryServiceRouting(existingTestDSID)
	if err != nil {
		t.Errorf("Could not ge Deliveryserivce Routing for %s reponse was: %v\n", existingTestDSID, err)
	}

	if apiDsRouting.CZ != clientDsRouting.CZ {
		t.Errorf("CZ -- Expected %v got %v", apiDsRouting.CZ, clientDsRouting.CZ)
	}

	if apiDsRouting.DSR != clientDsRouting.DSR {
		t.Errorf("DSR -- Expected %v got %v", apiDsRouting.DSR, clientDsRouting.DSR)
	}

	if apiDsRouting.Err != clientDsRouting.Err {
		t.Errorf("Err-- Expected %v got %v", apiDsRouting.Err, clientDsRouting.Err)
	}

	if apiDsRouting.Fed != clientDsRouting.Fed {
		t.Errorf("Fed -- Expected %v got %v", apiDsRouting.Fed, clientDsRouting.Fed)
	}

	if apiDsRouting.Geo != clientDsRouting.Geo {
		t.Errorf("Geo -- Expected %v got %v", apiDsRouting.Geo, clientDsRouting.Geo)
	}

	if apiDsRouting.Miss != clientDsRouting.Miss {
		t.Errorf("Miss -- Expected %v got %v", apiDsRouting.Miss, clientDsRouting.Miss)
	}

	if apiDsRouting.RegionalAlternate != clientDsRouting.RegionalAlternate {
		t.Errorf("RegionalAlternate -- Expected %v got %v", apiDsRouting.RegionalAlternate, clientDsRouting.RegionalAlternate)
	}

	if apiDsRouting.RegionalDenied != clientDsRouting.RegionalDenied {
		t.Errorf("RegionalDenied -- Expected %v got %v", apiDsRouting.RegionalDenied, clientDsRouting.RegionalDenied)
	}
}

func TestDeliveryServiceServer(t *testing.T) {
	resp, err := Request(*to, "GET", "/api/1.2/deliveryserviceserver.json?page=1&limit=1", nil)
	if err != nil {
		t.Errorf("Could not get deliveryserviceserver.json reponse was: %v\n", err)
	}

	defer resp.Body.Close()
	var apiDsServerRes traffic_ops.DeliveryServiceServerResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiDsServerRes); err != nil {
		t.Errorf("Could not decode DeliveryserviceServer reponse.  Error is: %v\n", err)
	}

	clientDsServerRes, err := to.DeliveryServiceServer("1", "1")

	if err != nil {
		t.Errorf("Could not get DeliveryserviceServer, reponse was: %v\n", err)
	}

	for _, apiDss := range apiDsServerRes.Response {
		match := false
		for _, clientDss := range clientDsServerRes {
			if clientDss.DeliveryService != apiDss.DeliveryService {
				continue
			}
			match = true
			if apiDss.LastUpdated != clientDss.LastUpdated {
				t.Errorf("LastUpdated -- Expected %v got %v", apiDss.LastUpdated, clientDss.LastUpdated)
			}
			if apiDss.Server != clientDss.Server {
				t.Errorf("Server -- Expected %v got %v", apiDss.Server, clientDss.Server)
			}
		}
		if match != true {
			t.Errorf("No match found for the Deliveryservice %v in DeliveryserviceServer response: %v\n", apiDss.DeliveryService, err)
		}
	}

}

func TestDeliveryServiceSSLKeysByID(t *testing.T) {
	if sslDs.ID > 0 {
		uri := fmt.Sprintf("/api/1.2/deliveryservices/xmlId/%s/sslkeys.json", sslDs.XMLID)
		resp, err := Request(*to, "GET", uri, nil)
		if err != nil {
			t.Errorf("Could not get %s reponse was: %v\n", uri, err)
		}

		defer resp.Body.Close()
		var apiSslRes traffic_ops.DeliveryServiceSSLKeysResponse
		if err := json.NewDecoder(resp.Body).Decode(&apiSslRes); err != nil {
			t.Errorf("Could not decode DeliveryServiceSSLKeysResponse reponse.  Error is: %v\n", err)
		}

		clientSslRes, err := to.DeliveryServiceSSLKeysByID(sslDs.XMLID)

		if err != nil {
			t.Errorf("Could not get DeliveryserviceSSLKeys, reponse was: %v\n", err)
		}
		compareSSLResponse(apiSslRes.Response, *clientSslRes, t)
	} else {
		t.Skip("Skipping TestDeliveryServiceSSLKeysByID because no Deliveryservice was found with SSL enabled")
	}
}

func TestDeliveryServiceSSLKeysByHostname(t *testing.T) {
	if sslDs.ID > 0 {
		var hostname string
		for _, exampleURL := range sslDs.ExampleURLs {
			if strings.Contains(exampleURL, "edge.") {
				u, err := url.Parse(exampleURL)
				if err != nil {
					t.Errorf("could not parse exampleURL %s\n", exampleURL)
					t.FailNow()
				}
				hostname = u.Host
			}
		}
		if hostname == "" {
			t.Skipf("could not find an example URL from Deliveryservice %s to use for testing\n", sslDs.XMLID)
			t.SkipNow()
		}

		uri := fmt.Sprintf("/api/1.2/deliveryservices/hostname/%s/sslkeys.json", hostname)
		resp, err := Request(*to, "GET", uri, nil)
		if err != nil {
			t.Errorf("Could not get %s reponse was: %v\n", uri, err)
		}

		defer resp.Body.Close()
		var apiSslRes traffic_ops.DeliveryServiceSSLKeysResponse
		if err := json.NewDecoder(resp.Body).Decode(&apiSslRes); err != nil {
			t.Errorf("Could not decode DeliveryServiceSSLKeysResponse reponse.  Error is: %v\n", err)
		}

		clientSslRes, err := to.DeliveryServiceSSLKeysByHostname(hostname)

		if err != nil {
			t.Errorf("Could not get DeliveryserviceSSLKeys, reponse was: %v\n", err)
		}
		compareSSLResponse(apiSslRes.Response, *clientSslRes, t)
	} else {
		t.Skip("Skipping TestDeliveryServiceSSLKeysByID because no Deliveryservice was found with SSL enabled")
	}
}

func compareDs(ds1 traffic_ops.DeliveryService, ds2 traffic_ops.DeliveryService, t *testing.T) {
	if ds1.Active != ds1.Active {
		t.Errorf("Active -- Expected %v, Got %v\n", ds1.Active, ds2.Active)
	}
	if ds1.CCRDNSTTL != ds2.CCRDNSTTL {
		t.Errorf("CCRDNSTTL -- Expected %v, Got %v\n", ds1.CCRDNSTTL, ds2.CCRDNSTTL)
	}
	if ds1.CDNName != ds2.CDNName {
		t.Errorf("CDNName -- Expected %v, Got %v\n", ds1.CDNName, ds2.CDNName)
	}
	if ds1.CDNID != ds2.CDNID {
		t.Errorf("CDNID -- Expected %v, Got %v\n", ds1.CDNID, ds2.CDNID)
	}
	if ds1.CacheURL != ds2.CacheURL {
		t.Errorf("CacheURL -- Expected %v, Got %v\n", ds1.CacheURL, ds2.CacheURL)
	}
	if ds1.CheckPath != ds2.CheckPath {
		t.Errorf("CheckPath -- Expected %v, Got %v\n", ds1.CheckPath, ds2.CheckPath)
	}
	if ds1.DNSBypassCname != ds2.DNSBypassCname {
		t.Errorf("DNSBypassCname -- Expected %v, Got %v\n", ds1.DNSBypassCname, ds2.DNSBypassCname)
	}
	if ds1.DNSBypassIP != ds2.DNSBypassIP {
		t.Errorf("DNSBypassIP -- Expected %v, Got %v\n", ds1.DNSBypassIP, ds2.DNSBypassIP)
	}
	if ds1.DNSBypassIP6 != ds2.DNSBypassIP6 {
		t.Errorf("DNSBypassIP6 -- Expected %v, Got %v\n", ds1.DNSBypassIP6, ds2.DNSBypassIP6)
	}
	if ds1.DNSBypassTTL != ds2.DNSBypassTTL {
		t.Errorf("DNSBypassTTL -- Expected %v, Got %v\n", ds1.DNSBypassTTL, ds2.DNSBypassTTL)
	}
	if ds1.DSCP != ds2.DSCP {
		t.Errorf("DSCP -- Expected %v, Got %v\n", ds1.DSCP, ds2.DSCP)
	}
	if ds1.DisplayName != ds2.DisplayName {
		t.Errorf("DisplayName -- Expected %v, Got %v\n", ds1.DisplayName, ds2.DisplayName)
	}
	if ds1.EdgeHeaderRewrite != ds2.EdgeHeaderRewrite {
		t.Errorf("EdgeHeaderRewrite -- Expected %v, Got %v\n", ds1.EdgeHeaderRewrite, ds2.EdgeHeaderRewrite)
	}
	if ds1.GeoLimit != ds2.GeoLimit {
		t.Errorf("GeoLimit -- Expected %v, Got %v\n", ds1.GeoLimit, ds2.GeoLimit)
	}
	if ds1.GeoProvider != ds2.GeoProvider {
		t.Errorf("GeoProvider -- Expected %v, Got %v\n", ds1.GeoProvider, ds2.GeoProvider)
	}
	if ds1.GlobalMaxMBPS != ds2.GlobalMaxMBPS {
		t.Errorf("GlobalMaxMBPS -- Expected %v, Got %v\n", ds1.GlobalMaxMBPS, ds2.GlobalMaxMBPS)
	}
	if ds1.GlobalMaxTPS != ds2.GlobalMaxTPS {
		t.Errorf("GlobalMaxTPS -- Expected %v, Got %v\n", ds1.GlobalMaxTPS, ds2.GlobalMaxTPS)
	}
	if ds1.HTTPBypassFQDN != ds2.HTTPBypassFQDN {
		t.Errorf("HTTPBypassFQDN -- Expected %v, Got %v\n", ds1.HTTPBypassFQDN, ds2.HTTPBypassFQDN)
	}
	if ds1.ID > 0 && ds1.ID != ds2.ID {
		t.Errorf("ID -- Expected %v, Got %v\n", ds1.ID, ds2.ID)
	}
	if ds1.IPV6RoutingEnabled != ds2.IPV6RoutingEnabled {
		t.Errorf("IPv6RoutingEnabled -- Expected %v, Got %v\n", ds1.IPV6RoutingEnabled, ds2.IPV6RoutingEnabled)
	}
	if ds1.InfoURL != ds2.InfoURL {
		t.Errorf("InfoURL -- Expected %v, Got %v\n", ds1.InfoURL, ds2.InfoURL)
	}
	if ds1.InitialDispersion != ds2.InitialDispersion {
		t.Errorf("InitialDispersion -- Expected %v, Got %v\n", ds1.InitialDispersion, ds2.InitialDispersion)
	}
	if ds1.LastUpdated != "" && ds1.LastUpdated != ds2.LastUpdated {
		t.Errorf("LastUpdated -- Expected %v, Got %v\n", ds1.LastUpdated, ds2.LastUpdated)
	}
	if ds1.LongDesc != ds2.LongDesc {
		t.Errorf("LongDesc -- Expected %v, Got %v\n", ds1.LongDesc, ds2.LongDesc)
	}
	if ds1.LongDesc1 != ds2.LongDesc1 {
		t.Errorf("LongDesc1 -- Expected %v, Got %v\n", ds1.LongDesc1, ds2.LongDesc1)
	}
	if ds1.LongDesc2 != ds2.LongDesc2 {
		t.Errorf("LongDesc2 -- Expected %v, Got %v\n", ds1.LongDesc2, ds2.LongDesc2)
	}
	if ds1.MaxDNSAnswers != ds2.MaxDNSAnswers {
		t.Errorf("MaxDNSAnswers-- Expected %v, Got %v\n", ds1.MaxDNSAnswers, ds2.MaxDNSAnswers)
	}
	if ds1.MidHeaderRewrite != ds2.MidHeaderRewrite {
		t.Errorf("MidHeaderRewrite -- Expected %v, Got %v\n", ds1.MidHeaderRewrite, ds2.MidHeaderRewrite)
	}
	if ds1.MissLat != ds2.MissLat {
		t.Errorf("MissLat -- Expected %v, Got %v\n", ds1.MissLat, ds2.MissLat)
	}
	if ds1.MissLong != ds2.MissLong {
		t.Errorf("MissLong -- Expected %v, Got %v\n", ds1.MissLong, ds2.MissLong)
	}
	if ds1.MultiSiteOrigin != ds2.MultiSiteOrigin {
		t.Errorf("MutiSiteOrigin -- Expected %v, Got %v\n", ds1.MultiSiteOrigin, ds2.MultiSiteOrigin)
	}
	if ds1.OrgServerFQDN != ds2.OrgServerFQDN {
		t.Errorf("OrgServerFQDN -- Expected %v, Got %v\n", ds1.OrgServerFQDN, ds2.OrgServerFQDN)
	}
	if ds1.ProfileDesc != ds2.ProfileDesc {
		t.Errorf("ProfileDesc -- Expected %v, Got %v\n", ds1.ProfileDesc, ds2.ProfileDesc)
	}
	if ds1.ProfileName != ds2.ProfileName {
		t.Errorf("ProfileName -- Expected %v, Got %v\n", ds1.ProfileName, ds2.ProfileName)
	}
	if ds1.Protocol != ds2.Protocol {
		t.Errorf("Protocol -- Expected %v, Got %v\n", ds1.Protocol, ds2.Protocol)
	}
	if ds1.QStringIgnore != ds2.QStringIgnore {
		t.Errorf("QStringIgnore -- Expected %v, Got %v\n", ds1.QStringIgnore, ds2.QStringIgnore)
	}
	if ds1.RangeRequestHandling != ds2.RangeRequestHandling {
		t.Errorf("RangeRequestHandling -- Expected %v, Got %v\n", ds1.RangeRequestHandling, ds2.RangeRequestHandling)
	}
	if ds1.RegexRemap != ds2.RegexRemap {
		t.Errorf("RegexRemap -- Expected %v, Got %v\n", ds1.RegexRemap, ds2.RegexRemap)
	}
	if ds1.RemapText != ds2.RemapText {
		t.Errorf("RemapText -- Expected %v, Got %v\n", ds1.RemapText, ds2.RemapText)
	}
	if ds1.Signed != ds2.Signed {
		t.Errorf("Signed -- Expected %v, Got %v\n", ds1.Signed, ds2.Signed)
	}
	if ds1.TRResponseHeaders != ds2.TRResponseHeaders {
		t.Errorf("TRResponseHeaders -- Expected %v, Got %v\n", ds1.TRResponseHeaders, ds2.TRResponseHeaders)
	}
	if ds1.Type != ds2.Type {
		t.Errorf("Type -- Expected %v, Got %v\n", ds1.Type, ds2.Type)
	}
	if ds1.RegionalGeoBlocking != ds2.RegionalGeoBlocking {
		t.Errorf("RegionalGeoBlocking -- Expected %v, Got %v\n", ds1.RegionalGeoBlocking, ds2.RegionalGeoBlocking)
	}
	if ds1.LogsEnabled != ds2.LogsEnabled {
		t.Errorf("LogsEnabled -- Expected %v, Got %v\n", ds1.LogsEnabled, ds2.LogsEnabled)
	}
	if len(ds1.MatchList) > 0 {
		for i, match := range ds1.MatchList {
			if match != ds2.MatchList[i] {
				t.Errorf("Matchlist -- Expected %v, Got %v\n", match, ds2.MatchList[i])
			}
		}
	}
	if len(ds1.ExampleURLs) > 0 {
		for i, url := range ds1.ExampleURLs {
			if url != ds2.ExampleURLs[i] {
				t.Errorf("ExampleURL -- Expected %v, Got %v\n", url, ds2.ExampleURLs[i])
			}
		}
	}
}

func compareSSLResponse(apiSslRes traffic_ops.DeliveryServiceSSLKeys, clientSslRes traffic_ops.DeliveryServiceSSLKeys, t *testing.T) {
	if apiSslRes.BusinessUnit != clientSslRes.BusinessUnit {
		t.Errorf("BusinessUnit -- Expected %v got %v", apiSslRes.BusinessUnit, clientSslRes.BusinessUnit)
	}
	if apiSslRes.CDN != clientSslRes.CDN {
		t.Errorf("CDN -- Expected %v got %v", apiSslRes.CDN, clientSslRes.CDN)
	}
	if apiSslRes.Certificate.CSR != clientSslRes.Certificate.CSR {
		t.Errorf("CSR -- Expected %v got %v", apiSslRes.Certificate.CSR, clientSslRes.Certificate.CSR)
	}
	if apiSslRes.Certificate.Crt != clientSslRes.Certificate.Crt {
		t.Errorf("CRT -- Expected %v got %v", apiSslRes.Certificate.Crt, clientSslRes.Certificate.Crt)
	}
	if apiSslRes.Certificate.Key != clientSslRes.Certificate.Key {
		t.Errorf("Key -- Expected %v got %v", apiSslRes.Certificate.Key, clientSslRes.Certificate.Key)
	}
	if apiSslRes.City != clientSslRes.City {
		t.Errorf("City -- Expected %v got %v", apiSslRes.City, clientSslRes.City)
	}
	if apiSslRes.Country != clientSslRes.Country {
		t.Errorf("Country -- Expected %v got %v", apiSslRes.Country, clientSslRes.Country)
	}
	if apiSslRes.DeliveryService != clientSslRes.DeliveryService {
		t.Errorf("DeliveryService -- Expected %v got %v", apiSslRes.DeliveryService, clientSslRes.DeliveryService)
	}
	if apiSslRes.Hostname != clientSslRes.Hostname {
		t.Errorf("Hostname -- Expected %v got %v", apiSslRes.Hostname, clientSslRes.Hostname)
	}
	if apiSslRes.Organization != clientSslRes.Organization {
		t.Errorf("Organization -- Expected %v got %v", apiSslRes.Organization, clientSslRes.Organization)
	}
}
