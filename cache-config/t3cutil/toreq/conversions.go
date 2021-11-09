package toreq

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

func serversToLatest(svs tc.ServersV4Response) ([]atscfg.Server, error) {
	return atscfg.ToServers(svs.Response), nil
}

func serverToLatest(oldSv *tc.ServerV40) (*atscfg.Server, error) {
	asv := atscfg.Server(*oldSv)
	return &asv, nil
}

func dsesToLatest(dses []tc.DeliveryServiceV40) []atscfg.DeliveryService {
	return atscfg.V40ToDeliveryServices(dses)
}

func jobsToLatest(jobs []tc.InvalidationJobV4) []atscfg.InvalidationJob {
	return atscfg.ToInvalidationJobs(jobs)
}
