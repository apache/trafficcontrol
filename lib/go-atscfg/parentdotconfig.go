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
	"bytes"
	"errors"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
)

const ContentTypeParentDotConfig = ContentTypeTextASCII
const LineCommentParentDotConfig = LineCommentHash

const ParentConfigParamQStringHandling = "psel.qstring_handling"
const ParentConfigParamMSOAlgorithm = "mso.algorithm"
const ParentConfigParamMSOParentRetry = "mso.parent_retry"
const ParentConfigParamMSOUnavailableServerRetryResponses = "mso.unavailable_server_retry_responses"
const ParentConfigParamMSOMaxSimpleRetries = "mso.max_simple_retries"
const ParentConfigParamMSOMaxUnavailableServerRetries = "mso.max_unavailable_server_retries"
const ParentConfigParamAlgorithm = "algorithm"
const ParentConfigParamQString = "qstring"
const ParentConfigParamSecondaryMode = "try_all_primaries_before_secondary"

const ParentConfigParamParentRetry = "parent_retry"
const ParentConfigParamUnavailableServerRetryResponses = "unavailable_server_retry_responses"
const ParentConfigParamMaxSimpleRetries = "max_simple_retries"
const ParentConfigParamMaxUnavailableServerRetries = "max_unavailable_server_retries"

const ParentConfigDSParamDefaultMSOAlgorithm = "consistent_hash"
const ParentConfigDSParamDefaultMSOParentRetry = "both"
const ParentConfigDSParamDefaultMSOUnavailableServerRetryResponses = ""
const ParentConfigDSParamDefaultMaxSimpleRetries = "1"
const ParentConfigDSParamDefaultMaxUnavailableServerRetries = "1"

const ParentConfigCacheParamWeight = "weight"
const ParentConfigCacheParamPort = "port"
const ParentConfigCacheParamUseIP = "use_ip_address"
const ParentConfigCacheParamRank = "rank"
const ParentConfigCacheParamNotAParent = "not_a_parent"

// TODO change, this is terrible practice, using a hard-coded key. What if there were a delivery service named "all_parents" (transliterated Perl)
const DeliveryServicesAllParentsKey = "all_parents"

type ParentConfigDS struct {
	Name                 tc.DeliveryServiceName
	QStringIgnore        tc.QStringIgnore
	OriginFQDN           string
	MultiSiteOrigin      bool
	OriginShield         string
	Type                 tc.DSType
	QStringHandling      string
	RequiredCapabilities map[ServerCapability]struct{}
	Topology             string
}

type ParentConfigDSTopLevel struct {
	ParentConfigDS
	MSOAlgorithm                       string
	MSOParentRetry                     string
	MSOUnavailableServerRetryResponses string
	MSOMaxSimpleRetries                string
	MSOMaxUnavailableServerRetries     string
}

type ParentInfo struct {
	Host            string
	Port            int
	Domain          string
	Weight          string
	UseIP           bool
	Rank            int
	IP              string
	PrimaryParent   bool
	SecondaryParent bool
	Capabilities    map[ServerCapability]struct{}
}

func (p ParentInfo) Format() string {
	host := ""
	if p.UseIP {
		host = p.IP
	} else {
		host = p.Host + "." + p.Domain
	}
	return host + ":" + strconv.Itoa(p.Port) + "|" + p.Weight + ";"
}

type OriginHost string
type OriginFQDN string

type ParentInfos map[OriginHost]ParentInfo

type ParentInfoSortByRank []ParentInfo

func (s ParentInfoSortByRank) Len() int      { return len(([]ParentInfo)(s)) }
func (s ParentInfoSortByRank) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s ParentInfoSortByRank) Less(i, j int) bool {
	if s[i].Rank != s[j].Rank {
		return s[i].Rank < s[j].Rank
	} else if s[i].Host != s[j].Host {
		return s[i].Host < s[j].Host
	} else if s[i].Domain != s[j].Domain {
		return s[i].Domain < s[j].Domain
	} else if s[i].Port != s[j].Port {
		return s[i].Port < s[j].Port
	}
	return s[i].IP < s[j].IP
}

type ServerWithParams struct {
	tc.ServerNullable
	Params ProfileCache
}

type ServersWithParamsSortByRank []ServerWithParams

func (ss ServersWithParamsSortByRank) Len() int      { return len(ss) }
func (ss ServersWithParamsSortByRank) Swap(i, j int) { ss[i], ss[j] = ss[j], ss[i] }
func (ss ServersWithParamsSortByRank) Less(i, j int) bool {
	if ss[i].Params.Rank != ss[j].Params.Rank {
		return ss[i].Params.Rank < ss[j].Params.Rank
	}

	if ss[i].HostName == nil {
		if ss[j].HostName != nil {
			return true
		}
	} else if ss[j].HostName == nil {
		return false
	} else if ss[i].HostName != ss[j].HostName {
		return *ss[i].HostName < *ss[j].HostName
	}

	if ss[i].DomainName == nil {
		if ss[j].DomainName != nil {
			return true
		}
	} else if ss[j].DomainName == nil {
		return false
	} else if ss[i].DomainName != ss[j].DomainName {
		return *ss[i].DomainName < *ss[j].DomainName
	}

	if ss[i].Params.Port != ss[j].Params.Port {
		return ss[i].Params.Port < ss[j].Params.Port
	}

	iIP := GetServerIPAddress(&ss[i].ServerNullable)
	jIP := GetServerIPAddress(&ss[j].ServerNullable)

	if iIP == nil {
		if jIP != nil {
			return true
		}
	} else if jIP == nil {
		return false
	}
	return bytes.Compare(iIP, jIP) <= 0
}

type ParentConfigDSTopLevelSortByName []ParentConfigDSTopLevel

func (s ParentConfigDSTopLevelSortByName) Len() int      { return len(([]ParentConfigDSTopLevel)(s)) }
func (s ParentConfigDSTopLevelSortByName) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s ParentConfigDSTopLevelSortByName) Less(i, j int) bool {
	return strings.Compare(string(s[i].Name), string(s[j].Name)) < 0
}

type DSesSortByName []tc.DeliveryServiceNullableV30

func (s DSesSortByName) Len() int      { return len(s) }
func (s DSesSortByName) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s DSesSortByName) Less(i, j int) bool {
	if s[i].XMLID == nil {
		return true
	}
	if s[j].XMLID == nil {
		return false
	}
	return *s[i].XMLID < *s[j].XMLID
}

type ProfileCache struct {
	Weight     string
	Port       int
	UseIP      bool
	Rank       int
	NotAParent bool
}

func DefaultProfileCache() ProfileCache {
	return ProfileCache{
		Weight:     "0.999",
		Port:       0,
		UseIP:      false,
		Rank:       1,
		NotAParent: false,
	}
}

// CGServer is the server table data needed when selecting the servers assigned to a cachegroup.
type CGServer struct {
	ServerID       ServerID
	ServerHost     string
	ServerIP       string
	ServerPort     int
	CacheGroupID   int
	CacheGroupName string
	Status         int
	Type           int
	ProfileID      ProfileID
	ProfileName    string
	CDN            int
	TypeName       string
	Domain         string
	Capabilities   map[ServerCapability]struct{}
}

type OriginURI struct {
	Scheme string
	Host   string
	Port   string
}

