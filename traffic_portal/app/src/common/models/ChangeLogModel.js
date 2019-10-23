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

var ChangeLogModel = function($rootScope, $interval, changeLogService, userModel) {

	var newLogCount = 0,
		pollingIntervalInSecs = 30,
		changeLogInterval;

	this.newLogCount = function() {
		return newLogCount;
	};

	var createChangeLogInterval = function() {
		killChangeLogInterval();
		changeLogInterval = $interval(function() { getNewLogCount() }, (pollingIntervalInSecs*1000)); // every X minutes
	};

	var killChangeLogInterval = function() {
		if (angular.isDefined(changeLogInterval)) {
			$interval.cancel(changeLogInterval);
			changeLogInterval = undefined;
		}
	};

	var getNewLogCount = function() {
		changeLogService.getNewLogCount()
			.then(function(response) {
				newLogCount = response.newLogcount;
			});
	};

	$rootScope.$on('authService::login', function() {
		getNewLogCount();
		createChangeLogInterval();
	});

	$rootScope.$on('trafficPortal::exit', function() {
		killChangeLogInterval();
	});

	$rootScope.$on('changeLogService::getChangeLogs', function() {
		newLogCount = 0;
	});

	var init = function () {
		if (userModel.loaded) {
			getNewLogCount();
			createChangeLogInterval();
		}
	};
	init();

};

ChangeLogModel.$inject = ['$rootScope', '$interval', 'changeLogService', 'userModel'];
module.exports = ChangeLogModel;
