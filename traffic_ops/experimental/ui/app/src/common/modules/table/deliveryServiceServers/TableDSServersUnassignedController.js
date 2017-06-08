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

var TableDSServersUnassignedController = function(deliveryService, eligibleServers, assignedServers, $scope, $uibModalInstance) {

	var selectedServerIds = [];

	var addServer = function(serverId) {
		if (_.indexOf(selectedServerIds, serverId) == -1) {
			selectedServerIds.push(serverId);
		}
	};

	var removeServer = function(serverId) {
		selectedServerIds = _.without(selectedServerIds, serverId);
	};

	var addAll = function() {
		markServers(true);
		selectedServerIds = _.pluck(eligibleServers, 'id');
	};

	var removeAll = function() {
		markServers(false);
		selectedServerIds = [];
	};

	var markServers = function(selected) {
		$scope.selectedServers = _.map(eligibleServers, function(server) {
			server['selected'] = selected;
			return server;
		});
	};

	$scope.deliveryService = deliveryService;

	$scope.selectedServers = _.map(eligibleServers, function(eligibleServer) {
		var isAssigned = _.find(assignedServers, function(assignedServer) { return assignedServer.id == eligibleServer.id });
		if (isAssigned) {
			eligibleServer['selected'] = true; // so the checkbox will be checked
			selectedServerIds.push(eligibleServer.id); // so the server is added to selected servers
		}
		return eligibleServer;
	});

	$scope.allSelected = function() {
		return eligibleServers.length == selectedServerIds.length;
	};

	$scope.selectAll = function($event) {
		var checkbox = $event.target;
		if (checkbox.checked) {
			addAll();
		} else {
			removeAll();
		}
	};

	$scope.updateServers = function($event, serverId) {
		var checkbox = $event.target;
		if (checkbox.checked) {
			addServer(serverId);
		} else {
			removeServer(serverId);
		}
	};

	$scope.submit = function() {
		$uibModalInstance.close(selectedServerIds);
	};

	$scope.cancel = function () {
		$uibModalInstance.dismiss('cancel');
	};

	angular.element(document).ready(function () {
		$('#dsServersUnassignedTable').dataTable({
			"aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
			"iDisplayLength": 25,
			"order": [[ 1, 'asc' ]],
			"columnDefs": [
				{ 'orderable': false, 'targets': 0 },
				{ "width": "5%", "targets": 0 }
			]
		});
	});

};

TableDSServersUnassignedController.$inject = ['deliveryService', 'eligibleServers', 'assignedServers', '$scope', '$uibModalInstance'];
module.exports = TableDSServersUnassignedController;
