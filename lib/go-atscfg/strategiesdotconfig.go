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
	//	"github.com/apache/trafficcontrol/lib/go-util"
)

const ContentTypeStrategiesDotYAML = ContentTypeYAML
const LineCommentStrategiesDotYAML = LineCommentHash

// StrategiesYAMLOpts contains settings to configure strategies.config generation options.
type StrategiesYAMLOpts struct {
	// VerboseComments is whether to add informative comments to the generated file, about what was generated and why.
	// Note this does not include the header comment, which is configured separately with HdrComment.
	// These comments are human-readable and not guarnateed to be consistent between versions. Automating anything based on them is strongly discouraged.
	VerboseComments bool

	// HdrComment is the header comment to include at the beginning of the file.
	// This should be the text desired, without comment syntax (like # or //). The file's comment syntax will be added.
	// To omit the header comment, pass the empty string.
	HdrComment string
}

func MakeStrategiesDotYAML(
	dses []DeliveryService,
	server *Server,
	servers []Server,
	topologies []tc.Topology,
	tcServerParams []tc.Parameter,
	tcParentConfigParams []tc.Parameter,
	serverCapabilities map[int]map[ServerCapability]struct{},
	dsRequiredCapabilities map[int]map[ServerCapability]struct{},
	cacheGroupArr []tc.CacheGroupNullable,
	dss []DeliveryServiceServer,
	cdn *tc.CDN,
	opt *StrategiesYAMLOpts,
) (Cfg, error) {
	if opt == nil {
		opt = &StrategiesYAMLOpts{}
	}
	parentAbstraction, warnings, err := makeParentDotConfigData(
		dses,
		server,
		servers,
		topologies,
		tcServerParams,
		tcParentConfigParams,
		serverCapabilities,
		dsRequiredCapabilities,
		cacheGroupArr,
		dss,
		cdn,
		&ParentConfigOpts{
			AddComments: opt.VerboseComments,
			HdrComment:  opt.HdrComment,
		}, // TODO change makeParentDotConfigData to its own opt?
	)
	if err != nil {
		return Cfg{}, makeErr(warnings, err.Error())
	}

	atsMajorVer, verWarns := getATSMajorVersion(tcServerParams)
	warnings = append(warnings, verWarns...)

	text, paWarns, err := parentAbstractionToStrategiesDotYaml(parentAbstraction, opt, atsMajorVer)
	warnings = append(warnings, paWarns...)
	if err != nil {
		return Cfg{}, makeErr(warnings, err.Error())
	}

	hdr := ""
	if opt.HdrComment != "" {
		hdr = makeHdrComment(opt.HdrComment)
	}

	return Cfg{
		Text:        hdr + text,
		ContentType: ContentTypeParentDotConfig,
		LineComment: LineCommentParentDotConfig,
		Warnings:    warnings,
	}, nil
}

// YAMLDocumentStart is the YAML document start directive
const YAMLDocumentStart = "---"
const YAMLDocumentEnd = "..."

func parentAbstractionToStrategiesDotYaml(pa *ParentAbstraction, opt *StrategiesYAMLOpts, atsMajorVersion int) (string, []string, error) {
	warnings := []string{}
	txt := YAMLDocumentStart +
		getStrategyHostsSection(pa) +
		getStrategyGroupsSection(pa) +
		getStrategyStrategiesSection(pa) +
		"\n" + YAMLDocumentEnd +
		"\n"
	return txt, warnings, nil
}

// getStrategyPolicy returns the strategies.config text for the retry policy.
// Returns the default policy if policy is invalid, without error.
func getStrategyPolicy(policy ParentAbstractionServiceRetryPolicy) string {
	switch policy {
	case ParentAbstractionServiceRetryPolicyConsistentHash:
		return `consistent_hash`
	case ParentAbstractionServiceRetryPolicyRoundRobinIP:
		return `rr_ip`
	case ParentAbstractionServiceRetryPolicyRoundRobinStrict:
		return `rr_strict`
	case ParentAbstractionServiceRetryPolicyFirst:
		return `first_live`
	case ParentAbstractionServiceRetryPolicyLatched:
		return `latched`
	default:
		return getStrategyPolicy(DefaultParentAbstractionServiceRetryPolicy)
	}
}

func getStrategyName(dsName string) string {
	return "strategy-" + dsName
}

func getStrategySecondaryMode(mode ParentAbstractionServiceParentSecondaryMode) string {
	switch mode {
	case ParentAbstractionServiceParentSecondaryModeExhaust:
		return `exhaust_ring`
	case ParentAbstractionServiceParentSecondaryModeAlternate:
		return `alternate_ring`
	case ParentAbstractionServiceParentSecondaryModePeering:
		return `peering_ring`
	default:
		return getStrategySecondaryMode(ParentAbstractionServiceParentSecondaryModeDefault)
	}
}

func getStrategyErrorCodes(codes []int) string {
	str := " ["
	join := " "
	for _, code := range codes {
		str += join + strconv.Itoa(code)
		join = ", "
	}
	str += " ]"
	return str
}

func getStrategyGroups(svc *ParentAbstractionService) string {
	txt := ""
	if getStrategySecondaryMode(svc.SecondaryMode) == "peering_ring" {
		txt += "\n" + `      - *peers_group`
		if len(svc.Parents) != 0 {
			txt += "\n" + `      - *group_parents_` + svc.Name
		}
	} else {
		if len(svc.Parents) != 0 {
			txt += "\n" + `      - *group_parents_` + svc.Name
		}
		if len(svc.SecondaryParents) != 0 {
			txt += "\n" + `      - *group_secondary_parents_` + svc.Name
		}
	}
	return txt
}