func MakeParentDotConfig(
	toToolName string, // tm.toolname global parameter (TODO: cache itself?)
	toURL string, // tm.url global parameter (TODO: cache itself?)
	dses []tc.DeliveryServiceNullableV30,
	server *tc.ServerNullable,
	servers []tc.ServerNullable,
	topologies []tc.Topology,
	tcServerParams []tc.Parameter,
	tcParentConfigParams []tc.Parameter,
	serverCapabilities map[int]map[ServerCapability]struct{},
	dsRequiredCapabilities map[int]map[ServerCapability]struct{},
	cacheGroupArr []tc.CacheGroupNullable,
	dss []tc.DeliveryServiceServer,
	cdn *tc.CDN,
) string {
	if server.HostName == nil || *server.HostName == "" {
		return "ERROR: server HostName missing"
	} else if server.CDNName == nil || *server.CDNName == "" {
		return "ERROR: server CDNName missing"
	} else if server.Cachegroup == nil || *server.Cachegroup == "" {
		return "ERROR: server Cachegroup missing"
	} else if server.Profile == nil || *server.Profile == "" {
		return "ERROR: server Profile missing"
	} else if server.TCPPort == nil {
		return "ERROR: server TCPPort missing"
	}

	atsMajorVer := getATSMajorVersion(tcServerParams)

	cacheGroups, err := MakeCGMap(cacheGroupArr)
	if err != nil {
		return "ERROR: making CacheGroup map, config will be malformed! : " + err.Error()
	}
	serverParentCGData, err := GetParentCacheGroupData(server, cacheGroups)
	if err != nil {
		log.Errorln("making parent.config, getting server parent cachegroup data, config will be malformed! : " + err.Error())
	}
	isTopLevelCache := IsTopLevelCache(serverParentCGData)
	serverCDNDomain := cdn.DomainName

	sort.Sort(DSesSortByName(dses))

	nameVersionStr := GetNameVersionStringFromToolNameAndURL(toToolName, toURL)
	hdr := HeaderCommentWithTOVersionStr(*server.HostName, nameVersionStr)

	textArr := []string{}
	processedOriginsToDSNames := map[string]tc.DeliveryServiceName{}

	parentConfigParamsWithProfiles, err := TCParamsToParamsWithProfiles(tcParentConfigParams)
	if err != nil {
		log.Errorln("parent.config generation: error getting profiles from Traffic Ops Parameters, Parameters will not be considered for generation! : " + err.Error())
		parentConfigParamsWithProfiles = []ParameterWithProfiles{}
	}
	parentConfigParams := ParameterWithProfilesToMap(parentConfigParamsWithProfiles)

	// this is an optimization, to avoid looping over all params, for every DS. Instead, we loop over all params only once, and put them in a profile map.
	profileParentConfigParams := map[string]map[string]string{} // map[profileName][paramName]paramVal
	for _, param := range parentConfigParamsWithProfiles {
		for _, profile := range param.ProfileNames {
			if _, ok := profileParentConfigParams[profile]; !ok {
				profileParentConfigParams[profile] = map[string]string{}
			}
			profileParentConfigParams[profile][param.Name] = param.Value
		}
	}

	// We only need parent.config params, don't need all the params on the server
	serverParams := map[string]string{}
	if server.Profile == nil || *server.Profile != "" { // TODO warn/error if false? Servers requires profiles
		for name, val := range profileParentConfigParams[*server.Profile] {
			if name == ParentConfigParamQStringHandling ||
				name == ParentConfigParamAlgorithm ||
				name == ParentConfigParamQString {
				serverParams[name] = val
			}
		}
	}

	parentCacheGroups := map[string]struct{}{}
	if isTopLevelCache {
		log.Infoln("This cache Is Top Level!")
		for _, cg := range cacheGroups {
			if cg.Type == nil {
				return "ERROR: cachegroup type is nil!"
			}
			if cg.Name == nil {
				return "ERROR: cachegroup name is nil!"
			}

			if *cg.Type != tc.CacheGroupOriginTypeName {
				continue
			}
			parentCacheGroups[*cg.Name] = struct{}{}
		}
	} else {
		for _, cg := range cacheGroups {
			if cg.Type == nil {
				return "ERROR: cachegroup type is nil!"
			}
			if cg.Name == nil {
				return "ERROR: cachegroup type is nil!"
			}

			if *cg.Name == *server.Cachegroup {
				if cg.ParentName != nil && *cg.ParentName != "" {
					parentCacheGroups[*cg.ParentName] = struct{}{}
				}
				if cg.SecondaryParentName != nil && *cg.SecondaryParentName != "" {
					parentCacheGroups[*cg.SecondaryParentName] = struct{}{}
				}
				break
			}
		}
	}

	nameTopologies := MakeTopologyNameMap(topologies)

	cgServers := map[int]tc.ServerNullable{} // map[serverID]server
	for _, sv := range servers {
		if sv.ID == nil {
			log.Errorln("parent.config generation: TO servers had server with missing ID, skipping!")
			continue
		} else if sv.CDNName == nil {
			log.Errorln("parent.config generation: TO servers had server with missing CDNName, skipping!")
			continue
		} else if sv.Cachegroup == nil || *sv.Cachegroup == "" {
			log.Errorln("parent.config generation: TO servers had server with missing Cachegroup, skipping!")
			continue
		} else if sv.Status == nil || *sv.Status == "" {
			log.Errorln("parent.config generation: TO servers had server with missing Status, skipping!")
			continue
		} else if sv.Type == "" {
			log.Errorln("parent.config generation: TO servers had server with missing Type, skipping!")
			continue
		}
		if *sv.CDNName != *server.CDNName {
			continue
		}
		if _, ok := parentCacheGroups[*sv.Cachegroup]; !ok {
			continue
		}
		if sv.Type != tc.OriginTypeName &&
			!strings.HasPrefix(sv.Type, tc.EdgeTypePrefix) &&
			!strings.HasPrefix(sv.Type, tc.MidTypePrefix) {
			continue
		}
		if *sv.Status != string(tc.CacheStatusReported) && *sv.Status != string(tc.CacheStatusOnline) {
			continue
		}
		cgServers[*sv.ID] = sv
	}

	cgServerIDs := map[int]struct{}{}
	for serverID, _ := range cgServers {
		cgServerIDs[serverID] = struct{}{}
	}
	cgServerIDs[*server.ID] = struct{}{}

	cgDSServers := FilterDSS(dss, nil, cgServerIDs)
	parentServerDSes := map[int]map[int]struct{}{} // map[serverID][dsID]
	for _, dss := range cgDSServers {
		if dss.Server == nil || dss.DeliveryService == nil {
			return "ERROR: getting parent.config cachegroup parent server delivery service servers: got dss with nil members!"
		}
		if parentServerDSes[*dss.Server] == nil {
			parentServerDSes[*dss.Server] = map[int]struct{}{}
		}
		parentServerDSes[*dss.Server][*dss.DeliveryService] = struct{}{}
	}

	originServers, profileCaches, err := GetOriginServersAndProfileCaches(cgServers, parentServerDSes, profileParentConfigParams, dses, serverCapabilities, dsRequiredCapabilities)
	if err != nil {
		return "ERROR getting origin servers and profile caches: " + err.Error()
	}

	parentInfos := MakeParentInfo(serverParentCGData, serverCDNDomain, profileCaches, originServers)

	dsOrigins := makeDSOrigins(dss, dses, servers)

	for _, ds := range dses {
		if ds.XMLID == nil || *ds.XMLID == "" {
			log.Errorln("parent.config got ds with missing XMLID, skipping!")
			continue
		} else if ds.ID == nil {
			log.Errorln("parent.config got ds with missing ID, skipping!")
			continue
		} else if ds.Type == nil {
			log.Errorln("parent.config got ds with missing Type, skipping!")
			continue
		}

		if !isTopLevelCache && ds.Topology == nil {
			if _, ok := parentServerDSes[*server.ID][*ds.ID]; !ok {
				continue // skip DSes not assigned to this server.
			}
		}

		if !ds.Type.IsHTTP() && !ds.Type.IsDNS() {
			continue // skip ANY_MAP, STEERING, etc
		}
		if ds.OrgServerFQDN == nil || *ds.OrgServerFQDN == "" {
			// this check needs to be after the HTTP|DNS check, because Steering DSes without origins are ok
			log.Errorln("ds  '" + *ds.XMLID + "' has no origin server! Skipping!")
			continue
		}

		// Note these Parameters are only used for MSO for legacy DeliveryServiceServers DeliveryServices (except QueryStringHandling which is used by all DeliveryServices).
		//      Topology DSes use them for all DSes, MSO and non-MSO.
		dsParams := getParentDSParams(ds, profileParentConfigParams)

		log.Infoln("parent.config processing ds '" + *ds.XMLID + "'")

		if existingDS, ok := processedOriginsToDSNames[*ds.OrgServerFQDN]; ok {
			log.Errorln("parent.config generation: duplicate origin! services '" + *ds.XMLID + "' and '" + string(existingDS) + "' share origin '" + *ds.OrgServerFQDN + "': skipping '" + *ds.XMLID + "'!")
			continue
		}

		// TODO put these in separate functions. No if-statement should be this long.
		if ds.Topology != nil && *ds.Topology != "" {
			log.Infoln("parent.config generating Topology line for ds '" + *ds.XMLID + "'")
			txt, err := GetTopologyParentConfigLine(
				server,
				servers,
				&ds,
				serverParams,
				parentConfigParams,
				nameTopologies,
				serverCapabilities,
				dsRequiredCapabilities,
				cacheGroups,
				dsParams,
				atsMajorVer,
				dsOrigins[DeliveryServiceID(*ds.ID)],
			)
			if err != nil {
				log.Errorln(err)
				continue
			}

			if txt != "" { // will be empty with no error if this server isn't in the Topology, or if it doesn't have the Required Capabilities
				textArr = append(textArr, txt)
			}
		} else if IsTopLevelCache(serverParentCGData) {
			log.Infoln("parent.config generating top level line for ds '" + *ds.XMLID + "'")
			parentQStr := "ignore"
			if dsParams.QueryStringHandling == "" && dsParams.Algorithm == tc.AlgorithmConsistentHash && ds.QStringIgnore != nil && tc.QStringIgnore(*ds.QStringIgnore) == tc.QStringIgnoreUseInCacheKeyAndPassUp {
				parentQStr = "consider"
			}

			orgURI, err := GetOriginURI(*ds.OrgServerFQDN)
			if err != nil {
				log.Errorln("Malformed ds '" + *ds.XMLID + "' origin  URI: '" + *ds.OrgServerFQDN + "': skipping!" + err.Error())
				continue
			}

			textLine := ""

			if ds.OriginShield != nil && *ds.OriginShield != "" {
				algorithm := ""
				if parentSelectAlg := serverParams[ParentConfigParamAlgorithm]; strings.TrimSpace(parentSelectAlg) != "" {
					algorithm = "round_robin=" + parentSelectAlg
				}
				textLine += "dest_domain=" + orgURI.Hostname() + " port=" + orgURI.Port() + " parent=" + *ds.OriginShield + " " + algorithm + " go_direct=true\n"
			} else if ds.MultiSiteOrigin != nil && *ds.MultiSiteOrigin {
				textLine += "dest_domain=" + orgURI.Hostname() + " port=" + orgURI.Port() + " "
				if len(parentInfos) == 0 {
				}

				if len(parentInfos[OriginHost(orgURI.Hostname())]) == 0 {
					// TODO error? emulates Perl
					log.Warnln("ParentInfo: delivery service " + *ds.XMLID + " has no parent servers")
				}

				parents, secondaryParents := getMSOParentStrs(&ds, parentInfos[OriginHost(orgURI.Hostname())], atsMajorVer, dsRequiredCapabilities, dsParams.Algorithm, dsParams.TryAllPrimariesBeforeSecondary)
				textLine += parents + secondaryParents + ` round_robin=` + dsParams.Algorithm + ` qstring=` + parentQStr + ` go_direct=false parent_is_proxy=false`
				textLine += getParentRetryStr(true, atsMajorVer, dsParams.ParentRetry, dsParams.UnavailableServerRetryResponses, dsParams.MaxSimpleRetries, dsParams.MaxUnavailableServerRetries)
				textLine += "\n" // TODO remove, and join later on "\n" instead of ""?
				textArr = append(textArr, textLine)
			}
		} else {
			log.Infoln("parent.config generating non-top level line for ds '" + *ds.XMLID + "'")
			queryStringHandling := serverParams[ParentConfigParamQStringHandling] // "qsh" in Perl

			roundRobin := `round_robin=consistent_hash`
			goDirect := `go_direct=false`

			parents, secondaryParents := getParentStrs(&ds, dsRequiredCapabilities, parentInfos[DeliveryServicesAllParentsKey], atsMajorVer, dsParams.TryAllPrimariesBeforeSecondary)

			text := ""
			orgURI, err := GetOriginURI(*ds.OrgServerFQDN)
			if err != nil {
				log.Errorln("Malformed ds '" + *ds.XMLID + "' origin  URI: '" + *ds.OrgServerFQDN + "': skipping!" + err.Error())
				continue
			}

			// TODO encode this in a DSType func, IsGoDirect() ?
			if *ds.Type == tc.DSTypeHTTPNoCache || *ds.Type == tc.DSTypeHTTPLive || *ds.Type == tc.DSTypeDNSLive {
				text += `dest_domain=` + orgURI.Hostname() + ` port=` + orgURI.Port() + ` go_direct=true` + "\n"
			} else {

				// check for profile psel.qstring_handling.  If this parameter is assigned to the server profile,
				// then edges will use the qstring handling value specified in the parameter for all profiles.

				// If there is no defined parameter in the profile, then check the delivery service profile.
				// If psel.qstring_handling exists in the DS profile, then we use that value for the specified DS only.
				// This is used only if not overridden by a server profile qstring handling parameter.

				// TODO refactor this logic, hard to understand (transliterated from Perl)
				dsQSH := queryStringHandling
				if dsQSH == "" {
					dsQSH = dsParams.QueryStringHandling
				}
				parentQStr := dsQSH
				if parentQStr == "" {
					parentQStr = "ignore"
				}
				if ds.QStringIgnore != nil && tc.QStringIgnore(*ds.QStringIgnore) == tc.QStringIgnoreUseInCacheKeyAndPassUp && dsQSH == "" {
					parentQStr = "consider"
				}

				text += `dest_domain=` + orgURI.Hostname() + ` port=` + orgURI.Port() + ` ` + parents + ` ` + secondaryParents + ` ` + roundRobin + ` ` + goDirect + ` qstring=` + parentQStr + "\n"
			}
			textArr = append(textArr, text)
		}
		processedOriginsToDSNames[*ds.OrgServerFQDN] = tc.DeliveryServiceName(*ds.XMLID)
	}

	// TODO determine if this is necessary. It's super-dangerous, and moreover ignores Server Capabilitites.
	defaultDestText := ""
	if !IsTopLevelCache(serverParentCGData) {
		invalidDS := &tc.DeliveryServiceNullableV30{}
		invalidDS.ID = util.IntPtr(-1)
		tryAllPrimariesBeforeSecondary := false
		parents, secondaryParents := getParentStrs(invalidDS, dsRequiredCapabilities, parentInfos[DeliveryServicesAllParentsKey], atsMajorVer, tryAllPrimariesBeforeSecondary)
		defaultDestText = `dest_domain=. ` + parents
		if serverParams[ParentConfigParamAlgorithm] == tc.AlgorithmConsistentHash {
			defaultDestText += secondaryParents
		}
		defaultDestText += ` round_robin=consistent_hash go_direct=false`

		if qStr := serverParams[ParentConfigParamQString]; qStr != "" {
			defaultDestText += ` qstring=` + qStr
		}
		defaultDestText += "\n"
	}

	sort.Sort(sort.StringSlice(textArr))
	text := hdr + strings.Join(textArr, "") + defaultDestText
	return text
}

