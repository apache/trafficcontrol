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

module.exports = angular.module('trafficPortal.private.physLocations.servers', [])
	.config(function($stateProvider, $urlRouterProvider) {
		$stateProvider
			.state('trafficPortal.private.physLocations.servers', {
				url: '/{physLocationId}/servers',
				views: {
					physLocationsContent: {
						templateUrl: 'common/modules/table/physLocationServers/table.physLocationServers.tpl.html',
						controller: 'TablePhysLocationServersController',
						resolve: {
							physLocation: function($stateParams, physLocationService) {
								return physLocationService.getPhysLocation($stateParams.physLocationId);
							},
							servers: function(physLocation, $stateParams, serverService) {
								return serverService.getServers({ physLocation: physLocation.id, orderby: 'hostName' });
							},
							filter: function(physLocation) {
								return {
									physLocation: {
										filterType: "text",
										type: "equals",
										filter: physLocation.name
									}
								}
							}
						}
					}
				}
			})
		;
		$urlRouterProvider.otherwise('/');
	});
