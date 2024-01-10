package cdni

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
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/lib/pq"
)

const (
	CapabilityQuery   = `SELECT id, type, ucdn FROM cdni_capabilities WHERE type = $1 AND ucdn = $2`
	AllFootprintQuery = `SELECT footprint_type, footprint_value::text[], capability_id FROM cdni_footprints`

	limitsQuery = `
SELECT limit_id, scope_type, scope_value, limit_type, maximum_hard, maximum_soft, cl.telemetry_id, cl.telemetry_metric, t.id, t.type, tm.name, cl.capability_id
FROM cdni_limits AS cl
LEFT JOIN cdni_telemetry as t ON telemetry_id = t.id
LEFT JOIN cdni_telemetry_metrics as tm ON telemetry_metric = tm.name`

	InsertCapabilityUpdateQuery     = `INSERT INTO cdni_capability_updates (ucdn, data, async_status_id, request_type, host) VALUES ($1, $2, $3, $4, $5)`
	SelectCapabilityUpdateQuery     = `SELECT ucdn, data, async_status_id, request_type, host FROM cdni_capability_updates WHERE id = $1`
	SelectAllCapabilityUpdatesQuery = `SELECT id, ucdn, data, request_type, host FROM cdni_capability_updates`

	DeleteCapabilityUpdateQuery               = `DELETE FROM cdni_capability_updates WHERE id = $1`
	UpdateLimitsByCapabilityAndLimitTypeQuery = `UPDATE cdni_limits SET maximum_hard = $1 WHERE capability_id = $2 AND limit_type = $3`
	hostQuery                                 = `SELECT count(*) FROM cdni_limits WHERE $1 = ANY(scope_value)`

	hostConfigLabel = "hostConfigUpdate"
)

// GetCapabilities returns the CDNi capability limits.
func GetCapabilities(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	if inf.Config.Cdni == nil || inf.Config.Secrets[0] == "" || inf.Config.Cdni.DCdnId == "" {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("cdn.conf does not contain CDNi information"))
		return
	}

	bearerToken := getBearerToken(r)

	ucdn, err := checkBearerToken(bearerToken, inf)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, err, nil)
		return
	}

	capacities, err := getCapacities(inf, ucdn)
	if err != nil {
		api.HandleErr(w, r, nil, http.StatusInternalServerError, err, nil)
		return
	}

	telemetries, err := getTelemetries(inf, ucdn)
	if err != nil {
		api.HandleErr(w, r, nil, http.StatusInternalServerError, err, nil)
		return
	}

	fciCaps := Capabilities{}
	capsList := make([]Capability, 0, len(capacities.Capabilities)+len(telemetries.Capabilities))
	capsList = append(capsList, capacities.Capabilities...)
	capsList = append(capsList, telemetries.Capabilities...)

	fciCaps.Capabilities = capsList

	api.WriteRespRaw(w, r, fciCaps)
}

func getBearerToken(r *http.Request) string {
	if r.Header.Get(rfc.Authorization) != "" && strings.Contains(r.Header.Get(rfc.Authorization), "Bearer") {
		givenTokenSplit := strings.Split(r.Header.Get(rfc.Authorization), " ")
		if len(givenTokenSplit) < 2 {
			return ""
		}

		return givenTokenSplit[1]
	}
	for _, cookie := range r.Cookies() {
		switch cookie.Name {
		case rfc.AccessToken:
			return cookie.Value
		}
	}
	return ""
}

