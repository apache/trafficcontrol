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

var TableDeliveryServiceServersController = function(deliveryService, servers, $controller, $scope, $uibModal, deliveryServiceService, serverUtils) {

	// extends the TableParentServersController to inherit common methods
	angular.extend(this, $controller('TableParentServersController', { servers: servers, $scope: $scope }));

	let dsServersTable;

	var removeServer = function(serverId) {
		deliveryServiceService.deleteDeliveryServiceServer($scope.deliveryService.id, serverId)
			.then(
				function() {
					$scope.refresh();
				}
			);
	};

	$scope.deliveryService = deliveryService;

	$scope.isEdge = serverUtils.isEdge;

	$scope.isOrigin = serverUtils.isOrigin;

	// adds some items to the base servers context menu
	$scope.contextMenuItems.splice(2, 0,
		{
			text: 'Unlink Server from Delivery Service',
			hasBottomDivider: function() {
				return true;
			},
			click: function ($itemScope, evt) {
				$scope.confirmRemoveServer($itemScope.s, evt);
			}
		}
	);

	$scope.selectServers = function() {
		var modalInstance = $uibModal.open({
			templateUrl: 'common/modules/table/deliveryServiceServers/table.assignDSServers.tpl.html',
			controller: 'TableAssignDSServersController',
			size: 'lg',
			resolve: {
				deliveryService: function() {
					return deliveryService;
				},
				servers: function(serverService) {
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

	$scope.confirmRemoveServer = function(server, $event) {
		if ($event) {
			$event.stopPropagation(); // this kills the click event so it doesn't trigger anything else
		}

		var params = {
			title: 'Remove Server from Delivery Service?',
			message: 'Are you sure you want to remove ' + server.hostName + ' from this delivery service?'
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
			removeServer(server.id);
		}, function () {
			// do nothing
		});
	};

	$scope.toggleVisibility = function(colName) {
		const col = dsServersTable.column(colName + ':name');
		col.visible(!col.visible());
		dsServersTable.rows().invalidate().draw();
	};

	$scope.columnFilterFn = function(column) {
		if (column.name === 'Action') {
			return false;
		}
		return true;
	};

	angular.element(document).ready(function () {
		dsServersTable = $('#dsServersTable').DataTable({
			"lengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
			"iDisplayLength": 25,
			"aaSorting": [],
			"columnDefs": [
				{ 'orderable': false, 'targets': 32 }
			],
			"columns": $scope.columns.concat([{ "name": "Action", "visible": true, "searchable": false }]),
			"initComplete": function(settings, json) {
				try {
					// need to create the show/hide column checkboxes and bind to the current visibility
					$scope.columns = JSON.parse(localStorage.getItem('DataTables_dsServersTable_/')).columns;
				} catch (e) {
					console.error("Failure to retrieve required column info from localStorage (key=DataTables_dsServersTable_/):", e);
				}
			}
		});
	});

};

TableDeliveryServiceServersController.$inject = ['deliveryService', 'servers', '$controller', '$scope', '$uibModal', 'deliveryServiceService', 'serverUtils'];
module.exports = TableDeliveryServiceServersController;
