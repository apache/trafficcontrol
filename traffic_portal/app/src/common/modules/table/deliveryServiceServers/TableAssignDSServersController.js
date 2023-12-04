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

/** @typedef { import('../../../../common/modules/table/agGrid/CommonGridController').CGC } CGC */

var TableAssignDSServersController = function(deliveryService, servers, assignedServers, $scope, $uibModalInstance) {

	$scope.selectedServers = [];

	/** @type CGC.ColumnDefinition */
	$scope.columns = [
		{
			headerName: "Host",
			field: "hostName",
			checkboxSelection: true,
			headerCheckboxSelection: true,
		},
		{
			headerName: "Cache Group",
			field: "cacheGroup",
		},
		{
			headerName: "Profile(s)",
			field: "profile",
			valueGetter:  function(params) {
				return params.data.profiles;
			},
			tooltipValueGetter: function(params) {
				return params.data.profiles.join(", ");
			}
		}
	];

	$scope.deliveryService = deliveryService;

	$scope.servers = servers.map(server => {
		let isAssigned = assignedServers.find(assignedServer => assignedServer.id === server.id);
		if (isAssigned) {
			server['selected'] = true;
		}
		return server;
	});

	$scope.submit = function() {
		const selectedServerIds = this.selectedServers.map(s => s["id"]);
		$uibModalInstance.close(selectedServerIds);
	};

	$scope.cancel = function () {
		$uibModalInstance.dismiss('cancel');
	};

	/** @type CGC.GridSettings */
	$scope.gridOptions = {
		selectRows: true,
		selectionProperty: "selected"
	};
};

TableAssignDSServersController.$inject = ['deliveryService', 'servers', 'assignedServers', '$scope', '$uibModalInstance'];
module.exports = TableAssignDSServersController;
