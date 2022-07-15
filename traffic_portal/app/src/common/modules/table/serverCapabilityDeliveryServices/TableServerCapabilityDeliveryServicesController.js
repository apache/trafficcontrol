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

var TableServerCapabilityDeliveryServicesController = function(serverCapability, deliveryServices, $scope, $state, $uibModal, $window, locationUtils, deliveryServiceService, messageModel) {

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
			text: 'Remove Capability from Delivery Service',
			click: function ($itemScope) {
				$scope.confirmRemoveCapability($itemScope.ds);
			}
		},
		null, // Divider
		{
			text: 'Edit Delivery Service',
			click: function ($itemScope) {
				$scope.editDeliveryService($itemScope.ds);
			}
		},
		{
			text: 'Manage Required Server Capabilities',
			click: function ($itemScope) {
				locationUtils.navigateToPath('/delivery-services/' + $itemScope.ds.deliveryServiceID + '/required-server-capabilities');
			}
		}
	];

	$scope.navigateToPath = locationUtils.navigateToPath;

	$scope.confirmRemoveCapability = function(ds, $event) {
		if ($event) {
			$event.stopPropagation(); // this kills the click event so it doesn't trigger anything else
		}

		const params = {
			title: 'Remove Required Server Capability from Delivery Service?',
			message: 'Are you sure you want to remove the ' + serverCapability.name + ' server capability requirement from the ' + ds.xmlID + ' delivery service?'
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
			removeCapability(ds.deliveryServiceID);
		});
	};

	$scope.editDeliveryService = function(ds) {
		deliveryServiceService.getDeliveryService(ds.deliveryServiceID)
			.then(function(result) {
				let path = '/delivery-services/' + result.id + '?dsType=' + result.type;
				locationUtils.navigateToPath(path);
			});
	};

	$scope.refresh = function() {
		$state.reload(); // reloads all the resolves for the view
	};

	angular.element(document).ready(function () {
		$('#serverCapabilityDeliveryServicesTable').dataTable({
			"lengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
			"iDisplayLength": 25,
			"aaSorting": []
		});
	});

};

TableServerCapabilityDeliveryServicesController.$inject = ['serverCapability', 'deliveryServices', '$scope', '$state', '$uibModal', '$window', 'locationUtils', 'deliveryServiceService', 'messageModel'];
module.exports = TableServerCapabilityDeliveryServicesController;
