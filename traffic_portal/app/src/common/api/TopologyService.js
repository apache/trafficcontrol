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

var TopologyService = function($http, ENV, locationUtils, messageModel, propertiesModel) {

	this.getTopologies = function(queryParams) {
		let top = propertiesModel.topology;
		return [ top ];
	};

	this.createTopology = function(topology) {
		return $http.post(ENV.api['root'] + 'topologies', topology).then(
			function(result) {
				messageModel.setMessages([ { level: 'success', text: 'Topology created' } ], true);
				locationUtils.navigateToPath('/topologies');
				return result;
			},
			function(err) {
				messageModel.setMessages(err.data.alerts, false);
				throw err;
			}
		);
	};

	this.updateTopology = function(topology) {
		console.log(topology);
		propertiesModel.setTopology(topology);
	};

	// this.updateTopology = function(topology) {
	// 	return $http.put(ENV.api['root'] + 'topologies/' + topology.id, topology).then(
	// 		function(result) {
	// 			messageModel.setMessages([ { level: 'success', text: 'Topology updated' } ], false);
	// 			return result;
	// 		},
	// 		function(err) {
	// 			messageModel.setMessages(err.data.alerts, false);
	// 			throw err;
	// 		}
	// 	);
	// };

	this.deleteTopology = function(id) {
		return $http.delete(ENV.api['root'] + "topologies/" + id).then(
			function(result) {
				messageModel.setMessages([ { level: 'success', text: 'Topology deleted' } ], true);
				return result;
			},
			function(err) {
				messageModel.setMessages(err.data.alerts, true);
				throw err;
			}
		);
	};

};

TopologyService.$inject = ['$http', 'ENV', 'locationUtils', 'messageModel', 'propertiesModel'];
module.exports = TopologyService;
