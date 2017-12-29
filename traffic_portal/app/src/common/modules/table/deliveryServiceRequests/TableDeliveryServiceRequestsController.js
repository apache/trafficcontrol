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

var TableDeliveryServicesRequestsController = function(dsRequests, $scope, $state, $uibModal, dateUtils, locationUtils, typeService, deliveryServiceRequestService, userModel) {

	$scope.DRAFT = 0;
	$scope.SUBMITTED = 1;
	$scope.REJECTED = 2;
	$scope.COMPLETE = 3;

	$scope.dsRequests = dsRequests;

	$scope.getRelativeTime = dateUtils.getRelativeTime;

	$scope.refresh = function() {
		$state.reload(); // reloads all the resolves for the view
	};

	$scope.assignRequest = function(request, assign, $event) {
		$event.stopPropagation(); // this kills the click event so it doesn't trigger anything else
		var params = {
			title: 'Assign Delivery Service Request',
			message: (assign) ? 'Are you sure you want to assign this delivery service request to yourself?' : 'Are you sure you want to unassign this delivery service request?'
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
			request.assigneeId = (assign) ? userModel.user.id : null;
			deliveryServiceRequestService.updateDeliveryServiceRequest(request.id, request).
				then(function() {
					$scope.refresh();
				});
		}, function () {
			// do nothing
		});
	};

	$scope.editStatus = function(request, $event) {
		$event.stopPropagation(); // this kills the click event so it doesn't trigger anything else
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
					return [
						{ id: $scope.DRAFT, name: 'Draft' },
						{ id: $scope.SUBMITTED, name: 'Submitted' },
						{ id: $scope.REJECTED, name: 'Rejected' },
						{ id: $scope.COMPLETE, name: 'Complete' }
					];
				}
			}
		});
		modalInstance.result.then(function(action) {
			switch (action.id) {
				case $scope.DRAFT:
					request.status = 'draft';
					break;
				case $scope.SUBMITTED:
					request.status = 'submitted';
					break;
				case $scope.REJECTED:
					request.status = 'rejected';
					break;
				case $scope.COMPLETE:
					request.status = 'complete';
			}
			deliveryServiceRequestService.updateDeliveryServiceRequest(request.id, request).
				then(function() {
					$scope.refresh();
				});
		}, function () {
			// do nothing
		});
	};

	$scope.deleteRequest = function(request, $event) {
		$event.stopPropagation(); // this kills the click event so it doesn't trigger anything else
		var params = {
			title: 'Delete ' + request.request.xmlId + ' ' + request.changeType + ' request?',
			key: request.request.xmlId + ' request'
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
			deliveryServiceRequestService.deleteDeliveryServiceRequest(request.id, false).
				then(function() {
					$scope.refresh();
				});
		}, function () {
			// do nothing
		});
	};

	$scope.editDeliveryServiceRequest = function(request) {
		var path = '/delivery-service-requests/' + request.id + '?type=';
		typeService.getType(request.request.typeId)
			.then(function(result) {
				path += result.name;
				locationUtils.navigateToPath(path);
			});
	};

	angular.element(document).ready(function () {
		$('#dsRequestsTable').dataTable({
			"aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
			"iDisplayLength": 25,
			"columnDefs": [
				{ 'orderable': false, 'targets': 6 }
			],
			"aaSorting": []
		});
	});

};

TableDeliveryServicesRequestsController.$inject = ['dsRequests', '$scope', '$state', '$uibModal', 'dateUtils', 'locationUtils', 'typeService', 'deliveryServiceRequestService', 'userModel'];
module.exports = TableDeliveryServicesRequestsController;
