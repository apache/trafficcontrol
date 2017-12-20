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
	"crypto/tls"
	"encoding/json"
	"errors"
	"testing"

	"github.com/basho/riak-go-client"
	"github.com/jmoiron/sqlx"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

const (
	goodKeyset = `
{
  "Kabletown URI Authority 1": {
		"renewal_kid": "Second Key",
    "keys": [
      {
        "alg": "HS256",
        "kid": "First Key",
        "kty": "oct",
        "k": "Kh_RkUMj-fzbD37qBnDf_3e_RvQ3RP9PaSmVEpE24AM"
      },
      {
        "alg": "HS256",
        "kid": "Second Key",
        "kty": "oct",
        "k": "fZBpDBNbk2GqhwoB_DGBAsBxqQZVix04rIoLJ7p_RlE"
      }
    ]
  },
  "Kabletown URI Authority 2": {
    "keys": [
      {
        "alg": "HS256",
        "kid": "First Key",
        "kty": "oct",
        "k": "Kh_RkUMj-fzbD37qBnDf_3e_RvQ3RP9PaSmVEpE24AM"
      },
      {
        "alg": "HS256",
        "kid": "Third Key",
        "kty": "oct",
        "k": "fZBpDBNbk2GqhwoB_DGBAsBxqQZVix04rIoLJ7p_RlE"
      }
    ]
  }
}
`
	badJSONKeySet = `

  "Kabletown URI Authority 1": {
		"renewal_kid": "Second Key",
    "keys": [
      {
        "alg": "HS256",
        "kid": "First Key",
        "kty": "oct",
        "k": "Kh_RkUMj-fzbD37qBnDf_3e_RvQ3RP9PaSmVEpE24AM"
      },
      {
        "alg": "HS256",
        "kid": "Second Key",
        "kty": "oct",
        "k": "fZBpDBNbk2GqhwoB_DGBAsBxqQZVix04rIoLJ7p_RlE"
      }
    ]
  }
}
`
	noRenewalKidKeyset = `
{
  "Kabletown URI Authority 1": {
    "keys": [
      {
        "alg": "HS256",
        "kid": "First Key",
        "kty": "oct",
        "k": "Kh_RkUMj-fzbD37qBnDf_3e_RvQ3RP9PaSmVEpE24AM"
      },
      {
        "alg": "HS256",
        "kid": "Second Key",
        "kty": "oct",
        "k": "fZBpDBNbk2GqhwoB_DGBAsBxqQZVix04rIoLJ7p_RlE"
      }
    ]
  }
}
`
	noMatchingRenewalKidKeyset = `
{
  "Kabletown URI Authority 1": {
		"renewal_kid": "Second Key",
    "keys": [
      {
        "alg": "HS256",
        "kid": "First Key",
        "kty": "oct",
        "k": "Kh_RkUMj-fzbD37qBnDf_3e_RvQ3RP9PaSmVEpE24AM"
      },
      {
        "alg": "HS256",
        "kid": "Other Key",
        "kty": "oct",
        "k": "fZBpDBNbk2GqhwoB_DGBAsBxqQZVix04rIoLJ7p_RlE"
      }
    ]
  }
}
`
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

	_, err := fetchObjectValues("myobject", "bucket", cluster)
	if err != nil {
		t.Error("expected nil error got ", err)
	}

	_, err = fetchObjectValues("", "bucket", cluster)
	if err == nil {
		t.Error("expected an error because key is empty but got no error")
	}

	_, err = fetchObjectValues("myobject", "", cluster)
	if err == nil {
		t.Error("expected an error because the bucket name is empty but got no error")
	}

	_, err = fetchObjectValues("myobject", "mybucket", nil)
	if err == nil {
		t.Error("expected an error because the cluster is nil but got no error")
	}
}

