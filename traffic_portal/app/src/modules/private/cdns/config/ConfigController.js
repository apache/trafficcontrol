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

var ConfigController = function(cdn, currentSnapshot, newSnapshot, $scope, $state, $uibModal, $window, locationUtils, cdnService, propertiesModel, ENV) {

	var oldConfig = currentSnapshot.config,
		newConfig = newSnapshot.config;

	var oldTrafficRouters = currentSnapshot.contentRouters,
		newTrafficRouters = newSnapshot.contentRouters;

	var oldTrafficMonitors = currentSnapshot.monitors,
		newTrafficMonitors = newSnapshot.monitors;

	var oldTrafficServers = currentSnapshot.contentServers,
		newTrafficServers = newSnapshot.contentServers;

	var oldDeliveryServices = currentSnapshot.deliveryServices,
		newDeliveryServices = newSnapshot.deliveryServices;

	var oldEdgeCacheGroups = currentSnapshot.edgeLocations,
		newEdgeCacheGroups = newSnapshot.edgeLocations;

	var oldTrafficRouterCacheGroups = currentSnapshot.trafficRouterLocations,
		newTrafficRouterCacheGroups = newSnapshot.trafficRouterLocations;

	var oldStats = currentSnapshot.stats,
		newStats = newSnapshot.stats;

	var performDiff = function(oldJSON, newJSON, destination) {
		var added = 0,
			removed = 0,
			updated = 0;

		var oldConfig = oldJSON || {},
			newConfig = newJSON || {};

		var diff = jsonpatch.compare(oldConfig, newConfig);
		diff.forEach(function(change){
			if (change.op == 'add') {
				added++;
			} else if (change.op == 'remove') {
				removed++;
			} else if (change.op == 'replace') {
				change.op = 'update'; // changing the name to 'update'
				updated++;
			}
		});

		$scope[destination + "Count"].added = added;
		$scope[destination + "Count"].removed = removed;
		$scope[destination + "Count"].updated = updated;
		$scope[destination + "Changes"] = diff;
	};

	var snapshot = function() {
		cdnService.snapshot(cdn);
	};

	$scope.cdn = cdn;

	$scope.expandLevel = propertiesModel.properties.snapshot.diff.expandLevel;

	$scope.configCount = {
		added: 0,
		removed: 0,
		updated: 0,
		templateUrl: 'configPopoverTemplate.html'
	};

	$scope.contentRoutersCount = {
		added: 0,
		removed: 0,
		updated: 0,
		templateUrl: 'crPopoverTemplate.html'
	};

	$scope.monitorsCount = {
		added: 0,
		removed: 0,
		updated: 0,
		templateUrl: 'mPopoverTemplate.html'
	};

	$scope.contentServersCount = {
		added: 0,
		removed: 0,
		updated: 0,
		templateUrl: 'csPopoverTemplate.html'
	};

	$scope.deliveryServicesCount = {
		added: 0,
		removed: 0,
		updated: 0,
		templateUrl: 'dsPopoverTemplate.html'
	};

	$scope.edgeLocationsCount = {
		added: 0,
		removed: 0,
		updated: 0,
		templateUrl: 'elPopoverTemplate.html'
	};

	$scope.trLocationsCount = {
		added: 0,
		removed: 0,
		updated: 0,
		templateUrl: 'tlPopoverTemplate.html'
	};

	$scope.statsCount = {
		added: 0,
		removed: 0,
		updated: 0,
		templateUrl: 'statsPopoverTemplate.html'
	};

	$scope.viewSnapshot = function(cdn, type) {
		var url = ENV.api['root'] + 'cdns/' + cdn.name + '/snapshot';
		if (type == 'pending') {
			url += '/new';
		}
		$window.open(url, '_blank');
	};

	$scope.confirmSnapshot = function(cdn) {
		var params = {
			title: 'Perform Snapshot',
			message: 'Are you sure you want to snapshot the ' + cdn.name + ' config?'
		};
		var modalInstance = $uibModal.open({
			templateUrl: 'common/modules/dialog/confirm/dialog.confirm.tpl.html',
			controller: 'DialogConfirmController',
			size: 'sm',
			resolve: {
				params: function () {
					return params;
				}
			}
		});
		modalInstance.result.then(function() {
			snapshot();
		}, function () {
			// do nothing
		});
	};

	$scope.navigateToPath = locationUtils.navigateToPath;

	angular.element(document).ready(function () {

		$('table.changes').dataTable({
			"aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
			"iDisplayLength": 25,
			"aaSorting": [[ 0, "asc" ]],
			"language": {
				"emptyTable": "No pending changes"
			},
			"columnDefs": [
				{ 'orderable': false, 'targets': [2,3] }
			]
		});

	});

	var init = function() {
		performDiff(oldConfig, newConfig, 'config');
		performDiff(oldTrafficRouters, newTrafficRouters, 'contentRouters');
		performDiff(oldTrafficMonitors, newTrafficMonitors, 'monitors');
		performDiff(oldTrafficServers, newTrafficServers, 'contentServers');
		performDiff(oldDeliveryServices, newDeliveryServices, 'deliveryServices');
		performDiff(oldEdgeCacheGroups, newEdgeCacheGroups, 'edgeLocations');
		performDiff(oldTrafficRouterCacheGroups, newTrafficRouterCacheGroups, 'trLocations');
		performDiff(oldStats, newStats, 'stats');
	};
	init();

};

ConfigController.$inject = ['cdn', 'currentSnapshot', 'newSnapshot', '$scope', '$state', '$uibModal', '$window', 'locationUtils', 'cdnService', 'propertiesModel', 'ENV'];
module.exports = ConfigController;
