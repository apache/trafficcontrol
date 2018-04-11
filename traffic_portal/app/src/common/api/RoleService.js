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

var RoleService = function(Restangular, $http, $q, messageModel, ENV) {

    this.getRoles = function(queryParams) {
        return Restangular.all('roles').getList(queryParams);
    };

    this.createRole = function(role) {
        var request = $q.defer();

        $http.post(ENV.api['root'] + "roles", role)
            .then(
                function(result) {
                    request.resolve(result.data);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                    request.reject(fault);
                }
            );

        return request.promise;
    };

    this.updateRole = function(role) {
        var request = $q.defer();

        $http.put(ENV.api['root'] + "roles?id=" + role.id, role)
            .then(
                function(result) {
                    request.resolve(result.data);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                    request.reject();
                }
            );

        return request.promise;
    };

    this.deleteRole = function(id) {
        var request = $q.defer();

        $http.delete(ENV.api['root'] + "roles?id=" + id)
            .then(
                function(result) {
                    request.resolve(result.data);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                    request.reject(fault);
                }
            );

        return request.promise;
    };

};

RoleService.$inject = ['Restangular', '$http', '$q', 'messageModel', 'ENV'];
module.exports = RoleService;
