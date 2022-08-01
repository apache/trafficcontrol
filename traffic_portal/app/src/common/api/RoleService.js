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

var RoleService = function($http, messageModel, ENV) {

    this.getRoles = function(queryParams) {
        return $http.get(ENV.api.unstable + 'roles', {params: queryParams}).then(
            function (result) {
                return result.data.response;
            },
            function (err) {
                throw err;
            }
        );
    };

    this.createRole = function(role) {
        return $http.post(ENV.api.unstable + "roles", role).then(
            function(result) {
                return result.data;
            },
            function(err) {
                messageModel.setMessages(err.data.alerts, false);
                throw err;
            }
        );
    };

    this.updateRole = function(role, oldRole) {
        return $http.put(ENV.api.unstable + "roles", role, {params: {name: oldRole}}).then(
            function(result) {
                return result.data;
            },
            function(err) {
                messageModel.setMessages(err.data.alerts, false);
                throw err;
            }
        );
    };

    this.deleteRole = function(name) {
        return $http.delete(ENV.api.unstable + "roles", {params: {name: name}}).then(
            function(result) {
                return result.data;
            },
            function(err) {
                messageModel.setMessages(err.data.alerts, false);
                throw err;
            }
        );
    };

};

RoleService.$inject = ['$http', 'messageModel', 'ENV'];
module.exports = RoleService;
