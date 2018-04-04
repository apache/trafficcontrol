package crconfig

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
	"net/url"
	"strconv"
	"time"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/config"

	"github.com/jmoiron/sqlx"
)

const PrivLevel = auth.PrivLevelAdmin

// Handler creates and serves the CRConfig from the raw SQL data.
// This MUST only be used for debugging or previewing, the raw un-snapshotted data MUST NOT be used by any component of the CDN.
func Handler(db *sqlx.DB, cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		handleErrs := tc.GetHandleErrorsFunc(w, r)

		params, err := api.GetCombinedParams(r)
		if err != nil {
			handleErrs(http.StatusInternalServerError, err)
			return
		}

		cdn, ok := params["cdn"]
		if !ok {
			handleErrs(http.StatusInternalServerError, errors.New("params missing CDN"))
			return
		}

		ctx := r.Context()
		user, err := auth.GetCurrentUser(ctx)
		if err != nil {
			handleErrs(http.StatusInternalServerError, err)
			return
		}

		crConfig, err := Make(db.DB, cdn, user.UserName, r.Host, r.URL.Path, cfg.Version)
		if err != nil {
			handleErrs(http.StatusInternalServerError, err)
			return
		}

		respBts, err := json.Marshal(crConfig)
		if err != nil {
			handleErrs(http.StatusInternalServerError, err)
			return
		}
		log.Infof("CRConfig time to generate: %+v\n", time.Since(start))
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s", respBts)
	}
}

// SnapshotGetHandler gets and serves the CRConfig from the snapshot table.
func SnapshotGetHandler(db *sqlx.DB, cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleErrs := tc.GetHandleErrorsFunc(w, r)
		params, err := api.GetCombinedParams(r)
		if err != nil {
			handleErrs(http.StatusInternalServerError, err)
			return
		}
		cdn, ok := params["cdn"]
		if !ok {
			handleErrs(http.StatusInternalServerError, errors.New("params missing CDN"))
			return
		}

		snapshot, cdnExists, err := GetSnapshot(db.DB, cdn)
		if err != nil {
			handleErrs(http.StatusInternalServerError, errors.New("getting snapshot: "+err.Error()))
			return
		}
		if !cdnExists {
			handleErrs(http.StatusNotFound, errors.New("CDN not found"))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s", []byte(snapshot))
	}
}

// SnapshotHandler creates the CRConfig JSON and writes it to the snapshot table in the database.
func SnapshotHandler(db *sqlx.DB, cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleErrs := tc.GetHandleErrorsFunc(w, r)
		params, err := api.GetCombinedParams(r)
		if err != nil {
			handleErrs(http.StatusInternalServerError, err)
			return
		}

		ctx := r.Context()
		user, err := auth.GetCurrentUser(ctx)
		if err != nil {
			handleErrs(http.StatusInternalServerError, err)
			return
		}

		cdn, ok := params["cdn"]
		if !ok {
			idStr, ok := params["id"]
			if !ok {
				handleErrs(http.StatusNotFound, errors.New("params missing CDN"))
				return
			}
			id, err := strconv.Atoi(idStr)
			if err != nil {
				handleErrs(http.StatusNotFound, errors.New("param CDN ID is not an integer"))
				return
			}
			name, ok, err := getCDNNameFromID(id, db.DB)
			if err != nil {
				handleErrs(http.StatusInternalServerError, errors.New("Error getting CDN name from ID: "+err.Error()))
				return
			}
			if !ok {
				handleErrs(http.StatusNotFound, errors.New("No CDN found with that ID"))
				return
			}
			cdn = name
		}

		crConfig, err := Make(db.DB, cdn, user.UserName, r.Host, r.URL.Path, cfg.Version)
		if err != nil {
			handleErrs(http.StatusInternalServerError, err)
			return
		}

		if err := Snapshot(db.DB, crConfig); err != nil {
			handleErrs(http.StatusInternalServerError, err)
			return
		}
		w.WriteHeader(http.StatusOK) // TODO change to 204 No Content in new version
	}
}

// SnapshotGUIHandler creates the CRConfig JSON and writes it to the snapshot table in the database. The response emulates the old Perl UI function. This should go away when the old Perl UI ceases to exist.
func SnapshotOldGUIHandler(db *sqlx.DB, cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params, err := api.GetCombinedParams(r)
		if err != nil {
			log.Errorln(r.RemoteAddr + " unable to get parameters from request: " + err.Error())
			writePerlHTMLErr(w, r, err)
			return
		}

		cdn, ok := params["cdn"]
		if !ok {
			err := errors.New("params missing CDN")
			log.Errorln(r.RemoteAddr + " " + err.Error())
			writePerlHTMLErr(w, r, err)
			return
		}

		log.Errorln("DEBUG calling crconfig.SnapshotOldGUIHandler CDN " + cdn)

		ctx := r.Context()
		user, err := auth.GetCurrentUser(ctx)
		if err != nil {
			log.Errorln(r.RemoteAddr + " getting user: " + err.Error())
			writePerlHTMLErr(w, r, err)
			return
		}

		crConfig, err := Make(db.DB, cdn, user.UserName, r.Host, r.URL.Path, cfg.Version)
		if err != nil {
			log.Errorln(r.RemoteAddr + " making CRConfig: " + err.Error())
			writePerlHTMLErr(w, r, err)
			return
		}

		if err := Snapshot(db.DB, crConfig); err != nil {
			log.Errorln(r.RemoteAddr + " making CRConfig: " + err.Error())
			writePerlHTMLErr(w, r, err)
			return
		}

		http.Redirect(w, r, "/tools/flash_and_close/"+url.PathEscape("Successfully wrote the CRConfig.json!"), http.StatusFound)
	}
}

func writePerlHTMLErr(w http.ResponseWriter, r *http.Request, err error) {
	http.Redirect(w, r, "/tools/flash_and_close/"+url.PathEscape("Error: "+err.Error()), http.StatusFound)
}
