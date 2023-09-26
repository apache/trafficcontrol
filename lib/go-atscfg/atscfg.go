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
	"fmt"
	"net"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
)

// InvalidID is used as a placeholder value when determining the ID of a parent
// Cache Group.
//
// This is not strictly invalid; when a Cache Group's Parent's ID is "detected"
// as this value, then it may simply mean that the Cache Group is "top-level",
// and has no parent.
//
// TODO: unexport this? It's highly specific to this package, and isn't used
// anywhere else.
const InvalidID = -1

// DefaultATSVersion is the version of Traffic Server assumed when one is not
// directly specified in a server's Parameters.
//
// TODO: Emulates Perl; change to 6? ATC no longer officially supports ATS 5.
const DefaultATSVersion = "5"

// HeaderCommentDateFormat is an unused format for dates and times.
//
// TODO: Remove this, it's unused.
const HeaderCommentDateFormat = "Mon Jan 2 15:04:05 MST 2006"

// ContentTypeTextASCII is the MIME type for plain text using an ASCII
// character set.
//
// TODO: Move to lib/go-rfc package?
const ContentTypeTextASCII = `text/plain; charset=us-ascii`

// LineCommentHash is the string that denotes the beginning of a line comment,
// for grammars that use the number sign for line comments.
//
// TODO: Unexport? This is only used, it seems, to group multiple usages into a
// single symbol not used outside this package. It's conceivable to want to
// know the line comment string for a given grammar, but it's unlikely that
// anyone importing this package would find an alias for "#" itself helpful,
// with no other attached semantics.
const LineCommentHash = "#"

// ConfigSuffix is a suffix (or "file extension") appended to many
// configuration file names, most notably those for the Header Rewrite Traffic
// Server plugin.
//
// TODO: Unexport? This is only used by the header rewrite configuration file
// generation functions, and doesn't have much practical purpose outside of
// that limited usage.
const ConfigSuffix = ".config"

// TsDefaultRequestHeaderMaxSize is the default maximum request header size
// that will be allowed by ATS unless otherwise specified on a server's
// Profile.
//
// For more information, refer to the ATS documentation for the reconds.config
// configuration option 'proxy.config.http.request_header_max_size'.
const TsDefaultRequestHeaderMaxSize = 131072

// DeliveryServiceID is used internally to represent a Delivery Service ID.
//
// TODO: Unexport? This isn't used outside of this package. Alternatively:
// remove? lib/go-tc and cache-config/* use ints, so all this is doing is
// causing spurious casting back and forth.
type DeliveryServiceID int

// ProfileID is used internally to represent a Profile ID.
//
// TODO: Unexport? This isn't used outside of this package. Alternatively:
// remove? lib/go-tc and cache-config/* use ints, so all this is doing is
// causing spurious casting back and forth.
type ProfileID int

// ServerID is used internally to represent a server ID.
//
// TODO: Unexport? This isn't used outside of this package. Alternatively:
// remove? lib/go-tc and cache-config/* use ints, so all this is doing is
// causing spurious casting back and forth.
type ServerID int

// ProfileName is used internally to represent a Profile Name.
type ProfileName string

// TopologyName is used internally to represent a Topology Name.
//
// Deprecated: github.com/apache/trafficcontrol/lib/go-tc provides the
// an identical type by the same name, but that is the type actually used
// by the Go client and Traffic Ops. Therefore, new code should use that
// type instead of this one.
type TopologyName string

// CacheGroupType is used internally to represent the Name of the Type of a
// Cache Group.
//
// TODO: Unexport? This isn't used outside of this package. Alternatively:
// remove? lib/go-tc and cache-config/* use strings, so all this is doing is
// causing spurious casting back and forth.
type CacheGroupType string

// ServerCapability is used to represent the Name of a Capability that a server
// has OR that a Delivery Service requires.
type ServerCapability string

// Server is a tc.Server for the latest lib/go-tc and traffic_ops/vx-client type.
// This allows atscfg to not have to change the type everywhere it's used, every time ATC changes the base type,
// but to only have to change it here, and the places where breaking symbol changes were made.
type Server tc.ServerV5

// DeliveryService is a tc.DeliveryService for the latest lib/go-tc and traffic_ops/vx-client type.
// This allows atscfg to not have to change the type everywhere it's used, every time ATC changes the base type,
// but to only have to change it here, and the places where breaking symbol changes were made.
type DeliveryService tc.DeliveryServiceV5

// InvalidationJob is a tc.InvalidationJob for the latest lib/go-tc and traffic_ops/vx-client type.
// This allows atscfg to not have to change the type everywhere it's used, every time ATC changes the base type,
// but to only have to change it here, and the places where breaking symbol changes were made.
type InvalidationJob tc.InvalidationJobV4

