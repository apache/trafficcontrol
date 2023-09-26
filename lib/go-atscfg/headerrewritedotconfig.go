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
	"math"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
)

// HeaderRewritePrefix is a prefix on filenames of configuration files for the
// Header Rewrite ATS plugin for edge-tier cache servers. The rest of the
// filename is dependent on the Delivery Service for which the file contains
// rewrite rules.
const HeaderRewritePrefix = "hdr_rw_"

// HeaderRewriteMidPrefix is a prefix on filenames of configuration files for
// the Header Rewrite ATS plugin for mid-tier cache servers. The rest of the
// filename is dependent on the Delivery Service for which the file contains
// rewrite rules.
const HeaderRewriteMidPrefix = "hdr_rw_mid_"

// ContentTypeHeaderRewriteDotConfig is the MIME type of the contents of a
// configuration file for the Header Rewrite ATS plugin.
const ContentTypeHeaderRewriteDotConfig = ContentTypeTextASCII

// LineCommentHeaderRewriteDotConfig is the string used to signify the beginning
// of a line comment in the grammar of configuration files for the Header
// Rewrite ATS plugin.
const LineCommentHeaderRewriteDotConfig = LineCommentHash

// ServiceCategoryHeader is the internal service category header for logging the service category.
// Note this is internal, and will never be set in an HTTP Request or Response by ATS.
const ServiceCategoryHeader = "@CDN-SVC"

// MaxOriginConnectionsNoMax is a value specially interpreted by ATS to mean "no
// maximum origin connections".
//
// TODO: Remove this? It's not used anywhere, not even internally in this
// package.
const MaxOriginConnectionsNoMax = 0 // 0 indicates no limit on origin connections

// HeaderRewriteFirstPrefix is a prefix on filenames of configuration files for
// the Header Rewrite ATS plugin for cache servers at the first (or "edge") tier
// of a Topology. The rest of the filename is dependent on the Delivery Service
// for which the file contains rewrite rules.
const HeaderRewriteFirstPrefix = HeaderRewritePrefix + "first_"

// HeaderRewriteInnerPrefix is a prefix on filenames of configuration files for
// the Header Rewrite ATS plugin for cache servers at any tier between the first
// and last tier of a Topology. The rest of the filename is dependent on the
// Delivery Service for which the file contains rewrite rules.
const HeaderRewriteInnerPrefix = HeaderRewritePrefix + "inner_"

// HeaderRewriteLastPrefix is a prefix on filenames of configuration files for
// the Header Rewrite ATS plugin for cache servers at the last tier of a
// Topology. The rest of the filename is dependent on the Delivery Service for
// which the file contains rewrite rules.
const HeaderRewriteLastPrefix = HeaderRewritePrefix + "last_"

// FirstHeaderRewriteConfigFileName returns the full name of a configuration
// file for the Header Rewrite ATS plugin for a cache server at the first tier
// of a Topology.
//
// The dsName passed in should NOT be the Delivery Service's Display Name, it
// should be its "XMLID".
func FirstHeaderRewriteConfigFileName(dsName string) string {
	return HeaderRewriteFirstPrefix + dsName + ConfigSuffix
}

// InnerHeaderRewriteConfigFileName returns the full name of a configuration
// file for the Header Rewrite ATS plugin for a cache server at the first tier
// of a Topology.
//
// The dsName passed in should NOT be the Delivery Service's Display Name, it
// should be its "XMLID".
func InnerHeaderRewriteConfigFileName(dsName string) string {
	return HeaderRewriteInnerPrefix + dsName + ConfigSuffix
}

// LastHeaderRewriteConfigFileName returns the full name of a configuration
// file for the Header Rewrite ATS plugin for a cache server at the first tier
// of a Topology.
//
// The dsName passed in should NOT be the Delivery Service's Display Name, it
// should be its "XMLID".
func LastHeaderRewriteConfigFileName(dsName string) string {
	return HeaderRewriteLastPrefix + dsName + ConfigSuffix
}

