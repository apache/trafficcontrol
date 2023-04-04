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

/** @typedef {import("jquery")} */

/**
 * @param {import("../../../api/DeliveryServiceService").DeliveryService} deliveryService
 * @param {import("../../../api/DeliveryServiceService").SteeringTarget[]} targets
 * @param {*} $scope
 * @param {*} $state
 * @param {import("../../../service/utils/LocationUtils")} locationUtils
 */
var TableDeliveryServiceTargetsController = function(deliveryService, targets, $scope, $state, locationUtils) {

	$scope.deliveryService = deliveryService;

	$scope.targets = targets;

	$scope.editTarget = function(dsId, targetId) {
		locationUtils.navigateToPath('/delivery-services/' + dsId + '/targets/' + targetId);
	};

	$scope.createTarget = function(dsId) {
		locationUtils.navigateToPath('/delivery-services/' + dsId + '/targets/new');
	};

	$scope.refresh = function() {
		$state.reload(); // reloads all the resolves for the view
	};

	$scope.navigateToPath = (path, unsavedChanges) => locationUtils.navigateToPath(path, unsavedChanges);

	angular.element(document).ready(function () {
		// Datatables...
		// @ts-ignore
		$('#targetsTable').dataTable({
			"aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
			"iDisplayLength": 25,
			"aaSorting": []
		});
	});

};

TableDeliveryServiceTargetsController.$inject = ['deliveryService', 'targets', '$scope', '$state', 'locationUtils'];
module.exports = TableDeliveryServiceTargetsController;
