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
	// "fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-rfc"
	"github.com/apache/trafficcontrol/lib/go-tc"
)

const ContentTypeSSLMultiCertDotConfig = ContentTypeTextASCII
const LineCommentSSLMultiCertDotConfig = LineCommentHash
const SSLMultiCertConfigFileName = `ssl_multicert.config`

// SSLMultiCertDotConfigOpts contains settings to configure generation options.
type SSLMultiCertDotConfigOpts struct {
	// HdrComment is the header comment to include at the beginning of the file.
	// This should be the text desired, without comment syntax (like # or //). The file's comment syntax will be added.
	// To omit the header comment, pass the empty string.
	HdrComment string

	// Certificates is a list of additional certificates to manually enter into ssl_multicert.config.
	// These certificates must already exist, and not be managed by Traffic Ops.
	//
	// The most common use for this is the client and server certificates for
	// intra-cdn child-parent communication a.k.a. "end-to-end ssl".
	Certificates []SSLMultiCertDotConfigCertInf

	// InternalHTTPS is whether to generate rules for internal https communication.
	// If omitted, the default is 'no'
	InternalHTTPS InternalHTTPS

	// E2ESSLCAPath is the file name and path to the Certificate Authority to use
	// for End-to-End SSL Client and Server Certificates.
	// If empty, no specific CA will be inserted, and ATS will use the primary CA bundle.
	E2ESSLCAPath string

	// E2ESSLServerKeyPath is the file name and path to the key to use
	// for End-to-End SSL Server Certificates.
	// This must exist if InternalHTTPS is true or no-child.
	E2ESSLServerKeyPath string
}

type E2ECertMetaData struct {
	DSName   tc.DeliveryServiceName `json:"delivery_service_name"`
	URI      *url.URL               `json:"uri"` // TODO marshal/unmarshal as string? built-in type names are awkward
	Internal bool                   `json:"internal"`
	Type     RemapMapType           `json:"type"`
	CertPath string                 `json:"cert_path"`
	KeyPath  string                 `json:"key_path"`
	CAPath   string                 `json:"ca_path"`
}

type RemapMapType string

const RemapMapTypeSource = RemapMapType("source")
const RemapMapTypeTarget = RemapMapType("target")

// SSLMultiCertDotConfigOpts contains metadata needed by config generation,
// in addition to the file itself.
type SSLMultiCertDotConfigMetaData struct {
	E2ECerts []E2ECertMetaData `json:"e2e_certs"`
}

// IsCfgMetaData implements CfgMetaData.
func (md SSLMultiCertDotConfigMetaData) IsCfgMetaData() {}

// SSLMultiCertDotConfigCertInf is the information for a certificate in the Apache Traffic Server ssl_multicert.config config file.
// The paths are relative to the ATS records.config proxy.config.ssl.server.cert.path, and certificates should generally be put there (typically etc/trafficserver/ssl/), and their relative filenames used here.
// The CAPath is the path to the Certificate Authority file. This is optional, and only necessary if the certificates are signed by a CA not in the system CA bundle.
type SSLMultiCertDotConfigCertInf struct {
	CertPath string
	KeyPath  string
	CAPath   string
}

