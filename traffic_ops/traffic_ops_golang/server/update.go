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
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
)

// UpdateHandler implements an http handler that updates a server's upd_pending and reval_pending values.
//
// Deprecated: As of V4, prefer to use UpdateHandlerV4. This provides legacy compatibility with V3 and lower.
func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"id-or-name"}, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	idOrName := inf.Params["id-or-name"]
	serverID, err := strconv.Atoi(idOrName)
	hostName := ""
	if err == nil {
		name, ok, err := dbhelpers.GetServerNameFromID(inf.Tx.Tx, int64(serverID))
		if err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting server name from id '"+idOrName+"': "+err.Error()))
			return
		} else if !ok {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusNotFound, errors.New("server ID '"+idOrName+"' not found"), nil)
			return
		}
		hostName = name
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
	} else {
		hostName = idOrName
		var ok bool
		serverID, ok, err = dbhelpers.GetServerIDFromName(hostName, inf.Tx.Tx)
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

	values := new(updateValues)

	if hasUpdated {
		updatedBool := strToBool(updated)
		values.configUpdateBool = &updatedBool
	}

	if hasRevalUpdated {
		revalUpdatedBool := strToBool(revalUpdated)
		values.revalUpdateBool = &revalUpdatedBool
	}

	if err := setUpdateStatuses(inf.Tx.Tx, int64(serverID), *values); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("setting updated statuses: "+err.Error()))
		return
	}

	respMsg := "successfully set server '" + hostName + "'"
	if hasUpdated {
		respMsg += " updated=" + updated
	}
	if hasRevalUpdated {
		respMsg += " reval_updated=" + revalUpdated
	}

	api.WriteAlerts(w, r, http.StatusOK, tc.CreateAlerts(tc.SuccessLevel, respMsg))
}

type updateValues struct {
	configUpdateBool *bool // Deprecated, prefer timestamps
	revalUpdateBool  *bool // Deprecated, prefer timestamps
	configApplyTime  *time.Time
	revalApplyTime   *time.Time
}

func parseQueryParams(params map[string]string) (*updateValues, error) {
	var paramValues updateValues

	// Verify query string parameters
	configUpdatedBoolParam, hasConfigUpdatedBoolParam := params["updated"]     // Deprecated, but still required for backwards compatibility
	revalUpdatedBoolParam, hasRevalUpdatedBoolParam := params["reval_updated"] // Deprecated, but still required for backwards compatibility
	configApplyTimeParam, hasConfigApplyTimeParam := params["config_apply_time"]
	revalidateApplyTimeParam, hasRevalidateApplyTimeParam := params["revalidate_apply_time"]

	if !hasConfigApplyTimeParam && !hasRevalidateApplyTimeParam &&
		!hasConfigUpdatedBoolParam && !hasRevalUpdatedBoolParam {
		return nil, errors.New("must pass at least one of the following query parameters: 'config_apply_time', 'revalidate_apply_time' ,'updated', 'reval_updated'")

	}
	// Prevent collision between booleans and timestamps
	if hasConfigApplyTimeParam && hasConfigUpdatedBoolParam {
		return nil, errors.New("conflicting parameters. may not pass 'updated' along with 'config_apply_time'")

	}
	if hasRevalidateApplyTimeParam && hasRevalUpdatedBoolParam {
		return nil, errors.New("conflicting parameters. may not pass 'reval_updated' along with 'revalidate_apply_time'")

	}

	// Validate and parse parameters before attempting to apply them (don't want to partially apply various status before an error)
	// Timestamps
	if hasConfigApplyTimeParam {
		configApplyTime, err := time.Parse(time.RFC3339Nano, configApplyTimeParam)
		if err != nil {
			return nil, errors.New("query parameter 'config_apply_time' must be valid RFC3339Nano format:" + err.Error())
		}
		paramValues.configApplyTime = &configApplyTime
	}

	if hasRevalidateApplyTimeParam {
		revalApplyTime, err := time.Parse(time.RFC3339Nano, revalidateApplyTimeParam)
		if err != nil {
			return nil, errors.New("query parameter 'revalidate_apply_time' must be valid RFC3339Nano format:" + err.Error())
		}
		paramValues.revalApplyTime = &revalApplyTime
	}

	// Booleans
	if hasConfigUpdatedBoolParam {
		updatedBool, err := strconv.ParseBool(configUpdatedBoolParam)
		if err != nil {
			return nil, errors.New("query parameter 'updated' must be a boolean")
		}
		paramValues.configUpdateBool = &updatedBool
	}
	if hasRevalUpdatedBoolParam {
		revalUpdatedBool, err := strconv.ParseBool(revalUpdatedBoolParam)
		if err != nil {
			return nil, errors.New("query parameter 'reval_updated' must be a boolean")
		}
		paramValues.revalUpdateBool = &revalUpdatedBool
	}
	return &paramValues, nil
}

