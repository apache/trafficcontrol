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
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
)

const ContentTypeParentDotConfig = ContentTypeTextASCII
const LineCommentParentDotConfig = LineCommentHash

const ParentConfigFileName = "parent.config"

const ParentConfigParamQStringHandling = "psel.qstring_handling"
const ParentConfigParamMSOAlgorithm = "mso.algorithm"
const ParentConfigParamMSOParentRetry = "mso.parent_retry"
const ParentConfigParamMSOUnavailableServerRetryResponses = "mso.unavailable_server_retry_responses"
const ParentConfigParamMSOMaxSimpleRetries = "mso.max_simple_retries"
const ParentConfigParamMSOMaxUnavailableServerRetries = "mso.max_unavailable_server_retries"
const ParentConfigParamMergeGroups = "merge_parent_groups"
const ParentConfigParamAlgorithm = "algorithm"
const ParentConfigParamQString = "qstring"
const ParentConfigParamSecondaryMode = "try_all_primaries_before_secondary"

const ParentConfigParamParentRetry = "parent_retry"
const ParentConfigParamUnavailableServerRetryResponses = "unavailable_server_retry_responses"
const ParentConfigParamMaxSimpleRetries = "max_simple_retries"
const ParentConfigParamMaxUnavailableServerRetries = "max_unavailable_server_retries"

const ParentConfigDSParamDefaultMSOAlgorithm = ParentAbstractionServiceRetryPolicyConsistentHash
const ParentConfigDSParamDefaultMSOParentRetry = "both"
const ParentConfigDSParamDefaultMSOUnavailableServerRetryResponses = ""
const ParentConfigDSParamDefaultMaxSimpleRetries = 1
const ParentConfigDSParamDefaultMaxUnavailableServerRetries = 1

const ParentConfigCacheParamWeight = "weight"
const ParentConfigCacheParamPort = "port"
const ParentConfigCacheParamUseIP = "use_ip_address"
const ParentConfigCacheParamRank = "rank"
const ParentConfigCacheParamNotAParent = "not_a_parent"

type OriginHost string
type OriginFQDN string

// ParentConfigOpts contains settings to configure parent.config generation options.
type ParentConfigOpts struct {
	// AddComments is whether to add informative comments to the generated file, about what was generated and why.
	// Note this does not include the header comment, which is configured separately with HdrComment.
	// These comments are human-readable and not guarnateed to be consistent between versions. Automating anything based on them is strongly discouraged.
	AddComments bool

	// HdrComment is the header comment to include at the beginning of the file.
	// This should be the text desired, without comment syntax (like # or //). The file's comment syntax will be added.
	// To omit the header comment, pass the empty string.
	HdrComment string
}

func MakeParentDotConfig(
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
	opt *ParentConfigOpts,
) (Cfg, error) {
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
		opt,
	)
	if err != nil {
		return Cfg{}, makeErr(warnings, err.Error())
	}

	atsMajorVer, verWarns := getATSMajorVersion(tcServerParams)
	warnings = append(warnings, verWarns...)

	text, paWarns, err := parentAbstractionToParentDotConfig(parentAbstraction, opt, atsMajorVer)
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