func MakeSSLMultiCertDotConfig(
	server *Server,
	deliveryServices []DeliveryService,
	dss []DeliveryServiceServer,
	dsRegexArr []tc.DeliveryServiceRegexes,
	cdn *tc.CDN,
	topologies []tc.Topology,
	cacheGroupArr []tc.CacheGroupNullable,
	serverCapabilities map[int]map[ServerCapability]struct{},
	dsRequiredCapabilities map[int]map[ServerCapability]struct{},
	opt *SSLMultiCertDotConfigOpts,
) (Cfg, error) {
	if opt == nil {
		opt = &SSLMultiCertDotConfigOpts{}
	}
	warnings := []string{}
	if server.CDNName == nil {
		return Cfg{}, makeErr(warnings, "server missing CDNName")
	}

	warnings = append(warnings, "DEBUG MakeSSLMultiCertDotConfig  opt.E2ESSLServerKeyPath '"+opt.E2ESSLServerKeyPath+"'")

	dses, dsWarns := DeliveryServicesToSSLMultiCertDSes(deliveryServices)
	warnings = append(warnings, dsWarns...)

	hdr := makeHdrComment(opt.HdrComment)

	dses = GetSSLMultiCertDotConfigDeliveryServices(dses)

	lines := []string{}

	// TODO fix so servers get DS certs if they're accepting client connections
	//      for that DS, per the DS and/or its Topology.
	//      Currently, "EDGE" types accept clients, and "MID" types accept other
	//      caches. That needs to change, Topologies shouldn't care about Types.
	//      But as of this writing, Topologies don't have that data yet, so it
	//      isn't possible in cache config until TO has that data.
	if tc.CacheType(server.Type) == tc.CacheTypeEdge {
		for dsName, ds := range dses {
			cerName, keyName := GetSSLMultiCertDotConfigCertAndKeyName(dsName, ds)
			lines = append(lines, `ssl_cert_name=`+cerName+"\t"+` ssl_key_name=`+keyName+"\n")
		}
	}

	e2eMetaData, e2eWarns, err := GetE2ESSLCertInf(
		server,
		deliveryServices,
		dss,
		dsRegexArr,
		cdn,
		topologies,
		cacheGroupArr,
		serverCapabilities,
		dsRequiredCapabilities,
		opt.InternalHTTPS,
		opt.E2ESSLCAPath,
		opt.E2ESSLServerKeyPath,
	)
	warnings = append(warnings, e2eWarns...)
	if err != nil {
		return Cfg{}, makeErr(warnings, "getting e2e cert data: "+err.Error())
	}

	// warnings = append(warnings, fmt.Sprintf("DEBUG GetE2ESSLCertInf returned len %+v", e2eMetaData))

	certs := []SSLMultiCertDotConfigCertInf{}

	for _, md := range e2eMetaData {
		if !md.Internal {
			continue // only internal routes need E2E certs
		}
		if md.Type != RemapMapTypeSource {
			continue // only Sources need E2E certs per-DS; targets use a single client cert
		}
		if md.URI.Scheme != rfc.SchemeHTTPS {
			continue // https remaps don't have certs
		}
		certs = append(certs, SSLMultiCertDotConfigCertInf{
			CertPath: md.CertPath,
			KeyPath:  md.KeyPath,
			CAPath:   md.CAPath,
		})
	}

	certs = append(certs, opt.Certificates...)

	for _, certInf := range certs {
		// TODO check that files exist, and error if they don't? Ideally that they're valid certs and keys?
		line := `ssl_cert_name=` + certInf.CertPath + "\t" + ` ssl_key_name=` + certInf.KeyPath
		if strings.TrimSpace(certInf.CAPath) != "" {
			line += "\t" + ` ssl_ca_name=` + certInf.CAPath
		}
		line += "\n"
		lines = append(lines, line)
	}

	sort.Strings(lines)

	txt := hdr + strings.Join(lines, "")

	return Cfg{
		Text:        txt,
		ContentType: ContentTypeSSLMultiCertDotConfig,
		LineComment: LineCommentSSLMultiCertDotConfig,
		Secure:      true,
		MetaData:    SSLMultiCertDotConfigMetaData{E2ECerts: e2eMetaData},
		Warnings:    warnings,
	}, nil
}

