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
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc/tovalidate"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/util/ims"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type TOStaticDNSEntry struct {
	api.APIInfoImpl `json:"-"`
	tc.StaticDNSEntryNullable
}

func (v *TOStaticDNSEntry) GetLastUpdated() (*time.Time, bool, error) {
	return api.GetLastUpdated(v.APIInfo().Tx, *v.ID, "staticdnsentry")
}

func (v *TOStaticDNSEntry) SetLastUpdated(t tc.TimeNoMod) { v.LastUpdated = &t }
func (v *TOStaticDNSEntry) InsertQuery() string           { return insertQuery() }
func (v *TOStaticDNSEntry) NewReadObj() interface{}       { return &tc.StaticDNSEntryNullable{} }
func (v *TOStaticDNSEntry) SelectQuery() string           { return selectQuery() }
func (v *TOStaticDNSEntry) ParamColumns() map[string]dbhelpers.WhereColumnInfo {
	return map[string]dbhelpers.WhereColumnInfo{
		"address":           dbhelpers.WhereColumnInfo{Column: "sde.address"},
		"cachegroup":        dbhelpers.WhereColumnInfo{Column: "cg.name"},
		"cachegroupId":      dbhelpers.WhereColumnInfo{Column: "cg.id"},
		"deliveryservice":   dbhelpers.WhereColumnInfo{Column: "ds.xml_id"},
		"deliveryserviceId": dbhelpers.WhereColumnInfo{Column: "sde.deliveryservice"},
		"host":              dbhelpers.WhereColumnInfo{Column: "sde.host"},
		"id":                dbhelpers.WhereColumnInfo{Column: "sde.id"},
		"ttl":               dbhelpers.WhereColumnInfo{Column: "sde.ttl"},
		"type":              dbhelpers.WhereColumnInfo{Column: "tp.name"},
		"typeId":            dbhelpers.WhereColumnInfo{Column: "tp.id"},
	}
}
func (v *TOStaticDNSEntry) UpdateQuery() string { return updateQuery() }
func (v *TOStaticDNSEntry) DeleteQuery() string { return deleteQuery() }

func (staticDNSEntry TOStaticDNSEntry) GetKeyFieldsInfo() []api.KeyFieldInfo {
	return []api.KeyFieldInfo{{Field: "id", Func: api.GetIntKey}}
}

// Implementation of the Identifier, Validator interface functions
func (staticDNSEntry TOStaticDNSEntry) GetKeys() (map[string]interface{}, bool) {
	if staticDNSEntry.ID == nil {
		return map[string]interface{}{"id": 0}, false
	}
	return map[string]interface{}{"id": *staticDNSEntry.ID}, true
}

func (staticDNSEntry TOStaticDNSEntry) GetAuditName() string {
	if staticDNSEntry.Host != nil {
		return *staticDNSEntry.Host
	}
	if staticDNSEntry.ID != nil {
		return strconv.Itoa(*staticDNSEntry.ID)
	}
	return "0"
}

func (staticDNSEntry TOStaticDNSEntry) GetType() string { return "staticDNSEntry" }

func (staticDNSEntry *TOStaticDNSEntry) SetKeys(keys map[string]interface{}) {
	i, _ := keys["id"].(int) //this utilizes the non panicking type assertion, if the thrown away ok variable is false i will be the zero of the type, 0 here.
	staticDNSEntry.ID = &i
}

