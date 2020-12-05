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

var DeliveryServiceService = function($http, locationUtils, messageModel, ENV) {

    this.getDeliveryServices = function(queryParams) {
        return $http.get(ENV.api['deliveryServices'] + 'deliveryservices', {params: queryParams}).then(
            function(result) {
                return result.data.response;
            },
            function(err) {
                throw err;
            }
        );
    };

    this.getDeliveryService = function(id) {
        return $http.get(ENV.api['deliveryServices'] + 'deliveryservices', {params: {id: id}}).then(
            function(result) {
                return result.data.response[0];
            },
            function(err) {
                throw err;
            }
        );
    };

    this.createDeliveryService = function(ds) {
        // strip out any falsy values or duplicates from consistentHashQueryParams
        ds.consistentHashQueryParams = Array.from(new Set(ds.consistentHashQueryParams)).filter(function(i){return i;});

        return $http.post(ENV.api['deliveryServices'] + "deliveryservices", ds).then(
            function(response) {
                return response;
            },
            function(err) {
                throw err;
            }
        );
    };

    this.updateDeliveryService = function(ds) {
        // strip out any falsy values or duplicates from consistentHashQueryParams
        ds.consistentHashQueryParams = Array.from(new Set(ds.consistentHashQueryParams)).filter(function(i){return i;});

        return $http.put(ENV.api['deliveryServices'] + "deliveryservices/" + ds.id, ds).then(
            function(response) {
                return response;
            },
            function(err) {
                throw err;
            }
        );
    };

    // todo: change to use query param when it is supported
    this.deleteDeliveryService = function(ds) {
        return $http.delete(ENV.api['deliveryServices'] + "deliveryservices/" + ds.id).then(
            function(response) {
                return response;
            },
            function(err) {
                throw err;
            }
        );
    };

    this.getServerCapabilities = function(id) {
        return $http.get(ENV.api['deliveryServices'] + 'deliveryservices_required_capabilities', { params: { deliveryServiceID: id } }).then(
            function (result) {
                return result.data.response;
            },
            function (err) {
                throw err;
            }
        )
    };

    this.addServerCapability = function(deliveryServiceId, capabilityName) {
        return $http.post(ENV.api['deliveryServices'] + 'deliveryservices_required_capabilities', { deliveryServiceID: deliveryServiceId, requiredCapability: capabilityName}).then(
            function(result) {
                return result.data;
            },
            function(err) {
                if (err.data && err.data.alerts) {
                    messageModel.setMessages(err.data.alerts, false);
                }
                throw err;
            }
        );
    };

    this.removeServerCapability = function(deliveryServiceId, capabilityName) {
        return $http.delete(ENV.api['deliveryServices'] + 'deliveryservices_required_capabilities', { params: { deliveryServiceID: deliveryServiceId, requiredCapability: capabilityName} }).then(
            function(result) {
                return result.data;
            },
            function(err) {
                if (err.data && err.data.alerts) {
                    messageModel.setMessages(err.data.alerts, false);
                }
                throw err;
            }
        );
    };

    this.getServerDeliveryServices = function(serverId) {
        return $http.get(ENV.api['deliveryServices'] + 'servers/' + serverId + '/deliveryservices').then(
            function(result) {
                return result.data.response;
            },
            function(err) {
                throw err;
            }
        );
    };

    this.getDeliveryServiceTargets = function(dsId) {
        return $http.get(ENV.api['deliveryServices'] + 'steering/' + dsId + '/targets').then(
            function(result) {
                return result.data.response;
            },
            function(err) {
                throw err;
            }
        );
    };

    this.getDeliveryServiceTarget = function(dsId, targetId) {
        return $http.get(ENV.api['deliveryServices'] + 'steering/' + dsId + '/targets', {params: {target: targetId}}).then(
            function(result) {
                return result.data.response[0];
            },
            function(err) {
                throw err;
            }
        );
    };

    this.updateDeliveryServiceTarget = function(dsId, targetId, target) {
        return $http.put(ENV.api['deliveryServices'] + "steering/" + dsId + "/targets/" + targetId, target).then(
            function(result) {
                messageModel.setMessages(result.data.alerts, true);
                locationUtils.navigateToPath('/delivery-services/' + dsId + '/targets');
                return result;
            },
            function(err) {
                messageModel.setMessages(err.data.alerts, false);
                throw err;
            }
        );
    };

    this.createDeliveryServiceTarget = function(dsId, target) {
        return $http.post(ENV.api['deliveryServices'] + 'steering/' + dsId + '/targets', target).then(
            function(result) {
                messageModel.setMessages(result.data.alerts, true);
                locationUtils.navigateToPath('/delivery-services/' + dsId + '/targets');
                return result;
            },
            function(err) {
                messageModel.setMessages(err.data.alerts, false);
                throw err;
            }
        );
    };

    this.deleteDeliveryServiceTarget = function(dsId, targetId) {
        return $http.delete(ENV.api['deliveryServices'] + 'steering/' + dsId + '/targets/' + targetId).then(
            function(result) {
                messageModel.setMessages(result.data.alerts, true);
                locationUtils.navigateToPath('/delivery-services/' + dsId + '/targets');
                return result;
            },
            function(err) {
                messageModel.setMessages(err.data.alerts, true);
                throw err;
            }
        );
    };

    this.getUserDeliveryServices = function(userId) {
        return $http.get(ENV.api['deliveryServices'] + 'users/' + userId + '/deliveryservices').then(
            function(result) {
                return result.data.response;
            },
            function(err) {
                throw err;
            }
        );
    };

    this.deleteDeliveryServiceServer = function(dsId, serverId) {
        return $http.delete(ENV.api['deliveryServices'] + 'deliveryserviceserver/' + dsId + '/' + serverId).then(
            function(result) {
                messageModel.setMessages(result.data.alerts, false);
                return result;
            },
            function(err) {
                messageModel.setMessages(err.data.alerts, true);
                throw err;
            }
        );
    };

    this.assignDeliveryServiceServers = function(dsId, servers) {
        return $http.post(ENV.api['deliveryServices'] + 'deliveryserviceserver',{ dsId: dsId, servers: servers, replace: true } ).then(
            function(result) {
                messageModel.setMessages(result.data.alerts, false);
                return result;
            },
            function(err) {
                messageModel.setMessages(err.data.alerts, false);
                throw err;
            }
        );
    };

    this.getConsistentHashResult = function (regex, requestPath, cdnId) {
        const url = ENV.api['deliveryServices'] + "consistenthash";
        const params = {regex: regex, requestPath: requestPath, cdnId: cdnId};

        return $http.post(url, params).then(
            function (result) {
                return result.data;
            },
            function (err) {
                messageModel.setMessages(err.data.alerts, false);
                throw err;
            }
        );
    };

};

DeliveryServiceService.$inject = ['$http', 'locationUtils', 'messageModel', 'ENV'];
module.exports = DeliveryServiceService;
