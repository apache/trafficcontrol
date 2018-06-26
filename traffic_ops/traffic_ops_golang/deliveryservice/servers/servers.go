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
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-tc/tovalidate"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/deliveryservice"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/tenant"

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

// ReadDSSHandler list all of the Deliveryservice Servers in response to requests to api/1.1/deliveryserviceserver$
func ReadDSSHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//create error function with ResponseWriter and Request
		handleErrs := tc.GetHandleErrorsFunc(w, r)

		ctx := r.Context()

		// Load the PathParams into the query parameters for pass through
		params, err := api.GetCombinedParams(r)
		if err != nil {
			log.Errorf("unable to get parameters from request: %s", err)
			handleErrs(http.StatusInternalServerError, err)
		}

		user, err := auth.GetCurrentUser(ctx)
		if err != nil {
			log.Errorf("unable to retrieve current user from context: %s", err)
			handleErrs(http.StatusInternalServerError, err)
			return
		}

		results, errs, errType := GetRefType().readDSS(db, params, *user)
		if len(errs) > 0 {
			tc.HandleErrorsWithType(errs, errType, handleErrs)
			return
		}
		respBts, err := json.Marshal(results)
		if err != nil {
			handleErrs(http.StatusInternalServerError, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s", respBts)
	}
}
func (dss *TODeliveryServiceServer) readDSS(db *sqlx.DB, params map[string]string, user auth.CurrentUser) (*tc.DeliveryServiceServerResponse, []error, tc.ApiErrorType) {
	limitstr := params["limit"]
	pagestr := params["page"]
	orderby := params["orderby"]
	limit := 20
	offset := 1
	page := 1
	var err error = nil

	if limitstr != "" {
		limit, err = strconv.Atoi(limitstr)

		if err != nil {
			log.Errorf("limit parameter is not an integer")
			return nil, []error{errors.New("limit parameter must be an integer.")}, tc.SystemError
		}
	}

	if pagestr != "" {
		offset, err = strconv.Atoi(pagestr)
		page, err = strconv.Atoi(pagestr)

		if err != nil {
			log.Errorf("page parameter is not an integer")
			return nil, []error{errors.New("page parameter must be an integer.")}, tc.SystemError
		}

		if offset > 0 {
			offset -= 1
		}

		offset *= limit
	}

	if orderby == "" {
		orderby = "deliveryService"
	}

	query, err := selectQuery(orderby, strconv.Itoa(limit), strconv.Itoa(offset))
	log.Debugln("Query is ", query)

	rows, err := db.NamedQuery(query, dss)
	if err != nil {
		log.Errorf("Error querying DeliveryserviceServers: %v", err)
		return nil, []error{tc.DBError}, tc.SystemError
	}
	defer rows.Close()

	servers := []tc.DeliveryServiceServer{}
	for rows.Next() {
		var s tc.DeliveryServiceServer
		if err = rows.StructScan(&s); err != nil {
			log.Errorf("error parsing dss rows: %v", err)
			return nil, []error{tc.DBError}, tc.SystemError
		}
		servers = append(servers, s)
	}

	return &tc.DeliveryServiceServerResponse{orderby, servers, page, limit}, []error{}, tc.NoError
}

//all implementations of Deleter should use transactions and return the proper errorType

//The Parameter implementation of the Deleter interface
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
		log.Errorln("could not begin transaction: %v", err)
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

func selectQuery(orderBy string, limit string, offset string) (string, error) {

	selectStmt := `SELECT
	s.deliveryService,
	s.server,
	s.last_updated
	FROM deliveryservice_server s`

	allowedOrderByCols := map[string]string{
		"":                "",
		"deliveryservice": "s.deliveryService",
		"server":          "s.server",
		"lastupdated":     "s.last_updated",
		"deliveryService": "s.deliveryService",
		"lastUpdated":     "s.last_updated",
		"last_updated":    "s.last_updated",
	}
	orderBy, ok := allowedOrderByCols[orderBy]
	if !ok {
		return "", errors.New("orderBy '" + orderBy + "' not permitted")
	}

	if orderBy != "" {
		selectStmt += ` ORDER BY ` + orderBy
	}

	selectStmt += ` LIMIT ` + limit + ` OFFSET ` + offset + ` ROWS`
	return selectStmt, nil
}

func deleteQuery() string {
	query := `DELETE FROM deliveryservice_server
	WHERE deliveryservice=:deliveryservice and server=:server`
	return query
}

type DSServerIds struct {
	DsId    *int  `json:"dsId" db:"deliveryservice"`
	Servers []int `json:"servers"`
	Replace *bool `json:"replace"`
}

