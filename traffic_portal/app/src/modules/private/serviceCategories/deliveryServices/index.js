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

module.exports = angular.module('trafficPortal.private.serviceCategories.deliveryServices', [])
	.config(function($stateProvider, $urlRouterProvider) {
		$stateProvider
			.state('trafficPortal.private.serviceCategories.deliveryServices', {
				url: '/{serviceCategory}/delivery-services',
				views: {
					serviceCategoriesContent: {
						templateUrl: 'common/modules/table/serviceCategoryDeliveryServices/table.serviceCategoryDeliveryServices.tpl.html',
						controller: 'TableServiceCategoryDeliveryServicesController',
						resolve: {
							serviceCategory: function($stateParams, serviceCategoryService) {
								return serviceCategoryService.getServiceCategory($stateParams.serviceCategory);
							},
							deliveryServices: function(serviceCategory, deliveryServiceService) {
								return deliveryServiceService.getDeliveryServices({ serviceCategory: serviceCategory.name });
							},
							steeringTargets: function (deliveryServiceService) {
								return deliveryServiceService.getSteering();
							},
							filter: function(serviceCategory) {
								return {
									serviceCategory: {
										filterType: "text",
										type: "equals",
										filter: serviceCategory.name
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