// HeaderRewriteDotConfigOpts contains settings to configure generation options.
type HeaderRewriteDotConfigOpts struct {
	// HdrComment is the header comment to include at the beginning of the file.
	// This should be the text desired, without comment syntax (like # or //). The file's comment syntax will be added.
	// To omit the header comment, pass the empty string.
	HdrComment string

	// ATSMajorVersion is the integral major version of Apache Traffic server,
	// used to generate the proper config for the proper version.
	//
	// If omitted or 0, the major version will be read from the Server's Profile Parameter config file 'package' name 'trafficserver'. If no such Parameter exists, the ATS version will default to 5.
	// This was the old Traffic Control behavior, before the version was specifiable externally.
	//
	ATSMajorVersion uint
}

// MakeHeaderRewriteDotConfig makes the header rewrite file for
// an Edge hdr_rw_ or Mid hdr_rw_mid_ or Topology hdr_rw_{first,inner,last} file,
// as generated by MakeMetaConfigFilesList.
func MakeHeaderRewriteDotConfig(
	fileName string,
	deliveryServices []DeliveryService,
	deliveryServiceServers []DeliveryServiceServer,
	server *Server,
	servers []Server,
	cacheGroupsArr []tc.CacheGroupNullableV5,
	tcServerParams []tc.ParameterV5,
	serverCapabilities map[int]map[ServerCapability]struct{},
	requiredCapabilities map[int]map[ServerCapability]struct{},
	topologiesArr []tc.TopologyV5,
	opt *HeaderRewriteDotConfigOpts,
) (Cfg, error) {
	if opt == nil {
		opt = &HeaderRewriteDotConfigOpts{}
	}
	warnings := []string{}
	if server.CDN == "" {
		return Cfg{}, makeErr(warnings, "this server missing CDNName")
	} else if tc.CacheType(server.Type) == tc.CacheTypeInvalid {
		return Cfg{}, makeErr(warnings, "this server missing Type")
	} else if server.HostName == "" {
		return Cfg{}, makeErr(warnings, "server missing HostName")
	} else if server.CacheGroup == "" {
		return Cfg{}, makeErr(warnings, "server missing Cachegroup")
	}

	cacheGroups, err := makeCGMap(cacheGroupsArr)
	if err != nil {
		return Cfg{}, makeErr(warnings, "making CacheGroup map: "+err.Error())
	}
	topologies := makeTopologyNameMap(topologiesArr)

	// TODO verify prefix and suffix? Perl doesn't
	dsName := fileName
	dsName = strings.TrimSuffix(dsName, ConfigSuffix)
	dsName = strings.TrimPrefix(dsName, HeaderRewriteFirstPrefix)
	dsName = strings.TrimPrefix(dsName, HeaderRewriteInnerPrefix)
	dsName = strings.TrimPrefix(dsName, HeaderRewriteLastPrefix)
	dsName = strings.TrimPrefix(dsName, HeaderRewriteMidPrefix)
	dsName = strings.TrimPrefix(dsName, HeaderRewritePrefix)

	ds := &DeliveryService{}
	for _, ids := range deliveryServices {
		if ids.Active == tc.DSActiveStateInactive {
			continue
		}
		if ids.XMLID == "" {
			warnings = append(warnings, "deliveryServices had DS with nil xmlId (name)")
			continue
		}
		if ids.XMLID != dsName {
			continue
		}
		ds = &ids
		break
	}
	if ds.ID == nil {
		return Cfg{}, makeErr(warnings, "ds '"+dsName+"' not found")
	} else if ds.CDNName == nil {
		return Cfg{}, makeErr(warnings, "ds '"+dsName+"' missing cdn")
	}

	if ds.Topology != nil && *ds.Topology != "" && headerRewriteTopologyTier(fileName) == TopologyCacheTierInvalid {
		// write a blank file, rather than an error. Because this usually means a bad location parameter,
		// we don't want to break all of config generation during the migration to Topologies
		warnings = append(warnings, "header rewrite file '"+fileName+"' for non-topology, but delivery service has a Topology. Do you have a location Parameter that needs deleted? Writing blank file!")
		return Cfg{
			Text:        "",
			ContentType: ContentTypeHeaderRewriteDotConfig,
			LineComment: LineCommentHeaderRewriteDotConfig,
			Warnings:    warnings,
		}, nil
	}

	topology := tc.TopologyV5{}
	if ds.Topology != nil && *ds.Topology != "" {
		topology = topologies[TopologyName(*ds.Topology)]
		if topology.Name == "" {
			return Cfg{}, makeErr(warnings, "DS "+ds.XMLID+" topology '"+*ds.Topology+"' not found in Topologies!")
		}
	}

	atsRqstMaxHdrSize, paramWarns := getMaxRequestHeaderParam(tcServerParams)
	warnings = append(warnings, paramWarns...)

	atsMajorVersion := getATSMajorVersion(opt.ATSMajorVersion, tcServerParams, &warnings)

	assignedTierPeers, assignWarns := getAssignedTierPeers(server, ds, topology, servers, deliveryServiceServers, cacheGroupsArr, serverCapabilities, requiredCapabilities[*ds.ID])
	warnings = append(warnings, assignWarns...)

	dsOnlinePeerCount := 0
	for _, sv := range assignedTierPeers {
		if sv.Status == "" {
			warnings = append(warnings, "got server with nil status! skipping!")
			continue
		}
		if tc.CacheStatus(sv.Status) == tc.CacheStatusReported || tc.CacheStatus(sv.Status) == tc.CacheStatusOnline {
			dsOnlinePeerCount++
		}
	}
	numLastTierServers := dsOnlinePeerCount

	serverIsLastTier, err := headerRewriteServerIsLastTier(server, ds, fileName, cacheGroups, topology)
	if err != nil {
		return Cfg{}, makeErr(warnings, "getting header rewrite tier from delivery service: "+err.Error())
	}

	headerRewriteTxt, err := getTierHeaderRewrite(server, ds, fileName)
	if err != nil {
		return Cfg{}, makeErr(warnings, "getting header rewrite text from delivery service: "+err.Error())
	}

	text := makeHdrComment(opt.HdrComment)

	// Add the TC directives (which may be empty).
	// NOTE!! Custom TC injections MUST NOT EVER have a `[L]`. Doing so will break custom header rewrites!
	// NOTE!! The TC injections MUST be come before custom rewrites (EdgeHeaderRewrite, InnerHeaderRewrite, etc).
	//        If they're placed after, custom rewrites with [L] directives will result in them being applied inconsistently and incorrectly.
	text += makeATCHeaderRewriteDirectives(ds, headerRewriteTxt, serverIsLastTier, numLastTierServers, atsMajorVersion, atsRqstMaxHdrSize)

	if headerRewriteTxt != nil && *headerRewriteTxt != "" {
		hdrRw := returnRe.ReplaceAllString(*headerRewriteTxt, "\n")
		hdrRw = strings.TrimSpace(hdrRw)
		text += `
` + hdrRw + `
`
	}

	return Cfg{
		Text:        text,
		ContentType: ContentTypeHeaderRewriteDotConfig,
		LineComment: LineCommentHeaderRewriteDotConfig,
		Warnings:    warnings,
	}, nil
}

