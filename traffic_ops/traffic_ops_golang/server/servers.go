package server

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
	"github.com/apache/incubator-trafficcontrol/lib/go-tc/common"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/tovalidate"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

//we need a type alias to define functions on
type TOServer tc.ServerNullable

//the refType is passed into the handlers where a copy of its type is used to decode the json.
var refType = TOServer(tc.ServerNullable{})

func GetRefType() *TOServer {
	return &refType
}

//Implementation of the Identifier, Validator interface functions
func (server *TOServer) GetID() (int, bool) {
	if server.ID == nil {
		return 0, false
	}
	return *server.ID, true
}

func (server *TOServer) GetAuditName() string {
	if server.DomainName != nil {
		return *server.DomainName
	}
	id, _ := server.GetID()
	return strconv.Itoa(id)
}

func (server *TOServer) GetType() string {
	return "server"
}

func (server *TOServer) SetID(i int) {
	server.ID = &i
}

func (server *TOServer) Validate(db *sqlx.DB) []error {

	noSpaces := validation.NewStringRule(tovalidate.NoSpaces, "cannot contain spaces")

	validateErrs := validation.Errors{
		"cachegroupId":   validation.Validate(server.CachegroupID, validation.NotNil),
		"cdnId":          validation.Validate(server.CDNID, validation.NotNil),
		"domainName":     validation.Validate(server.DomainName, validation.NotNil, noSpaces),
		"hostName":       validation.Validate(server.HostName, validation.NotNil, noSpaces),
		"httpsPort":      validation.Validate(server.HTTPSPort, validation.NotNil),
		"interfaceMtu":   validation.Validate(server.InterfaceMtu, validation.NotNil),
		"interfaceName":  validation.Validate(server.InterfaceName, validation.NotNil),
		"ipAddress":      validation.Validate(server.IPAddress, validation.NotNil, is.IPv4),
		"ipNetmask":      validation.Validate(server.IPNetmask, validation.NotNil),
		"ipGateway":      validation.Validate(server.IPGateway, validation.NotNil),
		"physLocationId": validation.Validate(server.PhysLocationID, validation.NotNil),
		"profileId":      validation.Validate(server.ProfileID, validation.NotNil),
		"statusId":       validation.Validate(server.StatusID, validation.NotNil),
		"tcpPort":        validation.Validate(server.TCPPort, validation.NotNil),
		"typeId":         validation.Validate(server.TypeID, validation.NotNil),
		"updPending":     validation.Validate(server.UpdPending, validation.NotNil),
	}
	errs := tovalidate.ToErrors(validateErrs)
	if len(errs) > 0 {
		return errs
	}

	rows, err := db.Query("select use_in_table from type where id=$1", server.TypeID)
	if err != nil {
		log.Error.Printf("could not execute select use_in_table from type: %s\n", err)
		errs = append(errs, common.DBError)
		return errs
	}
	defer rows.Close()
	var useInTable string
	for rows.Next() {
		if err := rows.Scan(&useInTable); err != nil {
			log.Error.Printf("could not scan use_in_table from type: %s\n", err)
			errs = append(errs, common.DBError)
			return errs
		}
	}
	if useInTable != "server" {
		errs = append(errs, errors.New("invalid server type"))
	}

	rows, err = db.Query("select cdn from profile where id=$1", server.ProfileID)
	if err != nil {
		log.Error.Printf("could not execute select cdnID from profile: %s\n", err)
		errs = append(errs, common.DBError)
		return errs
	}
	defer rows.Close()
	var cdnID int
	for rows.Next() {
		if err := rows.Scan(&cdnID); err != nil {
			log.Error.Printf("could not scan cdnID from profile: %s\n", err)
			errs = append(errs, common.DBError)
			return errs
		}
	}
	log.Infof("got cdn id: %d from profile and cdn id: %d from server", cdnID, *server.CDNID)
	if cdnID != *server.CDNID {
		errs = append(errs, errors.New("CDN of profile does not match Server CDN"))
	}
	return errs
}

