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

var UserService = function($http, locationUtils, userModel, messageModel, ENV) {

    this.getCurrentUser = function() {
        return $http.get(ENV.api.unstable + "user/current").then(
            function(result) {
                userModel.setUser(result.data.response);
                return result.data.response;
            },
            function(err) {
                throw err;
            }
        );
    };

    this.resetPassword = function(email) {
        return $http.post(ENV.api.unstable + "user/reset_password", { email: email }).then(
            function(result) {
                messageModel.setMessages(result.data.alerts, false);
                return result;
            },
            function(err) {
                messageModel.setMessages(err.data.alerts, false);
                throw err;
            }
        );
    };

    this.getUsers = function(queryParams) {
        return $http.get(ENV.api.unstable + 'users', {params: queryParams}).then(
            function (result) {
                return result.data.response;
            },
            function (err) {
                console.error(err);
                throw err;
            }
        )
    };

    this.getUser = function(id) {
        return $http.get(ENV.api.unstable + 'users', {params: {id: id}}).then(
            function (result) {
                return result.data.response[0];
            },
            function (err) {
                console.error(err);
                throw err;
            }
        )
    };

    this.createUser = function(user) {
        return $http.post(ENV.api.unstable + 'users', user).then(
            function(result) {
                messageModel.setMessages([ { level: 'success', text: 'User created' } ], true);
                locationUtils.navigateToPath('/users');
                return result;
            },
            function(err) {
                messageModel.setMessages(err.data.alerts, false);
                throw err;
            }
        );
    };

    this.updateCurrentUser = function(user) {
        // We should be using PUT 'user/current' to update the current user
        const currUser = { user };
        return $http.put(ENV.api.unstable + 'user/current', currUser).then(
            function(result) {
                userModel.setUser(user);
                messageModel.setMessages(result.data.alerts, false);
                return result;
            },
            function(err) {
                messageModel.setMessages(err.data.alerts, false);
                throw err;
            }
        );
    };

    // todo: change to use query param when it is supported
    this.updateUser = function(user) {
        if (userModel.user.id === user.id) {
            return this.updateCurrentUser(user);
        } else {
            return $http.put(ENV.api.unstable + "users/" + user.id, user).then(
                function(result) {
                    messageModel.setMessages(result.data.alerts, false);
                    return result;
                },
                function(err) {
                    messageModel.setMessages(err.data.alerts, false);
                    throw err;
                }
            );
        }
    };

    this.registerUser = function(registration) {
        return $http.post(ENV.api.unstable + "users/register", registration).then(
            function(result) {
                messageModel.setMessages(result.data.alerts, false);
                return result;
            },
            function(err) {
                messageModel.setMessages(err.data.alerts, false);
                throw err;
            }
        );
    };
};

UserService.$inject = ['$http', 'locationUtils', 'userModel', 'messageModel', 'ENV'];
module.exports = UserService;
