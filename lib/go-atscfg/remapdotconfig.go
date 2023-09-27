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
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"

	"github.com/cbroglie/mustache"
)

// RemapConfigDropQstringConfigFile is the configuration file for the
// `drop_qstring` plugin used in remap.config files.
const RemapConfigDropQstringConfigFile = `drop_qstring.config`

// ContentTypeRemapDotConfig is the MIME type of the contents of a remap.config
// ATS configuration file.
const ContentTypeRemapDotConfig = ContentTypeTextASCII

// LineCommentRemapDotConfig is the string used to signal the beginning of a
// line comment in the grammar of a remap.config ATS configuration file.
const LineCommentRemapDotConfig = LineCommentHash

// RemapConfigRangeDirective is a special string which, if found in a Delivery
// Service's Raw Remap Text, will be replaced by t3c during configuration
// generation with the parts of a remap.config file line that are necessary to
// implement said Delivery Service's declared Range Request Handling.
const RemapConfigRangeDirective = `__RANGE_DIRECTIVE__`

// RemapConfigCachekeyDirective is a special string which, if found in a
// Delivery Service's Raw Remap Text, will be replaced by t3c during
// configuration generation with the parts of a remap.config file line that call
// the 'cachekey.so' ATS plugin.
const RemapConfigCachekeyDirective = `__CACHEKEY_DIRECTIVE__`

// RemapConfigRegexRemapDirective is a special string which, if found in a
// Delivery Service's Raw Remap Text, will be replaced by t3c during
// configuration generation with the parts of a remap.config file line that call
// the 'regex_remap.so' ATS plugin.
const RemapConfigRegexRemapDirective = `__REGEX_REMAP_DIRECTIVE__`

const RemapConfigTemplateFirst = `template.first`
const RemapConfigTemplateInner = `template.inner`
const RemapConfigTemplateLast = `template.last`

const DefaultFirstRemapConfigTemplateString = `map {{{Source}}} {{{Destination}}} {{{Strategy}}} {{{Dscp}}} {{{HeaderRewrite}}} {{{DropQstring}}} {{{Signing}}} {{{RegexRemap}}} {{{Cachekey}}} {{{RangeRequests}}} {{{Pacing}}} {{{RawText}}}`
const DefaultLastRemapConfigTemplateString = `map {{{Source}}} {{{Destination}}} {{{Strategy}}} {{{HeaderRewrite}}} {{{Cachekey}}} {{{RangeRequests}}} {{{RawText}}}`
const DefaultInnerRemapConfigTemplateString = DefaultLastRemapConfigTemplateString
const selfHealParam = `no_self_healing`

type LineTemplates map[string]*mustache.Template

var RemapLineTemplates = LineTemplates{}

// This parses but also maintains a cache of parsed templates
func (lts *LineTemplates) parse(templateString string) (*mustache.Template, error) {
	if tmpl, ok := RemapLineTemplates[templateString]; ok {
		return tmpl, nil
	}
	tmpl, err := mustache.ParseString(templateString)
	if err != nil {
		RemapLineTemplates[templateString] = tmpl
		tmpl = nil
	}
	return tmpl, err
}

// RemapTags contains the mustache template fillers.
type RemapTags struct {
	Source        string
	Destination   string
	Strategy      string
	Dscp          string
	HeaderRewrite string
	DropQstring   string
	Signing       string
	RegexRemap    string
	Cachekey      string
	Pacing        string
	RangeRequests string
	RawText       string
}

func (rs *RemapTags) AnyMidOptionsSet() bool {
	return rs.Strategy != "" ||
		rs.HeaderRewrite != "" ||
		rs.Cachekey != "" ||
		rs.RangeRequests != ""
}

// RemapDotConfigOpts contains settings to configure generation options.
type RemapDotConfigOpts struct {
	// VerboseComments is whether to add informative comments to the generated file, about what was generated and why.
	// Note this does not include the header comment, which is configured separately with HdrComment.
	// These comments are human-readable and not guaranteed to be consistent between versions. Automating anything based on them is strongly discouraged.
	VerboseComments bool

	// HdrComment is the header comment to include at the beginning of the file.
	// This should be the text desired, without comment syntax (like # or //). The file's comment syntax will be added.
	// To omit the header comment, pass the empty string.
	HdrComment string

	// UseStrategies is whether to use strategies.yaml rather than parent.config.
	UseStrategies bool
	// UseCoreStrategies is whether to use the ATS core strategies, rather than the parent_select plugin.
	// This has no effect if UseStrategies is false.
	UseStrategiesCore bool

	// ATSMajorVersion is the integral major version of Apache Traffic server,
	// used to generate the proper config for the proper version.
	//
	// If omitted or 0, the major version will be read from the Server's Profile Parameter config file 'package' name 'trafficserver'. If no such Parameter exists, the ATS version will default to 5.
	// This was the old Traffic Control behavior, before the version was specifiable externally.
	//
	ATSMajorVersion uint
}

