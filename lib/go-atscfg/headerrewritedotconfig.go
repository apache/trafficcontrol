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
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
)

const HeaderRewritePrefix = "hdr_rw_"
const ContentTypeHeaderRewriteDotConfig = ContentTypeTextASCII
const LineCommentHeaderRewriteDotConfig = LineCommentHash

const ServiceCategoryHeader = "CDN-SVC"

const MaxOriginConnectionsNoMax = 0 // 0 indicates no limit on origin connections

func MakeHeaderRewriteDotConfig(
	fileName string,
	deliveryServices []DeliveryService,
	deliveryServiceServers []tc.DeliveryServiceServer,
	server *Server,
	servers []Server,
	hdrComment string,
) (Cfg, error) {
	warnings := []string{}

	dsName := strings.TrimSuffix(strings.TrimPrefix(fileName, HeaderRewritePrefix), ConfigSuffix) // TODO verify prefix and suffix? Perl doesn't

	tcDS := DeliveryService{}
	for _, ds := range deliveryServices {
		if ds.XMLID == nil {
			warnings = append(warnings, "deliveryServices had DS with nil xmlId (name)")
			continue
		}
		if *ds.XMLID != dsName {
			continue
		}
		tcDS = ds
		break
	}
	if tcDS.ID == nil {
		return Cfg{}, makeErr(warnings, "ds '"+dsName+"' not found")
	}

	if tcDS.CDNName == nil {
		return Cfg{}, makeErr(warnings, "ds '"+dsName+"' missing cdn")
	}

	ds, err := headerRewriteDSFromDS(&tcDS)
	if err != nil {
		return Cfg{}, makeErr(warnings, "converting ds to config ds: "+err.Error())
	}

	dsServers := filterDSS(deliveryServiceServers, map[int]struct{}{ds.ID: {}}, nil)

	dsServerIDs := map[int]struct{}{}
	for _, dss := range dsServers {
		if dss.Server == nil || dss.DeliveryService == nil {
			warnings = append(warnings, "deliveryservice-servers had entry with nil values, skipping!")
			continue
		}
		if *dss.DeliveryService != *tcDS.ID {
			continue
		}
		dsServerIDs[*dss.Server] = struct{}{}
	}

	assignedEdges := []headerRewriteServer{}
	for _, server := range servers {
		if server.CDNName == nil {
			warnings = append(warnings, "servers had server with missing cdnName, skipping!")
			continue
		}
		if server.ID == nil {
			warnings = append(warnings, "servers had server with missing id, skipping!")
			continue
		}
		if *server.CDNName != *tcDS.CDNName {
			continue
		}
		if _, ok := dsServerIDs[*server.ID]; !ok && tcDS.Topology == nil {
			continue
		}
		cfgServer, err := headerRewriteServerFromServer(server)
		if err != nil {
			warnings = append(warnings, "error getting header rewrite server, skipping: "+err.Error())
			continue
		}
		assignedEdges = append(assignedEdges, cfgServer)
	}

	if server.CDNName == nil {
		return Cfg{}, makeErr(warnings, "this server missing CDNName")
	}

	text := makeHdrComment(hdrComment)

	// write a header rewrite rule if maxOriginConnections > 0 and the ds does NOT use mids
	if ds.MaxOriginConnections > 0 && !ds.Type.UsesMidCache() {
		dsOnlineEdgeCount := 0
		for _, sv := range assignedEdges {
			if sv.Status == tc.CacheStatusReported || sv.Status == tc.CacheStatusOnline {
				dsOnlineEdgeCount++
			}
		}

		if dsOnlineEdgeCount > 0 {
			maxOriginConnectionsPerEdge := int(math.Round(float64(ds.MaxOriginConnections) / float64(dsOnlineEdgeCount)))
			text += "cond %{REMAP_PSEUDO_HOOK}\nset-config proxy.config.http.origin_max_connections " + strconv.Itoa(maxOriginConnectionsPerEdge)
			if ds.EdgeHeaderRewrite == "" {
				text += " [L]"
			} else {
				text += "\n"
			}
		}
	}

	// write the contents of ds.EdgeHeaderRewrite to hdr_rw_xml-id.config replacing any instances of __RETURN__ (surrounded by spaces or not) with \n
	if ds.EdgeHeaderRewrite != "" {
		re := regexp.MustCompile(`\s*__RETURN__\s*`)
		text += re.ReplaceAllString(ds.EdgeHeaderRewrite, "\n")
	}

	if !strings.Contains(text, ServiceCategoryHeader) && ds.ServiceCategory != "" {
		scHeaderVal := fmt.Sprintf("\nset-header %s \"%s|%s\" %s", ServiceCategoryHeader, dsName, ds.ServiceCategory, "[L]")
		if strings.Contains(text, "[L]") {
			text = strings.Replace(text, "[L]", scHeaderVal, 1)
		} else {
			text += scHeaderVal
		}
	}

	text += "\n"

	return Cfg{
		Text:        text,
		ContentType: ContentTypeHeaderRewriteDotConfig,
		LineComment: LineCommentHeaderRewriteDotConfig,
		Warnings:    warnings,
	}, nil
}

