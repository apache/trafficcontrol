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
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc/tovalidate"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/tenant"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

const (
	ServerCapabilityQueryParam = "serverCapability"
	ServerQueryParam           = "serverId"
	ServerHostNameQueryParam   = "serverHostName"
)

type TOServerServerCapabilityV5 struct {
	api.APIInfoImpl `json:"-"`
	tc.ServerServerCapabilityV5
}

func (ssc *TOServerServerCapabilityV5) SetLastUpdated(t tc.TimeNoMod) { ssc.LastUpdated = &t.Time }
func (ssc *TOServerServerCapabilityV5) NewReadObj() interface{} {
	return &tc.ServerServerCapabilityV5{}
}
func (ssc *TOServerServerCapabilityV5) SelectQuery() string { return scSelectQuery() }
func (ssc *TOServerServerCapabilityV5) ParamColumns() map[string]dbhelpers.WhereColumnInfo {
	return map[string]dbhelpers.WhereColumnInfo{
		ServerCapabilityQueryParam: dbhelpers.WhereColumnInfo{Column: "sc.server_capability"},
		ServerQueryParam:           dbhelpers.WhereColumnInfo{Column: "s.id", Checker: api.IsInt},
		ServerHostNameQueryParam:   dbhelpers.WhereColumnInfo{Column: "s.host_name"},
	}

}
func (ssc *TOServerServerCapabilityV5) DeleteQuery() string { return scDeleteQuery() }
func (ssc TOServerServerCapabilityV5) GetKeyFieldsInfo() []api.KeyFieldInfo {
	return []api.KeyFieldInfo{
		{Field: ServerQueryParam, Func: api.GetIntKey},
		{Field: ServerCapabilityQueryParam, Func: api.GetStringKey},
	}
}

// Need to satisfy Identifier interface but is a no-op as path does not have Update
func (ssc TOServerServerCapabilityV5) GetKeys() (map[string]interface{}, bool) {
	if ssc.ServerID == nil {
		return map[string]interface{}{ServerQueryParam: 0}, false
	}
	if ssc.ServerCapability == nil {
		return map[string]interface{}{ServerCapabilityQueryParam: 0}, false
	}
	return map[string]interface{}{
		ServerQueryParam:           *ssc.ServerID,
		ServerCapabilityQueryParam: *ssc.ServerCapability,
	}, true
}

func (ssc *TOServerServerCapabilityV5) SetKeys(keys map[string]interface{}) {
	sID, _ := keys[ServerQueryParam].(int)
	ssc.ServerID = &sID

	sc, _ := keys[ServerCapabilityQueryParam].(string)
	ssc.ServerCapability = &sc
}

func (ssc *TOServerServerCapabilityV5) GetAuditName() string {
	if ssc.ServerCapability != nil {
		return *ssc.ServerCapability
	}
	return "unknown"
}

func (ssc *TOServerServerCapabilityV5) GetType() string {
	return "server server_capability"
}

// Validate fulfills the api.Validator interface.
func (ssc TOServerServerCapabilityV5) Validate() (error, error) {
	errs := validation.Errors{
		ServerQueryParam:           validation.Validate(ssc.ServerID, validation.Required),
		ServerCapabilityQueryParam: validation.Validate(ssc.ServerCapability, validation.Required),
	}

	return util.JoinErrs(tovalidate.ToErrors(errs)), nil
}

func (ssc *TOServerServerCapabilityV5) Read(h http.Header, useIMS bool) ([]interface{}, error, error, int, *time.Time) {
	api.DefaultSort(ssc.APIInfo(), "serverHostName")
	return api.GenericRead(h, ssc, useIMS)
}
func (v *TOServerServerCapabilityV5) SelectMaxLastUpdatedQuery(where, orderBy, pagination, tableName string) string {
	return `SELECT max(t) from (
		SELECT max(sc.last_updated) as t from server_server_capability sc
JOIN server s ON sc.server = s.id ` + where + orderBy + pagination +
		` UNION ALL
	select max(last_updated) as t from last_deleted l where l.table_name='server_server_capability') as res`
}