// MakeRemapDotConfig constructs a remap.config ATS configuration file.
func MakeRemapDotConfig(
	server *Server,
	servers []Server,
	unfilteredDSes []DeliveryService,
	dss []DeliveryServiceServer,
	dsRegexArr []tc.DeliveryServiceRegexes,
	serverParams []tc.ParameterV5,
	cdn *tc.CDNV5,
	remapConfigParams []tc.ParameterV5, // includes cachekey.config
	topologies []tc.TopologyV5,
	cacheGroupArr []tc.CacheGroupNullableV5,
	serverCapabilities map[int]map[ServerCapability]struct{},
	dsRequiredCapabilities map[int]map[ServerCapability]struct{},
	configDir string,
	opt *RemapDotConfigOpts,
) (Cfg, error) {
	if opt == nil {
		opt = &RemapDotConfigOpts{}
	}
	warnings := []string{}

	if !opt.UseStrategies && opt.UseStrategiesCore {
		warnings = append(warnings, "opt.UseStrategies was false, but opt.UseStrategiesCore was set, which has no effect! Not using strategies, per opt.UseStrategies.")
	}

	if server.HostName == "" {
		return Cfg{}, makeErr(warnings, "server HostName missing")
	} else if server.ID == 0 {
		return Cfg{}, makeErr(warnings, "server ID missing")
	} else if server.CacheGroup == "" {
		return Cfg{}, makeErr(warnings, "server Cachegroup missing")
	} else if server.DomainName == "" {
		return Cfg{}, makeErr(warnings, "server DomainName missing")
	}

	cdnDomain := cdn.DomainName
	dsRegexes := MakeDSRegexMap(dsRegexArr)
	// Returned DSes are guaranteed to have a non-nil XMLID, Type, DSCP, ID, and Active.
	dses, dsWarns := remapFilterDSes(server, dss, unfilteredDSes)
	warnings = append(warnings, dsWarns...)

	dsProfilesConfigParams, paramWarns, err := makeDSProfilesConfigParams(server, dses, remapConfigParams)
	warnings = append(warnings, paramWarns...)
	if err != nil {
		warnings = append(warnings, "making Delivery Service Cache Key Params, cache key will be missing! : "+err.Error())
	}

	atsMajorVersion := getATSMajorVersion(opt.ATSMajorVersion, serverParams, &warnings)

	serverPackageParamData, paramWarns := makeServerPackageParamData(server, serverParams)
	warnings = append(warnings, paramWarns...)
	cacheGroups, err := makeCGMap(cacheGroupArr)
	if err != nil {
		return Cfg{}, makeErr(warnings, "making remap.config, config will be malformed! : "+err.Error())
	}

	nameTopologies := makeTopologyNameMap(topologies)
	anyCastPartners := GetAnyCastPartners(server, servers)

	hdr := makeHdrComment(opt.HdrComment)
	txt := ""
	typeWarns := []string{}
	if tc.CacheTypeFromString(server.Type) == tc.CacheTypeMid {
		txt, typeWarns, err = getServerConfigRemapDotConfigForMid(atsMajorVersion, dsProfilesConfigParams, dses, dsRegexes, hdr, server, nameTopologies, cacheGroups, serverCapabilities, dsRequiredCapabilities, configDir, opt)
	} else {
		txt, typeWarns, err = getServerConfigRemapDotConfigForEdge(dsProfilesConfigParams, serverPackageParamData, dses, dsRegexes, atsMajorVersion, hdr, server, anyCastPartners, nameTopologies, cacheGroups, serverCapabilities, dsRequiredCapabilities, cdnDomain, configDir, opt)
	}
	warnings = append(warnings, typeWarns...)
	if err != nil {
		return Cfg{}, makeErr(warnings, err.Error()) // the GetFor funcs include error context
	}

	return Cfg{
		Text:        txt,
		ContentType: ContentTypeRemapDotConfig,
		LineComment: LineCommentRemapDotConfig,
		Warnings:    warnings,
	}, nil
}

// This sticks the DS parameters in a map.
// remap.config parameters use "<plugin>.pparam" key
// cachekey.config parameters retain the 'cachekey.config' key.
func classifyConfigParams(configParams []tc.ParameterV5) (map[string][]tc.ParameterV5, bool) {
	configParamMap := map[string][]tc.ParameterV5{}
	selfHeal := true
	for _, param := range configParams {
		key := param.ConfigFile
		if "remap.config" == key {
			key = param.Name
			if param.Value == selfHealParam {
				selfHeal = false
				continue
			}
		}
		configParamMap[key] = append(configParamMap[key], param)
	}
	return configParamMap, selfHeal
}

// For general <plugin>.pparam parameters.
func paramsStringFor(parameters []tc.ParameterV5, warnings *[]string) (paramsString string) {
	uniquemap := map[string]int{}

	for _, param := range parameters {
		if strings.TrimSpace(param.Value) == "" {
			continue
		}
		paramsString += " @pparam=" + param.Value

		// Try to extract argument
		index := strings.IndexAny(param.Value, "= ")
		arg := ""
		if 0 < index {
			arg = param.Value[:index]
		} else {
			arg = param.Value
		}

		// Warn on detection, but don't correct
		if _, exists := uniquemap[arg]; !exists {
			uniquemap[arg] = 1
		} else {
			*warnings = append(*warnings, "Multiple repeated arguments '"+arg+"'")
		}
	}
	return
}

// for parameters that use 'cachekey.config' as their key.
func paramsStringOldFor(parameters []tc.ParameterV5, warnings *[]string) (paramsString string) {
	// check for duplicate parameters
	uniquemap := map[string]int{}
	paramKeyVals := []keyVal{}
	for _, param := range parameters {
		key := param.Name
		val := param.Value

		if _, exists := uniquemap[key]; !exists {
			uniquemap[key] = 1
			paramKeyVals = append(paramKeyVals, keyVal{Key: key, Val: val})
		} else {
			uniquemap[key]++
			*warnings = append(*warnings, "got multiple parameters for name '"+key+"' - ignoring '"+val+"'")
		}
	}

	sort.Sort(keyVals(paramKeyVals))
	for _, keyVal := range paramKeyVals {
		paramsString += " @pparam=--" + keyVal.Key + "=" + keyVal.Val
	}
	return
}

// Handles special case for cachekey.
func cachekeyArgsFor(configParamsMap map[string][]tc.ParameterV5, warnings *[]string) (argsString string) {

	hasCachekey := false

	if params, ok := configParamsMap["cachekey.pparam"]; ok {
		argsString += paramsStringFor(params, warnings)
		hasCachekey = true
	}

	// Add on the cachekey.config
	if params, ok := configParamsMap["cachekey.config"]; ok {
		if hasCachekey {
			*warnings = append(*warnings, "Both new cachekey.pparam and old cachekey.config parameters assigned")
		}
		argsString += paramsStringOldFor(params, warnings)
	}
	return
}

// lastPrePostRemapLinesFor Returns any pre or post raw remap lines.
func lastPrePostRemapLinesFor(dsConfigParamsMap map[string][]tc.ParameterV5, dsid string) ([]string, []string) {
	preRemapLines := []string{}
	postRemapLines := []string{}

	// Any raw pre pend
	if params, ok := dsConfigParamsMap["LastRawRemapPre"]; ok {
		for _, param := range params {
			preRemapLines = append(preRemapLines, param.Value+" # Raw: "+dsid+"\n")
		}
	}

	// Any raw post pend
	if params, ok := dsConfigParamsMap["LastRawRemapPost"]; ok {
		for _, param := range params {
			postRemapLines = append(postRemapLines, param.Value+" # Raw: "+dsid+"\n")
		}
	}

	return preRemapLines, postRemapLines
}

