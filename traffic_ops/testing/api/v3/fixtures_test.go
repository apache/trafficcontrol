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

package v3

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
)

var fixturesFilePath string

// LoadFixtures loads the testing fixture data.
func LoadFixtures(fixturesPath string) {
	fixturesFilePath = fixturesPath

	f, err := ioutil.ReadFile(fixturesPath)
	if err != nil {
		log.Errorf("Cannot unmarshal fixtures json %s", err)
		os.Exit(1)
	}
	err = json.Unmarshal(f, &testData)
	if err != nil {
		log.Errorf("Cannot unmarshal fixtures json %v", err)
		os.Exit(1)
	}
}

// Reloads the testing fixture data from the last location from which they were
// loaded with LoadFixtures - this MUST be called after LoadFixtures.
func ReloadFixtures() {
	testData = TrafficControl{
		ASNs:                                 []tc.ASN{},
		CDNs:                                 []tc.CDN{},
		CacheGroups:                          []tc.CacheGroupNullable{},
		CacheGroupParameterRequests:          []tc.CacheGroupParameterRequest{},
		Capabilities:                         []tc.Capability{},
		Coordinates:                          []tc.Coordinate{},
		DeliveryServicesRegexes:              []tc.DeliveryServiceRegexesTest{},
		DeliveryServiceRequests:              []tc.DeliveryServiceRequestV30{},
		DeliveryServiceRequestComments:       []tc.DeliveryServiceRequestComment{},
		DeliveryServices:                     []tc.DeliveryServiceV30{},
		DeliveryServicesRequiredCapabilities: []tc.DeliveryServicesRequiredCapability{},
		Divisions:                            []tc.Division{},
		Federations:                          []tc.CDNFederation{},
		FederationResolvers:                  []tc.FederationResolver{},
		Origins:                              []tc.Origin{},
		Profiles:                             []tc.Profile{},
		Parameters:                           []tc.Parameter{},
		ProfileParameters:                    []tc.ProfileParameter{},
		PhysLocations:                        []tc.PhysLocation{},
		Regions:                              []tc.Region{},
		Roles:                                []tc.Role{},
		Servers:                              []tc.ServerNullable{},
		ServerServerCapabilities:             []tc.ServerServerCapability{},
		ServerCapabilities:                   []tc.ServerCapability{},
		ServiceCategories:                    []tc.ServiceCategory{},
		Statuses:                             []tc.StatusNullable{},
		StaticDNSEntries:                     []tc.StaticDNSEntry{},
		StatsSummaries:                       []tc.StatsSummary{},
		Tenants:                              []tc.Tenant{},
		ServerCheckExtensions:                []tc.ServerCheckExtensionNullable{},
		Topologies:                           []tc.Topology{},
		Types:                                []tc.Type{},
		SteeringTargets:                      []tc.SteeringTargetNullable{},
		Serverchecks:                         []tc.ServercheckRequestNullable{},
		Users:                                []tc.User{},
		InvalidationJobs:                     []tc.InvalidationJobInput{},
	}
	LoadFixtures(fixturesFilePath)
}
