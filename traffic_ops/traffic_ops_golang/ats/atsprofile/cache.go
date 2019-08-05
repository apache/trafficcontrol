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

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/ats"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/config"
)

func GetCache(w http.ResponseWriter, r *http.Request) {
	WithProfileData(w, r, makeCache)
}

func makeCache(tx *sql.Tx, _ *config.Config, profile ats.ProfileData, _ string) (string, error) {
	lines := map[string]struct{}{} // use a "set" for lines, to avoid duplicates, since we're looking up by profile
	profileDSes, err := ats.GetProfileDS(tx, profile.ID)
	if err != nil {
		return "", errors.New("getting profile delivery services: " + err.Error())
	}

	for _, ds := range profileDSes {
		if ds.Type != tc.DSTypeHTTPNoCache {
			continue
		}
		if ds.OriginFQDN == nil || *ds.OriginFQDN == "" {
			log.Warnf("profileCacheDotConfig ds has no origin fqdn, skipping!") // TODO add ds name to data loaded, to put it in the error here?
			continue
		}
		originFQDN, originPort := getHostPortFromURI(*ds.OriginFQDN)
		if originPort != "" {
			l := "dest_domain=" + originFQDN + " port=" + originPort + " scheme=http action=never-cache\n"
			lines[l] = struct{}{}
		} else {
			l := "dest_domain=" + originFQDN + " scheme=http action=never-cache\n"
			lines[l] = struct{}{}
		}
	}

	text := ""
	for line, _ := range lines {
		text += line
	}
	if text == "" {
		text = "\n" // If no params exist, don't send "not found," but an empty file. We know the profile exists.
	}
	return text, nil
}

func getHostPortFromURI(uriStr string) (string, string) {
	originFQDN := uriStr
	originFQDN = strings.TrimPrefix(originFQDN, "http://")
	originFQDN = strings.TrimPrefix(originFQDN, "https://")

	slashPos := strings.Index(originFQDN, "/")
	if slashPos != -1 {
		originFQDN = originFQDN[:slashPos]
	}
	portPos := strings.Index(originFQDN, ":")
	portStr := ""
	if portPos != -1 {
		portStr = originFQDN[portPos+1:]
		originFQDN = originFQDN[:portPos]
	}
	return originFQDN, portStr
}
