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

var DeliveryServiceService = function(Restangular, $http, $q, locationUtils, httpService, messageModel, ENV) {

    this.getDeliveryServices = function(queryParams) {
        return Restangular.all('deliveryservices').getList(queryParams);
    };

    this.getDeliveryService = function(id) {
        return Restangular.one("deliveryservices", id).get();
    };

    this.createDeliveryService = function(ds) {
        return Restangular.service('deliveryservices').post(ds)
            .then(
                function(response) {
                    messageModel.setMessages([ { level: 'success', text: 'Delivery Service [ ' + ds.xmlId + ' ] created' } ], true);
                    locationUtils.navigateToPath('/delivery-services/' + response.id + '?type=' + response.type);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
            );
    };

    this.updateDeliveryService = function(ds, delay) {
        var request = $q.defer();

        $http.put(ENV.api['root'] + "deliveryservices/" + ds.id, ds)
            .then(
                function(response) {
                    var response2 = response.data.response[0];
                    messageModel.setMessages([ { level: 'success', text: 'Delivery Service [ ' + ds.xmlId + ' ] updated' } ], delay);
                    locationUtils.navigateToPath('/delivery-services/' + response2.id + '?type=' + response2.type);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
            );

        return request.promise;
    };

    this.deleteDeliveryService = function(ds, delay) {
        return Restangular.one("deliveryservices", ds.id).remove()
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Delivery service [ ' + ds.xmlId + ' ] deleted' } ], delay);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, true);
                }
            );
    };

    this.getServerDeliveryServices = function(serverId) {
        return Restangular.one('servers', serverId).getList('deliveryservices');
    };

    this.getDeliveryServiceTargets = function(dsId) {
        return Restangular.one('steering', dsId).getList('targets');
    };

    this.getDeliveryServiceTarget = function(dsId, targetId) {
        return Restangular.one('steering', dsId).one('targets', targetId).get();
    };

    this.updateDeliveryServiceTarget = function(dsId, targetId, target) {
        var request = $q.defer();

        $http.put(ENV.api['root'] + "steering/" + dsId + "/targets/" + targetId, target)
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Steering target updated' } ], false);
                    locationUtils.navigateToPath('/delivery-services/' + dsId + '/targets');
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
            );

        return request.promise;
    };

    this.createDeliveryServiceTarget = function(dsId, target) {
        return Restangular.one('steering', dsId).all('targets').post(target)
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Steering target created' } ], true);
                    locationUtils.navigateToPath('/delivery-services/' + dsId + '/targets');
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
            );
    };

    this.deleteDeliveryServiceTarget = function(dsId, targetId) {
        return Restangular.one('steering', dsId).one('targets', targetId).remove()
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Steering target deleted' } ], true);
                    locationUtils.navigateToPath('/delivery-services/' + dsId + '/targets');
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, true);
                }
            );
    };

    this.getUserDeliveryServices = function(userId) {
        return Restangular.one('users', userId).getList('deliveryservices');
    };

    this.deleteDeliveryServiceServer = function(dsId, serverId) {
        return httpService.delete(ENV.api['root'] + 'deliveryservice_server/' + dsId + '/' + serverId)
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Delivery service and server were unlinked.' } ], false);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, true);
                }
            );
    };

    this.assignDeliveryServiceServers = function(dsId, servers) {
        return Restangular.service('deliveryserviceserver').post( { dsId: dsId, servers: servers, replace: true } )
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Servers linked to delivery service' } ], false);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
            );
    };


};

DeliveryServiceService.$inject = ['Restangular', '$http', '$q', 'locationUtils', 'httpService', 'messageModel', 'ENV'];
module.exports = DeliveryServiceService;
