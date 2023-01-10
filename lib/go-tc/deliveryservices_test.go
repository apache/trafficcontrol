package tc

/*
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func compareV31DSes(a, b DeliveryServiceNullableV30, t *testing.T) {
	if (a.Active == nil && b.Active != nil) || (a.Active != nil && b.Active == nil) {
		t.Error("Mismatched 'Active' property; one was nil but the other was not.")
	} else if a.Active != nil && *a.Active != *b.Active {
		t.Errorf("Mismatched 'Active' property; one was '%v', the other was '%v'", *a.Active, *b.Active)
	}
	if (a.AnonymousBlockingEnabled == nil && b.AnonymousBlockingEnabled != nil) || (a.AnonymousBlockingEnabled != nil && b.AnonymousBlockingEnabled == nil) {
		t.Error("Mismatched 'AnonymousBlockingEnabled' property; one was nil but the other was not.")
	} else if a.AnonymousBlockingEnabled != nil && *a.AnonymousBlockingEnabled != *b.AnonymousBlockingEnabled {
		t.Errorf("Mismatched 'AnonymousBlockingEnabled' property; one was '%v', the other was '%v'", *a.AnonymousBlockingEnabled, *b.AnonymousBlockingEnabled)
	}
	if (a.CCRDNSTTL == nil && b.CCRDNSTTL != nil) || (a.CCRDNSTTL != nil && b.CCRDNSTTL == nil) {
		t.Error("Mismatched 'CCRDNSTTL' property; one was nil but the other was not.")
	} else if a.CCRDNSTTL != nil && *a.CCRDNSTTL != *b.CCRDNSTTL {
		t.Errorf("Mismatched 'CCRDNSTTL' property; one was '%v', the other was '%v'", *a.CCRDNSTTL, *b.CCRDNSTTL)
	}
	if (a.CDNID == nil && b.CDNID != nil) || (a.CDNID != nil && b.CDNID == nil) {
		t.Error("Mismatched 'CDNID' property; one was nil but the other was not.")
	} else if a.CDNID != nil && *a.CDNID != *b.CDNID {
		t.Errorf("Mismatched 'CDNID' property; one was '%v', the other was '%v'", *a.CDNID, *b.CDNID)
	}
	if (a.CDNName == nil && b.CDNName != nil) || (a.CDNName != nil && b.CDNName == nil) {
		t.Error("Mismatched 'CDNName' property; one was nil but the other was not.")
	} else if a.CDNName != nil && *a.CDNName != *b.CDNName {
		t.Errorf("Mismatched 'CDNName' property; one was '%v', the other was '%v'", *a.CDNName, *b.CDNName)
	}
	if (a.CheckPath == nil && b.CheckPath != nil) || (a.CheckPath != nil && b.CheckPath == nil) {
		t.Error("Mismatched 'CheckPath' property; one was nil but the other was not.")
	} else if a.CheckPath != nil && *a.CheckPath != *b.CheckPath {
		t.Errorf("Mismatched 'CheckPath' property; one was '%v', the other was '%v'", *a.CheckPath, *b.CheckPath)
	}
	if len(a.ConsistentHashQueryParams) != len(b.ConsistentHashQueryParams) {
		t.Errorf("Mismatched 'ConsistentHashQueryParams' property; one contained %d but the other contained %d.", len(a.ConsistentHashQueryParams), len(b.ConsistentHashQueryParams))
	} else {
		for i, qp := range a.ConsistentHashQueryParams {
			if qp != b.ConsistentHashQueryParams[i] {
				t.Errorf("Mismatched 'ConsistentHashQueryParams[%d]'; one was %s, but the other was %s", i, qp, b.ConsistentHashQueryParams[i])
			}
		}
	}
	if (a.ConsistentHashRegex == nil && b.ConsistentHashRegex != nil) || (a.ConsistentHashRegex != nil && b.ConsistentHashRegex == nil) {
		t.Error("Mismatched 'ConsistentHashRegex' property; one was nil but the other was not.")
	} else if a.ConsistentHashRegex != nil && *a.ConsistentHashRegex != *b.ConsistentHashRegex {
		t.Errorf("Mismatched 'ConsistentHashRegex' property; one was '%v', the other was '%v'", *a.ConsistentHashRegex, *b.ConsistentHashRegex)
	}
	if (a.DeepCachingType == nil && b.DeepCachingType != nil) || (a.DeepCachingType != nil && b.DeepCachingType == nil) {
		t.Error("Mismatched 'DeepCachingType' property; one was nil but the other was not.")
	} else if a.DeepCachingType != nil && *a.DeepCachingType != *b.DeepCachingType {
		t.Errorf("Mismatched 'DeepCachingType' property; one was '%v', the other was '%v'", *a.DeepCachingType, *b.DeepCachingType)
	}
	if (a.DisplayName == nil && b.DisplayName != nil) || (a.DisplayName != nil && b.DisplayName == nil) {
		t.Error("Mismatched 'DisplayName' property; one was nil but the other was not.")
	} else if a.DisplayName != nil && *a.DisplayName != *b.DisplayName {
		t.Errorf("Mismatched 'DisplayName' property; one was '%v', the other was '%v'", *a.DisplayName, *b.DisplayName)
	}
	if (a.DNSBypassCNAME == nil && b.DNSBypassCNAME != nil) || (a.DNSBypassCNAME != nil && b.DNSBypassCNAME == nil) {
		t.Error("Mismatched 'DNSBypassCNAME' property; one was nil but the other was not.")
	} else if a.DNSBypassCNAME != nil && *a.DNSBypassCNAME != *b.DNSBypassCNAME {
		t.Errorf("Mismatched 'DNSBypassCNAME' property; one was '%v', the other was '%v'", *a.DNSBypassCNAME, *b.DNSBypassCNAME)
	}
	if (a.DNSBypassIP6 == nil && b.DNSBypassIP6 != nil) || (a.DNSBypassIP6 != nil && b.DNSBypassIP6 == nil) {
		t.Error("Mismatched 'DNSBypassIP6' property; one was nil but the other was not.")
	} else if a.DNSBypassIP6 != nil && *a.DNSBypassIP6 != *b.DNSBypassIP6 {
		t.Errorf("Mismatched 'DNSBypassIP6' property; one was '%v', the other was '%v'", *a.DNSBypassIP6, *b.DNSBypassIP6)
	}
	if (a.DNSBypassIP == nil && b.DNSBypassIP != nil) || (a.DNSBypassIP != nil && b.DNSBypassIP == nil) {
		t.Error("Mismatched 'DNSBypassIP' property; one was nil but the other was not.")
	} else if a.DNSBypassIP != nil && *a.DNSBypassIP != *b.DNSBypassIP {
		t.Errorf("Mismatched 'DNSBypassIP' property; one was '%v', the other was '%v'", *a.DNSBypassIP, *b.DNSBypassIP)
	}
	if (a.DNSBypassTTL == nil && b.DNSBypassTTL != nil) || (a.DNSBypassTTL != nil && b.DNSBypassTTL == nil) {
		t.Error("Mismatched 'DNSBypassTTL' property; one was nil but the other was not.")
	} else if a.DNSBypassTTL != nil && *a.DNSBypassTTL != *b.DNSBypassTTL {
		t.Errorf("Mismatched 'DNSBypassTTL' property; one was '%v', the other was '%v'", *a.DNSBypassTTL, *b.DNSBypassTTL)
	}
	if (a.DSCP == nil && b.DSCP != nil) || (a.DSCP != nil && b.DSCP == nil) {
		t.Error("Mismatched 'DSCP' property; one was nil but the other was not.")
	} else if a.DSCP != nil && *a.DSCP != *b.DSCP {
		t.Errorf("Mismatched 'DSCP' property; one was '%v', the other was '%v'", *a.DSCP, *b.DSCP)
	}
	if a.EcsEnabled != b.EcsEnabled {
		t.Errorf("Mismatched 'EcsEnabled' property; one was '%v', the other was '%v'", a.EcsEnabled, b.EcsEnabled)
	}
	if (a.EdgeHeaderRewrite == nil && b.EdgeHeaderRewrite != nil) || (a.EdgeHeaderRewrite != nil && b.EdgeHeaderRewrite == nil) {
		t.Error("Mismatched 'EdgeHeaderRewrite' property; one was nil but the other was not.")
	} else if a.EdgeHeaderRewrite != nil && *a.EdgeHeaderRewrite != *b.EdgeHeaderRewrite {
		t.Errorf("Mismatched 'EdgeHeaderRewrite' property; one was '%v', the other was '%v'", *a.EdgeHeaderRewrite, *b.EdgeHeaderRewrite)
	}
	if len(a.ExampleURLs) != len(b.ExampleURLs) {
		t.Errorf("Mismatched 'ExampleURLs' property; one contained %d but the other contained %d", len(a.ExampleURLs), len(b.ExampleURLs))
	} else {
		for i, eu := range a.ExampleURLs {
			if eu != b.ExampleURLs[i] {
				t.Errorf("Mismatched 'ExampleURLs[%d]' property; one was '%v', the other was '%v'", i, eu, b.ExampleURLs[i])
			}
		}
	}
	if (a.FirstHeaderRewrite == nil && b.FirstHeaderRewrite != nil) || (a.FirstHeaderRewrite != nil && b.FirstHeaderRewrite == nil) {
		t.Error("Mismatched 'FirstHeaderRewrite' property; one was nil but the other was not.")
	} else if a.FirstHeaderRewrite != nil && *a.FirstHeaderRewrite != *b.FirstHeaderRewrite {
		t.Errorf("Mismatched 'FirstHeaderRewrite' property; one was '%v', the other was '%v'", *a.FirstHeaderRewrite, *b.FirstHeaderRewrite)
	}
	if (a.FQPacingRate == nil && b.FQPacingRate != nil) || (a.FQPacingRate != nil && b.FQPacingRate == nil) {
		t.Error("Mismatched 'FQPacingRate' property; one was nil but the other was not.")
	} else if a.FQPacingRate != nil && *a.FQPacingRate != *b.FQPacingRate {
		t.Errorf("Mismatched 'FQPacingRate' property; one was '%v', the other was '%v'", *a.FQPacingRate, *b.FQPacingRate)
	}
	if (a.GeoLimit == nil && b.GeoLimit != nil) || (a.GeoLimit != nil && b.GeoLimit == nil) {
		t.Error("Mismatched 'GeoLimit' property; one was nil but the other was not.")
	} else if a.GeoLimit != nil && *a.GeoLimit != *b.GeoLimit {
		t.Errorf("Mismatched 'GeoLimit' property; one was '%v', the other was '%v'", *a.GeoLimit, *b.GeoLimit)
	}
	if (a.GeoLimitCountries == nil && b.GeoLimitCountries != nil) || (a.GeoLimitCountries != nil && b.GeoLimitCountries == nil) {
		t.Error("Mismatched 'GeoLimitCountries' property; one was nil but the other was not.")
	} else if a.GeoLimitCountries != nil && *a.GeoLimitCountries != *b.GeoLimitCountries {
		t.Errorf("Mismatched 'GeoLimitCountries' property; one was '%v', the other was '%v'", *a.GeoLimitCountries, *b.GeoLimitCountries)
	}
	if (a.GeoLimitRedirectURL == nil && b.GeoLimitRedirectURL != nil) || (a.GeoLimitRedirectURL != nil && b.GeoLimitRedirectURL == nil) {
		t.Error("Mismatched 'GeoLimitRedirectURL' property; one was nil but the other was not.")
	} else if a.GeoLimitRedirectURL != nil && *a.GeoLimitRedirectURL != *b.GeoLimitRedirectURL {
		t.Errorf("Mismatched 'GeoLimitRedirectURL' property; one was '%v', the other was '%v'", *a.GeoLimitRedirectURL, *b.GeoLimitRedirectURL)
	}
	if (a.GeoProvider == nil && b.GeoProvider != nil) || (a.GeoProvider != nil && b.GeoProvider == nil) {
		t.Error("Mismatched 'GeoProvider' property; one was nil but the other was not.")
	} else if a.GeoProvider != nil && *a.GeoProvider != *b.GeoProvider {
		t.Errorf("Mismatched 'GeoProvider' property; one was '%v', the other was '%v'", *a.GeoProvider, *b.GeoProvider)
	}
	if (a.GlobalMaxMBPS == nil && b.GlobalMaxMBPS != nil) || (a.GlobalMaxMBPS != nil && b.GlobalMaxMBPS == nil) {
		t.Error("Mismatched 'GlobalMaxMBPS' property; one was nil but the other was not.")
	} else if a.GlobalMaxMBPS != nil && *a.GlobalMaxMBPS != *b.GlobalMaxMBPS {
		t.Errorf("Mismatched 'GlobalMaxMBPS' property; one was '%v', the other was '%v'", *a.GlobalMaxMBPS, *b.GlobalMaxMBPS)
	}
	if (a.GlobalMaxTPS == nil && b.GlobalMaxTPS != nil) || (a.GlobalMaxTPS != nil && b.GlobalMaxTPS == nil) {
		t.Error("Mismatched 'GlobalMaxTPS' property; one was nil but the other was not.")
	} else if a.GlobalMaxTPS != nil && *a.GlobalMaxTPS != *b.GlobalMaxTPS {
		t.Errorf("Mismatched 'GlobalMaxTPS' property; one was '%v', the other was '%v'", *a.GlobalMaxTPS, *b.GlobalMaxTPS)
	}
	if (a.HTTPBypassFQDN == nil && b.HTTPBypassFQDN != nil) || (a.HTTPBypassFQDN != nil && b.HTTPBypassFQDN == nil) {
		t.Error("Mismatched 'HTTPBypassFQDN' property; one was nil but the other was not.")
	} else if a.HTTPBypassFQDN != nil && *a.HTTPBypassFQDN != *b.HTTPBypassFQDN {
		t.Errorf("Mismatched 'HTTPBypassFQDN' property; one was '%v', the other was '%v'", *a.HTTPBypassFQDN, *b.HTTPBypassFQDN)
	}
	if (a.ID == nil && b.ID != nil) || (a.ID != nil && b.ID == nil) {
		t.Error("Mismatched 'ID' property; one was nil but the other was not.")
	} else if a.ID != nil && *a.ID != *b.ID {
		t.Errorf("Mismatched 'ID' property; one was '%v', the other was '%v'", *a.ID, *b.ID)
	}
	if (a.InfoURL == nil && b.InfoURL != nil) || (a.InfoURL != nil && b.InfoURL == nil) {
		t.Error("Mismatched 'InfoURL' property; one was nil but the other was not.")
	} else if a.InfoURL != nil && *a.InfoURL != *b.InfoURL {
		t.Errorf("Mismatched 'InfoURL' property; one was '%v', the other was '%v'", *a.InfoURL, *b.InfoURL)
	}
	if (a.InitialDispersion == nil && b.InitialDispersion != nil) || (a.InitialDispersion != nil && b.InitialDispersion == nil) {
		t.Error("Mismatched 'InitialDispersion' property; one was nil but the other was not.")
	} else if a.InitialDispersion != nil && *a.InitialDispersion != *b.InitialDispersion {
		t.Errorf("Mismatched 'InitialDispersion' property; one was '%v', the other was '%v'", *a.InitialDispersion, *b.InitialDispersion)
	}
	if (a.InnerHeaderRewrite == nil && b.InnerHeaderRewrite != nil) || (a.InnerHeaderRewrite != nil && b.InnerHeaderRewrite == nil) {
		t.Error("Mismatched 'InnerHeaderRewrite' property; one was nil but the other was not.")
	} else if a.InnerHeaderRewrite != nil && *a.InnerHeaderRewrite != *b.InnerHeaderRewrite {
		t.Errorf("Mismatched 'InnerHeaderRewrite' property; one was '%v', the other was '%v'", *a.InnerHeaderRewrite, *b.InnerHeaderRewrite)
	}
	if (a.IPV6RoutingEnabled == nil && b.IPV6RoutingEnabled != nil) || (a.IPV6RoutingEnabled != nil && b.IPV6RoutingEnabled == nil) {
		t.Error("Mismatched 'IPV6RoutingEnabled' property; one was nil but the other was not.")
	} else if a.IPV6RoutingEnabled != nil && *a.IPV6RoutingEnabled != *b.IPV6RoutingEnabled {
		t.Errorf("Mismatched 'IPV6RoutingEnabled' property; one was '%v', the other was '%v'", *a.IPV6RoutingEnabled, *b.IPV6RoutingEnabled)
	}
	if (a.LastHeaderRewrite == nil && b.LastHeaderRewrite != nil) || (a.LastHeaderRewrite != nil && b.LastHeaderRewrite == nil) {
		t.Error("Mismatched 'LastHeaderRewrite' property; one was nil but the other was not.")
	} else if a.LastHeaderRewrite != nil && *a.LastHeaderRewrite != *b.LastHeaderRewrite {
		t.Errorf("Mismatched 'LastHeaderRewrite' property; one was '%v', the other was '%v'", *a.LastHeaderRewrite, *b.LastHeaderRewrite)
	}
	if (a.LastUpdated == nil && b.LastUpdated != nil) || (a.LastUpdated != nil && b.LastUpdated == nil) {
		t.Error("Mismatched 'LastUpdated' property; one was nil but the other was not.")
	} else if a.LastUpdated != nil && *a.LastUpdated != *b.LastUpdated {
		t.Errorf("Mismatched 'LastUpdated' property; one was '%v', the other was '%v'", *a.LastUpdated, *b.LastUpdated)
	}
	if (a.LogsEnabled == nil && b.LogsEnabled != nil) || (a.LogsEnabled != nil && b.LogsEnabled == nil) {
		t.Error("Mismatched 'LogsEnabled' property; one was nil but the other was not.")
	} else if a.LogsEnabled != nil && *a.LogsEnabled != *b.LogsEnabled {
		t.Errorf("Mismatched 'LogsEnabled' property; one was '%v', the other was '%v'", *a.LogsEnabled, *b.LogsEnabled)
	}
	if (a.LongDesc1 == nil && b.LongDesc1 != nil) || (a.LongDesc1 != nil && b.LongDesc1 == nil) {
		t.Error("Mismatched 'LongDesc1' property; one was nil but the other was not.")
	} else if a.LongDesc1 != nil && *a.LongDesc1 != *b.LongDesc1 {
		t.Errorf("Mismatched 'LongDesc1' property; one was '%v', the other was '%v'", *a.LongDesc1, *b.LongDesc1)
	}
	if (a.LongDesc2 == nil && b.LongDesc2 != nil) || (a.LongDesc2 != nil && b.LongDesc2 == nil) {
		t.Error("Mismatched 'LongDesc2' property; one was nil but the other was not.")
	} else if a.LongDesc2 != nil && *a.LongDesc2 != *b.LongDesc2 {
		t.Errorf("Mismatched 'LongDesc2' property; one was '%v', the other was '%v'", *a.LongDesc2, *b.LongDesc2)
	}
	if (a.LongDesc == nil && b.LongDesc != nil) || (a.LongDesc != nil && b.LongDesc == nil) {
		t.Error("Mismatched 'LongDesc' property; one was nil but the other was not.")
	} else if a.LongDesc != nil && *a.LongDesc != *b.LongDesc {
		t.Errorf("Mismatched 'LongDesc' property; one was '%v', the other was '%v'", *a.LongDesc, *b.LongDesc)
	}
	if (a.MatchList != nil && b.MatchList == nil) || (a.MatchList == nil && b.MatchList != nil) {
		t.Error("Mismatched 'MatchList' property; one was nil but the other was not")
	} else if a.MatchList != nil {
		if len(*a.MatchList) != len(*b.MatchList) {
			t.Errorf("Mismatched 'MatchList' property; one contained %d but the other contained %d", len(*a.MatchList), len(*b.MatchList))
		} else {
			for i, m := range *a.MatchList {
				if m != (*b.MatchList)[i] {
					t.Errorf("Mismatched 'MatchList[%d]' property; one was '%v', the other was '%v'", i, m, (*b.MatchList)[i])
				}
			}
		}
	}
	if (a.MaxDNSAnswers == nil && b.MaxDNSAnswers != nil) || (a.MaxDNSAnswers != nil && b.MaxDNSAnswers == nil) {
		t.Error("Mismatched 'MaxDNSAnswers' property; one was nil but the other was not.")
	} else if a.MaxDNSAnswers != nil && *a.MaxDNSAnswers != *b.MaxDNSAnswers {
		t.Errorf("Mismatched 'MaxDNSAnswers' property; one was '%v', the other was '%v'", *a.MaxDNSAnswers, *b.MaxDNSAnswers)
	}
	if (a.MaxOriginConnections == nil && b.MaxOriginConnections != nil) || (a.MaxOriginConnections != nil && b.MaxOriginConnections == nil) {
		t.Error("Mismatched 'MaxOriginConnections' property; one was nil but the other was not.")
	} else if a.MaxOriginConnections != nil && *a.MaxOriginConnections != *b.MaxOriginConnections {
		t.Errorf("Mismatched 'MaxOriginConnections' property; one was '%v', the other was '%v'", *a.MaxOriginConnections, *b.MaxOriginConnections)
	}
	if (a.MaxRequestHeaderBytes == nil && b.MaxRequestHeaderBytes != nil) || (a.MaxRequestHeaderBytes != nil && b.MaxRequestHeaderBytes == nil) {
		t.Error("Mismatched 'MaxRequestHeaderBytes' property; one was nil but the other was not.")
	} else if a.MaxRequestHeaderBytes != nil && *a.MaxRequestHeaderBytes != *b.MaxRequestHeaderBytes {
		t.Errorf("Mismatched 'MaxRequestHeaderBytes' property; one was '%v', the other was '%v'", *a.MaxRequestHeaderBytes, *b.MaxRequestHeaderBytes)
	}
	if (a.MidHeaderRewrite == nil && b.MidHeaderRewrite != nil) || (a.MidHeaderRewrite != nil && b.MidHeaderRewrite == nil) {
		t.Error("Mismatched 'MidHeaderRewrite' property; one was nil but the other was not.")
	} else if a.MidHeaderRewrite != nil && *a.MidHeaderRewrite != *b.MidHeaderRewrite {
		t.Errorf("Mismatched 'MidHeaderRewrite' property; one was '%v', the other was '%v'", *a.MidHeaderRewrite, *b.MidHeaderRewrite)
	}
	if (a.MissLat == nil && b.MissLat != nil) || (a.MissLat != nil && b.MissLat == nil) {
		t.Error("Mismatched 'MissLat' property; one was nil but the other was not.")
	} else if a.MissLat != nil && *a.MissLat != *b.MissLat {
		t.Errorf("Mismatched 'MissLat' property; one was '%v', the other was '%v'", *a.MissLat, *b.MissLat)
	}
	if (a.MissLong == nil && b.MissLong != nil) || (a.MissLong != nil && b.MissLong == nil) {
		t.Error("Mismatched 'MissLong' property; one was nil but the other was not.")
	} else if a.MissLong != nil && *a.MissLong != *b.MissLong {
		t.Errorf("Mismatched 'MissLong' property; one was '%v', the other was '%v'", *a.MissLong, *b.MissLong)
	}
	if (a.MultiSiteOrigin == nil && b.MultiSiteOrigin != nil) || (a.MultiSiteOrigin != nil && b.MultiSiteOrigin == nil) {
		t.Error("Mismatched 'MultiSiteOrigin' property; one was nil but the other was not.")
	} else if a.MultiSiteOrigin != nil && *a.MultiSiteOrigin != *b.MultiSiteOrigin {
		t.Errorf("Mismatched 'MultiSiteOrigin' property; one was '%v', the other was '%v'", *a.MultiSiteOrigin, *b.MultiSiteOrigin)
	}
	if (a.OrgServerFQDN == nil && b.OrgServerFQDN != nil) || (a.OrgServerFQDN != nil && b.OrgServerFQDN == nil) {
		t.Error("Mismatched 'OrgServerFQDN' property; one was nil but the other was not.")
	} else if a.OrgServerFQDN != nil && *a.OrgServerFQDN != *b.OrgServerFQDN {
		t.Errorf("Mismatched 'OrgServerFQDN' property; one was '%v', the other was '%v'", *a.OrgServerFQDN, *b.OrgServerFQDN)
	}
	if (a.OriginShield == nil && b.OriginShield != nil) || (a.OriginShield != nil && b.OriginShield == nil) {
		t.Error("Mismatched 'OriginShield' property; one was nil but the other was not.")
	} else if a.OriginShield != nil && *a.OriginShield != *b.OriginShield {
		t.Errorf("Mismatched 'OriginShield' property; one was '%v', the other was '%v'", *a.OriginShield, *b.OriginShield)
	}
	if (a.ProfileDesc == nil && b.ProfileDesc != nil) || (a.ProfileDesc != nil && b.ProfileDesc == nil) {
		t.Error("Mismatched 'ProfileDesc' property; one was nil but the other was not.")
	} else if a.ProfileDesc != nil && *a.ProfileDesc != *b.ProfileDesc {
		t.Errorf("Mismatched 'ProfileDesc' property; one was '%v', the other was '%v'", *a.ProfileDesc, *b.ProfileDesc)
	}
	if (a.ProfileID == nil && b.ProfileID != nil) || (a.ProfileID != nil && b.ProfileID == nil) {
		t.Error("Mismatched 'ProfileID' property; one was nil but the other was not.")
	} else if a.ProfileID != nil && *a.ProfileID != *b.ProfileID {
		t.Errorf("Mismatched 'ProfileID' property; one was '%v', the other was '%v'", *a.ProfileID, *b.ProfileID)
	}
	if (a.ProfileName == nil && b.ProfileName != nil) || (a.ProfileName != nil && b.ProfileName == nil) {
		t.Error("Mismatched 'ProfileName' property; one was nil but the other was not.")
	} else if a.ProfileName != nil && *a.ProfileName != *b.ProfileName {
		t.Errorf("Mismatched 'ProfileName' property; one was '%v', the other was '%v'", *a.ProfileName, *b.ProfileName)
	}
	if (a.Protocol == nil && b.Protocol != nil) || (a.Protocol != nil && b.Protocol == nil) {
		t.Error("Mismatched 'Protocol' property; one was nil but the other was not.")
	} else if a.Protocol != nil && *a.Protocol != *b.Protocol {
		t.Errorf("Mismatched 'Protocol' property; one was '%v', the other was '%v'", *a.Protocol, *b.Protocol)
	}
	if (a.QStringIgnore == nil && b.QStringIgnore != nil) || (a.QStringIgnore != nil && b.QStringIgnore == nil) {
		t.Error("Mismatched 'QStringIgnore' property; one was nil but the other was not.")
	} else if a.QStringIgnore != nil && *a.QStringIgnore != *b.QStringIgnore {
		t.Errorf("Mismatched 'QStringIgnore' property; one was '%v', the other was '%v'", *a.QStringIgnore, *b.QStringIgnore)
	}
	if (a.RangeRequestHandling == nil && b.RangeRequestHandling != nil) || (a.RangeRequestHandling != nil && b.RangeRequestHandling == nil) {
		t.Error("Mismatched 'RangeRequestHandling' property; one was nil but the other was not.")
	} else if a.RangeRequestHandling != nil && *a.RangeRequestHandling != *b.RangeRequestHandling {
		t.Errorf("Mismatched 'RangeRequestHandling' property; one was '%v', the other was '%v'", *a.RangeRequestHandling, *b.RangeRequestHandling)
	}
	if (a.RangeSliceBlockSize == nil && b.RangeSliceBlockSize != nil) || (a.RangeSliceBlockSize != nil && b.RangeSliceBlockSize == nil) {
		t.Error("Mismatched 'RangeSliceBlockSize' property; one was nil but the other was not.")
	} else if a.RangeSliceBlockSize != nil && *a.RangeSliceBlockSize != *b.RangeSliceBlockSize {
		t.Errorf("Mismatched 'RangeSliceBlockSize' property; one was '%v', the other was '%v'", *a.RangeSliceBlockSize, *b.RangeSliceBlockSize)
	}
	if (a.RegexRemap == nil && b.RegexRemap != nil) || (a.RegexRemap != nil && b.RegexRemap == nil) {
		t.Error("Mismatched 'RegexRemap' property; one was nil but the other was not.")
	} else if a.RegexRemap != nil && *a.RegexRemap != *b.RegexRemap {
		t.Errorf("Mismatched 'RegexRemap' property; one was '%v', the other was '%v'", *a.RegexRemap, *b.RegexRemap)
	}
	if (a.RegionalGeoBlocking == nil && b.RegionalGeoBlocking != nil) || (a.RegionalGeoBlocking != nil && b.RegionalGeoBlocking == nil) {
		t.Error("Mismatched 'RegionalGeoBlocking' property; one was nil but the other was not.")
	} else if a.RegionalGeoBlocking != nil && *a.RegionalGeoBlocking != *b.RegionalGeoBlocking {
		t.Errorf("Mismatched 'RegionalGeoBlocking' property; one was '%v', the other was '%v'", *a.RegionalGeoBlocking, *b.RegionalGeoBlocking)
	}
	if (a.RemapText == nil && b.RemapText != nil) || (a.RemapText != nil && b.RemapText == nil) {
		t.Error("Mismatched 'RemapText' property; one was nil but the other was not.")
	} else if a.RemapText != nil && *a.RemapText != *b.RemapText {
		t.Errorf("Mismatched 'RemapText' property; one was '%v', the other was '%v'", *a.RemapText, *b.RemapText)
	}
	if (a.RoutingName == nil && b.RoutingName != nil) || (a.RoutingName != nil && b.RoutingName == nil) {
		t.Error("Mismatched 'RoutingName' property; one was nil but the other was not.")
	} else if a.RoutingName != nil && *a.RoutingName != *b.RoutingName {
		t.Errorf("Mismatched 'RoutingName' property; one was '%v', the other was '%v'", *a.RoutingName, *b.RoutingName)
	}
	if (a.ServiceCategory == nil && b.ServiceCategory != nil) || (a.ServiceCategory != nil && b.ServiceCategory == nil) {
		t.Error("Mismatched 'ServiceCategory' property; one was nil but the other was not.")
	} else if a.ServiceCategory != nil && *a.ServiceCategory != *b.ServiceCategory {
		t.Errorf("Mismatched 'ServiceCategory' property; one was '%v', the other was '%v'", *a.ServiceCategory, *b.ServiceCategory)
	}
	if a.Signed != b.Signed {
		t.Errorf("Mismatched 'Signed' property; one was '%v', the other was '%v'", a.Signed, b.Signed)
	}
	if (a.SigningAlgorithm == nil && b.SigningAlgorithm != nil) || (a.SigningAlgorithm != nil && b.SigningAlgorithm == nil) {
		t.Error("Mismatched 'SigningAlgorithm' property; one was nil but the other was not.")
	} else if a.SigningAlgorithm != nil && *a.SigningAlgorithm != *b.SigningAlgorithm {
		t.Errorf("Mismatched 'SigningAlgorithm' property; one was '%v', the other was '%v'", *a.SigningAlgorithm, *b.SigningAlgorithm)
	}
	if (a.SSLKeyVersion == nil && b.SSLKeyVersion != nil) || (a.SSLKeyVersion != nil && b.SSLKeyVersion == nil) {
		t.Error("Mismatched 'SSLKeyVersion' property; one was nil but the other was not.")
	} else if a.SSLKeyVersion != nil && *a.SSLKeyVersion != *b.SSLKeyVersion {
		t.Errorf("Mismatched 'SSLKeyVersion' property; one was '%v', the other was '%v'", *a.SSLKeyVersion, *b.SSLKeyVersion)
	}
	if (a.Tenant == nil && b.Tenant != nil) || (a.Tenant != nil && b.Tenant == nil) {
		t.Error("Mismatched 'Tenant' property; one was nil but the other was not.")
	} else if a.Tenant != nil && *a.Tenant != *b.Tenant {
		t.Errorf("Mismatched 'Tenant' property; one was '%v', the other was '%v'", *a.Tenant, *b.Tenant)
	}
	if (a.TenantID == nil && b.TenantID != nil) || (a.TenantID != nil && b.TenantID == nil) {
		t.Error("Mismatched 'TenantID' property; one was nil but the other was not.")
	} else if a.TenantID != nil && *a.TenantID != *b.TenantID {
		t.Errorf("Mismatched 'TenantID' property; one was '%v', the other was '%v'", *a.TenantID, *b.TenantID)
	}
	if (a.Topology == nil && b.Topology != nil) || (a.Topology != nil && b.Topology == nil) {
		t.Error("Mismatched 'Topology' property; one was nil but the other was not.")
	} else if a.Topology != nil && *a.Topology != *b.Topology {
		t.Errorf("Mismatched 'Topology' property; one was '%v', the other was '%v'", *a.Topology, *b.Topology)
	}
	if (a.TRRequestHeaders == nil && b.TRRequestHeaders != nil) || (a.TRRequestHeaders != nil && b.TRRequestHeaders == nil) {
		t.Error("Mismatched 'TRRequestHeaders' property; one was nil but the other was not.")
	} else if a.TRRequestHeaders != nil && *a.TRRequestHeaders != *b.TRRequestHeaders {
		t.Errorf("Mismatched 'TRRequestHeaders' property; one was '%v', the other was '%v'", *a.TRRequestHeaders, *b.TRRequestHeaders)
	}
	if (a.TRResponseHeaders == nil && b.TRResponseHeaders != nil) || (a.TRResponseHeaders != nil && b.TRResponseHeaders == nil) {
		t.Error("Mismatched 'TRResponseHeaders' property; one was nil but the other was not.")
	} else if a.TRResponseHeaders != nil && *a.TRResponseHeaders != *b.TRResponseHeaders {
		t.Errorf("Mismatched 'TRResponseHeaders' property; one was '%v', the other was '%v'", *a.TRResponseHeaders, *b.TRResponseHeaders)
	}
	if (a.Type == nil && b.Type != nil) || (a.Type != nil && b.Type == nil) {
		t.Error("Mismatched 'Type' property; one was nil but the other was not.")
	} else if a.Type != nil && *a.Type != *b.Type {
		t.Errorf("Mismatched 'Type' property; one was '%v', the other was '%v'", *a.Type, *b.Type)
	}
	if (a.TypeID == nil && b.TypeID != nil) || (a.TypeID != nil && b.TypeID == nil) {
		t.Error("Mismatched 'TypeID' property; one was nil but the other was not.")
	} else if a.TypeID != nil && *a.TypeID != *b.TypeID {
		t.Errorf("Mismatched 'TypeID' property; one was '%v', the other was '%v'", *a.TypeID, *b.TypeID)
	}
	if (a.XMLID == nil && b.XMLID != nil) || (a.XMLID != nil && b.XMLID == nil) {
		t.Error("Mismatched 'XMLID' property; one was nil but the other was not.")
	} else if a.XMLID != nil && *a.XMLID != *b.XMLID {
		t.Errorf("Mismatched 'XMLID' property; one was '%v', the other was '%v'", *a.XMLID, *b.XMLID)
	}
}

// This gets equivalent legacy and new Delivery Services, for testing comparisons.
func dsUpgradeAndDowngradeTestingPair() (DeliveryServiceNullableV30, DeliveryServiceV4) {
	anonymousBlockingEnabled := false
	cacheURL := "testquest"
	cCRDNSTTL := 42
	cdnID := -12
	cdnName := "cdnName"
	checkPath := "checkPath"
	consistentHashQueryParams := []string{"consistent", "hash", "query", "params"}
	consistentHashRegex := "consistentHashRegex"
	deepCachingType := DeepCachingTypeNever
	displayName := "displayName"
	dnsBypassCNAME := "dnsBypassCNAME"
	dnsBypassIP := "dnsBypassIP"
	dnsBypassIP6 := "dnsBypassIP6"
	dnsBypassTTL := 100
	dscp := -69
	ecsEnabled := true
	edgeHeaderRewrite := "edgeHeaderRewrite"
	exampleURLs := []string{"http://example", "https://URLs"}
	firstHeaderRewrite := "firstHeaderRewrite"
	fqPacingRate := 1337
	geoLimit := 2
	geoLimitCountries := []string{"geo", "Limit", "Countries"}
	geo := (GeoLimitCountriesType)(geoLimitCountries)
	geoStr := strings.Join(geoLimitCountries, ",")
	geoLimitRedirectURL := "wss://geoLimitRedirectURL"
	geoProvider := 1
	globalMaxMBPS := -72485
	globalMaxTPS := 867
	hTTPBypassFQDN := "hTTPBypassFQDN"
	id := -1551
	infoURL := "infoURL"
	initialDispersion := 65154
	innerHeaderRewrite := "innerHeaderRewrite"
	ipv6RoutingEnabled := false
	lastHeaderRewrite := "lastHeaderRewrite"
	lastUpdated := NewTimeNoMod()
	logsEnabled := true
	longDesc := "longDesc"
	longDesc1 := "longDesc1"
	longDesc2 := "longDesc2"
	matchList := []DeliveryServiceMatch{
		{
			SetNumber: 1,
			Type:      DSMatchTypeHostRegex,
			Pattern:   "this is a testing pattern",
		},
	}
	maxDNSAnswers := -76675
	maxOriginConnections := 6514684
	maxRequestHeaderBytes := 555
	midHeaderRewrite := "midHeaderRewrite"
	missLat := -98.171455
	missLong := 42.122167
	multiSiteOrigin := false
	originShield := "originShield"
	orgServerFQDN := "orgServerFQDN"
	profileDesc := "profileDesc"
	profileID := -4657
	profileName := "profileName"
	protocol := 87487
	qstringIgnore := -474
	rangeRequestHandling := 16716
	rangeSliceBlockSize := -92559
	regexRemap := "regexRemap"
	regionalGeoBlocking := true
	remapText := "remapText"
	routingName := "routingName"
	serviceCategory := "serviceCategory"
	signed := false
	signingAlgorithm := "signingAlgorithm"
	sSLKeyVersion := 721574
	tenant := "tenant"
	tenantID := -6551
	topology := "topology"
	trResponseHeaders := "trResponseHeaders"
	trRequestHeaders := "trRequestHeaders"
	typ := DSTypeDNS
	typeID := 22
	xmlid := "xmlid"
	requiredCapabilities := []string{"foo"}

	newDS := DeliveryServiceV4{}
	newDS.Active = new(bool)
	newDS.AnonymousBlockingEnabled = &anonymousBlockingEnabled
	newDS.CCRDNSTTL = &cCRDNSTTL
	newDS.CDNID = &cdnID
	newDS.CDNName = &cdnName
	newDS.CheckPath = &checkPath
	newDS.ConsistentHashQueryParams = consistentHashQueryParams
	newDS.ConsistentHashRegex = &consistentHashRegex
	newDS.DeepCachingType = &deepCachingType
	newDS.DisplayName = &displayName
	newDS.DNSBypassCNAME = &dnsBypassCNAME
	newDS.DNSBypassIP = &dnsBypassIP
	newDS.DNSBypassIP6 = &dnsBypassIP6
	newDS.DNSBypassTTL = &dnsBypassTTL
	newDS.DSCP = &dscp
	newDS.EcsEnabled = ecsEnabled
	newDS.EdgeHeaderRewrite = &edgeHeaderRewrite
	newDS.ExampleURLs = exampleURLs
	newDS.FirstHeaderRewrite = &firstHeaderRewrite
	newDS.FQPacingRate = &fqPacingRate
	newDS.GeoLimit = &geoLimit
	newDS.GeoLimitCountries = geo
	newDS.GeoLimitRedirectURL = &geoLimitRedirectURL
	newDS.GeoProvider = &geoProvider
	newDS.GlobalMaxMBPS = &globalMaxMBPS
	newDS.GlobalMaxTPS = &globalMaxTPS
	newDS.HTTPBypassFQDN = &hTTPBypassFQDN
	newDS.ID = &id
	newDS.InfoURL = &infoURL
	newDS.InitialDispersion = &initialDispersion
	newDS.InnerHeaderRewrite = &innerHeaderRewrite
	newDS.IPV6RoutingEnabled = &ipv6RoutingEnabled
	newDS.LastHeaderRewrite = &lastHeaderRewrite
	newDS.LastUpdated = lastUpdated
	newDS.LogsEnabled = &logsEnabled
	newDS.LongDesc = &longDesc
	newDS.LongDesc1 = &longDesc1
	newDS.LongDesc2 = &longDesc2
	newDS.MatchList = &matchList
	newDS.MaxDNSAnswers = &maxDNSAnswers
	newDS.MaxOriginConnections = &maxOriginConnections
	newDS.MaxRequestHeaderBytes = &maxRequestHeaderBytes
	newDS.MidHeaderRewrite = &midHeaderRewrite
	newDS.MissLat = &missLat
	newDS.MissLong = &missLong
	newDS.MultiSiteOrigin = &multiSiteOrigin
	newDS.OriginShield = &originShield
	newDS.OrgServerFQDN = &orgServerFQDN
	newDS.ProfileDesc = &profileDesc
	newDS.ProfileID = &profileID
	newDS.ProfileName = &profileName
	newDS.Protocol = &protocol
	newDS.QStringIgnore = &qstringIgnore
	newDS.RangeRequestHandling = &rangeRequestHandling
	newDS.RangeSliceBlockSize = &rangeSliceBlockSize
	newDS.RegexRemap = &regexRemap
	newDS.RegionalGeoBlocking = &regionalGeoBlocking
	newDS.RemapText = &remapText
	newDS.RoutingName = &routingName
	newDS.ServiceCategory = &serviceCategory
	newDS.Signed = signed
	newDS.SigningAlgorithm = &signingAlgorithm
	newDS.SSLKeyVersion = &sSLKeyVersion
	newDS.Tenant = &tenant
	newDS.TenantID = &tenantID
	newDS.TLSVersions = []string{"1.0", "1.1", "1.2", "1.3"}
	newDS.Topology = &topology
	newDS.TRResponseHeaders = &trResponseHeaders
	newDS.TRRequestHeaders = &trRequestHeaders
	newDS.Type = &typ
	newDS.TypeID = &typeID
	newDS.XMLID = &xmlid
	newDS.RequiredCapabilities = requiredCapabilities

	active := false
	oldDS := DeliveryServiceNullableV30{
		DeliveryServiceV30: DeliveryServiceV30{
			DeliveryServiceNullableV15: DeliveryServiceNullableV15{
				DeliveryServiceNullableV14: DeliveryServiceNullableV14{
					DeliveryServiceNullableV13: DeliveryServiceNullableV13{
						DeliveryServiceNullableV12: DeliveryServiceNullableV12{
							DeliveryServiceNullableV11: DeliveryServiceNullableV11{
								DeliveryServiceNullableFieldsV11: DeliveryServiceNullableFieldsV11{
									Active:                   &active,
									AnonymousBlockingEnabled: &anonymousBlockingEnabled,
									CCRDNSTTL:                &cCRDNSTTL,
									CDNID:                    &cdnID,
									CDNName:                  &cdnName,
									CheckPath:                &checkPath,
									DisplayName:              &displayName,
									DNSBypassCNAME:           &dnsBypassCNAME,
									DNSBypassIP:              &dnsBypassIP,
									DNSBypassIP6:             &dnsBypassIP6,
									DNSBypassTTL:             &dnsBypassTTL,
									DSCP:                     &dscp,
									EdgeHeaderRewrite:        &edgeHeaderRewrite,
									ExampleURLs:              exampleURLs,
									GeoLimit:                 &geoLimit,
									GeoLimitCountries:        &geoStr,
									GeoLimitRedirectURL:      &geoLimitRedirectURL,
									GeoProvider:              &geoProvider,
									GlobalMaxMBPS:            &globalMaxMBPS,
									GlobalMaxTPS:             &globalMaxTPS,
									HTTPBypassFQDN:           &hTTPBypassFQDN,
									ID:                       &id,
									InfoURL:                  &infoURL,
									InitialDispersion:        &initialDispersion,
									IPV6RoutingEnabled:       &ipv6RoutingEnabled,
									LastUpdated:              lastUpdated,
									LogsEnabled:              &logsEnabled,
									LongDesc:                 &longDesc,
									LongDesc1:                &longDesc1,
									LongDesc2:                &longDesc2,
									MatchList:                &matchList,
									MaxDNSAnswers:            &maxDNSAnswers,
									MidHeaderRewrite:         &midHeaderRewrite,
									MissLat:                  &missLat,
									MissLong:                 &missLong,
									MultiSiteOrigin:          &multiSiteOrigin,
									OriginShield:             &originShield,
									OrgServerFQDN:            &orgServerFQDN,
									ProfileDesc:              &profileDesc,
									ProfileID:                &profileID,
									ProfileName:              &profileName,
									Protocol:                 &protocol,
									QStringIgnore:            &qstringIgnore,
									RangeRequestHandling:     &rangeRequestHandling,
									RegexRemap:               &regexRemap,
									RegionalGeoBlocking:      &regionalGeoBlocking,
									RemapText:                &remapText,
									RoutingName:              &routingName,
									Signed:                   signed,
									SSLKeyVersion:            &sSLKeyVersion,
									TenantID:                 &tenantID,
									Type:                     &typ,
									TypeID:                   &typeID,
									XMLID:                    &xmlid,
								},
								DeliveryServiceRemovedFieldsV11: DeliveryServiceRemovedFieldsV11{
									CacheURL: &cacheURL,
								},
							},
						},
						DeliveryServiceFieldsV13: DeliveryServiceFieldsV13{
							DeepCachingType:   &deepCachingType,
							FQPacingRate:      &fqPacingRate,
							SigningAlgorithm:  &signingAlgorithm,
							Tenant:            &tenant,
							TRResponseHeaders: &trResponseHeaders,
							TRRequestHeaders:  &trRequestHeaders,
						},
					},
					DeliveryServiceFieldsV14: DeliveryServiceFieldsV14{
						ConsistentHashQueryParams: consistentHashQueryParams,
						ConsistentHashRegex:       &consistentHashRegex,
						MaxOriginConnections:      &maxOriginConnections,
					},
				},
				DeliveryServiceFieldsV15: DeliveryServiceFieldsV15{
					EcsEnabled:          ecsEnabled,
					RangeSliceBlockSize: &rangeSliceBlockSize,
				},
			},
			DeliveryServiceFieldsV30: DeliveryServiceFieldsV30{
				FirstHeaderRewrite: &firstHeaderRewrite,
				InnerHeaderRewrite: &innerHeaderRewrite,
				LastHeaderRewrite:  &lastHeaderRewrite,
				ServiceCategory:    &serviceCategory,
				Topology:           &topology,
			},
		},
		DeliveryServiceFieldsV31: DeliveryServiceFieldsV31{
			MaxRequestHeaderBytes: &maxRequestHeaderBytes,
		},
	}

	return oldDS, newDS
}

func TestDeliveryServiceUpgradeAndDowngrade(t *testing.T) {
	oldDS, newDS := dsUpgradeAndDowngradeTestingPair()
	compareV31DSes(oldDS, newDS.DowngradeToV31(), t)

	nullableOldDS := DeliveryServiceNullableV30(oldDS)
	upgraded := nullableOldDS.UpgradeToV4()
	compareV31DSes(upgraded.DowngradeToV31(), newDS.DowngradeToV31(), t)

	downgraded := newDS.DowngradeToV31()
	upgraded = downgraded.UpgradeToV4()
	downgraded = upgraded.DowngradeToV31()
	compareV31DSes(oldDS, downgraded, t)

	upgraded = nullableOldDS.UpgradeToV4()
	downgraded = newDS.DowngradeToV31()
	tmp := downgraded.UpgradeToV4()
	downgraded = tmp.DowngradeToV31()
	compareV31DSes(upgraded.DowngradeToV31(), downgraded, t)

	if oldDS.CacheURL == nil {
		oldDS.CacheURL = new(string)
		*oldDS.CacheURL = "testquest"
	}

	upgraded = oldDS.UpgradeToV4()
	downgraded = upgraded.DowngradeToV31()
	if downgraded.CacheURL != nil {
		t.Error("Expected 'cacheurl' to be null after upgrade then downgrade because it doesn't exist in APIv4, but it wasn't")
	}

	downgraded = newDS.DowngradeToV31()
	upgraded = downgraded.UpgradeToV4()
	if upgraded.TLSVersions != nil {
		t.Errorf("Expected 'tlsVersions' to be nil after upgrade, because all TLS versions are implicitly supported for an APIv3 DS; found: %v", upgraded.TLSVersions)
	}
}

func TestDeliveryServiceV4UpgradeAndV5Downgrade(t *testing.T) {
	_, ds := dsUpgradeAndDowngradeTestingPair()
	ds.RemoveLD1AndLD2()
	upDowned := ds.Upgrade().Downgrade()
	if !reflect.DeepEqual(ds, upDowned) {
		bts, err := json.MarshalIndent(ds, "", "\t")
		if err != nil {
			t.Fatalf("failed to encode original DS after upgrade/downgrade comparison failed: %v", err)
		}
		t.Logf("original: %s", bts)
		bts, err = json.MarshalIndent(upDowned, "", "\t")
		if err != nil {
			t.Fatalf("failed to encode upgraded-then-downgraded DS after upgrade/downgrade comparison failed: %v", err)
		}
		t.Logf("upgraded-then-downgraded: %s", bts)
		t.Error("Delivery Service upgrade followed by downgrade should result in exact copy")
	}

	ds.OrgServerFQDN = new(string)
	*ds.OrgServerFQDN = "testquest.org"
	if upDowned.OrgServerFQDN != nil && *upDowned.OrgServerFQDN == *ds.OrgServerFQDN {
		t.Error("Modifying a reference in the original should leave the copy unchanged")
	}

	ds.Active = new(bool)
	upgraded := ds.Upgrade()
	if upgraded.Active != DSActiveStatePrimed {
		t.Fatalf("Upgrading an Active=false APIv4 DS should yield Active='%s', got: %s", DSActiveStatePrimed, upgraded.Active)
	}
	*ds.Active = true
	if upgraded = ds.Upgrade(); upgraded.Active != DSActiveStateActive {
		t.Fatalf("Upgrading an Active=true APIv4 DS should yield Active='%s', got: %s", DSActiveStateActive, upgraded.Active)
	}

	if downgraded := upgraded.Downgrade(); downgraded.Active == nil {
		t.Error("Unexpected nil active state after downgrading APIv5 DS to APIv4")
	} else if !*downgraded.Active {
		t.Errorf("Downgrading an Active='%s' APIv5 DS should yield Active=true, got: false", DSActiveStateActive)
	}
	upgraded.Active = DSActiveStatePrimed
	if downgraded := upgraded.Downgrade(); downgraded.Active == nil {
		t.Error("Unexpected nil active state after downgrading APIv5 DS to APIv4")
	} else if *downgraded.Active {
		t.Errorf("Downgrading an Active='%s' APIv5 DS should yield Active=false, got: true", DSActiveStatePrimed)
	}
	upgraded.Active = DSActiveStateInactive
	if downgraded := upgraded.Downgrade(); downgraded.Active == nil {
		t.Error("Unexpected nil active state after downgrading APIv5 DS to APIv4")
	} else if *downgraded.Active {
		t.Errorf("Downgrading an Active='%s' APIv5 DS should yield Active=false, got: true", DSActiveStateInactive)
	}
}

func expectOnlyWarnings(a Alerts) func(*testing.T) {
	return func(t *testing.T) {
		found := 0
		for _, alert := range a.Alerts {
			if alert.Level != WarnLevel.String() {
				found++
			}
		}

		if found > 0 {
			t.Errorf("Expected only warning-level alerts, found %d that were not warning-level: %v", found, a.Alerts)
		}
	}
}

func containsWarning(a Alerts, expected string) func(*testing.T) {
	return func(t *testing.T) {
		for _, alert := range a.Alerts {
			if alert.Level == WarnLevel.String() && alert.Text == expected {
				return
			}
		}
		t.Errorf("Expected to find a warning-level Alert containing the message '%s', but didn't: %v", expected, a.Alerts)
	}
}

const nonNilTLSVersionsWarningMessage = "setting TLS Versions that are explicitly supported may break older clients that can't use the specified versions"

func TestTLSVersionsAlerts(t *testing.T) {
	var ds DeliveryServiceV40
	alerts := ds.TLSVersionsAlerts()
	if alerts.HasAlerts() {
		t.Errorf("nil versions should not produce any warnings, but these were generated: %v", alerts.Alerts)
	}
	ds.TLSVersions = make([]string, 0, 3)
	alerts = ds.TLSVersionsAlerts()
	if alerts.HasAlerts() {
		t.Errorf("empty versions should not produce any warnings, but these were generated: %v", alerts.Alerts)
	}

	ds.Protocol = new(int)
	*ds.Protocol = 0
	expected := "old TLS version 1.0 is allowed, but newer versions 1.1, 1.2, and 1.3 are disallowed; this configuration may be insecure"
	ds.TLSVersions = append(ds.TLSVersions, TLSVersion10)
	alerts = ds.TLSVersionsAlerts()
	if len(alerts.Alerts) != 3 {
		t.Errorf("expected allowing only TLS v1.0 to generate exactly three warnings, got %d: %v", len(alerts.Alerts), alerts)
	} else {
		t.Run("only TLS version 1.0 allowed - returns warnings", expectOnlyWarnings(alerts))
		t.Run("only TLS version 1.0 allowed - has expected warning", containsWarning(alerts, expected))
		expected = "tlsVersions has no effect on Delivery Services with Protocol '0' (HTTP_ONLY)"
		t.Run("only TLS version 1.0 allowed - has warning about Protocol '0'", containsWarning(alerts, expected))
		t.Run("only TLS version 1.0 allowed - has warning about non-nil tlsVersions", containsWarning(alerts, nonNilTLSVersionsWarningMessage))
	}

	ds.Protocol = nil
	expected = "old TLS version 1.0 is allowed, but newer versions 1.1, and 1.3 are disallowed; this configuration may be insecure"
	ds.TLSVersions = append(ds.TLSVersions, TLSVersion12)
	alerts = ds.TLSVersionsAlerts()
	if len(alerts.Alerts) != 2 {
		t.Errorf("expected allowing only TLS v1.0 and 1.2 to generate exactly two warnings, got %d: %v", len(alerts.Alerts), alerts)
	} else {
		t.Run("only TLS versions 1.0 and 1.2 allowed - returns warnings", expectOnlyWarnings(alerts))
		t.Run("only TLS versions 1.0 and 1.2 allowed - has expected warning", containsWarning(alerts, expected))
		t.Run("only TLS versions 1.0 and 1.2 allowed - has warning about non-nil tlsVersions", containsWarning(alerts, nonNilTLSVersionsWarningMessage))
	}

	expected = "old TLS version 1.0 is allowed, but newer version 1.3 is disallowed; this configuration may be insecure"
	ds.TLSVersions = append(ds.TLSVersions, TLSVersion11)
	alerts = ds.TLSVersionsAlerts()
	if len(alerts.Alerts) != 2 {
		t.Errorf("expected disallowing only TLS v1.3 to generate exactly two warnings, got %d: %v", len(alerts.Alerts), alerts)
	} else {
		t.Run("only TLS version 1.3 disallowed - returns warnings", expectOnlyWarnings(alerts))
		t.Run("only TLS version 1.3 disallowed - has expected warning", containsWarning(alerts, expected))
		t.Run("only TLS version 1.3 disallowed - has warning about non-nil tlsVersions", containsWarning(alerts, nonNilTLSVersionsWarningMessage))
	}

	expected = "old TLS version 1.1 is allowed, but newer versions 1.2, and 1.3 are disallowed; this configuration may be insecure"
	ds.TLSVersions = make([]string, 0, 2)
	ds.TLSVersions = append(ds.TLSVersions, TLSVersion11)
	alerts = ds.TLSVersionsAlerts()
	if len(alerts.Alerts) != 2 {
		t.Errorf("expected allowing only TLS v1.1 to generate exactly two warnings, got %d: %v", len(alerts.Alerts), alerts)
	} else {
		t.Run("only TLS version 1.1 allowed - returns warnings", expectOnlyWarnings(alerts))
		t.Run("only TLS version 1.1 allowed - has expected warning", containsWarning(alerts, expected))
		t.Run("only TLS version 1.1 allowed - has warning about non-nil tlsVersions", containsWarning(alerts, nonNilTLSVersionsWarningMessage))
	}

	expected = "old TLS version 1.1 is allowed, but newer version 1.3 is disallowed; this configuration may be insecure"
	ds.TLSVersions = append(ds.TLSVersions, TLSVersion12)
	alerts = ds.TLSVersionsAlerts()
	if len(alerts.Alerts) != 2 {
		t.Errorf("expected allowing only TLS v1.1 and 1.2 to generate exactly two warning, got %d: %v", len(alerts.Alerts), alerts)
	} else {
		t.Run("only TLS versions 1.1 and 1.2 allowed - returns warnings", expectOnlyWarnings(alerts))
		t.Run("only TLS versions 1.1 and 1.2 allowed - has expected warning", containsWarning(alerts, expected))
		t.Run("only TLS versions 1.1 and 1.2 allowed - has warning about non-nil tlsVersions", containsWarning(alerts, nonNilTLSVersionsWarningMessage))
	}

	expected = "old TLS version 1.2 is allowed, but newer version 1.3 is disallowed; this configuration may be insecure"
	ds.TLSVersions = make([]string, 0, 4)
	ds.TLSVersions = append(ds.TLSVersions, TLSVersion12)
	alerts = ds.TLSVersionsAlerts()
	if len(alerts.Alerts) != 2 {
		t.Errorf("expected allowing only TLS v1.2 to generate exactly two warnings, got %d: %v", len(alerts.Alerts), alerts)
	} else {
		t.Run("only TLS version 1.2 allowed - returns warnings", expectOnlyWarnings(alerts))
		t.Run("only TLS version 1.2 allowed - has expected warning", containsWarning(alerts, expected))
		t.Run("only TLS version 1.2 allowed - has warning about non-nil tlsVersions", containsWarning(alerts, nonNilTLSVersionsWarningMessage))
	}

	ds.TLSVersions = append(ds.TLSVersions, TLSVersion13)
	alerts = ds.TLSVersionsAlerts()
	if len(alerts.Alerts) != 1 {
		t.Errorf("Expected allowing TLS versions 1.2 and 1.3 to generate exactly one warning, got %d: %v", len(alerts.Alerts), alerts.Alerts)
	} else {
		t.Run("only TLS versions 1.2 and 1.3 allowed - has warning about non-nil tlsVersions", containsWarning(alerts, nonNilTLSVersionsWarningMessage))
	}

	expected = "unknown TLS version '13.37' - possible typo"
	ds.TLSVersions = append(ds.TLSVersions, "13.37")
	alerts = ds.TLSVersionsAlerts()
	if len(alerts.Alerts) != 2 {
		t.Errorf("expected allowing an unknown TLS version to generate exactly two warnings, got %d: %v", len(alerts.Alerts), alerts.Alerts)
	} else {
		t.Run("unknown TLS version allowed - returns warnings", expectOnlyWarnings(alerts))
		t.Run("unknown TLS version allowed - has expected warning", containsWarning(alerts, expected))
		t.Run("unknown TLS version allowed - has warning about non-nil tlsVersions", containsWarning(alerts, nonNilTLSVersionsWarningMessage))
	}

	ds.TLSVersions = append(ds.TLSVersions, TLSVersion10)
	alerts = ds.TLSVersionsAlerts()
	if len(alerts.Alerts) != 3 {
		t.Errorf("expected allowing an unknown TLS version and disallowing only version 1.1 to generate exactly three warnings, got %d: %v", len(alerts.Alerts), alerts.Alerts)
	} else {
		t.Run("unknown TLS version and disallowed v1.1 - returns warnings", expectOnlyWarnings(alerts))
		t.Run("unknown TLS version and disallowed v1.1 - has expected warning", containsWarning(alerts, expected))
		expected = "old TLS version 1.0 is allowed, but newer version 1.1 is disallowed; this configuration may be insecure"
		t.Run("unknown TLS version and disallowed v1.1 - has second expected warning", containsWarning(alerts, expected))
		t.Run("unknown TLS version and disallowed v1.1 - has warning about non-nil tlsVersions", containsWarning(alerts, nonNilTLSVersionsWarningMessage))
	}
}

func BenchmarkTLSVersionsAlerts(b *testing.B) {
	versions := make([]string, 0, 101)
	for major := 1; major <= 10; major++ {
		for minor := 0; minor < 10; minor++ {
			if major == 1 && minor == 1 {
				// skip this one to force generating some more warnings
				continue
			}
			versions = append(versions, fmt.Sprintf("%d.%d", major, minor))
		}
	}
	ds := DeliveryServiceV40{TLSVersions: versions}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ds.TLSVersionsAlerts()
	}
}
