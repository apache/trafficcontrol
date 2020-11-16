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
	"math"
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

const HeaderRewriteFirstPrefix = HeaderRewritePrefix + "first_"
const HeaderRewriteInnerPrefix = HeaderRewritePrefix + "inner_"
const HeaderRewriteLastPrefix = HeaderRewritePrefix + "last_"

func FirstHeaderRewriteConfigFileName(dsName string) string {
	return HeaderRewriteFirstPrefix + dsName + ConfigSuffix
}

func InnerHeaderRewriteConfigFileName(dsName string) string {
	return HeaderRewriteInnerPrefix + dsName + ConfigSuffix
}

func LastHeaderRewriteConfigFileName(dsName string) string {
	return HeaderRewriteLastPrefix + dsName + ConfigSuffix
}

func MakeTopologyHeaderRewriteDotConfig(
	fileName string,
	server *Server,
	servers []Server,
	deliveryServices []tc.DeliveryServiceNullableV30,
	topologies []tc.Topology,
	serverCapabilities map[int]map[ServerCapability]struct{},
	requiredCapabilities map[int]map[ServerCapability]struct{},
	hdrComment string,
) (Cfg, error) {
	warnings := []string{}

	if server.HostName == nil {
		return Cfg{}, makeErr(warnings, "server missing HostName")
	} else if server.Cachegroup == nil {
		return Cfg{}, makeErr(warnings, "server missing Cachegroup")
	}

	dsName := fileName
	dsName = strings.TrimSuffix(dsName, ConfigSuffix)
	dsName = strings.TrimPrefix(dsName, HeaderRewriteFirstPrefix)
	dsName = strings.TrimPrefix(dsName, HeaderRewriteInnerPrefix)
	dsName = strings.TrimPrefix(dsName, HeaderRewriteLastPrefix)

	tier := TopologyCacheTierInvalid
	switch {
	case strings.HasPrefix(fileName, HeaderRewriteFirstPrefix):
		tier = TopologyCacheTierFirst
	case strings.HasPrefix(fileName, HeaderRewriteInnerPrefix):
		tier = TopologyCacheTierInner
	case strings.HasPrefix(fileName, HeaderRewriteLastPrefix):
		tier = TopologyCacheTierLast
	default:
		return Cfg{}, makeErr(warnings, "topology header rewrite called for unknown tier: '"+fileName+"'")
	}

	ds := tc.DeliveryServiceNullableV30{}
	for _, ids := range deliveryServices {
		if ids.XMLID == nil || *ids.XMLID != dsName {
			continue
		}
		ds = ids
		break
	}
	if ds.ID == nil {
		return Cfg{}, makeErr(warnings, "topology ds '"+dsName+"' not found")
	}

	dsRequiredCapabilities := requiredCapabilities[*ds.ID]

	text := makeHdrComment(hdrComment)

	if ds.Topology == nil || *ds.Topology == "" {
		warnings = append(warnings, "Topology Header Rewrite called for DS '"+*ds.XMLID+"' with no Topology! This should never be called, a DS with no topology should never have a First, Inner, or Last Header Rewrite config in the list of config files! Returning blank config!")
		return Cfg{
			Text:        text,
			ContentType: ContentTypeHeaderRewriteDotConfig,
			LineComment: LineCommentHeaderRewriteDotConfig,
			Warnings:    warnings,
		}, nil
	}

	nameTopologies := makeTopologyNameMap(topologies)
	topology := nameTopologies[TopologyName(*ds.Topology)]
	if topology.Name == "" {
		warnings = append(warnings, "Topology Header Rewrite called for DS '"+*ds.XMLID+"' but its Topology '"+*ds.Topology+"' not found in list of topologies! Returning blank config!")
		return Cfg{
			Text:        text,
			ContentType: ContentTypeHeaderRewriteDotConfig,
			LineComment: LineCommentHeaderRewriteDotConfig,
			Warnings:    warnings,
		}, nil
	}

	headerRewrite := (*string)(nil)
	switch tier {
	case TopologyCacheTierFirst:
		headerRewrite = ds.FirstHeaderRewrite
	case TopologyCacheTierInner:
		headerRewrite = ds.InnerHeaderRewrite
	case TopologyCacheTierLast:
		headerRewrite = ds.LastHeaderRewrite
	default:
		warnings = append(warnings, "Topology Header Rewrite called for DS '"+*ds.XMLID+"' on server '"+*server.HostName+"' got unknown topology cache tier '"+string(tier)+"'! Returning blank config!")
		return Cfg{
			Text:        text,
			ContentType: ContentTypeHeaderRewriteDotConfig,
			LineComment: LineCommentHeaderRewriteDotConfig,
			Warnings:    warnings,
		}, nil
	}

	if tier == TopologyCacheTierLast && ds.MaxOriginConnections != nil && *ds.MaxOriginConnections > 0 {
		lastTierCacheCount, topoWarns := getTopologyDSServerCount(dsRequiredCapabilities, tc.CacheGroupName(*server.Cachegroup), servers, serverCapabilities)
		warnings = append(warnings, topoWarns...)

		maxOriginConnectionsPerServer := int(math.Round(float64(*ds.MaxOriginConnections) / float64(lastTierCacheCount)))
		if maxOriginConnectionsPerServer < 1 {
			maxOriginConnectionsPerServer = 1
		}

		text += "cond %{REMAP_PSEUDO_HOOK}\nset-config proxy.config.http.origin_max_connections " + strconv.Itoa(maxOriginConnectionsPerServer)
		if headerRewrite == nil || *headerRewrite == "" {
			text += " [L]"
		} else {
			text += "\n"
		}
	}

	if headerRewrite != nil && *headerRewrite != "" {
		text += *headerRewrite
	}
	text += "\n"

	return Cfg{
		Text:        text,
		ContentType: ContentTypeHeaderRewriteDotConfig,
		LineComment: LineCommentHeaderRewriteDotConfig,
		Warnings:    warnings,
	}, nil
}

// getTopologyDSServerCount returns the number of servers in cg which will be used to serve ds.
// This should only be used for DSes with Topologies.
// It returns all servers in CG with the Capabilities of ds in cg.
// It will not be the number of servers for Delivery Services not using Topologies, which use DeliveryService-Server assignments instead.
// Returns the server count, and any warnings.
func getTopologyDSServerCount(dsRequiredCapabilities map[ServerCapability]struct{}, cg tc.CacheGroupName, servers []Server, serverCapabilities map[int]map[ServerCapability]struct{}) (int, []string) {
	warnings := []string{}
	count := 0
	for _, sv := range servers {
		if sv.Cachegroup == nil {
			warnings = append(warnings, "Servers had server with nil cachegroup, skipping!")
			continue
		} else if sv.Status == nil {
			warnings = append(warnings, "Servers had server with nil status, skipping!")
			continue
		} else if sv.ID == nil {
			warnings = append(warnings, "Servers had server with nil id, skipping!")
			continue
		}

		if *sv.Cachegroup != string(cg) {
			continue
		}
		if *sv.Status != string(tc.CacheStatusReported) && *sv.Status != string(tc.CacheStatusOnline) {
			continue
		}
		if !hasRequiredCapabilities(serverCapabilities[*sv.ID], dsRequiredCapabilities) {
			continue
		}
		count++
	}
	return count, warnings
}
