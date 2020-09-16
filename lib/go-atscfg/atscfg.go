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
	"encoding/json"
	"errors"
	"net"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
)

const InvalidID = -1

const DefaultATSVersion = "5" // TODO Emulates Perl; change to 6? ATC no longer officially supports ATS 5.

const HeaderCommentDateFormat = "Mon Jan 2 15:04:05 MST 2006"

const ContentTypeTextASCII = `text/plain; charset=us-ascii`

const LineCommentHash = "#"

type TopologyName string
type CacheGroupType string
type ServerCapability string

type ServerInfo struct {
	CacheGroupID                  int
	CacheGroupName                string
	CDN                           tc.CDNName
	CDNID                         int
	DomainName                    string
	HostName                      string
	HTTPSPort                     int
	ID                            int
	IP                            string
	ParentCacheGroupID            int
	ParentCacheGroupType          string
	ProfileID                     ProfileID
	ProfileName                   string
	Port                          int
	SecondaryParentCacheGroupID   int
	SecondaryParentCacheGroupType string
	Type                          string
}

func (s *ServerInfo) IsTopLevelCache() bool {
	return (s.ParentCacheGroupType == tc.CacheGroupOriginTypeName || s.ParentCacheGroupID == InvalidID) &&
		(s.SecondaryParentCacheGroupType == tc.CacheGroupOriginTypeName || s.SecondaryParentCacheGroupID == InvalidID)
}

func MakeCGMap(cgs []tc.CacheGroupNullable) (map[tc.CacheGroupName]tc.CacheGroupNullable, error) {
	cgMap := map[tc.CacheGroupName]tc.CacheGroupNullable{}
	for _, cg := range cgs {
		if cg.Name == nil {
			return nil, errors.New("got cachegroup with nil name!'")
		}
		cgMap[tc.CacheGroupName(*cg.Name)] = cg
	}
	return cgMap, nil
}

type ServerParentCacheGroupData struct {
	ParentID            int
	ParentType          CacheGroupType
	SecondaryParentID   int
	SecondaryParentType CacheGroupType
}

// GetParentCacheGroupData returns the parent CacheGroup IDs and types for the given server.
// Takes a server and a CG map. To create a CGMap from an API CacheGroup slice, use MakeCGMap.
// If server's CacheGroup has no parent or secondary parent, returns InvalidID and "" with no error.
func GetParentCacheGroupData(server *tc.ServerNullable, cgMap map[tc.CacheGroupName]tc.CacheGroupNullable) (ServerParentCacheGroupData, error) {
	if server.Cachegroup == nil || *server.Cachegroup == "" {
		return ServerParentCacheGroupData{}, errors.New("server missing cachegroup")
	} else if server.HostName == nil || *server.HostName == "" {
		return ServerParentCacheGroupData{}, errors.New("server missing hostname")
	}
	serverCG, ok := cgMap[tc.CacheGroupName(*server.Cachegroup)]
	if !ok {
		return ServerParentCacheGroupData{}, errors.New("server '" + *server.HostName + "' cachegroup '" + *server.Cachegroup + "' not found in CacheGroups")
	}

	parentCGID := InvalidID
	parentCGType := ""
	if serverCG.ParentName != nil && *serverCG.ParentName != "" {
		parentCG, ok := cgMap[tc.CacheGroupName(*serverCG.ParentName)]
		if !ok {
			return ServerParentCacheGroupData{}, errors.New("server '" + *server.HostName + "' cachegroup '" + *server.Cachegroup + "' parent '" + *serverCG.ParentName + "' not found in CacheGroups")
		}
		if parentCG.ID == nil {
			return ServerParentCacheGroupData{}, errors.New("got cachegroup '" + *parentCG.Name + "' with nil ID!'")
		}
		parentCGID = *parentCG.ID

		if parentCG.Type == nil {
			return ServerParentCacheGroupData{}, errors.New("got cachegroup '" + *parentCG.Name + "' with nil Type!'")
		}
		parentCGType = *parentCG.Type
	}

	secondaryParentCGID := InvalidID
	secondaryParentCGType := ""
	if serverCG.SecondaryParentName != nil && *serverCG.SecondaryParentName != "" {
		parentCG, ok := cgMap[tc.CacheGroupName(*serverCG.SecondaryParentName)]
		if !ok {
			return ServerParentCacheGroupData{}, errors.New("server '" + *server.HostName + "' cachegroup '" + *server.Cachegroup + "' secondary parent '" + *serverCG.SecondaryParentName + "' not found in CacheGroups")
		}

		if parentCG.ID == nil {
			return ServerParentCacheGroupData{}, errors.New("got cachegroup '" + *parentCG.Name + "' with nil ID!'")
		}
		secondaryParentCGID = *parentCG.ID
		if parentCG.Type == nil {
			return ServerParentCacheGroupData{}, errors.New("got cachegroup '" + *parentCG.Name + "' with nil Type!'")
		}

		secondaryParentCGType = *parentCG.Type
	}

	return ServerParentCacheGroupData{
		ParentID:            parentCGID,
		ParentType:          CacheGroupType(parentCGType),
		SecondaryParentID:   secondaryParentCGID,
		SecondaryParentType: CacheGroupType(secondaryParentCGType),
	}, nil
}