func (ssc *TOServerServerCapabilityV5) Delete() (error, error, int) {
	tenantIDs, err := tenant.GetUserTenantIDListTx(ssc.APIInfo().Tx.Tx, ssc.APIInfo().User.TenantID)
	if err != nil {
		return nil, fmt.Errorf("deleting servers_server_capability: %v", err), http.StatusInternalServerError
	}
	accessibleTenants := make(map[int]struct{}, len(tenantIDs))
	for _, id := range tenantIDs {
		accessibleTenants[id] = struct{}{}
	}
	userErr, sysErr, status := checkTopologyBasedDSRequiredCapabilitiesV5(ssc, accessibleTenants)
	if userErr != nil || sysErr != nil {
		return userErr, sysErr, status
	}

	userErr, sysErr, status = checkDSRequiredCapabilitiesV5(ssc, accessibleTenants)
	if userErr != nil || sysErr != nil {
		return userErr, sysErr, status
	}

	if ssc.ServerID != nil {
		cdnName, err := dbhelpers.GetCDNNameFromServerID(ssc.APIInfo().Tx.Tx, int64(*ssc.ServerID))
		if err != nil {
			return nil, err, http.StatusInternalServerError
		}
		userErr, sysErr, errCode := dbhelpers.CheckIfCurrentUserCanModifyCDN(ssc.APIInfo().Tx.Tx, string(cdnName), ssc.APIInfo().User.UserName)
		if userErr != nil || sysErr != nil {
			return userErr, sysErr, errCode
		}
	}
	return api.GenericDelete(ssc)
}

func checkTopologyBasedDSRequiredCapabilitiesV5(ssc *TOServerServerCapabilityV5, accessibleTenants map[int]struct{}) (error, error, int) {
	dsRows, err := ssc.APIInfo().Tx.Tx.Query(getTopologyBasedDSesReqCapQuery(), ssc.ServerID, ssc.ServerCapability)
	if err != nil {
		return nil, fmt.Errorf("querying topology-based DSes with the required capability %s: %v", *ssc.ServerCapability, err), http.StatusInternalServerError
	}
	defer log.Close(dsRows, "closing dsRows in checkTopologyBasedDSRequiredCapabilitiesV5")

	xmlidToTopology := make(map[string]string)
	xmlidToTenantID := make(map[string]int)
	xmlidToReqCaps := make(map[string][]string)
	for dsRows.Next() {
		xmlID := ""
		topology := ""
		tenantID := 0
		reqCaps := []string{}
		if err := dsRows.Scan(&xmlID, &topology, &tenantID, pq.Array(&reqCaps)); err != nil {
			return nil, fmt.Errorf("scanning dsRows in checkTopologyBasedDSRequiredCapabilitiesV5: %v", err), http.StatusInternalServerError
		}
		xmlidToTenantID[xmlID] = tenantID
		xmlidToTopology[xmlID] = topology
		xmlidToReqCaps[xmlID] = reqCaps
	}
	if len(xmlidToTopology) == 0 {
		return nil, nil, http.StatusOK
	}

	serverRows, err := ssc.APIInfo().Tx.Tx.Query(getServerCapabilitiesOfCachegoupQuery(), ssc.ServerID, ssc.ServerCapability)
	if err != nil {
		return nil, fmt.Errorf("querying server capabilitites of server %d's cachegroup: %v", *ssc.ServerID, err), http.StatusInternalServerError
	}
	defer log.Close(serverRows, "closing serverRows in checkTopologyBasedDSRequiredCapabilitiesV5")

	serverIDToCapabilities := make(map[int]map[string]struct{})
	for serverRows.Next() {
		serverID := 0
		capabilities := []string{}
		if err := serverRows.Scan(&serverID, pq.Array(&capabilities)); err != nil {
			return nil, fmt.Errorf("scanning serverRows in checkTopologyBasedDSRequiredCapabilitiesV5: %v", err), http.StatusInternalServerError
		}
		serverIDToCapabilities[serverID] = make(map[string]struct{})
		for _, c := range capabilities {
			serverIDToCapabilities[serverID][c] = struct{}{}
		}
	}

	unsatisfiedDSes := []string{}
	for ds, dsReqCaps := range xmlidToReqCaps {
		dsIsSatisfied := false
		for _, serverCaps := range serverIDToCapabilities {
			serverHasCapabilities := true
			for _, dsReqCap := range dsReqCaps {
				if _, ok := serverCaps[dsReqCap]; !ok {
					serverHasCapabilities = false
					break
				}
			}
			if serverHasCapabilities {
				dsIsSatisfied = true
				break
			}
		}
		if !dsIsSatisfied {
			unsatisfiedDSes = append(unsatisfiedDSes, ds)
		}
	}
	if len(unsatisfiedDSes) == 0 {
		return nil, nil, http.StatusOK
	}

	dsStrings := make([]string, 0, len(unsatisfiedDSes))
	for _, ds := range unsatisfiedDSes {
		if _, ok := accessibleTenants[xmlidToTenantID[ds]]; ok {
			dsStrings = append(dsStrings, "(xml_id = "+ds+", topology = "+xmlidToTopology[ds]+")")
		}
	}
	return fmt.Errorf("this capability is required by delivery services, but there are no other servers in this server's cachegroup to satisfy them %s", strings.Join(dsStrings, ", ")), nil, http.StatusBadRequest
}

