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
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc/tovalidate"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/deliveryservice"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/tenant"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/util/ims"

	validation "github.com/go-ozzo/ozzo-validation"
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
	CDN                string        `json:"-" db:"cdn"`
}

func (dss TODeliveryServiceServer) GetKeyFieldsInfo() []api.KeyFieldInfo {
	return []api.KeyFieldInfo{{Field: "deliveryservice", Func: api.GetIntKey}, {Field: "server", Func: api.GetIntKey}}
}

// Implementation of the Identifier, Validator interface functions
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
	if dss.DeliveryService != nil && dss.Server != nil {
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

// ReadDSSHandler is the handler for GET requests to /deliveryserviceserver.
func ReadDSSHandler(w http.ResponseWriter, r *http.Request) {
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
	cfg, e := api.GetConfig(r.Context())
	useIMS := false
	if e == nil && cfg != nil {
		useIMS = cfg.UseIMS
	} else {
		log.Warnf("Couldn't get config %v", e)
	}

	results, err, maxTime := dss.readDSS(r.Header, inf.Tx, inf.User, inf.Params, inf.IntParams, dsIDs, serverIDs, useIMS)
	if maxTime != nil && api.SetLastModifiedHeader(r, useIMS) {
		// RFC1123
		date := maxTime.Format("Mon, 02 Jan 2006 15:04:05 MST")
		w.Header().Add(rfc.LastModified, date)
	}
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, err)
		return
	}
	// statusnotmodified
	if err == nil && results == nil {
		w.WriteHeader(http.StatusNotModified)
	}
	if inf.Version.GreaterThanOrEqualTo(&api.Version{
		Major: 5,
		Minor: 0,
	}) {
		var resultsV5 tc.DeliveryServiceServerResponseV5
		resultsV5.Limit = results.Limit
		resultsV5.Orderby = results.Orderby
		resultsV5.Size = results.Size
		resultsV5.Alerts = results.Alerts
		resultsV5.Response = *upgrade(results.Response)
		api.WriteRespRaw(w, r, resultsV5)
	} else {
		api.WriteRespRaw(w, r, results)
	}
}

func upgrade(dsServers []tc.DeliveryServiceServer) *[]tc.DeliveryServiceServerV5 {
	var err error
	dsServersV5 := make([]tc.DeliveryServiceServerV5, len(dsServers))
	for i, s := range dsServers {
		dsServersV5[i].Server = s.Server
		dsServersV5[i].DeliveryService = s.DeliveryService
		if dsServersV5[i].LastUpdated, err = util.ConvertTimeFormat(s.LastUpdated.Time, time.RFC3339); err != nil {
			dsServersV5[i].LastUpdated = &s.LastUpdated.Time
		}
	}
	return &dsServersV5
}

func (dss *TODeliveryServiceServer) readDSS(h http.Header, tx *sqlx.Tx, user *auth.CurrentUser, params map[string]string, intParams map[string]int, dsIDs []int64, serverIDs []int64, useIMS bool) (*tc.DeliveryServiceServerResponse, error, *time.Time) {
	var maxTime time.Time
	var runSecond bool
	// NOTE: if the 'orderby' query param exists but has an empty value, that means no ordering should be done.
	// If the 'orderby' query param does not exist, order by "deliveryService" by default. This allows clients to
	// specifically skip sorting if it's unnecessary, reducing load on the DB.
	orderby, ok := params["orderby"]
	if !ok {
		orderby = "deliveryService"
	}
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

	tenantIDs, err := tenant.GetUserTenantIDListTx(tx.Tx, user.TenantID)
	if err != nil {
		return nil, errors.New("getting user tenant ID list: " + err.Error()), nil
	}
	for _, id := range tenantIDs {
		dss.TenantIDs = append(dss.TenantIDs, int64(id))
	}
	dss.ServerIDs = serverIDs
	dss.DeliveryServiceIDs = dsIDs

	queryValues := map[string]interface{}{}
	cdn := ""
	if cdnName, ok := params["cdn"]; ok {
		cdn = cdnName
		queryValues["cdn"] = cdnName
	}
	dss.CDN = cdn
	query1, err := selectQuery(orderby, strconv.Itoa(limit), strconv.Itoa(offset), dsIDs, serverIDs, true, cdn)
	if err != nil {
		log.Warnf("Error getting the max last updated query %v", err)
	}
	if useIMS {
		queryValues["accessibleTenants"] = pq.Array(tenantIDs)
		runSecond, maxTime = ims.TryIfModifiedSinceQuery(tx, h, queryValues, query1)
		if !runSecond {
			log.Debugln("IMS HIT")
			return nil, nil, &maxTime
		}
		log.Debugln("IMS MISS")
	} else {
		log.Debugln("Non IMS request")
	}
	query, err := selectQuery(orderby, strconv.Itoa(limit), strconv.Itoa(offset), dsIDs, serverIDs, false, cdn)
	if err != nil {
		return nil, errors.New("creating query for DeliveryserviceServers: " + err.Error()), nil
	}
	log.Debugln("Query is ", query)

	rows, err := tx.NamedQuery(query, dss)
	if err != nil {
		return nil, errors.New("Error querying DeliveryserviceServers: " + err.Error()), nil
	}
	defer rows.Close()
	servers := []tc.DeliveryServiceServer{}
	for rows.Next() {
		s := tc.DeliveryServiceServer{}
		if err = rows.StructScan(&s); err != nil {
			return nil, errors.New("error parsing dss rows: " + err.Error()), nil
		}
		servers = append(servers, s)
	}
	return &tc.DeliveryServiceServerResponse{Orderby: orderby, Response: servers, Size: page, Limit: limit}, nil, &maxTime
}

