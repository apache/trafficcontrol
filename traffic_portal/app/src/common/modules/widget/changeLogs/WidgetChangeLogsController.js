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

/** @typedef {import("moment")} moment */

/**
 * @param {*} $scope
 * @param {import("../../../service/utils/LocationUtils")} locationUtils
 * @param {import("../../../api/ChangeLogService")} changeLogService
 */
var WidgetChangeLogsController = function($scope, locationUtils, changeLogService) {

	var getChangeLogs = function() {
		changeLogService.getChangeLogs({ limit: 5 })
			.then(function(response) {
				$scope.changeLogs = response;
			});
	};

	$scope.getRelativeTime = function(date) {
		return moment(date).fromNow();
	};

	$scope.navigateToPath = (path, unsavedChanges) => locationUtils.navigateToPath(path, unsavedChanges);

	var init = function() {
		getChangeLogs();
	};
	init();

};

WidgetChangeLogsController.$inject = ['$scope', 'locationUtils', 'changeLogService'];
module.exports = WidgetChangeLogsController;
