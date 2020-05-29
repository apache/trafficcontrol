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

var TableSelectTopologyCacheGroupsController = function(parent, topology, cacheGroups, usedCacheGroupNames, $scope, $uibModal, $uibModalInstance, serverService) {

	let selectedCacheGroups = [],
		usedCacheGroupCount = 0;

	let markVisibleCacheGroups = function(selected) {
		let visibleCacheGroupIds = $('#availableCacheGroupsTable tr.cg-row').map(
			function() {
				return parseInt($(this).attr('id'));
			}).get();
		$scope.cacheGroups = cacheGroups.map(function(cg) {
			if (cg['used'] === true) {
				return cg;
			}
			if (selected && visibleCacheGroupIds.includes(cg.id)) {
				cg['selected'] = true;
			} else {
				cg['selected'] = false;
			}
			return cg;
		});
		updateSelectedCount();
	};

	let decorateCacheGroups = function() {
		$scope.cacheGroups = cacheGroups.map(function(cg) {
			const isUsed = usedCacheGroupNames.find(function(usedCacheGroupName) { return usedCacheGroupName === cg.name });
			if (isUsed) {
				cg['selected'] = true;
				cg['used'] = true;
				usedCacheGroupCount++;
			}
			return cg;
		});
	};

	let updateSelectedCount = function() {
		let visibleCacheGroupIds = $('#availableCacheGroupsTable tr.cg-row').map(
			function() {
				return parseInt($(this).attr('id'));
			}).get();

		selectedCacheGroups = $scope.cacheGroups.filter(function(cg) { return visibleCacheGroupIds.includes(cg.id) && cg['selected'] === true && !cg['used'] } );
		$('div.selected-count').html('<strong><span class="text-success">' + selectedCacheGroups.length + ' selected</span><span> | ' + usedCacheGroupCount + ' already used in topology</span></strong>');
	};

	$scope.parent = parent;

	$scope.cacheGroups = cacheGroups.filter(function(cg) {
		// displayed cache groups are filtered based on parent cache group type
		if (parent.type === 'ROOT') {
			// the topology root can have children of any kind (EDGE_LOC, MID_LOC, ORG_LOC)
			return true;
		} else if (parent.type === 'EDGE_LOC') {
			// EDGE_LOC can only have EDGE_LOC children
			return cg.typeName === 'EDGE_LOC';
		} else {
			// only EDGE_LOC and MID_LOC can be added farther down the topology tree (not root)
			return (cg.typeName === 'EDGE_LOC' || cg.typeName === 'MID_LOC');
		}
	});

	$scope.selectAll = function($event) {
		const checkbox = $event.target;
		if (checkbox.checked) {
			markVisibleCacheGroups(true);
		} else {
			markVisibleCacheGroups(false);
		}
	};

	$scope.onChange = function(cg) {
		if (cg.used) return;

		cg.selected = !cg.selected;
		updateSelectedCount();
	};

	$scope.viewCacheGroupServers = function(cg, $event) {
		if ($event) {
			$event.stopPropagation(); // this kills the click event so it doesn't trigger anything else
		}
		$uibModal.open({
			templateUrl: 'common/modules/table/topologyCacheGroupServers/table.topologyCacheGroupServers.tpl.html',
			controller: 'TableTopologyCacheGroupServersController',
			size: 'lg',
			resolve: {
				cacheGroupName: function() {
					return cg.name;
				},
				cacheGroupServers: function(serverService) {
					return serverService.getServers({ cachegroup: cg.id });
				}
			}
		});
	};

	$scope.submit = function() {
		/*  Cache groups that can act as a second parent include:
			1. cache groups that are not currently acting as the primary parent
			2. cache groups that exist currently in the topology
		 */
		let eligibleSecParentCandidates = cacheGroups.filter(function(cg) {
			return (parent.cachegroup && parent.cachegroup !== cg.name) &&
				usedCacheGroupNames.includes(cg.name);
		});

		if (eligibleSecParentCandidates.length === 0) {
			$uibModalInstance.close({ selectedCacheGroups: selectedCacheGroups, parent: parent.cachegroup, secParent: '' });
			return;
		}
		let params = {
			title: 'Assign secondary parent?',
			message: 'Would you like to assign a secondary parent to the following cache groups?<br><br>primary parent = ' + parent.cachegroup + '<br><br>'
		};
		params.message += selectedCacheGroups.map(function(cg) { return cg.name }).join('<br>') + '<br><br>';
		let modalInstance = $uibModal.open({
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
			// user wants to select a secondary parent
			let params = {
				title: 'Select a secondary parent',
				message: 'Please select a secondary parent that is part of the ' + topology.name + ' topology',
				key: 'name'
			};
			let modalInstance = $uibModal.open({
				templateUrl: 'common/modules/dialog/select/dialog.select.tpl.html',
				controller: 'DialogSelectController',
				size: 'md',
				resolve: {
					params: function () {
						return params;
					},
					collection: function() {
						// cache groups that are eligible to be a secondary parent include cache groups that are:
						return eligibleSecParentCandidates;
					}
				}
			});
			modalInstance.result.then(function(cg) {
				// user selected a secondary parent
				$uibModalInstance.close({ selectedCacheGroups: selectedCacheGroups, parent: parent.cachegroup, secParent: cg.name });
			}, function () {
				// user apparently changed their mind and doesn't want to select a secondary parent
				$uibModalInstance.close({ selectedCacheGroups: selectedCacheGroups, parent: parent.cachegroup, secParent: '' });
			});
		}, function () {
			// user doesn't want to select a secondary parent
			$uibModalInstance.close({ selectedCacheGroups: selectedCacheGroups, parent: parent.cachegroup, secParent: '' });
		});
	};

	$scope.cancel = function () {
		$uibModalInstance.dismiss('cancel');
	};

	angular.element(document).ready(function () {
		decorateCacheGroups();

		$('#availableCacheGroupsTable').DataTable({
			"scrollY": "60vh",
			"paging": false,
			"order": [[ 1, 'asc' ]],
			"dom": '<"selected-count">frtip',
			"drawCallback": function() {
				updateSelectedCount();
			},
			"columnDefs": [
				{ 'orderable': false, 'targets': [0,5] },
				{ "width": "5%", "targets": [ 0 ] },
				{ "width": "35%", "targets": [ 1 ] },
				{ "width": "15%", "targets": [ 2,3,4,5 ] }
			],
			"stateSave": false
		});
	});

};

TableSelectTopologyCacheGroupsController.$inject = ['parent', 'topology', 'cacheGroups', 'usedCacheGroupNames', '$scope', '$uibModal', '$uibModalInstance', 'serverService'];
module.exports = TableSelectTopologyCacheGroupsController;
