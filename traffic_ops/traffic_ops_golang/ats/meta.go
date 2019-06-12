package ats

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
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
)

func GetConfigMetaData(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"server-name-or-id"}, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	serverName, userErr, sysErr, errCode := getServerNameFromNameOrID(inf.Tx.Tx, inf.Params["server-name-or-id"])
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}

	server, ok, err := getServerInfoByHost(inf.Tx.Tx, serverName)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("GetConfigMetaData getting server info: "+err.Error()))
		return
	} else if !ok {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusNotFound, nil, errors.New("server not found"))
		return
	}

	tmParams, err := GetTMParams(inf.Tx.Tx)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("GetConfigMetaData getting tm.url parameter: "+err.Error()))
		return
	}
	if tmParams.URL == "" {
		log.Errorln("ats.GetConfigMetadata: global tm.url parameter missing or empty! Setting empty in meta config!")
	}

	atsData := tc.ATSConfigMetaData{
		Info: tc.ATSConfigMetaDataInfo{
			ProfileID:         int(server.ProfileID),
			TOReverseProxyURL: tmParams.ReverseProxyURL,
			TOURL:             tmParams.URL,
			ServerIPv4:        server.IP,
			ServerPort:        server.Port,
			ServerName:        server.HostName,
			CDNID:             server.CDNID,
			CDNName:           string(server.CDN),
			ServerID:          server.ID,
			ProfileName:       server.ProfileName,
		},
		ConfigFiles: []tc.ATSConfigMetaDataConfigFile{},
	}

	locationParams, err := GetLocationParams(inf.Tx.Tx, int(server.ProfileID))
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("GetConfigMetaData getting location parameters: "+err.Error()))
		return
	}

	if locationParams["remap.config"].Location != "" {
		configLocation := locationParams["remap.config"].Location
		uriSignedDSes, err := GetServerURISignedDSes(inf.Tx.Tx, server.HostName, server.Port)
		if err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("GetConfigMetaData getting server uri-signed dses: "+err.Error()))
			return
		}
		for _, ds := range uriSignedDSes {
			cfgName := "uri_signing_" + string(ds) + ".config"
			// If there's already a parameter for it, don't clobber it. The user may wish to override the location.
			if _, ok := locationParams[cfgName]; !ok {
				p := locationParams[cfgName]
				p.FileNameOnDisk = cfgName
				p.Location = configLocation
			}
		}
	}

	for cfgFile, cfgParams := range locationParams {
		atsCfg := tc.ATSConfigMetaDataConfigFile{
			FileNameOnDisk: cfgParams.FileNameOnDisk,
			Location:       cfgParams.Location,
		}

		scope, err := getServerScope(inf.Tx.Tx, cfgFile, server.Type)
		if err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("GetConfigMetaData getting scope: "+err.Error()))
			return
		}

		if cfgParams.URL != "" {
			scope = tc.ATSConfigMetaDataConfigFileScopeCDNs
			atsCfg.URL = cfgParams.URL
		} else {
			scopeID := ""
			if scope == tc.ATSConfigMetaDataConfigFileScopeCDNs {
				scopeID = string(server.CDN)
			} else if scope == tc.ATSConfigMetaDataConfigFileScopeProfiles {
				scopeID = server.ProfileName
			} else { // ATSConfigMetaDataConfigFileScopeServers
				scopeID = server.HostName
			}
			atsCfg.APIURI = "/api/1.2/" + string(scope) + "/" + scopeID + "/configfiles/ats/" + cfgFile
		}

		atsCfg.Scope = string(scope)

		atsData.ConfigFiles = append(atsData.ConfigFiles, atsCfg)
	}

	api.WriteRespRaw(w, r, atsData)
}

// GetTMParams returns the global "tm.url" and "tm.rev_proxy.url" parameters, and any error. If either param doesn't exist, an empty string is returned without error.
func GetTMParams(tx *sql.Tx) (TMParams, error) {
	rows, err := tx.Query(`SELECT name, value from parameter where config_file = 'global' AND (name = 'tm.url' OR name = 'tm.rev_proxy.url')`)
	if err != nil {
		return TMParams{}, errors.New("querying: " + err.Error())
	}
	defer rows.Close()

	p := TMParams{}
	for rows.Next() {
		name := ""
		val := ""
		if err := rows.Scan(&name, &val); err != nil {
			return TMParams{}, errors.New("scanning: " + err.Error())
		}
		if name == "tm.url" {
			p.URL = val
		} else if name == "tm.rev_proxy.url" {
			p.ReverseProxyURL = val
		} else {
			return TMParams{}, errors.New("querying got unexpected parameter: " + name + " (value: '" + val + "')") // should never happen
		}
	}
	return p, nil
}

// GetLocationParams returns a map[configFile]locationParams, and any error. If either param doesn't exist, an empty string is returned without error.
func GetLocationParams(tx *sql.Tx, profileID int) (map[string]ConfigProfileParams, error) {
	qry := `
SELECT
  p.name,
  p.config_file,
  p.value
FROM
  parameter p
  JOIN profile_parameter pp ON pp.parameter = p.id
WHERE
  pp.profile = $1
`
	rows, err := tx.Query(qry, profileID)
	if err != nil {
		return nil, errors.New("querying: " + err.Error())
	}
	defer rows.Close()

	params := map[string]ConfigProfileParams{}
	for rows.Next() {
		name := ""
		file := ""
		val := ""
		if err := rows.Scan(&name, &file, &val); err != nil {
			return nil, errors.New("scanning: " + err.Error())
		}
		if name == "location" {
			p := params[file]
			p.FileNameOnDisk = file
			p.Location = val
			params[file] = p
		} else if name == "URL" {
			p := params[file]
			p.URL = val
			params[file] = p
		}
	}
	return params, nil
}

