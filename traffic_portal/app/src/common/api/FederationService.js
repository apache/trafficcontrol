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

var FederationService = function(Restangular, $http, $q, ENV, locationUtils, messageModel, httpService) {

	var service = this;

	this.getCDNFederations = function(cdnName) {
		return Restangular.one('cdns', cdnName).getList('federations');
	};

	this.getCDNFederation = function(cdnName, fedId) {
		return Restangular.one('cdns', cdnName).one('federations', fedId).get();
	};

	this.createFederation = function(cdn, fed) {
		return $http.post(ENV.api['root'] + 'cdns/' + cdn.name + '/federations', fed)
			.then(
				function(result) {
					var newFedId = result.data.response.id,
						alerts = result.data.alerts,
						promises = [];
					// after creating the federation, assign the selected user and ds to the federation
					promises.push(service.assignFederationUsers(newFedId, [ fed.userId ], false));
					promises.push(service.assignFederationDeliveryServices(newFedId, [ fed.dsId ], false));

					$q.all(promises)
						.then(
							function() {
								messageModel.setMessages(alerts, true);
								locationUtils.navigateToPath('/cdns/' + cdn.id + '/federations/' + newFedId);
							},
							function(fault) {
								// do nothing
							});

				},
				function(fault) {
					messageModel.setMessages(fault.data.alerts, false);
				}
			);
	};

	this.updateFederation = function(fed) {
		return fed.put()
			.then(
				function() {
					service.assignFederationDeliveryServices(fed.id, [ fed.dsId ], true)
						.then(
								function() {
									messageModel.setMessages([ { level: 'success', text: 'Federation updated' } ], false);
								}
							);
				},
				function(fault) {
					messageModel.setMessages(fault.data.alerts, false);
				}
			);
	};

	this.deleteFederation = function(cdnId, fedId) {
		return Restangular.one('cdns', cdnId).one('federations', fedId).remove()
			.then(
				function() {
					messageModel.setMessages([ { level: 'success', text: 'Federation deleted' } ], true);
				},
				function(fault) {
					messageModel.setMessages(fault.data.alerts, true);
				}
			);
	};

	this.getFederationUsers = function(fedId) {
		return Restangular.one('federations', fedId).getList('users');
	};

	this.assignFederationUsers = function(fedId, userIds, replace) {
		return $http.post(ENV.api['root'] + 'federations/' + fedId + '/users', { userIds: userIds, replace: replace })
			.then(
				function() {
					messageModel.setMessages([ { level: 'success', text: 'Users linked to federation' } ], false);
				},
				function(fault) {
					messageModel.setMessages(fault.data.alerts, false);
				}
			);
	};

	this.deleteFederationUser = function(fedId, userId) {
		return httpService.delete(ENV.api['root'] + 'federations/' + fedId + '/users/' + userId)
			.then(
				function() {
					messageModel.setMessages([ { level: 'success', text: 'Federation and user were unlinked.' } ], false);
				},
				function(fault) {
					messageModel.setMessages(fault.data.alerts, true);
				}
			);
	};

	this.getFederationDeliveryServices = function(fedId) {
		return Restangular.one('federations', fedId).getList('deliveryservices');
	};

	this.assignFederationDeliveryServices = function(fedId, dsIds, replace) {
		return $http.post(ENV.api['root'] + 'federations/' + fedId + '/deliveryservices', { dsIds: dsIds, replace: replace })
			.then(
				function() {
					messageModel.setMessages([ { level: 'success', text: 'Delivery services linked to federation' } ], false);
				},
				function(fault) {
					messageModel.setMessages(fault.data.alerts, false);
				}
			);
	};

	this.deleteFederationDeliveryService = function(fedId, dsId) {
		return httpService.delete(ENV.api['root'] + 'federations/' + fedId + '/deliveryservices/' + dsId)
			.then(
				function() {
					messageModel.setMessages([ { level: 'success', text: 'Federation and delivery service were unlinked.' } ], false);
				},
				function(fault) {
					messageModel.setMessages(fault.data.alerts, false);
				}
			);
	};

	this.getFederationFederationResolvers = function(fedId) {
		return Restangular.one('federations', fedId).getList('federation_resolvers');
	};

};

FederationService.$inject = ['Restangular', '$http', '$q', 'ENV', 'locationUtils', 'messageModel', 'httpService'];
module.exports = FederationService;