func selectQuery(orderBy string, limit string, offset string, dsIDs []int64, serverIDs []int64, getMaxQuery bool, cdn string) (string, error) {
	selectStmt := `SELECT
	s.deliveryService,
	s.server,
	s.last_updated
	FROM deliveryservice_server s`

	if getMaxQuery {
		selectStmt = `SELECT max(t) from ( (
		SELECT max(s.last_updated) as t FROM deliveryservice_server s`
	}
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
	if len(cdn) > 0 {
		selectStmt += `
AND d.cdn_id = (SELECT id FROM cdn WHERE name = :cdn)
`
	}

	if getMaxQuery {
		selectStmt += ` GROUP BY s.deliveryservice`
	}

	if orderBy != "" {
		selectStmt += ` ORDER BY ` + orderBy
	}

	selectStmt += ` LIMIT ` + limit + ` OFFSET ` + offset + ` ROWS `
	if getMaxQuery {
		return selectStmt + ` )
UNION ALL
select max(last_updated) as t from last_deleted l where l.table_name='deliveryservice_server') as res`, nil
	}
	return selectStmt, nil
}

type DSServerIds struct {
	DsId    *int  `json:"dsId" db:"deliveryservice"`
	Servers []int `json:"servers"`
	Replace *bool `json:"replace"`
}

type TODSServerIds DSServerIds

const checkPreExistingEdgeServersQuery = `
SELECT t.name AS name
FROM type t
JOIN server s ON t.id = s.type
JOIN status st ON s.status = st.id
JOIN deliveryservice_server dss ON dss.server = s.id
WHERE s.id = ANY(ARRAY(SELECT server FROM deliveryservice_server WHERE deliveryservice=$1))
AND (st.name = '` + string(tc.CacheStatusOnline) + `' OR st.name = '` + string(tc.CacheStatusReported) + `')
AND t.name like '` + string(tc.EdgeTypePrefix) + `%'
AND dss.deliveryservice=$1
LIMIT 1
`

func hasAvailableEdgesCurrentlyAssigned(tx *sql.Tx, dsID int) (bool, error) {
	t := ""
	if err := tx.QueryRow(checkPreExistingEdgeServersQuery, dsID).Scan(&t); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("couldn't query for pre existing edge servers for this DS: %s", err.Error())
	}
	return true, nil
}

