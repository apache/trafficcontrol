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

package v13

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/auth"
)

var (
	db *sql.DB
)

// OpenConnection ...
func OpenConnection() (*sql.DB, error) {
	var err error
	sslStr := "require"
	if !Config.TrafficOpsDB.SSL {
		sslStr = "disable"
	}

	db, err = sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s", Config.TrafficOpsDB.User, Config.TrafficOpsDB.Password, Config.TrafficOpsDB.Hostname, Config.TrafficOpsDB.Name, sslStr))

	if err != nil {
		log.Errorf("opening database: %v\n", err)
		return nil, fmt.Errorf("transaction failed: %s", err)
	}
	return db, err
}

// SetupTestData ...
func SetupTestData(*sql.DB) error {
	var err error

	err = SetupTenants(db)
	if err != nil {
		fmt.Printf("\nError setting up tenants %s - %s, %v\n", Config.TrafficOps.URL, Config.TrafficOps.Users.Admin, err)
		os.Exit(1)
	}

	err = SetupRoles(db)
	if err != nil {
		fmt.Printf("\nError setting up roles %s - %s, %v\n", Config.TrafficOps.URL, Config.TrafficOps.Users.Admin, err)
		os.Exit(1)
	}

	err = SetupTmusers(db)
	if err != nil {
		fmt.Printf("\nError setting up tm_user %s - %s, %v\n", Config.TrafficOps.URL, Config.TrafficOps.Users.Admin, err)
		os.Exit(1)
	}

	return err
}

// SetupRoles ...
func SetupRoles(db *sql.DB) error {

	sqlStmt := `
INSERT INTO role (id, name, description, priv_level) VALUES (1, 'disallowed','Block all access',0) ON CONFLICT DO NOTHING;
INSERT INTO role (id, name, description, priv_level) VALUES (2, 'read-only user','Block all access', 10) ON CONFLICT DO NOTHING;
INSERT INTO role (id, name, description, priv_level) VALUES (3, 'operations','Block all access', 20) ON CONFLICT DO NOTHING;
INSERT INTO role (id, name, description, priv_level) VALUES (4, 'admin','super-user', 30) ON CONFLICT DO NOTHING;
INSERT INTO role (id, name, description, priv_level) VALUES (5, 'portal','Portal User', 2) ON CONFLICT DO NOTHING;
INSERT INTO role (id, name, description, priv_level) VALUES (7, 'federation','Role for Secondary CZF', 15) ON CONFLICT DO NOTHING;
`
	err := execSQL(db, sqlStmt, "role")
	if err != nil {
		return fmt.Errorf("exec failed %v", err)
	}
	return nil
}

// SetupTmusers ...
func SetupTmusers(db *sql.DB) error {

	var err error
	encryptedPassword, err := auth.DerivePassword(Config.TrafficOps.UserPassword)
	if err != nil {
		return fmt.Errorf("password encryption failed %v", err)
	}

	// Creates users in different tenants
	sqlStmt := `
INSERT INTO tm_user (username, local_passwd, confirm_local_passwd, role, tenant_id) VALUES ('` + Config.TrafficOps.Users.Disallowed + `','` + encryptedPassword + `','` + encryptedPassword + `', 1, 3);
INSERT INTO tm_user (username, local_passwd, confirm_local_passwd, role, tenant_id) VALUES ('` + Config.TrafficOps.Users.ReadOnly + `','` + encryptedPassword + `','` + encryptedPassword + `', 2, 3);
INSERT INTO tm_user (username, local_passwd, confirm_local_passwd, role, tenant_id) VALUES ('` + Config.TrafficOps.Users.Operations + `','` + encryptedPassword + `','` + encryptedPassword + `', 3, 3);
INSERT INTO tm_user (username, local_passwd, confirm_local_passwd, role, tenant_id) VALUES ('` + Config.TrafficOps.Users.Admin + `','` + encryptedPassword + `','` + encryptedPassword + `', 4, 2);
INSERT INTO tm_user (username, local_passwd, confirm_local_passwd, role, tenant_id) VALUES ('` + Config.TrafficOps.Users.Portal + `','` + encryptedPassword + `','` + encryptedPassword + `', 5, 3);
INSERT INTO tm_user (username, local_passwd, confirm_local_passwd, role, tenant_id) VALUES ('` + Config.TrafficOps.Users.Federation + `','` + encryptedPassword + `','` + encryptedPassword + `', 7, 3);
`
	err = execSQL(db, sqlStmt, "tm_user")
	if err != nil {
		return fmt.Errorf("exec failed %v", err)
	}
	return nil
}

// SetupTenants ...
func SetupTenants(db *sql.DB) error {

	sqlStmt := `
INSERT INTO tenant (id, name, active, parent_id, last_updated) VALUES (1, 'root', true, null, '2018-01-19 19:01:21.327262');
INSERT INTO tenant (id, name, active, parent_id, last_updated) VALUES (2, 'grandparent tenant', true, 1, '2018-01-19 19:01:21.327262');
INSERT INTO tenant (id, name, active, parent_id, last_updated) VALUES (3, 'parent tenant', true, 2, '2018-01-19 19:01:21.327262');
INSERT INTO tenant (id, name, active, parent_id, last_updated) VALUES (4, 'child tenant', true, 3, '2018-01-19 19:01:21.327262');
`
	err := execSQL(db, sqlStmt, "tenant")
	if err != nil {
		return fmt.Errorf("exec failed %v", err)
	}
	return nil
}

