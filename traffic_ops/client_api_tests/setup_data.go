package client_tests

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"strconv"

	"golang.org/x/crypto/scrypt"

	log "github.com/apache/incubator-trafficcontrol/lib/go-log"
)

func prepareDatabaseTest(cfg *Config) {
	fmt.Printf("cfg ---> %v\n", cfg)
}
func prepareDatabase(cfg *Config) {

	fmt.Printf("cfg ---> %v\n", cfg)
	var db *sql.DB
	var err error

	sslStr := "require"
	if !cfg.DB.SSL {
		sslStr = "disable"
	}

	db, err = sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s", cfg.DB.User, cfg.DB.Password, cfg.DB.Hostname, cfg.DB.DBName, sslStr))
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
	encryptedPassword := EncryptPassword("password")
	fmt.Printf("encryptedPassword ---> %v\n", encryptedPassword)
	userInsert := `INSERT INTO tm_user (username, local_passwd, confirm_local_passwd, role) VALUES ('admin','` + encryptedPassword + `','` + encryptedPassword + `', 4)`
	fmt.Printf("userInsert ---> %v\n", userInsert)
	res, err = tx.Exec(userInsert)
	if err != nil {
		log.Errorf("Transaction Failed %v %v", err, res)
	}

	tx.Commit()
	if err != nil {
		log.Infof(err.Error())
	}
}

func EncryptPassword(password string) string {
	var salt []byte
	var err error
	salt, err = GenerateRandomBytes(64)
	if err != nil {
		log.Errorf(err.Error())
	}
	n := 16384
	r := 8
	p := 1
	keyLen := 64
	key, err := scrypt.Key([]byte(password), salt, n, r, p, keyLen)
	//key, err := scrypt.Key([]byte("laser1"), salt, 1<<15, 8, 1, 64)
	if err != nil {
		log.Errorf(err.Error())
	}
	nStr := strconv.Itoa(n)
	if err != nil {
		log.Errorf(err.Error())
	}
	rStr := strconv.Itoa(r)
	pStr := strconv.Itoa(p)
	saltBase64 := base64.StdEncoding.EncodeToString(salt)
	keyBase64 := base64.StdEncoding.EncodeToString(key)
	return "SCRYPT:" + nStr + ":" + rStr + ":" + pStr + ":" + saltBase64 + ":" + keyBase64
}

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}
