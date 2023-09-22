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
	"bytes"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/apache/trafficcontrol/v8/lib/go-util"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const PROPERTIES_FILE = "./reencrypt.conf"

func main() {
	previousKeyLocation := flag.String("previous-key", "/opt/traffic_ops/app/conf/aes.key", "(Optional) The file path for the previous base64 encoded AES key. Default is /opt/traffic_ops/app/conf/aes.key.")
	newKeyLocation := flag.String("new-key", "/opt/traffic_ops/app/conf/new.key", "(Optional) The file path for the new base64 encoded AES key. Default is /opt/traffic_ops/app/conf/new.key.")
	cfg := flag.String("cfg", PROPERTIES_FILE, "(Optional) The path for the configuration file. Default is "+PROPERTIES_FILE+".")
	help := flag.Bool("help", false, "(Optional) Print usage information and exit.")
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
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

	tx, err := db.Begin()
	if err != nil {
		die(fmt.Sprintf("transaction begin failed %v %v ", err, tx))
	}
	defer tx.Commit()

	if err = reEncryptSslKeys(tx, previousKey, newKey); err != nil {
		tx.Rollback()
		die("re-encrypting SSL Keys: " + err.Error())
	}
	if err = reEncryptUrlSigKeys(tx, previousKey, newKey); err != nil {
		tx.Rollback()
		die("re-encrypting URL Sig Keys: " + err.Error())
	}
	if err = reEncryptUriSigningKeys(tx, previousKey, newKey); err != nil {
		tx.Rollback()
		die("re-encrypting URI Signing Keys: " + err.Error())
	}
	if err = reEncryptDNSSECKeys(tx, previousKey, newKey); err != nil {
		tx.Rollback()
		die("re-encrypting DNSSEC Keys: " + err.Error())
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
		return []byte{}, fmt.Errorf("AES key cannot be decoded from base64: %w", err)
	}

	// verify the key works
	if err = util.ValidateAESKey(key); err != nil {
		return []byte{}, err
	}

	return key, nil
}

type sslInfo struct {
	xmlId        string
	version      string
	previousData []byte
	newData      []byte
}

func reEncryptSslKeys(tx *sql.Tx, previousKey []byte, newKey []byte) error {
	rows, err := tx.Query("SELECT deliveryservice, version, data FROM sslkey")
	if err != nil {
		return fmt.Errorf("querying: %w", err)
	}
	defer rows.Close()

	var sslKeyInfos []sslInfo

	for rows.Next() {
		sslKeyInfo := sslInfo{}

		if err = rows.Scan(&sslKeyInfo.xmlId, &sslKeyInfo.version, &sslKeyInfo.previousData); err != nil {
			return fmt.Errorf("getting SSL Keys: %w", err)
		}
		jsonKeys, err := util.AESDecrypt(sslKeyInfo.previousData, previousKey)
		if err != nil {
			return fmt.Errorf("reading SSL Keys: %w", err)
		}

		if !bytes.HasPrefix(jsonKeys, []byte("{")) {
			return fmt.Errorf("decrypted SSL Key did not have prefix '{' for xmlid %s", sslKeyInfo.xmlId)
		}

		sslKeyInfo.newData, err = util.AESEncrypt(jsonKeys, newKey)
		if err != nil {
			return fmt.Errorf("encrypting SSL Keys with new key: %w", err)
		}

		sslKeyInfos = append(sslKeyInfos, sslKeyInfo)
	}

	for _, sslKeyInfo := range sslKeyInfos {
		res, err := tx.Exec(`UPDATE sslkey SET data = $1 WHERE deliveryservice = $2 AND version = $3`, sslKeyInfo.newData, sslKeyInfo.xmlId, sslKeyInfo.version)
		if err != nil {
			return fmt.Errorf("updating SSL Keys for xmlid %s: %w", sslKeyInfo.xmlId, err)
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return fmt.Errorf("determining rows affected for reencrypting SSL Keys with xmlid %s: %w", sslKeyInfo.xmlId, err)
		}
		if rowsAffected == 0 {
			return fmt.Errorf("no rows updated for reencrypting SSL Keys for xmlid %s", sslKeyInfo.xmlId)
		}
	}

	return nil
}

