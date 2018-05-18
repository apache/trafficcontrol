package ip

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
	"strconv"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
	"github.com/apache/incubator-trafficcontrol/lib/go-tc/v13"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/tovalidate"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/utils"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

//we need a type alias to define functions on
type TOIP v13.IPNullable

//the refType is passed into the handlers where a copy of its type is used to decode the json.
var refType = TOIP{}

func GetRefType() *TOIP {
	return &refType
}

func (ip TOIP) GetKeyFieldsInfo() []api.KeyFieldInfo {
	return []api.KeyFieldInfo{{"id", api.GetIntKey}}
}

//Implementation of the Identifier, Validator interface functions
func (ip TOIP) GetKeys() (map[string]interface{}, bool) {
	if ip.ID == nil {
		return map[string]interface{}{"id": 0}, false
	}
	return map[string]interface{}{"id": *ip.ID}, true
}

func (ip *TOIP) SetKeys(keys map[string]interface{}) {
	i, _ := keys["id"].(int) //this utilizes the non panicking type assertion, if the thrown away ok variable is false i will be the zero of the type, 0 here.
	ip.ID = &i
}

//Implementation of the Identifier, Validator interface functions
func (ip *TOIP) GetID() (int, bool) {
	if ip.ID == nil {
		return 0, false
	}
	return *ip.ID, true
}

func (ip *TOIP) GetAuditName() string {
	id, _ := ip.GetID()
	return strconv.Itoa(id)
}

func (ip *TOIP) GetType() string {
	return "ip"
}

func (ip *TOIP) SetID(i int) {
	ip.ID = &i
}

func (ip *TOIP) Validate(db *sqlx.DB) []error {
	isIPv4Mask := validation.NewStringRule(tovalidate.IsIPv4Mask, "must be a valid IPv4 Mask")
	isIPv6CIDR := validation.NewStringRule(tovalidate.IsIPv6CIDR, "must be a valid IPv6 CIDR address")

	validateErrs := validation.Errors{
		"serverId":   validation.Validate(ip.ServerID, validation.NotNil),
		"typeId":     validation.Validate(ip.TypeID, validation.NotNil),
		"ipAddress":  validation.Validate(ip.IPAddress, is.IPv4),
		"ipNetmask":  validation.Validate(ip.IPNetmask, isIPv4Mask),
		"ipGateway":  validation.Validate(ip.IPGateway, is.IPv4),
		"ip6Address": validation.Validate(ip.IP6Address, isIPv6CIDR),
		"ip6Gateway": validation.Validate(ip.IP6Gateway, is.IPv6),
	}
	errs := tovalidate.ToErrors(validateErrs)
	if len(errs) > 0 {
		return errs
	}

	rows, err := db.Query("select id from server where id=$1", ip.ServerID)
	if err != nil {
		log.Error.Printf("could not execute select id from server: %s\n", err)
		errs = append(errs, tc.DBError)
		return errs
	}
	defer rows.Close()
	if !rows.Next() {
		errs = append(errs, errors.New("invalid server id"))
	}

	if ip.InterfaceID != nil && *ip.InterfaceID != 0 {
		rows, err := db.Query("select id from interface where id=$1", ip.InterfaceID)
		if err != nil {
			log.Error.Printf("could not execute select id from interface: %s\n", err)
			errs = append(errs, tc.DBError)
			return errs
		}
		defer rows.Close()
		if !rows.Next() {
			errs = append(errs, errors.New("invalid interface id"))
		}
	}

	rows, err = db.Query("select use_in_table from type where id=$1", ip.TypeID)
	if err != nil {
		log.Error.Printf("could not execute select use_in_table from type: %s\n", err)
		errs = append(errs, tc.DBError)
		return errs
	}
	defer rows.Close()
	var useInTable string
	for rows.Next() {
		if err := rows.Scan(&useInTable); err != nil {
			log.Error.Printf("could not scan use_in_table from type: %s\n", err)
			errs = append(errs, tc.DBError)
			return errs
		}
	}
	if useInTable != "ip" {
		errs = append(errs, errors.New("invalid ip type"))
	}

	return errs
}

func (ip *TOIP) Read(db *sqlx.DB, params map[string]string, user auth.CurrentUser) ([]interface{}, []error, tc.ApiErrorType) {
	var rows *sqlx.Rows
	var err error

	// Query Parameters to Database Query column mappings
	// see the fields mapped in the SQL query
	queryParamsToSQLCols := map[string]dbhelpers.WhereColumnInfo{
		"serverId": dbhelpers.WhereColumnInfo{"ip.server", api.IsInt},
		"id":       dbhelpers.WhereColumnInfo{"ip.id", api.IsInt},
	}

	where, orderBy, queryValues, errs := dbhelpers.BuildWhereAndOrderBy(params, queryParamsToSQLCols)
	if len(errs) > 0 {
		return nil, errs, tc.DataConflictError
	}

	query := SelectQuery() + where + orderBy
	log.Debugln("Query is ", query)

	rows, err = db.NamedQuery(query, queryValues)
	if err != nil {
		return nil, []error{fmt.Errorf("querying: %v", err)}, tc.SystemError
	}
	defer rows.Close()

	ips := []interface{}{}

	for rows.Next() {
		var s v13.IPNullable
		if err = rows.StructScan(&s); err != nil {
			return nil, []error{fmt.Errorf("getting IPs: %v", err)}, tc.SystemError
		}
		ips = append(ips, s)
	}
	return ips, nil, tc.NoError
}

