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

func CreateV15(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	foo := tc.FooV15{} // the struct to parse into would need to change for each specific minor version
	if err := json.NewDecoder(r.Body).Decode(&foo); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("decoding: "+err.Error()), nil)
	}
	res := createV15(w, r, inf, foo)
	if res != nil {
		api.WriteRespAlertObj(w, r, tc.SuccessLevel, "Foo creation was successful.", []tc.FooV15{*res})
	}
}

func createV15(w http.ResponseWriter, r *http.Request, info *api.APIInfo, foo tc.FooV15) *tc.FooV15 {
	fooV16 := tc.FooV16{FooV15: foo}
	res := createV16(w, r, info, fooV16)
	if res != nil {
		return &res.FooV15
	}
	return nil
}

func CreateV16(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	foo := tc.FooV16{} // TODO: this needed to change when we added a V17. Can we do better?
	if err := json.NewDecoder(r.Body).Decode(&foo); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("decoding: "+err.Error()), nil)
	}
	res := createV16(w, r, inf, foo)
	if res != nil {
		api.WriteRespAlertObj(w, r, tc.SuccessLevel, "Foo creation was successful.", []tc.FooV16{*res})
	}
}

func createV16(w http.ResponseWriter, r *http.Request, info *api.APIInfo, foo tc.FooV16) *tc.FooV16 {
	fooV17 := tc.FooV17{FooV16: foo}
	res := createV17(w, r, info, fooV17)
	if res != nil {
		return &res.FooV16
	}
	return nil
}

func CreateV17(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	foo := tc.FooV17{}
	if err := json.NewDecoder(r.Body).Decode(&foo); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("decoding: "+err.Error()), nil)
	}
	res := createV17(w, r, inf, foo)
	if res != nil {
		api.WriteRespAlertObj(w, r, tc.SuccessLevel, "Foo creation was successful.", []tc.FooV17{*res})
	}

}

func createV17(w http.ResponseWriter, r *http.Request, info *api.APIInfo, foo tc.FooV17) *tc.FooV17 {
	fooV18 := tc.FooV18{FooV17: foo}
	res := createV18(w, r, info, fooV18)
	if res != nil {
		return &res.FooV17
	}
	return nil
}

func CreateV18(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	foo := tc.FooV18{}
	if err := json.NewDecoder(r.Body).Decode(&foo); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("decoding: "+err.Error()), nil)
	}
	res := createV18(w, r, inf, foo)
	if res != nil {
		api.WriteRespAlertObj(w, r, tc.SuccessLevel, "Foo creation was successful.", []tc.FooV18{*res})
	}

}

func createV18(w http.ResponseWriter, r *http.Request, info *api.APIInfo, foo tc.FooV18) *tc.FooV18 {
	fooLatest := tc.Foo(foo) // NOTE: the "latest" subhandler might always have to do this?
	if err := fooLatest.Validate(info.Tx.Tx); err != nil {
		api.HandleErr(w, r, info.Tx.Tx, http.StatusBadRequest, errors.New("invalid request: "+err.Error()), nil)
		return nil
	}
	log.Infoln("here we would call tx.Query(insertQuery(), &foo.Name, &foo.A)")

	fooV18 := tc.FooV18(fooLatest)
	return &fooV18
}

// Read is the Foo implementation of the Reader interface. For multiple minor versions, might need to rework this to
// better handle minor versions with the Reader interface. Otherwise, it might make sense to just not use the Reader
// interface for endpoints with lots of minor versions. Or, have a "shared" Read handler that each minor version calls,
// then extracts its specific versioned struct out of the result from the shared handler. I think that is mainly how the
// deliveryservices Read handlers work today.
func (foo *TOFoo) Read() ([]interface{}, error, error, int) {
	returnable := []interface{}{}
	log.Infoln("here we could call tx.Query(selectQuery()) and rows.Scan(&foo.ID, &foo.Name, &foo.A)")
	foos := []tc.Foo{
		{},
		{},
		{},
	}

	for _, f := range foos {
		returnable = append(returnable, f)
	}
	return returnable, nil, nil, http.StatusOK
}