type ParentDSParams struct {
	Algorithm                       string
	ParentRetry                     string
	UnavailableServerRetryResponses string
	MaxSimpleRetries                string
	MaxUnavailableServerRetries     string
	QueryStringHandling             string
	TryAllPrimariesBeforeSecondary  bool
}

// getDSParams returns the Delivery Service Profile Parameters used in parent.config.
// If Parameters don't exist, defaults are returned. Non-MSO Delivery Services default to no custom retry logic (we should reevaluate that).
// Note these Parameters are only used for MSO for legacy DeliveryServiceServers DeliveryServices.
//      Topology DSes use them for all DSes, MSO and non-MSO.
func getParentDSParams(ds tc.DeliveryServiceNullableV30, profileParentConfigParams map[string]map[string]string) ParentDSParams {
	params := ParentDSParams{}
	isMSO := ds.MultiSiteOrigin != nil && *ds.MultiSiteOrigin
	if isMSO {
		params.Algorithm = ParentConfigDSParamDefaultMSOAlgorithm
		params.ParentRetry = ParentConfigDSParamDefaultMSOParentRetry
		params.UnavailableServerRetryResponses = ParentConfigDSParamDefaultMSOUnavailableServerRetryResponses
		params.MaxSimpleRetries = ParentConfigDSParamDefaultMaxSimpleRetries
		params.MaxUnavailableServerRetries = ParentConfigDSParamDefaultMaxUnavailableServerRetries
	}
	if ds.ProfileName == nil || *ds.ProfileName == "" {
		return params
	}
	dsParams, ok := profileParentConfigParams[*ds.ProfileName]
	if !ok {
		return params
	}

	params.QueryStringHandling = dsParams[ParentConfigParamQStringHandling] // may be blank, no default
	// TODO deprecate & remove "mso." Parameters - there was never a reason to restrict these settings to MSO.
	if isMSO {
		if v, ok := dsParams[ParentConfigParamMSOAlgorithm]; ok && strings.TrimSpace(v) != "" {
			params.Algorithm = v
		}
		if v, ok := dsParams[ParentConfigParamMSOParentRetry]; ok {
			params.ParentRetry = v
		}
		if v, ok := dsParams[ParentConfigParamMSOUnavailableServerRetryResponses]; ok {
			if v != "" && !unavailableServerRetryResponsesValid(v) {
				log.Errorln("Malformed " + ParentConfigParamMSOUnavailableServerRetryResponses + " parameter '" + v + "', not using!")
			} else if v != "" {
				params.UnavailableServerRetryResponses = v
			}
		}
		if v, ok := dsParams[ParentConfigParamMSOMaxSimpleRetries]; ok {
			params.MaxSimpleRetries = v
		}
		if v, ok := dsParams[ParentConfigParamMSOMaxUnavailableServerRetries]; ok {
			params.MaxUnavailableServerRetries = v
		}
	}

	// Even if the DS is MSO, non-"mso." Parameters override "mso." ones, because they're newer.
	if v, ok := dsParams[ParentConfigParamAlgorithm]; ok && strings.TrimSpace(v) != "" {
		params.Algorithm = v
	}
	if v, ok := dsParams[ParentConfigParamParentRetry]; ok {
		params.ParentRetry = v
	}
	if v, ok := dsParams[ParentConfigParamUnavailableServerRetryResponses]; ok {
		if v != "" && !unavailableServerRetryResponsesValid(v) {
			log.Errorln("Malformed " + ParentConfigParamUnavailableServerRetryResponses + " parameter '" + v + "', not using!")
		} else if v != "" {
			params.UnavailableServerRetryResponses = v
		}
	}
	if v, ok := dsParams[ParentConfigParamMaxSimpleRetries]; ok {
		params.MaxSimpleRetries = v
	}
	if v, ok := dsParams[ParentConfigParamMaxUnavailableServerRetries]; ok {
		params.MaxUnavailableServerRetries = v
	}
	if v, ok := dsParams[ParentConfigParamSecondaryMode]; ok {
		if v != "" {
			log.Errorln("parent.config generation: DS '" + *ds.XMLID + "' had Parameter " + ParentConfigParamSecondaryMode + " which is used if it exists, the value is ignored! Non-empty value '" + v + "' will be ignored!")
		}
		params.TryAllPrimariesBeforeSecondary = true
	}

	return params
}