func makeParentDotConfigData(
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
	opt *ParentConfigOpts,
) (*ParentAbstraction, []string, error) {
	if opt == nil {
		opt = &ParentConfigOpts{}
	}
	parentAbstraction := &ParentAbstraction{}
	warnings := []string{}

	if server.HostName == nil || *server.HostName == "" {
		return nil, warnings, errors.New("server HostName missing")
	} else if server.CDNName == nil || *server.CDNName == "" {
		return nil, warnings, errors.New("server CDNName missing")
	} else if server.Cachegroup == nil || *server.Cachegroup == "" {
		return nil, warnings, errors.New("server Cachegroup missing")
	} else if server.Profiles == nil || len(*server.Profiles) == 0 {
		return nil, warnings, errors.New("server Profile missing")
	} else if server.TCPPort == nil {
		return nil, warnings, errors.New("server TCPPort missing")
	}

	// TODO remove, the abstraction shouldn't depend on the ATS version
	atsMajorVer, verWarns := getATSMajorVersion(tcServerParams)
	warnings = append(warnings, verWarns...)

	cacheGroups, err := makeCGMap(cacheGroupArr)
	if err != nil {
		return nil, warnings, errors.New("making CacheGroup map: " + err.Error())
	}
	serverParentCGData, err := getParentCacheGroupData(server, cacheGroups)
	if err != nil {
		return nil, warnings, errors.New("getting server parent cachegroup data: " + err.Error())
	}
	cacheIsTopLevel := isTopLevelCache(serverParentCGData)
	serverCDNDomain := cdn.DomainName

	sort.Sort(dsesSortByName(dses))

	profileParentConfigParams, parentWarns := getProfileParentConfigParams(tcParentConfigParams)
	warnings = append(warnings, parentWarns...)

	parentConfigParamsWithProfiles, err := tcParamsToParamsWithProfiles(tcParentConfigParams)
	if err != nil {
		return nil, warnings, errors.New("adding profiles to parent config params: " + err.Error())
	}

	// parentConfigParams are the parent.config params for all profiles (needed for parents)
	parentConfigParams := parameterWithProfilesToMap(parentConfigParamsWithProfiles)

	// serverParams are the parent.config params for this particular server
	serverParams := getServerParentConfigParams(server, profileParentConfigParams)

	parentCacheGroups := map[string]struct{}{}
	if cacheIsTopLevel {
		for _, cg := range cacheGroups {
			if cg.Type == nil {
				return nil, warnings, errors.New("cachegroup type is nil!")
			}
			if cg.Name == nil {
				return nil, warnings, errors.New("cachegroup name is nil!")
			}

			if *cg.Type != tc.CacheGroupOriginTypeName {
				continue
			}
			parentCacheGroups[*cg.Name] = struct{}{}
		}
	} else {
		for _, cg := range cacheGroups {
			if cg.Type == nil {
				return nil, warnings, errors.New("cachegroup type is nil!")
			}
			if cg.Name == nil {
				return nil, warnings, errors.New("cachegroup name is nil!")
			}

			if *cg.Name == *server.Cachegroup {
				if cg.ParentName != nil && *cg.ParentName != "" {
					parentCacheGroups[*cg.ParentName] = struct{}{}
				}
				if cg.SecondaryParentName != nil && *cg.SecondaryParentName != "" {
					parentCacheGroups[*cg.SecondaryParentName] = struct{}{}
				}
				break
			}
		}
	}

	nameTopologies := makeTopologyNameMap(topologies)

	cgServers := map[int]Server{} // map[serverID]server
	for _, sv := range servers {
		if sv.ID == nil {
			warnings = append(warnings, "TO servers had server with missing ID, skipping!")
			continue
		} else if sv.CDNName == nil {
			warnings = append(warnings, "TO servers had server with missing CDNName, skipping!")
			continue
		} else if sv.Cachegroup == nil || *sv.Cachegroup == "" {
			warnings = append(warnings, "TO servers had server with missing Cachegroup, skipping!")
			continue
		} else if sv.Status == nil || *sv.Status == "" {
			warnings = append(warnings, "TO servers had server with missing Status, skipping!")
			continue
		} else if sv.Type == "" {
			warnings = append(warnings, "TO servers had server with missing Type, skipping!")
			continue
		}
		if *sv.CDNName != *server.CDNName {
			continue
		}
		if _, ok := parentCacheGroups[*sv.Cachegroup]; !ok {
			continue
		}
		if sv.Type != tc.OriginTypeName &&
			!strings.HasPrefix(sv.Type, tc.EdgeTypePrefix) &&
			!strings.HasPrefix(sv.Type, tc.MidTypePrefix) {
			continue
		}
		if *sv.Status != string(tc.CacheStatusReported) && *sv.Status != string(tc.CacheStatusOnline) {
			continue
		}
		cgServers[*sv.ID] = sv
	}

	cgServerIDs := map[int]struct{}{}
	for serverID, _ := range cgServers {
		cgServerIDs[serverID] = struct{}{}
	}
	cgServerIDs[*server.ID] = struct{}{}

	cgDSServers := filterDSS(dss, nil, cgServerIDs)
	parentServerDSes := map[int]map[int]struct{}{} // map[serverID][dsID]
	for _, dss := range cgDSServers {
		if parentServerDSes[dss.Server] == nil {
			parentServerDSes[dss.Server] = map[int]struct{}{}
		}
		parentServerDSes[dss.Server][dss.DeliveryService] = struct{}{}
	}

	originServers, profileCaches, orgProfWarns, err := getOriginServersAndProfileCaches(cgServers, parentServerDSes, profileParentConfigParams, dses, serverCapabilities)
	warnings = append(warnings, orgProfWarns...)
	if err != nil {
		return nil, warnings, errors.New("getting origin servers and profile caches: " + err.Error())
	}

	parentInfos, piWarns := makeParentInfo(serverParentCGData, serverCDNDomain, profileCaches, originServers)
	warnings = append(warnings, piWarns...)

	dsOrigins, dsOriginWarns := makeDSOrigins(dss, dses, servers)
	warnings = append(warnings, dsOriginWarns...)

	for _, ds := range dses {
		if ds.XMLID == nil || *ds.XMLID == "" {
			warnings = append(warnings, "got ds with missing XMLID, skipping!")
			continue
		} else if ds.ID == nil {
			warnings = append(warnings, "got ds with missing ID, skipping!")
			continue
		} else if ds.Type == nil {
			warnings = append(warnings, "got ds with missing Type, skipping!")
			continue
		}

		if !cacheIsTopLevel && ds.Topology == nil {
			if _, ok := parentServerDSes[*server.ID][*ds.ID]; !ok {
				continue // skip DSes not assigned to this server.
			}
		}

		if !ds.Type.IsHTTP() && !ds.Type.IsDNS() {
			continue // skip ANY_MAP, STEERING, etc
		}
		if ds.OrgServerFQDN == nil || *ds.OrgServerFQDN == "" {
			// this check needs to be after the HTTP|DNS check, because Steering DSes without origins are ok'
			warnings = append(warnings, "DS '"+*ds.XMLID+"' has no origin server! Skipping!")
			continue
		}

		// Note these Parameters are only used for MSO for legacy DeliveryServiceServers DeliveryServices (except QueryStringHandling which is used by all DeliveryServices).
		//      Topology DSes use them for all DSes, MSO and non-MSO.
		dsParams, dsParamsWarnings := getParentDSParams(ds, profileParentConfigParams)
		warnings = append(warnings, dsParamsWarnings...)

		// TODO put these in separate functions. No if-statement should be this long.
		if ds.Topology != nil && *ds.Topology != "" {
			txt, topoWarnings, err := getTopologyParentConfigLine(
				server,
				servers,
				&ds,
				serverParams,
				parentConfigParams,
				nameTopologies,
				serverCapabilities,
				dsRequiredCapabilities,
				cacheGroups,
				dsParams,
				atsMajorVer,
				dsOrigins[DeliveryServiceID(*ds.ID)],
				opt.AddComments,
			)
			warnings = append(warnings, topoWarnings...)
			if err != nil {
				// we don't want to fail generation with an error if one ds is malformed
				warnings = append(warnings, err.Error()) // getTopologyParentConfigLine includes error context
				continue
			}

			if txt != nil { // will be nil with no error if this server isn't in the Topology, or if it doesn't have the Required Capabilities
				parentAbstraction.Services = append(parentAbstraction.Services, txt)
			}
		} else if cacheIsTopLevel {
			parentQStr := false
			if dsParams.QueryStringHandling == "" && dsParams.Algorithm == tc.AlgorithmConsistentHash && ds.QStringIgnore != nil && tc.QStringIgnore(*ds.QStringIgnore) == tc.QStringIgnoreUseInCacheKeyAndPassUp {
				parentQStr = true
			}

			orgFQDNStr := *ds.OrgServerFQDN
			// if this cache isn't the last tier, i.e. we're not going to the origin, use http not https
			if isLastCacheTier := noTopologyServerIsLastCacheForDS(server, &ds); !isLastCacheTier {
				orgFQDNStr = strings.Replace(orgFQDNStr, `https://`, `http://`, -1)
			}
			orgURI, orgWarns, err := getOriginURI(orgFQDNStr)
			warnings = append(warnings, orgWarns...)
			if err != nil {
				warnings = append(warnings, "DS '"+*ds.XMLID+"' has malformed origin URI: '"+orgFQDNStr+"': skipping!"+err.Error())
				continue
			}

			textLine := &ParentAbstractionService{}
			textLine.Name = *ds.XMLID

			if ds.OriginShield != nil && *ds.OriginShield != "" {

				policy := ParentAbstractionServiceRetryPolicyConsistentHash
				if parentSelectAlg := serverParams[ParentConfigParamAlgorithm]; strings.TrimSpace(parentSelectAlg) != "" {
					paramPolicy := ParentSelectAlgorithmToParentAbstractionServiceRetryPolicy(parentSelectAlg)
					if paramPolicy != ParentAbstractionServiceRetryPolicyInvalid {
						policy = paramPolicy
					} else {
						warnings = append(warnings, "DS '"+*ds.XMLID+"' had malformed "+ParentConfigParamAlgorithm+" parameter '"+parentSelectAlg+"', not using!")
					}
				}
				textLine.Comment = makeParentComment(opt.AddComments, *ds.XMLID, "")
				textLine.DestDomain = orgURI.Hostname()
				textLine.Port, err = strconv.Atoi(orgURI.Port())
				if err != nil {
					if strings.ToLower(orgURI.Scheme) == "https" {
						textLine.Port = 443
					} else {
						textLine.Port = 80
					}
					warnings = append(warnings, "DS '"+*ds.XMLID+"' had malformed origin  port: '"+orgURI.Port()+"': using "+strconv.Itoa(textLine.Port)+"! : "+err.Error())
				}

				fqdnPort := strings.Split(*ds.OriginShield, ":")
				parent := &ParentAbstractionServiceParent{}
				parent.FQDN = fqdnPort[0]
				if len(fqdnPort) > 1 {
					parent.Port, err = strconv.Atoi(fqdnPort[1])
					if err != nil {
						parent.Port = 80
						warnings = append(warnings, "DS '"+*ds.XMLID+"' had malformed origin  port: '"+*ds.OriginShield+"': using "+strconv.Itoa(parent.Port)+"! : "+err.Error())
					}
				} else {
					parent.Port = 80
					warnings = append(warnings, "DS '"+*ds.XMLID+"' had no origin port: '"+*ds.OriginShield+"': using "+strconv.Itoa(parent.Port)+"!")
				}
				textLine.Parents = append(textLine.Parents, parent)
				textLine.RetryPolicy = policy
				textLine.GoDirect = true

				// textLine += "dest_domain=" + orgURI.Hostname() + " port=" + orgURI.Port() + " parent=" + *ds.OriginShield + " " + algorithm + " go_direct=true\n"

			} else if ds.MultiSiteOrigin != nil && *ds.MultiSiteOrigin {
				textLine.Comment = makeParentComment(opt.AddComments, *ds.XMLID, "")
				textLine.DestDomain = orgURI.Hostname()
				textLine.Port, err = strconv.Atoi(orgURI.Port())
				if err != nil {
					textLine.Port = 80
					warnings = append(warnings, "DS '"+*ds.XMLID+"' had malformed origin  port: '"+orgURI.Port()+"': using "+strconv.Itoa(textLine.Port)+"! : "+err.Error())
				}

				// textLine += "dest_domain=" + orgURI.Hostname() + " port=" + orgURI.Port() + " "
				if len(parentInfos) == 0 {
				}

				if len(parentInfos[OriginHost(orgURI.Hostname())]) == 0 {
					// TODO error? emulates Perl
					warnings = append(warnings, "DS "+*ds.XMLID+" has no parent servers")
				}

				parents, secondaryParents, secondaryMode, parentWarns := getMSOParentStrs(&ds, parentInfos[OriginHost(orgURI.Hostname())], atsMajorVer, dsParams.Algorithm, dsParams.TryAllPrimariesBeforeSecondary)
				warnings = append(warnings, parentWarns...)
				textLine.Parents = parents
				textLine.SecondaryParents = secondaryParents
				textLine.SecondaryMode = secondaryMode
				textLine.RetryPolicy = dsParams.Algorithm // TODO convert
				textLine.IgnoreQueryStringInParentSelection = !parentQStr
				textLine.GoDirect = true

				// textLine += parents + secondaryParents + ` round_robin=` + dsParams.Algorithm + ` qstring=` + parentQStr + ` go_direct=false parent_is_proxy=false`
				prWarns := []string{}
				textLine.MaxSimpleRetries, textLine.MaxMarkdownRetries, textLine.MarkdownResponseCodes, textLine.ErrorResponseCodes, prWarns = getParentRetryStr(true, atsMajorVer, dsParams.ParentRetry, dsParams.UnavailableServerRetryResponses, dsParams.MaxSimpleRetries, dsParams.MaxUnavailableServerRetries)
				warnings = append(warnings, prWarns...)

				parentAbstraction.Services = append(parentAbstraction.Services, textLine)
			}
		} else {
			queryStringHandling := ParentSelectParamQStringHandlingToBool(serverParams[ParentConfigParamQStringHandling]) // "qsh" in Perl
			if queryStringHandling == nil && serverParams[ParentConfigParamQStringHandling] != "" {
				warnings = append(warnings, "Server Parameter '"+ParentConfigParamQStringHandling+"' value '"+serverParams[ParentConfigParamQStringHandling]+"' malformed, not using!")
			}

			roundRobin := ParentAbstractionServiceRetryPolicyConsistentHash
			// roundRobin := `round_robin=consistent_hash`
			goDirect := false
			// goDirect := `go_direct=false`

			parents, secondaryParents, secondaryMode, parentWarns := getParentStrs(&ds, dsRequiredCapabilities, parentInfos[deliveryServicesAllParentsKey], atsMajorVer, dsParams.TryAllPrimariesBeforeSecondary)
			warnings = append(warnings, parentWarns...)

			text := &ParentAbstractionService{}
			text.Name = *ds.XMLID

			orgFQDNStr := *ds.OrgServerFQDN
			// if this cache isn't the last tier, i.e. we're not going to the origin, use http not https
			if isLastCacheTier := noTopologyServerIsLastCacheForDS(server, &ds); !isLastCacheTier {
				orgFQDNStr = strings.Replace(orgFQDNStr, `https://`, `http://`, -1)
			}
			orgURI, orgWarns, err := getOriginURI(orgFQDNStr)
			warnings = append(warnings, orgWarns...)
			if err != nil {
				warnings = append(warnings, "DS '"+*ds.XMLID+"' had malformed origin  URI: '"+*ds.OrgServerFQDN+"': skipping!"+err.Error())
				continue
			}

			text.Comment = makeParentComment(opt.AddComments, *ds.XMLID, "")

			// TODO encode this in a DSType func, IsGoDirect() ?
			if *ds.Type == tc.DSTypeHTTPNoCache || *ds.Type == tc.DSTypeHTTPLive || *ds.Type == tc.DSTypeDNSLive {
				text.DestDomain = orgURI.Hostname()
				text.Port, err = strconv.Atoi(orgURI.Port())
				if err != nil {
					if strings.ToLower(orgURI.Scheme) == "https" {
						text.Port = 443
					} else {
						text.Port = 80
					}
					warnings = append(warnings, "DS '"+*ds.XMLID+"' had malformed origin  port: '"+orgURI.Port()+"': using "+strconv.Itoa(text.Port)+"! : "+err.Error())
				}
				text.GoDirect = true
				// text += `dest_domain=` + orgURI.Hostname() + ` port=` + orgURI.Port() + ` go_direct=true` + "\n"
			} else {

				// check for profile psel.qstring_handling.  If this parameter is assigned to the server profile,
				// then edges will use the qstring handling value specified in the parameter for all profiles.

				// If there is no defined parameter in the profile, then check the delivery service profile.
				// If psel.qstring_handling exists in the DS profile, then we use that value for the specified DS only.
				// This is used only if not overridden by a server profile qstring handling parameter.

				// TODO refactor this logic, hard to understand (transliterated from Perl)
				dsQSH := queryStringHandling
				if dsQSH == nil {
					dsQSH = ParentSelectParamQStringHandlingToBool(dsParams.QueryStringHandling)
					if dsQSH == nil && dsParams.QueryStringHandling != "" {
						warnings = append(warnings, "Delivery Service parameter '"+ParentConfigParamQStringHandling+"' value '"+dsParams.QueryStringHandling+"' malformed, not using!")
					}

				}
				parentQStr := dsQSH
				if parentQStr == nil {
					v := false
					parentQStr = &v
				}
				if ds.QStringIgnore != nil && tc.QStringIgnore(*ds.QStringIgnore) == tc.QStringIgnoreUseInCacheKeyAndPassUp && dsQSH == nil {
					v := true
					parentQStr = &v
				}
				if parentQStr == nil {
					b := !DefaultIgnoreQueryStringInParentSelection
					parentQStr = &b
				}

				text.DestDomain = orgURI.Hostname()
				text.Port, err = strconv.Atoi(orgURI.Port())
				if err != nil {
					if strings.ToLower(orgURI.Scheme) == "https" {
						text.Port = 443
					} else {
						text.Port = 80
					}
					warnings = append(warnings, "DS '"+*ds.XMLID+"' had malformed origin  port: '"+orgURI.Port()+"': using "+strconv.Itoa(text.Port)+"! : "+err.Error())
				}
				text.Parents = parents
				text.SecondaryParents = secondaryParents
				text.SecondaryMode = secondaryMode
				text.RetryPolicy = roundRobin
				text.GoDirect = goDirect
				text.IgnoreQueryStringInParentSelection = !*parentQStr
				// text += `dest_domain=` + orgURI.Hostname() + ` port=` + orgURI.Port() + ` ` + parents + ` ` + secondaryParents + ` ` + roundRobin + ` ` + goDirect + ` qstring=` + parentQStr + "\n"
			}

			parentAbstraction.Services = append(parentAbstraction.Services, text)
		}
	}

	// TODO determine if this is necessary. It's super-dangerous, and moreover ignores Server Capabilitites.
	defaultDestText := (*ParentAbstractionService)(nil)
	if !isTopLevelCache(serverParentCGData) {
		defaultDestText = &ParentAbstractionService{}
		// magic uuid to prevent accidental DS name collision
		defaultDestText.Name = `default-destination-c3854be4-a859-41d6-815d-7b36297e48c6`
		invalidDS := &DeliveryService{}
		invalidDS.ID = util.IntPtr(-1)
		tryAllPrimariesBeforeSecondary := false
		parents, secondaryParents, secondaryMode, parentWarns := getParentStrs(invalidDS, dsRequiredCapabilities, parentInfos[deliveryServicesAllParentsKey], atsMajorVer, tryAllPrimariesBeforeSecondary)
		warnings = append(warnings, parentWarns...)

		defaultDestText.DestDomain = `.`
		defaultDestText.Parents = parents
		// defaultDestText = `dest_domain=. ` + parents
		if serverParams[ParentConfigParamAlgorithm] == tc.AlgorithmConsistentHash {
			defaultDestText.SecondaryParents = secondaryParents
			defaultDestText.SecondaryMode = secondaryMode
			// defaultDestText += secondaryParents
		}
		defaultDestText.RetryPolicy = ParentAbstractionServiceRetryPolicyConsistentHash
		defaultDestText.GoDirect = false
		// defaultDestText += ` round_robin=consistent_hash go_direct=false`

		if qStr := serverParams[ParentConfigParamQString]; qStr != "" {
			if v := ParentSelectParamQStringHandlingToBool(qStr); v != nil {
				defaultDestText.IgnoreQueryStringInParentSelection = !*v
			} else if qStr != "" {
				warnings = append(warnings, "Server parameter '"+ParentConfigParamQString+"' value '"+qStr+"' malformed, not using!")
			}
			// defaultDestText += ` qstring=` + qStr
		}
		defaultDestText.Comment = makeParentComment(opt.AddComments, "", "")
	}

	sort.Sort(ParentAbstractionServices(parentAbstraction.Services))
	if defaultDestText != nil {
		parentAbstraction.Services = append(parentAbstraction.Services, defaultDestText)
	}

	return parentAbstraction, warnings, nil
}