// PutHostConfiguration adds the requested CDNi configuration update for a specific host to the queue and adds an async status.
func PutHostConfiguration(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"host"}, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	host := inf.Params["host"]
	if errCode, userErr, sysErr := validateHostExists(host, inf.Tx.Tx); userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}

	if inf.Config.Cdni == nil || inf.Config.Secrets[0] == "" || inf.Config.Cdni.DCdnId == "" {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("cdn.conf does not contain CDNi information"))
		return
	}

	bearerToken := getBearerToken(r)
	ucdn, err := checkBearerToken(bearerToken, inf)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, err, nil)
		return
	}

	var genericHostRequest GenericHostMetadata
	err = json.NewDecoder(r.Body).Decode(&genericHostRequest)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("decoding host json request: %w", err))
		return
	}

	db, err := api.GetDB(r.Context())
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("getting async db: %w", err))
		return
	}
	asyncTx, err := db.Begin()
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("getting async tx: %w", err))
		return
	}
	logTx, err := db.Begin()
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("getting log tx: %w", err))
		return
	}
	defer logTx.Commit()

	asyncStatusId, errCode, userErr, sysErr := api.InsertAsyncStatus(asyncTx, "CDNi host configuration update request received.")
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}

	data := genericHostRequest.HostMetadata.Metadata

	_, err = inf.Tx.Tx.Query(InsertCapabilityUpdateQuery, ucdn, data, asyncStatusId, hostConfigLabel, host)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("inserting capability update request into queue: %w", err))
		return
	}

	msg := "CDNi configuration update request received. Status updates can be found here: " + api.CurrentAsyncEndpoint + strconv.Itoa(asyncStatusId)
	api.CreateChangeLogRawTx(api.ApiChange, msg, inf.User, logTx)

	var alerts tc.Alerts
	alerts.AddAlert(tc.Alert{
		Text:  msg,
		Level: tc.SuccessLevel.String(),
	})

	w.Header().Add(rfc.Location, api.CurrentAsyncEndpoint+strconv.Itoa(asyncStatusId))
	api.WriteAlerts(w, r, http.StatusAccepted, alerts)
}

// PutConfiguration adds the requested CDNi configuration update to the queue and adds an async status.
func PutConfiguration(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	if inf.Config.Cdni == nil || inf.Config.Secrets[0] == "" || inf.Config.Cdni.DCdnId == "" {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("cdn.conf does not contain CDNi information"))
		return
	}

	bearerToken := getBearerToken(r)
	ucdn, err := checkBearerToken(bearerToken, inf)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, err, nil)
		return
	}

	var genericRequest GenericRequestMetadata
	err = json.NewDecoder(r.Body).Decode(&genericRequest)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("decoding json request: %w", err))
		return
	}

	db, err := api.GetDB(r.Context())
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("getting async db: %w", err))
		return
	}
	asyncTx, err := db.Begin()
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("getting async tx: %w", err))
		return
	}
	logTx, err := db.Begin()
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("getting log tx: %w", err))
		return
	}
	defer logTx.Commit()

	asyncStatusId, errCode, userErr, sysErr := api.InsertAsyncStatus(asyncTx, "CDNi configuration update request received.")
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}

	data := genericRequest.Metadata

	_, err = inf.Tx.Tx.Query(InsertCapabilityUpdateQuery, ucdn, data, asyncStatusId, SupportedGenericMetadataType(genericRequest.Type), genericRequest.Host)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("inserting capability update request into queue: %w", err))
		return
	}

	msg := "CDNi configuration update request received. Status updates can be found here: " + api.CurrentAsyncEndpoint + strconv.Itoa(asyncStatusId)
	api.CreateChangeLogRawTx(api.ApiChange, msg, inf.User, logTx)
	var alerts tc.Alerts
	alerts.AddAlert(tc.Alert{
		Text:  msg,
		Level: tc.SuccessLevel.String(),
	})

	w.Header().Add(rfc.Location, api.CurrentAsyncEndpoint+strconv.Itoa(asyncStatusId))
	api.WriteAlerts(w, r, http.StatusAccepted, alerts)
}

