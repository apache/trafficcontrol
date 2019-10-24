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

	"github.com/apache/trafficcontrol/lib/go-atscfg"
	"github.com/apache/trafficcontrol/lib/go-log"
)

func GetConfigFileCDNRegexRemap(cfg TCCfg, cdnNameOrID string, fileName string) (string, error) {
	cdnName, err := GetCDNNameFromCDNNameOrID(cfg, cdnNameOrID)
	if err != nil {
		return "", errors.New("getting CDN name from '" + cdnNameOrID + "': " + err.Error())
	}

	toToolName, toURL, err := GetTOToolNameAndURLFromTO(cfg)
	if err != nil {
		return "", errors.New("getting global parameters: " + err.Error())
	}

	cdn, err := GetCDN(cfg, cdnName)
	if err != nil {
		return "", errors.New("getting cdn '" + string(cdnName) + "': " + err.Error())
	}

	dses, err := GetCDNDeliveryServices(cfg, cdn.ID)
	if err != nil {
		return "", errors.New("getting delivery services: " + err.Error())
	}

	log.Errorf("DEBUG gcfcrr dses %+v\n", len(dses))

	cfgDSes := atscfg.DeliveryServicesToCDNDSes(dses)

	log.Errorf("DEBUG gcfcrr cfgdses %+v\n", len(cfgDSes))

	txt := atscfg.MakeRegexRemapDotConfig(cdnName, toToolName, toURL, fileName, cfgDSes)
	return txt, nil
}
