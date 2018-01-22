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

package todb

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/testing/api/config"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/auth"
)

var (
	db *sql.DB
)

// OpenConnection ...
func OpenConnection(cfg *config.Config) (*sql.DB, error) {
	var err error
	sslStr := "require"
	if !cfg.TrafficOpsDB.SSL {
		sslStr = "disable"
	}

	db, err = sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s", cfg.TrafficOpsDB.User, cfg.TrafficOpsDB.Password, cfg.TrafficOpsDB.Hostname, cfg.TrafficOpsDB.Name, sslStr))

	if err != nil {
		log.Errorf("opening database: %v\n", err)
		return nil, fmt.Errorf("transaction failed: %s", err)
	}
	return db, err
}

// SetupTestData ...
func SetupTestData(cfg *config.Config, db *sql.DB) error {
	var err error

	err = SetupTenants(cfg, db)
	if err != nil {
		fmt.Printf("\nError setting up tenants %s - %s, %v\n", cfg.TrafficOps.URL, cfg.TrafficOps.User, err)
		os.Exit(1)
	}

	err = SetupCDNs(cfg, db)
	if err != nil {
		fmt.Printf("\nError setting up cdns %s - %s, %v\n", cfg.TrafficOps.URL, cfg.TrafficOps.User, err)
		os.Exit(1)
	}

	err = SetupRoles(cfg, db)
	if err != nil {
		fmt.Printf("\nError setting up roles %s - %s, %v\n", cfg.TrafficOps.URL, cfg.TrafficOps.User, err)
		os.Exit(1)
	}

	err = SetupTmusers(cfg, db)
	if err != nil {
		fmt.Printf("\nError setting up tm_user %s - %s, %v\n", cfg.TrafficOps.URL, cfg.TrafficOps.User, err)
		os.Exit(1)
	}

	err = SetupStatuses(cfg, db)
	if err != nil {
		fmt.Printf("\nError setting up status %s - %s, %v\n", cfg.TrafficOps.URL, cfg.TrafficOps.User, err)
		os.Exit(1)
	}

	err = SetupParameters(cfg, db)
	if err != nil {
		fmt.Printf("\nError setting up parameter %s - %s, %v\n", cfg.TrafficOps.URL, cfg.TrafficOps.User, err)
		os.Exit(1)
	}

	err = SetupProfiles(cfg, db)
	if err != nil {
		fmt.Printf("\nError setting up profile %s - %s, %v\n", cfg.TrafficOps.URL, cfg.TrafficOps.User, err)
		os.Exit(1)
	}

	err = SetupProfileParameters(cfg, db)
	if err != nil {
		fmt.Printf("\nError setting up parameter %s - %s, %v\n", cfg.TrafficOps.URL, cfg.TrafficOps.User, err)
		os.Exit(1)
	}
	err = SetupTypes(cfg, db)
	if err != nil {
		fmt.Printf("\nError setting up type %s - %s, %v\n", cfg.TrafficOps.URL, cfg.TrafficOps.User, err)
		os.Exit(1)
	}

	err = SetupCacheGroups(cfg, db)
	if err != nil {
		fmt.Printf("\nError setting up cachegroup %s - %s, %v\n", cfg.TrafficOps.URL, cfg.TrafficOps.User, err)
		os.Exit(1)
	}

	err = SetupDivisions(cfg, db)
	if err != nil {
		fmt.Printf("\nError setting up division %s - %s, %v\n", cfg.TrafficOps.URL, cfg.TrafficOps.User, err)
		os.Exit(1)
	}

	err = SetupRegions(cfg, db)
	if err != nil {
		fmt.Printf("\nError setting up region %s - %s, %v\n", cfg.TrafficOps.URL, cfg.TrafficOps.User, err)
		os.Exit(1)
	}

	err = SetupPhysLocations(cfg, db)
	if err != nil {
		fmt.Printf("\nError setting up phys_location %s - %s, %v\n", cfg.TrafficOps.URL, cfg.TrafficOps.User, err)
		os.Exit(1)
	}

	err = SetupServers(cfg, db)
	if err != nil {
		fmt.Printf("\nError setting up server %s - %s, %v\n", cfg.TrafficOps.URL, cfg.TrafficOps.User, err)
		os.Exit(1)
	}

	err = SetupAsns(cfg, db)
	if err != nil {
		fmt.Printf("\nError setting up asn %s - %s, %v\n", cfg.TrafficOps.URL, cfg.TrafficOps.User, err)
		os.Exit(1)
	}

	err = SetupDeliveryServices(cfg, db)
	if err != nil {
		fmt.Printf("\nError setting up deliveryservice %s - %s, %v\n", cfg.TrafficOps.URL, cfg.TrafficOps.User, err)
		os.Exit(1)
	}

	err = SetupRegexes(cfg, db)
	if err != nil {
		fmt.Printf("\nError setting up regex %s - %s, %v\n", cfg.TrafficOps.URL, cfg.TrafficOps.User, err)
		os.Exit(1)
	}

	err = SetupDeliveryServiceRegexes(cfg, db)
	if err != nil {
		fmt.Printf("\nError setting up deliveryservice_regex %s - %s, %v\n", cfg.TrafficOps.URL, cfg.TrafficOps.User, err)
		os.Exit(1)
	}

	err = SetupDeliveryServiceTmUsers(cfg, db)
	if err != nil {
		fmt.Printf("\nError setting up deliveryservice_tmuser %s - %s, %v\n", cfg.TrafficOps.URL, cfg.TrafficOps.User, err)
		os.Exit(1)
	}

	err = SetupDeliveryServiceServers(cfg, db)
	if err != nil {
		fmt.Printf("\nError setting up deliveryservice_server %s - %s, %v\n", cfg.TrafficOps.URL, cfg.TrafficOps.User, err)
		os.Exit(1)
	}

	err = SetupJobStatuses(cfg, db)
	if err != nil {
		fmt.Printf("\nError setting up job_status %s - %s, %v\n", cfg.TrafficOps.URL, cfg.TrafficOps.User, err)
		os.Exit(1)
	}

	err = SetupJobAgents(cfg, db)
	if err != nil {
		fmt.Printf("\nError setting up job_agent %s - %s, %v\n", cfg.TrafficOps.URL, cfg.TrafficOps.User, err)
		os.Exit(1)
	}

	err = SetupJobs(cfg, db)
	if err != nil {
		fmt.Printf("\nError setting up job %s - %s, %v\n", cfg.TrafficOps.URL, cfg.TrafficOps.User, err)
		os.Exit(1)
	}

	return err
}

// SetupRoles ...
func SetupRoles(cfg *config.Config, db *sql.DB) error {

	sqlStmt := `
INSERT INTO role (id, name, description, priv_level) VALUES (1, 'disallowed','Block all access',0) ON CONFLICT DO NOTHING;
INSERT INTO role (id, name, description, priv_level) VALUES (2, 'read-only user','Block all access', 10) ON CONFLICT DO NOTHING;
INSERT INTO role (id, name, description, priv_level) VALUES (3, 'operations','Block all access', 20) ON CONFLICT DO NOTHING;
INSERT INTO role (id, name, description, priv_level) VALUES (4, 'admin','super-user', 30) ON CONFLICT DO NOTHING;
INSERT INTO role (id, name, description, priv_level) VALUES (5, 'portal','Portal User', 2) ON CONFLICT DO NOTHING;
INSERT INTO role (id, name, description, priv_level) VALUES (7, 'federation','Role for Secondary CZF', 15) ON CONFLICT DO NOTHING;
`
	err := execSQL(cfg, db, sqlStmt, "role")
	if err != nil {
		return fmt.Errorf("exec failed %v", err)
	}
	return nil
}

// SetupTmusers ...
func SetupTmusers(cfg *config.Config, db *sql.DB) error {

	var err error
	encryptedPassword, err := auth.DerivePassword(cfg.TrafficOps.UserPassword)
	if err != nil {
		return fmt.Errorf("password encryption failed %v", err)
	}
	sqlStmt := `INSERT INTO tm_user (username, local_passwd, confirm_local_passwd, role) VALUES ('admin','` + encryptedPassword + `','` + encryptedPassword + `', 4)`
	err = execSQL(cfg, db, sqlStmt, "tm_user")
	if err != nil {
		return fmt.Errorf("exec failed %v", err)
	}
	return nil
}

// SetupTenants ...
func SetupTenants(cfg *config.Config, db *sql.DB) error {

	sqlStmt := `
INSERT INTO tenant (id, name, active, parent_id, last_updated) VALUES (1000000000, 'root', true, null, '2018-01-19 19:01:21.327262');
`
	err := execSQL(cfg, db, sqlStmt, "tenant")
	if err != nil {
		return fmt.Errorf("exec failed %v", err)
	}
	return nil
}

// SetupCDNs ...
func SetupCDNs(cfg *config.Config, db *sql.DB) error {

	sqlStmt := `
INSERT INTO cdn (id, name, last_updated, dnssec_enabled, domain_name) VALUES (100, 'cdn5', '2018-01-19 21:19:31.588795', false, 'cdn1.kabletown.net');
INSERT INTO cdn (id, name, last_updated, dnssec_enabled, domain_name) VALUES (200, 'cdn6', '2018-01-19 21:19:31.591457', false, 'cdn2.kabletown.net');
INSERT INTO cdn (id, name, last_updated, dnssec_enabled, domain_name) VALUES (300, 'cdn7', '2018-01-19 21:19:31.592700', false, 'cdn3.kabletown.net');
`
	err := execSQL(cfg, db, sqlStmt, "cdn")
	if err != nil {
		return fmt.Errorf("exec failed %v", err)
	}
	return nil
}

// SetupStatuses ...
func SetupStatuses(cfg *config.Config, db *sql.DB) error {

	sqlStmt := `
INSERT INTO status (id, name, description, last_updated) VALUES (1, 'OFFLINE', 'Edge: Puts server in CCR config file in this state, but CCR will never route traffic to it. Mid: Server will not be included in parent.config files for its edge caches', '2018-01-19 19:01:21.388399');
INSERT INTO status (id, name, description, last_updated) VALUES (2, 'ONLINE', 'Edge: Puts server in CCR config file in this state, and CCR will always route traffic to it. Mid: Server will be included in parent.config files for its edges', '2018-01-19 19:01:21.384459');
INSERT INTO status (id, name, description, last_updated) VALUES (3, 'REPORTED', 'Edge: Puts server in CCR config file in this state, and CCR will adhere to the health protocol. Mid: N/A for now', '2018-01-19 19:01:21.379811');
INSERT INTO status (id, name, description, last_updated) VALUES (4, 'ADMIN_DOWN', 'Temporary down. Edge: XMPP client will send status OFFLINE to CCR, otherwise similar to REPORTED. Mid: Server will not be included in parent.config files for its edge caches', '2018-01-19 19:01:21.385798');
INSERT INTO status (id, name, description, last_updated) VALUES (5, 'CCR_IGNORE', 'Edge: 12M will not include caches in this state in CCR config files. Mid: N/A for now', '2018-01-19 19:01:21.383085');
INSERT INTO status (id, name, description, last_updated) VALUES (6, 'PRE_PROD', 'Pre Production. Not active in any configuration.', '2018-01-19 19:01:21.387146');
`
	err := execSQL(cfg, db, sqlStmt, "status")
	if err != nil {
		return fmt.Errorf("exec failed %v", err)
	}
	return nil
}