// headerRewriteServerIsLastTier is whether the server is the last tier for the delivery service of this header rewrite.
// This should NOT be abstracted into a function that could be used by any other config.
// This is whether the server is the last tier for this header rewrite. Which may not be true for other rewrites or configs.
func headerRewriteServerIsLastTier(server *Server, ds *DeliveryService, fileName string, cacheGroups map[tc.CacheGroupName]tc.CacheGroupNullableV5, topology tc.TopologyV5) (bool, error) {
	if ds.Topology != nil {
		return headerRewriteTopologyTier(fileName) == TopologyCacheTierLast, nil
		// serverPlacement, err := getTopologyPlacement(tc.CacheGroupName(*server.Cachegroup), topology, cacheGroups, ds)
		// fmt.Printf("DEBUG ds '%v' topo placement %+v\n", *ds.XMLID, serverPlacement)
		// if err != nil {
		// 	return false, errors.New("getting topology placement: " + err.Error())
		// }
		// if !serverPlacement.InTopology {
		// 	return false, errors.New("server not in topology")
		// }
		// return serverPlacement.IsLastCacheTier, nil
	}

	serverIsMid := serverIsMid(server)
	dsUsesMids := tc.DSType(*ds.Type).UsesMidCache()
	dssIsLastTier := (!serverIsMid && !dsUsesMids) || (serverIsMid && dsUsesMids)
	return dssIsLastTier, nil

}

