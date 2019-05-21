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

var CacheGroupService = function($http, $q, Restangular, locationUtils, messageModel, ENV) {

    this.getCacheGroups = function(queryParams) {
        return Restangular.all('cachegroups').getList(queryParams);
    };

    this.getCacheGroup = function(id) {
        return Restangular.one("cachegroups", id).get();
    };

    this.createCacheGroup = function(cacheGroup) {
        return Restangular.service('cachegroups').post(cacheGroup)
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'CacheGroup created' } ], true);
                    locationUtils.navigateToPath('/cache-groups');
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
            );
    };

    this.updateCacheGroup = function(cacheGroup) {
        return cacheGroup.put()
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Cache group updated' } ], false);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
            );
    };

    this.deleteCacheGroup = function(id) {
        var request = $q.defer();

        $http.delete(ENV.api['root'] + "cachegroups/" + id)
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

    this.queueServerUpdates = function(cgId, cdnId) {
        return Restangular.one("cachegroups", cgId).customPOST( { action: "queue", cdnId: cdnId }, "queue_update" )
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Queued cache group server updates' } ], false);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
            );
    };

    this.clearServerUpdates = function(cgId, cdnId) {
        return Restangular.one("cachegroups", cgId).customPOST( { action: "dequeue", cdnId: cdnId}, "queue_update" )
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Cleared cache group server updates' } ], false);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
            );
    };

    this.getParameterCacheGroups = function(paramId) {
        // todo: this needs an api: /parameters/:id/cachegroups
        return Restangular.one('parameters', paramId).getList('cachegroups');
    };

    this.getCacheGroupHealth = function() {
        var deferred = $q.defer();

        $http.get(ENV.api['root'] + "cdns/health")
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

};

CacheGroupService.$inject = ['$http', '$q', 'Restangular', 'locationUtils', 'messageModel', 'ENV'];
module.exports = CacheGroupService;
