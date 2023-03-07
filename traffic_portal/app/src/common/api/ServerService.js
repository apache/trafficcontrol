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

var ServerService = function($http, messageModel, ENV) {

    this.getServers = function(queryParams) {
        return $http.get(ENV.api.unstable + 'servers', {params: queryParams}).then(
            function (result){
                return result.data.response;
            },
            function (err) {
                throw err;
            }
        )
    };

    this.createServer = function(server) {
        return $http.post(ENV.api.unstable + 'servers', server).then(
            function(result) {
                return result;
            },
            function(err) {
                throw err;
            }
        );
    };

    // todo: change to use query param when it is supported
    this.updateServer = function(server) {
        return $http.put(ENV.api.unstable + 'servers/' + server.id, server).then(
            function(result) {
                return result;
            },
            function(err) {
                throw err;
            }
        );
    };

    // todo: change to use query param when it is supported
    this.deleteServer = function(id) {
        return $http.delete(ENV.api.unstable + "servers/" + id).then(
            function(result) {
                return result.data;
            },
            function(err) {
                messageModel.setMessages(err.data.alerts, false);
                throw err;
            }
        );
    };

    this.getServerCapabilities = function(id) {
        return $http.get(ENV.api.unstable + 'server_server_capabilities', { params: { serverId: id } }).then(
            function (result) {
                return result.data.response;
            },
            function (err) {
                throw err;
            }
        )
    };

    this.addServerCapability = function(serverId, capabilityName) {
        return $http.post(ENV.api.unstable + 'server_server_capabilities', { serverId: serverId, serverCapability: capabilityName}).then(
            function(result) {
                return result.data;
            },
            function(err) {
                messageModel.setMessages(err.data.alerts, false);
                throw err;
            }
        );
    };

    this.removeServerCapability = function(serverId, capabilityName) {
        return $http.delete(ENV.api.unstable + 'server_server_capabilities', { params: { serverId: serverId, serverCapability: capabilityName} }).then(
            function(result) {
                return result.data;
            },
            function(err) {
                messageModel.setMessages(err.data.alerts, false);
                throw err;
            }
        );
    };

    this.getEligibleDeliveryServiceServers = function(dsId) {
        return $http.get(ENV.api.unstable + 'deliveryservices/' + dsId + '/servers/eligible').then(
            function (result) {
                return result.data.response;
            },
            function (err) {
                throw err;
            }
        )
    };

    this.assignDeliveryServices = function(server, dsIds, replace, delay) {
        return $http.post(ENV.api.unstable + 'servers/' + server.id + '/deliveryservices', dsIds, {params: {replace: replace}}).then(
            function(result) {
                messageModel.setMessages([ { level: 'success', text: dsIds.length + ' delivery services assigned to ' + server.hostName + '.' + server.domainName } ], delay);
                return result;
            },
            function(err) {
                messageModel.setMessages(err.data.alerts, false);
                throw err;
            }
        );
    };

    this.queueServerUpdates = function(id) {
        return $http.post(ENV.api.unstable + "servers/" + id + '/queue_update', { action: "queue"}).then(
            function(result) {
                messageModel.setMessages([ { level: 'success', text: 'Queued server updates' } ], false);
                return result;
            },
            function(err) {
                messageModel.setMessages(err.data.alerts, false);
                throw err;
            }
        );
    };

    this.clearServerUpdates = function(id) {
        return $http.post(ENV.api.unstable + "servers/" + id + '/queue_update', { action: "dequeue"}).then(
            function(result) {
                messageModel.setMessages([ { level: 'success', text: 'Cleared server updates' } ], false);
                return result;
            },
            function(err) {
                messageModel.setMessages(err.data.alerts, false);
                throw err;
            }
        );
    };

    this.getCacheStats = function() {
        return $http.get(ENV.api.unstable + "caches/stats").then(
            function(result) {
                return result.data.response;
            },
            function(err) {
                throw err;
            }
        );
    };

    this.getCacheChecks = function() {
        return $http.get(ENV.api.unstable + "servercheck").then(
            function(result) {
                return result.data.response;
            },
            function(err) {
                throw err;
            }
        );
    };

    this.updateStatus = function(id, payload) {
        return $http.put(ENV.api.unstable + "servers/" + id + "/status", payload).then(
            function(result) {
                return result;
            },
            function(err) {
                throw err;
            }
        );
    };

};

ServerService.$inject = ['$http', 'messageModel', 'ENV'];
module.exports = ServerService;
