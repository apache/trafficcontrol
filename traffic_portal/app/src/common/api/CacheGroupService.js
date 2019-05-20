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

var CacheGroupService = function($http, locationUtils, messageModel, ENV) {

    this.getCacheGroups = function(queryParams) {
        return $http.get(ENV.api['root'] + 'cachegroups', {params: queryParams}).then(
            function(result) {
                return result.data.response;
            },
            function(err) {
                console.error(err);
            }
        );
    };

    this.getCacheGroup = function(id) {
        return $http.get(ENV.api['root'] + 'cachegroups', {params: {'id': id}}).then(
            function (result) {
                return result.data.response[0];
            },
            function (err) {
                console.error(err);
                return err;
            }
        );
    };

    this.createCacheGroup = function(cacheGroup) {
        return $http.post(ENV.api['root'] + 'cachegroups', cacheGroup).then(
            function (result) {
                messageModel.setMessages(result.data.alerts, true);
                locationUtils.navigateToPath('/cache-groups');
            },
            function(err) {
                messageModel.setMessages(err.data.alerts, false);
            }
        );
    };

    this.updateCacheGroup = function(cacheGroup) {
        return $http.put(ENV.api['root'] + 'cachegroups/' + cacheGroup.id, cacheGroup).then(
            function(result) {
                messageModel.setMessages(result.data.alerts, false);
            },
            function(err) {
                messageModel.setMessages(err.data.alerts, false);
            }
        );
    };

    this.deleteCacheGroup = function(id) {
        return $http.delete(ENV.api['root'] + "cachegroups/" + id).then(
            function(result) {
                messageModel.setMessages(result.data.alerts, false);
            },
            function(err) {
                messageModel.setMessages(err.data.alerts, false);
            }
        );
    };

    this.queueServerUpdates = function(cgId, cdnId) {
        return $http.post(ENV.api['root'] + 'cachegroups/' + cgId + '/queue_update', {action: "queue", cdnId: cdnId}).then(
            function() {
                messageModel.setMessages([{level: 'success', text: 'Queued Cache Group server updates'}], false);
            },
            function(err) {
                messageModel.setMessages(err.data.alerts, false);
            }
        );
    };

    this.clearServerUpdates = function(cgId, cdnId) {
        return $http.post(ENV.api['root'] + 'cachegroups/' + cgId + '/queue_update', {action: "dequeue", cdnId: cdnId}).then(
            function() {
                messageModel.setMessages([{level: 'success', text: 'Cleared Cache Group server updates'}], false);
            },
            function(err) {
                messageModel.setMessages(err.data.alerts, false);
            }
        );
    };

    this.getCacheGroupHealth = function() {
        return $http.get(ENV.api['root'] + "cdns/health").then(
            function(result) {
                return result.data.response;
            },
            function(err) {
                console.error(err);
                return err;
            }
        );
    };

};

CacheGroupService.$inject = ['$http', 'locationUtils', 'messageModel', 'ENV'];
module.exports = CacheGroupService;
