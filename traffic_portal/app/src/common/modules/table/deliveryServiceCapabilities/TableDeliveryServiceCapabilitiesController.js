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
 * The controller for the table that lists the server capabilities required by a
 * Delivery Service.
 *
 * @param {import("../../../api/DeliveryServiceService").DeliveryService & {id: number}} deliveryService
 * @param {string[]} requiredCapabilities
 * @param {*} $scope
 * @param {*} $state
 * @param {{open: ({}) => {result: Promise<*>}}} $uibModal
 * @param {import("../../../service/utils/LocationUtils")} locationUtils
 * @param {import("../../../api/DeliveryServiceService")} deliveryServiceService
 * @param {import("../../../models/MessageModel")} messageModel
 */
var TableDeliveryServiceCapabilitiesController = function(deliveryService, requiredCapabilities, $scope, $state, $uibModal, locationUtils, deliveryServiceService, messageModel) {

	$scope.deliveryService = deliveryService;

	$scope.requiredCapabilities = requiredCapabilities;

	$scope.contextMenuItems = [
		{
			text: 'Remove Required Server Capability',
			click: function ($itemScope) {
				$scope.confirmRemoveCapability($itemScope.rq.requiredCapability);
			}
		}
	];

	$scope.addDeliveryServiceCapability = function() {
		const params = {
			title: 'Add Required Server Capability',
			message: "Please select a server capability required by this delivery service",
			key: "name"
		};
		const modalInstance = $uibModal.open({
			templateUrl: 'common/modules/dialog/select/dialog.select.tpl.html',
			controller: 'DialogSelectController',
			size: 'md',
			resolve: {
				params: function () {
					return params;
				},
				collection: function(serverCapabilityService) {
					return serverCapabilityService.getServerCapabilities();
				}
			}
		});
		modalInstance.result.then(function(serverCapability) {
			deliveryServiceService.addServerCapability(deliveryService.id, serverCapability.name)
				.then(
					function(result) {
						messageModel.setMessages(result.alerts, false);
						$scope.refresh(); // refresh the table
					}
				);
		});
	};

	$scope.confirmRemoveCapability = function(requiredCapability, $event) {
		if ($event) {
			$event.stopPropagation(); // this kills the click event so it doesn't trigger anything else
		}
		const params = {
			title: 'Remove Required Server Capability from Delivery Service?',
			message: 'Are you sure you want to remove the ' + requiredCapability + ' server capability requirement from this delivery service?'
		};
		const modalInstance = $uibModal.open({
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
			deliveryServiceService.removeServerCapability(deliveryService.id, requiredCapability)
				.then(
					function(result) {
						messageModel.setMessages(result.alerts, false);
						$scope.refresh(); // refresh the table
					}
				);
		});
	};

	$scope.editServerCapability = function(capabilityName) {
		locationUtils.navigateToPath('/server-capabilities/edit?name=' + capabilityName);
	};

	$scope.refresh = function() {
		$state.reload(); // reloads all the resolves for the view
	};

	$scope.navigateToPath = locationUtils.navigateToPath;

	angular.element(document).ready(function () {
		// Datatable types don't exist in the project, and they should all be
		// replaced with AG-Grid anyway.
		// @ts-ignore
		$('#deliveryServiceCapabilitiesTable').dataTable({
			"lengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
			"iDisplayLength": 25,
			"columnDefs": [
				{ "width": "5%", "targets": 1 },
				{ 'orderable': false, 'targets': 1 }
			],
			"aaSorting": []
		});
	});

};

TableDeliveryServiceCapabilitiesController.$inject = ['deliveryService', 'requiredCapabilities', '$scope', '$state', '$uibModal', 'locationUtils', 'deliveryServiceService', 'messageModel'];
module.exports = TableDeliveryServiceCapabilitiesController;
