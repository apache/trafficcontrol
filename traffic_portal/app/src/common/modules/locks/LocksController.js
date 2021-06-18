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

var LocksController = function($scope, $rootScope, $interval, $state, $uibModal, cdnService, userModel) {

	let interval;

	let getLocks = function() {
		cdnService.getLocks()
			.then(function(result) {
				$scope.locks = result;
			});
	};

	let createInterval = function() {
		interval = $interval(function() { getLocks() }, 30000 );
	};

	let killInterval = function() {
		if (angular.isDefined(interval)) {
			$interval.cancel(interval);
			interval = undefined;
		}
	};

	$scope.loggedInUser = userModel.user.username;

	$scope.confirmUnlock = function(lock) {
		const params = {
			title: 'Remove lock from: ' + lock.cdn,
			message: 'Are you sure you want to remove the lock from the ' + lock.cdn + ' CDN?'
		};
		const modalInstance = $uibModal.open({
			templateUrl: 'common/modules/dialog/confirm/dialog.confirm.tpl.html',
			controller: 'DialogConfirmController',
			size: 'md',
			resolve: {
				params: function () {
					return params;
				}
			}
		});
		modalInstance.result.then(function() {
			cdnService.deleteLock({ cdn: lock.cdn }).
				then(
					function() {
						$state.reload();
					}
				);
		}, function () {
			// do nothing
		});
	};

	$rootScope.$on('authService::login', function() {
		getLocks();
		createInterval();
	});

	$rootScope.$on('trafficPortal::exit', function() {
		killInterval();
	});

	let init = function () {
		if (userModel.loaded) {
			getLocks();
			createInterval();
		}
	};
	init();

};

LocksController.$inject = ['$scope', '$rootScope', '$interval', '$state', '$uibModal', 'cdnService', 'userModel'];
module.exports = LocksController;
