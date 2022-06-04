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

var DialogSelectLockController = function(cdns, $scope, $uibModalInstance, userService) {
	var getUsers = function() {
		userService.getUsers().then(
			function(response) {
				$scope.users = response;
			},
			function(err) {
				throw err;
			});
	};
	$scope.cdns = cdns.filter(
		function (cdn) {
			// you cannot apply a lock to the 'ALL' cdn
			return cdn.name !== 'ALL';
		}
	);

	$scope.lock = {
		cdn: null,
		soft: true,
		message: null
	};

	$scope.select = function() {
		$uibModalInstance.close($scope.lock);
	};

	$scope.cancel = function () {
		$uibModalInstance.dismiss('cancel');
	};
	var init = function () {
		getUsers();
	};
	init();
};

DialogSelectLockController.$inject = ['cdns', '$scope', '$uibModalInstance', 'userService'];
module.exports = DialogSelectLockController;