func checkDSRequiredCapabilitiesV5(ssc *TOServerServerCapabilityV5, accessibleTenants map[int]struct{}) (error, error, int) {
	// Ensure that the user is not removing a server capability from the server
	// that is required by the delivery services the server is assigned to (if applicable)
	dsIDs := []int64{}
	if err := ssc.APIInfo().Tx.Tx.QueryRow(checkDSReqCapQuery(), ssc.ServerID, ssc.ServerCapability).Scan(pq.Array(&dsIDs)); err != nil {
		return nil, fmt.Errorf("checking removing server server capability would still suffice delivery service requried capabilites: %v", err), http.StatusInternalServerError
	}

	if len(dsIDs) > 0 {
		return ssc.buildDSReqCapError(dsIDs, accessibleTenants)
	}
	return nil, nil, http.StatusOK
}

func (ssc *TOServerServerCapabilityV5) buildDSReqCapError(dsIDs []int64, accessibleTenants map[int]struct{}) (error, error, int) {

	dsTenantIDs, err := getDSTenantIDsByIDs(ssc.APIInfo().Tx, dsIDs)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	authDSIDs := []string{}

	for _, dsTenantID := range dsTenantIDs {
		if _, ok := accessibleTenants[dsTenantID.TenantID]; ok {
			if ok {
				authDSIDs = append(authDSIDs, strconv.Itoa(dsTenantID.ID))
			}
			continue
		}
	}

	dsStr := "delivery services"
	if len(authDSIDs) > 0 {
		dsStr = fmt.Sprintf("the delivery services %v", strings.Join(authDSIDs, ","))
	}
	return fmt.Errorf("cannot remove the capability %v from the server %v as the server is assigned to %v that require it", *ssc.ServerCapability, *ssc.ServerID, dsStr), nil, http.StatusBadRequest
}

func (ssc *TOServerServerCapabilityV5) Create() (error, error, int) {
	tx := ssc.APIInfo().Tx

	// Check existence prior to checking type
	_, exists, err := dbhelpers.GetServerNameFromID(tx.Tx, int64(*ssc.ServerID))
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}
	if !exists {
		return fmt.Errorf("server %v does not exist", *ssc.ServerID), nil, http.StatusNotFound
	}

	// Ensure type is correct
	var sidList []int64
	sidList = append(sidList, int64(*ssc.ServerID))
	errCode, userErr, sysErr := checkServerType(tx.Tx, sidList)
	if userErr != nil || sysErr != nil {
		return userErr, sysErr, errCode

	}

	cdnName, err := dbhelpers.GetCDNNameFromServerID(tx.Tx, int64(*ssc.ServerID))
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}
	userErr, sysErr, errCode = dbhelpers.CheckIfCurrentUserCanModifyCDN(tx.Tx, string(cdnName), ssc.APIInfo().User.UserName)
	if userErr != nil || sysErr != nil {
		return userErr, sysErr, errCode
	}

	resultRows, err := tx.NamedQuery(scInsertQuery(), ssc)
	if err != nil {
		return api.ParseDBError(err)
	}
	defer resultRows.Close()

	rowsAffected := 0
	for resultRows.Next() {
		rowsAffected++
		if err := resultRows.StructScan(&ssc); err != nil {
			return nil, errors.New(ssc.GetType() + " create scanning: " + err.Error()), http.StatusInternalServerError
		}
	}
	if rowsAffected == 0 {
		return nil, errors.New(ssc.GetType() + " create: no " + ssc.GetType() + " was inserted, no rows was returned"), http.StatusInternalServerError
	} else if rowsAffected > 1 {
		return nil, errors.New("too many rows returned from " + ssc.GetType() + " insert"), http.StatusInternalServerError
	}

	return nil, nil, http.StatusOK
}

type (
	TOServerServerCapability struct {
		api.APIInfoImpl `json:"-"`
		tc.ServerServerCapability
	}

	DSTenant struct {
		TenantID int `db:"tenant_id"`
		ID       int `db:"id"`
	}
)

