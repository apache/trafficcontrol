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
const ParentConfigParamMergeGroups = "merge_parent_groups"

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
const StrategyConfigUsePeering = "use_peering"

// same across DS
const ParentConfigParamQString = "qstring"

type ParentConfigRetryKeys struct {
	Algorithm                 string
	SecondaryMode             string
	ParentRetry               string
	MaxSimpleRetries          string
	SimpleRetryResponses      string
	MaxUnavailableRetries     string
	UnavailableRetryResponses string
}

func MakeParentConfigRetryKeysWithPrefix(prefix string) ParentConfigRetryKeys {
	return ParentConfigRetryKeys{
		Algorithm:                 prefix + "algorithm",
		SecondaryMode:             prefix + "try_all_primaries_before_secondary",
		ParentRetry:               prefix + "parent_retry",
		MaxSimpleRetries:          prefix + "max_simple_retries",
		SimpleRetryResponses:      prefix + "simple_server_retry_responses",
		MaxUnavailableRetries:     prefix + "max_unavailable_server_retries",
		UnavailableRetryResponses: prefix + "unavailable_server_retry_responses",
	}
}

var ParentConfigRetryKeysFirst = MakeParentConfigRetryKeysWithPrefix("first.")
var ParentConfigRetryKeysInner = MakeParentConfigRetryKeysWithPrefix("inner.")
var ParentConfigRetryKeysLast = MakeParentConfigRetryKeysWithPrefix("last.")

var ParentConfigRetryKeysMSO = MakeParentConfigRetryKeysWithPrefix("mso.")
var ParentConfigRetryKeysDefault = MakeParentConfigRetryKeysWithPrefix("")

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

	// ATSMajorVersion is the integral major version of Apache Traffic server,
	// used to generate the proper config for the proper version.
	//
	// If omitted or 0, the major version will be read from the Server's Profile Parameter config file 'package' name 'trafficserver'. If no such Parameter exists, the ATS version will default to 5.
	// This was the old Traffic Control behavior, before the version was specifiable externally.
	//
	ATSMajorVersion uint
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
	warnings := []string{}
	atsMajorVersion := getATSMajorVersion(opt.ATSMajorVersion, tcServerParams, &warnings)

	parentAbstraction, dataWarns, err := makeParentDotConfigData(
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
		atsMajorVersion,
	)
	warnings = append(warnings, dataWarns...)
	if err != nil {
		return Cfg{}, makeErr(warnings, err.Error())
	}

	text, paWarns, err := parentAbstractionToParentDotConfig(parentAbstraction, opt, atsMajorVersion)
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