// SetupCacheGroups ...
func SetupCacheGroups(cfg *config.Config, db *sql.DB) error {

	sqlStmt := `
INSERT INTO cachegroup (id, name, short_name, latitude, longitude, parent_cachegroup_id, secondary_parent_cachegroup_id, type, last_updated) VALUES (100, 'mid-northeast-group', 'ne', 120, 120, null, null, 2, '2018-01-19 21:19:32.041913');
INSERT INTO cachegroup (id, name, short_name, latitude, longitude, parent_cachegroup_id, secondary_parent_cachegroup_id, type, last_updated) VALUES (200, 'mid-northwest-group', 'nw', 100, 100, 100, null, 2, '2018-01-19 21:19:32.052005');
INSERT INTO cachegroup (id, name, short_name, latitude, longitude, parent_cachegroup_id, secondary_parent_cachegroup_id, type, last_updated) VALUES (800, 'mid_cg3', 'mid_cg3', 100, 100, null, null, 6, '2018-01-19 21:19:32.056908');
INSERT INTO cachegroup (id, name, short_name, latitude, longitude, parent_cachegroup_id, secondary_parent_cachegroup_id, type, last_updated) VALUES (900, 'edge_cg4', 'edge_cg4', 100, 100, 800, null, 5, '2018-01-19 21:19:32.059077');
INSERT INTO cachegroup (id, name, short_name, latitude, longitude, parent_cachegroup_id, secondary_parent_cachegroup_id, type, last_updated) VALUES (300, 'edge_atl_group', 'atl', 120, 120, 100, 200, 5, '2018-01-19 21:19:32.063375');
`
	err := execSQL(cfg, db, sqlStmt, "cachegroup")
	if err != nil {
		return fmt.Errorf("exec failed %v", err)
	}
	return nil
}

// SetupPhysLocations ...
func SetupPhysLocations(cfg *config.Config, db *sql.DB) error {

	sqlStmt := `
INSERT INTO phys_location (id, name, short_name, address, city, state, zip, poc, phone, email, comments, region, last_updated) VALUES (100, 'Denver', 'denver', '1234 mile high circle', 'Denver', 'CO', '80202', null, '303-111-1111', null, null, 100, '2018-01-19 21:19:32.081465');
INSERT INTO phys_location (id, name, short_name, address, city, state, zip, poc, phone, email, comments, region, last_updated) VALUES (200, 'Boulder', 'boulder', '1234 green way', 'Boulder', 'CO', '80301', null, '303-222-2222', null, null, 100, '2018-01-19 21:19:32.086195');
INSERT INTO phys_location (id, name, short_name, address, city, state, zip, poc, phone, email, comments, region, last_updated) VALUES (300, 'HotAtlanta', 'atlanta', '1234 southern way', 'Atlanta', 'GA', '30301', null, '404-222-2222', null, null, 100, '2018-01-19 21:19:32.089538');
`
	err := execSQL(cfg, db, sqlStmt, "phys_location")
	if err != nil {
		return fmt.Errorf("exec failed %v", err)
	}
	return nil
}

