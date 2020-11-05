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
	"strings"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

const ContentTypeCacheURLDotConfig = ContentTypeTextASCII
const LineCommentCacheURLDotConfig = LineCommentHash

func MakeCacheURLDotConfig(
	fileName string,
	server *tc.ServerNullable,
	deliveryServices []tc.DeliveryServiceNullableV30,
	deliveryServiceServers []tc.DeliveryServiceServer,
	hdrComment string,
) (Cfg, error) {
	warnings := []string{}

	if server.CDNName == nil {
		return Cfg{}, makeErr(warnings, "server missing CDNName")
	}

	dsIDs := map[int]struct{}{}
	for _, ds := range deliveryServices {
		if ds.ID != nil {
			dsIDs[*ds.ID] = struct{}{} // TODO warn?
		}
	}

	dss := filterDSS(deliveryServiceServers, dsIDs, nil)

	dssMap := map[int][]int{} // map[dsID]serverID
	for _, dss := range dss {
		if dss.Server == nil || dss.DeliveryService == nil {
			warnings = append(warnings, "Delivery Service Servers had nil entries, skipping!")
			continue
		}
		dssMap[*dss.DeliveryService] = append(dssMap[*dss.DeliveryService], *dss.Server)
	}

	dsesWithServers := []tc.DeliveryServiceNullableV30{}
	for _, ds := range deliveryServices {
		if ds.ID == nil {
			warnings = append(warnings, "Delivery Service had nil id, skipping!")
			continue
		}
		// ANY_MAP and STEERING DSes don't have origins, and thus can't be put into the cacheurl config.
		if ds.Type != nil && (*ds.Type == tc.DSTypeAnyMap || *ds.Type == tc.DSTypeSteering) {
			continue
		}
		if len(dssMap[*ds.ID]) == 0 && ds.Topology == nil {
			continue
		}
		dsesWithServers = append(dsesWithServers, ds)
	}

	dses, dsWarns := deliveryServicesToCacheURLDSes(dsesWithServers)
	warnings = append(warnings, dsWarns...)

	text := makeHdrComment(hdrComment)

	if fileName == "cacheurl_qstring.config" { // This is the per remap drop qstring w cacheurl use case, the file is the same for all remaps
		text += `http://([^?]+)(?:\?|$)  http://$1` + "\n"
		text += `https://([^?]+)(?:\?|$)  https://$1` + "\n"

		return Cfg{
			Text:        text,
			ContentType: ContentTypeCacheURLDotConfig,
			LineComment: LineCommentCacheURLDotConfig,
			Warnings:    warnings,
		}, nil
	}

	if fileName == "cacheurl.config" { // this is the global drop qstring w cacheurl use case
		seenOrigins := map[string]struct{}{}
		for dsName, ds := range dses {
			if ds.QStringIgnore != 1 {
				continue
			}
			if _, ok := seenOrigins[ds.OrgServerFQDN]; ok {
				continue
			}
			org := ds.OrgServerFQDN

			scheme := "https://"
			if !strings.HasPrefix(org, scheme) {
				scheme = "http://"
			}

			if !strings.HasPrefix(org, scheme) {
				// TODO determine if we should return an empty config here. A bad DS should not break config gen, and MUST NOT for self-service
				warnings = append(warnings, "ds '"+string(dsName)+"' origin '"+org+"' with no scheme! cacheurl.config will likely be malformed!")
			}

			fqdnPath := strings.TrimPrefix(org, scheme)

			text += scheme + `(` + fqdnPath + `/[^?]+)(?:\?|$)  ` + scheme + `$1` + "\n"

			seenOrigins[ds.OrgServerFQDN] = struct{}{}
		}
		text = strings.Replace(text, `__RETURN__`, "\n", -1)
		return Cfg{
			Text:        text,
			ContentType: ContentTypeCacheURLDotConfig,
			LineComment: LineCommentCacheURLDotConfig,
			Warnings:    warnings,
		}, nil
	}

	// TODO verify prefix and suffix exist, and warn if they don't? Perl doesn't
	dsName := tc.DeliveryServiceName(strings.TrimSuffix(strings.TrimPrefix(fileName, "cacheurl_"), ".config"))

	ds, ok := dses[dsName]
	if !ok {
		warnings = append(warnings, "ds '"+string(dsName)+"' not found, not creating in cacheurl config!")
	} else {
		text += ds.CacheURL + "\n"
		text = strings.Replace(text, `__RETURN__`, "\n", -1)
	}
	return Cfg{
		Text:        text,
		ContentType: ContentTypeCacheURLDotConfig,
		LineComment: LineCommentCacheURLDotConfig,
		Warnings:    warnings,
	}, nil
}

type cacheURLDS struct {
	OrgServerFQDN string
	QStringIgnore int
	CacheURL      string
}

// DeliveryServicesToCacheURLDSes returns the "CacheURLDS" map, and any warnings.
func deliveryServicesToCacheURLDSes(dses []tc.DeliveryServiceNullableV30) (map[tc.DeliveryServiceName]cacheURLDS, []string) {
	warnings := []string{}
	sDSes := map[tc.DeliveryServiceName]cacheURLDS{}
	for _, ds := range dses {
		if ds.OrgServerFQDN == nil || ds.QStringIgnore == nil || ds.XMLID == nil || ds.Active == nil {
			warnings = append(warnings, fmt.Sprintf("atscfg.DeliveryServicesToCacheURLDSes got DS %+v with nil values! Skipping!", ds))
			continue
		}
		if !*ds.Active {
			continue
		}
		sds := cacheURLDS{OrgServerFQDN: *ds.OrgServerFQDN, QStringIgnore: *ds.QStringIgnore}
		if ds.CacheURL != nil {
			sds.CacheURL = *ds.CacheURL
		}
		sDSes[tc.DeliveryServiceName(*ds.XMLID)] = sds
	}
	return sDSes, warnings
}