// SetupDeliveryServices ...
func SetupDeliveryServices(db *sql.DB) error {

	sqlStmt := `
INSERT INTO deliveryservice (id, xml_id, active, dscp, signing_algorithm, qstring_ignore, geo_limit, http_bypass_fqdn, dns_bypass_ip, dns_bypass_ip6, dns_bypass_ttl, org_server_fqdn, type, profile, cdn_id, ccr_dns_ttl, global_max_mbps, global_max_tps, long_desc, long_desc_1, long_desc_2, max_dns_answers, info_url, miss_lat, miss_long, check_path, last_updated, protocol, ssl_key_version, ipv6_routing_enabled, range_request_handling, edge_header_rewrite, origin_shield, mid_header_rewrite, regex_remap, cacheurl, remap_text, multi_site_origin, display_name, tr_response_headers, initial_dispersion, dns_bypass_cname, tr_request_headers, regional_geo_blocking, geo_provider, geo_limit_countries, logs_enabled, multi_site_origin_algorithm, geolimit_redirect_url, tenant_id, routing_name) VALUES (100, 'test-ds1', true, 40, null, 0, 0, '', '', '', -1, 'http://test-ds1.edge', 21, 100, 100, 3600, 0, 0, 'test-ds1 long_desc', 'test-ds1 long_desc_1', 'test-ds1 long_desc_2', 0, 'http://test-ds1.edge/info_url.html', 41.881944, -87.627778, '/crossdomain.xml', '2018-01-19 21:19:32.217466', 1, 0, true, 0, null, null, null, null, null, null, false, 'test-ds1-displayname', null, 1, null, null, true, 0, null, true, null, null, 2, 'foo');
INSERT INTO deliveryservice (id, xml_id, active, dscp, signing_algorithm, qstring_ignore, geo_limit, http_bypass_fqdn, dns_bypass_ip, dns_bypass_ip6, dns_bypass_ttl, org_server_fqdn, type, profile, cdn_id, ccr_dns_ttl, global_max_mbps, global_max_tps, long_desc, long_desc_1, long_desc_2, max_dns_answers, info_url, miss_lat, miss_long, check_path, last_updated, protocol, ssl_key_version, ipv6_routing_enabled, range_request_handling, edge_header_rewrite, origin_shield, mid_header_rewrite, regex_remap, cacheurl, remap_text, multi_site_origin, display_name, tr_response_headers, initial_dispersion, dns_bypass_cname, tr_request_headers, regional_geo_blocking, geo_provider, geo_limit_countries, logs_enabled, multi_site_origin_algorithm, geolimit_redirect_url, tenant_id, routing_name) VALUES (1000, 'steering-target-ds1', true, 40, null, 0, 0, '', 'hokeypokey', null, 10, 'http://target-ds1.edge', 21, 100, 100, 3600, 0, 0, 'target-ds1 long_desc', 'target-ds1 long_desc_1', 'target-ds1 long_desc_2', 0, 'http://target-ds1.edge/info_url.html', 41.881944, -87.627778, '/crossdomain.xml', '2018-01-19 21:19:32.226858', 1, 0, true, 0, null, null, null, null, null, null, false, 'target-ds1-displayname', null, 1, null, null, true, 0, null, false, null, null, 2, 'foo');
INSERT INTO deliveryservice (id, xml_id, active, dscp, signing_algorithm, qstring_ignore, geo_limit, http_bypass_fqdn, dns_bypass_ip, dns_bypass_ip6, dns_bypass_ttl, org_server_fqdn, type, profile, cdn_id, ccr_dns_ttl, global_max_mbps, global_max_tps, long_desc, long_desc_1, long_desc_2, max_dns_answers, info_url, miss_lat, miss_long, check_path, last_updated, protocol, ssl_key_version, ipv6_routing_enabled, range_request_handling, edge_header_rewrite, origin_shield, mid_header_rewrite, regex_remap, cacheurl, remap_text, multi_site_origin, display_name, tr_response_headers, initial_dispersion, dns_bypass_cname, tr_request_headers, regional_geo_blocking, geo_provider, geo_limit_countries, logs_enabled, multi_site_origin_algorithm, geolimit_redirect_url, tenant_id, routing_name) VALUES (1100, 'steering-target-ds2', true, 40, null, 0, 0, '', 'hokeypokey', null, 10, 'http://target-ds2.edge', 21, 100, 100, 3600, 0, 0, 'target-ds2 long_desc', 'target-ds2 long_desc_1', 'target-ds2 long_desc_2', 0, 'http://target-ds2.edge/info_url.html', 41.881944, -87.627778, '/crossdomain.xml', '2018-01-19 21:19:32.235025', 1, 0, true, 0, null, null, null, null, null, null, false, 'target-ds2-displayname', null, 1, null, null, true, 0, null, false, null, null, 2, 'foo');
INSERT INTO deliveryservice (id, xml_id, active, dscp, signing_algorithm, qstring_ignore, geo_limit, http_bypass_fqdn, dns_bypass_ip, dns_bypass_ip6, dns_bypass_ttl, org_server_fqdn, type, profile, cdn_id, ccr_dns_ttl, global_max_mbps, global_max_tps, long_desc, long_desc_1, long_desc_2, max_dns_answers, info_url, miss_lat, miss_long, check_path, last_updated, protocol, ssl_key_version, ipv6_routing_enabled, range_request_handling, edge_header_rewrite, origin_shield, mid_header_rewrite, regex_remap, cacheurl, remap_text, multi_site_origin, display_name, tr_response_headers, initial_dispersion, dns_bypass_cname, tr_request_headers, regional_geo_blocking, geo_provider, geo_limit_countries, logs_enabled, multi_site_origin_algorithm, geolimit_redirect_url, tenant_id, routing_name) VALUES (1200, 'steering-target-ds3', true, 40, null, 0, 0, '', 'hokeypokey', null, 10, 'http://target-ds3.edge', 21, 100, 100, 3600, 0, 0, 'target-ds3 long_desc', 'target-ds3 long_desc_1', 'target-ds3 long_desc_2', 0, 'http://target-ds3.edge/info_url.html', 41.881944, -87.627778, '/crossdomain.xml', '2018-01-19 21:19:32.241327', 1, 0, true, 0, null, null, null, null, null, null, false, 'target-ds3-displayname', null, 1, null, null, true, 0, null, false, null, null, 3, 'foo');
INSERT INTO deliveryservice (id, xml_id, active, dscp, signing_algorithm, qstring_ignore, geo_limit, http_bypass_fqdn, dns_bypass_ip, dns_bypass_ip6, dns_bypass_ttl, org_server_fqdn, type, profile, cdn_id, ccr_dns_ttl, global_max_mbps, global_max_tps, long_desc, long_desc_1, long_desc_2, max_dns_answers, info_url, miss_lat, miss_long, check_path, last_updated, protocol, ssl_key_version, ipv6_routing_enabled, range_request_handling, edge_header_rewrite, origin_shield, mid_header_rewrite, regex_remap, cacheurl, remap_text, multi_site_origin, display_name, tr_response_headers, initial_dispersion, dns_bypass_cname, tr_request_headers, regional_geo_blocking, geo_provider, geo_limit_countries, logs_enabled, multi_site_origin_algorithm, geolimit_redirect_url, tenant_id, routing_name) VALUES (1300, 'steering-target-ds4', true, 40, null, 0, 0, '', 'hokeypokey', null, 10, 'http://target-ds4.edge', 21, 100, 100, 3600, 0, 0, 'target-ds4 long_desc', 'target-ds4 long_desc_1', 'target-ds4 long_desc_2', 0, 'http://target-ds4.edge/info_url.html', 41.881944, -87.627778, '/crossdomain.xml', '2018-01-19 21:19:32.247731', 1, 0, true, 0, null, null, null, null, null, null, false, 'target-ds4-displayname', null, 1, null, null, true, 0, null, false, null, null, 3, 'foo');
INSERT INTO deliveryservice (id, xml_id, active, dscp, signing_algorithm, qstring_ignore, geo_limit, http_bypass_fqdn, dns_bypass_ip, dns_bypass_ip6, dns_bypass_ttl, org_server_fqdn, type, profile, cdn_id, ccr_dns_ttl, global_max_mbps, global_max_tps, long_desc, long_desc_1, long_desc_2, max_dns_answers, info_url, miss_lat, miss_long, check_path, last_updated, protocol, ssl_key_version, ipv6_routing_enabled, range_request_handling, edge_header_rewrite, origin_shield, mid_header_rewrite, regex_remap, cacheurl, remap_text, multi_site_origin, display_name, tr_response_headers, initial_dispersion, dns_bypass_cname, tr_request_headers, regional_geo_blocking, geo_provider, geo_limit_countries, logs_enabled, multi_site_origin_algorithm, geolimit_redirect_url, tenant_id, routing_name) VALUES (200, 'test-ds2', true, 40, null, 0, 0, '', '', '', -1, 'http://test-ds2.edge', 9, 100, 100, 3600, 0, 0, 'test-ds2 long_desc', 'test-ds2 long_desc_1', 'test-ds2 long_desc_2', 0, 'http://test-ds2.edge/info_url.html', 41.881944, -87.627778, '/crossdomain.xml', '2018-01-19 21:19:32.253469', 0, 0, false, 0, null, null, null, null, null, null, false, 'test-ds2-displayname', null, 1, null, null, false, 0, null, false, null, null, 3, 'foo');
INSERT INTO deliveryservice (id, xml_id, active, dscp, signing_algorithm, qstring_ignore, geo_limit, http_bypass_fqdn, dns_bypass_ip, dns_bypass_ip6, dns_bypass_ttl, org_server_fqdn, type, profile, cdn_id, ccr_dns_ttl, global_max_mbps, global_max_tps, long_desc, long_desc_1, long_desc_2, max_dns_answers, info_url, miss_lat, miss_long, check_path, last_updated, protocol, ssl_key_version, ipv6_routing_enabled, range_request_handling, edge_header_rewrite, origin_shield, mid_header_rewrite, regex_remap, cacheurl, remap_text, multi_site_origin, display_name, tr_response_headers, initial_dispersion, dns_bypass_cname, tr_request_headers, regional_geo_blocking, geo_provider, geo_limit_countries, logs_enabled, multi_site_origin_algorithm, geolimit_redirect_url, tenant_id, routing_name) VALUES (2100, 'test-ds1-root', true, 40, null, 0, 0, '', '', '', -1, 'http://test-ds1-root.edge', 21, 100, 100, 3600, 0, 0, 'test-ds1-root long_desc', 'test-ds1-root long_desc_1', 'test-ds1-root long_desc_2', 0, 'http://test-ds1-root.edge/info_url.html', 41.881944, -87.627778, '/crossdomain.xml', '2018-01-19 21:19:32.261210', 1, 0, true, 0, null, null, null, null, null, null, false, 'test-ds1-root-displayname', null, 1, null, null, true, 0, null, true, null, null, 1, 'foo');
INSERT INTO deliveryservice (id, xml_id, active, dscp, signing_algorithm, qstring_ignore, geo_limit, http_bypass_fqdn, dns_bypass_ip, dns_bypass_ip6, dns_bypass_ttl, org_server_fqdn, type, profile, cdn_id, ccr_dns_ttl, global_max_mbps, global_max_tps, long_desc, long_desc_1, long_desc_2, max_dns_answers, info_url, miss_lat, miss_long, check_path, last_updated, protocol, ssl_key_version, ipv6_routing_enabled, range_request_handling, edge_header_rewrite, origin_shield, mid_header_rewrite, regex_remap, cacheurl, remap_text, multi_site_origin, display_name, tr_response_headers, initial_dispersion, dns_bypass_cname, tr_request_headers, regional_geo_blocking, geo_provider, geo_limit_countries, logs_enabled, multi_site_origin_algorithm, geolimit_redirect_url, tenant_id, routing_name) VALUES (2200, 'xxfoo.bar', true, 40, null, 0, 0, '', '', '', -1, 'http://foo.bar.edge', 21, 100, 100, 3600, 0, 0, 'foo.bar long_desc', 'foo.bar long_desc_1', 'foo.bar long_desc_2', 0, 'http://foo.bar.edge/info_url.html', 41.881944, -87.627778, '/crossdomain.xml', '2018-01-19 21:19:32.265717', 1, 0, true, 0, null, null, null, null, null, null, false, 'foo.bar-displayname', null, 1, null, null, true, 0, null, true, null, null, 2, 'foo');
INSERT INTO deliveryservice (id, xml_id, active, dscp, signing_algorithm, qstring_ignore, geo_limit, http_bypass_fqdn, dns_bypass_ip, dns_bypass_ip6, dns_bypass_ttl, org_server_fqdn, type, profile, cdn_id, ccr_dns_ttl, global_max_mbps, global_max_tps, long_desc, long_desc_1, long_desc_2, max_dns_answers, info_url, miss_lat, miss_long, check_path, last_updated, protocol, ssl_key_version, ipv6_routing_enabled, range_request_handling, edge_header_rewrite, origin_shield, mid_header_rewrite, regex_remap, cacheurl, remap_text, multi_site_origin, display_name, tr_response_headers, initial_dispersion, dns_bypass_cname, tr_request_headers, regional_geo_blocking, geo_provider, geo_limit_countries, logs_enabled, multi_site_origin_algorithm, geolimit_redirect_url, tenant_id, routing_name) VALUES (300, 'test-ds3', true, 40, null, 0, 0, '', '', '', -1, 'http://test-ds3.edge', 9, 100, 100, 3600, 0, 0, 'test-ds3 long_desc', 'test-ds3 long_desc_1', 'test-ds3 long_desc_2', 0, 'http://test-ds3.edge/info_url.html', 41.881944, -87.627778, '/crossdomain.xml', '2018-01-19 21:19:32.269358', 0, 0, false, 0, null, null, null, null, null, null, false, 'test-ds3-displayname', null, 1, null, null, false, 0, null, false, null, null, 3, 'foo');
INSERT INTO deliveryservice (id, xml_id, active, dscp, signing_algorithm, qstring_ignore, geo_limit, http_bypass_fqdn, dns_bypass_ip, dns_bypass_ip6, dns_bypass_ttl, org_server_fqdn, type, profile, cdn_id, ccr_dns_ttl, global_max_mbps, global_max_tps, long_desc, long_desc_1, long_desc_2, max_dns_answers, info_url, miss_lat, miss_long, check_path, last_updated, protocol, ssl_key_version, ipv6_routing_enabled, range_request_handling, edge_header_rewrite, origin_shield, mid_header_rewrite, regex_remap, cacheurl, remap_text, multi_site_origin, display_name, tr_response_headers, initial_dispersion, dns_bypass_cname, tr_request_headers, regional_geo_blocking, geo_provider, geo_limit_countries, logs_enabled, multi_site_origin_algorithm, geolimit_redirect_url, tenant_id, routing_name) VALUES (400, 'test-ds4', true, 40, null, 0, 0, '', '', '', -1, 'http://test-ds4.edge', 9, 100, 100, 3600, 0, 0, 'test-ds4 long_desc', 'test-ds4 long_desc_1', 'test-ds4 long_desc_2', 0, 'http://test-ds4.edge/info_url.html', 41.881944, -87.627778, '/crossdomain.xml', '2018-01-19 21:19:32.272467', 0, 0, false, 0, null, null, null, null, null, null, false, 'test-ds4-displayname', null, 1, null, null, false, 0, null, true, null, null, 4, 'foo');
INSERT INTO deliveryservice (id, xml_id, active, dscp, signing_algorithm, qstring_ignore, geo_limit, http_bypass_fqdn, dns_bypass_ip, dns_bypass_ip6, dns_bypass_ttl, org_server_fqdn, type, profile, cdn_id, ccr_dns_ttl, global_max_mbps, global_max_tps, long_desc, long_desc_1, long_desc_2, max_dns_answers, info_url, miss_lat, miss_long, check_path, last_updated, protocol, ssl_key_version, ipv6_routing_enabled, range_request_handling, edge_header_rewrite, origin_shield, mid_header_rewrite, regex_remap, cacheurl, remap_text, multi_site_origin, display_name, tr_response_headers, initial_dispersion, dns_bypass_cname, tr_request_headers, regional_geo_blocking, geo_provider, geo_limit_countries, logs_enabled, multi_site_origin_algorithm, geolimit_redirect_url, tenant_id, routing_name) VALUES (500, 'test-ds5', true, 40, null, 0, 0, '', '', '', -1, 'http://test-ds5.edge', 7, 300, 100, 3600, 0, 0, 'test-ds5 long_desc', 'test-ds5 long_desc_1', 'test-ds5 long_desc_2', 0, 'http://test-ds5.edge/info_url.html', 41.881944, -87.627778, '/crossdomain.xml', '2018-01-19 21:19:32.275400', 0, 0, false, 0, null, null, null, null, null, null, false, 'test-ds5-displayname', null, 1, null, null, false, 0, null, false, null, null, 3, 'foo');
INSERT INTO deliveryservice (id, xml_id, active, dscp, signing_algorithm, qstring_ignore, geo_limit, http_bypass_fqdn, dns_bypass_ip, dns_bypass_ip6, dns_bypass_ttl, org_server_fqdn, type, profile, cdn_id, ccr_dns_ttl, global_max_mbps, global_max_tps, long_desc, long_desc_1, long_desc_2, max_dns_answers, info_url, miss_lat, miss_long, check_path, last_updated, protocol, ssl_key_version, ipv6_routing_enabled, range_request_handling, edge_header_rewrite, origin_shield, mid_header_rewrite, regex_remap, cacheurl, remap_text, multi_site_origin, display_name, tr_response_headers, initial_dispersion, dns_bypass_cname, tr_request_headers, regional_geo_blocking, geo_provider, geo_limit_countries, logs_enabled, multi_site_origin_algorithm, geolimit_redirect_url, tenant_id, routing_name) VALUES (600, 'test-ds6', true, 40, null, 0, 0, '', '', '', -1, 'http://test-ds6.edge', 9, 300, 100, 3600, 0, 0, 'test-ds6 long_desc', 'test-ds6 long_desc_1', 'test-ds6 long_desc_2', 0, 'http://test-ds6.edge/info_url.html', 41.881944, -87.627778, '/crossdomain.xml', '2018-01-19 21:19:32.278451', 0, 0, false, 0, null, null, null, null, null, null, false, 'test-ds6-displayname', null, 1, null, null, false, 0, null, false, null, null, 3, 'foo');
INSERT INTO deliveryservice (id, xml_id, active, dscp, signing_algorithm, qstring_ignore, geo_limit, http_bypass_fqdn, dns_bypass_ip, dns_bypass_ip6, dns_bypass_ttl, org_server_fqdn, type, profile, cdn_id, ccr_dns_ttl, global_max_mbps, global_max_tps, long_desc, long_desc_1, long_desc_2, max_dns_answers, info_url, miss_lat, miss_long, check_path, last_updated, protocol, ssl_key_version, ipv6_routing_enabled, range_request_handling, edge_header_rewrite, origin_shield, mid_header_rewrite, regex_remap, cacheurl, remap_text, multi_site_origin, display_name, tr_response_headers, initial_dispersion, dns_bypass_cname, tr_request_headers, regional_geo_blocking, geo_provider, geo_limit_countries, logs_enabled, multi_site_origin_algorithm, geolimit_redirect_url, tenant_id, routing_name) VALUES (700, 'steering-ds1', true, 40, null, 0, 0, '', 'hokeypokey', null, 10, 'http://steering-ds1.edge', 21, 100, 100, 3600, 0, 0, 'steering-ds1 long_desc', 'steering-ds1 long_desc_1', 'steering-ds1 long_desc_2', 0, 'http://steering-ds1.edge/info_url.html', 41.881944, -87.627778, '/crossdomain.xml', '2018-01-19 21:19:32.281466', 1, 0, true, 0, null, null, null, null, null, null, false, 'steering-ds1-displayname', null, 1, null, null, true, 0, null, false, null, null, 3, 'foo');
INSERT INTO deliveryservice (id, xml_id, active, dscp, signing_algorithm, qstring_ignore, geo_limit, http_bypass_fqdn, dns_bypass_ip, dns_bypass_ip6, dns_bypass_ttl, org_server_fqdn, type, profile, cdn_id, ccr_dns_ttl, global_max_mbps, global_max_tps, long_desc, long_desc_1, long_desc_2, max_dns_answers, info_url, miss_lat, miss_long, check_path, last_updated, protocol, ssl_key_version, ipv6_routing_enabled, range_request_handling, edge_header_rewrite, origin_shield, mid_header_rewrite, regex_remap, cacheurl, remap_text, multi_site_origin, display_name, tr_response_headers, initial_dispersion, dns_bypass_cname, tr_request_headers, regional_geo_blocking, geo_provider, geo_limit_countries, logs_enabled, multi_site_origin_algorithm, geolimit_redirect_url, tenant_id, routing_name) VALUES (800, 'steering-ds2', true, 40, null, 0, 0, '', 'hokeypokey', null, 10, 'http://steering-ds2.edge', 21, 100, 100, 3600, 0, 0, 'steering-ds2 long_desc', 'steering-ds2 long_desc_1', 'steering-ds2 long_desc_2', 0, 'http://steering-ds2.edge/info_url.html', 41.881944, -87.627778, '/crossdomain.xml', '2018-01-19 21:19:32.284567', 1, 0, true, 0, null, null, null, null, null, null, false, 'steering-ds2-displayname', null, 1, null, null, true, 0, null, false, null, null, 3, 'foo');
INSERT INTO deliveryservice (id, xml_id, active, dscp, signing_algorithm, qstring_ignore, geo_limit, http_bypass_fqdn, dns_bypass_ip, dns_bypass_ip6, dns_bypass_ttl, org_server_fqdn, type, profile, cdn_id, ccr_dns_ttl, global_max_mbps, global_max_tps, long_desc, long_desc_1, long_desc_2, max_dns_answers, info_url, miss_lat, miss_long, check_path, last_updated, protocol, ssl_key_version, ipv6_routing_enabled, range_request_handling, edge_header_rewrite, origin_shield, mid_header_rewrite, regex_remap, cacheurl, remap_text, multi_site_origin, display_name, tr_response_headers, initial_dispersion, dns_bypass_cname, tr_request_headers, regional_geo_blocking, geo_provider, geo_limit_countries, logs_enabled, multi_site_origin_algorithm, geolimit_redirect_url, tenant_id, routing_name) VALUES (900, 'steering-ds3', true, 40, null, 0, 0, '', 'hokeypokey', null, 10, 'http://new-steering-ds.edge', 21, 100, 100, 3600, 0, 0, 'new-steering-ds long_desc', 'new-steering-ds long_desc_1', 'new-steering-ds long_desc_2', 0, 'http://new-steering-ds.edge/info_url.html', 41.881944, -87.627778, '/crossdomain.xml', '2018-01-19 21:19:32.287726', 1, 0, true, 0, null, null, null, null, null, null, false, 'new-steering-ds-displayname', null, 1, null, null, true, 0, null, false, null, null, 4, 'foo');
`
	err := execSQL(db, sqlStmt, "deliveryservice")
	if err != nil {
		return fmt.Errorf("exec failed %v", err)
	}
	return nil
}