// makeParentComment creates the parent line comment and returns it.
// If addComments is false, returns the empty string. This exists for composability.
// Either dsName or topology may be the empty string.
// The returned comment includes a trailing newline.
func makeParentComment(addComments bool, dsName string, topology string) string {
	if !addComments {
		return ""
	}
	return "ds '" + dsName + "' topology '" + topology + "'"
}

type parentConfigDS struct {
	Name                 tc.DeliveryServiceName
	QStringIgnore        tc.QStringIgnore
	OriginFQDN           string
	MultiSiteOrigin      bool
	OriginShield         string
	Type                 tc.DSType
	QStringHandling      string
	RequiredCapabilities map[ServerCapability]struct{}
	Topology             string
}

type parentConfigDSTopLevel struct {
	parentConfigDS
	MSOAlgorithm                       string
	MSOParentRetry                     string
	MSOUnavailableServerRetryResponses string
	MSOMaxSimpleRetries                string
	MSOMaxUnavailableServerRetries     string
}

type parentInfo struct {
	Host            string
	Port            int
	Domain          string
	Weight          float64
	UseIP           bool
	Rank            int
	IP              string
	PrimaryParent   bool
	SecondaryParent bool
	Capabilities    map[ServerCapability]struct{}
}

func (p parentInfo) Format() string {
	host := ""
	if p.UseIP {
		host = p.IP
	} else {
		host = p.Host + "." + p.Domain
	}
	return host + ":" + strconv.Itoa(p.Port) + "|" + strconv.FormatFloat(p.Weight, 'f', 3, 64) + ";"
}

func (p parentInfo) ToAbstract() *ParentAbstractionServiceParent {
	host := ""
	if p.UseIP {
		host = p.IP
	} else {
		host = p.Host + "." + p.Domain
	}
	return &ParentAbstractionServiceParent{
		FQDN:   host,
		Port:   p.Port,
		Weight: p.Weight,
	}
}

type parentInfos map[OriginHost]parentInfo

type parentInfoSortByRank []parentInfo

func (s parentInfoSortByRank) Len() int      { return len(s) }
func (s parentInfoSortByRank) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s parentInfoSortByRank) Less(i, j int) bool {
	if s[i].Rank != s[j].Rank {
		return s[i].Rank < s[j].Rank
	} else if s[i].Host != s[j].Host {
		return s[i].Host < s[j].Host
	} else if s[i].Domain != s[j].Domain {
		return s[i].Domain < s[j].Domain
	} else if s[i].Port != s[j].Port {
		return s[i].Port < s[j].Port
	}
	return s[i].IP < s[j].IP
}

type serverWithParams struct {
	Server
	Params profileCache
}

type serversWithParamsSortByRank []serverWithParams

func (ss serversWithParamsSortByRank) Len() int      { return len(ss) }
func (ss serversWithParamsSortByRank) Swap(i, j int) { ss[i], ss[j] = ss[j], ss[i] }
func (ss serversWithParamsSortByRank) Less(i, j int) bool {
	if ss[i].Params.Rank != ss[j].Params.Rank {
		return ss[i].Params.Rank < ss[j].Params.Rank
	}

	if ss[i].HostName == nil {
		if ss[j].HostName != nil {
			return true
		}
	} else if ss[j].HostName == nil {
		return false
	} else if ss[i].HostName != ss[j].HostName {
		return *ss[i].HostName < *ss[j].HostName
	}

	if ss[i].DomainName == nil {
		if ss[j].DomainName != nil {
			return true
		}
	} else if ss[j].DomainName == nil {
		return false
	} else if ss[i].DomainName != ss[j].DomainName {
		return *ss[i].DomainName < *ss[j].DomainName
	}

	if ss[i].Params.Port != ss[j].Params.Port {
		return ss[i].Params.Port < ss[j].Params.Port
	}

	iIP := getServerIPAddress(&ss[i].Server)
	jIP := getServerIPAddress(&ss[j].Server)

	if iIP == nil {
		if jIP != nil {
			return true
		}
	} else if jIP == nil {
		return false
	}
	return bytes.Compare(iIP, jIP) <= 0
}

type parentConfigDSTopLevelSortByName []parentConfigDSTopLevel

func (s parentConfigDSTopLevelSortByName) Len() int      { return len(s) }
func (s parentConfigDSTopLevelSortByName) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s parentConfigDSTopLevelSortByName) Less(i, j int) bool {
	return strings.Compare(string(s[i].Name), string(s[j].Name)) < 0
}

type dsesSortByName []DeliveryService

func (s dsesSortByName) Len() int      { return len(s) }
func (s dsesSortByName) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s dsesSortByName) Less(i, j int) bool {
	if s[i].XMLID == nil {
		return true
	}
	if s[j].XMLID == nil {
		return false
	}
	return *s[i].XMLID < *s[j].XMLID
}

type profileCache struct {
	Weight     string
	Port       int
	UseIP      bool
	Rank       int
	NotAParent bool
}

const DefaultParentWeight = 0.999

func defaultProfileCache() profileCache {
	return profileCache{
		Weight:     strconv.FormatFloat(DefaultParentWeight, 'f', 3, 64),
		Port:       0,
		UseIP:      false,
		Rank:       1,
		NotAParent: false,
	}
}

// cgServer is the server table data needed when selecting the servers assigned to a cachegroup.
type cgServer struct {
	ServerID       ServerID
	ServerHost     string
	ServerIP       string
	ServerPort     int
	CacheGroupID   int
	CacheGroupName string
	Status         int
	Type           int
	ProfileID      ProfileID
	ProfileName    string
	CDN            int
	TypeName       string
	Domain         string
	Capabilities   map[ServerCapability]struct{}
}

type originURI struct {
	Scheme string
	Host   string
	Port   string
}

// TODO change, this is terrible practice, using a hard-coded key. What if there were a delivery service named "all_parents" (transliterated Perl)
const deliveryServicesAllParentsKey = "all_parents"

type parentDSParams struct {
	Algorithm                       ParentAbstractionServiceRetryPolicy
	ParentRetry                     string
	UnavailableServerRetryResponses string
	MaxSimpleRetries                string
	MaxUnavailableServerRetries     string
	QueryStringHandling             string
	TryAllPrimariesBeforeSecondary  bool
	MergeGroups                     []string
}