// headerRewriteTopologyTier returns the topology tier of this header rewrite file,
// or TopologyCacheTierInvalid if the file is not for a topology header rewrite.
func headerRewriteTopologyTier(fileName string) TopologyCacheTier {
	switch {
	case strings.HasPrefix(fileName, HeaderRewriteFirstPrefix):
		return TopologyCacheTierFirst
	case strings.HasPrefix(fileName, HeaderRewriteInnerPrefix):
		return TopologyCacheTierInner
	case strings.HasPrefix(fileName, HeaderRewriteLastPrefix):
		return TopologyCacheTierLast
	default:
		return TopologyCacheTierInvalid
	}
}

// getTierHeaderRewrite returns the ds MidHeaderRewrite if server is a MID, else the EdgeHeaderRewrite.
// Does not consider Topologies.
// May return nil, if the tier's HeaderRewrite (Edge or Mid) is nil.
func getTierHeaderRewrite(server *Server, ds *DeliveryService, fileName string) (*string, error) {
	if ds.Topology != nil {
		return getTierHeaderRewriteTopology(server, ds, fileName)
	}
	if serverIsMid(server) {
		return ds.MidHeaderRewrite, nil
	}
	return ds.EdgeHeaderRewrite, nil
}

func getTierHeaderRewriteTopology(server *Server, ds *DeliveryService, fileName string) (*string, error) {
	tier := headerRewriteTopologyTier(fileName)
	switch tier {
	case TopologyCacheTierFirst:
		return ds.FirstHeaderRewrite, nil
	case TopologyCacheTierInner:
		return ds.InnerHeaderRewrite, nil
	case TopologyCacheTierLast:
		return ds.LastHeaderRewrite, nil
	default:
		return nil, errors.New("Topology Header Rewrite called for DS '" + ds.XMLID + "' on server '" + server.HostName + "' file '" + fileName + "' had unknown topology cache tier '" + string(tier) + "'!")
	}
}

// serverIsMid returns true if server's type is Mid. The server.Type MUST NOT be nil. Does not consider Topologies.
func serverIsMid(server *Server) bool {
	return strings.HasPrefix(server.Type, tc.MidTypePrefix)
}

// getAssignedTierPeers returns all edges assigned to the DS if server is an edge,
// or all mids if server is a mid,
// or all the servers at the same tier, if the Delivery Service uses Topologies.
// Note this returns all servers of any status, not just ONLINE or REPORTED servers.
// Returns the list of assigned peers, and any warnings.
func getAssignedTierPeers(
	server *Server,
	ds *DeliveryService,
	topology tc.TopologyV5,
	servers []Server,
	deliveryServiceServers []DeliveryServiceServer,
	cacheGroups []tc.CacheGroupNullableV5,
	serverCapabilities map[int]map[ServerCapability]struct{},
	dsRequiredCapabilities map[ServerCapability]struct{},
) ([]Server, []string) {
	if ds.Topology != nil {
		return getTopologyTierServers(ds, dsRequiredCapabilities, tc.CacheGroupName(server.CacheGroup), topology, cacheGroups, servers, serverCapabilities)
	}
	if serverIsMid(server) {
		return getAssignedMids(server, ds, servers, deliveryServiceServers, cacheGroups)
	}
	return getAssignedEdges(ds, server, servers, deliveryServiceServers)
}

// getAssignedEdges returns all EDGE caches assigned to ds via DeliveryService-Service. Does not consider Topologies.
// Note this returns all servers of any status, not just ONLINE or REPORTED servers.
// Returns the list of assigned servers, and any warnings.
func getAssignedEdges(
	ds *DeliveryService,
	server *Server,
	servers []Server,
	deliveryServiceServers []DeliveryServiceServer,
) ([]Server, []string) {
	warnings := []string{}

	dsServers := filterDSS(deliveryServiceServers, map[int]struct{}{*ds.ID: {}}, nil)

	dsServerIDs := map[int]struct{}{}
	for _, dss := range dsServers {
		if dss.DeliveryService != *ds.ID {
			continue
		}
		dsServerIDs[dss.Server] = struct{}{}
	}

	assignedEdges := []Server{}
	for _, sv := range servers {
		if sv.CDN == "" {
			warnings = append(warnings, "servers had server with missing cdnName, skipping!")
			continue
		}
		if sv.ID == 0 {
			warnings = append(warnings, "servers had server with missing id, skipping!")
			continue
		}
		if sv.CDN != *ds.CDNName {
			continue
		}
		if _, ok := dsServerIDs[sv.ID]; !ok && ds.Topology == nil {
			continue
		}
		if ds != nil && ds.Regional && sv.CacheGroup != server.CacheGroup {
			continue
		}
		assignedEdges = append(assignedEdges, sv)
	}
	return assignedEdges, warnings
}