// SetupServers ...
func SetupServers(cfg *config.Config, db *sql.DB) error {

	sqlStmt := `
INSERT INTO server (id, host_name, domain_name, tcp_port, xmpp_id, xmpp_passwd, interface_name, ip_address, ip_netmask, ip_gateway, ip6_address, ip6_gateway, interface_mtu, phys_location, rack, cachegroup, type, status, offline_reason, upd_pending, profile, cdn_id, mgmt_ip_address, mgmt_ip_netmask, mgmt_ip_gateway, ilo_ip_address, ilo_ip_netmask, ilo_ip_gateway, ilo_username, ilo_password, router_host_name, router_port_name, guid, last_updated, https_port, reval_pending) VALUES (100, 'atlanta-edge-01', 'ga.atlanta.kabletown.net', 80, 'atlanta-edge-01\@ocdn.kabletown.net', 'X', 'bond0', '127.0.0.1', '255.255.255.252', '127.0.0.1', '2345:1234:12:8::2/64', '2345:1234:12:8::1', 9000, 100, 'RR 119.02', 300, 1, 2, null, false, 100, 100, '', '', '', '', '', '', '', '', '', '', null, '2018-01-19 21:19:32.094746', null, false);
INSERT INTO server (id, host_name, domain_name, tcp_port, xmpp_id, xmpp_passwd, interface_name, ip_address, ip_netmask, ip_gateway, ip6_address, ip6_gateway, interface_mtu, phys_location, rack, cachegroup, type, status, offline_reason, upd_pending, profile, cdn_id, mgmt_ip_address, mgmt_ip_netmask, mgmt_ip_gateway, ilo_ip_address, ilo_ip_netmask, ilo_ip_gateway, ilo_username, ilo_password, router_host_name, router_port_name, guid, last_updated, https_port, reval_pending) VALUES (1000, 'influxdb02', 'kabletown.net', 8086, '', '', 'eth1', '127.0.0.11', '255.255.252.0', '127.0.0.11', '127.0.0.11', '127.0.0.11', 1500, 300, 'RR 119.02', 100, 32, 2, null, false, 500, 100, '', '', '', '', '', '', '', '', '', '', null, '2018-01-19 21:19:32.115164', null, false);
INSERT INTO server (id, host_name, domain_name, tcp_port, xmpp_id, xmpp_passwd, interface_name, ip_address, ip_netmask, ip_gateway, ip6_address, ip6_gateway, interface_mtu, phys_location, rack, cachegroup, type, status, offline_reason, upd_pending, profile, cdn_id, mgmt_ip_address, mgmt_ip_netmask, mgmt_ip_gateway, ilo_ip_address, ilo_ip_netmask, ilo_ip_gateway, ilo_username, ilo_password, router_host_name, router_port_name, guid, last_updated, https_port, reval_pending) VALUES (1100, 'atlanta-router-01', 'ga.atlanta.kabletown.net', 80, 'atlanta-router-01\@ocdn.kabletown.net', 'X', 'bond0', '127.0.0.12', '255.255.255.252', '127.0.0.1', '2345:1234:12:8::10/64', '2345:1234:12:8::1', 9000, 100, 'RR 119.02', 300, 4, 2, null, false, 100, 100, '', '', '', '', '', '', '', '', '', '', null, '2018-01-19 21:19:32.125603', null, false);
INSERT INTO server (id, host_name, domain_name, tcp_port, xmpp_id, xmpp_passwd, interface_name, ip_address, ip_netmask, ip_gateway, ip6_address, ip6_gateway, interface_mtu, phys_location, rack, cachegroup, type, status, offline_reason, upd_pending, profile, cdn_id, mgmt_ip_address, mgmt_ip_netmask, mgmt_ip_gateway, ilo_ip_address, ilo_ip_netmask, ilo_ip_gateway, ilo_username, ilo_password, router_host_name, router_port_name, guid, last_updated, https_port, reval_pending) VALUES (1200, 'atlanta-edge-03', 'ga.atlanta.kabletown.net', 80, 'atlanta-edge-03\@ocdn.kabletown.net', 'X', 'bond0', '127.0.0.13', '255.255.255.252', '127.0.0.1', '2345:1234:12:2::2/64', '2345:1234:12:8::1', 9000, 100, 'RR 119.02', 300, 1, 3, null, false, 100, 100, '', '', '', '', '', '', '', '', '', '', null, '2018-01-19 21:19:32.135422', null, false);
INSERT INTO server (id, host_name, domain_name, tcp_port, xmpp_id, xmpp_passwd, interface_name, ip_address, ip_netmask, ip_gateway, ip6_address, ip6_gateway, interface_mtu, phys_location, rack, cachegroup, type, status, offline_reason, upd_pending, profile, cdn_id, mgmt_ip_address, mgmt_ip_netmask, mgmt_ip_gateway, ilo_ip_address, ilo_ip_netmask, ilo_ip_gateway, ilo_username, ilo_password, router_host_name, router_port_name, guid, last_updated, https_port, reval_pending) VALUES (1300, 'atlanta-edge-14', 'ga.atlanta.kabletown.net', 80, 'atlanta-edge-14\@ocdn.kabletown.net', 'X', 'bond0', '127.0.0.14', '255.255.255.252', '127.0.0.1', '2345:1234:12:8::14/64', '2345:1234:12:8::1', 9000, 100, 'RR 119.02', 900, 1, 2, null, false, 100, 100, '', '', '', '', '', '', '', '', '', '', null, '2018-01-19 21:19:32.145252', null, false);
INSERT INTO server (id, host_name, domain_name, tcp_port, xmpp_id, xmpp_passwd, interface_name, ip_address, ip_netmask, ip_gateway, ip6_address, ip6_gateway, interface_mtu, phys_location, rack, cachegroup, type, status, offline_reason, upd_pending, profile, cdn_id, mgmt_ip_address, mgmt_ip_netmask, mgmt_ip_gateway, ilo_ip_address, ilo_ip_netmask, ilo_ip_gateway, ilo_username, ilo_password, router_host_name, router_port_name, guid, last_updated, https_port, reval_pending) VALUES (1400, 'atlanta-edge-15', 'ga.atlanta.kabletown.net', 80, 'atlanta-edge-15\@ocdn.kabletown.net', 'X', 'bond0', '127.0.0.15', '255.255.255.252', '127.0.0.7', '2345:1234:12:d::15/64', '2345:1234:12:d::1', 9000, 100, 'RR 119.02', 900, 1, 2, null, false, 100, 100, '', '', '', '', '', '', '', '', '', '', null, '2018-01-19 21:19:32.155043', null, false);
INSERT INTO server (id, host_name, domain_name, tcp_port, xmpp_id, xmpp_passwd, interface_name, ip_address, ip_netmask, ip_gateway, ip6_address, ip6_gateway, interface_mtu, phys_location, rack, cachegroup, type, status, offline_reason, upd_pending, profile, cdn_id, mgmt_ip_address, mgmt_ip_netmask, mgmt_ip_gateway, ilo_ip_address, ilo_ip_netmask, ilo_ip_gateway, ilo_username, ilo_password, router_host_name, router_port_name, guid, last_updated, https_port, reval_pending) VALUES (1500, 'atlanta-mid-16', 'ga.atlanta.kabletown.net', 80, 'atlanta-mid-16\@ocdn.kabletown.net', 'X', 'bond0', '127.0.0.16', '255.255.255.252', '127.0.0.7', '2345:1234:12:d::16/64', '2345:1234:12:d::1', 9000, 100, 'RR 119.02', 800, 2, 2, null, false, 100, 100, '', '', '', '', '', '', '', '', '', '', null, '2018-01-19 21:19:32.164825', null, false);
INSERT INTO server (id, host_name, domain_name, tcp_port, xmpp_id, xmpp_passwd, interface_name, ip_address, ip_netmask, ip_gateway, ip6_address, ip6_gateway, interface_mtu, phys_location, rack, cachegroup, type, status, offline_reason, upd_pending, profile, cdn_id, mgmt_ip_address, mgmt_ip_netmask, mgmt_ip_gateway, ilo_ip_address, ilo_ip_netmask, ilo_ip_gateway, ilo_username, ilo_password, router_host_name, router_port_name, guid, last_updated, https_port, reval_pending) VALUES (1600, 'atlanta-org-1', 'ga.atlanta.kabletown.net', 80, 'atlanta-org-1\@ocdn.kabletown.net', 'X', 'bond0', '127.0.0.17', '255.255.255.252', '127.0.0.17', '2345:1234:12:d::17/64', '2345:1234:12:d::1', 9000, 100, 'RR 119.02', 800, 3, 2, null, false, 100, 100, '', '', '', '', '', '', '', '', '', '', null, '2018-01-19 21:19:32.167782', null, false);
INSERT INTO server (id, host_name, domain_name, tcp_port, xmpp_id, xmpp_passwd, interface_name, ip_address, ip_netmask, ip_gateway, ip6_address, ip6_gateway, interface_mtu, phys_location, rack, cachegroup, type, status, offline_reason, upd_pending, profile, cdn_id, mgmt_ip_address, mgmt_ip_netmask, mgmt_ip_gateway, ilo_ip_address, ilo_ip_netmask, ilo_ip_gateway, ilo_username, ilo_password, router_host_name, router_port_name, guid, last_updated, https_port, reval_pending) VALUES (1700, 'atlanta-org-2', 'ga.atlanta.kabletown.net', 80, 'atlanta-org-1\@ocdn.kabletown.net', 'X', 'bond0', '127.0.0.18', '255.255.255.252', '127.0.0.18', '2345:1234:12:d::18/64', '2345:1234:12:d::1', 9000, 100, 'RR 119.02', 800, 3, 2, null, false, 900, 200, '', '', '', '', '', '', '', '', '', '', null, '2018-01-19 21:19:32.170592', null, false);
INSERT INTO server (id, host_name, domain_name, tcp_port, xmpp_id, xmpp_passwd, interface_name, ip_address, ip_netmask, ip_gateway, ip6_address, ip6_gateway, interface_mtu, phys_location, rack, cachegroup, type, status, offline_reason, upd_pending, profile, cdn_id, mgmt_ip_address, mgmt_ip_netmask, mgmt_ip_gateway, ilo_ip_address, ilo_ip_netmask, ilo_ip_gateway, ilo_username, ilo_password, router_host_name, router_port_name, guid, last_updated, https_port, reval_pending) VALUES (200, 'atlanta-mid-01', 'ga.atlanta.kabletown.net', 80, 'atlanta-mid-01\@ocdn.kabletown.net', 'X', 'bond0', '127.0.0.2', '255.255.255.252', '127.0.0.2', '2345:1234:12:9::2/64', '2345:1234:12:9::1', 9000, 100, 'RR 119.02', 100, 2, 2, null, false, 200, 100, '', '', '', '', '', '', '', '', '', '', null, '2018-01-19 21:19:32.173304', null, false);
INSERT INTO server (id, host_name, domain_name, tcp_port, xmpp_id, xmpp_passwd, interface_name, ip_address, ip_netmask, ip_gateway, ip6_address, ip6_gateway, interface_mtu, phys_location, rack, cachegroup, type, status, offline_reason, upd_pending, profile, cdn_id, mgmt_ip_address, mgmt_ip_netmask, mgmt_ip_gateway, ilo_ip_address, ilo_ip_netmask, ilo_ip_gateway, ilo_username, ilo_password, router_host_name, router_port_name, guid, last_updated, https_port, reval_pending) VALUES (300, 'rascal01', 'kabletown.net', 81, 'rascal\@kabletown.net', 'X', 'bond0', '127.0.0.4', '255.255.255.252', '127.0.0.4', '2345:1234:12:b::2/64', '2345:1234:12:b::1', 9000, 100, 'RR 119.02', 100, 4, 2, null, false, 300, 200, '', '', '', '', '', '', '', '', '', '', null, '2018-01-19 21:19:32.176375', null, false);
INSERT INTO server (id, host_name, domain_name, tcp_port, xmpp_id, xmpp_passwd, interface_name, ip_address, ip_netmask, ip_gateway, ip6_address, ip6_gateway, interface_mtu, phys_location, rack, cachegroup, type, status, offline_reason, upd_pending, profile, cdn_id, mgmt_ip_address, mgmt_ip_netmask, mgmt_ip_gateway, ilo_ip_address, ilo_ip_netmask, ilo_ip_gateway, ilo_username, ilo_password, router_host_name, router_port_name, guid, last_updated, https_port, reval_pending) VALUES (400, 'riak01', 'kabletown.net', 8088, '', '', 'eth1', '127.0.0.5', '255.255.252.0', '127.0.0.5', '', '', 1500, 100, 'RR 119.02', 100, 31, 2, null, false, 500, 100, '', '', '', '', '', '', '', '', '', '', null, '2018-01-19 21:19:32.180698', null, false);
INSERT INTO server (id, host_name, domain_name, tcp_port, xmpp_id, xmpp_passwd, interface_name, ip_address, ip_netmask, ip_gateway, ip6_address, ip6_gateway, interface_mtu, phys_location, rack, cachegroup, type, status, offline_reason, upd_pending, profile, cdn_id, mgmt_ip_address, mgmt_ip_netmask, mgmt_ip_gateway, ilo_ip_address, ilo_ip_netmask, ilo_ip_gateway, ilo_username, ilo_password, router_host_name, router_port_name, guid, last_updated, https_port, reval_pending) VALUES (500, 'rascal02', 'kabletown.net', 81, 'rascal\@kabletown.net', 'X', 'bond0', '127.0.0.6', '255.255.255.252', '127.0.0.6', '2345:1234:12:c::2/64', '2345:1234:12:c::1', 9000, 100, 'RR 119.05', 100, 4, 2, null, false, 300, 200, '', '', '', '', '', '', '', '', '', '', null, '2018-01-19 21:19:32.184322', null, false);
INSERT INTO server (id, host_name, domain_name, tcp_port, xmpp_id, xmpp_passwd, interface_name, ip_address, ip_netmask, ip_gateway, ip6_address, ip6_gateway, interface_mtu, phys_location, rack, cachegroup, type, status, offline_reason, upd_pending, profile, cdn_id, mgmt_ip_address, mgmt_ip_netmask, mgmt_ip_gateway, ilo_ip_address, ilo_ip_netmask, ilo_ip_gateway, ilo_username, ilo_password, router_host_name, router_port_name, guid, last_updated, https_port, reval_pending) VALUES (600, 'atlanta-edge-02', 'ga.atlanta.kabletown.net', 80, 'atlanta-edge-02\@ocdn.kabletown.net', 'X', 'bond0', '127.0.0.7', '255.255.255.252', '127.0.0.7', '2345:1234:12:d::2/64', '2345:1234:12:d::1', 9000, 100, 'RR 119.02', 300, 1, 2, null, false, 100, 100, '', '', '', '', '', '', '', '', '', '', null, '2018-01-19 21:19:32.187856', null, false);
INSERT INTO server (id, host_name, domain_name, tcp_port, xmpp_id, xmpp_passwd, interface_name, ip_address, ip_netmask, ip_gateway, ip6_address, ip6_gateway, interface_mtu, phys_location, rack, cachegroup, type, status, offline_reason, upd_pending, profile, cdn_id, mgmt_ip_address, mgmt_ip_netmask, mgmt_ip_gateway, ilo_ip_address, ilo_ip_netmask, ilo_ip_gateway, ilo_username, ilo_password, router_host_name, router_port_name, guid, last_updated, https_port, reval_pending) VALUES (700, 'atlanta-mid-02', 'ga.atlanta.kabletown.net', 80, 'atlanta-mid-02\@ocdn.kabletown.net', 'X', 'bond0', '127.0.0.8', '255.255.255.252', '127.0.0.8', '2345:1234:12:e::2/64', '2345:1234:12:e::1', 9000, 200, 'RR 119.02', 200, 2, 2, null, false, 200, 200, '', '', '', '', '', '', '', '', '', '', null, '2018-01-19 21:19:32.191292', null, false);
INSERT INTO server (id, host_name, domain_name, tcp_port, xmpp_id, xmpp_passwd, interface_name, ip_address, ip_netmask, ip_gateway, ip6_address, ip6_gateway, interface_mtu, phys_location, rack, cachegroup, type, status, offline_reason, upd_pending, profile, cdn_id, mgmt_ip_address, mgmt_ip_netmask, mgmt_ip_gateway, ilo_ip_address, ilo_ip_netmask, ilo_ip_gateway, ilo_username, ilo_password, router_host_name, router_port_name, guid, last_updated, https_port, reval_pending) VALUES (800, 'riak02', 'kabletown.net', 8088, '', '', 'eth1', '127.0.0.9', '255.255.252.0', '127.0.0.9', '2345:1234:12:f::2/64', '2345:1234:12:f::1', 1500, 200, 'RR 119.02', 100, 31, 2, null, false, 500, 100, '', '', '', '', '', '', '', '', '', '', null, '2018-01-19 21:19:32.194538', null, false);
INSERT INTO server (id, host_name, domain_name, tcp_port, xmpp_id, xmpp_passwd, interface_name, ip_address, ip_netmask, ip_gateway, ip6_address, ip6_gateway, interface_mtu, phys_location, rack, cachegroup, type, status, offline_reason, upd_pending, profile, cdn_id, mgmt_ip_address, mgmt_ip_netmask, mgmt_ip_gateway, ilo_ip_address, ilo_ip_netmask, ilo_ip_gateway, ilo_username, ilo_password, router_host_name, router_port_name, guid, last_updated, https_port, reval_pending) VALUES (900, 'influxdb01', 'kabletown.net', 8086, '', '', 'eth1', '127.0.0.10', '255.255.252.0', '127.0.0.10', '127.0.0.10', '127.0.0.10', 1500, 300, 'RR 119.02', 100, 32, 2, null, false, 500, 100, '', '', '', '', '', '', '', '', '', '', null, '2018-01-19 21:19:32.197808', null, false);
`
	err := execSQL(cfg, db, sqlStmt, "servers")
	if err != nil {
		return fmt.Errorf("exec failed %v", err)
	}
	return nil
}

