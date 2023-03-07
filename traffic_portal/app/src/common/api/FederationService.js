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
 * @param {import("angular").IQService} $q
 * @param {{api: Record<PropertyKey, string>}} ENV
 * @param {import("../service/utils/LocationUtils")} locationUtils
 * @param {import("../models/MessageModel")} messageModel
 */
var FederationService = function($http, $q, ENV, locationUtils, messageModel) {

	const service = this;

	this.getCDNFederations = function(cdnName) {
		return $http.get(ENV.api.unstable + 'cdns/' + cdnName + '/federations').then(
			function (result) {
				return result.data.response;
			},
			function (err) {
				throw err;
			}
		);
	};

	this.getCDNFederation = function(cdnName, fedId) {
		return $http.get(ENV.api.unstable + 'cdns/' + cdnName + '/federations', {params: {id: fedId}}).then(
			function (result) {
				return result.data.response[0];
			},
			function (err) {
				throw err;
			}
		);
	};

	this.createFederation = function(cdn, fed) {
		return $http.post(ENV.api.unstable + 'cdns/' + cdn.name + '/federations', fed).then(
			function(result) {
				const newFedId = result.data.response.id;
				const alerts = result.data.alerts;
				const promises = [];

				// after creating the federation, assign the selected user and ds to the federation
				promises.push(service.assignFederationUsers(newFedId, [ fed.userId ], false));
				promises.push(service.assignFederationDeliveryServices(newFedId, [ fed.dsId ], false));

				$q.all(promises).then(
					function() {
						messageModel.setMessages(alerts, true);
						locationUtils.navigateToPath('/cdns/' + cdn.id + '/federations/' + newFedId);
					},
					function(err) {
						console.error(err);
					}
				);
			},
			function(err) {
				messageModel.setMessages(err.data.alerts, false);
				throw err;
			}
		);
	};

	this.updateFederation = function(cdnName, fed) {
		return $http.put(ENV.api.unstable + 'cdns/' + cdnName + '/federations/' + fed.id, fed).then(
			function() {
				service.assignFederationDeliveryServices(fed.id, [ fed.dsId ], true).then(
					function() {
						messageModel.setMessages([{level: 'success', text: 'Federation updated'}], false);
					},
					function(err) {
						console.error(err);
					}
				);
			},
			function(err) {
				messageModel.setMessages(err.data.alerts, false);
				throw err;
			}
		);
	};

	this.deleteFederation = function(cdnName, fedId) {
		return $http.delete(ENV.api.unstable + 'cdns/' + cdnName + '/federations/' + fedId).then(
			function(result) {
				messageModel.setMessages([{level: 'success', text: 'Federation deleted'}], true);
				return result;
			},
			function(err) {
				messageModel.setMessages(err.data.alerts, true);
				throw err;
			}
		);
	};

	this.getFederationUsers = function(fedId) {
		return $http.get(ENV.api.unstable + 'federations/' + fedId + '/users').then(
			function (result) {
				return result.data.response;
			},
			function (err) {
				throw err;
			}
		);
	};

	this.assignFederationUsers = function(fedId, userIds, replace) {
		return $http.post(ENV.api.unstable + 'federations/' + fedId + '/users', { userIds: userIds, replace: replace }).then(
			function(result) {
				messageModel.setMessages([{level: 'success', text: 'Users linked to federation'}], false);
				return result;
			},
			function(err) {
				messageModel.setMessages(err.data.alerts, false);
				throw err;
			}
		);
	};

	this.deleteFederationUser = function(fedId, userId) {
		return $http.delete(ENV.api.unstable + 'federations/' + fedId + '/users/' + userId).then(
				function(result) {
					messageModel.setMessages([ { level: 'success', text: 'Federation and user were unlinked.' } ], false);
					return result;
				},
				function(err) {
					messageModel.setMessages(err.data.alerts, true);
					throw err;
				}
			);
	};

	this.getFederationDeliveryServices = function(fedId) {
		return $http.get(ENV.api.unstable + 'federations/' + fedId + '/deliveryservices').then(
			function (result) {
				return result.data.response;
			},
			function (err) {
				throw err;
			}
		);
	};

	this.assignFederationDeliveryServices = function(fedId, dsIds, replace) {
		return $http.post(ENV.api.unstable + 'federations/' + fedId + '/deliveryservices', { dsIds: dsIds, replace: replace }).then(
			function(result) {
				messageModel.setMessages([{level: 'success', text: 'Delivery services linked to federation'}], false);
				return result;
			},
			function(err) {
				messageModel.setMessages(err.data.alerts, false);
				throw err;
			}
		);
	};

	this.deleteFederationDeliveryService = function(fedId, dsId) {
		return $http.delete(ENV.api.unstable + 'federations/' + fedId + '/deliveryservices/' + dsId).then(
			function(result) {
				messageModel.setMessages([{level: 'success', text: 'Federation and delivery service were unlinked.'}], false);
				return result;
			},
			function(err) {
				messageModel.setMessages(err.data.alerts, false);
				throw err;
			}
		);
	};

	this.getFederationFederationResolvers = function(fedId) {
		return $http.get(ENV.api.unstable + 'federations/' + fedId + '/federation_resolvers').then(
			function (result) {
				return result.data.response;
			},
			function (err) {
				throw err;
			}
		);
	};

};

FederationService.$inject = ['$http', '$q', 'ENV', 'locationUtils', 'messageModel'];
module.exports = FederationService;
