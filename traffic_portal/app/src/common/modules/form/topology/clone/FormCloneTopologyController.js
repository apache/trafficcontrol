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

var FormCloneTopologyController = function(topologies, cacheGroups, $scope, $controller) {

	// extends the FormNewTopologyController to inherit common methods
	angular.extend(this, $controller('FormNewTopologyController', { topology: topologies[0], cacheGroups: cacheGroups, $scope: $scope }));

	$scope.topologyName = angular.copy($scope.topology.name) + ' clone';

	$scope.settings = {
		isNew: true,
		saveLabel: 'Clone'
	};

	let init = function() {
		// cloned topology needs a new name
		$scope.topology.name = '';
	};
	init();

};

FormCloneTopologyController.$inject = ['topologies', 'cacheGroups', '$scope', '$controller'];
module.exports = FormCloneTopologyController;
