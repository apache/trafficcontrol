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

	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/ats"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/config"
)

func GetDropQString(w http.ResponseWriter, r *http.Request) {
	WithProfileData(w, r, makeDropQString)
}

func makeDropQString(tx *sql.Tx, _ *config.Config, profile ats.ProfileData, _ string) (string, error) {
	dropQStringVal, hasDropQStringParam, err := ats.GetProfileParamValue(tx, profile.ID, "drop_qstring.config", "content")
	if err != nil {
		return "", errors.New("getting profile param val: " + err.Error())
	}

	text := ""
	if hasDropQStringParam {
		text += dropQStringVal + "\n"
	} else {
		text += `/([^?]+) $s://$t/$1` + "\n"
	}
	return text, nil
}
