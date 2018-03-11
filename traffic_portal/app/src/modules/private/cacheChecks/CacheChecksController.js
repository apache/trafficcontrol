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

var CacheChecksController = function(cacheChecks, showReval, $scope, $state, $interval, locationUtils, serverUtils, propertiesModel) {

	var cacheChecksInterval,
		autoRefresh = false,
		refreshRateInMS = 10000;

	var createInterval = function() {
		killInterval();
		cacheChecksInterval = $interval(function() { $scope.refresh() }, refreshRateInMS );
	};

	var killInterval = function() {
		if (angular.isDefined(cacheChecksInterval)) {
			$interval.cancel(cacheChecksInterval);
			cacheChecksInterval = undefined;
		}
	};

	$scope.cacheChecks = cacheChecks;

	$scope.showReval = (showReval[0]) ? showReval[0].value : 0;

	$scope.ccDescriptions = propertiesModel.properties.cacheChecks.descriptions;

	$scope.ssh = serverUtils.ssh;

	$scope.editCache = function(id) {
		locationUtils.navigateToPath('/servers/' + id);
	};

	$scope.refresh = function() {
		$state.reload(); // reloads all the resolves for the view
	};

	$scope.$on("$destroy", function() {
		killInterval();
	});

	angular.element(document).ready(function () {
		$('#cacheChecksTable').dataTable({
			"aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
			"iDisplayLength": -1
		});
	});

	var init = function () {
		if (autoRefresh) {
			createInterval();
		}
	};
	init();

};

CacheChecksController.$inject = ['cacheChecks', 'showReval', '$scope', '$state', '$interval', 'locationUtils', 'serverUtils', 'propertiesModel'];
module.exports = CacheChecksController;
