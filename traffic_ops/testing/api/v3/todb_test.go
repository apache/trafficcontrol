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

package v3

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/auth"
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
// TODO error does not need returned as this function can never return a non-nil error
func SetupTestData(*sql.DB) error {
	var err error

	err = SetupRoles(db)
	if err != nil {
		fmt.Printf("\nError setting up roles %s - %s, %v\n", Config.TrafficOps.URL, Config.TrafficOps.Users.Admin, err)
		os.Exit(1)
	}

	err = SetupCapabilities(db)
	if err != nil {
		fmt.Printf("\nError setting up capabilities %s - %s, %v\n", Config.TrafficOps.URL, Config.TrafficOps.Users.Admin, err)
		os.Exit(1)
	}

	err = SetupRoleCapabilities(db)
	if err != nil {
		fmt.Printf("\nError setting up roleCapabilities %s - %s, %v\n", Config.TrafficOps.URL, Config.TrafficOps.Users.Admin, err)
		os.Exit(1)
	}

	err = SetupAPICapabilities(db)
	if err != nil {
		fmt.Printf("\nError setting up APICapabilities %s - %s, %v\n", Config.TrafficOps.URL, Config.TrafficOps.Users.Admin, err)
		os.Exit(1)
	}

	err = SetupTenants(db)
	if err != nil {
		fmt.Printf("\nError setting up tenant %s - %s, %v\n", Config.TrafficOps.URL, Config.TrafficOps.Users.Admin, err)
		os.Exit(1)
	}

	err = SetupTmusers(db)
	if err != nil {
		fmt.Printf("\nError setting up tm_user %s - %s, %v\n", Config.TrafficOps.URL, Config.TrafficOps.Users.Admin, err)
		os.Exit(1)
	}

	err = SetupJobAgents(db)
	if err != nil {
		fmt.Printf("\nError setting up job agents %s - %s, %v\n", Config.TrafficOps.URL, Config.TrafficOps.Users.Admin, err)
		os.Exit(1)
	}

	err = SetupJobStatuses(db)
	if err != nil {
		fmt.Printf("\nError setting up job agents %s - %s, %v\n", Config.TrafficOps.URL, Config.TrafficOps.Users.Admin, err)
		os.Exit(1)
	}

	err = SetupTypes(db)
	if err != nil {
		fmt.Printf("\nError setting up types %s - %s, %v\n", Config.TrafficOps.URL, Config.TrafficOps.Users.Admin, err)
		os.Exit(1)
	}

	err = SetupToExtensions(db)
	if err != nil {
		fmt.Printf("\nError setting up to_extensions %s - %s, %v\n", Config.TrafficOps.URL, Config.TrafficOps.Users.Admin, err)
		os.Exit(1)
	}

	return err
}

// SetupRoles ...
func SetupRoles(db *sql.DB) error {

	sqlStmt := `
INSERT INTO role (name, description, priv_level) VALUES ('disallowed','Block all access',0) ON CONFLICT DO NOTHING;
INSERT INTO role (name, description, priv_level) VALUES ('read-only user','Block all access', 10) ON CONFLICT DO NOTHING;
INSERT INTO role (name, description, priv_level) VALUES ('operations','Block all access', 20) ON CONFLICT DO NOTHING;
INSERT INTO role (name, description, priv_level) VALUES ('admin','super-user', 30) ON CONFLICT DO NOTHING;
INSERT INTO role (name, description, priv_level) VALUES ('portal','Portal User', 2) ON CONFLICT DO NOTHING;
INSERT INTO role (name, description, priv_level) VALUES ('steering','Steering User', 15) ON CONFLICT DO NOTHING;
INSERT INTO role (name, description, priv_level) VALUES ('federation','Role for Secondary CZF', 15) ON CONFLICT DO NOTHING;
`
	err := execSQL(db, sqlStmt)
	if err != nil {
		return fmt.Errorf("exec failed %v", err)
	}
	return nil
}

func SetupCapabilities(db *sql.DB) error {
	sqlStmt := `
INSERT INTO capability (name, description) VALUES ('all-read','Full read access') ON CONFLICT DO NOTHING;
INSERT INTO capability (name, description) VALUES ('all-write','Full write access') ON CONFLICT DO NOTHING;
INSERT INTO capability (name, description) VALUES ('cdn-read','View CDN configuration') ON CONFLICT DO NOTHING;
INSERT INTO capability (name, description) VALUES ('asns-read', 'Read ASNs') ON CONFLICT DO NOTHING;
INSERT INTO capability (name, description) VALUES ('asns-write', 'Write ASNs') ON CONFLICT DO NOTHING;
INSERT INTO capability (name, description) VALUES ('cache-groups-read', 'Read CGs') ON CONFLICT DO NOTHING;
`
	err := execSQL(db, sqlStmt)
	if err != nil {
		return fmt.Errorf("exec failed %v", err)
	}
	return nil
}

func SetupAPICapabilities(db *sql.DB) error {
	sqlStmt := `
INSERT INTO api_capability (http_method, route, capability) VALUES ('GET', '/asns', 'asns-read') ON CONFLICT DO NOTHING;
INSERT INTO api_capability (http_method, route, capability) VALUES ('POST', '/asns', 'asns-write') ON CONFLICT DO NOTHING;
INSERT INTO api_capability (http_method, route, capability) VALUES ('GET', '/cachegroups', 'cache-groups-read') ON CONFLICT DO NOTHING;
`

	err := execSQL(db, sqlStmt)
	if err != nil {
		return fmt.Errorf("exec failed %v", err)
	}
	return nil
}