func (server *TOServer) Read(db *sqlx.DB, params map[string]string, user auth.CurrentUser) ([]interface{}, []error, common.ApiErrorType) {
	returnable := []interface{}{}

	privLevel := user.PrivLevel

	servers, errs, errType := getServers(params, db, privLevel)
	if len(errs) > 0 {
		for _, err := range errs {
			if err.Error() == `id cannot parse to integer` {
				return nil, []error{errors.New("Resource not found.")}, common.DataMissingError //matches perl response
			}
		}
		return nil, errs, errType
	}

	for _, server := range servers {
		returnable = append(returnable, server)
	}

	return returnable, nil, common.NoError
}

func getServers(params map[string]string, db *sqlx.DB, privLevel int) ([]TOServer, []error, common.ApiErrorType) {
	var rows *sqlx.Rows
	var err error

	// Query Parameters to Database Query column mappings
	// see the fields mapped in the SQL query
	queryParamsToSQLCols := map[string]dbhelpers.WhereColumnInfo{
		"cachegroup":   dbhelpers.WhereColumnInfo{"s.cachegroup", api.IsInt},
		"cdn":          dbhelpers.WhereColumnInfo{"s.cdn_id", api.IsInt},
		"id":           dbhelpers.WhereColumnInfo{"s.id", api.IsInt},
		"hostName":     dbhelpers.WhereColumnInfo{"s.host_name", nil},
		"physLocation": dbhelpers.WhereColumnInfo{"s.phys_location", api.IsInt},
		"profileId":    dbhelpers.WhereColumnInfo{"s.profile", api.IsInt},
		"status":       dbhelpers.WhereColumnInfo{"st.name", nil},
		"type":         dbhelpers.WhereColumnInfo{"t.name", nil},
	}

	where, orderBy, queryValues, errs := dbhelpers.BuildWhereAndOrderBy(params, queryParamsToSQLCols)
	if len(errs) > 0 {
		return nil, errs, common.DataConflictError
	}

	query := selectQuery() + where + orderBy
	log.Debugln("Query is ", query)

	rows, err = db.NamedQuery(query, queryValues)
	if err != nil {
		return nil, []error{fmt.Errorf("querying: %v", err)}, common.SystemError
	}
	defer rows.Close()

	servers := []TOServer{}

	HiddenField := "********"

	for rows.Next() {
		var s TOServer
		if err = rows.StructScan(&s); err != nil {
			return nil, []error{fmt.Errorf("getting servers: %v", err)}, common.SystemError
		}
		if privLevel < auth.PrivLevelAdmin {
			s.ILOPassword = &HiddenField
			s.XMPPPasswd = &HiddenField
		}
		servers = append(servers, s)
	}
	return servers, nil, common.NoError
}

func selectQuery() string {

	const JumboFrameBPS = 9000

	// COALESCE is needed to default values that are nil in the database
	// because Go does not allow that to marshal into the struct
	selectStmt := `SELECT
cg.name as cachegroup,
s.cachegroup as cachegroup_id,
s.cdn_id,
cdn.name as cdn_name,
s.domain_name,
s.guid,
s.host_name,
s.https_port,
s.id,
s.ilo_ip_address,
s.ilo_ip_gateway,
s.ilo_ip_netmask,
s.ilo_password,
s.ilo_username,
COALESCE(s.interface_mtu, ` + strconv.Itoa(JumboFrameBPS) + `) as interface_mtu,
s.interface_name,
s.ip6_address,
s.ip6_gateway,
s.ip_address,
s.ip_gateway,
s.ip_netmask,
s.last_updated,
s.mgmt_ip_address,
s.mgmt_ip_gateway,
s.mgmt_ip_netmask,
s.offline_reason,
pl.name as phys_location,
s.phys_location as phys_location_id,
p.name as profile,
p.description as profile_desc,
s.profile as profile_id,
s.rack,
s.reval_pending,
s.router_host_name,
s.router_port_name,
st.name as status,
s.status as status_id,
s.tcp_port,
t.name as server_type,
s.type as server_type_id,
s.upd_pending as upd_pending,
s.xmpp_id,
s.xmpp_passwd

FROM server s

JOIN cachegroup cg ON s.cachegroup = cg.id
JOIN cdn cdn ON s.cdn_id = cdn.id
JOIN phys_location pl ON s.phys_location = pl.id
JOIN profile p ON s.profile = p.id
JOIN status st ON s.status = st.id
JOIN type t ON s.type = t.id`

	return selectStmt
}

