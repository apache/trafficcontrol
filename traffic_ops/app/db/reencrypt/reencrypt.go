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
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/trafficvault/backends/postgres"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

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
)

const PROPERTIES_FILE = "./reencrypt.conf"

func main() {
	previousKeyLocation := flag.String("previousKey", "", "The file path for the previous base64 encoded AES key.")
	newKeyLocation := flag.String("newKey", "", "The file path for the new base64 encoded AES key.")
	flag.Parse()

	if previousKeyLocation == nil || *previousKeyLocation == "" {
		fmt.Println("previousKey flag is required.")
		os.Exit(0)
	}
	if newKeyLocation == nil || *newKeyLocation == "" {
		fmt.Println("newKey flag is required.")
		os.Exit(0)
	}

	newKey, err := readKey(*newKeyLocation)
	if err != nil {
		fmt.Println("reading newKey: ", err.Error())
		os.Exit(0)
	}

	previousKey, err := readKey(*previousKeyLocation)
	if err != nil {
		fmt.Println("reading previousKey: ", err.Error())
		os.Exit(0)
	}

	dbConfBytes, err := ioutil.ReadFile(PROPERTIES_FILE)
	if err != nil {
		fmt.Println("reading db conf '", PROPERTIES_FILE, "': ", err.Error())
		os.Exit(0)
	}

	pgCfg := Config{}
	err = json.Unmarshal(dbConfBytes, &pgCfg)
	if err != nil {
		fmt.Println("unmarshalling '", PROPERTIES_FILE, "': ", err.Error())
		os.Exit(0)
	}

	sslStr := "require"
	if !pgCfg.SSL {
		sslStr = "disable"
	}
	db, err := sqlx.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s&fallback_application_name=trafficvault", pgCfg.User, pgCfg.Password, pgCfg.Hostname, pgCfg.Port, pgCfg.DBName, sslStr))
	if err != nil {
		fmt.Println("opening database: " + err.Error())
		os.Exit(0)
	}

	if err = reEncryptSslKeys(db, previousKey, newKey); err != nil {
		fmt.Println("re-encrypting SSL Keys: ", err.Error())
		os.Exit(0)
	}
	if err = reEncryptUrlSigKeys(db, previousKey, newKey); err != nil {
		fmt.Println("re-encrypting URL Sig Keys: ", err.Error())
		os.Exit(0)
	}
	if err = reEncryptUriSigningKeys(db, previousKey, newKey); err != nil {
		fmt.Println("re-encrypting URI Signing Keys: ", err.Error())
		os.Exit(0)
	}
	if err = reEncryptDNSSECKeys(db, previousKey, newKey); err != nil {
		fmt.Println("re-encrypting DNSSEC Keys: ", err.Error())
		os.Exit(0)
	}

	if err = updateKeyFile(previousKey, *previousKeyLocation, newKey); err != nil {
		fmt.Println("updating the key file: ", err.Error())
		os.Exit(0)
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
		return []byte{}, errors.New("AES key cannot be decoded from base64")
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
		fmt.Println("transaction begin failed ", err, " ", tx)
		os.Exit(0)
	}
	defer tx.Commit()

	rows, err := tx.Query("SELECT id, data FROM sslkey")
	if err != nil {
		fmt.Println("querying: ", err)
		os.Exit(0)
	}
	defer rows.Close()

	for rows.Next() {
		updateTx, err := db.Begin()
		if err != nil {
			fmt.Println("transaction begin failed ", err, " ", updateTx)
			os.Exit(0)
		}
		defer updateTx.Commit()

		id := 0
		var encryptedSslKeys []byte
		if err = rows.Scan(&id, &encryptedSslKeys); err != nil {
			fmt.Println("getting SSL Keys: ", err)
			return err
		}
		jsonKeys, err := postgres.AesDecrypt(encryptedSslKeys, previousKey)
		if err != nil {
			fmt.Println("reading SSL Keys: ", err)
			return err
		}

		reencryptedKeys, err := postgres.AesEncrypt(jsonKeys, newKey)
		if err != nil {
			fmt.Println("encrypting SSL Keys with new key: ", err)
			return err
		}

		res, err := updateTx.Exec(`UPDATE sslkey SET data = $1 WHERE id = $2`, []byte(reencryptedKeys), id)
		if err != nil {
			fmt.Println("updating SSL Keys for id ", id, ": ", err)
			return err
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			fmt.Println("error determining rows affected for reencrypting SSL Keys with id ", id)
			return errors.New(fmt.Sprintf("determining rows affected for reencrypting SSL Keys with id %d", id))
		}
		if rowsAffected == 0 {
			fmt.Println("no rows updated for reencrypting SSL Keys for id ", id)
			return errors.New(fmt.Sprintf("no rows updated for reencrypting SSL Keys for id %d", id))
		}

	}

	return nil
}

