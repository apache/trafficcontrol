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
var DivisionService = function($http, ENV, locationUtils, messageModel) {

    this.getDivisions = function(queryParams) {
        return $http.get(ENV.api.unstable + 'divisions', {params: queryParams}).then(
            function(result) {
                return result.data.response;
            },
            function(err) {
                throw err;
            }
        );
    };

    this.getDivision = function(id) {
        return $http.get(ENV.api.unstable + 'divisions', {params: {id: id}}).then(
            function(result) {
                return result.data.response[0];
            },
            function(err) {
                throw err;
            }
        );
    };

    this.createDivision = function(division) {
        return $http.post(ENV.api.unstable + 'divisions', division).then(
            function(result) {
                messageModel.setMessages(result.data.alerts, true);
                locationUtils.navigateToPath('/divisions');
                return result;
            },
            function(err) {
                messageModel.setMessages(err.data.alerts, false);
                throw err;
            }
        );
    };

    // todo: change to use query param when it is supported
    this.updateDivision = function(division) {
        return $http.put(ENV.api.unstable + 'divisions/' + division.id, division).then(
            function(result) {
                messageModel.setMessages(result.data.alerts, false);
                return result;            },
            function(err) {
                messageModel.setMessages(err.data.alerts, false);
                throw err;
            }
        );
    };

    // todo: change to use query param when it is supported
    this.deleteDivision = function(id) {
        return $http.delete(ENV.api.unstable + 'divisions/' + id).then(
                function(result) {
                    messageModel.setMessages(result.data.alerts, true);
                    return result;
                },
                function(err) {
                    messageModel.setMessages(err.data.alerts, false);
                    throw err;
                }
            );
    };

};

DivisionService.$inject = ['$http', 'ENV', 'locationUtils', 'messageModel'];
module.exports = DivisionService;
