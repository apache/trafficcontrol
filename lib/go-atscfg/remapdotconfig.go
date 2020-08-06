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
	"github.com/apache/trafficcontrol/lib/go-util"
)

const CacheURLParameterConfigFile = "cacheurl.config"
const CacheKeyParameterConfigFile = "cachekey.config"
const ContentTypeRemapDotConfig = ContentTypeTextASCII
const LineCommentRemapDotConfig = LineCommentHash

func MakeRemapDotConfig(
	server *tc.Server,
	unfilteredDSes []tc.DeliveryServiceNullable,
	dss []tc.DeliveryServiceServer,
	dsRegexArr []tc.DeliveryServiceRegexes,
	serverParams []tc.Parameter,
	cdn *tc.CDN,
	toToolName string, // tm.toolname global parameter (TODO: cache itself?)
	toURL string, // tm.url global parameter (TODO: cache itself?)
	cacheKeyParams []tc.Parameter,
	topologies []tc.Topology,
	cacheGroupArr []tc.CacheGroupNullable,
	serverCapabilities map[int]map[ServerCapability]struct{},
	dsRequiredCapabilities map[int]map[ServerCapability]struct{},
) string {
	cdnDomain := cdn.DomainName
	dsRegexes := makeDSRegexMap(dsRegexArr)
	// Returned DSes are guaranteed to have a non-nil XMLID, Type, DSCP, ID, and Active.
	dses := remapFilterDSes(server, dss, unfilteredDSes, cacheKeyParams)

	dsProfilesCacheKeyConfigParams, err := makeDSProfilesCacheKeyConfigParams(server, dses, cacheKeyParams)
	if err != nil {
		log.Errorln("Error making Delivery Service Cache Key Params, cache key will be missing! : " + err.Error())
	}

	atsMajorVersion := getATSMajorVersion(serverParams)
	serverPackageParamData := makeServerPackageParamData(server, serverParams)
	cacheURLConfigParams := ParamsToMap(FilterParams(serverParams, CacheURLParameterConfigFile, "", "", ""))
	cacheGroups := MakeCGMap(cacheGroupArr)
	nameTopologies := MakeTopologyNameMap(topologies)

	hdr := GenericHeaderComment(server.HostName, toToolName, toURL)
	if tc.CacheTypeFromString(server.Type) == tc.CacheTypeMid {
		return GetServerConfigRemapDotConfigForMid(atsMajorVersion, dsProfilesCacheKeyConfigParams, dses, dsRegexes, hdr, server, nameTopologies, cacheGroups, serverCapabilities, dsRequiredCapabilities)
	}
	return GetServerConfigRemapDotConfigForEdge(cacheURLConfigParams, dsProfilesCacheKeyConfigParams, serverPackageParamData, dses, dsRegexes, atsMajorVersion, hdr, server, nameTopologies, cacheGroups, serverCapabilities, dsRequiredCapabilities, cdnDomain)
}

