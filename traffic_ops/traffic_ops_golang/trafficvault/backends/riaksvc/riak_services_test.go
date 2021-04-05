package riaksvc

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
	"context"
	"crypto/tls"
	"errors"
	"testing"
	"time"

	"github.com/basho/riak-go-client"
	"github.com/jmoiron/sqlx"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

type MockStorageCluster struct {
	Bucket  string
	Key     string
	Running bool
}

func (mc MockStorageCluster) Stop() error {
	if mc.Running == true {
		mc.Running = false
	} else {
		return errors.New("the cluster is not started")
	}
	return nil
}

func (mc MockStorageCluster) Start() error {
	if mc.Running == false {
		mc.Running = true
	} else {
		return errors.New("the cluster is already started")
	}
	return nil
}

func (mc MockStorageCluster) Execute(command riak.Command) error {
	return nil
}

func TestFetchObjectValues(t *testing.T) {
	cluster := &MockStorageCluster{
		Running: true,
	}

	_, err := FetchObjectValues("myobject", "bucket", cluster)
	if err != nil {
		t.Error("expected nil error got ", err)
	}

	_, err = FetchObjectValues("", "bucket", cluster)
	if err == nil {
		t.Error("expected an error because key is empty but got no error")
	}

	_, err = FetchObjectValues("myobject", "", cluster)
	if err == nil {
		t.Error("expected an error because the bucket name is empty but got no error")
	}

	_, err = FetchObjectValues("myobject", "mybucket", nil)
	if err == nil {
		t.Error("expected an error because the cluster is nil but got no error")
	}
}

func TestSaveObject(t *testing.T) {
	cluster := &MockStorageCluster{
		Running: true,
	}

	obj := riak.Object{}

	err := SaveObject(&obj, "bucket", cluster)
	if err != nil {
		t.Error("expected nil error got ", err)
	}

	err = SaveObject(&obj, "", cluster)
	if err == nil {
		t.Error("expected an error due to empty bucket name but got a nil error")
	}

	err = SaveObject(nil, "bucket", cluster)
	if err == nil {
		t.Error("expected an error because the obj is nil but go no error", err)
	}

	err = SaveObject(&obj, "bucket", nil)
	if err == nil {
		t.Error("expected an error because the cluster is nil but go no error", err)
	}
}

func TestDeleteObject(t *testing.T) {
	cluster := &MockStorageCluster{
		Running: true,
	}

	err := DeleteObject("myobject", "bucket", cluster)
	if err != nil {
		t.Error("expected nil error got ", err)
	}

	err = DeleteObject("", "bucket", cluster)
	if err == nil {
		t.Error("expected an empty key error but got nil error", err)
	}

	err = DeleteObject("myobject", "", cluster)
	if err == nil {
		t.Error("expected an empty bucket error but got nil error", err)
	}

	err = DeleteObject("myobject", "bucket", nil)
	if err == nil {
		t.Error("expected an nil cluster error but got nil error", err)
	}
}

func TestSetTLSVersion(t *testing.T) {
	riakConfig := &TOAuthOptions{
		AuthOptions:   riak.AuthOptions{TlsConfig: &tls.Config{}},
		MaxTLSVersion: nil,
	}
	validVersion := "1.1"
	riakConfig.MaxTLSVersion = &validVersion
	if err := setMaxTLSVersion(riakConfig); err != nil {
		t.Error("expected nil but got ", err)
	}
	if riakConfig.TlsConfig.MaxVersion != tls.VersionTLS11 {
		t.Errorf("expected the TlsConfig's max version to be set to %v, but instead got %v.", tls.VersionTLS11, riakConfig.TlsConfig.MaxVersion)
	}

	invalidVersion := "1.a"
	riakConfig.MaxTLSVersion = &invalidVersion
	if err := setMaxTLSVersion(riakConfig); err == nil {
		t.Error("expected error due to an invalid TLS version but got no error.")
	}

	riakConfig.TlsConfig.MaxVersion = 0
	riakConfig.MaxTLSVersion = nil
	_ = setMaxTLSVersion(riakConfig)
	if riakConfig.TlsConfig.MaxVersion != tls.VersionTLS11 {
		t.Errorf("by default, expected the TlsConfig's max version to be set to %v, but instead got %v.", tls.VersionTLS11, riakConfig.TlsConfig.MaxVersion)
	}

}

func TestGetRiakCluster(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	mock.ExpectBegin()
	rows1 := sqlmock.NewRows([]string{"fqdn"})
	rows1.AddRow("www.devnull.com")

	dbCtx, _ := context.WithTimeout(context.TODO(), time.Duration(10)*time.Second)
	tx, err := db.BeginTx(dbCtx, nil)
	if err != nil {
		t.Fatalf("creating transaction: %v", err)
	}
	defer tx.Commit()

	mock.ExpectQuery("SELECT").WillReturnError(errors.New("foo"))
	if _, err := GetRiakServers(tx, nil); err == nil {
		t.Errorf("expected an error retrieving nil servers.")
	}

	mock.ExpectQuery("SELECT").WillReturnRows(rows1)
	servers, err := GetRiakServers(tx, nil)
	if err != nil {
		t.Errorf("expected to receive servers: %v", err)
	}

	if _, err := GetRiakCluster(servers, nil); err == nil {
		t.Errorf("expected an error due to nil RiakAuthoptions in the config but, go no error.")
	}

	authOptions := riak.AuthOptions{
		User:      "riakuser",
		Password:  "password",
		TlsConfig: &tls.Config{},
	}

	if _, err := GetRiakCluster(servers, &authOptions); err != nil {
		t.Errorf("expected no errors, actual: %v", err)
	}

	rows2 := sqlmock.NewRows([]string{"s.host_name", "s.domain_name"})
	mock.ExpectQuery("SELECT").WillReturnRows(rows2)

	if _, err := GetPooledCluster(tx, &authOptions, nil); err == nil {
		t.Errorf("expected an error due to no available riak servers.")
	}
}
