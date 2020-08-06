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
	"sort"
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

const HostingConfigFileName = `hosting.config`
const HostingConfigParamConfigFile = `storage.config`
const ContentTypeHostingDotConfig = ContentTypeTextASCII
const LineCommentHostingDotConfig = LineCommentHash

const ParamDrivePrefix = "Drive_Prefix"
const ParamRAMDrivePrefix = "RAM_Drive_Prefix"

const ServerHostingDotConfigMidIncludeInactive = false
const ServerHostingDotConfigEdgeIncludeInactive = true

func MakeHostingDotConfig(
	server *tc.Server,
	toToolName string, // tm.toolname global parameter (TODO: cache itself?)
	toURL string, // tm.url global parameter (TODO: cache itself?)
	params map[string]string, // map[name]value - config file should always be storage.config
	dses []tc.DeliveryServiceNullable,
	topologies []tc.Topology,
) string {
	text := GenericHeaderComment(server.HostName, toToolName, toURL)

	nameTopologies := MakeTopologyNameMap(topologies)

	lines := []string{}
	if _, ok := params[ParamRAMDrivePrefix]; ok {
		nextVolume := 1
		if _, ok := params[ParamDrivePrefix]; ok {
			diskVolume := nextVolume
			text += `# TRAFFIC OPS NOTE: volume ` + strconv.Itoa(diskVolume) + ` is the Disk volume` + "\n"
			nextVolume++
		}
		ramVolume := nextVolume
		text += `# TRAFFIC OPS NOTE: volume ` + strconv.Itoa(ramVolume) + ` is the RAM volume` + "\n"

		seenOrigins := map[string]struct{}{}
		for _, ds := range dses {
			if ds.OrgServerFQDN == nil || ds.XMLID == nil || ds.Active == nil {
				continue // TODO warn?
			}

			origin := *ds.OrgServerFQDN
			if _, ok := seenOrigins[origin]; ok {
				continue
			}

			if ds.Topology != nil && *ds.Topology != "" {
				topology, hasTopology := nameTopologies[TopologyName(*ds.Topology)]
				if hasTopology && !topologyIncludesServer(topology, server) {
					continue
				}
			}

			seenOrigins[origin] = struct{}{}
			origin = strings.TrimPrefix(origin, `http://`)
			origin = strings.TrimPrefix(origin, `https://`)
			lines = append(lines, `hostname=`+origin+` volume=`+strconv.Itoa(ramVolume)+"\n")
		}
	}
	diskVolume := 1 // note this will actually be the RAM (RAM_Drive_Prefix) volume if there is no Drive_Prefix parameter.

	lines = append(lines, `hostname=*   volume=`+strconv.Itoa(diskVolume)+"\n")

	sort.Strings(lines)
	text += strings.Join(lines, "")
	return text
}
