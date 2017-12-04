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

package api

import (
	"database/sql"
	"fmt"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/auth"
)

var (
	db *sql.DB
)

func openConnection(cfg *Config) (*sql.DB, error) {
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

func setupUserData(cfg *Config, db *sql.DB) error {

	log.Debugln("Setting up initial user data")
	var err error

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("transaction begin failed %v %v ", err, tx)
	}

	res, err := tx.Exec("INSERT INTO role (id, name, description, priv_level) VALUES (1, 'disallowed','Block all access',0) ON CONFLICT DO NOTHING")
	if err != nil {
		return fmt.Errorf("exec failed %v %v", err, res)
	}

	res, err = tx.Exec("INSERT INTO role (id, name, description, priv_level) VALUES (2, 'read-only user','Block all access', 10) ON CONFLICT DO NOTHING")
	if err != nil {
		return fmt.Errorf("exec failed %v %v", err, res)
	}

	res, err = tx.Exec("INSERT INTO role (id, name, description, priv_level) VALUES (3, 'operations','Block all access', 20) ON CONFLICT DO NOTHING")
	if err != nil {
		return fmt.Errorf("exec failed %v %v", err, res)
	}
	res, err = tx.Exec("INSERT INTO role (id, name, description, priv_level) VALUES (4, 'admin','super-user', 30) ON CONFLICT DO NOTHING")
	if err != nil {
		return fmt.Errorf("exec failed %v %v", err, res)
	}
	res, err = tx.Exec("INSERT INTO role (id, name, description, priv_level) VALUES (5, 'portal','Portal User', 2) ON CONFLICT DO NOTHING")
	if err != nil {
		return fmt.Errorf("exec failed %v %v", err, res)
	}
	res, err = tx.Exec("INSERT INTO role (id, name, description, priv_level) VALUES (6, 'migrations','database migrations user - DO NOT REMOVE', 20) ON CONFLICT DO NOTHING")
	if err != nil {
		return fmt.Errorf("exec failed %v %v", err, res)
	}
	res, err = tx.Exec("INSERT INTO role (id, name, description, priv_level) VALUES (7, 'federation','Role for Secondary CZF', 15) ON CONFLICT DO NOTHING")
	if err != nil {
		return fmt.Errorf("exec failed %v %v", err, res)
	}

	encryptedPassword, err := auth.DerivePassword(cfg.TrafficOps.UserPassword)
	if err != nil {
		return fmt.Errorf("password encryption failed %v", err)
	}
	itm := `INSERT INTO tm_user (username, local_passwd, confirm_local_passwd, role) VALUES ('admin','` + encryptedPassword + `','` + encryptedPassword + `', 4)`
	res, err = tx.Exec(itm)
	if err != nil {
		return fmt.Errorf("exec failed %v %v", err, res)
	}

	tx.Commit()
	if err != nil {
		return fmt.Errorf("commit failed %v %v", err, res)
	}
	return nil
}

// ensures that the data is cleaned up for a fresh run
func teardownData(cfg *Config, db *sql.DB) error {
	log.Debugln("Tearing down data")
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("transaction begin failed %v %v ", err, tx)
	}

	res, err := tx.Exec("DELETE FROM TO_EXTENSION")
	if err != nil {
		return fmt.Errorf("exec failed %v %v", err, res)
	}

	res, err = tx.Exec("DELETE FROM STATICDNSENTRY")
	if err != nil {
		return fmt.Errorf("exec failed %v %v", err, res)
	}
	res, err = tx.Exec("DELETE FROM JOB")
	if err != nil {
		return fmt.Errorf("exec failed %v %v", err, res)
	}
	res, err = tx.Exec("DELETE FROM JOB_AGENT")
	if err != nil {
		return fmt.Errorf("exec failed %v %v", err, res)
	}
	res, err = tx.Exec("DELETE FROM JOB_STATUS")
	if err != nil {
		return fmt.Errorf("exec failed %v %v", err, res)
	}
	res, err = tx.Exec("DELETE FROM LOG")
	if err != nil {
		return fmt.Errorf("exec failed %v %v", err, res)
	}
	res, err = tx.Exec("DELETE FROM ASN")
	if err != nil {
		return fmt.Errorf("exec failed %v %v", err, res)
	}
	res, err = tx.Exec("DELETE FROM DELIVERYSERVICE_TMUSER")
	if err != nil {
		return fmt.Errorf("exec failed %v %v", err, res)
	}
	res, err = tx.Exec("DELETE FROM TM_USER")
	if err != nil {
		return fmt.Errorf("exec failed %v %v", err, res)
	}
	res, err = tx.Exec("DELETE FROM ROLE")
	if err != nil {
		return fmt.Errorf("exec failed %v %v", err, res)
	}
	res, err = tx.Exec("DELETE FROM DELIVERYSERVICE_REGEX")
	if err != nil {
		return fmt.Errorf("exec failed %v %v", err, res)
	}
	res, err = tx.Exec("DELETE FROM REGEX")
	if err != nil {
		return fmt.Errorf("exec failed %v %v", err, res)
	}
	res, err = tx.Exec("DELETE FROM DELIVERYSERVICE_SERVER")
	if err != nil {
		return fmt.Errorf("exec failed %v %v", err, res)
	}
	res, err = tx.Exec("DELETE FROM DELIVERYSERVICE")
	if err != nil {
		return fmt.Errorf("exec failed %v %v", err, res)
	}
	res, err = tx.Exec("DELETE FROM SERVER")
	if err != nil {
		return fmt.Errorf("exec failed %v %v", err, res)
	}
	res, err = tx.Exec("DELETE FROM PHYS_LOCATION")
	if err != nil {
		return fmt.Errorf("exec failed %v %v", err, res)
	}
	res, err = tx.Exec("DELETE FROM REGION")
	if err != nil {
		return fmt.Errorf("exec failed %v %v", err, res)
	}
	res, err = tx.Exec("DELETE FROM DIVISION")
	if err != nil {
		return fmt.Errorf("exec failed %v %v", err, res)
	}
	res, err = tx.Exec("DELETE FROM PROFILE")
	if err != nil {
		return fmt.Errorf("exec failed %v %v", err, res)
	}
	res, err = tx.Exec("DELETE FROM PARAMETER")
	if err != nil {
		return fmt.Errorf("exec failed %v %v", err, res)
	}
	res, err = tx.Exec("DELETE FROM PROFILE_PARAMETER")
	if err != nil {
		return fmt.Errorf("exec failed %v %v", err, res)
	}
	res, err = tx.Exec("DELETE FROM CACHEGROUP")
	if err != nil {
		return fmt.Errorf("exec failed %v %v", err, res)
	}
	res, err = tx.Exec("DELETE FROM TYPE")
	if err != nil {
		return fmt.Errorf("exec failed %v %v", err, res)
	}
	res, err = tx.Exec("DELETE FROM STATUS")
	if err != nil {
		return fmt.Errorf("exec failed %v %v", err, res)
	}
	res, err = tx.Exec("DELETE FROM SNAPSHOT")
	if err != nil {
		return fmt.Errorf("exec failed %v %v", err, res)
	}
	res, err = tx.Exec("DELETE FROM CDN")
	if err != nil {
		return fmt.Errorf("exec failed %v %v", err, res)
	}
	res, err = tx.Exec("DELETE FROM TENANT")
	if err != nil {
		return fmt.Errorf("exec failed %v %v", err, res)
	}

	tx.Commit()
	if err != nil {
		return fmt.Errorf("commit failed %v %v", err, res)
	}
	return err
}
