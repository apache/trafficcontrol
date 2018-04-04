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

var TableDeliveryServicesRequestsController = function(dsRequests, $scope, $state, $uibModal, $anchorScroll, $q, $location, dateUtils, locationUtils, typeService, deliveryServiceService, deliveryServiceRequestService, messageModel, userModel) {

	var createDeliveryServiceDeleteRequest = function(deliveryService) {
		var params = {
			title: "Delivery Service Delete Request",
			message: 'All delivery service deletions must be reviewed.'
		};
		var modalInstance = $uibModal.open({
			templateUrl: 'common/modules/dialog/deliveryServiceRequest/dialog.deliveryServiceRequest.tpl.html',
			controller: 'DialogDeliveryServiceRequestController',
			size: 'md',
			resolve: {
				params: function () {
					return params;
				},
				statuses: function() {
					return [
						{ id: $scope.DRAFT, name: 'Save as Draft' },
						{ id: $scope.SUBMITTED, name: 'Submit for Review and Deployment' }
					];
				}
			}
		});
		modalInstance.result.then(function(options) {
			var dsRequest = {
				changeType: 'delete',
				status: (options.status.id == $scope.SUBMITTED) ? 'submitted' : 'draft',
				deliveryService: deliveryService
			};
			deliveryServiceRequestService.createDeliveryServiceRequest(dsRequest).
				then(
					function(response) {
						var comment = {
							deliveryServiceRequestId: response.id,
							value: options.comment
						};
						deliveryServiceRequestService.createDeliveryServiceRequestComment(comment).
							then(
								function() {
									messageModel.setMessages([ { level: 'success', text: 'Created request to ' + dsRequest.changeType + ' the ' + dsRequest.deliveryService.xmlId + ' delivery service' } ], false);
									$scope.refresh();
								}
							);
					}
				);
		}, function () {
			// do nothing
		});
	};

	var createComment = function(request, placeholder) {
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
		modalInstance.result.then(function(commentValue) {
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

	$scope.closed = function(request) {
		return (request.status == 'rejected' || request.status == 'complete');
	};

	$scope.compareRequests = function() {
		var params = {
			title: 'Compare Delivery Service Requests',
			message: "Please select 2 delivery service requests to compare",
			labelFunction: function (item) {
				return item['deliveryService']['xmlId'] + ' ' + item['changeType'] + ' (' + item['author'] + ' created on ' + item['createdAt'] + ')'
			}
		};
		var modalInstance = $uibModal.open({
			templateUrl: 'common/modules/dialog/compare/dialog.compare.tpl.html',
			controller: 'DialogCompareController',
			size: 'lg',
			resolve: {
				params: function () {
					return params;
				},
				collection: function(deliveryServiceRequestService) {
					return deliveryServiceRequestService.getDeliveryServiceRequests();
				}
			}
		});
		modalInstance.result.then(function(requests) {
			$location.path($location.path() + '/compare/' + requests[0].id + '/' + requests[1].id);
		}, function () {
			// do nothing
		});
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
			var assigneeId = (assign) ? userModel.user.id : null;
			deliveryServiceRequestService.assignDeliveryServiceRequest(request.id, assigneeId).
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
			var status = (action.id == $scope.DRAFT) ? 'draft' : 'submitted';
			deliveryServiceRequestService.updateDeliveryServiceRequestStatus(request.id, status).
				then(function() {
					$scope.refresh();
				});
		}, function () {
			// do nothing
		});
	};

	$scope.viewComments = function(request, $event) {
		$event.stopPropagation(); // this kills the click event so it doesn't trigger anything else
		$location.path($location.path() + '/' + request.id + '/comments');
	};

	$scope.rejectRequest = function(request, $event) {
		$event.stopPropagation(); // this kills the click event so it doesn't trigger anything else

		if (request.assigneeId != userModel.user.id) {
			messageModel.setMessages([ { level: 'error', text: 'Only the assignee can mark a delivery service request as rejected' } ], false);
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
		modalInstance.result.then(function() {
			deliveryServiceRequestService.updateDeliveryServiceRequestStatus(request.id, 'rejected').
				then(
					function() {
						$scope.refresh();
						createComment(request, 'Enter rejection reason...');
					});
		}, function () {
			// do nothing
		});
	};

	$scope.completeRequest = function(request, $event) {
		$event.stopPropagation(); // this kills the click event so it doesn't trigger anything else

		if (request.assigneeId != userModel.user.id) {
			messageModel.setMessages([ { level: 'error', text: 'Only the assignee can mark a delivery service request as complete' } ], false);
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
			deliveryServiceRequestService.updateDeliveryServiceRequestStatus(request.id, 'complete').
				then(function() {
					$scope.refresh();
					createComment(request, 'Enter comment...');
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
			deliveryServiceRequestService.deleteDeliveryServiceRequest(request.id).
				then(function() {
					messageModel.setMessages([ { level: 'success', text: 'Delivery service request deleted' } ], false);
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
				{ 'orderable': false, 'targets': 7 }
			],
			"aaSorting": []
		});
		$('div.filter-checkbox').html('<input id="showClosed" type="checkbox"><label for="showClosed">&nbsp;&nbsp;Show Closed</label>');

		// only show "open" ds requests on render
		dsRequestsTable.fnFilter('draft|submitted|pending', 2, true, false);

		$('#showClosed').click(function() {
			var checked = $('#showClosed').is(':checked');
			if (checked) {
				dsRequestsTable.fnFilter('', 2, true, false);
			} else {
				dsRequestsTable.fnFilter('draft|submitted|pending', 2, true, false);
			}
		});

	});

};

TableDeliveryServicesRequestsController.$inject = ['dsRequests', '$scope', '$state', '$uibModal', '$anchorScroll', '$q', '$location', 'dateUtils', 'locationUtils', 'typeService', 'deliveryServiceService', 'deliveryServiceRequestService', 'messageModel', 'userModel'];
module.exports = TableDeliveryServicesRequestsController;
