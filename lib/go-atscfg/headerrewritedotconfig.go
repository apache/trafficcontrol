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
	"math"
	"regexp"
	"strconv"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
)

const HeaderRewritePrefix = "hdr_rw_"

const MaxOriginConnectionsNoMax = 0 // 0 indicates no limit on origin connections

type HeaderRewriteDS struct {
	EdgeHeaderRewrite    string
	ID                   int
	MaxOriginConnections int
	MidHeaderRewrite     string
	Type                 tc.DSType
}

type HeaderRewriteServer struct {
	HostName   string
	DomainName string
	Port       int
	Status     tc.CacheStatus
}

func HeaderRewriteServersFromServers(servers []tc.ServerNullable) ([]HeaderRewriteServer, error) {
	hServers := []HeaderRewriteServer{}
	for _, sv := range servers {
		hsv, err := HeaderRewriteServerFromServer(sv)
		if err != nil {
			return nil, err
		}
		hServers = append(hServers, hsv)
	}
	return hServers, nil
}

func HeaderRewriteServerFromServer(sv tc.ServerNullable) (HeaderRewriteServer, error) {
	if sv.HostName == nil {
		return HeaderRewriteServer{}, errors.New("server host name must not be nil")
	}
	if sv.DomainName == nil {
		return HeaderRewriteServer{}, errors.New("server domain name must not be nil")
	}
	if sv.TCPPort == nil {
		return HeaderRewriteServer{}, errors.New("server port must not be nil")
	}
	if sv.Status == nil {
		return HeaderRewriteServer{}, errors.New("server status must not be nil")
	}
	status := tc.CacheStatusFromString(*sv.Status)
	if status == tc.CacheStatusInvalid {
		return HeaderRewriteServer{}, errors.New("server status '" + *sv.Status + "' invalid")
	}
	return HeaderRewriteServer{Status: status, HostName: *sv.HostName, DomainName: *sv.DomainName, Port: *sv.TCPPort}, nil
}

func HeaderRewriteServerFromServerNotNullable(sv tc.Server) (HeaderRewriteServer, error) {
	if sv.HostName == "" {
		return HeaderRewriteServer{}, errors.New("server host name must not be nil")
	}
	if sv.DomainName == "" {
		return HeaderRewriteServer{}, errors.New("server domain name must not be nil")
	}
	if sv.TCPPort == 0 {
		return HeaderRewriteServer{}, errors.New("server port must not be nil")
	}
	status := tc.CacheStatusFromString(sv.Status)
	if status == tc.CacheStatusInvalid {
		return HeaderRewriteServer{}, errors.New("server status '" + sv.Status + "' invalid")
	}
	return HeaderRewriteServer{Status: status, HostName: sv.HostName, DomainName: sv.DomainName, Port: sv.TCPPort}, nil
}

func HeaderRewriteDSFromDS(ds *tc.DeliveryServiceNullable) (HeaderRewriteDS, error) {
	errs := []error{}
	if ds.ID == nil {
		errs = append(errs, errors.New("ID cannot be nil"))
	}
	if ds.Type == nil {
		errs = append(errs, errors.New("Type cannot be nil"))
	}
	if len(errs) > 0 {
		return HeaderRewriteDS{}, util.JoinErrs(errs)
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

	return HeaderRewriteDS{
		EdgeHeaderRewrite:    *ds.EdgeHeaderRewrite,
		ID:                   *ds.ID,
		MaxOriginConnections: *ds.MaxOriginConnections,
		MidHeaderRewrite:     *ds.MidHeaderRewrite,
		Type:                 *ds.Type,
	}, nil
}

func MakeHeaderRewriteDotConfig(
	cdnName tc.CDNName,
	toToolName string, // tm.toolname global parameter (TODO: cache itself?)
	toURL string, // tm.url global parameter (TODO: cache itself?)
	ds HeaderRewriteDS,
	assignedEdges []HeaderRewriteServer, // the edges assigned to ds
) string {
	text := GenericHeaderComment(string(cdnName), toToolName, toURL)

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

	text += "\n"
	return text
}
