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
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

const SNIDotYAMLFileName = "sni.yaml"

const ContentTypeSNIDotYAML = ContentTypeYAML
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

	// E2ESSLData has data for adding and enforce certificates for intra-cache communication.
	// For parents, this means requiring children present client certificates validated by the given CA.
	// For children, this means passing a client certificate, and requiring the parent cache present a valid origin certificate
	//
	// Note this does not set whether HTTPS will be used. That is, if remap.config and parent.config are not configured
	// to use https, these settings will not be used, and will not prevent http from succeeding.
	//
	// To not add end-to-end inter-cdn TLS certificate data, pass a nil pointer.
	E2ESSLData SNIDotYAMLE2EInf
}

type SNIDotYAMLE2EInf struct {
	// ClientCAPath is the path to the Certificate Authority used to sign client certificates which will be presented by child caches to parent caches.
	// This must be relative to the ATS config directory (e.g. etc/trafficserver/), _not_ the proxy.config.ssl.server.cert.path (e.g. etc/trafficserver/ssl/).
	ClientCAPath string
	// ClientCertPath is the path to the client certificate presented by child caches to parent caches.
	// This is relative to records.config proxy.config.ssl.client.cert.path.
	ClientCertPath string
	// ClientCertKeyPath is the path to the key for the client certificate presented by child caches to parent caches.
	// This is relative to records.config proxy.config.ssl.client.private_key.path.
	ClientCertKeyPath string
}

func MakeSNIDotYAML(
	server *Server,
	dses []DeliveryService,
	dss []DeliveryServiceServer,
	dsRegexArr []tc.DeliveryServiceRegexes,
	tcParentConfigParams []tc.Parameter,
	cdn *tc.CDN,
	topologies []tc.Topology,
	cacheGroupArr []tc.CacheGroupNullable,
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

		if sslData.IsEdge {
			for _, requestFQDN := range sslData.RequestFQDNs {
				// TODO let active DSes take precedence?
				if _, ok := seenFQDNs[requestFQDN]; ok {
					warnings = append(warnings, "ds '"+sslData.DSName+"' had the same FQDN '"+requestFQDN+"' as some other delivery service, skipping!")
					continue
				}
				seenFQDNs[requestFQDN] = struct{}{}

				dsTxt := "\n"
				if opt.VerboseComments {
					dsTxt += LineCommentYAML + ` ds '` + sslData.DSName + `' from-external-client` + "\n"
				}
				dsTxt += `- fqdn: '` + requestFQDN + `'`
				dsTxt += "\n" + `  disable_h2: ` + strconv.FormatBool(!sslData.EnableH2)
				dsTxt += "\n" + `  valid_tls_versions_in: [` + strings.Join(tlsVersionsATS, `,`) + `]`
				txt += dsTxt + "\n"
			}
		}

		if opt.E2ESSLData != (SNIDotYAMLE2EInf{}) && sslData.OriginFQDN != "" {
			if sslData.IsParent || sslData.IsChild {
				dsTxt := "\n"
				if opt.VerboseComments {
					dsTxt += LineCommentYAML + ` ds '` + sslData.DSName + `'`
					if sslData.IsParent {
						dsTxt += ` from-internal-child`
					}
					if sslData.IsChild {
						dsTxt += ` to-internal-parent`
					}
					dsTxt += "\n"
				}
				dsTxt += `- fqdn: '` + sslData.OriginFQDN + `'`
				if sslData.IsParent {
					verifyClientTxt := `STRICT`
					if sslData.DisableInternalValidation {
						verifyClientTxt = `NONE`
						warnings = append(warnings, "ds '"+sslData.DSName+"' had parameter to disable internal certificate validation! Internal encrypted traffic will not be verified! It is strongly encouraged to fix the problem and re-enable TLS validation!!")
					}
					dsTxt += "\n" + `  verify_client: '` + verifyClientTxt + `'`
					dsTxt += "\n" + `  verify_client_ca_certs: '` + opt.E2ESSLData.ClientCAPath + `'`
				}
				if sslData.IsChild {
					verifyServerTxt := `ENFORCED`
					if sslData.DisableInternalValidation {
						verifyServerTxt = `DISABLED`
						warnings = append(warnings, "ds '"+sslData.DSName+"' had parameter to disable internal certificate validation! Internal encrypted traffic will not be verified! It is strongly encouraged to fix the problem and re-enable TLS validation!!")
					}
					dsTxt += "\n" + `  verify_server_policy: '` + verifyServerTxt + `'`
					dsTxt += "\n" + `  client_cert: '` + opt.E2ESSLData.ClientCertPath + `'`
					dsTxt += "\n" + `  client_key: '` + opt.E2ESSLData.ClientCertKeyPath + `'`
				}
				txt += dsTxt + "\n"
			}
		}
	}

	return Cfg{
		Text:        txt,
		ContentType: ContentTypeSNIDotYAML,
		LineComment: LineCommentSNIDotYAML,
		Warnings:    warnings,
	}, nil
}
