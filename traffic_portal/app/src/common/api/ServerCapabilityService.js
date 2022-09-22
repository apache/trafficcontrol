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

var ServerCapabilityService = function($http, ENV, locationUtils, messageModel) {

	this.getServerCapabilities = function(queryParams) {
		return $http.get(ENV.api.unstable + 'server_capabilities', {params: queryParams}).then(
			function(result) {
				return result.data.response;
			},
			function(err) {
				throw err;
			}
		);
	};

	this.getServerCapability = function(name) {
		return $http.get(ENV.api.unstable + 'server_capabilities', {params: {"name": name}}).then(
			function(result) {
				return result.data.response[0];
			},
			function(err) {
				throw err;
			}
		);
	};

	this.createServerCapability = function(serverCap) {
		return $http.post(ENV.api.unstable + 'server_capabilities', serverCap).then(
			function(result) {
				messageModel.setMessages(result.data.alerts, true);
				locationUtils.navigateToPath('/server-capabilities');
				return result;
			},
			function(err) {
				messageModel.setMessages(err.data.alerts, false);
				throw err;
			}
		);
	};

	this.deleteServerCapability = function(name) {
		return $http.delete(ENV.api.unstable + 'server_capabilities', {params: {"name": name}}).then(
			function(result) {
				messageModel.setMessages(result.data.alerts, true);
				return result;
			},
			function(err) {
				messageModel.setMessages(err.data.alerts, false);
				throw err;
			}
		);
	};

	this.assignSCsServer = function(serverId, serverCapabilities) {
		return $http.put(ENV.api.unstable + 'multiple_server_capabilities',{ serverId: serverId, serverCapabilities: serverCapabilities, replace: true } ).then(
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

	this.assignServersPerSC = function(serverCapability, serverIds) {
		return $http.put(ENV.api.unstable + 'multiple_servers_per_capability',{ serverCapability: serverCapability.name, serverIds: serverIds, replace: true } ).then(
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

	this.updateServerCapability = function(currentName, serverCapability) {
		return $http.put(ENV.api.unstable + 'server_capabilities', serverCapability, {params: {"name": currentName}}).then(
			function(result) {
				return result;
			},
			function(err) {
				messageModel.setMessages(err.data.alerts, false);
				throw err;
			}
		);
	};

	this.getServerCapabilityServers = function(capabilityName) {
		return $http.get(ENV.api.unstable + 'server_server_capabilities', { params: { serverCapability: capabilityName } }).then(
			function (result) {
				return result.data.response;
			},
			function (err) {
				throw err;
			}
		);
	};

	this.getServerCapabilityDeliveryServices = function(capabilityName) {
		return $http.get(ENV.api.unstable + 'deliveryservices_required_capabilities', { params: { requiredCapability: capabilityName } }).then(
			function (result) {
				return result.data.response;
			},
			function (err) {
				throw err;
			}
		);
	};

};

ServerCapabilityService.$inject = ['$http', 'ENV', 'locationUtils', 'messageModel'];
module.exports = ServerCapabilityService;
