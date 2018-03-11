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

module.exports = angular.module('trafficPortal.private.cdns.federations.edit', [])
	.config(function($stateProvider, $urlRouterProvider) {
		$stateProvider
			.state('trafficPortal.private.cdns.federations.edit', {
				url: '/{fedId:[0-9]{1,8}}',
				views: {
					cdnFederationsContent: {
						templateUrl: 'common/modules/form/federation/form.federation.tpl.html',
						controller: 'FormEditFederationController',
						resolve: {
							cdn: function($stateParams, cdnService) {
								return cdnService.getCDN($stateParams.cdnId);
							},
							federation: function(cdn, $stateParams, federationService) {
								return federationService.getCDNFederation(cdn.name, $stateParams.fedId);
							},
							resolvers: function(federation, federationService) {
								return federationService.getFederationFederationResolvers(federation.id);
							},
							deliveryServices: function(cdn, deliveryServiceService) {
								return deliveryServiceService.getDeliveryServices({ cdn: cdn.id });
							},
							federationDeliveryServices: function(federation, federationService) {
								return federationService.getFederationDeliveryServices(federation.id);
							}
						}
					}
				}
			})
		;
		$urlRouterProvider.otherwise('/');
	});