func (ssc *TOServerServerCapability) SetLastUpdated(t tc.TimeNoMod) { ssc.LastUpdated = &t }
func (ssc *TOServerServerCapability) NewReadObj() interface{} {
	return &tc.ServerServerCapability{}
}
func (ssc *TOServerServerCapability) SelectQuery() string { return scSelectQuery() }
func (ssc *TOServerServerCapability) ParamColumns() map[string]dbhelpers.WhereColumnInfo {
	return map[string]dbhelpers.WhereColumnInfo{
		ServerCapabilityQueryParam: dbhelpers.WhereColumnInfo{Column: "sc.server_capability"},
		ServerQueryParam:           dbhelpers.WhereColumnInfo{Column: "s.id", Checker: api.IsInt},
		ServerHostNameQueryParam:   dbhelpers.WhereColumnInfo{Column: "s.host_name"},
	}

}
func (ssc *TOServerServerCapability) DeleteQuery() string { return scDeleteQuery() }
func (ssc TOServerServerCapability) GetKeyFieldsInfo() []api.KeyFieldInfo {
	return []api.KeyFieldInfo{
		{Field: ServerQueryParam, Func: api.GetIntKey},
		{Field: ServerCapabilityQueryParam, Func: api.GetStringKey},
	}
}

// Need to satisfy Identifier interface but is a no-op as path does not have Update
func (ssc TOServerServerCapability) GetKeys() (map[string]interface{}, bool) {
	if ssc.ServerID == nil {
		return map[string]interface{}{ServerQueryParam: 0}, false
	}
	if ssc.ServerCapability == nil {
		return map[string]interface{}{ServerCapabilityQueryParam: 0}, false
	}
	return map[string]interface{}{
		ServerQueryParam:           *ssc.ServerID,
		ServerCapabilityQueryParam: *ssc.ServerCapability,
	}, true
}

func (ssc *TOServerServerCapability) SetKeys(keys map[string]interface{}) {
	sID, _ := keys[ServerQueryParam].(int)
	ssc.ServerID = &sID

	sc, _ := keys[ServerCapabilityQueryParam].(string)
	ssc.ServerCapability = &sc
}

func (ssc *TOServerServerCapability) GetAuditName() string {
	if ssc.ServerCapability != nil {
		return *ssc.ServerCapability
	}
	return "unknown"
}

func (ssc *TOServerServerCapability) GetType() string {
	return "server server_capability"
}

// Validate fulfills the api.Validator interface.
func (ssc TOServerServerCapability) Validate() (error, error) {
	errs := validation.Errors{
		ServerQueryParam:           validation.Validate(ssc.ServerID, validation.Required),
		ServerCapabilityQueryParam: validation.Validate(ssc.ServerCapability, validation.Required),
	}

	return util.JoinErrs(tovalidate.ToErrors(errs)), nil
}

func (ssc *TOServerServerCapability) Read(h http.Header, useIMS bool) ([]interface{}, error, error, int, *time.Time) {
	api.DefaultSort(ssc.APIInfo(), "serverHostName")
	return api.GenericRead(h, ssc, useIMS)
}
func (v *TOServerServerCapability) SelectMaxLastUpdatedQuery(where, orderBy, pagination, tableName string) string {
	return `SELECT max(t) from (
		SELECT max(sc.last_updated) as t from server_server_capability sc
JOIN server s ON sc.server = s.id ` + where + orderBy + pagination +
		` UNION ALL
	select max(last_updated) as t from last_deleted l where l.table_name='server_server_capability') as res`
}

func (ssc *TOServerServerCapability) Delete() (error, error, int) {
	tenantIDs, err := tenant.GetUserTenantIDListTx(ssc.APIInfo().Tx.Tx, ssc.APIInfo().User.TenantID)
	if err != nil {
		return nil, fmt.Errorf("deleting servers_server_capability: %v", err), http.StatusInternalServerError
	}
	accessibleTenants := make(map[int]struct{}, len(tenantIDs))
	for _, id := range tenantIDs {
		accessibleTenants[id] = struct{}{}
	}
	userErr, sysErr, status := checkTopologyBasedDSRequiredCapabilities(ssc, accessibleTenants)
	if userErr != nil || sysErr != nil {
		return userErr, sysErr, status
	}

	userErr, sysErr, status = checkDSRequiredCapabilities(ssc, accessibleTenants)
	if userErr != nil || sysErr != nil {
		return userErr, sysErr, status
	}

	if ssc.ServerID != nil {
		cdnName, err := dbhelpers.GetCDNNameFromServerID(ssc.APIInfo().Tx.Tx, int64(*ssc.ServerID))
		if err != nil {
			return nil, err, http.StatusInternalServerError
		}
		userErr, sysErr, errCode := dbhelpers.CheckIfCurrentUserCanModifyCDN(ssc.APIInfo().Tx.Tx, string(cdnName), ssc.APIInfo().User.UserName)
		if userErr != nil || sysErr != nil {
			return userErr, sysErr, errCode
		}
	}
	return api.GenericDelete(ssc)
}