// SetupRegexes ...
func SetupRegexes(db *sql.DB) error {

	sqlStmt := `
INSERT INTO regex (id, pattern, type, last_updated) VALUES (100, '.*\.omg-01\..*', 19, '2018-01-19 21:58:36.120746');
INSERT INTO regex (id, pattern, type, last_updated) VALUES (1000, '.*\.target-ds1\..*', 19, '2018-01-19 21:58:36.125624');
INSERT INTO regex (id, pattern, type, last_updated) VALUES (1100, '.*\.target-ds2\..*', 19, '2018-01-19 21:58:36.128372');
INSERT INTO regex (id, pattern, type, last_updated) VALUES (1200, '.*\.target-ds3\..*', 19, '2018-01-19 21:58:36.130749');
INSERT INTO regex (id, pattern, type, last_updated) VALUES (1300, '.*\.target-ds4\..*', 19, '2018-01-19 21:58:36.133992');
INSERT INTO regex (id, pattern, type, last_updated) VALUES (1400, '.*\.target-ds5\..*', 19, '2018-01-19 21:58:36.136503');
INSERT INTO regex (id, pattern, type, last_updated) VALUES (1500, '.*\.target-ds6\..*', 19, '2018-01-19 21:58:36.138890');
INSERT INTO regex (id, pattern, type, last_updated) VALUES (1600, '.*\.target-ds7\..*', 19, '2018-01-19 21:58:36.140495');
INSERT INTO regex (id, pattern, type, last_updated) VALUES (1700, '.*\.target-ds8\..*', 19, '2018-01-19 21:58:36.142473');
INSERT INTO regex (id, pattern, type, last_updated) VALUES (1800, '.*\.target-ds9\..*', 19, '2018-01-19 21:58:36.144087');
INSERT INTO regex (id, pattern, type, last_updated) VALUES (1900, '.*\.target-ds10\..*', 19, '2018-01-19 21:58:36.145505');
INSERT INTO regex (id, pattern, type, last_updated) VALUES (200, '.*\.foo\..*', 19, '2018-01-19 21:58:36.146953');
INSERT INTO regex (id, pattern, type, last_updated) VALUES (300, '.*/force-to-one/.*', 20, '2018-01-19 21:58:36.149052');
INSERT INTO regex (id, pattern, type, last_updated) VALUES (400, '.*/force-to-one-also/.*', 20, '2018-01-19 21:58:36.150904');
INSERT INTO regex (id, pattern, type, last_updated) VALUES (500, '.*/go-to-four/.*', 20, '2018-01-19 21:58:36.152416');
INSERT INTO regex (id, pattern, type, last_updated) VALUES (600, '.*/use-three/.*', 20, '2018-01-19 21:58:36.153884');
INSERT INTO regex (id, pattern, type, last_updated) VALUES (700, '.*\.new-steering-ds\..*', 19, '2018-01-19 21:58:36.155352');
INSERT INTO regex (id, pattern, type, last_updated) VALUES (800, '.*\.steering-ds1\..*', 19, '2018-01-19 21:58:36.156867');
INSERT INTO regex (id, pattern, type, last_updated) VALUES (900, '.*\.steering-ds2\..*', 19, '2018-01-19 21:58:36.158999');
`
	err := execSQL(db, sqlStmt, "regex")
	if err != nil {
		return fmt.Errorf("exec failed %v", err)
	}
	return nil
}