func reEncryptUrlSigKeys(db *sqlx.DB, previousKey []byte, newKey []byte) error {
	tx, err := db.Begin()
	if err != nil {
		fmt.Println("transaction begin failed ", err, " ", tx)
		os.Exit(0)
	}
	defer tx.Commit()

	rows, err := tx.Query("SELECT deliveryservice, data FROM url_sig_key")
	if err != nil {
		fmt.Println("querying: ", err)
		os.Exit(0)
	}
	defer rows.Close()

	for rows.Next() {
		updateTx, err := db.Begin()
		if err != nil {
			fmt.Println("transaction begin failed ", err, " ", updateTx)
			os.Exit(0)
		}
		defer updateTx.Commit()

		ds := ""
		var encryptedSslKeys []byte
		if err = rows.Scan(&ds, &encryptedSslKeys); err != nil {
			fmt.Println("getting URL Sig Keys: ", err)
			return err
		}
		jsonKeys, err := postgres.AesDecrypt(encryptedSslKeys, previousKey)
		if err != nil {
			fmt.Println("reading URL Sig Keys: ", err)
			return err
		}

		reencryptedKeys, err := postgres.AesEncrypt(jsonKeys, newKey)
		if err != nil {
			fmt.Println("encrypting URL Sig Keys with new key: ", err)
			return err
		}

		res, err := updateTx.Exec(`UPDATE url_sig_key SET data = $1 WHERE deliveryservice = $2`, []byte(reencryptedKeys), ds)
		if err != nil {
			fmt.Println("updating URL Sig Keys for deliveryservice ", ds, ": ", err)
			return err
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			fmt.Println("error determining rows affected for reencrypting URL Sig Keys with deliveryservice ", ds)
			return errors.New(fmt.Sprintf("determining rows affected for reencrypting URL Sig Keys with deliveryservice %d", ds))
		}
		if rowsAffected == 0 {
			fmt.Println("no rows updated for reencrypting URL Sig Keys for id ", ds)
			return errors.New(fmt.Sprintf("no rows updated for reencrypting URL Sig Keys for deliveryservice %d", ds))
		}

	}

	return nil
}

func reEncryptUriSigningKeys(db *sqlx.DB, previousKey []byte, newKey []byte) error {
	tx, err := db.Begin()
	if err != nil {
		fmt.Println("transaction begin failed ", err, " ", tx)
		os.Exit(0)
	}
	defer tx.Commit()

	rows, err := tx.Query("SELECT deliveryservice, data FROM uri_signing_key")
	if err != nil {
		fmt.Println("querying: ", err)
		os.Exit(0)
	}
	defer rows.Close()

	for rows.Next() {
		updateTx, err := db.Begin()
		if err != nil {
			fmt.Println("transaction begin failed ", err, " ", updateTx)
			os.Exit(0)
		}
		defer updateTx.Commit()

		ds := ""
		var encryptedSslKeys []byte
		if err = rows.Scan(&ds, &encryptedSslKeys); err != nil {
			fmt.Println("getting URI Signing Keys: ", err)
			return err
		}
		jsonKeys, err := postgres.AesDecrypt(encryptedSslKeys, previousKey)
		if err != nil {
			fmt.Println("reading URI Signing Keys: ", err)
			return err
		}

		reencryptedKeys, err := postgres.AesEncrypt(jsonKeys, newKey)
		if err != nil {
			fmt.Println("encrypting URI Signing Keys with new key: ", err)
			return err
		}

		res, err := updateTx.Exec(`UPDATE uri_signing_key SET data = $1 WHERE deliveryservice = $2`, []byte(reencryptedKeys), ds)
		if err != nil {
			fmt.Println("updating URI Signing Keys for deliveryservice ", ds, ": ", err)
			return err
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			fmt.Println("error determining rows affected for reencrypting URI Signing Keys with deliveryservice ", ds)
			return errors.New(fmt.Sprintf("determining rows affected for reencrypting URI Signing Keys with deliveryservice %d", ds))
		}
		if rowsAffected == 0 {
			fmt.Println("no rows updated for reencrypting URI Signing Keys for id ", ds)
			return errors.New(fmt.Sprintf("no rows updated for reencrypting URI Signing Keys for deliveryservice %d", ds))
		}

	}

	return nil
}

func reEncryptDNSSECKeys(db *sqlx.DB, previousKey []byte, newKey []byte) error {
	tx, err := db.Begin()
	if err != nil {
		fmt.Println("transaction begin failed ", err, " ", tx)
		os.Exit(0)
	}
	defer tx.Commit()

	rows, err := tx.Query("SELECT cdn, data FROM dnssec")
	if err != nil {
		fmt.Println("querying: ", err)
		os.Exit(0)
	}
	defer rows.Close()

	for rows.Next() {
		updateTx, err := db.Begin()
		if err != nil {
			fmt.Println("transaction begin failed ", err, " ", updateTx)
			os.Exit(0)
		}
		defer updateTx.Commit()

		ds := ""
		var encryptedSslKeys []byte
		if err = rows.Scan(&ds, &encryptedSslKeys); err != nil {
			fmt.Println("getting DNSSEC Keys: ", err)
			return err
		}
		jsonKeys, err := postgres.AesDecrypt(encryptedSslKeys, previousKey)
		if err != nil {
			fmt.Println("reading DNSSEC Keys: ", err)
			return err
		}

		reencryptedKeys, err := postgres.AesEncrypt(jsonKeys, newKey)
		if err != nil {
			fmt.Println("encrypting DNSSEC Keys with new key: ", err)
			return err
		}

		res, err := updateTx.Exec(`UPDATE dnssec SET data = $1 WHERE cdn = $2`, []byte(reencryptedKeys), ds)
		if err != nil {
			fmt.Println("updating DNSSEC Keys for deliveryservice ", ds, ": ", err)
			return err
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			fmt.Println("error determining rows affected for reencrypting DNSSEC Keys with deliveryservice ", ds)
			return errors.New(fmt.Sprintf("determining rows affected for reencrypting DNSSEC Keys with deliveryservice %d", ds))
		}
		if rowsAffected == 0 {
			fmt.Println("no rows updated for reencrypting DNSSEC Keys for id ", ds)
			return errors.New(fmt.Sprintf("no rows updated for reencrypting DNSSEC Keys for deliveryservice %d", ds))
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
