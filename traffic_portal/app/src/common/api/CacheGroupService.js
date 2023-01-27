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
 * @param {import("angular").IHttpService} $http
 * @param {import("../service/utils/LocationUtils")} locationUtils
 * @param {import("../models/MessageModel")} messageModel
 * @param {{api: Record<PropertyKey, string>}} ENV
 */
var CacheGroupService = function($http, locationUtils, messageModel, ENV) {

    this.getCacheGroups = function(queryParams) {
        return $http.get(ENV.api.unstable + 'cachegroups', {params: queryParams}).then(
            function(result) {
                return result.data.response;
            },
            function(err) {
                throw err;
            }
        );
    };

    this.getCacheGroup = function(id) {
        return $http.get(ENV.api.unstable + 'cachegroups', {params: {'id': id}}).then(
            function (result) {
                return result.data.response[0];
            },
            function (err) {
                throw err;
            }
        );
    };

    this.createCacheGroup = function(cacheGroup) {
        return $http.post(ENV.api.unstable + 'cachegroups', cacheGroup).then(
            function (result) {
                messageModel.setMessages(result.data.alerts, true);
                locationUtils.navigateToPath('/cache-groups');
                return result;
            },
            function(err) {
                messageModel.setMessages(err.data.alerts, false);
                throw err;
            }
        );
    };

    // todo: change to use query param when it is supported
    this.updateCacheGroup = function(cacheGroup) {
        return $http.put(ENV.api.unstable + 'cachegroups/' + cacheGroup.id, cacheGroup).then(
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

    // todo: change to use query param when it is supported
    this.deleteCacheGroup = function(id) {
        return $http.delete(ENV.api.unstable + "cachegroups/" + id).then(
            function(result) {
                return result.data;
            },
            function(err) {
                messageModel.setMessages(err.data.alerts, false);
                throw err;
            }
        );
    };

    this.queueServerUpdates = function(cgId, cdnId) {
        return $http.post(ENV.api.unstable + 'cachegroups/' + cgId + '/queue_update', {action: "queue", cdnId: cdnId}).then(
            function(result) {
                messageModel.setMessages([{level: 'success', text: 'Queued Cache Group server updates'}], false);
                return result;
            },
            function(err) {
                messageModel.setMessages(err.data.alerts, false);
                throw err;
            }
        );
    };

    this.clearServerUpdates = function(cgId, cdnId) {
        return $http.post(ENV.api.unstable + 'cachegroups/' + cgId + '/queue_update', {action: "dequeue", cdnId: cdnId}).then(
            function(result) {
                messageModel.setMessages([{level: 'success', text: 'Cleared Cache Group server updates'}], false);
                return result;
            },
            function(err) {
                messageModel.setMessages(err.data.alerts, false);
                throw err;
            }
        );
    };

    this.getCacheGroupHealth = function() {
        return $http.get(ENV.api.unstable + "cdns/health").then(
            function(result) {
                return result.data.response;
            },
            function(err) {
                throw err;
            }
        );
    };

};

CacheGroupService.$inject = ['$http', 'locationUtils', 'messageModel', 'ENV'];
module.exports = CacheGroupService;
