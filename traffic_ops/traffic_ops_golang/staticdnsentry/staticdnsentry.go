package staticdnsentry

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
	"errors"
	"fmt"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type TOStaticDNSEntry struct {
	ReqInfo *api.APIInfo `json:"-"`
	tc.StaticDNSEntry
}

func GetReaderSingleton() func(reqInfo *api.APIInfo) api.Reader {
	return func(reqInfo *api.APIInfo) api.Reader {
		toReturn := TOStaticDNSEntry{reqInfo, tc.StaticDNSEntry{}}
		return &toReturn
	}
}

func (staticDNSEntry *TOStaticDNSEntry) SetKeys(keys map[string]interface{}) {
	i, _ := keys["id"].(int) //this utilizes the non panicking type assertion, if the thrown away ok variable is false i will be the zero of the type, 0 here.
	staticDNSEntry.ID = &i
}

func (staticDNSEntry *TOStaticDNSEntry) Read(parameters map[string]string) ([]interface{}, []error, tc.ApiErrorType) {
	queryParamsToQueryCols := map[string]dbhelpers.WhereColumnInfo{
		"deliveryservice": dbhelpers.WhereColumnInfo{"deliveryservice", nil}, // order by
	}
	where, orderBy, queryValues, errs := dbhelpers.BuildWhereAndOrderBy(parameters, queryParamsToQueryCols)
	if len(errs) > 0 {
		log.Errorf("Data Conflict Error")
		return nil, errs, tc.DataConflictError
	}
	query := selectQuery() + where + orderBy
	log.Debugln("Query is ", query)
	rows, err := staticDNSEntry.ReqInfo.Tx.NamedQuery(query, queryValues)
	if err != nil {
		log.Errorf("Error querying StaticDNSEntries: %v", err)
		return nil, []error{tc.DBError}, tc.SystemError
	}
	defer rows.Close()
	staticDNSEntries := []interface{}{}
	for rows.Next() {
		s := tc.StaticDNSEntry{}
		if err = rows.StructScan(&s); err != nil {
			log.Errorln("error parsing StaticDNSEntry rows: " + err.Error())
			return nil, []error{tc.DBError}, tc.SystemError
		}
		staticDNSEntries = append(staticDNSEntries, s)
	}
	return staticDNSEntries, []error{}, tc.NoError
}

//The TOStaticDNSEntry implementation of the Creator interface
//all implementations of Creator should use transactions and return the proper errorType
//ParsePQUniqueConstraintError is used to determine if a staticDNSEntry with conflicting values exists
//if so, it will return an errorType of DataConflict and the type should be appended to the
//generic error message returned
//The insert sql returns the id and lastUpdated values of the newly inserted staticDNSEntry and have
//to be added to the struct
func (staticDNSEntry *TOStaticDNSEntry) Create(db *sqlx.DB, user auth.CurrentUser) (error, tc.ApiErrorType) {
	rollbackTransaction := true
	tx, err := db.Beginx()
	defer func() {
		if tx == nil || !rollbackTransaction {
			return
		}
		err := tx.Rollback()
		if err != nil {
			log.Errorln(errors.New("rolling back transaction: " + err.Error()))
		}
	}()

	if err != nil {
		log.Error.Printf("could not begin transaction: %v", err)
		return tc.DBError, tc.SystemError
	}
	// make sure that staticDNSEntry.DomainName is lowercase
	resultRows, err := tx.NamedQuery(insertQuery(), staticDNSEntry)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			err, eType := dbhelpers.ParsePQUniqueConstraintError(pqErr)
			if eType == tc.DataConflictError {
				return errors.New("a staticDNSEntry with " + err.Error()), eType
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
		err = errors.New("no staticDNSEntry was inserted, no id was returned")
		log.Errorln(err)
		return tc.DBError, tc.SystemError
	} else if rowsAffected > 1 {
		err = errors.New("too many ids returned from staticDNSEntry insert")
		log.Errorln(err)
		return tc.DBError, tc.SystemError
	}
	staticDNSEntry.SetKeys(map[string]interface{}{"id": id})
	staticDNSEntry.LastUpdated = &lastUpdated
	err = tx.Commit()
	if err != nil {
		log.Errorln("Could not commit transaction: ", err)
		return tc.DBError, tc.SystemError
	}
	rollbackTransaction = false
	return nil, tc.NoError
}