func GetE2ESSLCertInf(
	server *Server,
	unfilteredDSes []DeliveryService,
	dss []DeliveryServiceServer,
	dsRegexArr []tc.DeliveryServiceRegexes,
	cdn *tc.CDN,
	topologies []tc.Topology,
	cacheGroupArr []tc.CacheGroupNullable,
	serverCapabilities map[int]map[ServerCapability]struct{},
	dsRequiredCapabilities map[int]map[ServerCapability]struct{},
	internalHTTPS InternalHTTPS,
	e2eSSLCAPath string,
	e2eSSLServerKeyPath string,
) ([]E2ECertMetaData, []string, error) {
	warnings := []string{}

	//	warnings = append(warnings, "DEBUG GetE2ESSLCertInf calling")

	dsRemapURIDatas, dsWarns, err := GetSourceAndTargetURIs(
		server,
		unfilteredDSes,
		dss,
		dsRegexArr,
		cdn,
		topologies,
		cacheGroupArr,
		serverCapabilities,
		dsRequiredCapabilities,
		internalHTTPS,
	)
	warnings = append(warnings, dsWarns...)
	if err != nil {
		return nil, warnings, errors.New("getting internal source and target URIs: " + err.Error())
	}

	// warnings = append(warnings, fmt.Sprintf("DEBUG GetE2ESSLCertInf got datas len %+v", dsRemapURIDatas))

	certInf := []E2ECertMetaData{}

	for _, dsRemapURIData := range dsRemapURIDatas {
		for _, line := range dsRemapURIData.RemapLines {
			if e2eSSLServerKeyPath == "" {
				return nil, warnings, errors.New("ds '" + string(dsRemapURIData.DS) + "' requires E2E Server Certs, but opt server key path was empty")
			}
			certInf = append(certInf, E2ECertMetaData{
				DSName:   dsRemapURIData.DS,
				URI:      line.Source.URI,
				Type:     RemapMapTypeSource,
				Internal: line.Source.Internal,
				CertPath: GetDSE2ESSLCertFileName(dsRemapURIData.DS, line.Source.URI) + DSE2ESSLCertFileNameExtensionCert,
				KeyPath:  e2eSSLServerKeyPath,
				CAPath:   e2eSSLCAPath,
			})
			certInf = append(certInf, E2ECertMetaData{
				DSName:   dsRemapURIData.DS,
				URI:      line.Target.URI,
				Type:     RemapMapTypeTarget,
				Internal: line.Target.Internal,
				CertPath: GetDSE2ESSLCertFileName(dsRemapURIData.DS, line.Target.URI) + DSE2ESSLCertFileNameExtensionCert,
				KeyPath:  e2eSSLServerKeyPath,
				CAPath:   e2eSSLCAPath,
			})
		}
	}

	// warnings = append(warnings, fmt.Sprintf("DEBUG GetE2ESSLCertInf made inf len %+v", certInf))
	return certInf, warnings, nil
}

// GetDSE2ESSLCertFileName returns the file name to be used for the given internal URL
// for End-to-End SSL Certificates for internal HTTPS traffic.
// File name does not include the extension. See DSE2ESSLCertFileNameExtensionCert and
// DSE2ESSLCertFileNameExtensionKey.
func GetDSE2ESSLCertFileName(ds tc.DeliveryServiceName, uri *url.URL) string {
	return "e2e_" + string(ds) + "_" + uri.Hostname() + "_" + rfc.URLPortOrDefault(uri)
}

const DSE2ESSLCertFileNameExtensionCert = `.cert`
const DSE2ESSLCertFileNameExtensionKey = `.key`

// func GetInternalServerCertFileName(internalRemapURI string) (string, error) {
// }

type sslMultiCertDS struct {
	XMLID       string
	Type        tc.DSType
	Protocol    int
	ExampleURLs []string
}

// deliveryServicesToSSLMultiCertDSes returns the "SSLMultiCertDS" map, and any warnings.
func DeliveryServicesToSSLMultiCertDSes(dses []DeliveryService) (map[tc.DeliveryServiceName]sslMultiCertDS, []string) {
	warnings := []string{}
	sDSes := map[tc.DeliveryServiceName]sslMultiCertDS{}
	for _, ds := range dses {
		if ds.Type == nil || ds.Protocol == nil || ds.XMLID == nil {
			if ds.XMLID == nil {
				warnings = append(warnings, "got unknown DS with nil values! Skipping!")
			} else {
				warnings = append(warnings, "got DS '"+*ds.XMLID+"' with nil values! Skipping!")
			}
			continue
		}
		sDSes[tc.DeliveryServiceName(*ds.XMLID)] = sslMultiCertDS{Type: *ds.Type, Protocol: *ds.Protocol, ExampleURLs: ds.ExampleURLs}
	}
	return sDSes, warnings
}

// GetSSLMultiCertDotConfigCertAndKeyName returns the cert file name and key file name for the given delivery service.
func GetSSLMultiCertDotConfigCertAndKeyName(dsName tc.DeliveryServiceName, ds sslMultiCertDS) (string, string) {
	hostName := ds.ExampleURLs[0] // first one is the one we want

	scheme := "https://"
	if !strings.HasPrefix(hostName, scheme) {
		scheme = "http://"
	}
	newHost := hostName
	if len(hostName) < len(scheme) {
		log.Errorln("MakeSSLMultiCertDotConfig got ds '" + string(dsName) + "' example url '" + hostName + "' with no scheme! ssl_multicert.config will likely be malformed!")
	} else {
		newHost = hostName[len(scheme):]
	}
	keyName := newHost + ".key"

	newHost = strings.Replace(newHost, ".", "_", -1)

	cerName := newHost + "_cert.cer"
	return cerName, keyName
}