// SetupTypes ...
func SetupTypes(cfg *config.Config, db *sql.DB) error {

	sqlStmt := `
INSERT INTO type (id, name, description, use_in_table, last_updated) VALUES (1, 'EDGE', 'Edge Cache', 'server', '2018-01-19 19:01:21.815104');
INSERT INTO type (id, name, description, use_in_table, last_updated) VALUES (2, 'MID', 'Mid Tier Cache', 'server', '2018-01-19 19:01:21.794365');
INSERT INTO type (id, name, description, use_in_table, last_updated) VALUES (3, 'ORG', 'Origin', 'server', '2018-01-19 19:01:21.779521');
INSERT INTO type (id, name, description, use_in_table, last_updated) VALUES (4, 'CCR', 'Kabletown Content Router', 'server', '2018-01-19 19:01:21.801776');
INSERT INTO type (id, name, description, use_in_table, last_updated) VALUES (5, 'EDGE_LOC', 'Edge Cachegroup', 'cachegroup', '2018-01-19 19:01:21.817872');
INSERT INTO type (id, name, description, use_in_table, last_updated) VALUES (6, 'MID_LOC', 'Mid Cachegroup', 'cachegroup', '2018-01-19 19:01:21.789240');
INSERT INTO type (id, name, description, use_in_table, last_updated) VALUES (7, 'DNS', 'DNS Content Routing', 'deliveryservice', '2018-01-19 19:01:21.805605');
INSERT INTO type (id, name, description, use_in_table, last_updated) VALUES (8, 'OTHER_CDN', 'Other CDN (CDS-IS, Akamai, etc)', 'server', '2018-01-19 19:01:21.807236');
INSERT INTO type (id, name, description, use_in_table, last_updated) VALUES (9, 'HTTP_NO_CACHE', 'HTTP Content Routing, no caching', 'deliveryservice', '2018-01-19 19:01:21.787978');
INSERT INTO type (id, name, description, use_in_table, last_updated) VALUES (12, 'HTTP_LIVE', 'HTTP Content routing cache in RAM ', 'deliveryservice', '2018-01-19 19:01:21.774250');
INSERT INTO type (id, name, description, use_in_table, last_updated) VALUES (14, 'RASCAL', 'Rascal health polling & reporting', 'server', '2018-01-19 19:01:21.786631');
INSERT INTO type (id, name, description, use_in_table, last_updated) VALUES (19, 'HOST_REGEXP', 'Host header regular expression', 'regex', '2018-01-19 19:01:21.804297');
INSERT INTO type (id, name, description, use_in_table, last_updated) VALUES (20, 'PATH_REGEXP', 'Path regular expression', 'regex', '2018-01-19 19:01:21.816413');
INSERT INTO type (id, name, description, use_in_table, last_updated) VALUES (21, 'A_RECORD', 'Static DNS A entry', 'staticdnsentry', '2018-01-19 19:01:21.791667');
INSERT INTO type (id, name, description, use_in_table, last_updated) VALUES (22, 'AAAA_RECORD', 'Static DNS AAAA entry', 'staticdnsentry', '2018-01-19 19:01:21.783942');
INSERT INTO type (id, name, description, use_in_table, last_updated) VALUES (23, 'CNAME_RECORD', 'Static DNS CNAME entry', 'staticdnsentry', '2018-01-19 19:01:21.795551');
INSERT INTO type (id, name, description, use_in_table, last_updated) VALUES (24, 'HTTP_LIVE_NATNL', 'HTTP Content routing, RAM cache, National', 'deliveryservice', '2018-01-19 19:01:21.796784');
INSERT INTO type (id, name, description, use_in_table, last_updated) VALUES (27, 'DNS_LIVE_NATNL', 'DNS Content routing, RAM cache, National', 'deliveryservice', '2018-01-19 19:01:21.790471');
INSERT INTO type (id, name, description, use_in_table, last_updated) VALUES (28, 'LOCAL', 'Local User', 'tm_user', '2018-01-19 19:01:21.808970');
INSERT INTO type (id, name, description, use_in_table, last_updated) VALUES (29, 'ACTIVE_DIRECTORY', 'Active Directory User', 'tm_user', '2018-01-19 19:01:21.799178');
INSERT INTO type (id, name, description, use_in_table, last_updated) VALUES (30, 'TOOLS_SERVER', 'Ops hosts for management', 'server', '2018-01-19 19:01:21.797996');
INSERT INTO type (id, name, description, use_in_table, last_updated) VALUES (31, 'RIAK', 'riak type', 'server', '2018-01-19 19:01:21.819171');
INSERT INTO type (id, name, description, use_in_table, last_updated) VALUES (32, 'INFLUXDB', 'influxdb type', 'server', '2018-01-19 19:01:21.803064');
INSERT INTO type (id, name, description, use_in_table, last_updated) VALUES (33, 'RESOLVE4', 'federation type resolve4', 'federation', '2018-01-19 19:01:21.800497');
INSERT INTO type (id, name, description, use_in_table, last_updated) VALUES (34, 'RESOLVE6', 'federation type resolve6', 'federation', '2018-01-19 19:01:21.792993');
INSERT INTO type (id, name, description, use_in_table, last_updated) VALUES (35, 'ANY_MAP', 'any_map type', 'deliveryservice', '2018-01-19 19:01:21.780894');
INSERT INTO type (id, name, description, use_in_table, last_updated) VALUES (36, 'HTTP', 'HTTP Content routing cache ', 'deliveryservice', '2018-01-19 19:01:21.813811');
INSERT INTO type (id, name, description, use_in_table, last_updated) VALUES (37, 'STEERING', 'Steering Delivery Service', 'deliveryservice', '2018-01-19 19:01:21.785303');
INSERT INTO type (id, name, description, use_in_table, last_updated) VALUES (38, 'CLIENT_STEERING', 'Client-Controlled Steering Delivery Service', 'deliveryservice', '2018-01-19 19:01:21.782467');
INSERT INTO type (id, name, description, use_in_table, last_updated) VALUES (39, 'STEERING_WEIGHT', 'Weighted steering target', 'steering_target', '2018-01-19 19:01:21.812447');
INSERT INTO type (id, name, description, use_in_table, last_updated) VALUES (40, 'STEERING_ORDER', 'Ordered steering target', 'steering_target', '2018-01-19 19:01:21.810875');
`
	err := execSQL(cfg, db, sqlStmt, "type")
	if err != nil {
		return fmt.Errorf("exec failed %v", err)
	}

	return nil
}