func GetTopologyParentConfigLine(
	server *tc.ServerNullable,
	servers []tc.ServerNullable,
	ds *tc.DeliveryServiceNullableV30,
	serverParams map[string]string,
	parentConfigParams []ParameterWithProfilesMap, // all params with configFile parent.config
	nameTopologies map[TopologyName]tc.Topology,
	serverCapabilities map[int]map[ServerCapability]struct{},
	dsRequiredCapabilities map[int]map[ServerCapability]struct{},
	cacheGroups map[tc.CacheGroupName]tc.CacheGroupNullable,
	dsParams ParentDSParams,
	atsMajorVer int,
	dsOrigins map[ServerID]struct{},
) (string, error) {
	txt := ""

	if !HasRequiredCapabilities(serverCapabilities[*server.ID], dsRequiredCapabilities[*ds.ID]) {
		return "", nil
	}

	orgURI, err := GetOriginURI(*ds.OrgServerFQDN)
	if err != nil {
		return "", errors.New("Malformed ds '" + *ds.XMLID + "' origin  URI: '" + *ds.OrgServerFQDN + "': skipping!" + err.Error())
	}

	topology := nameTopologies[TopologyName(*ds.Topology)]
	if topology.Name == "" {
		return "", errors.New("DS " + *ds.XMLID + " topology '" + *ds.Topology + "' not found in Topologies!")
	}

	txt += "dest_domain=" + orgURI.Hostname() + " port=" + orgURI.Port()

	serverPlacement := getTopologyPlacement(tc.CacheGroupName(*server.Cachegroup), topology, cacheGroups, ds)
	if !serverPlacement.InTopology {
		return "", nil // server isn't in topology, no error
	}
	// TODO add Topology/Capabilities to remap.config

	parents, secondaryParents, err := GetTopologyParents(server, ds, servers, parentConfigParams, topology, serverPlacement.IsLastTier, serverCapabilities, dsRequiredCapabilities, dsOrigins)
	if err != nil {
		return "", errors.New("getting topology parents for '" + *ds.XMLID + "': skipping! " + err.Error())
	}
	if len(parents) == 0 {
		return "", errors.New("getting topology parents for '" + *ds.XMLID + "': no parents found! skipping! (Does your Topology have a CacheGroup with no servers in it?)")
	}

	txt += ` parent="` + strings.Join(parents, `;`) + `"`
	if len(secondaryParents) > 0 {
		txt += ` secondary_parent="` + strings.Join(secondaryParents, `;`) + `"`
		txt += getSecondaryModeStr(dsParams.TryAllPrimariesBeforeSecondary, atsMajorVer, tc.DeliveryServiceName(*ds.XMLID))
	}
	txt += ` round_robin=` + getTopologyRoundRobin(ds, serverParams, serverPlacement.IsLastCacheTier, dsParams.Algorithm)
	txt += ` go_direct=` + getTopologyGoDirect(ds, serverPlacement.IsLastTier)
	txt += ` qstring=` + getTopologyQueryString(ds, serverParams, serverPlacement.IsLastCacheTier, dsParams.Algorithm, dsParams.QueryStringHandling)
	txt += getTopologyParentIsProxyStr(serverPlacement.IsLastCacheTier)
	txt += getParentRetryStr(serverPlacement.IsLastCacheTier, atsMajorVer, dsParams.ParentRetry, dsParams.UnavailableServerRetryResponses, dsParams.MaxSimpleRetries, dsParams.MaxUnavailableServerRetries)
	txt += " # topology '" + *ds.Topology + "'"
	txt += "\n"
	return txt, nil
}