// GetRequests returns the CDNi configuration update requests.
func GetRequests(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	var rows *sql.Rows
	var err error

	idParam := inf.Params["id"]
	if idParam != "" {
		id, parseErr := strconv.Atoi(idParam)
		if parseErr != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("id must be an integer"), nil)
			return
		}
		rows, err = inf.Tx.Tx.Query(SelectAllCapabilityUpdatesQuery+" WHERE id = $1", id)
	} else {
		rows, err = inf.Tx.Tx.Query(SelectAllCapabilityUpdatesQuery)
	}

	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("querying for capability update requests: %w", err))
		return
	}
	defer log.Close(rows, "closing capabilities update query")
	requests := []ConfigurationUpdateRequest{}
	for rows.Next() {
		var request ConfigurationUpdateRequest
		if err := rows.Scan(&request.ID, &request.UCDN, &request.Data, &request.RequestType, &request.Host); err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("scanning db rows: %w", err))
			return
		}
		requests = append(requests, request)
	}

	api.WriteResp(w, r, requests)

}

// PutConfigurationResponse approves or denies a CDNi configuration request and updates the configuration and async status appropriately.
func PutConfigurationResponse(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"approved", "id"}, []string{"id"})
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	reqId := inf.IntParams["id"]
	approvedString := inf.Params["approved"]
	approved, err := strconv.ParseBool(approvedString)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("approved parameter must be a boolean"), nil)
		return
	}

	db, err := api.GetDB(r.Context())
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("getting async db: %w", err))
		return
	}

	logTx, err := db.Begin()
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("getting log tx: %w", err))
		return
	}
	defer logTx.Commit()

	rows, err := inf.Tx.Tx.Query(SelectCapabilityUpdateQuery, reqId)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("querying for capability update request: %w", err))
		return
	}
	defer log.Close(rows, "closing capabilities update query")
	var ucdn string
	var data json.RawMessage
	var host string
	var asyncId int
	var requestType string
	count := 0
	for rows.Next() {
		if err := rows.Scan(&ucdn, &data, &asyncId, &requestType, &host); err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("scanning db rows: %w", err))
			return
		}
		count++
	}
	if count == 0 {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusNotFound, fmt.Errorf("no configuration request for that id"), nil)
		return
	}

	if !approved {
		if asyncErr := api.UpdateAsyncStatus(db, api.AsyncFailed, "Requested configuration update has been denied.", asyncId, true); asyncErr != nil {
			log.Errorf("updating async status for id %d: %s", asyncId, asyncErr.Error())
		}
		status, err := deleteCapabilityRequest(reqId, inf.Tx.Tx)
		if err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, status, nil, fmt.Errorf("deleting configuration request from queue: %w", err))
			return
		}
		msg := "Successfully denied configuration update request."
		api.CreateChangeLogRawTx(api.ApiChange, msg, inf.User, inf.Tx.Tx)
		api.WriteResp(w, r, msg)
		return
	}

	var updatedDataList []GenericMetadata
	if err = json.Unmarshal(data, &updatedDataList); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("unmarshalling data for configuration update: %w", err))
		return
	}

	var unsupportedTypes []string
	for _, updatedData := range updatedDataList {
		if !updatedData.Type.isValid() {
			unsupportedTypes = append(unsupportedTypes, string(updatedData.Type))
		}
	}

	if len(unsupportedTypes) != 0 {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, fmt.Errorf("unsupported generic metadata types found: %v", strings.Join(unsupportedTypes, ", ")), nil)
		return
	}

	for _, updatedData := range updatedDataList {
		switch updatedData.Type {
		case MiRequestedCapacityLimits:
			var capacityRequestedLimits CapacityRequestedLimits
			if err = json.Unmarshal(updatedData.Value, &capacityRequestedLimits); err != nil {
				api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("unmarshalling data for configuration update: %w", err))
				return
			}
			for _, capLim := range capacityRequestedLimits.RequestedLimits {
				capId, err := getCapabilityIdFromFootprints(capLim, ucdn, inf)
				if err != nil {
					api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("finding capability for given information: %w", err))
					return
				}

				query := UpdateLimitsByCapabilityAndLimitTypeQuery
				queryParams := []interface{}{capLim.LimitValue, capId, capLim.LimitType}
				if host != "" {
					query = query + " AND $4 = ANY(scope_value)"
					queryParams = []interface{}{capLim.LimitValue, capId, capLim.LimitType, host}
				}

				result, err := inf.Tx.Tx.Exec(query, queryParams...)
				if err != nil {
					api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("updating capacity: %w", err))
					return
				}

				if rowsAffected, err := result.RowsAffected(); err != nil {
					api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("updating capacity: getting rows affected: %w", err))
					return
				} else if rowsAffected < 1 {
					api.HandleErr(w, r, inf.Tx.Tx, http.StatusNotFound, fmt.Errorf("no capacity found for update: host: %s, type: %s, limit: %v", host, updatedData.Type, capLim), nil)
					return
				} else if rowsAffected > 1 {
					api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("capacity update affected too many rows: %d", rowsAffected))
					return
				}
			}
		}
	}

	if asyncErr := api.UpdateAsyncStatus(db, api.AsyncSucceeded, "Capacity requested update has been completed.", asyncId, true); asyncErr != nil {
		log.Errorf("updating async status for id %v: %v", asyncId, asyncErr)
	}
	status, err := deleteCapabilityRequest(reqId, inf.Tx.Tx)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, status, nil, fmt.Errorf("deleting capacity request from queue: %w", err))
		return
	}
	msg := "Successfully updated configuration."
	api.CreateChangeLogRawTx(api.ApiChange, msg, inf.User, logTx)
	api.WriteResp(w, r, msg)
}

