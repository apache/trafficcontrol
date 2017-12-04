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

module.exports = angular.module('trafficPortal.private.types.staticDnsEntries', [])
	.config(function($stateProvider, $urlRouterProvider) {
		$stateProvider
			.state('trafficPortal.private.types.staticDnsEntries', {
				url: '/{typeId}/static-dns-entries',
				views: {
					typesContent: {
						templateUrl: 'common/modules/table/typeStaticDnsEntries/table.typeStaticDnsEntries.tpl.html',
						controller: 'TableTypeStaticDnsEntriesController',
						resolve: {
							type: function($stateParams, typeService) {
								return typeService.getType($stateParams.typeId);
							},
							staticDnsEntries: function(type, staticDnsEntryService) {
								return staticDnsEntryService.getStaticDnsEntries({ type: type.name });
							}
						}
					}
				}
			})
		;
		$urlRouterProvider.otherwise('/');
	});