// SetupDeliveryServiceRegexes ...
func SetupDeliveryServiceRegexes(db *sql.DB) error {

	sqlStmt := `
INSERT INTO deliveryservice_regex (deliveryservice, regex, set_number) VALUES (200, 100, 0);
INSERT INTO deliveryservice_regex (deliveryservice, regex, set_number) VALUES (400, 100, 0);
INSERT INTO deliveryservice_regex (deliveryservice, regex, set_number) VALUES (400, 1000, 0);
INSERT INTO deliveryservice_regex (deliveryservice, regex, set_number) VALUES (500, 1100, 0);
INSERT INTO deliveryservice_regex (deliveryservice, regex, set_number) VALUES (600, 1200, 0);
INSERT INTO deliveryservice_regex (deliveryservice, regex, set_number) VALUES (700, 1300, 0);
INSERT INTO deliveryservice_regex (deliveryservice, regex, set_number) VALUES (800, 1400, 0);
INSERT INTO deliveryservice_regex (deliveryservice, regex, set_number) VALUES (900, 1500, 0);
INSERT INTO deliveryservice_regex (deliveryservice, regex, set_number) VALUES (1000, 1600, 0);
INSERT INTO deliveryservice_regex (deliveryservice, regex, set_number) VALUES (1100, 1700, 0);
INSERT INTO deliveryservice_regex (deliveryservice, regex, set_number) VALUES (1200, 1800, 0);
INSERT INTO deliveryservice_regex (deliveryservice, regex, set_number) VALUES (1300, 1900, 0);
INSERT INTO deliveryservice_regex (deliveryservice, regex, set_number) VALUES (100, 200, 0);
INSERT INTO deliveryservice_regex (deliveryservice, regex, set_number) VALUES (400, 200, 0);
INSERT INTO deliveryservice_regex (deliveryservice, regex, set_number) VALUES (700, 300, 0);
INSERT INTO deliveryservice_regex (deliveryservice, regex, set_number) VALUES (600, 400, 0);
INSERT INTO deliveryservice_regex (deliveryservice, regex, set_number) VALUES (300, 600, 0);
INSERT INTO deliveryservice_regex (deliveryservice, regex, set_number) VALUES (100, 800, 1);
INSERT INTO deliveryservice_regex (deliveryservice, regex, set_number) VALUES (200, 900, 0);
`
	err := execSQL(db, sqlStmt, "deliveryservice_regex")
	if err != nil {
		return fmt.Errorf("exec failed %v", err)
	}
	return nil
}

