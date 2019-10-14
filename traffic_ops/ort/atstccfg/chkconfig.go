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

	"github.com/apache/trafficcontrol/lib/go-atscfg"
	"github.com/apache/trafficcontrol/lib/go-tc"
)

func GetConfigFileServerChkconfig(cfg TCCfg, serverNameOrID string) (string, error) {
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

	profileParams, err := GetProfileParameters(cfg, server.Profile)
	if err != nil {
		return "", errors.New("getting profile '" + server.Profile + "' parameters: " + err.Error())
	}
	if len(profileParams) == 0 {
		// The TO endpoint behind toclient.GetParametersByProfileName returns an empty object with a 200, if the Profile doesn't exist.
		// So we act as though we got a 404 if there are no params, to make ORT behave correctly.
		return "", ErrNotFound
	}

	fileParams := map[string][]string{}
	for _, param := range profileParams {
		if param.ConfigFile != atscfg.ChkconfigParamConfigFile {
			continue
		}
		fileParams[param.Name] = append(fileParams[param.Name], param.Value)
	}

	txt := atscfg.MakeChkconfig(fileParams)
	return txt, nil
}