func reEncryptUrlSigKeys(tx *sql.Tx, previousKey []byte, newKey []byte) error {
	rows, err := tx.Query("SELECT deliveryservice, data FROM url_sig_key")
	if err != nil {
		return fmt.Errorf("querying: %w", err)
	}
	defer rows.Close()

	urlSigKeysMap := map[string][]byte{}

	for rows.Next() {
		ds := ""
		var encryptedSslKeys []byte
		if err = rows.Scan(&ds, &encryptedSslKeys); err != nil {
			return fmt.Errorf("getting URL Sig Keys: %w", err)
		}
		jsonKeys, err := util.AESDecrypt(encryptedSslKeys, previousKey)
		if err != nil {
			return fmt.Errorf("reading URL Sig Keys: %w", err)
		}

		if !bytes.HasPrefix(jsonKeys, []byte("{")) {
			return fmt.Errorf("decrypted URL Sig Key did not have prefix '{' for ds: %s", ds)
		}

		reencryptedKeys, err := util.AESEncrypt(jsonKeys, newKey)
		if err != nil {
			return fmt.Errorf("encrypting URL Sig Keys with new key: %w", err)
		}

		urlSigKeysMap[ds] = reencryptedKeys
	}

	for ds, reencryptedKeys := range urlSigKeysMap {
		res, err := tx.Exec(`UPDATE url_sig_key SET data = $1 WHERE deliveryservice = $2`, reencryptedKeys, ds)
		if err != nil {
			return fmt.Errorf("updating URL Sig Keys for deliveryservice %s: %w", ds, err)
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return fmt.Errorf("determining rows affected for reencrypting URL Sig Keys with deliveryservice %s: %w", ds, err)
		}
		if rowsAffected == 0 {
			return fmt.Errorf("no rows updated for reencrypting URL Sig Keys for deliveryservice %s", ds)
		}
	}

	return nil
}

func reEncryptUriSigningKeys(tx *sql.Tx, previousKey []byte, newKey []byte) error {
	rows, err := tx.Query("SELECT deliveryservice, data FROM uri_signing_key")
	if err != nil {
		return fmt.Errorf("querying: %w", err)
	}
	defer rows.Close()

	uriSigningKeyMap := map[string][]byte{}

	for rows.Next() {
		ds := ""
		var encryptedSslKeys []byte
		if err = rows.Scan(&ds, &encryptedSslKeys); err != nil {
			return fmt.Errorf("getting URI Signing Keys: %w", err)
		}
		jsonKeys, err := util.AESDecrypt(encryptedSslKeys, previousKey)
		if err != nil {
			return fmt.Errorf("reading URI Signing Keys: %w", err)
		}

		if !bytes.HasPrefix(jsonKeys, []byte("{")) {
			return fmt.Errorf("decrypted URI Signing Key did not have prefix '{' for ds: %s", ds)
		}

		reencryptedKeys, err := util.AESEncrypt(jsonKeys, newKey)
		if err != nil {
			return fmt.Errorf("encrypting URI Signing Keys with new key: %w", err)
		}

		uriSigningKeyMap[ds] = reencryptedKeys
	}

	for ds, reencryptedKeys := range uriSigningKeyMap {
		res, err := tx.Exec(`UPDATE uri_signing_key SET data = $1 WHERE deliveryservice = $2`, reencryptedKeys, ds)
		if err != nil {
			return fmt.Errorf("updating URI Signing Keys for deliveryservice %s: %w", ds, err)
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return fmt.Errorf("determining rows affected for reencrypting URI Signing Keys with deliveryservice %s: %w", ds, err)
		}
		if rowsAffected == 0 {
			return fmt.Errorf("no rows updated for reencrypting URI Signing Keys for deliveryservice %s", ds)
		}
	}

	return nil
}

func reEncryptDNSSECKeys(tx *sql.Tx, previousKey []byte, newKey []byte) error {
	rows, err := tx.Query("SELECT cdn, data FROM dnssec")
	if err != nil {
		return fmt.Errorf("querying: %w", err)
	}
	defer rows.Close()

	dnssecKeyMap := map[string][]byte{}

	for rows.Next() {
		cdn := ""
		var encryptedSslKeys []byte
		if err = rows.Scan(&cdn, &encryptedSslKeys); err != nil {
			return fmt.Errorf("getting DNSSEC Keys: %w", err)
		}
		jsonKeys, err := util.AESDecrypt(encryptedSslKeys, previousKey)
		if err != nil {
			return fmt.Errorf("reading DNSSEC Keys: %w", err)
		}

		if !bytes.HasPrefix(jsonKeys, []byte("{")) {
			return fmt.Errorf("decrypted DNSSEC Key did not have prefix '{' for cdn: %s", cdn)
		}

		reencryptedKeys, err := util.AESEncrypt(jsonKeys, newKey)
		if err != nil {
			return fmt.Errorf("encrypting DNSSEC Keys with new key: %w", err)
		}

		dnssecKeyMap[cdn] = reencryptedKeys
	}

	for cdn, reencryptedKeys := range dnssecKeyMap {
		res, err := tx.Exec(`UPDATE dnssec SET data = $1 WHERE cdn = $2`, reencryptedKeys, cdn)
		if err != nil {
			return fmt.Errorf("updating DNSSEC Keys for cdn %s: %w", cdn, err)
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return fmt.Errorf("determining rows affected for reencrypting DNSSEC Keys with cdn %s: %w", cdn, err)
		}
		if rowsAffected == 0 {
			return fmt.Errorf("no rows updated for reencrypting DNSSEC Keys for cdn %s", cdn)
		}
	}

	return nil
}

func die(message string) {
	fmt.Fprintln(os.Stderr, message)
	os.Exit(1)
}