// ServerUpdateStatus is a tc.ServerUpdateStatus for the latest lib/go-tc and
// traffic_ops/vx-client type. This allows atscfg to not have to change the type
// everywhere it's used, every time ATC changes the base type, but to only have
// to change it here, and the places where breaking symbol changes were made.
type ServerUpdateStatus tc.ServerUpdateStatusV5

//type CDN tc.CDNV5

func ToCDNs(cdns []tc.CDN) []tc.CDNV5 {
	ac := make([]tc.CDNV5, 0, len(cdns))
	for _, cdn := range cdns {
		ac = append(ac, ToCDN(cdn))
	}
	return ac
}

func ToCDN(cdn tc.CDN) tc.CDNV5 {
	return tc.CDNV5{
		DNSSECEnabled: cdn.DNSSECEnabled,
		DomainName:    cdn.DomainName,
		ID:            cdn.ID,
		LastUpdated:   cdn.LastUpdated.Time,
		Name:          cdn.Name,
		TTLOverride:   &cdn.TTLOverride,
	}
}

func ToCacheGroups(cacheGroups []tc.CacheGroupNullable) []tc.CacheGroupNullableV5 {
	ag := make([]tc.CacheGroupNullableV5, 0, len(cacheGroups))
	for _, cg := range cacheGroups {
		ag = append(ag, ToCacheGroup(cg))
	}
	return ag
}

func ToCacheGroup(cacheGroup tc.CacheGroupNullable) tc.CacheGroupNullableV5 {
	return tc.CacheGroupNullableV5{
		ID:                          cacheGroup.ID,
		Name:                        cacheGroup.Name,
		ShortName:                   cacheGroup.ShortName,
		Latitude:                    cacheGroup.Latitude,
		Longitude:                   cacheGroup.Longitude,
		ParentName:                  cacheGroup.ParentName,
		ParentCachegroupID:          cacheGroup.ParentCachegroupID,
		SecondaryParentName:         cacheGroup.SecondaryParentName,
		SecondaryParentCachegroupID: cacheGroup.SecondaryParentCachegroupID,
		FallbackToClosest:           cacheGroup.FallbackToClosest,
		LocalizationMethods:         cacheGroup.LocalizationMethods,
		Type:                        cacheGroup.Type,
		TypeID:                      cacheGroup.TypeID,
		LastUpdated:                 &cacheGroup.LastUpdated.Time,
		Fallbacks:                   cacheGroup.Fallbacks,
	}

}

func ToTopologies(topologies []tc.Topology) []tc.TopologyV5 {
	tg := make([]tc.TopologyV5, 0, len(topologies))
	for _, topology := range topologies {
		tg = append(tg, ToTopology(topology))
	}
	return tg
}

func ToTopology(topology tc.Topology) tc.TopologyV5 {
	nodes := []tc.TopologyNodeV5{}
	for _, n := range topology.Nodes {
		nodes = append(nodes, tc.TopologyNodeV5{
			Cachegroup: n.Cachegroup,
			Parents:    n.Parents,
		})
	}
	return tc.TopologyV5{
		Description: topology.Description,
		Name:        topology.Name,
		Nodes:       nodes,
		LastUpdated: &topology.LastUpdated.Time,
	}
}

// ToDeliveryServices converts a slice of the latest lib/go-tc and traffic_ops/vx-client type to the local alias.
func ToDeliveryServices(dses []tc.DeliveryServiceV5) []DeliveryService {
	ad := make([]DeliveryService, 0, len(dses))
	for _, ds := range dses {
		ad = append(ad, DeliveryService(ds))
	}
	return ad
}

// V50ToDeliveryServices converts a slice of the traffic_ops/v5-client type to the local alias.
func V5ToDeliveryServices(dses []tc.DeliveryServiceV5) []DeliveryService {
	ad := make([]DeliveryService, 0, len(dses))
	for _, ds := range dses {
		if ds.Active != tc.DSActiveStateInactive {
			ad = append(ad, DeliveryService(ds))
		}
	}
	return ad
}

// ToInvalidationJobs converts a slice of the latest lib/go-tc and traffic_ops/vx-client type to the local alias.
func ToInvalidationJobs(jobs []tc.InvalidationJobV4) []InvalidationJob {
	aj := make([]InvalidationJob, 0, len(jobs))
	for _, job := range jobs {
		aj = append(aj, InvalidationJob(job))
	}
	return aj
}

// ToServers converts a slice of the latest lib/go-tc and traffic_ops/vx-client type to the local alias.
func ToServers(servers []tc.ServerV5) []Server {
	as := make([]Server, 0, len(servers))
	for _, sv := range servers {
		as = append(as, Server(sv))
	}
	return as
}

// ToServerUpdateStatuses converts a slice of the latest lib/go-tc and traffic_ops/vx-client type to the local alias.
func ToServerUpdateStatuses(statuses []tc.ServerUpdateStatusV50) []ServerUpdateStatus {
	sus := make([]ServerUpdateStatus, 0, len(statuses))
	for _, st := range statuses {
		sus = append(sus, ServerUpdateStatus(st))
	}
	return sus
}

