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
	if !cfg.DB.SSL {
		sslStr = "disable"
	}

	db, err = sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s", cfg.DB.User, cfg.DB.Password, cfg.DB.Hostname, cfg.DB.Name, sslStr))
	if err != nil {
		log.Errorf("opening database: %v\n", err)
		return nil, fmt.Errorf("Transaction Failed: %s", err)
	}
	return db, err
}

func setupData(cfg *Config, db *sql.DB) error {

	log.Debugln("Setting up Data")
	var err error

	//res, err = tx.Exec("CREATE EXTENSION IF NOT EXISTS pgcrypto")
	//if err != nil {
	//log.Errorf("Transaction Failed %v %v", err, res)
	//}

	tx, err := db.Begin()

	res, err := tx.Exec("INSERT INTO role (id, name, description, priv_level) VALUES (1, 'disallowed','Block all access',0) ON CONFLICT DO NOTHING")
	if err != nil {
		return fmt.Errorf("Transaction Failed %v %v", err, res)
	}

	res, err = tx.Exec("INSERT INTO role (id, name, description, priv_level) VALUES (2, 'read-only user','Block all access', 10) ON CONFLICT DO NOTHING")
	if err != nil {
		return fmt.Errorf("Transaction Failed %v %v", err, res)
	}

	res, err = tx.Exec("INSERT INTO role (id, name, description, priv_level) VALUES (3, 'operations','Block all access', 20) ON CONFLICT DO NOTHING")
	if err != nil {
		return fmt.Errorf("Transaction Failed %v %v", err, res)
	}
	res, err = tx.Exec("INSERT INTO role (id, name, description, priv_level) VALUES (4, 'admin','super-user', 30) ON CONFLICT DO NOTHING")
	if err != nil {
		return fmt.Errorf("Transaction Failed %v %v", err, res)
	}
	res, err = tx.Exec("INSERT INTO role (id, name, description, priv_level) VALUES (5, 'portal','Portal User', 2) ON CONFLICT DO NOTHING")
	if err != nil {
		return fmt.Errorf("Transaction Failed %v %v", err, res)
	}
	res, err = tx.Exec("INSERT INTO role (id, name, description, priv_level) VALUES (6, 'migrations','database migrations user - DO NOT REMOVE', 20) ON CONFLICT DO NOTHING")
	if err != nil {
		return fmt.Errorf("Transaction Failed %v %v", err, res)
	}
	res, err = tx.Exec("INSERT INTO role (id, name, description, priv_level) VALUES (7, 'federation','Role for Secondary CZF', 15) ON CONFLICT DO NOTHING")
	if err != nil {
		return fmt.Errorf("Transaction Failed %v %v", err, res)
	}

	encryptedPassword, err := auth.DerivePassword(cfg.TOUserPassword)
	if err != nil {
		return fmt.Errorf("Password encryption failed %v %v", err)
	}
	userInsert := `INSERT INTO tm_user (username, local_passwd, confirm_local_passwd, role) VALUES ('admin','` + encryptedPassword + `','` + encryptedPassword + `', 4)`
	res, err = tx.Exec(userInsert)
	if err != nil {
		return fmt.Errorf("Transaction Failed %v %v", err, res)
	}

	tx.Commit()
	if err != nil {
		return fmt.Errorf("Commit Failed %v %v", err, res)
	}
	return nil
}

