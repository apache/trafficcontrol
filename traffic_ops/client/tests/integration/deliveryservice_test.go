package integration

import (
	"encoding/json"
	"strconv"
	"testing"
	"time"

	traffic_ops "github.com/apache/incubator-trafficcontrol/traffic_ops/client"
)

// TestDeliveryServices compares the results of the Deliveryservices api and Deliveryservices client
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

var testDsID string

func TestCreateDs(t *testing.T) {
	//create a DS and validate response
	cdn, err := GetCdn()
	if err != nil {
		t.Errorf("TestCreateDs -- Could not get CDNs from TO...%v\n", err)
	}

	profile, err := GetProfile()
	if err != nil {
		t.Errorf("TestCreateDs -- Could not get Profiles from TO...%v\n", err)
	}

	newDs := new(traffic_ops.DeliveryService)
	newDs.Active = false
	newDs.CCRDNSTTL = 30
	newDs.CDNName = cdn.Name
	newDs.CacheURL = "cacheURL"
	newDs.CheckPath = "CheckPath"
	newDs.DNSBypassCname = "DNSBypassCNAME"
	newDs.DNSBypassIP = "10.10.10.10"
	newDs.DNSBypassIP6 = "FF01:0:0:0:0:0:0:FB"
	newDs.DNSBypassTTL = 30
	newDs.DSCP = 0
	newDs.DisplayName = "DisplayName"
	newDs.EdgeHeaderRewrite = "EdgeHeaderRewrite"
	newDs.GeoLimit = 5
	newDs.GeoProvider = 1
	newDs.GlobalMaxMBPS = 15000
	newDs.GlobalMaxTPS = 15000
	newDs.HTTPBypassFQDN = "HTTPBypassFQDN"
	newDs.IPV6RoutingEnabled = true
	newDs.InfoURL = "InfoUrl"
	newDs.InitialDispersion = 5
	newDs.LongDesc = "LongDesc"
	newDs.LongDesc1 = "LongDesc1"
	newDs.LongDesc2 = "LongDesc2"
	newDs.MaxDNSAnswers = 5
	newDs.MidHeaderRewrite = "MidHeaderRewrite"
	newDs.MissLat = 5.555
	newDs.MissLong = -50.5050
	newDs.MultiSiteOrigin = true
	newDs.OrgServerFQDN = "http://OrgServerFQDN"
	newDs.ProfileDesc = profile.Description
	newDs.ProfileName = profile.Name
	newDs.Protocol = 1
	newDs.QStringIgnore = 1
	newDs.RangeRequestHandling = 0
	newDs.RegexRemap = "regexRemap"
	newDs.RemapText = "remapText"
	newDs.Signed = false
	newDs.TRResponseHeaders = "TRResponseHeaders"
	newDs.Type = "HTTP"
	newDs.XMLID = "Test-DS-" + strconv.FormatInt(time.Now().Unix(), 10)
	newDs.RegionalGeoBlocking = false
	newDs.LogsEnabled = false

	//Create currently does not write regexes...
	// newDsMatch1 := new(traffic_ops.DeliveryServiceMatch)
	// newDsMatch1.Pattern = "Pattern1"
	// newDsMatch1.SetNumber = "0"
	// newDsMatch1.Type = "HOST"

	// newDsMatch2 := new(traffic_ops.DeliveryServiceMatch)
	// newDsMatch2.Pattern = "Pattern2"
	// newDsMatch2.SetNumber = "1"
	// newDsMatch2.Type = "HOST"

	// newDs.MatchList = append(newDs.MatchList, *newDsMatch1)
	// newDs.MatchList = append(newDs.MatchList, *newDsMatch2)

	res, err := to.CreateDeliveryService(newDs)
	if err != nil {
		t.Error("Failed to create deliveryservice!  Error is: ", err)
	} else {
		compareDs(*newDs, res.Response[0], t)
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
	if ds1.LastUpdated != ds2.LastUpdated {
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
}
