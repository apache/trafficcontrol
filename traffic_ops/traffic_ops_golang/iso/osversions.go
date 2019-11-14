package iso

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
	"io/ioutil"
	"net/http"
	"path"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/jmoiron/sqlx"
)

func GetOSVersions(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	f := api.RespWriter(w, r, inf.Tx.Tx)

	cfgPath, err := osversionCfgPath(inf.Tx)
	if err != nil {
		f(nil, err)
		return
	}

	b, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		f(nil, err)
		return
	}

	var data tc.OSVersionsResponse
	if err = json.Unmarshal(b, &data); err != nil {
		f(nil, err)
		return
	}

	f(data, nil)
}

const (
	ksFilesParamName       = "kickstart.files.location"
	ksFilesParamConfigFile = "mkisofs"
	//defaultCfgDir          = "/var/www/files"
	defaultCfgDir = "."
	cfgFilename   = "osversions.json"
)

func osversionCfgPath(tx *sqlx.Tx) (string, error) {
	var cfgDir string
	if err := tx.QueryRow(`SELECT value FROM parameter WHERE name = $1 AND config_file = $2 LIMIT 1`, ksFilesParamName, ksFilesParamConfigFile).Scan(&cfgDir); err != nil {
		if err != sql.ErrNoRows {
			return "", err
		}
	}
	if cfgDir == "" {
		cfgDir = defaultCfgDir
	}

	return path.Join(cfgDir, cfgFilename), nil
}