type headerRewriteDS struct {
	EdgeHeaderRewrite    string
	ID                   int
	MaxOriginConnections int
	MidHeaderRewrite     string
	Type                 tc.DSType
	ServiceCategory      string
}

type headerRewriteServer struct {
	HostName   string
	DomainName string
	Port       int
	Status     tc.CacheStatus
}

func headerRewriteServersFromServers(servers []Server) ([]headerRewriteServer, error) {
	hServers := []headerRewriteServer{}
	for _, sv := range servers {
		hsv, err := headerRewriteServerFromServer(sv)
		if err != nil {
			return nil, err
		}
		hServers = append(hServers, hsv)
	}
	return hServers, nil
}

func headerRewriteServerFromServer(sv Server) (headerRewriteServer, error) {
	if sv.HostName == nil {
		return headerRewriteServer{}, errors.New("server host name must not be nil")
	}
	if sv.DomainName == nil {
		return headerRewriteServer{}, errors.New("server domain name must not be nil")
	}
	if sv.TCPPort == nil {
		return headerRewriteServer{}, errors.New("server port must not be nil")
	}
	if sv.Status == nil {
		return headerRewriteServer{}, errors.New("server status must not be nil")
	}
	status := tc.CacheStatusFromString(*sv.Status)
	if status == tc.CacheStatusInvalid {
		return headerRewriteServer{}, errors.New("server status '" + *sv.Status + "' invalid")
	}
	return headerRewriteServer{Status: status, HostName: *sv.HostName, DomainName: *sv.DomainName, Port: *sv.TCPPort}, nil
}

func headerRewriteServerFromServerNotNullable(sv tc.Server) (headerRewriteServer, error) {
	if sv.HostName == "" {
		return headerRewriteServer{}, errors.New("server host name must not be nil")
	}
	if sv.DomainName == "" {
		return headerRewriteServer{}, errors.New("server domain name must not be nil")
	}
	if sv.TCPPort == 0 {
		return headerRewriteServer{}, errors.New("server port must not be nil")
	}
	status := tc.CacheStatusFromString(sv.Status)
	if status == tc.CacheStatusInvalid {
		return headerRewriteServer{}, errors.New("server status '" + sv.Status + "' invalid")
	}
	return headerRewriteServer{Status: status, HostName: sv.HostName, DomainName: sv.DomainName, Port: sv.TCPPort}, nil
}

func headerRewriteDSFromDS(ds *DeliveryService) (headerRewriteDS, error) {
	errs := []error{}
	if ds.ID == nil {
		errs = append(errs, errors.New("ID cannot be nil"))
	}
	if ds.Type == nil {
		errs = append(errs, errors.New("Type cannot be nil"))
	}
	if len(errs) > 0 {
		return headerRewriteDS{}, util.JoinErrs(errs)
	}

	if ds.MaxOriginConnections == nil {
		ds.MaxOriginConnections = util.IntPtr(MaxOriginConnectionsNoMax)
	}
	if ds.EdgeHeaderRewrite == nil {
		ds.EdgeHeaderRewrite = util.StrPtr("")
	}
	if ds.MidHeaderRewrite == nil {
		ds.MidHeaderRewrite = util.StrPtr("")
	}
	if ds.ServiceCategory == nil {
		ds.ServiceCategory = new(string)
	}

	return headerRewriteDS{
		EdgeHeaderRewrite:    *ds.EdgeHeaderRewrite,
		ID:                   *ds.ID,
		MaxOriginConnections: *ds.MaxOriginConnections,
		MidHeaderRewrite:     *ds.MidHeaderRewrite,
		Type:                 *ds.Type,
		ServiceCategory:      *ds.ServiceCategory,
	}, nil
}