// setUpdateStatuses set timestamps for config update/apply and revalidation
// update/apply. If any value is nil, no changes occur
func setUpdateStatuses(tx *sql.Tx, serverID int64, values updateValues) error {
	if values.configApplyTime != nil {
		if err := dbhelpers.SetApplyUpdateForServerWithTime(tx, serverID, *values.configApplyTime); err != nil {
			return fmt.Errorf("setting config apply time: %w", err)
		}
	}

	if values.revalApplyTime != nil {
		if err := dbhelpers.SetApplyRevalForServerWithTime(tx, serverID, *values.revalApplyTime); err != nil {
			return fmt.Errorf("setting reval apply time: %w", err)
		}
	}

	if values.configUpdateBool != nil {
		if *values.configUpdateBool {
			if err := dbhelpers.QueueUpdateForServer(tx, serverID); err != nil {
				return err
			}
		} else {
			if err := dbhelpers.SetApplyUpdateForServer(tx, serverID); err != nil {
				return err
			}
		}
	}

	if values.revalUpdateBool != nil {
		if *values.revalUpdateBool {
			if err := dbhelpers.QueueRevalForServer(tx, serverID); err != nil {
				return err
			}
		} else {
			if err := dbhelpers.SetApplyRevalForServer(tx, serverID); err != nil {
				return err
			}
		}
	}
	return nil
}

func responseMessage(idOrName string, values updateValues) string {
	respMsg := "successfully set server '" + idOrName + "'"

	if values.configUpdateBool != nil {
		respMsg += " updated=" + strconv.FormatBool(*values.configUpdateBool)
	}
	if values.revalUpdateBool != nil {
		respMsg += " reval_updated=" + strconv.FormatBool(*values.revalUpdateBool)
	}

	if values.configApplyTime != nil {
		respMsg += " config_apply_time=" + (*values.configApplyTime).Format(time.RFC3339Nano)
	}
	if values.revalApplyTime != nil {
		respMsg += " revalidate_apply_time=" + (*values.revalApplyTime).Format(time.RFC3339Nano)
	}

	return respMsg
}

// UpdateHandlerV4 implements an http handler that updates a server's config update and reval times.
func UpdateHandlerV4(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"id-or-name"}, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	idOrName := inf.Params["id-or-name"]
	serverID, err := strconv.ParseInt(idOrName, 10, 64)
	if err != nil {
		id, ok, err := dbhelpers.GetServerIDFromName(idOrName, inf.Tx.Tx)
		if err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("getting server id from name '%v': %w", idOrName, err))
			return
		} else if !ok {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusNotFound, errors.New("server name '"+idOrName+"' not found"), nil)
			return
		}
		serverID = int64(id)
	}

	cdnName, err := dbhelpers.GetCDNNameFromServerID(inf.Tx.Tx, serverID)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, err)
		return
	}
	userErr, sysErr, statusCode := dbhelpers.CheckIfCurrentUserHasCdnLock(inf.Tx.Tx, string(cdnName), inf.User.UserName)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, statusCode, userErr, sysErr)
		return
	}

	// TODO parse JSON body to trump Query Params?
	updateValues, err := parseQueryParams(inf.Params)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, err, nil)
		return
	}

	if updateValues == nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("updateValues should not be nil"))
		return
	}

	err = setUpdateStatuses(inf.Tx.Tx, serverID, *updateValues)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, err)
		return
	}

	respMsg := responseMessage(idOrName, *updateValues)

	api.WriteAlerts(w, r, http.StatusOK, tc.CreateAlerts(tc.SuccessLevel, respMsg))
}
