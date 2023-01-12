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
 * @param {*} $scope
 * @param {import("angular").IIntervalService} $interval
 * @param {import("../../../api/CacheGroupService")} cacheGroupService
 * @param {import("../../../service/utils/LocationUtils")} locationUtils
 * @param {import("../../../models/PropertiesModel")} propertiesModel
 */
var WidgetCacheGroupsController = function($scope, $interval, cacheGroupService, locationUtils, propertiesModel) {

	var interval,
		autoRefresh = propertiesModel.properties.dashboard.autoRefresh;

	var getCacheGroupHealth = function() {
		cacheGroupService.getCacheGroupHealth()
			.then(function(result) {
				$scope.cacheGroupHealth = result;
			});
	};

	var createInterval = function() {
		killInterval();
		interval = $interval(function() { getCacheGroupHealth() }, propertiesModel.properties.dashboard.cacheGroupHealth.refreshRateInMS );
	};

	var killInterval = function() {
		if (angular.isDefined(interval)) {
			$interval.cancel(interval);
			interval = undefined;
		}
	};

	// pagination
	$scope.currentCacheGroupsPage = 1;
	$scope.cacheGroupsPerPage = 10;

	$scope.navigateToPath = (path, unsavedChanges) => locationUtils.navigateToPath(path, unsavedChanges);

	$scope.onlinePercent = function(location) {
		return (location.online / (location.online + location.offline)) * 100;
	};

	$scope.$on("$destroy", function() {
		killInterval();
	});

	var init = function() {
		getCacheGroupHealth();
		if (autoRefresh) {
			createInterval();
		}
	};
	init();

};

WidgetCacheGroupsController.$inject = ['$scope', '$interval', 'cacheGroupService', 'locationUtils', 'propertiesModel'];
module.exports = WidgetCacheGroupsController;