// GetReplaceHandler is the handler for POST requests to the /deliveryserviceserver API endpoint.
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
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("deliveryserviceserver getting delivery service info for ID %d: %v", *dsId, err))
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
	if ds.CDNID != nil {
		cdn, ok, err := dbhelpers.GetCDNNameFromID(inf.Tx.Tx, int64(*ds.CDNID))
		if err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, err)
			return
		} else if !ok {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusNotFound, nil, nil)
			return
		}
		userErr, sysErr, statusCode := dbhelpers.CheckIfCurrentUserCanModifyCDN(inf.Tx.Tx, string(cdn), inf.User.UserName)
		if userErr != nil || sysErr != nil {
			api.HandleErr(w, r, inf.Tx.Tx, statusCode, userErr, sysErr)
			return
		}
	}
	serverInfos, err := dbhelpers.GetServerInfosFromIDs(inf.Tx.Tx, servers)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, err)
		return
	}

	userErr, sysErr, status := validateDSSAssignments(inf.Tx.Tx, ds, serverInfos, *payload.Replace)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, status, userErr, sysErr)
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

	if err := deliveryservice.EnsureParams(inf.Tx.Tx, *dsId, ds.Name, ds.EdgeHeaderRewrite, ds.MidHeaderRewrite, ds.RegexRemap, ds.SigningAlgorithm, ds.Type, ds.MaxOriginConnections); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("deliveryservice_server replace ensuring ds parameters: "+err.Error()))
		return
	}
	if err := deliveryservice.EnsureCacheURLParams(inf.Tx.Tx, ds.ID, ds.Name, ds.CacheURL); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("deliveryservice_server replace ensuring ds parameters: "+err.Error()))
		return
	}
	api.CreateChangeLogRawTx(api.ApiChange, "DS: "+ds.Name+", ID: "+strconv.Itoa(*dsId)+", ACTION: Replace existing servers assigned to delivery service", inf.User, inf.Tx.Tx)
	api.WriteRespAlertObj(w, r, tc.SuccessLevel, "server assignments complete", tc.DSSMapResponse{DsId: *dsId, Replace: *payload.Replace, Servers: respServers})
}

type TODeliveryServiceServers tc.DeliveryServiceServers

// GetCreateHandler is the handler for POST requests to /deliveryservices/{{XMLID}}/servers.
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
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("ds servers getting delivery service info for xmlID %s: %v", dsName, err))
		return
	} else if !ok {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusNotFound, nil, errors.New("delivery service not found"))
		return
	}

	if ds.CDNID != nil {
		cdn, ok, err := dbhelpers.GetCDNNameFromID(inf.Tx.Tx, int64(*ds.CDNID))
		if err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, err)
			return
		} else if !ok {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusNotFound, nil, nil)
			return
		}
		userErr, sysErr, statusCode := dbhelpers.CheckIfCurrentUserCanModifyCDN(inf.Tx.Tx, string(cdn), inf.User.UserName)
		if userErr != nil || sysErr != nil {
			api.HandleErr(w, r, inf.Tx.Tx, statusCode, userErr, sysErr)
			return
		}
	}

	// get list of server Ids to insert
	payload := tc.DeliveryServiceServers{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("malformed JSON"), nil)
		return
	}
	payload.XmlId = dsName
	serverNames := payload.ServerNames

	serverInfos, err := dbhelpers.GetServerInfosFromHostNames(inf.Tx.Tx, serverNames)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, err)
		return
	}

	userErr, sysErr, status := validateDSSAssignments(inf.Tx.Tx, ds, serverInfos, false)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, status, userErr, sysErr)
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

	if err := deliveryservice.EnsureParams(inf.Tx.Tx, ds.ID, ds.Name, ds.EdgeHeaderRewrite, ds.MidHeaderRewrite, ds.RegexRemap, ds.SigningAlgorithm, ds.Type, ds.MaxOriginConnections); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("deliveryservice_server replace ensuring ds parameters: "+err.Error()))
		return
	}
	if err := deliveryservice.EnsureCacheURLParams(inf.Tx.Tx, ds.ID, ds.Name, ds.CacheURL); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("deliveryservice_server replace ensuring ds parameters: "+err.Error()))
		return
	}
	api.CreateChangeLogRawTx(api.ApiChange, "DS: "+dsName+", ID: "+strconv.Itoa(ds.ID)+", ACTION: Assigned servers "+strings.Join(serverNames, ", ")+" to delivery service", inf.User, inf.Tx.Tx)
	api.WriteResp(w, r, tc.DeliveryServiceServers{ServerNames: payload.ServerNames, XmlId: payload.XmlId})
}

