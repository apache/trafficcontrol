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

var TableCacheGroupsServersController = function(cacheGroup, servers, $controller, $scope, $state, $uibModal, cacheGroupService) {

	// extends the TableServersController to inherit common methods
	angular.extend(this, $controller('TableServersController', { servers: servers, $scope: $scope }));

	let cacheGroupServersTable;

	$scope.cacheGroup = cacheGroup;

	var queueCacheGroupServerUpdates = function(cacheGroup, cdnId) {
		cacheGroupService.queueServerUpdates(cacheGroup.id, cdnId)
			.then(
				function() {
					$scope.refresh();
				}
			);
	};

	var clearCacheGroupServerUpdates = function(cacheGroup, cdnId) {
		cacheGroupService.clearServerUpdates(cacheGroup.id, cdnId)
			.then(
				function() {
					$scope.refresh();
				}
			);
	};

	$scope.confirmCacheGroupQueueServerUpdates = function(cacheGroup) {
		var params = {
			title: 'Queue Server Updates: ' + cacheGroup.name,
			message: "Please select a CDN"
		};
		var modalInstance = $uibModal.open({
			templateUrl: 'common/modules/dialog/select/dialog.select.tpl.html',
			controller: 'DialogSelectController',
			size: 'md',
			resolve: {
				params: function () {
					return params;
				},
				collection: function(cdnService) {
					return cdnService.getCDNs();
				}
			}
		});
		modalInstance.result.then(function(cdn) {
			queueCacheGroupServerUpdates(cacheGroup, cdn.id);
		}, function () {
			// do nothing
		});
	};

	$scope.confirmCacheGroupClearServerUpdates = function(cacheGroup) {
		var params = {
			title: 'Clear Server Updates: ' + cacheGroup.name,
			message: "Please select a CDN"
		};
		var modalInstance = $uibModal.open({
			templateUrl: 'common/modules/dialog/select/dialog.select.tpl.html',
			controller: 'DialogSelectController',
			size: 'md',
			resolve: {
				params: function () {
					return params;
				},
				collection: function(cdnService) {
					return cdnService.getCDNs();
				}
			}
		});
		modalInstance.result.then(function(cdn) {
			clearCacheGroupServerUpdates(cacheGroup, cdn.id);
		}, function () {
			// do nothing
		});
	};

	$scope.toggleVisibility = function(colName) {
		const col = cacheGroupServersTable.column(colName + ':name');
		col.visible(!col.visible());
		cacheGroupServersTable.rows().invalidate().draw();
	};

	angular.element(document).ready(function () {
		cacheGroupServersTable = $('#cacheGroupServersTable').DataTable({
			"lengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
			"iDisplayLength": 25,
			"aaSorting": [],
			"columns": $scope.columns,
			"colReorder": {
				realtime: false
			},
			"initComplete": function(settings, json) {
				try {
					// need to create the show/hide column checkboxes and bind to the current visibility
					$scope.columns = JSON.parse(localStorage.getItem('DataTables_cacheGroupServersTable_/')).columns;
				} catch (e) {
					console.error("Failure to retrieve required column info from localStorage (key=DataTables_cacheGroupServersTable_/):", e);
				}
			}
		});
	});

};

TableCacheGroupsServersController.$inject = ['cacheGroup', 'servers', '$controller', '$scope', '$state', '$uibModal', 'cacheGroupService'];
module.exports = TableCacheGroupsServersController;
