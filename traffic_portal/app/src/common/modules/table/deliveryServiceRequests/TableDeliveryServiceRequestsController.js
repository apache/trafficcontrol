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

var TableDeliveryServicesRequestsController = function(dsRequests, $scope, $state, $uibModal, $anchorScroll, $q, dateUtils, locationUtils, typeService, deliveryServiceService, deliveryServiceRequestService, messageModel, userModel) {

	var createDeliveryServiceDeleteRequest = function(deliveryService) {
		var params = {
			title: "Delivery Service Delete Request",
			message: 'All delivery service deletions must be reviewed.<br><br>Are you sure you want to submit a request to delete the ' + deliveryService.xmlId + ' delivery service?'
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
			var dsRequest = {
				changeType: 'delete',
				status: 'submitted',
				deliveryService: deliveryService
			};
			deliveryServiceRequestService.createDeliveryServiceRequest(dsRequest, false).
				then(function() {
					$scope.refresh();
				});
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

	$scope.refresh = function() {
		$state.reload(); // reloads all the resolves for the view
	};

	$scope.fulfillable = function(request) {
		return request.status == 'submitted';
	};

	$scope.rejectable = function(request) {
		return request.status == 'submitted';
	};

	$scope.completeable = function(request) {
		return request.status == 'pending';
	};

	$scope.open = function(request) {
		return (request.status == 'draft' || request.status == 'submitted');
	};

	$scope.assignable = function(request) {
		return (request.status == 'submitted');
	};

	$scope.deleteable = function(request) {
		return (request.status == 'draft' || request.status == 'submitted' || request.status == 'rejected');
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
						{ id: $scope.DRAFT, name: 'Save as Draft' },
						{ id: $scope.SUBMITTED, name: 'Submit for Review / Deployment' }
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
			}
			// todo jeremy: this needs to call the api to update ds request status
			deliveryServiceRequestService.updateDeliveryServiceRequest(request.id, request).
				then(function() {
					$scope.refresh();
				});
		}, function () {
			// do nothing
		});
	};

	$scope.rejectRequest = function(request, $event) {
		$event.stopPropagation(); // this kills the click event so it doesn't trigger anything else
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
		modalInstance.result.then(function() {
			request.assigneeId = userModel.user.id;
			request.status = 'rejected';
			deliveryServiceRequestService.updateDeliveryServiceRequest(request.id, request).
				then(function() {
					$scope.refresh();
				});
		}, function () {
			// do nothing
		});
	};

	$scope.completeRequest = function(request, $event) {
		$event.stopPropagation(); // this kills the click event so it doesn't trigger anything else

		if (request.assigneeId != userModel.user.id) {
			messageModel.setMessages([ { level: 'error', text: 'Only the Assignee can mark a delivery service request as complete' } ], false);
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
		modalInstance.result.then(function() {
			request.status = 'complete';
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
		modalInstance.result.then(function() {
			deliveryServiceRequestService.deleteDeliveryServiceRequest(request.id, false).
				then(function() {
					$scope.refresh();
				});
		}, function () {
			// do nothing
		});
	};

	$scope.createRequest = function() {
		var CREATE = 1,
			UPDATE = 2,
			DELETE = 3;

		var params = {
			title: 'Create Delivery Service Request',
			message: 'What kind of delivery service request would you like to create?'
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
						{ id: CREATE, name: 'A request for a new delivery service' },
						{ id: UPDATE, name: 'A request to update an existing delivery service' },
						{ id: DELETE, name: 'A request to delete an existing delivery service' }
					];
				}
			}
		});
		modalInstance.result.then(function(action) {
			var params,
				modalInstance;

			if (action.id == CREATE) {
				params = {
					title: 'Create Delivery Service',
					message: "Please select a content routing category"
				};
				modalInstance = $uibModal.open({
					templateUrl: 'common/modules/dialog/select/dialog.select.tpl.html',
					controller: 'DialogSelectController',
					size: 'md',
					resolve: {
						params: function () {
							return params;
						},
						collection: function() {
							// the following represent the 4 categories of delivery services
							// the ids are arbitrary but the dialog.select dropdown needs them
							return [
								{ id: 1, name: 'ANY_MAP' },
								{ id: 2, name: 'DNS' },
								{ id: 3, name: 'HTTP' },
								{ id: 4, name: 'STEERING' }
							];
						}
					}
				});
				modalInstance.result.then(function(type) {
					var path = '/delivery-services/new?type=' + type.name;
					locationUtils.navigateToPath(path);
				}, function () {
					// do nothing on cancel
				});
			} else if (action.id == UPDATE) {
				params = {
					title: 'Update Delivery Service',
					message: "Please select a delivery service to update",
					labelFunction: function (item) {
						return item['xmlId']
					}
				};
				modalInstance = $uibModal.open({
					templateUrl: 'common/modules/dialog/select/dialog.select.tpl.html',
					controller: 'DialogSelectController',
					size: 'md',
					resolve: {
						params: function () {
							return params;
						},
						collection: function (deliveryServiceService) {
							return deliveryServiceService.getDeliveryServices();
						}
					}
				});
				modalInstance.result.then(function (ds) {
					locationUtils.navigateToPath('/delivery-services/' + ds.id + '?type=' + ds.type);
				}, function () {
					// do nothing on cancel
				});
			} else if (action.id == DELETE) {
				params = {
					title: 'Delete Delivery Service',
					message: "Please select a delivery service to delete",
					labelFunction: function (item) {
						return item['xmlId']
					}
				};
				modalInstance = $uibModal.open({
					templateUrl: 'common/modules/dialog/select/dialog.select.tpl.html',
					controller: 'DialogSelectController',
					size: 'md',
					resolve: {
						params: function () {
							return params;
						},
						collection: function(deliveryServiceService) {
							return deliveryServiceService.getDeliveryServices();
						}
					}
				});
				modalInstance.result.then(function(ds) {
					createDeliveryServiceDeleteRequest(ds);
				}, function () {
					// do nothing on cancel
				});
			}

		}, function () {
			// do nothing on cancel
		});
	};

	$scope.editDeliveryServiceRequest = function(request) {
		var path = '/delivery-service-requests/' + request.id + '?type=';
		typeService.getType(request.deliveryService.typeId)
			.then(function(result) {
				path += result.name;
				locationUtils.navigateToPath(path);
			});
	};

	$scope.fulfillRequest = function(request, $event) {
		$event.stopPropagation(); // this kills the click event so it doesn't trigger anything else
		var path = '/delivery-service-requests/' + request.id + '?type=';
		typeService.getType(request.deliveryService.typeId)
			.then(function(result) {
				path += result.name;
				locationUtils.navigateToPath(path);
			});
	};

	angular.element(document).ready(function () {
		var dsRequestsTable = $('#dsRequestsTable').dataTable({
			"paging": false,
			"dom": '<"filter-checkbox">frtip',
			"columnDefs": [
				{ 'orderable': false, 'targets': 6 }
			],
			"aaSorting": []
		});
		$('div.filter-checkbox').html('<input id="showClosed" type="checkbox"><label for="showClosed">&nbsp;&nbsp;Show Closed</label>');

		$('#showClosed').click(function() {
			var checked = $('#showClosed').is(':checked');
			localStorage.setItem('showClosed', checked);
			if (checked) {
				dsRequestsTable.fnFilter('', 2, true, false);
			} else {
				dsRequestsTable.fnFilter('draft|submitted|pending', 2, true, false);
			}
		});

		if (localStorage.showClosed == 'true') {
			$('#showClosed').attr('checked', true).triggerHandler('click');
		}

	});

};

TableDeliveryServicesRequestsController.$inject = ['dsRequests', '$scope', '$state', '$uibModal', '$anchorScroll', '$q', 'dateUtils', 'locationUtils', 'typeService', 'deliveryServiceService', 'deliveryServiceRequestService', 'messageModel', 'userModel'];
module.exports = TableDeliveryServicesRequestsController;
