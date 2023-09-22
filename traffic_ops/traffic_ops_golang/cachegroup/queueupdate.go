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
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/dbhelpers"
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
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("cdn is a required field"), nil)
		return
	}

	var cdnID int64
	if reqObj.CDNID != nil && (reqObj.CDN == nil || *reqObj.CDN == "") {
		cdnName, ok, sysErr := dbhelpers.GetCDNNameFromID(inf.Tx.Tx, int64(*reqObj.CDNID))
		if sysErr != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, sysErr)
			return
		}
		if !ok {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, fmt.Errorf("cdn id %d does not exist", *reqObj.CDNID), nil)
			return
		}
		cdnID = int64(*reqObj.CDNID)
		reqObj.CDN = &cdnName
	}

	if reqObj.CDNID == nil && (reqObj.CDN != nil && *reqObj.CDN != "") {
		var found bool
		var cdnIDint int
		cdnIDint, found, sysErr = dbhelpers.GetCDNIDFromName(inf.Tx.Tx, *reqObj.CDN)
		cdnID = int64(cdnIDint)
		if sysErr != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, sysErr)
			return
		}
		if !found {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, fmt.Errorf("cdn %s does not exist", *reqObj.CDN), nil)
			return
		}
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

	queue := reqObj.Action == "queue"
	if queue {
		userErr, sysErr, statusCode := dbhelpers.CheckIfCurrentUserHasCdnLock(inf.Tx.Tx, string(*reqObj.CDN), inf.User.UserName)
		if userErr != nil || sysErr != nil {
			api.HandleErr(w, r, inf.Tx.Tx, statusCode, userErr, sysErr)
			return
		}
	}

	// Queue updates
	var updatedCaches []tc.CacheName
	if reqObj.Action == "queue" {
		updatedCaches, err = dbhelpers.QueueUpdateForServerWithCachegroupCDN(inf.Tx.Tx, cgID, cdnID)
	} else {
		updatedCaches, err = dbhelpers.DequeueUpdateForServerWithCachegroupCDN(inf.Tx.Tx, cgID, cdnID)
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