// SetupDeliveryServiceTmUsers ...
func SetupDeliveryServiceTmUsers(db *sql.DB) error {

	sqlStmt := `
INSERT INTO deliveryservice_tmuser (deliveryservice, tm_user_id, last_updated) VALUES (100, (SELECT id FROM tm_user where username = 'admin') , '2018-01-19 21:19:32.372969');
`
	err := execSQL(db, sqlStmt, "deliveryservice_tmuser")
	if err != nil {
		return fmt.Errorf("exec failed %v", err)
	}
	return nil
}

// SetupDeliveryServiceServers ...
func SetupDeliveryServiceServers(db *sql.DB) error {

	sqlStmt := `
INSERT INTO deliveryservice_server (deliveryservice, server, last_updated) VALUES (100, 300, '2018-01-19 21:19:32.396609');
INSERT INTO deliveryservice_server (deliveryservice, server, last_updated) VALUES (100, 1300, '2018-01-19 21:19:32.408819');
INSERT INTO deliveryservice_server (deliveryservice, server, last_updated) VALUES (100, 100, '2018-01-19 21:19:32.414612');
INSERT INTO deliveryservice_server (deliveryservice, server, last_updated) VALUES (200, 800, '2018-01-19 21:19:32.420745');
INSERT INTO deliveryservice_server (deliveryservice, server, last_updated) VALUES (200, 700, '2018-01-19 21:19:32.426505');
INSERT INTO deliveryservice_server (deliveryservice, server, last_updated) VALUES (500, 1500, '2018-01-19 21:19:32.434097');
INSERT INTO deliveryservice_server (deliveryservice, server, last_updated) VALUES (500, 1400, '2018-01-19 21:19:32.439622');
INSERT INTO deliveryservice_server (deliveryservice, server, last_updated) VALUES (600, 1400, '2018-01-19 21:19:32.440831');
INSERT INTO deliveryservice_server (deliveryservice, server, last_updated) VALUES (600, 1500, '2018-01-19 21:19:32.442121');
INSERT INTO deliveryservice_server (deliveryservice, server, last_updated) VALUES (700, 900, '2018-01-19 21:19:32.443372');
`
	err := execSQL(db, sqlStmt, "deliveryservice_server")
	if err != nil {
		return fmt.Errorf("exec failed %v", err)
	}
	return nil
}