// getParentRetryStr builds the parent retry directive(s).
// If atsMajorVer < 6, "" is returned (ATS 5 and below don't support retry directives).
// If isLastCacheTier is false, "" is returned. This argument exists to simplify usage.
// If parentRetry is "", "" is returned (because the other directives are unused if parent_retry doesn't exist). This is allowed to simplify usage.
// If unavailableServerRetryResponses is not "", it must be valid. Use unavailableServerRetryResponsesValid to check.
// If maxSimpleRetries is "", ParentConfigDSParamDefaultMaxSimpleRetries will be used.
// If maxUnavailableServerRetries is "", ParentConfigDSParamDefaultMaxUnavailableServerRetries will be used.
func getParentRetryStr(isLastCacheTier bool, atsMajorVer int, parentRetry string, unavailableServerRetryResponses string, maxSimpleRetries string, maxUnavailableServerRetries string) string {
	if !isLastCacheTier || // allow !isLastCacheTier, to simplify usage.
		parentRetry == "" || // allow parentRetry to be empty, to simplify usage.
		atsMajorVer < 6 { // ATS 5 and below don't support parent_retry directives
		return ""
	}

	if maxSimpleRetries == "" {
		maxSimpleRetries = ParentConfigDSParamDefaultMaxSimpleRetries
	}
	if maxUnavailableServerRetries == "" {
		maxUnavailableServerRetries = ParentConfigDSParamDefaultMaxUnavailableServerRetries
	}

	txt := ` parent_retry=` + parentRetry
	if unavailableServerRetryResponses != "" {
		txt += ` unavailable_server_retry_responses=` + unavailableServerRetryResponses
	}
	txt += ` max_simple_retries=` + maxSimpleRetries + ` max_unavailable_server_retries=` + maxUnavailableServerRetries
	return txt
}

func getSecondaryModeStr(tryAllPrimariesBeforeSecondary bool, atsMajorVer int, ds tc.DeliveryServiceName) string {
	if !tryAllPrimariesBeforeSecondary {
		return ""
	}
	if atsMajorVer < 8 {
		log.Errorln("Delivery Service '" + string(ds) + "' had Parameter " + ParentConfigParamSecondaryMode + " but this cache is " + strconv.Itoa(atsMajorVer) + " and secondary_mode isn't supported in ATS until 8. Not using!")
		return ""
	}
	return ` secondary_mode=2` // See https://docs.trafficserver.apache.org/en/8.0.x/admin-guide/files/parent.config.en.html
}

func getTopologyParentIsProxyStr(serverIsLastCacheTier bool) string {
	if serverIsLastCacheTier {
		return ` parent_is_proxy=false`
	}
	return ""
}

func getTopologyRoundRobin(
	ds *tc.DeliveryServiceNullableV30,
	serverParams map[string]string,
	serverIsLastTier bool,
	algorithm string,
) string {
	roundRobinConsistentHash := "consistent_hash"
	if !serverIsLastTier {
		return roundRobinConsistentHash
	}
	if parentSelectAlg := serverParams[ParentConfigParamAlgorithm]; ds.OriginShield != nil && *ds.OriginShield != "" && strings.TrimSpace(parentSelectAlg) != "" {
		return parentSelectAlg
	}
	if algorithm != "" {
		return algorithm
	}
	return roundRobinConsistentHash
}

func getTopologyGoDirect(ds *tc.DeliveryServiceNullableV30, serverIsLastTier bool) string {
	if !serverIsLastTier {
		return "false"
	}
	if ds.OriginShield != nil && *ds.OriginShield != "" {
		return "true"
	}
	if ds.MultiSiteOrigin != nil && *ds.MultiSiteOrigin {
		return "false"
	}
	return "true"
}

func getTopologyQueryString(
	ds *tc.DeliveryServiceNullableV30,
	serverParams map[string]string,
	serverIsLastTier bool,
	algorithm string,
	qStringHandling string,
) string {
	if serverIsLastTier {
		if ds.MultiSiteOrigin != nil && *ds.MultiSiteOrigin && qStringHandling == "" && algorithm == tc.AlgorithmConsistentHash && ds.QStringIgnore != nil && tc.QStringIgnore(*ds.QStringIgnore) == tc.QStringIgnoreUseInCacheKeyAndPassUp {
			return "consider"
		}
		return "ignore"
	}

	if param := serverParams[ParentConfigParamQStringHandling]; param != "" {
		return param
	}
	if qStringHandling != "" {
		return qStringHandling
	}
	if ds.QStringIgnore != nil && tc.QStringIgnore(*ds.QStringIgnore) == tc.QStringIgnoreUseInCacheKeyAndPassUp {
		return "consider"
	}
	return "ignore"
}

// serverParentageParams gets the Parameters used for parent= line, or defaults if they don't exist
// Returns the Parameters used for parent= lines, for the given server.
func serverParentageParams(sv *tc.ServerNullable, params []ParameterWithProfilesMap) ProfileCache {
	// TODO deduplicate with atstccfg/parentdotconfig.go
	profileCache := DefaultProfileCache()
	if sv.TCPPort != nil {
		profileCache.Port = *sv.TCPPort
	}
	for _, param := range params {
		if _, ok := param.ProfileNames[*sv.Profile]; !ok {
			continue
		}
		switch param.Name {
		case ParentConfigCacheParamWeight:
			profileCache.Weight = param.Value
		case ParentConfigCacheParamPort:
			i, err := strconv.ParseInt(param.Value, 10, 64)
			if err != nil {
				log.Errorln("parent.config generation: port param is not an integer, skipping! : " + err.Error())
			} else {
				profileCache.Port = int(i)
			}
		case ParentConfigCacheParamUseIP:
			profileCache.UseIP = param.Value == "1"
		case ParentConfigCacheParamRank:
			i, err := strconv.ParseInt(param.Value, 10, 64)
			if err != nil {
				log.Errorln("parent.config generation: rank param is not an integer, skipping! : " + err.Error())
			} else {
				profileCache.Rank = int(i)
			}
		case ParentConfigCacheParamNotAParent:
			profileCache.NotAParent = param.Value != "false"
		}
	}
	return profileCache
}

func serverParentStr(sv *tc.ServerNullable, svParams ProfileCache) (string, error) {
	if svParams.NotAParent {
		return "", nil
	}
	host := ""
	if svParams.UseIP {
		// TODO get service interface here
		ip := GetServerIPAddress(sv)
		if ip == nil {
			return "", errors.New("server params Use IP, but has no valid IPv4 Service Address")
		}
		host = ip.String()
	} else {
		host = *sv.HostName + "." + *sv.DomainName
	}
	return host + ":" + strconv.Itoa(svParams.Port) + "|" + svParams.Weight, nil
}