//The TOServer implementation of the Updater interface
//all implementations of Updater should use transactions and return the proper errorType
//ParsePQUniqueConstraintError is used to determine if a cdn with conflicting values exists
//if so, it will return an errorType of DataConflict and the type should be appended to the
//generic error message returned
func (server *TOServer) Update(db *sqlx.DB, user auth.CurrentUser) (error, common.ApiErrorType) {
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
		return common.DBError, common.SystemError
	}

	log.Debugf("about to run exec query: %s with server: %++v", updateQuery(), server)
	resultRows, err := tx.NamedQuery(updateQuery(), server)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			err, eType := dbhelpers.ParsePQUniqueConstraintError(pqErr)
			if eType == common.DataConflictError {
				return errors.New("a server with " + err.Error()), eType
			}
			return err, eType
		} else {
			log.Errorf("received error: %++v from update execution", err)
			return common.DBError, common.SystemError
		}
	}
	defer resultRows.Close()

	var lastUpdated common.TimeNoMod
	rowsAffected := 0
	for resultRows.Next() {
		rowsAffected++
		if err := resultRows.Scan(&lastUpdated); err != nil {
			log.Error.Printf("could not scan lastUpdated from insert: %s\n", err)
			return common.DBError, common.SystemError
		}
	}
	log.Debugf("lastUpdated: %++v", lastUpdated)
	server.LastUpdated = lastUpdated
	if rowsAffected != 1 {
		if rowsAffected < 1 {
			return errors.New("no cdn found with this id"), common.DataMissingError
		} else {
			return fmt.Errorf("this update affected too many rows: %d", rowsAffected), common.SystemError
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Errorln("Could not commit transaction: ", err)
		return common.DBError, common.SystemError
	}
	rollbackTransaction = false
	return nil, common.NoError
}

func updateQuery() string {
	query := `UPDATE
server SET
cachegroup=:cachegroup_id,
cdn_id=:cdn_id,
domain_name=:domain_name,
host_name=:host_name,
https_port=:https_port,
ilo_ip_address=:ilo_ip_address,
ilo_ip_netmask=:ilo_ip_netmask,
ilo_ip_gateway=:ilo_ip_gateway,
ilo_username=:ilo_username,
ilo_password=:ilo_password,
interface_mtu=:interface_mtu,
interface_name=:interface_name,
ip6_address=:ip6_address,
ip6_gateway=:ip6_gateway,
ip_address=:ip_address,
ip_netmask=:ip_netmask,
ip_gateway=:ip_gateway,
mgmt_ip_address=:mgmt_ip_address,
mgmt_ip_netmask=:mgmt_ip_netmask,
mgmt_ip_gateway=:mgmt_ip_gateway,
offline_reason=:offline_reason,
phys_location=:phys_location_id,
profile=:profile_id,
rack=:rack,
router_host_name=:router_host_name,
router_port_name=:router_port_name,
status=:status_id,
tcp_port=:tcp_port,
type=:server_type_id,
upd_pending=:upd_pending
WHERE id=:id RETURNING last_updated`
	return query
}

