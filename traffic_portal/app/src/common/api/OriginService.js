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

var OriginService = function($http, locationUtils, messageModel, ENV) {

    this.getOrigins = function(queryParams) {
        return $http.get(ENV.api.unstable + 'origins', {params: queryParams}).then(
            function (result) {
                return result.data.response;
            },
            function (err) {
                throw err;
            }
        );
    };

    this.createOrigin = function(origin) {
        return $http.post(ENV.api.unstable + 'origins', origin).then(
            function(response) {
                messageModel.setMessages(response.data.alerts, true);
                locationUtils.navigateToPath('/origins');
                return response;
            },
            function(err) {
                messageModel.setMessages(err.data.alerts, false)
                throw err;
            }
        );
    };

    this.updateOrigin = function(id, origin) {
        return $http.put(ENV.api.unstable + "origins", origin, {params: {id: id}}).then(
            function(response) {
                messageModel.setMessages(response.data.alerts, false);
                return response;
            },
            function(err) {
                messageModel.setMessages(err.data.alerts, false);
                throw err;
            }
        );
    };

    this.deleteOrigin = function(id) {
        return $http.delete(ENV.api.unstable + "origins", {params: {id: id}}).then(
            function(response) {
                messageModel.setMessages(response.data.alerts, true);
                return response;
            },
            function(err) {
                messageModel.setMessages(err.data.alerts, false);
                throw err;
            }
        );
    };

};

OriginService.$inject = ['$http', 'locationUtils', 'messageModel', 'ENV'];
module.exports = OriginService;