// getDSParams returns the Delivery Service Profile Parameters used in parent.config, and any warnings.
// If Parameters don't exist, defaults are returned. Non-MSO Delivery Services default to no custom retry logic (we should reevaluate that).
// Note these Parameters are only used for MSO for legacy DeliveryServiceServers DeliveryServices.
//      Topology DSes use them for all DSes, MSO and non-MSO.
func getParentDSParams(ds DeliveryService, profileParentConfigParams map[string]map[string]string) (parentDSParams, []string) {
	warnings := []string{}
	params := parentDSParams{}
	isMSO := ds.MultiSiteOrigin != nil && *ds.MultiSiteOrigin
	if isMSO {
		params.Algorithm = ParentConfigDSParamDefaultMSOAlgorithm
		params.ParentRetry = ParentConfigDSParamDefaultMSOParentRetry
		params.UnavailableServerRetryResponses = ParentConfigDSParamDefaultMSOUnavailableServerRetryResponses
		params.MaxSimpleRetries = strconv.Itoa(ParentConfigDSParamDefaultMaxSimpleRetries)
		params.MaxUnavailableServerRetries = strconv.Itoa(ParentConfigDSParamDefaultMaxUnavailableServerRetries)
	}
	if ds.ProfileName == nil || *ds.ProfileName == "" {
		return params, warnings
	}
	dsParams, ok := profileParentConfigParams[*ds.ProfileName]
	if !ok {
		return params, warnings
	}

	// the following may be blank, no default
	params.QueryStringHandling = dsParams[ParentConfigParamQStringHandling]
	params.MergeGroups = strings.Split(dsParams[ParentConfigParamMergeGroups], " ")
	if 0 < len(params.MergeGroups) {
		sort.Strings(params.MergeGroups)
	}

	// TODO deprecate & remove "mso." Parameters - there was never a reason to restrict these settings to MSO.
	if isMSO {
		if v, ok := dsParams[ParentConfigParamMSOAlgorithm]; ok && strings.TrimSpace(v) != "" {
			policy := ParentSelectAlgorithmToParentAbstractionServiceRetryPolicy(v)
			if policy != ParentAbstractionServiceRetryPolicyInvalid {
				params.Algorithm = policy
			} else {
				warnings = append(warnings, "DS '"+*ds.XMLID+"' had malformed "+ParentConfigParamMSOAlgorithm+" parameter '"+v+"', not using!")
			}
		}
		if v, ok := dsParams[ParentConfigParamMSOParentRetry]; ok {
			params.ParentRetry = v
		}
		if v, ok := dsParams[ParentConfigParamMSOUnavailableServerRetryResponses]; ok {
			if v != "" && !unavailableServerRetryResponsesValid(v) {
				warnings = append(warnings, "DS '"+*ds.XMLID+"' had malformed "+ParentConfigParamMSOUnavailableServerRetryResponses+" parameter '"+v+"', not using!")
			} else if v != "" {
				params.UnavailableServerRetryResponses = v
			}
		}
		if v, ok := dsParams[ParentConfigParamMSOMaxSimpleRetries]; ok {
			params.MaxSimpleRetries = v
		}
		if v, ok := dsParams[ParentConfigParamMSOMaxUnavailableServerRetries]; ok {
			params.MaxUnavailableServerRetries = v
		}
	}

	// Even if the DS is MSO, non-"mso." Parameters override "mso." ones, because they're newer.
	if v, ok := dsParams[ParentConfigParamAlgorithm]; ok && strings.TrimSpace(v) != "" {
		policy := ParentSelectAlgorithmToParentAbstractionServiceRetryPolicy(v)
		if policy != ParentAbstractionServiceRetryPolicyInvalid {
			params.Algorithm = policy
		} else {
			warnings = append(warnings, "DS '"+*ds.XMLID+"' had malformed "+ParentConfigParamAlgorithm+" parameter '"+v+"', not using!")
		}
	}
	if v, ok := dsParams[ParentConfigParamParentRetry]; ok {
		params.ParentRetry = v
	}
	if v, ok := dsParams[ParentConfigParamUnavailableServerRetryResponses]; ok {
		if v != "" && !unavailableServerRetryResponsesValid(v) {
			warnings = append(warnings, "DS '"+*ds.XMLID+"' had malformed "+ParentConfigParamUnavailableServerRetryResponses+" parameter '"+v+"', not using!")
		} else if v != "" {
			params.UnavailableServerRetryResponses = v
		}
	}
	if v, ok := dsParams[ParentConfigParamMaxSimpleRetries]; ok {
		params.MaxSimpleRetries = v
	}
	if v, ok := dsParams[ParentConfigParamMaxUnavailableServerRetries]; ok {
		params.MaxUnavailableServerRetries = v
	}
	if v, ok := dsParams[ParentConfigParamSecondaryMode]; ok {
		if v != "" {
			warnings = append(warnings, "DS '"+*ds.XMLID+"' had Parameter "+ParentConfigParamSecondaryMode+" which is used if it exists, the value is ignored! Non-empty value '"+v+"' will be ignored!")
		}
		params.TryAllPrimariesBeforeSecondary = true
	}

	return params, warnings
}

// getTopologyParentConfigLine returns the topology parent.config line, any warnings, and any error
// If the given DS is not used by the server, returns a nil ParentAbstractionService and nil error.
func getTopologyParentConfigLine(
	server *Server,
	servers []Server,
	ds *DeliveryService,
	serverParams map[string]string,
	parentConfigParams []parameterWithProfilesMap, // all params with configFile parent.config
	nameTopologies map[TopologyName]tc.Topology,
	serverCapabilities map[int]map[ServerCapability]struct{},
	dsRequiredCapabilities map[int]map[ServerCapability]struct{},
	cacheGroups map[tc.CacheGroupName]tc.CacheGroupNullable,
	dsParams parentDSParams,
	atsMajorVer int,
	dsOrigins map[ServerID]struct{},
	addComments bool,
) (*ParentAbstractionService, []string, error) {
	warnings := []string{}
	if !hasRequiredCapabilities(serverCapabilities[*server.ID], dsRequiredCapabilities[*ds.ID]) {
		return nil, warnings, nil
	}

	topology := nameTopologies[TopologyName(*ds.Topology)]
	if topology.Name == "" {
		return nil, warnings, errors.New("DS " + *ds.XMLID + " topology '" + *ds.Topology + "' not found in Topologies!")
	}

	serverPlacement, err := getTopologyPlacement(tc.CacheGroupName(*server.Cachegroup), topology, cacheGroups, ds)

	if err != nil {
		return nil, warnings, errors.New("getting topology placement: " + err.Error())
	}
	if !serverPlacement.InTopology {
		return nil, warnings, nil // server isn't in topology, no error
	}

	orgFQDNStr := *ds.OrgServerFQDN
	// if this cache isn't the last tier, i.e. we're not going to the origin, use http not https
	if !serverPlacement.IsLastCacheTier {
		orgFQDNStr = strings.Replace(orgFQDNStr, `https://`, `http://`, -1)
	}
	orgURI, orgWarns, err := getOriginURI(orgFQDNStr)
	warnings = append(warnings, orgWarns...)
	if err != nil {
		return nil, warnings, errors.New("DS '" + *ds.XMLID + "' has malformed origin URI: '" + *ds.OrgServerFQDN + "': skipping!" + err.Error())
	}

	txt := &ParentAbstractionService{}
	txt.Name = *ds.XMLID
	txt.Comment = makeParentComment(addComments, *ds.XMLID, *ds.Topology)
	txt.DestDomain = orgURI.Hostname()
	txt.Port, err = strconv.Atoi(orgURI.Port())
	if err != nil {
		return nil, warnings, fmt.Errorf("parent %v port '%v' was not an integer", orgURI, orgURI.Port())
	}
	// txt += "dest_domain=" + orgURI.Hostname() + " port=" + orgURI.Port()

	parents, secondaryParents, parentWarnings, err := getTopologyParents(server, ds, servers, parentConfigParams, topology, serverPlacement.IsLastTier, serverCapabilities, dsRequiredCapabilities, dsOrigins, dsParams.MergeGroups)
	warnings = append(warnings, parentWarnings...)
	if err != nil {
		return nil, warnings, errors.New("getting topology parents for '" + *ds.XMLID + "': skipping! " + err.Error())
	}
	if len(parents) == 0 {
		return nil, warnings, errors.New("getting topology parents for '" + *ds.XMLID + "': no parents found! skipping! (Does your Topology have a CacheGroup with no servers in it?)")
	}

	txt.Parents = parents
	// txt += ` parent="` + strings.Join(parents, `;`) + `"`
	if len(secondaryParents) > 0 {
		txt.SecondaryParents = secondaryParents
		// txt += ` secondary_parent="` + strings.Join(secondaryParents, `;`) + `"`

		secondaryModeStr, secondaryModeWarnings := getSecondaryModeStr(dsParams.TryAllPrimariesBeforeSecondary, atsMajorVer, tc.DeliveryServiceName(*ds.XMLID))
		warnings = append(warnings, secondaryModeWarnings...)
		// txt += secondaryModeStr
		txt.SecondaryMode = secondaryModeStr // TODO convert
	}

	txt.RetryPolicy = getTopologyRoundRobin(ds, serverParams, serverPlacement.IsLastCacheTier, dsParams.Algorithm)
	// txt += ` round_robin=` + getTopologyRoundRobin(ds, serverParams, serverPlacement.IsLastCacheTier, dsParams.Algorithm)

	txt.GoDirect = getTopologyGoDirect(ds, serverPlacement.IsLastTier, serverPlacement.IsLastCacheTier)
	// txt += ` go_direct=` + getTopologyGoDirect(ds, serverPlacement.IsLastTier)

	// TODO convert
	useQueryStringInParentSelection := (*bool)(nil)
	if dsParams.QueryStringHandling != "" {
		qs := ParentSelectParamQStringHandlingToBool(dsParams.QueryStringHandling)
		if qs != nil {
			useQueryStringInParentSelection = qs
		} else if dsParams.QueryStringHandling != "" {
			warnings = append(warnings, fmt.Sprintf("DS '"+*ds.XMLID+"' has malformed query string handling param '"+dsParams.QueryStringHandling+"', using default %v", useQueryStringInParentSelection))
		}
	}

	tqWarns := []string{}
	txt.IgnoreQueryStringInParentSelection, tqWarns = getTopologyQueryStringIgnore(ds, serverParams, serverPlacement.IsLastCacheTier, dsParams.Algorithm, useQueryStringInParentSelection)
	warnings = append(warnings, tqWarns...)

	// txt += ` qstring=` + getTopologyQueryString(ds, serverParams, serverPlacement.IsLastCacheTier, dsParams.Algorithm, dsParams.QueryStringHandling)

	// TODO ensure value is always !goDirect, and determine what to do if not
	// txt += getTopologyParentIsProxyStr(serverPlacement.IsLastCacheTier)

	// TODO convert
	prWarns := []string{}
	txt.MaxSimpleRetries, txt.MaxMarkdownRetries, txt.MarkdownResponseCodes, txt.ErrorResponseCodes, prWarns = getParentRetryStr(serverPlacement.IsLastCacheTier, atsMajorVer, dsParams.ParentRetry, dsParams.UnavailableServerRetryResponses, dsParams.MaxSimpleRetries, dsParams.MaxUnavailableServerRetries)
	warnings = append(warnings, prWarns...)

	// txt += getParentRetryStr(serverPlacement.IsLastCacheTier, atsMajorVer, dsParams.ParentRetry, dsParams.UnavailableServerRetryResponses, dsParams.MaxSimpleRetries, dsParams.MaxUnavailableServerRetries)
	// txt += "\n"

	return txt, warnings, nil
}

