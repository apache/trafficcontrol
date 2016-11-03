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

var FormEditUserController = function(user, $scope, $controller, $uibModal, $anchorScroll, locationUtils, userService) {

    // extends the FormUserController to inherit common methods
    angular.extend(this, $controller('FormUserController', { user: user, $scope: $scope }));

    var saveUser = function(user) {
        userService.updateUser(user).
            then(function() {
                $scope.userName = angular.copy(user.username);
                $anchorScroll(); // scrolls window to top
            });
    };

    var deleteUser = function(user) {
        userService.deleteUser(user.id)
            .then(function() {
                locationUtils.navigateToPath('/admin/users');
            });
    };

    $scope.userName = angular.copy(user.username);

    $scope.settings = {
        showDelete: true,
        saveLabel: 'Update'
    };

    $scope.confirmSave = function(user, usernameField) {
        saveUser(user);
    };

    $scope.confirmDelete = function(user) {
        var params = {
            title: 'Delete User: ' + user.username,
            key: user.username
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
            deleteUser(user);
        }, function () {
            // do nothing
        });
    };

};

FormEditUserController.$inject = ['user', '$scope', '$controller', '$uibModal', '$anchorScroll', 'locationUtils', 'userService'];
module.exports = FormEditUserController;
