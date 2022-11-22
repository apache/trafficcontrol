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
 *
 * @param {*} roles
 * @param {*} $scope
 * @param {import("angular").IControllerService} $controller
 * @param {import("../../../../service/utils/angular.ui.bootstrap").IModalService} $uibModal
 * @param {import("angular").IAnchorScrollService} $anchorScroll
 * @param {import("angular").ILocationService} $location
 * @param {import("../../../../service/utils/LocationUtils")} locationUtils
 * @param {import("../../../../api/RoleService")} roleService
 * @param {import("../../../../models/MessageModel")} messageModel
 * @param {import("../../../../models/PropertiesModel")} propertiesModel
 */
var FormEditRoleController = function(roles, $scope, $controller, $uibModal, $anchorScroll, $location, locationUtils, roleService, messageModel, propertiesModel) {

	// extends the FormRoleController to inherit common methods
	angular.extend(this, $controller('FormRoleController', { roles: roles, $scope: $scope }));

	var deleteRole = function(role) {
		roleService.deleteRole(role.name)
			.then(function(result) {
				messageModel.setMessages(result.alerts, true);
				locationUtils.navigateToPath('/roles');
			});
	};

	var save = function(role) {
		roleService.updateRole(role, $scope.roleName).
			then(function(result) {
				$scope.roleName = angular.copy(role.name);
				messageModel.setMessages(result.alerts, false);
				$anchorScroll(); // scrolls window to top
			});
	};

	$scope.enforceCapabilities = propertiesModel.properties.enforceCapabilities;

	$scope.roleName = angular.copy($scope.role.name);

	$scope.settings = {
		isNew: false,
		saveLabel: 'Update'
	};

	$scope.viewCapabilities = function() {
		$location.path($location.path() + '/capabilities');
	};

	$scope.viewUsers = function() {
		$location.path($location.path() + '/users');
	};

	$scope.confirmSave = function(role) {
		var params = {
			title: 'Update Role?',
			message: 'Are you sure you want to update the role?'
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
			save(role);
		}, function () {
			// do nothing
		});
	};

	$scope.confirmDelete = function(role) {
		var params = {
			title: 'Delete Role: ' + role.name,
			key: role.name
		};
		var modalInstance = $uibModal.open({
			templateUrl: 'common/modules/dialog/delete/dialog.delete.tpl.html',
			controller: 'DialogDeleteController',
			size: 'md',
			resolve: {
				params: function () {
					return params;
				}
			}
		});
		modalInstance.result.then(function() {
			deleteRole(role);
		}, function () {
			// do nothing
		});
	};

};

FormEditRoleController.$inject = ['roles', '$scope', '$controller', '$uibModal', '$anchorScroll', '$location', 'locationUtils', 'roleService', 'messageModel', 'propertiesModel'];
module.exports = FormEditRoleController;