// validateDSSAssignments returns an error if the given servers cannot be assigned to the given delivery service.
func validateDSSAssignments(tx *sql.Tx, ds DSInfo, serverInfos []tc.ServerInfo, replace bool) (error, error, int) {
	anyAvailableServers := false
	userErr, sysErr, status := validateDSS(tx, ds, serverInfos)
	if userErr != nil || sysErr != nil {
		return userErr, sysErr, status
	}

	if ds.Active && replace {
		ids := make([]int, 0, len(serverInfos))
		newOrgCount := 0
		newAvailableEdgeCount := 0
		for _, inf := range serverInfos {
			ids = append(ids, inf.ID)
			// We dont check for the cache type to be = EDGE here because if this is a new DS, and we want to assign an online/ reported ORG to it,
			// we should be able to do that.
			if inf.Status == string(tc.CacheStatusOnline) || inf.Status == string(tc.CacheStatusReported) {
				anyAvailableServers = true
				if inf.Type == tc.OriginTypeName {
					newOrgCount++
				} else if strings.HasPrefix(inf.Type, tc.CacheTypeEdge.String()) {
					newAvailableEdgeCount++
				}
			}
		}
		// Prevent the user from deleting all the servers in an active, non-topology-based DS
		if len(ids) == 0 && ds.Topology == nil {
			return fmt.Errorf("this server assignment leaves Active Delivery Service #%d without any '%s' or '%s' servers", ds.ID, tc.CacheStatusOnline, tc.CacheStatusReported), nil, http.StatusConflict
		}
		// prevent the user from deleting all ORG servers from an active, MSO-enabled DS
		if ds.UseMultiSiteOrigin && newOrgCount < 1 {
			return fmt.Errorf("this server assignment leaves Active, MSO-enabled Delivery Service #%d without any '%s' or '%s' %s servers", ds.ID, tc.CacheStatusOnline, tc.CacheStatusReported, tc.OriginTypeName), nil, http.StatusConflict
		}
		if ds.Topology == nil {
			// The following check is necessary because of the following:
			// Consider a brand new active DS that has no server assignments.
			// Now, you wish to assign an online/ reported ORG server to it.
			// Since this is a new DS and it didnt have any "pre existing" online/ reported EDGEs, this should be possible.
			// However, if that DS had a couple of online/ reported EDGEs assigned to it, and now if you wanted to "replace"
			// that assignment with the new assignment of an online/ reported ORG, this should be prohibited by TO.
			currentlyHasAvailableEdgesAssigned, err := hasAvailableEdgesCurrentlyAssigned(tx, ds.ID)
			if err != nil {
				return nil, fmt.Errorf("checking for pre existing ONLINE/ REPORTED EDGES: %v", err), http.StatusInternalServerError
			}
			if (currentlyHasAvailableEdgesAssigned && newAvailableEdgeCount < 1) || !anyAvailableServers {
				return fmt.Errorf("this server assignment leaves Active Delivery Service #%d without any '%s' or '%s' servers", ds.ID, tc.CacheStatusOnline, tc.CacheStatusReported), nil, http.StatusConflict
			}
		}
	}

	userErr, sysErr, status = ValidateServerCapabilities(tx, ds.ID, serverInfos)
	if userErr != nil || sysErr != nil {
		return userErr, sysErr, status
	}
	return nil, nil, http.StatusOK
}

func validateDSS(tx *sql.Tx, ds DSInfo, servers []tc.ServerInfo) (error, error, int) {
	if ds.Topology == nil {
		for _, s := range servers {
			if ds.CDNID != nil && s.CDNID != *ds.CDNID {
				return errors.New("server and delivery service CDNs do not match"), nil, http.StatusBadRequest
			}
		}
		return nil, nil, http.StatusOK
	}
	for _, s := range servers {
		if s.Type != tc.OriginTypeName {
			return fmt.Errorf("only servers of type %s may be assigned to topology-based delivery services", tc.OriginTypeName), nil, http.StatusBadRequest
		}
	}

	_, cachegroups, sysErr := dbhelpers.GetTopologyCachegroups(tx, *ds.Topology)
	if sysErr != nil {
		return nil, fmt.Errorf("validating %s servers in topology %s: %v", tc.OriginTypeName, *ds.Topology, sysErr), http.StatusInternalServerError
	}
	userErr := CheckServersInCachegroups(servers, cachegroups)
	if userErr != nil {
		return fmt.Errorf("validating %s servers in topology %s: %v", tc.OriginTypeName, *ds.Topology, userErr), nil, http.StatusBadRequest
	}
	return nil, nil, http.StatusOK
}