func UpdateV15(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, []string{"id"})
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	foo := tc.FooV15{}
	if err := json.NewDecoder(r.Body).Decode(&foo); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("malformed JSON: "+err.Error()), nil)
		return
	}

	res := updateV15(w, r, inf, foo)
	if res != nil {
		api.WriteRespAlertObj(w, r, tc.SuccessLevel, "Foo update was successful.", []tc.FooV15{*res})
	}
}

func updateV15(w http.ResponseWriter, r *http.Request, inf *api.APIInfo, foo tc.FooV15) *tc.FooV15 {
	fooV16 := tc.FooV16{FooV15: foo}

	log.Infoln("here we would query the DB for the existing B value (a 1.6 field) to populate a 1.6 request, essentially upgrading this 1.5 request into a 1.6 request")

	res := updateV16(w, r, inf, fooV16)
	if res != nil {
		return &res.FooV15
	}
	return nil
}

func UpdateV16(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, []string{"id"})
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	foo := tc.FooV16{}
	if err := json.NewDecoder(r.Body).Decode(&foo); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("malformed JSON: "+err.Error()), nil)
		return
	}

	res := updateV16(w, r, inf, foo)
	if res != nil {
		api.WriteRespAlertObj(w, r, tc.SuccessLevel, "Foo update was successful.", []tc.FooV16{*res})
	}
}

func updateV16(w http.ResponseWriter, r *http.Request, inf *api.APIInfo, foo tc.FooV16) *tc.FooV16 {
	fooV17 := tc.FooV17{FooV16: foo}

	log.Infoln("here we would query the DB for the existing C value (a 1.7 field) to populate a 1.7 request, essentially upgrading this 1.6 request into a 1.7 request")

	res := updateV17(w, r, inf, fooV17)
	if res != nil {
		return &res.FooV16
	}
	return nil
}

func UpdateV17(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, []string{"id"})
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	foo := tc.FooV17{}
	if err := json.NewDecoder(r.Body).Decode(&foo); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("malformed JSON: "+err.Error()), nil)
		return
	}

	res := updateV17(w, r, inf, foo)
	if res != nil {
		api.WriteRespAlertObj(w, r, tc.SuccessLevel, "Foo update was successful.", []tc.FooV17{*res})
	}
}

func updateV17(w http.ResponseWriter, r *http.Request, inf *api.APIInfo, foo tc.FooV17) *tc.FooV17 {
	fooV18 := tc.FooV18{FooV17: foo}

	log.Infoln("here we would query the DB for the existing D value (a 1.8 field) to populate a 1.8 request, essentially upgrading this 1.7 request into a 1.8 request")

	res := updateV18(w, r, inf, fooV18)
	if res != nil {
		return &res.FooV17
	}
	return nil
}

func UpdateV18(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, []string{"id"})
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	foo := tc.FooV18{}
	if err := json.NewDecoder(r.Body).Decode(&foo); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("malformed JSON: "+err.Error()), nil)
		return
	}

	res := updateV18(w, r, inf, foo)
	if res != nil {
		api.WriteRespAlertObj(w, r, tc.SuccessLevel, "Foo update was successful.", []tc.FooV18{*res})
	}
}

func updateV18(w http.ResponseWriter, r *http.Request, inf *api.APIInfo, foo tc.FooV18) *tc.FooV18 {
	fooLatest := tc.Foo(foo)
	id := inf.IntParams["id"]

	foo.ID = &id

	if err := fooLatest.Validate(inf.Tx.Tx); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("invalid request: "+err.Error()), nil)
		return nil
	}

	log.Infoln("here we would call tx.Query(updateQuery(), &foo.Name, &foo.A, &foo.B, &foo.C, &foo.D)")

	fooV18 := tc.FooV18(fooLatest)
	return &fooV18
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