// getAssignedMids returns all MID caches with a child EDGE assigned to ds via DeliveryService-Service. Does not consider Topologies.
// Note this returns all servers of any status, not just ONLINE or REPORTED servers.
// Returns the list of assigned servers, and any warnings.
func getAssignedMids(
	server *Server,
	ds *DeliveryService,
	servers []Server,
	deliveryServiceServers []DeliveryServiceServer,
	cacheGroups []tc.CacheGroupNullableV5,
) ([]Server, []string) {
	warnings := []string{}
	assignedServers := map[int]struct{}{}
	for _, dss := range deliveryServiceServers {
		if dss.DeliveryService != *ds.ID {
			continue
		}
		assignedServers[dss.Server] = struct{}{}
	}

	serverCGs := map[tc.CacheGroupName]struct{}{}
	for _, sv := range servers {
		if sv.CDN == "" {
			warnings = append(warnings, "TO returned Servers sv with missing CDNName, skipping!")
			continue
		} else if sv.ID == 0 {
			warnings = append(warnings, "TO returned Servers sv with missing ID, skipping!")
			continue
		} else if sv.Status == "" {
			warnings = append(warnings, "TO returned Servers sv with missing Status, skipping!")
			continue
		} else if sv.CacheGroup == "" {
			warnings = append(warnings, "TO returned Servers sv with missing Cachegroup, skipping!")
			continue
		}

		if sv.CDN != server.CDN {
			continue
		}
		if _, ok := assignedServers[sv.ID]; !ok && (ds.Topology == nil || *ds.Topology == "") {
			continue
		}
		if tc.CacheStatus(sv.Status) != tc.CacheStatusReported && tc.CacheStatus(sv.Status) != tc.CacheStatusOnline {
			continue
		}
		serverCGs[tc.CacheGroupName(sv.CacheGroup)] = struct{}{}
	}

	parentCGs := map[string]struct{}{} // names of cachegroups which are parent cachegroups of the cachegroup of any edge assigned to the given DS
	for _, cg := range cacheGroups {
		if cg.Name == nil {
			warnings = append(warnings, "cachegroups had cachegroup with nil name, skipping!")
			continue
		}
		if cg.ParentName == nil {
			continue // this is normal for top-level cachegroups
		}
		if _, ok := serverCGs[tc.CacheGroupName(*cg.Name)]; !ok {
			continue
		}
		parentCGs[*cg.ParentName] = struct{}{}
	}

	assignedMids := []Server{}
	for _, sv := range servers {
		if sv.CDN == "" {
			warnings = append(warnings, "TO returned Servers server with missing CDNName, skipping!")
			continue
		}
		if sv.CacheGroup == "" {
			warnings = append(warnings, "TO returned Servers server with missing Cachegroup, skipping!")
			continue
		}
		if sv.CDN != *ds.CDNName {
			continue
		}
		if _, ok := parentCGs[sv.CacheGroup]; !ok {
			continue
		}
		if ds != nil && ds.Regional && sv.CacheGroup != server.CacheGroup {
			continue
		}
		assignedMids = append(assignedMids, sv)
	}

	return assignedMids, warnings
}

// getTopologyTierServers returns the servers in the same tier as cg which will be used to serve ds.
// This should only be used for DSes with Topologies.
// It returns all servers in with the Capabilities of ds in the same tier as cg.
// Returns the servers, and any warnings.
func getTopologyTierServers(ds *DeliveryService, dsRequiredCapabilities map[ServerCapability]struct{}, cg tc.CacheGroupName, topology tc.TopologyV5, cacheGroups []tc.CacheGroupNullableV5, servers []Server, serverCapabilities map[int]map[ServerCapability]struct{}) ([]Server, []string) {
	warnings := []string{}
	topoServers := []Server{}
	cacheGroupsInSameTier := getCachegroupsInSameTopologyTier(string(cg), cacheGroups, topology)
	for _, sv := range servers {
		if sv.CacheGroup == "" {
			warnings = append(warnings, "Servers had server with nil cachegroup, skipping!")
			continue
		} else if sv.ID == 0 {
			warnings = append(warnings, "Servers had server with nil id, skipping!")
			continue
		}

		if !cacheGroupsInSameTier[sv.CacheGroup] {
			continue
		}

		if !hasRequiredCapabilities(serverCapabilities[sv.ID], dsRequiredCapabilities) {
			continue
		}
		if ds != nil && ds.Regional && sv.CacheGroup != string(cg) {
			continue
		}
		topoServers = append(topoServers, sv)
	}
	return topoServers, warnings
}

