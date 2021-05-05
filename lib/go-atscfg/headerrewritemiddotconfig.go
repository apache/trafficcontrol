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
	"regexp"
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

const HeaderRewriteMidPrefix = "hdr_rw_mid_"

func MakeHeaderRewriteMidDotConfig(
	fileName string,
	deliveryServices []DeliveryService,
	deliveryServiceServers []tc.DeliveryServiceServer,
	server *Server,
	servers []Server,
	cacheGroups []tc.CacheGroupNullable,
	hdrComment string,
) (Cfg, error) {
	warnings := []string{}
	if server.CDNName == nil {
		return Cfg{}, makeErr(warnings, "this server missing CDNName")
	}

	dsName := strings.TrimSuffix(strings.TrimPrefix(fileName, HeaderRewriteMidPrefix), ConfigSuffix) // TODO verify prefix and suffix? Perl doesn't

	tcDS := DeliveryService{}
	for _, ds := range deliveryServices {
		if ds.XMLID == nil || *ds.XMLID != dsName {
			continue
		}
		tcDS = ds
		break
	}
	if tcDS.ID == nil {
		return Cfg{}, makeErr(warnings, "ds '"+dsName+"' not found")
	}

	if tcDS.CDNName == nil {
		return Cfg{}, makeErr(warnings, "ds '"+dsName+"' missing cdn")
	}

	ds, err := headerRewriteDSFromDS(&tcDS)
	if err != nil {
		return Cfg{}, makeErr(warnings, "converting ds to config ds: "+err.Error())
	}

	assignedServers := map[int]struct{}{}
	for _, dss := range deliveryServiceServers {
		if dss.Server == nil || dss.DeliveryService == nil {
			continue
		}
		if *dss.DeliveryService != *tcDS.ID {
			continue
		}
		assignedServers[*dss.Server] = struct{}{}
	}

	serverCGs := map[tc.CacheGroupName]struct{}{}
	for _, sv := range servers {
		if sv.CDNName == nil {
			warnings = append(warnings, "TO returned Servers server with missing CDNName, skipping!")
			continue
		} else if sv.ID == nil {
			warnings = append(warnings, "TO returned Servers server with missing ID, skipping!")
			continue
		} else if sv.Status == nil {
			warnings = append(warnings, "TO returned Servers server with missing Status, skipping!")
			continue
		} else if sv.Cachegroup == nil {
			warnings = append(warnings, "TO returned Servers server with missing Cachegroup, skipping!")
			continue
		}

		if *sv.CDNName != *server.CDNName {
			continue
		}
		if _, ok := assignedServers[*sv.ID]; !ok && (tcDS.Topology == nil || *tcDS.Topology == "") {
			continue
		}
		if tc.CacheStatus(*sv.Status) != tc.CacheStatusReported && tc.CacheStatus(*sv.Status) != tc.CacheStatusOnline {
			continue
		}
		serverCGs[tc.CacheGroupName(*sv.Cachegroup)] = struct{}{}
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

	assignedMids := []headerRewriteServer{}
	for _, server := range servers {
		if server.CDNName == nil {
			warnings = append(warnings, "TO returned Servers server with missing CDNName, skipping!")
			continue
		}
		if server.Cachegroup == nil {
			warnings = append(warnings, "TO returned Servers server with missing Cachegroup, skipping!")
			continue
		}
		if *server.CDNName != *tcDS.CDNName {
			continue
		}
		if _, ok := parentCGs[*server.Cachegroup]; !ok {
			continue
		}
		cfgServer, err := headerRewriteServerFromServer(server)
		if err != nil {
			warnings = append(warnings, "failed to make header rewrite server,skipping! : "+err.Error())
			continue
		}
		assignedMids = append(assignedMids, cfgServer)
	}

	text := makeHdrComment(hdrComment)

	// write a header rewrite rule if maxOriginConnections > 0 and the ds DOES use mids
	if ds.MaxOriginConnections > 0 && ds.Type.UsesMidCache() {
		dsOnlineMidCount := 0
		for _, sv := range assignedMids {
			if sv.Status == tc.CacheStatusReported || sv.Status == tc.CacheStatusOnline {
				dsOnlineMidCount++
			}
		}
		if dsOnlineMidCount > 0 {
			maxOriginConnectionsPerMid := int(math.Round(float64(ds.MaxOriginConnections) / float64(dsOnlineMidCount)))
			text += "cond %{REMAP_PSEUDO_HOOK}\nset-config proxy.config.http.origin_max_connections " + strconv.Itoa(maxOriginConnectionsPerMid) + "\n"
		}
	}

	if !strings.Contains(ds.MidHeaderRewrite, ServiceCategoryHeader) && ds.ServiceCategory != "" {
		text += "cond %{REMAP_PSEUDO_HOOK}\nset-header " + ServiceCategoryHeader + ` "` + dsName + "|" + ds.ServiceCategory + `"` + "\n"
	}

	// write the contents of ds.MidHeaderRewrite to hdr_rw_mid_xml-id.config replacing any instances of __RETURN__ (surrounded by spaces or not) with \n
	if ds.MidHeaderRewrite != "" {
		re := regexp.MustCompile(`\s*__RETURN__\s*`)
		text += re.ReplaceAllString(ds.MidHeaderRewrite, "\n")
	}

	text += "\n"

	return Cfg{
		Text:        text,
		ContentType: ContentTypeHeaderRewriteDotConfig,
		LineComment: LineCommentHeaderRewriteDotConfig,
		Warnings:    warnings,
	}, nil
}