// CfgFile is all the information necessary to create an ATS config file, including the file name, path, data, and metadata.
// This is provided as a convenience and unified structure for users. The lib/go-atscfg library doesn't actually use or return this. See ATSConfigFileData.
type CfgFile struct {
	Name string
	Path string
	Cfg
}

// Cfg is the data and metadata for an ATS Config File.
//
// This includes the text, the content type (which is necessary for HTTP, multipart, and other things), and the line comment syntax if any.
//
// This is what is generated by the lib/go-atscfg library. Note it does not include the file name or path, which this library doesn't have enough information to return and is not part of generation. That information should be fetched from Traffic Ops, along with the data used to generate config files, or else generated from the machine. See CfgFile.
type Cfg struct {
	Text        string
	ContentType string
	LineComment string
	Secure      bool
	Warnings    []string
}

func makeCGMap(cgs []tc.CacheGroupNullableV5) (map[tc.CacheGroupName]tc.CacheGroupNullableV5, error) {
	cgMap := map[tc.CacheGroupName]tc.CacheGroupNullableV5{}
	for _, cg := range cgs {
		if cg.Name == nil {
			return nil, errors.New("got cachegroup with nil name!'")
		}
		cgMap[tc.CacheGroupName(*cg.Name)] = cg
	}
	return cgMap, nil
}

type serverParentCacheGroupData struct {
	ParentID            int
	ParentType          CacheGroupType
	SecondaryParentID   int
	SecondaryParentType CacheGroupType
}

// getParentCacheGroupData returns the parent CacheGroup IDs and types for the given server.
// Takes a server and a CG map. To create a CGMap from an API CacheGroup slice, use MakeCGMap.
// If server's CacheGroup has no parent or secondary parent, returns InvalidID and "" with no error.
func getParentCacheGroupData(server *Server, cgMap map[tc.CacheGroupName]tc.CacheGroupNullableV5) (serverParentCacheGroupData, error) {
	if server.CacheGroup == "" {
		return serverParentCacheGroupData{}, errors.New("server missing cachegroup")
	} else if server.HostName == "" {
		return serverParentCacheGroupData{}, errors.New("server missing hostname")
	}
	serverCG, ok := cgMap[tc.CacheGroupName(server.CacheGroup)]
	if !ok {
		return serverParentCacheGroupData{}, errors.New("server '" + server.HostName + "' cachegroup '" + server.CacheGroup + "' not found in CacheGroups")
	}

	parentCGID := InvalidID
	parentCGType := ""
	if serverCG.ParentName != nil && *serverCG.ParentName != "" {
		parentCG, ok := cgMap[tc.CacheGroupName(*serverCG.ParentName)]
		if !ok {
			return serverParentCacheGroupData{}, errors.New("server '" + server.HostName + "' cachegroup '" + server.CacheGroup + "' parent '" + *serverCG.ParentName + "' not found in CacheGroups")
		}
		if parentCG.ID == nil {
			return serverParentCacheGroupData{}, errors.New("got cachegroup '" + *parentCG.Name + "' with nil ID!'")
		}
		parentCGID = *parentCG.ID

		if parentCG.Type == nil {
			return serverParentCacheGroupData{}, errors.New("got cachegroup '" + *parentCG.Name + "' with nil Type!'")
		}
		parentCGType = *parentCG.Type
	}

	secondaryParentCGID := InvalidID
	secondaryParentCGType := ""
	if serverCG.SecondaryParentName != nil && *serverCG.SecondaryParentName != "" {
		parentCG, ok := cgMap[tc.CacheGroupName(*serverCG.SecondaryParentName)]
		if !ok {
			return serverParentCacheGroupData{}, errors.New("server '" + server.HostName + "' cachegroup '" + server.CacheGroup + "' secondary parent '" + *serverCG.SecondaryParentName + "' not found in CacheGroups")
		}

		if parentCG.ID == nil {
			return serverParentCacheGroupData{}, errors.New("got cachegroup '" + *parentCG.Name + "' with nil ID!'")
		}
		secondaryParentCGID = *parentCG.ID
		if parentCG.Type == nil {
			return serverParentCacheGroupData{}, errors.New("got cachegroup '" + *parentCG.Name + "' with nil Type!'")
		}

		secondaryParentCGType = *parentCG.Type
	}

	return serverParentCacheGroupData{
		ParentID:            parentCGID,
		ParentType:          CacheGroupType(parentCGType),
		SecondaryParentID:   secondaryParentCGID,
		SecondaryParentType: CacheGroupType(secondaryParentCGType),
	}, nil
}

