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
	"time"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
)

// UpdateHandler implements an http handler that updates a server's config update and reval times.
func UpdateHandler(w http.ResponseWriter, r *http.Request) {
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
	updated, hasUpdated := inf.Params["updated"]
	revalUpdated, hasRevalUpdated := inf.Params["reval_updated"]
	config_update_time, hasConfig_update_time := inf.Params["config_update_time"]
	config_apply_time, hasConfig_apply_time := inf.Params["config_apply_time"]
	revalidate_update_time, hasRevalidate_update_time := inf.Params["revalidate_update_time"]
	revalidate_apply_time, hasRevalidate_apply_time := inf.Params["revalidate_apply_time"]
	if !hasUpdated && !hasRevalUpdated && !hasConfig_update_time && !hasConfig_apply_time && !hasRevalidate_update_time && !hasRevalidate_apply_time {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("must pass at least one query parameter: 'updated', 'reval_updated', 'config_update_time', 'config_apply_time', 'revalidate_update_time', 'revalidate_apply_time'"), nil)
		return
	}

	// Legacy, prefer RFC3339 query parameters (*_time) over boolean flags
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

	// Validate parameters before attempting to apply them (don't want to partially apply various status before an error)
	var configUpdateTime time.Time
	if hasConfig_update_time {
		configUpdateTime, err = time.Parse(time.RFC3339Nano, config_update_time)
		if err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("query parameter 'config_update_time' must be valid RFC3339 format"), nil)
			return
		}
	}

	var configApplyTime time.Time
	if hasConfig_apply_time {
		configApplyTime, err = time.Parse(time.RFC3339Nano, config_apply_time)
		if err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("query parameter 'config_apply_time' must be valid RFC3339 format"), nil)
			return
		}
	}

	var revalUpdateTime time.Time
	if hasRevalidate_update_time {
		revalUpdateTime, err = time.Parse(time.RFC3339Nano, revalidate_update_time)
		if err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("query parameter 'revalidate_update_time' must be valid RFC3339 format"), nil)
			return
		}
	}

	var revalApplyTime time.Time
	if hasRevalidate_apply_time {
		revalApplyTime, err = time.Parse(time.RFC3339Nano, revalidate_apply_time)
		if err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("query parameter 'revalidate_apply_time' must be valid RFC3339 format"), nil)
			return
		}
	}

	if hasUpdated && (hasConfig_update_time || hasConfig_apply_time) {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("potential conflict between params. use 'config_update_time' or 'config_apply_time' over 'updated'"), nil)
		return
	}

	if hasRevalUpdated && (hasRevalidate_update_time || hasRevalidate_apply_time) {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("potential conflict between params. use 'reval_update_time' or 'reval_apply_time' over 'reval_updated'"), nil)
		return
	}

	strToBool := func(s string) bool {
		return !strings.HasPrefix(s, "f")
	}

	if hasUpdated {
		updatedBool := strToBool(updated)
		if updatedBool {
			if err = dbhelpers.QueueUpdateForServer(inf.Tx.Tx, serverID); err != nil {
				api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("setting update status: %w", err))
				return
			}
		} else {
			if err = dbhelpers.SetApplyUpdateForServer(inf.Tx.Tx, serverID); err != nil {
				api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("setting update status: %w", err))
				return
			}
		}
	}

	if hasRevalUpdated {
		revalUpdatedBool := strToBool(revalUpdated)
		if revalUpdatedBool {
			if err = dbhelpers.QueueRevalForServer(inf.Tx.Tx, serverID); err != nil {
				api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("setting reval status: %w", err))
				return
			}
		} else {
			if err = dbhelpers.SetApplyRevalForServer(inf.Tx.Tx, serverID); err != nil {
				api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("setting reval status: %w", err))
				return
			}
		}
	}

	if hasConfig_update_time {
		if err = dbhelpers.QueueUpdateForServerWithTime(inf.Tx.Tx, serverID, configUpdateTime); err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("setting config update time: %w", err))
			return
		}
	}

	if hasConfig_apply_time {
		if err = dbhelpers.SetApplyUpdateForServerWithTime(inf.Tx.Tx, serverID, configApplyTime); err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("setting config apply time: %w", err))
			return
		}
	}

	if hasRevalidate_update_time {
		if err = dbhelpers.QueueRevalForServerWithTime(inf.Tx.Tx, serverID, revalUpdateTime); err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("setting reval update time: %w", err))
			return
		}
	}

	if hasConfig_apply_time {
		if err = dbhelpers.SetApplyUpdateForServerWithTime(inf.Tx.Tx, serverID, revalApplyTime); err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("setting reval apply time: %w", err))
			return
		}
	}

	hostName, _, err := dbhelpers.GetServerNameFromID(inf.Tx.Tx, serverID)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("querying for server name with id %d: %w", serverID, err))
		return
	}

	respMsg := "successfully set server '" + hostName + "'"
	if hasUpdated {
		respMsg += " updated=" + updated
	}
	if hasRevalUpdated {
		respMsg += " reval_updated=" + revalUpdated
	}
	if hasConfig_update_time {
		respMsg += " config_update_time=" + config_update_time
	}
	if hasConfig_apply_time {
		respMsg += " config_apply_time=" + config_apply_time
	}
	if hasRevalidate_update_time {
		respMsg += " revalidate_update_time=" + revalidate_update_time
	}
	if hasRevalidate_apply_time {
		respMsg += " revalidate_apply_time=" + revalidate_apply_time
	}

	api.WriteAlerts(w, r, http.StatusOK, tc.CreateAlerts(tc.SuccessLevel, respMsg))
}
