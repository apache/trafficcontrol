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
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-tc/tovalidate"
	"github.com/apache/trafficcontrol/lib/go-util"
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
type TODeliveryServiceServer struct {
	api.APIInfoImpl `json:"-"`
	tc.DeliveryServiceServer
	TenantIDs          pq.Int64Array `json:"-" db:"accessibleTenants"`
	DeliveryServiceIDs pq.Int64Array `json:"-" db:"dsids"`
	ServerIDs          pq.Int64Array `json:"-" db:"serverids"`
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
func (dss *TODeliveryServiceServer) Validate(tx *sql.Tx) error {

	errs := validation.Errors{
		"deliveryservice": validation.Validate(dss.DeliveryService, validation.Required),
		"server":          validation.Validate(dss.Server, validation.Required),
	}

	return util.JoinErrs(tovalidate.ToErrors(errs))
}

// ReadDSSHandler list all of the Deliveryservice Servers in response to requests to api/1.1/deliveryserviceserver$
func ReadDSSHandler(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, []string{"limit", "page"})
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	dss := TODeliveryServiceServer{}
	dss.SetInfo(inf)
	results, err := dss.readDSS(inf.Tx, inf.User, inf.Params, inf.IntParams, nil, nil)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, err)
		return
	}
	api.WriteRespRaw(w, r, results)
}

// ReadDSSHandler list all of the Deliveryservice Servers in response to requests to api/1.1/deliveryserviceserver$
func ReadDSSHandlerV14(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, []string{"limit", "page"})
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	dsIDs := []int64{}
	dsIDStrs := strings.Split(inf.Params["deliveryserviceids"], ",")
	for _, dsIDStr := range dsIDStrs {
		dsIDStr = strings.TrimSpace(dsIDStr)
		if dsIDStr == "" {
			continue
		}
		dsID, err := strconv.Atoi(dsIDStr)
		if err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, 400, errors.New("deliveryserviceids query parameter must be a comma-delimited list of integers, got '"+inf.Params["deliveryserviceids"]+"'"), nil)
			return
		}
		dsIDs = append(dsIDs, int64(dsID))
	}

	serverIDs := []int64{}
	serverIDStrs := strings.Split(inf.Params["serverids"], ",")
	for _, serverIDStr := range serverIDStrs {
		serverIDStr = strings.TrimSpace(serverIDStr)
		if serverIDStr == "" {
			continue
		}
		serverID, err := strconv.Atoi(serverIDStr)
		if err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, 400, errors.New("serverids query parameter must be a comma-delimited list of integers, got '"+inf.Params["serverids"]+"'"), nil)
			return
		}
		serverIDs = append(serverIDs, int64(serverID))
	}

	dss := TODeliveryServiceServer{}
	dss.SetInfo(inf)
	results, err := dss.readDSS(inf.Tx, inf.User, inf.Params, inf.IntParams, dsIDs, serverIDs)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, err)
		return
	}
	api.WriteRespRaw(w, r, results)
}

func (dss *TODeliveryServiceServer) readDSS(tx *sqlx.Tx, user *auth.CurrentUser, params map[string]string, intParams map[string]int, dsIDs []int64, serverIDs []int64) (*tc.DeliveryServiceServerResponse, error) {
	orderby := params["orderby"]
	limit := 20
	offset := 0
	page := 0
	err := error(nil)
	if plimit, ok := intParams["limit"]; ok {
		limit = plimit
	}
	if ppage, ok := intParams["page"]; ok {
		page = ppage
		offset = page
		if offset > 0 {
			offset -= 1
		}
		offset *= limit
	}
	if orderby == "" {
		orderby = "deliveryService"
	}

	tenantIDs, err := tenant.GetUserTenantIDListTx(tx.Tx, user.TenantID)
	if err != nil {
		return nil, errors.New("getting user tenant ID list: " + err.Error())
	}
	for _, id := range tenantIDs {
		dss.TenantIDs = append(dss.TenantIDs, int64(id))
	}
	dss.ServerIDs = serverIDs
	dss.DeliveryServiceIDs = dsIDs

	query, err := selectQuery(orderby, strconv.Itoa(limit), strconv.Itoa(offset), dsIDs, serverIDs)
	if err != nil {
		return nil, errors.New("creating query for DeliveryserviceServers: " + err.Error())
	}
	log.Debugln("Query is ", query)

	rows, err := tx.NamedQuery(query, dss)
	if err != nil {
		return nil, errors.New("Error querying DeliveryserviceServers: " + err.Error())
	}
	defer rows.Close()
	servers := []tc.DeliveryServiceServer{}
	for rows.Next() {
		s := tc.DeliveryServiceServer{}
		if err = rows.StructScan(&s); err != nil {
			return nil, errors.New("error parsing dss rows: " + err.Error())
		}
		servers = append(servers, s)
	}
	return &tc.DeliveryServiceServerResponse{orderby, servers, page, limit}, nil
}

