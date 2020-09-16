package cfgfile

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
	"strings"

	"github.com/apache/trafficcontrol/lib/go-atscfg"
	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops_ort/atstccfg/config"
)

func GetConfigFileCDNHeaderRewrite(toData *config.TOData, fileName string) (string, string, string, error) {
	dsName := strings.TrimSuffix(strings.TrimPrefix(fileName, atscfg.HeaderRewritePrefix), atscfg.ConfigSuffix) // TODO verify prefix and suffix? Perl doesn't

	tcDS := tc.DeliveryServiceNullableV30{}
	for _, ds := range toData.DeliveryServices {
		if ds.XMLID == nil || *ds.XMLID != dsName {
			continue
		}
		tcDS = ds
		break
	}
	if tcDS.ID == nil {
		return "", "", "", errors.New("ds '" + dsName + "' not found")
	}

	if tcDS.CDNName == nil {
		return "", "", "", errors.New("ds '" + dsName + "' missing cdn")
	}

	cfgDS, err := atscfg.HeaderRewriteDSFromDS(&tcDS)
	if err != nil {
		return "", "", "", errors.New("converting ds to config ds: " + err.Error())
	}

	dsServers := FilterDSS(toData.DeliveryServiceServers, map[int]struct{}{cfgDS.ID: {}}, nil)

	dsServerIDs := map[int]struct{}{}
	for _, dss := range dsServers {
		if dss.Server == nil || dss.DeliveryService == nil {
			continue // TODO warn?
		}
		if *dss.DeliveryService != *tcDS.ID {
			continue
		}
		dsServerIDs[*dss.Server] = struct{}{}
	}

	assignedEdges := []atscfg.HeaderRewriteServer{}
	for _, server := range toData.Servers {
		if server.CDNName == nil {
			log.Errorln("TO returned Servers server with missing CDNName, skipping!")
			continue
		}
		if server.ID == nil {
			log.Errorln("TO returned Servers server with missing ID, skipping!")
			continue
		}
		if *server.CDNName != *tcDS.CDNName {
			continue
		}
		if _, ok := dsServerIDs[*server.ID]; !ok && tcDS.Topology == nil {
			continue
		}
		cfgServer, err := atscfg.HeaderRewriteServerFromServer(server)
		if err != nil {
			continue // TODO warn?
		}
		assignedEdges = append(assignedEdges, cfgServer)
	}

	if toData.Server.CDNName == nil {
		return "", "", "", errors.New("this server missing CDNName")
	}

	return atscfg.MakeHeaderRewriteDotConfig(tc.CDNName(*toData.Server.CDNName), toData.TOToolName, toData.TOURL, cfgDS, assignedEdges, dsName), atscfg.ContentTypeHeaderRewriteDotConfig, atscfg.LineCommentHeaderRewriteDotConfig, nil
}
