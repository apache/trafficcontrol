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

var TableAssignDeliveryServicesController = function(server, deliveryServices, assignedDeliveryServices, $scope, $uibModalInstance, serverUtils) {

	var selectedDeliveryServices = [];

	var addAll = function() {
		markVisibleDeliveryServices(true);
	};

	var removeAll = function() {
		markVisibleDeliveryServices(false);
	};

	var markVisibleDeliveryServices = function(selected) {
		var visibleDSIds = $('#assignDSTable tr.ds-row').map(
			function() {
				return parseInt($(this).attr('id'));
			}).get();
		$scope.selectedDeliveryServices = _.map(deliveryServices, function(ds) {
			if (ds.topology && $scope.isCache(server)) {
				return ds;
			}
			if (visibleDSIds.includes(ds.id)) {
				ds['selected'] = selected;
			}
			return ds;
		});
		updateSelectedCount();
	};

	var updateSelectedCount = function() {
		selectedDeliveryServices = _.filter($scope.selectedDeliveryServices, function(ds) { return ds['selected'] == true; } );
		$('div.selected-count').html('<b>' + selectedDeliveryServices.length + ' delivery services selected</b>');
	};

	$scope.server = server;

	$scope.isCache = serverUtils.isCache;

	$scope.selectedDeliveryServices = _.map(deliveryServices, function(ds) {
		var isAssigned = _.find(assignedDeliveryServices, function(assignedDS) { return assignedDS.id == ds.id });
		if (isAssigned) {
			ds['selected'] = true;
		}
		return ds;
	});

	$scope.toggleRow = function(ds) {
		// a ds w/ a topology has no use being assigned to cache servers
		if (ds.topology && $scope.isCache(server)) {
			return;
		}
		ds.selected = !ds.selected;
		$scope.onChange();
	};

	$scope.selectAll = function($event) {
		var checkbox = $event.target;
		if (checkbox.checked) {
			addAll();
		} else {
			removeAll();
		}
	};

	$scope.onChange = function() {
		updateSelectedCount();
	};

	$scope.submit = function() {
		var selectedDSIds = _.pluck(selectedDeliveryServices, 'id');
		$uibModalInstance.close(selectedDSIds);
	};

	$scope.cancel = function () {
		$uibModalInstance.dismiss('cancel');
	};

	angular.element(document).ready(function () {
		var assignDSTable = $('#assignDSTable').dataTable({
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
		assignDSTable.on( 'search.dt', function () {
			$("#selectAllCB").removeAttr("checked"); // uncheck the all box when filtering
		} );
		updateSelectedCount();
	});

};

TableAssignDeliveryServicesController.$inject = ['server', 'deliveryServices', 'assignedDeliveryServices', '$scope', '$uibModalInstance', 'serverUtils'];
module.exports = TableAssignDeliveryServicesController;