func selectQuery(orderBy string, limit string, offset string, dsIDs []int64, serverIDs []int64) (string, error) {
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

	// TODO refactor to use dbhelpers.AddTenancyCheck
	selectStmt += `
JOIN deliveryservice d on s.deliveryservice = d.id
WHERE d.tenant_id = ANY(CAST(:accessibleTenants AS bigint[]))
`
	if len(dsIDs) > 0 {
		selectStmt += `
AND s.deliveryservice = ANY(:dsids)
`
	}
	if len(serverIDs) > 0 {
		selectStmt += `
AND s.server = ANY(:serverids)
`
	}

	if orderBy != "" {
		selectStmt += ` ORDER BY ` + orderBy
	}

	selectStmt += ` LIMIT ` + limit + ` OFFSET ` + offset + ` ROWS`
	return selectStmt, nil
}

type DSServerIds struct {
	DsId    *int  `json:"dsId" db:"deliveryservice"`
	Servers []int `json:"servers"`
	Replace *bool `json:"replace"`
}

type TODSServerIds DSServerIds

func GetReplaceHandler(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, []string{"limit", "page"})
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	payload := DSServerIds{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("malformed JSON"), nil)
		return
	}

	servers := payload.Servers
	dsId := payload.DsId
	if servers == nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("servers must exist in post"), nil)
		return
	}
	if dsId == nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("dsid must exist in post"), nil)
		return
	}
	if payload.Replace == nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("replace must exist in post"), nil)
		return
	}

	ds, ok, err := GetDSInfo(inf.Tx.Tx, *dsId)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("deliveryserviceserver getting XMLID: "+err.Error()))
		return
	}
	if !ok {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("no delivery service with that ID exists"), nil)
		return
	}
	if userErr, sysErr, errCode := tenant.Check(inf.User, ds.Name, inf.Tx.Tx); userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	serverNamesCdnIdAndTypes, err := dbhelpers.GetServerHostNamesAndTypesFromIDs(inf.Tx.Tx, servers)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, err, nil)
		return
	}
	userErr = ValidateDSSAssignments(ds, serverNamesCdnIdAndTypes)
	if userErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, userErr, nil)
		return
	}

	usrErr, sysErr, status := ValidateServerCapabilities(ds.ID, serverNamesCdnIdAndTypes, inf.Tx.Tx)
	if usrErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, status, usrErr, sysErr)
		return
	}

	if *payload.Replace {
		// delete existing
		_, err := inf.Tx.Tx.Exec("DELETE FROM deliveryservice_server WHERE deliveryservice = $1", *dsId)
		if err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("unable to remove the existing servers assigned to the delivery service: "+err.Error()))
			return
		}
	}

	respServers := []int{}
	for _, server := range servers {
		dtos := map[string]interface{}{"id": dsId, "server": server}
		if _, err := inf.Tx.NamedExec(insertIdsQuery(), dtos); err != nil {
			usrErr, sysErr, code := api.ParseDBError(err)
			api.HandleErr(w, r, inf.Tx.Tx, code, usrErr, sysErr)
			return
		}
		respServers = append(respServers, server)
	}

	if err := deliveryservice.EnsureParams(inf.Tx.Tx, *dsId, ds.Name, ds.EdgeHeaderRewrite, ds.MidHeaderRewrite, ds.RegexRemap, ds.CacheURL, ds.SigningAlgorithm, ds.Type, ds.MaxOriginConnections); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("deliveryservice_server replace ensuring ds parameters: "+err.Error()))
		return
	}
	api.CreateChangeLogRawTx(api.ApiChange, "DS: "+ds.Name+", ID: "+strconv.Itoa(*dsId)+", ACTION: Replace existing servers assigned to delivery service", inf.User, inf.Tx.Tx)
	api.WriteRespAlertObj(w, r, tc.SuccessLevel, "server assignements complete", tc.DSSMapResponse{*dsId, *payload.Replace, respServers})
}

