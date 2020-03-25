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

var FormTopologyController = function(topology, $scope, $location, formUtils, locationUtils, serverService) {

	$scope.topology = topology;

	$scope.massaged = [];

	$scope.treeOptions = {
		beforeDrop: function(evt) {
			console.log('drop');
			console.log(evt);
			let node = evt.source.nodeScope.$modelValue.cachegroup,
				parent = (evt.dest.nodesScope.$parent.$modelValue) ? evt.dest.nodesScope.$parent.$modelValue.cachegroup : 'root'
			return confirm("Move " + node + " under " + parent + "?");
		}
	};

	var massage = function(topology) {
		var roots = []; // things without parent

		// make them accessible by guid on this map
		var all = {};

		topology.nodes.forEach(function(node, index) {
			all[index] = node;
		});

		// connect childrens to its parent, and split roots apart
		Object.keys(all).forEach(function (guid) {
			var item = all[guid];
			if (!('children' in item)) {
				item.children = []
			}
			if (item.parents.length === 0) {
				roots.push(item)
			} else if (item.parents[0] in all) {
				var p = all[item.parents[0]]
				if (!('children' in p)) {
					p.children = []
				}
				p.children.push(item);
				// add secParent to each node
				if (item.parents.length === 2 && item.parents[1] in all) {
					item.secParent = all[item.parents[1]].cachegroup;
				}
			}
		});

		console.log(roots);

		$scope.massaged = roots;
	};

	$scope.target = [
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
	];

	$scope.navigateToPath = locationUtils.navigateToPath;

	$scope.hasError = formUtils.hasError;

	$scope.hasPropertyError = formUtils.hasPropertyError;

	$scope.second = function() {
		alert('add 2nd parent');
	};

	$scope.remove = function(scope) {
		console.log(scope);
		// if (confirm("Remove " + evt.$modelValue.cachegroup + " and all its children?")){
		// 	scope.remove();
		// }
	};

	$scope.toggle = function(scope) {
		scope.toggle();
	};

	$scope.newSubItem = function(scope) {
		var nodeData = scope.$modelValue;
		nodeData.children.push({
			id: nodeData.id * 10 + nodeData.children.length,
			secParent: "",
			cachegroup: nodeData.cachegroup + '.' + (nodeData.children.length + 1),
			children: []
		});
	};

	$scope.viewServers = function(scope) {
		var nodeData = scope.$modelValue;
		alert('open dialog with cachegroup servers');
		// serverService.getServers({ cachegroup: nodeData.id })
		// 	.then(function(result) {
		// 		debugger;
		// 		$scope.cacheGroupServers = result;
		// 	});

	};

	let init = function() {
		massage(angular.copy($scope.topology));
	};
	init();


};

FormTopologyController.$inject = ['topology', '$scope', '$location', 'formUtils', 'locationUtils', 'serverService'];
module.exports = FormTopologyController;
