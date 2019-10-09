package main

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
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-atscfg"
	"github.com/apache/trafficcontrol/lib/go-tc"
)

const ServerCacheDotConfigIncludeInactiveDSes = false // TODO move to lib/go-atscfg

func GetConfigFileServerCacheDotConfig(cfg TCCfg, serverNameOrID string) (string, error) {
	// TODO TOAPI add /servers?cdn=1 query param
	servers, err := GetServers(cfg)
	if err != nil {
		return "", errors.New("getting servers: " + err.Error())
	}

	server := tc.Server{ID: atscfg.InvalidID}
	if serverID, err := strconv.Atoi(serverNameOrID); err == nil {
		for _, toServer := range servers {
			if toServer.ID == serverID {
				server = toServer
				break
			}
		}
	} else {
		serverName := serverNameOrID
		for _, toServer := range servers {
			if toServer.HostName == serverName {
				server = toServer
				break
			}
		}
	}
	if server.ID == atscfg.InvalidID {
		return "", errors.New("server '" + serverNameOrID + " not found in servers")
	}

	if !strings.HasPrefix(string(server.Type), tc.MidTypePrefix) {
		// emulates Perl
		return "", errors.New("Error - incorrect file scope for route used.  Please use the profiles route.")
	}

	dses, err := GetCDNDeliveryServices(cfg, server.CDNID)
	if err != nil {
		return "", errors.New("getting delivery services: " + err.Error())
	}

	dsData := map[tc.DeliveryServiceName]atscfg.ServerCacheConfigDS{}
	for _, ds := range dses {
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

	serverName := tc.CacheName(server.HostName)

	toToolName, toURL, err := GetTOToolNameAndURLFromTO(cfg)
	if err != nil {
		return "", errors.New("getting global parameters: " + err.Error())
	}

	txt := atscfg.MakeServerCacheDotConfig(serverName, toToolName, toURL, dsData)
	return txt, nil
}