type TODeliveryServiceServers tc.DeliveryServiceServers

// GetCreateHandler assigns an existing Server to and existing Deliveryservice in response to api/1.1/deliveryservices/{xml_id}/servers
func GetCreateHandler(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"xml_id"}, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	dsName := inf.Params["xml_id"]

	if userErr, sysErr, errCode := tenant.Check(inf.User, dsName, inf.Tx.Tx); userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}

	ds, ok, err := GetDSInfoByName(inf.Tx.Tx, dsName)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("ds servers create scanning: "+err.Error()))
		return
	} else if !ok {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusNotFound, nil, errors.New("delivery service not found"))
		return
	}

	// get list of server Ids to insert
	payload := tc.DeliveryServiceServers{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("malformed JSON"), nil)
		return
	}
	payload.XmlId = dsName
	serverNames := payload.ServerNames

	serverNamesCdnIdAndTypes, err := dbhelpers.GetServerTypesCdnIdFromHostNames(inf.Tx.Tx, serverNames)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, err, nil)
		return
	}

	userErr = ValidateDSSAssignments(ds, serverNamesCdnIdAndTypes)
	if userErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, userErr, nil)
		return
	}

	usrErr, sysErr, status := ValidateServerCapabilities(ds.ID, serverNamesCdnIdAndTypes, inf.Tx.Tx)
	if usrErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, status, usrErr, sysErr)
		return
	}

	res, err := inf.Tx.Tx.Exec(`INSERT INTO deliveryservice_server (deliveryservice, server) SELECT $1, id FROM server WHERE host_name = ANY($2::text[])`, ds.ID, pq.Array(serverNames))
	if err != nil {

		usrErr, sysErr, code := api.ParseDBError(err)
		api.HandleErr(w, r, inf.Tx.Tx, code, usrErr, sysErr)
		return
	}

	if rowsAffected, err := res.RowsAffected(); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("ds servers inserting for create delivery service servers: getting rows affected: "+err.Error()))
		return
	} else if int(rowsAffected) != len(serverNames) {
		// this happens when the names they gave don't exist
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusNotFound, errors.New("servers not found"), nil)
		return
	}

	if err := deliveryservice.EnsureParams(inf.Tx.Tx, ds.ID, ds.Name, ds.EdgeHeaderRewrite, ds.MidHeaderRewrite, ds.RegexRemap, ds.CacheURL, ds.SigningAlgorithm, ds.Type, ds.MaxOriginConnections); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("deliveryservice_server replace ensuring ds parameters: "+err.Error()))
		return
	}
	api.CreateChangeLogRawTx(api.ApiChange, "DS: "+dsName+", ID: "+strconv.Itoa(ds.ID)+", ACTION: Assigned servers "+strings.Join(serverNames, ", ")+" to delivery service", inf.User, inf.Tx.Tx)
	api.WriteResp(w, r, tc.DeliveryServiceServers{payload.ServerNames, payload.XmlId})
}