func checkTopologyBasedDSRequiredCapabilities(ssc *TOServerServerCapability, accessibleTenants map[int]struct{}) (error, error, int) {
	dsRows, err := ssc.APIInfo().Tx.Tx.Query(getTopologyBasedDSesReqCapQuery(), ssc.ServerID, ssc.ServerCapability)
	if err != nil {
		return nil, fmt.Errorf("querying topology-based DSes with the required capability %s: %v", *ssc.ServerCapability, err), http.StatusInternalServerError
	}
	defer log.Close(dsRows, "closing dsRows in checkTopologyBasedDSRequiredCapabilities")

	xmlidToTopology := make(map[string]string)
	xmlidToTenantID := make(map[string]int)
	xmlidToReqCaps := make(map[string][]string)
	for dsRows.Next() {
		xmlID := ""
		topology := ""
		tenantID := 0
		reqCaps := []string{}
		if err := dsRows.Scan(&xmlID, &topology, &tenantID, pq.Array(&reqCaps)); err != nil {
			return nil, fmt.Errorf("scanning dsRows in checkTopologyBasedDSRequiredCapabilities: %v", err), http.StatusInternalServerError
		}
		xmlidToTenantID[xmlID] = tenantID
		xmlidToTopology[xmlID] = topology
		xmlidToReqCaps[xmlID] = reqCaps
	}
	if len(xmlidToTopology) == 0 {
		return nil, nil, http.StatusOK
	}

	serverRows, err := ssc.APIInfo().Tx.Tx.Query(getServerCapabilitiesOfCachegoupQuery(), ssc.ServerID, ssc.ServerCapability)
	if err != nil {
		return nil, fmt.Errorf("querying server capabilitites of server %d's cachegroup: %v", *ssc.ServerID, err), http.StatusInternalServerError
	}
	defer log.Close(serverRows, "closing serverRows in checkTopologyBasedDSRequiredCapabilities")

	serverIDToCapabilities := make(map[int]map[string]struct{})
	for serverRows.Next() {
		serverID := 0
		capabilities := []string{}
		if err := serverRows.Scan(&serverID, pq.Array(&capabilities)); err != nil {
			return nil, fmt.Errorf("scanning serverRows in checkTopologyBasedDSRequiredCapabilities: %v", err), http.StatusInternalServerError
		}
		serverIDToCapabilities[serverID] = make(map[string]struct{})
		for _, c := range capabilities {
			serverIDToCapabilities[serverID][c] = struct{}{}
		}
	}

	unsatisfiedDSes := []string{}
	for ds, dsReqCaps := range xmlidToReqCaps {
		dsIsSatisfied := false
		for _, serverCaps := range serverIDToCapabilities {
			serverHasCapabilities := true
			for _, dsReqCap := range dsReqCaps {
				if _, ok := serverCaps[dsReqCap]; !ok {
					serverHasCapabilities = false
					break
				}
			}
			if serverHasCapabilities {
				dsIsSatisfied = true
				break
			}
		}
		if !dsIsSatisfied {
			unsatisfiedDSes = append(unsatisfiedDSes, ds)
		}
	}
	if len(unsatisfiedDSes) == 0 {
		return nil, nil, http.StatusOK
	}

	dsStrings := make([]string, 0, len(unsatisfiedDSes))
	for _, ds := range unsatisfiedDSes {
		if _, ok := accessibleTenants[xmlidToTenantID[ds]]; ok {
			dsStrings = append(dsStrings, "(xml_id = "+ds+", topology = "+xmlidToTopology[ds]+")")
		}
	}
	return fmt.Errorf("this capability is required by delivery services, but there are no other servers in this server's cachegroup to satisfy them %s", strings.Join(dsStrings, ", ")), nil, http.StatusBadRequest
}

func checkDSRequiredCapabilities(ssc *TOServerServerCapability, accessibleTenants map[int]struct{}) (error, error, int) {
	// Ensure that the user is not removing a server capability from the server
	// that is required by the delivery services the server is assigned to (if applicable)
	dsIDs := []int64{}
	if err := ssc.APIInfo().Tx.Tx.QueryRow(checkDSReqCapQuery(), ssc.ServerID, ssc.ServerCapability).Scan(pq.Array(&dsIDs)); err != nil {
		return nil, fmt.Errorf("checking removing server server capability would still suffice delivery service requried capabilites: %v", err), http.StatusInternalServerError
	}

	if len(dsIDs) > 0 {
		return ssc.buildDSReqCapError(dsIDs, accessibleTenants)
	}
	return nil, nil, http.StatusOK
}