// IsTopLevelCache returns whether server is a top-level cache, as defined by traditional CacheGroup parentage.
// This does not consider Topologies, and should not be used if the Delivery Service being considered has a Topology.
// Takes a ServerParentCacheGroupData, which may be created via GetParentCacheGroupData.
func IsTopLevelCache(s ServerParentCacheGroupData) bool {
	return (s.ParentType == tc.CacheGroupOriginTypeName || s.ParentID == InvalidID) &&
		(s.SecondaryParentType == tc.CacheGroupOriginTypeName || s.SecondaryParentID == InvalidID)
}

func HeaderCommentWithTOVersionStr(name string, nameVersionStr string) string {
	return "# DO NOT EDIT - Generated for " + name + " by " + nameVersionStr + " on " + time.Now().UTC().Format(HeaderCommentDateFormat) + "\n"
}

func GetNameVersionStringFromToolNameAndURL(toolName string, url string) string {
	return toolName + " (" + url + ")"
}

func GenericHeaderComment(name string, toolName string, url string) string {
	return HeaderCommentWithTOVersionStr(name, GetNameVersionStringFromToolNameAndURL(toolName, url))
}

// GetATSMajorVersionFromATSVersion returns the major version of the given profile's package trafficserver parameter.
// The atsVersion is typically a Parameter on the Server's Profile, with the configFile "package" name "trafficserver".
// Returns an error if atsVersion is empty or does not start with an unsigned integer followed by a period or nothing.
func GetATSMajorVersionFromATSVersion(atsVersion string) (int, error) {
	dotPos := strings.Index(atsVersion, ".")
	if dotPos == -1 {
		dotPos = len(atsVersion) // if there's no '.' then assume the whole string is just a major version.
	}
	majorVerStr := atsVersion[:dotPos]

	majorVer, err := strconv.ParseUint(majorVerStr, 10, 64)
	if err != nil {
		return 0, errors.New("unexpected version format, expected e.g. '7.1.2.whatever'")
	}
	return int(majorVer), nil
}

type DeliveryServiceID int
type ProfileID int
type ServerID int

// GenericProfileConfig generates a generic profile config text, from the profile's parameters with the given config file name.
// This does not include a header comment, because a generic config may not use a number sign as a comment.
// If you need a header comment, it can be added manually via ats.HeaderComment, or automatically with WithProfileDataHdr.
func GenericProfileConfig(
	paramData map[string]string, // GetProfileParamData(tx, profileID, fileName)
	separator string,
) string {
	text := ""

	lines := []string{}
	for name, val := range paramData {
		name = trimParamUnderscoreNumSuffix(name)
		lines = append(lines, name+separator+val+"\n")
	}
	sort.Strings(lines)
	text = strings.Join(lines, "")
	return text
}

// trimParamUnderscoreNumSuffix removes any trailing "__[0-9]+" and returns the trimmed string.
func trimParamUnderscoreNumSuffix(paramName string) string {
	underscorePos := strings.LastIndex(paramName, `__`)
	if underscorePos == -1 {
		return paramName
	}
	if _, err := strconv.ParseFloat(paramName[underscorePos+2:], 64); err != nil {
		return paramName
	}
	return paramName[:underscorePos]
}