func getCapabilityIdFromFootprints(updatedData CapacityLimit, ucdn string, inf *api.Info) (int, error) {
	tableAbbr := ""
	selectClause := ""
	whereClause := ""
	var queryParams []interface{}
	paramCount := 1

	for i, footprint := range updatedData.Footprints {
		if i == 0 {
			tableAbbr = "f"
			selectClause = "SELECT " + tableAbbr + ".capability_id FROM cdni_footprints as " + tableAbbr
			whereClause = " WHERE " + tableAbbr + ".ucdn = $" + strconv.Itoa(paramCount) + " AND " + tableAbbr + ".footprint_type = $" + strconv.Itoa(paramCount+1) + " AND " + tableAbbr + ".footprint_value = $" + strconv.Itoa(paramCount+2) + "::text[]"
		} else {
			oldTableAbbr := tableAbbr
			tableAbbr = tableAbbr + "f"
			selectClause = selectClause + " JOIN cdni_footprints as " + tableAbbr + " on " + tableAbbr + ".capability_id = " + oldTableAbbr + ".capability_id"
			whereClause = whereClause + " AND " + tableAbbr + ".ucdn = $" + strconv.Itoa(paramCount) + " AND " + tableAbbr + ".footprint_type = $" + strconv.Itoa(paramCount+1) + " AND " + tableAbbr + ".footprint_value = $" + strconv.Itoa(paramCount+2) + "::text[]"
		}
		paramCount = paramCount + 3
		queryParams = append(queryParams, ucdn)
		queryParams = append(queryParams, footprint.FootprintType)
		queryParams = append(queryParams, pq.Array(footprint.FootprintValue))
	}

	selectQuery := selectClause + whereClause + " AND (SELECT count(*) from cdni_footprints as c where c.capability_id = f.capability_id) = " + strconv.Itoa(len(updatedData.Footprints))
	rows, err := inf.Tx.Tx.Query(selectQuery, queryParams...)
	if err != nil {
		return 0, fmt.Errorf("querying for capacity update request: %w", err)
	}
	defer log.Close(rows, "closing footprints query")
	var capabilityIds []int
	rowCount := 0
	for rows.Next() {
		var capabilityId int
		if err := rows.Scan(&capabilityId); err != nil {
			return 0, fmt.Errorf("scanning db rows: %w", err)
		}
		rowCount++
		capabilityIds = append(capabilityIds, capabilityId)
	}

	if len(capabilityIds) == 0 {
		return 0, fmt.Errorf("no capabilities found that match all footprints: %v", updatedData.Footprints)
	}
	if len(capabilityIds) > 1 {
		return 0, fmt.Errorf("more than 1 capability found that match all footprints: %v", updatedData.Footprints)
	}
	return capabilityIds[0], nil
}

