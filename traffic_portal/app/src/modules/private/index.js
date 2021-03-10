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

module.exports = angular.module('trafficPortal.private', [])
    .controller('PrivateController', require('./PrivateController'))
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficPortal.private', {
                url: '',
                abstract: true,
                views: {
                    navigation: {
                        templateUrl: 'common/modules/navigation/navigation.tpl.html',
                        controller: 'NavigationController'
                    },
                    header: {
                        templateUrl: 'common/modules/header/header.tpl.html',
                        controller: 'HeaderController'
                    },
                    locks: {
                        templateUrl: 'common/modules/locks/locks.tpl.html',
                        controller: 'LocksController'
                    },
                    notifications: {
                        templateUrl: 'common/modules/notifications/notifications.tpl.html',
                        controller: 'NotificationsController'
                    },
                    message: {
                        templateUrl: 'common/modules/message/message.tpl.html',
                        controller: 'MessageController'
                    },
                    content: {
                        templateUrl: 'modules/private/private.tpl.html',
                        controller: 'PrivateController'
                    }
                },
                resolve: {
                    tokenLogin: function($location, authService) {
                        var token = $location.search().token; // if there is a token query param, attempt to login with it
                        if (angular.isDefined(token)) {
                            return authService.tokenLogin(token);
                        }
                    },
                    currentUser: function(tokenLogin, $state, userService, userModel) {
                        if (userModel.loaded) {
                            return userModel.user;
                        } else {
                            return userService.getCurrentUser();
                        }
                    }
                }
            });
        $urlRouterProvider.otherwise('/');
    });
