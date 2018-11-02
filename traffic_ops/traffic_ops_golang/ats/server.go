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
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
)

func GetServerConfigRemap(w http.ResponseWriter, r *http.Request) {
	// TODO accept names or ids
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"id"}, []string{"id"})
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	serverID := inf.IntParams["id"]

	serverInfo, ok, err := getServerInfo(inf.Tx.Tx, serverID)
	if !ok {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusNotFound, errors.New("server not found"), nil)
		return
	}

	hdr, err := headerComment(inf.Tx.Tx, serverInfo.HostName)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("Getting header comment: "+err.Error()))
		return
	}

	remapDSData, err := GetRemapDSData(inf.Tx.Tx, serverInfo)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("Getting remap ds data: "+err.Error()))
		return
	}

	text := ""
	if tc.CacheTypeFromString(serverInfo.Type) == tc.CacheTypeMid {
		text, err = GetServerConfigRemapDotConfigForMid(inf.Tx.Tx, serverInfo, remapDSData, hdr)
	} else {
		text, err = GetServerConfigRemapDotConfigForEdge(inf.Tx.Tx, serverInfo, remapDSData, hdr)
	}

	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("Getting server config remap dot config: "+err.Error()))
		return
	}

	// TODO figure out why Commit() hangs
	if err := inf.Tx.Tx.Rollback(); err != nil && err != sql.ErrTxDone {
		log.Errorln("rolling back transaction: " + err.Error())
	}
	w.Header().Set("Content-Type", "text/plain")
	io.WriteString(w, text)
}

func DSProfileIDs(dses []RemapConfigDSData) []int {
	dsProfileIDs := []int{}
	for _, ds := range dses {
		if ds.ProfileID != nil {
			// TODO determine if this is right, if the DS has no profile
			dsProfileIDs = append(dsProfileIDs, *ds.ProfileID)
		}
	}
	return dsProfileIDs
}

func GetServerConfigRemapDotConfigForMid(tx *sql.Tx, server *ServerInfo, dses []RemapConfigDSData, header string) (string, error) {

	profilesCacheKeyConfigParams, err := GetProfilesParamData(tx, DSProfileIDs(dses), "cachekey.config") // (map[int]map[string]string, error) {
	if err != nil {
		return "", errors.New("getting profiles param data: " + err.Error())
	}

	midRemaps := map[string]string{}
	for _, ds := range dses {
		if ds.Type.IsLive() && ds.Type.IsNational() {
			continue // Live local delivery services skip mids
		}
		if ds.OriginFQDN == nil || *ds.OriginFQDN != "" {
			log.Warnf("GetServerConfigRemapDotConfigForMid ds '" + ds.Name + "' has no origin fqdn, skipping!") // TODO confirm - Perl uses without checking!
			continue
		}

		if midRemaps[*ds.OriginFQDN] != "" {
			continue // skip remap rules from extra HOST_REGEXP entries
		}

		midRemap := ""
		if ds.MidHeaderRewrite != nil && *ds.MidHeaderRewrite != "" {
			midRemap += ` @plugin=header_rewrite.so @pparam=` + MidHeaderRewriteConfigFileName(ds.Name)
		}
		if ds.QStringIgnore != nil && *ds.QStringIgnore == tc.QueryStringIgnoreIgnoreInCacheKeyAndPassUp {
			atsMajorVersion, err := GetATSMajorVersion(tx, server.ProfileID)
			if err != nil {
				return "", errors.New("getting ATS major version: " + err.Error())
			}
			midRemap += GetQStringIgnoreRemap(atsMajorVersion)
		}
		if ds.CacheURL != nil && *ds.CacheURL != "" {
			midRemap += ` @plugin=cacheurl.so @pparam=` + CacheURLConfigFileName(ds.Name)
		}

		if ds.ProfileID != nil && len(profilesCacheKeyConfigParams[*ds.ProfileID]) > 0 {
			for name, val := range profilesCacheKeyConfigParams[*ds.ProfileID] {
				midRemap += ` @pparam=--` + name + "=" + val
			}
		}
		if ds.RangeRequestHandling != nil && *ds.RangeRequestHandling == tc.RangeRequestHandlingCacheRangeRequest {
			midRemap += ` @plugin=cache_range_requests.so`
		}
		midRemaps[*ds.OriginFQDN] = midRemap
	}

	textLines := []string{}
	for originFQDN, midRemap := range midRemaps {
		textLines = append(textLines, "map "+originFQDN+" "+originFQDN+midRemap+"\n")
	}
	sort.Sort(sort.StringSlice(textLines))

	text := header
	for _, line := range textLines {
		text += line
	}
	return text, nil
}

func GetServerConfigRemapDotConfigForEdge(tx *sql.Tx, server *ServerInfo, dses []RemapConfigDSData, header string) (string, error) {
	pData, err := serverParamData(tx, server.ProfileID, "package", server.HostName, server.DomainName)
	if err != nil {
		return "", errors.New("getting server param data: " + err.Error())
	}

	profilesCacheKeyConfigParams, err := GetProfilesParamData(tx, DSProfileIDs(dses), "cachekey.config") // (map[int]map[string]string, error) {

	textLines := []string{}

	// DEBUG
	log.Errorln("DEBUG GetServerConfigRemapDotConfigForEdge dses {{")
	for _, ds := range dses {
		log.Errorln(ds.Name)
	}
	log.Errorln("}}")

	for _, ds := range dses {
		remapText := ""
		if ds.Type == tc.DSTypeAnyMap {
			if ds.RemapText == nil {
				log.Errorln("ds '" + ds.Name + "' is ANY_MAP, but has no remap text - skipping")
				continue
			}
			remapText = *ds.RemapText + "\n"
			textLines = append(textLines, remapText)
			continue
		}

		remapLines, err := MakeEdgeDSDataRemapLines(ds, server)
		if err != nil {
			return "", errors.New("making remap lines: " + err.Error()) // TODO log and continue?
		}
		for _, line := range remapLines {
			profilecacheKeyConfigParams := (map[string]string)(nil)
			if ds.ProfileID != nil {
				profilecacheKeyConfigParams = profilesCacheKeyConfigParams[*ds.ProfileID]
			}
			remapText, err = BuildRemapLine(tx, server, pData, remapText, ds, line.From, line.To, profilecacheKeyConfigParams)
			if err != nil {
				return "", errors.New("ds '" + ds.Name + "' building remap line: " + err.Error()) // TODO log and continue?
			}
		}
		textLines = append(textLines, remapText)
	}

	text := header
	sort.Sort(sort.StringSlice(textLines))
	for _, line := range textLines {
		text += line
	}
	return text, nil
}

