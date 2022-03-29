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
	"github.com/apache/trafficcontrol/lib/go-atscfg"
	"github.com/apache/trafficcontrol/lib/go-tc"
)

// serversToLatest converts a []tc.Server to []tc.ServerV30.
// This is necessary, because the old Traffic Ops client doesn't return the same type as the latest client.
func serversToLatest(svs tc.ServersV3Response) ([]atscfg.Server, error) {
	nss := []atscfg.Server{}
	for _, sv := range svs.Response {
		svLatest, err := serverToLatest(&sv)
		if err != nil {
			return nil, err // serverToLatest adds context
		}
		nss = append(nss, atscfg.Server(*svLatest))
	}
	return nss, nil
}

// serverToLatest converts a tc.Server to tc.ServerV30.
// This is necessary, because the old Traffic Ops client doesn't return the same type as the latest client.
func serverToLatest(oldSv *tc.ServerV30) (*atscfg.Server, error) {
	sv, err := oldSv.UpgradeToV40([]string{*oldSv.Profile})
	if err != nil {
		return nil, err
	}
	asv := atscfg.Server(sv)
	return &asv, nil
}

func dsesToLatest(dses []tc.DeliveryServiceNullableV30) []atscfg.DeliveryService {
	newDSes := []tc.DeliveryServiceV40{}
	for _, ds := range dses {
		newDSes = append(newDSes, ds.UpgradeToV4())
	}
	return atscfg.ToDeliveryServices(newDSes)
}
