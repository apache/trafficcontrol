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
	"net/http"
	"strconv"
	"strings"

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
	serverID, err := strconv.Atoi(idOrName)
	if err != nil {
		id, ok, err := dbhelpers.GetServerIDFromName(idOrName, inf.Tx.Tx)
		if err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting server id from name '"+idOrName+"': "+err.Error()))
			return
		} else if !ok {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusNotFound, errors.New("server name '"+idOrName+"' not found"), nil)
			return
		}
		serverID = id
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

	updated, hasUpdated := inf.Params["updated"]
	revalUpdated, hasRevalUpdated := inf.Params["reval_updated"]
	if !hasUpdated && !hasRevalUpdated {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("must pass at least one query parameter of 'updated' or 'reval_updated'"), nil)
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
		return !strings.HasPrefix(s, "f")
	}

	if hasUpdated {
		updatedBool := strToBool(updated)
		if updatedBool {
			if err = dbhelpers.QueueUpdateForServer(inf.Tx.Tx, int64(serverID)); err != nil {
				api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("setting update status: "+err.Error()))
				return
			}
		}
	}

	if hasRevalUpdated {
		revalUpdatedBool := strToBool(revalUpdated)
		if revalUpdatedBool {
			if err = dbhelpers.QueueRevalForServer(inf.Tx.Tx, int64(serverID)); err != nil {
				api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("setting reval status: "+err.Error()))
				return
			}
		}
	}

	hostName, _, _ := dbhelpers.GetServerNameFromID(inf.Tx.Tx, serverID)

	respMsg := "successfully set server '" + hostName + "'"
	if hasUpdated {
		respMsg += " updated=" + updated
	}
	if hasRevalUpdated {
		respMsg += " reval_updated=" + revalUpdated
	}

	api.WriteAlerts(w, r, http.StatusOK, tc.CreateAlerts(tc.SuccessLevel, respMsg))
}
