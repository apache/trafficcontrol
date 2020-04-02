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

var TableSelectTopologyCacheGroupsController = function(node, topology, cacheGroups, usedCacheGroupNames, $scope, $uibModalInstance) {

	var selectedCacheGroups = [],
		usedCacheGroupCount = 0;

	var addAll = function() {
		markVisibleCacheGroups(true);
	};

	var removeAll = function() {
		markVisibleCacheGroups(false);
	};

	var markVisibleCacheGroups = function(selected) {
		var visibleCacheGroupNames = $('#availableCacheGroupsTable tr.cg-row').map(
			function() {
				return parseInt($(this).attr('id'));
			}).get();
		$scope.cacheGroups = _.map(cacheGroups, function(cg) {
			if (visibleCacheGroupNames.includes(cg.id)) {
				cg['selected'] = selected;
			}
			return cg;
		});
		updateSelectedCount();
	};

	var decorateCacheGroups = function() {
		$scope.cacheGroups = _.map(cacheGroups, function(cg) {
			var isUsed = _.find(usedCacheGroupNames, function(usedCacheGroupName) { return usedCacheGroupName == cg.name });
			if (isUsed) {
				cg['selected'] = true;
				cg['used'] = true;
				usedCacheGroupCount++;
			}
			return cg;
		});
	};

	var updateSelectedCount = function() {
		selectedCacheGroups = _.filter($scope.cacheGroups, function(cg) { return cg['selected'] == true && !cg['used'] } );
		$('div.selected-count').html('<strong><span class="text-success">' + selectedCacheGroups.length + ' selected</span><span> | ' + usedCacheGroupCount + ' currently used</span></strong>');
	};

	$scope.topology = topology;

	$scope.cacheGroups = _.filter(cacheGroups, function(cg) {
		// all cg types (ORG_LOC, MID_LOC, EDGE_LOC) can be added to the top of a topology
		// but only EDGE_LOC and MID_LOC can be added farther down the topology tree
		if (node.type === 'ORIGIN_LAYER') return (cg.typeName === 'EDGE_LOC' || cg.typeName === 'MID_LOC' || cg.typeName === 'ORG_LOC');
		return (cg.typeName === 'EDGE_LOC' || cg.typeName === 'MID_LOC');
	});

	$scope.selectAll = function($event) {
		// todo:
		alert('select/unselect all visible/not used cgs')
		// var checkbox = $event.target;
		// if (checkbox.checked) {
		// 	addAll();
		// } else {
		// 	removeAll();
		// }
	};

	$scope.onChange = function(cg) {
		if (cg.used) return;

		cg.selected = !cg.selected;
		updateSelectedCount();
	};

	$scope.submit = function() {
		$uibModalInstance.close(selectedCacheGroups);
	};

	$scope.cancel = function () {
		$uibModalInstance.dismiss('cancel');
	};

	angular.element(document).ready(function () {
		var availableCacheGroupsTable = $('#availableCacheGroupsTable').dataTable({
			"scrollY": "60vh",
			"paging": false,
			"order": [[ 1, 'asc' ]],
			"dom": '<"selected-count">frtip',
			"columnDefs": [
				{ 'orderable': false, 'targets': 0 },
				{ "width": "5%", "targets": 0 }
			],
			"stateSave": false
		});
		availableCacheGroupsTable.on( 'search.dt', function () {
			$("#selectAllCB").removeAttr("checked"); // uncheck the all box when filtering
		} );
		decorateCacheGroups();
		updateSelectedCount();
	});

};

TableSelectTopologyCacheGroupsController.$inject = ['node', 'topology', 'cacheGroups', 'usedCacheGroupNames', '$scope', '$uibModalInstance'];
module.exports = TableSelectTopologyCacheGroupsController;
