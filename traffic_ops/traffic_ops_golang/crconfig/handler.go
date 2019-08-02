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
	"database/sql"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/deliveryservice"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/monitoring"
)

// Handler creates and serves the CRConfig from the raw SQL data.
// This MUST only be used for debugging or previewing, the raw un-snapshotted data MUST NOT be used by any component of the CDN.
func Handler(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"cdn"}, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	start := time.Now()
	crConfig, err := Make(inf.Tx.Tx, inf.Params["cdn"], inf.User.UserName, r.Host, r.URL.Path, inf.Config.Version, inf.Config.CRConfigUseRequestHost, inf.Config.CRConfigEmulateOldPath)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, err)
		return
	}
	log.Infof("CRConfig time to generate: %+v\n", time.Since(start))
	api.WriteResp(w, r, crConfig)
}

// SnapshotGetHandler gets and serves the CRConfig from the snapshot table.
func SnapshotGetHandler(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"cdn"}, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	snapshot, cdnExists, err := GetSnapshot(inf.Tx.Tx, inf.Params["cdn"])
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting snapshot: "+err.Error()))
		return
	}
	if !cdnExists {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusNotFound, errors.New("CDN not found"), nil)
		return
	}
	w.Header().Set(tc.ContentType, tc.ApplicationJson)
	w.Write([]byte(`{"response":` + snapshot + `}`))
}

// SnapshotGetMonitoringHandler gets and serves the CRConfig from the snapshot table.
func SnapshotGetMonitoringHandler(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"cdn"}, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	snapshot, cdnExists, err := GetSnapshotMonitoring(inf.Tx.Tx, inf.Params["cdn"])
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting snapshot: "+err.Error()))
		return
	}
	if !cdnExists {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusNotFound, errors.New("CDN not found"), nil)
		return
	}
	w.Header().Set(tc.ContentType, tc.ApplicationJson)
	w.Write([]byte(`{"response":` + snapshot + `}`))
}

// SnapshotOldGetHandler gets and serves the CRConfig from the snapshot table, not wrapped in response to match the old non-API CRConfig-Snapshots endpoint
func SnapshotOldGetHandler(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"cdn"}, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	snapshot, cdnExists, err := GetSnapshot(inf.Tx.Tx, inf.Params["cdn"])
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting snapshot: "+err.Error()))
		return
	}
	if !cdnExists {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusNotFound, errors.New("CDN not found"), nil)
		return
	}
	w.Header().Set(tc.ContentType, tc.ApplicationJson)
	w.Write([]byte(snapshot))
}

// SnapshotHandler creates the CRConfig JSON and writes it to the snapshot table in the database.
func SnapshotHandler(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, []string{"id"})
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	db, err := api.GetDB(r.Context())
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("SnapshotHandler getting db from context: "+err.Error()))
		return
	}

	cdn, ok := inf.Params["cdn"]
	if !ok {
		id, ok := inf.IntParams["id"]
		if !ok {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusNotFound, errors.New("params missing CDN"), nil)
			return
		}
		name, ok, err := getCDNNameFromID(id, inf.Tx.Tx)
		if err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("Error getting CDN name from ID: "+err.Error()))
			return
		}
		if !ok {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusNotFound, errors.New("No CDN found with that ID"), nil)
			return
		}
		cdn = name
	}

	crConfig, err := Make(inf.Tx.Tx, cdn, inf.User.UserName, r.Host, r.URL.Path, inf.Config.Version, inf.Config.CRConfigUseRequestHost, inf.Config.CRConfigEmulateOldPath)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, err)
		return
	}

	monitoringJSON, err := monitoring.GetMonitoringJSON(inf.Tx.Tx, cdn)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New(r.RemoteAddr+" getting monitoring.json data: "+err.Error()))
		return
	}

	if err := Snapshot(inf.Tx.Tx, crConfig, monitoringJSON); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New(r.RemoteAddr+" snaphsotting CRConfig and Monitoring: "+err.Error()))
		return
	}

	if err := deliveryservice.DeleteOldCerts(db.DB, inf.Tx.Tx, inf.Config, tc.CDNName(cdn)); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New(r.RemoteAddr+" snapshotting CRConfig and Monitoring: starting old certificate deletion job: "+err.Error()))
		return
	}

	api.CreateChangeLogRawTx(api.ApiChange, "CDN: "+cdn+", ID: "+strconv.Itoa(inf.IntParams["id"])+", ACTION: Snapshot of CRConfig and Monitor", inf.User, inf.Tx.Tx)
	api.WriteResp(w, r, "SUCCESS")
}

// SnapshotGUIHandler creates the CRConfig JSON and writes it to the snapshot table in the database. The response emulates the old Perl UI function. This should go away when the old Perl UI ceases to exist.
func SnapshotOldGUIHandler(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, _ := api.NewInfo(r, []string{"cdn"}, nil)
	if userErr != nil || sysErr != nil {
		writePerlHTMLErr(w, r, inf.Tx.Tx, errors.New(r.RemoteAddr+" unable to get info from request: "+sysErr.Error()), userErr)
		return
	}
	defer inf.Close()

	db, err := api.GetDB(r.Context())
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("SnapshotHandler getting db from context: "+err.Error()))
		return
	}

	cdn := inf.Params["cdn"]

	crConfig, err := Make(inf.Tx.Tx, cdn, inf.User.UserName, r.Host, r.URL.Path, inf.Config.Version, inf.Config.CRConfigUseRequestHost, inf.Config.CRConfigEmulateOldPath)
	if err != nil {
		writePerlHTMLErr(w, r, inf.Tx.Tx, errors.New(r.RemoteAddr+" making CRConfig: "+err.Error()), err)
		return
	}

	tm, err := monitoring.GetMonitoringJSON(inf.Tx.Tx, cdn)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New(r.RemoteAddr+" getting monitoring.json data: "+err.Error()))
		return
	}

	if err := Snapshot(inf.Tx.Tx, crConfig, tm); err != nil {
		writePerlHTMLErr(w, r, inf.Tx.Tx, errors.New(r.RemoteAddr+" making CRConfig: "+err.Error()), err)
		return
	}

	if err := deliveryservice.DeleteOldCerts(db.DB, inf.Tx.Tx, inf.Config, tc.CDNName(cdn)); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New(r.RemoteAddr+" old snapshotting CRConfig and Monitoring: starting old certificate deletion job: "+err.Error()))
		return
	}

	api.CreateChangeLogRawTx(api.ApiChange, "Snapshot of CRConfig performed for "+cdn, inf.User, inf.Tx.Tx)
	http.Redirect(w, r, "/tools/flash_and_close/"+url.PathEscape("Successfully wrote the CRConfig.json!"), http.StatusFound)
}

func writePerlHTMLErr(w http.ResponseWriter, r *http.Request, tx *sql.Tx, logErr error, err error) {
	log.Errorln(logErr.Error())
	tx.Rollback()
	http.Redirect(w, r, "/tools/flash_and_close/"+url.PathEscape("Error: "+err.Error()), http.StatusFound)
}
