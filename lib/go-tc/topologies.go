package tc

import "time"

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

type TopologyV5 TopologyV50

type TopologyV50 struct {
	Description string           `json:"description" db:"description"`
	Name        string           `json:"name" db:"name"`
	Nodes       []TopologyNodeV5 `json:"nodes"`
	LastUpdated *time.Time       `json:"lastUpdated" db:"last_updated"`
}

// Topology holds the name and set of TopologyNodes that comprise a flexible topology.
type Topology struct {
	Description string         `json:"description" db:"description"`
	Name        string         `json:"name" db:"name"`
	Nodes       []TopologyNode `json:"nodes"`
	LastUpdated *TimeNoMod     `json:"lastUpdated" db:"last_updated"`
}

// TopologyNode holds a reference to a cachegroup and the indices of up to 2 parent
// nodes in the same topology's array of nodes.
type TopologyNode struct {
	Id          int        `json:"-" db:"id"`
	Cachegroup  string     `json:"cachegroup" db:"cachegroup"`
	Parents     []int      `json:"parents"`
	LastUpdated *TimeNoMod `json:"-" db:"last_updated"`
}

type TopologyNodeV5 TopologyNodeV50

type TopologyNodeV50 struct {
	Id          int        `json:"-" db:"id"`
	Cachegroup  string     `json:"cachegroup" db:"cachegroup"`
	Parents     []int      `json:"parents"`
	LastUpdated *time.Time `json:"-" db:"last_updated"`
}

type TopologyResponseV5 TopologyResponseV50

type TopologyResponseV50 struct {
	Response TopologyV5 `json:"response"`
	Alerts
}

// TopologyResponse models the JSON object returned for a single Topology in a
// response from the Traffic Ops API.
type TopologyResponse struct {
	Response Topology `json:"response"`
	Alerts
}

type TopologiesResponseV5 TopologiesResponseV50

type TopologiesResponseV50 struct {
	Response []TopologyV5 `json:"response"`
	Alerts
}

// TopologiesResponse models the JSON object returned for a list of topologies in a response.
type TopologiesResponse struct {
	Response []Topology `json:"response"`
	Alerts
}

// TopologiesQueueUpdateRequest encodes the request data for the POST
// topologies/{{name}}/queue_update endpoint.
type TopologiesQueueUpdateRequest struct {
	Action string `json:"action"`
	CDNID  int64  `json:"cdnId"`
}

// TopologiesQueueUpdateResponse encodes the response data for the POST
// topologies/{{name}}/queue_update endpoint.
type TopologiesQueueUpdateResponse struct {
	TopologiesQueueUpdate `json:"response"`
	Alerts
}

// TopologiesQueueUpdate decodes the update data from the POST
// topologies/{{name}}/queue_update endpoint.
type TopologiesQueueUpdate struct {
	Action   string       `json:"action"`
	CDNID    int64        `json:"cdnId"`
	Topology TopologyName `json:"topology"`
}