func GetTopologyParents(
	server *tc.ServerNullable,
	ds *tc.DeliveryServiceNullableV30,
	servers []tc.ServerNullable,
	parentConfigParams []ParameterWithProfilesMap, // all params with configFile parent.confign
	topology tc.Topology,
	serverIsLastTier bool,
	serverCapabilities map[int]map[ServerCapability]struct{},
	dsRequiredCapabilities map[int]map[ServerCapability]struct{},
	dsOrigins map[ServerID]struct{}, // for Topology DSes, MSO still needs DeliveryServiceServer assignments.
) ([]string, []string, error) {
	// If it's the last tier, then the parent is the origin.
	// Note this doesn't include MSO, whose final tier cachegroup points to the origin cachegroup.
	if serverIsLastTier {
		orgURI, err := GetOriginURI(*ds.OrgServerFQDN) // TODO pass, instead of calling again
		if err != nil {
			return nil, nil, err
		}
		return []string{orgURI.Host}, nil, nil
	}

	svNode := tc.TopologyNode{}
	for _, node := range topology.Nodes {
		if node.Cachegroup == *server.Cachegroup {
			svNode = node
			break
		}
	}
	if svNode.Cachegroup == "" {
		return nil, nil, errors.New("This server '" + *server.HostName + "' not in DS " + *ds.XMLID + " topology, skipping")
	}

	if len(svNode.Parents) == 0 {
		return nil, nil, errors.New("DS " + *ds.XMLID + " topology '" + *ds.Topology + "' is last tier, but NonLastTier called! Should never happen")
	}
	if numParents := len(svNode.Parents); numParents > 2 {
		log.Errorln("DS " + *ds.XMLID + " topology '" + *ds.Topology + "' has " + strconv.Itoa(numParents) + " parent nodes, but Apache Traffic Server only supports Primary and Secondary (2) lists of parents. CacheGroup nodes after the first 2 will be ignored!")
	}
	if len(topology.Nodes) <= svNode.Parents[0] {
		return nil, nil, errors.New("DS " + *ds.XMLID + " topology '" + *ds.Topology + "' node parent " + strconv.Itoa(svNode.Parents[0]) + " greater than number of topology nodes " + strconv.Itoa(len(topology.Nodes)) + ". Cannot create parents!")
	}
	if len(svNode.Parents) > 1 && len(topology.Nodes) <= svNode.Parents[1] {
		log.Errorln("DS " + *ds.XMLID + " topology '" + *ds.Topology + "' node secondary parent " + strconv.Itoa(svNode.Parents[1]) + " greater than number of topology nodes " + strconv.Itoa(len(topology.Nodes)) + ". Secondary parent will be ignored!")
	}

	parentCG := topology.Nodes[svNode.Parents[0]].Cachegroup
	secondaryParentCG := ""
	if len(svNode.Parents) > 1 && len(topology.Nodes) > svNode.Parents[1] {
		secondaryParentCG = topology.Nodes[svNode.Parents[1]].Cachegroup
	}

	if parentCG == "" {
		return nil, nil, errors.New("Server '" + *server.HostName + "' DS " + *ds.XMLID + " topology '" + *ds.Topology + "' cachegroup '" + *server.Cachegroup + "' topology node parent " + strconv.Itoa(svNode.Parents[0]) + " is not in the topology!")
	}

	parentStrs := []string{}
	secondaryParentStrs := []string{}

	serversWithParams := []ServerWithParams{}
	for _, sv := range servers {
		serversWithParams = append(serversWithParams, ServerWithParams{
			ServerNullable: sv,
			Params:         serverParentageParams(&sv, parentConfigParams),
		})
	}
	sort.Sort(ServersWithParamsSortByRank(serversWithParams))

	for _, sv := range serversWithParams {
		if sv.ID == nil {
			log.Errorln("TO Servers server had nil ID, skipping")
			continue
		} else if sv.Cachegroup == nil {
			log.Errorln("TO Servers server had nil Cachegroup, skipping")
			continue
		} else if sv.CDNName == nil {
			log.Errorln("parent.config generation: TO servers had server with missing CDNName, skipping!")
			continue
		} else if sv.Status == nil || *sv.Status == "" {
			log.Errorln("parent.config generation: TO servers had server with missing Status, skipping!")
			continue
		}

		if tc.CacheType(sv.Type) != tc.CacheTypeEdge && tc.CacheType(sv.Type) != tc.CacheTypeMid && sv.Type != tc.OriginTypeName {
			continue // only consider edges, mids, and origins in the CacheGroup.
		}
		if _, dsHasOrigin := dsOrigins[ServerID(*sv.ID)]; sv.Type == tc.OriginTypeName && !dsHasOrigin {
			continue
		}
		if *sv.CDNName != *server.CDNName {
			continue
		}
		if *sv.Status != string(tc.CacheStatusReported) && *sv.Status != string(tc.CacheStatusOnline) {
			continue
		}

		if !HasRequiredCapabilities(serverCapabilities[*sv.ID], dsRequiredCapabilities[*ds.ID]) {
			continue
		}
		if *sv.Cachegroup == parentCG {
			parentStr, err := serverParentStr(&sv.ServerNullable, sv.Params)
			if err != nil {
				return nil, nil, errors.New("getting server parent string: " + err.Error())
			}
			if parentStr != "" { // will be empty if server is not_a_parent (possibly other reasons)
				parentStrs = append(parentStrs, parentStr)
			}
		}
		if *sv.Cachegroup == secondaryParentCG {
			parentStr, err := serverParentStr(&sv.ServerNullable, sv.Params)
			if err != nil {
				return nil, nil, errors.New("getting server parent string: " + err.Error())
			}
			secondaryParentStrs = append(secondaryParentStrs, parentStr)
		}
	}

	return parentStrs, secondaryParentStrs, nil
}

func GetOriginURI(fqdn string) (*url.URL, error) {
	orgURI, err := url.Parse(fqdn) // TODO verify origin is always a host:port
	if err != nil {
		return nil, errors.New("parsing: " + err.Error())
	}
	if orgURI.Port() == "" {
		if orgURI.Scheme == "http" {
			orgURI.Host += ":80"
		} else if orgURI.Scheme == "https" {
			orgURI.Host += ":443"
		} else {
			log.Errorln("parent.config generation non-top-level: origin '" + fqdn + "' is unknown scheme '" + orgURI.Scheme + "', but has no port! Using as-is! ")
		}
	}
	return orgURI, nil
}

// getParentStrs returns the parents= and secondary_parents= strings for ATS parent.config lines.
func getParentStrs(
	ds *tc.DeliveryServiceNullableV30,
	dsRequiredCapabilities map[int]map[ServerCapability]struct{},
	parentInfos []ParentInfo,
	atsMajorVer int,
	tryAllPrimariesBeforeSecondary bool,
) (string, string) {
	parentInfo := []string{}
	secondaryParentInfo := []string{}

	sort.Sort(ParentInfoSortByRank(parentInfos))

	for _, parent := range parentInfos { // TODO fix magic key
		if !HasRequiredCapabilities(parent.Capabilities, dsRequiredCapabilities[*ds.ID]) {
			continue
		}

		pTxt := parent.Format()
		if parent.PrimaryParent {
			parentInfo = append(parentInfo, pTxt)
		} else if parent.SecondaryParent {
			secondaryParentInfo = append(secondaryParentInfo, pTxt)
		}
	}

	if len(parentInfo) == 0 {
		parentInfo = secondaryParentInfo
		secondaryParentInfo = []string{}
	}

	// TODO remove duplicate code with top level if block
	seen := map[string]struct{}{} // TODO change to host+port? host isn't unique
	parentInfo, seen = util.RemoveStrDuplicates(parentInfo, seen)
	secondaryParentInfo, seen = util.RemoveStrDuplicates(secondaryParentInfo, seen)

	dsName := tc.DeliveryServiceName("")
	if ds != nil && ds.XMLID != nil {
		dsName = tc.DeliveryServiceName(*ds.XMLID)
	}

	parents := ""
	secondaryParents := "" // "secparents" in Perl

	if atsMajorVer >= 6 && len(secondaryParentInfo) > 0 {
		parents = `parent="` + strings.Join(parentInfo, "") + `"`
		secondaryParents = ` secondary_parent="` + strings.Join(secondaryParentInfo, "") + `"`
		secondaryParents += getSecondaryModeStr(tryAllPrimariesBeforeSecondary, atsMajorVer, dsName)
	} else {
		parents = `parent="` + strings.Join(parentInfo, "") + strings.Join(secondaryParentInfo, "") + `"`
	}

	return parents, secondaryParents
}

