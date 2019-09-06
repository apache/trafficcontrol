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
	"strings"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
)

type CacheURLDS struct {
	OrgServerFQDN string
	QStringIgnore int
	CacheURL      string
}

func DeliveryServicesToCacheURLDSes(dses []tc.DeliveryServiceNullable) map[tc.DeliveryServiceName]CacheURLDS {
	sDSes := map[tc.DeliveryServiceName]CacheURLDS{}
	for _, ds := range dses {
		if ds.OrgServerFQDN == nil || ds.QStringIgnore == nil || ds.XMLID == nil || ds.Active == nil {
			log.Errorf("atscfg.DeliveryServicesToCacheURLDSes got DS %+v with nil values! Skipping!", ds)
			continue
		}
		if !*ds.Active {
			continue
		}
		sds := CacheURLDS{OrgServerFQDN: *ds.OrgServerFQDN, QStringIgnore: *ds.QStringIgnore}
		if ds.CacheURL != nil {
			sds.CacheURL = *ds.CacheURL
		}
		sDSes[tc.DeliveryServiceName(*ds.XMLID)] = sds
	}
	return sDSes
}

func MakeCacheURLDotConfig(
	cdnName tc.CDNName,
	toToolName string, // tm.toolname global parameter (TODO: cache itself?)
	toURL string, // tm.url global parameter (TODO: cache itself?)
	fileName string,
	dses map[tc.DeliveryServiceName]CacheURLDS,
) string {
	text := GenericHeaderComment(string(cdnName), toToolName, toURL)

	if fileName == "cacheurl_qstring.config" { // This is the per remap drop qstring w cacheurl use case, the file is the same for all remaps
		text += `http://([^?]+)(?:\?|$)  http://$1` + "\n"
		text += `https://([^?]+)(?:\?|$)  https://$1` + "\n"
		return text
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
				log.Errorln("MakeCacheURLDotConfig got ds '" + string(dsName) + "' origin '" + org + "' with no scheme! cacheurl.config will likely be malformed!")
			}

			fqdnPath := strings.TrimPrefix(org, scheme)

			text += scheme + `(` + fqdnPath + `/[^?]+)(?:\?|$)  ` + scheme + `$1` + "\n"

			seenOrigins[ds.OrgServerFQDN] = struct{}{}
		}
		text = strings.Replace(text, `__RETURN__`, "\n", -1)
		return text
	}

	// TODO verify prefix and suffix exist, and warn if they don't? Perl doesn't
	dsName := tc.DeliveryServiceName(strings.TrimSuffix(strings.TrimPrefix(fileName, "cacheurl_"), ".config"))

	ds, ok := dses[dsName]
	if !ok {
		return text // TODO warn? Perl doesn't
	}
	text += ds.CacheURL + "\n"
	text = strings.Replace(text, `__RETURN__`, "\n", -1)
	return text
}
