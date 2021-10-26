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

var RegionService = function($http, ENV, messageModel) {

    this.getRegions = function(queryParams) {
        return $http.get(ENV.api.unstable + 'regions', {params: queryParams}).then(
            function (result) {
                return result.data.response;
            },
            function (err) {
                throw err;
            }
        );
    };

    this.getRegion = function(id) {
        return $http.get(ENV.api.unstable + 'regions', {params: {id: id}}).then(
            function (result) {
                return result.data.response[0];
            },
            function (err) {
                throw err;
            }
        )
    };

    this.createRegion = function(region) {
        return $http.post(ENV.api.unstable + 'regions', region).then(
            function(result) {
                messageModel.setMessages([ { level: 'success', text: 'Region created' } ], true);
                return result;
            },
            function(err) {
                messageModel.setMessages(err.data.alerts, false);
                throw err;
            }
        );
    };

    // todo: change to use query param when it is supported
    this.updateRegion = function(region) {
        return $http.put(ENV.api.unstable + 'regions/' + region.id, region).then(
            function(result) {
                messageModel.setMessages([ { level: 'success', text: 'Region updated' } ], false);
                return result;
            },
            function(err) {
                messageModel.setMessages(err.data.alerts, false);
                throw err;
            }
        );
    };

    // todo: change to use query param when it is supported
    this.deleteRegion = function(id) {
        return $http.delete(ENV.api.unstable + "regions", {params: {id: id}}).then(
            function(result) {
                messageModel.setMessages([ { level: 'success', text: 'Region deleted' } ], true);
                return result;
            },
            function(err) {
                messageModel.setMessages(err.data.alerts, false);
                throw err;
            }
        );
    };

};

RegionService.$inject = ['$http', 'ENV', 'messageModel'];
module.exports = RegionService;
