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
	"strings"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

const ContentTypeCacheDotConfig = ContentTypeTextASCII
const LineCommentCacheDotConfig = LineCommentHash

func MakeCacheDotConfig(
	server *Server,
	servers []Server,
	deliveryServices []tc.DeliveryServiceNullableV30,
	deliveryServiceServers []tc.DeliveryServiceServer,
	hdrComment string,
) (Cfg, error) {
	if tc.CacheTypeFromString(server.Type) == tc.CacheTypeMid {
		return makeCacheDotConfigMid(server, deliveryServices, hdrComment)
	} else {
		return makeCacheDotConfigEdge(server, servers, deliveryServices, deliveryServiceServers, hdrComment)
	}
}

// MakeCacheDotConfig makes the ATS cache.config config file.
// profileDSes must be the list of delivery services, which are assigned to severs, for which this profile is assigned. It MUST NOT contain any other delivery services. Note DSesToProfileDSes may be helpful if you have a []tc.DeliveryServiceNullable, for example from traffic_ops/client.
func makeCacheDotConfigEdge(
	server *Server,
	servers []Server,
	deliveryServices []tc.DeliveryServiceNullableV30,
	deliveryServiceServers []tc.DeliveryServiceServer,
	hdrComment string,
) (Cfg, error) {
	warnings := []string{}

	if server.Profile == nil {
		return Cfg{}, makeErr(warnings, "server missing profile")
	}

	profileServerIDsMap := map[int]struct{}{}
	for _, sv := range servers {
		if sv.Profile == nil {
			warnings = append(warnings, "servers had server with nil profile, skipping!")
			continue
		}
		if sv.ID == nil {
			warnings = append(warnings, "servers had server with nil id, skipping!")
			continue
		}
		if *sv.Profile != *server.Profile {
			continue
		}
		profileServerIDsMap[*sv.ID] = struct{}{}
	}

	dsServers := filterDSS(deliveryServiceServers, nil, profileServerIDsMap)

	dsIDs := map[int]struct{}{}
	for _, dss := range dsServers {
		if dss.Server == nil || dss.DeliveryService == nil {
			continue // TODO warn? err?
		}
		if _, ok := profileServerIDsMap[*dss.Server]; !ok {
			continue
		}
		dsIDs[*dss.DeliveryService] = struct{}{}
	}

	profileDSes := []profileDS{}
	for _, ds := range deliveryServices {
		if ds.ID == nil || ds.Type == nil || ds.OrgServerFQDN == nil {
			continue // TODO warn? err?
		}
		if *ds.Type == tc.DSTypeInvalid {
			continue // TODO warn? err?
		}
		if *ds.OrgServerFQDN == "" {
			continue // TODO warn? err?
		}
		if _, ok := dsIDs[*ds.ID]; !ok && ds.Topology == nil {
			continue
		}
		origin := *ds.OrgServerFQDN
		profileDSes = append(profileDSes, profileDS{Type: *ds.Type, OriginFQDN: &origin})
	}

	lines := map[string]struct{}{} // use a "set" for lines, to avoid duplicates, since we're looking up by profile
	for _, ds := range profileDSes {
		if ds.Type != tc.DSTypeHTTPNoCache {
			continue
		}
		if ds.OriginFQDN == nil || *ds.OriginFQDN == "" {
			warnings = append(warnings, "profileCacheDotConfig ds has no origin fqdn, skipping!") // TODO add ds name to data loaded, to put it in the error here?
			continue
		}
		originFQDN, originPort := getHostPortFromURI(*ds.OriginFQDN)
		if originPort != "" {
			l := "dest_domain=" + originFQDN + " port=" + originPort + " scheme=http action=never-cache\n"
			lines[l] = struct{}{}
		} else {
			l := "dest_domain=" + originFQDN + " scheme=http action=never-cache\n"
			lines[l] = struct{}{}
		}
	}

	linesArr := []string{}
	for line, _ := range lines {
		linesArr = append(linesArr, line)
	}
	sort.Strings(linesArr)
	text := strings.Join(linesArr, "")
	if text == "" {
		text = "\n" // If no params exist, don't send "not found," but an empty file. We know the profile exists.
	}
	hdr := makeHdrComment(hdrComment)
	text = hdr + text

	return Cfg{
		Text:        text,
		ContentType: ContentTypeCacheDotConfig,
		LineComment: LineCommentCacheDotConfig,
		Warnings:    warnings,
	}, nil
}

type profileDS struct {
	Type       tc.DSType
	OriginFQDN *string
}

// dsesToProfileDSes is a helper function to convert a []tc.DeliveryServiceNullable to []ProfileDS.
// Note this does not check for nil values. If any DeliveryService's Type or OrgServerFQDN may be nil, the returned ProfileDS should be checked for DSTypeInvalid and nil, respectively.
func dsesToProfileDSes(dses []tc.DeliveryServiceNullable) []profileDS {
	pdses := []profileDS{}
	for _, ds := range dses {
		pds := profileDS{}
		if ds.Type != nil {
			pds.Type = *ds.Type
		}
		if ds.OrgServerFQDN != nil && *ds.OrgServerFQDN != "" {
			org := *ds.OrgServerFQDN
			pds.OriginFQDN = &org
		}
		pdses = append(pdses, pds)
	}
	return pdses
}

func getHostPortFromURI(uriStr string) (string, string) {
	originFQDN := uriStr
	originFQDN = strings.TrimPrefix(originFQDN, "http://")
	originFQDN = strings.TrimPrefix(originFQDN, "https://")

	slashPos := strings.Index(originFQDN, "/")
	if slashPos != -1 {
		originFQDN = originFQDN[:slashPos]
	}
	portPos := strings.Index(originFQDN, ":")
	portStr := ""
	if portPos != -1 {
		portStr = originFQDN[portPos+1:]
		originFQDN = originFQDN[:portPos]
	}
	return originFQDN, portStr
}
