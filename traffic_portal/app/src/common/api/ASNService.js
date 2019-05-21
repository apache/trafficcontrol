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

var ASNService = function(Restangular, locationUtils, messageModel) {

    this.getASNs = function(queryParams) {
        return Restangular.all('asns').getList(queryParams);
    };

    this.getASN = function(id) {
        return Restangular.one("asns", id).get();
    };

    this.createASN = function(asn) {
        return Restangular.service('asns').post(asn)
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'ASN created' } ], true);
                    locationUtils.navigateToPath('/asns');
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
            );
    };

    this.updateASN = function(asn) {
        return asn.put()
            .then(
            function() {
                messageModel.setMessages([ { level: 'success', text: 'ASN updated' } ], false);
            },
            function(fault) {
                messageModel.setMessages(fault.data.alerts, false);
            }
        );
    };

    this.deleteASN = function(id) {
        return Restangular.one("asns", id).remove()
            .then(
            function() {
                messageModel.setMessages([ { level: 'success', text: 'ASN deleted' } ], true);
            },
            function(fault) {
                messageModel.setMessages(fault.data.alerts, true);
            }
        );
    };

};

ASNService.$inject = ['Restangular', 'locationUtils', 'messageModel'];
module.exports = ASNService;