// GetSSLMultiCertDotConfigDeliveryServices takes a list of delivery services, and returns the delivery services which will be inserted into the config by MakeSSLMultiCertDotConfig.
// This is public, so users can see which Delivery Services are used, without parsing the config file.
// For example, this is useful to determine which certificates are needed.
func GetSSLMultiCertDotConfigDeliveryServices(dses map[tc.DeliveryServiceName]sslMultiCertDS) map[tc.DeliveryServiceName]sslMultiCertDS {
	usedDSes := map[tc.DeliveryServiceName]sslMultiCertDS{}
	for dsName, ds := range dses {
		if ds.Type == tc.DSTypeAnyMap {
			continue
		}
		if ds.Type.IsSteering() {
			continue // Steering delivery service SSLs should not be on the edges.
		}
		if ds.Protocol == 0 {
			continue
		}
		if len(ds.ExampleURLs) == 0 {
			continue // TODO warn? error? Perl doesn't
		}
		usedDSes[dsName] = ds
	}
	return usedDSes
}

// DSRemapURIData contains info about the given Delivery Service's remaps for a particular server.
type DSRemapURIData struct {
	DS         tc.DeliveryServiceName
	RemapLines []DSRemapLine
}

type DSRemapLine struct {
	Source RemapURIData
	Target RemapURIData
}

type RemapURIData struct {
	// URI is the URI of this remap source or target. This will typically only include the Scheme and Host.
	URI *url.URL
	// Internal is whether this URI is internal to the CDN or not.
	// For remap sources, false means incoming requests are clients, true means incoming requests are child caches.
	// For remap targets, false means outgoing requests are origins, true means outgoing requests are parent caches.
	Internal bool
}

// GetSourceAndTargetURIs returns, for the given server,
// the internal CDN source URIs (sources on this server's remap config
// which will be requested of this cache by other caches),
// targetURIs (targets on this server's remap config which this cache will request of other caches,
// any warnings, and any error.
func GetSourceAndTargetURIs(
	server *Server,
	unfilteredDSes []DeliveryService,
	dss []DeliveryServiceServer,
	dsRegexArr []tc.DeliveryServiceRegexes,
	cdn *tc.CDN,
	topologies []tc.Topology,
	cacheGroupArr []tc.CacheGroupNullable,
	serverCapabilities map[int]map[ServerCapability]struct{},
	dsRequiredCapabilities map[int]map[ServerCapability]struct{},
	internalHTTPS InternalHTTPS,
) ([]DSRemapURIData, []string, error) {
	// TODO fix duplication with remap.config.
	//      There's no obvious answer. Will probably require a complete refactor of
	//      remap.config, to use these funcs, modified to include all remaps not just
	//      internal ones.
	warnings := []string{}

	cdnDomain := cdn.DomainName
	dsRegexes := makeDSRegexMap(dsRegexArr)
	// Returned DSes are guaranteed to have a non-nil XMLID, Type, DSCP, ID, and Active.
	dses, dsWarns := remapFilterDSes(server, dss, unfilteredDSes)
	warnings = append(warnings, dsWarns...)

	cacheGroups, err := makeCGMap(cacheGroupArr)
	if err != nil {
		return nil, warnings, errors.New("making cachegroup map: " + err.Error())
	}

	nameTopologies := makeTopologyNameMap(topologies)

	if tc.CacheTypeFromString(server.Type) == tc.CacheTypeMid {
		return getSourceAndTargetURIsForMid(dses, dsRegexes, server, nameTopologies, cacheGroups, serverCapabilities, dsRequiredCapabilities, internalHTTPS)
	}
	return getSourceAndTargetURIsForEdge(dses, dsRegexes, server, nameTopologies, cacheGroups, serverCapabilities, dsRequiredCapabilities, cdnDomain, internalHTTPS)
}

