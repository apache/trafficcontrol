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

var FormEditDeliveryServiceRequestController = function(deliveryServiceRequest, deliveryService, type, types, $scope, $state, $stateParams, $controller, $uibModal, $anchorScroll, $q, $location, locationUtils, deliveryServiceService, deliveryServiceRequestService, messageModel, userModel) {

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

	$scope.viewComments = function() {
		$location.path($location.path() + '/comments');
	};

	$scope.editStatus = function(status) {
		var params = {
			title: 'Change Delivery Service Request Status',
			message: "Are you sure you want to change the status of the delivery service request to '" + status + "'?"
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
			if (status == 'complete' && dsRequest.assigneeId != userModel.user.id) {
				messageModel.setMessages([ { level: 'error', text: 'Only the Assignee can mark a delivery service request as complete' } ], false);
				$anchorScroll(); // scrolls window to top
				return;
			}
			deliveryServiceRequestService.updateDeliveryServiceRequestStatus(dsRequest.id, status).
				then(function() {
					$state.reload();
				});
		}, function () {
			// do nothing
		});
	};

	var updateDeliveryServiceRequest = function() {
		var promises = [];
		// update the ds request if the ds request actually changed
		if ($scope.deliveryServiceForm.$dirty) {
			promises.push(deliveryServiceRequestService.updateDeliveryServiceRequest(dsRequest.id, dsRequest));
		}
		// make sure the ds request is assigned to the user that is fulfilling the request
		promises.push(deliveryServiceRequestService.assignDeliveryServiceRequest(dsRequest.id, userModel.user.id));
		// set the status to 'pending'
		promises.push(deliveryServiceRequestService.updateDeliveryServiceRequestStatus(dsRequest.id, 'pending'));
	};

	$scope.fulfillRequest = function(ds) {
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
			// update the ds request if the ds request actually changed
			if ($scope.deliveryServiceForm.$dirty) {
				promises.push(deliveryServiceRequestService.updateDeliveryServiceRequest(dsRequest.id, dsRequest));
			}
			// make sure the ds request is assigned to the user that is fulfilling the request
			promises.push(deliveryServiceRequestService.assignDeliveryServiceRequest(dsRequest.id, userModel.user.id));
			// set the status to 'pending'
			promises.push(deliveryServiceRequestService.updateDeliveryServiceRequestStatus(dsRequest.id, 'pending'));

			// create, update or delete the ds per the ds request
			if ($scope.changeType == 'create') {
				deliveryServiceService.createDeliveryService(ds).
					then(function() {
						$q.all(promises); // after a successful create, update the ds request, assignee and status
					});
			} else if ($scope.changeType == 'update') {
				deliveryServiceService.updateDeliveryService(ds, true).
					then(function() {
						$q.all(promises); // after a successful update, update the ds request, assignee and status
					});
			} else if ($scope.changeType == 'delete') {
				// and we're going to ask even again if they really want to delete but this time they need to enter the ds name to confirm the delete
				params = {
					title: 'Delete Delivery Service: ' + ds.xmlId,
					key: ds.xmlId
				};
				modalInstance = $uibModal.open({
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
					deliveryServiceService.deleteDeliveryService(ds).
						then(
							function() {
								updateDeliveryServiceRequest(); // after a successful delete, update the ds request, assignee and status and navigate to ds requests page
								messageModel.setMessages([ { level: 'success', text: 'Delivery service [ ' + ds.xmlId + ' ] deleted' } ], true);
								locationUtils.navigateToPath('/delivery-service-requests');
							},
							function(fault) {
								$anchorScroll(); // scrolls window to top
								messageModel.setMessages(fault.data.alerts, false);
							}
						);
				}, function () {
					// do nothing
				});
			}
		}, function () {
			// do nothing
		});
	};

	$scope.save = function(deliveryService) {
		var params = {
			title: 'Update Delivery Service Request',
			statusMessage: 'Please select the status of your delivery service request.',
			commentMessage: 'Why is this request being changed?'
		};
		var modalInstance = $uibModal.open({
			templateUrl: 'common/modules/dialog/deliveryServiceRequest/dialog.deliveryServiceRequest.tpl.html',
			controller: 'DialogDeliveryServiceRequestController',
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

FormEditDeliveryServiceRequestController.$inject = ['deliveryServiceRequest', 'deliveryService', 'type', 'types', '$scope', '$state', '$stateParams', '$controller', '$uibModal', '$anchorScroll', '$q', '$location', 'locationUtils', 'deliveryServiceService', 'deliveryServiceRequestService', 'messageModel', 'userModel'];
module.exports = FormEditDeliveryServiceRequestController;
