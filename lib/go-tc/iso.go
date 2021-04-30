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

// OSVersionsResponse is the JSON representation of the
// OS versions data for ISO generation.
type OSVersionsResponse map[string]string

// OSVersionsAPIResponse is the type of a response from Traffic Ops to a
// request to its /osversions endpoint.
type OSVersionsAPIResponse struct {
	// Structure of this map:
	//  key:   Name of OS
	//  value: Directory where the ISO source can be found
	Response map[string]string `json:"response"`
	Alerts
}
