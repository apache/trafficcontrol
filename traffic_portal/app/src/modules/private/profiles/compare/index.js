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

module.exports = angular.module('trafficPortal.private.profiles.compare', [])
	.config(function($stateProvider, $urlRouterProvider) {
		$stateProvider
			.state('trafficPortal.private.profiles.compare', {
				url: '/compare/{profile1Id}/{profile2Id}',
				views: {
					profilesContent: {
						templateUrl: 'common/modules/compare/compare.tpl.html',
						controller: 'CompareController',
						resolve: {
							profile1: function($stateParams, profileService) {
								return profileService.getProfile($stateParams.profile1Id, { includeParams: true });
							},
							profile2: function($stateParams, profileService) {
								return profileService.getProfile($stateParams.profile2Id, { includeParams: true });
							},
							item1Name: function(profile1) {
								return profile1.name;
							},
							item2Name: function(profile2) {
								return profile2.name;
							},
							item1: function(profile1) {
								return profile1.params;
							},
							item2: function(profile2) {
								return profile2.params;
							}
						}
					}
				}
			})
		;
		$urlRouterProvider.otherwise('/');
	});
