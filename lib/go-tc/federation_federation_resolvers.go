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

// FederationFederationResolversResponse represents an API response of association between a federation and a federation_resolver.
type FederationFederationResolversResponse struct {
	Response []FederationResolver `json:"response"`
}

// AssignFederationFederationResolvers represents an API response of federation_resolver assignment to a federation.
type AssignFederationFederationResolversResponse struct {
	Alerts   []tc.Alerts                      `json:"alerts"`
	Response AssignFederationResolversRequest `json:"response"`
}

// AssignFederationResolversRequest represents an API request/response for assigning federation_resolvers to a federation.
type AssignFederationResolversRequest struct {
	Replace        bool  `json:"replace"`
	FedResolverIDs []int `json:"fedResolverIds"`
}
