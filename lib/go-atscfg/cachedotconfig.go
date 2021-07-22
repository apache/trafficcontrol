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
	"fmt"
	"sort"
	"strings"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
)

// ContentTypeCacheDotConfig is the MIME type of the contents of a cache.config
// file.
const ContentTypeCacheDotConfig = ContentTypeTextASCII

// LineCommentCacheDotConfig is the string that signifies the beginning of a
// line comment in the grammar of a cache.config file.
const LineCommentCacheDotConfig = LineCommentHash

// CacheDotConfigOpts contains settings to configure generation options.
type CacheDotConfigOpts struct {
	// HdrComment is the header comment to include at the beginning of the file.
	// This should be the text desired, without comment syntax (like # or //). The file's comment syntax will be added.
	// To omit the header comment, pass the empty string.
	HdrComment string
}

// MakeCacheDotConfig makes the ATS cache.config config file.
func MakeCacheDotConfig(
	server *Server,
	servers []Server,
	deliveryServices []DeliveryService,
	deliveryServiceServers []DeliveryServiceServer,
	opt *CacheDotConfigOpts,
) (Cfg, error) {
	if opt == nil {
		opt = &CacheDotConfigOpts{}
	}
	if tc.CacheTypeFromString(server.Type) == tc.CacheTypeMid {
		return makeCacheDotConfigMid(server, deliveryServices, opt)
	} else {
		return makeCacheDotConfigEdge(server, servers, deliveryServices, deliveryServiceServers, opt)
	}
}

func makeCacheDotConfigEdge(
	server *Server,
	servers []Server,
	deliveryServices []DeliveryService,
	deliveryServiceServers []DeliveryServiceServer,
	opt *CacheDotConfigOpts,
) (Cfg, error) {
	if opt == nil {
		opt = &CacheDotConfigOpts{}
	}
	warnings := []string{}

	if len(server.Profiles) == 0 {
		return Cfg{}, makeErr(warnings, "server missing profiles")
	}

	profileDSes, dsWarnings := GetProfileDSes(server, servers, deliveryServices, deliveryServiceServers)
	warnings = append(warnings, dsWarnings...)

	lines := map[string]struct{}{} // use a "set" for lines, to avoid duplicates, since we're looking up by profile
	for _, ds := range profileDSes {
		if ds.Type != tc.DSTypeHTTPNoCache {
			continue
		}
		if ds.OriginFQDN == "" {
			warnings = append(warnings, fmt.Sprintf("profileCacheDotConfig ds %s has no origin fqdn, skipping!", ds.Name))
			continue
		}
		originFQDN, originPort := GetHostPortFromURI(ds.OriginFQDN)
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
	hdr := makeHdrComment(opt.HdrComment)
	text = hdr + text

	return Cfg{
		Text:        text,
		ContentType: ContentTypeCacheDotConfig,
		LineComment: LineCommentCacheDotConfig,
		Warnings:    warnings,
	}, nil
}

// ProfileDS struct for filtered delivery services.
type ProfileDS struct {
	Name       string
	Type       tc.DSType
	OriginFQDN string
}

// GetProfileDSes filters delivery services and return delivery services with valid type and non-empty FQDN.
func GetProfileDSes(server *Server,
	servers []Server,
	deliveryServices []DeliveryService,
	deliveryServiceServers []DeliveryServiceServer,
) ([]ProfileDS, []string) {
	warnings := make([]string, 0)
	profileServerIDsMap := map[int]struct{}{}
	for _, sv := range servers {
		if len(sv.Profiles) == 0 {
			warnings = append(warnings, "servers had server with nil profile, skipping!")
			continue
		}
		if sv.ID == 0 {
			warnings = append(warnings, "servers had server with nil id, skipping!")
			continue
		}
		if !ServerProfilesMatch(server, &sv) {
			continue
		}
		profileServerIDsMap[sv.ID] = struct{}{}
	}

	dsServers := filterDSS(deliveryServiceServers, nil, profileServerIDsMap)

	dsIDs := map[int]struct{}{}
	for _, dss := range dsServers {
		if _, ok := profileServerIDsMap[dss.Server]; !ok {
			continue
		}
		dsIDs[dss.DeliveryService] = struct{}{}
	}

	profileDSes := []ProfileDS{}
	for _, ds := range deliveryServices {
		if ds.Active == tc.DSActiveStateInactive {
			continue
		}
		if ds.ID == nil {
			warnings = append(warnings, "deliveryservices had ds with nil id, skipping!")
			continue
		}
		if ds.Type == nil {
			warnings = append(warnings, "deliveryservices had ds with nil type, skipping!")
			continue
		}
		if ds.OrgServerFQDN == nil {
			continue // this is normal for steering and anymap dses
		}
		if ds.XMLID == "" {
			warnings = append(warnings, "got ds with missing XMLID, skipping!")
			continue
		}
		if *ds.Type == tc.DSTypeInvalid.String() {
			warnings = append(warnings, "deliveryservices had ds with invalid type, skipping!")
			continue
		}
		if *ds.OrgServerFQDN == "" {
			warnings = append(warnings, "deliveryservices had ds with empty origin, skipping!")
			continue
		}
		if _, ok := dsIDs[*ds.ID]; !ok && ds.Topology == nil {
			continue
		}
		profileDSes = append(profileDSes, ProfileDS{Name: ds.XMLID, Type: tc.DSType(*ds.Type), OriginFQDN: *ds.OrgServerFQDN})
	}

	return profileDSes, warnings
}

// dsesToProfileDSes is a helper function to convert a []tc.DeliveryServiceNullable to []ProfileDS.
// Note this does not check for nil values. If any DeliveryService's Type or OrgServerFQDN may be nil, the returned ProfileDS should be checked for DSTypeInvalid and nil, respectively.
func dsesToProfileDSes(dses []tc.DeliveryServiceNullable) []ProfileDS {
	pdses := []ProfileDS{}
	for _, ds := range dses {
		pds := ProfileDS{}
		if ds.Type != nil {
			pds.Type = *ds.Type
		}
		if ds.OrgServerFQDN != nil && *ds.OrgServerFQDN != "" {
			pds.OriginFQDN = *ds.OrgServerFQDN
		}
		pdses = append(pdses, pds)
	}
	return pdses
}

// GetHostPortFromURI strips HTTP(s) scheme and path and return host with port (if found).
func GetHostPortFromURI(uriStr string) (string, string) {
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
