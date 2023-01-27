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
 * @param {*} topology
 * @param {*} cacheGroups
 * @param {*} $scope
 * @param {import("angular").IControllerService} $controller
 * @param {import("../../../../service/utils/LocationUtils")} locationUtils
 * @param {import("../../../../api/TopologyService")} topologyService
 * @param {import("../../../../models/MessageModel")} messageModel
 * @param {import("../../../../service/utils/TopologyUtils")} topologyUtils
 */
var FormNewTopologyController = function(topology, cacheGroups, $scope, $controller, locationUtils, topologyService, messageModel, topologyUtils) {

	// extends the FormTopologyController to inherit common methods
	angular.extend(this, $controller('FormTopologyController', { topology: topology, cacheGroups: cacheGroups, $scope: $scope }));

	$scope.topologyName = 'New';

	$scope.settings = {
		isNew: true,
		saveLabel: 'Create'
	};

	$scope.save = function(name, description, topologyTree) {
		let normalizedTopology = topologyUtils.getNormalizedTopology(name, description, topologyTree);
		topologyService.createTopology(normalizedTopology).
			then(function(result) {
				messageModel.setMessages(result.data.alerts, true);
				locationUtils.navigateToPath('/topologies');
			});
	};

};

FormNewTopologyController.$inject = ['topology', 'cacheGroups', '$scope', '$controller', 'locationUtils', 'topologyService', 'messageModel', 'topologyUtils'];
module.exports = FormNewTopologyController;
