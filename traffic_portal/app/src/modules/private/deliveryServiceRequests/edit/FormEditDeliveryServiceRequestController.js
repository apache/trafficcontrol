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

var FormEditDeliveryServiceRequestController = function(deliveryServiceRequest, type, types, $scope, $state, $stateParams, $controller, $uibModal, $anchorScroll, locationUtils, deliveryServiceService, deliveryServiceRequestService, userModel) {

	var dsRequest = deliveryServiceRequest[0];
		
	// extends the FormDeliveryServiceController to inherit common methods
	angular.extend(this, $controller('FormDeliveryServiceController', { deliveryService: dsRequest.request, type: type, types: types, $scope: $scope }));
	
	$scope.changeType = dsRequest.changeType;

	$scope.deliveryServiceName = angular.copy(dsRequest.request.xmlId);

	$scope.advancedShowing = true;

	$scope.settings = {
		isNew: false,
		isRequest: true,
		saveLabel: 'Update Request',
		deleteLabel: 'Delete Request'
	};

	$scope.fulfill = function(deliveryService) {
		var params = {
			title: 'Delivery Service ' + $scope.changeType + ': ' + deliveryService.xmlId,
			message: 'Are you sure you want to fulfill this delivery service request and ' + $scope.changeType + ' the ' + deliveryService.xmlId + ' delivery service?'
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
			// make sure the ds request is assigned to the user that is fulfilling the request
			dsRequest.assigneeId = userModel.user.id;
			deliveryServiceRequestService.updateDeliveryServiceRequest(dsRequest.id, dsRequest);
			// now update or create the ds per the ds request
			if ($scope.changeType == 'update') {
				deliveryServiceService.updateDeliveryService(deliveryService);
			} else if ($scope.changeType == 'create') {
				deliveryServiceService.createDeliveryService(deliveryService);
			}
		}, function () {
			// do nothing
		});
	};

	$scope.save = function(deliveryService) {
		var params = {
			title: 'Edit Delivery Service Request',
			message: 'All delivery service changes must be reviewed for completeness and accuracy before deployment. Please select the status of your request.'
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
						{ id: $scope.SUBMITTED, name: 'Submit Request for Review / Deployment' }
					];
				}
			}
		});
		modalInstance.result.then(function(action) {
			dsRequest.status = (action.id == $scope.SUBMITTED) ? 'submitted' : 'draft';
			dsRequest.request = deliveryService;
			deliveryServiceRequestService.updateDeliveryServiceRequest(dsRequest.id, dsRequest).
				then(function() {
					$anchorScroll(); // scrolls window to top
				});

		}, function () {
			// do nothing
		});

	};

	$scope.confirmDelete = function(deliveryService) {
		var params = {
			title: 'Delete ' + deliveryService.xmlId + ' ' + dsRequest.changeType + ' request?',
			key: deliveryService.xmlId + ' request'
		};
		var modalInstance = $uibModal.open({
			templateUrl: 'common/modules/dialog/delete/dialog.delete.tpl.html',
			controller: 'DialogDeleteController',
			size: 'md',
			resolve: {
				params: function () {
					return params;
				}
			}
		});
		modalInstance.result.then(function() {
			deliveryServiceRequestService.deleteDeliveryServiceRequest($stateParams.deliveryServiceRequestId, true);
		}, function () {
			// do nothing
		});
	};

};

FormEditDeliveryServiceRequestController.$inject = ['deliveryServiceRequest', 'type', 'types', '$scope', '$state', '$stateParams', '$controller', '$uibModal', '$anchorScroll', 'locationUtils', 'deliveryServiceService', 'deliveryServiceRequestService', 'userModel'];
module.exports = FormEditDeliveryServiceRequestController;
