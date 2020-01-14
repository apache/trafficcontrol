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
	"github.com/apache/trafficcontrol/lib/go-tc/tce"

	"github.com/apache/trafficcontrol/lib/go-atscfg"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/ort/atstccfg/config"
	"github.com/apache/trafficcontrol/traffic_ops/ort/atstccfg/toreq"
)

func GetConfigFileCDNSSLMultiCertDotConfig(cfg config.TCCfg, cdnNameOrID string) (string, error) {
	cdnName, err := toreq.GetCDNNameFromCDNNameOrID(cfg, cdnNameOrID)
	if err != nil {
		return "", errors.New("getting CDN name from '" + cdnNameOrID + "': " + err.Error())
	}

	toToolName, toURL, err := toreq.GetTOToolNameAndURLFromTO(cfg)
	if err != nil {
		return "", errors.New("getting global parameters: " + err.Error())
	}

	cdn, err := toreq.GetCDN(cfg, cdnName)
	if err != nil {
		return "", errors.New("getting cdn '" + string(cdnName) + "': " + err.Error())
	}

	dses, err := toreq.GetCDNDeliveryServices(cfg, cdn.ID)
	if err != nil {
		return "", errors.New("getting delivery services: " + err.Error())
	}

	filteredDSes := []tc.DeliveryServiceNullable{}
	for _, ds := range dses {
		// ANY_MAP and STEERING DSes don't have origins, and thus can't be put into the ssl config.
		if ds.Type != nil && (*ds.Type == tce.DSTypeAnyMap || *ds.Type == tce.DSTypeSteering) {
			continue
		}
		filteredDSes = append(filteredDSes, ds)
	}

	cfgDSes := atscfg.DeliveryServicesToSSLMultiCertDSes(filteredDSes)

	txt := atscfg.MakeSSLMultiCertDotConfig(cdnName, toToolName, toURL, cfgDSes)
	return txt, nil
}
