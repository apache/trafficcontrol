package ats

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/apache/trafficcontrol/lib/go-atscfg"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
)

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

const configSuffix = ".config"

const HeaderRewritePrefix = "hdr_rw_"
const RegexRemapPrefix = "regex_remap_"
const CacheUrlPrefix = "cacheurl_"

const RemapFile = "remap.config"

func GetConfigFile(prefix string, xmlId string) string {
	return prefix + xmlId + configSuffix
}

func GetNameVersionString(tx *sql.Tx) (string, error) {
	toolName, url, err := GetToolNameAndURL(tx)
	if err != nil {
		return "", errors.New("getting toolname and url parameters: " + err.Error())
	}
	return atscfg.GetNameVersionStringFromToolNameAndURL(toolName, url), nil
}

func GetToolNameAndURL(tx *sql.Tx) (string, string, error) {
	qry := `
SELECT
  p.name,
  p.value
FROM
  parameter p
WHERE
  (p.name = 'tm.toolname' OR p.name = 'tm.url') AND p.config_file = 'global'
`
	rows, err := tx.Query(qry)
	if err != nil {
		return "", "", errors.New("querying: " + err.Error())
	}
	defer rows.Close()

	toolName := ""
	url := ""
	for rows.Next() {
		name := ""
		val := ""
		if err := rows.Scan(&name, &val); err != nil {
			return "", "", errors.New("scanning: " + err.Error())
		}
		if name == "tm.toolname" {
			toolName = val
		} else if name == "tm.url" {
			url = val
		}
	}
	return toolName, url, nil
}

// getCDNNameFromNameOrID returns the CDN name from a parameter which may be the name or ID.
// This also checks and verifies the existence of the given CDN, and returns an appropriate user error if it doesn't exist.
// Returns the name, any user error, any system error, and any error code.
func getCDNNameFromNameOrID(tx *sql.Tx, cdnNameOrID string) (string, error, error, int) {
	if cdnID, err := strconv.Atoi(cdnNameOrID); err == nil {
		cdnName, ok, err := dbhelpers.GetCDNNameFromID(tx, int64(cdnID))
		if err != nil {
			return "", nil, fmt.Errorf("getting CDN name from id %v: %v", cdnID, err), http.StatusInternalServerError
		} else if !ok {
			return "", errors.New("cdn not found"), nil, http.StatusNotFound
		}
		return string(cdnName), nil, nil, http.StatusOK
	}

	cdnName := cdnNameOrID
	if ok, err := dbhelpers.CDNExists(cdnName, tx); err != nil {
		return "", nil, fmt.Errorf("checking CDN name '%v' existence: %v", cdnName, err), http.StatusInternalServerError
	} else if !ok {
		return "", errors.New("cdn not found"), nil, http.StatusNotFound
	}
	return cdnName, nil, nil, http.StatusOK
}

// getServerNameFromNameOrID returns the server name from a parameter which may be the name or ID.
// This also checks and verifies the existence of the given server, and returns an appropriate user error if it doesn't exist.
// Returns the name, any user error, any system error, and any error code.
func getServerNameFromNameOrID(tx *sql.Tx, serverNameOrID string) (string, error, error, int) {
	if serverID, err := strconv.Atoi(serverNameOrID); err == nil {
		serverName, ok, err := dbhelpers.GetServerNameFromID(tx, int64(serverID))
		if err != nil {
			return "", nil, fmt.Errorf("getting server name from id %v: %v", serverID, err), http.StatusInternalServerError
		} else if !ok {
			return "", errors.New("server not found"), nil, http.StatusNotFound
		}
		return string(serverName), nil, nil, http.StatusOK
	}

	serverName := serverNameOrID
	if ok, err := dbhelpers.ServerExists(serverName, tx); err != nil {
		return "", nil, fmt.Errorf("checking server name '%v' existence: %v", serverName, err), http.StatusInternalServerError
	} else if !ok {
		return "", errors.New("server not found"), nil, http.StatusNotFound
	}
	return serverName, nil, nil, http.StatusOK
}

func headerComment(tx *sql.Tx, name string) (string, error) {
	nameVersionStr, err := GetNameVersionString(tx)
	if err != nil {
		return "", errors.New("getting name version string: " + err.Error())
	}
	return atscfg.HeaderCommentWithTOVersionStr(name, nameVersionStr), nil
}
