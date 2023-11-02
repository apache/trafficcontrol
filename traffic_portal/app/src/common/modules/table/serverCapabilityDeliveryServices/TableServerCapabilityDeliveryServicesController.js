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
 * The controller for the table that lists the Delivery Services that require a
 * particular Server Capability.
 *
 * @param {{name: string}} serverCapability
 * @param {import("../../../api/DeliveryServiceService").DeliveryService[]} deliveryServices
 * @param {*} $scope
 * @param {*} $state
 * @param {import("../../../service/utils/angular.ui.bootstrap").IModalService} $uibModal
 * @param {import("../../../service/utils/LocationUtils")} locationUtils
 * @param {import("../../../api/DeliveryServiceService")} deliveryServiceService
 * @param {import("../../../models/MessageModel")} messageModel
 */
var TableServerCapabilityDeliveryServicesController = function(serverCapability, deliveryServices, $scope, $state, $uibModal, locationUtils, deliveryServiceService, messageModel) {

	var removeCapability = function(dsId) {
		deliveryServiceService.removeServerCapability(dsId, serverCapability.name)
			.then(
				function(result) {
					messageModel.setMessages(result.alerts, false);
					$scope.refresh();
				}
			);
	};

	$scope.serverCapability = serverCapability;

	$scope.deliveryServices = deliveryServices;

	$scope.contextMenuItems = [
		{
			text: 'Edit Delivery Service',
			click: function ($itemScope) {
				$scope.editDeliveryService($itemScope.ds);
			}
		},
	];

	$scope.navigateToPath = (path, unsavedChanges) => locationUtils.navigateToPath(path, unsavedChanges);

	$scope.editDeliveryService = function(ds) {
		deliveryServiceService.getDeliveryService(ds.id)
			.then(function(result) {
				let path = '/delivery-services/' + result.id + '?dsType=' + result.type;
				locationUtils.navigateToPath(path);
			});
	};

	$scope.refresh = function() {
		$state.reload(); // reloads all the resolves for the view
	};

	angular.element(document).ready(function () {
		// Datatables...
		// @ts-ignore
		$('#serverCapabilityDeliveryServicesTable').dataTable({
			"lengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
			"iDisplayLength": 25,
			"aaSorting": []
		});
	});

};

TableServerCapabilityDeliveryServicesController.$inject = ['serverCapability', 'deliveryServices', '$scope', '$state', '$uibModal', 'locationUtils', 'deliveryServiceService', 'messageModel'];
module.exports = TableServerCapabilityDeliveryServicesController;
