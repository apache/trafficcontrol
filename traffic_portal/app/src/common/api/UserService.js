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

var UserService = function(Restangular, $http, $location, $q, authService, httpService, locationUtils, userModel, messageModel, ENV) {

    var service = this;

    this.getCurrentUser = function() {
        var deferred = $q.defer();

        $http.get(ENV.api['root'] + "user/current")
            .then(
                function(result) {
                    userModel.setUser(result.data.response);
                    deferred.resolve(result.data.response);
                },
                function(fault) {
                    deferred.reject(fault);
                }
            );

        return deferred.promise;
    };

    this.resetPassword = function(email) {
        return $http.post(ENV.api['root'] + "user/reset_password", { email: email })
            .then(
                function(result) {
                    messageModel.setMessages(result.data.alerts, false);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
            );
    };

    this.getUsers = function(queryParams) {
        return Restangular.all('users').getList(queryParams);
    };

    this.getUser = function(id) {
        return Restangular.one("users", id).get();
    };

    this.createUser = function(user) {
        return Restangular.service('users').post(user)
            .then(
                function(result) {
                    messageModel.setMessages([ { level: 'success', text: 'User created' } ], true);
                    locationUtils.navigateToPath('/users');
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
            );
    };

    this.updateUser = function(user) {
        return $http.put(ENV.api['root'] + "users/" + user.id, user)
            .then(
                function(result) {
                    if (userModel.user.id == user.id) {
                        // if you are updating the currently logged in user...
                        userModel.setUser(user);
                    }
                    messageModel.setMessages(result.data.alerts, false);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
            );
    };

    this.updateCurrentUser = function(user) {
        return $http.post(ENV.api['root'] + "user/current/update", { user: user })
            .then(
                function() {
                    userModel.setUser(user);
                    messageModel.setMessages([ { level: 'success', text: 'Current user updated' } ], false);
                },
                function() {
                    messageModel.setMessages([ { level: 'error', text: 'Current user updated failed' } ], false);
                }
            );
    };

    this.getUnassignedUserDeliveryServices = function(userId) {
        var deferred = $q.defer();

        $http.get(ENV.api['root'] + "user/" + userId + "/deliveryservices/available")
            .then(
                function(result) {
                    deferred.resolve(result.data.response);
                },
                function(fault) {
                    deferred.reject(fault);
                }
            );

        return deferred.promise;
    };

    this.deleteUserDeliveryService = function(userId, dsId) {
        return httpService.delete(ENV.api['root'] + 'deliveryservice_user/' + dsId + '/' + userId)
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'User and delivery service were unlinked.' } ], false);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, true);
                }
            );
    };

    this.assignUserDeliveryServices = function(userId, deliveryServices) {
        return Restangular.service('deliveryservice_user').post( { userId: userId, deliveryServices: deliveryServices, replace: true } )
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Delivery services linked to user' } ], false);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
            );
    };

    this.registerUser = function(registration) {
        return $http.post(ENV.api['root'] + "users/register", registration)
            .then(
                function(result) {
                    messageModel.setMessages(result.data.alerts, false);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
            );
    };


};

UserService.$inject = ['Restangular', '$http', '$location', '$q', 'authService', 'httpService', 'locationUtils', 'userModel', 'messageModel', 'ENV'];
module.exports = UserService;
