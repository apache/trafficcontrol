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
	"net/http"
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/tenant"
	"github.com/jmoiron/sqlx"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-tc/tovalidate"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/lib/pq"
)

const (
	ServerCapabilityQueryParam = "serverCapability"
	ServerQueryParam           = "serverId"
	ServerHostNameQueryParam   = "serverHostName"
)

type (
	TOServerServerCapability struct {
		api.APIInfoImpl `json:"-"`
		tc.ServerServerCapability
	}

	DSTenant struct {
		TenantID int64 `db:"tenant_id"`
		ID       int64 `db:"id"`
	}
)

func (ssc *TOServerServerCapability) SetLastUpdated(t tc.TimeNoMod) { ssc.LastUpdated = &t }
func (ssc *TOServerServerCapability) NewReadObj() interface{} {
	return &tc.ServerServerCapability{}
}
func (ssc *TOServerServerCapability) SelectQuery() string { return scSelectQuery() }
func (ssc *TOServerServerCapability) ParamColumns() map[string]dbhelpers.WhereColumnInfo {
	return map[string]dbhelpers.WhereColumnInfo{
		ServerCapabilityQueryParam: dbhelpers.WhereColumnInfo{"sc.server_capability", nil},
		ServerQueryParam:           dbhelpers.WhereColumnInfo{"s.id", api.IsInt},
		ServerHostNameQueryParam:   dbhelpers.WhereColumnInfo{"s.host_name", nil},
	}

}
func (ssc *TOServerServerCapability) DeleteQuery() string { return scDeleteQuery() }

func (ssc TOServerServerCapability) GetKeyFieldsInfo() []api.KeyFieldInfo {
	return []api.KeyFieldInfo{
		{ServerQueryParam, api.GetIntKey},
		{ServerCapabilityQueryParam, api.GetStringKey},
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

// Validate fulfills the api.Validator interface
func (ssc TOServerServerCapability) Validate() error {
	errs := validation.Errors{
		ServerQueryParam:           validation.Validate(ssc.ServerID, validation.Required),
		ServerCapabilityQueryParam: validation.Validate(ssc.ServerCapability, validation.Required),
	}

	return util.JoinErrs(tovalidate.ToErrors(errs))
}

func (ssc *TOServerServerCapability) Update() (error, error, int) {
	return nil, nil, http.StatusNotImplemented
}

func (ssc *TOServerServerCapability) Read() ([]interface{}, error, error, int) {
	return api.GenericRead(ssc)
}

func (ssc *TOServerServerCapability) Delete() (error, error, int) {
	// Ensure that the user is not removing a server capability from the server
	// that is required by the delivery services the server is assigned to (if applicable)
	dsIDs := []int64{}
	if err := ssc.APIInfo().Tx.QueryRow(checkDSReqCapQuery(), ssc.ServerID, ssc.ServerCapability).Scan(pq.Array(&dsIDs)); err != nil {
		return nil, fmt.Errorf("checking removing server server capability would still suffice delivery service requried capabilites: %v", err), http.StatusInternalServerError
	}

	if len(dsIDs) > 0 {
		return ssc.buildDSReqCapError(dsIDs)
	}

	// Delete association
	return api.GenericDelete(ssc)
}

func (ssc *TOServerServerCapability) buildDSReqCapError(dsIDs []int64) (error, error, int) {

	dsTenantIDs, err := getDSTenantIDsByIDs(ssc.APIInfo().Tx, dsIDs)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	authDSIDs := []string{}
	checkedTenants := map[int64]bool{}

	for _, dsTenantID := range dsTenantIDs {
		if auth, ok := checkedTenants[dsTenantID.TenantID]; ok { // No need to check tenant again
			if auth {
				authDSIDs = append(authDSIDs, strconv.FormatInt(dsTenantID.ID, 10))
			}
			continue
		}
		authorized, err := tenant.IsResourceAuthorizedToUserTx(int(dsTenantID.TenantID), ssc.APIInfo().User, ssc.APIInfo().Tx.Tx)
		if err != nil {
			return nil, fmt.Errorf("checking tenancy on delivery service: %v", err), http.StatusInternalServerError
		}
		if authorized {
			authDSIDs = append(authDSIDs, strconv.FormatInt(dsTenantID.ID, 10))
		}
		checkedTenants[dsTenantID.TenantID] = authorized
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
	_, exists, err := dbhelpers.GetServerNameFromID(tx.Tx, *ssc.ServerID)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}
	if !exists {
		return fmt.Errorf("server %v does not exist", *ssc.ServerID), nil, http.StatusNotFound
	}

	// Ensure type is correct
	correctType := true
	if err := tx.Tx.QueryRow(scCheckServerTypeQuery(), ssc.ServerID).Scan(&correctType); err != nil {
		return nil, fmt.Errorf("checking server type: %v", err), http.StatusInternalServerError
	}
	if !correctType {
		return fmt.Errorf("server %v has an incorrect server type. Server capabilities can only be assigned to EDGE or MID servers", *ssc.ServerID), nil, http.StatusBadRequest
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

func scCheckServerTypeQuery() string {
	return `
SELECT EXISTS (
	SELECT s.id
	FROM server s
	JOIN type t ON s.type = t.id
	WHERE s.id = $1
	AND t.use_in_table = 'server'
	AND (t.name LIKE 'MID%' OR t.name LIKE 'EDGE%'))`
}

func checkDSReqCapQuery() string {
	return `
SELECT ARRAY(
	SELECT dsrc.deliveryservice_id
	FROM deliveryservices_required_capability as dsrc
	WHERE deliveryservice_id IN (
		SELECT deliveryservice 
		FROM deliveryservice_server
		WHERE server = $1)
	AND dsrc.required_capability = $2)`
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

	for resultRows.Next() {
		dsTenantID := DSTenant{}
		if err := resultRows.StructScan(&dsTenantID); err != nil {
			return nil, errors.New("scanning delivery service tenant ID: " + err.Error())
		}
		dsTenantIDs = append(dsTenantIDs, dsTenantID)
	}

	return dsTenantIDs, nil
}
