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

var LoginController = function($scope, $log, $uibModal, $location, authService, userService, propertiesModel) {

    $scope.oAuthEnabled = propertiesModel.properties.oAuth.enabled;

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

    $scope.loginOauth = function() {
        const redirectUriParamKey = propertiesModel.properties.oAuth.redirectUriParameterOverride !== '' ? propertiesModel.properties.oAuth.redirectUriParameterOverride : 'redirect_uri';
        const redirectParam = $location.search()['redirect'] !== undefined ? $location.search()['redirect'] : '';

        // Builds redirect_uri parameter value to be sent with request to OAuth provider.  This will redirect to the /sso page with any previous redirect information
        var redirectUriParam = new URL(window.location.href.replace(window.location.hash, '') + 'sso');

        // Builds the URL to redirect to the OAuth provider including the redirect_uri (or override), client_id, and response_type fields
        var continueURL = new URL(propertiesModel.properties.oAuth.oAuthUrl);
        continueURL.searchParams.append(redirectUriParamKey, redirectUriParam);
        continueURL.searchParams.append('client_id', propertiesModel.properties.oAuth.clientId);
        continueURL.searchParams.append('response_type', 'code');
        continueURL.searchParams.append('scope', 'openid profile email');

        localStorage.setItem('redirectUri', redirectUriParam.toString());
        localStorage.setItem('redirectParam', redirectParam);

        window.location.href = continueURL.href;
    };

    var init = function() {};
    init();
};

LoginController.$inject = ['$scope', '$log', '$uibModal', '$location', 'authService', 'userService', 'propertiesModel'];
module.exports = LoginController;