// getMSOParentStrs returns the parents= and secondary_parents= strings for ATS parent.config lines, for MSO.
func getMSOParentStrs(
	ds *tc.DeliveryServiceNullableV30,
	parentInfos []ParentInfo,
	atsMajorVer int,
	dsRequiredCapabilities map[int]map[ServerCapability]struct{},
	msoAlgorithm string,
	tryAllPrimariesBeforeSecondary bool,
) (string, string) {
	// TODO determine why MSO is different, and if possible, combine with getParentAndSecondaryParentStrs.

	rankedParents := ParentInfoSortByRank(parentInfos)
	sort.Sort(rankedParents)

	parentInfo := []string{}
	secondaryParentInfo := []string{}
	nullParentInfo := []string{}
	for _, parent := range ([]ParentInfo)(rankedParents) {
		if !HasRequiredCapabilities(parent.Capabilities, dsRequiredCapabilities[*ds.ID]) {
			continue
		}

		if parent.PrimaryParent {
			parentInfo = append(parentInfo, parent.Format())
		} else if parent.SecondaryParent {
			secondaryParentInfo = append(secondaryParentInfo, parent.Format())
		} else {
			nullParentInfo = append(nullParentInfo, parent.Format())
		}
	}

	if len(parentInfo) == 0 {
		// If no parents are found in the secondary parent either, then set the null parent list (parents in neither secondary or primary)
		// as the secondary parent list and clear the null parent list.
		if len(secondaryParentInfo) == 0 {
			secondaryParentInfo = nullParentInfo
			nullParentInfo = []string{}
		}
		parentInfo = secondaryParentInfo
		secondaryParentInfo = []string{} // TODO should thi be '= secondary'? Currently emulates Perl
	}

	// TODO benchmark, verify this isn't slow. if it is, it could easily be made faster
	seen := map[string]struct{}{} // TODO change to host+port? host isn't unique
	parentInfo, seen = util.RemoveStrDuplicates(parentInfo, seen)
	secondaryParentInfo, seen = util.RemoveStrDuplicates(secondaryParentInfo, seen)
	nullParentInfo, seen = util.RemoveStrDuplicates(nullParentInfo, seen)

	secondaryParentStr := strings.Join(secondaryParentInfo, "") + strings.Join(nullParentInfo, "")

	dsName := tc.DeliveryServiceName("")
	if ds != nil && ds.XMLID != nil {
		dsName = tc.DeliveryServiceName(*ds.XMLID)
	}

	// If the ats version supports it and the algorithm is consistent hash, put secondary and non-primary parents into secondary parent group.
	// This will ensure that secondary and tertiary parents will be unused unless all hosts in the primary group are unavailable.

	parents := ""
	secondaryParents := ""

	if atsMajorVer >= 6 && msoAlgorithm == "consistent_hash" && len(secondaryParentStr) > 0 {
		parents = `parent="` + strings.Join(parentInfo, "") + `"`
		secondaryParents = ` secondary_parent="` + secondaryParentStr + `"`
		secondaryParents += getSecondaryModeStr(tryAllPrimariesBeforeSecondary, atsMajorVer, dsName)
	} else {
		parents = `parent="` + strings.Join(parentInfo, "") + secondaryParentStr + `"`
	}
	return parents, secondaryParents
}

func MakeParentInfo(
	serverParentCGData ServerParentCacheGroupData,
	serverDomain string, // getCDNDomainByProfileID(tx, server.ProfileID)
	profileCaches map[ProfileID]ProfileCache, // getServerParentCacheGroupProfiles(tx, server)
	originServers map[OriginHost][]CGServer, // getServerParentCacheGroupProfiles(tx, server)
) map[OriginHost][]ParentInfo {
	parentInfos := map[OriginHost][]ParentInfo{}

	// note servers also contains an "all" key
	for originHost, servers := range originServers {
		for _, row := range servers {
			profile := profileCaches[row.ProfileID]
			if profile.NotAParent {
				continue
			}
			// Perl has this check, but we only select servers ("deliveryServices" in Perl) with the right CDN in the first place
			// if profile.Domain != serverDomain {
			// 	continue
			// }

			parentInf := ParentInfo{
				Host:            row.ServerHost,
				Port:            profile.Port,
				Domain:          row.Domain,
				Weight:          profile.Weight,
				UseIP:           profile.UseIP,
				Rank:            profile.Rank,
				IP:              row.ServerIP,
				PrimaryParent:   serverParentCGData.ParentID == row.CacheGroupID,
				SecondaryParent: serverParentCGData.SecondaryParentID == row.CacheGroupID,
				Capabilities:    row.Capabilities,
			}
			if parentInf.Port < 1 {
				parentInf.Port = row.ServerPort
			}
			parentInfos[originHost] = append(parentInfos[originHost], parentInf)
		}
	}
	return parentInfos
}

// unavailableServerRetryResponsesValid returns whether a unavailable_server_retry_responses parameter is valid for an ATS parent rule.
func unavailableServerRetryResponsesValid(s string) bool {
	// optimization if param is empty
	if s == "" {
		return false
	}
	re := regexp.MustCompile(`^"(:?\d{3},)+\d{3}"\s*$`) // TODO benchmark, cache if performance matters
	return re.MatchString(s)
}

// HasRequiredCapabilities returns whether the given caps has all the required capabilities in the given reqCaps.
func HasRequiredCapabilities(caps map[ServerCapability]struct{}, reqCaps map[ServerCapability]struct{}) bool {
	for reqCap, _ := range reqCaps {
		if _, ok := caps[reqCap]; !ok {
			return false
		}
	}
	return true
}

func GetOriginServersAndProfileCaches(
	cgServers map[int]tc.ServerNullable,
	parentServerDSes map[int]map[int]struct{},
	profileParentConfigParams map[string]map[string]string, // map[profileName][paramName]paramVal
	dses []tc.DeliveryServiceNullableV30,
	serverCapabilities map[int]map[ServerCapability]struct{},
	dsRequiredCapabilities map[int]map[ServerCapability]struct{},
) (map[OriginHost][]CGServer, map[ProfileID]ProfileCache, error) {
	originServers := map[OriginHost][]CGServer{}  // "deliveryServices" in Perl
	profileCaches := map[ProfileID]ProfileCache{} // map[profileID]ProfileCache

	dsIDMap := map[int]tc.DeliveryServiceNullableV30{}
	for _, ds := range dses {
		if ds.ID == nil {
			return nil, nil, errors.New("delivery services got nil ID!")
		}
		if !ds.Type.IsHTTP() && !ds.Type.IsDNS() {
			continue // skip ANY_MAP, STEERING, etc
		}
		dsIDMap[*ds.ID] = ds
	}

	allDSMap := map[int]tc.DeliveryServiceNullableV30{} // all DSes for this server, NOT all dses in TO
	for _, dsIDs := range parentServerDSes {
		for dsID, _ := range dsIDs {
			if _, ok := dsIDMap[dsID]; !ok {
				// this is normal if the TO was too old to understand our /deliveryserviceserver?servers= query param
				// In which case, the DSS will include DSes from other CDNs, which aren't in the dsIDMap
				// If the server was new enough to respect the params, this should never happen.
				// log.Warnln("getting delivery services: parent server DS %v not in dsIDMap\n", dsID)
				continue
			}
			if _, ok := allDSMap[dsID]; !ok {
				allDSMap[dsID] = dsIDMap[dsID]
			}
		}
	}

	dsOrigins, err := GetDSOrigins(allDSMap)
	if err != nil {
		return nil, nil, errors.New("getting DS origins: " + err.Error())
	}

	profileParams := GetParentConfigProfileParams(cgServers, profileParentConfigParams)

	for _, cgServer := range cgServers {
		if cgServer.ID == nil {
			log.Errorln("parent.config getting origin servers: got server with nil ID, skipping!")
			continue
		} else if cgServer.HostName == nil {
			log.Errorln("parent.config getting origin servers: got server with nil HostName, skipping!")
			continue
		} else if cgServer.TCPPort == nil {
			log.Errorln("parent.config getting origin servers: got server with nil TCPPort, skipping!")
			continue
		} else if cgServer.CachegroupID == nil {
			log.Errorln("parent.config getting origin servers: got server with nil CachegroupID, skipping!")
			continue
		} else if cgServer.StatusID == nil {
			log.Errorln("parent.config getting origin servers: got server with nil StatusID, skipping!")
			continue
		} else if cgServer.TypeID == nil {
			log.Errorln("parent.config getting origin servers: got server with nil TypeID, skipping!")
			continue
		} else if cgServer.ProfileID == nil {
			log.Errorln("parent.config getting origin servers: got server with nil ProfileID, skipping!")
			continue
		} else if cgServer.CDNID == nil {
			log.Errorln("parent.config getting origin servers: got server with nil CDNID, skipping!")
			continue
		} else if cgServer.DomainName == nil {
			log.Errorln("parent.config getting origin servers: got server with nil DomainName, skipping!")
			continue
		}

		ipAddr := GetServerIPAddress(&cgServer)
		if ipAddr == nil {
			log.Errorln("parent.config getting origin servers: got server with no valid IP Address, skipping!")
			continue
		}

		realCGServer := CGServer{
			ServerID:     ServerID(*cgServer.ID),
			ServerHost:   *cgServer.HostName,
			ServerIP:     ipAddr.String(),
			ServerPort:   *cgServer.TCPPort,
			CacheGroupID: *cgServer.CachegroupID,
			Status:       *cgServer.StatusID,
			Type:         *cgServer.TypeID,
			ProfileID:    ProfileID(*cgServer.ProfileID),
			CDN:          *cgServer.CDNID,
			TypeName:     cgServer.Type,
			Domain:       *cgServer.DomainName,
			Capabilities: serverCapabilities[*cgServer.ID],
		}

		if cgServer.Type == tc.OriginTypeName {
			for dsID, _ := range parentServerDSes[*cgServer.ID] { // map[serverID][]dsID
				orgURI := dsOrigins[dsID]
				if orgURI == nil {
					// log.Warnln("ds %v has no origins! Skipping!\n", dsID) // TODO determine if this is normal
					continue
				}
				if HasRequiredCapabilities(serverCapabilities[*cgServer.ID], dsRequiredCapabilities[dsID]) {
					orgHost := OriginHost(orgURI.Host)
					originServers[orgHost] = append(originServers[orgHost], realCGServer)
				} else {
					log.Errorf("ds %v server %v missing required caps, skipping!\n", dsID, orgURI.Host)
				}
			}
		} else {
			originServers[DeliveryServicesAllParentsKey] = append(originServers[DeliveryServicesAllParentsKey], realCGServer)
		}

		if _, profileCachesHasProfile := profileCaches[realCGServer.ProfileID]; !profileCachesHasProfile {
			if profileCache, profileParamsHasProfile := profileParams[*cgServer.Profile]; !profileParamsHasProfile {
				log.Warnf("cachegroup has server with profile %+v but that profile has no parameters\n", *cgServer.ProfileID)
				profileCaches[realCGServer.ProfileID] = DefaultProfileCache()
			} else {
				profileCaches[realCGServer.ProfileID] = profileCache
			}
		}
	}

	return originServers, profileCaches, nil
}

