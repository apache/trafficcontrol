package v13

import tc "github.com/apache/incubator-trafficcontrol/lib/go-tc"

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

type InterfacesResponse struct {
	Response []Interface `json:"response"`
}

// IntfResponse is the JSON object returned for a single intf
type InterfaceResponse struct {
	Response Interface `json:"response"`
}

type Interface struct {
	ID            int          `json:"id" db:"id"`
	Server        string       `json:"server" db:"server"`
	ServerID      int          `json:"serverId" db:"server_id"`
	InterfaceName string       `json:"interfaceName" db:"interface_name"`
	InterfaceMtu  int          `json:"interfaceMtu" db:"interface_mtu"`
	LastUpdated   tc.TimeNoMod `json:"lastUpdated" db:"last_updated"`
}

type InterfaceNullable struct {
	ID            *int          `json:"id" db:"id"`
	Server        *string       `json:"server" db:"server"`
	ServerID      *int          `json:"serverId" db:"server_id"`
	InterfaceName *string       `json:"interfaceName" db:"interface_name"`
	InterfaceMtu  *int          `json:"interfaceMtu" db:"interface_mtu"`
	LastUpdated   *tc.TimeNoMod `json:"lastUpdated" db:"last_updated"`
}