//The TOServer implementation of the Inserter interface
//all implementations of Inserter should use transactions and return the proper errorType
//ParsePQUniqueConstraintError is used to determine if a server with conflicting values exists
//if so, it will return an errorType of DataConflict and the type should be appended to the
//generic error message returned
//The insert sql returns the id and lastUpdated values of the newly inserted server and have
//to be added to the struct
func (server *TOServer) Create(db *sqlx.DB, user auth.CurrentUser) (error, common.ApiErrorType) {
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
		return common.DBError, common.SystemError
	}
	if server.XMPPID == nil || *server.XMPPID == "" {
		server.XMPPID = server.HostName
	}

	resultRows, err := tx.NamedQuery(insertQuery(), server)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			err, eType := dbhelpers.ParsePQUniqueConstraintError(pqErr)
			if eType == common.DataConflictError {
				return errors.New("a server with " + err.Error()), eType
			}
			return err, eType
		} else {
			log.Errorf("received non pq error: %++v from create execution", err)
			return common.DBError, common.SystemError
		}
	}
	defer resultRows.Close()

	var id int
	var lastUpdated common.TimeNoMod
	rowsAffected := 0
	for resultRows.Next() {
		rowsAffected++
		if err := resultRows.Scan(&id, &lastUpdated); err != nil {
			log.Error.Printf("could not scan id from insert: %s\n", err)
			return common.DBError, common.SystemError
		}
	}
	if rowsAffected == 0 {
		err = errors.New("no server was inserted, no id was returned")
		log.Errorln(err)
		return common.DBError, common.SystemError
	} else if rowsAffected > 1 {
		err = errors.New("too many ids returned from server insert")
		log.Errorln(err)
		return common.DBError, common.SystemError
	}
	server.SetID(id)
	server.LastUpdated = lastUpdated
	err = tx.Commit()
	if err != nil {
		log.Errorln("Could not commit transaction: ", err)
		return common.DBError, common.SystemError
	}
	rollbackTransaction = false
	return nil, common.NoError
}

func insertQuery() string {
	query := `INSERT INTO server (
cachegroup,
cdn_id,
domain_name,
host_name,
https_port,
ilo_ip_address,
ilo_ip_netmask,
ilo_ip_gateway,
ilo_username,
ilo_password,
interface_mtu,
interface_name,
ip6_address,
ip6_gateway,
ip_address,
ip_netmask,
ip_gateway,
mgmt_ip_address,
mgmt_ip_netmask,
mgmt_ip_gateway,
offline_reason,
phys_location,
profile,
rack,
router_host_name,
router_port_name,
status,
tcp_port,
type,
upd_pending) VALUES (
:cachegroup_id,
:cdn_id,
:domain_name,
:host_name,
:https_port,
:ilo_ip_address,
:ilo_ip_netmask,
:ilo_ip_gateway,
:ilo_username,
:ilo_password,
:interface_mtu,
:interface_name,
:ip6_address,
:ip6_gateway,
:ip_address,
:ip_netmask,
:ip_gateway,
:mgmt_ip_address,
:mgmt_ip_netmask,
:mgmt_ip_gateway,
:offline_reason,
:phys_location_id,
:profile_id,
:rack,
:router_host_name,
:router_port_name,
:status_id,
:tcp_port,
:server_type_id,
:upd_pending) RETURNING id,last_updated`
	return query
}

//The Server implementation of the Deleter interface
//all implementations of Deleter should use transactions and return the proper errorType
func (server *TOServer) Delete(db *sqlx.DB, user auth.CurrentUser) (error, common.ApiErrorType) {
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
		return common.DBError, common.SystemError
	}
	log.Debugf("about to run exec query: %s with server: %++v", deleteServerQuery(), server)
	result, err := tx.NamedExec(deleteServerQuery(), server)
	if err != nil {
		log.Errorf("received error: %++v from delete execution", err)
		return common.DBError, common.SystemError
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return common.DBError, common.SystemError
	}
	if rowsAffected != 1 {
		if rowsAffected < 1 {
			return errors.New("no server with that id found"), common.DataMissingError
		} else {
			return fmt.Errorf("this create affected too many rows: %d", rowsAffected), common.SystemError
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Errorln("Could not commit transaction: ", err)
		return common.DBError, common.SystemError
	}
	rollbackTransaction = false
	return nil, common.NoError
}

func deleteServerQuery() string {
	query := `DELETE FROM server
WHERE id=:id`
	return query
}