func SetupRoleCapabilities(db *sql.DB) error {
	sqlStmt := `
INSERT INTO role_capability (role_id, cap_name) VALUES (4,'all-write') ON CONFLICT DO NOTHING;
INSERT INTO role_capability (role_id, cap_name) VALUES (4,'all-read') ON CONFLICT DO NOTHING;
`
	err := execSQL(db, sqlStmt)
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
INSERT INTO tm_user (username, local_passwd, role, tenant_id) VALUES ('` + Config.TrafficOps.Users.Disallowed + `','` + encryptedPassword + `', 1, 1);
INSERT INTO tm_user (username, local_passwd, role, tenant_id) VALUES ('` + Config.TrafficOps.Users.ReadOnly + `','` + encryptedPassword + `', 2, 1);
INSERT INTO tm_user (username, local_passwd, role, tenant_id) VALUES ('` + Config.TrafficOps.Users.Operations + `','` + encryptedPassword + `', 3, 1);
INSERT INTO tm_user (username, local_passwd, role, tenant_id) VALUES ('` + Config.TrafficOps.Users.Admin + `','` + encryptedPassword + `', 4, 1);
INSERT INTO tm_user (username, local_passwd, role, tenant_id) VALUES ('` + Config.TrafficOps.Users.Portal + `','` + encryptedPassword + `', 5, 1);
INSERT INTO tm_user (username, local_passwd, role, tenant_id) VALUES ('` + Config.TrafficOps.Users.Federation + `','` + encryptedPassword + `', 6, 1);
INSERT INTO tm_user (username, local_passwd, role, tenant_id) VALUES ('` + Config.TrafficOps.Users.Extension + `','` + encryptedPassword + `', 3, 1);
`
	err = execSQL(db, sqlStmt)
	if err != nil {
		return fmt.Errorf("exec failed %v", err)
	}
	return nil
}

// SetupTenants ...
func SetupTenants(db *sql.DB) error {

	// TODO: root tenant must be present in initial database.  "badtenant" is needed for now so tests can be done
	// with a tenant outside the user's tenant.  That should be removed once User API tests are in place rather than the SetupUsers defined above.
	sqlStmt := `
INSERT INTO tenant (name, active, parent_id, last_updated) VALUES ('root', true, null, '2018-01-19 19:01:21.327262');
INSERT INTO tenant (name, active, parent_id, last_updated) VALUES ('badtenant', true, 1, '2018-01-19 19:01:21.327262');
`
	err := execSQL(db, sqlStmt)
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
	err := execSQL(db, sqlStmt)
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
	err := execSQL(db, sqlStmt)
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
	err := execSQL(db, sqlStmt)
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
	err := execSQL(db, sqlStmt)
	if err != nil {
		return fmt.Errorf("exec failed %v", err)
	}
	return nil
}

// SetupTypes Set up to_extension types
func SetupTypes(db *sql.DB) error {

	sqlStmt := `
INSERT INTO type (name, description, use_in_table) VALUES ('CHECK_EXTENSION_BOOL', 'Extension for checkmark in Server Check', 'to_extension');
INSERT INTO type (name, description, use_in_table) VALUES ('CHECK_EXTENSION_NUM', 'Extension for int value in Server Check', 'to_extension');
INSERT INTO type (name, description, use_in_table) VALUES ('CHECK_EXTENSION_OPEN_SLOT', 'Open slot for check in Server Status', 'to_extension');
`
	err := execSQL(db, sqlStmt)
	if err != nil {
		return fmt.Errorf("exec failed %v", err)
	}
	return nil
}

// SetupToExtensions setup open slot in to_extension table
func SetupToExtensions(db *sql.DB) error {

	sqlStmt := `
INSERT INTO to_extension (name, version, info_url, isactive, script_file, servercheck_column_name, type) VALUES ('OPEN', '1.0.0', '-', false, '', 'aa', (SELECT id FROM type WHERE name = 'CHECK_EXTENSION_OPEN_SLOT'));
INSERT INTO to_extension (name, version, info_url, isactive, script_file, servercheck_column_name, type) VALUES ('OPEN', '1.0.0', '-', false, '', 'ab', (SELECT id FROM type WHERE name = 'CHECK_EXTENSION_OPEN_SLOT'));
	`
	err := execSQL(db, sqlStmt)
	if err != nil {
		return fmt.Errorf("exec failed %v", err)
	}
	return nil
}

// Teardown - ensures that the data is cleaned up for a fresh run
func Teardown(db *sql.DB) error {

	sqlStmt := `
	DELETE FROM api_capability;
	DELETE FROM deliveryservices_required_capability;
	DELETE FROM server_server_capability;
	DELETE FROM server_capability;
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
	DELETE FROM capability;
	ALTER SEQUENCE role_id_seq RESTART WITH 1;
	DELETE FROM deliveryservice_regex;
	DELETE FROM regex;
	DELETE FROM deliveryservice_server;
	DELETE FROM deliveryservice;
	DELETE FROM origin;
	DELETE FROM ip_address;
	DELETE FROM interface;
	DELETE FROM server;
	DELETE FROM phys_location;
	DELETE FROM region;
	DELETE FROM division;
	DELETE FROM profile;
	DELETE FROM parameter;
	DELETE FROM profile_parameter;
	DELETE FROM topology_cachegroup_parents;
	DELETE FROM topology_cachegroup;
	DELETE FROM topology;
	DELETE FROM cachegroup;
	DELETE FROM coordinate;
	DELETE FROM type;
	DELETE FROM status;
	DELETE FROM snapshot;
	DELETE FROM cdn;
	DELETE FROM service_category;
	DELETE FROM tenant;
	ALTER SEQUENCE tenant_id_seq RESTART WITH 1;
`
	err := execSQL(db, sqlStmt)
	if err != nil {
		return fmt.Errorf("exec failed %v", err)
	}
	return err
}

// execSQL ...
func execSQL(db *sql.DB, sqlStmt string) error {
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