// Validate fulfills the api.Validator interface.
func (staticDNSEntry TOStaticDNSEntry) Validate() (error, error) {
	typeStr, err := tc.ValidateTypeID(staticDNSEntry.ReqInfo.Tx.Tx, &staticDNSEntry.TypeID, "staticdnsentry")
	if err != nil {
		return err, nil
	}

	var addressErr, ttlErr error
	switch typeStr {
	case "A_RECORD":
		addressErr = validation.Validate(staticDNSEntry.Address, validation.Required, is.IPv4)
	case "AAAA_RECORD":
		addressErr = validation.Validate(staticDNSEntry.Address, validation.Required, is.IPv6)
	case "CNAME_RECORD":
		addressErr = validation.Validate(staticDNSEntry.Address, validation.Required, is.DNSName)
		address := *staticDNSEntry.Address
		if addressErr == nil {
			lastChar := address[len(address)-1:]
			if lastChar != "." {
				addressErr = fmt.Errorf("for type: CNAME_RECORD must have a trailing period")
			}
		}
	default:
		addressErr = validation.Validate(staticDNSEntry.Address, validation.Required)
	}

	if staticDNSEntry.TTL != nil {
		if *staticDNSEntry.TTL == 0 {
			ttlErr = validation.Validate(staticDNSEntry.TTL, is.Digit)
		}
	} else {
		ttlErr = validation.Validate(staticDNSEntry.TTL, validation.Required)
	}

	errs := validation.Errors{
		"host":              validation.Validate(staticDNSEntry.Host, validation.Required, is.DNSName),
		"address":           addressErr,
		"deliveryserviceId": validation.Validate(staticDNSEntry.DeliveryServiceID, validation.Required),
		"ttl":               ttlErr,
		"typeId":            validation.Validate(staticDNSEntry.TypeID, validation.Required),
	}
	return util.JoinErrs(tovalidate.ToErrors(errs)), nil
}

func (en *TOStaticDNSEntry) Read(h http.Header, useIMS bool) ([]interface{}, error, error, int, *time.Time) {
	api.DefaultSort(en.APIInfo(), "host")
	return api.GenericRead(h, en, useIMS)
}
func (en *TOStaticDNSEntry) Create() (error, error, int) {
	var cdnName tc.CDNName
	var err error
	if en.DeliveryServiceID != nil {
		_, cdnName, _, err = dbhelpers.GetDSNameAndCDNFromID(en.ReqInfo.Tx.Tx, *en.DeliveryServiceID)
		if err != nil {
			return nil, err, http.StatusInternalServerError
		}
		userErr, sysErr, errCode := dbhelpers.CheckIfCurrentUserCanModifyCDN(en.ReqInfo.Tx.Tx, string(cdnName), en.ReqInfo.User.UserName)
		if userErr != nil || sysErr != nil {
			return userErr, sysErr, errCode
		}
	}
	return api.GenericCreate(en)
}
func (en *TOStaticDNSEntry) Update(h http.Header) (error, error, int) {
	var cdnName tc.CDNName
	var err error
	if en.DeliveryServiceID != nil {
		_, cdnName, _, err = dbhelpers.GetDSNameAndCDNFromID(en.ReqInfo.Tx.Tx, *en.DeliveryServiceID)
		if err != nil {
			return nil, err, http.StatusInternalServerError
		}
		userErr, sysErr, errCode := dbhelpers.CheckIfCurrentUserCanModifyCDN(en.ReqInfo.Tx.Tx, string(cdnName), en.ReqInfo.User.UserName)
		if userErr != nil || sysErr != nil {
			return userErr, sysErr, errCode
		}
	}
	return api.GenericUpdate(h, en)
}
func (en *TOStaticDNSEntry) Delete() (error, error, int) {
	var cdnName tc.CDNName
	var err error
	var dsID int
	if en.DeliveryServiceID != nil {
		dsID = *en.DeliveryServiceID
	} else if en.ID != nil {
		dsID, err = dbhelpers.GetDSIDFromStaticDNSEntry(en.ReqInfo.Tx.Tx, *en.ID)
		if err != nil {
			return nil, errors.New("couldn't get DS ID from static dns entry ID: " + err.Error()), http.StatusInternalServerError
		}
	}
	_, cdnName, _, err = dbhelpers.GetDSNameAndCDNFromID(en.ReqInfo.Tx.Tx, dsID)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}
	userErr, sysErr, errCode := dbhelpers.CheckIfCurrentUserCanModifyCDN(en.ReqInfo.Tx.Tx, string(cdnName), en.ReqInfo.User.UserName)
	if userErr != nil || sysErr != nil {
		return userErr, sysErr, errCode
	}
	return api.GenericDelete(en)
}
func (v *TOStaticDNSEntry) SelectMaxLastUpdatedQuery(where, orderBy, pagination, tableName string) string {
	return `SELECT max(t) from (
		SELECT max(sde.last_updated) as t FROM staticdnsentry as sde
JOIN type as tp on sde.type = tp.id
LEFT JOIN cachegroup as cg ON sde.cachegroup = cg.id
JOIN deliveryservice as ds on sde.deliveryservice = ds.id ` + where + orderBy + pagination +
		` UNION ALL
	select max(last_updated) as t from last_deleted l where l.table_name='staticdnsentry') as res`
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
:deliveryservice_id,
:cachegroup_id,
:host,
:type_id,
:ttl) RETURNING id,last_updated`
	return query
}

func updateQuery() string {
	query := `UPDATE
