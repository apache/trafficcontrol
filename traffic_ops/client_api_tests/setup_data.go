package client_tests

import (
	"database/sql"
	"fmt"

	log "github.com/apache/incubator-trafficcontrol/lib/go-log"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/auth"
)

func prepareDatabase(cfg *Config) {

	log.Debugln("Setting up Data")
	var db *sql.DB
	var err error

	sslStr := "require"
	if !cfg.DB.SSL {
		sslStr = "disable"
	}

	db, err = sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s", cfg.DB.User, cfg.DB.Password, cfg.DB.Hostname, cfg.DB.Name, sslStr))
	if err != nil {
		log.Errorf("opening database: %v\n", err)
		return
	}

	//res, err = tx.Exec("CREATE EXTENSION IF NOT EXISTS pgcrypto")
	//if err != nil {
	//log.Errorf("Transaction Failed %v %v", err, res)
	//}

	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		log.Errorf("Transaction Failed", err)
	}
	res, err := tx.Exec("DELETE FROM JOB")
	if err != nil {
		log.Errorf("Transaction Failed %v %v", err, res)
	}
	//db.MustExec("insert into role (name, description, priv_level) VALUES ($1, $2, $3, $4)", "disallowed", "Block all access", "0")
	// Role
	res, err = tx.Exec("DELETE FROM ROLE")
	if err != nil {
		log.Errorf("Transaction Failed %v %v", err, res)
	}

	res, err = tx.Exec("insert into role (id, name, description, priv_level) values (1, 'disallowed','Block all access',0) ON CONFLICT DO NOTHING")
	if err != nil {
		log.Errorf("Transaction Failed %v %v", err, res)
	}

	res, err = tx.Exec("insert into role (id, name, description, priv_level) values (2, 'read-only user','Block all access', 10) ON CONFLICT DO NOTHING")
	if err != nil {
		log.Errorf("Transaction Failed %v %v", err, res)
	}

	res, err = tx.Exec("insert into role (id, name, description, priv_level) values (3, 'operations','Block all access', 20) ON CONFLICT DO NOTHING")
	if err != nil {
		log.Errorf("Transaction Failed %v %v", err, res)
	}
	res, err = tx.Exec("insert into role (id, name, description, priv_level) values (4, 'admin','super-user', 30) ON CONFLICT DO NOTHING")
	if err != nil {
		log.Errorf("Transaction Failed %v %v", err, res)
	}
	res, err = tx.Exec("insert into role (id, name, description, priv_level) values (5, 'portal','Portal User', 2) ON CONFLICT DO NOTHING")
	if err != nil {
		log.Errorf("Transaction Failed %v %v", err, res)
	}
	res, err = tx.Exec("insert into role (id, name, description, priv_level) values (6, 'migrations','database migrations user - DO NOT REMOVE', 20) ON CONFLICT DO NOTHING")
	if err != nil {
		log.Errorf("Transaction Failed %v %v", err, res)
	}
	res, err = tx.Exec("insert into role (id, name, description, priv_level) values (7, 'federation','Role for Secondary CZF', 15) ON CONFLICT DO NOTHING")
	if err != nil {
		log.Errorf("Transaction Failed %v %v", err, res)
	}

	// Tmuser
	res, err = tx.Exec("DELETE FROM LOG")
	if err != nil {
		log.Errorf("Transaction Failed %v %v", err, res)
	}
	res, err = tx.Exec("DELETE FROM TM_USER")
	if err != nil {
		log.Errorf("Transaction Failed %v %v", err, res)
	}
	encryptedPassword, err := auth.DerivePassword(cfg.TOUserPassword)
	if err != nil {
		log.Errorf("Password encryption failed %v %v", err)
	}
	userInsert := `INSERT INTO tm_user (username, local_passwd, confirm_local_passwd, role) VALUES ('admin','` + encryptedPassword + `','` + encryptedPassword + `', 4)`
	res, err = tx.Exec(userInsert)
	if err != nil {
		log.Errorf("Transaction Failed %v %v", err, res)
	}

	tx.Commit()
	if err != nil {
		log.Infof(err.Error())
	}
}
