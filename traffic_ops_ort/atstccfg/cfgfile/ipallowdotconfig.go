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
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops_ort/atstccfg/config"
)

func GetConfigFileServerIPAllowDotConfig(toData *config.TOData) (string, string, string, error) {
	fileParams := ParamsToMultiMap(FilterParams(toData.ServerParams, atscfg.IPAllowConfigFileName, "", "", ""))

	cgMap := map[string]tc.CacheGroupNullable{}
	for _, cg := range toData.CacheGroups {
		if cg.Name == nil {
			return "", "", "", errors.New("got cachegroup with nil name!'")
		}
		cgMap[*cg.Name] = cg
	}

	serverCG, ok := cgMap[toData.Server.Cachegroup]
	if !ok {
		return "", "", "", errors.New("server cachegroup not in cachegroups!")
	}

	childCGs := map[string]tc.CacheGroupNullable{}
	for cgName, cg := range cgMap {
		if (cg.ParentName != nil && *cg.ParentName == *serverCG.Name) || (cg.SecondaryParentName != nil && *cg.SecondaryParentName == *serverCG.Name) {
			childCGs[cgName] = cg
		}
	}

	childServers := map[tc.CacheName]atscfg.IPAllowServer{}
	for _, sv := range toData.Servers {
		_, ok := childCGs[sv.Cachegroup]
		if ok || (strings.HasPrefix(toData.Server.Type, tc.MidTypePrefix) && string(sv.Type) == tc.MonitorTypeName) {
			childServers[tc.CacheName(sv.HostName)] = atscfg.IPAllowServer{IPAddress: sv.IPAddress, IP6Address: sv.IP6Address}
		}
	}

	return atscfg.MakeIPAllowDotConfig(tc.CacheName(toData.Server.HostName), tc.CacheType(toData.Server.Type), toData.TOToolName, toData.TOURL, fileParams, childServers), atscfg.ContentTypeIPAllowDotConfig, atscfg.LineCommentIPAllowDotConfig, nil
}
