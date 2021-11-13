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

var FederationResolverService = function($http, ENV, locationUtils, messageModel) {

	this.getFederationResolvers = function(queryParams) {
		return $http.get(ENV.api.unstable + 'federation_resolvers', {params: queryParams}).then(
			function (result) {
				return result.data.response;
			},
			function (err) {
				throw err;
			}
		);
	};

	this.createFederationResolver = function(fedResolver) {
		return $http.post(ENV.api.unstable + 'federation_resolvers', fedResolver).then(
			function(result) {
				return result;
			},
			function(err) {
				throw err;
			}
		);
	};

	this.assignFederationResolvers = function(fedId, fedResIds, replace) {
		return $http.post(ENV.api.unstable + 'federations/' + fedId + '/federation_resolvers', { fedResolverIds: fedResIds, replace: replace }).then(
			function(result) {
				messageModel.setMessages([ { level: 'success', text: fedResIds.length + ' resolver(s) assigned to federation' } ], false);
				return result;
			},
			function(err) {
				messageModel.setMessages(err.data.alerts, false);
				throw err;
			}
		);
	};


};

FederationResolverService.$inject = ['$http', 'ENV', 'locationUtils', 'messageModel'];
module.exports = FederationResolverService;