func deleteCapabilityRequest(reqId int, tx *sql.Tx) (int, error) {
	result, err := tx.Exec(DeleteCapabilityUpdateQuery, reqId)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("deleting configuration update: %w", err)
	}

	if rowsAffected, err := result.RowsAffected(); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("deleting configuration update: getting rows affected: %w", err)
	} else if rowsAffected < 1 {
		return http.StatusNotFound, errors.New("no configuration update with that key found")
	} else if rowsAffected > 1 {
		return http.StatusInternalServerError, fmt.Errorf("delete affected too many rows: %d", rowsAffected)
	}

	return http.StatusOK, nil
}

func validateHostExists(host string, tx *sql.Tx) (int, error, error) {
	count := 0
	if err := tx.QueryRow(hostQuery, host).Scan(&count); err != nil {
		return http.StatusInternalServerError, nil, fmt.Errorf("querying if host %s exists: %w", host, err)
	}
	if count == 0 {
		return http.StatusBadRequest, fmt.Errorf("No data found for host: %s", host), nil
	}
	return http.StatusOK, nil, nil
}

func checkBearerToken(bearerToken string, inf *api.Info) (string, error) {
	if bearerToken == "" {
		return "", errors.New("bearer token is required")
	}

	token, err := jwt.Parse([]byte(bearerToken),
		jwt.WithVerify(jwa.HS256, []byte(inf.Config.Secrets[0])),
	)
	if err != nil {
		return "", fmt.Errorf("invalid token: %w", err)
	}

	if token.Expiration().Unix() < time.Now().Unix() {
		return "", errors.New("token is expired")
	}

	if token.Audience() == nil || len(token.Audience()) == 0 {
		return "", errors.New("invalid token - dcdn must be defined in audience claim")
	}
	if token.Audience()[0] != inf.Config.Cdni.DCdnId {
		return "", errors.New("invalid token - incorrect dcdn")
	}

	ucdn := token.Issuer()
	if ucdn != inf.User.UCDN {
		return "", errors.New("user ucdn did not match token ucdn")
	}

	if ucdn == "" {
		if inf.User.Can(tc.PermICDNUCDNOverride) {
			ucdn = inf.Params["ucdn"]
			if ucdn == "" {
				return "", errors.New("admin level ucdn requests require a ucdn query parameter")
			}
		} else {
			return "", errors.New("invalid token - empty ucdn field")
		}
	}

	return ucdn, nil
}

func getFootprintMap(tx *sql.Tx) (map[int][]Footprint, error) {
	footRows, err := tx.Query(AllFootprintQuery)
	if err != nil {
		return nil, fmt.Errorf("querying footprints: %w", err)
	}
	defer log.Close(footRows, "closing foorpint query")
	footprintMap := map[int][]Footprint{}
	for footRows.Next() {
		var footprint Footprint
		if err := footRows.Scan(&footprint.FootprintType, pq.Array(&footprint.FootprintValue), &footprint.CapabilityId); err != nil {
			return nil, fmt.Errorf("scanning db rows: %w", err)
		}

		footprintMap[footprint.CapabilityId] = append(footprintMap[footprint.CapabilityId], footprint)
	}

	return footprintMap, nil
}

