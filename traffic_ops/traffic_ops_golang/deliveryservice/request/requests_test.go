package request

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import (
	"net/http"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/auth"

	"github.com/jmoiron/sqlx"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestInsert(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("opening mock database: %v", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	mock.ExpectBegin()
	inf := api.Info{
		Params:    nil,
		IntParams: nil,
		User: &auth.CurrentUser{
			UserName:     "testUser",
			ID:           1,
			PrivLevel:    10,
			TenantID:     1,
			Role:         1,
			RoleName:     "testRole",
			Capabilities: nil,
			UCDN:         "",
		},
		ReqID:    0,
		Version:  nil,
		Tx:       db.MustBegin(),
		CancelTx: nil,
		Vault:    nil,
		Config:   nil,
	}
	dsr := tc.DeliveryServiceRequestV5{
		Assignee:       util.StrPtr("assignee"),
		AssigneeID:     util.IntPtr(25),
		Author:         "test",
		AuthorID:       util.IntPtr(35),
		ChangeType:     tc.DSRChangeTypeUpdate,
		CreatedAt:      time.Now(),
		LastEditedBy:   "test",
		LastEditedByID: util.IntPtr(35),
		LastUpdated:    time.Now(),
		Original:       nil,
		Requested:      &tc.DeliveryServiceV5{},
		Status:         tc.RequestStatusDraft,
		XMLID:          "dsXMLID",
	}

	rows := sqlmock.NewRows([]string{"id", "last_updated", "created_at"})
	rows.AddRow(1, time.Now(), time.Now())
	mock.ExpectQuery("INSERT INTO deliveryservice_request*").WillReturnRows(rows)

	rows2 := sqlmock.NewRows([]string{"active", "anonymous_blocking_enabled", "ccr_dns_ttl", "cdn_id", "cdnname", "check_path",
		"consistent_hash_regex", "deep_caching_type", "display_name", "dns_bypass_cname", "dns_bypass_ip", "dns_bypass_ip6",
		"dns_bypass_ttl", "dscp", "ecs_enabled", "edge_header_rewrite", "first_header_rewrite", "geolimit_redirect_url",
		"geo_limit", "geo_limit_countries", "geo_provider", "global_max_mbps", "global_max_tps", "fq_pacing_rate", "http_bypass_fqdn",
		"id", "info_url", "initial_dispersion", "inner_header_rewrite", "ipv6_routing_enabled", "last_header_rewrite", "last_updated",
		"logs_enabled", "long_desc", "long_desc_1", "long_desc_2", "max_dns_answers", "max_origin_connections", "max_request_header_bytes",
		"mid_header_rewrite", "miss_lat", "miss_long", "multi_site_origin", "org_server_fqdn ", "origin_shield", "profileid", "profile_name",
		"profile_description", "protocol", "qstring_ignore", "query_keys", "range_request_handling", "regex_remap", "regional", "regional_geo_blocking",
		"remap_text", "required_capabilities", "routing_name", "service_category", "signing_algorithm", "range_slice_block_size", "ssl_key_version", "tenant_id",
		"name", "tls_versions", "topology", "tr_request_headers", "tr_response_headers", "name", "type_id", "xml_id", "cdn_domain",
	})
	ds := tc.DeliveryServiceV5{
		Active:                   tc.DeliveryServiceActiveState("PRIMED"),
		AnonymousBlockingEnabled: false,
		CCRDNSTTL:                util.IntPtr(20),
		CDNID:                    11,
		CDNName:                  util.StrPtr("testCDN"),
		CheckPath:                util.StrPtr("blah"),
		ConsistentHashRegex:      nil,
		DeepCachingType:          tc.DeepCachingTypeNever,
		DisplayName:              "ds",
		DNSBypassCNAME:           nil,
		DNSBypassIP:              nil,
		DNSBypassIP6:             nil,
		DNSBypassTTL:             nil,
		DSCP:                     0,
		EcsEnabled:               false,
		EdgeHeaderRewrite:        nil,
		FirstHeaderRewrite:       nil,
		GeoLimitRedirectURL:      nil,
		GeoLimit:                 0,
		GeoLimitCountries:        nil,
		GeoProvider:              0,
		GlobalMaxMBPS:            nil,
		GlobalMaxTPS:             nil,
		FQPacingRate:             nil,
		HTTPBypassFQDN:           nil,
		ID:                       util.IntPtr(1),
		InfoURL:                  nil,
		InitialDispersion:        util.IntPtr(1),
		InnerHeaderRewrite:       nil,
		IPV6RoutingEnabled:       util.BoolPtr(true),
		LastHeaderRewrite:        nil,
		LastUpdated:              time.Now(),
		LogsEnabled:              true,
		LongDesc:                 "",
		MaxDNSAnswers:            util.IntPtr(5),
		MaxOriginConnections:     util.IntPtr(2),
		MaxRequestHeaderBytes:    util.IntPtr(0),
		MidHeaderRewrite:         nil,
		MissLat:                  util.FloatPtr(0.0),
		MissLong:                 util.FloatPtr(0.0),
		MultiSiteOrigin:          false,
		OrgServerFQDN:            util.StrPtr("http://1.2.3.4"),
		OriginShield:             nil,
		ProfileID:                util.IntPtr(99),
		ProfileName:              util.StrPtr("profile99"),
		ProfileDesc:              nil,
		Protocol:                 util.IntPtr(1),
		QStringIgnore:            nil,
		RangeRequestHandling:     nil,
		RegexRemap:               nil,
		Regional:                 false,
		RegionalGeoBlocking:      false,
		RemapText:                nil,
		RequiredCapabilities:     nil,
		RoutingName:              "",
		ServiceCategory:          nil,
		SigningAlgorithm:         nil,
		RangeSliceBlockSize:      nil,
		SSLKeyVersion:            nil,
		TenantID:                 100,
		Tenant:                   util.StrPtr("tenant100"),
		TLSVersions:              nil,
		Topology:                 nil,
		TRRequestHeaders:         nil,
		TRResponseHeaders:        nil,
		Type:                     util.StrPtr("type101"),
		TypeID:                   101,
		XMLID:                    "dsXMLID",
	}

	rows2.AddRow(
		ds.Active,
		ds.AnonymousBlockingEnabled,
		ds.CCRDNSTTL,
		ds.CDNID,
		ds.CDNName,
		ds.CheckPath,
		ds.ConsistentHashRegex,
		ds.DeepCachingType,
		ds.DisplayName,
		ds.DNSBypassCNAME,
		ds.DNSBypassIP,
		ds.DNSBypassIP6,
		ds.DNSBypassTTL,
		ds.DSCP,
		ds.EcsEnabled,
		ds.EdgeHeaderRewrite,
		ds.FirstHeaderRewrite,
		ds.GeoLimitRedirectURL,
		ds.GeoLimit,
		nil,
		ds.GeoProvider,
		ds.GlobalMaxMBPS,
		ds.GlobalMaxTPS,
		ds.FQPacingRate,
		ds.HTTPBypassFQDN,
		ds.ID,
		ds.InfoURL,
		ds.InitialDispersion,
		ds.InnerHeaderRewrite,
		ds.IPV6RoutingEnabled,
		ds.LastHeaderRewrite,
		ds.LastUpdated,
		ds.LogsEnabled,
		ds.LongDesc,
		ds.LongDesc,
		ds.LongDesc,
		ds.MaxDNSAnswers,
		ds.MaxOriginConnections,
		ds.MaxRequestHeaderBytes,
		ds.MidHeaderRewrite,
		ds.MissLat,
		ds.MissLong,
		ds.MultiSiteOrigin,
		ds.OrgServerFQDN,
		ds.OriginShield,
		ds.ProfileID,
		ds.ProfileName,
		ds.ProfileDesc,
		ds.Protocol,
		ds.QStringIgnore,
		nil,
		ds.RangeRequestHandling,
		ds.RegexRemap,
		ds.Regional,
		ds.RegionalGeoBlocking,
		ds.RemapText,
		nil,
		ds.RoutingName,
		ds.ServiceCategory,
		ds.SigningAlgorithm,
		ds.RangeSliceBlockSize,
		ds.SSLKeyVersion,
		ds.TenantID,
		ds.Tenant,
		nil,
		ds.Topology,
		ds.TRRequestHeaders,
		ds.TRResponseHeaders,
		ds.Type,
		ds.TypeID,
		ds.XMLID,
		"cdn_domain_name")

	mock.ExpectQuery("SELECT ds.active*").WillReturnRows(rows2)

	rows3 := sqlmock.NewRows([]string{
		"ds_name",
		"type",
		"pattern",
		"set_number",
	})
	rows3.AddRow(
		"dsXMLID",
		"HOST_REGEXP",
		".*\\.dsXMLID\\..*",
		0)
	mock.ExpectQuery("SELECT ds.xml_id as ds_name*").WillReturnRows(rows3)

	sc, userErr, sysErr := insert(&dsr, &inf)

	if userErr != nil || sysErr != nil {
		t.Fatalf("expected no error, but got userErr: %v, sysErr: %v", userErr, sysErr)
	}
	if sc != http.StatusOK {
		t.Fatalf("expected a 200 status code, but got %d", sc)
	}
	if dsr.Original == nil {
		t.Fatalf("expected original to be a valid delivery service, but got nothing")
	}
	if dsr.Original.XMLID != "dsXMLID" {
		t.Fatalf("expected original to have a DS with XMLID 'dsXMLID', but got %s", dsr.Original.XMLID)
	}

}
func TestGetOriginals(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("opening mock database: %v", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	mock.ExpectBegin()
	ID := 66
	ids := []int{ID}
	needOriginals := make(map[int][]*tc.DeliveryServiceRequestV5)
	dsr := tc.DeliveryServiceRequestV5{
		Assignee:       util.StrPtr("assignee"),
		AssigneeID:     util.IntPtr(25),
		Author:         "test",
		AuthorID:       util.IntPtr(35),
		ChangeType:     tc.DSRChangeTypeUpdate,
		CreatedAt:      time.Now(),
		ID:             util.IntPtr(1),
		LastEditedBy:   "test",
		LastEditedByID: util.IntPtr(35),
		LastUpdated:    time.Now(),
		Original:       nil,
		Requested:      nil,
		Status:         tc.RequestStatusDraft,
		XMLID:          "dsXMLID",
	}
	needOriginals[ID] = []*tc.DeliveryServiceRequestV5{&dsr}

	rows := sqlmock.NewRows([]string{"active", "anonymous_blocking_enabled", "ccr_dns_ttl", "cdn_id", "cdnname", "check_path",
		"consistent_hash_regex", "deep_caching_type", "display_name", "dns_bypass_cname", "dns_bypass_ip", "dns_bypass_ip6",
		"dns_bypass_ttl", "dscp", "ecs_enabled", "edge_header_rewrite", "first_header_rewrite", "geolimit_redirect_url",
		"geo_limit", "geo_limit_countries", "geo_provider", "global_max_mbps", "global_max_tps", "fq_pacing_rate", "http_bypass_fqdn",
		"id", "info_url", "initial_dispersion", "inner_header_rewrite", "ipv6_routing_enabled", "last_header_rewrite", "last_updated",
		"logs_enabled", "long_desc", "long_desc_1", "long_desc_2", "max_dns_answers", "max_origin_connections", "max_request_header_bytes",
		"mid_header_rewrite", "miss_lat", "miss_long", "multi_site_origin", "org_server_fqdn ", "origin_shield", "profileid", "profile_name",
		"profile_description", "protocol", "qstring_ignore", "query_keys", "range_request_handling", "regex_remap", "regional", "regional_geo_blocking",
		"remap_text", "required_capabilities", "routing_name", "service_category", "signing_algorithm", "range_slice_block_size", "ssl_key_version", "tenant_id",
		"name", "tls_versions", "topology", "tr_request_headers", "tr_response_headers", "name", "type_id", "xml_id", "cdn_domain",
	})
	ds := tc.DeliveryServiceV5{
		Active:                   tc.DeliveryServiceActiveState("PRIMED"),
		AnonymousBlockingEnabled: false,
		CCRDNSTTL:                util.IntPtr(20),
		CDNID:                    11,
		CDNName:                  util.StrPtr("testCDN"),
		CheckPath:                util.StrPtr("blah"),
		ConsistentHashRegex:      nil,
		DeepCachingType:          tc.DeepCachingTypeNever,
		DisplayName:              "ds",
		DNSBypassCNAME:           nil,
		DNSBypassIP:              nil,
		DNSBypassIP6:             nil,
		DNSBypassTTL:             nil,
		DSCP:                     0,
		EcsEnabled:               false,
		EdgeHeaderRewrite:        nil,
		FirstHeaderRewrite:       nil,
		GeoLimitRedirectURL:      nil,
		GeoLimit:                 0,
		GeoLimitCountries:        nil,
		GeoProvider:              0,
		GlobalMaxMBPS:            nil,
		GlobalMaxTPS:             nil,
		FQPacingRate:             nil,
		HTTPBypassFQDN:           nil,
		ID:                       util.IntPtr(ID),
		InfoURL:                  nil,
		InitialDispersion:        util.IntPtr(1),
		InnerHeaderRewrite:       nil,
		IPV6RoutingEnabled:       util.BoolPtr(true),
		LastHeaderRewrite:        nil,
		LastUpdated:              time.Now(),
		LogsEnabled:              true,
		LongDesc:                 "",
		MaxDNSAnswers:            util.IntPtr(5),
		MaxOriginConnections:     util.IntPtr(2),
		MaxRequestHeaderBytes:    util.IntPtr(0),
		MidHeaderRewrite:         nil,
		MissLat:                  util.FloatPtr(0.0),
		MissLong:                 util.FloatPtr(0.0),
		MultiSiteOrigin:          false,
		OrgServerFQDN:            util.StrPtr("http://1.2.3.4"),
		OriginShield:             nil,
		ProfileID:                util.IntPtr(99),
		ProfileName:              util.StrPtr("profile99"),
		ProfileDesc:              nil,
		Protocol:                 util.IntPtr(1),
		QStringIgnore:            nil,
		RangeRequestHandling:     nil,
		RegexRemap:               nil,
		Regional:                 false,
		RegionalGeoBlocking:      false,
		RemapText:                nil,
		RequiredCapabilities:     nil,
		RoutingName:              "",
		ServiceCategory:          nil,
		SigningAlgorithm:         nil,
		RangeSliceBlockSize:      nil,
		SSLKeyVersion:            nil,
		TenantID:                 100,
		Tenant:                   util.StrPtr("tenant100"),
		TLSVersions:              nil,
		Topology:                 nil,
		TRRequestHeaders:         nil,
		TRResponseHeaders:        nil,
		Type:                     util.StrPtr("type101"),
		TypeID:                   101,
		XMLID:                    "dsXMLID",
	}
	rows.AddRow(
		ds.Active,
		ds.AnonymousBlockingEnabled,
		ds.CCRDNSTTL,
		ds.CDNID,
		ds.CDNName,
		ds.CheckPath,
		ds.ConsistentHashRegex,
		ds.DeepCachingType,
		ds.DisplayName,
		ds.DNSBypassCNAME,
		ds.DNSBypassIP,
		ds.DNSBypassIP6,
		ds.DNSBypassTTL,
		ds.DSCP,
		ds.EcsEnabled,
		ds.EdgeHeaderRewrite,
		ds.FirstHeaderRewrite,
		ds.GeoLimitRedirectURL,
		ds.GeoLimit,
		nil,
		ds.GeoProvider,
		ds.GlobalMaxMBPS,
		ds.GlobalMaxTPS,
		ds.FQPacingRate,
		ds.HTTPBypassFQDN,
		ds.ID,
		ds.InfoURL,
		ds.InitialDispersion,
		ds.InnerHeaderRewrite,
		ds.IPV6RoutingEnabled,
		ds.LastHeaderRewrite,
		ds.LastUpdated,
		ds.LogsEnabled,
		ds.LongDesc,
		ds.LongDesc,
		ds.LongDesc,
		ds.MaxDNSAnswers,
		ds.MaxOriginConnections,
		ds.MaxRequestHeaderBytes,
		ds.MidHeaderRewrite,
		ds.MissLat,
		ds.MissLong,
		ds.MultiSiteOrigin,
		ds.OrgServerFQDN,
		ds.OriginShield,
		ds.ProfileID,
		ds.ProfileName,
		ds.ProfileDesc,
		ds.Protocol,
		ds.QStringIgnore,
		nil,
		ds.RangeRequestHandling,
		ds.RegexRemap,
		ds.Regional,
		ds.RegionalGeoBlocking,
		ds.RemapText,
		nil,
		ds.RoutingName,
		ds.ServiceCategory,
		ds.SigningAlgorithm,
		ds.RangeSliceBlockSize,
		ds.SSLKeyVersion,
		ds.TenantID,
		ds.Tenant,
		nil,
		ds.Topology,
		ds.TRRequestHeaders,
		ds.TRResponseHeaders,
		ds.Type,
		ds.TypeID,
		ds.XMLID,
		"cdn_domain_name")

	mock.ExpectQuery("SELECT ds.active*").WillReturnRows(rows)

	rows2 := sqlmock.NewRows([]string{
		"ds_name",
		"type",
		"pattern",
		"set_number",
	})
	rows2.AddRow(
		"dsXMLID",
		"HOST_REGEXP",
		".*\\.dsXMLID\\..*",
		0)
	mock.ExpectQuery("SELECT ds.xml_id as ds_name*").WillReturnRows(rows2)

	if needOriginals[ID][0].Original != nil {
		t.Errorf("expected original to be initially empty")
	}
	sc, userErr, sysErr := getOriginals(ids, db.MustBegin(), needOriginals)
	if userErr != nil || sysErr != nil {
		t.Fatalf("expected no error, but got userErr: %v, sysErr: %v", userErr, sysErr)
	}
	if sc != http.StatusOK {
		t.Fatalf("expected a 200 status code, but got %d", sc)
	}
	if needOriginals[ID][0].Original == nil {
		t.Fatalf("expected original to be a valid delivery service, but got nothing")
	}
	if needOriginals[ID][0].Original.XMLID != "dsXMLID" {
		t.Fatalf("expected original to have a DS with XMLID 'dsXMLID', but got %s", needOriginals[ID][0].Original.XMLID)
	}
}

func TestGetAssignee(t *testing.T) {
	req := assignmentRequest{
		AssigneeID: nil,
		Assignee:   nil,
	}

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("opening mock database: %v", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	// check simple case, no Assignee or ID means no change
	mock.ExpectBegin()
	_, _, userErr, sysErr := getAssignee(&req, "test", db.MustBegin().Tx)
	if userErr != nil {
		t.Errorf("unexpected user error: %v", userErr)
	}
	if sysErr != nil {
		t.Errorf("unexpected system error: %v", sysErr)
	}
	if req.AssigneeID != nil {
		t.Errorf("assignee ID was somehow set to: %d", *req.AssigneeID)
	}
	if req.Assignee != nil {
		t.Errorf("assignee was somehow set to: %s", *req.Assignee)
	}

	expectID := 12
	expectName := "test assignee"

	req.Assignee = &expectName

	rows := sqlmock.NewRows([]string{"id"})
	rows.AddRow(expectID)
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT id").WillReturnRows(rows)

	// check case where getting Assignee ID from username
	_, _, userErr, sysErr = getAssignee(&req, "test", db.MustBegin().Tx)
	if userErr != nil {
		t.Errorf("unexpected user error: %v", userErr)
	}
	if sysErr != nil {
		t.Errorf("unexpected system error: %v", sysErr)
	}

	if req.Assignee == nil {
		t.Error("Expected assignee to not be nil after getting assignee")
	} else if *req.Assignee != expectName {
		t.Errorf("Incorrect assignee; expected: '%s', got: '%s'", expectName, *req.Assignee)
	}

	if req.AssigneeID == nil {
		t.Error("Expected assignee ID to not be nil after getting assignee")
	} else if *req.AssigneeID != expectID {
		t.Errorf("Incorrect assignee ID; expected: %d, got: %d", expectID, *req.AssigneeID)
	}

	req.Assignee = nil
	req.AssigneeID = &expectID

	rows = sqlmock.NewRows([]string{"username"})
	rows.AddRow(expectName)
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT username").WillReturnRows(rows)

	// check case where getting username from Assignee ID
	_, _, userErr, sysErr = getAssignee(&req, "test", db.MustBegin().Tx)
	if userErr != nil {
		t.Errorf("unexpected user error: %v", userErr)
	}
	if sysErr != nil {
		t.Errorf("unexpected system error: %v", sysErr)
	}

	if req.Assignee == nil {
		t.Error("Expected assignee to not be nil after getting assignee")
	} else if *req.Assignee != expectName {
		t.Errorf("Incorrect assignee; expected: '%s', got: '%s'", expectName, *req.Assignee)
	}

	if req.AssigneeID == nil {
		t.Error("Expected assignee ID to not be nil after getting assignee")
	} else if *req.AssigneeID != expectID {
		t.Errorf("Incorrect assignee ID; expected: %d, got: %d", expectID, *req.AssigneeID)
	}

	req.Assignee = new(string)
	*req.Assignee = expectName + " - but not actually"
	req.AssigneeID = &expectID
	rows = sqlmock.NewRows([]string{"username"})
	rows.AddRow(expectName)
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT username").WillReturnRows(rows)

	// check that Assignee ID has precedence over Assignee
	_, _, userErr, sysErr = getAssignee(&req, "test", db.MustBegin().Tx)
	if userErr != nil {
		t.Errorf("unexpected user error: %v", userErr)
	}
	if sysErr != nil {
		t.Errorf("unexpected system error: %v", sysErr)
	}

	if req.Assignee == nil {
		t.Error("Expected assignee to not be nil after getting assignee")
	} else if *req.Assignee != expectName {
		t.Errorf("Incorrect assignee; expected: '%s', got: '%s'", expectName, *req.Assignee)
	}

	if req.AssigneeID == nil {
		t.Error("Expected assignee ID to not be nil after getting assignee")
	} else if *req.AssigneeID != expectID {
		t.Errorf("Incorrect assignee ID; expected: %d, got: %d", expectID, *req.AssigneeID)
	}

	req.Assignee = nil
	req.AssigneeID = &expectID
	rows = sqlmock.NewRows([]string{"username"})
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT username").WillReturnRows(rows)

	// check that looking for ID of non-existent username is an error
	_, _, userErr, sysErr = getAssignee(&req, "test", db.MustBegin().Tx)
	if userErr == nil {
		t.Error("Expected a user error, but didn't get one")
	}

	req.Assignee = &expectName
	req.AssigneeID = nil
	rows = sqlmock.NewRows([]string{"id"})
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT id").WillReturnRows(rows)

	// check that looking for username of non-existent Assignee is an error
	_, _, userErr, sysErr = getAssignee(&req, "test", db.MustBegin().Tx)
	if userErr == nil {
		t.Error("Expected a user error, but didn't get one")
	}
}