func SelectQuery() string {
	selectStmt := `SELECT
ip.id,
ip.server as server_id,
s.host_name as server,
ip.type as type_id,
t.name as type,
ip.interface as interface_id,
if.interface_name as interface,
ip.ipv6,
ip.ipv6_gateway,
host(ip.ipv4) as ip_address,
netmask(ip.ipv4) as ip_netmask,
ip.ipv4_gateway,
ip.last_updated

FROM ip ip

JOIN type t ON ip.type = t.id
JOIN server s ON ip.server = s.id
LEFT JOIN interface if ON ip.interface = if.id`

	return selectStmt
}

//The TOIP implementation of the Updater interface
//all implementations of Updater should use transactions and return the proper errorType
//ParsePQUniqueConstraintError is used to determine if a cdn with conflicting values exists
//if so, it will return an errorType of DataConflict and the type should be appended to the
//generic error message returned
func (ip *TOIP) Update(db *sqlx.DB, user auth.CurrentUser) (error, tc.ApiErrorType) {
	// not allowed to change the ip type
	rows, err := db.Query("select type from ip where id=$1", ip.ID)
	if err != nil {
		log.Error.Printf("could not execute select type from ip: %s\n", err)
		return tc.DBError, tc.SystemError
	}
	defer rows.Close()
	var typeId int
	for rows.Next() {
		if err := rows.Scan(&typeId); err != nil {
			log.Error.Printf("could not scan type from ip: %s\n", err)
			return tc.DBError, tc.SystemError
		}
	}
	if typeId != *ip.TypeID {
		log.Error.Printf("not allowed to change the ip type\n")
		return errors.New("not allowed to change the ip type"), tc.ForbiddenError
	}

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

	err, errType := ip.UpdateExecAndCheck(tx)
	if err != nil {
		return err, errType
	}

	err = tx.Commit()
	if err != nil {
		log.Errorln("Could not commit transaction: ", err)
		return tc.DBError, tc.SystemError
	}
	rollbackTransaction = false
	return nil, tc.NoError
}