type TODSServerIds DSServerIds

func createServersForDsIdRef() *TODSServerIds {
	var dsserversRef = TODSServerIds(DSServerIds{})
	return &dsserversRef
}

func GetReplaceHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		handleErrs := tc.GetHandleErrorsFunc(w, r)
		ctx := r.Context()
		user, err := auth.GetCurrentUser(ctx)
		if err != nil {
			log.Errorf("unable to retrieve current user from context: %s", err)
			handleErrs(http.StatusInternalServerError, err)
			return
		}

		// get list of server Ids to insert
		payload := createServersForDsIdRef()

		if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
			log.Errorf("Error trying to decode the request body: %s", err)
			handleErrs(http.StatusInternalServerError, err)
			return
		}

		servers := payload.Servers
		dsId := payload.DsId

		if servers == nil {
			api.HandleErr(w, r, http.StatusBadRequest, errors.New("servers must exist in post"), nil)
			return
		}

		if dsId == nil {
			api.HandleErr(w, r, http.StatusBadRequest, errors.New("dsid must exist in post"), nil)
			return
		}

		if payload.Replace == nil {
			api.HandleErr(w, r, http.StatusBadRequest, errors.New("replace must exist in post"), nil)
			return
		}

		// perform the insert transaction
		rollbackTransaction := true
		tx, err := db.Beginx()
		if err != nil {
			log.Errorln("could not begin transaction: %v", err)
			handleErrs(http.StatusInternalServerError, err)
			return
		}
		defer func() {
			if tx == nil || !rollbackTransaction {
				return
			}
			err := tx.Rollback()
			if err != nil {
				log.Errorln(errors.New("rolling back transaction: " + err.Error()))
			}
		}()

		xmlID, ok, err := deliveryservice.GetXMLID(tx.Tx, *dsId)
		if err != nil {
			api.HandleErr(w, r, http.StatusInternalServerError, nil, errors.New("deliveryserviceserver getting XMLID: "+err.Error()))
			return
		}
		if !ok {
			api.HandleErr(w, r, http.StatusBadRequest, errors.New("no delivery service with that ID exists"), nil)
			return
		}
		if userErr, sysErr, errCode := tenant.Check(user, xmlID, tx.Tx); userErr != nil || sysErr != nil {
			api.HandleErr(w, r, errCode, userErr, sysErr)
			return
		}

		if *payload.Replace {
			// delete existing
			rows, err := db.Queryx("DELETE FROM deliveryservice_server WHERE deliveryservice = $1", *dsId)
			if err != nil {
				log.Errorf("unable to remove the existing servers assigned to the delivery service: %s", err)
				handleErrs(http.StatusInternalServerError, err)
				return
			}

			defer rows.Close()
		}

		i := 0
		respServers := []int{}

		for _, server := range servers {
			dtos := map[string]interface{}{"id": dsId, "server": server}
			resultRows, err := tx.NamedQuery(insertIdsQuery(), dtos)
			if err != nil {
				if pqErr, ok := err.(*pq.Error); ok {
					err, eType := dbhelpers.ParsePQUniqueConstraintError(pqErr)
					log.Errorln("could not begin transaction: %v", err)
					if eType == tc.DataConflictError {
						handleErrs(http.StatusInternalServerError, err)
						return
					}
					handleErrs(http.StatusInternalServerError, err)
					return
				}
				log.Errorf("received non pq error: %++v from create execution", err)
				return
			}
			respServers = append(respServers, server)
			resultRows.Next()
			i++
			defer resultRows.Close()
		}

		err = tx.Commit()
		if err != nil {
			log.Errorln("Could not commit transaction: ", err)
			return
		}
		resAlerts := []tc.Alert{tc.Alert{"server assignements complete", "success"}}
		repRes := tc.DSSReplaceResponse{resAlerts, tc.DSSMapResponse{*dsId, *payload.Replace, respServers}}

		// marshal the results to the response stream
		respBts, err := json.Marshal(repRes)
		if err != nil {
			log.Errorln("Could not marshal the response as expected: ", err)
			handleErrs(http.StatusInternalServerError, err)
			return
		}

		rollbackTransaction = false
		w.Header().Set(tc.ContentType, tc.ApplicationJson)
		fmt.Fprintf(w, "%s", respBts)
		return
	}
}

type TODeliveryServiceServers tc.DeliveryServiceServers

func createServersRef() *TODeliveryServiceServers {
	serversRef := TODeliveryServiceServers(tc.DeliveryServiceServers{})
	return &serversRef
}

