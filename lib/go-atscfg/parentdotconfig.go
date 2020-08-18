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
	return strings.Compare(string(s[i].Name), string(s[j].Name)) < 0
}

type DSesSortByName []tc.DeliveryServiceNullable

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
	dses []tc.DeliveryServiceNullable,
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

		msoAlgorithm := ParentConfigDSParamDefaultMSOAlgorithm
		msoParentRetry := ParentConfigDSParamDefaultMSOParentRetry
		msoUnavailableServerRetryResponses := ParentConfigDSParamDefaultMSOUnavailableServerRetryResponses
		msoMaxSimpleRetries := ParentConfigDSParamDefaultMaxSimpleRetries
		msoMaxUnavailableServerRetries := ParentConfigDSParamDefaultMaxUnavailableServerRetries
		qStringHandling := ""
		if ds.ProfileName != nil && *ds.ProfileName != "" {
			if dsParams, ok := profileParentConfigParams[*ds.ProfileName]; ok {
				qStringHandling = dsParams[ParentConfigParamQStringHandling] // may be blank, no default
				if v, ok := dsParams[ParentConfigParamMSOAlgorithm]; ok && strings.TrimSpace(v) != "" {
					msoAlgorithm = v
				}
				if v, ok := dsParams[ParentConfigParamMSOParentRetry]; ok {
					msoParentRetry = v
				}
				if v, ok := dsParams[ParentConfigParamUnavailableServerRetryResponses]; ok {
					msoUnavailableServerRetryResponses = v
				}
				if v, ok := dsParams[ParentConfigParamMaxSimpleRetries]; ok {
					msoMaxSimpleRetries = v
				}
				if v, ok := dsParams[ParentConfigParamMaxUnavailableServerRetries]; ok {
					msoMaxUnavailableServerRetries = v
				}
			}
		}

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
				msoAlgorithm,
				msoParentRetry,
				msoUnavailableServerRetryResponses,
				msoMaxSimpleRetries,
				msoMaxUnavailableServerRetries,
				qStringHandling,
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
			if qStringHandling == "" && msoAlgorithm == tc.AlgorithmConsistentHash && ds.QStringIgnore != nil && tc.QStringIgnore(*ds.QStringIgnore) == tc.QStringIgnoreUseInCacheKeyAndPassUp {
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

				parents, secondaryParents := getMSOParentStrs(&ds, parentInfos[OriginHost(orgURI.Hostname())], atsMajorVer, dsRequiredCapabilities, msoAlgorithm)
				textLine += parents + secondaryParents + ` round_robin=` + msoAlgorithm + ` qstring=` + parentQStr + ` go_direct=false parent_is_proxy=false`

				parentRetry := msoParentRetry
				if atsMajorVer >= 6 && parentRetry != "" {
					if unavailableServerRetryResponsesValid(msoUnavailableServerRetryResponses) {
						textLine += ` parent_retry=` + parentRetry + ` unavailable_server_retry_responses=` + msoUnavailableServerRetryResponses
					} else {
						if msoUnavailableServerRetryResponses != "" {
							log.Errorln("Malformed unavailable_server_retry_responses parameter '" + msoUnavailableServerRetryResponses + "', not using!")
						}
						textLine += ` parent_retry=` + parentRetry
					}
					textLine += ` max_simple_retries=` + msoMaxSimpleRetries + ` max_unavailable_server_retries=` + msoMaxUnavailableServerRetries
				}
				textLine += "\n" // TODO remove, and join later on "\n" instead of ""?
				textArr = append(textArr, textLine)
			}
		} else {
			log.Infoln("parent.config generating non-top level line for ds '" + *ds.XMLID + "'")
			queryStringHandling := serverParams[ParentConfigParamQStringHandling] // "qsh" in Perl

			roundRobin := `round_robin=consistent_hash`
			goDirect := `go_direct=false`

			parents, secondaryParents := getParentStrs(&ds, dsRequiredCapabilities, parentInfos[DeliveryServicesAllParentsKey], atsMajorVer)

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
					dsQSH = qStringHandling
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
		invalidDS := &tc.DeliveryServiceNullable{}
		invalidDS.ID = util.IntPtr(-1)
		parents, secondaryParents := getParentStrs(invalidDS, dsRequiredCapabilities, parentInfos[DeliveryServicesAllParentsKey], atsMajorVer)
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
	server *tc.ServerNullable,
	servers []tc.ServerNullable,
	ds *tc.DeliveryServiceNullable,
	serverParams map[string]string,
	parentConfigParams []ParameterWithProfilesMap, // all params with configFile parent.config
	nameTopologies map[TopologyName]tc.Topology,
	serverCapabilities map[int]map[ServerCapability]struct{},
	dsRequiredCapabilities map[int]map[ServerCapability]struct{},
	cacheGroups map[tc.CacheGroupName]tc.CacheGroupNullable,
	msoAlgorithm string,
	msoParentRetry string,
	msoUnavailableServerRetryResponses string,
	msoMaxSimpleRetries string,
	msoMaxUnavailableServerRetries string,
	qStringHandling string,
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

	serverPlacement := getTopologyPlacement(tc.CacheGroupName(*server.Cachegroup), topology, cacheGroups)
	if !serverPlacement.InTopology {
		return "", nil // server isn't in topology, no error
	}
	// TODO add Topology/Capabilities to remap.config

	parents, secondaryParents, err := GetTopologyParents(server, ds, servers, parentConfigParams, topology, serverPlacement.IsLastTier, serverCapabilities, dsRequiredCapabilities)
	if err != nil {
		return "", errors.New("getting topology parents for '" + *ds.XMLID + "': skipping! " + err.Error())
	}
	if len(parents) == 0 {
		return "", errors.New("getting topology parents for '" + *ds.XMLID + "': no parents found! skipping! (Does your Topology have a CacheGroup with no servers in it?)")
	}

	txt += ` parent="` + strings.Join(parents, `;`) + `"`
	if len(secondaryParents) > 0 {
		txt += ` secondary_parent="` + strings.Join(secondaryParents, `;`) + `"`
	}
	txt += ` round_robin=` + getTopologyRoundRobin(ds, serverParams, serverPlacement.IsLastTier, msoAlgorithm)
	txt += ` go_direct=` + getTopologyGoDirect(ds, serverPlacement.IsLastTier)
	txt += ` qstring=` + getTopologyQueryString(ds, serverParams, serverPlacement.IsLastTier, msoAlgorithm, qStringHandling)
	txt += getTopologyParentIsProxyStr(serverPlacement.IsLastTier)
	txt += " # topology '" + *ds.Topology + "'"
	txt += "\n"
	return txt, nil
}