staticdnsentry SET
id=:id,
address=:address,
deliveryservice=:deliveryservice_id,
cachegroup=:cachegroup_id,
host=:host,
type=:type_id,
ttl=:ttl
WHERE id=:id RETURNING last_updated`
	return query
}

func selectQuery() string {
	return `SELECT
ds.xml_id as dsname,
sde.host,
sde.id as id,
sde.deliveryservice as deliveryservice_id,
sde.ttl,
sde.address,
sde.last_updated,
tp.id as type_id,
tp.name as type,
cg.id as cachegroup_id,
cg.name as cachegroup
FROM staticdnsentry as sde
JOIN type as tp on sde.type = tp.id
LEFT JOIN cachegroup as cg ON sde.cachegroup = cg.id
JOIN deliveryservice as ds on sde.deliveryservice = ds.id
`
}

func deleteQuery() string {
	query := `DELETE FROM staticdnsentry
WHERE id=:id`
	return query
}

// Validate validates a statisDNSEntry entity and makes sure that all the supplied fields are valid.
func Validate(staticDNSEntry tc.StaticDNSEntryV5, tx *sql.Tx) (error, error) {
	typeStr, err := tc.ValidateTypeID(tx, staticDNSEntry.TypeID, "staticdnsentry")
	if err != nil {
		return err, nil
	}

	var addressErr, ttlErr error
	switch typeStr {
	case "A_RECORD":
		addressErr = validation.Validate(staticDNSEntry.Address, validation.Required, is.IPv4)
	case "AAAA_RECORD":
		addressErr = validation.Validate(staticDNSEntry.Address, validation.Required, is.IPv6)
	case "CNAME_RECORD":
		addressErr = validation.Validate(staticDNSEntry.Address, validation.Required, is.DNSName)
		address := *staticDNSEntry.Address
		if addressErr == nil {
			lastChar := address[len(address)-1:]
			if lastChar != "." {
				addressErr = fmt.Errorf("for type: CNAME_RECORD must have a trailing period")
			}
		}
	default:
		addressErr = validation.Validate(staticDNSEntry.Address, validation.Required)
	}

	if staticDNSEntry.TTL != nil {
		if *staticDNSEntry.TTL == 0 {
			ttlErr = validation.Validate(staticDNSEntry.TTL, is.Digit)
		}
	} else {
		ttlErr = validation.Validate(staticDNSEntry.TTL, validation.Required)
	}

	errs := validation.Errors{
		"host":              validation.Validate(staticDNSEntry.Host, validation.Required, is.DNSName),
		"address":           addressErr,
		"deliveryserviceId": validation.Validate(staticDNSEntry.DeliveryServiceID, validation.Required),
		"ttl":               ttlErr,
		"typeId":            validation.Validate(staticDNSEntry.TypeID, validation.Required),
	}
	return util.JoinErrs(tovalidate.ToErrors(errs)), nil
}

// Update will modify an existing StaticDNSEntry entity in the database, for api v5.0.
func Update(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	tx := inf.Tx.Tx

	defer r.Body.Close()
	var staticDNSEntry tc.StaticDNSEntryV5

	if err := json.NewDecoder(r.Body).Decode(&staticDNSEntry); err != nil {
		api.HandleErr(w, r, nil, http.StatusInternalServerError, nil, err)
		return
	}

	userErr, sysErr = Validate(staticDNSEntry, tx)
	if userErr != nil || sysErr != nil {
		code := http.StatusBadRequest
		if sysErr != nil {
			code = http.StatusInternalServerError
		}
		api.HandleErr(w, r, tx, code, userErr, sysErr)
		return
	}

	var cdnName tc.CDNName
	if id, ok := inf.Params["id"]; !ok {
		api.HandleErr(w, r, tx, http.StatusBadRequest, errors.New("missing key: id"), nil)
		return
	} else {
		idNum, err := strconv.Atoi(id)
		if err != nil {
			api.HandleErr(w, r, tx, http.StatusBadRequest, errors.New("couldn't convert ID into a numeric value: "+err.Error()), nil)
			return
		}
		staticDNSEntry.ID = &idNum
		if staticDNSEntry.DeliveryServiceID != nil {
			_, cdnName, _, err = dbhelpers.GetDSNameAndCDNFromID(tx, *staticDNSEntry.DeliveryServiceID)
			if err != nil {
				api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, err)
				return
			}
			userErr, sysErr, errCode := dbhelpers.CheckIfCurrentUserCanModifyCDN(tx, string(cdnName), inf.User.UserName)
			if userErr != nil || sysErr != nil {
				api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
				return
			}
		}

		existingLastUpdated, found, err := api.GetLastUpdated(inf.Tx, idNum, "staticdnsentry")
		if err == nil && found == false {
			api.HandleErr(w, r, tx, http.StatusNotFound, errors.New("no staticDNSEntry found with this id"), nil)
			return
		}
		if err != nil {
			api.HandleErr(w, r, tx, http.StatusNotFound, nil, err)
			return
		}
		if !api.IsUnmodified(r.Header, *existingLastUpdated) {
			api.HandleErr(w, r, tx, http.StatusPreconditionFailed, api.ResourceModifiedError, nil)
			return
		}

		rows, err := inf.Tx.NamedQuery(updateQuery(), staticDNSEntry)
		if err != nil {
			userErr, sysErr, errCode = api.ParseDBError(err)
			api.HandleErr(w, r, tx, errCode, userErr, sysErr)
			return
		}
		defer rows.Close()

		if !rows.Next() {
			api.HandleErr(w, r, tx, http.StatusNotFound, errors.New("no staticDNSEntry found with this id"), nil)
			return
		}
		lastUpdated := time.Time{}
		if err := rows.Scan(&lastUpdated); err != nil {
			api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, errors.New("scanning lastUpdated from staticDNSEntry insert: "+err.Error()))
			return
		}
		staticDNSEntry.LastUpdated = &lastUpdated
		if rows.Next() {
			api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, errors.New("staticDNSEntry update affected too many rows: >1"))
			return
		}

		alerts := tc.CreateAlerts(tc.SuccessLevel, "staticDNSEntry was updated.")
		api.WriteAlertsObj(w, r, http.StatusOK, alerts, staticDNSEntry)
		changeLogMsg := fmt.Sprintf("STATICDNSENTRY: %s, ID: %d, ACTION: Updated staticDNSEntry", *staticDNSEntry.Host, *staticDNSEntry.ID)
		api.CreateChangeLogRawTx(api.ApiChange, changeLogMsg, inf.User, tx)
		return
	}
}

// Create will add a new StaticDNSEntry entity into the database, for api v5.0.
func Create(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	tx := inf.Tx.Tx

	defer r.Body.Close()
	var staticDNSEntry tc.StaticDNSEntryV5

	if err := json.NewDecoder(r.Body).Decode(&staticDNSEntry); err != nil {
		api.HandleErr(w, r, nil, http.StatusInternalServerError, nil, err)
		return
	}

	userErr, sysErr = Validate(staticDNSEntry, tx)
	if userErr != nil || sysErr != nil {
		code := http.StatusBadRequest
		if sysErr != nil {
			code = http.StatusInternalServerError
		}
		api.HandleErr(w, r, tx, code, userErr, sysErr)
		return
	}

	var cdnName tc.CDNName
	var err error
	if staticDNSEntry.DeliveryServiceID != nil {
		_, cdnName, _, err = dbhelpers.GetDSNameAndCDNFromID(tx, *staticDNSEntry.DeliveryServiceID)
		if err != nil {
			api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, err)
			return
		}
		userErr, sysErr, errCode := dbhelpers.CheckIfCurrentUserCanModifyCDN(tx, string(cdnName), inf.User.UserName)
		if userErr != nil || sysErr != nil {
			api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
			return
		}
	}

	resultRows, err := inf.Tx.NamedQuery(insertQuery(), staticDNSEntry)
	if err != nil {
		userErr, sysErr, errCode = api.ParseDBError(err)
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	defer resultRows.Close()

	var id int
	lastUpdated := time.Time{}
	rowsAffected := 0
	for resultRows.Next() {
		rowsAffected++
		if err := resultRows.Scan(&id, &lastUpdated); err != nil {
			api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, errors.New("staticDNSEntry create scanning: "+err.Error()))
			return
		}
	}

	if rowsAffected == 0 {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, errors.New("staticDNSEntry create: no staticDNSEntry was inserted, no id was returned"))
		return
	} else if rowsAffected > 1 {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, errors.New("too many ids returned from staticDNSEntry insert"))
		return
	}
	staticDNSEntry.ID = &id
	staticDNSEntry.LastUpdated = &lastUpdated

	alerts := tc.CreateAlerts(tc.SuccessLevel, "staticDNSEntry was created.")
	api.WriteAlertsObj(w, r, http.StatusOK, alerts, staticDNSEntry)
	changeLogMsg := fmt.Sprintf("STATICDNSENTRY: %s, ID: %d, ACTION: Created staticDNSEntry", *staticDNSEntry.Host, *staticDNSEntry.ID)
	api.CreateChangeLogRawTx(api.ApiChange, changeLogMsg, inf.User, tx)
	return
}

// Delete removes a staticDNSEntry from the database, for api v5.0.
func Delete(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	tx := inf.Tx.Tx
	if id, ok := inf.Params["id"]; !ok {
		api.HandleErr(w, r, tx, http.StatusBadRequest, errors.New("missing key: id"), nil)
		return
	} else {
		idNum, err := strconv.Atoi(id)
		if err != nil {
			api.HandleErr(w, r, tx, http.StatusBadRequest, errors.New("couldn't convert ID into a numeric value: "+err.Error()), nil)
			return
		}
		staticDNSEntry := tc.StaticDNSEntryV5{
			ID: &idNum,
		}
		var cdnName tc.CDNName
		var dsID int
		if staticDNSEntry.DeliveryServiceID != nil {
			dsID = *staticDNSEntry.DeliveryServiceID
		} else if staticDNSEntry.ID != nil {
			dsID, err = dbhelpers.GetDSIDFromStaticDNSEntry(tx, *staticDNSEntry.ID)
			if err != nil {
				api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, errors.New("couldn't get DS ID from static dns entry ID: "+err.Error()))
				return
			}
		}
		_, cdnName, _, err = dbhelpers.GetDSNameAndCDNFromID(tx, dsID)
		if err != nil {
			api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, err)
			return
		}
		userErr, sysErr, errCode := dbhelpers.CheckIfCurrentUserCanModifyCDN(tx, string(cdnName), inf.User.UserName)
		if userErr != nil || sysErr != nil {
			api.HandleErr(w, r, tx, errCode, userErr, sysErr)
			return
		}
		result, err := inf.Tx.NamedExec(deleteQuery(), staticDNSEntry)
		if err != nil {
			userErr, sysErr, errCode = api.ParseDBError(err)
			api.HandleErr(w, r, tx, errCode, userErr, sysErr)
			return
		}

		if rowsAffected, err := result.RowsAffected(); err != nil {
			api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, errors.New("deleting staticDNSEntry: getting rows affected: "+err.Error()))
			return
		} else if rowsAffected < 1 {
			api.HandleErr(w, r, tx, http.StatusBadRequest, errors.New("no staticDNSEntry with that key found"), nil)
			return
		} else if rowsAffected > 1 {
			api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("staticDNSEntry delete affected too many rows: %d", rowsAffected))
			return
		}
		log.Debugf("changelog for delete on staticDNSEntry")

		alerts := tc.CreateAlerts(tc.SuccessLevel, "staticDNSEntry was deleted.")
		api.WriteAlerts(w, r, http.StatusOK, alerts)
		changeLogMsg := fmt.Sprintf("STATICDNSENTRY: %d, ID: %d, ACTION: Deleted staticDNSEntry", *staticDNSEntry.ID, *staticDNSEntry.ID)
		api.CreateChangeLogRawTx(api.ApiChange, changeLogMsg, inf.User, tx)
		return
	}
}

// Get : function to read the staticDNSEntries, for api version 5.0.
func Get(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	api.DefaultSort(inf, "host")
	useIMS := false
	config, e := api.GetConfig(r.Context())
	if e == nil && config != nil {
		useIMS = config.UseIMS
	} else {
		log.Warnf("Couldn't get config %v", e)
	}

	var maxTime time.Time
	var runSecond bool

	scList := make([]tc.StaticDNSEntryV5, 0)

	tx := inf.Tx

	where, orderBy, pagination, queryValues, errs := dbhelpers.BuildWhereAndOrderByAndPagination(inf.Params, map[string]dbhelpers.WhereColumnInfo{
		"address":           dbhelpers.WhereColumnInfo{Column: "sde.address"},
		"cachegroup":        dbhelpers.WhereColumnInfo{Column: "cg.name"},
		"cachegroupId":      dbhelpers.WhereColumnInfo{Column: "cg.id"},
		"deliveryservice":   dbhelpers.WhereColumnInfo{Column: "ds.xml_id"},
		"deliveryserviceId": dbhelpers.WhereColumnInfo{Column: "sde.deliveryservice"},
		"host":              dbhelpers.WhereColumnInfo{Column: "sde.host"},
		"id":                dbhelpers.WhereColumnInfo{Column: "sde.id"},
		"ttl":               dbhelpers.WhereColumnInfo{Column: "sde.ttl"},
		"type":              dbhelpers.WhereColumnInfo{Column: "tp.name"},
		"typeId":            dbhelpers.WhereColumnInfo{Column: "tp.id"},
	})
	if len(errs) > 0 {
		api.HandleErr(w, r, tx.Tx, http.StatusBadRequest, nil, util.JoinErrs(errs))
		return
	}
	if useIMS {
		runSecond, maxTime = ims.TryIfModifiedSinceQuery(inf.Tx, r.Header, queryValues, selectMaxLastUpdatedQuery(where))
		if !runSecond {
			log.Debugln("IMS HIT")
			api.AddLastModifiedHdr(w, maxTime)
			w.WriteHeader(http.StatusNotModified)
			return
		}
		log.Debugln("IMS MISS")
	} else {
		log.Debugln("Non IMS request")
	}

	// Case where we need to run the second query
	query := selectQuery() + where + orderBy + pagination
	rows, err := tx.NamedQuery(query, queryValues)
	if err != nil {
		api.HandleErr(w, r, tx.Tx, http.StatusInternalServerError, nil, errors.New("querying static DNS entries: "+err.Error()))
		return
	}
	defer rows.Close()

	for rows.Next() {
		v := tc.StaticDNSEntryV5{}
		if err = rows.StructScan(&v); err != nil {
			api.HandleErr(w, r, tx.Tx, http.StatusInternalServerError, nil, errors.New("scanning static DNS entries: "+err.Error()))
			return
		}
		scList = append(scList, v)
	}

	api.WriteResp(w, r, scList)
}

func selectMaxLastUpdatedQuery(where string) string {
	return `SELECT max(t) from (
		SELECT max(sde.last_updated) as t FROM staticdnsentry as sde
JOIN type as tp on sde.type = tp.id
LEFT JOIN cachegroup as cg ON sde.cachegroup = cg.id
JOIN deliveryservice as ds on sde.deliveryservice = ds.id ` + where +
		` UNION ALL
	select max(last_updated) as t from last_deleted l where l.table_name='staticdnsentry') as res`
}