// ValidateDSSAssignments returns an error if the given servers cannot be assigned to the given delivery service.
func ValidateDSSAssignments(ds DSInfo, servers []dbhelpers.ServerHostNameCDNIDAndType) error {
	if ds.Topology == nil {
		for _, s := range servers {
			if ds.CDNID != nil && s.CDNID != *ds.CDNID {
				return errors.New("server and delivery service CDNs do not match")
			}
		}
		return nil
	}
	for _, s := range servers {
		if s.Type != tc.OriginTypeName {
			return errors.New("only servers of type ORG may be assigned to topology-based delivery services")
		}
	}
	return nil
}

// ValidateServerCapabilities checks that the delivery service's requirements are met by each server to be assigned.
func ValidateServerCapabilities(dsID int, serverNamesAndTypes []dbhelpers.ServerHostNameCDNIDAndType, tx *sql.Tx) (error, error, int) {
	nonOriginServerNames := []string{}
	for _, s := range serverNamesAndTypes {
		if strings.HasPrefix(s.Type, tc.EdgeTypePrefix) {
			nonOriginServerNames = append(nonOriginServerNames, s.HostName)
		}
	}

	var sCaps []string
	dsCaps, err := dbhelpers.GetDSRequiredCapabilitiesFromID(dsID, tx)

	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	for _, name := range nonOriginServerNames {
		sCaps, err = dbhelpers.GetServerCapabilitiesFromName(name, tx)
		if err != nil {
			return nil, err, http.StatusInternalServerError
		}
		for _, dsc := range dsCaps {
			if !util.ContainsStr(sCaps, dsc) {
				return fmt.Errorf("Caching server cannot be assigned to this delivery service without having the required delivery service capabilities: [%v] for server %s", dsCaps, name), nil, http.StatusBadRequest
			}
		}
	}

	return nil, nil, 0
}

func insertIdsQuery() string {
	query := `INSERT INTO deliveryservice_server (deliveryservice, server)
VALUES (:id, :server )`
	return query
}

// GetReadAssigned retrieves lists of servers  based in the filter identified in the request: api/1.1/deliveryservices/{id}/servers|unassigned_servers|eligible
func GetReadAssigned(w http.ResponseWriter, r *http.Request) {
	getRead(w, r, false, tc.Alerts{})
}

// GetReadUnassigned retrieves lists of servers  based in the filter identified in the request: api/1.1/deliveryservices/{id}/servers|unassigned_servers|eligible
func GetReadUnassigned(w http.ResponseWriter, r *http.Request) {
	alerts := api.CreateDeprecationAlerts(nil)
	getRead(w, r, true, alerts)
}

func getRead(w http.ResponseWriter, r *http.Request, unassigned bool, alerts tc.Alerts) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"id"}, []string{"id"})
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	servers, err := read(inf.Tx, inf.IntParams["id"], inf.User, unassigned)
	if err != nil {
		alerts.AddNewAlert(tc.ErrorLevel, err.Error())
		api.WriteAlerts(w, r, http.StatusInternalServerError, alerts)
		return
	}

	if inf.Version.Major < 3 {
		v11ServerList := []tc.DSServerV11{}
		for _, srv := range servers {
			v11server := tc.DSServerV11{}
			v11server.DSServer = srv.DSServer

			interfaces := *srv.ServerInterfaces
			legacyInterface, err := tc.ConvertInterfaceInfotoV11(interfaces)
			if err != nil {
				api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("converting to server detail v11: "+err.Error()))
				return
			}
			v11server.LegacyInterfaceDetails = legacyInterface

			v11ServerList = append(v11ServerList, v11server)
		}
		api.WriteAlertsObj(w, r, http.StatusOK, alerts, v11ServerList)
		return
	}

	api.WriteAlertsObj(w, r, http.StatusOK, alerts, servers)
}

