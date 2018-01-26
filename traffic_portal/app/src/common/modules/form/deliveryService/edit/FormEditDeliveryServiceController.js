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

var FormEditDeliveryServiceController = function(deliveryService, type, types, $scope, $state, $controller, $uibModal, locationUtils, deliveryServiceService, deliveryServiceRequestService) {

	// extends the FormDeliveryServiceController to inherit common methods
	angular.extend(this, $controller('FormDeliveryServiceController', { deliveryService: deliveryService, dsCurrent: deliveryService, type: type, types: types, $scope: $scope }));

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
			deliveryServiceRequestService.createDeliveryServiceRequest(dsRequest, true);
		}, function () {
			// do nothing
		});
	};

	$scope.deliveryServiceName = angular.copy(deliveryService.xmlId);

	$scope.settings = {
		isNew: false,
		isRequest: false,
		saveLabel: 'Update',
		deleteLabel: 'Delete'
	};

	$scope.save = function(deliveryService) {
		if ($scope.dsRequestsEnabled) {
			var params = {
				title: "Delivery Service Update Request",
				message: 'All delivery service updates must be reviewed for completeness and accuracy before deployment.<br><br>Please select the status of your delivery service update request.'
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
							{ id: $scope.SUBMITTED, name: 'Submit for Review and Deployment' }
						];
					}
				}
			});
			modalInstance.result.then(function(action) {
				var dsRequest = {
					changeType: 'update',
					status: (action.id == $scope.SUBMITTED) ? 'submitted' : 'draft',
					deliveryService: deliveryService
				};
				deliveryServiceRequestService.createDeliveryServiceRequest(dsRequest, true);
			}, function () {
				// do nothing
			});
		} else {
			deliveryServiceService.updateDeliveryService(deliveryService, false).
				then(function() {
					$state.reload(); // reloads all the resolves for the view
				});
		}
	};

	$scope.confirmDelete = function(deliveryService) {
		if ($scope.dsRequestsEnabled) {
			createDeliveryServiceDeleteRequest(deliveryService);
		} else {
			var params = {
				title: 'Delete Delivery Service: ' + deliveryService.displayName,
				key: deliveryService.xmlId
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
				deliveryServiceService.deleteDeliveryService(deliveryService, true)
					.then(function() {
						locationUtils.navigateToPath('/delivery-services');
					});
			}, function () {
				// do nothing
			});
		}

	};

};

FormEditDeliveryServiceController.$inject = ['deliveryService', 'type', 'types', '$scope', '$state', '$controller', '$uibModal', 'locationUtils', 'deliveryServiceService', 'deliveryServiceRequestService'];
module.exports = FormEditDeliveryServiceController;
