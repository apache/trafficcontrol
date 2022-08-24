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

var TableServerServerCapabilitiesController = function(server, serverCapabilities, $scope, $state, $uibModal, locationUtils, serverUtils, serverService, messageModel, serverCapabilityService) {

	$scope.server = server[0];

	$scope.serverCapabilities = serverCapabilities;

	$scope.contextMenuItems = [
		{
			text: 'Remove Server Capability',
			click: function ($itemScope) {
				$scope.confirmRemoveCapability($itemScope.sc.serverCapability);
			}
		}
	];

	$scope.selectSCs = function () {
		var modalInstance = $uibModal.open({
			templateUrl: 'common/modules/table/serverServerCapabilities/table.assignServerSCs.tpl.html',
			controller: 'TableAssignServerSCsController',
			size: 'md',
			resolve: {
				server: function() {
					return server;
				},
				serverCapabilities: function(serverCapabilityService) {
					return serverCapabilityService.getServerCapabilities();
				},
				assignedSCs: function() {
					return serverCapabilities
				}
			}
		});
		modalInstance.result.then(function(selectedSCs) {
			serverCapabilityService.assignSCsServer(server[0].id, selectedSCs)
				.then(
					function() {
						$scope.refresh();
					}
				);
		}, function () {
			// do nothing
		});
	};

	$scope.addServerCapability = function() {
		const params = {
			title: 'Add Server Capability',
			message: "Please select a capability to add to this server",
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
			serverService.addServerCapability($scope.server.id, serverCapability.name)
				.then(
					function(result) {
						messageModel.setMessages(result.alerts, false);
						$scope.refresh(); // refresh the profile parameters table
					}
				);
		});
	};

	$scope.confirmRemoveCapability = function(serverCapability, $event) {
		if ($event) {
			$event.stopPropagation(); // this kills the click event so it doesn't trigger anything else
		}
		const params = {
			title: 'Remove Capability from Server?',
			message: 'Are you sure you want to remove the ' + serverCapability + ' capability from this server?'
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
			serverService.removeServerCapability($scope.server.id, serverCapability)
				.then(
					function(result) {
						messageModel.setMessages(result.alerts, false);
						$scope.refresh(); // refresh the profile parameters table
					}
				);
		});
	};

	$scope.editServerCapability = function(capabilityName) {
		locationUtils.navigateToPath('/server-capabilities/' + capabilityName);
	};

	$scope.refresh = function() {
		$state.reload(); // reloads all the resolves for the view
	};

	$scope.navigateToPath = locationUtils.navigateToPath;

	angular.element(document).ready(function () {
		$('#serverCapabilitiesTable').dataTable({
			"lengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
			"iDisplayLength": 25,
			"columnDefs": [
				{ "width": "5%", "targets": 1 },
				{ 'orderable': false, 'targets': 1 }
			],
			"aaSorting": []
		});
	});

	$scope.isCache = serverUtils.isCache;
};

TableServerServerCapabilitiesController.$inject = ['server', 'serverCapabilities', '$scope', '$state', '$uibModal', 'locationUtils', 'serverUtils', 'serverService', 'messageModel', 'serverCapabilityService'];
module.exports = TableServerServerCapabilitiesController;
