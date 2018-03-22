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

var TableDeliveryServicesRequestsController = function(request, comments, $scope, $state, $stateParams, $uibModal, dateUtils, locationUtils, deliveryServiceRequestService, messageModel) {

	$scope.request = request[0];

	$scope.comments = comments;

	$scope.type = $stateParams.type;

	$scope.refresh = function() {
		$state.reload(); // reloads all the resolves for the view
	};

	$scope.createComment = function() {
		var params = {
			title: 'Add Comment',
			placeholder: "Enter comment...",
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
				deliveryServiceRequestId: $scope.request.id,
				value: commentValue
			};
			deliveryServiceRequestService.createDeliveryServiceRequestComment(comment).
				then(
					function() {
						messageModel.setMessages([ { level: 'success', text: 'Delivery service request comment created' } ], false);
						$scope.refresh();
					}
			);
		}, function () {
			// do nothing
		});
	};

	$scope.editComment = function(comment) {
		var params = {
			title: 'Edit Comment',
			text: comment.value
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
		modalInstance.result.then(function(newValue) {
			var editedComment = {
				id: comment.id,
				deliveryServiceRequestId: comment.deliveryServiceRequestId,
				value: newValue
			};
			deliveryServiceRequestService.updateDeliveryServiceRequestComment(editedComment).
				then(
					function() {
						messageModel.setMessages([ { level: 'success', text: 'Delivery service request comment updated' } ], false);
						$scope.refresh();
					}
				);
		}, function () {
			// do nothing
		});
	};

	$scope.deleteComment = function(comment, $event) {
		$event.stopPropagation(); // this kills the click event so it doesn't trigger anything else
		var params = {
			title: 'Delete Comment',
			message: 'Are you sure you want to delete this comment?'
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
			deliveryServiceRequestService.deleteDeliveryServiceRequestComment(comment).
				then(
					function() {
						messageModel.setMessages([ { level: 'success', text: 'Delivery service request comment deleted' } ], false);
						$scope.refresh();
					}
				);
		}, function () {
			// do nothing
		});
	};

	$scope.getRelativeTime = dateUtils.getRelativeTime;

	$scope.navigateToPath = locationUtils.navigateToPath;

	angular.element(document).ready(function () {
		$('#dsRequestCommentsTable').dataTable({
			"aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
			"iDisplayLength": -1,
			"ordering": false,
			"columnDefs": [
				{ "width": "3%", "targets": 3 }
			]
		});
	});

};

TableDeliveryServicesRequestsController.$inject = ['request', 'comments', '$scope', '$state', '$stateParams', '$uibModal', 'dateUtils', 'locationUtils', 'deliveryServiceRequestService', 'messageModel'];
module.exports = TableDeliveryServicesRequestsController;