// CheckServersInCachegroups checks whether or not all the given server cachegroups belong to the topology
// and returns a user error (if any).
func CheckServersInCachegroups(servers []tc.ServerInfo, cachegroups []string) error {
	cgSet := make(map[string]struct{}, len(cachegroups))
	for _, c := range cachegroups {
		cgSet[c] = struct{}{}
	}
	invalid := []string{}
	for _, s := range servers {
		if _, ok := cgSet[s.Cachegroup]; !ok {
			invalid = append(invalid, s.HostName+" ("+s.Cachegroup+")")
		}
	}
	if len(invalid) > 0 {
		return fmt.Errorf("the following servers are not in any of the given cachegroups (%s): %s", strings.Join(cachegroups, ", "), strings.Join(invalid, ", "))
	}
	return nil
}

// ValidateServerCapabilities checks that the delivery service's requirements are met by each server to be assigned.
func ValidateServerCapabilities(tx *sql.Tx, dsID int, serverNamesAndTypes []tc.ServerInfo) (error, error, int) {
	nonOriginServerNames := []string{}
	for _, s := range serverNamesAndTypes {
		if strings.HasPrefix(s.Type, tc.EdgeTypePrefix) {
			nonOriginServerNames = append(nonOriginServerNames, s.HostName)
		}
	}

	dsCaps, err := dbhelpers.GetDSRequiredCapabilitiesFromID(dsID, tx)

	if err != nil {
		return nil, fmt.Errorf("validating server capabilities: %v", err), http.StatusInternalServerError
	}

	serverCaps, err := dbhelpers.GetServerCapabilitiesOfServers(nonOriginServerNames, tx)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}
	for hostname, caps := range serverCaps {
		for _, dsc := range dsCaps {
			if !util.ContainsStr(caps, dsc) {
				return fmt.Errorf("the cache %s cannot be assigned to this delivery service without having the required delivery service capabilities: %v", hostname, dsCaps), nil, http.StatusBadRequest
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

// GetReadAssigned is the handler for GET requests to /deliveryservices/{id}/servers.
func GetReadAssigned(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"id"}, []string{"id"})
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	alerts := tc.Alerts{}
	servers, err := read(inf)
	if err != nil {
		alerts.AddNewAlert(tc.ErrorLevel, err.Error())
		api.WriteAlerts(w, r, http.StatusInternalServerError, alerts)
		return
	}

	if inf.Version.Major == 3 {
		v3ServerList := []tc.DSServer{}
		for _, srv := range servers {
			routerHostName := ""
			routerPort := ""
			interfaces := *srv.ServerInterfaces
			// All interfaces should have the same router name/port when they were upgraded from v1/2/3 to v4, so we can just choose any of them
			if len(interfaces) != 0 {
				routerHostName = interfaces[0].RouterHostName
				routerPort = interfaces[0].RouterPortName
			}
			v3Interfaces, err := tc.V4InterfaceInfoToV3Interfaces(interfaces)
			if err != nil {
				api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("converting to server detail v11: "+err.Error()))
				return
			}
			v3server := tc.DSServer{}
			pid, pdesc := dbhelpers.GetProfileIDDesc(inf.Tx.Tx, srv.ProfileNames[0])
			v3server.DSServerBase = srv.DSServerBaseV4.ToDSServerBase(&routerHostName, &routerPort, &pdesc, &pid)

			v3server.ServerInterfaces = &v3Interfaces

			v3ServerList = append(v3ServerList, v3server)
		}
		api.WriteAlertsObj(w, r, http.StatusOK, alerts, v3ServerList)
		return
	}

	// Based on version we load Delivery Service Server - for version 5 and above we use DSServerV5
	if inf.Version.GreaterThanOrEqualTo(&api.Version{Major: 5, Minor: 0}) {

		newServerList := make([]tc.DSServerV5, len(servers))

		for i, server := range servers {
			newServerList[i] = server.Upgrade()
		}

		api.WriteAlertsObj(w, r, http.StatusOK, alerts, newServerList)
		return
	}

	api.WriteAlertsObj(w, r, http.StatusOK, alerts, servers)
}

