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

module.exports = angular.module('trafficPortal.private.deliveryServices.compare', [])
	.config(function($stateProvider, $urlRouterProvider) {
		$stateProvider
			.state('trafficPortal.private.deliveryServices.compare', {
				url: '/compare/{ds1Id}/{ds2Id}',
				views: {
					deliveryServicesContent: {
						templateUrl: 'common/modules/compare/compare.tpl.html',
						controller: 'CompareController',
						resolve: {
							ds1: function($stateParams, deliveryServiceService) {
								return deliveryServiceService.getDeliveryService($stateParams.ds1Id);
							},
							ds2: function($stateParams, deliveryServiceService) {
								return deliveryServiceService.getDeliveryService($stateParams.ds2Id);
							},
							item1Name: function(ds1) {
								return ds1.xmlId;
							},
							item2Name: function(ds2) {
								return ds2.xmlId;
							},
							item1: function(ds1) {
								return ds1;
							},
							item2: function(ds2) {
								return ds2;
							}
						}
					}
				}
			})
		;
		$urlRouterProvider.otherwise('/');
	});