func (ssc *TOServerServerCapability) buildDSReqCapError(dsIDs []int64, accessibleTenants map[int]struct{}) (error, error, int) {

	dsTenantIDs, err := getDSTenantIDsByIDs(ssc.APIInfo().Tx, dsIDs)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	authDSIDs := []string{}

	for _, dsTenantID := range dsTenantIDs {
		if _, ok := accessibleTenants[dsTenantID.TenantID]; ok {
			if ok {
				authDSIDs = append(authDSIDs, strconv.Itoa(dsTenantID.ID))
			}
			continue
		}
	}

	dsStr := "delivery services"
	if len(authDSIDs) > 0 {
		dsStr = fmt.Sprintf("the delivery services %v", strings.Join(authDSIDs, ","))
	}
	return fmt.Errorf("cannot remove the capability %v from the server %v as the server is assigned to %v that require it", *ssc.ServerCapability, *ssc.ServerID, dsStr), nil, http.StatusBadRequest
}

func (ssc *TOServerServerCapability) Create() (error, error, int) {
	tx := ssc.APIInfo().Tx

	// Check existence prior to checking type
	_, exists, err := dbhelpers.GetServerNameFromID(tx.Tx, int64(*ssc.ServerID))
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}
	if !exists {
		return fmt.Errorf("server %v does not exist", *ssc.ServerID), nil, http.StatusNotFound
	}

	// Ensure type is correct
	var sidList []int64
	sidList = append(sidList, int64(*ssc.ServerID))
	errCode, userErr, sysErr := checkServerType(tx.Tx, sidList)
	if userErr != nil || sysErr != nil {
		return userErr, sysErr, errCode

	}

	cdnName, err := dbhelpers.GetCDNNameFromServerID(tx.Tx, int64(*ssc.ServerID))
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}
	userErr, sysErr, errCode = dbhelpers.CheckIfCurrentUserCanModifyCDN(tx.Tx, string(cdnName), ssc.APIInfo().User.UserName)
	if userErr != nil || sysErr != nil {
		return userErr, sysErr, errCode
	}

	resultRows, err := tx.NamedQuery(scInsertQuery(), ssc)
	if err != nil {
		return api.ParseDBError(err)
	}
	defer resultRows.Close()

	rowsAffected := 0
	for resultRows.Next() {
		rowsAffected++
		if err := resultRows.StructScan(&ssc); err != nil {
			return nil, errors.New(ssc.GetType() + " create scanning: " + err.Error()), http.StatusInternalServerError
		}
	}
	if rowsAffected == 0 {
		return nil, errors.New(ssc.GetType() + " create: no " + ssc.GetType() + " was inserted, no rows was returned"), http.StatusInternalServerError
	} else if rowsAffected > 1 {
		return nil, errors.New("too many rows returned from " + ssc.GetType() + " insert"), http.StatusInternalServerError
	}

	return nil, nil, http.StatusOK
}

func scSelectQuery() string {
	return `SELECT
sc.server_capability,
sc.server,
sc.last_updated,
s.host_name as host_name
FROM server_server_capability sc
JOIN server s ON sc.server = s.id`
}

func scDeleteQuery() string {
	return `DELETE FROM server_server_capability
WHERE server = :server AND server_capability = :server_capability`
}

func scInsertQuery() string {
	return `INSERT INTO server_server_capability (
server_capability,
server) VALUES (
:server_capability,
:server) RETURNING server, server_capability, last_updated`
}

func checkDSReqCapQuery() string {
	return `
SELECT ARRAY(
	SELECT ds.id
	FROM deliveryservice as ds
	WHERE id IN (
		SELECT deliveryservice
		FROM deliveryservice_server
		WHERE server = $1)
	AND $2 = ANY(ds.required_capabilities))`
}

// get the topology-based DSes (with all their required capabilities) that a given
// server is assigned to, filtered by the given capability
func getTopologyBasedDSesReqCapQuery() string {
	return `
SELECT
  ds.xml_id,
  ds.topology,
  ds.tenant_id,
  ds.required_capabilities AS req_caps
FROM server s
JOIN cachegroup c ON s.cachegroup = c.id
JOIN topology_cachegroup tc ON c.name = tc.cachegroup
JOIN deliveryservice ds ON ds.topology = tc.topology
WHERE s.id = $1
GROUP BY ds.xml_id, ds.tenant_id, ds.topology, ds.required_capabilities
HAVING $2 = ANY(ds.required_capabilities)
`
}

