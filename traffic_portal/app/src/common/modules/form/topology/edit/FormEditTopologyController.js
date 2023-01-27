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
 * @param {*} topologies
 * @param {*} cacheGroups
 * @param {*} $scope
 * @param {import("angular").IControllerService} $controller
 * @param {import("../../../../service/utils/angular.ui.bootstrap").IModalService} $uibModal
 * @param {import("../../../../service/utils/LocationUtils")} locationUtils
 * @param {import("../../../../api/TopologyService")} topologyService
 * @param {import("../../../../models/MessageModel")} messageModel
 * @param {import("../../../../service/utils/TopologyUtils")} topologyUtils
 */
var FormEditTopologyController = function(topologies, cacheGroups, $scope, $controller, $uibModal, locationUtils, topologyService, messageModel, topologyUtils) {

	// extends the FormTopologyController to inherit common methods
	angular.extend(this, $controller('FormTopologyController', { topology: topologies[0], cacheGroups: cacheGroups, $scope: $scope }));

	let deleteTopology = function(topology) {
		topologyService.deleteTopology(topology)
			.then(function() {
				messageModel.setMessages([ { level: 'success', text: 'Topology deleted' } ], true);
				locationUtils.navigateToPath('/topologies');
			});
	};

	$scope.topologyName = angular.copy($scope.topology.name);

	$scope.settings = {
		isNew: false,
		saveLabel: 'Update'
	};

	$scope.save = function(currentName, newName, description, topologyTree) {
		let normalizedTopology = topologyUtils.getNormalizedTopology(newName, description, topologyTree);
		topologyService.updateTopology(normalizedTopology, currentName).
			then(function(result) {
				messageModel.setMessages(result.data.alerts, currentName !== newName);
				locationUtils.navigateToPath('/topologies/edit?name=' + result.data.response.name);
			});
	};

	$scope.confirmDelete = function(topology) {
		let params = {
			title: 'Delete Topology: ' + topology.name,
			key: topology.name
		};
		let modalInstance = $uibModal.open({
			templateUrl: 'common/modules/dialog/delete/dialog.delete.tpl.html',
			controller: 'DialogDeleteController',
			size: 'md',
			resolve: {
				params: function () {
					return params;
				}
			}
		});
		modalInstance.result.then(function() {
			deleteTopology(topology);
		});
	};

	$scope.confirmTopologyQueueServerUpdates = function(topology) {
		const params = {
			title: 'Queue Server Updates: ' + topology.name,
			message: "Please select a CDN"
		};
		const modalInstance = $uibModal.open({
			templateUrl: 'common/modules/dialog/select/dialog.select.tpl.html',
			controller: 'DialogSelectController',
			size: 'md',
			resolve: {
				params: function () {
					return params;
				},
				collection: function(cdnService) {
					return cdnService.getCDNs();
				}
			}
		});
		modalInstance.result.then(function(cdn) {
			topologyService.queueServerUpdates(topology.name, cdn.id);
		}, function () {
			// do nothing
		});
	};

	$scope.confirmTopologyClearServerUpdates = function(topology) {
		const params = {
			title: 'Clear Server Updates: ' + topology.name,
			message: "Please select a CDN"
		};
		const modalInstance = $uibModal.open({
			templateUrl: 'common/modules/dialog/select/dialog.select.tpl.html',
			controller: 'DialogSelectController',
			size: 'md',
			resolve: {
				params: function () {
					return params;
				},
				collection: function(cdnService) {
					return cdnService.getCDNs();
				}
			}
		});
		modalInstance.result.then(function(cdn) {
			topologyService.clearServerUpdates(topology.name, cdn.id);
		}, function () {
			// do nothing
		});
	};

};

FormEditTopologyController.$inject = ['topologies', 'cacheGroups', '$scope', '$controller', '$uibModal', 'locationUtils', 'topologyService', 'messageModel', 'topologyUtils'];
module.exports = FormEditTopologyController;
