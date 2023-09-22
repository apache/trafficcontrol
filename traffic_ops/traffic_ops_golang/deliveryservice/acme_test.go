package deliveryservice

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
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"testing"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/config"
	"github.com/go-acme/lego/challenge/dns01"

	"github.com/jmoiron/sqlx"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestGetStoredAcmeAccountInfo(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("expected no error while generating key, but got %v", err)
	}
	keyBuf := bytes.Buffer{}
	keyDer := x509.MarshalPKCS1PrivateKey(priv)
	err = pem.Encode(&keyBuf, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: keyDer})
	if err != nil {
		t.Fatalf("expected no error while encoding key, but got %v", err)
	}

	mock.ExpectBegin()
	rows := sqlmock.NewRows([]string{"email", "private_key", "uri"})
	rows.AddRow("testuser@blah.com", keyBuf.Bytes(), "https://uri.com")
	mock.ExpectQuery("SELECT email, private_key, uri").WithArgs("testuser@blah.com", "Lets Encrypt").WillReturnRows(rows)

	info, err := getStoredAcmeAccountInfo(db.MustBegin().Tx, "testuser@blah.com", "Lets Encrypt")
	if err != nil {
		t.Errorf("expected no error while getting stored acme account into, but got %v", err)
	}
	if info == nil {
		t.Fatalf("expected valid acme account info in response, but got nothing")
	}
	if info.Email != "testuser@blah.com" {
		t.Errorf("expected email to be testuser@blah.com, but got %s", info.Email)
	}
	if info.Key != string(keyBuf.Bytes()) {
		t.Errorf("expected key to be %s, but got %s", string(keyBuf.Bytes()), info.Key)
	}
	if info.URI != "https://uri.com" {
		t.Errorf("expected uri to be https://uri.com, but got %s", info.URI)
	}
}

func TestGetAcmeAccountConfig(t *testing.T) {
	cfgAcmeAccounts := make([]config.ConfigAcmeAccount, 0)
	cfg := config.Config{
		AcmeAccounts: cfgAcmeAccounts,
		ConfigLetsEncrypt: config.ConfigLetsEncrypt{
			Email:       "testuser@apache.org",
			Environment: "production",
		},
	}
	c := GetAcmeAccountConfig(&cfg, tc.LetsEncryptAuthType)
	if c == nil {
		t.Fatalf("expected a valid Acme Account Config in response, but got nothing")
	}
	if c.UserEmail != cfg.Email {
		t.Errorf("expected user email to be %s, but got %s", cfg.Email, c.UserEmail)
	}
	if c.AcmeUrl != "https://acme-v02.api.letsencrypt.org/directory" {
		t.Errorf("expected AcmeProvider to be https://acme-v02.api.letsencrypt.org/directory, but got %s", c.AcmeUrl)
	}
	if c.AcmeProvider != tc.LetsEncryptAuthType {
		t.Errorf("expected AcmeProvider to be Lets Encrypt, but got %s", c.AcmeProvider)
	}
}

func TestDNSProviderTrafficRouter_Present(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	d := DNSProviderTrafficRouter{
		db:    db,
		xmlId: util.Ptr("dsXMLID"),
	}
	keyAuthShaBytes := sha256.Sum256([]byte("blah"))
	value := base64.RawURLEncoding.EncodeToString(keyAuthShaBytes[:sha256.Size])
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO dnschallenges").WithArgs("_acme-challenge.test.", value, *d.xmlId).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	err = d.Present("test", "token", "blah")
	if err != nil {
		t.Errorf("expected no error, but got %v", err)
	}
}

func TestDNSProviderTrafficRouter_Cleanup(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	d := DNSProviderTrafficRouter{
		db:    db,
		xmlId: util.Ptr("dsXMLID"),
	}
	keyAuthShaBytes := sha256.Sum256([]byte("blah"))
	value := base64.RawURLEncoding.EncodeToString(keyAuthShaBytes[:sha256.Size])
	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM dnschallenges").WithArgs("_acme-challenge.test.", value).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	err = d.CleanUp("test", "token", "blah")
	if err != nil {
		t.Errorf("expected no error, but got %v", err)
	}
}

func TestGetRecord(t *testing.T) {
	fqdn, val := dns01.GetRecord("test", "blah")
	keyAuthShaBytes := sha256.Sum256([]byte("blah"))
	value := base64.RawURLEncoding.EncodeToString(keyAuthShaBytes[:sha256.Size])
	if fqdn != "_acme-challenge.test." {
		t.Errorf("expected fqdn to be _acme-challenge.test., but got %s", fqdn)
	}
	if val != value {
		t.Errorf("expected returned value to be %s, but got %s", value, val)
	}
}