const ConfigSuffix = ".config"

func GetConfigFile(prefix string, xmlId string) string {
	return prefix + xmlId + ConfigSuffix
}

// topologyIncludesServer returns whether the given topology includes the given server.
func topologyIncludesServer(topology tc.Topology, server *tc.Server) bool {
	for _, node := range topology.Nodes {
		if node.Cachegroup == server.Cachegroup {
			return true
		}
	}
	return false
}

// topologyIncludesServerNullable returns whether the given topology includes the given server.
func topologyIncludesServerNullable(topology tc.Topology, server *tc.ServerNullable) (bool, error) {
	if server.Cachegroup == nil {
		return false, errors.New("server missing Cachegroup")
	}
	for _, node := range topology.Nodes {
		if node.Cachegroup == *server.Cachegroup {
			return true, nil
		}
	}
	return false, nil
}

// TopologyCacheTier is the position of a cache in the topology.
// Note this is the cache tier itself, notwithstanding MSO. So for an MSO service,
// Caches immediately before the origin are the TopologyCacheTierLast, even for MSO.
type TopologyCacheTier string

const (
	TopologyCacheTierFirst   = TopologyCacheTier("first")
	TopologyCacheTierInner   = TopologyCacheTier("inner")
	TopologyCacheTierLast    = TopologyCacheTier("last")
	TopologyCacheTierInvalid = TopologyCacheTier("")
)

// TopologyPlacement contains data about the placement of a server in a topology.
type TopologyPlacement struct {
	// InTopology is whether the server is in the topology at all.
	InTopology bool
	// IsLastTier is whether the server is the last tier in the topology.
	// Note this is different for MSO vs non-MSO. For MSO, the last tier is the Origin. For non-MSO, the last tier is the last cache tier.
	IsLastTier bool
	// CacheTier is the position of the cache in the topology.
	// Note this is whether the cache is the last cache, even if it has parents in the topology who are origins (MSO).
	// Thus, it's possible for a server to be CacheTierLast and not IsLastTier.
	CacheTier TopologyCacheTier
}

// getTopologyPlacement returns information about the cachegroup's placement in the topology.
// - Whether the cachegroup is the last tier in the topology.
// - Whether the cachegroup is in the topology at all.
// - Whether it's the first, inner, or last cache tier before the Origin.
func getTopologyPlacement(cacheGroup tc.CacheGroupName, topology tc.Topology, cacheGroups map[tc.CacheGroupName]tc.CacheGroupNullable) TopologyPlacement {
	serverNode := tc.TopologyNode{}
	serverNodeIndex := -1
	for nodeI, node := range topology.Nodes {
		if node.Cachegroup == string(cacheGroup) {
			serverNode = node
			serverNodeIndex = nodeI
			break
		}
	}
	if serverNode.Cachegroup == "" {
		return TopologyPlacement{InTopology: false}
	}

	topologyNodeHasChildren := false
nodeFor:
	for _, node := range topology.Nodes {
		for _, parent := range node.Parents {
			if parent == serverNodeIndex {
				topologyNodeHasChildren = true
				break nodeFor
			}
		}
	}

	cacheTier := TopologyCacheTierFirst
	if topologyNodeHasChildren {
		cacheTier = TopologyCacheTierInner
	}

	isLastTier := len(serverNode.Parents) == 0

	if isLastTier {
		cacheTier = TopologyCacheTierLast
	}
	// Check if the parent is an Origin, and if so, set to Last
	if cacheTier == TopologyCacheTierInner {
		// TODO extra safety: check other parents, and warn if parents have different types?
		parentI := serverNode.Parents[0]
		if parentI >= len(topology.Nodes) {
			log.Errorln("ATS config generation: topology '" + topology.Name + "' has node with parent larger than nodes size! Config Generation will be malformed!")
		} else {
			parentNode := topology.Nodes[parentI]
			cg, ok := cacheGroups[tc.CacheGroupName(parentNode.Cachegroup)]
			if !ok {
				log.Errorln("ATS config generation: topology '" + topology.Name + "' has node with cachegroup '" + parentNode.Cachegroup + "' that wasn't found in cachegroups! Config Generation will be malformed!")
			} else if cg.Type == nil {
				log.Errorln("ATS config generation: cachegroup '" + parentNode.Cachegroup + "' with nil type! Config Generation will be malformed!")
			} else if *cg.Type == tc.CacheGroupOriginTypeName {
				// this server's parent in the topology is an Origin, so this server is the last cache tier.
				cacheTier = TopologyCacheTierLast
			}
		}
	}
	return TopologyPlacement{InTopology: true, IsLastTier: isLastTier, CacheTier: cacheTier}
}