func getLimitsMap(tx *sql.Tx) (map[int][]LimitsQueryResponse, error) {
	rows, err := tx.Query(limitsQuery)
	if err != nil {
		return nil, fmt.Errorf("querying limits: %w", err)
	}

	defer log.Close(rows, "closing capacity limits query")
	limitsMap := map[int][]LimitsQueryResponse{}
	for rows.Next() {
		var limit LimitsQueryResponse
		var scope LimitScope
		if err := rows.Scan(&limit.LimitId, &scope.ScopeType, pq.Array(&scope.ScopeValue), &limit.LimitType, &limit.MaximumHard, &limit.MaximumSoft, &limit.TelemetryId, &limit.TelemetryMetic, &limit.Id, &limit.Type, &limit.Name, &limit.CapabilityId); err != nil {
			return nil, fmt.Errorf("scanning db rows: %w", err)
		}
		if scope.ScopeType != nil {
			limit.Scope = &scope
		}

		limitsMap[limit.CapabilityId] = append(limitsMap[limit.CapabilityId], limit)
	}

	return limitsMap, nil
}

func getTelemetriesMap(tx *sql.Tx) (map[int][]Telemetry, error) {
	rows, err := tx.Query(`SELECT id, type, capability_id, configuration_url FROM cdni_telemetry`)
	if err != nil {
		return nil, errors.New("querying cdni telemetry: " + err.Error())
	}
	defer log.Close(rows, "closing telemetry query")

	telemetryMap := map[int][]Telemetry{}
	for rows.Next() {
		telemetry := Telemetry{}
		if err := rows.Scan(&telemetry.Id, &telemetry.Type, &telemetry.CapabilityId, &telemetry.Configuration.Url); err != nil {
			return nil, errors.New("scanning telemetry: " + err.Error())
		}

		telemetryMap[telemetry.CapabilityId] = append(telemetryMap[telemetry.CapabilityId], telemetry)
	}

	return telemetryMap, nil
}

func getTelemetryMetricsMap(tx *sql.Tx) (map[string][]Metric, error) {
	tmRows, err := tx.Query(`SELECT name, time_granularity, data_percentile, latency, telemetry_id FROM cdni_telemetry_metrics`)
	if err != nil {
		return nil, errors.New("querying cdni telemetry metrics: " + err.Error())
	}
	defer log.Close(tmRows, "closing telemetry metrics query")

	telemetryMetricMap := map[string][]Metric{}
	for tmRows.Next() {
		metric := Metric{}
		if err := tmRows.Scan(&metric.Name, &metric.TimeGranularity, &metric.DataPercentile, &metric.Latency, &metric.TelemetryId); err != nil {
			return nil, errors.New("scanning telemetry metric: " + err.Error())
		}

		telemetryMetricMap[metric.TelemetryId] = append(telemetryMetricMap[metric.TelemetryId], metric)
	}

	return telemetryMetricMap, nil
}

// Capabilities contains an array of CDNi capabilities.
type Capabilities struct {
	Capabilities []Capability `json:"capabilities"`
}

// Capability contains information about a CDNi capability.
type Capability struct {
	CapabilityType  SupportedCapabilities `json:"capability-type"`
	CapabilityValue interface{}           `json:"capability-value"`
	Footprints      []Footprint           `json:"footprints"`
}

// CapacityCapabilityValue contains the total and host capability limits.
type CapacityCapabilityValue struct {
	Limits []Limit `json:"limits"`
}

// Limit contains the information for a capacity limit.
type Limit struct {
	Id              string            `json:"id"`
	Scope           *LimitScope       `json:"scope,omitempty"`
	LimitType       CapacityLimitType `json:"limit-type"`
	MaximumHard     int64             `json:"maximum-hard"`
	MaximumSoft     int64             `json:"maximum-soft"`
	TelemetrySource TelemetrySource   `json:"telemetry-source"`
}

// TelemetrySource contains the information for a telemetry source.
type TelemetrySource struct {
	Id     string `json:"id"`
	Metric string `json:"metric"`
}

// TelemetryCapabilityValue contains an array of telemetry sources.
type TelemetryCapabilityValue struct {
	Sources []Telemetry `json:"sources"`
}

// Telemetry contains the information for a telemetry metric.
type Telemetry struct {
	Id            string                 `json:"id"`
	Type          TelemetrySourceType    `json:"type"`
	CapabilityId  int                    `json:"-"`
	Metrics       []Metric               `json:"metrics"`
	Configuration TelemetryConfiguration `json:"configuration"`
}

