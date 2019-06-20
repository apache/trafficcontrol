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

var OriginService = function($http, $q, Restangular, locationUtils, messageModel, ENV) {

    this.getOrigins = function(queryParams) {
        return Restangular.all('origins').getList(queryParams);
    };

    this.createOrigin = function(origin) {
        var request = $q.defer();

        $http.post(ENV.api['root'] + "origins", origin)
            .then(
                function(response) {
                    messageModel.setMessages(response.data.alerts, true);
                    locationUtils.navigateToPath('/origins');
                    request.resolve(response);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false)
                    request.reject(fault);
                }
            );

        return request.promise;
    };

    this.updateOrigin = function(id, origin) {
        var request = $q.defer();

        $http.put(ENV.api['root'] + "origins?id=" + id, origin)
            .then(
                function(response) {
                    messageModel.setMessages(response.data.alerts, false);
                    request.resolve();
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                    request.reject();
                }
            );
        return request.promise;
    };

    this.deleteOrigin = function(id) {
        var deferred = $q.defer();

        $http.delete(ENV.api['root'] + "origins?id=" + id)
            .then(
                function(response) {
                    messageModel.setMessages(response.data.alerts, true);
                    deferred.resolve(response);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                    deferred.reject(fault);
                }
            );
        return deferred.promise;
    };

};

OriginService.$inject = ['$http', '$q', 'Restangular', 'locationUtils', 'messageModel', 'ENV'];
module.exports = OriginService;