// getServerConfigRemapDotConfigForMid returns the remap lines, any warnings, and any error.
func getServerConfigRemapDotConfigForMid(
	atsMajorVersion uint,
	profilesConfigParams map[int][]tc.ParameterV5,
	dses []DeliveryService,
	dsRegexes map[tc.DeliveryServiceName][]tc.DeliveryServiceRegex,
	header string,
	server *Server,
	nameTopologies map[TopologyName]tc.TopologyV5,
	cacheGroups map[tc.CacheGroupName]tc.CacheGroupNullableV5,
	serverCapabilities map[int]map[ServerCapability]struct{},
	dsRequiredCapabilities map[int]map[ServerCapability]struct{},
	configDir string,
	opts *RemapDotConfigOpts,
) (string, []string, error) {
	warnings := []string{}
	midRemaps := map[string]string{}
	preRemapLines := []string{}
	postRemapLines := []string{}

	for _, ds := range dses {
		if !hasRequiredCapabilities(serverCapabilities[server.ID], dsRequiredCapabilities[*ds.ID]) {
			continue
		}

		topology, hasTopology := nameTopologies[TopologyName(*ds.Topology)]
		if *ds.Topology != "" && hasTopology {
			topoIncludesServer, err := topologyIncludesServerNullable(topology, server)
			if err != nil {
				return "", warnings, errors.New("getting Topology Server inclusion: " + err.Error())
			}
			if !topoIncludesServer {
				continue
			}
		}

		if ds.Type != nil && !tc.DSType(*ds.Type).UsesMidCache() && (!hasTopology || *ds.Topology == "") {
			continue // Live local delivery services skip mids (except Topologies ignore DS types)
		}

		if ds.OrgServerFQDN == nil || *ds.OrgServerFQDN == "" {
			warnings = append(warnings, "ds '"+ds.XMLID+"' has no origin fqdn, skipping!") // TODO confirm - Perl uses without checking!
			continue
		}

		remapFrom := strings.Replace(*ds.OrgServerFQDN, `https://`, `http://`, -1)

		if midRemaps[remapFrom] != "" {
			continue // skip remap rules from extra HOST_REGEXP entries
		}

		//remapTags := NewRemapTags()
		remapTags := RemapTags{}

		// template fill for moustache
		remapTags.Source = remapFrom
		strategy := strategyDirective(getStrategyName(ds.XMLID), configDir, opts)
		if strategy != "" {
			remapTags.Strategy = strategy
		}

		if *ds.Topology != "" {
			topoTxt, err := makeDSTopologyHeaderRewriteTxt(ds, tc.CacheGroupName(server.CacheGroup), topology, cacheGroups)
			if err != nil {
				return "", warnings, err
			}
			remapTags.HeaderRewrite = topoTxt
		} else if (ds.MidHeaderRewrite != nil && *ds.MidHeaderRewrite != "") || (ds.MaxOriginConnections != nil && *ds.MaxOriginConnections > 0) || (ds.ServiceCategory != nil && *ds.ServiceCategory != "") {
			remapTags.HeaderRewrite = `@plugin=header_rewrite.so @pparam=` + midHeaderRewriteConfigFileName(ds.XMLID)
		}

		// Logic for handling cachekey params
		cachekeyArgs := ``

		// qstring ignore
		if ds.QStringIgnore != nil && *ds.QStringIgnore == tc.QueryStringIgnoreIgnoreInCacheKeyAndPassUp {
			cachekeyArgs = getQStringIgnoreRemap(atsMajorVersion)
		}

		selfHeal := true
		dsConfigParamsMap := map[string][]tc.ParameterV5{}
		if nil != ds.ProfileID {
			dsConfigParamsMap, selfHeal = classifyConfigParams(profilesConfigParams[*ds.ProfileID])
		}

		if len(dsConfigParamsMap) > 0 {
			cachekeyArgs += cachekeyArgsFor(dsConfigParamsMap, &warnings)
		}

		if cachekeyArgs != "" {
			remapTags.Cachekey = `@plugin=cachekey.so` + cachekeyArgs
		}

		if ds.RangeRequestHandling != nil && (*ds.RangeRequestHandling == tc.RangeRequestHandlingCacheRangeRequest || *ds.RangeRequestHandling == tc.RangeRequestHandlingSlice) {
			crrParam := paramsStringFor(dsConfigParamsMap["cache_range_requests.pparam"], &warnings)
			remapTags.RangeRequests = `@plugin=cache_range_requests.so` + crrParam
			if *ds.RangeRequestHandling == tc.RangeRequestHandlingSlice && !strings.Contains(crrParam, "--consider-ims") && selfHeal {
				remapTags.RangeRequests += ` @pparam=--consider-ims`
			}
		}

		isLastCache, err := serverIsLastCacheForDS(server, &ds, nameTopologies, cacheGroups)
		if err != nil {
			return "", warnings, errors.New("determining if cache is the last tier: " + err.Error())
		}

		mapTo := *ds.OrgServerFQDN

		// if this remap is going to a parent, use http not https.
		// cache-to-cache communication inside the CDN is always http (though that's likely to change in the future)
		if !isLastCache {
			mapTo = strings.Replace(mapTo, `https://`, `http://`, -1)
		}

		remapTags.Destination = mapTo

		// look for override template
		defaultParamName := RemapConfigTemplateLast
		if !isLastCache {
			defaultParamName = RemapConfigTemplateInner
		}
		tmplparams, tmplok := dsConfigParamsMap[defaultParamName]

		// check for fallthrough condition or override template
		if remapTags.AnyMidOptionsSet() || tmplok {

			// Check for ds parameter template override
			var template *mustache.Template
			if tmplok {
				if 0 < len(tmplparams) {
					templateString := tmplparams[0].Value
					template, err = RemapLineTemplates.parse(templateString)
					if err != nil {
						warnings = append(warnings, "Error decoding override "+defaultParamName+" "+templateString+" for DS "+ds.XMLID+", falling back!")
					}
				}
			}

			// fall back to default remap config line template
			if template == nil {
				templateString := DefaultLastRemapConfigTemplateString
				if !isLastCache {
					templateString = DefaultInnerRemapConfigTemplateString
				}

				template, err = RemapLineTemplates.parse(templateString)
				if err != nil {
					return "", warnings, errors.New("unable to load default template string: " + templateString + " " + err.Error())
				}
			}

			// Expand the moustache
			var linebuf bytes.Buffer
			err := template.FRender(&linebuf, remapTags)
			if err != nil {
				return "", warnings, errors.New("Error filling mustache template: " + err.Error())
			}
			//midRemaps[remapFrom] = removeDeleteMes(linebuf.String())
			midRemaps[remapFrom] = linebuf.String()
		}

		// Any raw pre or post pend
		dsPreRemaps, dsPostRemaps := lastPrePostRemapLinesFor(dsConfigParamsMap, ds.XMLID)

		// Add to pre/post remap lines if this is last tier
		if len(dsPreRemaps) > 0 || len(dsPostRemaps) > 0 {
			if isLastCache {
				preRemapLines = append(preRemapLines, dsPreRemaps...)
				postRemapLines = append(postRemapLines, dsPostRemaps...)
			}
		}
	}

	textLines := []string{}

	for _, remapLine := range midRemaps {
		textLines = append(textLines, remapLine+"\n")
	}

	sort.Strings(preRemapLines)
	sort.Strings(textLines)
	sort.Strings(postRemapLines)

	// Prepend any pre remap lines
	remapLinesAll := append(preRemapLines, textLines...)

	// Append on any post raw remap lines
	remapLinesAll = append(remapLinesAll, postRemapLines...)

	text := header + strings.Join(remapLinesAll, "")
	return text, warnings, nil
}