// getParentRetryStr builds the parent retry directive(s).
//
// Returns the MaxSimpleRetries, MaxMarkdownRetries, MarkdownResponseCodes, and ErrorResponseCodes.
//
// If atsMajorVer < 6, "" is returned (ATS 5 and below don't support retry directives).
// If isLastCacheTier is false, "" is returned. This argument exists to simplify usage.
// If parentRetry is "", "" is returned (because the other directives are unused if parent_retry doesn't exist). This is allowed to simplify usage.
// If unavailableServerRetryResponses is not "", it must be valid. Use unavailableServerRetryResponsesValid to check.
// If maxSimpleRetries is "", ParentConfigDSParamDefaultMaxSimpleRetries will be used.
// If maxUnavailableServerRetries is "", ParentConfigDSParamDefaultMaxUnavailableServerRetries will be used.
//
// Does not return errors. If any input is malformed, warnings are returned and that value is set to -1.
//
func getParentRetryStr(isLastCacheTier bool, atsMajorVer int, parentRetry string, unavailableServerRetryResponses string, maxSimpleRetries string, maxUnavailableServerRetries string) (int, int, []int, []int, []string) {
	warnings := []string{}
	if !isLastCacheTier || // allow !isLastCacheTier, to simplify usage.
		parentRetry == "" || // allow parentRetry to be empty, to simplify usage.
		atsMajorVer < 6 { // ATS 5 and below don't support parent_retry directives
		// warnings = append(warnings, "ATS 5 doesn't support parent retry, not using parent retry values")
		return -1, -1, nil, nil, warnings // TODO move to formatter?
	}

	err := error(nil)

	maxSimpleRetriesInt := ParentConfigDSParamDefaultMaxSimpleRetries
	if maxSimpleRetries != "" {
		maxSimpleRetriesInt, err = strconv.Atoi(maxSimpleRetries)
		if err != nil {
			maxSimpleRetriesInt = ParentConfigDSParamDefaultMaxSimpleRetries
			warnings = append(warnings, "malformed maxSimpleRetries '"+maxSimpleRetries+"', using default "+strconv.Itoa(maxSimpleRetriesInt))
		}
	}

	maxUnavailableServerRetriesInt := ParentConfigDSParamDefaultMaxUnavailableServerRetries
	if maxUnavailableServerRetries != "" {
		maxUnavailableServerRetriesInt, err = strconv.Atoi(maxUnavailableServerRetries)
		if err != nil {
			maxUnavailableServerRetriesInt = ParentConfigDSParamDefaultMaxUnavailableServerRetries
			warnings = append(warnings, "malformed maxUnavailableServerRetries '"+maxUnavailableServerRetries+"', using default "+strconv.Itoa(maxUnavailableServerRetriesInt))
		}
	}

	unavailableServerRetryResponsesInts, err := ParseRetryResponses(unavailableServerRetryResponses)
	if err != nil {
		warnings = append(warnings, "malformed unavailableServerRetryResponses '"+unavailableServerRetryResponses+"', using default")
		unavailableServerRetryResponsesInts = []int{}
	}

	simpleRetryResponsesInts := []int{}
	// TODO add support for 9.1
	// simpleRetryResponsesInts, err := ParseRetryResponses(simpleRetryResponses)
	// if err != nil {
	// 	warnings = append(warnings, "malformed simpleRetryResponses '"+simpleRetryResponses+"', using default")
	// 	simpleRetryResponsesInts = []int{}
	// }

	// TODO make consts
	switch strings.ToLower(strings.TrimSpace(parentRetry)) {
	case "simple_retry":
		unavailableServerRetryResponsesInts = []int{}
		if len(simpleRetryResponsesInts) == 0 {
			simpleRetryResponsesInts = append(simpleRetryResponsesInts, DefaultSimpleRetryCodes...)
		}
	case "unavailable_server_retry":
		simpleRetryResponsesInts = []int{}
		if len(unavailableServerRetryResponsesInts) == 0 {
			unavailableServerRetryResponsesInts = append(unavailableServerRetryResponsesInts, DefaultUnavailableServerRetryCodes...)
		}
	case "both":
		if len(unavailableServerRetryResponsesInts) == 0 {
			unavailableServerRetryResponsesInts = append(unavailableServerRetryResponsesInts, DefaultUnavailableServerRetryCodes...)
		}
		if len(simpleRetryResponsesInts) == 0 {
			simpleRetryResponsesInts = append(simpleRetryResponsesInts, DefaultSimpleRetryCodes...)
		}
	default:
		unavailableServerRetryResponsesInts = []int{}
		simpleRetryResponsesInts = []int{}
	}

	// txt := ` parent_retry=` + parentRetry
	// if unavailableServerRetryResponses != "" {
	// 	txt += ` unavailable_server_retry_responses=` + unavailableServerRetryResponses
	// }
	// txt += ` max_simple_retries=` + maxSimpleRetries + ` max_unavailable_server_retries=` + maxUnavailableServerRetries
	return maxSimpleRetriesInt, maxUnavailableServerRetriesInt, unavailableServerRetryResponsesInts, simpleRetryResponsesInts, warnings
}

// getSecondaryModeStr returns the secondary_mode string, and any warnings.
func getSecondaryModeStr(tryAllPrimariesBeforeSecondary bool, atsMajorVer int, ds tc.DeliveryServiceName) (ParentAbstractionServiceParentSecondaryMode, []string) {
	warnings := []string{}
	if !tryAllPrimariesBeforeSecondary {
		return ParentAbstractionServiceParentSecondaryModeDefault, warnings
	}
	if atsMajorVer < 8 {
		warnings = append(warnings, "DS '"+string(ds)+"' had Parameter "+ParentConfigParamSecondaryMode+" but this cache is "+strconv.Itoa(atsMajorVer)+" and secondary_mode isn't supported in ATS until 8. Not using!")
		return ParentAbstractionServiceParentSecondaryModeDefault, warnings
	}

	// See https://docs.trafficserver.apache.org/en/8.0.x/admin-guide/files/parent.config.en.html
	return ParentAbstractionServiceParentSecondaryModeExhaust, warnings
}

func getTopologyParentIsProxyStr(serverIsLastCacheTier bool) string {
	if serverIsLastCacheTier {
		return ` parent_is_proxy=false`
	}
	return ""
}

// RetryPolicy
func getTopologyRoundRobin(
	ds *DeliveryService,
	serverParams map[string]string,
	serverIsLastTier bool,
	algorithm ParentAbstractionServiceRetryPolicy,
) ParentAbstractionServiceRetryPolicy {
	if !serverIsLastTier {
		return ParentAbstractionServiceRetryPolicyConsistentHash
	}
	if parentSelectAlg := serverParams[ParentConfigParamAlgorithm]; ds.OriginShield != nil && *ds.OriginShield != "" && strings.TrimSpace(parentSelectAlg) != "" {
		if policy := ParentSelectAlgorithmToParentAbstractionServiceRetryPolicy(parentSelectAlg); policy != ParentAbstractionServiceRetryPolicyInvalid {
			return policy
		}
	}
	if algorithm != "" {
		return algorithm
	}
	return ParentAbstractionServiceRetryPolicyConsistentHash
}

func getTopologyGoDirect(ds *DeliveryService, serverIsLastTier bool, serverIsLastCacheTier bool) bool {
	if serverIsLastCacheTier {
		return true
	}
	if !serverIsLastTier {
		return false
	}
	if ds.OriginShield != nil && *ds.OriginShield != "" {
		return true
	}
	if ds.MultiSiteOrigin != nil && *ds.MultiSiteOrigin {
		return false
	}
	return true
}

func getTopologyQueryStringIgnore(
	ds *DeliveryService,
	serverParams map[string]string,
	serverIsLastTier bool,
	algorithm ParentAbstractionServiceRetryPolicy,
	qStringHandling *bool,
) (bool, []string) {
	warnings := []string{}
	if serverIsLastTier {
		if ds.MultiSiteOrigin != nil && *ds.MultiSiteOrigin && qStringHandling == nil && algorithm == ParentAbstractionServiceRetryPolicyConsistentHash && ds.QStringIgnore != nil && tc.QStringIgnore(*ds.QStringIgnore) == tc.QStringIgnoreUseInCacheKeyAndPassUp {
			return false, warnings
		}

		if qStringHandling != nil {
			return !(*qStringHandling), warnings
		}

		return true, warnings
	}

	if param := serverParams[ParentConfigParamQStringHandling]; param != "" {
		if useQStr := ParentSelectParamQStringHandlingToBool(param); useQStr != nil {
			return !(*useQStr), warnings
		} else if param != "" {
			warnings = append(warnings, "Server param '"+ParentConfigParamQStringHandling+"' value '"+param+"' malformed, not using!")
		}
		// TODO warn if parsing fails?
	}
	if qStringHandling != nil {
		return !(*qStringHandling), warnings
	}
	if ds.QStringIgnore != nil && tc.QStringIgnore(*ds.QStringIgnore) == tc.QStringIgnoreUseInCacheKeyAndPassUp {
		return false, warnings
	}
	return true, warnings
}