func read(inf *api.Info) ([]tc.DSServerV4, error) {
	queryDataString :=
		`,
cg.name as cachegroup,
s.cachegroup as cachegroup_id,
s.cdn_id,
cdn.name as cdn_name,
s.domain_name,
s.guid,
s.host_name,
s.https_port,
s.ilo_ip_address,
s.ilo_ip_gateway,
s.ilo_ip_netmask,
s.ilo_password,
s.ilo_username,
s.last_updated,
s.mgmt_ip_address,
s.mgmt_ip_gateway,
s.mgmt_ip_netmask,
s.offline_reason,
pl.name as phys_location,
s.phys_location as phys_location_id,
(SELECT ARRAY_AGG(profile_name) FROM server_profile WHERE server_profile.server=s.id) as profile_name,
s.rack,
st.name as status,
s.status as status_id,
s.tcp_port,
t.name as server_type,
s.type as server_type_id,
s.config_update_time > s.config_apply_time AS upd_pending
`

	queryFormatString := `
SELECT
	s.id
	%v
FROM server s
JOIN cachegroup cg ON s.cachegroup = cg.id
JOIN cdn cdn ON s.cdn_id = cdn.id
JOIN phys_location pl ON s.phys_location = pl.id
JOIN profile p ON s.profile = p.id
JOIN status st ON s.status = st.id
JOIN type t ON s.type = t.id
WHERE s.id in (select server from deliveryservice_server where deliveryservice = $1)`

	dsID := inf.IntParams["id"]
	idRows, err := inf.Tx.Queryx(fmt.Sprintf(queryFormatString, ""), dsID)
	if err != nil {
		return nil, errors.New("error querying dss ids: " + err.Error())
	}
	var serverIDs []int
	for idRows.Next() {
		var serverID int
		err := idRows.Scan(&serverID)
		if err != nil {
			return nil, errors.New("error scanning dss id rows: " + err.Error())
		}

		serverIDs = append(serverIDs, serverID)
	}

	serversMap, err := dbhelpers.GetServersInterfaces(serverIDs, inf.Tx.Tx)
	if err != nil {
		return nil, errors.New("unable to get server interfaces: " + err.Error())
	}

	rows, err := inf.Tx.Queryx(fmt.Sprintf(queryFormatString, queryDataString), dsID)
	if err != nil {
		return nil, errors.New("error querying dss rows: " + err.Error())
	}
	defer rows.Close()

	servers := []tc.DSServerV4{}
	for rows.Next() {
		s := tc.DSServerV4{}
		err := rows.Scan(
			&s.ID,
			&s.Cachegroup,
			&s.CachegroupID,
			&s.CDNID,
			&s.CDNName,
			&s.DomainName,
			&s.GUID,
			&s.HostName,
			&s.HTTPSPort,
			&s.ILOIPAddress,
			&s.ILOIPGateway,
			&s.ILOIPNetmask,
			&s.ILOPassword,
			&s.ILOUsername,
			&s.LastUpdated,
			&s.MgmtIPAddress,
			&s.MgmtIPGateway,
			&s.MgmtIPNetmask,
			&s.OfflineReason,
			&s.PhysLocation,
			&s.PhysLocationID,
			pq.Array(&s.ProfileNames),
			&s.Rack,
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
		s.ServerInterfaces = &[]tc.ServerInterfaceInfoV40{}
		if interfacesMap, ok := serversMap[*s.ID]; ok {
			for _, interfaceInfo := range interfacesMap {
				*s.ServerInterfaces = append(*s.ServerInterfaces, interfaceInfo)
			}
		}

		canViewILOPswd := false
		if (inf.Version.GreaterThanOrEqualTo(&api.Version{Major: 4}) && inf.Config.RoleBasedPermissions) || inf.Version.GreaterThanOrEqualTo(&api.Version{Major: 5}) {
			canViewILOPswd = inf.User.Can(tc.PermSecureServerRead)
		} else {
			canViewILOPswd = inf.User.PrivLevel == auth.PrivLevelAdmin
		}

		if !canViewILOPswd {
			s.ILOPassword = util.StrPtr("")
		}
		servers = append(servers, s)
	}
	return servers, nil
}

type TODSSDeliveryService struct {
	api.APIInfoImpl `json:"-"`
	tc.DeliveryServiceNullable
}

// Read shows all of the delivery services associated with the specified server.
func (dss *TODSSDeliveryService) Read(h http.Header, useIMS bool) ([]interface{}, error, error, int, *time.Time) {
	version := dss.APIInfo().Version
	if version == nil {
		return nil, nil, errors.New("TODSSDeliveryService.Read called with nil API version"), http.StatusInternalServerError, nil
	}
	if version.Major == 1 && version.Minor < 1 {
		return nil, nil, fmt.Errorf("TODSSDeliveryService.Read called with invalid API version: %d.%d", version.Major, version.Minor), http.StatusInternalServerError, nil
	}
	var maxTime time.Time
	var runSecond bool
	params := dss.APIInfo().Params
	tx := dss.APIInfo().Tx.Tx
	user := dss.APIInfo().User

	if err := api.IsInt(params["id"]); err != nil {
		return nil, err, nil, http.StatusBadRequest, nil
	}

	if _, ok := params["orderby"]; !ok {
		params["orderby"] = "xml_id"
	}

	// Query Parameters to Database Query column mappings
	// see the fields mapped in the SQL query
	queryParamsToSQLCols := map[string]dbhelpers.WhereColumnInfo{
		"xml_id": dbhelpers.WhereColumnInfo{Column: "ds.xml_id"},
		"xmlId":  dbhelpers.WhereColumnInfo{Column: "ds.xml_id"},
	}
	where, orderBy, pagination, queryValues, errs := dbhelpers.BuildWhereAndOrderByAndPagination(params, queryParamsToSQLCols)
	if len(errs) > 0 {
		return nil, nil, errors.New("reading server dses: " + util.JoinErrsStr(errs)), http.StatusInternalServerError, nil
	}

	if where != "" {
		where = where + " AND "
	} else {
		where = "WHERE "
	}

	serverID, _ := strconv.Atoi(params["id"])
	serverInfo, exists, err := dbhelpers.GetServerInfo(serverID, tx)
	if err != nil {
		return nil, nil, err, http.StatusInternalServerError, nil
	}
	if !exists {
		return nil, fmt.Errorf("server with ID %d doesn't exist", serverID), nil, http.StatusNotFound, nil
	}
	if serverInfo.Type == tc.OriginTypeName {
		where += `ds.id in (SELECT deliveryservice FROM deliveryservice_server WHERE server = :server)`
	} else {
		where += `
(ds.id in (
	SELECT deliveryservice FROM deliveryservice_server WHERE server = :server
) OR ds.id in (
	SELECT d.id FROM deliveryservice d
	JOIN cdn c ON d.cdn_id = c.id
	WHERE d.topology in (
		SELECT topology FROM topology_cachegroup
		WHERE cachegroup = (
			SELECT name FROM cachegroup
			WHERE id = (
				SELECT cachegroup FROM server WHERE id = :server
			)))
	AND d.cdn_id = (SELECT cdn_id FROM server WHERE id = :server)))
AND
((
(SELECT COALESCE(ARRAY_AGG(ssc.server_capability), '{}') 
FROM server_server_capability ssc 
WHERE ssc."server" = :server) 
@>
(
SELECT COALESCE(ds.required_capabilities, '{}')
)))
`
	}

	tenantIDs, err := tenant.GetUserTenantIDListTx(tx, user.TenantID)
	if err != nil {
		log.Errorln("received error querying for user's tenants: " + err.Error())
		return nil, nil, err, http.StatusInternalServerError, nil
	}
	where, queryValues = dbhelpers.AddTenancyCheck(where, queryValues, "ds.tenant_id", tenantIDs)
	query := deliveryservice.SelectDeliveryServicesQuery + where + orderBy + pagination
	queryValues["server"] = dss.APIInfo().Params["id"]

	if useIMS {
		runSecond, maxTime = ims.TryIfModifiedSinceQuery(dss.APIInfo().Tx, h, queryValues, selectMaxLastUpdatedQuery(where))
		if !runSecond {
			log.Debugln("IMS HIT")
			return nil, nil, nil, http.StatusNotModified, &maxTime
		}
		log.Debugln("IMS MISS")
	} else {
		log.Debugln("Non IMS request")
	}
	log.Debugln("generated deliveryServices query: " + query)
	log.Debugf("executing with values: %++v\n", queryValues)

	dses, userErr, sysErr, _ := deliveryservice.GetDeliveryServices(query, queryValues, dss.APIInfo().Tx)
	if sysErr != nil {
		sysErr = fmt.Errorf("reading server dses: %v ", sysErr)
	}
	if userErr != nil || sysErr != nil {
		return nil, userErr, sysErr, http.StatusInternalServerError, nil
	}

	returnable := make([]interface{}, 0, len(dses))
	for _, d := range dses {
		ds := d.DS
		if version.Major > 4 {
			returnable = append(returnable, ds)
		} else if version.Major > 3 && version.Minor >= 0 {
			returnable = append(returnable, ds.Downgrade())
		} else {
			legacyDS := ds.Downgrade()
			if version.Minor > 0 {
				dsV31 := legacyDS.DowngradeToV31()
				returnable = append(returnable, dsV31)
			} else {
				dsV30 := legacyDS.DowngradeToV31().DeliveryServiceV30
				returnable = append(returnable, dsV30)
			}
		}
	}
	return returnable, nil, nil, http.StatusOK, &maxTime
}

func selectMaxLastUpdatedQuery(where string) string {
	return `SELECT max(t) from (
		SELECT max(dss.last_updated) as t from deliveryservice_server dss JOIN deliveryservice ds ON ds.id = dss.deliveryservice ` + where +
		` UNION ALL
	select max(last_updated) as t from last_deleted l where l.table_name='deliveryservice_server') as res`
}

type DSInfo struct {
	Active               bool
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
	UseMultiSiteOrigin   bool
}

// language=sql
const getDSInfoBaseQuery = `
SELECT
  ds.active,
  ds.id,
  ds.xml_id,
  tp.name as type,
  ds.edge_header_rewrite,
  ds.mid_header_rewrite,
  ds.regex_remap,
  ds.signing_algorithm,
  ds.cacheurl,
  ds.max_origin_connections,
  ds.topology,
  ds.cdn_id,
  ds.multi_site_origin
FROM
  deliveryservice ds
  JOIN type tp ON ds.type = tp.id
`

func scanDSInfoRow(row *sql.Row) (DSInfo, bool, error) {
	di := DSInfo{}
	var useMSO *bool
	var active tc.DeliveryServiceActiveState
	if err := row.Scan(
		&active,
		&di.ID,
		&di.Name,
		&di.Type,
		&di.EdgeHeaderRewrite,
		&di.MidHeaderRewrite,
		&di.RegexRemap,
		&di.SigningAlgorithm,
		&di.CacheURL,
		&di.MaxOriginConnections,
		&di.Topology,
		&di.CDNID,
		&useMSO,
	); err != nil {
		if err == sql.ErrNoRows {
			return DSInfo{}, false, nil
		}
		return DSInfo{}, false, fmt.Errorf("querying delivery service server ds info: %v", err)
	}
	di.Active = active == tc.DSActiveStateActive
	di.Type = tc.DSTypeFromString(string(di.Type))
	if useMSO != nil {
		di.UseMultiSiteOrigin = *useMSO
	}
	return di, true, nil
}

// GetDSInfo loads the DeliveryService fields needed by Delivery Service Servers from the database, from the ID. Returns the data, whether the delivery service was found, and any error.
func GetDSInfo(tx *sql.Tx, id int) (DSInfo, bool, error) {
	qry := getDSInfoBaseQuery + `
WHERE ds.id = $1
`
	row := tx.QueryRow(qry, id)
	return scanDSInfoRow(row)
}

// GetDSInfoByName loads the DeliveryService fields needed by Delivery Service Servers from the database, from the name (xml_id). Returns the data, whether the delivery service was found, and any error.
func GetDSInfoByName(tx *sql.Tx, dsName string) (DSInfo, bool, error) {
	qry := getDSInfoBaseQuery + `
WHERE ds.xml_id = $1
`
	row := tx.QueryRow(qry, dsName)
	return scanDSInfoRow(row)
}
