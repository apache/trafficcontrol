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

var FormEditDeliveryServiceRequestController = function(deliveryServiceRequest, type, types, $scope, $state, $stateParams, $controller, $uibModal, locationUtils, deliveryServiceService, deliveryServiceRequestService, userModel) {

	// extends the FormDeliveryServiceController to inherit common methods
	angular.extend(this, $controller('FormDeliveryServiceController', { deliveryService: deliveryServiceRequest.request, type: type, types: types, $scope: $scope }));

	var deleteDeliveryServiceRequest = function(requestId) {
		deliveryServiceRequestService.deleteDeliveryServiceRequest(requestId)
			.then(function() {
				locationUtils.navigateToPath('/delivery-service-requests');
			});
	};

	$scope.requestType = deliveryServiceRequest.requestType;

	$scope.deliveryServiceName = angular.copy(deliveryServiceRequest.request.xmlId);

	$scope.advancedShowing = true;

	$scope.settings = {
		isNew: false,
		isRequest: true,
		saveLabel: 'Update Request',
		deleteLabel: 'Delete Request'
	};

	$scope.fulfill = function(deliveryService) {
		var params = {
			title: $scope.requestType + ' Delivery Service: ' + deliveryService.xmlId,
			message: 'Are you sure you want to fulfill this delivery service request and ' + $scope.requestType + ' ' + deliveryService.xmlId + '?'
		};
		var modalInstance = $uibModal.open({
			templateUrl: 'common/modules/dialog/confirm/dialog.confirm.tpl.html',
			controller: 'DialogConfirmController',
			size: 'md',
			resolve: {
				params: function () {
					return params;
				}
			}
		});
		modalInstance.result.then(function() {
			if ($scope.requestType == 'UPDATE') {
				deliveryServiceService.updateDeliveryService(deliveryService);
			} else if ($scope.requestType == 'CREATE') {
				deliveryServiceService.createDeliveryService(deliveryService);
			}
		}, function () {
			// do nothing
		});
	};

	$scope.save = function(deliveryService) {
		var params = {
			title: 'Edit Delivery Service Request',
			message: 'Delivery services changes must be reviewed for completeness and accuracy before deployment.'
		};
		var modalInstance = $uibModal.open({
			templateUrl: 'common/modules/dialog/select/dialog.select.tpl.html',
			controller: 'DialogSelectController',
			size: 'md',
			resolve: {
				params: function () {
					return params;
				},
				collection: function() {
					return [
						{ id: $scope.DRAFT, name: 'Save Request as Draft' },
						{ id: $scope.SUBMITTED, name: 'Submit Request for Review' }
					];
				}
			}
		});
		modalInstance.result.then(function(action) {
			var dsRequest = {
				id: Math.floor((Math.random() * 100) + 1),
				xmlId: deliveryService.xmlId,
				serviceType: 'DNS',
				tenantId: deliveryService.tenantId,
				cdnId: deliveryService.cdnId,
				requestType: 'UPDATE',
				status: (action.id == $scope.SUBMITTED) ? 'SUBMITTED' : 'DRAFT',
				author: userModel.user.id,
				request: deliveryService,
				lastUpdated: moment()
			};

			deliveryServiceRequestService.updateDeliveryServiceRequest(dsRequest);
		}, function () {
			// do nothing
		});

	};

	$scope.confirmDelete = function(deliveryService) {
		var params = {
			title: 'Delete Delivery Service Request for ' + deliveryService.xmlId + '?',
			message: 'Are you sure you want to delete this delivery service request?'
		};
		var modalInstance = $uibModal.open({
			templateUrl: 'common/modules/dialog/confirm/dialog.confirm.tpl.html',
			controller: 'DialogConfirmController',
			size: 'md',
			resolve: {
				params: function () {
					return params;
				}
			}
		});
		modalInstance.result.then(function() {
			deleteDeliveryServiceRequest($stateParams.deliveryServiceRequestId);
		}, function () {
			// do nothing
		});
	};


};

FormEditDeliveryServiceRequestController.$inject = ['deliveryServiceRequest', 'type', 'types', '$scope', '$state', '$stateParams', '$controller', '$uibModal', 'locationUtils', 'deliveryServiceService', 'deliveryServiceRequestService', 'userModel'];
module.exports = FormEditDeliveryServiceRequestController;
