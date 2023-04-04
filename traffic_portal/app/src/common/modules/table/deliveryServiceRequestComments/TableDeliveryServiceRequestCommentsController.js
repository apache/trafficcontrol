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

/** @typedef {import("angular")} angular */

/**
 * @param {import("../../../api/DeliveryServiceRequestService").DeliveryServiceRequest} request
 * @param {*} $scope
 * @param {*} $stateParams
 * @param {import("../../../service/utils/angular.ui.bootstrap").IModalService} $uibModal
 * @param {import("angular").IAnchorScrollService} $anchorScroll
 * @param {import("../../../service/utils/DateUtils")} dateUtils
 * @param {import("../../../service/utils/LocationUtils")} locationUtils
 * @param {import("../../../api/DeliveryServiceRequestService")} deliveryServiceRequestService
 * @param {import("../../../models/MessageModel")} messageModel
 */
var TableDeliveryServicesRequestsController = function (request, $scope, $stateParams, $uibModal, $anchorScroll, dateUtils, locationUtils, deliveryServiceRequestService, messageModel) {

	$scope.request = request[0];
	$scope.type = $stateParams.dsType;
	$scope.defaultParams = {
		placeholder: '',
		text: null,
		buttonText: '',
		type: 'default',
		comment: null
	};
	$scope.params = angular.copy($scope.defaultParams);

	$scope.getComments = function () {
		deliveryServiceRequestService.getDeliveryServiceRequestComments({
			deliveryServiceRequestId: $stateParams.deliveryServiceRequestId,
			orderby: 'id'
		}).then(
			function (comments) {
				$scope.comments = comments;
			}
		);
	};

	$scope.createComment = function () {
		$scope.params = {
			placeholder: 'Enter your new comment',
			text: null,
			buttonText: 'Create Comment',
			callback: $scope.submitComment,
			type: 'add'
		};
	};

	$scope.editComment = function (comment) {
		$scope.params = {
			placeholder: '',
			text: comment.value,
			buttonText: 'Update Comment',
			type: 'edit',
			comment: comment
		};
		$anchorScroll();
	};

	$scope.submitComment = function () {
		switch ($scope.params.type) {
			case 'add' :
				var comment = {
					deliveryServiceRequestId: $scope.request.id,
					value: $scope.params.text
				};
				deliveryServiceRequestService.createDeliveryServiceRequestComment(comment).then(
					function () {
						messageModel.setMessages([{
							level: 'success',
							text: 'Delivery service request comment created'
						}], false);
						$scope.updateCommentsView();
					}
				);
				break;

			case 'edit' :
				var editedComment = {
					id: $scope.params.comment.id,
					deliveryServiceRequestId: $scope.params.comment.deliveryServiceRequestId,
					value: $scope.params.text
				};
				deliveryServiceRequestService.updateDeliveryServiceRequestComment(editedComment).then(
					function () {
						messageModel.setMessages([{
							level: 'success',
							text: 'Delivery service request comment updated'
						}], false);
						$scope.updateCommentsView();
					}
				);
				break;
		}
	};

	$scope.updateCommentsView = function () {
		$scope.getComments();
		$scope.params = angular.copy($scope.defaultParams);
	};

	$scope.deleteComment = function (comment, $event) {
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
		modalInstance.result.then(function () {
			deliveryServiceRequestService.deleteDeliveryServiceRequestComment(comment).then(
				function () {
					messageModel.setMessages([{
						level: 'success',
						text: 'Delivery service request comment deleted'
					}], false);
					$scope.getComments();
				}
			);
		}, function () {
			// do nothing
		});
	};

	$scope.cancel = function () {
		$scope.params = angular.copy($scope.defaultParams);
	};

	$scope.getRelativeTime = dateUtils.getRelativeTime;

	$scope.navigateToPath = (path, unsavedChanges) => locationUtils.navigateToPath(path, unsavedChanges);

	$scope.getComments();
};

TableDeliveryServicesRequestsController.$inject = ['request', '$scope', '$stateParams', '$uibModal', '$anchorScroll', 'dateUtils', 'locationUtils', 'deliveryServiceRequestService', 'messageModel'];
module.exports = TableDeliveryServicesRequestsController;
