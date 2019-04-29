package foo

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
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
)

//we need a type alias to define functions on

type TOFoo struct {
	api.APIInfoImpl `json:"-"`
	tc.Foo
}

func (foo *TOFoo) GetKeys() (map[string]interface{}, bool) {
	if foo.ID == nil {
		return map[string]interface{}{"id": 0}, false
	}
	return map[string]interface{}{"id": *foo.ID}, true
}

func (foo *TOFoo) GetAuditName() string {
	if foo.Name != nil {
		return *foo.Name
	}
	if foo.ID != nil {
		return strconv.Itoa(*foo.ID)
	}
	return "unknown"
}

func (foo *TOFoo) GetKeyFieldsInfo() []api.KeyFieldInfo {
	return []api.KeyFieldInfo{{"id", api.GetIntKey}}
}

func (foo *TOFoo) APIInfo() *api.APIInfo { return foo.ReqInfo }

func (foo *TOFoo) SetKeys(keys map[string]interface{}) {
	i, _ := keys["id"].(int) //this utilizes the non panicking type assertion, if the thrown away ok variable is false i will be the zero of the type, 0 here.
	foo.ID = &i
}

func (foo *TOFoo) GetType() string {
	return "foo"
}

func Create(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	foo := tc.Foo{} // the struct to parse into would need to change for each specific minor version
	if err := api.Parse(r.Body, inf.Tx.Tx, &foo); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("decoding: "+err.Error()), nil)
		return
	}

	if err := foo.Validate(inf.Tx.Tx); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("invalid request: "+err.Error()), nil)
		return
	}
	log.Infoln("here we would call tx.Query(insertQuery(), &foo.Name, &foo.A)")
	api.WriteRespAlertObj(w, r, tc.SuccessLevel, "Foo creation was successful.", []tc.Foo{foo})
}

// Read is the Foo implementation of the Reader interface
func (foo *TOFoo) Read() ([]interface{}, error, error, int) {
	returnable := []interface{}{}
	log.Infoln("here we could call tx.Query(selectQuery()) and rows.Scan(&foo.ID, &foo.Name, &foo.A)")
	foos := []tc.Foo{
		{
			ID:   util.IntPtr(1),
			Name: util.StrPtr("one"),
			A:    util.StrPtr("A1"),
		},
		{
			ID:   util.IntPtr(2),
			Name: util.StrPtr("two"),
			A:    util.StrPtr("A2"),
		},
		{
			ID:   util.IntPtr(3),
			Name: util.StrPtr("three"),
			A:    util.StrPtr("A3"),
		},
	}

	for _, f := range foos {
		returnable = append(returnable, f)
	}
	return returnable, nil, nil, http.StatusOK
}

func Update(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, []string{"id"})
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	id := inf.IntParams["id"]

	foo := tc.Foo{}
	if err := json.NewDecoder(r.Body).Decode(&foo); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("malformed JSON: "+err.Error()), nil)
		return
	}
	foo.ID = &id

	if err := foo.Validate(inf.Tx.Tx); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("invalid request: "+err.Error()), nil)
		return
	}

	log.Infoln("here we would call tx.Query(updateQuery(), &foo.Name, &foo.A)")
	api.WriteRespAlertObj(w, r, tc.SuccessLevel, "Foo update was successful.", []tc.Foo{foo})
}

// Delete is the Foo implementation of the Deleter interface.
func (foo *TOFoo) Delete() (error, error, int) {
	if foo.ID == nil {
		return errors.New("missing id"), nil, http.StatusBadRequest
	}
	return nil, nil, http.StatusOK
}

func selectQuery() string {
	return `
SELECT
f.name,
f.A
FROM foo f
`
}

func updateFooQuery() string {
	return `
UPDATE foo SET
name=$1,
A=$2,
WHERE id=$3
`
}

func insertQuery() string {
	return `
INSERT INTO foo (
name,
A)
VALUES ($1,$2)
RETURNING id
`
}

func deleteQuery() string {
	return `
DELETE FROM foo
WHERE id = $1
`
}
