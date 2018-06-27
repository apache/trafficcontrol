package asn

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
	"fmt"
	"net/http"
	"strconv"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-tc/tovalidate"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// ASNsPrivLevel ...
const ASNsPrivLevel = 10

//we need a type alias to define functions on
type TOASNV11 struct {
	ReqInfo *api.APIInfo `json:"-"`
	tc.ASNNullable
}

func GetTypeSingleton() api.CRUDFactory {
	return func(reqInfo *api.APIInfo) api.CRUDer {
		toReturn := TOASNV11{reqInfo, tc.ASNNullable{}}
		return &toReturn
	}
}

func (asn TOASNV11) GetKeyFieldsInfo() []api.KeyFieldInfo {
	return []api.KeyFieldInfo{{"id", api.GetIntKey}}
}

// func (asn TOASNV12) GetKeyFieldsInfo() []api.KeyFieldInfo { return TOASNV11(asn).GetKeyFieldsInfo() }

//Implementation of the Identifier, Validator interface functions
func (asn TOASNV11) GetKeys() (map[string]interface{}, bool) {
	if asn.ID == nil {
		return map[string]interface{}{"id": 0}, false
	}
	return map[string]interface{}{"id": *asn.ID}, true
}

func (asn *TOASNV11) SetKeys(keys map[string]interface{}) {
	i, _ := keys["id"].(int) //this utilizes the non panicking type assertion, if the thrown away ok variable is false i will be the zero of the type, 0 here.
	asn.ID = &i
}

func (asn TOASNV11) GetAuditName() string {
	if asn.ASN != nil {
		return strconv.Itoa(*asn.ASN)
	}
	if asn.ID != nil {
		return strconv.Itoa(*asn.ID)
	}
	return "unknown"
}

func (asn TOASNV11) GetType() string {
	return "asn"
}

func (asn TOASNV11) Validate() []error {
	errs := validation.Errors{
		"asn":          validation.Validate(asn.ASN, validation.NotNil, validation.Min(0)),
		"cachegroupId": validation.Validate(asn.CachegroupID, validation.NotNil, validation.Min(0)),
	}
	return tovalidate.ToErrors(errs)
}

//The TOASNV11 implementation of the Creator interface
//all implementations of Creator should use transactions and return the proper errorType
//ParsePQUniqueConstraintError is used to determine if a asn with conflicting values exists
//if so, it will return an errorType of DataConflict and the type should be appended to the
//generic error message returned
//The insert sql returns the id and lastUpdated values of the newly inserted asn and have
//to be added to the struct
func (asn *TOASNV11) Create() (error, tc.ApiErrorType) {
	resultRows, err := asn.ReqInfo.Tx.NamedQuery(insertQuery(), asn)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			err, eType := dbhelpers.ParsePQUniqueConstraintError(pqErr)
			if eType == tc.DataConflictError {
				return errors.New("a asn with " + err.Error()), eType
			}
			return err, eType
		} else {
			log.Errorf("received non pq error: %++v from create execution", err)
			return tc.DBError, tc.SystemError
		}
	}
	defer resultRows.Close()

	var id int
	var lastUpdated tc.TimeNoMod
	rowsAffected := 0
	for resultRows.Next() {
		rowsAffected++
		if err := resultRows.Scan(&id, &lastUpdated); err != nil {
			log.Error.Printf("could not scan id from insert: %s\n", err)
			return tc.DBError, tc.SystemError
		}
	}
	if rowsAffected == 0 {
		err = errors.New("no asn was inserted, no id was returned")
		log.Errorln(err)
		return tc.DBError, tc.SystemError
	} else if rowsAffected > 1 {
		err = errors.New("too many ids returned from asn insert")
		log.Errorln(err)
		return tc.DBError, tc.SystemError
	}
	asn.SetKeys(map[string]interface{}{"id": id})
	asn.LastUpdated = &lastUpdated
	return nil, tc.NoError
}

// Read implements the /api/1.1/asns/id route for reading individual ASNs.
// Note this does NOT correctly implement the 1.1 API for all ASNs, because that route is in a different format than the CRUD utilities and all other routes.
// The /api/1.1/asns route MUST call V11ReadAll, not this function, to correctly implement the 1.1 API.
func (asn *TOASNV11) Read(parameters map[string]string) ([]interface{}, []error, tc.ApiErrorType) {
	asns, err, errType := read(asn.ReqInfo.Tx, parameters)
	if len(err) > 0 {
		return nil, err, errType
	}
	*asn.ReqInfo.CommitTx = true
	iasns := make([]interface{}, len(asns), len(asns))
	for i, readASN := range asns {
		iasns[i] = readASN
	}
	return iasns, err, errType
}

// V11ReadAll implements the asns 1.1 route, which is different from the 1.1 route for a single ASN and from 1.2+ routes, in that it wraps the content in an additional "asns" object.
func V11ReadAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleErrs := tc.GetHandleErrorsFunc(w, r)

		inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
		if userErr != nil || sysErr != nil {
			api.HandleErr(w, r, errCode, userErr, sysErr)
			return
		}
		defer inf.Close()

		params, err := api.GetCombinedParams(r)
		if err != nil {
			handleErrs(http.StatusInternalServerError, err)
			return
		}

		asns, errs, errType := read(inf.Tx, params)
		if len(errs) > 0 {
			tc.HandleErrorsWithType(errs, errType, handleErrs)
			return
		}
		*inf.CommitTx = true
		resp := struct {
			Response struct {
				ASNs []tc.ASNNullable `json:"asns"`
			} `json:"response"`
		}{Response: struct {
			ASNs []tc.ASNNullable `json:"asns"`
		}{ASNs: asns}}

		respBts, err := json.Marshal(resp)
		if err != nil {
			handleErrs(http.StatusInternalServerError, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s", respBts)
	}
}