func getSourceAndTargetURIsForMid(
	dses []DeliveryService,
	dsRegexes map[tc.DeliveryServiceName][]tc.DeliveryServiceRegex,
	server *Server,
	nameTopologies map[TopologyName]tc.Topology,
	cacheGroups map[tc.CacheGroupName]tc.CacheGroupNullable,
	serverCapabilities map[int]map[ServerCapability]struct{},
	dsRequiredCapabilities map[int]map[ServerCapability]struct{},
	internalHTTPS InternalHTTPS,
) ([]DSRemapURIData, []string, error) {

	warnings := []string{}
	remapURIDatas := []DSRemapURIData{}

	midRemaps := map[string]string{}
	for _, ds := range dses {
		dsURIData := DSRemapURIData{DS: tc.DeliveryServiceName(*ds.XMLID)}

		if !hasRequiredCapabilities(serverCapabilities[*server.ID], dsRequiredCapabilities[*ds.ID]) {
			continue
		}

		_, isChild, _, err := dsCacheIsParentChildEdge(&ds, server, nameTopologies, cacheGroups)
		if err != nil {
			return nil, warnings, errors.New("checking ds '" + *ds.XMLID + "' server parentage: " + err.Error())
		}

		topology, hasTopology := nameTopologies[TopologyName(*ds.Topology)]
		if *ds.Topology != "" && hasTopology {
			topoIncludesServer, err := topologyIncludesServerNullable(topology, server)
			if err != nil {
				return nil, warnings, errors.New("getting Topology Server inclusion: " + err.Error())
			}
			if !topoIncludesServer {
				continue
			}
		}

		if !ds.Type.UsesMidCache() && (!hasTopology || *ds.Topology == "") {
			continue // Live local delivery services skip mids (except Topologies ignore DS types)
		}

		if ds.OrgServerFQDN == nil || *ds.OrgServerFQDN == "" {
			warnings = append(warnings, "ds '"+*ds.XMLID+"' has no origin fqdn, skipping!") // TODO confirm - Perl uses without checking!
			continue
		}

		sourceURI, err := url.Parse(*ds.OrgServerFQDN)
		if err != nil {
			warnings = append(warnings, "ds '"+*ds.XMLID+"' remap source '"+*ds.OrgServerFQDN+"' not a valid URI, skipping! Parse error: "+err.Error())
			continue
		}

		targetURI, err := url.Parse(*ds.OrgServerFQDN) // may be the same as the source, but it needs to be a copy, because we're going to modify them
		if err != nil {
			warnings = append(warnings, "ds '"+*ds.XMLID+"' origin '"+*ds.OrgServerFQDN+"' not a valid URI, skipping! Parse error: "+err.Error())
			continue
		}

		sourceURI.Scheme = rfc.SchemeHTTP // set to http for now, we'll determine if we need http, https, or both below (may need multiple remaps)

		if midRemaps[sourceURI.String()] != "" {
			continue // skip remap rules from extra HOST_REGEXP entries
		}
		midRemaps[sourceURI.String()] = targetURI.String()

		if isChild {
			// if we're a child requesting of a parent, make the target scheme the internal HTTPS setting
			if internalHTTPS == InternalHTTPSNo || internalHTTPS == InternalHTTPSNoChild {
				targetURI.Scheme = rfc.SchemeHTTP
			} else {
				targetURI.Scheme = rfc.SchemeHTTPS
			}
		}

		remapTarget := RemapURIData{
			URI:      targetURI,
			Internal: isChild,
		}

		// Above, we made originURI always be http.
		// Now, we decide whether to insert http, https, or both, based on the InternalHTTPS setting
		if internalHTTPS == InternalHTTPSNo || internalHTTPS == InternalHTTPSNoChild {
			dsURIData.RemapLines = append(dsURIData.RemapLines, DSRemapLine{
				Source: RemapURIData{
					URI:      sourceURI,
					Internal: true, // TODO: MID sources are always children (and EDGE sources are always clients). That needs fixed, but will require a much larger remap.config refactor.
				},
				Target: remapTarget,
			})
		}
		if internalHTTPS == InternalHTTPSYes || internalHTTPS == InternalHTTPSNoChild {
			sourceURIHTTPS := *sourceURI // need to copy, so we don't modify the http uri
			sourceURIHTTPS.Scheme = rfc.SchemeHTTPS
			dsURIData.RemapLines = append(dsURIData.RemapLines, DSRemapLine{
				Source: RemapURIData{
					URI:      &sourceURIHTTPS,
					Internal: true, // TODO: MID sources are always children (and EDGE sources are always clients). That needs fixed, but will require a much larger remap.config refactor.
				},
				Target: remapTarget,
			})
		}
		remapURIDatas = append(remapURIDatas, dsURIData)
	}

	return remapURIDatas, warnings, nil
}

