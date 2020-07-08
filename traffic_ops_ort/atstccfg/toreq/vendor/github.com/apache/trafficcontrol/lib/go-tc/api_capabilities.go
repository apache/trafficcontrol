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

// APICapability represents an association between a Traffic Ops API route and a required capability.
type APICapability struct {
	ID          int       `json:"id" db:"id"`
	HTTPMethod  string    `json:"httpMethod" db:"http_method"`
	Route       string    `json:"httpRoute" db:"route"`
	Capability  string    `json:"capability" db:"capability"`
	LastUpdated TimeNoMod `json:"lastUpdated" db:"last_updated"`
}

// APICapabilityResponse represents an HTTP response to an API Capability request.
type APICapabilityResponse struct {
	Response []APICapability `json:"response"`
	Alerts
}
