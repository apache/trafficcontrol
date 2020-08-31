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

	"github.com/apache/trafficcontrol/lib/go-atscfg"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops_ort/atstccfg/config"
)

func GetMeta(toData *config.TOData, dir string) (*tc.ATSConfigMetaData, error) {
	if toData.Server.Cachegroup == nil {
		return nil, errors.New("this server missing Cachegroup")
	} else if toData.Server.HostName == nil {
		return nil, errors.New("this server missing HostName")
	} else if toData.Server.CachegroupID == nil {
		return nil, errors.New("this server missing CachegroupID")
	} else if toData.Server.Cachegroup == nil {
		return nil, errors.New("this server missing Cachegroup")
	} else if toData.Server.CDNID == nil {
		return nil, errors.New("this server missing CDNID")
	} else if toData.Server.CDNName == nil {
		return nil, errors.New("this server missing CDNName")
	} else if toData.Server.ID == nil {
		return nil, errors.New("this server missing ID")
	}

	toReverseProxyURL := ""
	toURL := ""
	for _, param := range toData.GlobalParams {
		if param.Name == "tm.rev_proxy.url" {
			toReverseProxyURL = param.Value
		} else if param.Name == "tm.url" {
			toURL = param.Value
		}
		if toReverseProxyURL != "" && toURL != "" {
			break
		}
	}

	scopeParams := ParamsToMap(toData.ScopeParams)

	locationParams := map[string]atscfg.ConfigProfileParams{}
	for _, param := range toData.ServerParams {
		if param.Name == "location" {
			p := locationParams[param.ConfigFile]
			p.FileNameOnDisk = param.ConfigFile
			p.Location = param.Value
			locationParams[param.ConfigFile] = p
		} else if param.Name == "URL" {
			p := locationParams[param.ConfigFile]
			p.URL = param.Value
			locationParams[param.ConfigFile] = p
		}
	}

	dses := map[tc.DeliveryServiceName]tc.DeliveryServiceV30{}
	if tc.CacheTypeFromString(toData.Server.Type) != tc.CacheTypeMid {
		dsIDs := map[int]struct{}{}
		for _, ds := range toData.DeliveryServices {
			if ds.ID == nil {
				// TODO log error?
				continue
			}
			dsIDs[*ds.ID] = struct{}{}
		}

		// TODO verify?
		//		serverIDs := []int{toData.Server.ID}

		dssMap := map[int]struct{}{}
		for _, dss := range toData.DeliveryServiceServers {
			if dss.Server == nil || dss.DeliveryService == nil {
				continue // TODO warn?
			}
			if *dss.Server != *toData.Server.ID {
				continue
			}
			if _, ok := dsIDs[*dss.DeliveryService]; !ok {
				continue
			}
			dssMap[*dss.DeliveryService] = struct{}{}
		}

		for _, ds := range toData.DeliveryServices {
			if ds.ID == nil {
				continue
			}
			if ds.XMLID == nil {
				continue // TODO log?
			}
			if _, ok := dssMap[*ds.ID]; !ok && ds.Topology == nil {
				continue
			}
			dses[tc.DeliveryServiceName(*ds.XMLID)] = ds
		}
	} else {
		for _, ds := range toData.DeliveryServices {
			if ds.ID == nil {
				continue
			}
			if ds.XMLID == nil {
				continue // TODO log?
			}
			if ds.CDNID == nil || *ds.CDNID != *toData.Server.CDNID {
				continue
			}
			dses[tc.DeliveryServiceName(*ds.XMLID)] = ds
		}
	}

	uriSignedDSes := []tc.DeliveryServiceName{}
	for _, ds := range toData.DeliveryServices {
		if ds.ID == nil {
			continue
		}
		if ds.XMLID == nil {
			continue // TODO log?
		}
		if _, ok := dses[tc.DeliveryServiceName(*ds.XMLID)]; !ok {
			continue
		}
		if ds.SigningAlgorithm == nil || *ds.SigningAlgorithm != tc.SigningAlgorithmURISigning {
			continue
		}
		uriSignedDSes = append(uriSignedDSes, tc.DeliveryServiceName(*ds.XMLID))
	}

	metaObj, err := atscfg.MakeMetaObj(toData.Server, toURL, toReverseProxyURL, locationParams, uriSignedDSes, scopeParams, dses, toData.CacheGroups, toData.Topologies, dir)
	if err != nil {
		return nil, errors.New("generating: " + err.Error())
	}
	return &metaObj, nil
}
