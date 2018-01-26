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

var FormEditDeliveryServiceRequestController = function(deliveryServiceRequest, deliveryService, type, types, $scope, $state, $stateParams, $controller, $uibModal, $anchorScroll, $q, locationUtils, deliveryServiceService, deliveryServiceRequestService, messageModel, userModel) {

	var dsRequest = deliveryServiceRequest[0];
		
	// extends the FormDeliveryServiceController to inherit common methods
	angular.extend(this, $controller('FormDeliveryServiceController', { deliveryService: dsRequest.deliveryService, dsCurrent: deliveryService, type: type, types: types, $scope: $scope }));

	$scope.changeType = dsRequest.changeType;

	$scope.requestStatus = dsRequest.status;

	$scope.deliveryServiceName = angular.copy(dsRequest.deliveryService.xmlId);

	$scope.advancedShowing = true;

	$scope.settings = {
		isNew: false,
		isRequest: true,
		saveLabel: 'Update Request',
		deleteLabel: 'Delete Request'
	};

	$scope.saveable = function() {
		return (dsRequest.status == 'draft' || dsRequest.status == 'submitted');
	};

	$scope.deletable = function() {
		return (dsRequest.status == 'draft' || dsRequest.status == 'submitted');
	};

	$scope.fulfillable = function() {
		return dsRequest.status == 'submitted';
	};

	$scope.open = function() {
		return (dsRequest.status == 'draft' || dsRequest.status == 'submitted' || dsRequest.status == 'pending');
	};

	$scope.magicNumberLabel = function(collection, magicNumber) {
		var item = _.findWhere(collection, { value: magicNumber });
		return item.label;
	};

	$scope.editStatus = function() {
		var params = {
			title: "Edit Delivery Service Request Status",
			message: 'Please select the appropriate status for this request.'
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
					var statuses = [];
					if (dsRequest.status == 'draft' || dsRequest.status == 'submitted') {
						statuses.push({ id: $scope.DRAFT, name: 'Save as Draft' });
						statuses.push({ id: $scope.SUBMITTED, name: 'Submit for Review / Deployment' });
					} else if (dsRequest.status == 'pending') {
						statuses.push({ id: $scope.COMPLETE, name: 'Complete' });
					}
					return statuses;
				}
			}
		});
		modalInstance.result.then(function(action) {
			switch (action.id) {
				case $scope.DRAFT:
					dsRequest.status = 'draft';
					break;
				case $scope.SUBMITTED:
					dsRequest.status = 'submitted';
					break;
				case $scope.COMPLETE:
					if (dsRequest.assigneeId != userModel.user.id) {
						messageModel.setMessages([ { level: 'error', text: 'Only the Assignee can mark a delivery service request as complete' } ], false);
						$anchorScroll(); // scrolls window to top
						return;
					}
					dsRequest.status = 'complete';
			}
			// todo jeremy: this needs to call the api to update ds request status
			deliveryServiceRequestService.updateDeliveryServiceRequest(dsRequest.id, dsRequest).
				then(function() {
					$state.reload();
				});
		}, function () {
			// do nothing
		});
	};

	$scope.fulfillRequest = function(ds) {
		var promises = [];
		var params = {
			title: 'Delivery Service ' + $scope.changeType + ': ' + ds.xmlId,
			message: 'Are you sure you want to fulfill this delivery service request and ' + $scope.changeType + ' the ' + ds.xmlId + ' delivery service'
		};
		params['message'] += ($scope.changeType == 'create' || $scope.changeType == 'update') ? ' with these configuration settings?' : '?';
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
			// set the status to 'pending'
			// todo jeremy: this needs to call the api to update ds request status
			dsRequest.status = 'pending';
			// update the ds request
			promises.push(deliveryServiceRequestService.updateDeliveryServiceRequest(dsRequest.id, dsRequest));
			// now create, update or delete the ds per the ds request
			if ($scope.changeType == 'create') {
				promises.push(deliveryServiceService.createDeliveryService(ds));
			} else if ($scope.changeType == 'update') {
				promises.push(deliveryServiceService.updateDeliveryService(ds, true));
			} else if ($scope.changeType == 'delete') {
				promises.push(deliveryServiceService.deleteDeliveryService(ds, true));
			}

			$q.all(promises)
				.then(
					function() {
						if ($scope.changeType == 'delete') {
							locationUtils.navigateToPath('/delivery-service-requests');
						}
					});

		}, function () {
			// do nothing
		});
	};

	$scope.save = function(deliveryService) {
		var params = {
			title: 'Delivery Service Request Status',
			message: 'Please select the status of your delivery service request.'
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
						{ id: $scope.DRAFT, name: 'Save as Draft' },
						{ id: $scope.SUBMITTED, name: 'Submit for Review / Deployment' }
					];
				}
			}
		});
		modalInstance.result.then(function(action) {
			dsRequest.status = (action.id == $scope.SUBMITTED) ? 'submitted' : 'draft';
			dsRequest.deliveryService = deliveryService;
			deliveryServiceRequestService.updateDeliveryServiceRequest(dsRequest.id, dsRequest).
			then(function() {
				messageModel.setMessages([ { level: 'success', text: 'Updated delivery service request for ' + dsRequest.deliveryService.xmlId + ' and set status to ' + dsRequest.status } ], false);
				$anchorScroll(); // scrolls window to top
				$state.reload();
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

FormEditDeliveryServiceRequestController.$inject = ['deliveryServiceRequest', 'deliveryService', 'type', 'types', '$scope', '$state', '$stateParams', '$controller', '$uibModal', '$anchorScroll', '$q', 'locationUtils', 'deliveryServiceService', 'deliveryServiceRequestService', 'messageModel', 'userModel'];
module.exports = FormEditDeliveryServiceRequestController;
