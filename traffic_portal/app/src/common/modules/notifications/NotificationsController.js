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

var NotificationsController = function($rootScope, $scope, $interval, userModel, cdnService) {

	let interval;

	let getNotifications = function() {
		cdnService.getNotifications()
			.then(function(result) {
				$scope.notifications = result;
				$rootScope.$broadcast('notificationsController::refreshNotifications', { notifications: result });
			});
	};

	let createInterval = function() {
		interval = $interval(function() { getNotifications() }, 30000 );
	};

	let killInterval = function() {
		if (angular.isDefined(interval)) {
			$interval.cancel(interval);
			interval = undefined;
		}
	};

	$scope.dismissedNotifications = JSON.parse(localStorage.getItem("dismissed_notification_ids")) || [];

	$scope.dismissNotification = function(notification) {
		$scope.dismissedNotifications.push(notification.id);
		localStorage.setItem("dismissed_notification_ids", JSON.stringify($scope.dismissedNotifications));
	};

	$rootScope.$on('authService::login', function() {
		getNotifications();
		createInterval();
	});

	$rootScope.$on('trafficPortal::exit', function() {
		killInterval();
	});

	var init = function () {
		if (userModel.loaded) {
			getNotifications();
			createInterval();
		}
	};
	init();

};

NotificationsController.$inject = ['$rootScope', '$scope', '$interval', 'userModel', 'cdnService'];
module.exports = NotificationsController;
