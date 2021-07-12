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

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"

	"github.com/jmoiron/sqlx"
)

// UpdateHandler implements an http handler that updates a server's upd_pending and reval_pending values.
func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"id-or-name"}, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	idOrName := inf.Params["id-or-name"]
	id, err := strconv.Atoi(idOrName)
	hostName := ""
	if err == nil {
		name, ok, err := dbhelpers.GetServerNameFromID(inf.Tx.Tx, id)
		if err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting server name from id '"+idOrName+"': "+err.Error()))
			return
		} else if !ok {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusNotFound, errors.New("server ID '"+idOrName+"' not found"), nil)
			return
		}
		hostName = name
		cdnName, err := dbhelpers.GetCDNNameFromServerID(inf.Tx.Tx, int64(id))
		if err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, err)
			return
		}
		userErr, sysErr, statusCode := dbhelpers.CheckIfCurrentUserHasCdnLock(inf.Tx.Tx, string(cdnName), inf.User.UserName)
		if userErr != nil || sysErr != nil {
			api.HandleErr(w, r, inf.Tx.Tx, statusCode, userErr, sysErr)
			return
		}
	} else {
		hostName = idOrName
		serverID, ok, err := dbhelpers.GetServerIDFromName(hostName, inf.Tx.Tx)
		if err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting server id from name '"+idOrName+"': "+err.Error()))
			return
		} else if !ok {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusNotFound, errors.New("server name '"+idOrName+"' not found"), nil)
			return
		}
		cdnName, err := dbhelpers.GetCDNNameFromServerID(inf.Tx.Tx, int64(serverID))
		if err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, err)
			return
		}
		userErr, sysErr, statusCode := dbhelpers.CheckIfCurrentUserHasCdnLock(inf.Tx.Tx, string(cdnName), inf.User.UserName)
		if userErr != nil || sysErr != nil {
			api.HandleErr(w, r, inf.Tx.Tx, statusCode, userErr, sysErr)
			return
		}
	}

	updated, hasUpdated := inf.Params["updated"]
	revalUpdated, hasRevalUpdated := inf.Params["reval_updated"]
	if !hasUpdated && !hasRevalUpdated {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("Must pass at least one query paramter of 'updated' or 'reval_updated'"), nil)
		return
	}
	updated = strings.ToLower(updated)
	revalUpdated = strings.ToLower(revalUpdated)

	if hasUpdated && updated != `t` && updated != `true` && updated != `f` && updated != `false` {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("query parameter 'updated' must be 'true' or 'false'"), nil)
		return
	}
	if hasRevalUpdated && revalUpdated != `t` && revalUpdated != `true` && revalUpdated != `f` && revalUpdated != `false` {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("query parameter 'reval_updated' must be 'true' or 'false'"), nil)
		return
	}

	strToBool := func(s string) bool {
		return !strings.HasPrefix(strings.ToLower(s), "f")
	}

	updatedPtr := (*bool)(nil)
	if hasUpdated {
		updatedBool := strToBool(updated)
		updatedPtr = &updatedBool
	}
	revalUpdatedPtr := (*bool)(nil)
	if hasRevalUpdated {
		revalUpdatedBool := strToBool(revalUpdated)
		revalUpdatedPtr = &revalUpdatedBool
	}

	if err := setUpdateStatuses(inf.Tx.Tx, hostName, updatedPtr, revalUpdatedPtr); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("setting updated statuses: "+err.Error()))
		return
	}

	respMsg := "successfully set server '" + hostName + "'"
	if hasUpdated {
		respMsg += " updated=" + strconv.FormatBool(strToBool(updated))
	}
	if hasRevalUpdated {
		respMsg += " reval_updated=" + strconv.FormatBool(strToBool(revalUpdated))
	}

	api.WriteAlerts(w, r, http.StatusOK, tc.CreateAlerts(tc.SuccessLevel, respMsg))
}

// setUpdateStatuses sets the upd_pending and reval_pending columns of a server.
// If updatePending or revalPending is nil, that value is not changed.
func setUpdateStatuses(tx *sql.Tx, hostName string, updatePending *bool, revalPending *bool) error {
	if updatePending == nil && revalPending == nil {
		return errors.New("either updatePending or revalPending must not be nil")
	}
	qry := `UPDATE server SET `
	updateStrs := []string{}
	nextI := 1
	qryVals := []interface{}{}
	if updatePending != nil {
		updateStrs = append(updateStrs, `upd_pending = $`+strconv.Itoa(nextI))
		nextI++
		qryVals = append(qryVals, *updatePending)
	}
	if revalPending != nil {
		updateStrs = append(updateStrs, `reval_pending = $`+strconv.Itoa(nextI))
		nextI++
		qryVals = append(qryVals, *revalPending)
	}
	qry += strings.Join(updateStrs, ", ") + ` WHERE host_name = $` + strconv.Itoa(nextI)
	qryVals = append(qryVals, hostName)

	if _, err := tx.Exec(qry, qryVals...); err != nil {
		return errors.New("executing: " + err.Error())
	}
	return nil
}