// getServerConfigRemapDotConfigForEdge returns the remap lines, any warnings, and any error.
func getServerConfigRemapDotConfigForEdge(
	profilesRemapConfigParams map[int][]tc.ParameterV5,
	serverPackageParamData map[string]string, // map[paramName]paramVal for this server, config file 'package'
	dses []DeliveryService,
	dsRegexes map[tc.DeliveryServiceName][]tc.DeliveryServiceRegex,
	atsMajorVersion uint,
	header string,
	server *Server,
	anyCastPartners map[string][]string,
	nameTopologies map[TopologyName]tc.TopologyV5,
	cacheGroups map[tc.CacheGroupName]tc.CacheGroupNullableV5,
	serverCapabilities map[int]map[ServerCapability]struct{},
	dsRequiredCapabilities map[int]map[ServerCapability]struct{},
	cdnDomain string,
	configDir string,
	opts *RemapDotConfigOpts,
) (string, []string, error) {
	warnings := []string{}
	textLines := []string{}
	preRemapLines := []string{}
	postRemapLines := []string{}

	for _, ds := range dses {
		if !hasRequiredCapabilities(serverCapabilities[server.ID], dsRequiredCapabilities[*ds.ID]) {
			continue
		}

		topology, hasTopology := nameTopologies[TopologyName(*ds.Topology)]
		if *ds.Topology != "" && hasTopology {
			topoIncludesServer, err := topologyIncludesServerNullable(topology, server)
			if err != nil {
				return "", warnings, errors.New("getting topology server inclusion: " + err.Error())
			}
			if !topoIncludesServer {
				continue
			}
		}
		remapText := ""
		if tc.DSType(*ds.Type) == tc.DSTypeAnyMap {
			if ds.RemapText == nil {
				warnings = append(warnings, "ds '"+ds.XMLID+"' is ANY_MAP, but has no remap text - skipping")
				continue
			}
			remapText = *ds.RemapText + "\n"
			textLines = append(textLines, remapText)
			continue
		}

		requestFQDNs, err := GetDSRequestFQDNs(&ds, dsRegexes[tc.DeliveryServiceName(ds.XMLID)], server, anyCastPartners, cdnDomain)
		if err != nil {
			warnings = append(warnings, "error getting ds '"+ds.XMLID+"' request fqdns, skipping! Error: "+err.Error())
			continue
		}

		for _, requestFQDN := range requestFQDNs {
			remapLines, err := makeEdgeDSDataRemapLines(ds, requestFQDN, server, cdnDomain)
			if err != nil {
				warnings = append(warnings, "DS '"+ds.XMLID+"' - skipping! : "+err.Error())
				continue
			}

			for _, line := range remapLines {
				profileremapConfigParams := []tc.ParameterV5{}
				if ds.ProfileID != nil {
					profileremapConfigParams = profilesRemapConfigParams[*ds.ProfileID]
				}
				dsLines, remapWarns, err := buildEdgeRemapLine(atsMajorVersion, server, serverPackageParamData, remapText, ds, line.From, line.To, profileremapConfigParams, cacheGroups, nameTopologies, configDir, opts)
				warnings = append(warnings, remapWarns...)
				remapText += dsLines.Text

				// Add to pre/post remap lines if this is last tier
				if len(dsLines.Pre) > 0 || len(dsLines.Post) > 0 {
					preRemapLines = append(preRemapLines, dsLines.Pre...)
					postRemapLines = append(postRemapLines, dsLines.Post...)
				}

				if err != nil {
					return "", warnings, err
				}
				remapText += ` # ds '` + ds.XMLID + `' topology '`
				if hasTopology {
					remapText += topology.Name
				}
				remapText += `'` + "\n"
			}
		}

		textLines = append(textLines, remapText)
	}

	sort.Strings(preRemapLines)
	sort.Strings(textLines)
	sort.Strings(postRemapLines)

	remapLinesAll := append(preRemapLines, textLines...)
	remapLinesAll = append(remapLinesAll, postRemapLines...)

	text := header
	text += strings.Join(remapLinesAll, "")
	return text, warnings, nil
}

// RemapLines represents the contents of a remap.config Apache Traffic Server
// configuration file for a given Delivery Service.
type RemapLines struct {
	// Pre contains lines to appear before the "main" text of a Delivery
	// Service's remap.config content.
	Pre []string
	// Text contains the "main" text of a Delivery Service's remap.config
	// content.
	Text string
	// Post contains lines to appear after the "main" text of a Delivery
	// Service's remap.config content.
	Post []string
}

