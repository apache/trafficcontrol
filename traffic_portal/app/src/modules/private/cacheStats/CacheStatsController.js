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

var CacheStatsController = function(cacheStats, $scope, $state, $interval, numberUtils, serverUtils) {

	var cacheStatsInterval,
		autoRefresh = false,
		refreshRateInMS = 10000;

	var createInterval = function() {
		killInterval();
		cacheStatsInterval = $interval(function() { $scope.refresh() }, refreshRateInMS );
	};

	var killInterval = function() {
		if (angular.isDefined(cacheStatsInterval)) {
			$interval.cancel(cacheStatsInterval);
			cacheStatsInterval = undefined;
		}
	};

	$scope.cacheStats = cacheStats;

	$scope.ssh = serverUtils.ssh;

	$scope.bandwidth = function(kbps, unit) {
		return numberUtils.addCommas(numberUtils.convertTo(kbps, unit));
	};

	$scope.connections = function(amt) {
		return numberUtils.addCommas(amt);
	};

	$scope.refresh = function() {
		$state.reload(); // reloads all the resolves for the view
	};

	$scope.$on("$destroy", function() {
		killInterval();
	});

	angular.element(document).ready(function () {
		$('#cacheStatsTable').dataTable({
			"aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
			"iDisplayLength": 25,
			"order": [ [ 6, "desc" ] ] // sort by bandwidth, descending
		});
	});

	var init = function () {
		if (autoRefresh) {
			createInterval();
		}
	};
	init();

};

CacheStatsController.$inject = ['cacheStats', '$scope', '$state', '$interval', 'numberUtils', 'serverUtils'];
module.exports = CacheStatsController;
