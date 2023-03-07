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
var ASNService = function($http, locationUtils, messageModel, ENV) {

    this.getASNs = function(queryParams) {
        return $http.get(ENV.api.unstable + 'asns', {params: queryParams}).then(
            function(result) {
                return result.data.response;
            },
            function(err) {
                console.error(err);
                throw err;
            }
        );
    };

    this.getASN = function(id) {
        return $http.get(ENV.api.unstable + 'asns', {params: {id: id}}).then(
            function(result) {
                return result.data.response[0];
            },
            function(err) {
                console.error(err);
                throw err;
            }
        );
    };

    this.createASN = function(asn) {
        return $http.post(ENV.api.unstable + 'asns', asn).then(
            function(result) {
                messageModel.setMessages([{level: 'success', text: 'ASN created' }], true);
                console.info("created new ASN: ", result.data.response);
                locationUtils.navigateToPath('/asns');
                return result;
            },
            function(err) {
                messageModel.setMessages(err.data.alerts, false);
                throw err;
            }
        );
    };

    this.updateASN = function(asn) {
        return $http.put(ENV.api.unstable + 'asns', asn, {params: {id: asn.id}}).then(
            function(result) {
                messageModel.setMessages([{level: 'success', text: 'ASN updated'}], false);
                console.info('updated ASN: ', result.data.response);
                return result;
            },
            function(err) {
                messageModel.setMessages(err.data.alerts, false);
                throw err;
            }
        );
    };

    this.deleteASN = function(id) {
        return $http.delete(ENV.api.unstable + 'asns', {params: {id: id}}).then(
            function(result) {
                messageModel.setMessages(result.data.alerts, true);
                return result;
            },
            function(err) {
                messageModel.setMessages(err.data.alerts, true);
                throw err;
            }
        );
    };

};

ASNService.$inject = ['$http', 'locationUtils', 'messageModel', 'ENV'];
module.exports = ASNService;
