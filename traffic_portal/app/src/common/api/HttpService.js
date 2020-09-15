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

var HttpService = function($http, $q) {

    this.get = function(resource) {
        const deferred = $q.defer();

        $http.get(resource)
            .then(
                function(result) {
                    deferred.resolve(result);
                },
                function(fault) {
                    deferred.reject(fault);
                }
            );

        return deferred.promise;
    };

    this.post = function(resource, payload) {
        const deferred = $q.defer();

        $http.post(resource, payload)
            .then(
                function(result) {
                    deferred.resolve(result);
                },
                function(fault) {
                    deferred.reject(fault);
                }
            );

        return deferred.promise;
    };

    this.put = function(resource, payload) {
        const deferred = $q.defer();

        $http.put(resource, payload)
            .then(
                function(result) {
                    deferred.resolve(result.response);
                },
                function(fault) {
                    deferred.reject(fault);
                }
            );

        return deferred.promise;
    };

    this.delete = function(resource) {
        const deferred = $q.defer();

        $http.delete(resource)
            .then(
                function(result) {
                    deferred.resolve(result.response);
                },
                function(fault) {
                    deferred.reject(fault);
                }
            );

        return deferred.promise;
    };

};

HttpService.$inject = ['$http', '$q'];
module.exports = HttpService;