// GetCreateHandler assigns an existing Server to and existing Deliveryservice in response to api/1.1/deliveryservices/{xml_id}/servers
func GetCreateHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleErrs := tc.GetHandleErrorsFunc(w, r)

		// find the delivery service Id dsId matching the xml_id
		params, err := api.GetCombinedParams(r)
		if err != nil {
			log.Errorf("unable to get parameters from request: %s", err)
			handleErrs(http.StatusInternalServerError, err)
			return
		}

		xmlId, ok := params["xml_id"]
		if !ok {
			log.Errorf("unable to get xml_id parameter from request: %s", err)
			handleErrs(http.StatusInternalServerError, err)
			return
		}

		ctx := r.Context()
		user, err := auth.GetCurrentUser(ctx)
		if err != nil {
			log.Errorf("unable to retrieve current user from context: %s", err)
			handleErrs(http.StatusInternalServerError, err)
			return
		}

		rollbackTransaction := true
		tx, err := db.Beginx()
		if err != nil {
			log.Errorln("could not begin transaction: %v", err)
			handleErrs(http.StatusInternalServerError, err)
			return
		}
		defer func() {
			if tx == nil || !rollbackTransaction {
				return
			}
			err := tx.Rollback()
			if err != nil {
				log.Errorln(errors.New("rolling back transaction: " + err.Error()))
			}
		}()

		if userErr, sysErr, errCode := tenant.Check(user, xmlId, tx.Tx); userErr != nil || sysErr != nil {
			api.HandleErr(w, r, errCode, userErr, sysErr)
			return
		}

		row := db.QueryRow(selectDeliveryService(), xmlId)
		var dsId int
		row.Scan(&dsId)

		// get list of server Ids to insert
		defer r.Body.Close()
		payload := createServersRef()

		if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
			log.Errorf("Error trying to decode the request body: %s", err)
			handleErrs(http.StatusInternalServerError, err)
			return
		}

		payload.XmlId = xmlId
		serverNames := payload.ServerNames
		q, arg, err := sqlx.In(selectServerIds(), serverNames)

		if err != nil {
			log.Errorln("Could not form IN query : %v", err)
			handleErrs(http.StatusInternalServerError, err)
			return
		}
		q = sqlx.Rebind(sqlx.DOLLAR, q)
		serverIds, err := db.Query(q, arg...)
		defer serverIds.Close()
		if err != nil {
			log.Errorln("Could not select the ServerIds: %v", err)
			handleErrs(http.StatusInternalServerError, err)
			return
		}

		// We have to get the server Ids and iterate through them because of a bug in the Go
		// transaction which returns an error if you perform a Select after an Insert in
		// the same transaction
		for serverIds.Next() {
			var serverId int
			err := serverIds.Scan(&serverId)
			dtos := map[string]interface{}{"id": dsId, "server": serverId}
			resultRows, err := tx.NamedQuery(insertIdsQuery(), dtos)
			if err != nil {
				if pqErr, ok := err.(*pq.Error); ok {
					err, eType := dbhelpers.ParsePQUniqueConstraintError(pqErr)
					log.Errorln("could not begin transaction: %v", err)
					if eType == tc.DataConflictError {
						handleErrs(http.StatusInternalServerError, err)
						return
					}
					handleErrs(http.StatusInternalServerError, err)
					return
				}
				log.Errorf("received non pq error: %++v from create execution", err)
				return
			}
			resultRows.Next()
		}

		err = tx.Commit()
		if err != nil {
			log.Errorln("Could not commit transaction: ", err)
			return
		}

		// marshal the results to the response stream
		tcPayload := tc.DeliveryServiceServers{payload.ServerNames, payload.XmlId}
		payloadResp := tc.DSServersResponse{tcPayload}
		respBts, err := json.Marshal(payloadResp)
		if err != nil {
			log.Errorln("Could not marshal the response as expected: ", err)
			handleErrs(http.StatusInternalServerError, err)
			return
		}

		rollbackTransaction = false
		w.Header().Set(tc.ContentType, tc.ApplicationJson)
		fmt.Fprintf(w, "%s", respBts)
		return
	}
}

func selectDeliveryService() string {
	query := `SELECT id FROM deliveryservice WHERE xml_id = $1`
	return query
}

func insertIdsQuery() string {
	query := `INSERT INTO deliveryservice_server (deliveryservice, server) 
VALUES (:id, :server )`
	return query
}

func selectServerIds() string {
	query := `SELECT id FROM server WHERE host_name in (?)`
	return query
}