func insertQuery() string {
	query := `INSERT INTO staticdnsentry (
address,
deliveryservice,
cachegroup,
host,
type,
ttl) VALUES (
:address,
:deliveryservice,
:cachegroup,
:host,
:type,
:ttl) RETURNING id,last_updated`
	return query
}

//The TOStaticDNSEntry implementation of the Updater interface
//all implementations of Updater should use transactions and return the proper errorType
//ParsePQUniqueConstraintError is used to determine if a staticDNSEntry with conflicting values exists
//if so, it will return an errorType of DataConflict and the type should be appended to the
//generic error message returned
func (staticDNSEntry *TOStaticDNSEntry) Update(db *sqlx.DB, user auth.CurrentUser) (error, tc.ApiErrorType) {
	rollbackTransaction := true
	tx, err := db.Beginx()
	defer func() {
		if tx == nil || !rollbackTransaction {
			return
		}
		err := tx.Rollback()
		if err != nil {
			log.Errorln(errors.New("rolling back transaction: " + err.Error()))
		}
	}()

	if err != nil {
		log.Error.Printf("could not begin transaction: %v", err)
		return tc.DBError, tc.SystemError
	}
	log.Debugf("about to run exec query: %s with staticDNSEntry: %++v", updateQuery(), staticDNSEntry)
	resultRows, err := tx.NamedQuery(updateQuery(), staticDNSEntry)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			err, eType := dbhelpers.ParsePQUniqueConstraintError(pqErr)
			if eType == tc.DataConflictError {
				return errors.New("a staticDNSEntry with " + err.Error()), eType
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
	staticDNSEntry.LastUpdated = &lastUpdated
	if rowsAffected != 1 {
		if rowsAffected < 1 {
			return errors.New("no staticDNSEntry found with this id"), tc.DataMissingError
		} else {
			return fmt.Errorf("this update affected too many rows: %d", rowsAffected), tc.SystemError
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Errorln("Could not commit transaction: ", err)
		return tc.DBError, tc.SystemError
	}
	rollbackTransaction = false
	return nil, tc.NoError
}

func updateQuery() string {
	query := `UPDATE
staticdnsentry SET
address=:address,
deliveryservice=:deliveryservice,
cachegroup=:cachegroup,
host=:host,
type=:type,
ttl=:ttl
WHERE id=:id RETURNING last_updated`
	return query
}

//The StaticDNSEntry implementation of the Deleter interface
//all implementations of Deleter should use transactions and return the proper errorType
func (staticDNSEntry *TOStaticDNSEntry) Delete(db *sqlx.DB, user auth.CurrentUser) (error, tc.ApiErrorType) {
	rollbackTransaction := true
	tx, err := db.Beginx()
	defer func() {
		if tx == nil || !rollbackTransaction {
			return
		}
		err := tx.Rollback()
		if err != nil {
			log.Errorln(errors.New("rolling back transaction: " + err.Error()))
		}
	}()

	if err != nil {
		log.Error.Printf("could not begin transaction: %v", err)
		return tc.DBError, tc.SystemError
	}
	log.Debugf("about to run exec query: %s with staticDNSEntry: %++v", deleteQuery(), staticDNSEntry)
	result, err := tx.NamedExec(deleteQuery(), staticDNSEntry)
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
			return errors.New("no staticDNSEntry with that id found"), tc.DataMissingError
		} else {
			return fmt.Errorf("this create affected too many rows: %d", rowsAffected), tc.SystemError
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Errorln("Could not commit transaction: ", err)
		return tc.DBError, tc.SystemError
	}
	rollbackTransaction = false
	return nil, tc.NoError
}

func selectQuery() string {
	return `SELECT
ds.xml_id as dsname,
sde.host,
sde.ttl,
sde.address,
tp.name as type,
cg.name as cachegroup
FROM staticdnsentry as sde
JOIN type as tp on sde.type = tp.id
JOIN cachegroup as cg ON sde.cachegroup = cg.id
JOIN deliveryservice as ds on sde.deliveryservice = ds.id
`
}

func deleteQuery() string {
	query := `DELETE FROM staticdnsentry
WHERE id=:id`
	return query
}
