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

var TableParameterCacheGroupsController = function(parameter, cacheGroups, $scope, $state, locationUtils) {

	$scope.parameter = parameter;

	$scope.cacheGroups = cacheGroups;

	$scope.addCacheGroup = function() {
		alert('not hooked up yet: add cg to parameter');
	};

	$scope.removeCacheGroup = function() {
		alert('not hooked up yet: remove cg from parameter');
	};

	$scope.refresh = function() {
		$state.reload(); // reloads all the resolves for the view
	};

	$scope.navigateToPath = locationUtils.navigateToPath;

	angular.element(document).ready(function () {
		$('#cacheGroupsTable').dataTable({
			"aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
			"iDisplayLength": 25,
			"aaSorting": []
		});
	});

};

TableParameterCacheGroupsController.$inject = ['parameter', 'cacheGroups', '$scope', '$state', 'locationUtils'];
module.exports = TableParameterCacheGroupsController;