// serverParentageParams gets the Parameters used for parent= line, or defaults if they don't exist
// Returns the Parameters used for parent= lines for the given server, and any warnings.
func serverParentageParams(sv *Server, params []parameterWithProfilesMap) (profileCache, []string) {
	warnings := []string{}
	// TODO deduplicate with atstccfg/parentdotconfig.go
	profileCache := defaultProfileCache()
	if sv.TCPPort != nil {
		profileCache.Port = *sv.TCPPort
	}
	for _, param := range params {
		if _, ok := param.ProfileNames[*sv.Profile]; !ok {
			continue
		}
		switch param.Name {
		case ParentConfigCacheParamWeight:
			profileCache.Weight = param.Value
		case ParentConfigCacheParamPort:
			if i, err := strconv.Atoi(param.Value); err != nil {
				warnings = append(warnings, "port param is not an integer, skipping! : "+err.Error())
			} else {
				profileCache.Port = i
			}
		case ParentConfigCacheParamUseIP:
			profileCache.UseIP = param.Value == "1"
		case ParentConfigCacheParamRank:

			if i, err := strconv.Atoi(param.Value); err != nil {
				warnings = append(warnings, "rank param is not an integer, skipping! : "+err.Error())
			} else {
				profileCache.Rank = i
			}
		case ParentConfigCacheParamNotAParent:
			profileCache.NotAParent = param.Value != "false"
		}
	}

	return profileCache, warnings
}

func serverParentStr(sv *Server, svParams profileCache) (*ParentAbstractionServiceParent, error) {
	if svParams.NotAParent {
		return nil, nil
	}
	host := ""
	if svParams.UseIP {
		// TODO get service interface here
		ip := getServerIPAddress(sv)
		if ip == nil {
			return nil, errors.New("server params Use IP, but has no valid IPv4 Service Address")
		}
		host = ip.String()
	} else {
		host = *sv.HostName + "." + *sv.DomainName
	}

	weight, err := strconv.ParseFloat(svParams.Weight, 64)
	if err != nil {
		// TODO warn? error?
		weight = DefaultParentWeight
	}

	return &ParentAbstractionServiceParent{
		FQDN:   host,
		Port:   svParams.Port,
		Weight: weight,
	}, nil
}

// GetTopologyParents returns the parents, secondary parents, any warnings, and any error.
func getTopologyParents(
	server *Server,
	ds *DeliveryService,
	servers []Server,
	parentConfigParams []parameterWithProfilesMap, // all params with configFile parent.config
	topology tc.Topology,
	serverIsLastTier bool,
	serverCapabilities map[int]map[ServerCapability]struct{},
	dsRequiredCapabilities map[int]map[ServerCapability]struct{},
	dsOrigins map[ServerID]struct{}, // for Topology DSes, MSO still needs DeliveryServiceServer assignments.
	dsMergeGroups []string, // sorted parent merge groups for this ds
) ([]*ParentAbstractionServiceParent, []*ParentAbstractionServiceParent, []string, error) {
	warnings := []string{}
	// If it's the last tier, then the parent is the origin.
	// Note this doesn't include MSO, whose final tier cachegroup points to the origin cachegroup.
	if serverIsLastTier {
		orgURI, orgWarns, err := getOriginURI(*ds.OrgServerFQDN) // TODO pass, instead of calling again
		warnings = append(warnings, orgWarns...)
		if err != nil {
			return nil, nil, warnings, err
		}

		orgPort, err := strconv.Atoi(orgURI.Port())
		if err != nil {
			warnings = append(warnings, "DS "+*ds.XMLID+" origin '"+*ds.OrgServerFQDN+"' failed to parse port, using 80!")
			orgPort = 80
		}
		parent := &ParentAbstractionServiceParent{
			FQDN:   orgURI.Hostname(),
			Port:   orgPort,
			Weight: DefaultParentWeight,
		}

		return []*ParentAbstractionServiceParent{parent}, nil, warnings, nil
	}

	svNode := tc.TopologyNode{}
	for _, node := range topology.Nodes {
		if node.Cachegroup == *server.Cachegroup {
			svNode = node
			break
		}
	}
	if svNode.Cachegroup == "" {
		return nil, nil, warnings, errors.New("This server '" + *server.HostName + "' not in DS " + *ds.XMLID + " topology, skipping")
	}

	if len(svNode.Parents) == 0 {
		return nil, nil, warnings, errors.New("DS " + *ds.XMLID + " topology '" + *ds.Topology + "' is last tier, but NonLastTier called! Should never happen")
	}
	if numParents := len(svNode.Parents); numParents > 2 {
		warnings = append(warnings, "DS "+*ds.XMLID+" topology '"+*ds.Topology+"' has "+strconv.Itoa(numParents)+" parent nodes, but Apache Traffic Server only supports Primary and Secondary (2) lists of parents. CacheGroup nodes after the first 2 will be ignored!")
	}
	if len(topology.Nodes) <= svNode.Parents[0] {
		return nil, nil, warnings, errors.New("DS " + *ds.XMLID + " topology '" + *ds.Topology + "' node parent " + strconv.Itoa(svNode.Parents[0]) + " greater than number of topology nodes " + strconv.Itoa(len(topology.Nodes)) + ". Cannot create parents!")
	}
	if len(svNode.Parents) > 1 && len(topology.Nodes) <= svNode.Parents[1] {
		warnings = append(warnings, "DS "+*ds.XMLID+" topology '"+*ds.Topology+"' node secondary parent "+strconv.Itoa(svNode.Parents[1])+" greater than number of topology nodes "+strconv.Itoa(len(topology.Nodes))+". Secondary parent will be ignored!")
	}

	parentCG := topology.Nodes[svNode.Parents[0]].Cachegroup
	secondaryParentCG := ""
	if len(svNode.Parents) > 1 && len(topology.Nodes) > svNode.Parents[1] {
		secondaryParentCG = topology.Nodes[svNode.Parents[1]].Cachegroup
	}

	if parentCG == "" {
		return nil, nil, warnings, errors.New("Server '" + *server.HostName + "' DS " + *ds.XMLID + " topology '" + *ds.Topology + "' cachegroup '" + *server.Cachegroup + "' topology node parent " + strconv.Itoa(svNode.Parents[0]) + " is not in the topology!")
	}

	parentStrs := []*ParentAbstractionServiceParent{}
	secondaryParentStrs := []*ParentAbstractionServiceParent{}

	serversWithParams := []serverWithParams{}
	for _, sv := range servers {
		serverParentParams, parentWarns := serverParentageParams(&sv, parentConfigParams)
		warnings = append(warnings, parentWarns...)
		serversWithParams = append(serversWithParams, serverWithParams{
			Server: sv,
			Params: serverParentParams,
		})
	}
	sort.Sort(serversWithParamsSortByRank(serversWithParams))

	for _, sv := range serversWithParams {
		if sv.ID == nil {
			warnings = append(warnings, "TO Servers server had nil ID, skipping")
			continue
		} else if sv.Cachegroup == nil {
			warnings = append(warnings, "TO Servers server had nil Cachegroup, skipping")
			continue
		} else if sv.CDNName == nil {
			warnings = append(warnings, "TO servers had server with missing CDNName, skipping!")
			continue
		} else if sv.Status == nil || *sv.Status == "" {
			warnings = append(warnings, "TO servers had server with missing Status, skipping!")
			continue
		}

		if !strings.HasPrefix(sv.Type, tc.EdgeTypePrefix) && !strings.HasPrefix(sv.Type, tc.MidTypePrefix) && sv.Type != tc.OriginTypeName {
			continue // only consider edges, mids, and origins in the CacheGroup.
		}
		if _, dsHasOrigin := dsOrigins[ServerID(*sv.ID)]; sv.Type == tc.OriginTypeName && !dsHasOrigin {
			continue
		}
		if *sv.CDNName != *server.CDNName {
			continue
		}
		if *sv.Status != string(tc.CacheStatusReported) && *sv.Status != string(tc.CacheStatusOnline) {
			continue
		}

		if sv.Type != tc.OriginTypeName && !hasRequiredCapabilities(serverCapabilities[*sv.ID], dsRequiredCapabilities[*ds.ID]) {
			continue
		}
		if *sv.Cachegroup == parentCG {
			parentStr, err := serverParentStr(&sv.Server, sv.Params)
			if err != nil {
				return nil, nil, warnings, errors.New("getting server parent string: " + err.Error())
			}
			if parentStr != nil { // will be nil if server is not_a_parent (possibly other reasons)
				parentStrs = append(parentStrs, parentStr)
			}
		}
		if *sv.Cachegroup == secondaryParentCG {
			parentStr, err := serverParentStr(&sv.Server, sv.Params)
			if err != nil {
				return nil, nil, warnings, errors.New("getting server parent string: " + err.Error())
			}
			secondaryParentStrs = append(secondaryParentStrs, parentStr)
		}
	}

	if 0 < len(dsMergeGroups) && 0 < len(secondaryParentStrs) {
		if sort.SearchStrings(dsMergeGroups, parentCG) < len(dsMergeGroups) &&
			sort.SearchStrings(dsMergeGroups, secondaryParentCG) < len(dsMergeGroups) {
			parentStrs = append(parentStrs, secondaryParentStrs...)
			secondaryParentStrs = nil
		}
	}

	return parentStrs, secondaryParentStrs, warnings, nil
}

// getOriginURI returns the URL, any warnings, and any error.
func getOriginURI(fqdn string) (*url.URL, []string, error) {
	warnings := []string{}

	orgURI, err := url.Parse(fqdn) // TODO verify origin is always a host:port
	if err != nil {
		return nil, warnings, errors.New("parsing: " + err.Error())
	}
	if orgURI.Port() == "" {
		if orgURI.Scheme == "http" {
			orgURI.Host += ":80"
		} else if orgURI.Scheme == "https" {
			orgURI.Host += ":443"
		} else {
			warnings = append(warnings, "non-top-level: origin '"+fqdn+"' is unknown scheme '"+orgURI.Scheme+"', but has no port! Using as-is! ")
		}
	}
	return orgURI, warnings, nil
}