func GetServerConfigRemapDotConfigForMid(
	atsMajorVersion int,
	profilesCacheKeyConfigParams map[int]map[string]string,
	dses []tc.DeliveryServiceNullable,
	dsRegexes map[tc.DeliveryServiceName][]tc.DeliveryServiceRegex,
	header string,
	server *tc.Server,
	nameTopologies map[TopologyName]tc.Topology,
	cacheGroups map[tc.CacheGroupName]tc.CacheGroupNullable,
	serverCapabilities map[int]map[ServerCapability]struct{},
	dsRequiredCapabilities map[int]map[ServerCapability]struct{},
) string {
	midRemaps := map[string]string{}
	for _, ds := range dses {
		if !HasRequiredCapabilities(serverCapabilities[server.ID], dsRequiredCapabilities[*ds.ID]) {
			continue
		}

		topology, hasTopology := nameTopologies[TopologyName(*ds.Topology)]
		if *ds.Topology != "" && hasTopology && !topologyIncludesServer(topology, server) {
			continue
		}
		if ds.Type.IsLive() && !ds.Type.IsNational() && !hasTopology {
			continue // Live local delivery services skip mids (except Topologies ignore DS types)
		}

		if ds.OrgServerFQDN == nil || *ds.OrgServerFQDN == "" {
			log.Warnf("GetServerConfigRemapDotConfigForMid ds '" + *ds.XMLID + "' has no origin fqdn, skipping!") // TODO confirm - Perl uses without checking!
			continue
		}

		if midRemaps[*ds.OrgServerFQDN] != "" {
			continue // skip remap rules from extra HOST_REGEXP entries
		}

		// multiple uses of cacheurl and cachekey plugins don't work right in ATS, but Perl has always done it.
		// So for now, keep track of it, so we can log an error when it happens.
		hasCacheURL := false
		hasCacheKey := false

		midRemap := ""

		if *ds.Topology != "" {
			midRemap += MakeDSTopologyHeaderRewriteTxt(ds, tc.CacheGroupName(server.Cachegroup), topology, cacheGroups)
		} else if ds.MidHeaderRewrite != nil && *ds.MidHeaderRewrite != "" {
			midRemap += ` @plugin=header_rewrite.so @pparam=` + MidHeaderRewriteConfigFileName(*ds.XMLID)
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
				log.Errorln("Making remap.config for Delivery Service '" + *ds.XMLID + "': qstring_ignore and cacheurl both add cacheurl, but ATS cacheurl doesn't work correctly with multiple entries! Adding anyway!")
			}
			midRemap += ` @plugin=cacheurl.so @pparam=` + CacheURLConfigFileName(*ds.XMLID)
		}

		if ds.ProfileID != nil && len(profilesCacheKeyConfigParams[*ds.ProfileID]) > 0 {
			if hasCacheKey {
				log.Errorln("Making remap.config for Delivery Service '" + *ds.XMLID + "': qstring_ignore and cachekey params both add cachekey, but ATS cachekey doesn't work correctly with multiple entries! Adding anyway!")
			}
			midRemap += ` @plugin=cachekey.so`
			for name, val := range profilesCacheKeyConfigParams[*ds.ProfileID] {
				midRemap += ` @pparam=--` + name + "=" + val
			}
		}
		if ds.RangeRequestHandling != nil && (*ds.RangeRequestHandling == tc.RangeRequestHandlingCacheRangeRequest || *ds.RangeRequestHandling == tc.RangeRequestHandlingSlice) {
			midRemap += ` @plugin=cache_range_requests.so`
		}

		if midRemap != "" {
			midRemaps[*ds.OrgServerFQDN] = midRemap
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
	dses []tc.DeliveryServiceNullable,
	dsRegexes map[tc.DeliveryServiceName][]tc.DeliveryServiceRegex,
	atsMajorVersion int,
	header string,
	server *tc.Server,
	nameTopologies map[TopologyName]tc.Topology,
	cacheGroups map[tc.CacheGroupName]tc.CacheGroupNullable,
	serverCapabilities map[int]map[ServerCapability]struct{},
	dsRequiredCapabilities map[int]map[ServerCapability]struct{},
	cdnDomain string,
) string {
	textLines := []string{}

	for _, ds := range dses {
		for _, dsRegex := range dsRegexes[tc.DeliveryServiceName(*ds.XMLID)] {
			if !HasRequiredCapabilities(serverCapabilities[server.ID], dsRequiredCapabilities[*ds.ID]) {
				continue
			}

			topology, hasTopology := nameTopologies[TopologyName(*ds.Topology)]
			if *ds.Topology != "" && hasTopology && !topologyIncludesServer(topology, server) {
				continue
			}
			remapText := ""
			if *ds.Type == tc.DSTypeAnyMap {
				if ds.RemapText == nil {
					log.Errorln("ds '" + *ds.XMLID + "' is ANY_MAP, but has no remap text - skipping")
					continue
				}
				remapText = *ds.RemapText + "\n"
				textLines = append(textLines, remapText)
				continue
			}

			remapLines, err := MakeEdgeDSDataRemapLines(ds, dsRegex, server, cdnDomain)
			if err != nil {
				log.Errorln("making remap lines for DS '" + *ds.XMLID + "' - skipping! : " + err.Error())
				continue
			}

			for _, line := range remapLines {
				profilecacheKeyConfigParams := (map[string]string)(nil)
				if ds.ProfileID != nil {
					profilecacheKeyConfigParams = profilesCacheKeyConfigParams[*ds.ProfileID]
				}
				remapText = BuildEdgeRemapLine(cacheURLConfigParams, atsMajorVersion, server, serverPackageParamData, remapText, ds, dsRegex, line.From, line.To, profilecacheKeyConfigParams, cacheGroups, nameTopologies)
				if hasTopology {
					remapText += " # topology '" + topology.Name + "'"
				}
				remapText += "\n"
			}
			textLines = append(textLines, remapText)
		}
	}

	text := header
	sort.Strings(textLines)
	text += strings.Join(textLines, "")
	return text
}

const RemapConfigRangeDirective = `__RANGE_DIRECTIVE__`

// BuildEdgeRemapLine builds the remap line for the given server and delivery service.
// The cacheKeyConfigParams map may be nil, if this ds profile had no cache key config params.
func BuildEdgeRemapLine(
	cacheURLConfigParams map[string]string,
	atsMajorVersion int,
	server *tc.Server,
	pData map[string]string,
	text string,
	ds tc.DeliveryServiceNullable,
	dsRegex tc.DeliveryServiceRegex,
	mapFrom string,
	mapTo string,
	cacheKeyConfigParams map[string]string,
	cacheGroups map[tc.CacheGroupName]tc.CacheGroupNullable,
	nameTopologies map[TopologyName]tc.Topology,
) string {
	// ds = 'remap' in perl
	mapFrom = strings.Replace(mapFrom, `__http__`, server.HostName, -1)

	if _, hasDSCPRemap := pData["dscp_remap"]; hasDSCPRemap {
		text += "map	" + mapFrom + "     " + mapTo + ` @plugin=dscp_remap.so @pparam=` + strconv.Itoa(*ds.DSCP)
	} else {
		text += "map	" + mapFrom + "     " + mapTo + ` @plugin=header_rewrite.so @pparam=dscp/set_dscp_` + strconv.Itoa(*ds.DSCP) + ".config"
	}

	if *ds.Topology != "" {
		text += MakeDSTopologyHeaderRewriteTxt(ds, tc.CacheGroupName(server.Cachegroup), nameTopologies[TopologyName(*ds.Topology)], cacheGroups)
	} else if ds.EdgeHeaderRewrite != nil && *ds.EdgeHeaderRewrite != "" {
		text += ` @plugin=header_rewrite.so @pparam=` + EdgeHeaderRewriteConfigFileName(*ds.XMLID)
	}

	if ds.SigningAlgorithm != nil && *ds.SigningAlgorithm != "" {
		if *ds.SigningAlgorithm == tc.SigningAlgorithmURLSig {
			text += ` @plugin=url_sig.so @pparam=url_sig_` + *ds.XMLID + ".config"
		} else if *ds.SigningAlgorithm == tc.SigningAlgorithmURISigning {
			text += ` @plugin=uri_signing.so @pparam=uri_signing_` + *ds.XMLID + ".config"
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
				log.Warnln("Making remap.config for Delivery Service '" + *ds.XMLID + "': qstring_ignore == 1, but global cacheurl.config param exists, so skipping remap rename config_file=cacheurl.config parameter")
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
			log.Errorln("Making remap.config for Delivery Service '" + *ds.XMLID + "': qstring_ignore and cacheurl both add cacheurl, but ATS cacheurl doesn't work correctly with multiple entries! Adding anyway!")
		}
		text += ` @plugin=cacheurl.so @pparam=` + CacheURLConfigFileName(*ds.XMLID)
	}

	if len(cacheKeyConfigParams) > 0 {
		if hasCacheKey {
			log.Errorln("Making remap.config for Delivery Service '" + *ds.XMLID + "': qstring_ignore and params both add cachekey, but ATS cachekey doesn't work correctly with multiple entries! Adding anyway!")
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
		text += ` @plugin=regex_remap.so @pparam=regex_remap_` + *ds.XMLID + ".config"
	}

	rangeReqTxt := ""
	if ds.RangeRequestHandling != nil {
		if *ds.RangeRequestHandling == tc.RangeRequestHandlingBackgroundFetch {
			rangeReqTxt = ` @plugin=background_fetch.so @pparam=bg_fetch.config`
		} else if *ds.RangeRequestHandling == tc.RangeRequestHandlingCacheRangeRequest {
			rangeReqTxt = ` @plugin=cache_range_requests.so `
		} else if *ds.RangeRequestHandling == tc.RangeRequestHandlingSlice && ds.RangeSliceBlockSize != nil {
			rangeReqTxt = ` @plugin=slice.so @pparam=--blockbytes=` + strconv.Itoa(*ds.RangeSliceBlockSize) + ` @plugin=cache_range_requests.so	`
		}
	}

	remapText := ""
	if ds.RemapText != nil {
		remapText = *ds.RemapText
	}

	if strings.Contains(remapText, RemapConfigRangeDirective) {
		remapText = strings.Replace(remapText, RemapConfigRangeDirective, rangeReqTxt, 1)
	} else {
		text += rangeReqTxt
	}

	if remapText != "" {
		text += " " + remapText
	}

	if ds.FQPacingRate != nil && *ds.FQPacingRate > 0 {
		text += ` @plugin=fq_pacing.so @pparam=--rate=` + strconv.Itoa(*ds.FQPacingRate)
	}
	return text
}

// MakeDSTopologyHeaderRewriteTxt returns the appropriate header rewrite remap line text for the given DS on the given server.
// May be empty, if the DS has no header rewrite for the server's position in the topology.
func MakeDSTopologyHeaderRewriteTxt(ds tc.DeliveryServiceNullable, cg tc.CacheGroupName, topology tc.Topology, cacheGroups map[tc.CacheGroupName]tc.CacheGroupNullable) string {
	placement := getTopologyPlacement(cg, topology, cacheGroups)
	log.Errorf("DEBUG topo MakeDSTopologyHeaderRewriteTxt calling getTopologyPlacement cg '"+string(cg)+"' placement %+v\n", placement)
	const pluginTxt = ` @plugin=header_rewrite.so @pparam=`
	switch placement.CacheTier {
	case TopologyCacheTierFirst:
		if *ds.FirstHeaderRewrite == "" {
			return ""
		}
		log.Errorf("DEBUG topo MakeDSTopologyHeaderRewriteTxt first returning '" + pluginTxt + FirstHeaderRewriteConfigFileName(*ds.XMLID) + "'")
		return pluginTxt + FirstHeaderRewriteConfigFileName(*ds.XMLID)
	case TopologyCacheTierInner:
		if *ds.InnerHeaderRewrite == "" {
			return ""
		}
		return pluginTxt + InnerHeaderRewriteConfigFileName(*ds.XMLID)
	case TopologyCacheTierLast:
		if *ds.LastHeaderRewrite == "" {
			return ""
		}
		return pluginTxt + LastHeaderRewriteConfigFileName(*ds.XMLID)
	default:
		log.Errorln("Making topology header rewrite text: got unknown cache tier '" + placement.CacheTier + "' - not setting!")
		return ""
	}
}

type RemapLine struct {
	From string
	To   string
}

// MakeEdgeDSDataRemapLines returns the remap lines for the given server and delivery service.
// Returns nil, if the given server and ds have no remap lines, i.e. the DS match is not a host regex, or has no origin FQDN.
func MakeEdgeDSDataRemapLines(ds tc.DeliveryServiceNullable, dsRegex tc.DeliveryServiceRegex, server *tc.Server, cdnDomain string) ([]RemapLine, error) {
	if tc.DSMatchType(dsRegex.Type) != tc.DSMatchTypeHostRegex || ds.OrgServerFQDN == nil || *ds.OrgServerFQDN == "" {
		return nil, nil
	}
	if dsRegex.Pattern == "" {
		return nil, errors.New("ds missing regex pattern")
	}
	if ds.Protocol == nil {
		return nil, errors.New("ds missing protocol")
	}
	if cdnDomain == "" {
		return nil, errors.New("ds missing domain")
	}

	remapLines := []RemapLine{}
	hostRegex := dsRegex.Pattern
	mapTo := *ds.OrgServerFQDN + "/"

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
		if hName == "__http__" && server.TCPPort > 0 && server.TCPPort != 80 {
			portStr = ":" + strconv.Itoa(server.TCPPort)
		}

		httpsPortStr := ""
		if hName == "__http__" && server.HTTPSPort > 0 && server.HTTPSPort != 443 {
			httpsPortStr = ":" + strconv.Itoa(server.HTTPSPort)
		}

		mapFromHTTP = "http://" + hName + re + cdnDomain + portStr + "/"
		mapFromHTTPS = "https://" + hName + re + cdnDomain + httpsPortStr + "/"
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

// makeServerPackageParamData returns a map[paramName]paramVal for this server, config file 'package'.
func makeServerPackageParamData(server *tc.Server, serverParams []tc.Parameter) map[string]string {
	serverPackageParamData := map[string]string{}
	for _, param := range serverParams {
		if param.ConfigFile != "package" { // TODO put in const
			continue
		}
		if param.Name == "location" { // TODO put in const
			continue
		}

		paramName := param.Name
		// some files have multiple lines with the same key... handle that with param id.
		if _, ok := serverPackageParamData[param.Name]; ok {
			paramName += "__" + strconv.Itoa(param.ID)
		}
		paramValue := param.Value
		if paramValue == "STRING __HOSTNAME__" {
			paramValue = server.HostName + "." + server.DomainName // TODO strings.Replace to replace all anywhere, instead of just an exact match?
		}

		if val, ok := serverPackageParamData[paramName]; ok {
			if val < paramValue {
				log.Errorln("remap config generation got multiple parameters for server package name '" + paramName + "' - ignoring '" + paramValue + "'")
				continue
			} else {
				log.Errorln("config generation got multiple parameters for server package name '" + paramName + "' - ignoring '" + val + "'")
			}
		}
		serverPackageParamData[paramName] = paramValue
	}
	return serverPackageParamData
}

// remapFilterDSes filters Delivery Services to be used to generate remap.config for the given server.
// Returned DSes are guaranteed to have a non-nil XMLID, Type, DSCP, ID, Active, and Topology.
// If a DS has a nil Topology, OrgServerFQDN, FirstHeaderRewrite, InnerHeaderRewrite, or LastHeaderRewrite, "" is assigned.
func remapFilterDSes(server *tc.Server, dss []tc.DeliveryServiceServer, dses []tc.DeliveryServiceNullable, cacheKeyParams []tc.Parameter) []tc.DeliveryServiceNullable {
	isMid := strings.HasPrefix(server.Type, string(tc.CacheTypeMid))

	serverIDs := map[int]struct{}{}
	if !isMid {
		// mids use all servers, so pass empty=all. Edges only use this current server
		serverIDs[server.ID] = struct{}{}
	}

	dsIDs := map[int]struct{}{}
	for _, ds := range dses {
		if ds.ID == nil {
			// TODO log error?
			continue
		}
		dsIDs[*ds.ID] = struct{}{}
	}

	dsServers := FilterDSS(dss, dsIDs, serverIDs)

	dssMap := map[int]map[int]struct{}{} // set of map[dsID][serverID]
	for _, dss := range dsServers {
		if dss.Server == nil || dss.DeliveryService == nil {
			continue // TODO log?
		}
		if dssMap[*dss.DeliveryService] == nil {
			dssMap[*dss.DeliveryService] = map[int]struct{}{}
		}
		dssMap[*dss.DeliveryService][*dss.Server] = struct{}{}
	}

	useInactive := false
	if !isMid {
		// mids get inactive DSes, edges don't. This is how it's always behaved, not necessarily how it should.
		useInactive = true
	}

	filteredDSes := []tc.DeliveryServiceNullable{}
	for _, ds := range dses {
		if ds.Topology == nil {
			ds.Topology = util.StrPtr("")
		}
		if ds.OrgServerFQDN == nil {
			ds.OrgServerFQDN = util.StrPtr("")
		}
		if ds.FirstHeaderRewrite == nil {
			ds.FirstHeaderRewrite = util.StrPtr("")
		}
		if ds.InnerHeaderRewrite == nil {
			ds.InnerHeaderRewrite = util.StrPtr("")
		}
		if ds.LastHeaderRewrite == nil {
			ds.LastHeaderRewrite = util.StrPtr("")
		}
		if ds.XMLID == nil {
			log.Errorln("Remap config gen got Delivery Service with nil XMLID, skipping!")
			continue
		} else if ds.Type == nil {
			log.Errorln("Remap config gen got Delivery Service '" + *ds.XMLID + "'  with nil Type, skipping!")
			continue
		} else if ds.DSCP == nil {
			log.Errorln("Remap config gen got Delivery Service '" + *ds.XMLID + "'  with nil DSCP, skipping!")
			continue
		} else if ds.ID == nil {
			log.Errorln("Remap config gen got Delivery Service '" + *ds.XMLID + "'  with nil ID, skipping!")
			continue
		} else if ds.Active == nil {
			log.Errorln("Remap config gen got Delivery Service '" + *ds.XMLID + "'  with nil Active, skipping!")
			continue
		} else if _, ok := dssMap[*ds.ID]; !ok && *ds.Topology == "" {
			continue // normal, not an error, this DS just isn't assigned to our Cache
		} else if !useInactive && !*ds.Active {
			continue // normal, not an error, DS just isn't active and we aren't including inactive DSes
		}
		filteredDSes = append(filteredDSes, ds)
	}
	return filteredDSes
}

// getATSMajorVersion returns the ATS major version from the config_file 'package' name 'trafficserver' Parameter on the given Server Profile Parameters.
// If no Parameter is found, or the value is malformed, a warning or error is logged and DefaultATSVersion is returned.
func getATSMajorVersion(serverParams []tc.Parameter) int {
	atsVersionParam := ""
	for _, param := range serverParams {
		if param.ConfigFile != "package" || param.Name != "trafficserver" {
			continue
		}
		atsVersionParam = param.Value
		break
	}
	if atsVersionParam == "" {
		log.Warnln("ATS version Parameter (config_file 'package' name 'trafficserver') not found on Server Profile, using default")
		atsVersionParam = DefaultATSVersion
	}

	atsMajorVer, err := GetATSMajorVersionFromATSVersion(atsVersionParam)
	if err != nil {
		log.Errorln("getting ATS major version from server Profile Parameter, using default: " + err.Error())
		atsMajorVer, err = GetATSMajorVersionFromATSVersion(DefaultATSVersion)
		if err != nil {
			// should never happen
			log.Errorln("getting ATS major version from default version! Should never happen! Using 0, config will be malformed! : " + err.Error())
		}
	}
	return atsMajorVer
}

// makeDSProfilesCacheKeyConfigParams returns a map[ProfileID][ParamName]ParamValue for the cache key params for each profile.
func makeDSProfilesCacheKeyConfigParams(server *tc.Server, dses []tc.DeliveryServiceNullable, cacheKeyParams []tc.Parameter) (map[int]map[string]string, error) {
	cacheKeyParamsWithProfiles, err := TCParamsToParamsWithProfiles(cacheKeyParams)
	if err != nil {
		return nil, errors.New("decoding cache key parameter profiles: " + err.Error())
	}

	cacheKeyParamsWithProfilesMap := ParameterWithProfilesToMap(cacheKeyParamsWithProfiles)

	dsProfileNamesToIDs := map[string]int{}
	for _, ds := range dses {
		if ds.ProfileID == nil || ds.ProfileName == nil {
			continue // TODO log
		}
		dsProfileNamesToIDs[*ds.ProfileName] = *ds.ProfileID
	}

	dsProfilesCacheKeyConfigParams := map[int]map[string]string{}
	for _, param := range cacheKeyParamsWithProfilesMap {
		for dsProfileName, dsProfileID := range dsProfileNamesToIDs {
			if _, ok := param.ProfileNames[dsProfileName]; ok {
				if _, ok := dsProfilesCacheKeyConfigParams[dsProfileID]; !ok {
					dsProfilesCacheKeyConfigParams[dsProfileID] = map[string]string{}
				}
				if val, ok := dsProfilesCacheKeyConfigParams[dsProfileID][param.Name]; ok {
					if val < param.Value {
						log.Errorln("remap config generation got multiple parameters for name '" + param.Name + "' - ignoring '" + param.Value + "'")
						continue
					} else {
						log.Errorln("remap config generation got multiple parameters for name '" + param.Name + "' - ignoring '" + val + "'")
					}
				}
				dsProfilesCacheKeyConfigParams[dsProfileID][param.Name] = param.Value
			}
		}
	}
	return dsProfilesCacheKeyConfigParams, nil
}

type DeliveryServiceRegexesSortByTypeThenSetNum []tc.DeliveryServiceRegex

func (r DeliveryServiceRegexesSortByTypeThenSetNum) Len() int { return len(r) }
func (r DeliveryServiceRegexesSortByTypeThenSetNum) Less(i, j int) bool {
	if rc := strings.Compare(r[i].Type, r[j].Type); rc != 0 {
		return rc < 0
	}
	return r[i].SetNumber < r[j].SetNumber
}
func (r DeliveryServiceRegexesSortByTypeThenSetNum) Swap(i, j int) { r[i], r[j] = r[j], r[i] }

func makeDSRegexMap(regexes []tc.DeliveryServiceRegexes) map[tc.DeliveryServiceName][]tc.DeliveryServiceRegex {
	dsRegexMap := map[tc.DeliveryServiceName][]tc.DeliveryServiceRegex{}
	for _, dsRegex := range regexes {
		sort.Sort(DeliveryServiceRegexesSortByTypeThenSetNum(dsRegex.Regexes))
		dsRegexMap[tc.DeliveryServiceName(dsRegex.DSName)] = dsRegex.Regexes
	}
	return dsRegexMap
}
