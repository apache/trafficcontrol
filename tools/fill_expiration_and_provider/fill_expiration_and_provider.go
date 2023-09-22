package main

/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with this
 * work for additional information regarding copyright ownership.  The ASF
 * licenses this file to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.  See the
 * License for the specific language governing permissions and limitations under
 * the License.
 */

import (
	"crypto/x509"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"flag"
	"fmt"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/jmoiron/sqlx"
	"io/ioutil"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

const PROPERTIES_FILE = "./fill_expiration_and_provider_conf.json"

func main() {
	aesKeyLocation := flag.String("aes-key", "/opt/traffic_ops/app/conf/aes.key", "The file path for the previous base64 encoded AES key. Default is /opt/traffic_ops/app/conf/aes.key.")
	cfg := flag.String("cfg", PROPERTIES_FILE, "The path for the configuration file. Default is "+PROPERTIES_FILE+".")
	help := flag.Bool("help", false, "Print usage information and exit.")
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	aesKey, err := readKey(*aesKeyLocation)
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

	rows, err := tx.Query("SELECT deliveryservice, cdn, version, data, provider, expiration FROM sslkey")
	if err != nil {
		die("querying: " + err.Error())
	}
	defer rows.Close()

	type expiryAndProvider struct {
		Provider   string
		Expiration time.Time
	}
	sslKeyMap := map[string]expiryAndProvider{}

	for rows.Next() {
		var ds string
		var cdn string
		var version string
		var encryptedSslKeys []byte
		provider := sql.NullString{}
		var expiration time.Time
		if err = rows.Scan(&ds, &cdn, &version, &encryptedSslKeys, &provider, &expiration); err != nil {
			die("getting SSL Keys: " + err.Error())
		}
		id := strings.Join([]string{ds, cdn, version}, ", ")
		jsonKeys, err := util.AESDecrypt(encryptedSslKeys, aesKey)
		if err != nil {
			die("reading SSL Keys: " + err.Error())
		}

		sslKey := tc.DeliveryServiceSSLKeysV15{}
		err = json.Unmarshal([]byte(jsonKeys), &sslKey)
		if err != nil {
			die("unmarshalling ssl keys: " + err.Error())
		}

		parsedCert := sslKey.Certificate
		err = Base64DecodeCertificate(&parsedCert)
		if err != nil {
			die("getting SSL keys for ID '" + id + "': " + err.Error())
		}

		block, _ := pem.Decode([]byte(parsedCert.Crt))
		if block == nil {
			die("Error decoding cert to parse expiration")
		}

		x509cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			die("Error parsing cert to get expiration - " + err.Error())
		}

		sslKeyMap[id] = expiryAndProvider{
			Provider:   sslKey.AuthType,
			Expiration: x509cert.NotAfter,
		}
	}

	for id, info := range sslKeyMap {
		if strings.Count(id, ",") != 2 {
			die("found id that does not contain 2 commas: " + id)
		}
		idParts := strings.Split(id, ", ")
		if len(idParts) != 3 {
			die(fmt.Sprintf("expected cert id string (ds, cdn, version) to have 3 parts but found %d in %s", len(idParts), idParts))
		}
		ds := idParts[0]
		cdn := idParts[1]
		version := idParts[2]
		res, err := tx.Exec(`UPDATE sslkey SET provider = $1, expiration = $2 WHERE deliveryservice = $3 AND cdn = $4 AND version = $5`, info.Provider, info.Expiration, ds, cdn, version)
		if err != nil {
			die(fmt.Sprintf("updating SSL Keys for %s, %s", id, err))
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			die(fmt.Sprintf("determining rows affected for expiration and provider in SSL Keys: %s: %s", id, err.Error()))
		}
		if rowsAffected == 0 {
			die(fmt.Sprintf("no rows updated for expiration and provider in SSL Keys for %s", id))
		}
	}
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
		return []byte{}, fmt.Errorf("reading file '"+keyLocation+"': %s", err)
	}
	keyBase64 = string(keyBase64Bytes)

	key, err := base64.StdEncoding.DecodeString(keyBase64)
	if err != nil {
		return []byte{}, fmt.Errorf("AES key cannot be decoded from base64: %s", err)
	}

	// verify the key works
	if err = util.ValidateAESKey(key); err != nil {
		return []byte{}, err
	}

	return key, nil
}

func die(message string) {
	fmt.Fprintln(os.Stderr, message)
	os.Exit(1)
}

func Base64DecodeCertificate(cert *tc.DeliveryServiceSSLKeysCertificate) error {
	csrDec, err := base64.StdEncoding.DecodeString(cert.CSR)
	if err != nil {
		return errors.New("base64 decoding csr: " + err.Error())
	}
	cert.CSR = string(csrDec)
	crtDec, err := base64.StdEncoding.DecodeString(cert.Crt)
	if err != nil {
		return errors.New("base64 decoding crt: " + err.Error())
	}
	cert.Crt = string(crtDec)
	keyDec, err := base64.StdEncoding.DecodeString(cert.Key)
	if err != nil {
		return errors.New("base64 decoding key: " + err.Error())
	}
	cert.Key = string(keyDec)
	return nil
}