// getParentStrs returns the primary parents, secondary parents, the secondary mode, and any warnings.
func getParentStrs(
	ds *DeliveryService,
	dsRequiredCapabilities map[int]map[ServerCapability]struct{},
	parentInfos []parentInfo,
	atsMajorVer int,
	tryAllPrimariesBeforeSecondary bool,
) ([]*ParentAbstractionServiceParent, []*ParentAbstractionServiceParent, ParentAbstractionServiceParentSecondaryMode, []string) {
	warnings := []string{}
	parentInfo := []*ParentAbstractionServiceParent{}
	secondaryParentInfo := []*ParentAbstractionServiceParent{}

	sort.Sort(parentInfoSortByRank(parentInfos))

	for _, parent := range parentInfos { // TODO fix magic key
		if !hasRequiredCapabilities(parent.Capabilities, dsRequiredCapabilities[*ds.ID]) {
			continue
		}

		pTxt := parent.ToAbstract()
		if parent.PrimaryParent {
			parentInfo = append(parentInfo, pTxt)
		} else if parent.SecondaryParent {
			secondaryParentInfo = append(secondaryParentInfo, pTxt)
		}
	}

	if len(parentInfo) == 0 {
		parentInfo = secondaryParentInfo
		secondaryParentInfo = []*ParentAbstractionServiceParent{}
	}

	// TODO remove duplicate code with top level if block
	seen := map[string]struct{}{} // TODO change to host+port? host isn't unique
	parentInfo, seen = RemoveParentDuplicates(parentInfo, seen)
	secondaryParentInfo, seen = RemoveParentDuplicates(secondaryParentInfo, seen)

	dsName := tc.DeliveryServiceName("")
	if ds != nil && ds.XMLID != nil {
		dsName = tc.DeliveryServiceName(*ds.XMLID)
	}

	// parents := ""
	// secondaryParents := "" // "secparents" in Perl

	// TODO the abstract->text needs to take this into account
	// if atsMajorVer >= 6 && len(secondaryParentInfo) > 0 {
	// parents = `parent="` + strings.Join(parentInfo, "") + `"`
	// secondaryParents = ` secondary_parent="` + strings.Join(secondaryParentInfo, "") + `"`
	secondaryMode, secondaryModeWarnings := getSecondaryModeStr(tryAllPrimariesBeforeSecondary, atsMajorVer, dsName)
	warnings = append(warnings, secondaryModeWarnings...)
	// 	secondaryParents += secondaryModeStr
	// } else {
	// 	parents = `parent="` + strings.Join(parentInfo, "") + strings.Join(secondaryParentInfo, "") + `"`
	// }

	return parentInfo, secondaryParentInfo, secondaryMode, warnings
}

// getMSOParentStrs returns the parents= and secondary_parents= strings for ATS parent.config lines for MSO, and any warnings.
func getMSOParentStrs(
	ds *DeliveryService,
	parentInfos []parentInfo,
	atsMajorVer int,
	msoAlgorithm ParentAbstractionServiceRetryPolicy,
	tryAllPrimariesBeforeSecondary bool,
) ([]*ParentAbstractionServiceParent, []*ParentAbstractionServiceParent, ParentAbstractionServiceParentSecondaryMode, []string) {
	warnings := []string{}
	// TODO determine why MSO is different, and if possible, combine with getParentAndSecondaryParentStrs.

	rankedParents := parentInfoSortByRank(parentInfos)
	sort.Sort(rankedParents)

	parentInfoTxt := []*ParentAbstractionServiceParent{}
	secondaryParentInfo := []*ParentAbstractionServiceParent{}
	nullParentInfo := []*ParentAbstractionServiceParent{}
	for _, parent := range ([]parentInfo)(rankedParents) {
		if parent.PrimaryParent {
			parentInfoTxt = append(parentInfoTxt, parent.ToAbstract())
		} else if parent.SecondaryParent {
			secondaryParentInfo = append(secondaryParentInfo, parent.ToAbstract())
		} else {
			nullParentInfo = append(nullParentInfo, parent.ToAbstract())
		}
	}

	if len(parentInfoTxt) == 0 {
		// If no parents are found in the secondary parent either, then set the null parent list (parents in neither secondary or primary)
		// as the secondary parent list and clear the null parent list.
		if len(secondaryParentInfo) == 0 {
			secondaryParentInfo = nullParentInfo
			nullParentInfo = []*ParentAbstractionServiceParent{}
		}
		parentInfoTxt = secondaryParentInfo
		secondaryParentInfo = []*ParentAbstractionServiceParent{} // TODO should this be '= secondary'? Currently emulates Perl
	}

	// TODO benchmark, verify this isn't slow. if it is, it could easily be made faster
	seen := map[string]struct{}{} // TODO change to host+port? host isn't unique
	parentInfoTxt, seen = RemoveParentDuplicates(parentInfoTxt, seen)
	secondaryParentInfo, seen = RemoveParentDuplicates(secondaryParentInfo, seen)
	nullParentInfo, seen = RemoveParentDuplicates(nullParentInfo, seen)

	// secondaryParentStr := strings.Join(secondaryParentInfo, "") + strings.Join(nullParentInfo, "")
	secondaryParentInfo = append(secondaryParentInfo, nullParentInfo...)

	dsName := tc.DeliveryServiceName("")
	if ds != nil && ds.XMLID != nil {
		dsName = tc.DeliveryServiceName(*ds.XMLID)
	}

	// If the ats version supports it and the algorithm is consistent hash, put secondary and non-primary parents into secondary parent group.
	// This will ensure that secondary and tertiary parents will be unused unless all hosts in the primary group are unavailable.

	// parents := ""
	// secondaryParents := ""

	// TODO add this logic to the abstraction->text converter
	// if atsMajorVer >= 6 && msoAlgorithm == "consistent_hash" && len(secondaryParentStr) > 0 {
	// parents = `parent="` + strings.Join(parentInfoTxt, "") + `"`
	// secondaryParents = ` secondary_parent="` + secondaryParentStr + `"`
	secondaryMode, secondaryModeWarnings := getSecondaryModeStr(tryAllPrimariesBeforeSecondary, atsMajorVer, dsName)
	warnings = append(warnings, secondaryModeWarnings...)
	// 	secondaryParents += secondaryModeStr
	// } else {
	// 	parents = `parent="` + strings.Join(parentInfoTxt, "") + secondaryParentStr + `"`
	// }
	return parentInfoTxt, secondaryParentInfo, secondaryMode, warnings
}

// makeParentInfo returns the parent info and any warnings
func makeParentInfo(
	serverParentCGData serverParentCacheGroupData,
	serverDomain string, // getCDNDomainByProfileID(tx, server.ProfileID)
	profileCaches map[ProfileID]profileCache, // getServerParentCacheGroupProfiles(tx, server)
	originServers map[OriginHost][]cgServer, // getServerParentCacheGroupProfiles(tx, server)
) (map[OriginHost][]parentInfo, []string) {
	warnings := []string{}
	parentInfos := map[OriginHost][]parentInfo{}

	// note servers also contains an "all" key
	for originHost, servers := range originServers {
		for _, row := range servers {
			profile := profileCaches[row.ProfileID]
			if profile.NotAParent {
				continue
			}
			// Perl has this check, but we only select servers ("deliveryServices" in Perl) with the right CDN in the first place
			// if profile.Domain != serverDomain {
			// 	continue
			// }

			weight, err := strconv.ParseFloat(profile.Weight, 64)
			if err != nil {
				warnings = append(warnings, "profile "+strconv.Itoa(int(row.ProfileID))+" had malformed weight, using default!")
				weight = DefaultParentWeight
			}

			parentInf := parentInfo{
				Host:            row.ServerHost,
				Port:            profile.Port,
				Domain:          row.Domain,
				Weight:          weight,
				UseIP:           profile.UseIP,
				Rank:            profile.Rank,
				IP:              row.ServerIP,
				PrimaryParent:   serverParentCGData.ParentID == row.CacheGroupID,
				SecondaryParent: serverParentCGData.SecondaryParentID == row.CacheGroupID,
				Capabilities:    row.Capabilities,
			}
			if parentInf.Port < 1 {
				parentInf.Port = row.ServerPort
			}
			parentInfos[originHost] = append(parentInfos[originHost], parentInf)
		}
	}
	return parentInfos, warnings
}

// unavailableServerRetryResponsesValid returns whether a unavailable_server_retry_responses parameter is valid for an ATS parent rule.
func unavailableServerRetryResponsesValid(s string) bool {
	// optimization if param is empty
	if s == "" {
		return false
	}
	re := regexp.MustCompile(`^"(:?\d{3},)+\d{3}"\s*$`) // TODO benchmark, cache if performance matters
	return re.MatchString(s)
}

// getOriginServersAndProfileCaches returns the origin servers, ProfileCaches, any warnings, and any error.
func getOriginServersAndProfileCaches(
	cgServers map[int]Server,
	parentServerDSes map[int]map[int]struct{},
	profileParentConfigParams map[string]map[string]string, // map[profileName][paramName]paramVal
	dses []DeliveryService,
	serverCapabilities map[int]map[ServerCapability]struct{},
) (map[OriginHost][]cgServer, map[ProfileID]profileCache, []string, error) {
	warnings := []string{}
	originServers := map[OriginHost][]cgServer{}  // "deliveryServices" in Perl
	profileCaches := map[ProfileID]profileCache{} // map[profileID]ProfileCache

	dsIDMap := map[int]DeliveryService{}
	for _, ds := range dses {
		if ds.ID == nil {
			return nil, nil, warnings, errors.New("delivery services got nil ID!")
		}
		if !ds.Type.IsHTTP() && !ds.Type.IsDNS() {
			continue // skip ANY_MAP, STEERING, etc
		}
		dsIDMap[*ds.ID] = ds
	}

	allDSMap := map[int]DeliveryService{} // all DSes for this server, NOT all dses in TO
	for _, dsIDs := range parentServerDSes {
		for dsID, _ := range dsIDs {
			if _, ok := dsIDMap[dsID]; !ok {
				// this is normal if the TO was too old to understand our /deliveryserviceserver?servers= query param
				// In which case, the DSS will include DSes from other CDNs, which aren't in the dsIDMap
				// If the server was new enough to respect the params, this should never happen.
				// warnings = append(warnings, ("getting delivery services: parent server DS %v not in dsIDMap\n", dsID)
				continue
			}
			if _, ok := allDSMap[dsID]; !ok {
				allDSMap[dsID] = dsIDMap[dsID]
			}
		}
	}

	dsOrigins, dsOriginWarns, err := getDSOrigins(allDSMap)
	warnings = append(warnings, dsOriginWarns...)
	if err != nil {
		return nil, nil, warnings, errors.New("getting DS origins: " + err.Error())
	}

	profileParams, profParamWarns := getParentConfigProfileParams(cgServers, profileParentConfigParams)
	warnings = append(warnings, profParamWarns...)

	for _, cgSv := range cgServers {
		if cgSv.ID == nil {
			warnings = append(warnings, "getting origin servers: got server with nil ID, skipping!")
			continue
		} else if cgSv.HostName == nil {
			warnings = append(warnings, "getting origin servers: got server with nil HostName, skipping!")
			continue
		} else if cgSv.TCPPort == nil {
			warnings = append(warnings, "getting origin servers: got server with nil TCPPort, skipping!")
			continue
		} else if cgSv.CachegroupID == nil {
			warnings = append(warnings, "getting origin servers: got server with nil CachegroupID, skipping!")
			continue
		} else if cgSv.StatusID == nil {
			warnings = append(warnings, "getting origin servers: got server with nil StatusID, skipping!")
			continue
		} else if cgSv.TypeID == nil {
			warnings = append(warnings, "getting origin servers: got server with nil TypeID, skipping!")
			continue
		} else if cgSv.CDNID == nil {
			warnings = append(warnings, "getting origin servers: got server with nil CDNID, skipping!")
			continue
		} else if cgSv.DomainName == nil {
			warnings = append(warnings, "getting origin servers: got server with nil DomainName, skipping!")
			continue
		}

		ipAddr := getServerIPAddress(&cgSv)
		if ipAddr == nil {
			warnings = append(warnings, "getting origin servers: got server with no valid IP Address, skipping!")
			continue
		}

		realCGServer := cgServer{
			ServerID:     ServerID(*cgSv.ID),
			ServerHost:   *cgSv.HostName,
			ServerIP:     ipAddr.String(),
			ServerPort:   *cgSv.TCPPort,
			CacheGroupID: *cgSv.CachegroupID,
			Status:       *cgSv.StatusID,
			Type:         *cgSv.TypeID,
			CDN:          *cgSv.CDNID,
			TypeName:     cgSv.Type,
			Domain:       *cgSv.DomainName,
			Capabilities: serverCapabilities[*cgSv.ID],
		}

		if cgSv.Type == tc.OriginTypeName {
			for dsID, _ := range parentServerDSes[*cgSv.ID] { // map[serverID][]dsID
				orgURI := dsOrigins[dsID]
				if orgURI == nil {
					// warnings = append(warnings, fmt.Sprintf(("ds %v has no origins! Skipping!\n", dsID) // TODO determine if this is normal
					continue
				}
				orgHost := OriginHost(orgURI.Host)
				originServers[orgHost] = append(originServers[orgHost], realCGServer)
			}
		} else {
			originServers[deliveryServicesAllParentsKey] = append(originServers[deliveryServicesAllParentsKey], realCGServer)
		}

		if _, profileCachesHasProfile := profileCaches[realCGServer.ProfileID]; !profileCachesHasProfile {
			if profileCache, profileParamsHasProfile := profileParams[(*cgSv.Profiles)[0]]; !profileParamsHasProfile {
				warnings = append(warnings, fmt.Sprintf("cachegroup has server with profile %+v but that profile has no parameters\n", (*cgSv.Profiles)[0]))
				profileCaches[realCGServer.ProfileID] = defaultProfileCache()
			} else {
				profileCaches[realCGServer.ProfileID] = profileCache
			}
		}
	}

	return originServers, profileCaches, warnings, nil
}