// buildEdgeRemapLine builds the remap line for the given server and delivery service.
// The cacheKeyConfigParams map may be nil, if this ds profile had no cache key config params.
// Returns the remap line, any warnings, and any error.
func buildEdgeRemapLine(
	atsMajorVersion uint,
	server *Server,
	pData map[string]string,
	text string,
	ds DeliveryService,
	mapFrom string,
	mapTo string,
	remapConfigParams []tc.ParameterV5,
	cacheGroups map[tc.CacheGroupName]tc.CacheGroupNullableV5,
	nameTopologies map[TopologyName]tc.TopologyV5,
	configDir string,
	opts *RemapDotConfigOpts,
) (RemapLines, []string, error) {
	warnings := []string{}
	remapLines := RemapLines{}

	// ds = 'remap' in perl
	mapFrom = strings.Replace(mapFrom, `__http__`, server.HostName, -1)

	isLastCache, err := serverIsLastCacheForDS(server, &ds, nameTopologies, cacheGroups)
	if err != nil {
		return remapLines, warnings, errors.New("determining if cache is the last tier: " + err.Error())
	}

	// if this remap is going to a parent, use http not https.
	// cache-to-cache communication inside the CDN is always http (though that's likely to change in the future)
	if !isLastCache {
		mapTo = strings.Replace(mapTo, `https://`, `http://`, -1)
	}

	// template fill for moustache
	//remapTags := NewRemapTags()
	remapTags := RemapTags{}
	remapTags.Source = mapFrom
	remapTags.Destination = mapTo
	strategy := strategyDirective(getStrategyName(ds.XMLID), configDir, opts)
	if strategy != "" {
		remapTags.Strategy = strategy
	}

	if _, hasDSCPRemap := pData["dscp_remap"]; hasDSCPRemap {
		remapTags.Dscp = `@plugin=dscp_remap.so @pparam=` + strconv.Itoa(ds.DSCP)
	} else {
		remapTags.Dscp = `@plugin=header_rewrite.so @pparam=dscp/set_dscp_` + strconv.Itoa(ds.DSCP) + ".config"
	}

	if *ds.Topology != "" {
		topoTxt, err := makeDSTopologyHeaderRewriteTxt(ds, tc.CacheGroupName(server.CacheGroup), nameTopologies[TopologyName(*ds.Topology)], cacheGroups)
		if err != nil {
			return remapLines, warnings, err
		}
		remapTags.HeaderRewrite = topoTxt
	} else if (ds.EdgeHeaderRewrite != nil && *ds.EdgeHeaderRewrite != "") || (ds.ServiceCategory != nil && *ds.ServiceCategory != "") || (ds.MaxOriginConnections != nil && *ds.MaxOriginConnections != 0) {
		remapTags.HeaderRewrite = `@plugin=header_rewrite.so @pparam=` + edgeHeaderRewriteConfigFileName(ds.XMLID)
	}

	dsConfigParamsMap, selfHeal := classifyConfigParams(remapConfigParams)

	if ds.SigningAlgorithm != nil && *ds.SigningAlgorithm != "" {
		if *ds.SigningAlgorithm == tc.SigningAlgorithmURLSig {
			remapTags.Signing = `@plugin=url_sig.so @pparam=url_sig_` + ds.XMLID + ".config" +
				paramsStringFor(dsConfigParamsMap["url_sig.pparam"], &warnings)
		} else if *ds.SigningAlgorithm == tc.SigningAlgorithmURISigning {
			remapTags.Signing = `@plugin=uri_signing.so @pparam=uri_signing_` + ds.XMLID + ".config" +
				paramsStringFor(dsConfigParamsMap["uri_signing.pparam"], &warnings)
		}
	}

	// Raw remap text, this allows the directive hacks
	rawRemapText := ""
	if ds.RemapText != nil {
		rawRemapText = *ds.RemapText
	}

	// Form the cachekey args string, qstring ignore, then
	// remap.config then cachekey.config
	cachekeyTxt := ""
	cachekeyArgs := ""

	if ds.QStringIgnore != nil {
		if *ds.QStringIgnore == tc.QueryStringIgnoreDropAtEdge {
			remapTags.DropQstring = `@plugin=regex_remap.so @pparam=` + RemapConfigDropQstringConfigFile
		} else if *ds.QStringIgnore == tc.QueryStringIgnoreIgnoreInCacheKeyAndPassUp {
			cachekeyArgs = getQStringIgnoreRemap(atsMajorVersion)
		}
	}

	if len(dsConfigParamsMap) > 0 {
		cachekeyArgs += cachekeyArgsFor(dsConfigParamsMap, &warnings)
	}

	if cachekeyArgs != "" {
		cachekeyTxt = `@plugin=cachekey.so` + cachekeyArgs
	}

	// Hack for moving the cachekey directive into the raw remap text
	if strings.Contains(rawRemapText, RemapConfigCachekeyDirective) {
		rawRemapText = strings.Replace(rawRemapText, RemapConfigCachekeyDirective, cachekeyTxt, 1)
	} else if cachekeyTxt != "" {
		remapTags.Cachekey = cachekeyTxt
	}

	regexRemapTxt := ""

	// Note: should use full path here?
	if ds.RegexRemap != nil && *ds.RegexRemap != "" {
		regexRemapTxt = `@plugin=regex_remap.so @pparam=regex_remap_` + ds.XMLID + ".config"
	}

	// Hack for moving the regex_remap directive into the raw remap text
	if strings.Contains(rawRemapText, RemapConfigRegexRemapDirective) {
		rawRemapText = strings.Replace(rawRemapText, RemapConfigRegexRemapDirective, regexRemapTxt, 1)
	} else if regexRemapTxt != "" {
		remapTags.RegexRemap = regexRemapTxt
	}

	rangeReqTxt := ""
	if ds.RangeRequestHandling != nil {
		crr := false

		if *ds.RangeRequestHandling == tc.RangeRequestHandlingBackgroundFetch {
			rangeReqTxt = `@plugin=background_fetch.so @pparam=--config=bg_fetch.config` +
				paramsStringFor(dsConfigParamsMap["background_fetch.pparam"], &warnings)
		} else if *ds.RangeRequestHandling == tc.RangeRequestHandlingSlice && ds.RangeSliceBlockSize != nil {

			rangeReqTxt = `@plugin=slice.so @pparam=--blockbytes=` + strconv.Itoa(*ds.RangeSliceBlockSize) +
				paramsStringFor(dsConfigParamsMap["slice.pparam"], &warnings)
			rangeReqTxt += ` `
			crr = true
		} else if *ds.RangeRequestHandling == tc.RangeRequestHandlingCacheRangeRequest {
			crr = true
		}

		if crr {
			rangeReqTxt += `@plugin=cache_range_requests.so ` + paramsStringFor(dsConfigParamsMap["cache_range_requests.pparam"], &warnings)
			if *ds.RangeRequestHandling == tc.RangeRequestHandlingSlice && !strings.Contains(rangeReqTxt, "--consider-ims") && selfHeal {
				rangeReqTxt += ` @pparam=--consider-ims`
			}
		}
	}

	// Hack for moving the range directive into the raw remap text
	if strings.Contains(rawRemapText, RemapConfigRangeDirective) {
		rawRemapText = strings.Replace(rawRemapText, RemapConfigRangeDirective, rangeReqTxt, 1)
	} else if rangeReqTxt != "" {
		remapTags.RangeRequests = rangeReqTxt
	}

	if rawRemapText != "" {
		remapTags.RawText = rawRemapText
	}

	if ds.FQPacingRate != nil && *ds.FQPacingRate > 0 {
		remapTags.Pacing = `@plugin=fq_pacing.so @pparam=--rate=` + strconv.Itoa(*ds.FQPacingRate)
	}

	// Check for ds parameter template override
	var template *mustache.Template
	if params, ok := dsConfigParamsMap[RemapConfigTemplateFirst]; ok {
		if 0 < len(params) {
			templateString := params[0].Value
			template, err = RemapLineTemplates.parse(templateString)
			if err != nil {
				warnings = append(warnings, "Error decoding override "+RemapConfigTemplateFirst+" "+templateString+" for DS "+ds.XMLID+", falling back!")
			}
		}
	}

	// fall back to default remap config line template
	if template == nil {
		template, err = RemapLineTemplates.parse(DefaultFirstRemapConfigTemplateString)
		if err != nil {
			return remapLines, warnings, errors.New("unable to load default template string: " + DefaultFirstRemapConfigTemplateString + " " + err.Error())
		}
	}

	// Expand the moustache
	var linebytes bytes.Buffer
	err = template.FRender(&linebytes, remapTags)
	if err != nil {
		return remapLines, warnings, errors.New("Error filling edge remap template: " + err.Error())
	}

	//remapLines.Text = removeDeleteMes(linebytes.String())
	remapLines.Text = linebytes.String()

	// Any raw pre or post pend lines?
	if isLastCache {
		remapLines.Pre, remapLines.Post = lastPrePostRemapLinesFor(dsConfigParamsMap, ds.XMLID)
	}

	return remapLines, warnings, nil
}

