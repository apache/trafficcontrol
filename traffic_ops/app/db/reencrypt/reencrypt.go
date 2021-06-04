/*

Name
	reencrypt

Synopsis
	reencrypt [--previous-key value] [--new-key value] [--cfg value]

Description
  The reencrypt app is used to re-encrypt all data in the Postgres Traffic Vault
  using a new base64-encoded AES key.

Options
	--previous-key
        The file path for the previous base64-encoded AES key.

	--new-key
        The file path for the new base64-encoded AES key.

	--cfg
        The path for the configuration file. Default is `./reencrypt.conf`.

*/

package main

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
	"crypto/aes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/trafficvault/backends/postgres"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const PROPERTIES_FILE = "./reencrypt.conf"

func main() {
	previousKeyLocation := flag.String("previous-key", "", "The file path for the previous base64 encoded AES key.")
	newKeyLocation := flag.String("new-key", "", "The file path for the new base64 encoded AES key.")
	cfg := flag.String("cfg", PROPERTIES_FILE, "The path for the configuration file. Default is "+PROPERTIES_FILE+".")
	flag.Parse()

	if len(os.Args) < 4 {
		flag.Usage()
		os.Exit(1)
	}

	if previousKeyLocation == nil || *previousKeyLocation == "" {
		die("previous-key flag is required.")
	}
	if newKeyLocation == nil || *newKeyLocation == "" {
		die("new-key flag is required.")
	}

	newKey, err := readKey(*newKeyLocation)
	if err != nil {
		die("reading new-key: " + err.Error())
	}

	previousKey, err := readKey(*previousKeyLocation)
	if err != nil {
		die("reading previous-key: " + err.Error())
	}

	dbConfBytes, err := ioutil.ReadFile(*cfg)
	if err != nil {
		die("reading db conf '" + *cfg + "': " + err.Error())
	}

	pgCfg := Config{}
	err = json.Unmarshal(dbConfBytes, &pgCfg)
	if err != nil {
		die("unmarshalling '" + *cfg + "': " + err.Error())
	}

	sslStr := "require"
	if !pgCfg.SSL {
		sslStr = "disable"
	}
	db, err := sqlx.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s&fallback_application_name=trafficvault", pgCfg.User, pgCfg.Password, pgCfg.Hostname, pgCfg.Port, pgCfg.DBName, sslStr))
	if err != nil {
		die("opening database: " + err.Error())
	}

	if err = reEncryptSslKeys(db, previousKey, newKey); err != nil {
		die("re-encrypting SSL Keys: " + err.Error())
	}
	if err = reEncryptUrlSigKeys(db, previousKey, newKey); err != nil {
		die("re-encrypting URL Sig Keys: " + err.Error())
	}
	if err = reEncryptUriSigningKeys(db, previousKey, newKey); err != nil {
		die("re-encrypting URI Signing Keys: " + err.Error())
	}
	if err = reEncryptDNSSECKeys(db, previousKey, newKey); err != nil {
		die("re-encrypting DNSSEC Keys: " + err.Error())
	}

	if err = updateKeyFile(previousKey, *previousKeyLocation, newKey); err != nil {
		die("updating the key file: " + err.Error())
	}

	fmt.Println("Successfully re-encrypted keys for SSL Keys, URL Sig Keys, URI Signing Keys, and DNSSEC Keys.")
}

type Config struct {
	DBName   string `json:"dbname"`
	Hostname string `json:"hostname"`
	User     string `json:"user"`
	Password string `json:"password"`
	Port     int    `json:"port"`
	SSL      bool   `json:"ssl"`
}

func readKey(keyLocation string) ([]byte, error) {
	var keyBase64 string
	keyBase64Bytes, err := ioutil.ReadFile(keyLocation)
	if err != nil {
		return []byte{}, errors.New("reading file '" + keyLocation + "':" + err.Error())
	}
	keyBase64 = string(keyBase64Bytes)

	key, err := base64.StdEncoding.DecodeString(keyBase64)
	if err != nil {
		return []byte{}, errors.New(fmt.Sprintf("AES key cannot be decoded from base64: %v", err))
	}

	// verify the key works
	_, err = aes.NewCipher(key)
	if err != nil {
		return []byte{}, err
	}

	return key, nil
}

