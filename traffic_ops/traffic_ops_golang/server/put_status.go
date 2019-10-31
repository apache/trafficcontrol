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
	"strings"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
)

func UpdateStatusHandler(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"id"}, []string{"id"})
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()
	reqObj := tc.ServerPutStatus{}
	if err := json.NewDecoder(r.Body).Decode(&reqObj); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("malformed JSON: "+err.Error()), nil)
		return
	}

	serverInfo, exists, err := dbhelpers.GetServerInfo(inf.IntParams["id"], inf.Tx.Tx)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, err)
		return
	}
	if !exists {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusNotFound, fmt.Errorf("server ID %d not found", inf.IntParams["id"]), nil)
		return
	}

	status := tc.StatusNullable{}
	statusExists := false
	if reqObj.Status.Name != nil {
		status, statusExists, err = dbhelpers.GetStatusByName(*reqObj.Status.Name, inf.Tx.Tx)
	} else if reqObj.Status.ID != nil {
		status, statusExists, err = dbhelpers.GetStatusByID(*reqObj.Status.ID, inf.Tx.Tx)
	} else {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("status is required"), nil)
		return
	}
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, err)
		return
	}
	if !statusExists {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("invalid status (does not exist)"), nil)
		return
	}

	if *status.Name == tc.CacheStatusAdminDown.String() || *status.Name == tc.CacheStatusOffline.String() {
		if reqObj.OfflineReason == nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("offlineReason is required for "+tc.CacheStatusAdminDown.String()+" or "+tc.CacheStatusOffline.String()+" status"), nil)
			return
		}
		*reqObj.OfflineReason = inf.User.UserName + ": " + *reqObj.OfflineReason
	} else {
		reqObj.OfflineReason = nil
	}
	if err := updateServerStatusAndOfflineReason(inf.IntParams["id"], *status.ID, reqObj.OfflineReason, inf.Tx.Tx); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, err)
		return
	}
	offlineReason := ""
	if reqObj.OfflineReason != nil {
		offlineReason = *reqObj.OfflineReason
	}
	msg := "Updated status [ " + *status.Name + " ] for " + serverInfo.HostName + "." + serverInfo.DomainName + " [ " + offlineReason + " ]"

	// queue updates on child servers if server is ^EDGE or ^MID
	if strings.HasPrefix(serverInfo.Type, tc.CacheTypeEdge.String()) || strings.HasPrefix(serverInfo.Type, tc.CacheTypeMid.String()) {
		if err := queueUpdatesOnChildCaches(inf.Tx.Tx, serverInfo.CDNID, serverInfo.CachegroupID); err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, err)
			return
		}
		msg += " and queued updates on all child caches"
	}
	api.CreateChangeLogRawTx(api.ApiChange, msg, inf.User, inf.Tx.Tx)
	api.WriteRespAlert(w, r, tc.SuccessLevel, msg)
}

// queueUpdatesOnChildCaches queues updates on child caches of the given cdnID and parentCachegroupID and returns an error (if one occurs).
func queueUpdatesOnChildCaches(tx *sql.Tx, cdnID, parentCachegroupID int) error {
	q := `
UPDATE server
SET    upd_pending = TRUE
WHERE  server.cdn_id = $1
       AND server.cachegroup IN (SELECT id
                                 FROM   cachegroup
                                 WHERE  parent_cachegroup_id = $2
                                        OR secondary_parent_cachegroup_id = $2)
`
	if _, err := tx.Exec(q, cdnID, parentCachegroupID); err != nil {
		return errors.New("queueing updates on child caches: " + err.Error())
	}
	return nil
}

// updateServerStatusAndOfflineReason updates a server's status and offline_reason and returns an error (if one occurs).
func updateServerStatusAndOfflineReason(serverID, statusID int, offlineReason *string, tx *sql.Tx) error {
	q := `
UPDATE server
SET    status = $1,
       offline_reason = $2
WHERE  id = $3
`
	if _, err := tx.Exec(q, statusID, offlineReason, serverID); err != nil {
		return errors.New("updating server status and offline_reason: " + err.Error())
	}
	return nil
}