// SetupJobStatuses ...
func SetupJobStatuses(db *sql.DB) error {

	sqlStmt := `
INSERT INTO job_status (id, name, description, last_updated) VALUES (1, 'PENDING', 'Job is queued, but has not been picked up by any agents yet', '2018-01-19 21:19:32.444857');
`
	err := execSQL(db, sqlStmt, "job_status")
	if err != nil {
		return fmt.Errorf("exec failed %v", err)
	}
	return nil
}

// SetupJobAgents ...
func SetupJobAgents(db *sql.DB) error {

	sqlStmt := `
INSERT INTO job_agent (id, name, description, active, last_updated) VALUES (1, 'agent1', 'Test Agent1', 0, '2018-01-19 21:19:32.448076');
`
	err := execSQL(db, sqlStmt, "job_agent")
	if err != nil {
		return fmt.Errorf("exec failed %v", err)
	}
	return nil
}

// SetupJobs ...
func SetupJobs(db *sql.DB) error {

	sqlStmt := `
INSERT INTO job (id, agent, object_type, object_name, keyword, parameters, asset_url, asset_type, status, start_time, entered_time, job_user, last_updated, job_deliveryservice) VALUES (100, 1, null, null, 'PURGE', 'TTL:48h', 'http://cdn2.edge/job1/.*', 'file', 1, '2018-01-19 21:01:14.000000', '2018-01-19 21:01:14.000000', (SELECT id FROM tm_user where username = 'admin'), '2018-01-19 21:19:32.468643', 100);
INSERT INTO job (id, agent, object_type, object_name, keyword, parameters, asset_url, asset_type, status, start_time, entered_time, job_user, last_updated, job_deliveryservice) VALUES (200, 1, null, null, 'PURGE', 'TTL:48h', 'http://cdn2.edge/job2/.*', 'file', 1, '2018-01-19 21:09:34.000000', '2018-01-19 21:09:34.000000', (SELECT id FROM tm_user where username = 'admin'), '2018-01-19 21:19:32.450915', 200);
INSERT INTO job (id, agent, object_type, object_name, keyword, parameters, asset_url, asset_type, status, start_time, entered_time, job_user, last_updated, job_deliveryservice) VALUES (300, 1, null, null, 'PURGE', 'TTL:48h', 'http://cdn2.edge/job3/.*', 'file', 1, '2018-01-19 21:14:34.000000', '2018-01-19 21:14:34.000000', (SELECT id FROM tm_user where username = 'admin'), '2018-01-19 21:19:32.460870', 100);
`
	err := execSQL(db, sqlStmt, "job")
	if err != nil {
		return fmt.Errorf("exec failed %v", err)
	}
	return nil
}

