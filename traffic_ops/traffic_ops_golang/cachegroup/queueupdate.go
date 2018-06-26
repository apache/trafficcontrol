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
	"net/http"
	"strconv"

	"github.com/apache/trafficcontrol/lib/go-tc"
	tcv13 "github.com/apache/trafficcontrol/lib/go-tc/v13"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/auth"
)

func QueueUpdates(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.GetCurrentUser(r.Context())
		if err != nil {
			api.HandleErr(w, r, http.StatusInternalServerError, nil, errors.New("getting user: "+err.Error()))
			return
		}
		params, intParams, userErr, sysErr, errCode := api.AllParams(r, []string{"id"}, []string{"id"})
		if userErr != nil || sysErr != nil {
			api.HandleErr(w, r, errCode, userErr, sysErr)
			return
		}
		reqObj := tcv13.CachegroupQueueUpdatesRequest{}
		if err := json.NewDecoder(r.Body).Decode(&reqObj); err != nil {
			api.HandleErr(w, r, http.StatusBadRequest, errors.New("malformed JSON: "+err.Error()), nil)
			return
		}
		if reqObj.Action != "queue" && reqObj.Action != "dequeue" {
			api.HandleErr(w, r, http.StatusBadRequest, errors.New("action must be 'queue' or 'dequeue'"), nil)
			return
		}
		if reqObj.CDN == nil && reqObj.CDNID == nil {
			api.HandleErr(w, r, http.StatusBadRequest, errors.New("cdn does not exist"), nil)
			return
		}
		if reqObj.CDN == nil || *reqObj.CDN == "" {
			cdn, ok, err := getCDNNameFromID(db, int64(*reqObj.CDNID))
			if err != nil {
				api.HandleErr(w, r, http.StatusInternalServerError, nil, errors.New("getting CDN name from ID '"+strconv.Itoa(int(*reqObj.CDNID))+"': "+err.Error()))
				return
			}
			if !ok {
				api.HandleErr(w, r, http.StatusBadRequest, errors.New("cdn "+strconv.Itoa(int(*reqObj.CDNID))+" does not exist"), nil)
				return
			}
			reqObj.CDN = &cdn
		}
		cgID := int64(intParams["id"])
		cgName, ok, err := getCGNameFromID(db, cgID)
		if err != nil {
			api.HandleErr(w, r, http.StatusInternalServerError, nil, errors.New("getting cachegroup name from ID '"+params["id"]+"': "+err.Error()))
			return
		}
		if !ok {
			api.HandleErr(w, r, http.StatusBadRequest, errors.New("cachegroup "+params["id"]+" does not exist"), nil)
			return
		}
		queue := reqObj.Action == "queue"
		updatedCaches, err := queueUpdates(db, cgID, *reqObj.CDN, queue)
		if err != nil {
			api.HandleErr(w, r, http.StatusInternalServerError, nil, errors.New("queueing updates: "+err.Error()))
			return
		}
		api.WriteResp(w, r, QueueUpdatesResp{
			CacheGroupName: cgName,
			Action:         reqObj.Action,
			ServerNames:    updatedCaches,
			CDN:            *reqObj.CDN,
			CacheGroupID:   cgID,
		})
		api.CreateChangeLogRaw(api.ApiChange, "Server updates "+reqObj.Action+"d for "+string(cgName), user, db)
	}
}

type QueueUpdatesResp struct {
	CacheGroupName tc.CacheGroupName `json:"cachegroupName"`
	Action         string            `json:"action"`
	ServerNames    []tc.CacheName    `json:"serverNames"`
	CDN            tc.CDNName        `json:"cdn"`
	CacheGroupID   int64             `json:"cachegroupID"`
}

func getCDNNameFromID(db *sql.DB, id int64) (tc.CDNName, bool, error) {
	name := ""
	if err := db.QueryRow(`SELECT name FROM cdn WHERE id = $1`, id).Scan(&name); err != nil {
		if err == sql.ErrNoRows {
			return "", false, nil
		}
		return "", false, errors.New("querying CDN ID: " + err.Error())
	}
	return tc.CDNName(name), true, nil
}

func getCGNameFromID(db *sql.DB, id int64) (tc.CacheGroupName, bool, error) {
	name := ""
	if err := db.QueryRow(`SELECT name FROM cachegroup WHERE id = $1`, id).Scan(&name); err != nil {
		if err == sql.ErrNoRows {
			return "", false, nil
		}
		return "", false, errors.New("querying cachegroup ID: " + err.Error())
	}
	return tc.CacheGroupName(name), true, nil
}

func queueUpdates(db *sql.DB, cgID int64, cdn tc.CDNName, queue bool) ([]tc.CacheName, error) {
	q := `
UPDATE server SET upd_pending = $1
WHERE server.cachegroup = $2
AND server.cdn_id = (select id from cdn where name = $3)
RETURNING server.host_name
`
	rows, err := db.Query(q, queue, cgID, cdn)
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