// GetParentConfigProfileParams returns the parent config profile params, and any warnings.
func getParentConfigProfileParams(
	cgServers map[int]Server,
	profileParentConfigParams map[string]map[string]string, // map[profileName][paramName]paramVal
) (map[string]profileCache, []string) {
	warnings := []string{}
	parentConfigServerCacheProfileParams := map[string]profileCache{} // map[profileName]ProfileCache
	for _, cgServer := range cgServers {
		if cgServer.Profiles == nil {
			warnings = append(warnings, "getting parent config profile params: server has nil profiles, skipping!")
			continue
		}
		profileCache, ok := parentConfigServerCacheProfileParams[(*cgServer.Profiles)[0]]
		if !ok {
			profileCache = defaultProfileCache()
		}
		params, ok := profileParentConfigParams[(*cgServer.Profiles)[0]]
		if !ok {
			parentConfigServerCacheProfileParams[(*cgServer.Profiles)[0]] = profileCache
			continue
		}
		for name, val := range params {
			switch name {
			case ParentConfigCacheParamWeight:
				// f, err := strconv.ParseFloat(param.Val, 64)
				// if err != nil {
				// 	warnings = append(warnings, "parent.config generation: weight param is not a float, skipping! : " + err.Error())
				// } else {
				// 	profileCache.Weight = f
				// }
				// TODO validate float?
				profileCache.Weight = val
			case ParentConfigCacheParamPort:
				if i, err := strconv.Atoi(val); err != nil {
					warnings = append(warnings, "port param is not an integer, skipping! : "+err.Error())
				} else {
					profileCache.Port = i
				}
			case ParentConfigCacheParamUseIP:
				profileCache.UseIP = val == "1"
			case ParentConfigCacheParamRank:
				if i, err := strconv.Atoi(val); err != nil {
					warnings = append(warnings, "rank param is not an integer, skipping! : "+err.Error())
				} else {
					profileCache.Rank = i
				}
			case ParentConfigCacheParamNotAParent:
				profileCache.NotAParent = val != "false"
			}
		}
		parentConfigServerCacheProfileParams[(*cgServer.Profiles)[0]] = profileCache
	}
	return parentConfigServerCacheProfileParams, warnings
}

// getDSOrigins takes a map[deliveryServiceID]DeliveryService, and returns a map[DeliveryServiceID]OriginURI, any warnings, and any error.
func getDSOrigins(dses map[int]DeliveryService) (map[int]*originURI, []string, error) {
	warnings := []string{}
	dsOrigins := map[int]*originURI{}
	for _, ds := range dses {
		if ds.ID == nil {
			return nil, warnings, errors.New("ds has nil ID")
		}
		if ds.XMLID == nil {
			return nil, warnings, errors.New("ds has nil XMLID")
		}
		if ds.OrgServerFQDN == nil {
			warnings = append(warnings, fmt.Sprintf("GetDSOrigins ds %v got nil OrgServerFQDN, skipping!\n", *ds.XMLID))
			continue
		}
		orgURL, err := url.Parse(*ds.OrgServerFQDN)
		if err != nil {
			return nil, warnings, errors.New("parsing ds '" + *ds.XMLID + "' OrgServerFQDN '" + *ds.OrgServerFQDN + "': " + err.Error())
		}
		if orgURL.Scheme == "" {
			return nil, warnings, errors.New("parsing ds '" + *ds.XMLID + "' OrgServerFQDN '" + *ds.OrgServerFQDN + "': " + "missing scheme")
		}
		if orgURL.Host == "" {
			return nil, warnings, errors.New("parsing ds '" + *ds.XMLID + "' OrgServerFQDN '" + *ds.OrgServerFQDN + "': " + "missing scheme")
		}

		scheme := orgURL.Scheme
		host := orgURL.Hostname()
		port := orgURL.Port()
		if port == "" {
			if scheme == "http" {
				port = "80"
			} else if scheme == "https" {
				port = "443"
			} else {
				warnings = append(warnings, "parsing ds '"+*ds.XMLID+"' OrgServerFQDN '"+*ds.OrgServerFQDN+"': "+"unknown scheme '"+scheme+"' and no port, leaving port empty!")
			}
		}
		dsOrigins[*ds.ID] = &originURI{Scheme: scheme, Host: host, Port: port}
	}
	return dsOrigins, warnings, nil
}

// makeDSOrigins returns the DS Origins and any warnings.
func makeDSOrigins(dsses []DeliveryServiceServer, dses []DeliveryService, servers []Server) (map[DeliveryServiceID]map[ServerID]struct{}, []string) {
	warnings := []string{}
	dssMap := map[DeliveryServiceID]map[ServerID]struct{}{}
	for _, dss := range dsses {
		dsID := DeliveryServiceID(dss.DeliveryService)
		serverID := ServerID(dss.Server)
		if dssMap[dsID] == nil {
			dssMap[dsID] = map[ServerID]struct{}{}
		}
		dssMap[dsID][serverID] = struct{}{}
	}

	svMap := map[ServerID]Server{}
	for _, sv := range servers {
		if sv.ID == nil {
			warnings = append(warnings, "got server with missing ID, skipping!")
		}
		svMap[ServerID(*sv.ID)] = sv
	}

	dsOrigins := map[DeliveryServiceID]map[ServerID]struct{}{}
	for _, ds := range dses {
		if ds.ID == nil {
			warnings = append(warnings, "got ds with missing ID, skipping!")
			continue
		}
		dsID := DeliveryServiceID(*ds.ID)
		assignedServers := dssMap[dsID]
		for svID, _ := range assignedServers {
			sv := svMap[svID]
			if sv.Type != tc.OriginTypeName {
				continue
			}
			if dsOrigins[dsID] == nil {
				dsOrigins[dsID] = map[ServerID]struct{}{}
			}
			dsOrigins[dsID][svID] = struct{}{}
		}
	}
	return dsOrigins, warnings
}

// getProfileParentConfigParams returns a map[profileName][paramName]paramVal and any warnings
func getProfileParentConfigParams(tcParentConfigParams []tc.Parameter) (map[string]map[string]string, []string) {
	warnings := []string{}
	parentConfigParamsWithProfiles, err := tcParamsToParamsWithProfiles(tcParentConfigParams)
	if err != nil {
		warnings = append(warnings, "error getting profiles from Traffic Ops Parameters, Parameters will not be considered for generation! : "+err.Error())
		parentConfigParamsWithProfiles = []parameterWithProfiles{}
	}
	// parentConfigParams := parameterWithProfilesToMap(parentConfigParamsWithProfiles)

	// this is an optimization, to avoid looping over all params, for every DS. Instead, we loop over all params only once, and put them in a profile map.
	profileParentConfigParams := map[string]map[string]string{} // map[profileName][paramName]paramVal
	for _, param := range parentConfigParamsWithProfiles {
		for _, profile := range param.ProfileNames {
			if _, ok := profileParentConfigParams[profile]; !ok {
				profileParentConfigParams[profile] = map[string]string{}
			}
			profileParentConfigParams[profile][param.Name] = param.Value
		}
	}
	return profileParentConfigParams, warnings
}

// getServerParentConfigParams returns a map[name]value.
// Intended to be called with the result of getProfileParentConfigParams.
func getServerParentConfigParams(server *Server, profileParentConfigParams map[string]map[string]string) map[string]string {
	// We only need parent.config params, don't need all the params on the server
	serverParams := map[string]string{}
	if server.Profiles == nil || (*server.Profiles)[0] != "" { // TODO warn/error if false? Servers requires profiles
		for name, val := range profileParentConfigParams[(*server.Profiles)[0]] {
			if name == ParentConfigParamQStringHandling ||
				name == ParentConfigParamAlgorithm ||
				name == ParentConfigParamQString {
				serverParams[name] = val
			}
		}
	}
	return serverParams
}
