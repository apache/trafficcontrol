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

module.exports = angular.module('trafficPortal.private.deliveryServiceRequests.compare', [])
	.config(function($stateProvider, $urlRouterProvider) {
		$stateProvider
			.state('trafficPortal.private.deliveryServiceRequests.compare', {
				url: '/compare/{dsr1Id}/{dsr2Id}',
				views: {
					deliveryServiceRequestsContent: {
						templateUrl: 'common/modules/compare/compare.tpl.html',
						controller: 'CompareController',
						resolve: {
							dsr1: function($stateParams, deliveryServiceRequestService) {
								return deliveryServiceRequestService.getDeliveryServiceRequests({ id: $stateParams.dsr1Id });
							},
							dsr2: function($stateParams, deliveryServiceRequestService) {
								return deliveryServiceRequestService.getDeliveryServiceRequests({ id: $stateParams.dsr2Id });
							},
							item1Name: function(dsr1) {
								return dsr1[0].deliveryService.xmlId + ' ' + dsr1[0].changeType + ' (' + dsr1[0].author + ' created on ' + dsr1[0].createdAt + ')';
							},
							item2Name: function(dsr2) {
								return dsr2[0].deliveryService.xmlId + ' ' + dsr2[0].changeType + ' (' + dsr2[0].author + ' created on ' + dsr2[0].createdAt + ')';
							},
							item1: function(dsr1) {
								return dsr1[0].deliveryService;
							},
							item2: function(dsr2) {
								return dsr2[0].deliveryService;
							}
						}
					}
				}
			})
		;
		$urlRouterProvider.otherwise('/');
	});