// GetReadHandler retrieves lists of servers  based in the filter identified in the request: api/1.1/deliveryservices/{id}/servers|unassigned_servers|eligible
func GetReadHandler(db *sqlx.DB, filter tc.Filter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleErrs := tc.GetHandleErrorsFunc(w, r)
		params, err := api.GetCombinedParams(r)
		if err != nil {
			log.Errorf("unable to get parameters from request: %s", err)
			handleErrs(http.StatusInternalServerError, err)
		}

		where := `WHERE s.id in (select server from deliveryservice_server where deliveryservice = $1)`

		if filter == tc.Unassigned {
			where = `WHERE s.id not in (select server from deliveryservice_server where deliveryservice = $1)`
		}

		servers, errors, etype := read(db, params, auth.CurrentUser{}, where)

		if len(errors) > 0 {
			tc.HandleErrorsWithType(errors, etype, handleErrs)
			return
		}

		dssres := tc.DSServersAttrResponse{servers}
		respBts, err := json.Marshal(dssres)
		if err != nil {
			handleErrs(http.StatusInternalServerError, err)
			return
		}

		w.Header().Set(tc.ContentType, tc.ApplicationJson)
		fmt.Fprintf(w, "%s", respBts)
	}
}

func read(db *sqlx.DB, params map[string]string, user auth.CurrentUser, where string) ([]tc.DSServer, []error, tc.ApiErrorType) {
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

	query := dssSelectQuery() + where
	log.Debugln("Query is ", query)

	rows, err := db.Queryx(query, id)
	if err != nil {
		log.Errorf("Error querying DeliveryserviceServers: %v", err)
		return nil, []error{tc.DBError}, tc.SystemError
	}
	defer rows.Close()

	servers := []tc.DSServer{}
	for rows.Next() {
		var s tc.DSServer
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

func dssSelectQuery() string {

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
	JOIN type t ON s.type = t.id `

	return selectStmt
}

type TODSSDeliveryService tc.DSSDeliveryService

var dserviceRef = TODSSDeliveryService(tc.DSSDeliveryService{})

func GetDServiceRef() *TODSSDeliveryService {
	return &dserviceRef
}

// Read shows all of the delivery services associated with the specified server.
func (dss *TODSSDeliveryService) Read(db *sqlx.DB, params map[string]string, user auth.CurrentUser) ([]interface{}, []error, tc.ApiErrorType) {
	var err error = nil
	orderby := params["orderby"]
	serverId := params["id"]

	if orderby == "" {
		orderby = "deliveryService"
	}

	query := SDSSelectQuery()
	log.Debugln("Query is ", query)

	rows, err := db.Queryx(query, serverId)
	if err != nil {
		log.Errorf("Error querying DeliveryserviceServers: %v", err)
		return nil, []error{tc.DBError}, tc.SystemError
	}
	defer rows.Close()

	services := []interface{}{}
	for rows.Next() {
		var s tc.DSSDeliveryService
		if err = rows.StructScan(&s); err != nil {
			log.Errorf("error parsing dss rows: %v", err)
			return nil, []error{tc.DBError}, tc.SystemError
		}
		services = append(services, s)
	}

	return services, []error{}, tc.NoError
}

func SDSSelectQuery() string {

	selectStmt := `SELECT
 		active,
		ccr_dns_ttl,
		cdn_id,
		cacheurl,
		check_path,
		dns_bypass_cname,
		dns_bypass_ip,
		dns_bypass_ip6,
		dns_bypass_ttl,
		dscp,
		display_name,
		edge_header_rewrite,
		geo_limit,
		geo_limit_countries,
		geolimit_redirect_url,
		geo_provider,
		global_max_mbps,
		global_max_tps,
		http_bypass_fqdn,
		id,
		ipv6_routing_enabled,
		info_url,
		initial_dispersion,
		last_updated,
		logs_enabled,
		long_desc,
		long_desc_1,
		long_desc_2,
		max_dns_answers,
		mid_header_rewrite,
		miss_lat,
		miss_long,
		multi_site_origin,
		multi_site_origin_algorithm,
		(SELECT o.protocol::text || '://' || o.fqdn || rtrim(concat(':', o.port::text), ':')
		FROM origin o
		WHERE o.deliveryservice = d.id
		AND o.is_primary) as org_server_fqdn,
		origin_shield,
		profile,
		protocol,
		qstring_ignore,
		range_request_handling,
		regex_remap,
		regional_geo_blocking,
		remap_text,
		routing_name,
		ssl_key_version,
		signing_algorithm,
		tr_request_headers,
		tr_response_headers,
		tenant_id,
		type,
		xml_id
	FROM deliveryservice d
		WHERE id in (SELECT deliveryService FROM deliveryservice_server where server = $1)`
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