// get all the capabilities of the servers in a given server's cachegroup
// that have a given capability
func getServerCapabilitiesOfCachegoupQuery() string {
	return `
SELECT s.id, ARRAY_AGG(ssc.server_capability) AS capabilities
FROM server s
JOIN cachegroup c ON c.id = s.cachegroup AND c.id = (SELECT cachegroup FROM server WHERE server.id = $1)
JOIN server_server_capability ssc ON ssc.server = s.id
WHERE
  s.cdn_id = (SELECT cdn_id FROM server WHERE server.id = $1)
  AND s.id != $1
GROUP BY s.id
HAVING $2 = ANY(ARRAY_AGG(ssc.server_capability));
`
}

func getDSTenantIDsByIDs(tx *sqlx.Tx, dsIDs []int64) ([]DSTenant, error) {
	dsTenantIDs := []DSTenant{}

	query, args, err := sqlx.In("SELECT id, tenant_id FROM deliveryservice where id IN (?);", dsIDs)
	if err != nil {
		return nil, fmt.Errorf("building query for getting delivery services' tenants: %v", err)
	}
	query = tx.Rebind(query)
	resultRows, err := tx.Queryx(query, args...)
	if err != nil {
		return nil, fmt.Errorf("querying tenant IDs for delivery service IDs: %v", err)
	}
	defer log.Close(resultRows, "closing resultRows in getDSTenantIDsByIDs")

	for resultRows.Next() {
		dsTenantID := DSTenant{}
		if err := resultRows.StructScan(&dsTenantID); err != nil {
			return nil, errors.New("scanning delivery service tenant ID: " + err.Error())
		}
		dsTenantIDs = append(dsTenantIDs, dsTenantID)
	}

	return dsTenantIDs, nil
}

// AssignMultipleServersCapabilities assigns multiple servers to a capability or multiple server capabilities to a server
func AssignMultipleServersCapabilities(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	tx := inf.Tx.Tx
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	var mssc tc.MultipleServersCapabilities
	if err := json.NewDecoder(r.Body).Decode(&mssc); err != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, fmt.Errorf("error decoding POST request body into MultipleServersCapabilities struct %w", err), nil)
		return
	}

	// validate JSON body.
	errs := tovalidate.ToErrors(validation.Errors{
		"serverIds":          validation.Validate(mssc.ServerIDs, validation.Required),
		"serverCapabilities": validation.Validate(mssc.ServerCapabilities, validation.Required),
		"pageType":           validation.Validate(mssc.PageType, validation.Required),
	})

	if len(errs) > 0 {
		api.HandleErr(w, r, tx, http.StatusBadRequest, util.JoinErrs(errs), nil)
		return
	}

	if len(mssc.ServerIDs) > 1 && len(mssc.ServerCapabilities) > 1 {
		api.HandleErr(w, r, tx, http.StatusBadRequest, fmt.Errorf("not allowed to have many:many association between server and server capability. "+
			"Only associations allowed are; 1:1, 1:many or many:1"), nil)
		return
	}

	if len(mssc.ServerIDs) >= 1 {
		errCode, userErr, sysErr = checkExistingServer(tx, mssc.ServerIDs, inf.User.UserName)
		if userErr != nil || sysErr != nil {
			api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
			return
		}
	}

	// Ensure type is correct
	errCode, userErr, sysErr = checkServerType(tx, mssc.ServerIDs)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}

	// Insert rows in DB
	sid := make([]int64, len(mssc.ServerCapabilities))
	scs := make([]string, len(mssc.ServerIDs))
	switch mssc.PageType {
	case "sc":
		for i := range mssc.ServerIDs {
			scs[i] = mssc.ServerCapabilities[0]
		}
		sid = mssc.ServerIDs
	case "server":
		for i := range mssc.ServerCapabilities {
			sid[i] = mssc.ServerIDs[0]
		}
		scs = mssc.ServerCapabilities
	default:
		api.HandleErr(w, r, tx, http.StatusBadRequest, fmt.Errorf("incorrect page type: '%s'. Should be 'sc' or 'server'", mssc.PageType), nil)
		return
	}

	msscQuery := `INSERT INTO server_server_capability
			select "server_capability", "server"
			FROM UNNEST($1::text[], $2::int[]) AS tmp("server_capability", "server")`
	_, err := tx.Query(msscQuery, pq.Array(scs), pq.Array(sid))
	if err != nil {
		useErr, sysErr, statusCode := api.ParseDBError(err)
		api.HandleErr(w, r, tx, statusCode, useErr, sysErr)
		return
	}

	var alerts tc.Alerts
	if mssc.PageType == "sc" {
		alerts = tc.CreateAlerts(tc.SuccessLevel, "Assign Server(s) to a capability")
	} else {
		alerts = tc.CreateAlerts(tc.SuccessLevel, "Assign Server Capability(ies) to a server")
	}
	api.WriteAlertsObj(w, r, http.StatusOK, alerts, mssc)
	return
}

