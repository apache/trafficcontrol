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

	"github.com/apache/trafficcontrol/lib/go-log"
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
	server *tc.ServerNullable,
	toToolName string, // tm.toolname global parameter (TODO: cache itself?)
	toURL string, // tm.url global parameter (TODO: cache itself?)
	ds tc.DeliveryServiceV30,
	topologies []tc.Topology,
	servers []tc.ServerNullable,
	serverCapabilities map[int]map[ServerCapability]struct{},
	dsRequiredCapabilities map[ServerCapability]struct{},
	tier TopologyCacheTier,
) string {
	if server.HostName == nil {
		return "ERROR: server missing HostName"
	} else if server.Cachegroup == nil {
		return "ERROR: server missing Cachegroup"
	}

	text := GenericHeaderComment(*server.HostName, toToolName, toURL)

	if ds.Topology == nil || *ds.Topology == "" {
		log.Errorln("Config generation: Topology Header Rewrite called for DS '" + *ds.XMLID + "' with no Topology! This should never be called, a DS with no topology should never have a First, Inner, or Last Header Rewrite config in the list of config files! Returning blank config!")
		return text
	}

	nameTopologies := MakeTopologyNameMap(topologies)
	topology := nameTopologies[TopologyName(*ds.Topology)]
	if topology.Name == "" {
		log.Errorln("Config generation: Topology Header Rewrite called for DS '" + *ds.XMLID + "' but its Topology '" + *ds.Topology + "' not found in list of topologies! Returning blank config!")
		return text
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
		log.Errorln("Config generation: Topology Header Rewrite called for DS '" + *ds.XMLID + "' on server '" + *server.HostName + "' got unknown topology cache tier '" + string(tier) + "'! Returning blank config!")
		return text
	}

	if tier == TopologyCacheTierLast && ds.MaxOriginConnections != nil && *ds.MaxOriginConnections > 0 {
		lastTierCacheCount := GetTopologyDSServerCount(dsRequiredCapabilities, tc.CacheGroupName(*server.Cachegroup), servers, serverCapabilities)

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
	return text
}

// GetTopologyDSServerCount returns the number of servers in cg which will be used to serve ds.
// This should only be used for DSes with Topologies.
// It returns all servers in CG with the Capabilities of ds in cg.
// It will not be the number of servers for Delivery Services not using Topologies, which use DeliveryService-Server assignments instead.
func GetTopologyDSServerCount(dsRequiredCapabilities map[ServerCapability]struct{}, cg tc.CacheGroupName, servers []tc.ServerNullable, serverCapabilities map[int]map[ServerCapability]struct{}) int {
	count := 0
	for _, sv := range servers {
		if sv.Cachegroup == nil {
			log.Errorln("TO Servers had nil cachegroup, skipping!")
			continue
		} else if sv.Status == nil {
			log.Errorln("TO Servers had nil status, skipping!")
			continue
		} else if sv.ID == nil {
			log.Errorln("TO Servers had nil id, skipping!")
			continue
		}

		if *sv.Cachegroup != string(cg) {
			continue
		}
		if *sv.Status != string(tc.CacheStatusReported) && *sv.Status != string(tc.CacheStatusOnline) {
			continue
		}
		if !HasRequiredCapabilities(serverCapabilities[*sv.ID], dsRequiredCapabilities) {
			continue
		}
		count++
	}
	return count
}
