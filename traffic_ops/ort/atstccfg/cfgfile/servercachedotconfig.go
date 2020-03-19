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
	"github.com/apache/trafficcontrol/traffic_ops/ort/atstccfg/config"
)

const ServerCacheDotConfigIncludeInactiveDSes = false // TODO move to lib/go-atscfg

func GetConfigFileServerCacheDotConfig(toData *config.TOData) (string, string, error) {
	// TODO TOAPI add /servers?cdn=1 query param

	// TODO remove this, we generated the scope, we know it's right? Or should we have an extra safety check?
	if !strings.HasPrefix(string(toData.Server.Type), tc.MidTypePrefix) {
		// emulates Perl
		return "", "", errors.New("Error - incorrect file scope for route used.  Please use the profiles route.")
	}

	dsData := map[tc.DeliveryServiceName]atscfg.ServerCacheConfigDS{}
	for _, ds := range toData.DeliveryServices {
		if ds.XMLID == nil || ds.Active == nil || ds.OrgServerFQDN == nil || ds.Type == nil {
			// TODO orgserverfqdn is nil for some DSes - MSO? Verify.
			continue
			//			return "", fmt.Errorf("getting delivery services: got DS with nil values! '%v' %v %+v\n", *ds.XMLID, *ds.ID, ds)
		}
		if !ServerCacheDotConfigIncludeInactiveDSes && !*ds.Active {
			continue
		}
		dsData[tc.DeliveryServiceName(*ds.XMLID)] = atscfg.ServerCacheConfigDS{OrgServerFQDN: *ds.OrgServerFQDN, Type: *ds.Type}
	}

	serverName := tc.CacheName(toData.Server.HostName)

	return atscfg.MakeServerCacheDotConfig(serverName, toData.TOToolName, toData.TOURL, dsData), atscfg.ContentTypeCacheDotConfig, nil
}