func strategyDirective(strategyName string, configDir string, opt *RemapDotConfigOpts) string {
	if !opt.UseStrategies {
		return ""
	}
	if !opt.UseStrategiesCore {
		return `@plugin=parent_select.so @pparam=` + filepath.Join(configDir, "strategies.yaml") + ` @pparam=` + strategyName
	}
	return `@strategy=` + strategyName
}

// makeDSTopologyHeaderRewriteTxt returns the appropriate header rewrite remap line text for the given DS on the given server, and any error.
// May be empty, if the DS has no header rewrite for the server's position in the topology.
func makeDSTopologyHeaderRewriteTxt(ds DeliveryService, cg tc.CacheGroupName, topology tc.TopologyV5, cacheGroups map[tc.CacheGroupName]tc.CacheGroupNullableV5) (string, error) {
	placement, err := getTopologyPlacement(cg, topology, cacheGroups, &ds)
	if err != nil {
		return "", errors.New("getting topology placement: " + err.Error())
	}
	txt := ""
	const pluginTxt = `@plugin=header_rewrite.so @pparam=`
	if placement.IsFirstCacheTier && ((ds.FirstHeaderRewrite != nil && *ds.FirstHeaderRewrite != "") || (ds.ServiceCategory != nil && *ds.ServiceCategory != "")) {
		txt += pluginTxt + FirstHeaderRewriteConfigFileName(ds.XMLID) + ` `
	}
	if placement.IsInnerCacheTier && ((ds.InnerHeaderRewrite != nil && *ds.InnerHeaderRewrite != "") || (ds.ServiceCategory != nil && *ds.ServiceCategory != "")) {
		txt += pluginTxt + InnerHeaderRewriteConfigFileName(ds.XMLID) + ` `
	}
	if placement.IsLastCacheTier && ((ds.LastHeaderRewrite != nil && *ds.LastHeaderRewrite != "") || (ds.ServiceCategory != nil && *ds.ServiceCategory != "") || (ds.MaxOriginConnections != nil && *ds.MaxOriginConnections != 0)) {
		txt += pluginTxt + LastHeaderRewriteConfigFileName(ds.XMLID) + ` `
	}
	return txt, nil
}

type remapLine struct {
	From string
	To   string
}

// makeEdgeDSDataRemapLines returns the remap lines for the given server and delivery service.
// Returns nil, if the given server and ds have no remap lines, i.e. the DS match is not a host regex, or has no origin FQDN.
func makeEdgeDSDataRemapLines(
	ds DeliveryService,
	requestFQDN string,
	//	dsRegex tc.DeliveryServiceRegex,
	server *Server,
	cdnDomain string,
) ([]remapLine, error) {
	if ds.Protocol == nil {
		return nil, errors.New("ds had nil protocol")
	}

	remapLines := []remapLine{}
	mapTo := *ds.OrgServerFQDN + "/"

	portStr := ""
	if !tc.DSType(*ds.Type).IsDNS() && server.TCPPort != nil && *server.TCPPort > 0 && *server.TCPPort != 80 {
		portStr = ":" + strconv.Itoa(*server.TCPPort)
	}

	httpsPortStr := ""
	if !tc.DSType(*ds.Type).IsDNS() && server.HTTPSPort != nil && *server.HTTPSPort > 0 && *server.HTTPSPort != 443 {
		httpsPortStr = ":" + strconv.Itoa(*server.HTTPSPort)
	}

	mapFromHTTP := "http://" + requestFQDN + portStr + "/"
	mapFromHTTPS := "https://" + requestFQDN + httpsPortStr + "/"

	if *ds.Protocol == tc.DSProtocolHTTP || *ds.Protocol == tc.DSProtocolHTTPAndHTTPS {
		remapLines = append(remapLines, remapLine{From: mapFromHTTP, To: mapTo})
	}
	if *ds.Protocol == tc.DSProtocolHTTPS || *ds.Protocol == tc.DSProtocolHTTPToHTTPS || *ds.Protocol == tc.DSProtocolHTTPAndHTTPS {
		remapLines = append(remapLines, remapLine{From: mapFromHTTPS, To: mapTo})
	}

	return remapLines, nil
}

func edgeHeaderRewriteConfigFileName(dsName string) string {
	return "hdr_rw_" + dsName + ".config"
}

func midHeaderRewriteConfigFileName(dsName string) string {
	return "hdr_rw_mid_" + dsName + ".config"
}

// getQStringIgnoreRemap returns the remap, whether cachekey was added.
func getQStringIgnoreRemap(atsMajorVersion uint) string {
	if atsMajorVersion < 7 {
		log.Errorf("Unsupport version of ats found %v", atsMajorVersion)
		return ""
	}
	return ` @pparam=--separator= @pparam=--remove-all-params=true @pparam=--remove-path=true @pparam=--capture-prefix-uri=/^([^?]*)/$1/`
}

// makeServerPackageParamData returns a map[paramName]paramVal for this server, config file 'package'.
// Returns the param data, and any warnings.
func makeServerPackageParamData(server *Server, serverParams []tc.ParameterV5) (map[string]string, []string) {
	warnings := []string{}

	serverPackageParamData := map[string]string{}
	for _, param := range serverParams {
		if param.ConfigFile != "package" { // TODO put in const
			continue
		}
		if param.Name == "location" { // TODO put in const
			continue
		}

		paramName := param.Name
		// some files have multiple lines with the same key... handle that with param id.
		if _, ok := serverPackageParamData[param.Name]; ok {
			paramName += "__" + strconv.Itoa(param.ID)
		}
		paramValue := param.Value
		if paramValue == "STRING __HOSTNAME__" {
			paramValue = server.HostName + "." + server.DomainName // TODO strings.Replace to replace all anywhere, instead of just an exact match?
		}

		if val, ok := serverPackageParamData[paramName]; ok {
			if val < paramValue {
				warnings = append(warnings, "got multiple parameters for server package name '"+paramName+"' - ignoring '"+paramValue+"'")
				continue
			} else {
				warnings = append(warnings, "got multiple parameters for server package name '"+paramName+"' - ignoring '"+val+"'")
			}
		}
		serverPackageParamData[paramName] = paramValue
	}
	return serverPackageParamData, warnings
}

