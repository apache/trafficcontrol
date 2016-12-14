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

var ServerService = function(Restangular, locationUtils, messageModel) {

    this.getServers = function(dsId, profileId) {
        return Restangular.all('servers').getList({ dsId: dsId, profileId: profileId });
    };

    this.getServer = function(id) {
        return Restangular.one("servers", id).get();
    };

    this.createServer = function(server) {
        return Restangular.service('servers').post(server)
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Server created' } ], true);
                    locationUtils.navigateToPath('/configure/servers');
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
            );
    };

    this.updateServer = function(server) {
        return server.put()
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Server updated' } ], false);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
            );
    };

    this.deleteServer = function(id) {
        return Restangular.one("servers", id).remove()
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Server deleted' } ], true);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, true);
                }
            );
    };

};

ServerService.$inject = ['Restangular', 'locationUtils', 'messageModel'];
module.exports = ServerService;
