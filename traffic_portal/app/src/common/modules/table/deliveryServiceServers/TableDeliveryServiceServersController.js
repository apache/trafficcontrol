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

var TableDeliveryServiceServersController = function(deliveryService, servers, $scope, $state, $uibModal, locationUtils, serverUtils, deliveryServiceService) {

	$scope.deliveryService = deliveryService;

	$scope.servers = servers;

	$scope.removeServer = function(dsId, serverId) {
		deliveryServiceService.deleteDeliveryServiceServer(dsId, serverId)
			.then(
				function() {
					$scope.refresh();
				}
			);
	};

	$scope.refresh = function() {
		$state.reload(); // reloads all the resolves for the view
	};

	$scope.selectServers = function() {
		var modalInstance = $uibModal.open({
			templateUrl: 'common/modules/table/deliveryServiceServers/table.dsServersUnassigned.tpl.html',
			controller: 'TableDSServersUnassignedController',
			size: 'lg',
			resolve: {
				deliveryService: function() {
					return deliveryService;
				},
				eligibleServers: function(serverService) {
					return serverService.getEligibleDeliveryServiceServers(deliveryService.id);
				},
				assignedServers: function() {
					return servers;
				}
			}
		});
		modalInstance.result.then(function(selectedServerIds) {
			deliveryServiceService.assignDeliveryServiceServers(deliveryService.id, selectedServerIds)
				.then(
					function() {
						$scope.refresh();
					}
				);
		}, function () {
			// do nothing
		});
	};


	$scope.isOffline = serverUtils.isOffline;

	$scope.offlineReason = serverUtils.offlineReason;

	$scope.navigateToPath = locationUtils.navigateToPath;

	angular.element(document).ready(function () {
		$('#serversTable').dataTable({
			"aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
			"iDisplayLength": 25,
			"aaSorting": []
		});
	});

};

TableDeliveryServiceServersController.$inject = ['deliveryService', 'servers', '$scope', '$state', '$uibModal', 'locationUtils', 'serverUtils', 'deliveryServiceService'];
module.exports = TableDeliveryServiceServersController;