// ensures that the data is cleaned up for a fresh run
func teardownData(cfg *Config, db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("Transaction Failed %v %v ", err, tx)
	}

	res, err := tx.Exec("DELETE FROM TO_EXTENSION")
	if err != nil {
		return fmt.Errorf("Transaction Failed %v %v", err, res)
	}

	res, err = tx.Exec("DELETE FROM STATICDNSENTRY")
	if err != nil {
		return fmt.Errorf("Transaction Failed %v %v", err, res)
	}
	res, err = tx.Exec("DELETE FROM JOB")
	if err != nil {
		return fmt.Errorf("Transaction Failed %v %v", err, res)
	}
	res, err = tx.Exec("DELETE FROM JOB_AGENT")
	if err != nil {
		return fmt.Errorf("Transaction Failed %v %v", err, res)
	}
	res, err = tx.Exec("DELETE FROM JOB_STATUS")
	if err != nil {
		return fmt.Errorf("Transaction Failed %v %v", err, res)
	}
	res, err = tx.Exec("DELETE FROM LOG")
	if err != nil {
		return fmt.Errorf("Transaction Failed %v %v", err, res)
	}
	res, err = tx.Exec("DELETE FROM ASN")
	if err != nil {
		return fmt.Errorf("Transaction Failed %v %v", err, res)
	}
	res, err = tx.Exec("DELETE FROM DELIVERYSERVICE_TMUSER")
	if err != nil {
		return fmt.Errorf("Transaction Failed %v %v", err, res)
	}
	res, err = tx.Exec("DELETE FROM TM_USER")
	if err != nil {
		return fmt.Errorf("Transaction Failed %v %v", err, res)
	}
	res, err = tx.Exec("DELETE FROM ROLE")
	if err != nil {
		return fmt.Errorf("Transaction Failed %v %v", err, res)
	}
	res, err = tx.Exec("DELETE FROM DELIVERYSERVICE_REGEX")
	if err != nil {
		return fmt.Errorf("Transaction Failed %v %v", err, res)
	}
	res, err = tx.Exec("DELETE FROM REGEX")
	if err != nil {
		return fmt.Errorf("Transaction Failed %v %v", err, res)
	}
	res, err = tx.Exec("DELETE FROM DELIVERYSERVICE_SERVER")
	if err != nil {
		return fmt.Errorf("Transaction Failed %v %v", err, res)
	}
	res, err = tx.Exec("DELETE FROM DELIVERYSERVICE")
	if err != nil {
		return fmt.Errorf("Transaction Failed %v %v", err, res)
	}
	res, err = tx.Exec("DELETE FROM SERVER")
	if err != nil {
		return fmt.Errorf("Transaction Failed %v %v", err, res)
	}
	res, err = tx.Exec("DELETE FROM PHYS_LOCATION")
	if err != nil {
		return fmt.Errorf("Transaction Failed %v %v", err, res)
	}
	res, err = tx.Exec("DELETE FROM REGION")
	if err != nil {
		return fmt.Errorf("Transaction Failed %v %v", err, res)
	}
	res, err = tx.Exec("DELETE FROM DIVISION")
	if err != nil {
		return fmt.Errorf("Transaction Failed %v %v", err, res)
	}
	res, err = tx.Exec("DELETE FROM PROFILE")
	if err != nil {
		return fmt.Errorf("Transaction Failed %v %v", err, res)
	}
	res, err = tx.Exec("DELETE FROM PARAMETER")
	if err != nil {
		return fmt.Errorf("Transaction Failed %v %v", err, res)
	}
	res, err = tx.Exec("DELETE FROM PROFILE_PARAMETER")
	if err != nil {
		return fmt.Errorf("Transaction Failed %v %v", err, res)
	}
	res, err = tx.Exec("DELETE FROM CACHEGROUP")
	if err != nil {
		return fmt.Errorf("Transaction Failed %v %v", err, res)
	}
	res, err = tx.Exec("DELETE FROM TYPE")
	if err != nil {
		return fmt.Errorf("Transaction Failed %v %v", err, res)
	}
	res, err = tx.Exec("DELETE FROM STATUS")
	if err != nil {
		return fmt.Errorf("Transaction Failed %v %v", err, res)
	}
	res, err = tx.Exec("DELETE FROM SNAPSHOT")
	if err != nil {
		return fmt.Errorf("Transaction Failed %v %v", err, res)
	}
	res, err = tx.Exec("DELETE FROM CDN")
	if err != nil {
		return fmt.Errorf("Transaction Failed %v %v", err, res)
	}
	res, err = tx.Exec("DELETE FROM TENANT")
	if err != nil {
		return fmt.Errorf("Transaction Failed %v %v", err, res)
	}

	tx.Commit()
	if err != nil {
		return fmt.Errorf("Commit Failed %v %v", err, res)
	}
	return err
}
