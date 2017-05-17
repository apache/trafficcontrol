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

var DeliveryServiceService = function(Restangular, locationUtils, httpService, messageModel, ENV) {

    this.getDeliveryServices = function(queryParams) {
        return Restangular.all('deliveryservices').getList(queryParams);
    };

    this.getDeliveryService = function(id) {
        return Restangular.one("deliveryservices", id).get();
    };

    this.createDeliveryService = function(deliveryService) {
        return Restangular.service('deliveryservices').post(deliveryService)
            .then(
                function(response) {
                    messageModel.setMessages([ { level: 'success', text: 'DeliveryService created' } ], true);
                    locationUtils.navigateToPath('/configure/delivery-services/' + response.id + '?type=' + response.type);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
            );
    };

    this.updateDeliveryService = function(deliveryService) {
        return deliveryService.put()
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Delivery service updated' } ], false);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
            );
    };

    this.deleteDeliveryService = function(id) {
        return Restangular.one("deliveryservices", id).remove()
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Delivery service deleted' } ], true);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, true);
                }
            );
    };

    this.getServerDeliveryServices = function(serverId) {
        return Restangular.one('servers', serverId).getList('deliveryservices');
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

};

DeliveryServiceService.$inject = ['Restangular', 'locationUtils', 'httpService', 'messageModel', 'ENV'];
module.exports = DeliveryServiceService;
