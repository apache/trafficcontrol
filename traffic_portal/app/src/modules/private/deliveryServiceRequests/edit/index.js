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

module.exports = angular.module('trafficPortal.private.deliveryServiceRequests.edit', [])
	.controller('FormEditDeliveryServiceRequestController', require('./FormEditDeliveryServiceRequestController'))
	.config(function ($stateProvider, $urlRouterProvider) {
		$stateProvider
			.state('trafficPortal.private.deliveryServiceRequests.edit', {
				url: '/{deliveryServiceRequestId:[0-9]{1,8}}?dsType&fulfill',
				views: {
					deliveryServiceRequestsContent: {
						templateUrl: function ($stateParams) {
							var type = $stateParams.dsType,
								template;

							if (type.indexOf('ANY_MAP') != -1) {
								template = 'common/modules/form/deliveryService/form.deliveryService.anyMap.tpl.html';
							} else if (type.indexOf('DNS') != -1) {
								template = 'common/modules/form/deliveryService/form.deliveryService.DNS.tpl.html';
							} else if (type.indexOf('HTTP') != -1) {
								template = 'common/modules/form/deliveryService/form.deliveryService.HTTP.tpl.html';
							} else if (type.indexOf('STEERING') != -1) {
								template = 'common/modules/form/deliveryService/form.deliveryService.Steering.tpl.html';
							} else {

							}

							return template;
						},
						controller: 'FormEditDeliveryServiceRequestController',
						resolve: {
							deliveryServiceRequest: function ($stateParams, deliveryServiceRequestService) {
								return deliveryServiceRequestService.getDeliveryServiceRequests({id: $stateParams.deliveryServiceRequestId});
							},
							dsCurrent: function (deliveryServiceRequest) {
								let dsRequest = deliveryServiceRequest[0];
								if (dsRequest.changeType == 'create') {
									// on create, there is nothing to compare so just return requested so there will be no diff
									return dsRequest.requested;
								} else if (dsRequest.changeType == 'update') {
									// on update, fetch the ds's original value so it can be compared against requested
									return dsRequest.original;
								} else if (dsRequest.changeType == 'delete') {
									// on delete, there is nothing to compare so just return original so there will be no diff
									return dsRequest.original;
								}
							},
							origin: function () {
								return [{}];
							},
							topologies: function(topologyService) {
								return topologyService.getTopologies();
							},
							type: function ($stateParams) {
								return $stateParams.dsType;
							},
							types: function (typeService) {
								return typeService.getTypes({useInTable: 'deliveryservice'});
							}
						}
					},
					deliveryServiceRequestsComments: {
						templateUrl: 'common/modules/table/deliveryServiceRequestComments/table.deliveryServiceRequestComments.tpl.html',
						controller: 'TableDeliveryServiceRequestCommentsController',
						resolve: {
							request: function ($stateParams, deliveryServiceRequestService) {
								return deliveryServiceRequestService.getDeliveryServiceRequests({id: $stateParams.deliveryServiceRequestId});
							},
							comments: function ($stateParams, deliveryServiceRequestService) {
								return deliveryServiceRequestService.getDeliveryServiceRequestComments({
									deliveryServiceRequestId: $stateParams.deliveryServiceRequestId,
									orderby: 'id'
								});
							}
						}
					}
				}
			})
		;
		$urlRouterProvider.otherwise('/');
	});
