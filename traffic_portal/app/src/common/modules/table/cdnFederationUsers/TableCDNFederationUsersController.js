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
 * @param {*} cdn
 * @param {*} federation
 * @param {*} users
 * @param {*} $scope
 * @param {*} $state
 * @param {import("../../../service/utils/angular.ui.bootstrap").IModalService} $uibModal
 * @param {import("../../../service/utils/LocationUtils")} locationUtils
 * @param {import("../../../api/FederationService")} federationService
 */
var TableCDNFederationUsersController = function(cdn, federation, users, $scope, $state, $uibModal, locationUtils, federationService) {

	var removeUser = function(userId) {
		federationService.deleteFederationUser($scope.federation.id, userId)
			.then(
				function() {
					$scope.refresh();
				}
			);
	};

	$scope.cdn = cdn;

	$scope.federation = federation;

	$scope.users = users;

	$scope.selectUsers = function() {
		var modalInstance = $uibModal.open({
			templateUrl: 'common/modules/table/cdnFederationUsers/table.assignFederationUsers.tpl.html',
			controller: 'TableAssignFederationUsersController',
			size: 'lg',
			resolve: {
				federation: function() {
					return federation;
				},
				users: function(userService) {
					return userService.getUsers();
				},
				assignedUsers: function() {
					return users;
				}
			}
		});
		modalInstance.result.then(function(selectedUserIds) {
			federationService.assignFederationUsers(federation.id, selectedUserIds, true)
				.then(
					function() {
						$scope.refresh();
					}
				);
		}, function () {
			// do nothing
		});
	};

	$scope.confirmRemoveUser = function(user) {
		var params = {
			title: 'Remove User from Federation?',
			message: 'Are you sure you want to remove ' + user.username + ' from this federation?'
		};
		var modalInstance = $uibModal.open({
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
			removeUser(user.id);
		}, function () {
			// do nothing
		});
	};

	$scope.refresh = function() {
		$state.reload(); // reloads all the resolves for the view
	};

	$scope.navigateToPath = (path, unsavedChanges) => locationUtils.navigateToPath(path, unsavedChanges);

	angular.element(document).ready(function () {
		$('#federationUsersTable').dataTable({
			"aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
			"iDisplayLength": 25,
			"columnDefs": [
				{ 'orderable': false, 'targets': 5 },
				{ "width": "5%", "targets": 5 }
			],
			"aaSorting": []
		});
	});

};

TableCDNFederationUsersController.$inject = ['cdn', 'federation', 'users', '$scope', '$state', '$uibModal', 'locationUtils', 'federationService'];
module.exports = TableCDNFederationUsersController;