func GetParentConfigProfileParams(
	cgServers map[int]tc.ServerNullable,
	profileParentConfigParams map[string]map[string]string, // map[profileName][paramName]paramVal
) map[string]ProfileCache {
	parentConfigServerCacheProfileParams := map[string]ProfileCache{} // map[profileName]ProfileCache
	for _, cgServer := range cgServers {
		if cgServer.Profile == nil {
			log.Errorln("getting parent config profile params: server has nil profile, skipping!")
			continue
		}
		profileCache, ok := parentConfigServerCacheProfileParams[*cgServer.Profile]
		if !ok {
			profileCache = DefaultProfileCache()
		}
		params, ok := profileParentConfigParams[*cgServer.Profile]
		if !ok {
			parentConfigServerCacheProfileParams[*cgServer.Profile] = profileCache
			continue
		}
		for name, val := range params {
			switch name {
			case ParentConfigCacheParamWeight:
				// f, err := strconv.ParseFloat(param.Val, 64)
				// if err != nil {
				// 	log.Errorln("parent.config generation: weight param is not a float, skipping! : " + err.Error())
				// } else {
				// 	profileCache.Weight = f
				// }
				// TODO validate float?
				profileCache.Weight = val
			case ParentConfigCacheParamPort:
				i, err := strconv.ParseInt(val, 10, 64)
				if err != nil {
					log.Errorln("parent.config generation: port param is not an integer, skipping! : " + err.Error())
				} else {
					profileCache.Port = int(i)
				}
			case ParentConfigCacheParamUseIP:
				profileCache.UseIP = val == "1"
			case ParentConfigCacheParamRank:
				i, err := strconv.ParseInt(val, 10, 64)
				if err != nil {
					log.Errorln("parent.config generation: rank param is not an integer, skipping! : " + err.Error())
				} else {
					profileCache.Rank = int(i)
				}
			case ParentConfigCacheParamNotAParent:
				profileCache.NotAParent = val != "false"
			}
		}
		parentConfigServerCacheProfileParams[*cgServer.Profile] = profileCache
	}
	return parentConfigServerCacheProfileParams
}

// GetDSOrigins takes a map[deliveryServiceID]DeliveryService, and returns a map[DeliveryServiceID]OriginURI.
func GetDSOrigins(dses map[int]tc.DeliveryServiceNullableV30) (map[int]*OriginURI, error) {
	dsOrigins := map[int]*OriginURI{}
	for _, ds := range dses {
		if ds.ID == nil {
			return nil, errors.New("ds has nil ID")
		}
		if ds.XMLID == nil {
			return nil, errors.New("ds has nil XMLID")
		}
		if ds.OrgServerFQDN == nil {
			log.Warnf("GetDSOrigins ds %v got nil OrgServerFQDN, skipping!\n", *ds.XMLID)
			continue
		}
		orgURL, err := url.Parse(*ds.OrgServerFQDN)
		if err != nil {
			return nil, errors.New("parsing ds '" + *ds.XMLID + "' OrgServerFQDN '" + *ds.OrgServerFQDN + "': " + err.Error())
		}
		if orgURL.Scheme == "" {
			return nil, errors.New("parsing ds '" + *ds.XMLID + "' OrgServerFQDN '" + *ds.OrgServerFQDN + "': " + "missing scheme")
		}
		if orgURL.Host == "" {
			return nil, errors.New("parsing ds '" + *ds.XMLID + "' OrgServerFQDN '" + *ds.OrgServerFQDN + "': " + "missing scheme")
		}

		scheme := orgURL.Scheme
		host := orgURL.Hostname()
		port := orgURL.Port()
		if port == "" {
			if scheme == "http" {
				port = "80"
			} else if scheme == "https" {
				port = "443"
			} else {
				log.Warnln("parsing ds '" + *ds.XMLID + "' OrgServerFQDN '" + *ds.OrgServerFQDN + "': " + "unknown scheme '" + scheme + "' and no port, leaving port empty!")
			}
		}
		dsOrigins[*ds.ID] = &OriginURI{Scheme: scheme, Host: host, Port: port}
	}
	return dsOrigins, nil
}

func makeDSOrigins(dsses []tc.DeliveryServiceServer, dses []tc.DeliveryServiceNullableV30, servers []tc.ServerNullable) map[DeliveryServiceID]map[ServerID]struct{} {
	dssMap := map[DeliveryServiceID]map[ServerID]struct{}{}
	for _, dss := range dsses {
		if dss.Server == nil || dss.DeliveryService == nil {
			log.Errorln("making parent.config, got deliveryserviceserver with nil values, skipping!")
			continue
		}
		dsID := DeliveryServiceID(*dss.DeliveryService)
		serverID := ServerID(*dss.Server)
		if dssMap[dsID] == nil {
			dssMap[dsID] = map[ServerID]struct{}{}
		}
		dssMap[dsID][serverID] = struct{}{}
	}

	svMap := map[ServerID]tc.ServerNullable{}
	for _, sv := range servers {
		if sv.ID == nil {
			log.Errorln("parent.config got server with missing ID, skipping!")
		}
		svMap[ServerID(*sv.ID)] = sv
	}

	dsOrigins := map[DeliveryServiceID]map[ServerID]struct{}{}
	for _, ds := range dses {
		if ds.ID == nil {
			log.Errorln("parent.config got ds with missing ID, skipping!")
			continue
		}
		dsID := DeliveryServiceID(*ds.ID)
		assignedServers := dssMap[dsID]
		for svID, _ := range assignedServers {
			sv := svMap[svID]
			if sv.Type != tc.OriginTypeName {
				continue
			}
			if dsOrigins[dsID] == nil {
				dsOrigins[dsID] = map[ServerID]struct{}{}
			}
			dsOrigins[dsID][svID] = struct{}{}
		}
	}
	return dsOrigins
}
