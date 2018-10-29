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

// import (
// 	"database/sql"
// 	"errors"
// 	"net/http"
// 	"strings"

// 	"github.com/apache/trafficcontrol/lib/go-log"
// 	"github.com/apache/trafficcontrol/lib/go-tc"
// 	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
// 	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/config"
// )

// func GetServerConfig(w http.ResponseWriter, r *http.Request) {
// 	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"id", "file"}, []string{"id"})
// 	if userErr != nil || sysErr != nil {
// 		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
// 		return
// 	}
// 	defer inf.Close()

// 	serverID := inf.IntParams["id"]
// 	fileName := inf.Params["file"].TrimSuffix(fileName, ".json")

// 	serverInfo, ok, err := getServerInfo(inf.Tx.Tx, serverID)
// 	if !ok {
// 		api.HandleErr(w, r, inf.Tx.Tx, http.StatusNotFound, "server not found", nil)
// 		return
// 	}

// 	scope := getServerScope(inf.Tx.Tx, fileName, serverInfo.TypeName)
// 	if scope != tc.ATSConfigMetaDataConfigFileScopeServers {
// 		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, "Error - incorrect file scope for route used.  Please use the "+string(scope)+" route.", nil)
// 		return
// 	}

// 	fileContents := ""

// 	// TODO add routes for each of these, rather than dispatching ourselves
// 	switch {
// 	case fileName == "cache.config":
// 		fileContents = serverCacheDotConfig(inf.Tx.Tx, serverInfo, fileName)
// 	case fileName == "ip_allow.config":
// 		fileContents = ipAllowDotConfig(inf.Tx.Tx, serverInfo, fileName)
// 	case fileName == "parent.config":
// 		fileContents = parentDotConfig(inf.Tx.Tx, serverInfo, fileName)
// 	case fileName == "hosting.config":
// 		fileContents = hostingDotConfig(inf.Tx.Tx, serverInfo, fileName)
// 	case fileName == "remap.config":
// 		fileContents = remapDotConfig(inf.Tx.Tx, serverInfo, fileName)
// 	case fileName == "packages":
// 		packageVersion = getPackageVersions(inf.Tx.Tx, serverInfo, fileName)
// 		packageJSON, err := json.Marshal(&packageVersion)
// 		if err != nil {
// 			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("marshalling package: "+err.Error()))
// 			return
// 		}
// 		fileContents = string(packageJSON)
// 	case fileName == "chkconfig":
// 		chkConfig = getChkConfig(inf.Tx.Tx, serverInfo, fileName)
// 		chkConfigJSON, err := json.Marshal(&chkConfig)
// 		if err != nil {
// 			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("marshalling chkconfig: "+err.Error()))
// 			return
// 		}
// 		fileContents = string(chkConfigJSON)
// 	case strings.HasPrefix(fileName, "to_ext_") && strings.HasSuffix(fileName, ".config"):
// 		fileContents = toExtDotConfig(inf.Tx.Tx, serverInfo, fileName)
// 	default:
// 		// TODO move to func, "getUnknownServerConfig"
// 		fileParam, ok, err := getConfigParam(inf.Tx.Tx, fileName)
// 		if err != nil {
// 			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("GetServerConfig getting config param: "+err.Error()))
// 			return
// 		} else if !ok {
// 			api.HandleErr(w, r, inf.Tx.Tx, http.StatusNotFound, "not found", nil)
// 			return
// 		}
// 		fileContents = takeAndBakeServer(inf.Tx.Tx, serverInfo.ProfileID, fileName)
// 	}
// 	if fileContents == "" {
// 		// TODO replicates old Perl; verify required.
// 		api.HandleErr(w, r, inf.Tx.Tx, http.StatusNotFound, "not found", nil)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "text/plain")
// 	w.Write([]byte{fileContents})
// }

// func takeAndBakeServer(tx *sql.Tx, profileName string, configFile string) (string, error) {
// 	paramData, err := serverParamData(tx, profileName, configFile)
// 	if err != nil {
// 		return "", errors.New("getting server param data: " + err.Error())
// 	}

// 	hdr, err := headerComment(profileName)
// 	if err != nil {
// 		return "", errors.New("getting header comment: " + err.Error())
// 	}
// 	for paramName, paramVal := range params {
// 		if paramName == "header" {
// 			if paramVal == "none" {
// 				hdr = ""
// 			} else {
// 				hdr = paramVal + "\n"
// 			}
// 		} else {
// 			text += paramVal + "\n"
// 		}
// 	}
// 	text = strings.Replace(text, "__RETURN__", "\n")
// 	return header + text
// }

// func serverCacheDotConfig(tx *sql.Tx, serverInfo *ServerInfo, fileName string) (string, error) {
// 	// TODO fix duplicate code with profileCacheDotConfig
// 	text, err := headerComment(profile.Name)
// 	if err != nil {
// 		return "", errors.New("getting header comment: " + err.Error())
// 	}

// 	dsData, err := getDSData(tx, serverInfo.ID)
// 	if err != nil {
// 		return "", errors.New("getting ds data: " + err.Error())
// 	}