// remapFilterDSes filters Delivery Services to be used to generate remap.config for the given server.
// Returned DSes are guaranteed to have a non-nil XMLID, Type, DSCP, ID, Active, and Topology.
// If a DS has a nil Topology, OrgServerFQDN, FirstHeaderRewrite, InnerHeaderRewrite, or LastHeaderRewrite, "" is assigned.
// Returns the filtered delivery services, and any warnings.
func remapFilterDSes(server *Server, dss []DeliveryServiceServer, dses []DeliveryService) ([]DeliveryService, []string) {
	warnings := []string{}
	isMid := strings.HasPrefix(server.Type, string(tc.CacheTypeMid))

	serverIDs := map[int]struct{}{}
	if !isMid {
		// mids use all servers, so pass empty=all. Edges only use this current server
		serverIDs[server.ID] = struct{}{}
	}

	dsIDs := map[int]struct{}{}
	for _, ds := range dses {
		if ds.ID == nil {
			// TODO log error?
			continue
		}
		dsIDs[*ds.ID] = struct{}{}
	}

	dsServers := filterDSS(dss, dsIDs, serverIDs)

	dssMap := map[int]map[int]struct{}{} // set of map[dsID][serverID]
	for _, dss := range dsServers {
		if dssMap[dss.DeliveryService] == nil {
			dssMap[dss.DeliveryService] = map[int]struct{}{}
		}
		dssMap[dss.DeliveryService][dss.Server] = struct{}{}
	}

	useInactive := false
	if !isMid {
		// mids get inactive DSes, edges don't. This is how it's always behaved, not necessarily how it should.
		useInactive = true
	}

	filteredDSes := []DeliveryService{}
	for _, ds := range dses {
		if ds.Topology == nil {
			ds.Topology = util.StrPtr("")
		}
		if ds.OrgServerFQDN == nil {
			ds.OrgServerFQDN = util.StrPtr("")
		}
		if ds.FirstHeaderRewrite == nil {
			ds.FirstHeaderRewrite = util.StrPtr("")
		}
		if ds.InnerHeaderRewrite == nil {
			ds.InnerHeaderRewrite = util.StrPtr("")
		}
		if ds.LastHeaderRewrite == nil {
			ds.LastHeaderRewrite = util.StrPtr("")
		}
		if ds.XMLID == "" {
			warnings = append(warnings, "got Delivery Service with nil XMLID, skipping!")
			continue
		} else if ds.Type == nil {
			warnings = append(warnings, "got Delivery Service '"+ds.XMLID+"'  with nil Type, skipping!")
			continue
		} else if ds.ID == nil {
			warnings = append(warnings, "got Delivery Service '"+ds.XMLID+"'  with nil ID, skipping!")
			continue
		} else if ds.Active == "" {
			warnings = append(warnings, "got Delivery Service '"+ds.XMLID+"'  with nil Active, skipping!")
			continue
		} else if _, ok := dssMap[*ds.ID]; !ok && *ds.Topology == "" {
			continue // normal, not an error, this DS just isn't assigned to our Cache
		} else if !useInactive && ds.Active != tc.DSActiveStateActive {
			continue // normal, not an error, DS just isn't active and we aren't including inactive DSes
		}
		filteredDSes = append(filteredDSes, ds)
	}
	return filteredDSes, warnings
}

// makeDSProfilesConfigParams returns a map[ProfileID][ParamName]ParamValue for the cache key params for each profile.
// Returns the params, any warnings, and any error.
func makeDSProfilesConfigParams(server *Server, dses []DeliveryService, remapConfigParams []tc.ParameterV5) (map[int][]tc.ParameterV5, []string, error) {
	warnings := []string{}
	dsConfigParamsWithProfiles, err := tcParamsToParamsWithProfiles(remapConfigParams)
	if err != nil {
		return nil, warnings, errors.New("decoding cache key parameter profiles: " + err.Error())
	}

	configParamsWithProfilesMap := parameterWithProfilesToMap(dsConfigParamsWithProfiles)

	dsProfileNamesToIDs := map[string]int{}
	for _, ds := range dses {
		if ds.ProfileID == nil || ds.ProfileName == nil {
			continue // TODO log
		}
		dsProfileNamesToIDs[*ds.ProfileName] = *ds.ProfileID
	}

	dsProfilesConfigParams := map[int][]tc.ParameterV5{}
	for _, param := range configParamsWithProfilesMap {
		for dsProfileName, dsProfileID := range dsProfileNamesToIDs {
			if _, ok := param.ProfileNames[dsProfileName]; ok {
				if _, ok := dsProfilesConfigParams[dsProfileID]; !ok {
					dsProfilesConfigParams[dsProfileID] = []tc.ParameterV5{}
				}
				dsProfilesConfigParams[dsProfileID] = append(dsProfilesConfigParams[dsProfileID], param.ParameterV5)
			}
		}
	}
	return dsProfilesConfigParams, warnings, nil
}

type deliveryServiceRegexesSortByTypeThenSetNum []tc.DeliveryServiceRegex

func (r deliveryServiceRegexesSortByTypeThenSetNum) Len() int { return len(r) }
func (r deliveryServiceRegexesSortByTypeThenSetNum) Less(i, j int) bool {
	if rc := strings.Compare(r[i].Type, r[j].Type); rc != 0 {
		return rc < 0
	}
	return r[i].SetNumber < r[j].SetNumber
}
func (r deliveryServiceRegexesSortByTypeThenSetNum) Swap(i, j int) { r[i], r[j] = r[j], r[i] }

func MakeDSRegexMap(regexes []tc.DeliveryServiceRegexes) map[tc.DeliveryServiceName][]tc.DeliveryServiceRegex {
	dsRegexMap := map[tc.DeliveryServiceName][]tc.DeliveryServiceRegex{}
	for _, dsRegex := range regexes {
		sort.Sort(deliveryServiceRegexesSortByTypeThenSetNum(dsRegex.Regexes))
		dsRegexMap[tc.DeliveryServiceName(dsRegex.DSName)] = dsRegex.Regexes
	}
	return dsRegexMap
}

