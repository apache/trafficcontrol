package toreqold

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
	"github.com/apache/trafficcontrol/v8/lib/go-atscfg"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
)

// serversToLatest converts a []tc.Server to []tc.ServerV30.
// This is necessary, because the old Traffic Ops client doesn't return the same type as the latest client.
func serversToLatest(svs tc.ServersV4Response) ([]atscfg.Server, error) {
	nss := []atscfg.Server{}
	for _, sv := range svs.Response {
		svLatest, err := serverToLatest(sv)
		if err != nil {
			return nil, err // serverToLatest adds context
		}
		nss = append(nss, atscfg.Server(*svLatest))
	}
	return nss, nil
}

// serverToLatest converts a tc.Server to tc.ServerV30.
// This is necessary, because the old Traffic Ops client doesn't return the same type as the latest client.
func serverToLatest(oldSv tc.ServerV40) (*atscfg.Server, error) {
	asv := atscfg.Server(oldSv.Upgrade())
	return &asv, nil
}

func dsesToLatest(dses []tc.DeliveryServiceV4) []atscfg.DeliveryService {
	v5DSes := []tc.DeliveryServiceV5{}
	for _, ds := range dses {
		v5DSes = append(v5DSes, ds.Upgrade())
	}
	return atscfg.ToDeliveryServices(v5DSes)
}

func serverUpdateStatusToLatest(status []tc.ServerUpdateStatusV40) []atscfg.ServerUpdateStatus {
	nStats := []tc.ServerUpdateStatusV50{}
	for _, stat := range status {
		nStat := tc.ServerUpdateStatusV50{
			HostName:             stat.HostName,
			UpdatePending:        stat.UpdatePending,
			RevalPending:         stat.RevalPending,
			UseRevalPending:      stat.UseRevalPending,
			HostId:               stat.HostId,
			Status:               stat.Status,
			ParentPending:        stat.ParentPending,
			ConfigUpdateTime:     stat.ConfigUpdateTime,
			ConfigApplyTime:      stat.ConfigApplyTime,
			RevalidateUpdateTime: stat.RevalidateUpdateTime,
			RevalidateApplyTime:  stat.RevalidateApplyTime,
		}
		nStats = append(nStats, nStat)
	}
	return atscfg.ToServerUpdateStatuses(nStats)
}

func deliveryServiceServersToLatest(dsServers []tc.DeliveryServiceServer) []tc.DeliveryServiceServerV5 {
	nDss := []tc.DeliveryServiceServerV5{}
	for _, dss := range dsServers {
		dsServer := tc.DeliveryServiceServerV5{
			Server:          dss.Server,
			DeliveryService: dss.DeliveryService,
			LastUpdated:     &dss.LastUpdated.Time,
		}
		nDss = append(nDss, dsServer)
	}
	return nDss
}

func parametersToLatest(params []tc.Parameter) []tc.ParameterV5 {
	nParams := []tc.ParameterV5{}
	for _, param := range params {
		np := tc.ParameterV5{
			ConfigFile:  param.ConfigFile,
			ID:          param.ID,
			LastUpdated: param.LastUpdated.Time,
			Name:        param.Name,
			Profiles:    param.Profiles,
			Secure:      param.Secure,
			Value:       param.Value,
		}
		nParams = append(nParams, np)
	}
	return nParams
}