// DeleteMultipleServersCapabilities deletes multiple servers to a capability or multiple server capabilities to a server
func DeleteMultipleServersCapabilities(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	tx := inf.Tx.Tx
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	var mssc tc.MultipleServersCapabilities
	if err := json.NewDecoder(r.Body).Decode(&mssc); err != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, fmt.Errorf("error decoding DELETE request body into MultipleServersCapabilities struct %w", err), nil)
		return
	}

	if len(mssc.ServerIDs) >= 1 {
		errCode, userErr, sysErr = checkExistingServer(tx, mssc.ServerIDs, inf.User.UserName)
		if userErr != nil || sysErr != nil {
			api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
			return
		}
	}

	//Delete existing rows from server_server_capability for a given server or for a given capability
	const delQuery = `DELETE FROM server_server_capability ssc WHERE `
	var dq string
	var alerts tc.Alerts
	var result sql.Result
	var err error
	switch mssc.PageType {
	case "sc":
		dq = delQuery + `ssc.server_capability=$1`
		if len(mssc.ServerIDs) == 1 {
			dq = dq + ` AND ssc.server=$2`
			result, err = tx.Exec(dq, mssc.ServerCapabilities[0], mssc.ServerIDs[0])
		} else {
			result, err = tx.Exec(dq, mssc.ServerCapabilities[0])
		}
	case "server":
		dq = delQuery + `ssc.server=$1`
		if len(mssc.ServerCapabilities) == 1 {
			dq = dq + ` AND ssc.server_capability=$2`
			result, err = tx.Exec(dq, mssc.ServerIDs[0], mssc.ServerCapabilities[0])
		} else {
			result, err = tx.Exec(dq, mssc.ServerIDs[0])
		}
	default:
		api.HandleErr(w, r, tx, http.StatusBadRequest, fmt.Errorf("incorrect page type:'%s'. Should be 'sc' or 'server'", mssc.PageType), nil)
		return
	}

	if err != nil {
		useErr, sysErr, statusCode := api.ParseDBError(err)
		api.HandleErr(w, r, tx, statusCode, useErr, sysErr)
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, fmt.Errorf("no rows were deleted from server_server_capability table: %w", err), sysErr)
		return
	}
	if rowsAffected >= 1 {
		if mssc.PageType == "sc" {
			alerts = tc.CreateAlerts(tc.SuccessLevel, "Removed Server(s) associated with a capability")
		} else {
			alerts = tc.CreateAlerts(tc.SuccessLevel, "Removed Server Capability(ies) associated with a server")
		}
	}

	api.WriteAlertsObj(w, r, http.StatusOK, alerts, mssc)
	return
}

// checkExistingServer checks server existence
func checkExistingServer(tx *sql.Tx, sidList []int64, uName string) (int, error, error) {
	for _, sid := range sidList {
		_, exists, err := dbhelpers.GetServerNameFromID(tx, sid)
		if err != nil {
			return http.StatusInternalServerError, nil, err
		}
		if !exists {
			userErr := fmt.Errorf("server %d does not exist", sid)
			return http.StatusNotFound, userErr, nil
		}

		cdnName, err := dbhelpers.GetCDNNameFromServerID(tx, sid)
		if err != nil {
			return http.StatusInternalServerError, nil, err
		}

		userErr, sysErr, errCode := dbhelpers.CheckIfCurrentUserCanModifyCDN(tx, string(cdnName), uName)
		if userErr != nil || sysErr != nil {
			return errCode, userErr, sysErr
		}
	}
	return http.StatusOK, nil, nil
}

// checkServerType checks if the server type is MID and/or EDGE
func checkServerType(tx *sql.Tx, sids []int64) (int, error, error) {
	var servArray []int64
	queryType := `SELECT array_agg(s.id) 
		FROM server s
		JOIN type t ON s.type = t.id
		WHERE s.id = any ($1)
		AND t.use_in_table = 'server'
		AND (t.name LIKE 'MID%' OR t.name LIKE 'EDGE%')`
	if err := tx.QueryRow(queryType, pq.Array(sids)).Scan(pq.Array(&servArray)); err != nil {
		return http.StatusInternalServerError, nil, fmt.Errorf("checking server type: %w", err)
	}
	cmp := make(map[int64]bool)
	for _, item := range servArray {
		cmp[item] = true
	}
	for _, sid := range sids {
		if _, ok := cmp[sid]; !ok {
			userErr := fmt.Errorf("server id: %d has an incorrect server type. Server capabilities can only be assigned to EDGE or MID servers", sid)
			return http.StatusBadRequest, userErr, nil
		}
	}
	return http.StatusOK, nil, nil
}
