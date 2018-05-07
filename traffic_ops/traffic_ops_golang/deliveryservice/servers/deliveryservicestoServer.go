package servers

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
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/tovalidate"
	"github.com/go-ozzo/ozzo-validation"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// TODeliveryServiceRequest provides a type alias to define functions on
type TODeliveryServiceServer tc.DeliveryServiceServer

//the refType is passed into the handlers where a copy of its type is used to decode the json.
var refType = TODeliveryServiceServer(tc.DeliveryServiceServer{})

func GetRefType() *TODeliveryServiceServer {
	return &refType
}

/*
# get all delivery services associated with a server (from deliveryservice_server table)
$r->get( "/api/$version/servers/:id/deliveryservices" => [ id => qr/\d+/ ] )->over( authenticated => 1, not_ldap => 1 )->to( 'Deliveryservice#get_deliveryservices_by_serverId', namespace => $namespace );

# delivery service / server assignments
$r->post("/api/$version/deliveryservices/:xml_id/servers")->over( authenticated => 1, not_ldap => 1 )
->to( 'Deliveryservice#assign_servers', namespace => $namespace );
$r->delete("/api/$version/deliveryservice_server/:dsId/:serverId" => [ dsId => qr/\d+/, serverId => qr/\d+/ ] )->over( authenticated => 1, not_ldap => 1 )->to( 'DeliveryServiceServer#remove_server_from_ds', namespace => $namespace );
	# -- DELIVERYSERVICES: SERVERS
	# Supports ?orderby=key
	$r->get("/api/$version/deliveryserviceserver")->over( authenticated => 1, not_ldap => 1 )->to( 'DeliveryServiceServer#index', namespace => $namespace );
	$r->post("/api/$version/deliveryserviceserver")->over( authenticated => 1, not_ldap => 1 )->to( 'DeliveryServiceServer#assign_servers_to_ds', namespace => $namespace );

		{1.2, http.MethodGet, `deliveryservices/{id}/servers$`, api.ReadHandler(dsserver.GetRefType(), d.DB),auth.PrivLevelReadOnly, Authenticated, nil},
		{1.2, http.MethodGet, `deliveryservices/{id}/unassigned_servers$`, api.ReadHandler(dsserver.GetRefType(), d.DB),auth.PrivLevelReadOnly, Authenticated, nil},
		{1.2, http.MethodGet, `deliveryservices/{id}/servers/eligible$`, api.ReadHandler(dsserver.GetRefType(), d.DB),auth.PrivLevelReadOnly, Authenticated, nil},

*/

func (dss TODeliveryServiceServer) GetKeyFieldsInfo() []api.KeyFieldInfo {
	return []api.KeyFieldInfo{{"deliveryservice", api.GetIntKey}, {"server", api.GetIntKey}}
}

//Implementation of the Identifier, Validator interface functions
func (dss TODeliveryServiceServer) GetKeys() (map[string]interface{}, bool) {
	if dss.DeliveryService == nil {
		return map[string]interface{}{"deliveryservice": 0}, false
	}
	if dss.Server == nil {
		return map[string]interface{}{"server": 0}, false
	}
	keys := make(map[string]interface{})
	ds_id := *dss.DeliveryService
	server_id := *dss.Server

	keys["deliveryservice"] = ds_id
	keys["server"] = server_id
	return keys, true
}

func (dss *TODeliveryServiceServer) GetAuditName() string {
	if dss.DeliveryService != nil {
		return strconv.Itoa(*dss.DeliveryService) + "-" + strconv.Itoa(*dss.Server)
	}
	return "unknown"
}

func (dss *TODeliveryServiceServer) GetType() string {
	return "deliveryserviceServers"
}

func (dss *TODeliveryServiceServer) SetKeys(keys map[string]interface{}) {
	ds_id, _ := keys["deliveryservice"].(int) //this utilizes the non panicking type assertion, if the thrown away ok variable is false i will be the zero of the type, 0 here.
	dss.DeliveryService = &ds_id

	server_id, _ := keys["server"].(int) //this utilizes the non panicking type assertion, if the thrown away ok variable is false i will be the zero of the type, 0 here.
	dss.Server = &server_id
}

// Validate fulfills the api.Validator interface
func (dss *TODeliveryServiceServer) Validate(db *sqlx.DB) []error {

	errs := validation.Errors{
		"deliveryservice": validation.Validate(dss.DeliveryService, validation.Required),
		"server":          validation.Validate(dss.Server, validation.Required),
	}

	return tovalidate.ToErrors(errs)
}

