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

module.exports = angular.module('trafficPortal.private.deliveryServices.regexes.edit', [])
	.config(function($stateProvider, $urlRouterProvider) {
		$stateProvider
			.state('trafficPortal.private.deliveryServices.regexes.edit', {
				url: '/{regexId:[0-9]{1,8}}',
				views: {
					deliveryServiceRegexesContent: {
						templateUrl: 'common/modules/form/deliveryServiceRegex/form.deliveryServiceRegex.tpl.html',
						controller: 'FormEditDeliveryServiceRegexController',
						resolve: {
							deliveryService: function($stateParams, deliveryServiceService) {
								return deliveryServiceService.getDeliveryService($stateParams.deliveryServiceId);
							},
							regex: function($stateParams, deliveryServiceRegexService) {
								return deliveryServiceRegexService.getDeliveryServiceRegex($stateParams.deliveryServiceId, $stateParams.regexId);
							}
						}
					}
				}
			})
		;
		$urlRouterProvider.otherwise('/');
	});
