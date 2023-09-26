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

// SNIDotYAMLFileName is an unused constant that probably defines either the
// name on disk of an sni.yaml ATS configuration file, the ConfigFile value of
// Parameters that may influence the generation of such files, or both.
//
// TODO: replace instances of sni.yaml with this constant.
const SNIDotYAMLFileName = "sni.yaml"

// ContentTypeSNIDotYAML is the MIME type of the contents of an sni.yaml ATS
// configuration file.
const ContentTypeSNIDotYAML = ContentTypeYAML

// LineCommentSNIDotYAML is the string used by the grammar of an sni.yaml file
// to indicate the start of a line comment.
const LineCommentSNIDotYAML = LineCommentYAML

// SNIDotYAMLOpts contains settings to configure sni.yaml generation options.
type SNIDotYAMLOpts struct {
	// VerboseComments is whether to add informative comments to the generated file, about what was generated and why.
	// Note this does not include the header comment, which is configured separately with HdrComment.
	// These comments are human-readable and not guaranteed to be consistent between versions. Automating anything based on them is strongly discouraged.
	VerboseComments bool

	// HdrComment is the header comment to include at the beginning of the file.
	// This should be the text desired, without comment syntax (like # or //). The file's comment syntax will be added.
	// To omit the header comment, pass the empty string.
	HdrComment string

	// DefaultTLSVersions is the list of TLS versions to enable on delivery services with no Parameter.
	DefaultTLSVersions []TLSVersion

	// DefaultEnableH2 is whether to disable H2 on delivery services with no Parameter.
	DefaultEnableH2 bool
}

// MakeSNIDotYAML constructs an sni.yaml ATS configuration file.
func MakeSNIDotYAML(
	server *Server,
	servers []Server,
	dses []DeliveryService,
	dss []DeliveryServiceServer,
	dsRegexArr []tc.DeliveryServiceRegexes,
	tcParentConfigParams []tc.ParameterV5,
	cdn *tc.CDNV5,
	topologies []tc.TopologyV5,
	cacheGroupArr []tc.CacheGroupNullableV5,
	serverCapabilities map[int]map[ServerCapability]struct{},
	dsRequiredCapabilities map[int]map[ServerCapability]struct{},
	opt *SNIDotYAMLOpts,
) (Cfg, error) {
	if opt == nil {
		opt = &SNIDotYAMLOpts{}
	}
	if len(opt.DefaultTLSVersions) == 0 {
		opt.DefaultTLSVersions = DefaultDefaultTLSVersions
	}

	sslDatas, warnings, err := GetServerSSLData(
		server,
		servers,
		dses,
		dss,
		dsRegexArr,
		tcParentConfigParams,
		cdn,
		topologies,
		cacheGroupArr,
		serverCapabilities,
		dsRequiredCapabilities,
		opt.DefaultTLSVersions,
		opt.DefaultEnableH2,
	)
	if err != nil {
		return Cfg{}, makeErr(warnings, "getting ssl data: "+err.Error())
	}

	txt := ""
	if opt.HdrComment != "" {
		txt += makeHdrComment(opt.HdrComment)
	}

	txt += `sni:` + "\n"

	seenFQDNs := map[string]struct{}{}

	for _, sslData := range sslDatas {
		tlsVersionsATS := []string{}
		for _, tlsVersion := range sslData.TLSVersions {
			tlsVersionsATS = append(tlsVersionsATS, `'`+tlsVersionsToATS[tlsVersion]+`'`)
		}

		for _, requestFQDN := range sslData.RequestFQDNs {
			// TODO let active DSes take precedence?
			if _, ok := seenFQDNs[requestFQDN]; ok {
				warnings = append(warnings, "ds '"+sslData.DSName+"' had the same FQDN '"+requestFQDN+"' as some other delivery service, skipping!")
				continue
			}
			seenFQDNs[requestFQDN] = struct{}{}

			dsTxt := "\n"
			if opt.VerboseComments {
				dsTxt += LineCommentYAML + ` ds '` + sslData.DSName + `'` + "\n"
			}
			dsTxt += `- fqdn: '` + requestFQDN + `'`
			dsTxt += "\n" + `  http2: ` + BoolOnOff(sslData.EnableH2)
			dsTxt += "\n" + `  valid_tls_versions_in: [` + strings.Join(tlsVersionsATS, `,`) + `]`

			txt += dsTxt + "\n"
		}
	}

	return Cfg{
		Text:        txt,
		ContentType: ContentTypeSNIDotYAML,
		LineComment: LineCommentSNIDotYAML,
		Warnings:    warnings,
	}, nil
}