func getCachegroupsInSameTopologyTier(cg string, cacheGroups []tc.CacheGroupNullableV5, topology tc.TopologyV5) map[string]bool {
	cacheGroupMap := make(map[string]tc.CacheGroupNullableV5)
	originCacheGroups := make(map[string]bool)
	for _, cg := range cacheGroups {
		if cg.Name == nil || cg.Type == nil {
			continue
		}
		cacheGroupMap[*cg.Name] = cg
		if *cg.Type == tc.CacheGroupOriginTypeName {
			originCacheGroups[*cg.Name] = true
		}
	}
	originNodes := make(map[int]bool)
	nodeIndex := -1
	for i, node := range topology.Nodes {
		if node.Cachegroup == cg {
			nodeIndex = i
		}
		if originCacheGroups[node.Cachegroup] {
			originNodes[i] = true
		}
	}
	nodesWithSameOriginDistances := getNodesWithSameOriginDistances(nodeIndex, topology, originNodes)
	cacheGroupsInSameTopologyTier := make(map[string]bool)
	for _, nodeI := range nodesWithSameOriginDistances {
		cacheGroupsInSameTopologyTier[topology.Nodes[nodeI].Cachegroup] = true
	}
	return cacheGroupsInSameTopologyTier
}

func getNodesWithSameOriginDistances(nodeIndex int, topology tc.TopologyV5, originNodes map[int]bool) []int {
	originDistances := make(map[int]int)
	nodeDistance := -1
	for i := range topology.Nodes {
		d := getOriginDistance(topology, i, originNodes, originDistances)
		if nodeIndex == i {
			nodeDistance = d
		}
	}
	sameDistances := make([]int, 0)
	for i, d := range originDistances {
		if d == nodeDistance {
			sameDistances = append(sameDistances, i)
		}
	}
	return sameDistances
}

func getOriginDistance(topology tc.TopologyV5, nodeIndex int, originNodes map[int]bool, originDistances map[int]int) int {
	if originDistance, ok := originDistances[nodeIndex]; ok {
		return originDistance
	}
	parents := topology.Nodes[nodeIndex].Parents
	if len(parents) == 0 {
		originDistances[nodeIndex] = 1
		return originDistances[nodeIndex]
	}
	for _, p := range parents {
		if originNodes[p] {
			originDistances[nodeIndex] = 1
			return originDistances[nodeIndex]
		}
	}
	originDistances[nodeIndex] = 1 + getOriginDistance(topology, parents[0], originNodes, originDistances)
	return originDistances[nodeIndex]
}

var returnRe = regexp.MustCompile(`\s*__RETURN__\s*`)

// makeATCHeaderRewriteDirectives returns the Header Rewrite text for all per-Delivery-Service Traffic Control directives, such as MaxOriginConnections and ServiceCategory.
// These should be prepended to any custom Header Rewrites, in order to prevent [L] directive errors.
// The returned text may be empty, if no directives are configured.
//
// NOTE!! Custom TC injections MUST NOT ever have a `[L]`. Doing so will break custom header rewrites!
// NOTE!! The TC injections MUST be come before custom rewrites (EdgeHeaderRewrite, InnerHeaderRewrite, etc).
//
//	If they're placed after, custom rewrites with [L] directives will result in them being applied inconsistently and incorrectly.
//
// The headerRewriteTxt is the custom header rewrite from the Delivery Service. This should be used for any logic that depends on it. The various header rewrite fields (EdgeHeaderRewrite, InnerHeaderRewrite, etc should never be used inside this function, since this function doesn't know what tier the server is at. This function should not insert the headerRewriteText, but may use it to make decisions about what to insert.
func makeATCHeaderRewriteDirectives(ds *DeliveryService, headerRewriteTxt *string, serverIsLastTier bool, numLastTierServers int, atsMajorVersion uint, atsRqstMaxHdrSize int) string {
	return makeATCHeaderRewriteDirectiveMaxOriginConns(ds, headerRewriteTxt, serverIsLastTier, numLastTierServers, atsMajorVersion) +
		makeATCHeaderRewriteDirectiveServiceCategoryHdr(ds, headerRewriteTxt) + makeATCHeaderRewriteDirectiveMaxRequestHeaderSize(ds, serverIsLastTier, atsRqstMaxHdrSize)
}

