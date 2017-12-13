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

var ConfigController = function(cdn, currentSnapshot, newSnapshot, $scope, $state, $timeout, $uibModal, locationUtils, cdnService, propertiesModel) {

	$scope.cdn = cdn;

	var diffSettings = propertiesModel.properties.snapshot.diff;

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

	var oldStats = currentSnapshot.stats,
		newStats = newSnapshot.stats;

	var performDiff = function(oldJSON, newJSON, destination) {
		var div = null,
			added = 0,
			removed = 0,
			context = diffSettings.context; // only show X lines of context (around added or removed parts) as defined in traffic_portal_properties.json;

		var display = document.getElementById(destination),
			fragment = document.createDocumentFragment();

		if (oldJSON) {
			var diff = JsDiff.diffJson(oldJSON, newJSON);
			diff.forEach(function(part){
				var partChanged = part.added || part.removed;
				if (!partChanged && part.count > (context * 2)) {
					var partArr = part.value.split("\n"),
						newArr = [];
					_.each(partArr, function(element, index) {
						if (index < context || index > partArr.length - context) {
							newArr.push(element);
						}
					});
					newArr.splice(context, 0, "\n*****************\n*   TRUNCATED   *\n*****************\n");
					part.value = newArr.join("\n");
				}
				if (part.added) {
					added++;
				} else if (part.removed) {
					removed++;
				}
				div = document.createElement('div');
				div.className = part.added ? 'added' : part.removed ? 'removed' : 'no-change';

				if (partChanged) {
					var prepend = part.added ? '++' : '--';
					part.value = prepend + part.value.slice(2);
				}

				div.appendChild(document.createTextNode(part.value));
				fragment.appendChild(div);
			});

			$scope[destination + "Count"].added = added;
			$scope[destination + "Count"].removed = removed;
			display.innerHTML = '';
			display.appendChild(fragment);
		} else {
			display.innerHTML = 'Diff failed. You may need to perform your first snapshot.';
		}

	};

	var snapshot = function() {
		cdnService.snapshot(cdn);
	};

	$scope.configCount = {
		added: 0,
		removed: 0,
		templateUrl: 'configPopoverTemplate.html'
	};

	$scope.contentRoutersCount = {
		added: 0,
		removed: 0,
		templateUrl: 'crPopoverTemplate.html'
	};

	$scope.monitorsCount = {
		added: 0,
		removed: 0,
		templateUrl: 'mPopoverTemplate.html'
	};

	$scope.contentServersCount = {
		added: 0,
		removed: 0,
		templateUrl: 'csPopoverTemplate.html'
	};

	$scope.deliveryServicesCount = {
		added: 0,
		removed: 0,
		templateUrl: 'dsPopoverTemplate.html'
	};

	$scope.edgeLocationsCount = {
		added: 0,
		removed: 0,
		templateUrl: 'elPopoverTemplate.html'
	};

	$scope.statsCount = {
		added: 0,
		removed: 0,
		templateUrl: 'statsPopoverTemplate.html'
	};

	$scope.diffConfig = function(timeout) {
		$('#config').html('<i class="fa fa-refresh fa-spin fa-1x fa-fw"></i> Generating diff...');
		$timeout(function() {
			performDiff(oldConfig, newConfig, 'config');
		}, timeout);
	};

	$scope.diffContentRouters = function(timeout) {
		$('#contentRouters').html('<i class="fa fa-refresh fa-spin fa-1x fa-fw"></i> Generating diff...');
		$timeout(function() {
			performDiff(oldTrafficRouters, newTrafficRouters, 'contentRouters');
		}, timeout);
	};

	$scope.diffMonitors = function(timeout) {
		$('#monitors').html('<i class="fa fa-refresh fa-spin fa-1x fa-fw"></i> Generating diff...');
		$timeout(function() {
			performDiff(oldTrafficMonitors, newTrafficMonitors, 'monitors');
		}, timeout);
	};

	$scope.diffContentServers = function(timeout) {
		$('#contentServers').html('<i class="fa fa-refresh fa-spin fa-1x fa-fw"></i> Generating diff...');
		$timeout(function() {
			performDiff(oldTrafficServers, newTrafficServers, 'contentServers');
		}, timeout);
	};

	$scope.diffDeliveryServices = function(timeout) {
		$('#deliveryServices').html('<i class="fa fa-refresh fa-spin fa-1x fa-fw"></i> Generating diff...');
		$timeout(function() {
			performDiff(oldDeliveryServices, newDeliveryServices, 'deliveryServices');
		}, timeout);
	};

	$scope.diffEdgeLocations = function(timeout) {
		$('#edgeLocations').html('<i class="fa fa-refresh fa-spin fa-1x fa-fw"></i> Generating diff...');
		$timeout(function() {
			performDiff(oldEdgeCacheGroups, newEdgeCacheGroups, 'edgeLocations');
		}, timeout);
	};

	$scope.diffStats = function(timeout) {
		$('#stats').html('<i class="fa fa-refresh fa-spin fa-1x fa-fw"></i> Generating diff...');
		$timeout(function() {
			performDiff(oldStats, newStats, 'stats');
		}, timeout);
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
		$scope.diffConfig(0);
		$scope.diffContentRouters(0);
		$scope.diffMonitors(0);
		$scope.diffContentServers(0);
		$scope.diffDeliveryServices(0);
		$scope.diffEdgeLocations(0);
		$scope.diffStats(0);
	});

};

ConfigController.$inject = ['cdn', 'currentSnapshot', 'newSnapshot', '$scope', '$state', '$timeout', '$uibModal', 'locationUtils', 'cdnService', 'propertiesModel'];
module.exports = ConfigController;
