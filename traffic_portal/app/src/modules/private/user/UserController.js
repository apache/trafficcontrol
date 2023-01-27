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
 * @param {import("angular").ILocationService} $location
 * @param {import("../../../common/service/utils/angular.ui.bootstrap").IModalService} $uibModal
 * @param {import("../../../common/service/utils/FormUtils")} formUtils
 * @param {import("../../../common/service/utils/LocationUtils")} locationUtils
 * @param {import("../../../common/service/utils/TenantUtils")} tenantUtils
 * @param {import("../../../common/api/UserService")} userService
 * @param {import("../../../common/api/AuthService")} authService
 * @param {import("../../../common/api/RoleService")} roleService
 * @param {import("../../../common/api/TenantService")} tenantService
 * @param {import("../../../common/models/UserModel")} userModel
 */
var UserController = function($scope, $location, $uibModal, formUtils, locationUtils, tenantUtils, userService, authService, roleService, tenantService, userModel) {

    var updateUser = function(user, options) {
        userService.updateUser(user).
            then(function() {
                if (options.signout) {
                    authService.logout();
                }
            });
    };

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

    $scope.userName = angular.copy(userModel.user.username);

    $scope.user = userModel.user;

    $scope.confirmSave = function(user, usernameField) {
        if (usernameField === undefined) {
            usernameField = user.username;
        }
        if (usernameField.$dirty) {
            var params = {
                title: 'Reauthentication Required',
                message: 'Changing your username to ' + user.username + ' will require you to reauthenticate. Is that OK?'
            };
            var modalInstance = $uibModal.open({
                templateUrl: 'common/modules/dialog/confirm/dialog.confirm.tpl.html',
                controller: 'DialogConfirmController',
                size: 'sm',
                resolve: {
                    params: function () {
                        return params;
                    }
                }
            });
            modalInstance.result.then(function() {
                updateUser(user, { signout : true });
            }, function () {
                // do nothing
            });
        } else {
            updateUser(user, { signout : false });
        }
    };

    $scope.viewDeliveryServices = function() {
        $location.path('/users/' + $scope.user.id + '/delivery-services');
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

UserController.$inject = ['$scope', '$location', '$uibModal', 'formUtils', 'locationUtils', 'tenantUtils', 'userService', 'authService', 'roleService', 'tenantService', 'userModel'];
module.exports = UserController;
