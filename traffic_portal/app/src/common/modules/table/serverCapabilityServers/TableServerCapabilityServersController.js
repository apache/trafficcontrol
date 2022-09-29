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

var TableServerCapabilityServersController = function(serverCapability, servers, $scope, $state, $uibModal, $window, locationUtils, serverService, messageModel, serverCapabilityService) {

	var removeCapability = function(serverId) {
		serverService.removeServerCapability(serverId, serverCapability.name)
			.then(
				function(result) {
					messageModel.setMessages(result.alerts, false);
					$scope.refresh();
				}
			);
	};

	$scope.servers = servers;

	$scope.serverCapability = serverCapability;

	$scope.contextMenuItems = [
		{
			text: 'Open Server in New Tab',
			click: function ($itemScope) {
				$window.open('/#!/servers/' + $itemScope.s.serverId, '_blank');
			}
		},
		null, // Divider
		{
			text: 'Remove Capability from Server',
			click: function ($itemScope) {
				$scope.confirmRemoveCapability($itemScope.s.serverId);
			}
		},
		null, // Divider
		{
			text: 'Edit Server',
			click: function ($itemScope) {
				$scope.editServer($itemScope.s.serverId);
			}
		},
		{
			text: 'Manage Server Capabilities',
			click: function ($itemScope) {
				locationUtils.navigateToPath('/servers/' + $itemScope.s.serverId + '/capabilities');
			}
		}
	];

	$scope.editServer = function(id) {
		locationUtils.navigateToPath('/servers/' + id);
	};

	$scope.selectServers = function () {
		const modalInstance = $uibModal.open({
			templateUrl: 'common/modules/table/serverCapabilityServers/table.assignServersPerCapability.tpl.html',
			controller: 'TableAssignServersPerCapabilityController',
			size: 'md',
			resolve: {
				serverCapability: function() {
					return serverCapability;
				},
				servers: function(serverService) {
					return serverService.getServers();
				},
				assignedServers: function() {
					return servers;
				}
			}
		});
		modalInstance.result.then(function(selectedServers) {
			serverCapabilityService.assignServersPerSC(serverCapability, selectedServers)
				.then(
					function() {
						$scope.refresh();
					}
				);
		}, function () {
			// do nothing
		});
	};

	$scope.confirmRemoveCapability = function(serverId, $event) {
		if ($event) {
			$event.stopPropagation(); // this kills the click event so it doesn't trigger anything else
		}

		const params = {
			title: 'Remove Capability from Server?',
			message: 'Are you sure you want to remove the ' + serverCapability.name + ' capability from this server?'
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
			removeCapability(serverId);
		});
	};

	$scope.refresh = function() {
		$state.reload(); // reloads all the resolves for the view
	};

	angular.element(document).ready(function () {
		$('#serverCapabilityServersTable').dataTable({
			"lengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
			"iDisplayLength": 25,
			"aaSorting": []
		});
	});

};

TableServerCapabilityServersController.$inject = ['serverCapability', 'servers', '$scope', '$state', '$uibModal', '$window', 'locationUtils', 'serverService', 'messageModel', 'serverCapabilityService'];
module.exports = TableServerCapabilityServersController;
