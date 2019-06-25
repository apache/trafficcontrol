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

var LoginController = function($scope, $log, $uibModal, authService, userService) {

    $scope.credentials = {
        username: '',
        password: ''
    };

    $scope.login = function(event, credentials) {
        event.stopImmediatePropagation();
        const btn = event.currentTarget;
        btn.disabled = true; // disable the login button to prevent multiple clicks
        authService.login(credentials.username, credentials.password).finally(
            function() {
                btn.disabled = false; // re-enable it
            }
        );
    };

    $scope.resetPassword = function() {

        var modalInstance = $uibModal.open({
            templateUrl: 'common/modules/dialog/reset/dialog.reset.tpl.html',
            controller: 'DialogResetController'
        });

        modalInstance.result.then(function(email) {
            userService.resetPassword(email);
        }, function () {
        });
    };

    var init = function() {};
    init();
};

LoginController.$inject = ['$scope', '$log', '$uibModal', 'authService', 'userService'];
module.exports = LoginController;
