/*


 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

 http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.

 */

var DeliveryServiceService = function($http, $q, deliveryServicesModel, messageModel, propertiesModel, ENV) {

    var capacityRequest,
        routingMethodsRequest,
        cacheGroupHealthRequest;

    var displayTimoutError = function(options) {
        var msg = (angular.isDefined(options.message)) ? options.message : 'Request timeout';
        if (options.status.toString().match(/^5\d[24]$/)) {
            // 502 or 504
            messageModel.setMessages([ { level: 'error', text: msg } ], false);
        }
    };

    this.getDeliveryServices = function() {
        var promise = $http.get(ENV.apiEndpoint['1.2'] + "deliveryservices.json")
            .success(function(result) {
                deliveryServicesModel.setDeliveryServices(result.response);
                return result.response;
            })
            .error(function(fault) {
            });

        return promise;
    };

    this.getDeliveryService = function(deliveryServiceId) {
        var deferred = $q.defer();
        $http.get(ENV.apiEndpoint['1.2'] + "deliveryservices/" + deliveryServiceId + ".json")
            .success(function(result) {
                deferred.resolve(result.response[0]);
            })
            .error(function(fault) {
                deferred.resolve(null);
            });

        return deferred.promise;
    };

    this.getState = function(deliveryServiceId, ignoreLoadingBar) {
        var deferred = $q.defer();
        $http.get(ENV.apiEndpoint['1.2'] + "deliveryservices/" + deliveryServiceId + "/state.json", { ignoreLoadingBar: ignoreLoadingBar })
            .success(function(result) {
                deferred.resolve(result.response);
            })
            .error(function(fault, status) {
                displayTimoutError({ status: status });
                deferred.reject();
            });

        return deferred.promise;
    };

    this.getCapacity = function(deliveryServiceId, ignoreLoadingBar, showTimeoutError) {
        if (capacityRequest) {
            capacityRequest.reject();
        }
        capacityRequest = $q.defer();

        $http.get(ENV.apiEndpoint['1.2'] + "deliveryservices/" + deliveryServiceId + "/capacity.json",
            { timeout: capacityRequest.promise, ignoreLoadingBar: ignoreLoadingBar })
            .success(function(result) {
                capacityRequest.resolve(result.response);
            })
            .error(function(fault, status) {
                if (showTimeoutError) displayTimoutError({ status: status });
                capacityRequest.reject();
            });

        return capacityRequest.promise;
    };


    this.getRoutingMethods = function(deliveryServiceId, ignoreLoadingBar, showTimeoutError) {
        if (routingMethodsRequest) {
            routingMethodsRequest.reject();
        }
        routingMethodsRequest = $q.defer();

        $http.get(ENV.apiEndpoint['1.2'] + "deliveryservices/" + deliveryServiceId + "/routing.json",
            { timeout: routingMethodsRequest.promise, ignoreLoadingBar: ignoreLoadingBar })
            .success(function(result) {
                routingMethodsRequest.resolve(result.response);
            })
            .error(function(fault, status) {
                if (showTimeoutError) displayTimoutError({ status: status });
                routingMethodsRequest.reject();
            });

        return routingMethodsRequest.promise;
    };

    this.getCacheGroupHealth = function(deliveryServiceId, ignoreLoadingBar, showTimeoutError) {
        if (cacheGroupHealthRequest) {
            cacheGroupHealthRequest.reject();
        }
        cacheGroupHealthRequest = $q.defer();

        $http.get(ENV.apiEndpoint['1.2'] + "deliveryservices/" + deliveryServiceId + "/health.json",
            { timeout: cacheGroupHealthRequest.promise, ignoreLoadingBar: ignoreLoadingBar })
            .success(function(result) {
                cacheGroupHealthRequest.resolve(result.response);
            })
            .error(function(fault, status) {
                if (showTimeoutError) displayTimoutError({ status: status });
                cacheGroupHealthRequest.reject();
            });

        return cacheGroupHealthRequest.promise;
    };

    this.getPurgeJobs = function(deliveryServiceId, ignoreLoadingBar) {
        var deferred = $q.defer();
        $http.get(ENV.apiEndpoint['1.2'] + "user/current/jobs.json?dsId=" + deliveryServiceId + "&keyword=PURGE", { ignoreLoadingBar: ignoreLoadingBar })
            .success(function(result) {
                deferred.resolve(result.response);
            })
            .error(function(fault) {
                deferred.reject();
            });

        return deferred.promise;
    };

    this.getServers = function(deliveryServiceId, ignoreLoadingBar) {
        var deferred = $q.defer();
        $http.get(ENV.apiEndpoint['1.2'] + "servers.json?orderby=type&dsId=" + deliveryServiceId, { ignoreLoadingBar: ignoreLoadingBar })
            .success(function(result) {
                deferred.resolve(result.response);
            })
            .error(function(fault) {
                deferred.reject();
            });

        return deferred.promise;
    };

    this.createDSRequest = function(dsData) {
        var deferred = $q.defer();
        $http.post(ENV.apiEndpoint['1.2'] + "deliveryservices/request", { emailTo: propertiesModel.properties.deliveryService.request.email, details: dsData } )
            .success(function(result) {
                if (angular.isDefined(result.alerts)) {
                    messageModel.setMessages(result.alerts, false);
                }
                deferred.resolve(result);
            })
            .error(function(fault) {
                if (angular.isDefined(fault.alerts)) {
                    messageModel.setMessages(fault.alerts, false);
                }
                deferred.reject();
            });

        return deferred.promise;
    };

};

DeliveryServiceService.$inject = ['$http', '$q', 'deliveryServicesModel', 'messageModel', 'propertiesModel', 'ENV'];
module.exports = DeliveryServiceService;