// isTopLevelCache returns whether server is a top-level cache, as defined by traditional CacheGroup parentage.
// This does not consider Topologies, and should not be used if the Delivery Service being considered has a Topology.
// Takes a ServerParentCacheGroupData, which may be created via GetParentCacheGroupData.
func isTopLevelCache(s serverParentCacheGroupData) bool {
	return (s.ParentType == tc.CacheGroupOriginTypeName || s.ParentID == InvalidID) &&
		(s.SecondaryParentType == tc.CacheGroupOriginTypeName || s.SecondaryParentID == InvalidID)
}

func makeHdrComment(hdrComment string) string {
	return "# " + hdrComment + "\n\n"
}

// GetATSMajorVersionFromATSVersion returns the major version of the given profile's package trafficserver parameter.
// The atsVersion is typically a Parameter on the Server's Profile, with the configFile "package" name "trafficserver".
// Returns an error if atsVersion is empty or does not start with an unsigned integer followed by a period or nothing.
func GetATSMajorVersionFromATSVersion(atsVersion string) (uint, error) {
	dotPos := strings.Index(atsVersion, ".")
	if dotPos == -1 {
		dotPos = len(atsVersion) // if there's no '.' then assume the whole string is just a major version.
	}
	majorVerStr := atsVersion[:dotPos]

	majorVer, err := strconv.ParseUint(majorVerStr, 10, 64)
	if err != nil || majorVer == 0 || majorVer > 9999 {
		return 0, errors.New("unexpected version format '" + majorVerStr + "', expected e.g. '7.1.2.whatever'")
	}
	return uint(majorVer), nil
}

// genericProfileConfig generates a generic profile config text, from the profile's parameters with the given config file name.
// This does not include a header comment, because a generic config may not use a number sign as a comment.
// If you need a header comment, it can be added manually via ats.HeaderComment, or automatically with WithProfileDataHdr.
func genericProfileConfig(
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
func topologyIncludesServerNullable(topology tc.TopologyV5, server *Server) (bool, error) {
	if server.CacheGroup == "" {
		return false, errors.New("server missing Cachegroup")
	}
	for _, node := range topology.Nodes {
		if node.Cachegroup == server.CacheGroup {
			return true, nil
		}
	}
	return false, nil
}

// TopologyCacheTier is the position of a cache in the topology.
// Note this is the cache tier itself, notwithstanding MSO. So for an MSO service,
// Caches immediately before the origin are the TopologyCacheTierLast, even for MSO.
type TopologyCacheTier string

// These are the allowed values for a TopologyCacheTier.
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
	// IsFirstCacheTier is whether the server is the first cache tier. Note if the server is the only tier, both IsFirstCacheTier and IsLastCacheTier may be true.
	IsFirstCacheTier bool
	// IsInnerCacheTier is whether the server is an inner cache tier. This will never be true if IsFirstCacheTier or IsLstCacheTier are true.
	IsInnerCacheTier bool
	// IsLastCacheTier is whether the server is the last cache tier. Note if the server is the only tier, both IsFirstCacheTier and IsLastCacheTier may be true.
	// Note this is distinct from IsLastTier, which will be false for the last cache tier of MSO.
	IsLastCacheTier bool
	// IsLastTier is whether the server is the last tier in the topology.
	// Note this is different than IsLastCacheTier for MSO vs non-MSO. For MSO, the last tier is the Origin. For non-MSO, the last tier is the last cache tier.
	IsLastTier bool
}

// getTopologyPlacement returns information about the cachegroup's placement in the topology, and any error.
// - Whether the cachegroup is the last tier in the topology.
// - Whether the cachegroup is in the topology at all.
// - Whether it's the first, inner, or last cache tier before the Origin.
func getTopologyPlacement(cacheGroup tc.CacheGroupName, topology tc.TopologyV5, cacheGroups map[tc.CacheGroupName]tc.CacheGroupNullableV5, ds *DeliveryService) (TopologyPlacement, error) {
	isMSO := ds.MultiSiteOrigin

	serverNode := tc.TopologyNodeV5{}
	serverNodeIndex := -1
	for nodeI, node := range topology.Nodes {
		if node.Cachegroup == string(cacheGroup) {
			serverNode = node
			serverNodeIndex = nodeI
			break
		}
	}
	if serverNode.Cachegroup == "" {
		return TopologyPlacement{InTopology: false}, nil
	}

	hasChildren := false
nodeFor:
	for _, node := range topology.Nodes {
		for _, parent := range node.Parents {
			if parent == serverNodeIndex {
				hasChildren = true
				break nodeFor
			}
		}
	}

	hasParents := len(serverNode.Parents) > 0

	parentIsOrigin := false
	if hasParents {
		// TODO extra safety: check other parents, and warn if parents have different types?
		parentI := serverNode.Parents[0]
		if parentI >= len(topology.Nodes) {
			return TopologyPlacement{}, errors.New("topology '" + topology.Name + "' has node with parent larger than nodes size! Config Generation will be malformed!")
		}
		parentNode := topology.Nodes[parentI]
		parentCG, ok := cacheGroups[tc.CacheGroupName(parentNode.Cachegroup)]
		if !ok {
			return TopologyPlacement{}, errors.New("topology '" + topology.Name + "' has node with cachegroup '" + parentNode.Cachegroup + "' that wasn't found in cachegroups! Config Generation will be malformed!")
		} else if parentCG.Type == nil {
			return TopologyPlacement{}, errors.New("ATS config generation: cachegroup '" + parentNode.Cachegroup + "' with nil type! Config Generation will be malformed!")
		}
		parentIsOrigin = *parentCG.Type == tc.CacheGroupOriginTypeName
	}

	return TopologyPlacement{
		InTopology:       true,
		IsFirstCacheTier: !hasChildren,
		IsInnerCacheTier: hasChildren && hasParents && !parentIsOrigin,
		IsLastCacheTier:  !hasParents || parentIsOrigin,
		IsLastTier:       !hasParents || (parentIsOrigin && !isMSO), // If the parent CG is an Origin CG, but this DS is not MSO, then ignore the Topology and declare this the last tier
	}, nil
}

