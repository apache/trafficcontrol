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

var TableDSServersUnassignedController = function(deliveryService, servers, $scope, $uibModalInstance) {

	var selectedServers = [];

	$scope.deliveryService = deliveryService;

	$scope.unassignedServers = servers;

	var addServer = function(serverId) {
		if (_.indexOf(selectedServers, serverId) == -1) {
			selectedServers.push(serverId);
		}
	};

	var removeServer = function(serverId) {
		selectedServers = _.without(selectedServers, serverId);
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
		$uibModalInstance.close(selectedServers);
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

TableDSServersUnassignedController.$inject = ['deliveryService', 'servers', '$scope', '$uibModalInstance'];
module.exports = TableDSServersUnassignedController;