// SetupParameters ...
func SetupParameters(cfg *config.Config, db *sql.DB) error {

	sqlStmt := `
INSERT INTO parameter (id, name, config_file, value, last_updated, secure) VALUES (4, 'health.threshold.loadavg', 'rascal.properties', '25.0', '2018-01-19 19:01:21.455131', false);
INSERT INTO parameter (id, name, config_file, value, last_updated, secure) VALUES (5, 'health.threshold.availableBandwidthInKbps', 'rascal.properties', '>1750000', '2018-01-19 19:01:21.472279', false);
INSERT INTO parameter (id, name, config_file, value, last_updated, secure) VALUES (6, 'history.count', 'rascal.properties', '30', '2018-01-19 19:01:21.489534', false);
INSERT INTO parameter (id, name, config_file, value, last_updated, secure) VALUES (7, 'key0', 'url_sig_cdl-c2.config', 'HOOJ3Ghq1x4gChp3iQkqVTcPlOj8UCi3', '2018-01-19 19:01:21.503311', false);
INSERT INTO parameter (id, name, config_file, value, last_updated, secure) VALUES (8, 'key1', 'url_sig_cdl-c2.config', '_9LZYkRnfCS0rCBF7fTQzM9Scwlp2FhO', '2018-01-19 19:01:21.505157', false);
INSERT INTO parameter (id, name, config_file, value, last_updated, secure) VALUES (9, 'key2', 'url_sig_cdl-c2.config', 'AFpkxfc4oTiyFSqtY6_ohjt3V80aAIxS', '2018-01-19 19:01:21.508548', false);
INSERT INTO parameter (id, name, config_file, value, last_updated, secure) VALUES (10, 'key3', 'url_sig_cdl-c2.config', 'AL9kzs_SXaRZjPWH8G5e2m4ByTTzkzlc', '2018-01-19 19:01:21.401781', false);
INSERT INTO parameter (id, name, config_file, value, last_updated, secure) VALUES (11, 'key4', 'url_sig_cdl-c2.config', 'poP3n3szbD1U4vx1xQXV65BvkVgWzfN8', '2018-01-19 19:01:21.406601', false);
INSERT INTO parameter (id, name, config_file, value, last_updated, secure) VALUES (12, 'key5', 'url_sig_cdl-c2.config', '1ir32ng4C4w137p5oq72kd2wqmIZUrya', '2018-01-19 19:01:21.408784', false);
INSERT INTO parameter (id, name, config_file, value, last_updated, secure) VALUES (13, 'key6', 'url_sig_cdl-c2.config', 'B1qLptn2T1b_iXeTCWDcVuYvANtH139f', '2018-01-19 19:01:21.410854', false);
INSERT INTO parameter (id, name, config_file, value, last_updated, secure) VALUES (14, 'key7', 'url_sig_cdl-c2.config', 'PiCV_5OODMzBbsNFMWsBxcQ8v1sK0TYE', '2018-01-19 19:01:21.412716', false);
INSERT INTO parameter (id, name, config_file, value, last_updated, secure) VALUES (15, 'key8', 'url_sig_cdl-c2.config', 'Ggpv6DqXDvt2s1CETPBpNKwaLk4fTM9l', '2018-01-19 19:01:21.414638', false);
INSERT INTO parameter (id, name, config_file, value, last_updated, secure) VALUES (16, 'key9', 'url_sig_cdl-c2.config', 'qPlVT_s6kL37aqb6hipDm4Bt55S72mI7', '2018-01-19 19:01:21.416551', false);
INSERT INTO parameter (id, name, config_file, value, last_updated, secure) VALUES (17, 'key10', 'url_sig_cdl-c2.config', 'BsI5A9EmWrobIS1FeuOs1z9fm2t2WSBe', '2018-01-19 19:01:21.418689', false);
INSERT INTO parameter (id, name, config_file, value, last_updated, secure) VALUES (18, 'key11', 'url_sig_cdl-c2.config', 'A54y66NCIj897GjS4yA9RrsSPtCUnQXP', '2018-01-19 19:01:21.420467', false);
INSERT INTO parameter (id, name, config_file, value, last_updated, secure) VALUES (19, 'key12', 'url_sig_cdl-c2.config', '2jZH0NDPSJttIr4c2KP510f47EKqTQAu', '2018-01-19 19:01:21.422414', false);
INSERT INTO parameter (id, name, config_file, value, last_updated, secure) VALUES (20, 'key13', 'url_sig_cdl-c2.config', 'XduT2FBjBmmVID5JRB5LEf9oR5QDtBgC', '2018-01-19 19:01:21.424435', false);
INSERT INTO parameter (id, name, config_file, value, last_updated, secure) VALUES (21, 'key14', 'url_sig_cdl-c2.config', 'D9nH0SvK_0kP5w8QNd1UFJ28ulFkFKPn', '2018-01-19 19:01:21.426125', false);
INSERT INTO parameter (id, name, config_file, value, last_updated, secure) VALUES (22, 'key15', 'url_sig_cdl-c2.config', 'udKXWYNwbXXweaaLzaKDGl57OixnIIcm', '2018-01-19 19:01:21.427797', false);
INSERT INTO parameter (id, name, config_file, value, last_updated, secure) VALUES (23, 'location', 'url_sig_cdl-c2.config', '/opt/trafficserver/etc/trafficserver', '2018-01-19 19:01:21.429365', false);
INSERT INTO parameter (id, name, config_file, value, last_updated, secure) VALUES (24, 'error_url', 'url_sig_cdl-c2.config', '403', '2018-01-19 19:01:21.431062', false);
INSERT INTO parameter (id, name, config_file, value, last_updated, secure) VALUES (25, 'CONFIG proxy.config.allocator.debug_filter', 'records.config', 'INT 0', '2018-01-19 19:01:21.432692', false);
INSERT INTO parameter (id, name, config_file, value, last_updated, secure) VALUES (26, 'CONFIG proxy.config.allocator.enable_reclaim', 'records.config', 'INT 0', '2018-01-19 19:01:21.434425', false);
INSERT INTO parameter (id, name, config_file, value, last_updated, secure) VALUES (27, 'CONFIG proxy.config.allocator.max_overage', 'records.config', 'INT 3', '2018-01-19 19:01:21.435957', false);
INSERT INTO parameter (id, name, config_file, value, last_updated, secure) VALUES (28, 'CONFIG proxy.config.diags.show_location', 'records.config', 'INT 0', '2018-01-19 19:01:21.437496', false);
INSERT INTO parameter (id, name, config_file, value, last_updated, secure) VALUES (29, 'CONFIG proxy.config.http.cache.allow_empty_doc', 'records.config', 'INT 0', '2018-01-19 19:01:21.439033', false);
INSERT INTO parameter (id, name, config_file, value, last_updated, secure) VALUES (30, 'LOCAL proxy.config.cache.interim.storage', 'records.config', 'STRING NULL', '2018-01-19 19:01:21.440502', false);
INSERT INTO parameter (id, name, config_file, value, last_updated, secure) VALUES (31, 'CONFIG proxy.config.http.parent_proxy.file', 'records.config', 'STRING parent.config', '2018-01-19 19:01:21.441933', false);
INSERT INTO parameter (id, name, config_file, value, last_updated, secure) VALUES (32, 'location', '12M_facts', '/opt/ort', '2018-01-19 19:01:21.443436', false);
INSERT INTO parameter (id, name, config_file, value, last_updated, secure) VALUES (33, 'location', 'cacheurl.config', '/opt/trafficserver/etc/trafficserver/', '2018-01-19 19:01:21.444898', false);
INSERT INTO parameter (id, name, config_file, value, last_updated, secure) VALUES (34, 'location', 'ip_allow.config', '/opt/trafficserver/etc/trafficserver', '2018-01-19 19:01:21.446396', false);
INSERT INTO parameter (id, name, config_file, value, last_updated, secure) VALUES (35, 'astats_over_http.so', 'plugin.config', '_astats 33.101.99.100,172.39.19.39,172.39.19.49,172.39.19.49,172.39.29.49', '2018-01-19 19:01:21.447837', false);
INSERT INTO parameter (id, name, config_file, value, last_updated, secure) VALUES (36, 'location', 'crontab_root', '/var/spool/cron', '2018-01-19 19:01:21.449259', false);
INSERT INTO parameter (id, name, config_file, value, last_updated, secure) VALUES (37, 'location', 'hdr_rw_cdl-c2.config', '/opt/trafficserver/etc/trafficserver', '2018-01-19 19:01:21.450778', false);
INSERT INTO parameter (id, name, config_file, value, last_updated, secure) VALUES (38, 'location', '50-ats.rules', '/etc/udev/rules.d/', '2018-01-19 19:01:21.452196', false);
INSERT INTO parameter (id, name, config_file, value, last_updated, secure) VALUES (39, 'location', 'parent.config', '/opt/trafficserver/etc/trafficserver/', '2018-01-19 19:01:21.453716', false);
INSERT INTO parameter (id, name, config_file, value, last_updated, secure) VALUES (40, 'location', 'remap.config', '/opt/trafficserver/etc/trafficserver/', '2018-01-19 19:01:21.456753', false);
INSERT INTO parameter (id, name, config_file, value, last_updated, secure) VALUES (41, 'location', 'drop_qstring.config', '/opt/trafficserver/etc/trafficserver', '2018-01-19 19:01:21.458350', false);
INSERT INTO parameter (id, name, config_file, value, last_updated, secure) VALUES (42, 'LogFormat.Format', 'logs_xml.config', '%<cqtq> chi=%<chi> phn=%<phn> shn=%<shn> url=%<cquuc> cqhm=%<cqhm> cqhv=%<cqhv> pssc=%<pssc> ttms=%<ttms> b=%<pscl> sssc=%<sssc> sscl=%<sscl> cfsc=%<cfsc> pfsc=%<pfsc> crc=%<crc> phr=%<phr> uas="%<{User-Agent}cqh>"', '2018-01-19 19:01:21.459788', false);
INSERT INTO parameter (id, name, config_file, value, last_updated, secure) VALUES (43, 'LogFormat.Name', 'logs_xml.config', 'custom_ats_2', '2018-01-19 19:01:21.461206', false);
INSERT INTO parameter (id, name, config_file, value, last_updated, secure) VALUES (44, 'LogObject.Format', 'logs_xml.config', 'custom_ats_2', '2018-01-19 19:01:21.462772', false);
INSERT INTO parameter (id, name, config_file, value, last_updated, secure) VALUES (45, 'LogObject.Filename', 'logs_xml.config', 'custom_ats_2', '2018-01-19 19:01:21.464259', false);
INSERT INTO parameter (id, name, config_file, value, last_updated, secure) VALUES (46, 'location', 'cache.config', '/opt/trafficserver/etc/trafficserver/', '2018-01-19 19:01:21.465717', false);
INSERT INTO parameter (id, name, config_file, value, last_updated, secure) VALUES (47, 'CONFIG proxy.config.cache.control.filename', 'records.config', 'STRING cache.config', '2018-01-19 19:01:21.467349', false);
INSERT INTO parameter (id, name, config_file, value, last_updated, secure) VALUES (48, 'regex_revalidate.so', 'plugin.config', '--config regex_revalidate.config', '2018-01-19 19:01:21.469075', false);
INSERT INTO parameter (id, name, config_file, value, last_updated, secure) VALUES (49, 'location', 'regex_revalidate.config', '/opt/trafficserver/etc/trafficserver', '2018-01-19 19:01:21.470677', false);
INSERT INTO parameter (id, name, config_file, value, last_updated, secure) VALUES (50, 'location', 'hosting.config', '/opt/trafficserver/etc/trafficserver/', '2018-01-19 19:01:21.474023', false);
INSERT INTO parameter (id, name, config_file, value, last_updated, secure) VALUES (51, 'location', 'volume.config', '/opt/trafficserver/etc/trafficserver/', '2018-01-19 19:01:21.475515', false);
INSERT INTO parameter (id, name, config_file, value, last_updated, secure) VALUES (52, 'allow_ip', 'astats.config', '127.0.0.1,172.39.0.0/16,33.101.99.0/24', '2018-01-19 19:01:21.477074', false);
INSERT INTO parameter (id, name, config_file, value, last_updated, secure) VALUES (53, 'allow_ip6', 'astats.config', '::1,2033:D011:3300::336/64,2033:D011:3300::335/64,2033:D021:3300::333/64,2033:D021:3300::334/64', '2018-01-19 19:01:21.478516', false);
INSERT INTO parameter (id, name, config_file, value, last_updated, secure) VALUES (54, 'record_types', 'astats.config', '144', '2018-01-19 19:01:21.480143', false);
INSERT INTO parameter (id, name, config_file, value, last_updated, secure) VALUES (55, 'location', 'astats.config', '/opt/trafficserver/etc/trafficserver', '2018-01-19 19:01:21.481582', false);
INSERT INTO parameter (id, name, config_file, value, last_updated, secure) VALUES (56, 'path', 'astats.config', '_astats', '2018-01-19 19:01:21.482959', false);
INSERT INTO parameter (id, name, config_file, value, last_updated, secure) VALUES (57, 'location', 'storage.config', '/opt/trafficserver/etc/trafficserver/', '2018-01-19 19:01:21.484501', false);
INSERT INTO parameter (id, name, config_file, value, last_updated, secure) VALUES (58, 'Drive_Prefix', 'storage.config', '/dev/sd', '2018-01-19 19:01:21.486250', false);
INSERT INTO parameter (id, name, config_file, value, last_updated, secure) VALUES (59, 'Drive_Letters', 'storage.config', 'b,c,d,e,f,g,h,i,j,k,l,m,n,o', '2018-01-19 19:01:21.487958', false);
INSERT INTO parameter (id, name, config_file, value, last_updated, secure) VALUES (60, 'Disk_Volume', 'storage.config', '1', '2018-01-19 19:01:21.491181', false);
INSERT INTO parameter (id, name, config_file, value, last_updated, secure) VALUES (61, 'CONFIG proxy.config.hostdb.storage_size', 'records.config', 'INT 33554432', '2018-01-19 19:01:21.492850', false);
INSERT INTO parameter (id, name, config_file, value, last_updated, secure) VALUES (63, 'maxRevalDurationDays', 'regex_revalidate.config', '3', '2018-01-19 19:01:21.494468', false);
INSERT INTO parameter (id, name, config_file, value, last_updated, secure) VALUES (64, 'maxRevalDurationDays', 'regex_revalidate.config', '90', '2018-01-19 19:01:21.496195', false);
INSERT INTO parameter (id, name, config_file, value, last_updated, secure) VALUES (65, 'unassigned_parameter_1', 'whaterver.config', '852', '2018-01-19 19:01:21.497838', false);
INSERT INTO parameter (id, name, config_file, value, last_updated, secure) VALUES (66, 'trafficserver', 'package', '5.3.2-765.f4354b9.el7.centos.x86_64', '2018-01-19 19:01:21.499423', false);
INSERT INTO parameter (id, name, config_file, value, last_updated, secure) VALUES (67, 'use_tenancy', 'global', '1', '2018-01-19 19:01:21.501151', false);
`
	err := execSQL(cfg, db, sqlStmt, "parameter")
	if err != nil {
		return fmt.Errorf("exec failed %v", err)
	}

	return nil
}

// SetupProfiles ...
func SetupProfiles(cfg *config.Config, db *sql.DB) error {

	sqlStmt := `
INSERT INTO profile (id, name, description, last_updated, type, cdn, routing_disabled) VALUES (100, 'EDGE1', 'edge description', '2018-01-19 19:01:21.512005', 'ATS_PROFILE', 100, false);
INSERT INTO profile (id, name, description, last_updated, type, cdn, routing_disabled) VALUES (200, 'MID1', 'mid description', '2018-01-19 19:01:21.517781', 'ATS_PROFILE', 100, false);
INSERT INTO profile (id, name, description, last_updated, type, cdn, routing_disabled) VALUES (300, 'CCR1', 'ccr description', '2018-01-19 19:01:21.521121', 'TR_PROFILE', 100, false);
INSERT INTO profile (id, name, description, last_updated, type, cdn, routing_disabled) VALUES (301, 'CCR2', 'ccr description', '2018-01-19 19:01:21.524584', 'TR_PROFILE', 200, false);
INSERT INTO profile (id, name, description, last_updated, type, cdn, routing_disabled) VALUES (500, 'RIAK1', 'riak description', '2018-01-19 19:01:21.528911', 'RIAK_PROFILE', 100, false);
INSERT INTO profile (id, name, description, last_updated, type, cdn, routing_disabled) VALUES (600, 'RASCAL1', 'rascal description', '2018-01-19 19:01:21.532539', 'TM_PROFILE', 100, false);
INSERT INTO profile (id, name, description, last_updated, type, cdn, routing_disabled) VALUES (700, 'RASCAL2', 'rascal2 description', '2018-01-19 19:01:21.536447', 'TM_PROFILE', 200, false);
INSERT INTO profile (id, name, description, last_updated, type, cdn, routing_disabled) VALUES (8, 'MISC', 'misc profile description', '2018-01-19 19:01:21.539022', 'UNK_PROFILE', null, false);
INSERT INTO profile (id, name, description, last_updated, type, cdn, routing_disabled) VALUES (900, 'EDGE2', 'edge description', '2018-01-19 19:01:21.541300', 'ATS_PROFILE', 200, false);
`
	err := execSQL(cfg, db, sqlStmt, "profile")
	if err != nil {
		return fmt.Errorf("exec failed %v", err)
	}
	return nil
}