func makeTopologyNameMap(topologies []tc.TopologyV5) map[TopologyName]tc.TopologyV5 {
	topoNames := map[TopologyName]tc.TopologyV5{}
	for _, to := range topologies {
		topoNames[TopologyName(to.Name)] = to
	}
	return topoNames
}

// getTopologyDirectChildren returns the cachegroups which are immediate children of the given cachegroup in any topology.
func getTopologyDirectChildren(
	cg tc.CacheGroupName,
	topologies []tc.TopologyV5,
) map[tc.CacheGroupName]struct{} {
	children := map[tc.CacheGroupName]struct{}{}

	for _, topo := range topologies {
		svNodeI := -1
		for nodeI, node := range topo.Nodes {
			if node.Cachegroup == string(cg) {
				svNodeI = nodeI
				break
			}
		}
		if svNodeI < 0 {
			continue // this cg wasn't in the topology
		}
		for _, node := range topo.Nodes {
			for _, parent := range node.Parents {
				if parent == svNodeI {
					children[tc.CacheGroupName(node.Cachegroup)] = struct{}{}
					break
				}
			}
		}
	}
	return children
}

type parameterWithProfiles struct {
	tc.ParameterV5
	ProfileNames []string
}

type parameterWithProfilesMap struct {
	tc.ParameterV5
	ProfileNames map[string]struct{}
}