func MakeTopologyNameMap(topologies []tc.Topology) map[TopologyName]tc.Topology {
	topoNames := map[TopologyName]tc.Topology{}
	for _, to := range topologies {
		topoNames[TopologyName(to.Name)] = to
	}
	return topoNames
}

type ParameterWithProfiles struct {
	tc.Parameter
	ProfileNames []string
}

type ParameterWithProfilesMap struct {
	tc.Parameter
	ProfileNames map[string]struct{}
}

// TCParamsToParamsWithProfiles unmarshals the Profiles that the tc struct doesn't.
func TCParamsToParamsWithProfiles(tcParams []tc.Parameter) ([]ParameterWithProfiles, error) {
	params := make([]ParameterWithProfiles, 0, len(tcParams))
	for _, tcParam := range tcParams {
		param := ParameterWithProfiles{Parameter: tcParam}

		profiles := []string{}
		if err := json.Unmarshal(tcParam.Profiles, &profiles); err != nil {
			return nil, errors.New("unmarshalling JSON from parameter '" + strconv.Itoa(param.ID) + "': " + err.Error())
		}
		param.ProfileNames = profiles
		param.Profiles = nil
		params = append(params, param)
	}
	return params, nil
}

func ParameterWithProfilesToMap(tcParams []ParameterWithProfiles) []ParameterWithProfilesMap {
	params := []ParameterWithProfilesMap{}
	for _, tcParam := range tcParams {
		param := ParameterWithProfilesMap{Parameter: tcParam.Parameter, ProfileNames: map[string]struct{}{}}
		for _, profile := range tcParam.ProfileNames {
			param.ProfileNames[profile] = struct{}{}
		}
		params = append(params, param)
	}
	return params
}

func FilterDSS(dsses []tc.DeliveryServiceServer, dsIDs map[int]struct{}, serverIDs map[int]struct{}) []tc.DeliveryServiceServer {
	// TODO filter only DSes on this server's CDN? Does anything ever needs DSS cross-CDN? Surely not.
	//      Then, we can remove a bunch of config files that filter only DSes on the current cdn.
	filtered := []tc.DeliveryServiceServer{}
	for _, dss := range dsses {
		if dss.Server == nil || dss.DeliveryService == nil {
			continue // TODO warn?
		}
		if len(dsIDs) > 0 {
			if _, ok := dsIDs[*dss.DeliveryService]; !ok {
				continue
			}
		}
		if len(serverIDs) > 0 {
			if _, ok := serverIDs[*dss.Server]; !ok {
				continue
			}
		}
		filtered = append(filtered, dss)
	}
	return filtered
}

// FilterParams filters params and returns only the parameters which match configFile, name, and value.
// If configFile, name, or value is the empty string, it is not filtered.
// Returns a slice of parameters.
func FilterParams(params []tc.Parameter, configFile string, name string, value string, omitName string) []tc.Parameter {
	filtered := []tc.Parameter{}
	for _, param := range params {
		if configFile != "" && param.ConfigFile != configFile {
			continue
		}
		if name != "" && param.Name != name {
			continue
		}
		if value != "" && param.Value != value {
			continue
		}
		if omitName != "" && param.Name == omitName {
			continue
		}
		filtered = append(filtered, param)
	}
	return filtered
}

