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
 * @param {import("../../../../service/utils/FormUtils")} formUtils
 * @param {import("../../../../service/utils/TenantUtils")} tenantUtils
 * @param {import("../../../../service/utils/LocationUtils")} locationUtils
 * @param {import("../../../../api/RoleService")} roleService
 * @param {import("../../../../api/TenantService")} tenantService
 * @param {import("../../../../api/UserService")} userService
 * @param {import("../../../../models/UserModel")} userModel
 */
var FormRegisterUserController = function($scope, formUtils, tenantUtils, locationUtils, roleService, tenantService, userService, userModel) {

	var getRoles = function() {
		roleService.getRoles()
			.then(function(result) {
				$scope.roles = _.sortBy(result, 'name');
			});
	};

	var getTenants = function() {
		tenantService.getTenant(userModel.user.tenantId)
			.then(function(tenant) {
				tenantService.getTenants()
					.then(function(tenants) {
						$scope.tenants = tenantUtils.hierarchySort(tenantUtils.groupTenantsByParent(tenants), tenant.parentId, []);
						tenantUtils.addLevels($scope.tenants);
					});
			});
	};

	$scope.registration = {};

	$scope.register = function(registration) {
		userService.registerUser(registration);
	};

	$scope.roleLabel = function(role) {
		return role.name;
	};

	$scope.navigateToPath = (path, unsavedChanges) => locationUtils.navigateToPath(path, unsavedChanges);

	$scope.hasError = formUtils.hasError;

	$scope.hasPropertyError = formUtils.hasPropertyError;

	var init = function () {
		getRoles();
		getTenants();
	};
	init();

};

FormRegisterUserController.$inject = ['$scope', 'formUtils', 'tenantUtils', 'locationUtils', 'roleService', 'tenantService', 'userService', 'userModel'];
module.exports = FormRegisterUserController;