// GetServerURISignedDSes returns a list of delivery service names which have the given server assigned and have URI signing enabled, and any error.
func GetServerURISignedDSes(tx *sql.Tx, serverHostName string, serverPort int) ([]tc.DeliveryServiceName, error) {
	qry := `
SELECT
  ds.xml_id
FROM
  deliveryservice ds
  JOIN deliveryservice_server dss ON ds.id = dss.deliveryservice
  JOIN server s ON s.id = dss.server
WHERE
  s.host_name = $1
  AND s.tcp_port = $2
  AND ds.signing_algorithm = 'uri_signing'
`
	rows, err := tx.Query(qry, serverHostName, serverPort)
	if err != nil {
		return nil, errors.New("querying: " + err.Error())
	}
	defer rows.Close()

	dses := []tc.DeliveryServiceName{}
	for rows.Next() {
		ds := tc.DeliveryServiceName("")
		if err := rows.Scan(&ds); err != nil {
			return nil, errors.New("scanning: " + err.Error())
		}
		dses = append(dses, ds)
	}
	return dses, nil
}

func getServerScope(tx *sql.Tx, cfgFile string, serverType string) (tc.ATSConfigMetaDataConfigFileScope, error) {
	switch {
	case cfgFile == "cache.config" && tc.CacheTypeFromString(serverType) == tc.CacheTypeMid:
		return tc.ATSConfigMetaDataConfigFileScopeServers, nil
	default:
		return getScope(tx, cfgFile)
	}
}

// getScope returns the ATSConfigMetaDataConfigFileScope for the given config file, and potentially the given server. If the config is not a Server scope, i.e. was part of an endpoint which does not include a server name or id, the server may be nil.
func getScope(tx *sql.Tx, cfgFile string) (tc.ATSConfigMetaDataConfigFileScope, error) {
	switch {
	case cfgFile == "ip_allow.config":
		fallthrough
	case cfgFile == "parent.config":
		fallthrough
	case cfgFile == "hosting.config":
		fallthrough
	case cfgFile == "packages":
		fallthrough
	case cfgFile == "chkconfig":
		fallthrough
	case cfgFile == "remap.config":
		fallthrough
	case strings.HasPrefix(cfgFile, "to_ext_") && strings.HasSuffix(cfgFile, ".config"):
		return tc.ATSConfigMetaDataConfigFileScopeServers, nil
	case cfgFile == "12M_facts":
		fallthrough
	case cfgFile == "50-ats.rules":
		fallthrough
	case cfgFile == "astats.config":
		fallthrough
	case cfgFile == "cache.config":
		fallthrough
	case cfgFile == "drop_qstring.config":
		fallthrough
	case cfgFile == "logs_xml.config":
		fallthrough
	case cfgFile == "logging.config":
		fallthrough
	case cfgFile == "plugin.config":
		fallthrough
	case cfgFile == "records.config":
		fallthrough
	case cfgFile == "storage.config":
		fallthrough
	case cfgFile == "volume.config":
		fallthrough
	case cfgFile == "sysctl.conf":
		fallthrough
	case strings.HasPrefix(cfgFile, "url_sig_") && strings.HasSuffix(cfgFile, ".config"):
		fallthrough
	case strings.HasPrefix(cfgFile, "uri_signing_") && strings.HasSuffix(cfgFile, ".config"):
		return tc.ATSConfigMetaDataConfigFileScopeProfiles, nil

	case cfgFile == "bg_fetch.config":
		fallthrough
	case cfgFile == "regex_revalidate.config":
		fallthrough
	case cfgFile == "ssl_multicert.config":
		fallthrough
	case strings.HasPrefix(cfgFile, "cacheurl") && strings.HasSuffix(cfgFile, ".config"):
		fallthrough
	case strings.HasPrefix(cfgFile, "hdr_rw_") && strings.HasSuffix(cfgFile, ".config"):
		fallthrough
	case strings.HasPrefix(cfgFile, "regex_remap_") && strings.HasSuffix(cfgFile, ".config"):
		fallthrough
	case strings.HasPrefix(cfgFile, "set_dscp_") && strings.HasSuffix(cfgFile, ".config"):
		return tc.ATSConfigMetaDataConfigFileScopeCDNs, nil
	}

	scope, ok, err := GetFirstScopeParameter(tx, cfgFile)
	if err != nil {
		return tc.ATSConfigMetaDataConfigFileScopeInvalid, errors.New("getting scope parameter: " + err.Error())
	}
	if !ok {
		scope = string(tc.ATSConfigMetaDataConfigFileScopeServers)
	}
	return tc.ATSConfigMetaDataConfigFileScope(scope), nil
}

type TMParams struct {
	URL             string
	ReverseProxyURL string
}

type ConfigProfileParams struct {
	FileNameOnDisk string
	Location       string
	URL            string
	APIURI         string
}

// GetFirstScopeParameter returns the value of the arbitrarily-first parameter with the name 'scope' and the given config file, whether a parameter was found, and any error.
func GetFirstScopeParameter(tx *sql.Tx, cfgFile string) (string, bool, error) {
	v := ""
	if err := tx.QueryRow(`SELECT p.value FROM parameter p WHERE p.config_file = $1 AND p.name = 'scope'`, cfgFile).Scan(&v); err != nil {
		if err == sql.ErrNoRows {
			return "", false, nil
		}
		return "", false, errors.New("querying first scope parameter: " + err.Error())
	}
	return v, true, nil
}