// tcParamsToParamsWithProfiles unmarshals the Profiles that the tc struct doesn't.
func tcParamsToParamsWithProfiles(tcParams []tc.ParameterV5) ([]parameterWithProfiles, error) {
	params := make([]parameterWithProfiles, 0, len(tcParams))
	for _, tcParam := range tcParams {
		param := parameterWithProfiles{ParameterV5: tcParam}

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

func parameterWithProfilesToMap(tcParams []parameterWithProfiles) []parameterWithProfilesMap {
	params := []parameterWithProfilesMap{}
	for _, tcParam := range tcParams {
		param := parameterWithProfilesMap{ParameterV5: tcParam.ParameterV5, ProfileNames: map[string]struct{}{}}
		for _, profile := range tcParam.ProfileNames {
			param.ProfileNames[profile] = struct{}{}
		}
		params = append(params, param)
	}
	return params
}

func filterDSS(dsses []DeliveryServiceServer, dsIDs map[int]struct{}, serverIDs map[int]struct{}) []DeliveryServiceServer {
	// TODO filter only DSes on this server's CDN? Does anything ever needs DSS cross-CDN? Surely not.
	//      Then, we can remove a bunch of config files that filter only DSes on the current cdn.
	filtered := []DeliveryServiceServer{}
	for _, dss := range dsses {
		if len(dsIDs) > 0 {
			if _, ok := dsIDs[dss.DeliveryService]; !ok {
				continue
			}
		}
		if len(serverIDs) > 0 {
			if _, ok := serverIDs[dss.Server]; !ok {
				continue
			}
		}
		filtered = append(filtered, dss)
	}
	return filtered
}

// filterParams filters params and returns only the parameters which match configFile, name, and value.
// If configFile, name, or value is the empty string, it is not filtered.
// Returns a slice of parameters.
func filterParams(params []tc.ParameterV5, configFile string, name string, value string, omitName string) []tc.ParameterV5 {
	filtered := []tc.ParameterV5{}
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

// paramsToMap converts a []tc.Parameter to a map[paramName]paramValue.
// If multiple params have the same value, the first one in params will be used an an error will be logged.
// Warnings will be returned if any parameters have the same name but different values.
// Returns the parameter map, and any warnings.
// See ParamArrToMultiMap.
func paramsToMap(params []tc.ParameterV5) (map[string]string, []string) {
	warnings := []string{}
	mp := map[string]string{}
	for _, param := range params {
		if val, ok := mp[param.Name]; ok {
			if val < param.Value {
				warnings = append(warnings, "got multiple parameters for name '"+param.Name+"' - ignoring '"+param.Value+"'")
				continue
			} else {
				warnings = append(warnings, "config generation got multiple parameters for name '"+param.Name+"' - ignoring '"+val+"'")
			}
		}
		mp[param.Name] = param.Value
	}
	return mp, warnings
}

// paramArrToMultiMap converts a []tc.Parameter to a map[paramName][]paramValue.
func paramsToMultiMap(params []tc.ParameterV5) map[string][]string {
	mp := map[string][]string{}
	for _, param := range params {
		mp[param.Name] = append(mp[param.Name], param.Value)
	}
	return mp
}

// getServerIPAddress gets the old IPv4 tc.Server.IPAddress from the new tc.Server.Interfaces.
// If no IPv4 address set as a ServiceAddress exists, returns nil
// Malformed addresses are ignored and skipped.
func getServerIPAddress(sv *Server) net.IP {
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

// getServiceAddresses returns the first "service" addresses for IPv4 and IPv6 that it finds.
// If an IPv4 or IPv6 "service" address is not found, returns nil for that IP.
// If no IPv4 address set as a ServiceAddress exists, returns nil
// Malformed addresses are ignored and skipped.
func getServiceAddresses(sv *Server) (net.IP, net.IP) {
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

// getATSMajorVersion takes a config variable of the version, the Server Parameters, and a pointer to warnings to populate.
// This allows callers to use the config variable if it was given, or get the ATS version from the Server Parameters if it wasn't, and add to a warnings variable, in a single line.
//
// If more flexibility is needed, getATSMajorVersionFromParams may be called directly;
// but it should generally be avoided, functions should always take a config variable for the
// ATS version, in case a user wants to manage the ATS package outside ATC.
func getATSMajorVersion(atsMajorVersion uint, serverParams []tc.ParameterV5, warnings *[]string) uint {
	if atsMajorVersion != 0 {
		return atsMajorVersion
	}
	verWarns := []string{}
	atsMajorVersion, verWarns = getATSMajorVersionFromParams(serverParams)
	*warnings = append(*warnings, verWarns...)
	return atsMajorVersion
}

// getATSMajorVersionFromParams should generally not be called directly. Rather, functions should always take the version as a config parameter which may be omitted, and call getATSMajorVersion.
//
// It returns the ATS major version from the config_file 'package' name 'trafficserver' Parameter on the given Server Profile Parameters.
// If no Parameter is found, or the value is malformed, a warning or error is logged and DefaultATSVersion is returned.
// Returns the ATS major version, and any warnings
func getATSMajorVersionFromParams(serverParams []tc.ParameterV5) (uint, []string) {
	warnings := []string{}
	atsVersionParam := ""
	for _, param := range serverParams {
		if param.ConfigFile != "package" || param.Name != "trafficserver" {
			continue
		}
		atsVersionParam = param.Value
		break
	}
	if atsVersionParam == "" {
		warnings = append(warnings, "ATS version Parameter (config_file 'package' name 'trafficserver') not found on Server Profile, using default")
		atsVersionParam = DefaultATSVersion
	}

	atsMajorVer, err := GetATSMajorVersionFromATSVersion(atsVersionParam)
	if err != nil {
		warnings = append(warnings, "getting ATS major version from server Profile Parameter, using default: "+err.Error())
		atsMajorVer, err = GetATSMajorVersionFromATSVersion(DefaultATSVersion)
		if err != nil {
			// should never happen
			warnings = append(warnings, "getting ATS major version from default version! Should never happen! Using 0, config will be malformed! : "+err.Error())
		}
	}
	return atsMajorVer, warnings
}

// getMaxRequestHeaderParam returns the 'CONFIG proxy.config.http.request_header_max_size' if configured in the Server Profile Parameters.
// If the parameter is not configured it will return the traffic server default request header max size.
func getMaxRequestHeaderParam(serverParams []tc.ParameterV5) (int, []string) {
	warnings := []string{}
	globalRequestHeaderMaxSize := TsDefaultRequestHeaderMaxSize
	params, paramWarns := paramsToMap(filterParams(serverParams, RecordsFileName, "", "", "location"))
	warnings = append(warnings, strings.Join(paramWarns, " "))
	if val, ok := params["CONFIG proxy.config.http.request_header_max_size"]; ok {
		size := strings.Fields(val)
		sizeI, err := strconv.Atoi(size[1])
		if err != nil {
			warnings = append(warnings, "Couldn't convert string to int for max_req_header_size")
		} else {
			globalRequestHeaderMaxSize = sizeI
		}
	}
	return globalRequestHeaderMaxSize, warnings
}

// hasRequiredCapabilities returns whether the given caps has all the required capabilities in the given reqCaps.
func hasRequiredCapabilities(caps map[ServerCapability]struct{}, reqCaps map[ServerCapability]struct{}) bool {
	for reqCap, _ := range reqCaps {
		if _, ok := caps[reqCap]; !ok {
			return false
		}
	}
	return true
}

// makeErr takes a list of warnings and an error string, and combines them to a single error.
// Configs typically generate a list of warnings as they go. When an error is encountered, we want to combine the warnings encountered and include them in the returned error message, since they're likely hints as to why the error occurred.
func makeErr(warnings []string, err string) error {
	if len(warnings) == 0 {
		return errors.New(err)
	}
	return errors.New(`(warnings: ` + strings.Join(warnings, `, `) + `) ` + err)
}

// makeErrf is a convenience for formatting errors for makeErr.
// todo also unused, maybe remove?
func makeErrf(warnings []string, format string, v ...interface{}) error {
	return makeErr(warnings, fmt.Sprintf(format, v...))
}

// DeliveryServiceServer is a compact version of DeliveryServiceServer.
// The Delivery Service Servers is massive on large CDNs not using Topologies, compacting it in JSON and dropping the timestamp drastically reduces the size.
// The t3c apps will also drop any DSS from Traffic Ops with null values, which are invalid and useless, to avoid pointers and further reduce size.
type DeliveryServiceServer struct {
	Server          int `json:"s"`
	DeliveryService int `json:"d"`
}

// JobsToInvalidationJobs converts a set of tc.Job structures into their
// equivalent InvalidationJob structures, by calling JobToInvalidationJob on
// each of them in turn.
//
// This will bail on the first error encountered when converting, so other jobs
// afterward may also have problems, but they won't be reported.
func JobsToInvalidationJobs(oldJobs []tc.Job) ([]InvalidationJob, error) {
	jobs := make([]InvalidationJob, len(oldJobs), len(oldJobs))
	err := error(nil)
	for i, oldJob := range oldJobs {
		jobs[i], err = JobToInvalidationJob(oldJob)
		if err != nil {
			return nil, errors.New("converting old tc.Job to tc.InvalidationJob: " + err.Error())
		}
	}
	return jobs, nil
}

// These are the formats for the "lastUpdated" timestamps of Content
// Invalidation Jobs in different Traffic Ops API versions.
const (
	// The format in API versions 4.0 and later.
	JobV4TimeFormat = time.RFC3339Nano
	// The format in API versions prior to 4.0.
	JobLegacyTimeFormat = "2006-01-02 15:04:05-07"
)

// These special suffixes when found on a rule override the default type of a
// Content Invalidation Job.
// Deprecated: these are no longer necessary in Traffic Ops API versions 4.0 and
// later.
const (
	JobLegacyRefetchSuffix = `##REFETCH##`
	JobLegacyRefreshSuffix = `##REFRESH##`
)

// These strings are used to parse legacy representations of Content
// Invalidation Jobs.
// Deprecated: When the Traffic Ops API version 4.0 representation of Content
// Invalidation Jobs becomes standard (i.e. when 3.x API version support is
// dropped) these will have no meaning or purpose.
const (
	JobLegacyParamPrefix = "TTL:"
	JobLegacyParamSuffix = "h"
	JobLegacyKeyword     = "PURGE"
)

// JobToInvalidationJob converts a tc.Job to an InvalidationJob. If the
// conversion fails, an error is returned.
func JobToInvalidationJob(jb tc.Job) (InvalidationJob, error) {
	startTime := tc.Time{}
	if err := json.Unmarshal([]byte(`"`+jb.StartTime+`"`), &startTime); err != nil {
		return InvalidationJob{}, errors.New("unmarshalling time: " + err.Error())
	}
	ttl, err := strconv.Atoi(strings.TrimSuffix(strings.TrimPrefix(jb.Parameters, "TTL:"), "h"))
	if err != nil {
		return InvalidationJob{}, errors.New("unmarshalling ttl: " + err.Error())
	}
	invalType := tc.REFRESH
	if strings.HasSuffix(jb.AssetURL, JobLegacyRefetchSuffix) {
		invalType = tc.REFETCH
	}

	return InvalidationJob{
		AssetURL:         jb.AssetURL,
		CreatedBy:        jb.CreatedBy,
		DeliveryService:  jb.DeliveryService,
		ID:               uint64(jb.ID),
		TTLHours:         uint(ttl),
		InvalidationType: invalType,
		StartTime:        startTime.Time,
	}, nil
}

// FilterServers returns the servers for which filter returns true.
func FilterServers(servers []Server, filter func(sv *Server) bool) []Server {
	// TODO add warning/error feature?
	filteredServers := []Server{}
	for _, sv := range servers {
		if filter(&sv) {
			filteredServers = append(filteredServers, sv)
		}
	}
	return filteredServers
}

// BoolOnOff returns 'on' if b, else 'off'.
// This is a helper func for some ATS config files that use "on" and "off" for boolean values.
func BoolOnOff(b bool) string {
	if b {
		return "on"
	}
	return "off"
}

// GetDSParameters returns the parameters for the given Delivery Service.
func GetDSParameters(
	ds *DeliveryService,
	params []tc.ParameterV5, // from v4-client.GetParameters -> /4.0/parameters
) ([]tc.ParameterV5, error) {
	profileNames := []string{}
	if ds.ProfileName != nil {
		profileNames = append(profileNames, *ds.ProfileName)
	}
	return LayerProfiles(profileNames, params)
}

// GetServerParameters returns the parameters for the given Server, per the Layered Profiles feature.
// See LayerProfiles.
func GetServerParameters(
	server *Server,
	params []tc.ParameterV5, // from v4-client.GetParameters -> /4.0/parameters
) ([]tc.ParameterV5, error) {
	return LayerProfiles(server.Profiles, params)
}

// LayerProfiles takes an ordered list of profile names (presumably from a Server or Delivery Service),
// and the Parameters from Traffic Ops (which includes Profile-Parameters data),
// and layers the parameters according to the ordered list of profiles.
//
// Returns the appropriate parameters for the Server, Delivery Service,
// or other object containing an ordered list of profiles.
func LayerProfiles(
	profileNames []string, // from a Server, Delivery Service, or other object with "layered profiles".
	tcParams []tc.ParameterV5, // from v4-client.GetParameters -> /4.0/parameters
) ([]tc.ParameterV5, error) {
	params, err := tcParamsToParamsWithProfiles(tcParams)
	if err != nil {
		return nil, errors.New("parsing parameters profiles: " + err.Error())
	}
	return layerProfilesFromWith(profileNames, params), nil
}

// layerProfilesFromWith is like LayerProfiles if you already have a []parameterWithProfiles.
func layerProfilesFromWith(profileNames []string, params []parameterWithProfiles) []tc.ParameterV5 {
	paramsMap := parameterWithProfilesToMap(params)
	return layerProfilesFromMap(profileNames, paramsMap)
}

// layerProfilesFromMap is like LayerProfiles if you already have a []parameterWithProfilesMap.
func layerProfilesFromMap(profileNames []string, params []parameterWithProfilesMap) []tc.ParameterV5 {
	// ParamKey is the key for a Parameter, which
	// if there's another Parameter with the same key in a subsequent profile
	// in the ordered list, the last Parameter with this key will be used.
	type ParamKey struct {
		Name       string
		ConfigFile string
	}

	getParamKey := func(pa tc.ParameterV5) ParamKey { return ParamKey{Name: pa.Name, ConfigFile: pa.ConfigFile} }

	allProfileParams := map[string][]tc.ParameterV5{}

	for _, param := range params {
		for profile, _ := range param.ProfileNames {
			allProfileParams[profile] = append(allProfileParams[profile], param.ParameterV5)
		}
	}

	layeredParamMap := map[ParamKey]tc.ParameterV5{}
	// profileNames is ordered, we need to iterate through this backwards
	// because we need subsequent params on other profiles to override previous ones,
	// this will provide the proper "layering" that we want.
	for i := len(profileNames) - 1; i >= 0; i-- {
		profileParams := allProfileParams[profileNames[i]]
		for _, param := range profileParams {
			paramkey := getParamKey(param)
			layeredParamMap[paramkey] = param
		}
	}

	layeredParams := []tc.ParameterV5{}
	for _, param := range layeredParamMap {
		layeredParams = append(layeredParams, param)
	}
	return layeredParams
}

// ServerProfilesMatch returns whether both servers have the same Profiles in the same order,
// and thus will have the same Parameters.
func ServerProfilesMatch(sa *Server, sb *Server) bool {
	return ProfilesMatch(sa.Profiles, sb.Profiles)
}

// ProfilesMatch takes two ordered lists of profile names (such as from Servers or Delivery Services)
// and returns whether they contain the same profiles in the same order,
// and thus whether they will contain the same Parameters.
func ProfilesMatch(pa []string, pb []string) bool {
	if len(pa) != len(pb) {
		return false
	}
	for i, _ := range pa {
		if pa[i] != pb[i] {
			return false
		}
	}
	return true
}

// IsGoDirect checks if this ds type is edge only.
func IsGoDirect(ds DeliveryService) bool {
	return *ds.Type == tc.DSTypeHTTPNoCache.String() || *ds.Type == tc.DSTypeHTTPLive.String() || *ds.Type == tc.DSTypeDNSLive.String()
}
