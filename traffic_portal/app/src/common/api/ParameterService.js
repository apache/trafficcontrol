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
var ParameterService = function($http, locationUtils, messageModel, ENV) {

    this.getParameters = function(queryParams) {
        return $http.get(ENV.api.unstable + 'parameters', {params: queryParams}).then(
            function (result) {
                return result.data.response
            },
            function (err) {
                throw err;
            }
        );
    };

    this.getParameter = function(id) {
        return $http.get(ENV.api.unstable + 'parameters', {params: {id: id}}).then(
            function (result) {
                return result.data.response[0];
            },
            function (err) {
                throw err;
            }
        );
    };

    this.createParameter = function(parameter) {
        return $http.post(ENV.api.unstable + 'parameters', parameter).then(
            function(result) {
                messageModel.setMessages([ { level: 'success', text: 'Parameter created' } ], true);
                locationUtils.navigateToPath('/parameters/' + result.data.response.id + '/profiles');
                return result;
            },
            function(err) {
                messageModel.setMessages(err.data.alerts, false);
                throw err;
            }
        );
    };

    // todo: change to use query param when it is supported
    this.updateParameter = function(parameter) {
        return $http.put(ENV.api.unstable + 'parameters/' + parameter.id, parameter).then(
            function(result) {
                messageModel.setMessages([ { level: 'success', text: 'Parameter updated' } ], false);
                return result;
            },
            function(err) {
                messageModel.setMessages(err.data.alerts, false);
                throw err;
            }
        );
    };

    // todo: change to use query param when it is supported
    this.deleteParameter = function(id) {
        return $http.delete(ENV.api.unstable + "parameters/" + id).then(
            function(result) {
                return result.data;
            },
            function(err) {
                messageModel.setMessages(err.data.alerts, false);
                throw err;
            }
        );
    };


    this.getProfileParameters = function(profileId) {
        return $http.get(ENV.api.unstable + 'profiles/' + profileId + '/parameters').then(
            function (result) {
                return result.data.response;
            },
            function (err) {
                throw err;
            }
        );
    };

};

ParameterService.$inject = ['$http', 'locationUtils', 'messageModel', 'ENV'];
module.exports = ParameterService;