func getStrategyStrategiesSection(pa *ParentAbstraction) string {
	txt := "\n" + `strategies:`
	for _, svc := range pa.Services {
		txt += "\n" + `  - strategy: '` + getStrategyName(svc.Name) + `'`
		txt += "\n" + `    policy: ` + getStrategyPolicy(svc.RetryPolicy)
		if !svc.IgnoreQueryStringInParentSelection {
			txt += "\n" + `    hash_key: path+query`
		} else {
			txt += "\n" + `    hash_key: path`
		}
		txt += "\n" + `    go_direct: ` + strconv.FormatBool(svc.GoDirect)
		if getStrategySecondaryMode(svc.SecondaryMode) == "peering_ring" {
			txt += "\n" + `    cache_peer_result: ` + strconv.FormatBool(svc.CachePeerResult)
		}
		txt += "\n" + `    groups:`
		txt += getStrategyGroups(svc)
		// TODO make strategies for both? add to parent?
		//      Ask John why this exists. Since it's specified on the remap, shouldn't it use the
		//      remap target's scheme?
		// txt += "\n" + `    scheme: http`
		txt += "\n" + `    failover:`
		txt += "\n" + `      ring_mode: ` + getStrategySecondaryMode(svc.SecondaryMode)
		if len(svc.ErrorResponseCodes) > 0 {
			txt += "\n" + `      max_simple_retries: ` + strconv.Itoa(svc.MaxSimpleRetries)
			txt += "\n" + `      response_codes:`
			txt += getStrategyErrorCodes(svc.ErrorResponseCodes)
		}
		if len(svc.MarkdownResponseCodes) > 0 {
			txt += "\n" + `      max_unavailable_retries: ` + strconv.Itoa(svc.MaxMarkdownRetries)
			txt += "\n" + `      markdown_codes:`
			txt += getStrategyErrorCodes(svc.MarkdownResponseCodes)
		}
		txt += "\n" + `      health_check:`
		txt += "\n" + `        - passive`
	}
	return txt
}

func serviceGroupParentsName(pa *ParentAbstractionService) string {
	anchorSvcName := pa.Name
	anchorSvcName = strings.Replace(anchorSvcName, `.`, `-dot-`, -1)
	return `group_parents_` + anchorSvcName
}

func serviceGroupSecondaryParentsName(pa *ParentAbstractionService) string {
	anchorSvcName := pa.Name
	anchorSvcName = strings.Replace(anchorSvcName, `.`, `-dot-`, -1)
	return `group_secondary_parents_` + anchorSvcName
}

func getStrategyGroupsSection(pa *ParentAbstraction) string {
	txt := "\n" + `groups:`
	for _, svc := range pa.Services {
		if len(svc.Parents) != 0 {
			txt += "\n" + `  - &` + serviceGroupParentsName(svc)
		}
		for _, parent := range svc.Parents {
			txt += "\n" + `    - <<: *` + getStrategyParentHostEntryName(svc, parent)
			txt += "\n" + `      weight: ` + strconv.FormatFloat(parent.Weight, 'f', 3, 64)
		}

		if len(svc.SecondaryParents) != 0 {
			txt += "\n" + `  - &` + serviceGroupSecondaryParentsName(svc)
		}
		for _, parent := range svc.SecondaryParents {
			txt += "\n" + `    - <<: *` + getStrategyParentHostEntryName(svc, parent)
			txt += "\n" + `      weight: ` + strconv.FormatFloat(parent.Weight, 'f', 3, 64)
		}
	}
	if len(pa.Peers) != 0 {
		txt += "\n" + `  - &peers_group`
		for i, peer := range pa.Peers {
			txt += "\n" + `    - <<: *peer` + strconv.Itoa(i+1)
			txt += "\n" + `      weight: ` + strconv.FormatFloat(peer.Weight, 'f', 3, 64)
		}
	}
	return txt
}

func getStrategyParentHostEntryName(svc *ParentAbstractionService, parent *ParentAbstractionServiceParent) string {
	// fqdn characters are valid yaml anchor characters except for '.'
	anchorFQDN := strings.Replace(parent.FQDN, `.`, `-dot-`, -1)
	return `host__` + svc.Name + `__parent__` + anchorFQDN + "__" + strconv.Itoa(parent.Port)
}

func getStrategyHostsSection(pa *ParentAbstraction) string {
	txt := "\n" + `hosts:`
	for _, svc := range pa.Services {
		for _, parent := range svc.Parents {
			txt += getStrategyHostsSectionHost(svc, parent)
		}
		for _, parent := range svc.SecondaryParents {
			txt += getStrategyHostsSectionHost(svc, parent)
		}
	}
	if len(pa.Peers) != 0 {
		count := 1
		for _, peer := range pa.Peers {
			txt += getStrategyHostsSectionPeer(peer, count)
			count++
		}
	}
	return txt
}

func getStrategyHostsSectionPeer(peer *ParentAbstractionServiceParent, count int) string {
	txt := ""
	txt += "\n" + `  - &peer` + strconv.Itoa(count)
	txt += "\n" + `    host: ` + peer.FQDN
	txt += "\n" + `    protocol:`
	txt += "\n" + `      - port: ` + strconv.Itoa(peer.Port)
	return txt
}

func getStrategyHostsSectionHost(svc *ParentAbstractionService, parent *ParentAbstractionServiceParent) string {
	txt := ""
	txt += "\n" + `  - &` + getStrategyParentHostEntryName(svc, parent)
	txt += "\n" + `    host: ` + parent.FQDN
	txt += "\n" + `    protocol:`
	// txt += "\n" + `      - scheme: http` // TODO fix?
	txt += "\n" + `      - port: ` + strconv.Itoa(parent.Port)
	return txt
}