func TestSaveObject(t *testing.T) {
	cluster := &MockStorageCluster{
		Running: true,
	}

	obj := riak.Object{}

	err := saveObject(&obj, "bucket", cluster)
	if err != nil {
		t.Error("expected nil error got ", err)
	}

	err = saveObject(&obj, "", cluster)
	if err == nil {
		t.Error("expected an error due to empty bucket name but got a nil error")
	}

	err = saveObject(nil, "bucket", cluster)
	if err == nil {
		t.Error("expected an error because the obj is nil but go no error", err)
	}

	err = saveObject(&obj, "bucket", nil)
	if err == nil {
		t.Error("expected an error because the cluster is nil but go no error", err)
	}
}

func TestDeleteObject(t *testing.T) {
	cluster := &MockStorageCluster{
		Running: true,
	}

	err := deleteObject("myobject", "bucket", cluster)
	if err != nil {
		t.Error("expected nil error got ", err)
	}

	err = deleteObject("", "bucket", cluster)
	if err == nil {
		t.Error("expected an empty key error but got nil error", err)
	}

	err = deleteObject("myobject", "", cluster)
	if err == nil {
		t.Error("expected an empty bucket error but got nil error", err)
	}

	err = deleteObject("myobject", "bucket", nil)
	if err == nil {
		t.Error("expected an nil cluster error but got nil error", err)
	}
}

func TestValidateURIKeyset(t *testing.T) {
	var keyset map[string]URISignerKeyset

	// unmarshal a good URISignerKeyset
	if err := json.Unmarshal([]byte(goodKeyset), &keyset); err != nil {
		t.Errorf("json.UnMarshal(): expected nil error, actual: %v", err)
	}

	// now validate it
	if err := validateURIKeyset(keyset); err != nil {
		t.Errorf("validateURIKeyset(): expected nil error, actual: %v", err)
	}

	// unmarshal a bad URISignerKeySet
	if err := json.Unmarshal([]byte(badJSONKeySet), &keyset); err == nil {
		t.Errorf("json.UnMarshal(): expected an error")
	}

	// unmarshal a good URISignerKeyset with a missing renewal_kid
	if err := json.Unmarshal([]byte(noRenewalKidKeyset), &keyset); err != nil {
		t.Errorf("json.UnMarshal(): expected nil error, actual: %v", err)
	}

	// now validate it, expect an erro due to missing renewal_kid
	if err := validateURIKeyset(keyset); err == nil {
		t.Errorf("validateURIKeyset(): expected an error")
	}

	// unmarshal a good URISignerKeyset with no matching kid to the renewal_kid
	if err := json.Unmarshal([]byte(noMatchingRenewalKidKeyset), &keyset); err != nil {
		t.Errorf("json.UnMarshal(): expected nil error, actual: %v", err)
	}

	// now validate it, expect an erro due to missing a matching key kid to renewal_kid
	if err := validateURIKeyset(keyset); err == nil {
		t.Errorf("validateURIKeyset(): expected an error")
	}
}

func TestGetRiakCluster(t *testing.T) {
	var cfg Config
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	rows1 := sqlmock.NewRows([]string{"s.host_name", "s.domain_name"})
	rows1.AddRow("www", "devnull.com")
	mock.ExpectQuery("SELECT").WillReturnRows(rows1)

	cfg.RiakAuthOptions = nil
	if _, err := getRiakCluster(db, cfg); err == nil {
		t.Errorf("expected an error due to nil RiakAuthoptions in the config but, go no error.")
	}

	cfg.RiakAuthOptions = &riak.AuthOptions{
		User:      "riakuser",
		Password:  "password",
		TlsConfig: &tls.Config{},
	}

	if _, err := getRiakCluster(db, cfg); err != nil {
		t.Errorf("expected no errors, actual: %s.", err)
	}

	rows2 := sqlmock.NewRows([]string{"s.host_name", "s.domain_name"})
	mock.ExpectQuery("SELECT").WillReturnRows(rows2)

	if _, err := getRiakCluster(db, cfg); err == nil {
		t.Errorf("expected an error due to no available riak servers.")
	}
}