// makeATCHeaderRewriteDirectiveMaxOriginConns generates the Max Origin Connections header rewrite text, which may be empty.
func makeATCHeaderRewriteDirectiveMaxOriginConns(ds *DeliveryService, headerRewriteTxt *string, serverIsLastTier bool, numLastTierServers int, atsMajorVersion uint) string {
	if !serverIsLastTier ||
		(ds.MaxOriginConnections == nil || *ds.MaxOriginConnections < 1) ||
		numLastTierServers < 1 {
		return ""
	}

	maxOriginConnectionsPerServer := int(math.Round(float64(*ds.MaxOriginConnections) / float64(numLastTierServers)))
	if maxOriginConnectionsPerServer < 1 {
		maxOriginConnectionsPerServer = 1
	}

	if atsMajorVersion < 9 {
		return `
cond %{REMAP_PSEUDO_HOOK}
set-config proxy.config.http.origin_max_connections ` + strconv.Itoa(maxOriginConnectionsPerServer) + `
`
	}

	// if the DS doesn't specify a match, use host. This will make ATS treat different hostnames (but not IPs) as different, for max origin connections.
	// Which is what we want. It's common for different DSes to CNAME the same origin, such as a cloud provider.
	// In that case, we want to give each hostname=remap=deliveryservice its own max
	maybeMatch := ""
	if headerRewriteTxt == nil || !strings.Contains(*headerRewriteTxt, `proxy.config.http.per_server.connection.match`) {
		maybeMatch += `set-config proxy.config.http.per_server.connection.match host
`
	}
	return `
cond %{REMAP_PSEUDO_HOOK}
` + maybeMatch + `set-config proxy.config.http.per_server.connection.max ` + strconv.Itoa(maxOriginConnectionsPerServer) + `
`
}

func makeATCHeaderRewriteDirectiveServiceCategoryHdr(ds *DeliveryService, headerRewriteTxt *string) string {
	if (ds.ServiceCategory == nil || *ds.ServiceCategory == "") ||
		(headerRewriteTxt != nil && strings.Contains(*headerRewriteTxt, ServiceCategoryHeader)) { // if the custom header rewrite already contains the service category header, don't add another one
		return ""
	}
	// Escape the ServiceCategory, which is user input, to prevent exploits. No Delivery Service should be able to break the CDN.
	// This is more conservative than necessary, but Go doesn't have a HeaderEscape, and valid path characters are a subset of valid header values.
	escapedServiceCategory := url.PathEscape(*ds.ServiceCategory)
	return `
cond %{REMAP_PSEUDO_HOOK}
set-header ` + ServiceCategoryHeader + ` "` + ds.XMLID + `|` + escapedServiceCategory + `"
`
}

func makeATCHeaderRewriteDirectiveMaxRequestHeaderSize(ds *DeliveryService, serverIsLastTier bool, atsRqstMaxHdrSize int) string {
	if serverIsLastTier || ds.MaxRequestHeaderBytes == nil || *ds.MaxRequestHeaderBytes < 1 {
		return ""
	}
	hdrTxt := "cond %{REMAP_PSEUDO_HOOK}\ncond % cqhl > " + strconv.Itoa(*ds.MaxRequestHeaderBytes) + "\nset-status 400"
	warnTxt := "#TO Max Request Header Size: " + strconv.Itoa(*ds.MaxRequestHeaderBytes) +
		",is larger than or equal to the global setting of " + strconv.Itoa(atsRqstMaxHdrSize) + ", header rw will be ignored.\n"
	if *ds.MaxRequestHeaderBytes >= atsRqstMaxHdrSize {
		return warnTxt + hdrTxt
	} else {
		return hdrTxt
	}
}