func read(tx *sqlx.Tx, dsID int, user *auth.CurrentUser, unassigned bool) ([]tc.DSServerV30, error) {
	where := `WHERE s.id in (select server from deliveryservice_server where deliveryservice = $1)`
	if unassigned {
		where = `WHERE s.id not in (select server from deliveryservice_server where deliveryservice = $1)`
	}
	query := dssSelectQuery() + where
	log.Debugln("Query is ", query)
	rows, err := tx.Queryx(query, dsID)
	if err != nil {
		return nil, errors.New("error querying dss rows: " + err.Error())
	}
	defer rows.Close()

	serverInterfaceInfo := []tc.ServerInterfaceInfo{}
	servers := []tc.DSServerV30{}
	for rows.Next() {
		s := tc.DSServerV30{}
		err := rows.Scan(
			&s.Cachegroup,
			&s.CachegroupID,
			&s.CDNID,
			&s.CDNName,
			&s.DomainName,
			&s.GUID,
			&s.HostName,
			&s.HTTPSPort,
			&s.ID,
			&s.ILOIPAddress,
			&s.ILOIPGateway,
			&s.ILOIPNetmask,
			&s.ILOPassword,
			&s.ILOUsername,
			pq.Array(&serverInterfaceInfo),
			&s.LastUpdated,
			&s.OfflineReason,
			&s.PhysLocation,
			&s.PhysLocationID,
			&s.Profile,
			&s.ProfileDesc,
			&s.ProfileID,
			&s.Rack,
			&s.RouterHostName,
			&s.RouterPortName,
			&s.Status,
			&s.StatusID,
			&s.TCPPort,
			&s.Type,
			&s.TypeID,
			&s.UpdPending,
		)
		if err != nil {
			return nil, errors.New("error scanning dss rows: " + err.Error())
		}
		s.ServerInterfaces = &serverInterfaceInfo

		if user.PrivLevel < auth.PrivLevelAdmin {
			s.ILOPassword = util.StrPtr("")
		}
		servers = append(servers, s)
	}
	return servers, nil
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
	ARRAY (
SELECT ( json_build_object (
'ipAddresses', ARRAY (
SELECT ( json_build_object (
'address', ip_address.address,
'gateway', ip_address.gateway,
'service_address', ip_address.service_address
))
FROM ip_address
WHERE ip_address.interface = interface.name
AND ip_address.server = s.id
),
'max_bandwidth', interface.max_bandwidth,
'monitor', interface.monitor,
'mtu', COALESCE (interface.mtu, 9000),
'name', interface.name
))
FROM interface
WHERE interface.server = s.id
) AS interfaces,
	s.last_updated,
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

type TODSSDeliveryService struct {
	api.APIInfoImpl `json:"-"`
	tc.DeliveryServiceNullable
}

// Read shows all of the delivery services associated with the specified server.
func (dss *TODSSDeliveryService) Read() ([]interface{}, error, error, int) {
	returnable := []interface{}{}
	params := dss.APIInfo().Params
	tx := dss.APIInfo().Tx.Tx
	user := dss.APIInfo().User

	if err := api.IsInt(params["id"]); err != nil {
		return nil, err, nil, http.StatusBadRequest
	}

	if _, ok := params["orderby"]; !ok {
		params["orderby"] = "xml_id"
	}

	// Query Parameters to Database Query column mappings
	// see the fields mapped in the SQL query
	queryParamsToSQLCols := map[string]dbhelpers.WhereColumnInfo{
		"xml_id": dbhelpers.WhereColumnInfo{"ds.xml_id", nil},
		"xmlId":  dbhelpers.WhereColumnInfo{"ds.xml_id", nil},
	}
	where, orderBy, pagination, queryValues, errs := dbhelpers.BuildWhereAndOrderByAndPagination(params, queryParamsToSQLCols)
	if len(errs) > 0 {
		return nil, nil, errors.New("reading server dses: " + util.JoinErrsStr(errs)), http.StatusInternalServerError
	}

	if where != "" {
		where = where + " AND "
	} else {
		where = "WHERE "
	}
	where += "ds.id in (SELECT deliveryService FROM deliveryservice_server where server = :server)"

	tenantIDs, err := tenant.GetUserTenantIDListTx(tx, user.TenantID)
	if err != nil {
		log.Errorln("received error querying for user's tenants: " + err.Error())
		return nil, nil, err, http.StatusInternalServerError
	}
	where, queryValues = dbhelpers.AddTenancyCheck(where, queryValues, "ds.tenant_id", tenantIDs)

	query := deliveryservice.GetDSSelectQuery() + where + orderBy + pagination
	queryValues["server"] = dss.APIInfo().Params["id"]
	log.Debugln("generated deliveryServices query: " + query)
	log.Debugf("executing with values: %++v\n", queryValues)

	dses, userErr, sysErr, _ := deliveryservice.GetDeliveryServices(query, queryValues, dss.APIInfo().Tx)
	if sysErr != nil {
		sysErr = fmt.Errorf("reading server dses: %v ", sysErr)
	}
	if userErr != nil || sysErr != nil {
		return nil, userErr, sysErr, http.StatusInternalServerError
	}

	for _, ds := range dses {
		returnable = append(returnable, ds)
	}
	return returnable, nil, nil, http.StatusOK
}

type DSInfo struct {
	ID                   int
	Name                 string
	Type                 tc.DSType
	EdgeHeaderRewrite    *string
	MidHeaderRewrite     *string
	RegexRemap           *string
	SigningAlgorithm     *string
	CacheURL             *string
	MaxOriginConnections *int
	Topology             *string
	CDNID                *int
}

// GetDSInfo loads the DeliveryService fields needed by Delivery Service Servers from the database, from the ID. Returns the data, whether the delivery service was found, and any error.
func GetDSInfo(tx *sql.Tx, id int) (DSInfo, bool, error) {
	qry := `
SELECT
  ds.xml_id,
  tp.name as type,
  ds.edge_header_rewrite,
  ds.mid_header_rewrite,
  ds.regex_remap,
  ds.signing_algorithm,
  ds.cacheurl,
  ds.max_origin_connections,
  ds.topology,
  ds.cdn_id
FROM
  deliveryservice ds
  JOIN type tp ON ds.type = tp.id
WHERE
  ds.id = $1
`
	di := DSInfo{ID: id}
	if err := tx.QueryRow(qry, id).Scan(&di.Name, &di.Type, &di.EdgeHeaderRewrite, &di.MidHeaderRewrite, &di.RegexRemap, &di.SigningAlgorithm, &di.CacheURL, &di.MaxOriginConnections, &di.Topology, &di.CDNID); err != nil {
		if err == sql.ErrNoRows {
			return DSInfo{}, false, nil
		}
		return DSInfo{}, false, fmt.Errorf("querying delivery service server ds info '%v': %v", id, err)
	}
	di.Type = tc.DSTypeFromString(string(di.Type))
	return di, true, nil
}

// GetDSInfoByName loads the DeliveryService fields needed by Delivery Service Servers from the database, from the ID. Returns the data, whether the delivery service was found, and any error.
func GetDSInfoByName(tx *sql.Tx, dsName string) (DSInfo, bool, error) {
	qry := `
SELECT
  ds.id,
  tp.name as type,
  ds.edge_header_rewrite,
  ds.mid_header_rewrite,
  ds.regex_remap,
  ds.signing_algorithm,
  ds.cacheurl,
  ds.max_origin_connections,
  ds.topology,
  ds.cdn_id
FROM
  deliveryservice ds
  JOIN type tp ON ds.type = tp.id
WHERE
  ds.xml_id = $1
`
	di := DSInfo{Name: dsName}
	if err := tx.QueryRow(qry, dsName).Scan(&di.ID, &di.Type, &di.EdgeHeaderRewrite, &di.MidHeaderRewrite, &di.RegexRemap, &di.SigningAlgorithm, &di.CacheURL, &di.MaxOriginConnections, &di.Topology, &di.CDNID); err != nil {
		if err == sql.ErrNoRows {
			return DSInfo{}, false, nil
		}
		return DSInfo{}, false, fmt.Errorf("querying delivery service server ds info by name '%v': %v", dsName, err)
	}
	di.Type = tc.DSTypeFromString(string(di.Type))
	return di, true, nil
}
