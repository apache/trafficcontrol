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

/**
 *
 * @param {*} cdn
 * @param {*} currentSnapshot
 * @param {*} newSnapshot
 * @param {*} $scope
 * @param {import("../../../../common/service/utils/angular.ui.bootstrap").IModalService} $uibModal
 * @param {import("../../../../common/service/utils/LocationUtils")} locationUtils
 * @param {import("../../../../common/service/utils/CollectionUtils")} collectionUtils
 * @param {import("../../../../common/api/CDNService")} cdnService
 * @param {import("../../../../common/models/PropertiesModel")} propertiesModel
 */
let ConfigController = function (cdn, currentSnapshot, newSnapshot, $scope, $uibModal, locationUtils, collectionUtils, cdnService, propertiesModel) {

	const oldConfig = currentSnapshot.config,
		newConfig = newSnapshot.config;

	const oldTrafficRouters = currentSnapshot.contentRouters,
		newTrafficRouters = newSnapshot.contentRouters;

	const oldTrafficMonitors = currentSnapshot.monitors,
		newTrafficMonitors = newSnapshot.monitors;

	const oldTrafficServers = currentSnapshot.contentServers,
		newTrafficServers = newSnapshot.contentServers;

	const oldDeliveryServices = currentSnapshot.deliveryServices,
		newDeliveryServices = newSnapshot.deliveryServices;

	const oldEdgeCacheGroups = currentSnapshot.edgeLocations,
		newEdgeCacheGroups = newSnapshot.edgeLocations;

	const oldTrafficRouterCacheGroups = currentSnapshot.trafficRouterLocations,
		newTrafficRouterCacheGroups = newSnapshot.trafficRouterLocations;

	const oldTopologies = currentSnapshot.topologies,
		newTopologies = newSnapshot.topologies;

	const oldStats = currentSnapshot.stats,
		newStats = newSnapshot.stats;

	let performDiff = function (oldJSON, newJSON, destination) {
		let added = 0,
			removed = 0,
			updated = 0;

		let oldConfig = oldJSON || {},
			newConfig = newJSON || {};

		let diff = jsonpatch.compare(oldConfig, newConfig);
		diff.forEach(function (change) {
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

	function minimizeServerCapabilitiesDiff(oldTrafficServers, newTrafficServers) {
		if (!(oldTrafficServers instanceof Object) || !(newTrafficServers instanceof Object)) {
			return;
		}
		const oldServersIterator = Object.entries(oldTrafficServers).entries();
		const newServersIterator = Object.entries(newTrafficServers).entries();
		const capabilitiesKey = "capabilities";
		for (let oldServersNext = oldServersIterator.next(), newServersNext = newServersIterator.next(); !(oldServersNext.done || newServersNext.done);) {
			const [, [oldHostname, oldServer]] = oldServersNext.value;
			const [, [newHostname, newServer]] = newServersNext.value;
			if (oldHostname < newHostname) {
				oldServersNext = oldServersIterator.next();
				continue;
			} else if (oldHostname > newHostname) {
				newServersNext = newServersIterator.next();
				continue;
			}
			const oldCapabilities = oldServer[capabilitiesKey];
			const newCapabilities = newServer[capabilitiesKey];
			if (oldCapabilities instanceof Object && newCapabilities instanceof Object) {
				newServer[capabilitiesKey] = collectionUtils.minimizeArrayDiff(oldCapabilities, newCapabilities);
			}
			oldServersNext = oldServersIterator.next();
			newServersNext = newServersIterator.next();
		}

	}

	let snapshot = function () {
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

	$scope.topologiesCount = {
		added: 0,
		removed: 0,
		updated: 0,
		templateUrl: 'topPopoverTemplate.html'
	};

	$scope.statsCount = {
		added: 0,
		removed: 0,
		updated: 0,
		templateUrl: 'statsPopoverTemplate.html'
	};

	$scope.confirmSnapshot = function (cdn) {
		let params = {
			title: 'Perform Snapshot',
			message: 'Are you sure you want to snapshot the ' + cdn.name + ' config?'
		};
		let modalInstance = $uibModal.open({
			templateUrl: 'common/modules/dialog/confirm/dialog.confirm.tpl.html',
			controller: 'DialogConfirmController',
			size: 'sm',
			resolve: {
				params: function () {
					return params;
				}
			}
		});
		modalInstance.result.then(function () {
			snapshot();
		}, function () {
			// do nothing
		});
	};

	$scope.tabSelected = function() {
		// issue 3863 - adjust column headers when tab is selected and data table is visible. hacky...sorry...
		window.setTimeout(function() {
			$($.fn.dataTable.tables(true)).DataTable()
				.columns.adjust();
			},100);
	};

	$scope.navigateToPath = (path, unsavedChanges) => locationUtils.navigateToPath(path, unsavedChanges);

	angular.element(document).ready(function () {

		$('table.changes').dataTable({
			"lengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
			"pageLength": 25,
			"order": [[0, "asc"]],
			"language": {
				"emptyTable": "No pending changes"
			},
			"buttons": [],
			"columnDefs": [
				{ 'orderable': false, 'targets': [2, 3] }
			]
		});

	});

	let init = function () {
		performDiff(oldConfig, newConfig, 'config');
		performDiff(oldTrafficRouters, newTrafficRouters, 'contentRouters');
		performDiff(oldTrafficMonitors, newTrafficMonitors, 'monitors');
		minimizeServerCapabilitiesDiff(oldTrafficServers, newTrafficServers);
		performDiff(oldTrafficServers, newTrafficServers, 'contentServers');
		performDiff(oldDeliveryServices, newDeliveryServices, 'deliveryServices');
		performDiff(oldEdgeCacheGroups, newEdgeCacheGroups, 'edgeLocations');
		performDiff(oldTrafficRouterCacheGroups, newTrafficRouterCacheGroups, 'trLocations');
		performDiff(oldTopologies, newTopologies, 'topologies');
		performDiff(oldStats, newStats, 'stats');
	};
	init();

};

ConfigController.$inject = ['cdn', 'currentSnapshot', 'newSnapshot', '$scope', '$uibModal', 'locationUtils', 'collectionUtils', 'cdnService', 'propertiesModel'];
module.exports = ConfigController;