// 	for _, ds := range dsData {
// 		if ds.Type != tc.DSTypeHTTPNoCache {
// 			continue
// 		}
// 		originFQDN, originPort := getHostPortFromURI(ds.OriginFQDN)
// 		if originPort != "" {
// 			l := "dest_domain=" + originFQDN + " port=" + originPort + " scheme=http action=never-cache\n"
// 			lines[l] = struct{}{}
// 		} else {
// 			l := "dest_domain=" + originFQDN + " scheme=http action=never-cache\n"
// 			lines[l] = struct{}{}
// 		}
// 	}

// 	text := ""
// 	for line, _ := range lines {
// 		text = line + "\n"
// 	}
// 	return text, nil
// }

// func ipAllowDotConfig(tx *sql.Tx, serverInfo *ServerInfo, fileName string) (string, error) {
// 	text, err := headerComment(profile.Name)
// 	if err != nil {
// 		return "", errors.New("getting header comment: " + err.Error())
// 	}

// 	ipAllowData, err := getIPAllowData(tx, serverInfo, fileName)
// 	if err != nil {
// 		return "", errors.New("getting ip allow data: " + err.Error())
// 	}

// 	allowedV4s := []string{}
// 	allowedV6s := []string{}
// 	for _, server := range ipAllowData {
// 		allowedV4s := append(allowedV4s, server.IPv4)
// 		allowedV6s := append(allowedV6s, server.IPv6)
// 		// TODO netmask?
// 	}

// 	compactV4s, err := GetIPv4CIDRs(allowedV4s)
// 	if err != nil {
// 		return "", errors.New("parsing IPv4s: " + err.Error())
// 	}

// 	compactV6s, err := GetIPv6CIDRs(allowedV6s)
// 	if err != nil {
// 		return "", errors.New("parsing IPv6s: " + err.Error())
// 	}

// 	alloweds := []IPAllowAccess{}
// 	for _, ip := range compactV4s {
// 		alloweds = append(alloweds, IPAllowAccess{SourceIP: ip, Action: "ip_allow", Method: "ALL"})
// 	}
// 	for _, ip := range compactV6s {
// 		alloweds = append(alloweds, IPAllowAccess{SourceIP: ip, Action: "ip_allow", Method: "ALL"})
// 	}

// 	// allow RFC 1918 server space - TODO JvD: parameterize
// 	alloweds = append(alloweds, IPAllowAccess{SourceIP: "10.0.0.0-10.255.255.255", Action: "ip_allow", Method: "ALL"})
// 	alloweds = append(alloweds, IPAllowAccess{SourceIP: "172.16.0.0-172.31.255.255", Action: "ip_allow", Method: "ALL"})
// 	alloweds = append(alloweds, IPAllowAccess{SourceIP: "192.168.0.0-192.168.255.255", Action: "ip_allow", Method: "ALL"})

// 	//end with a deny
// 	alloweds = append(alloweds, IPAllowAccess{SourceIP: "0.0.0.0-255.255.255.255", Action: "ip_deny", Method: "ALL"})
// 	alloweds = append(alloweds, IPAllowAccess{SourceIP: "::-ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff", Action: "ip_deny", Method: "ALL"})

// 	// for _, access := range ipAllowData {
// 	// 	// TODO verify format string, perl == go
// 	// 	text += fmt.Sprintf("src_ip=%-70s action=%-10s method=%-20s\n", access.SourceIP, access.Action, access.Method)
// 	// }
// 	return text, nil
// }

// func parentDotConfig(tx *sql.Tx, server *ServerInfo, fileName string) (string, error) {
// 	isTopLevelCache := (server.ParentCacheGroupType == tc.OriginLocationType || server.ParentCacheGroupID == -1) && (server.SecondaryParentCacheGroupType == tc.OriginLocationType || server.SecondaryParentCacheGroupID == -1)

// 	atsVer, ok, err := GetProfileParamValue(tx, server.ProfileID, "package", "trafficserver")
// 	if err != nil {
// 		return "", errors.New("getting ats version parameter: " + err.Error())
// 	}
// 	if !ok {
// 		return "", errors.New("no ats version parameter on this profile (config_file 'package', name 'trafficserver')")
// 	}
// 	if len(atsVer) == 0 {
// 		return "", errors.New("empty ats version parameter on this profile (config_file 'package', name 'trafficserver')")
// 	}
// 	atsMajorVer, err := strconv.Atoi(atsVer[:1])
// 	if err != nil {
// 		return "", errors.New("ats version parameter '" + atsVer + "' on this profile is not a number (config_file 'package', name 'trafficserver')")
// 	}

// 	hdr, err := headerComment(profileName)
// 	if err != nil {
// 		return "", errors.New("getting header comment: " + err.Error())
// 	}

// 	if isTopLevelCache {
// 		origins := map[string]struct{}{}

// 		parentData := getParentDSData(server, isTopLevelCache)

// 	}

// }

// func hostingDotConfig(tx *sql.Tx, serverInfo *ServerInfo, fileName string) {

// }

// func remapDotConfig(tx *sql.Tx, serverInfo *ServerInfo, fileName string) {

// }

// func getPackageVersions(tx *sql.Tx, serverInfo *ServerInfo, fileName string) {

// }

// func getChkConfig(tx *sql.Tx, serverInfo *ServerInfo, fileName string) {

// }

// func toExtDotConfig(tx *sql.Tx, serverInfo *ServerInfo, fileName string) {

// }