// SetupProfileParameters ...
func SetupProfileParameters(cfg *config.Config, db *sql.DB) error {

	sqlStmt := `
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (200, 43, '2018-01-19 19:01:21.556526');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (100, 19, '2018-01-19 19:01:21.566442');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (100, 35, '2018-01-19 19:01:21.571364');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (200, 49, '2018-01-19 19:01:21.575178');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (200, 61, '2018-01-19 19:01:21.578744');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (100, 9, '2018-01-19 19:01:21.582534');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (200, 46, '2018-01-19 19:01:21.586388');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (200, 35, '2018-01-19 19:01:21.588145');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (200, 16, '2018-01-19 19:01:21.589542');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (200, 57, '2018-01-19 19:01:21.591061');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (200, 48, '2018-01-19 19:01:21.592700');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (100, 60, '2018-01-19 19:01:21.594185');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (200, 31, '2018-01-19 19:01:21.595700');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (100, 49, '2018-01-19 19:01:21.597212');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (200, 4, '2018-01-19 19:01:21.598744');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (100, 36, '2018-01-19 19:01:21.600582');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (200, 27, '2018-01-19 19:01:21.602214');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (100, 41, '2018-01-19 19:01:21.604015');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (100, 16, '2018-01-19 19:01:21.605612');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (200, 17, '2018-01-19 19:01:21.607234');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (100, 21, '2018-01-19 19:01:21.609358');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (200, 41, '2018-01-19 19:01:21.611101');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (100, 32, '2018-01-19 19:01:21.613078');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (200, 32, '2018-01-19 19:01:21.614943');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (200, 28, '2018-01-19 19:01:21.616641');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (200, 6, '2018-01-19 19:01:21.618677');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (100, 66, '2018-01-19 19:01:21.620617');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (200, 58, '2018-01-19 19:01:21.622399');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (100, 28, '2018-01-19 19:01:21.623955');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (200, 56, '2018-01-19 19:01:21.625664');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (200, 23, '2018-01-19 19:01:21.627471');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (100, 11, '2018-01-19 19:01:21.629284');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (200, 30, '2018-01-19 19:01:21.630989');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (100, 22, '2018-01-19 19:01:21.632523');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (100, 23, '2018-01-19 19:01:21.634278');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (100, 37, '2018-01-19 19:01:21.635945');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (100, 25, '2018-01-19 19:01:21.637627');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (100, 38, '2018-01-19 19:01:21.639252');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (100, 52, '2018-01-19 19:01:21.640775');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (100, 29, '2018-01-19 19:01:21.642278');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (200, 12, '2018-01-19 19:01:21.644071');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (200, 45, '2018-01-19 19:01:21.645614');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (200, 60, '2018-01-19 19:01:21.647126');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (100, 26, '2018-01-19 19:01:21.648787');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (100, 57, '2018-01-19 19:01:21.650507');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (100, 13, '2018-01-19 19:01:21.652142');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (100, 27, '2018-01-19 19:01:21.653714');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (200, 26, '2018-01-19 19:01:21.655383');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (100, 39, '2018-01-19 19:01:21.657078');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (100, 12, '2018-01-19 19:01:21.658901');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (200, 25, '2018-01-19 19:01:21.661010');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (200, 21, '2018-01-19 19:01:21.662865');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (100, 33, '2018-01-19 19:01:21.664561');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (200, 38, '2018-01-19 19:01:21.666336');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (100, 34, '2018-01-19 19:01:21.668286');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (100, 58, '2018-01-19 19:01:21.670053');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (100, 24, '2018-01-19 19:01:21.671744');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (100, 43, '2018-01-19 19:01:21.673493');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (200, 5, '2018-01-19 19:01:21.675218');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (200, 37, '2018-01-19 19:01:21.676721');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (100, 51, '2018-01-19 19:01:21.678334');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (200, 19, '2018-01-19 19:01:21.679937');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (100, 6, '2018-01-19 19:01:21.681398');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (100, 18, '2018-01-19 19:01:21.682983');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (100, 42, '2018-01-19 19:01:21.684568');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (100, 7, '2018-01-19 19:01:21.686083');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (100, 56, '2018-01-19 19:01:21.687549');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (200, 13, '2018-01-19 19:01:21.689131');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (200, 22, '2018-01-19 19:01:21.690719');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (200, 36, '2018-01-19 19:01:21.692254');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (200, 53, '2018-01-19 19:01:21.693745');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (200, 40, '2018-01-19 19:01:21.695556');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (200, 51, '2018-01-19 19:01:21.697784');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (100, 46, '2018-01-19 19:01:21.699385');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (200, 11, '2018-01-19 19:01:21.701103');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (100, 54, '2018-01-19 19:01:21.702727');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (100, 17, '2018-01-19 19:01:21.704304');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (100, 53, '2018-01-19 19:01:21.705942');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (200, 10, '2018-01-19 19:01:21.707676');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (100, 8, '2018-01-19 19:01:21.709391');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (200, 39, '2018-01-19 19:01:21.711213');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (100, 40, '2018-01-19 19:01:21.713199');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (200, 29, '2018-01-19 19:01:21.715051');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (100, 59, '2018-01-19 19:01:21.716817');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (100, 47, '2018-01-19 19:01:21.718642');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (100, 44, '2018-01-19 19:01:21.720315');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (200, 9, '2018-01-19 19:01:21.722063');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (200, 8, '2018-01-19 19:01:21.723607');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (100, 20, '2018-01-19 19:01:21.725403');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (100, 48, '2018-01-19 19:01:21.727060');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (200, 55, '2018-01-19 19:01:21.728640');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (100, 10, '2018-01-19 19:01:21.730182');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (100, 45, '2018-01-19 19:01:21.731780');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (200, 15, '2018-01-19 19:01:21.733368');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (200, 33, '2018-01-19 19:01:21.734950');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (200, 50, '2018-01-19 19:01:21.736646');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (200, 52, '2018-01-19 19:01:21.738319');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (100, 14, '2018-01-19 19:01:21.739900');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (200, 14, '2018-01-19 19:01:21.741450');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (200, 18, '2018-01-19 19:01:21.743105');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (100, 61, '2018-01-19 19:01:21.744826');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (200, 44, '2018-01-19 19:01:21.746391');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (100, 55, '2018-01-19 19:01:21.747999');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (200, 59, '2018-01-19 19:01:21.749519');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (200, 34, '2018-01-19 19:01:21.751253');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (200, 24, '2018-01-19 19:01:21.753005');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (200, 7, '2018-01-19 19:01:21.754576');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (100, 15, '2018-01-19 19:01:21.757250');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (200, 47, '2018-01-19 19:01:21.759781');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (200, 54, '2018-01-19 19:01:21.761829');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (200, 42, '2018-01-19 19:01:21.763902');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (100, 50, '2018-01-19 19:01:21.765912');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (100, 31, '2018-01-19 19:01:21.767998');
INSERT INTO profile_parameter (profile, parameter, last_updated) VALUES (200, 20, '2018-01-19 19:01:21.769919');
`
	err := execSQL(cfg, db, sqlStmt, "profile_parameter")
	if err != nil {
		return fmt.Errorf("exec failed %v", err)
	}
	return nil
}

// SetupDivisions ...
func SetupDivisions(cfg *config.Config, db *sql.DB) error {

	sqlStmt := `
INSERT INTO division (id, name, last_updated) VALUES (100, 'mountain', '2018-01-19 19:01:21.851102');
`
	err := execSQL(cfg, db, sqlStmt, "division")
	if err != nil {
		return fmt.Errorf("exec failed %v", err)
	}
	return nil
}

// SetupRegions ...
func SetupRegions(cfg *config.Config, db *sql.DB) error {

	sqlStmt := `
INSERT INTO region (id, name, division, last_updated) VALUES (100, 'Denver Region', 100, '2018-01-19 19:01:21.859430');
INSERT INTO region (id, name, division, last_updated) VALUES (200, 'Boulder Region', 100, '2018-01-19 19:01:21.854509');
`
	err := execSQL(cfg, db, sqlStmt, "region")
	if err != nil {
		return fmt.Errorf("exec failed %v", err)
	}
	return nil
}

// SetupAsns ...
func SetupAsns(cfg *config.Config, db *sql.DB) error {

	sqlStmt := `
INSERT INTO asn (id, asn, cachegroup, last_updated) VALUES (100, 9939, 100, '2018-01-19 19:01:21.995075');
INSERT INTO asn (id, asn, cachegroup, last_updated) VALUES (200, 9940, 200, '2018-01-19 19:01:22.005683');
`
	err := execSQL(cfg, db, sqlStmt, "asn")
	if err != nil {
		return fmt.Errorf("exec failed %v", err)
	}
	return nil
}