type TelemetryConfiguration struct {
	Url string `json:"url"`
}

// Metric contains the metric information for a telemetry metric.
type Metric struct {
	Name            string `json:"name"`
	TimeGranularity int    `json:"time-granularity"`
	DataPercentile  int    `json:"data-percentile"`
	Latency         int    `json:"latency"`
	TelemetryId     string `json:"-"`
}

// Footprint contains the information for a footprint.
type Footprint struct {
	FootprintType  FootprintType `json:"footprint-type" db:"footprint_type"`
	FootprintValue []string      `json:"footprint-value" db:"footprint_value"`
	CapabilityId   int           `json:"-"`
}

// CapacityLimitType is a string of the capacity limit type.
type CapacityLimitType string

const (
	Egress         CapacityLimitType = "egress"
	Requests                         = "requests"
	StorageSize                      = "storage-size"
	StorageObjects                   = "storage-objects"
	Sessions                         = "sessions"
	CacheSize                        = "cache-size"
)

// SupportedCapabilities is a string of the supported capabilities.
type SupportedCapabilities string

const (
	FciTelemetry      SupportedCapabilities = "FCI.Telemetry"
	FciCapacityLimits                       = "FCI.CapacityLimits"
)

// SupportedGenericMetadataType is a string of the supported metadata type.
type SupportedGenericMetadataType string

const (
	MiRequestedCapacityLimits SupportedGenericMetadataType = "MI.RequestedCapacityLimits"
)

func (s SupportedGenericMetadataType) isValid() bool {
	switch s {
	case MiRequestedCapacityLimits:
		return true
	}
	return false
}

// TelemetrySourceType is a string of the telemetry source type. Right now only "generic" is supported.
type TelemetrySourceType string

const (
	Generic TelemetrySourceType = "generic"
)

// FootprintType is a string of the footprint type.
type FootprintType string

const (
	Ipv4Cidr    FootprintType = "ipv4cidr"
	Ipv6Cidr                  = "ipv6cidr"
	Asn                       = "asn"
	CountryCode               = "countrycode"
)

// GenericHostMetadata contains the generic CDNi metadata for a requested update to a specific host.
type GenericHostMetadata struct {
	Host         string           `json:"host"`
	HostMetadata HostMetadataList `json:"host-metadata"`
}

// GenericRequestMetadata contains the generic CDNi metadata for a requested update.
type GenericRequestMetadata struct {
	Type     string          `json:"type"`
	Metadata json.RawMessage `json:"metadata"`
	Host     string          `json:"host,omitempty"`
}

// HostMetadataList contains CDNi metadata for a specific host.
type HostMetadataList struct {
	Metadata json.RawMessage `json:"metadata"`
}

// GenericMetadata contains generic CDNi metadata.
type GenericMetadata struct {
	Type  SupportedGenericMetadataType `json:"generic-metadata-type"`
	Value json.RawMessage              `json:"generic-metadata-value"`
}

// CapacityRequestedLimits contains the requested capacity limits.
type CapacityRequestedLimits struct {
	RequestedLimits []CapacityLimit `json:"requested-limits"`
}

// CapacityLimit contains the limit information for a given footprint.
type CapacityLimit struct {
	LimitType  string      `json:"limit-type"`
	LimitValue int64       `json:"limit-value"`
	Footprints []Footprint `json:"footprints"`
}

// ConfigurationUpdateRequest contains information about a requested CDNi configuration update request.
type ConfigurationUpdateRequest struct {
	ID            int             `json:"id"`
	UCDN          string          `json:"ucdn"`
	Data          json.RawMessage `json:"data"`
	Host          string          `json:"host"`
	RequestType   string          `json:"requestType" db:"request_type"`
	AsyncStatusID int             `json:"asyncStatusId" db:"async_status_id"`
}