func getTopologyParentIsProxyStr(serverIsLastTier bool) string {
	if serverIsLastTier {
		return ` parent_is_proxy=false`
	}
	return ""
}

func getTopologyRoundRobin(
	ds *tc.DeliveryServiceNullable,
	serverParams map[string]string,
	serverIsLastTier bool,
	msoAlgorithm string,
) string {
	roundRobinConsistentHash := "consistent_hash"
	if !serverIsLastTier {
		return roundRobinConsistentHash
	}
	if parentSelectAlg := serverParams[ParentConfigParamAlgorithm]; ds.OriginShield != nil && *ds.OriginShield != "" && strings.TrimSpace(parentSelectAlg) != "" {
		return parentSelectAlg
	}
	if ds.MultiSiteOrigin != nil && *ds.MultiSiteOrigin {
		return msoAlgorithm
	}
	return roundRobinConsistentHash
}

func getTopologyGoDirect(ds *tc.DeliveryServiceNullable, serverIsLastTier bool) string {
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
	ds *tc.DeliveryServiceNullable,
	serverParams map[string]string,
	serverIsLastTier bool,
	msoAlgorithm string,
	qStringHandling string,
) string {
	if serverIsLastTier {
		if ds.MultiSiteOrigin != nil && *ds.MultiSiteOrigin && qStringHandling == "" && msoAlgorithm == tc.AlgorithmConsistentHash && ds.QStringIgnore != nil && tc.QStringIgnore(*ds.QStringIgnore) == tc.QStringIgnoreUseInCacheKeyAndPassUp {
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

func serverParentStr(sv *tc.ServerNullable, params []ParameterWithProfilesMap) (string, error) {
	svParams := serverParentageParams(sv, params)
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
	ds *tc.DeliveryServiceNullable,
	servers []tc.ServerNullable,
	parentConfigParams []ParameterWithProfilesMap, // all params with configFile parent.confign
	topology tc.Topology,
	serverIsLastTier bool,
	serverCapabilities map[int]map[ServerCapability]struct{},
	dsRequiredCapabilities map[int]map[ServerCapability]struct{},
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
	for _, sv := range servers {
		if sv.ID == nil {
			log.Errorln("TO Servers server had nil ID, skipping")
			continue
		} else if sv.Cachegroup == nil {
			log.Errorln("TO Servers server had nil Cachegroup, skipping")
			continue
		}

		if tc.CacheType(sv.Type) != tc.CacheTypeEdge && tc.CacheType(sv.Type) != tc.CacheTypeMid && sv.Type != tc.OriginTypeName {
			continue // only consider edges, mids, and origins in the CacheGroup.
		}
		if !HasRequiredCapabilities(serverCapabilities[*sv.ID], dsRequiredCapabilities[*ds.ID]) {
			continue
		}
		if *sv.Cachegroup == parentCG {
			parentStr, err := serverParentStr(&sv, parentConfigParams)
			if err != nil {
				return nil, nil, errors.New("getting server parent string: " + err.Error())
			}
			if parentStr != "" { // will be empty if server is not_a_parent (possibly other reasons)
				parentStrs = append(parentStrs, parentStr)
			}
		}
		if *sv.Cachegroup == secondaryParentCG {
			parentStr, err := serverParentStr(&sv, parentConfigParams)
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
	ds *tc.DeliveryServiceNullable,
	dsRequiredCapabilities map[int]map[ServerCapability]struct{},
	parentInfos []ParentInfo,
	atsMajorVer int,
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
func getMSOParentStrs(
	ds *tc.DeliveryServiceNullable,
	parentInfos []ParentInfo,
	atsMajorVer int,
	dsRequiredCapabilities map[int]map[ServerCapability]struct{},
	msoAlgorithm string,
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

	// If the ats version supports it and the algorithm is consistent hash, put secondary and non-primary parents into secondary parent group.
	// This will ensure that secondary and tertiary parents will be unused unless all hosts in the primary group are unavailable.

	parents := ""
	secondaryParents := ""

	if atsMajorVer >= 6 && msoAlgorithm == "consistent_hash" && len(secondaryParentStr) > 0 {
		parents = `parent="` + strings.Join(parentInfo, "") + `"`
		secondaryParents = ` secondary_parent="` + secondaryParentStr + `"`
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
	dses []tc.DeliveryServiceNullable,
	serverCapabilities map[int]map[ServerCapability]struct{},
	dsRequiredCapabilities map[int]map[ServerCapability]struct{},
) (map[OriginHost][]CGServer, map[ProfileID]ProfileCache, error) {
	originServers := map[OriginHost][]CGServer{}  // "deliveryServices" in Perl
	profileCaches := map[ProfileID]ProfileCache{} // map[profileID]ProfileCache

	dsIDMap := map[int]tc.DeliveryServiceNullable{}
	for _, ds := range dses {
		if ds.ID == nil {
			return nil, nil, errors.New("delivery services got nil ID!")
		}
		if !ds.Type.IsHTTP() && !ds.Type.IsDNS() {
			continue // skip ANY_MAP, STEERING, etc
		}
		dsIDMap[*ds.ID] = ds
	}

	allDSMap := map[int]tc.DeliveryServiceNullable{} // all DSes for this server, NOT all dses in TO
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
func GetDSOrigins(dses map[int]tc.DeliveryServiceNullable) (map[int]*OriginURI, error) {
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
