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
	"errors"
	"sort"
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
)

// ContentTypeYAML is the MIME type of YAML-format files and text.
//
// Note YAML has no IANA standard mime type. This is one of several common
// usages, and is likely to be the standardized value. If you're reading this,
// please check IANA to see if YAML has been added, and change this to the IANA
// definition if so. Also note we include 'charset=us-ascii' because YAML is
// commonly UTF-8, but ATS is likely to be unable to handle UTF.
//
// TODO: Move this into lib/go-rfc?
const ContentTypeYAML = "application/yaml; charset=us-ascii"

// LineCommentYAML is the string used to start a line comment in YAML syntax.
const LineCommentYAML = LineCommentHash

// SSLServerNameYAMLFileName is an unused constant that probably refers to the
// name of an ATS ssl_server_name.yaml file on disk, the ConfigFile value of
// Parameters that can affect the generation of such a file, or both.
const SSLServerNameYAMLFileName = "ssl_server_name.yaml"

// ContentTypeSSLServerNameYAML is the MIME type of the contents of an
// ssl_server_name.yaml ATS configuration file.
const ContentTypeSSLServerNameYAML = ContentTypeYAML

// LineCommentSSLServerNameYAML is the string used to indicate the start of a
// line comment in the grammar of an ssl_server_name.yaml ATS configuration
// file.
const LineCommentSSLServerNameYAML = LineCommentYAML

// DefaultDefaultEnableH2 is whether Delivery Services will have HTTP/2 enabled by default if they don't have an explicit Parameter, and no Opt is passed to the Make func.
// We disable by default, to prevent potentially enabling broken clients.
const DefaultDefaultEnableH2 = false

// A TLSVersion is a version of the TLS protocol.
//
// TODO: consolidate with the lib/go-tc TLS version constants.
type TLSVersion string

// These are valid, recognized versions of TLS.
const (
	TLSVersion1p0 TLSVersion = "1.0"
	TLSVersion1p1 TLSVersion = "1.1"
	TLSVersion1p2 TLSVersion = "1.2"
	TLSVersion1p3 TLSVersion = "1.3"
)

// TLSVersionInvalid is used to represent an unrecognized or unsupported TLS
// version.
const TLSVersionInvalid = TLSVersion("")

// StringToTLSVersion returns the TLSVersion or TLSVersionInvalid if the string is not a TLS Version enum.
func StringToTLSVersion(st string) TLSVersion {
	switch TLSVersion(st) {
	case TLSVersion1p0:
		return TLSVersion1p0
	case TLSVersion1p1:
		return TLSVersion1p1
	case TLSVersion1p2:
		return TLSVersion1p2
	case TLSVersion1p3:
		return TLSVersion1p3
	}
	return TLSVersionInvalid
}

// tlsVersionsToATS maps TLS version strings to the string used by ATS in
// ssl_server_name.yaml.
//
// TODO: could this be a String method to make TLSVersion a fmt.Stringer and
// remove a package-scope mutable variable?
var tlsVersionsToATS = map[TLSVersion]string{
	TLSVersion1p0: "TLSv1",
	TLSVersion1p1: "TLSv1_1",
	TLSVersion1p2: "TLSv1_2",
	TLSVersion1p3: "TLSv1_3",
}

// SSLServerNameYAMLParamEnableH2 is the Name of a Parameter which, if found on
// the Profile of a Delivery Service and having the ConfigFile value
// "parent.config" - NOT ssl_server_name.yaml, which would cause it to be
// ignored - will enable the use of the HTTP/2 protocol for Delivery Services if
// its value begins with 'T', 't', 'Y', or 'y'.
//
// Examples of Values for these Parameters that would cause HTTP/2 to be enabled
// are:
//
//	true
//	TrUe
//	Yes
//	yno
//	yFalse
//	tFalse
//	talse
//	ts"<>an177otehucau$%*(@YDU)
//
// Examples of Values for these Parameters that would cause HTTP/2 to be
// disabled are:
//
//	false
//	No
//	off
//	on
//	use
//	enable
//	1
//	affirmative
//	I wish to use this feature
//	frue
//	nYes
//	s"<>an177otehucau$%*(@YDU)
const SSLServerNameYAMLParamEnableH2 = "enable_h2"

// SSLServerNameYAMLParamTLSVersions is the Name of a Parameter which, if found
// on the Profile of a Delivery Service and having the ConfigFile value
// "parent.config" - NOT ssl_server_name.yaml, which would cause it to be
// ignored - will set the explicitly allowed TLS versions for communication with
// edge-tier cache servers.
//
// The Value of such Parameters should be a list of TLS versions, which may be
// separated by spaces, new lines, semicolons, or commas (but only one thereof,
// and only consistently a single of those delimeters, not a mixture).
//
// Deprecated: Now that Delivery Services support setting TLS versions to be
// explicitly allowed directly on their configuration, that feature should be
// used rather than this Parameter.
const SSLServerNameYAMLParamTLSVersions = "tls_versions"

