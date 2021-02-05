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

    var sendRegistration = function(user) {
        userService.registerUser(user).
            then(function() {
                $scope.userEmail = angular.copy(user.email);
                $anchorScroll(); // scrolls window to top
            });
    };

    $scope.userName = angular.copy(user.username);

    $scope.settings = {
        isNew: false,
        saveLabel: 'Update'
    };

    $scope.confirmSave = function(user) {
        saveUser(user);
    };

    $scope.resendRegistration = function(user) {
        sendRegistration(user);
    };

};

FormEditUserController.$inject = ['user', '$scope', '$controller', '$uibModal', '$anchorScroll', 'locationUtils', 'userService'];
module.exports = FormEditUserController;
