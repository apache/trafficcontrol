package atsprofile

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

	"github.com/apache/trafficcontrol/lib/go-atscfg"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/ats"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/config"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/riaksvc"
)

func GetURLSig(w http.ResponseWriter, r *http.Request) {
	WithProfileData(w, r, tc.ContentTypeTextPlain, urlSigDotConfig)
}

func urlSigDotConfig(tx *sql.Tx, cfg *config.Config, profile ats.ProfileData, fileName string) (string, error) {
	fileName = "url_sig_" + fileName + ".config" // the fileName from the http router is just the DS, missing "url_sig_" and ".config" - add them back now

	urlSigKeys, _, err := riaksvc.GetURLSigKeysFromConfigFileKey(tx, cfg.RiakAuthOptions, cfg.RiakPort, fileName)
	if err != nil {
		return "", errors.New("getting url sig keys from Riak: " + err.Error())
	}

	paramData, err := ats.GetProfileParamData(tx, profile.ID, fileName)
	if err != nil {
		return "", errors.New("getting profile param data: " + err.Error())
	}

	toolName, toURL, err := ats.GetToolNameAndURL(tx)
	if err != nil {
		return "", errors.New("getting tool name and URL: " + err.Error())
	}

	return atscfg.MakeURLSigConfig(profile.Name, urlSigKeys, paramData, toolName, toURL), nil
}
