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

module.exports = angular.module('trafficPortal.private.certExpirations.list', [])
	.config(function($stateProvider, $urlRouterProvider) {
		$stateProvider
			.state('trafficPortal.private.certExpirations.list', {
				url: '',
				views: {
					certExpirationsContent: {
						templateUrl: 'common/modules/table/certExpirations/table.certExpirations.tpl.html',
						controller: 'TableCertExpirationsController',
						resolve: {
							certExpirations: function(certExpirationsService) {
								return certExpirationsService.getCertExpirations();
							},
							tableName: function() {
								return 'cert-expirations';
							},
							deliveryservices: function(deliveryServiceService) {
								return deliveryServiceService.getDeliveryServices();
							}
						}
					}
				}
			})
		;
		$urlRouterProvider.otherwise('/');
	});