//The TODeliveryServiceServer implementation of the Creator interface
//all implementations of Creator should use transactions and return the proper errorType
//ParsePQUniqueConstraintError is used to determine if a profileparameter with conflicting values exists
//if so, it will return an errorType of DataConflict and the type should be appended to the
//generic error message returned
//The insert sql returns the profile and lastUpdated values of the newly inserted profileparameter and have
//to be added to the struct
func (dss *TODeliveryServiceServer) Create(db *sqlx.DB, user auth.CurrentUser) (error, tc.ApiErrorType) {
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
	resultRows, err := tx.NamedQuery(insertQuery(), dss)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			err, eType := dbhelpers.ParsePQUniqueConstraintError(pqErr)
			if eType == tc.DataConflictError {
				return errors.New("a parameter with " + err.Error()), eType
			}
			return err, eType
		}
		log.Errorf("received non pq error: %++v from create execution", err)
		return tc.DBError, tc.SystemError
	}
	defer resultRows.Close()

	var ds_id int
	var server_id int
	var lastUpdated tc.TimeNoMod
	rowsAffected := 0
	for resultRows.Next() {
		rowsAffected++
		if err := resultRows.Scan(&ds_id, &server_id, &lastUpdated); err != nil {
			log.Error.Printf("could not scan dss from insert: %s\n", err)
			return tc.DBError, tc.SystemError
		}
	}
	if rowsAffected == 0 {
		err = errors.New("no deliveryServiceServer was inserted, nothing to return")
		log.Errorln(err)
		return tc.DBError, tc.SystemError
	}
	if rowsAffected > 1 {
		err = errors.New("too many ids returned from parameter insert")
		log.Errorln(err)
		return tc.DBError, tc.SystemError
	}

	dss.SetKeys(map[string]interface{}{"deliveryservice": ds_id, "server": server_id})
	dss.LastUpdated = &lastUpdated
	err = tx.Commit()
	if err != nil {
		log.Errorln("Could not commit transaction: ", err)
		return tc.DBError, tc.SystemError
	}
	rollbackTransaction = false
	return nil, tc.NoError
}

func insertQuery() string {
	query := `INSERT INTO deliveryservice_server (
deliveryservice,
server) VALUES (
:ds_id,
:server_id) RETURNING deliveryservice, server, last_updated`
	return query
}

func (dss *TODeliveryServiceServer) Read(db *sqlx.DB, params map[string]string, user auth.CurrentUser) ([]interface{}, []error, tc.ApiErrorType) {
	idstr, ok := params["id"]

	if !ok {
		log.Errorf("Deliveryservice Server Id missing")
		return nil, []error{errors.New("Deliverservice id is required.")}, tc.DataMissingError
	}
	id, err := strconv.Atoi(idstr)

	if err != nil {
		log.Errorf("Deliveryservice Server Id is not an integer")
		return nil, []error{errors.New("Deliverservice id is not an integer.")}, tc.SystemError
	}

	query := selectQuery()
	log.Debugln("Query is ", query)

	rows, err := db.Queryx(query, id)
	if err != nil {
		log.Errorf("Error querying DeliveryserviceServers: %v", err)
		return nil, []error{tc.DBError}, tc.SystemError
	}
	defer rows.Close()

	servers := []interface{}{}
	for rows.Next() {
		var s tc.DssServer
		if err = rows.StructScan(&s); err != nil {
			log.Errorf("error parsing dss rows: %v", err)
			return nil, []error{tc.DBError}, tc.SystemError
		}
		hiddenField := ""
		if user.PrivLevel < auth.PrivLevelAdmin {
			s.ILOPassword = &hiddenField
		}
		servers = append(servers, s)
	}

	return servers, []error{}, tc.NoError

}

//The Parameter implementation of the Deleter interface
//all implementations of Deleter should use transactions and return the proper errorType
func (dss *TODeliveryServiceServer) Delete(db *sqlx.DB, user auth.CurrentUser) (error, tc.ApiErrorType) {
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
	log.Debugf("about to run exec query: %s with parameter: %++v", deleteQuery(), dss)
	result, err := tx.NamedExec(deleteQuery(), dss)
	if err != nil {
		log.Errorf("received error: %++v from delete execution", err)
		return tc.DBError, tc.SystemError
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return tc.DBError, tc.SystemError
	}
	if rowsAffected < 1 {
		return errors.New("no parameter with that id found"), tc.DataMissingError
	}
	if rowsAffected > 1 {
		return fmt.Errorf("this create affected too many rows: %d", rowsAffected), tc.SystemError
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
	s.router_host_name,
	s.router_port_name,
	st.name as status,
	s.status as status_id,
	s.tcp_port,
	t.name as server_type,
	s.type as server_type_id,
	s.upd_pending as upd_pending
	FROM server s
	JOIN cachegroup cg ON s.cachegroup = cg.id
	JOIN cdn cdn ON s.cdn_id = cdn.id
	JOIN phys_location pl ON s.phys_location = pl.id
	JOIN profile p ON s.profile = p.id
	JOIN status st ON s.status = st.id
	JOIN type t ON s.type = t.id
	WHERE s.id in (select server from deliveryservice_server where deliveryservice = $1)`

	return selectStmt
}

func updateQuery() string {
	query := `UPDATE
	profile_parameter SET
	profile=:profile_id,
	parameter=:parameter_id
	WHERE profile=:profile_id AND 
      parameter = :parameter_id 
      RETURNING last_updated`
	return query
}

func deleteQuery() string {
	query := `DELETE FROM profile_parameter
	WHERE profile=:profile_id and parameter=:parameter_id`
	return query
}
