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
const ParentConfigParamUnavailableServerRetryResponses = "mso.unavailable_server_retry_responses"
const ParentConfigParamMaxSimpleRetries = "mso.max_simple_retries"
const ParentConfigParamMaxUnavailableServerRetries = "mso.max_unavailable_server_retries"
const ParentConfigParamAlgorithm = "algorithm"
const ParentConfigParamQString = "qstring"

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

func (s ParentInfoSortByRank) Len() int           { return len(([]ParentInfo)(s)) }
func (s ParentInfoSortByRank) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s ParentInfoSortByRank) Less(i, j int) bool { return s[i].Rank < s[j].Rank }

type ParentConfigDSTopLevelSortByName []ParentConfigDSTopLevel

func (s ParentConfigDSTopLevelSortByName) Len() int      { return len(([]ParentConfigDSTopLevel)(s)) }
func (s ParentConfigDSTopLevelSortByName) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s ParentConfigDSTopLevelSortByName) Less(i, j int) bool {
	// TODO make this match the Perl sort "foreach my $ds ( sort @{ $data->{dslist} } )" ?
	return strings.Compare(string(s[i].Name), string(s[j].Name)) < 0
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
	serverInfo *ServerInfo, // getServerInfoByHost OR getServerInfoByID
	atsMajorVer int, // GetATSMajorVersion (TODO: determine if the cache itself [ORT via Yum] should produce this data, rather than asking TO?)
	toToolName string, // tm.toolname global parameter (TODO: cache itself?)
	toURL string, // tm.url global parameter (TODO: cache itself?)
	parentConfigDSes []ParentConfigDSTopLevel, // getParentConfigDSTopLevel(cdn) OR getParentConfigDS(server) (TODO determine how to handle non-top missing MSO?)
	serverParams map[string]string, // getParentConfigServerProfileParams(serverID)
	parentInfos map[OriginHost][]ParentInfo, // getParentInfo(profileID, parentCachegroupID, secondaryParentCachegroupID)
	server tc.Server,
	servers []tc.Server,
	topologies []tc.Topology,
	tcParentConfigParams []tc.Parameter,
	serverCapabilities map[int]map[ServerCapability]struct{},
	cacheGroupArr []tc.CacheGroupNullable,
) string {
	cacheGroups := MakeCGMap(cacheGroupArr)

	sort.Sort(ParentConfigDSTopLevelSortByName(parentConfigDSes))

	nameVersionStr := GetNameVersionStringFromToolNameAndURL(toToolName, toURL)
	hdr := HeaderCommentWithTOVersionStr(serverInfo.HostName, nameVersionStr)

	textArr := []string{}
	processedOriginsToDSNames := map[string]tc.DeliveryServiceName{}

	parentConfigParamsWithProfiles, err := TCParamsToParamsWithProfiles(tcParentConfigParams)
	if err != nil {
		log.Errorln("parent.config generation: error getting profiles from Traffic Ops Parameters, Parameters will not be considered for generation! : " + err.Error())
		parentConfigParamsWithProfiles = []ParameterWithProfiles{}
	}
	parentConfigParams := ParameterWithProfilesToMap(parentConfigParamsWithProfiles)

	for _, ds := range parentConfigDSes {
		log.Infoln("parent.config processing ds '" + ds.Name + "'")

		if existingDS, ok := processedOriginsToDSNames[ds.OriginFQDN]; ok {
			log.Errorln("parent.config generation: duplicate origin! services '" + string(ds.Name) + "' and '" + string(existingDS) + "' share origin '" + ds.OriginFQDN + "': skipping '" + string(ds.Name) + "'!")
			continue
		}

		// TODO put these in separate functions. No if-statement should be this long.
		if ds.Topology != "" {
			log.Infoln("parent.config generating Topology line for ds '" + ds.Name + "'")
			txt, err := GetTopologyParentConfigLine(server, servers, ds, serverParams, parentConfigParams, topologies, serverCapabilities, cacheGroups)
			if err != nil {
				log.Errorln(err)
				continue
			}
			if txt != "" { // will be empty with no error if this server isn't in the Topology, or if it doesn't have the Required Capabilities
				textArr = append(textArr, txt)
			}
		} else if serverInfo.IsTopLevelCache() {
			log.Infoln("parent.config generating top level line for ds '" + ds.Name + "'")
			parentQStr := "ignore"
			if ds.QStringHandling == "" && ds.MSOAlgorithm == tc.AlgorithmConsistentHash && ds.QStringIgnore == tc.QStringIgnoreUseInCacheKeyAndPassUp {
				parentQStr = "consider"
			}

			orgURI, err := GetOriginURI(ds.OriginFQDN)
			if err != nil {
				log.Errorln("Malformed ds '" + string(ds.Name) + "' origin  URI: '" + ds.OriginFQDN + "': skipping!" + err.Error())
				continue
			}

			textLine := ""

			if ds.OriginShield != "" {
				algorithm := ""
				if parentSelectAlg := serverParams[ParentConfigParamAlgorithm]; strings.TrimSpace(parentSelectAlg) != "" {
					algorithm = "round_robin=" + parentSelectAlg
				}
				textLine += "dest_domain=" + orgURI.Hostname() + " port=" + orgURI.Port() + " parent=" + ds.OriginShield + " " + algorithm + " go_direct=true\n"
			} else if ds.MultiSiteOrigin {
				textLine += "dest_domain=" + orgURI.Hostname() + " port=" + orgURI.Port() + " "
				if len(parentInfos) == 0 {
				}

				if len(parentInfos[OriginHost(orgURI.Hostname())]) == 0 {
					// TODO error? emulates Perl
					log.Warnln("ParentInfo: delivery service " + ds.Name + " has no parent servers")
				}

				parents, secondaryParents := getMSOParentStrs(ds, parentInfos[OriginHost(orgURI.Hostname())], atsMajorVer)
				textLine += parents + secondaryParents + ` round_robin=` + ds.MSOAlgorithm + ` qstring=` + parentQStr + ` go_direct=false parent_is_proxy=false`

				parentRetry := ds.MSOParentRetry
				if atsMajorVer >= 6 && parentRetry != "" {
					if unavailableServerRetryResponsesValid(ds.MSOUnavailableServerRetryResponses) {
						textLine += ` parent_retry=` + parentRetry + ` unavailable_server_retry_responses=` + ds.MSOUnavailableServerRetryResponses
					} else {
						if ds.MSOUnavailableServerRetryResponses != "" {
							log.Errorln("Malformed unavailable_server_retry_responses parameter '" + ds.MSOUnavailableServerRetryResponses + "', not using!")
						}
						textLine += ` parent_retry=` + parentRetry
					}
					textLine += ` max_simple_retries=` + ds.MSOMaxSimpleRetries + ` max_unavailable_server_retries=` + ds.MSOMaxUnavailableServerRetries
				}
				textLine += "\n" // TODO remove, and join later on "\n" instead of ""?
				textArr = append(textArr, textLine)
			}
		} else {
			log.Infoln("parent.config generating non-top level line for ds '" + ds.Name + "'")
			queryStringHandling := serverParams[ParentConfigParamQStringHandling] // "qsh" in Perl

			roundRobin := `round_robin=consistent_hash`
			goDirect := `go_direct=false`

			parents, secondaryParents := getParentStrs(ds, parentInfos[DeliveryServicesAllParentsKey], atsMajorVer)

			text := ""
			if ds.OriginFQDN == "" {
				continue // TODO warn? (Perl doesn't)
			}
			orgURI, err := GetOriginURI(ds.OriginFQDN)
			if err != nil {
				log.Errorln("Malformed ds '" + string(ds.Name) + "' origin  URI: '" + ds.OriginFQDN + "': skipping!" + err.Error())
				continue
			}

			// TODO encode this in a DSType func, IsGoDirect() ?
			if dsType := tc.DSType(ds.Type); dsType == tc.DSTypeHTTPNoCache || dsType == tc.DSTypeHTTPLive || dsType == tc.DSTypeDNSLive {
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
					dsQSH = ds.QStringHandling
				}
				parentQStr := dsQSH
				if parentQStr == "" {
					parentQStr = "ignore"
				}
				if ds.QStringIgnore == tc.QStringIgnoreUseInCacheKeyAndPassUp && dsQSH == "" {
					parentQStr = "consider"
				}

				text += `dest_domain=` + orgURI.Hostname() + ` port=` + orgURI.Port() + ` ` + parents + ` ` + secondaryParents + ` ` + roundRobin + ` ` + goDirect + ` qstring=` + parentQStr + "\n"
			}
			textArr = append(textArr, text)
		}
		processedOriginsToDSNames[ds.OriginFQDN] = ds.Name
	}

	// TODO determine if this is necessary. It's super-dangerous, and moreover ignores Server Capabilitites.
	defaultDestText := ""
	if !serverInfo.IsTopLevelCache() {
		parents, secondaryParents := getParentStrs(ParentConfigDSTopLevel{}, parentInfos[DeliveryServicesAllParentsKey], atsMajorVer)
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

func GetTopologyParentConfigLine(
	server tc.Server,
	servers []tc.Server,
	ds ParentConfigDSTopLevel,
	serverParams map[string]string,
	parentConfigParams []ParameterWithProfilesMap, // all params with configFile parent.config
	topologies []tc.Topology,
	serverCapabilities map[int]map[ServerCapability]struct{},
	cacheGroups map[tc.CacheGroupName]tc.CacheGroupNullable,
) (string, error) {
	txt := ""

	if !HasRequiredCapabilities(serverCapabilities[server.ID], ds.RequiredCapabilities) {
		return "", nil
	}

	orgURI, err := GetOriginURI(ds.OriginFQDN)
	if err != nil {
		return "", errors.New("Malformed ds '" + string(ds.Name) + "' origin  URI: '" + ds.OriginFQDN + "': skipping!" + err.Error())
	}

	// This could be put in a map beforehand to only iterate once, if performance mattered
	topology := tc.Topology{}
	for _, to := range topologies {
		if to.Name == ds.Topology {
			topology = to
			break
		}
	}
	if topology.Name == "" {
		return "", errors.New("DS " + string(ds.Name) + " topology '" + ds.Topology + "' not found in Topologies!")
	}

	txt += "dest_domain=" + orgURI.Hostname() + " port=" + orgURI.Port()

	log.Errorf("DEBUG topo GetTopologyParentConfigLine calling getTopologyPlacement cg '" + server.Cachegroup + "'\n")
	serverPlacement := getTopologyPlacement(tc.CacheGroupName(server.Cachegroup), topology, cacheGroups)
	if !serverPlacement.InTopology {
		return "", nil // server isn't in topology, no error
	}
	// TODO add Topology/Capabilities to remap.config

	parents, secondaryParents, err := GetTopologyParents(server, ds, servers, parentConfigParams, topology, serverPlacement.IsLastTier, serverCapabilities)
	if err != nil {
		return "", errors.New("getting topology parents for '" + string(ds.Name) + "': skipping! " + err.Error())
	}
	txt += ` parent="` + strings.Join(parents, `;`) + `"`
	if len(secondaryParents) > 0 {
		txt += ` secondary_parent="` + strings.Join(secondaryParents, `;`) + `"`
	}
	txt += ` round_robin=` + getTopologyRoundRobin(ds, serverParams, serverPlacement.IsLastTier)
	txt += ` go_direct=` + getTopologyGoDirect(ds, serverPlacement.IsLastTier)
	txt += ` qstring=` + getTopologyQueryString(ds, serverParams, serverPlacement.IsLastTier)
	txt += getTopologyParentIsProxyStr(serverPlacement.IsLastTier)
	txt += " # topology '" + ds.Topology + "'"
	txt += "\n"
	return txt, nil
}

func getTopologyParentIsProxyStr(serverIsLastTier bool) string {
	if serverIsLastTier {
		return ` parent_is_proxy=false`
	}
	return ""
}

func getTopologyRoundRobin(ds ParentConfigDSTopLevel, serverParams map[string]string, serverIsLastTier bool) string {
	roundRobinConsistentHash := "consistent_hash"
	if !serverIsLastTier {
		return roundRobinConsistentHash
	}
	if parentSelectAlg := serverParams[ParentConfigParamAlgorithm]; ds.OriginShield != "" && strings.TrimSpace(parentSelectAlg) != "" {
		return parentSelectAlg
	}
	if ds.MultiSiteOrigin {
		return ds.MSOAlgorithm
	}
	return roundRobinConsistentHash
}

func getTopologyGoDirect(ds ParentConfigDSTopLevel, serverIsLastTier bool) string {
	if !serverIsLastTier {
		return "false"
	}
	if ds.OriginShield != "" {
		return "true"
	}
	if ds.MultiSiteOrigin {
		return "false"
	}
	return "true"
}

func getTopologyQueryString(ds ParentConfigDSTopLevel, serverParams map[string]string, serverIsLastTier bool) string {
	if serverIsLastTier {
		if ds.MultiSiteOrigin && ds.QStringHandling == "" && ds.MSOAlgorithm == tc.AlgorithmConsistentHash && ds.QStringIgnore == tc.QStringIgnoreUseInCacheKeyAndPassUp {
			return "consider"
		}
		return "ignore"
	}

	if param := serverParams[ParentConfigParamQStringHandling]; param != "" {
		return param
	}
	if ds.QStringHandling != "" {
		return ds.QStringHandling
	}
	if ds.QStringIgnore == tc.QStringIgnoreUseInCacheKeyAndPassUp {
		return "consider"
	}
	return "ignore"
}

// serverParentageParams gets the Parameters used for parent= line, or defaults if they don't exist
// Returns the Parameters used for parent= lines, for the given server.
func serverParentageParams(sv tc.Server, params []ParameterWithProfilesMap) ProfileCache {
	// TODO deduplicate with atstccfg/parentdotconfig.go
	profileCache := DefaultProfileCache()
	profileCache.Port = sv.TCPPort
	for _, param := range params {
		if _, ok := param.ProfileNames[sv.Profile]; !ok {
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

func serverParentStr(sv tc.Server, params []ParameterWithProfilesMap) string {
	svParams := serverParentageParams(sv, params)
	if svParams.NotAParent {
		return ""
	}
	host := ""
	if svParams.UseIP {
		host = sv.IPAddress
	} else {
		host = sv.HostName + "." + sv.DomainName
	}
	return host + ":" + strconv.Itoa(svParams.Port) + "|" + svParams.Weight
}

func GetTopologyParents(
	server tc.Server,
	ds ParentConfigDSTopLevel,
	servers []tc.Server,
	parentConfigParams []ParameterWithProfilesMap, // all params with configFile parent.confign
	topology tc.Topology,
	serverIsLastTier bool,
	serverCapabilities map[int]map[ServerCapability]struct{},
) ([]string, []string, error) {
	// If it's the last tier, then the parent is the origin.
	// Note this doesn't include MSO, whose final tier cachegroup points to the origin cachegroup.
	if serverIsLastTier {
		orgURI, err := GetOriginURI(ds.OriginFQDN) // TODO pass, instead of calling again
		if err != nil {
			return nil, nil, err
		}
		return []string{orgURI.Host}, nil, nil
	}

	svNode := tc.TopologyNode{}
	for _, node := range topology.Nodes {
		if node.Cachegroup == server.Cachegroup {
			svNode = node
			break
		}
	}
	if svNode.Cachegroup == "" {
		return nil, nil, errors.New("This server '" + server.HostName + "' not in DS " + string(ds.Name) + " topology, skipping")
	}

	if len(svNode.Parents) == 0 {
		return nil, nil, errors.New("DS " + string(ds.Name) + " topology '" + ds.Topology + "' is last tier, but NonLastTier called! Should never happen")
	}
	if numParents := len(svNode.Parents); numParents > 2 {
		log.Errorln("DS " + string(ds.Name) + " topology '" + ds.Topology + "' has " + strconv.Itoa(numParents) + " parent nodes, but Apache Traffic Server only supports Primary and Secondary (2) lists of parents. CacheGroup nodes after the first 2 will be ignored!")
	}
	if len(topology.Nodes) <= svNode.Parents[0] {
		return nil, nil, errors.New("DS " + string(ds.Name) + " topology '" + ds.Topology + "' node parent " + strconv.Itoa(svNode.Parents[0]) + " greater than number of topology nodes " + strconv.Itoa(len(topology.Nodes)) + ". Cannot create parents!")
	}
	if len(svNode.Parents) > 1 && len(topology.Nodes) <= svNode.Parents[1] {
		log.Errorln("DS " + string(ds.Name) + " topology '" + ds.Topology + "' node secondary parent " + strconv.Itoa(svNode.Parents[1]) + " greater than number of topology nodes " + strconv.Itoa(len(topology.Nodes)) + ". Secondary parent will be ignored!")
	}

	parentCG := topology.Nodes[svNode.Parents[0]].Cachegroup
	secondaryParentCG := ""
	if len(svNode.Parents) > 1 && len(topology.Nodes) > svNode.Parents[1] {
		secondaryParentCG = topology.Nodes[svNode.Parents[1]].Cachegroup
	}

	if parentCG == "" {
		return nil, nil, errors.New("Server '" + server.HostName + "' DS " + string(ds.Name) + " topology '" + ds.Topology + "' cachegroup '" + server.Cachegroup + "' topology node parent " + strconv.Itoa(svNode.Parents[0]) + " is not in the topology!")
	}

	parentStrs := []string{}
	secondaryParentStrs := []string{}
	for _, sv := range servers {
		if tc.CacheType(sv.Type) != tc.CacheTypeEdge && tc.CacheType(sv.Type) != tc.CacheTypeMid && sv.Type != tc.OriginTypeName {
			continue // only consider edges, mids, and origins in the CacheGroup.
		}
		if !HasRequiredCapabilities(serverCapabilities[sv.ID], ds.RequiredCapabilities) {
			continue
		}
		if sv.Cachegroup == parentCG {
			parentStr := serverParentStr(sv, parentConfigParams)
			if parentStr != "" { // will be empty if server is not_a_parent (possibly other reasons)
				parentStrs = append(parentStrs, parentStr)
			}
		}
		if sv.Cachegroup == secondaryParentCG {
			secondaryParentStrs = append(secondaryParentStrs, serverParentStr(sv, parentConfigParams))
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
func getParentStrs(ds ParentConfigDSTopLevel, parentInfos []ParentInfo, atsMajorVer int) (string, string) {
	parentInfo := []string{}
	secondaryParentInfo := []string{}

	sort.Sort(ParentInfoSortByRank(parentInfos))

	for _, parent := range parentInfos { // TODO fix magic key
		if !HasRequiredCapabilities(parent.Capabilities, ds.RequiredCapabilities) {
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

	parents := ""
	secondaryParents := "" // "secparents" in Perl
	sort.Sort(sort.StringSlice(parentInfo))
	sort.Sort(sort.StringSlice(secondaryParentInfo))

	if atsMajorVer >= 6 && len(secondaryParentInfo) > 0 {
		parents = `parent="` + strings.Join(parentInfo, "") + `"`
		secondaryParents = ` secondary_parent="` + strings.Join(secondaryParentInfo, "") + `"`
	} else {
		parents = `parent="` + strings.Join(parentInfo, "") + strings.Join(secondaryParentInfo, "") + `"`
	}
	return parents, secondaryParents
}

// getMSOParentStrs returns the parents= and secondary_parents= strings for ATS parent.config lines, for MSO.
func getMSOParentStrs(ds ParentConfigDSTopLevel, parentInfos []ParentInfo, atsMajorVer int) (string, string) {
	// TODO determine why MSO is different, and if possible, combine with getParentAndSecondaryParentStrs.

	rankedParents := ParentInfoSortByRank(parentInfos)
	sort.Sort(rankedParents)

	parentInfo := []string{}
	secondaryParentInfo := []string{}
	nullParentInfo := []string{}
	for _, parent := range ([]ParentInfo)(rankedParents) {
		if !HasRequiredCapabilities(parent.Capabilities, ds.RequiredCapabilities) {
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

	// If the ats version supports it and the algorithm is consistent hash, put secondary and non-primary parents into secondary parent group.
	// This will ensure that secondary and tertiary parents will be unused unless all hosts in the primary group are unavailable.

	parents := ""
	secondaryParents := ""

	if atsMajorVer >= 6 && ds.MSOAlgorithm == "consistent_hash" && len(secondaryParentStr) > 0 {
		parents = `parent="` + strings.Join(parentInfo, "") + `"`
		secondaryParents = ` secondary_parent="` + secondaryParentStr + `"`
	} else {
		parents = `parent="` + strings.Join(parentInfo, "") + secondaryParentStr + `"`
	}
	return parents, secondaryParents
}

func MakeParentInfo(
	server *ServerInfo,
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
				PrimaryParent:   server.ParentCacheGroupID == row.CacheGroupID,
				SecondaryParent: server.SecondaryParentCacheGroupID == row.CacheGroupID,
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
