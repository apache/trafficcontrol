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

/**
 * This is the controller for the form used to modify a Target of a Steering
 * Delivery Service.
 *
 * @param {import("../../../../api/DeliveryServiceService").DeliveryService} deliveryService
 * @param {import("../../../../api/DeliveryServiceService").SteeringTarget[]} currentTargets
 * @param {import("../../../../api/DeliveryServiceService").SteeringTarget} target
 * @param {*} $scope
 * @param {import("angular").IControllerService} $controller
 * @param {*} $uibModal
 * @param {import("../../../../api/DeliveryServiceService")} deliveryServiceService
 */
var FormEditDeliveryServiceTargetController = function(deliveryService, currentTargets, target, $scope, $controller, $uibModal, deliveryServiceService) {

	// extends the FormDeliveryServiceTargetController to inherit common methods
	angular.extend(this, $controller('FormDeliveryServiceTargetController', { deliveryService: deliveryService, currentTargets: currentTargets, target: target, $scope: $scope }));

	var deleteTarget = function(target) {
		deliveryServiceService.deleteDeliveryServiceTarget(target.deliveryServiceId, target.targetId)
			.then(function() {

			});
	};

	$scope.settings = {
		isNew: false,
		saveLabel: 'Update'
	};

	$scope.targetName = angular.copy(target.target);

	$scope.save = function(dsId, targetId, target) {
		deliveryServiceService.updateDeliveryServiceTarget(dsId, targetId, target);
	};

	$scope.confirmDelete = function(target) {
		var params = {
			title: 'Delete Target',
			message: 'Are you sure you want to delete the steering target?'
		};
		var modalInstance = $uibModal.open({
			templateUrl: 'common/modules/dialog/confirm/dialog.confirm.tpl.html',
			controller: 'DialogConfirmController',
			size: 'sm',
			resolve: {
				params: function () {
					return params;
				}
			}
		});
		modalInstance.result.then(function() {
			deleteTarget(target);
		}, function () {
			// do nothing
		});
	};


};

FormEditDeliveryServiceTargetController.$inject = ['deliveryService', 'currentTargets', 'target', '$scope', '$controller', '$uibModal', 'deliveryServiceService'];
module.exports = FormEditDeliveryServiceTargetController;
