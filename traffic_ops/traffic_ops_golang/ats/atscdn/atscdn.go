package atscdn

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
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/ats"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/config"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
)

// GenericProfileConfig generates a generic profile config text, from the profile's parameters with the given config file name.
// This does not include a header comment, because a generic config may not use a number sign as a comment.
// If you need a header comment, it can be added manually via ats.HeaderComment, or automatically with WithProfileDataHdr.
func GenericProfileConfig(tx *sql.Tx, profile ats.ProfileData, fileName string, separator string) (string, error) {
	profileParamData, err := ats.GetProfileParamData(tx, profile.ID, fileName)
	if err != nil {
		return "", errors.New("getting profile param data: " + err.Error())
	}
	text := ""
	for name, val := range profileParamData {
		name = trimParamUnderscoreNumSuffix(name)
		text += name + separator + val + "\n"
	}
	return text, nil
}

// trimParamUnderscoreNumSuffix removes any trailing "__[0-9]+" and returns the trimmed string.
func trimParamUnderscoreNumSuffix(paramName string) string {
	underscorePos := strings.LastIndex(paramName, `__`)
	if underscorePos == -1 {
		return paramName
	}
	if _, err := strconv.ParseFloat(paramName[underscorePos+2:], 64); err != nil {
		return paramName
	}
	return paramName[:underscorePos]
}

type MakeCfgFunc func(tx *sql.Tx, cfg *config.Config, profile ats.ProfileData, fileName string)

// WithProfileData takes a makeCfg function which takes the ProfileData and returns the config text or any error.
//
// Most profile config files need the same data and write the same text file, so this can be used to reduce duplicate boilerplate code.
//
// This also adds HeaderComment with the profile name to the top of the config text.
//
// The route must include an "id" parameter.
//
// The route may include a "file" parameter, and if so, it will be passed to makeCfg as fileName. If not, fileName will be the empty string.
//
// If makeCfg returns a nil error and the empty string, a 404 Not Found will be returned to the client.
//
// If you need to avoid adding the standard header comment, or use a Content-Type other than text/plain, use WithProfileDataHdr.
//
func WithProfileData(w http.ResponseWriter, r *http.Request, makeCfg func(tx *sql.Tx, cfg *config.Config, profile ats.ProfileData, fileName string) (string, error)) {
	addHdr := true
	WithProfileDataHdr(w, r, addHdr, tc.ContentTypeTextPlain, makeCfg)
}

func WithProfileDataHdr(w http.ResponseWriter, r *http.Request, addHdr bool, contentType string, makeCfg func(tx *sql.Tx, cfg *config.Config, profile ats.ProfileData, fileName string) (string, error)) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"profile-name-or-id"}, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	profileNameOrID := inf.Params["profile-name-or-id"]
	profileID, err := strconv.Atoi(profileNameOrID)
	if err != nil {
		profileName := profileNameOrID
		ok := false
		if profileID, ok, err = ats.GetProfileIDFromName(inf.Tx.Tx, profileName); err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting profile id from name: "+err.Error()))
			return
		} else if !ok {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusNotFound, errors.New("Resource not found."), nil)
			return
		}
	}

	fileName := strings.TrimSuffix(inf.Params["file"], ".json")

	profileData, ok, err := ats.GetProfileData(inf.Tx.Tx, profileID)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting profile data: "+err.Error()))
		return
	}
	if !ok {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusNotFound, errors.New("not found"), nil)
		return
	}

	hdr := ""
	if addHdr {
		if hdr, err = ats.HeaderComment(inf.Tx.Tx, profileData.Name); err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting file contents: "+err.Error()))
			return
		}
	}

	text, err := makeCfg(inf.Tx.Tx, inf.Config, profileData, fileName)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("making config: "+err.Error()))
		return
	}

	if text == "" {
		// TODO replicates old Perl; verify required.
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusNotFound, errors.New("not found"), nil)
		return
	}

	if contentType != "" {
		w.Header().Set(tc.ContentType, contentType)
	}
	w.Write([]byte(hdr + text))
}

// GetCDNNameFromNameOrID takes a string which is a CDN name or ID, and returns the CDN name, whether the CDN exists, and any error.
// Note if cdnNameOrID is the CDN name, it will not be checked for existence!
func GetCDNNameFromNameOrID(tx *sql.Tx, nameOrID string) (tc.CDNName, bool, error) {
	if cdnID, err := strconv.Atoi(nameOrID); err == nil {
		return dbhelpers.GetCDNNameFromID(tx, int64(cdnID))
	}
	return tc.CDNName(nameOrID), true, nil
}