// CreateTopology creates an on the fly topology for this server and non topology delivery service.
func CreateTopology(server *Server, ds DeliveryService, nameTopologies map[TopologyName]tc.Topology, ocgmap map[OriginHost][]string) (string, tc.Topology, []string) {

	topoName := ""
	topo := tc.Topology{}
	warns := []string{}

	orgFQDNStr := *ds.OrgServerFQDN
	orgURI, orgWarns, err := getOriginURI(orgFQDNStr)
	warns = append(warns, orgWarns...)
	if err != nil {
		warns = append(warns, "DS '"+*ds.XMLID+"' has malformed origin URI: '"+orgFQDNStr+"': skipping!"+err.Error())
		return topoName, topo, warns
	}

	// use the topology name for the fqdn
	cgnames, ok := ocgmap[OriginHost(orgURI.Hostname())]
	if !ok {
		cgnames, ok = ocgmap[OriginHost(deliveryServicesAllParentsKey)]
		if !ok {
			warns = append(warns, "DS '"+*ds.XMLID+"' has no parent cache groups! Skipping!")
			return topoName, topo, warns
		}
	}

	// Manufactured topology
	topoName = "otf_" + *ds.XMLID

	// ensure name is unique
	if _, ok := nameTopologies[TopologyName(topoName)]; ok {
		warns = append(warns, "Found collision for topo name '"+topoName+"' for ds: '", *ds.XMLID+"'")
		topoName = topoName + "_"
	}

	topo = tc.Topology{Name: topoName}

	if IsGoDirect(ds) {
		node := tc.TopologyNode{
			Cachegroup: *server.Cachegroup,
		}
		topo.Nodes = append(topo.Nodes, node)
	} else {
		// If mid cache group, insert fake edge cache group.
		// This is incorrect if there are multiple MID tiers.
		pind := 1
		if strings.HasPrefix(server.Type, tc.MidTypePrefix) {
			parents := []int{pind}
			pind++
			edgeNode := tc.TopologyNode{
				Cachegroup: "fake_edgecg",
				Parents:    parents,
			}
			topo.Nodes = append(topo.Nodes, edgeNode)
		}

		parents := []int{}
		for ind := 0; ind < len(cgnames); ind++ {
			parents = append(parents, pind)
			pind++
		}

		node := tc.TopologyNode{
			Cachegroup: *server.Cachegroup,
			Parents:    parents,
		}
		topo.Nodes = append(topo.Nodes, node)

		for _, cg := range cgnames {
			topo.Nodes = append(topo.Nodes, tc.TopologyNode{Cachegroup: cg})
		}
	}
	return topoName, topo, warns
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
	atsMajorVersion uint,
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
	} else if len(server.ProfileNames) == 0 {
		return nil, warnings, errors.New("server Profile missing")
	} else if server.TCPPort == nil {
		return nil, warnings, errors.New("server TCPPort missing")
	}

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

	// serverParams are the parent.config params for this particular server
	serverParams := getServerParentConfigParams(server, parentConfigParams)

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

	cgPeers := map[int]serverWithParams{}   // map[serverID]server
	cgServers := map[int]serverWithParams{} // map[serverID]server
	for _, sv := range serversWithParams {
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
		// save cachegroup peer servers
		if *sv.CDNName == *server.CDNName && *sv.Cachegroup == *server.Cachegroup {
			if *sv.Status == string(tc.CacheStatusReported) || *sv.Status == string(tc.CacheStatusOnline) {
				if _, ok := cgPeers[*sv.ID]; !ok {
					cgPeers[*sv.ID] = sv
				}
			}
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

	// save the cache group peers
	for _, v := range cgPeers {
		peer := &ParentAbstractionServiceParent{}
		peer.FQDN = *v.HostName + "." + *v.DomainName
		peer.Port = *v.TCPPort
		peer.Weight = 0.999
		parentAbstraction.Peers = append(parentAbstraction.Peers, peer)
	}

	sort.Sort(peersSort(parentAbstraction.Peers))

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

	originServers, orgProfWarns, err := getOriginServers(cgServers, parentServerDSes, dses, serverCapabilities)
	warnings = append(warnings, orgProfWarns...)
	if err != nil {
		return nil, warnings, errors.New("getting origin servers and profile caches: " + err.Error())
	}

	parentInfos, piWarns := makeParentInfo(serverParentCGData, serverCDNDomain, originServers, serverCapabilities)
	warnings = append(warnings, piWarns...)

	dsOrigins, dsOriginWarns := makeDSOrigins(dss, dses, servers)
	warnings = append(warnings, dsOriginWarns...)

	ocgmap := map[OriginHost][]string{}

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

		// manufacture a topology for this DS.
		if ds.Topology == nil || *ds.Topology == "" {

			// only populate if there are non topology ds's
			if len(ocgmap) == 0 {
				ocgmap = makeOCGMap(parentInfos)
				if len(ocgmap) == 0 {
					ocgmap[""] = []string{}
				}
			}

			topoName, topo, warns := CreateTopology(server, ds, nameTopologies, ocgmap)

			warnings = append(warnings, warns...)
			if topoName == "" {
				continue
			}

			// check if topology already exists
			nameTopologies[TopologyName(topoName)] = topo
			ds.Topology = util.StrPtr(topoName)
		}

		isMSO := ds.MultiSiteOrigin != nil && *ds.MultiSiteOrigin

		pasvc, topoWarnings, err := getTopologyParentConfigLine(
			server,
			serversWithParams,
			&ds,
			serverParams,
			parentConfigParams,
			nameTopologies,
			serverCapabilities,
			dsRequiredCapabilities,
			cacheGroups,
			profileParentConfigParams,
			isMSO,
			atsMajorVersion,
			dsOrigins[DeliveryServiceID(*ds.ID)],
			opt.AddComments,
		)
		warnings = append(warnings, topoWarnings...)
		if err != nil {
			// we don't want to fail generation with an error if one ds is malformed
			warnings = append(warnings, err.Error()) // getTopologyParentConfigLine includes error context
			continue
		}

		if pasvc != nil { // will be nil with no error if this server isn't in the Topology, or if it doesn't have the Required Capabilities
			parentAbstraction.Services = append(parentAbstraction.Services, pasvc)
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
		parents, secondaryParents, secondaryMode, parentWarns := getParentStrs(invalidDS, dsRequiredCapabilities, parentInfos[deliveryServicesAllParentsKey], atsMajorVersion, tryAllPrimariesBeforeSecondary)
		warnings = append(warnings, parentWarns...)

		defaultDestText.DestDomain = `.`
		defaultDestText.Parents = parents
		// defaultDestText = `dest_domain=. ` + parents

		if serverParams[ParentConfigRetryKeysDefault.Algorithm] == tc.AlgorithmConsistentHash {
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

type parentInfo struct {
	Host            string
	Port            int
	Domain          string
	Weight          float64
	UseIP           bool
	Rank            int
	IP              string
	Cachegroup      string
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

// Returns a map of parent cache groups names per origin host.
func makeOCGMap(opis map[OriginHost][]parentInfo) map[OriginHost][]string {
	ocgnames := map[OriginHost][]string{}

	for host, pis := range opis {
		cgnames := make(map[string]struct{})
		for _, pi := range pis {
			cgnames[string(pi.Cachegroup)] = struct{}{}
		}

		for cg, _ := range cgnames {
			ocgnames[host] = append(ocgnames[host], cg)
		}
	}

	return ocgnames
}

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
	Params parentServerParams
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

type parentServerParams struct {
	Weight     string
	Port       int
	UseIP      bool
	Rank       int
	NotAParent bool
}

const DefaultParentWeight = 0.999

func defaultParentServerParams() parentServerParams {
	return parentServerParams{
		Weight:     strconv.FormatFloat(DefaultParentWeight, 'f', 3, 64),
		Port:       0,
		UseIP:      false,
		Rank:       1,
		NotAParent: false,
	}
}

type originURI struct {
	Scheme string
	Host   string
	Port   string
}

// TODO change, this is terrible practice, using a hard-coded key. What if there were a delivery service named "all_parents" (transliterated Perl)
const deliveryServicesAllParentsKey = "all_parents"

type parentDSParams struct {
	Algorithm           ParentAbstractionServiceRetryPolicy
	QueryStringHandling string

	HasRetryParams                  bool
	ParentRetry                     string
	MaxSimpleRetries                string
	MaxUnavailableServerRetries     string
	SimpleServerRetryResponses      string
	UnavailableServerRetryResponses string
	TryAllPrimariesBeforeSecondary  bool

	UsePeering  bool
	MergeGroups []string
}

// FillParentRetries populates the parentDSParams retries values from the ds parameters map for given ds parameter keys.
// Returns if any params found and any warnings.
func (dsp *parentDSParams) FillParentRetries(keys ParentConfigRetryKeys, dsParams map[string]string, dsid string) (bool, []string) {
	var warnings []string
	hasValues := false

	if v, ok := dsParams[keys.Algorithm]; ok && strings.TrimSpace(v) != "" {
		policy := ParentSelectAlgorithmToParentAbstractionServiceRetryPolicy(v)
		if policy != ParentAbstractionServiceRetryPolicyInvalid {
			dsp.Algorithm = policy
			hasValues = true
		} else {
			warnings = append(warnings, "DS '"+dsid+"' had malformed "+keys.Algorithm+" parameter '"+v+"', not using!")
		}
	}

	if v, ok := dsParams[keys.ParentRetry]; ok {
		dsp.ParentRetry = v
		hasValues = true
	}

	if v, ok := dsParams[keys.MaxSimpleRetries]; ok {
		dsp.MaxSimpleRetries = v
		hasValues = true
	}
	if v, ok := dsParams[keys.MaxUnavailableRetries]; ok {
		dsp.MaxUnavailableServerRetries = v
		hasValues = true
	}

	if v, ok := dsParams[keys.SimpleRetryResponses]; ok {
		if unavailableServerRetryResponsesValid(v) {
			dsp.SimpleServerRetryResponses = v
			hasValues = true
		} else {
			warnings = append(warnings, "DS '"+dsid+"' had malformed "+keys.SimpleRetryResponses+" parameter '"+v+"', not using!")
		}
	}

	if v, ok := dsParams[keys.UnavailableRetryResponses]; ok {
		if unavailableServerRetryResponsesValid(v) {
			dsp.UnavailableServerRetryResponses = v
			hasValues = true
		} else {
			warnings = append(warnings, "DS '"+dsid+"' had malformed "+keys.UnavailableRetryResponses+" parameter '"+v+"', not using!")
		}
	}

	if v, ok := dsParams[keys.SecondaryMode]; ok {
		if v == "false" {
			dsp.TryAllPrimariesBeforeSecondary = false
		} else {
			dsp.TryAllPrimariesBeforeSecondary = true
			if v != "" {
				warnings = append(warnings, "DS '"+dsid+"' had Parameter "+keys.SecondaryMode+" which is used if it exists, the value is ignored (unless false) ! Non-empty value '"+v+"' will be ignored!")
			}
		}
	}

	return hasValues, warnings
}

// getDSParams returns the Delivery Service Profile Parameters used in parent.config, and any warnings.
// If Parameters don't exist, defaults are returned. Non-MSO Delivery Services default to no custom retry logic (we should reevaluate that).
// Note these Parameters are only used for MSO for legacy DeliveryServiceServers DeliveryServices.
//
//	Topology DSes use them for all DSes, MSO and non-MSO.
func getParentDSParams(ds DeliveryService, profileParentConfigParams map[string]map[string]string, serverPlacement TopologyPlacement, isMSO bool) (parentDSParams, []string) {
	warnings := []string{}
	params := parentDSParams{}

	// Default values for origin facing MSO tier
	if serverPlacement.IsLastCacheTier && isMSO {
		params.Algorithm = ParentConfigDSParamDefaultMSOAlgorithm
		params.ParentRetry = ParentConfigDSParamDefaultMSOParentRetry
		params.UnavailableServerRetryResponses = ParentConfigDSParamDefaultMSOUnavailableServerRetryResponses
		params.MaxSimpleRetries = strconv.Itoa(ParentConfigDSParamDefaultMaxSimpleRetries)
		params.MaxUnavailableServerRetries = strconv.Itoa(ParentConfigDSParamDefaultMaxUnavailableServerRetries)
		params.HasRetryParams = true
	}

	if ds.ProfileName == nil || *ds.ProfileName == "" {
		return params, warnings
	}
	dsParams, ok := profileParentConfigParams[*ds.ProfileName]
	if !ok {
		return params, warnings
	}
	if val, ok := dsParams[StrategyConfigUsePeering]; ok {
		if val == "true" {
			params.UsePeering = true
		}
	}

	// the following may be blank, no default
	params.QueryStringHandling = dsParams[ParentConfigParamQStringHandling]
	params.MergeGroups = strings.Split(dsParams[ParentConfigParamMergeGroups], " ")

	// progressively fill in the params
	if serverPlacement.IsLastCacheTier {
		// mso. prefix lowest priority
		if isMSO {
			hasVals, warns := params.FillParentRetries(ParentConfigRetryKeysMSO, dsParams, *ds.XMLID)
			params.HasRetryParams = params.HasRetryParams || hasVals
			warnings = append(warnings, warns...)
		}

		// no prefix next priority
		hasVals, warns := params.FillParentRetries(ParentConfigRetryKeysDefault, dsParams, *ds.XMLID)
		warnings = append(warnings, warns...)
		params.HasRetryParams = params.HasRetryParams || hasVals

		// last. prefix highest priority
		hasVals, warns = params.FillParentRetries(ParentConfigRetryKeysLast, dsParams, *ds.XMLID)
		warnings = append(warnings, warns...)
		params.HasRetryParams = params.HasRetryParams || hasVals

	} else if serverPlacement.IsInnerCacheTier {
		hasVals, warns := params.FillParentRetries(ParentConfigRetryKeysInner, dsParams, *ds.XMLID)
		warnings = append(warnings, warns...)

		// Normal inner behavior has no parent retry strings
		params.HasRetryParams = hasVals
	} else { // if serverPlacement.IsFirstCacheTier {
		hasVals, warns := params.FillParentRetries(ParentConfigRetryKeysFirst, dsParams, *ds.XMLID)
		warnings = append(warnings, warns...)

		// Normal first behavior has no parent retry strings
		params.HasRetryParams = hasVals
	}

	return params, warnings
}

// getTopologyParentConfigLine returns the topology parent.config line, any warnings, and any error
// If the given DS is not used by the server, returns a nil ParentAbstractionService and nil error.
func getTopologyParentConfigLine(
	server *Server,
	serversWithParams []serverWithParams,
	ds *DeliveryService,
	serverParams map[string]string,
	parentConfigParams []parameterWithProfilesMap, // all params with configFile parent.config
	nameTopologies map[TopologyName]tc.Topology,
	serverCapabilities map[int]map[ServerCapability]struct{},
	dsRequiredCapabilities map[int]map[ServerCapability]struct{},
	cacheGroups map[tc.CacheGroupName]tc.CacheGroupNullable,
	profileParentConfigParams map[string]map[string]string,
	isMSO bool,
	atsMajorVersion uint,
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

	dsParams, dswarns := getParentDSParams(*ds, profileParentConfigParams, serverPlacement, isMSO)
	warnings = append(warnings, dswarns...)

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

	pasvc := &ParentAbstractionService{}
	pasvc.Name = *ds.XMLID
	pasvc.Comment = makeParentComment(addComments, *ds.XMLID, *ds.Topology)
	pasvc.DestDomain = orgURI.Hostname()
	pasvc.Port, err = strconv.Atoi(orgURI.Port())
	if err != nil {
		return nil, warnings, fmt.Errorf("parent %v port '%v' was not an integer", orgURI, orgURI.Port())
	}
	// txt += "dest_domain=" + orgURI.Hostname() + " port=" + orgURI.Port()

	parents, secondaryParents, parentWarnings, err := getTopologyParents(server, ds, serversWithParams, parentConfigParams, topology, serverPlacement.IsLastTier, serverCapabilities, dsRequiredCapabilities, dsOrigins, dsParams.MergeGroups)
	warnings = append(warnings, parentWarnings...)

	if err != nil {
		return nil, warnings, errors.New("getting topology parents for '" + *ds.XMLID + "': skipping! " + err.Error())
	}
	if len(parents) == 0 {
		if len(secondaryParents) > 0 {
			warnings = append(warnings, "getting topology parents for '"+*ds.XMLID+"': no parents found! using secondary parents")
			parents = secondaryParents
			secondaryParents = nil
		} else {
			return nil, warnings, errors.New("getting topology parents for '" + *ds.XMLID + "': no parents found! skipping! (Does your Topology have a CacheGroup with no servers in it?)")
		}
	}

	pasvc.Parents = parents
	// txt += ` parent="` + strings.Join(parents, `;`) + `"`
	if len(secondaryParents) > 0 {
		pasvc.SecondaryParents = secondaryParents
		// txt += ` secondary_parent="` + strings.Join(secondaryParents, `;`) + `"`

		secondaryModeStr, secondaryModeWarnings := getSecondaryModeStr(dsParams.TryAllPrimariesBeforeSecondary, atsMajorVersion, tc.DeliveryServiceName(*ds.XMLID))
		warnings = append(warnings, secondaryModeWarnings...)
		// txt += secondaryModeStr
		pasvc.SecondaryMode = secondaryModeStr // TODO convert
	}

	pasvc.RetryPolicy = getTopologyRoundRobin(ds, serverParams, serverPlacement.IsLastCacheTier, dsParams.Algorithm)
	// txt += ` round_robin=` + getTopologyRoundRobin(ds, serverParams, serverPlacement.IsLastCacheTier, dsParams.Algorithm)

	pasvc.GoDirect = getTopologyGoDirect(ds, serverPlacement)
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
	pasvc.IgnoreQueryStringInParentSelection, tqWarns = getTopologyQueryStringIgnore(ds, serverParams, serverPlacement, dsParams.Algorithm, useQueryStringInParentSelection)
	warnings = append(warnings, tqWarns...)

	// txt += ` qstring=` + getTopologyQueryString(ds, serverParams, serverPlacement.IsLastCacheTier, dsParams.Algorithm, dsParams.QueryStringHandling)

	// TODO ensure value is always !goDirect, and determine what to do if not
	// txt += getTopologyParentIsProxyStr(serverPlacement.IsLastCacheTier)

	// TODO convert
	prWarns := dsParams.FillParentSvcRetries(serverPlacement.IsLastCacheTier, atsMajorVersion, pasvc)
	warnings = append(warnings, prWarns...)

	// txt += getParentRetryStr(serverPlacement.IsLastCacheTier, atsMajorVer, dsParams.ParentRetry, dsParams.UnavailableServerRetryResponses, dsParams.MaxSimpleRetries, dsParams.MaxUnavailableServerRetries)
	// txt += "\n"

	if dsParams.UsePeering {
		pasvc.SecondaryMode = ParentAbstractionServiceParentSecondaryModePeering
	}

	return pasvc, warnings, nil
}

// getParentRetries builds the parent retry directive(s).
//
// Returns the MaxSimpleRetries, MaxMarkdownRetries, ErrorResponseCodes, MarkdownResponseCodes.
//
// If atsMajorVersion < 6, "" is returned (ATS 5 and below don't support retry directives).
// If isLastCacheTier is false, "" is returned. This argument exists to simplify usage.
// If parentRetry is "", "" is returned (because the other directives are unused if parent_retry doesn't exist). This is allowed to simplify usage.
// If unavailableServerRetryResponses is not "", it must be valid. Use unavailableServerRetryResponsesValid to check.
// If maxSimpleRetries is "", ParentConfigDSParamDefaultMaxSimpleRetries will be used.
// If maxUnavailableServerRetries is "", ParentConfigDSParamDefaultMaxUnavailableServerRetries will be used.
//
// Does not return errors. If any input is malformed, warnings are returned and that value is set to -1.
func (dsparams parentDSParams) FillParentSvcRetries(isLastCacheTier bool, atsMajorVersion uint, pasvc *ParentAbstractionService) []string {
	warnings := []string{}

	if !dsparams.HasRetryParams || // allow parentRetry to be empty, to simplify usage.
		atsMajorVersion < 6 { // ATS 5 and below don't support parent_retry directives
		// warnings = append(warnings, "ATS 5 doesn't support parent retry, not using parent retry values")
		pasvc.MaxSimpleRetries = -1
		pasvc.MaxMarkdownRetries = -1
		pasvc.ErrorResponseCodes = nil
		pasvc.MarkdownResponseCodes = nil
		return warnings
	}

	// Set initial defaults
	pasvc.MaxSimpleRetries = 0
	pasvc.MaxMarkdownRetries = 0
	pasvc.ErrorResponseCodes = []int{}
	pasvc.MarkdownResponseCodes = []int{}

	if isLastCacheTier {
		pasvc.MaxSimpleRetries = ParentConfigDSParamDefaultMaxSimpleRetries
		pasvc.MaxMarkdownRetries = ParentConfigDSParamDefaultMaxUnavailableServerRetries
	}

	if dsparams.MaxSimpleRetries != "" {
		if retriesint, err := strconv.Atoi(dsparams.MaxSimpleRetries); err == nil {
			pasvc.MaxSimpleRetries = retriesint
		} else {
			warnings = append(warnings, "malformed maxSimpleRetries '"+dsparams.MaxSimpleRetries+"', using default "+strconv.Itoa(pasvc.MaxSimpleRetries))
		}
	}

	if dsparams.MaxUnavailableServerRetries != "" {
		if retriesint, err := strconv.Atoi(dsparams.MaxUnavailableServerRetries); err == nil {
			pasvc.MaxMarkdownRetries = retriesint
		} else {
			warnings = append(warnings, "malformed maxUnavailableServerRetries '"+dsparams.MaxUnavailableServerRetries+"', using default "+strconv.Itoa(pasvc.MaxMarkdownRetries))
		}
	}

	// simple retry responses only supported int ATS for 9.1.x and above
	if simpleResponseInts, err := ParseRetryResponses(dsparams.SimpleServerRetryResponses); err == nil {
		pasvc.ErrorResponseCodes = simpleResponseInts
	} else {
		warnings = append(warnings, "malformed simpleServerRetryResponses '"+dsparams.SimpleServerRetryResponses+"', using default (parse err: "+err.Error()+")")
	}

	if unavailResponseInts, err := ParseRetryResponses(dsparams.UnavailableServerRetryResponses); err == nil {
		pasvc.MarkdownResponseCodes = unavailResponseInts
	} else {
		warnings = append(warnings, "malformed unavailableServerRetryResponses '"+dsparams.UnavailableServerRetryResponses+"', using default (parse err: "+err.Error()+")")
	}

	// TODO make consts
	switch strings.ToLower(strings.TrimSpace(dsparams.ParentRetry)) {
	case "simple_retry":
		if len(pasvc.ErrorResponseCodes) == 0 {
			pasvc.ErrorResponseCodes = DefaultSimpleRetryCodes
		}
	case "unavailable_server_retry":
		if len(pasvc.MarkdownResponseCodes) == 0 {
			pasvc.MarkdownResponseCodes = DefaultUnavailableServerRetryCodes
		}
	case "both":
		if len(pasvc.ErrorResponseCodes) == 0 {
			pasvc.ErrorResponseCodes = DefaultSimpleRetryCodes
		}
		if len(pasvc.MarkdownResponseCodes) == 0 {
			pasvc.MarkdownResponseCodes = DefaultUnavailableServerRetryCodes
		}
	default:
	}

	// txt := ` parent_retry=` + parentRetry
	// if unavailableServerRetryResponses != "" {
	// 	txt += ` unavailable_server_retry_responses=` + unavailableServerRetryResponses
	// }
	// txt += ` max_simple_retries=` + maxSimpleRetries + ` max_unavailable_server_retries=` + maxUnavailableServerRetries
	return warnings
}

// getSecondaryModeStr returns the secondary_mode string, and any warnings.
func getSecondaryModeStr(tryAllPrimariesBeforeSecondary bool, atsMajorVersion uint, ds tc.DeliveryServiceName) (ParentAbstractionServiceParentSecondaryMode, []string) {
	warnings := []string{}
	if !tryAllPrimariesBeforeSecondary {
		return ParentAbstractionServiceParentSecondaryModeDefault, warnings
	}
	if atsMajorVersion < 8 {
		warnings = append(warnings, "DS '"+string(ds)+"' had Parameter "+ParentConfigRetryKeysDefault.SecondaryMode+" but this cache is "+strconv.FormatUint(uint64(atsMajorVersion), 10)+" and secondary_mode isn't supported in ATS until 8. Not using!")
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
	if parentSelectAlg := serverParams[ParentConfigRetryKeysDefault.Algorithm]; ds.OriginShield != nil && *ds.OriginShield != "" && strings.TrimSpace(parentSelectAlg) != "" {
		if policy := ParentSelectAlgorithmToParentAbstractionServiceRetryPolicy(parentSelectAlg); policy != ParentAbstractionServiceRetryPolicyInvalid {
			return policy
		}
	}
	if algorithm != "" {
		return algorithm
	}
	return ParentAbstractionServiceRetryPolicyConsistentHash
}

func getTopologyGoDirect(ds *DeliveryService, serverPlacement TopologyPlacement) bool {
	if serverPlacement.IsLastCacheTier {
		return true
	} else if !serverPlacement.IsLastTier {
		return false
	} else if ds.OriginShield != nil && *ds.OriginShield != "" {
		return true
	} else if ds.MultiSiteOrigin != nil && *ds.MultiSiteOrigin {
		return false
	}
	return true
}

func getTopologyQueryStringIgnore(
	ds *DeliveryService,
	serverParams map[string]string,
	serverPlacement TopologyPlacement,
	algorithm ParentAbstractionServiceRetryPolicy,
	qStringHandling *bool,
) (bool, []string) {
	warnings := []string{}
	if serverPlacement.IsLastCacheTier {
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
func serverParentageParams(sv *Server, allParentConfigParams []parameterWithProfilesMap) (parentServerParams, []string) {
	warnings := []string{}
	// TODO deduplicate with atstccfg/parentdotconfig.go
	parentServerParams := defaultParentServerParams()
	if sv.TCPPort != nil {
		parentServerParams.Port = *sv.TCPPort
	}

	serverParentConfigParams := layerProfilesFromMap(sv.ProfileNames, allParentConfigParams)
	for _, param := range serverParentConfigParams {
		switch param.Name {
		case ParentConfigCacheParamWeight:
			parentServerParams.Weight = param.Value
		case ParentConfigCacheParamPort:
			if i, err := strconv.Atoi(param.Value); err != nil {
				warnings = append(warnings, "port param is not an integer, skipping! : "+err.Error())
			} else {
				parentServerParams.Port = i
			}
		case ParentConfigCacheParamUseIP:
			parentServerParams.UseIP = param.Value == "1"
		case ParentConfigCacheParamRank:

			if i, err := strconv.Atoi(param.Value); err != nil {
				warnings = append(warnings, "rank param is not an integer, skipping! : "+err.Error())
			} else {
				parentServerParams.Rank = i
			}
		case ParentConfigCacheParamNotAParent:
			parentServerParams.NotAParent = param.Value != "false"
		}
	}

	return parentServerParams, warnings
}

func serverParentStr(sv *Server, svParams parentServerParams) (*ParentAbstractionServiceParent, error) {
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

// getTopologyParents returns the parents, secondary parents, any warnings, and any error.
func getTopologyParents(
	server *Server,
	ds *DeliveryService,
	serversWithParams []serverWithParams,
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
		if util.ContainsStr(dsMergeGroups, parentCG) && util.ContainsStr(dsMergeGroups, secondaryParentCG) {
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
	atsMajorVersion uint,
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
	// if atsMajorVersion >= 6 && len(secondaryParentInfo) > 0 {
	// parents = `parent="` + strings.Join(parentInfo, "") + `"`
	// secondaryParents = ` secondary_parent="` + strings.Join(secondaryParentInfo, "") + `"`
	secondaryMode, secondaryModeWarnings := getSecondaryModeStr(tryAllPrimariesBeforeSecondary, atsMajorVersion, dsName)
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
	atsMajorVersion uint,
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
	// if atsMajorVersion >= 6 && msoAlgorithm == "consistent_hash" && len(secondaryParentStr) > 0 {
	// parents = `parent="` + strings.Join(parentInfoTxt, "") + `"`
	// secondaryParents = ` secondary_parent="` + secondaryParentStr + `"`
	secondaryMode, secondaryModeWarnings := getSecondaryModeStr(tryAllPrimariesBeforeSecondary, atsMajorVersion, dsName)
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
	originServers map[OriginHost][]serverWithParams,
	serverCapabilities map[int]map[ServerCapability]struct{},
) (map[OriginHost][]parentInfo, []string) {
	warnings := []string{}
	parentInfos := map[OriginHost][]parentInfo{}

	// note servers also contains an "all" key
	for originHost, servers := range originServers {
		for _, sv := range servers {
			if sv.Params.NotAParent {
				continue
			}
			// Perl has this check, but we only select servers ("deliveryServices" in Perl) with the right CDN in the first place
			// if profile.Domain != serverDomain {
			// 	continue
			// }

			weight, err := strconv.ParseFloat(sv.Params.Weight, 64)
			if err != nil {
				warnings = append(warnings, "server "+*sv.HostName+" profile had malformed weight, using default!")
				weight = DefaultParentWeight
			}

			ipAddr := getServerIPAddress(&sv.Server)
			if ipAddr == nil {
				warnings = append(warnings, "making parent info: got server with no valid IP Address, skipping!")
				continue
			}

			parentInf := parentInfo{
				Host:            *sv.HostName,
				Port:            sv.Params.Port,
				Domain:          *sv.DomainName,
				Weight:          weight,
				UseIP:           sv.Params.UseIP,
				Rank:            sv.Params.Rank,
				IP:              ipAddr.String(),
				Cachegroup:      *sv.Cachegroup,
				PrimaryParent:   serverParentCGData.ParentID == *sv.CachegroupID,
				SecondaryParent: serverParentCGData.SecondaryParentID == *sv.CachegroupID,
				Capabilities:    serverCapabilities[*sv.ID],
			}
			if parentInf.Port < 1 {
				parentInf.Port = *sv.TCPPort
			}

			parentInfos[originHost] = append(parentInfos[originHost], parentInf)
		}
	}
	return parentInfos, warnings
}

// unavailableServerRetryResponsesValid returns whether a unavailable_server_retry_responses parameter is valid for an ATS parent rule.
func unavailableServerRetryResponsesValid(s string) bool {
	_, err := ParseRetryResponses(s)
	return err == nil
}

// getOriginServers returns the origin servers with parameters, any warnings, and any error.
func getOriginServers(
	cgServers map[int]serverWithParams,
	parentServerDSes map[int]map[int]struct{},
	dses []DeliveryService,
	serverCapabilities map[int]map[ServerCapability]struct{},
) (map[OriginHost][]serverWithParams, []string, error) {
	warnings := []string{}
	originServers := map[OriginHost][]serverWithParams{}

	dsIDMap := map[int]DeliveryService{}
	for _, ds := range dses {
		if ds.ID == nil {
			return nil, warnings, errors.New("delivery services got nil ID!")
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
		return nil, warnings, errors.New("getting DS origins: " + err.Error())
	}

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
		} else if len(cgSv.ProfileNames) == 0 {
			warnings = append(warnings, "getting origin servers: got server with no profile names, skipping!")
			continue
		} else if cgSv.CDNID == nil {
			warnings = append(warnings, "getting origin servers: got server with nil CDNID, skipping!")
			continue
		} else if cgSv.DomainName == nil {
			warnings = append(warnings, "getting origin servers: got server with nil DomainName, skipping!")
			continue
		}

		ipAddr := getServerIPAddress(&cgSv.Server)
		if ipAddr == nil {
			warnings = append(warnings, "getting origin servers: got server with no valid IP Address, skipping!")
			continue
		}

		if cgSv.Type == tc.OriginTypeName {
			for dsID, _ := range parentServerDSes[*cgSv.ID] { // map[serverID][]dsID
				orgURI := dsOrigins[dsID]
				if orgURI == nil {
					// warnings = append(warnings, fmt.Sprintf(("ds %v has no origins! Skipping!\n", dsID) // TODO determine if this is normal
					continue
				}
				orgHost := OriginHost(orgURI.Host)
				originServers[orgHost] = append(originServers[orgHost], cgSv)
			}
		} else {
			originServers[deliveryServicesAllParentsKey] = append(originServers[deliveryServicesAllParentsKey], cgSv)
		}
	}

	return originServers, warnings, nil
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
func getServerParentConfigParams(server *Server, allParentConfigParams []parameterWithProfilesMap) map[string]string {
	// We only need parent.config params, don't need all the params on the server
	serverParams := map[string]string{}
	serverParentConfigParams := layerProfilesFromMap(server.ProfileNames, allParentConfigParams)
	for _, pa := range serverParentConfigParams {
		name := pa.Name
		val := pa.Value
		if name == ParentConfigParamQStringHandling ||
			name == ParentConfigRetryKeysDefault.Algorithm ||
			name == ParentConfigParamQString {
			serverParams[name] = val
		}
	}
	return serverParams
}