func read(tx *sqlx.Tx, parameters map[string]string) ([]tc.ASNNullable, []error, tc.ApiErrorType) {
	// Query Parameters to Database Query column mappings
	// see the fields mapped in the SQL query
	queryParamsToQueryCols := map[string]dbhelpers.WhereColumnInfo{
		"asn":            dbhelpers.WhereColumnInfo{"a.asn", nil},
		"cachegroup":     dbhelpers.WhereColumnInfo{"c.id", nil},
		"id":             dbhelpers.WhereColumnInfo{"a.id", api.IsInt},
		"cachegroupName": dbhelpers.WhereColumnInfo{"c.name", nil},
	}
	where, orderBy, queryValues, errs := dbhelpers.BuildWhereAndOrderBy(parameters, queryParamsToQueryCols)
	if len(errs) > 0 {
		return nil, errs, tc.DataConflictError
	}

	query := selectQuery() + where + orderBy
	log.Debugln("Query is ", query)

	rows, err := tx.NamedQuery(query, queryValues)
	if err != nil {
		log.Errorf("Error querying ASNs: %v", err)
		return nil, []error{err}, tc.SystemError
	}
	defer rows.Close()

	ASNs := []tc.ASNNullable{}
	for rows.Next() {
		var s tc.ASNNullable
		if err = rows.StructScan(&s); err != nil {
			log.Errorf("error parsing ASN rows: %v", err)
			return nil, []error{err}, tc.SystemError
		}
		ASNs = append(ASNs, s)
	}

	return ASNs, []error{}, tc.NoError
}

func selectQuery() string {
	query := `SELECT
a.id,
a.asn,
a.last_updated,
a.cachegroup AS cachegroup_id,
c.name AS cachegroup

FROM asn a JOIN cachegroup c ON a.cachegroup = c.id`
	return query
}

//The TOASNV11 implementation of the Updater interface
//all implementations of Updater should use transactions and return the proper errorType
//ParsePQUniqueConstraintError is used to determine if a asn with conflicting values exists
//if so, it will return an errorType of DataConflict and the type should be appended to the
//generic error message returned
func (asn *TOASNV11) Update() (error, tc.ApiErrorType) {
	log.Debugf("about to run exec query: %s with asn: %++v", updateQuery(), asn)
	resultRows, err := asn.ReqInfo.Tx.NamedQuery(updateQuery(), asn)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			err, eType := dbhelpers.ParsePQUniqueConstraintError(pqErr)
			if eType == tc.DataConflictError {
				return errors.New("a asn with " + err.Error()), eType
			}
			return err, eType
		} else {
			log.Errorf("received error: %++v from update execution", err)
			return tc.DBError, tc.SystemError
		}
	}
	defer resultRows.Close()

	var lastUpdated tc.TimeNoMod
	rowsAffected := 0
	for resultRows.Next() {
		rowsAffected++
		if err := resultRows.Scan(&lastUpdated); err != nil {
			log.Error.Printf("could not scan lastUpdated from insert: %s\n", err)
			return tc.DBError, tc.SystemError
		}
	}
	log.Debugf("lastUpdated: %++v", lastUpdated)
	asn.LastUpdated = &lastUpdated
	if rowsAffected != 1 {
		if rowsAffected < 1 {
			return errors.New("no asn found with this id"), tc.DataMissingError
		} else {
			return fmt.Errorf("this update affected too many rows: %d", rowsAffected), tc.SystemError
		}
	}

	return nil, tc.NoError
}

//The ASN implementation of the Deleter interface
//all implementations of Deleter should use transactions and return the proper errorType
func (asn *TOASNV11) Delete() (error, tc.ApiErrorType) {
	log.Debugf("about to run exec query: %s with asn: %++v", deleteQuery(), asn)
	result, err := asn.ReqInfo.Tx.NamedExec(deleteQuery(), asn)
	if err != nil {
		log.Errorf("received error: %++v from delete execution", err)
		return tc.DBError, tc.SystemError
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return tc.DBError, tc.SystemError
	}
	if rowsAffected != 1 {
		if rowsAffected < 1 {
			return errors.New("no asn with that id found"), tc.DataMissingError
		} else {
			return fmt.Errorf("this create affected too many rows: %d", rowsAffected), tc.SystemError
		}
	}

	return nil, tc.NoError
}

func insertQuery() string {
	query := `INSERT INTO asn (
asn,
cachegroup) 
VALUES (
:asn,
:cachegroup_id
)
RETURNING id,last_updated`
	return query
}

func updateQuery() string {
	query := `UPDATE
asn SET
asn=:asn,
cachegroup=:cachegroup_id
WHERE id=:id RETURNING last_updated`
	return query
}

func deleteQuery() string {
	query := `DELETE FROM asn
WHERE id=:id`
	return query
}
