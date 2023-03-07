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
 * @param {import("angular").IRootScopeService} $rootScope
 * @param {import("angular").IHttpService} $http
 * @param {*} $state
 * @param {import("angular").ILocationService} $location
 * @param {import("../models/UserModel")} userModel
 * @param {import("../models/MessageModel")} messageModel
 * @param {{api: Record<PropertyKey, string>}} ENV
 * @param {import("../service/utils/LocationUtils")} locationUtils
 */
var AuthService = function($rootScope, $http, $state, $location, userModel, messageModel, ENV, locationUtils) {

    this.login = function(username, password) {
        userModel.resetUser();
        return $http.post(ENV.api.unstable + 'user/login', { u: username, p: password }).then(
            function(result) {
                $rootScope.$broadcast('authService::login');
                const redirect = decodeURIComponent($location.search().redirect);
                if (redirect !== 'undefined') {
                    $location.search('redirect', null); // remove the redirect query param
                    $location.url(redirect);
                } else {
                    $location.url('/');
                }
                return result;
            },
            function(err) {
                throw err;
            }
        );
    };

    this.tokenLogin = function(token) {
        userModel.resetUser();
        return $http.post(ENV.api.unstable + "user/login/token", { t: token }).then(
            function(result) {
                $rootScope.$broadcast('authService::login');
                return result;
            },
            function(err) {
                throw err;
            }
        );
    };

    this.oauthLogin = function(authCodeTokenUrl, code, clientId, redirectUri) {
        return $http.post(ENV.api.unstable + 'user/login/oauth', { authCodeTokenUrl: authCodeTokenUrl, code: code, clientId: clientId, redirectUri: redirectUri})
            .then(
                function() {
                    $rootScope.$broadcast('authService::login');
                    let redirect = localStorage.getItem('redirectParam');
                    localStorage.clear();
                    if (!redirect) {
                        redirect = decodeURIComponent($location.search().redirect);
                    }
                    if (redirect !== undefined) {
                        $location.search('redirect', null); // remove the redirect query param
                        $location.url(redirect);
                    } else {
                        $location.url('/');
                    }
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, true);
                    locationUtils.navigateToPath('/login');
                }
            );
    };

    this.logout = function() {
        userModel.resetUser();
        return $http.post(`${ENV.api.unstable}user/logout`, undefined).then(
            function(result) {
                $rootScope.$broadcast('trafficPortal::exit');
                if ($state.current.name == 'trafficPortal.public.login') {
                    messageModel.setMessages(result.data.alerts, false);
                } else {
                    messageModel.setMessages(result.data.alerts, true);
                    $state.go('trafficPortal.public.login');
                }
                return result;
            },
            function(err) {
                throw err;
            }
        );
    };

};

AuthService.$inject = ['$rootScope', '$http', '$state', '$location', 'userModel', 'messageModel', 'ENV', 'locationUtils'];
module.exports = AuthService;
