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

var TableCDNServersController = function(cdn, servers, $scope, $state, $uibModal, locationUtils, serverUtils, cdnService, propertiesModel) {

	var queueServerUpdates = function(cdn) {
		cdnService.queueServerUpdates(cdn.id)
			.then(
				function() {
					$scope.refresh();
				}
			);
	};

	var clearServerUpdates = function(cdn) {
		cdnService.clearServerUpdates(cdn.id)
			.then(
				function() {
					$scope.refresh();
				}
			);
	};

	$scope.cdn = cdn;

	$scope.servers = servers;

	$scope.editServer = function(id) {
		locationUtils.navigateToPath('/servers/' + id);
	};

	$scope.queueServerUpdates = function(cdn) {
		var params = {
			title: 'Queue Server Updates: ' + cdn.name,
			message: 'Are you sure you want to queue server updates for all ' + cdn.name + ' servers?'
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
			queueServerUpdates(cdn);
		}, function () {
			// do nothing
		});
	};

	$scope.clearServerUpdates = function(cdn) {
		var params = {
			title: 'Clear Server Updates: ' + cdn.name,
			message: 'Are you sure you want to clear server updates for all ' + cdn.name + ' servers?'
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
			clearServerUpdates(cdn);
		}, function () {
			// do nothing
		});
	};

	$scope.refresh = function() {
		$state.reload(); // reloads all the resolves for the view
	};

	$scope.showChartsButton = propertiesModel.properties.servers.charts.show;

	$scope.ssh = serverUtils.ssh;

	$scope.gotoMonitor = serverUtils.gotoMonitor;

	$scope.openCharts = serverUtils.openCharts;

	$scope.isOffline = serverUtils.isOffline;

	$scope.offlineReason = serverUtils.offlineReason;

	$scope.navigateToPath = locationUtils.navigateToPath;

	angular.element(document).ready(function () {
		$('#serversTable').dataTable({
			"aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
			"iDisplayLength": 25,
			"columnDefs": [
				{ 'orderable': false, 'targets': 11 }
			],
			"aaSorting": []
		});
	});

};

TableCDNServersController.$inject = ['cdn', 'servers', '$scope', '$state', '$uibModal', 'locationUtils', 'serverUtils', 'cdnService', 'propertiesModel'];
module.exports = TableCDNServersController;
