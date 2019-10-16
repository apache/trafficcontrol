package atscfg

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
	"errors"
	"sort"
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
)

const CacheURLParameterConfigFile = "cacheurl.config"
const CacheKeyParameterConfigFile = "cachekey.config"

type RemapConfigDSData struct {
	ID                       int
	Type                     tc.DSType
	OriginFQDN               *string
	MidHeaderRewrite         *string
	CacheURL                 *string
	RangeRequestHandling     *int
	CacheKeyConfigParams     map[string]string
	RemapText                *string
	EdgeHeaderRewrite        *string
	SigningAlgorithm         *string
	Name                     string
	QStringIgnore            *int
	RegexRemap               *string
	FQPacingRate             *int
	DSCP                     int
	RoutingName              *string
	MultiSiteOrigin          *string
	Pattern                  *string
	RegexType                *string
	Domain                   *string
	RegexSetNumber           *string
	OriginShield             *string
	ProfileID                *int
	Protocol                 *int
	AnonymousBlockingEnabled *bool
	Active                   bool
}

func MakeRemapDotConfig(
	serverName tc.CacheName,
	toToolName string, // tm.toolname global parameter (TODO: cache itself?)
	toURL string, // tm.url global parameter (TODO: cache itself?)
	atsMajorVersion int,
	cacheURLConfigParams map[string]string, // map[name]value for this server's profile and config file 'cacheurl.config'
	dsProfilesCacheKeyConfigParams map[int]map[string]string, // map[profileID][paramName]paramVal for this server's profile, config file 'cachekey.config'
	serverPackageParamData map[string]string, // map[paramName]paramVal for this server, config file 'package'
	serverInfo *ServerInfo, // ServerInfo for this server
	remapDSData []RemapConfigDSData,
) string {
	hdr := GenericHeaderComment(string(serverName), toToolName, toURL)
	text := ""
	if tc.CacheTypeFromString(serverInfo.Type) == tc.CacheTypeMid {
		text = GetServerConfigRemapDotConfigForMid(atsMajorVersion, dsProfilesCacheKeyConfigParams, serverInfo, remapDSData, hdr)
	} else {
		text = GetServerConfigRemapDotConfigForEdge(cacheURLConfigParams, dsProfilesCacheKeyConfigParams, serverPackageParamData, serverInfo, remapDSData, atsMajorVersion, hdr)
	}
	return text
}

func GetServerConfigRemapDotConfigForMid(
	atsMajorVersion int,
	profilesCacheKeyConfigParams map[int]map[string]string,
	server *ServerInfo,
	dses []RemapConfigDSData,
	header string,
) string {
	midRemaps := map[string]string{}
	for _, ds := range dses {
		if ds.Type.IsLive() && !ds.Type.IsNational() {
			continue // Live local delivery services skip mids
		}

		if ds.OriginFQDN == nil || *ds.OriginFQDN == "" {
			log.Warnf("GetServerConfigRemapDotConfigForMid ds '" + ds.Name + "' has no origin fqdn, skipping!") // TODO confirm - Perl uses without checking!
			continue
		}

		if midRemaps[*ds.OriginFQDN] != "" {
			continue // skip remap rules from extra HOST_REGEXP entries
		}

		// multiple uses of cacheurl and cachekey plugins don't work right in ATS, but Perl has always done it.
		// So for now, keep track of it, so we can log an error when it happens.
		hasCacheURL := false
		hasCacheKey := false

		midRemap := ""
		if ds.MidHeaderRewrite != nil && *ds.MidHeaderRewrite != "" {
			midRemap += ` @plugin=header_rewrite.so @pparam=` + MidHeaderRewriteConfigFileName(ds.Name)
		}
		if ds.QStringIgnore != nil && *ds.QStringIgnore == tc.QueryStringIgnoreIgnoreInCacheKeyAndPassUp {
			qstr, addedCacheURL, addedCacheKey := GetQStringIgnoreRemap(atsMajorVersion)
			if addedCacheURL {
				hasCacheURL = true
			}
			if addedCacheKey {
				hasCacheKey = true
			}
			midRemap += qstr
		}
		if ds.CacheURL != nil && *ds.CacheURL != "" {
			if hasCacheURL {
				log.Errorln("DeliveryService qstring_ignore and cacheurl both add cacheurl, but ATS cacheurl doesn't work correctly with multiple entries! Adding anyway!")
			}
			midRemap += ` @plugin=cacheurl.so @pparam=` + CacheURLConfigFileName(ds.Name)
		}

		if ds.ProfileID != nil && len(profilesCacheKeyConfigParams[*ds.ProfileID]) > 0 {
			if hasCacheKey {
				log.Errorln("DeliveryService qstring_ignore and cachekey params both add cachekey, but ATS cachekey doesn't work correctly with multiple entries! Adding anyway!")
			}
			midRemap += ` @plugin=cachekey.so`
			for name, val := range profilesCacheKeyConfigParams[*ds.ProfileID] {
				midRemap += ` @pparam=--` + name + "=" + val
			}
		}
		if ds.RangeRequestHandling != nil && *ds.RangeRequestHandling == tc.RangeRequestHandlingCacheRangeRequest {
			midRemap += ` @plugin=cache_range_requests.so`
		}

		if midRemap != "" {
			midRemaps[*ds.OriginFQDN] = midRemap
		}
	}

	textLines := []string{}
	for originFQDN, midRemap := range midRemaps {
		textLines = append(textLines, "map "+originFQDN+" "+originFQDN+midRemap+"\n")
	}
	sort.Strings(textLines)

	text := header
	text += strings.Join(textLines, "")
	return text
}

