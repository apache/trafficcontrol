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

var TableUserDeliveryServicesController = function(user, userDeliveryServices, $scope, $state, $uibModal, locationUtils, userService) {

	$scope.user = user;

	$scope.userDeliveryServices = userDeliveryServices;

	$scope.removeDS = function(dsId) {
		userService.deleteUserDeliveryService(user.id, dsId)
			.then(
				function() {
					$scope.refresh();
				}
			);
	};

	$scope.selectDSs = function() {
		var modalInstance = $uibModal.open({
			templateUrl: 'common/modules/table/userDeliveryServices/table.userDSUnassigned.tpl.html',
			controller: 'TableUserDSUnassignedController',
			size: 'lg',
			resolve: {
				user: function() {
					return user;
				},
				deliveryServices: function(userService) {
					return userService.getUnassignedUserDeliveryServices(user.id);
				}
			}
		});
		modalInstance.result.then(function(selectedDSIds) {
			console.log(selectedDSIds);
			var userDSAssignments = { userId: user.id, deliveryServices: selectedDSIds };
			userService.assignUserDeliveryServices(userDSAssignments)
				.then(
					function() {
						$scope.refresh();
					}
				);
		}, function () {
			// do nothing
		});
	};

	$scope.refresh = function() {
		$state.reload(); // reloads all the resolves for the view
	};

	$scope.navigateToPath = locationUtils.navigateToPath;

	angular.element(document).ready(function () {
		$('#deliveryServicesTable').dataTable({
			"aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
			"iDisplayLength": 25,
			"aaSorting": []
		});
	});

};

TableUserDeliveryServicesController.$inject = ['user', 'userDeliveryServices', '$scope', '$state', '$uibModal', 'locationUtils', 'userService'];
module.exports = TableUserDeliveryServicesController;