// SetupDeliveryServices ...
func SetupDeliveryServices(cfg *config.Config, db *sql.DB) error {

	sqlStmt := `
INSERT INTO deliveryservice (id, xml_id, active, dscp, signing_algorithm, qstring_ignore, geo_limit, http_bypass_fqdn, dns_bypass_ip, dns_bypass_ip6, dns_bypass_ttl, org_server_fqdn, type, profile, cdn_id, ccr_dns_ttl, global_max_mbps, global_max_tps, long_desc, long_desc_1, long_desc_2, max_dns_answers, info_url, miss_lat, miss_long, check_path, last_updated, protocol, ssl_key_version, ipv6_routing_enabled, range_request_handling, edge_header_rewrite, origin_shield, mid_header_rewrite, regex_remap, cacheurl, remap_text, multi_site_origin, display_name, tr_response_headers, initial_dispersion, dns_bypass_cname, tr_request_headers, regional_geo_blocking, geo_provider, geo_limit_countries, logs_enabled, multi_site_origin_algorithm, geolimit_redirect_url, tenant_id, routing_name) VALUES (100, 'test-ds1', true, 40, null, 0, 0, '', '', '', -1, 'http://test-ds1.edge', 21, 100, 100, 3600, 0, 0, 'test-ds1 long_desc', 'test-ds1 long_desc_1', 'test-ds1 long_desc_2', 0, 'http://test-ds1.edge/info_url.html', 41.881944, -87.627778, '/crossdomain.xml', '2018-01-19 21:19:32.217466', 1, 0, true, 0, null, null, null, null, null, null, false, 'test-ds1-displayname', null, 1, null, null, true, 0, null, true, null, null, null, 'foo');
INSERT INTO deliveryservice (id, xml_id, active, dscp, signing_algorithm, qstring_ignore, geo_limit, http_bypass_fqdn, dns_bypass_ip, dns_bypass_ip6, dns_bypass_ttl, org_server_fqdn, type, profile, cdn_id, ccr_dns_ttl, global_max_mbps, global_max_tps, long_desc, long_desc_1, long_desc_2, max_dns_answers, info_url, miss_lat, miss_long, check_path, last_updated, protocol, ssl_key_version, ipv6_routing_enabled, range_request_handling, edge_header_rewrite, origin_shield, mid_header_rewrite, regex_remap, cacheurl, remap_text, multi_site_origin, display_name, tr_response_headers, initial_dispersion, dns_bypass_cname, tr_request_headers, regional_geo_blocking, geo_provider, geo_limit_countries, logs_enabled, multi_site_origin_algorithm, geolimit_redirect_url, tenant_id, routing_name) VALUES (1000, 'steering-target-ds1', true, 40, null, 0, 0, '', 'hokeypokey', null, 10, 'http://target-ds1.edge', 21, 100, 100, 3600, 0, 0, 'target-ds1 long_desc', 'target-ds1 long_desc_1', 'target-ds1 long_desc_2', 0, 'http://target-ds1.edge/info_url.html', 41.881944, -87.627778, '/crossdomain.xml', '2018-01-19 21:19:32.226858', 1, 0, true, 0, null, null, null, null, null, null, false, 'target-ds1-displayname', null, 1, null, null, true, 0, null, false, null, null, null, 'foo');
INSERT INTO deliveryservice (id, xml_id, active, dscp, signing_algorithm, qstring_ignore, geo_limit, http_bypass_fqdn, dns_bypass_ip, dns_bypass_ip6, dns_bypass_ttl, org_server_fqdn, type, profile, cdn_id, ccr_dns_ttl, global_max_mbps, global_max_tps, long_desc, long_desc_1, long_desc_2, max_dns_answers, info_url, miss_lat, miss_long, check_path, last_updated, protocol, ssl_key_version, ipv6_routing_enabled, range_request_handling, edge_header_rewrite, origin_shield, mid_header_rewrite, regex_remap, cacheurl, remap_text, multi_site_origin, display_name, tr_response_headers, initial_dispersion, dns_bypass_cname, tr_request_headers, regional_geo_blocking, geo_provider, geo_limit_countries, logs_enabled, multi_site_origin_algorithm, geolimit_redirect_url, tenant_id, routing_name) VALUES (1100, 'steering-target-ds2', true, 40, null, 0, 0, '', 'hokeypokey', null, 10, 'http://target-ds2.edge', 21, 100, 100, 3600, 0, 0, 'target-ds2 long_desc', 'target-ds2 long_desc_1', 'target-ds2 long_desc_2', 0, 'http://target-ds2.edge/info_url.html', 41.881944, -87.627778, '/crossdomain.xml', '2018-01-19 21:19:32.235025', 1, 0, true, 0, null, null, null, null, null, null, false, 'target-ds2-displayname', null, 1, null, null, true, 0, null, false, null, null, null, 'foo');
INSERT INTO deliveryservice (id, xml_id, active, dscp, signing_algorithm, qstring_ignore, geo_limit, http_bypass_fqdn, dns_bypass_ip, dns_bypass_ip6, dns_bypass_ttl, org_server_fqdn, type, profile, cdn_id, ccr_dns_ttl, global_max_mbps, global_max_tps, long_desc, long_desc_1, long_desc_2, max_dns_answers, info_url, miss_lat, miss_long, check_path, last_updated, protocol, ssl_key_version, ipv6_routing_enabled, range_request_handling, edge_header_rewrite, origin_shield, mid_header_rewrite, regex_remap, cacheurl, remap_text, multi_site_origin, display_name, tr_response_headers, initial_dispersion, dns_bypass_cname, tr_request_headers, regional_geo_blocking, geo_provider, geo_limit_countries, logs_enabled, multi_site_origin_algorithm, geolimit_redirect_url, tenant_id, routing_name) VALUES (1200, 'steering-target-ds3', true, 40, null, 0, 0, '', 'hokeypokey', null, 10, 'http://target-ds3.edge', 21, 100, 100, 3600, 0, 0, 'target-ds3 long_desc', 'target-ds3 long_desc_1', 'target-ds3 long_desc_2', 0, 'http://target-ds3.edge/info_url.html', 41.881944, -87.627778, '/crossdomain.xml', '2018-01-19 21:19:32.241327', 1, 0, true, 0, null, null, null, null, null, null, false, 'target-ds3-displayname', null, 1, null, null, true, 0, null, false, null, null, null, 'foo');
INSERT INTO deliveryservice (id, xml_id, active, dscp, signing_algorithm, qstring_ignore, geo_limit, http_bypass_fqdn, dns_bypass_ip, dns_bypass_ip6, dns_bypass_ttl, org_server_fqdn, type, profile, cdn_id, ccr_dns_ttl, global_max_mbps, global_max_tps, long_desc, long_desc_1, long_desc_2, max_dns_answers, info_url, miss_lat, miss_long, check_path, last_updated, protocol, ssl_key_version, ipv6_routing_enabled, range_request_handling, edge_header_rewrite, origin_shield, mid_header_rewrite, regex_remap, cacheurl, remap_text, multi_site_origin, display_name, tr_response_headers, initial_dispersion, dns_bypass_cname, tr_request_headers, regional_geo_blocking, geo_provider, geo_limit_countries, logs_enabled, multi_site_origin_algorithm, geolimit_redirect_url, tenant_id, routing_name) VALUES (1300, 'steering-target-ds4', true, 40, null, 0, 0, '', 'hokeypokey', null, 10, 'http://target-ds4.edge', 21, 100, 100, 3600, 0, 0, 'target-ds4 long_desc', 'target-ds4 long_desc_1', 'target-ds4 long_desc_2', 0, 'http://target-ds4.edge/info_url.html', 41.881944, -87.627778, '/crossdomain.xml', '2018-01-19 21:19:32.247731', 1, 0, true, 0, null, null, null, null, null, null, false, 'target-ds4-displayname', null, 1, null, null, true, 0, null, false, null, null, null, 'foo');
INSERT INTO deliveryservice (id, xml_id, active, dscp, signing_algorithm, qstring_ignore, geo_limit, http_bypass_fqdn, dns_bypass_ip, dns_bypass_ip6, dns_bypass_ttl, org_server_fqdn, type, profile, cdn_id, ccr_dns_ttl, global_max_mbps, global_max_tps, long_desc, long_desc_1, long_desc_2, max_dns_answers, info_url, miss_lat, miss_long, check_path, last_updated, protocol, ssl_key_version, ipv6_routing_enabled, range_request_handling, edge_header_rewrite, origin_shield, mid_header_rewrite, regex_remap, cacheurl, remap_text, multi_site_origin, display_name, tr_response_headers, initial_dispersion, dns_bypass_cname, tr_request_headers, regional_geo_blocking, geo_provider, geo_limit_countries, logs_enabled, multi_site_origin_algorithm, geolimit_redirect_url, tenant_id, routing_name) VALUES (200, 'test-ds2', true, 40, null, 0, 0, '', '', '', -1, 'http://test-ds2.edge', 9, 100, 100, 3600, 0, 0, 'test-ds2 long_desc', 'test-ds2 long_desc_1', 'test-ds2 long_desc_2', 0, 'http://test-ds2.edge/info_url.html', 41.881944, -87.627778, '/crossdomain.xml', '2018-01-19 21:19:32.253469', 0, 0, false, 0, null, null, null, null, null, null, false, 'test-ds2-displayname', null, 1, null, null, false, 0, null, false, null, null, null, 'foo');
INSERT INTO deliveryservice (id, xml_id, active, dscp, signing_algorithm, qstring_ignore, geo_limit, http_bypass_fqdn, dns_bypass_ip, dns_bypass_ip6, dns_bypass_ttl, org_server_fqdn, type, profile, cdn_id, ccr_dns_ttl, global_max_mbps, global_max_tps, long_desc, long_desc_1, long_desc_2, max_dns_answers, info_url, miss_lat, miss_long, check_path, last_updated, protocol, ssl_key_version, ipv6_routing_enabled, range_request_handling, edge_header_rewrite, origin_shield, mid_header_rewrite, regex_remap, cacheurl, remap_text, multi_site_origin, display_name, tr_response_headers, initial_dispersion, dns_bypass_cname, tr_request_headers, regional_geo_blocking, geo_provider, geo_limit_countries, logs_enabled, multi_site_origin_algorithm, geolimit_redirect_url, tenant_id, routing_name) VALUES (2100, 'test-ds1-root', true, 40, null, 0, 0, '', '', '', -1, 'http://test-ds1-root.edge', 21, 100, 100, 3600, 0, 0, 'test-ds1-root long_desc', 'test-ds1-root long_desc_1', 'test-ds1-root long_desc_2', 0, 'http://test-ds1-root.edge/info_url.html', 41.881944, -87.627778, '/crossdomain.xml', '2018-01-19 21:19:32.261210', 1, 0, true, 0, null, null, null, null, null, null, false, 'test-ds1-root-displayname', null, 1, null, null, true, 0, null, true, null, null, 1000000000, 'foo');
INSERT INTO deliveryservice (id, xml_id, active, dscp, signing_algorithm, qstring_ignore, geo_limit, http_bypass_fqdn, dns_bypass_ip, dns_bypass_ip6, dns_bypass_ttl, org_server_fqdn, type, profile, cdn_id, ccr_dns_ttl, global_max_mbps, global_max_tps, long_desc, long_desc_1, long_desc_2, max_dns_answers, info_url, miss_lat, miss_long, check_path, last_updated, protocol, ssl_key_version, ipv6_routing_enabled, range_request_handling, edge_header_rewrite, origin_shield, mid_header_rewrite, regex_remap, cacheurl, remap_text, multi_site_origin, display_name, tr_response_headers, initial_dispersion, dns_bypass_cname, tr_request_headers, regional_geo_blocking, geo_provider, geo_limit_countries, logs_enabled, multi_site_origin_algorithm, geolimit_redirect_url, tenant_id, routing_name) VALUES (2200, 'xxfoo.bar', true, 40, null, 0, 0, '', '', '', -1, 'http://foo.bar.edge', 21, 100, 100, 3600, 0, 0, 'foo.bar long_desc', 'foo.bar long_desc_1', 'foo.bar long_desc_2', 0, 'http://foo.bar.edge/info_url.html', 41.881944, -87.627778, '/crossdomain.xml', '2018-01-19 21:19:32.265717', 1, 0, true, 0, null, null, null, null, null, null, false, 'foo.bar-displayname', null, 1, null, null, true, 0, null, true, null, null, null, 'foo');
INSERT INTO deliveryservice (id, xml_id, active, dscp, signing_algorithm, qstring_ignore, geo_limit, http_bypass_fqdn, dns_bypass_ip, dns_bypass_ip6, dns_bypass_ttl, org_server_fqdn, type, profile, cdn_id, ccr_dns_ttl, global_max_mbps, global_max_tps, long_desc, long_desc_1, long_desc_2, max_dns_answers, info_url, miss_lat, miss_long, check_path, last_updated, protocol, ssl_key_version, ipv6_routing_enabled, range_request_handling, edge_header_rewrite, origin_shield, mid_header_rewrite, regex_remap, cacheurl, remap_text, multi_site_origin, display_name, tr_response_headers, initial_dispersion, dns_bypass_cname, tr_request_headers, regional_geo_blocking, geo_provider, geo_limit_countries, logs_enabled, multi_site_origin_algorithm, geolimit_redirect_url, tenant_id, routing_name) VALUES (300, 'test-ds3', true, 40, null, 0, 0, '', '', '', -1, 'http://test-ds3.edge', 9, 100, 100, 3600, 0, 0, 'test-ds3 long_desc', 'test-ds3 long_desc_1', 'test-ds3 long_desc_2', 0, 'http://test-ds3.edge/info_url.html', 41.881944, -87.627778, '/crossdomain.xml', '2018-01-19 21:19:32.269358', 0, 0, false, 0, null, null, null, null, null, null, false, 'test-ds3-displayname', null, 1, null, null, false, 0, null, false, null, null, null, 'foo');
INSERT INTO deliveryservice (id, xml_id, active, dscp, signing_algorithm, qstring_ignore, geo_limit, http_bypass_fqdn, dns_bypass_ip, dns_bypass_ip6, dns_bypass_ttl, org_server_fqdn, type, profile, cdn_id, ccr_dns_ttl, global_max_mbps, global_max_tps, long_desc, long_desc_1, long_desc_2, max_dns_answers, info_url, miss_lat, miss_long, check_path, last_updated, protocol, ssl_key_version, ipv6_routing_enabled, range_request_handling, edge_header_rewrite, origin_shield, mid_header_rewrite, regex_remap, cacheurl, remap_text, multi_site_origin, display_name, tr_response_headers, initial_dispersion, dns_bypass_cname, tr_request_headers, regional_geo_blocking, geo_provider, geo_limit_countries, logs_enabled, multi_site_origin_algorithm, geolimit_redirect_url, tenant_id, routing_name) VALUES (400, 'test-ds4', true, 40, null, 0, 0, '', '', '', -1, 'http://test-ds4.edge', 9, 100, 100, 3600, 0, 0, 'test-ds4 long_desc', 'test-ds4 long_desc_1', 'test-ds4 long_desc_2', 0, 'http://test-ds4.edge/info_url.html', 41.881944, -87.627778, '/crossdomain.xml', '2018-01-19 21:19:32.272467', 0, 0, false, 0, null, null, null, null, null, null, false, 'test-ds4-displayname', null, 1, null, null, false, 0, null, true, null, null, null, 'foo');
INSERT INTO deliveryservice (id, xml_id, active, dscp, signing_algorithm, qstring_ignore, geo_limit, http_bypass_fqdn, dns_bypass_ip, dns_bypass_ip6, dns_bypass_ttl, org_server_fqdn, type, profile, cdn_id, ccr_dns_ttl, global_max_mbps, global_max_tps, long_desc, long_desc_1, long_desc_2, max_dns_answers, info_url, miss_lat, miss_long, check_path, last_updated, protocol, ssl_key_version, ipv6_routing_enabled, range_request_handling, edge_header_rewrite, origin_shield, mid_header_rewrite, regex_remap, cacheurl, remap_text, multi_site_origin, display_name, tr_response_headers, initial_dispersion, dns_bypass_cname, tr_request_headers, regional_geo_blocking, geo_provider, geo_limit_countries, logs_enabled, multi_site_origin_algorithm, geolimit_redirect_url, tenant_id, routing_name) VALUES (500, 'test-ds5', true, 40, null, 0, 0, '', '', '', -1, 'http://test-ds5.edge', 7, 300, 100, 3600, 0, 0, 'test-ds5 long_desc', 'test-ds5 long_desc_1', 'test-ds5 long_desc_2', 0, 'http://test-ds5.edge/info_url.html', 41.881944, -87.627778, '/crossdomain.xml', '2018-01-19 21:19:32.275400', 0, 0, false, 0, null, null, null, null, null, null, false, 'test-ds5-displayname', null, 1, null, null, false, 0, null, false, null, null, null, 'foo');
INSERT INTO deliveryservice (id, xml_id, active, dscp, signing_algorithm, qstring_ignore, geo_limit, http_bypass_fqdn, dns_bypass_ip, dns_bypass_ip6, dns_bypass_ttl, org_server_fqdn, type, profile, cdn_id, ccr_dns_ttl, global_max_mbps, global_max_tps, long_desc, long_desc_1, long_desc_2, max_dns_answers, info_url, miss_lat, miss_long, check_path, last_updated, protocol, ssl_key_version, ipv6_routing_enabled, range_request_handling, edge_header_rewrite, origin_shield, mid_header_rewrite, regex_remap, cacheurl, remap_text, multi_site_origin, display_name, tr_response_headers, initial_dispersion, dns_bypass_cname, tr_request_headers, regional_geo_blocking, geo_provider, geo_limit_countries, logs_enabled, multi_site_origin_algorithm, geolimit_redirect_url, tenant_id, routing_name) VALUES (600, 'test-ds6', true, 40, null, 0, 0, '', '', '', -1, 'http://test-ds6.edge', 9, 300, 100, 3600, 0, 0, 'test-ds6 long_desc', 'test-ds6 long_desc_1', 'test-ds6 long_desc_2', 0, 'http://test-ds6.edge/info_url.html', 41.881944, -87.627778, '/crossdomain.xml', '2018-01-19 21:19:32.278451', 0, 0, false, 0, null, null, null, null, null, null, false, 'test-ds6-displayname', null, 1, null, null, false, 0, null, false, null, null, null, 'foo');
INSERT INTO deliveryservice (id, xml_id, active, dscp, signing_algorithm, qstring_ignore, geo_limit, http_bypass_fqdn, dns_bypass_ip, dns_bypass_ip6, dns_bypass_ttl, org_server_fqdn, type, profile, cdn_id, ccr_dns_ttl, global_max_mbps, global_max_tps, long_desc, long_desc_1, long_desc_2, max_dns_answers, info_url, miss_lat, miss_long, check_path, last_updated, protocol, ssl_key_version, ipv6_routing_enabled, range_request_handling, edge_header_rewrite, origin_shield, mid_header_rewrite, regex_remap, cacheurl, remap_text, multi_site_origin, display_name, tr_response_headers, initial_dispersion, dns_bypass_cname, tr_request_headers, regional_geo_blocking, geo_provider, geo_limit_countries, logs_enabled, multi_site_origin_algorithm, geolimit_redirect_url, tenant_id, routing_name) VALUES (700, 'steering-ds1', true, 40, null, 0, 0, '', 'hokeypokey', null, 10, 'http://steering-ds1.edge', 21, 100, 100, 3600, 0, 0, 'steering-ds1 long_desc', 'steering-ds1 long_desc_1', 'steering-ds1 long_desc_2', 0, 'http://steering-ds1.edge/info_url.html', 41.881944, -87.627778, '/crossdomain.xml', '2018-01-19 21:19:32.281466', 1, 0, true, 0, null, null, null, null, null, null, false, 'steering-ds1-displayname', null, 1, null, null, true, 0, null, false, null, null, null, 'foo');
INSERT INTO deliveryservice (id, xml_id, active, dscp, signing_algorithm, qstring_ignore, geo_limit, http_bypass_fqdn, dns_bypass_ip, dns_bypass_ip6, dns_bypass_ttl, org_server_fqdn, type, profile, cdn_id, ccr_dns_ttl, global_max_mbps, global_max_tps, long_desc, long_desc_1, long_desc_2, max_dns_answers, info_url, miss_lat, miss_long, check_path, last_updated, protocol, ssl_key_version, ipv6_routing_enabled, range_request_handling, edge_header_rewrite, origin_shield, mid_header_rewrite, regex_remap, cacheurl, remap_text, multi_site_origin, display_name, tr_response_headers, initial_dispersion, dns_bypass_cname, tr_request_headers, regional_geo_blocking, geo_provider, geo_limit_countries, logs_enabled, multi_site_origin_algorithm, geolimit_redirect_url, tenant_id, routing_name) VALUES (800, 'steering-ds2', true, 40, null, 0, 0, '', 'hokeypokey', null, 10, 'http://steering-ds2.edge', 21, 100, 100, 3600, 0, 0, 'steering-ds2 long_desc', 'steering-ds2 long_desc_1', 'steering-ds2 long_desc_2', 0, 'http://steering-ds2.edge/info_url.html', 41.881944, -87.627778, '/crossdomain.xml', '2018-01-19 21:19:32.284567', 1, 0, true, 0, null, null, null, null, null, null, false, 'steering-ds2-displayname', null, 1, null, null, true, 0, null, false, null, null, null, 'foo');
INSERT INTO deliveryservice (id, xml_id, active, dscp, signing_algorithm, qstring_ignore, geo_limit, http_bypass_fqdn, dns_bypass_ip, dns_bypass_ip6, dns_bypass_ttl, org_server_fqdn, type, profile, cdn_id, ccr_dns_ttl, global_max_mbps, global_max_tps, long_desc, long_desc_1, long_desc_2, max_dns_answers, info_url, miss_lat, miss_long, check_path, last_updated, protocol, ssl_key_version, ipv6_routing_enabled, range_request_handling, edge_header_rewrite, origin_shield, mid_header_rewrite, regex_remap, cacheurl, remap_text, multi_site_origin, display_name, tr_response_headers, initial_dispersion, dns_bypass_cname, tr_request_headers, regional_geo_blocking, geo_provider, geo_limit_countries, logs_enabled, multi_site_origin_algorithm, geolimit_redirect_url, tenant_id, routing_name) VALUES (900, 'steering-ds3', true, 40, null, 0, 0, '', 'hokeypokey', null, 10, 'http://new-steering-ds.edge', 21, 100, 100, 3600, 0, 0, 'new-steering-ds long_desc', 'new-steering-ds long_desc_1', 'new-steering-ds long_desc_2', 0, 'http://new-steering-ds.edge/info_url.html', 41.881944, -87.627778, '/crossdomain.xml', '2018-01-19 21:19:32.287726', 1, 0, true, 0, null, null, null, null, null, null, false, 'new-steering-ds-displayname', null, 1, null, null, true, 0, null, false, null, null, null, 'foo');
`
	err := execSQL(cfg, db, sqlStmt, "deliveryservice")
	if err != nil {
		return fmt.Errorf("exec failed %v", err)
	}
	return nil
}

