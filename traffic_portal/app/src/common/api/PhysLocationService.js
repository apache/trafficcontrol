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
 * @param {{api: Record<PropertyKey, string>}} ENV
 * @param {import("../service/utils/LocationUtils")} locationUtils
 * @param {import("../models/MessageModel")} messageModel
 */
var PhysLocationService = function($http, ENV, locationUtils, messageModel) {

    this.getPhysLocations = function(queryParams) {
        return $http.get(ENV.api.unstable + 'phys_locations', {params: queryParams}).then(
            function (result) {
                return result.data.response;
            },
            function (err) {
                throw err;
            }
        );
    };

    this.getPhysLocation = function(id) {
        return $http.get(ENV.api.unstable + 'phys_locations', {params: {id: id}}).then(
            function(result) {
                return result.data.response[0];
            },
            function(err) {
                throw err;
            }
        );
    };

    this.createPhysLocation = function(physLocation) {
        return $http.post(ENV.api.unstable + 'phys_locations', physLocation).then(
            function(result) {
                messageModel.setMessages([ { level: 'success', text: 'Physical location created' } ], true);
                locationUtils.navigateToPath('/phys-locations');
                return result;
            },
            function(err) {
                messageModel.setMessages(err.data.alerts, false);
                throw err;
            }
        );
    };

    // todo: change to use query param when it is supported
    this.updatePhysLocation = function(physLocation) {
        return $http.put(ENV.api.unstable + 'phys_locations/' + physLocation.id, physLocation).then(
            function(result) {
                messageModel.setMessages([ { level: 'success', text: 'Physical location updated' } ], false);
                return result;
            },
            function(err) {
                messageModel.setMessages(err.data.alerts, false);
                throw err;
            }
        );
    };

    // todo: change to use query param when it is supported
    this.deletePhysLocation = function(id) {
        return $http.delete(ENV.api.unstable + 'phys_locations/' + id).then(
            function(result) {
                messageModel.setMessages([ { level: 'success', text: 'Physical location deleted' } ], true);
                return result;
            },
            function(err) {
                messageModel.setMessages(err.data.alerts, false);
                throw err;
            }
        );
    };

};

PhysLocationService.$inject = ['$http', 'ENV', 'locationUtils', 'messageModel'];
module.exports = PhysLocationService;
