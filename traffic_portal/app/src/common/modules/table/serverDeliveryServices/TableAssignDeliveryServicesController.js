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

var TableAssignDeliveryServicesController = function(server, deliveryServices, assignedDeliveryServices, $scope, $uibModalInstance) {

	var selectedDsIds = [];

	var addDS = function(dsId) {
		if (_.indexOf(selectedDsIds, dsId) == -1) {
			selectedDsIds.push(dsId);
		}
	};

	var removeDS = function(dsId) {
		selectedDsIds = _.without(selectedDsIds, dsId);
	};

	var addAll = function() {
		markDeliveryServices(true);
		selectedDsIds = _.pluck(deliveryServices, 'id');
	};

	var removeAll = function() {
		markDeliveryServices(false);
		selectedDsIds = [];
	};

	var markDeliveryServices = function(selected) {
		$scope.selectedDeliveryServices = _.map(deliveryServices, function(ds) {
			ds['selected'] = selected;
			return ds;
		});
	};

	$scope.server = server;

	$scope.selectedDeliveryServices = _.map(deliveryServices, function(ds) {
		var isAssigned = _.find(assignedDeliveryServices, function(assignedDS) { return assignedDS.id == ds.id });
		if (isAssigned) {
			ds['selected'] = true; // so the checkbox will be checked
			selectedDsIds.push(ds.id); // so the ds is added to selected dsIds
		}
		return ds;
	});

	$scope.allSelected = function() {
		return deliveryServices.length == selectedDsIds.length;
	};

	$scope.selectAll = function($event) {
		var checkbox = $event.target;
		if (checkbox.checked) {
			addAll();
		} else {
			removeAll();
		}
	};

	$scope.updateDeliveryServices = function($event, dsId) {
		var checkbox = $event.target;
		if (checkbox.checked) {
			addDS(dsId);
		} else {
			removeDS(dsId);
		}
	};

	$scope.submit = function() {
		$uibModalInstance.close(selectedDsIds);
	};

	$scope.cancel = function () {
		$uibModalInstance.dismiss('cancel');
	};

	angular.element(document).ready(function () {
		var assignDeliveryServicesTable = $('#assignDeliveryServicesTable').dataTable({
			"aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
			"iDisplayLength": 25,
			"order": [[ 1, 'asc' ]],
			"columnDefs": [
				{ 'orderable': false, 'targets': 0 },
				{ "width": "5%", "targets": 0 }
			],
			"stateSave": false
		});
		assignDeliveryServicesTable.on( 'search.dt', function () {
			var search = $('#assignDeliveryServicesTable_filter input').val();
			if (search.length > 0) {
				$("#selectAllCB").attr("disabled", true);
			} else {
				$("#selectAllCB").removeAttr("disabled");
			}
		} );
	});

};

TableAssignDeliveryServicesController.$inject = ['server', 'deliveryServices', 'assignedDeliveryServices', '$scope', '$uibModalInstance'];
module.exports = TableAssignDeliveryServicesController;
