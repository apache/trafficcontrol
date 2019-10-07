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

var TableDeliveryServicesRequestsController = function (dsRequests, $scope, $state, $uibModal, $anchorScroll, $q, $location, dateUtils, locationUtils, typeService, deliveryServiceService, deliveryServiceRequestService, messageModel, propertiesModel, userModel) {

	var createComment = function (request, placeholder) {
		var params = {
			title: 'Add Comment',
			placeholder: placeholder,
			text: null
		};
		var modalInstance = $uibModal.open({
			templateUrl: 'common/modules/dialog/textarea/dialog.textarea.tpl.html',
			controller: 'DialogTextareaController',
			size: 'md',
			resolve: {
				params: function () {
					return params;
				}
			}
		});
		modalInstance.result.then(function (commentValue) {
			var comment = {
				deliveryServiceRequestId: request.id,
				value: commentValue
			};
			deliveryServiceRequestService.createDeliveryServiceRequestComment(comment);
		}, function () {
			// do nothing
		});
	};

	$scope.DRAFT = 0;
	$scope.SUBMITTED = 1;
	$scope.REJECTED = 2;
	$scope.PENDING = 3;
	$scope.COMPLETE = 4;

	$scope.dsRequests = dsRequests;

	$scope.getRelativeTime = dateUtils.getRelativeTime;

	$scope.refresh = function () {
		$state.reload(); // reloads all the resolves for the view
	};

	$scope.fulfillable = function (request) {
		return request.status == 'submitted';
	};

	$scope.rejectable = function (request) {
		return request.status == 'submitted';
	};

	$scope.completeable = function (request) {
		return request.status == 'pending';
	};

	$scope.open = function (request) {
		return (request.status == 'draft' || request.status == 'submitted');
	};

	$scope.closed = function (request) {
		return (request.status == 'rejected' || request.status == 'complete');
	};

	$scope.assignRequest = function (request, assign, $event) {
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
		modalInstance.result.then(function () {
			var assigneeId = (assign) ? userModel.user.id : null;
			deliveryServiceRequestService.assignDeliveryServiceRequest(request.id, assigneeId).then(function () {
				$scope.refresh();
			});
		}, function () {
			// do nothing
		});
	};

	$scope.editStatus = function (request, $event) {
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
				collection: function () {
					return [
						{id: $scope.DRAFT, name: 'Save as Draft'},
						{id: $scope.SUBMITTED, name: 'Submit for Review / Deployment'}
					];
				}
			}
		});
		modalInstance.result.then(function (action) {
			var status = (action.id == $scope.DRAFT) ? 'draft' : 'submitted';
			deliveryServiceRequestService.updateDeliveryServiceRequestStatus(request.id, status).then(function () {
				$scope.refresh();
			});
		}, function () {
			// do nothing
		});
	};

	$scope.rejectRequest = function (request, $event) {
		$event.stopPropagation(); // this kills the click event so it doesn't trigger anything else

		// only the user assigned to the request can mark it as rejected (unless the user has override capabilities)
		if ((request.assigneeId != userModel.user.id) && (userModel.user.roleName != propertiesModel.properties.dsRequests.overrideRole)) {
			messageModel.setMessages([{
				level: 'error',
				text: 'Only the assignee can mark a delivery service request as rejected'
			}], false);
			$anchorScroll(); // scrolls window to top
			return;
		}

		var params = {
			title: 'Reject Delivery Service Request',
			message: 'Are you sure you want to reject this delivery service request?'
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
		modalInstance.result.then(function () {
			deliveryServiceRequestService.updateDeliveryServiceRequestStatus(request.id, 'rejected').then(
				function () {
					$scope.refresh();
					createComment(request, 'Enter rejection reason...');
				});
		}, function () {
			// do nothing
		});
	};

	$scope.completeRequest = function (request, $event) {
		$event.stopPropagation(); // this kills the click event so it doesn't trigger anything else

		// only the user assigned to the request can mark it as complete (unless the user has override capabilities)
		if ((request.assigneeId != userModel.user.id) && (userModel.user.roleName != propertiesModel.properties.dsRequests.overrideRole)) {
			messageModel.setMessages([{
				level: 'error',
				text: 'Only the assignee can mark a delivery service request as complete'
			}], false);
			$anchorScroll(); // scrolls window to top
			return;
		}

		var params = {
			title: 'Complete Delivery Service Request',
			message: 'Are you sure you want to mark this delivery service request as complete?'
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
		modalInstance.result.then(function () {
			deliveryServiceRequestService.updateDeliveryServiceRequestStatus(request.id, 'complete').then(function () {
				$scope.refresh();
				createComment(request, 'Enter comment...');
			});
		}, function () {
			// do nothing
		});
	};

	$scope.deleteRequest = function (request, $event) {
		$event.stopPropagation(); // this kills the click event so it doesn't trigger anything else
		var params = {
			title: 'Delete the ' + request.deliveryService.xmlId + ' ' + request.changeType + ' request?',
			key: request.deliveryService.xmlId + ' ' + request.changeType + ' request'
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
		modalInstance.result.then(function () {
			deliveryServiceRequestService.deleteDeliveryServiceRequest(request.id).then(function () {
				messageModel.setMessages([{level: 'success', text: 'Delivery service request deleted'}], false);
				$scope.refresh();
			});
		}, function () {
			// do nothing
		});
	};

	$scope.editDeliveryServiceRequest = function (request) {
		var path = '/delivery-service-requests/' + request.id + '?type=';
		typeService.getType(request.deliveryService.typeId)
			.then(function (result) {
				path += result.name;
				locationUtils.navigateToPath(path);
			});
	};

	$scope.fulfillRequest = function (request, $event) {
		$event.stopPropagation(); // this kills the click event so it doesn't trigger anything else
		var path = '/delivery-service-requests/' + request.id + '?type=';
		typeService.getType(request.deliveryService.typeId)
			.then(function (result) {
				path += result.name;
				locationUtils.navigateToPath(path);
			});
	};

	angular.element(document).ready(function () {
		var dsRequestsTable = $('#dsRequestsTable').dataTable({
			"paging": false,
			"dom": '<"filter-checkbox">frtip',
			"columnDefs": [
				{'orderable': false, 'targets': 7}
			],
			"aaSorting": []
		});
		$('div.filter-checkbox').html('<input id="showClosed" type="checkbox"><label for="showClosed">&nbsp;&nbsp;Show Closed</label>');

		// only show "open" ds requests on render
		dsRequestsTable.fnFilter('draft|submitted|pending', 2, true, false);

		$('#showClosed').click(function () {
			var checked = $('#showClosed').is(':checked');
			if (checked) {
				dsRequestsTable.fnFilter('', 2, true, false);
			} else {
				dsRequestsTable.fnFilter('draft|submitted|pending', 2, true, false);
			}
		});

	});

};

TableDeliveryServicesRequestsController.$inject = ['dsRequests', '$scope', '$state', '$uibModal', '$anchorScroll', '$q', '$location', 'dateUtils', 'locationUtils', 'typeService', 'deliveryServiceService', 'deliveryServiceRequestService', 'messageModel', 'propertiesModel', 'userModel'];
module.exports = TableDeliveryServicesRequestsController;
