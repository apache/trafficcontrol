package tc

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

// HWInfoResponse is a list of HWInfos as a response.
type HWInfoResponse struct {
	Response []HWInfo `json:"response"`
}

// HWInfo can be used to return information about a server's hardware, but the
// corresponding Traffic Ops API route is deprecated and unusable without
// alteration.
type HWInfo struct {
	Description    string    `json:"description" db:"description"`
	ID             int       `json:"-" db:"id"`
	LastUpdated    TimeNoMod `json:"lastUpdated" db:"last_updated"`
	ServerHostName string    `json:"serverHostName" db:"serverhostname"`
	ServerID       int       `json:"serverId" db:"serverid"`
	Val            string    `json:"val" db:"val"`
}
