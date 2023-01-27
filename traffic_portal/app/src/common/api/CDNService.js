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

/**
 * @param {import("angular").IHttpService} $http
 * @param {import("../service/utils/LocationUtils")} locationUtils
 * @param {import("../models/MessageModel")} messageModel
 * @param {{api: Record<PropertyKey, string>}} ENV
 */
var CDNService = function($http, locationUtils, messageModel, ENV) {

    this.getCDNs = function(all) {
        return $http.get(ENV.api.unstable + 'cdns').then(
            function(result) {
                let response;
                if (all) { // there is a CDN called "ALL" that is not really a CDN but you might want it...
                    response = result.data.response;
                } else {
                    response = result.data.response.filter(function(cdn) {
                        return cdn.name != 'ALL';
                    });
                }
                return response;
            },
            function(err) {
                throw err;
            }
        );
    };


    this.getCDN = function(id) {
        return $http.get(ENV.api.unstable + 'cdns', {params: {id: id}}).then(
            function(result) {
                return result.data.response[0];
            },
            function(err) {
                throw err;
            }
        );
    };

    this.createCDN = function(cdn) {
        return $http.post(ENV.api.unstable + 'cdns', cdn).then(
            function(result) {
                messageModel.setMessages(result.data.alerts, true);
                locationUtils.navigateToPath('/cdns');
                return result;
            },
            function(err) {
                messageModel.setMessages(err.data.alerts, false);
                throw err;
            }
        );
    };

    // todo: change to use query param when it is supported
    this.updateCDN = function(cdn) {
        return $http.put(ENV.api.unstable + 'cdns/' + cdn.id, cdn).then(
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

    // todo: change to use query param when it is supported
    this.deleteCDN = function(id) {
        return $http.delete(ENV.api.unstable + 'cdns/' + id).then(
            function(result) {
                return result.data;
            },
            function(err) {
                messageModel.setMessages(err.data.alerts, false);
                throw err;
            }
        );
    };

    this.queueServerUpdates = function(id) {
        return $http.post(ENV.api.unstable + 'cdns/' + id + '/queue_update', {action: "queue"}).then(
            function(result) {
                messageModel.setMessages([{level: 'success', text: 'Queued CDN server updates'}], false);
                return result;
            },
            function(err) {
                messageModel.setMessages(err.data.alerts, false);
                throw err;
            }
        );
    };

    this.clearServerUpdates = function(id) {
        return $http.post(ENV.api.unstable + 'cdns/' + id + '/queue_update', {action: "dequeue"}).then(
            function(result) {
                messageModel.setMessages([{ level: 'success', text: 'Cleared CDN server updates'}], false);
                return result;
            },
            function(err) {
                messageModel.setMessages(err.data.alerts, false);
                throw err;
            }
        );
    };

    this.getCapacity = function() {
        return $http.get(ENV.api.unstable + 'cdns/capacity').then(
            function(result) {
                return result.data.response;
            },
            function(err) {
                throw err;
            }
        );
    };

    this.getRoutingMethods = function() {
        return $http.get(ENV.api.unstable + 'cdns/routing').then(
            function(result) {
                return result.data.response;
            },
            function(err) {
                throw err;
            }
        );
    };

    this.getCurrentStats = function() {
        return $http.get(ENV.api.unstable + 'current_stats').then(
            function(result) {
                if (result) {
                    return result.data.response;
                }
                console.warn("Failed to fetch stats: ", result);
            },
            function(err) {
                throw err;
            }
        );
    };

    this.getCurrentSnapshot = function(cdnName) {
       return $http.get(ENV.api.unstable + 'cdns/' + cdnName + '/snapshot').then(
            function(result) {
                return result.data.response;
            },
            function(err) {
                throw err;
            }
        );
    };

    this.getNewSnapshot = function(cdnName) {
        return $http.get(ENV.api.unstable + 'cdns/' + cdnName + '/snapshot/new').then(
            function(result) {
                return result.data.response;
            },
            function(err) {
                throw err;
            }
        );
    };

    this.snapshot = function(cdn) {
        return $http.put(ENV.api.unstable + 'snapshot', undefined, {params: {cdnID: cdn.id}}).then(
            function(result) {
                messageModel.setMessages([{level: 'success', text: 'Snapshot performed'}], true);
                locationUtils.navigateToPath('/cdns/' + cdn.id);
                return result;
            },
            function(err) {
                messageModel.setMessages(err.data.alerts, false);
                throw err;
            }
        );
    };

    this.getDNSSECKeys = function(cdnName) {
        return $http.get(ENV.api.unstable + 'cdns/name/' + cdnName + '/dnsseckeys').then(
            function(result) {
                return result.data.response;
            },
            function(err) {
                throw err;
            }
        );
    };

    this.generateDNSSECKeys = function(dnssecKeysRequest) {
        return $http.post(ENV.api.unstable + 'cdns/dnsseckeys/generate', dnssecKeysRequest).then(
            function(result) {
                return result;
            },
            function(err) {
                messageModel.setMessages(err.data.alerts, false);
                throw err;
            }
        );
    };

	this.regenerateKSK = function(kskRequest, cdnKey) {
		return $http.post(ENV.api.unstable + 'cdns/' + cdnKey + '/dnsseckeys/ksk/generate', kskRequest).then(
			function(result) {
				return result;
			},
			function(err) {
				messageModel.setMessages(err.data.alerts, false);
                throw err;
			}
		);
	}

    this.getNotifications = function(queryParams) {
        return $http.get(ENV.api.unstable + 'cdn_notifications', { params: queryParams }).then(
            function(result) {
                return result.data.response;
            },
            function(err) {
                throw err;
            }
        );
    };

    this.createNotification = function(cdn, notification) {
        return $http.post(ENV.api.unstable + 'cdn_notifications', { cdn: cdn.name, notification: notification}).then(
            function(result) {
                return result;
            },
            function(err) {
                messageModel.setMessages(err.data.alerts, false);
                throw err;
            }
        );
    };

    this.deleteNotification = function(queryParams) {
        return $http.delete(ENV.api.unstable + 'cdn_notifications', { params: queryParams }).then(
            function(result) {
                return result;
            },
            function(err) {
                messageModel.setMessages(err.data.alerts, false);
                throw err;
            }
        );
    };

    this.getLocks = function(queryParams) {
        return $http.get(ENV.api.unstable + 'cdn_locks', { params: queryParams }).then(
            function(result) {
                return result.data.response;
            },
            function(err) {
                throw err;
            }
        );
    };

    this.createLock = function(lock) {
        return $http.post(ENV.api.unstable + 'cdn_locks', lock).then(
            function(result) {
                return result;
            },
            function(err) {
                messageModel.setMessages(err.data.alerts, false);
                throw err;
            }
        );
    };

    this.deleteLock = function(queryParams) {
        return $http.delete(ENV.api.unstable + 'cdn_locks', { params: queryParams }).then(
            function(result) {
                return result;
            },
            function(err) {
                messageModel.setMessages(err.data.alerts, false);
                throw err;
            }
        );
    };

};

CDNService.$inject = ['$http', 'locationUtils', 'messageModel', 'ENV'];
module.exports = CDNService;
