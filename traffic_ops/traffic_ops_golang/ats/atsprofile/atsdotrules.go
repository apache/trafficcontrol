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
	"strings"

	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/ats"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/config"
)

func GetATSDotRules(w http.ResponseWriter, r *http.Request) {
	WithProfileData(w, r, makeATSDotRules)
}

func makeATSDotRules(tx *sql.Tx, _ *config.Config, profile ats.ProfileData, fileName string) (string, error) {
	// TODO add more efficient db func to only get drive params?
	paramData, err := ats.GetProfileParamData(tx, profile.ID, "storage.config") // ats.rules is based on the storage.config params
	if err != nil {
		return "", errors.New("getting profile param data: " + err.Error())
	}

	drivePrefix := strings.TrimPrefix(paramData["Drive_Prefix"], `/dev/`)
	drivePostfix := strings.Split(paramData["Drive_Letters"], ",")

	text := ""
	for _, l := range drivePostfix {
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}
		text += `KERNEL=="` + drivePrefix + l + `", OWNER="ats"` + "\n"
	}
	if ramPrefix, ok := paramData["RAM_Drive_Prefix"]; ok {
		ramPrefix = strings.TrimPrefix(ramPrefix, `/dev/`)
		ramPostfix := strings.Split(paramData["RAM_Drive_Letters"], ",")
		for _, l := range ramPostfix {
			text += `KERNEL=="` + ramPrefix + l + `", OWNER="ats"` + "\n"
		}
	}
	if text == "" {
		text = "\n" // prevents it being flagged as "not found"
	}
	return text, nil
}