func (ip *TOIP) UpdateExecAndCheck(tx *sqlx.Tx) (error, tc.ApiErrorType) {
	ip.ConvertFormatApiToDb()

	log.Debugf("about to run exec query: %s with ip: %++v", updateQuery(), ip)
	resultRows, err := tx.NamedQuery(updateQuery(), ip)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			err, eType := dbhelpers.ParsePQUniqueConstraintError(pqErr)
			if eType == tc.DataConflictError {
				return errors.New("an ip with " + err.Error()), eType
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
	ip.LastUpdated = &lastUpdated
	if rowsAffected != 1 {
		if rowsAffected < 1 {
			return errors.New("no ip found with this id"), tc.DataMissingError
		} else {
			return fmt.Errorf("this update affected too many rows: %d", rowsAffected), tc.SystemError
		}
	}

	return nil, tc.NoError
}

func updateQuery() string {
	query := `UPDATE
ip SET
interface=:interface_id,
ipv6=:ipv6,
ipv6_gateway=:ipv6_gateway,
ipv4=:ipv4,
ipv4_gateway=:ipv4_gateway
WHERE id=:id RETURNING last_updated`
	return query
}

//The TOIP implementation of the Inserter interface
//all implementations of Inserter should use transactions and return the proper errorType
//ParsePQUniqueConstraintError is used to determine if an ip with conflicting values exists
//if so, it will return an errorType of DataConflict and the type should be appended to the
//generic error message returned
//The insert sql returns the id and lastUpdated values of the newly inserted ip and have
//to be added to the struct
func (ip *TOIP) Create(db *sqlx.DB, user auth.CurrentUser) (error, tc.ApiErrorType) {
	// only secondary type ip is allowed, other type ip (primary, management and ilo) is NOT allowed
	rows, err := db.Query("select name from type where id=$1", ip.TypeID)
	if err != nil {
		log.Error.Printf("could not execute select name from type: %s\n", err)
		return tc.DBError, tc.SystemError
	}
	defer rows.Close()
	var typeName string
	for rows.Next() {
		if err := rows.Scan(&typeName); err != nil {
			log.Error.Printf("could not scan name from type: %s\n", err)
			return tc.DBError, tc.SystemError
		}
	}
	if typeName != "IP_SECONDARY" {
		log.Error.Printf("%s type ip can not be created by this API\n", typeName)
		return errors.New("only IP_SECONDARY type ip can be created by this API"), tc.ForbiddenError
	}

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

	err, errType := ip.InsertExecAndCheck(tx)
	if err != nil {
		return err, errType
	}

	err = tx.Commit()
	if err != nil {
		log.Errorln("Could not commit transaction: ", err)
		return tc.DBError, tc.SystemError
	}
	rollbackTransaction = false
	return nil, tc.NoError
}

func (ip *TOIP) InsertExecAndCheck(tx *sqlx.Tx) (error, tc.ApiErrorType) {
	ip.ConvertFormatApiToDb()

	resultRows, err := tx.NamedQuery(insertQuery(), ip)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			err, eType := dbhelpers.ParsePQUniqueConstraintError(pqErr)
			if eType == tc.DataConflictError {
				return errors.New("an ip with " + err.Error()), eType
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
		err = errors.New("no ip was inserted, no id was returned")
		log.Errorln(err)
		return tc.DBError, tc.SystemError
	} else if rowsAffected > 1 {
		err = errors.New("too many ids returned from ip insert")
		log.Errorln(err)
		return tc.DBError, tc.SystemError
	}
	ip.SetID(id)
	ip.LastUpdated = &lastUpdated
	return nil, tc.NoError
}

func insertQuery() string {
	query := `INSERT INTO ip (
server,
type,
interface,
ipv6,
ipv6_gateway,
ipv4,
ipv4_gateway) VALUES (
:server_id,
:type_id,
:interface_id,
:ipv6,
:ipv6_gateway,
:ipv4,
:ipv4_gateway) RETURNING id,last_updated`
	return query
}

//The TOIP implementation of the Deleter interface
//all implementations of Deleter should use transactions and return the proper errorType
func (ip *TOIP) Delete(db *sqlx.DB, user auth.CurrentUser) (error, tc.ApiErrorType) {
	// only secondary type ip is allowed, other type ip (primary, management and ilo) is NOT allowed
	rows, err := db.Query("select t.name from ip ip join type t on ip.type=t.id where ip.id=$1", ip.ID)
	if err != nil {
		log.Error.Printf("could not execute select t.name from ip ip join type t on ip.type=t.id: %s\n", err)
		return tc.DBError, tc.SystemError
	}
	defer rows.Close()
	var typeName string
	for rows.Next() {
		if err := rows.Scan(&typeName); err != nil {
			log.Error.Printf("could not scan t.name from ip ip join type t on ip.type=t.id: %s\n", err)
			return tc.DBError, tc.SystemError
		}
	}
	if typeName != "IP_SECONDARY" {
		log.Error.Printf("%s type ip can not be deleted by this API\n", typeName)
		return errors.New("only IP_SECONDARY type ip can be deleted by this API"), tc.ForbiddenError
	}

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

	err, errType := ip.DeleteExecAndCheck(tx)
	if err != nil {
		return err, errType
	}

	err = tx.Commit()
	if err != nil {
		log.Errorln("Could not commit transaction: ", err)
		return tc.DBError, tc.SystemError
	}
	rollbackTransaction = false
	return nil, tc.NoError
}

func (ip *TOIP) DeleteExecAndCheck(tx *sqlx.Tx) (error, tc.ApiErrorType) {
	log.Debugf("about to run exec query: %s with ip: %++v", deleteQuery(), ip)
	result, err := tx.NamedExec(deleteQuery(), ip)
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
			return errors.New("no ip with that id found"), tc.DataMissingError
		} else {
			return fmt.Errorf("this create affected too many rows: %d", rowsAffected), tc.SystemError
		}
	}

	return nil, tc.NoError
}

func deleteQuery() string {
	query := `DELETE FROM ip
WHERE id=:id`
	return query
}

func (ip *TOIP) ConvertFormatApiToDb() error {
	if ip.IPAddress != nil && *ip.IPAddress == "" {
		ip.IPAddress = nil
	}
	if ip.IPNetmask != nil && *ip.IPNetmask == "" {
		ip.IPNetmask = nil
	}
	if ip.IPGateway != nil && *ip.IPGateway == "" {
		ip.IPGateway = nil
	}
	if ip.IPAddress == nil {
		ip.IPAddressNetmask = nil
	} else {
		leadingOnes := utils.GetIPv4MaskLeadingOnes(*ip.IPNetmask)
		ipAddressNetmask := fmt.Sprintf("%s/%d", *ip.IPAddress, leadingOnes)
		ip.IPAddressNetmask = &ipAddressNetmask
	}

	if ip.IP6Address != nil && *ip.IP6Address == "" {
		ip.IP6Address = nil
	}
	if ip.IP6Gateway != nil && *ip.IP6Gateway == "" {
		ip.IP6Gateway = nil
	}

	if ip.InterfaceID != nil && *ip.InterfaceID == 0 {
		ip.InterfaceID = nil
	}

	return nil
}