func reEncryptSslKeys(db *sqlx.DB, previousKey []byte, newKey []byte) error {
	tx, err := db.Begin()
	if err != nil {
		return errors.New(fmt.Sprintf("transaction begin failed %v %v ", err, tx))
	}
	defer tx.Commit()

	rows, err := tx.Query("SELECT id, data FROM sslkey")
	if err != nil {
		return errors.New(fmt.Sprintf("querying: %v", err))
	}
	defer rows.Close()

	for rows.Next() {
		updateTx, err := db.Begin()
		if err != nil {
			return errors.New(fmt.Sprintf("transaction begin failed %v %v", err, updateTx))
		}
		defer updateTx.Commit()

		id := 0
		var encryptedSslKeys []byte
		if err = rows.Scan(&id, &encryptedSslKeys); err != nil {
			return errors.New(fmt.Sprintf("getting SSL Keys: %v", err))
		}
		jsonKeys, err := postgres.AESDecrypt(encryptedSslKeys, previousKey)
		if err != nil {
			return errors.New(fmt.Sprintf("reading SSL Keys: %v", err))
		}

		reencryptedKeys, err := postgres.AESEncrypt(jsonKeys, newKey)
		if err != nil {
			return errors.New(fmt.Sprintf("encrypting SSL Keys with new key: %v", err))
		}

		res, err := updateTx.Exec(`UPDATE sslkey SET data = $1 WHERE id = $2`, []byte(reencryptedKeys), id)
		if err != nil {
			return errors.New(fmt.Sprintf("updating SSL Keys for id %d: %v", id, err))
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return errors.New(fmt.Sprintf("determining rows affected for reencrypting SSL Keys with id %d", id))
		}
		if rowsAffected == 0 {
			return errors.New(fmt.Sprintf("no rows updated for reencrypting SSL Keys for id %d", id))
		}

	}

	return nil
}

func reEncryptUrlSigKeys(db *sqlx.DB, previousKey []byte, newKey []byte) error {
	tx, err := db.Begin()
	if err != nil {
		return errors.New(fmt.Sprintf("transaction begin failed %v %v", err, tx))
	}
	defer tx.Commit()

	rows, err := tx.Query("SELECT deliveryservice, data FROM url_sig_key")
	if err != nil {
		return errors.New(fmt.Sprintf("querying: %v", err))
	}
	defer rows.Close()

	for rows.Next() {
		updateTx, err := db.Begin()
		if err != nil {
			return errors.New(fmt.Sprintf("transaction begin failed %v %v", err, updateTx))
		}
		defer updateTx.Commit()

		ds := ""
		var encryptedSslKeys []byte
		if err = rows.Scan(&ds, &encryptedSslKeys); err != nil {
			return errors.New(fmt.Sprintf("getting URL Sig Keys: %v", err))
		}
		jsonKeys, err := postgres.AESDecrypt(encryptedSslKeys, previousKey)
		if err != nil {
			return errors.New(fmt.Sprintf("reading URL Sig Keys: %v", err))
		}

		reencryptedKeys, err := postgres.AESEncrypt(jsonKeys, newKey)
		if err != nil {
			return errors.New(fmt.Sprintf("encrypting URL Sig Keys with new key: %v", err))
		}

		res, err := updateTx.Exec(`UPDATE url_sig_key SET data = $1 WHERE deliveryservice = $2`, []byte(reencryptedKeys), ds)
		if err != nil {
			return errors.New(fmt.Sprintf("updating URL Sig Keys for deliveryservice %s: %v", ds, err))
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return errors.New(fmt.Sprintf("determining rows affected for reencrypting URL Sig Keys with deliveryservice %s", ds))
		}
		if rowsAffected == 0 {
			return errors.New(fmt.Sprintf("no rows updated for reencrypting URL Sig Keys for deliveryservice %s", ds))
		}

	}

	return nil
}

