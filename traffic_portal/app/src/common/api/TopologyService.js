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

var TopologyService = function($http, ENV, messageModel) {

	this.getTopologies = function(queryParams) {
		return $http.get(ENV.api.unstable + 'topologies', { params: queryParams }).then(
			function(result) {
				return result.data.response;
			},
			function(err) {
				throw err;
			}
		);
	};

	this.createTopology = function(topology) {
		return $http.post(ENV.api.unstable + 'topologies', topology).then(
			function(result) {
				return result;
			},
			function(err) {
				messageModel.setMessages(err.data.alerts, false);
				throw err;
			}
		);
	};

	this.updateTopology = function(topology, currentName) {
		return $http.put(ENV.api.unstable + 'topologies', topology, { params: { name: currentName } }).then(
			function(result) {
				return result;
			},
			function(err) {
				messageModel.setMessages(err.data.alerts, false);
				throw err;
			}
		);
	};

	this.deleteTopology = function(topology) {
		return $http.delete(ENV.api.unstable + "topologies", { params: { name: topology.name } }).then(
			function(result) {
				return result;
			},
			function(err) {
				messageModel.setMessages(err.data.alerts, false);
				throw err;
			}
		);
	};

	this.queueServerUpdates = function(topology, cdnId) {
		return $http.post(ENV.api.unstable + 'topologies/' + topology + '/queue_update', {action: "queue", cdnId: cdnId}).then(
			function(result) {
				messageModel.setMessages([{level: 'success', text: 'Queued topology server updates'}], false);
				return result;
			},
			function(err) {
				messageModel.setMessages(err.data.alerts, false);
				throw err;
			}
		);
	};

	this.clearServerUpdates = function(topology, cdnId) {
		return $http.post(ENV.api.unstable + 'topologies/' + topology + '/queue_update', {action: "dequeue", cdnId: cdnId}).then(
			function(result) {
				messageModel.setMessages([{level: 'success', text: 'Cleared topology server updates'}], false);
				return result;
			},
			function(err) {
				messageModel.setMessages(err.data.alerts, false);
				throw err;
			}
		);
	};


};

TopologyService.$inject = ['$http', 'ENV', 'messageModel'];
module.exports = TopologyService;