// Teardown - ensures that the data is cleaned up for a fresh run
func Teardown(db *sql.DB) error {

	sqlStmt := `
	DELETE FROM to_extension;
	DELETE FROM staticdnsentry;
	DELETE FROM job;
	DELETE FROM job_agent;
	DELETE FROM job_status;
	DELETE FROM log;
	DELETE FROM asn;
	DELETE FROM deliveryservice_tmuser;
	DELETE FROM tm_user;
	DELETE FROM role;
	DELETE FROM deliveryservice_regex;
	DELETE FROM regex;
	DELETE FROM deliveryservice_server;
	DELETE FROM deliveryservice;
	DELETE FROM server;
	DELETE FROM phys_location;
	DELETE FROM region;
	DELETE FROM division;
	DELETE FROM profile;
	DELETE FROM parameter;
	DELETE FROM profile_parameter;
	DELETE FROM cachegroup;
	DELETE FROM type;
	DELETE FROM status;
	DELETE FROM snapshot;
	DELETE FROM cdn;
	DELETE FROM tenant;
`
	err := execSQL(db, sqlStmt, "Tearing down")
	if err != nil {
		return fmt.Errorf("exec failed %v", err)
	}
	return err
}

// execSQL ...
func execSQL(db *sql.DB, sqlStmt string, dbTable string) error {

	log.Debugln(dbTable + " data")
	var err error

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("transaction begin failed %v %v ", err, tx)
	}

	res, err := tx.Exec(sqlStmt)
	if err != nil {
		return fmt.Errorf("exec failed %v %v", err, res)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commit failed %v %v", err, res)
	}
	return nil
}