func getSourceAndTargetURIsForEdge(
	dses []DeliveryService,
	dsRegexes map[tc.DeliveryServiceName][]tc.DeliveryServiceRegex,
	server *Server,
	nameTopologies map[TopologyName]tc.Topology,
	cacheGroups map[tc.CacheGroupName]tc.CacheGroupNullable,
	serverCapabilities map[int]map[ServerCapability]struct{},
	dsRequiredCapabilities map[int]map[ServerCapability]struct{},
	cdnDomain string,
	internalHTTPS InternalHTTPS,
) ([]DSRemapURIData, []string, error) {
	warnings := []string{}
	remapURIDatas := []DSRemapURIData{}

	for _, ds := range dses {
		if !hasRequiredCapabilities(serverCapabilities[*server.ID], dsRequiredCapabilities[*ds.ID]) {
			continue
		}

		topology, hasTopology := nameTopologies[TopologyName(*ds.Topology)]
		if *ds.Topology != "" && hasTopology {
			topoIncludesServer, err := topologyIncludesServerNullable(topology, server)
			if err != nil {
				return nil, warnings, errors.New("getting topology server inclusion: " + err.Error())
			}
			if !topoIncludesServer {
				continue
			}
		}

		requestFQDNs, err := getDSRequestFQDNs(&ds, dsRegexes[tc.DeliveryServiceName(*ds.XMLID)], server, cdnDomain)
		if err != nil {
			warnings = append(warnings, "error getting ds '"+*ds.XMLID+"' request fqdns, skipping! Error: "+err.Error())
			continue
		}

		dsURIData := DSRemapURIData{DS: tc.DeliveryServiceName(*ds.XMLID)}

		for _, requestFQDN := range requestFQDNs {
			remapLines, err := makeEdgeDSDataRemapLines(ds, requestFQDN, server, cdnDomain)
			if err != nil {
				warnings = append(warnings, "DS '"+*ds.XMLID+"' - skipping! : "+err.Error())
				continue
			}

			for _, line := range remapLines {
				dsRemapLines, dsWarns, err := getSourceAndTargetURIsForEdgeForDS(server, &ds, line.From, line.To, cacheGroups, nameTopologies, internalHTTPS)
				warnings = append(warnings, dsWarns...)
				if err != nil {
					return nil, warnings, errors.New("building internal remap uris for ds '" + *ds.XMLID + "': " + err.Error())
				}
				dsURIData.RemapLines = append(dsURIData.RemapLines, dsRemapLines...)
			}
		}
		remapURIDatas = append(remapURIDatas, dsURIData)
	}
	return remapURIDatas, warnings, nil
}

func getSourceAndTargetURIsForEdgeForDS(
	server *Server,
	ds *DeliveryService,
	mapFrom string,
	mapTo string,
	cacheGroups map[tc.CacheGroupName]tc.CacheGroupNullable,
	nameTopologies map[TopologyName]tc.Topology,
	internalHTTPS InternalHTTPS,
) ([]DSRemapLine, []string, error) {
	warnings := []string{}

	isLastCache, err := serverIsLastCacheForDS(server, ds, nameTopologies, cacheGroups)
	if err != nil {
		return nil, warnings, errors.New("determining if cache is the last tier: " + err.Error())
	}

	// must replace before parsing, because __http__ might do odd things to the parse (underscores aren't valid in FQDNs)
	mapFrom = strings.Replace(mapFrom, `__http__`, *server.HostName, -1)

	sourceURI, err := url.Parse(mapFrom)
	if err != nil {
		return nil, warnings, errors.New("ds '" + *ds.XMLID + "' edge remap source '" + mapFrom + "' URI parse error: " + err.Error())
	}

	targetURI, err := url.Parse(mapTo)
	if err != nil {
		return nil, warnings, errors.New("ds '" + *ds.XMLID + "' edge remap target '" + mapTo + "' URI parse error: " + err.Error())
	}

	// if this remap is going to a parent, use the internal protocol, not the origin protocol
	if !isLastCache {
		if internalHTTPS == InternalHTTPSNo || internalHTTPS == InternalHTTPSNoChild {
			targetURI.Scheme = rfc.SchemeHTTP
		} else {
			targetURI.Scheme = rfc.SchemeHTTPS
		}
	}

	return []DSRemapLine{
		{
			Source: RemapURIData{
				URI:      sourceURI,
				Internal: false, // TODO: EDGE sources are always clients (and MID sources are always child caches). That needs fixed, but will require a much larger remap.config refactor.
			},
			Target: RemapURIData{
				URI:      targetURI,
				Internal: isLastCache,
			},
		},
	}, warnings, nil
}