// DefaultDefaultTLSVersions is the list of TLS versions to enable by default, if no Parameter exists and no Opt is passed to the Make func.
// By default, we enable all, even insecure versions.
// As a CDN, Traffic Control assumes it should not break clients, and it's the client's responsibility to use secure protocols.
// Note this enables certain downgrade attacks. Operators or tenants concerned about these attacks should disable older TLS versions.
var DefaultDefaultTLSVersions = []TLSVersion{
	TLSVersion1p0,
	TLSVersion1p1,
	TLSVersion1p2,
	TLSVersion1p3,
}

// SSLServerNameYAMLOpts contains settings to configure ssl_server_name.yaml generation options.
type SSLServerNameYAMLOpts struct {
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

// MakeSSLServerNameYAML constructs an ssl_server_name.yaml ATS configuration
// file.
func MakeSSLServerNameYAML(
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
	opt *SSLServerNameYAMLOpts,
) (Cfg, error) {
	if opt == nil {
		opt = &SSLServerNameYAMLOpts{}
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
			dsTxt += "\n" + `  disable_h2: ` + strconv.FormatBool(!sslData.EnableH2)
			dsTxt += "\n" + `  valid_tls_versions_in: [` + strings.Join(tlsVersionsATS, `,`) + `]`

			txt += dsTxt + "\n"
		}

	}

	return Cfg{
		Text:        txt,
		ContentType: ContentTypeSSLServerNameYAML,
		LineComment: LineCommentSSLServerNameYAML,
		Warnings:    warnings,
	}, nil
}

// SSLData has the DS data needed for both sni.yaml (ATS 9+)  and ssl_server_name.yaml (ATS 8).
type SSLData struct {
	DSName       string
	RequestFQDNs []string
	EnableH2     bool
	TLSVersions  []TLSVersion
}

// GetServerSSLData gets the SSLData for all Delivery Services assigned to the given Server, any warnings, and any error.
func GetServerSSLData(
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
	defaultTLSVersions []TLSVersion,
	defaultEnableH2 bool,
) ([]SSLData, []string, error) {
	warnings := []string{}

	if len(server.Profiles) == 0 {
		return nil, warnings, errors.New("this server missing Profiles")
	}

	dsRegexes := MakeDSRegexMap(dsRegexArr)

	parentConfigParamsWithProfiles, err := tcParamsToParamsWithProfiles(tcParentConfigParams)
	if err != nil {
		warnings = append(warnings, "error getting profiles from Traffic Ops Parameters, Parameters will not be considered for generation! : "+err.Error())
		parentConfigParamsWithProfiles = []parameterWithProfiles{}
	}

	profileParentConfigParams := map[string]map[string]string{} // map[profileName][paramName]paramVal
	for _, param := range parentConfigParamsWithProfiles {
		for _, profile := range param.ProfileNames {
			if _, ok := profileParentConfigParams[profile]; !ok {
				profileParentConfigParams[profile] = map[string]string{}
			}
			profileParentConfigParams[profile][param.Name] = param.Value
		}
	}

	cacheGroups, err := makeCGMap(cacheGroupArr)
	if err != nil {
		return nil, warnings, errors.New("making cachegroup map: " + err.Error())
	}

	nameTopologies := makeTopologyNameMap(topologies)
	anyCastPartners := GetAnyCastPartners(server, servers)

	sort.Sort(dsesSortByName(dses))

	sslDatas := []SSLData{}

	for _, ds := range dses {
		hasDS, err := dsUsesServer(&ds, server, dss, nameTopologies, cacheGroups, serverCapabilities, dsRequiredCapabilities)
		if err != nil {
			warnings = append(warnings, "error checking if ds uses this server, considering false! Error: "+err.Error())
			continue
		}
		if !hasDS {
			continue
		}

		dsParentConfigParams := map[string]string{}
		if ds.ProfileName != nil {
			dsParentConfigParams = profileParentConfigParams[*ds.ProfileName]
		}

		requestFQDNs, err := GetDSRequestFQDNs(&ds, dsRegexes[tc.DeliveryServiceName(ds.XMLID)], server, anyCastPartners, cdn.DomainName)
		if err != nil {
			warnings = append(warnings, "error getting ds '"+ds.XMLID+"' request fqdns, skipping! Error: "+err.Error())
			continue
		}

		enableH2 := defaultEnableH2
		tlsVersions := defaultTLSVersions

		dsTLSVersions := []TLSVersion{}
		for _, tlsVersion := range ds.TLSVersions {
			if _, ok := tlsVersionsToATS[TLSVersion(tlsVersion)]; !ok {
				warnings = append(warnings, "ds '"+ds.XMLID+"' had unknown TLS Version '"+tlsVersion+"' - ignoring!")
				continue
			}
			dsTLSVersions = append(dsTLSVersions, TLSVersion(tlsVersion))
		}
		if len(dsTLSVersions) > 0 {
			tlsVersions = dsTLSVersions
		}

		paramValEnableH2 := dsParentConfigParams[SSLServerNameYAMLParamEnableH2]
		paramValEnableH2 = strings.TrimSpace(paramValEnableH2)
		paramValEnableH2 = strings.ToLower(paramValEnableH2)

		if paramValEnableH2 != "" {
			enableH2 = strings.HasPrefix(paramValEnableH2, "t") || strings.HasPrefix(paramValEnableH2, "y")
		}

		paramValTLSVersions := dsParentConfigParams[SSLServerNameYAMLParamTLSVersions]
		paramValTLSVersions = strings.Replace(paramValTLSVersions, " ", "", -1)
		paramValTLSVersions = strings.TrimSpace(paramValTLSVersions)

		paramTLSVersions := []TLSVersion{}
		if paramValTLSVersions != "" {
			// Allow delimiting with commas, semicolons, spaces, or newlines.
			delim := ","
			if !strings.Contains(paramValTLSVersions, delim) {
				delim = ";"
			}
			if !strings.Contains(paramValTLSVersions, delim) {
				delim = " "
			}
			if !strings.Contains(paramValTLSVersions, delim) {
				delim = "\n"
			}

			tlsVersionsParamArr := strings.Split(paramValTLSVersions, delim)
			for _, tlsVersion := range tlsVersionsParamArr {
				if _, ok := tlsVersionsToATS[TLSVersion(tlsVersion)]; !ok {
					warnings = append(warnings, "ds '"+ds.XMLID+"' had unknown "+SSLServerNameYAMLParamTLSVersions+" parameter '"+tlsVersion+"' - ignoring!")
					continue
				}
				paramTLSVersions = append(paramTLSVersions, TLSVersion(tlsVersion))
			}
		}

		// let Parameters override the Delivery Service field, for backward-compatibility,
		// and also because this lets tenants who own multiple DSes set them all in a single
		// place, instead of duplicating for every DS (which will be even more de-duplicated
		// when Layered Profiles are implemented)
		if len(paramTLSVersions) != 0 {
			tlsVersions = paramTLSVersions
		}

		sslDatas = append(sslDatas, SSLData{
			DSName:       ds.XMLID,
			RequestFQDNs: requestFQDNs,
			EnableH2:     enableH2,
			TLSVersions:  tlsVersions,
		})
	}

	return sslDatas, warnings, nil
}

