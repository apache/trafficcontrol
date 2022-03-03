package cachegroup

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
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
)

const queue = "queue"
const dequeue = "dequeue"

func QueueUpdates(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"id"}, []string{"id"})
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	// Validate body and params
	reqObj := tc.CachegroupQueueUpdatesRequest{}
	if err := json.NewDecoder(r.Body).Decode(&reqObj); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("malformed JSON: "+err.Error()), nil)
		return
	}
	if reqObj.Action != queue && reqObj.Action != dequeue {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("action must be 'queue' or 'dequeue'"), nil)
		return
	}
	if reqObj.CDNID == nil && (reqObj.CDN == nil || *reqObj.CDN == "") {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("cdn or cdnId is required"), nil)
		return
	}

	if reqObj.CDNID != nil && (reqObj.CDN == nil || *reqObj.CDN == "") {
		cdnName, ok, sysErr := dbhelpers.GetCDNNameFromID(inf.Tx.Tx, int64(*reqObj.CDNID))
		if userErr != nil || sysErr != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, sysErr)
			return
		}
		if !ok {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, fmt.Errorf("cdn %d does not exist", *reqObj.CDNID), nil)
			return
		}
		reqObj.CDN = &cdnName
	}

	if reqObj.CDNID == nil && (reqObj.CDN != nil && *reqObj.CDN != "") {
		cdnID, ok, sysErr := dbhelpers.GetCDNIDFromName(inf.Tx.Tx, *reqObj.CDN)
		if userErr != nil || sysErr != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, sysErr)
			return
		}
		if !ok {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, fmt.Errorf("cdn %s does not exist", *reqObj.CDN), nil)
			return
		}
		reqObj.CDNID = &cdnID
	}

	cgID := inf.IntParams["id"]
	cgName, ok, err := dbhelpers.GetCacheGroupNameFromID(inf.Tx.Tx, cgID)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting cachegroup name from ID '"+inf.Params["id"]+"': "+err.Error()))
		return
	} else if !ok {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusNotFound, nil, nil)
		return
	}

	// Verify rights
	userErr, sysErr, statusCode := dbhelpers.CheckIfCurrentUserHasCdnLock(inf.Tx.Tx, string(*reqObj.CDN), inf.User.UserName)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, statusCode, userErr, sysErr)
		return
	}

	// Queue updates
	var updatedCaches []tc.CacheName
	if reqObj.Action == queue {
		updatedCaches, err = queueUpdates(inf.Tx.Tx, cgID, *reqObj.CDNID)
	} else {
		updatedCaches, err = dequeueUpdates(inf.Tx.Tx, cgID, *reqObj.CDNID)
	}
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("queueing updates: "+err.Error()))
		return
	}

	api.WriteResp(w, r, QueueUpdatesResp{
		CacheGroupName: cgName,
		Action:         reqObj.Action,
		ServerNames:    updatedCaches,
		CDN:            *reqObj.CDN,
		CacheGroupID:   cgID,
	})
	api.CreateChangeLogRawTx(api.ApiChange, "CACHEGROUP: "+string(cgName)+", ID: "+strconv.Itoa(cgID)+", ACTION: "+strings.Title(reqObj.Action)+"d CacheGroup server updates to the "+string(*reqObj.CDN)+" CDN", inf.User, inf.Tx.Tx)
}

type QueueUpdatesResp struct {
	CacheGroupName tc.CacheGroupName `json:"cachegroupName"`
	Action         string            `json:"action"`
	ServerNames    []tc.CacheName    `json:"serverNames"`
	CDN            tc.CDNName        `json:"cdn"`
	CacheGroupID   int               `json:"cachegroupID"`
}

func queueUpdates(tx *sql.Tx, cgID int, cdnID int) ([]tc.CacheName, error) {
	q := `
INSERT INTO public.server_config_update (server_id, config_update_time) 
SELECT s.id, now() FROM "server" s WHERE s.cachegroup = $1 AND s.cdn_id = $2
ON CONFLICT (server_id)
DO UPDATE SET config_update_time = now()
RETURNING (SELECT s.host_name FROM "server" s WHERE s.id = server_id);
	`
	rows, err := tx.Query(q, cgID, cdnID)
	if err != nil {
		return nil, errors.New("querying : " + err.Error())
	}
	defer rows.Close()
	names := []tc.CacheName{}
	for rows.Next() {
		name := ""
		if err := rows.Scan(&name); err != nil {
			return nil, errors.New("scanning queue updates: " + err.Error())
		}
		names = append(names, tc.CacheName(name))
	}
	return names, nil
}

func dequeueUpdates(tx *sql.Tx, cgID int, cdnID int) ([]tc.CacheName, error) {
	q := `
UPDATE public.server_config_update
SET config_apply_time = config_update_time
WHERE server_id IN (SELECT s.id FROM "server" s WHERE s.cachegroup = $1 AND s.cdn_id = $2)
RETURNING (SELECT s.host_name FROM "server" s WHERE s.id = server_id);
`
	rows, err := tx.Query(q, cgID, cdnID)
	if err != nil {
		return nil, errors.New("querying queue updates: " + err.Error())
	}
	defer rows.Close()
	names := []tc.CacheName{}
	for rows.Next() {
		name := ""
		if err := rows.Scan(&name); err != nil {
			return nil, errors.New("scanning queue updates: " + err.Error())
		}
		names = append(names, tc.CacheName(name))
	}
	return names, nil
}
