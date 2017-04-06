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

var UserService = function(Restangular, $http, $location, $q, authService, locationUtils, userModel, messageModel, ENV) {

    var service = this;

    this.getCurrentUser = function() {
        var token = $location.search().token,
            deferred = $q.defer();

        if (angular.isDefined(token)) {
            $location.search('token', null); // remove the token query param
            authService.tokenLogin(token)
                .then(
                    function(response) {
                        service.getCurrentUser();
                    }
                );
        } else {
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
        }
    };


    this.updateCurrentUser = function(user) {
        return user.put()
            .then(
                function() {
                    userModel.setUser(user);
                    messageModel.setMessages([ { level: 'success', text: 'User updated' } ], false);
                },
                function() {
                    messageModel.setMessages([ { level: 'error', text: 'User updated failed' } ], false);
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
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'User created' } ], true);
                    locationUtils.navigateToPath('/admin/users');
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
            );
    };

    this.updateUser = function(user) {
        return user.put()
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'User updated' } ], false);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
            );
    };

    this.deleteUser = function(id) {
        return Restangular.one("users", id).remove()
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'User deleted' } ], true);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, true);
                }
            );
    };

};

UserService.$inject = ['Restangular', '$http', '$location', '$q', 'authService', 'locationUtils', 'userModel', 'messageModel', 'ENV'];
module.exports = UserService;
