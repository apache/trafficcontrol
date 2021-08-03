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

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
)

// ContentTypeRegexRemapDotConfig is the MIME type of the contents of a
// regex_remap.config ATS configuration file.
const ContentTypeRegexRemapDotConfig = ContentTypeTextASCII

// LineCommentRegexRemapDotConfig is the string used to indicate the start of a
// line comment in the grammar of a regex_remap.config ATS configuration file.
const LineCommentRegexRemapDotConfig = LineCommentHash

// RegexRemapDotConfigOpts contains settings to configure generation options.
type RegexRemapDotConfigOpts struct {
	// HdrComment is the header comment to include at the beginning of the file.
	// This should be the text desired, without comment syntax (like # or //). The file's comment syntax will be added.
	// To omit the header comment, pass the empty string.
	HdrComment string
}

// RegexRemapPrefix is a prefix applied to regex_remap.config ATS configuration
// files for a sepecific Delivery Service. The rest of the name is made up of
// the Delivery Service's XMLID, followed by '.config' as a suffix (or "file
// extension").
const RegexRemapPrefix = "regex_remap_"

// MakeRegexRemapDotConfig constructs a regex_remap.config file (specifically
// the name is given by 'fileName') for a cache server.
func MakeRegexRemapDotConfig(
	fileName string,
	server *Server,
	deliveryServices []DeliveryService,
	opt *RegexRemapDotConfigOpts,
) (Cfg, error) {
	if opt == nil {
		opt = &RegexRemapDotConfigOpts{}
	}
	warnings := []string{}
	if server.CDN == "" {
		return Cfg{}, makeErr(warnings, "server CDNName missing")
	}

	configSuffix := `.config`
	if !strings.HasPrefix(fileName, RegexRemapPrefix) || !strings.HasSuffix(fileName, configSuffix) {
		return Cfg{}, makeErr(warnings, "file '"+fileName+"' not of the form 'regex_remap_*.config! Please file a bug with Traffic Control, this should never happen")
	}

	dsName := strings.TrimSuffix(strings.TrimPrefix(fileName, RegexRemapPrefix), configSuffix)
	if dsName == "" {
		return Cfg{}, makeErr(warnings, "file '"+fileName+"' has no delivery service name!")
	}

	// only send the requested DS to atscfg. The atscfg.Make will work correctly even if we send it other DSes, but this will prevent deliveryServicesToCDNDSes from logging errors about AnyMap and Steering DSes without origins.
	ds := DeliveryService{}
	for _, dsesDS := range deliveryServices {
		if dsesDS.XMLID == "" {
			continue // TODO log?
		}
		if dsesDS.XMLID != dsName {
			continue
		}
		ds = dsesDS
	}
	if ds.ID == nil {
		return Cfg{}, makeErr(warnings, "delivery service '"+dsName+"' not found! Do you have a regex_remap_*.config location Parameter for a delivery service that doesn't exist?")
	}

	dses, dsWarns := deliveryServicesToCDNDSes([]DeliveryService{ds})
	warnings = append(warnings, dsWarns...)

	text := makeHdrComment(opt.HdrComment)

	cdnDS, ok := dses[tc.DeliveryServiceName(dsName)]
	if !ok {
		warnings = append(warnings, "ds '"+dsName+"' not in dses, skipping!")
	} else {
		text += cdnDS.RegexRemap + "\n"
		text = strings.Replace(text, `__RETURN__`, "\n", -1)
	}

	return Cfg{
		Text:        text,
		ContentType: ContentTypeRegexRemapDotConfig,
		LineComment: LineCommentRegexRemapDotConfig,
		Warnings:    warnings,
	}, nil
}

type cdnDS struct {
	OrgServerFQDN string
	QStringIgnore int
	RegexRemap    string
}

// deliveryServicesToCDNDSes returns the CDNDSes and any warnings.
func deliveryServicesToCDNDSes(dses []DeliveryService) (map[tc.DeliveryServiceName]cdnDS, []string) {
	warnings := []string{}
	sDSes := map[tc.DeliveryServiceName]cdnDS{}
	for _, ds := range dses {
		if ds.OrgServerFQDN == nil || ds.QStringIgnore == nil || ds.XMLID == "" {
			if ds.XMLID == "" {
				warnings = append(warnings, "got unknown DS with nil values! Skipping!")
			} else {
				warnings = append(warnings, "got DS '"+ds.XMLID+"' with nil values! Skipping!")
			}
			continue
		}
		sds := cdnDS{OrgServerFQDN: *ds.OrgServerFQDN, QStringIgnore: *ds.QStringIgnore}
		if ds.RegexRemap != nil {
			sds.RegexRemap = *ds.RegexRemap
		}
		sDSes[tc.DeliveryServiceName(ds.XMLID)] = sds
	}
	return sDSes, warnings
}
