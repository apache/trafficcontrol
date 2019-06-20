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

var ServerService = function($http, $q, Restangular, locationUtils, messageModel, ENV) {

    this.getServers = function(queryParams) {
        return Restangular.all('servers').getList(queryParams);
    };

    this.getServer = function(id) {
        return Restangular.one("servers", id).get();
    };

    this.createServer = function(server) {
        return Restangular.service('servers').post(server)
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Server created' } ], true);
                    locationUtils.navigateToPath('/servers');
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
        var request = $q.defer();

        $http.delete(ENV.api['root'] + "servers/" + id)
            .then(
                function(result) {
                    request.resolve(result.data);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                    request.reject(fault);
                }
            );

        return request.promise;
    };

    this.getServerConfigFiles = function(id) {
        return Restangular.one("servers", id).customGET('configfiles/ats');
    };

    this.getServerConfigFile = function(url) {
        var request = $q.defer();

        $http.get(url)
            .then(
                function(result) {
                    request.resolve(result.data);
                },
                function() {
                    request.reject();
                }
            );

        return request.promise;
    };

    this.getDeliveryServiceServers = function(dsId) {
        return Restangular.one('deliveryservices', dsId).getList('servers');
    };

    this.getUnassignedDeliveryServiceServers = function(dsId) {
        return Restangular.one('deliveryservices', dsId).getList('servers/unassigned');
    };

    this.getEligibleDeliveryServiceServers = function(dsId) {
        return Restangular.one('deliveryservices', dsId).getList('servers/eligible');
    };

    this.assignDeliveryServices = function(server, dsIds, replace, delay) {
        return Restangular.service('servers/' + server.id + '/deliveryservices?replace=' + replace).post( dsIds )
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: dsIds.length + ' delivery services assigned to ' + server.hostName + '.' + server.domainName } ], delay);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
            );
    };

    this.queueServerUpdates = function(id) {
        return Restangular.one("servers", id).customPOST( { action: "queue"}, "queue_update" )
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Queued server updates' } ], false);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
            );
    };

    this.clearServerUpdates = function(id) {
        return Restangular.one("servers", id).customPOST( { action: "dequeue"}, "queue_update" )
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Cleared server updates' } ], false);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
            );
    };

    this.getEdgeStatusCount = function() {
        var request = $q.defer();

        $http.get(ENV.api['root'] + "servers/status?type=EDGE")
            .then(
                function(result) {
                    request.resolve(result.data.response);
                },
                function() {
                    request.reject();
                }
            );

        return request.promise;
    };

    this.getCacheStats = function() {
        var request = $q.defer();

        $http.get(ENV.api['root'] + "caches/stats")
            .then(
                function(result) {
                    request.resolve(result.data.response);
                },
                function() {
                    request.reject();
                }
            );

        return request.promise;
    };

    this.getCacheChecks = function() {
        var request = $q.defer();

        $http.get(ENV.api['root'] + "servers/checks")
            .then(
                function(result) {
                    request.resolve(result.data.response);
                },
                function() {
                    request.reject();
                }
            );

        return request.promise;
    };

    this.updateStatus = function(id, payload) {
        var request = $q.defer();

        $http.put(ENV.api['root'] + "servers/" + id + "/status", payload)
            .then(
                function(result) {
                    request.resolve(result);
                },
                function(fault) {
                    request.reject(fault);
                }
            );

        return request.promise;
    };

};

ServerService.$inject = ['$http', '$q', 'Restangular', 'locationUtils', 'messageModel', 'ENV'];
module.exports = ServerService;