// BuildRemapLine builds the remap line for the given server and delivery service.
// The cacheKeyConfigParams map may be nil, if this ds profile had no cache key config params.
func BuildRemapLine(tx *sql.Tx, server *ServerInfo, pData map[string]string, text string, ds RemapConfigDSData, mapFrom string, mapTo string, cacheKeyConfigParams map[string]string) (string, error) {
	// ds = 'remap' in perl

	mapFrom = strings.Replace(mapFrom, `__http__`, server.HostName, -1)

	if _, hasDSCPRemap := pData["dscp_remap"]; hasDSCPRemap {
		text += "map	" + mapFrom + "     " + mapTo + ` @plugin=dscp_remap.so @pparam=` + strconv.Itoa(ds.DSCP)
	} else {
		text += "map	" + mapFrom + "     " + mapTo + ` @plugin=header_rewrite.so @pparam=dscp/set_dscp_` + strconv.Itoa(ds.DSCP) + ".config"
	}

	if ds.EdgeHeaderRewrite != nil && *ds.EdgeHeaderRewrite != "" {
		text += ` @plugin=header_rewrite.so @pparam=` + EdgeHeaderRewriteConfigFileName(ds.Name)
	}

	if ds.SigningAlgorithm != nil && *ds.SigningAlgorithm != "" {
		if *ds.SigningAlgorithm == tc.SigningAlgorithmURLSig {
			text += ` @plugin=url_sig.so @pparam=url_sig_` + ds.Name + ".config"
		} else if *ds.SigningAlgorithm == tc.SigningAlgorithmURISigning {
			text += ` @plugin=uri_signing.so @pparam=uri_signing_` + ds.Name + ".config"
		}
	}

	if ds.QStringIgnore != nil {
		if *ds.QStringIgnore == tc.QueryStringIgnoreDropAtEdge {
			dqsFile := "drop_qstring.config"
			text += ` @plugin=regex_remap.so @pparam=` + dqsFile
		} else if *ds.QStringIgnore == tc.QueryStringIgnoreIgnoreInCacheKeyAndPassUp {
			_, globalExists, err := GetProfileParamValue(tx, server.ProfileID, "cacheurl.config", "location")
			if err != nil {
				return "", errors.New("getting profile param value for cacheurl.config location: " + err.Error())
			}
			if globalExists {
				log.Debugln("qstring_ignore == 1, but global cacheurl.config param exists, so skipping remap rename config_file=cacheurl.config parameter if you want to change") // TODO warn? Perl was a debug
			} else {
				atsMajorVersion, err := GetATSMajorVersion(tx, server.ProfileID) // TODO only call once, pass to this func
				if err != nil {
					return "", errors.New("getting ATS major version: " + err.Error())
				}
				text += GetQStringIgnoreRemap(atsMajorVersion)
			}
		}
	}

	if ds.CacheURL != nil && *ds.CacheURL != "" {
		text += ` @plugin=cacheurl.so @pparam=` + CacheURLConfigFileName(ds.Name)
	}

	// DEBUG - sort to diff identical to perl - remove once diff is verified
	if len(cacheKeyConfigParams) > 0 {
		text += ` @plugin=cachekey.so`

		keys := []string{}
		for key, _ := range cacheKeyConfigParams {
			keys = append(keys, key)
		}
		sort.Sort(sort.StringSlice(keys))

		for _, key := range keys {
			text += ` @pparam=--` + key + "=" + cacheKeyConfigParams[key]
		}
	}

	// if len(cacheKeyConfigParams) > 0 {
	// 	text += ` @plugin=cachekey.so`
	// 	for key, val := range cacheKeyConfigParams {
	// 		text += ` @pparam=--` + key + "=" + val
	// 	}
	// }

	// Note: should use full path here?
	if ds.RegexRemap != nil && *ds.RegexRemap != "" {
		text += ` @plugin=regex_remap.so @pparam=regex_remap_` + ds.Name + ".config"
	}
	if ds.RangeRequestHandling != nil {
		if *ds.RangeRequestHandling == tc.RangeRequestHandlingBackgroundFetch {
			text += ` @plugin=background_fetch.so @pparam=bg_fetch.config`
		} else if *ds.RangeRequestHandling == tc.RangeRequestHandlingCacheRangeRequest {
			text += ` @plugin=cache_range_requests.so `
		}
	}
	if ds.RemapText != nil && *ds.RemapText != "" {
		text += " " + *ds.RemapText
	}

	if ds.FQPacingRate != nil && *ds.FQPacingRate > 0 {
		text += ` @plugin=fq_pacing.so @pparam=--rate=` + strconv.Itoa(*ds.FQPacingRate)
	}
	text += "\n"
	return text, nil
}

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
