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

var TopologyService = function($http, ENV, locationUtils, messageModel) {

	this.getTopologies = function(queryParams) {
		return [
			{
				"name": "FooTopology",
				"desc": "a topology for Foo DSes",
				"nodes": [
					{
						"id": 0,
						"cachegroup": "aurora",
						"parents": [
							1,
							2
						]
					},
					{
						"id": 1,
						"cachegroup": "denver",
						"parents": [
							5,
							6
						]
					},
					{
						"id": 2,
						"cachegroup": "sac",
						"parents": [
							5,
							6
						]
					},
					{
						"id": 3,
						"cachegroup": "boston",
						"parents": [
							6,
							5
						]
					},
					{
						"id": 4,
						"cachegroup": "albany",
						"parents": [
							5,
							6
						]
					},
					{
						"id": 5,
						"cachegroup": "mid-west",
						"parents": []
					},
					{
						"id": 6,
						"cachegroup": "mid-east",
						"parents": []
					}
				]
			}
		];
		// return $http.get(ENV.api['root'] + 'topologies', {params: queryParams}).then(
		// 	function (result) {
		// 		return result.data.response;
		// 	},
		// 	function (err) {
		// 		throw err;
		// 	}
		// )
	};

	this.getTopology = function(id) {
		return {
			"name": "Topology 1",
			"desc": "a topology for Foo DSes",
			"children": [
				{
					"id": 3,
					"depth": 0,
					"secParent": "",
					"title": "mid-west",
					"name": "mid-west",
					"type": "MID_LOC",
					"nodrop": true,
					"children": [
						{
							"id": 18,
							"depth": 1,
							"secParent": "mid-east",
							"title": "denver",
							"name": "denver",
							"type": "MID_LOC",
							"children": [
								{
									"id": 41,
									"depth": 2,
									"secParent": "sacramento",
									"size": 100,
									"title": "aurora",
									"name": "aurora",
									"type": "EDGE_LOC",
									"children": []
								}
							]
						},
						{
							"id": 1,
							"depth": 1,
							"secParent": "mid-east",
							"title": "sacramento",
							"name": "sacramento",
							"children": []
						}
					]
				},
				{
					"id": 2,
					"depth": 0,
					"secParent": "",
					"title": "mid-east",
					"name": "mid-east",
					"nodrop": true,
					"children": [
						{
							"id": 21,
							"depth": 1,
							"secParent": "mid-west",
							"size": 100,
							"title": "boston",
							"name": "boston",
							"children": []
						},
						{
							"id": 22,
							"depth": 1,
							"secParent": "mid-west",
							"size": 200,
							"title": "albany",
							"name": "albany",
							"children": []
						}
					]
				}
			]
		};
		// return $http.get(ENV.api['root'] + 'topologies', {params: {id: id}}).then(
		// 	function (result) {
		// 		return result.data.response[0];
		// 	},
		// 	function (err) {
		// 		throw err;
		// 	}
		// )
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
		return $http.put(ENV.api['root'] + 'topologies/' + topology.id, topology).then(
			function(result) {
				messageModel.setMessages([ { level: 'success', text: 'Topology updated' } ], false);
				return result;
			},
			function(err) {
				messageModel.setMessages(err.data.alerts, false);
				throw err;
			}
		);
	};

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

TopologyService.$inject = ['$http', 'ENV', 'locationUtils', 'messageModel'];
module.exports = TopologyService;
