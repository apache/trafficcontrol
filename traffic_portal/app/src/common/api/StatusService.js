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

var StatusService = function(Restangular, locationUtils, messageModel) {

    this.getStatuses = function(queryParams) {
        return Restangular.all('statuses').getList(queryParams);
    };

    this.getStatus = function(id) {
        return Restangular.one("statuses", id).get();
    };

    this.createStatus = function(status) {
        return Restangular.service('statuses').post(status)
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Status created' } ], true);
                    locationUtils.navigateToPath('/statuses');
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
            );
    };

    this.updateStatus = function(status) {
        return status.put()
            .then(
            function() {
                messageModel.setMessages([ { level: 'success', text: 'Status updated' } ], false);
            },
            function(fault) {
                messageModel.setMessages(fault.data.alerts, false);
            }
        );
    };

    this.deleteStatus = function(id) {
        return Restangular.one("statuses", id).remove()
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Status deleted' } ], true);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, true);
                }
        );
    };

};

StatusService.$inject = ['Restangular', 'locationUtils', 'messageModel'];
module.exports = StatusService;
