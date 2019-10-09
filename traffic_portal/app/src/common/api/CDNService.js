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

var CDNService = function($http, $q, Restangular, locationUtils, messageModel, ENV) {

    this.getCDNs = function(all) {
        var request = $q.defer();

        $http.get(ENV.api['root'] + "cdns")
            .then(
                function(result) {
                    var response;
                    if (all) { // there is a CDN called "ALL" that is not really a CDN but you might want it...
                        response = result.data.response;
                    } else {
                        response = _.filter(result.data.response, function(cdn) {
                            return cdn.name != 'ALL';
                        });
                    }
                    request.resolve(response);
                },
                function(fault) {
                    request.reject();
                }
            );

        return request.promise;
    };


    this.getCDN = function(id) {
        return Restangular.one("cdns", id).get();
    };

    this.createCDN = function(cdn) {
        return Restangular.service('cdns').post(cdn)
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'CDN created' } ], true);
                    locationUtils.navigateToPath('/cdns');
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
            );
    };

    this.updateCDN = function(cdn) {
        return cdn.put()
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'CDN updated' } ], false);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
            );
    };

    this.deleteCDN = function(id) {
        var request = $q.defer();

        $http.delete(ENV.api['root'] + "cdns/" + id)
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

    this.queueServerUpdates = function(id) {
        return Restangular.one("cdns", id).customPOST( { action: "queue"}, "queue_update" )
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Queued CDN server updates' } ], false);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
            );
    };

    this.clearServerUpdates = function(id) {
        return Restangular.one("cdns", id).customPOST( { action: "dequeue"}, "queue_update" )
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Cleared CDN server updates' } ], false);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
            );
    };

    this.getCapacity = function() {
        var request = $q.defer();

        $http.get(ENV.api['root'] + "cdns/capacity")
            .then(
                function(result) {
                    request.resolve(result.data.response);
                },
                function(fault) {
                    request.reject();
                }
            );

        return request.promise;
    };

    this.getRoutingMethods = function() {
        var request = $q.defer();

        $http.get(ENV.api['root'] + "cdns/routing")
            .then(
                function(result) {
                    request.resolve(result.data.response);
                },
                function(fault) {
                    request.reject();
                }
            );

        return request.promise;
    };

    this.getCurrentStats = function() {
        var request = $q.defer();

        $http.get(ENV.api['root'] + "current_stats")
            .then(
                function(result) {
                    request.resolve(result.data.response);
                },
                function(fault) {
                    request.reject();
                }
            );

        return request.promise;
    };

    this.getCurrentSnapshot = function(cdnName) {
        var request = $q.defer();

        $http.get(ENV.api['root'] + "cdns/" + cdnName + "/snapshot")
            .then(
                function(result) {
                    request.resolve(result.data.response);
                },
                function(fault) {
                    request.reject();
                }
            );

        return request.promise;
    };

    this.getNewSnapshot = function(cdnName) {
        var request = $q.defer();

        $http.get(ENV.api['root'] + "cdns/" + cdnName + "/snapshot/new")
            .then(
                function(result) {
                    request.resolve(result.data.response);
                },
                function(fault) {
                    request.reject();
                }
            );

        return request.promise;
    };

    this.snapshot = function(cdn) {
        var request = $q.defer();

        $http.put(ENV.api['root'] + "cdns/" + cdn.id + "/snapshot")
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Snapshot performed' } ], true);
                    locationUtils.navigateToPath('/cdns/' + cdn.id);

                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
            );

        return request.promise;
    };

    this.getDNSSECKeys = function(cdnName) {
        var request = $q.defer();

        $http.get(ENV.api['root'] + "cdns/name/" + cdnName + "/dnsseckeys")
            .then(
                function(result) {
                    request.resolve(result.data.response);
                },
                function(fault) {
                    request.reject();
                }
            );

        return request.promise;
    };

    this.generateDNSSECKeys = function(dnssecKeysRequest) {
        var request = $q.defer();

        $http.post(ENV.api['root'] + "cdns/dnsseckeys/generate", dnssecKeysRequest)
            .then(
                function(result) {
                    request.resolve(result);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                    request.reject();
                }
            );

        return request.promise;
    };

	this.regenerateKSK = function(kskRequest, cdnKey) {
		var request = $q.defer();

		$http.post(ENV.api['root'] + "cdns/" + cdnKey + "/dnsseckeys/ksk/generate", kskRequest)
			.then(
				function(result) {
					request.resolve(result);
				},
				function(fault) {
					messageModel.setMessages(fault.data.alerts, false);
					request.reject();
				}
			);

		return request.promise;
	};

};

CDNService.$inject = ['$http', '$q', 'Restangular', 'locationUtils', 'messageModel', 'ENV'];
module.exports = CDNService;
