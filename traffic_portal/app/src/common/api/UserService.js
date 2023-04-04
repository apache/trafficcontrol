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
 * @typedef Alert
 * @property {"error" | "info" | "success" | "warning"} level
 * @property {string} text
 */

/**
 * @typedef User
 * @property {string | null | undefined} addressLine1
 * @property {string | null | undefined} addressLine2
 * @property {number | null | undefined} changeLogCount
 * @property {string | null | undefined} city
 * @property {string | null | undefined} company
 * @property {string | null | undefined} country
 * @property {string} email
 * @property {string} fullName
 * @property {number | null | undefined} gid
 * @property {number | null | undefined} id
 * @property {string | null | undefined} lastAuthenticated
 * @property {string | null | undefined} lastUpdated
 * @property {boolean} newUser
 * @property {string | null | undefined} postalCode
 * @property {string | null | undefined} phoneNumber
 * @property {string | null | undefined} publicSshKey
 * @property {string | null | undefined} registrationSent
 * @property {string} role
 * @property {string | null | undefined} stateOrProvince
 * @property {string | null | undefined} tenant
 * @property {number} tenantId
 * @property {string} ucdn
 * @property {number | null | undefined} uid
 * @property {string} username
 */

/**
 * @typedef UserResponse
 * @property {User} response
 * @property {Alert[] | undefined} alerts
 */

/**
 * @param {import("angular").IHttpService} $http
 * @param {import("../service/utils/LocationUtils")} locationUtils
 * @param {import("../models/UserModel")} userModel
 * @param {import("../models/MessageModel")} messageModel
 * @param {{api: Record<PropertyKey, string>}} ENV
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

    /**
     * Updates the current user to match the one passed in.
     *
     * @param {User} user
     * @returns {Promise<{data: UserResponse & {changeLogCount: number; id: number; lastUpdated: string}}>}
     */
    async function updateCurrentUser(user) {
        let result;
        try {
            result = await $http.put(`${ENV.api.unstable}user/current`, user);
        } catch (err) {
            messageModel.setMessages(err.data.alerts, false);
            throw err;
        }
        userModel.setUser(user);
        messageModel.setMessages(result.data.alerts, false);
        return result;
    }

    /** @type {typeof updateCurrentUser} */
    this.updateCurrentUser = updateCurrentUser;

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
