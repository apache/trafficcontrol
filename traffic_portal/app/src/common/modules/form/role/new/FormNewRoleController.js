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
 * @param {*} roles
 * @param {*} $scope
 * @param {import("angular").IControllerService} $controller
 * @param {import("../../../../service/utils/LocationUtils")} locationUtils
 * @param {import("../../../../api/RoleService")} roleService
 * @param {import("../../../../models/MessageModel")} messageModel
 */
var FormNewRoleController = function(roles, $scope, $controller, locationUtils, roleService, messageModel) {

	// extends the FormRoleController to inherit common methods
	angular.extend(this, $controller('FormRoleController', { roles: roles, $scope: $scope }));

	$scope.roleName = 'New';

	$scope.settings = {
		isNew: true,
		saveLabel: 'Create'
	};

	$scope.confirmSave = function(role) {
		roleService.createRole(role).
			then(function(result) {
				messageModel.setMessages(result.alerts, true);
				locationUtils.navigateToPath('/roles');
		});
	};

};

FormNewRoleController.$inject = ['roles', '$scope', '$controller', 'locationUtils', 'roleService', 'messageModel'];
module.exports = FormNewRoleController;
