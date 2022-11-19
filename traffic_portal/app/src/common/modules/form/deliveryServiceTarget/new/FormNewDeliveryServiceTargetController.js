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
 * This is the controller for the form used to add a new Target to a Steering
 * Delivery Service.
 *
 * @param {import("../../../../api/DeliveryServiceService").DeliveryService} deliveryService
 * @param {import("../../../../api/DeliveryServiceService").SteeringTarget[]} currentTargets
 * @param {import("../../../../api/DeliveryServiceService").SteeringTarget} target
 * @param {*} $scope
 * @param {import("angular").IControllerService} $controller
 * @param {import("../../../../api/DeliveryServiceService")} deliveryServiceService
 */
var FormNewDeliveryServiceTargetController = function(deliveryService, currentTargets, target, $scope, $controller, deliveryServiceService) {

	// extends the FormDeliveryServiceTargetController to inherit common methods
	angular.extend(this, $controller('FormDeliveryServiceTargetController', { deliveryService: deliveryService, currentTargets: currentTargets, target: target, $scope: $scope }));

	$scope.targetName = 'New';

	$scope.settings = {
		isNew: true,
		saveLabel: 'Create'
	};

	$scope.save = function(dsId, targetId, target) {
		deliveryServiceService.createDeliveryServiceTarget(dsId, target);
	};

};

FormNewDeliveryServiceTargetController.$inject = ['deliveryService', 'currentTargets', 'target', '$scope', '$controller', 'deliveryServiceService'];
module.exports = FormNewDeliveryServiceTargetController;
