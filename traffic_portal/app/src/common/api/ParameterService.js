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

var ParameterService = function(Restangular, $http, $q, messageModel, ENV) {

    this.getParameters = function(queryParams) {
        return Restangular.all('parameters').getList(queryParams);
    };

    this.getParameter = function(id) {
        return Restangular.one("parameters", id).get();
    };

    this.createParameter = function(parameter) {
        return Restangular.service('parameters').post(parameter)
            .then(
            function(response) {
                messageModel.setMessages([ { level: 'success', text: 'Parameter created' } ], true);
                return response;
            },
            function(fault) {
                messageModel.setMessages(fault.data.alerts, false);
            }
        );
    };

    this.updateParameter = function(parameter) {
        return parameter.put()
            .then(
            function() {
                messageModel.setMessages([ { level: 'success', text: 'Parameter updated' } ], false);
            },
            function(fault) {
                messageModel.setMessages(fault.data.alerts, false);
            }
        );
    };

    this.deleteParameter = function(id) {
        var request = $q.defer();

        $http.delete(ENV.api['root'] + "parameters/" + id)
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


    this.getProfileParameters = function(profileId) {
        return Restangular.one('profiles', profileId).getList('parameters');
    };

    this.getProfileUnassignedParams = function(profileId) {
        return Restangular.one('profiles', profileId).getList('unassigned_parameters');
    };

    this.getCacheGroupUnassignedParams = function(cgId) {
        return Restangular.one('cachegroups', cgId).getList('unassigned_parameters');
    };

};

ParameterService.$inject = ['Restangular', '$http', '$q', 'messageModel', 'ENV'];
module.exports = ParameterService;
