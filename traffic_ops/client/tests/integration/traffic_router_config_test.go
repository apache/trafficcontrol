package integration

import (
	"encoding/json"
	"fmt"
	"testing"

	traffic_ops "github.com/apache/incubator-trafficcontrol/traffic_ops/client"
)

func TestTrafficRouterConfig(t *testing.T) {
	cdn, err := GetCdn()
	if err != nil {
		t.Errorf("Could not get CDN, error was: %v\n", err)
	}
	uri := fmt.Sprintf("/api/1.2/cdns/%s/configs/routing.json", cdn.Name)
	resp, err := Request(*to, "GET", uri, nil)
	if err != nil {
		t.Errorf("Could not get %s reponse was: %v\n", uri, err)
		t.FailNow()
	}

	defer resp.Body.Close()
	var apiTRConfigRes traffic_ops.TRConfigResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiTRConfigRes); err != nil {
		t.Errorf("Could not decode Traffic Router Config response.  Error is: %v\n", err)
		t.FailNow()
	}
	apiTRConfig := apiTRConfigRes.Response

	clientTRConfig, err := to.TrafficRouterConfig(cdn.Name)
	if err != nil {
		t.Errorf("Could not get Traffic Router Config from client.  Error is: %v\n", err)
		t.FailNow()
	}

	if len(apiTRConfig.CacheGroups) != len(clientTRConfig.CacheGroups) {
		t.Errorf("Length of Traffic Router config cachegroups do not match! Expected %v, got %v\n", len(apiTRConfig.CacheGroups), len(clientTRConfig.CacheGroups))
	}

	for _, apiCg := range apiTRConfig.CacheGroups {
		match := false
		for _, clientCg := range clientTRConfig.CacheGroups {
			if apiCg == clientCg {
				match = true
			}
		}
		if !match {
			t.Errorf("Did not get a cachegroup matching %+v\n", apiCg)
		}
	}

	if len(apiTRConfig.DeliveryServices) != len(clientTRConfig.DeliveryServices) {
		t.Errorf("Length of Traffic Router config deliveryserivces do not match! Expected %v, got %v\n", len(apiTRConfig.DeliveryServices), len(clientTRConfig.DeliveryServices))
	}

	for _, apiDs := range apiTRConfig.DeliveryServices {
		match := false
		for _, clientDs := range clientTRConfig.DeliveryServices {
			if apiDs.XMLID == clientDs.XMLID {
				match = true
				if apiDs.BypassDestination != clientDs.BypassDestination {
					t.Errorf("BypassDestination -- Expected %v, got %v\n", apiDs.BypassDestination, clientDs.BypassDestination)
				}
				if apiDs.CoverageZoneOnly != clientDs.CoverageZoneOnly {
					t.Errorf("CZ Only -- Expected %v, got %v\n", apiDs.CoverageZoneOnly, clientDs.CoverageZoneOnly)
				}
				if len(apiDs.Domains) != len(clientDs.Domains) {
					t.Errorf("len Domains -- Expected %v, got %v\n", len(apiDs.Domains), len(clientDs.Domains))
				}
				for _, apiDomain := range apiDs.Domains {
					domainMatch := false
					for _, clientDomain := range clientDs.Domains {
						if apiDomain == clientDomain {
							domainMatch = true
						}
					}
					if !domainMatch {
						t.Errorf("Domains -- Did not find a match for %v\n", apiDomain)
					}
				}
				if len(apiDs.MatchSets) != len(clientDs.MatchSets) {
					t.Errorf("len Matchsets -- Expected %v, got %v\n", len(apiDs.MatchSets), len(clientDs.MatchSets))
				}
				for _, apiMatch := range apiDs.MatchSets {
					foundMatch := false
					for _, clientMatch := range clientDs.MatchSets {
						if apiMatch.Protocol == clientMatch.Protocol && len(apiMatch.MatchList) == len(clientMatch.MatchList) {
							foundMatch = true
						}
					}
					if !foundMatch {
						t.Errorf("Matchsets -- Did not find a match for %+v\n", apiMatch)
					}
				}
				if apiDs.MissLocation != clientDs.MissLocation {
					t.Errorf("MissLocation -- Expected %v, got %v\n", apiDs.MissLocation, clientDs.MissLocation)
				}
				if apiDs.Soa != clientDs.Soa {
					t.Errorf("Soa-- Expected %v, got %v\n", apiDs.Soa, clientDs.Soa)
				}
				if apiDs.TTL != clientDs.TTL {
					t.Errorf("TTL -- Expected %v, got %v\n", apiDs.TTL, clientDs.TTL)
				}
				if apiDs.TTLs != clientDs.TTLs {
					t.Errorf("TTLs -- Expected %v, got %v\n", apiDs.TTLs, clientDs.TTLs)
				}
				if len(apiDs.StatcDNSEntries) != len(clientDs.StatcDNSEntries) {
					t.Errorf("len StaticDNSEntries -- Expected %v, got %v\n", len(apiDs.StatcDNSEntries), len(clientDs.StatcDNSEntries))
				}
				for _, apiEntry := range apiDs.StatcDNSEntries {
					found := false
					for _, clientEntry := range clientDs.StatcDNSEntries {
						if apiEntry == clientEntry {
							found = true
						}
					}
					if !found {
						t.Errorf("Static DNS -- Did not find a match for %+v\n", apiEntry)
					}
				}

			}
		}
		if !match {
			t.Errorf("Did not get a Deliveryservice matching %+v\n", apiDs)
		}
	}

	if len(apiTRConfig.TrafficMonitors) != len(clientTRConfig.TrafficMonitors) {
		t.Errorf("Length of Traffic Router config Traffic Routers does not match! Expected %v, got %v\n", len(apiTRConfig.TrafficMonitors), len(clientTRConfig.TrafficMonitors))
	}

	for _, apiTM := range apiTRConfig.TrafficMonitors {
		match := false
		for _, clientTM := range clientTRConfig.TrafficMonitors {
			if apiTM == clientTM {
				match = true
			}
		}
		if !match {
			t.Errorf("Did not get a Traffic Router matching %+v\n", apiTM)
		}
	}

	if len(apiTRConfig.TrafficServers) != len(clientTRConfig.TrafficServers) {
		t.Errorf("Length of Traffic Router config traffic servers does not match! Expected %v, got %v\n", len(apiTRConfig.TrafficServers), len(clientTRConfig.TrafficServers))
	}

	for _, apiTS := range apiTRConfig.TrafficServers {
		match := false
		for _, clientTS := range clientTRConfig.TrafficServers {
			if apiTS.HostName == clientTS.HostName {
				match = true
				if apiTS.CacheGroup != clientTS.CacheGroup {
					t.Errorf("Cachegroup -- Expected %v, got %v\n", apiTS.CacheGroup, clientTS.CacheGroup)
				}
				if len(apiTS.DeliveryServices) != len(clientTS.DeliveryServices) {
					t.Errorf("len DeliveryServices -- Expected %v, got %v\n", len(apiTS.DeliveryServices), len(clientTS.DeliveryServices))
				}
				for _, apiDS := range apiTS.DeliveryServices {
					dsMatch := false
					for _, clientDS := range clientTS.DeliveryServices {
						if apiDS.Xmlid == clientDS.Xmlid && len(apiDS.Remaps) == len(clientDS.Remaps) {
							dsMatch = true
						}
					}
					if !dsMatch {
						t.Errorf("Could not finding a matching DS for %v\n", apiDS.Xmlid)
					}
				}
				if apiTS.FQDN != clientTS.FQDN {
					t.Errorf("FQDN -- Expected %v, got %v\n", apiTS.FQDN, clientTS.FQDN)
				}
				if apiTS.HashID != clientTS.HashID {
					t.Errorf("HashID -- Expected %v, got %v\n", apiTS.HashID, clientTS.HashID)
				}
				if apiTS.IP != clientTS.IP {
					t.Errorf("IP -- Expected %v, got %v\n", apiTS.IP, clientTS.IP)
				}
				if apiTS.IP6 != clientTS.IP6 {
					t.Errorf("IP6 -- Expected %v, got %v\n", apiTS.IP6, clientTS.IP6)
				}
				if apiTS.InterfaceName != clientTS.InterfaceName {
					t.Errorf("Interface Name -- Expected %v, got %v\n", apiTS.InterfaceName, clientTS.InterfaceName)
				}
				if apiTS.Port != clientTS.Port {
					t.Errorf("Port -- Expected %v, got %v\n", apiTS.Port, clientTS.Port)
				}
				if apiTS.Profile != clientTS.Profile {
					t.Errorf("Profile -- Expected %v, got %v\n", apiTS.Profile, clientTS.Profile)
				}
				if apiTS.Status != clientTS.Status {
					t.Errorf("Status -- Expected %v, got %v\n", apiTS.Status, clientTS.Status)
				}
				if apiTS.Type != clientTS.Type {
					t.Errorf("Type -- Expected %v, got %v\n", apiTS.Type, clientTS.Type)
				}
			}
		}
		if !match {
			t.Errorf("Did not get a Traffic Server matching %+v\n", apiTS)
		}
	}
}
