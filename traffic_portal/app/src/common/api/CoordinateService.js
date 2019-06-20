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

var CoordinateService = function($http, $q, Restangular, locationUtils, messageModel, ENV) {

    this.getCoordinates = function(queryParams) {
        return Restangular.all('coordinates').getList(queryParams);
    };

    this.createCoordinate = function(coordinate) {
        var request = $q.defer();

        $http.post(ENV.api['root'] + "coordinates", coordinate)
            .then(
                function(response) {
                    messageModel.setMessages(response.data.alerts, true);
                    locationUtils.navigateToPath('/coordinates');
                    request.resolve(response);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false)
                    request.reject(fault);
                }
            );

        return request.promise;
    };

    this.updateCoordinate = function(id, coordinate) {
        var request = $q.defer();

        $http.put(ENV.api['root'] + "coordinates?id=" + id, coordinate)
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

    this.deleteCoordinate = function(id) {
        var deferred = $q.defer();

        $http.delete(ENV.api['root'] + "coordinates?id=" + id)
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

CoordinateService.$inject = ['$http', '$q', 'Restangular', 'locationUtils', 'messageModel', 'ENV'];
module.exports = CoordinateService;
