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

var FederationResolverService = function(Restangular, $http, $q, ENV, locationUtils, messageModel) {

	this.getFederationResolvers = function(queryParams) {
		return Restangular.all('federation_resolvers').getList(queryParams);
	};

	this.createFederationResolver = function(fedResolver) {
		var deferred = $q.defer();

		$http.post(ENV.api['root'] + 'federation_resolvers', fedResolver)
			.then(
				function(result) {
					deferred.resolve(result);
				},
				function(fault) {
					deferred.reject(fault);
				}
			);

		return deferred.promise;
	};

	this.assignFederationResolvers = function(fedId, fedResIds, replace) {
		return $http.post(ENV.api['root'] + 'federations/' + fedId + '/federation_resolvers', { fedResolverIds: fedResIds, replace: replace })
			.then(
				function() {
					messageModel.setMessages([ { level: 'success', text: fedResIds.length + ' resolver(s) assigned to federation' } ], false);
				},
				function(fault) {
					messageModel.setMessages(fault.data.alerts, false);
				}
			);
	};


};

FederationResolverService.$inject = ['Restangular', '$http', '$q', 'ENV', 'locationUtils', 'messageModel'];
module.exports = FederationResolverService;
