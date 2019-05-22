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

func GetUnknown(w http.ResponseWriter, r *http.Request) {
	addHdr := false
	WithProfileDataHdr(w, r, addHdr, "text/plain", makeUnknown)
}

func makeUnknown(tx *sql.Tx, _ *config.Config, profile ats.ProfileData, fileName string) (string, error) {
	params, err := ats.GetProfileParamData(tx, profile.ID, fileName)
	if err != nil {
		return "", errors.New("getting profile param data: " + err.Error())
	}
	fileContents, err := takeAndBakeProfile(tx, profile.Name, params)
	if err != nil {
		return "", errors.New("GetProfileConfig: takeAndBakeProfile '" + fileName + "': " + err.Error())
	}
	return fileContents, nil
}

func takeAndBakeProfile(tx *sql.Tx, profileName string, params map[string]string) (string, error) {
	hdr, err := ats.HeaderComment(tx, profileName)
	if err != nil {
		return "", errors.New("getting header comment: " + err.Error())
	}
	text := ""
	for paramName, paramVal := range params {
		if paramName == "header" {
			if paramVal == "none" {
				hdr = ""
			} else {
				hdr = paramVal + "\n"
			}
		} else {
			text += paramVal + "\n"
		}
	}
	text = strings.Replace(text, "__RETURN__", "\n", -1)
	return hdr + text, nil
}