func dsUsesServer(
	ds *DeliveryService,
	server *Server,
	dss []DeliveryServiceServer,
	nameTopologies map[TopologyName]tc.TopologyV5,
	cacheGroups map[tc.CacheGroupName]tc.CacheGroupNullableV5,
	serverCapabilities map[int]map[ServerCapability]struct{},
	dsRequiredCapabilities map[int]map[ServerCapability]struct{},
) (bool, error) {
	if ds.XMLID == "" {
		return false, errors.New("ds missing xmlId")
	} else if ds.ID == nil {
		return false, errors.New("ds missing id")
	} else if server.ID == 0 {
		return false, errors.New("server missing id")
	} else if ds.Type == nil {
		return false, errors.New("ds missing type")
	}

	if !hasRequiredCapabilities(serverCapabilities[server.ID], dsRequiredCapabilities[*ds.ID]) {
		return false, nil
	}

	serverParentCGData, err := getParentCacheGroupData(server, cacheGroups)
	if err != nil {
		return false, errors.New("getting server parent cachegroup data: " + err.Error())
	}
	cacheIsTopLevel := isTopLevelCache(serverParentCGData)

	if !cacheIsTopLevel && (ds.Topology == nil || *ds.Topology == "") {
		if !dsAssignedServer(*ds.ID, server.ID, dss) {
			return false, nil
		}
	}

	if ds.Topology != nil && *ds.Topology != "" {
		topology, ok := nameTopologies[TopologyName(*ds.Topology)]
		if !ok {
			return false, errors.New("ds topology '" + *ds.Topology + "' not found in topologies")
		}

		serverPlacement, err := getTopologyPlacement(tc.CacheGroupName(server.CacheGroup), topology, cacheGroups, ds)
		if err != nil {
			return false, errors.New("getting topology placement: " + err.Error())
		}
		if !serverPlacement.InTopology {
			return false, nil
		}
	}

	return true, nil
}

// dsAssignedServer returns whether the Delivery Service Servers has an assignment between the server and the DS.
// Does not check Topologies, or parentage. Only useful for Edges and pre-topology DSS.
func dsAssignedServer(dsID int, serverID int, dsses []DeliveryServiceServer) bool {
	for _, dss := range dsses {
		if dss.Server == serverID && dss.DeliveryService == dsID {
			return true
		}
	}
	return false
}
