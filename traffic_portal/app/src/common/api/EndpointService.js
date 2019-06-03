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

var EndpointService = function($http, ENV, locationUtils, messageModel) {

	this.getEndpoints = function(queryParams) {
		return $http.get(ENV.api['root'] + 'api_capabilities', {params: queryParams}).then(
			function (result) {
				return result.data.response;
			},
			function (err) {
				throw err;
			}
		);
	};

	this.getEndpoint = function(id) {
		return $http.get(ENV.api['root'] + 'api_capabilities/' + id).then(
			function (result) {
				return result.data.response[0];
			},
			function (err) {
				throw err;
			}
		);
	};

	this.createEndpoint = function(endpoint) {
		return $http.post(ENV.api['root'] + 'api_capabilities', endpoint).then(
			function(result) {
				messageModel.setMessages([{level: 'success', text: 'Endpoint created'}], true);
				locationUtils.navigateToPath('/endpoints');
				return result;
			},
			function(err) {
				messageModel.setMessages(err.data.alerts, false);
				throw err;
			}
		);
	};

	this.updateEndpoint = function(endpoint) {
		return $http.put(ENV.api['root'] + 'api_capabilities/' + endpoint.id, endpoint).then(
			function(result) {
				messageModel.setMessages([{level: 'success', text: 'Endpoint updated'}], false);
				return result;
			},
			function(err) {
				messageModel.setMessages(err.data.alerts, false);
				throw err;
			}
		);
	};

	this.deleteEndpoint = function(id) {
		return $http.delete(ENV.api['root'] + 'api_capabilities/' + id).then(
			function(result) {
				messageModel.setMessages([{level: 'success', text: 'Endpoint deleted'}], true);
				return result;
			},
			function(err) {
				messageModel.setMessages(err.data.alerts, true);
				throw err;
			}
		);
	};


};

EndpointService.$inject = ['$http', 'ENV', 'locationUtils', 'messageModel'];
module.exports = EndpointService;
