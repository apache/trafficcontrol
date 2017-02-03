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

var TableCacheGroupsServersController = function(cacheGroup, servers, $scope, locationUtils) {

	$scope.cacheGroup = cacheGroup;

	$scope.servers = servers;

	$scope.editServer = function(id) {
		locationUtils.navigateToPath('/configure/servers/' + id);
	};

	$scope.queueUpdates = function() {
		alert('not hooked up yet: queuing updates for all cachegroup servers');
	};

	$scope.dequeueUpdates = function() {
		alert('not hooked up yet: dequeuing updates for all cachegroup servers');
	};

	$scope.navigateToPath = locationUtils.navigateToPath;

	angular.element(document).ready(function () {
		$('#serversTable').dataTable({
			"aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
			"iDisplayLength": 100
		});
	});

};

TableCacheGroupsServersController.$inject = ['cacheGroup', 'servers', '$scope', 'locationUtils'];
module.exports = TableCacheGroupsServersController;