// SetupRegexes ...
func SetupRegexes(cfg *config.Config, db *sql.DB) error {

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
	err := execSQL(cfg, db, sqlStmt, "regex")
	if err != nil {
		return fmt.Errorf("exec failed %v", err)
	}
	return nil
}

// SetupDeliveryServiceRegexes ...
func SetupDeliveryServiceRegexes(cfg *config.Config, db *sql.DB) error {

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
	err := execSQL(cfg, db, sqlStmt, "deliveryservice_regex")
	if err != nil {
		return fmt.Errorf("exec failed %v", err)
	}
	return nil
}

// SetupDeliveryServiceTmUsers ...
func SetupDeliveryServiceTmUsers(cfg *config.Config, db *sql.DB) error {

	sqlStmt := `
INSERT INTO deliveryservice_tmuser (deliveryservice, tm_user_id, last_updated) VALUES (100, (SELECT id FROM tm_user where username = 'admin') , '2018-01-19 21:19:32.372969');
`
	err := execSQL(cfg, db, sqlStmt, "deliveryservice_tmuser")
	if err != nil {
		return fmt.Errorf("exec failed %v", err)
	}
	return nil
}

// SetupDeliveryServiceServers ...
func SetupDeliveryServiceServers(cfg *config.Config, db *sql.DB) error {

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
	err := execSQL(cfg, db, sqlStmt, "deliveryservice_server")
	if err != nil {
		return fmt.Errorf("exec failed %v", err)
	}
	return nil
}

// SetupJobStatuses ...
func SetupJobStatuses(cfg *config.Config, db *sql.DB) error {

	sqlStmt := `
INSERT INTO job_status (id, name, description, last_updated) VALUES (1, 'PENDING', 'Job is queued, but has not been picked up by any agents yet', '2018-01-19 21:19:32.444857');
`
	err := execSQL(cfg, db, sqlStmt, "job_status")
	if err != nil {
		return fmt.Errorf("exec failed %v", err)
	}
	return nil
}

// SetupJobAgents ...
func SetupJobAgents(cfg *config.Config, db *sql.DB) error {

	sqlStmt := `
INSERT INTO job_agent (id, name, description, active, last_updated) VALUES (1, 'agent1', 'Test Agent1', 0, '2018-01-19 21:19:32.448076');
`
	err := execSQL(cfg, db, sqlStmt, "job_agent")
	if err != nil {
		return fmt.Errorf("exec failed %v", err)
	}
	return nil
}

// SetupJobs ...
func SetupJobs(cfg *config.Config, db *sql.DB) error {

	sqlStmt := `
INSERT INTO job (id, agent, object_type, object_name, keyword, parameters, asset_url, asset_type, status, start_time, entered_time, job_user, last_updated, job_deliveryservice) VALUES (100, 1, null, null, 'PURGE', 'TTL:48h', 'http://cdn2.edge/job1/.*', 'file', 1, '2018-01-19 21:01:14.000000', '2018-01-19 21:01:14.000000', (SELECT id FROM tm_user where username = 'admin'), '2018-01-19 21:19:32.468643', 100);
INSERT INTO job (id, agent, object_type, object_name, keyword, parameters, asset_url, asset_type, status, start_time, entered_time, job_user, last_updated, job_deliveryservice) VALUES (200, 1, null, null, 'PURGE', 'TTL:48h', 'http://cdn2.edge/job2/.*', 'file', 1, '2018-01-19 21:09:34.000000', '2018-01-19 21:09:34.000000', (SELECT id FROM tm_user where username = 'admin'), '2018-01-19 21:19:32.450915', 200);
INSERT INTO job (id, agent, object_type, object_name, keyword, parameters, asset_url, asset_type, status, start_time, entered_time, job_user, last_updated, job_deliveryservice) VALUES (300, 1, null, null, 'PURGE', 'TTL:48h', 'http://cdn2.edge/job3/.*', 'file', 1, '2018-01-19 21:14:34.000000', '2018-01-19 21:14:34.000000', (SELECT id FROM tm_user where username = 'admin'), '2018-01-19 21:19:32.460870', 100);
`
	err := execSQL(cfg, db, sqlStmt, "job")
	if err != nil {
		return fmt.Errorf("exec failed %v", err)
	}
	return nil
}

// Teardown - ensures that the data is cleaned up for a fresh run
func Teardown(cfg *config.Config, db *sql.DB) error {

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
	err := execSQL(cfg, db, sqlStmt, "Tearing down")
	if err != nil {
		return fmt.Errorf("exec failed %v", err)
	}
	return err
}

// execSQL ...
func execSQL(cfg *config.Config, db *sql.DB, sqlStmt string, dbTable string) error {

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