// ProfileAndTypeQueueUpdateHandler queues/ dequeues updates on servers for a particular CDN, filtered by type and/ or profile.
func ProfileAndTypeQueueUpdateHandler(w http.ResponseWriter, r *http.Request) {
	var cdnID int
	var typeID int
	var profileID = -1
	var ok bool
	var err error
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"cdn"}, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	cols := map[string]dbhelpers.WhereColumnInfo{
		"cdnID":     {Column: "server.cdn_id", Checker: nil},
		"typeID":    {Column: "server.type", Checker: nil},
		"profileID": {Column: "server.profile", Checker: nil},
	}

	cdn := inf.Params["cdn"]
	typeName := inf.Params["type"]
	profile := inf.Params["profile"]

	var reqObj tc.ServerQueueUpdateRequest

	if err := json.NewDecoder(r.Body).Decode(&reqObj); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, fmt.Errorf("malformed JSON: %v", err), nil)
		return
	}

	if reqObj.Action != "queue" && reqObj.Action != "dequeue" {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("action must be 'queue' or 'dequeue'"), nil)
		return
	}
	queue := reqObj.Action == "queue"

	// get cdn ID
	cdnID, ok, err = dbhelpers.GetCDNIDFromName(inf.Tx.Tx, tc.CDNName(cdn))
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("error getting CDN ID from name: "+err.Error()))
		return
	}
	if !ok {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusNotFound, errors.New("no CDN ID found with that name"), nil)
		return
	}
	delete(inf.Params, "cdn")
	inf.Params["cdnID"] = strconv.Itoa(cdnID)

	// get type ID
	if typeName != "" {
		typeID, ok, err = dbhelpers.GetTypeIDByName(typeName, inf.Tx.Tx)
		if err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("error getting type ID from name: "+err.Error()))
			return
		}
		if !ok {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusNotFound, errors.New("no type ID found with that name"), nil)
			return
		}
		inf.Params["typeID"] = strconv.Itoa(typeID)
	}
	delete(inf.Params, "type")

	// get profile ID
	if profile != "" {
		profileID, ok, err = dbhelpers.GetProfileIDFromName(profile, inf.Tx.Tx)
		if err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("error getting profile ID from name: "+err.Error()))
			return
		}
		if !ok {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusNotFound, errors.New("no profile ID found with that name"), nil)
			return
		}
		inf.Params["profileID"] = strconv.Itoa(profileID)
	}
	delete(inf.Params, "profile")

	where, orderBy, pagination, queryValues, errs := dbhelpers.BuildWhereAndOrderByAndPagination(inf.Params, cols)
	if len(errs) > 0 {
		errCode = http.StatusBadRequest
		userErr = util.JoinErrs(errs)
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, nil)
		return
	}

	userErr, sysErr, statusCode := dbhelpers.CheckIfCurrentUserHasCdnLock(inf.Tx.Tx, string(cdn), inf.User.UserName)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, statusCode, userErr, sysErr)
		return
	}
	query := `UPDATE server SET upd_pending = :upd_pending`
	query = query + where + orderBy + pagination
	queryValues["upd_pending"] = queue
	ok, err = queueUpdatesByTypeOrProfile(inf.Tx, queryValues, query)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("queueing updates: %v", err))
		return
	}
	if !ok {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusNotFound, fmt.Errorf("no server with the given combination found"), nil)
		return
	}
	api.CreateChangeLogRawTx(api.ApiChange, "CDN: "+cdn+", Type: "+typeName+", ACTION: CDN server updates "+reqObj.Action+"d", inf.User, inf.Tx.Tx)
	api.WriteResp(w, r, tc.ServerGenericQueueUpdateResponse{Action: reqObj.Action, CDNID: cdnID, TypeID: typeID, ProfileID: profileID})
}

// queueUpdatesByTypeOrProfile is the helper function to queue/ dequeue updates on servers for a CDN, filtered by type and/ or profile
func queueUpdatesByTypeOrProfile(tx *sqlx.Tx, queryValues map[string]interface{}, query string) (bool, error) {
	result, err := tx.NamedExec(query, queryValues)
	if err != nil {
		return false, errors.New("querying generic queue updates: " + err.Error())
	} else if rc, err := result.RowsAffected(); err != nil {
		return false, fmt.Errorf("checking rows updated: %v", err)
	} else {
		return rc > 0, nil
	}
}
