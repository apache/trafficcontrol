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

	if err := setUpdateStatuses(inf.Tx.Tx, int64(serverID), updatedPtr, revalUpdatedPtr); err != nil {
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
func setUpdateStatuses(tx *sql.Tx, serverID int64, updatePending *bool, revalPending *bool) error {
	if updatePending == nil && revalPending == nil {
		return errors.New("either updatePending or revalPending must not be nil")
	}

	if updatePending != nil {
		if *updatePending {
			if err := dbhelpers.QueueUpdateForServer(tx, serverID); err != nil {
				return err
			}
		} else {
			if err := dbhelpers.SetApplyUpdateForServer(tx, serverID); err != nil {
				return err
			}
		}
	}

	if revalPending != nil {
		if *revalPending {
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

// UpdateHandler implements an http handler that updates a server's config update and reval times.
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

	// Verify query string parameters
	configUpdateTimeParam, hasConfigUpdateTimeParam := inf.Params["config_update_time"]
	configApplyTimeParam, hasConfigApplyTimeParam := inf.Params["config_apply_time"]
	revalidateUpdateTimeParam, hasRevalidateUpdateTimeParam := inf.Params["revalidate_update_time"]
	revalidateApplyTimeParam, hasRevalidateApplyTimeParam := inf.Params["revalidate_apply_time"]
	if !hasConfigUpdateTimeParam && !hasConfigApplyTimeParam && !hasRevalidateUpdateTimeParam && !hasRevalidateApplyTimeParam {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("must pass at least one query parameter: 'config_update_time', 'config_apply_time', 'revalidate_update_time', 'revalidate_apply_time'"), nil)
		return
	}

	// Validate parameters before attempting to apply them (don't want to partially apply various status before an error)
	var configUpdateTime time.Time
	if hasConfigUpdateTimeParam {
		configUpdateTime, err = time.Parse(time.RFC3339Nano, configUpdateTimeParam)
		if err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("query parameter 'config_update_time' must be valid RFC3339Nano format"), nil)
			return
		}
	}

	var configApplyTime time.Time
	if hasConfigApplyTimeParam {
		configApplyTime, err = time.Parse(time.RFC3339Nano, configApplyTimeParam)
		if err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("query parameter 'config_apply_time' must be valid RFC3339Nano format"), nil)
			return
		}
	}

	var revalUpdateTime time.Time
	if hasRevalidateUpdateTimeParam {
		revalUpdateTime, err = time.Parse(time.RFC3339Nano, revalidateUpdateTimeParam)
		if err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("query parameter 'revalidate_update_time' must be valid RFC3339Nano format"), nil)
			return
		}
	}

	var revalApplyTime time.Time
	if hasRevalidateApplyTimeParam {
		revalApplyTime, err = time.Parse(time.RFC3339Nano, revalidateApplyTimeParam)
		if err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("query parameter 'revalidate_apply_time' must be valid RFC3339Nano format"), nil)
			return
		}
	}

	if hasConfigUpdateTimeParam {
		if err = dbhelpers.QueueUpdateForServerWithTime(inf.Tx.Tx, serverID, configUpdateTime); err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("setting config update time: %w", err))
			return
		}
	}

	if hasConfigApplyTimeParam {
		if err = dbhelpers.SetApplyUpdateForServerWithTime(inf.Tx.Tx, serverID, configApplyTime); err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("setting config apply time: %w", err))
			return
		}
	}

	if hasRevalidateUpdateTimeParam {
		if err = dbhelpers.QueueRevalForServerWithTime(inf.Tx.Tx, serverID, revalUpdateTime); err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("setting reval update time: %w", err))
			return
		}
	}

	if hasRevalidateApplyTimeParam {
		if err = dbhelpers.SetApplyRevalForServerWithTime(inf.Tx.Tx, serverID, revalApplyTime); err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("setting reval apply time: %w", err))
			return
		}
	}

	respMsg := "successfully set server '" + idOrName + "'"

	if hasConfigUpdateTimeParam {
		respMsg += " config_update_time=" + configUpdateTimeParam
	}
	if hasConfigApplyTimeParam {
		respMsg += " config_apply_time=" + configApplyTimeParam
	}
	if hasRevalidateUpdateTimeParam {
		respMsg += " revalidate_update_time=" + revalidateUpdateTimeParam
	}
	if hasRevalidateApplyTimeParam {
		respMsg += " revalidate_apply_time=" + revalidateApplyTimeParam
	}

	api.WriteAlerts(w, r, http.StatusOK, tc.CreateAlerts(tc.SuccessLevel, respMsg))
}
