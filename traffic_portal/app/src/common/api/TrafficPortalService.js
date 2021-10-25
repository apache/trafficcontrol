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

var TrafficPortalService = function($http, messageModel, ENV) {

    this.getReleaseVersionInfo = function() {
        return $http.get('traffic_portal_release.json').then(
            function(result) {
                return result;
            },
            function(err) {
                throw err;
            }
        );
    };

    this.getProperties = function() {
        return $http.get('traffic_portal_properties.json').then(
            function(result) {
                return result.data.properties;
            },
            function (err) {
                throw err;
            }
        );
    };

    this.dbDump = function() {
        /*
        responseType=arraybuffer is important if you want to create a blob of your data
        See: https://developer.mozilla.org/en-US/docs/Web/API/XMLHttpRequest/Sending_and_Receiving_Binary_Data
        */
        return $http.get(ENV.api.unstable + 'dbdump', { responseType:'arraybuffer' } ).then(
            function(result) {
                download(result.data, moment().format() + '.pg_dump');
                return result;
            },
            function(err) {
                if (err && err.alerts && err.alerts.length > 0) {
                    messageModel.setMessages(err.alerts, false);
                } else {
                    messageModel.setMessages([ { level: 'error', text: err.status.toString() + ': ' + err.statusText } ], false);
                }
                throw err;
            }
        );
    };

};

TrafficPortalService.$inject = ['$http', 'messageModel', 'ENV'];
module.exports = TrafficPortalService;