func GetAnyCastPartners(server *Server, servers []Server) map[string][]string {
	anyCastIPs := make(map[string][]string)
	for _, int := range server.Interfaces {
		if int.Name == "lo" {
			for _, addr := range int.IPAddresses {
				anyCastIPs[addr.Address] = []string{}
			}
		}
	}
	for _, srv := range servers {
		if server.HostName == srv.HostName {
			continue
		}
		for _, int := range srv.Interfaces {
			if int.Name == "lo" {
				for _, address := range int.IPAddresses {
					if _, ok := anyCastIPs[address.Address]; ok && address.ServiceAddress {
						anyCastIPs[address.Address] = append(anyCastIPs[address.Address], srv.HostName)
					}
				}
			}
		}
	}
	return anyCastIPs
}

type keyVal struct {
	Key string
	Val string
}

type keyVals []keyVal

func (ks keyVals) Len() int      { return len(ks) }
func (ks keyVals) Swap(i, j int) { ks[i], ks[j] = ks[j], ks[i] }
func (ks keyVals) Less(i, j int) bool {
	if ks[i].Key != ks[j].Key {
		return ks[i].Key < ks[j].Key
	}
	return ks[i].Val < ks[j].Val
}

// GetDSRequestFQDNs returns the FQDNs that clients will request from the edge.
func GetDSRequestFQDNs(ds *DeliveryService, regexes []tc.DeliveryServiceRegex, server *Server, anyCastPartners map[string][]string, cdnDomain string) ([]string, error) {
	if server.HostName == "" {
		return nil, errors.New("server missing hostname")
	}

	fqdns := []string{}
	seenFQDNs := map[string]struct{}{}
	for _, dsRegex := range regexes {
		if tc.DSMatchType(dsRegex.Type) != tc.DSMatchTypeHostRegex || ds.OrgServerFQDN == nil || *ds.OrgServerFQDN == "" {
			continue
		}
		if dsRegex.Pattern == "" {
			return nil, errors.New("ds missing regex pattern")
		}
		if ds.Protocol == nil {
			return nil, errors.New("ds missing protocol")
		}
		if cdnDomain == "" {
			return nil, errors.New("ds missing domain")
		}

		hostRegex := dsRegex.Pattern

		fqdn, err := makeFQDN(hostRegex, ds, server.HostName, cdnDomain)
		if err != nil {
			return nil, err
		}
		fqdns = append(fqdns, fqdn)

		if tc.DSType(*ds.Type).IsHTTP() && strings.HasSuffix(hostRegex, `.*`) {
			for _, ip := range anyCastPartners {
				for _, hn := range ip {
					fqdn, err := makeFQDN(hostRegex, ds, hn, cdnDomain)
					if err != nil {
						return nil, err
					}
					if _, ok := seenFQDNs[fqdn]; ok {
						continue
					}
					seenFQDNs[fqdn] = struct{}{}
					fqdns = append(fqdns, fqdn)
				}
			}
		}
	}
	return fqdns, nil
}

func makeFQDN(hostRegex string, ds *DeliveryService, server string, cdnDomain string) (string, error) {
	fqdn := hostRegex
	if strings.HasSuffix(hostRegex, `.*`) {
		re := hostRegex
		re = strings.Replace(re, `\`, ``, -1)
		re = strings.Replace(re, `.*`, ``, -1)

		hName := server
		if ds.Type != nil && tc.DSType(*ds.Type).IsDNS() {
			if ds.RoutingName == "" {
				return "", errors.New("ds is dns, but missing routing name")
			}
			hName = ds.RoutingName
		}
		fqdn = hName + re + cdnDomain
	}
	return fqdn, nil
}

func serverIsLastCacheForDS(server *Server, ds *DeliveryService, topologies map[TopologyName]tc.TopologyV5, cacheGroups map[tc.CacheGroupName]tc.CacheGroupNullableV5) (bool, error) {
	if ds.Topology != nil && strings.TrimSpace(*ds.Topology) != "" {
		if server.CacheGroup == "" {
			return false, errors.New("Server has no CacheGroup")
		}
		topology, ok := topologies[TopologyName(*ds.Topology)]
		if !ok {
			return false, errors.New("DS topology '" + *ds.Topology + "' not found in topologies")
		}
		topoPlacement, err := getTopologyPlacement(tc.CacheGroupName(server.CacheGroup), topology, cacheGroups, ds)
		if err != nil {
			return false, errors.New("getting topology placement: " + err.Error())
		}
		return topoPlacement.IsLastCacheTier, nil
	}

	return noTopologyServerIsLastCacheForDS(server, ds, cacheGroups), nil
}

// noTopologyServerIsLastCacheForDS returns whether the server is the last tier for the DS, if the DS has no Topology.
// This helper MUST NOT be called if the DS has a Topology. It does not check.
func noTopologyServerIsLastCacheForDS(server *Server, ds *DeliveryService, cgs map[tc.CacheGroupName]tc.CacheGroupNullableV5) bool {
	if strings.HasPrefix(server.Type, tc.MidTypePrefix) {
		return true // if the type is "MID" it's always the last cache for non-topologies
	}
	if !tc.DSType(*ds.Type).UsesMidCache() {
		return true // if this DS type never uses mids (for pre-topology type-parentage), it's the last cache
	}

	// pre-topology parentage is based on Cachegroups

	if server.CacheGroup == "" {
		// if the server has no CG (which TO shouldn't allow), it can't possibly have parents.
		return true
	}

	cg, ok := cgs[tc.CacheGroupName(server.CacheGroup)]
	if !ok {
		// if the server's CG doesn't exist (which TO shouldn't allow), it can't possibly have parents.
		return true
	}

	if cg.ParentName == nil || *cg.ParentName == "" {
		// if the server's CG has no parents, it's going direct to the origin
		return true
	}

	parentCG, ok := cgs[tc.CacheGroupName(*cg.ParentName)]
	if !ok {
		// if the server's parent CG doesn't exist (which TO shouldn't allow), it can't possibly have parents.
		return true
	}

	if parentCG.Type == nil {
		// if the server's parent CG has no type (which TO shouldn't allow), then it must not be a cache, so this server is the last cache tier
		return true
	}

	if *parentCG.Type != tc.CacheGroupEdgeTypeName && *parentCG.Type != tc.CacheGroupMidTypeName {
		// if the server's parent CG isn't a cache, then this server is the last cache tier.
		return true
	}

	// at this point, this server's CG has a cache parent,
	// so this server isn't the last cache tier
	return false
}