func GetServerConfigRemapDotConfigForEdge(
	cacheURLConfigParams map[string]string,
	profilesCacheKeyConfigParams map[int]map[string]string,
	serverPackageParamData map[string]string, // map[paramName]paramVal for this server, config file 'package'
	server *ServerInfo,
	dses []RemapConfigDSData,
	atsMajorVersion int,
	header string,
) string {
	textLines := []string{}

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
			log.Errorln("making remap lines for DS '" + ds.Name + "' - skipping! : " + err.Error())
			continue
		}

		for _, line := range remapLines {
			profilecacheKeyConfigParams := (map[string]string)(nil)
			if ds.ProfileID != nil {
				profilecacheKeyConfigParams = profilesCacheKeyConfigParams[*ds.ProfileID]
			}
			remapText = BuildRemapLine(cacheURLConfigParams, atsMajorVersion, server, serverPackageParamData, remapText, ds, line.From, line.To, profilecacheKeyConfigParams)
		}
		textLines = append(textLines, remapText)
	}

	text := header
	sort.Strings(textLines)
	text += strings.Join(textLines, "")
	return text
}

// BuildRemapLine builds the remap line for the given server and delivery service.
// The cacheKeyConfigParams map may be nil, if this ds profile had no cache key config params.
func BuildRemapLine(cacheURLConfigParams map[string]string, atsMajorVersion int, server *ServerInfo, pData map[string]string, text string, ds RemapConfigDSData, mapFrom string, mapTo string, cacheKeyConfigParams map[string]string) string {
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

	// multiple uses of cacheurl and cachekey plugins don't work right in ATS, but Perl has always done it.
	// So for now, keep track of it, so we can log an error when it happens.
	hasCacheURL := false
	hasCacheKey := false

	if ds.QStringIgnore != nil {
		if *ds.QStringIgnore == tc.QueryStringIgnoreDropAtEdge {
			dqsFile := "drop_qstring.config"
			text += ` @plugin=regex_remap.so @pparam=` + dqsFile
		} else if *ds.QStringIgnore == tc.QueryStringIgnoreIgnoreInCacheKeyAndPassUp {
			if _, globalExists := cacheURLConfigParams["location"]; globalExists {
				log.Warnln("qstring_ignore == 1, but global cacheurl.config param exists, so skipping remap rename config_file=cacheurl.config parameter")
			} else {
				qstr, addedCacheURL, addedCacheKey := GetQStringIgnoreRemap(atsMajorVersion)
				if addedCacheURL {
					hasCacheURL = true
				}
				if addedCacheKey {
					hasCacheKey = true
				}
				text += qstr
			}
		}
	}

	if ds.CacheURL != nil && *ds.CacheURL != "" {
		if hasCacheURL {
			log.Errorln("DeliveryService qstring_ignore and cacheurl both add cacheurl, but ATS cacheurl doesn't work correctly with multiple entries! Adding anyway!")
		}
		text += ` @plugin=cacheurl.so @pparam=` + CacheURLConfigFileName(ds.Name)
	}

	if len(cacheKeyConfigParams) > 0 {
		if hasCacheKey {
			log.Errorln("DeliveryService qstring_ignore and params both add cachekey, but ATS cachekey doesn't work correctly with multiple entries! Adding anyway!")
		}
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
	return text
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

type RemapLine struct {
	From string
	To   string
}

// MakeEdgeDSDataRemapLines returns the remap lines for the given server and delivery service.
// Returns nil, if the given server and ds have no remap lines, i.e. the DS match is not a host regex, or has no origin FQDN.
func MakeEdgeDSDataRemapLines(ds RemapConfigDSData, server *ServerInfo) ([]RemapLine, error) {
	if ds.RegexType == nil || tc.DSMatchType(*ds.RegexType) != tc.DSMatchTypeHostRegex || ds.OriginFQDN == nil || *ds.OriginFQDN == "" {
		return nil, nil
	}
	if ds.Pattern == nil {
		return nil, errors.New("ds missing regex pattern")
	}
	if ds.Protocol == nil {
		return nil, errors.New("ds missing protocol")
	}
	if ds.Domain == nil {
		return nil, errors.New("ds missing domain")
	}

	remapLines := []RemapLine{}
	hostRegex := *ds.Pattern
	mapTo := *ds.OriginFQDN + "/"

	mapFromHTTP := "http://" + hostRegex + "/"
	mapFromHTTPS := "https://" + hostRegex + "/"
	if strings.HasSuffix(hostRegex, `.*`) {
		re := hostRegex
		re = strings.Replace(re, `\`, ``, -1)
		re = strings.Replace(re, `.*`, ``, -1)

		hName := "__http__"
		if ds.Type.IsDNS() {
			if ds.RoutingName == nil {
				return nil, errors.New("ds is dns, but missing routing name")
			}
			hName = *ds.RoutingName
		}

		portStr := ""
		if hName == "__http__" && server.Port > 0 && server.Port != 80 {
			portStr = ":" + strconv.Itoa(server.Port)
		}

		httpsPortStr := ""
		if hName == "__http__" && server.HTTPSPort > 0 && server.HTTPSPort != 443 {
			httpsPortStr = ":" + strconv.Itoa(server.HTTPSPort)
		}

		mapFromHTTP = "http://" + hName + re + *ds.Domain + portStr + "/"
		mapFromHTTPS = "https://" + hName + re + *ds.Domain + httpsPortStr + "/"
	}

	if *ds.Protocol == tc.DSProtocolHTTP || *ds.Protocol == tc.DSProtocolHTTPAndHTTPS {
		remapLines = append(remapLines, RemapLine{From: mapFromHTTP, To: mapTo})
	}
	if *ds.Protocol == tc.DSProtocolHTTPS || *ds.Protocol == tc.DSProtocolHTTPToHTTPS || *ds.Protocol == tc.DSProtocolHTTPAndHTTPS {
		remapLines = append(remapLines, RemapLine{From: mapFromHTTPS, To: mapTo})
	}

	return remapLines, nil
}

func EdgeHeaderRewriteConfigFileName(dsName string) string {
	return "hdr_rw_" + dsName + ".config"
}

func MidHeaderRewriteConfigFileName(dsName string) string {
	return "hdr_rw_mid_" + dsName + ".config"
}

func CacheURLConfigFileName(dsName string) string {
	return "cacheurl_" + dsName + ".config"
}

// GetQStringIgnoreRemap returns the remap, whether cacheurl was added, and whether cachekey was added.
func GetQStringIgnoreRemap(atsMajorVersion int) (string, bool, bool) {
	if atsMajorVersion >= 6 {
		addingCacheURL := false
		addingCacheKey := true
		return ` @plugin=cachekey.so @pparam=--separator= @pparam=--remove-all-params=true @pparam=--remove-path=true @pparam=--capture-prefix-uri=/^([^?]*)/$1/`, addingCacheURL, addingCacheKey
	} else {
		addingCacheURL := true
		addingCacheKey := false
		return ` @plugin=cacheurl.so @pparam=cacheurl_qstring.config`, addingCacheURL, addingCacheKey
	}
}
