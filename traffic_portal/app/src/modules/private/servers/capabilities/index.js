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

module.exports = angular.module('trafficPortal.private.servers.capabilities', [])
	.config(function($stateProvider, $urlRouterProvider) {
		$stateProvider
			.state('trafficPortal.private.servers.capabilities', {
				url: '/{serverId}/capabilities',
				views: {
					serversContent: {
						templateUrl: 'common/modules/table/serverServerCapabilities/table.serverServerCapabilities.tpl.html',
						controller: 'TableServerServerCapabilitiesController',
						resolve: {
							server: function($stateParams, serverService) {
								return serverService.getServers({id: $stateParams.serverId});
							},
							serverCapabilities: function($stateParams, serverService) {
								return serverService.getServerCapabilities($stateParams.serverId);
							}
						}
					}
				}
			})
		;
		$urlRouterProvider.otherwise('/');
	});
