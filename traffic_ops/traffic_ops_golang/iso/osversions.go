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

// GetOSVersions returns a map of available Operating System (OS) versions for ISO generation,
// as well as the name of the directory where the "kickstarter" files are found.
// The returned data comes from a configuration file. There's a default location of the config file
// which can be overridden via a Parameter database entry.
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

	cfgDefaultDir = "/var/www/files"
	cfgFilename   = "osversions.json"
)

// osversionsCfgPath returns a path to the configuration file
// containing the OS versions data. The name of the configuration
// file is constant, but a specific Parameter database entry can update the
// directory where the file is expected.
func osversionCfgPath(tx *sqlx.Tx) (string, error) {
	var cfgDir string
	err := tx.QueryRow(
		`SELECT value FROM parameter WHERE name = $1 AND config_file = $2 LIMIT 1`,
		ksFilesParamName,
		ksFilesParamConfigFile,
	).Scan(&cfgDir)
	if err != nil {
		if err != sql.ErrNoRows {
			return "", err
		}
	}
	if cfgDir == "" {
		cfgDir = cfgDefaultDir
	}

	return path.Join(cfgDir, cfgFilename), nil
}