// ParamsToMap converts a []tc.Parameter to a map[paramName]paramValue.
// If multiple params have the same value, the first one in params will be used an an error will be logged.
// See ParamArrToMultiMap.
func ParamsToMap(params []tc.Parameter) map[string]string {
	mp := map[string]string{}
	for _, param := range params {
		if val, ok := mp[param.Name]; ok {
			if val < param.Value {
				log.Errorln("config generation got multiple parameters for name '" + param.Name + "' - ignoring '" + param.Value + "'")
				continue
			} else {
				log.Errorln("config generation got multiple parameters for name '" + param.Name + "' - ignoring '" + val + "'")
			}
		}
		mp[param.Name] = param.Value
	}
	return mp
}

// ParamArrToMultiMap converts a []tc.Parameter to a map[paramName][]paramValue.
func ParamsToMultiMap(params []tc.Parameter) map[string][]string {
	mp := map[string][]string{}
	for _, param := range params {
		mp[param.Name] = append(mp[param.Name], param.Value)
	}
	return mp
}

// GetServerIPAddress gets the old IPv4 tc.Server.IPAddress from the new tc.Server.Interfaces.
// If no IPv4 address set as a ServiceAddress exists, returns nil
// Malformed addresses are ignored and skipped.
func GetServerIPAddress(sv *tc.ServerNullable) net.IP {
	for _, iFace := range sv.Interfaces {
		for _, addr := range iFace.IPAddresses {
			if !addr.ServiceAddress {
				continue
			}
			if ip := net.ParseIP(addr.Address); ip != nil {
				if ip4 := ip.To4(); ip4 != nil {
					return ip4 // Valid IPv4, return it
				}
				continue // IP, but not v4, keep looking
			}
			// not an IP, try a CIDR
			ip, _, err := net.ParseCIDR(addr.Address)
			if err != nil || ip == nil {
				continue // TODO log? Not a CIDR or IP, keep looking
			}
			// got a valid CIDR
			if ip4 := ip.To4(); ip4 != nil {
				return ip4 // CIDR is V4, return its IP
			}
			continue // valid CIDR, but not v4, keep looking
		}
	}
	return nil
}

// GetServerServiceAddresses returns the first "service" addresses for IPv4 and IPv6 that it finds.
// If an IPv4 or IPv6 "service" address is not found, returns nil for that IP.
// If no IPv4 address set as a ServiceAddress exists, returns nil
// Malformed addresses are ignored and skipped.
func getServiceAddresses(sv *tc.ServerNullable) (net.IP, net.IP) {
	v4 := net.IP(nil)
	v6 := net.IP(nil)
	for _, iFace := range sv.Interfaces {
		for _, addr := range iFace.IPAddresses {
			if !addr.ServiceAddress {
				continue
			}
			if ip := net.ParseIP(addr.Address); ip != nil {
				if ip4 := ip.To4(); ip4 != nil {
					if v4 == nil {
						v4 = ip4
					}
				} else {
					if v6 == nil {
						v6 = ip
					}
				}
				if v4 != nil && v6 != nil {
					return v4, v6
				}
				continue
			}

			// not an IP, try a CIDR
			ip, _, err := net.ParseCIDR(addr.Address)
			if err != nil || ip == nil {
				continue // TODO log error? Not an IP or CIDR
			}
			if ip4 := ip.To4(); ip4 != nil {
				if v4 == nil {
					v4 = ip4
				}
			} else {
				if v6 == nil {
					v6 = ip
				}
			}
			if v4 != nil && v6 != nil {
				return v4, v6
			}
			continue
		}
	}
	return v4, v6
}

// GetTOToolNameAndURL takes the Global Parameters and returns the Traffic Ops Tool Name and URL, as set in the tc.GlobalProfileName Profile 'tm.toolname' and 'tm.url' name Parameters.
func GetTOToolNameAndURL(globalParams []tc.Parameter) (string, string) {
	// TODO move somewhere generic
	toToolName := ""
	toURL := ""
	for _, param := range globalParams {
		if param.Name == "tm.toolname" {
			toToolName = param.Value
		} else if param.Name == "tm.url" {
			toURL = param.Value
		}
		if toToolName != "" && toURL != "" {
			break
		}
	}
	// TODO error here? Perl doesn't.
	if toToolName == "" {
		log.Warnln("Global Parameter tm.toolname not found, config may not be constructed properly!")
	}
	if toURL == "" {
		log.Warnln("Global Parameter tm.url not found, config may not be constructed properly!")
	}
	return toToolName, toURL
}