func reEncryptUriSigningKeys(db *sqlx.DB, previousKey []byte, newKey []byte) error {
	tx, err := db.Begin()
	if err != nil {
		return errors.New(fmt.Sprintf("transaction begin failed %v %v", err, tx))
	}
	defer tx.Commit()

	rows, err := tx.Query("SELECT deliveryservice, data FROM uri_signing_key")
	if err != nil {
		return errors.New(fmt.Sprintf("querying: %v", err))
	}
	defer rows.Close()

	for rows.Next() {
		updateTx, err := db.Begin()
		if err != nil {
			return errors.New(fmt.Sprintf("transaction begin failed %v %v", err, updateTx))
		}
		defer updateTx.Commit()

		ds := ""
		var encryptedSslKeys []byte
		if err = rows.Scan(&ds, &encryptedSslKeys); err != nil {
			return errors.New(fmt.Sprintf("getting URI Signing Keys: %v", err))
		}
		jsonKeys, err := postgres.AESDecrypt(encryptedSslKeys, previousKey)
		if err != nil {
			return errors.New(fmt.Sprintf("reading URI Signing Keys: %v", err))
		}

		reencryptedKeys, err := postgres.AESEncrypt(jsonKeys, newKey)
		if err != nil {
			return errors.New(fmt.Sprintf("encrypting URI Signing Keys with new key: %v", err))
		}

		res, err := updateTx.Exec(`UPDATE uri_signing_key SET data = $1 WHERE deliveryservice = $2`, []byte(reencryptedKeys), ds)
		if err != nil {
			return errors.New(fmt.Sprintf("updating URI Signing Keys for deliveryservice %s: %v", ds, err))
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return errors.New(fmt.Sprintf("determining rows affected for reencrypting URI Signing Keys with deliveryservice %s", ds))
		}
		if rowsAffected == 0 {
			return errors.New(fmt.Sprintf("no rows updated for reencrypting URI Signing Keys for deliveryservice %s", ds))
		}

	}

	return nil
}

func reEncryptDNSSECKeys(db *sqlx.DB, previousKey []byte, newKey []byte) error {
	tx, err := db.Begin()
	if err != nil {
		return errors.New(fmt.Sprintf("transaction begin failed %v %v", err, tx))
	}
	defer tx.Commit()

	rows, err := tx.Query("SELECT cdn, data FROM dnssec")
	if err != nil {
		return errors.New(fmt.Sprintf("querying: %v", err))
	}
	defer rows.Close()

	for rows.Next() {
		updateTx, err := db.Begin()
		if err != nil {
			return errors.New(fmt.Sprintf("transaction begin failed %v %v", err, updateTx))
		}
		defer updateTx.Commit()

		ds := ""
		var encryptedSslKeys []byte
		if err = rows.Scan(&ds, &encryptedSslKeys); err != nil {
			return errors.New(fmt.Sprintf("getting DNSSEC Keys: %v", err))
		}
		jsonKeys, err := postgres.AESDecrypt(encryptedSslKeys, previousKey)
		if err != nil {
			return errors.New(fmt.Sprintf("reading DNSSEC Keys: %v", err))
		}

		reencryptedKeys, err := postgres.AESEncrypt(jsonKeys, newKey)
		if err != nil {
			return errors.New(fmt.Sprintf("encrypting DNSSEC Keys with new key: %v", err))
		}

		res, err := updateTx.Exec(`UPDATE dnssec SET data = $1 WHERE cdn = $2`, []byte(reencryptedKeys), ds)
		if err != nil {
			return errors.New(fmt.Sprintf("updating DNSSEC Keys for deliveryservice %s: %v", ds, err))
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return errors.New(fmt.Sprintf("determining rows affected for reencrypting DNSSEC Keys with deliveryservice %s", ds))
		}
		if rowsAffected == 0 {
			return errors.New(fmt.Sprintf("no rows updated for reencrypting DNSSEC Keys for deliveryservice %s", ds))
		}

	}

	return nil
}

// updateKeyFile saves the new key in the same location as the old one so TO can use the new one.
// It also saves the old one as a new file with a date to indicate when the re-encrypt was performed.
func updateKeyFile(previousKey []byte, previousKeyLocation string, newKey []byte) error {
	newKeyStr := base64.StdEncoding.EncodeToString(newKey)
	if err := ioutil.WriteFile(previousKeyLocation, []byte(newKeyStr), 0644); err != nil {
		return errors.New("writing key file " + previousKeyLocation + ": " + err.Error())
	}

	previousKeyStr := base64.StdEncoding.EncodeToString(previousKey)
	savedKeyLocation := fmt.Sprintf("%s/savedEncryptionKey-%s.txt", filepath.Dir(previousKeyLocation), time.Now().Format(time.RFC3339))
	if err := ioutil.WriteFile(savedKeyLocation, []byte(previousKeyStr), 0644); err != nil {
		return errors.New("writing key file " + savedKeyLocation + ": " + err.Error())
	}

	return nil
}

func die(message string) {
	fmt.Println(message)
	os.Exit(1)
}
