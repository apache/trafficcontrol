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

var FormTopologyController = function(topology, cacheGroups, $scope, $location, $uibModal, formUtils, locationUtils, objectUtils) {

	$scope.topology = topology;

	$scope.massaged = [];

	$scope.treeOptions = {
		beforeDrop: function(evt) {
			console.log(evt);

			let node = evt.source.nodeScope.$modelValue,
				parent = evt.dest.nodesScope.$parent.$modelValue;

			let parentName = (parent) ? parent.cachegroup : 'root';

			console.log(node);
			console.log(parent);


			if (parent && parent.type === 'EDGE_LOC') {
				alert('sorry');
				return false;
			}

			return confirm("Move " + node.cachegroup + " under " + parentName + "?");
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

		$scope.massaged = roots;

		// console.log(_.flatten(_.map(roots, _.values)) );


		// traverse(roots[0]);
	};

	const traverse = function(obj) {
		_.each(obj, function (val, key, obj) {
			if (_.isArray(val)) {
				val.forEach(function(el) {
					traverse(el);
				});
			} else if (_.isObject(val)) {
				traverse(val);
			} else {
				console.log('i am a leaf');
				console.log(val);
			}
		});
	};

	var addCacheGroupTypeToTopology = function() {
		topology.nodes.forEach(function(node) {
			var cg = _.findWhere(cacheGroups, { name: node.cachegroup} );
			_.extend(node, { type: cg.typeName });
		})
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

	$scope.deleteCacheGroup = function(scope) {
		console.log(scope);
		let cg = scope.$nodeScope.$modelValue.cachegroup;
		if (confirm("Remove " + cg + " and all its children?")){
			scope.remove();
		}
	};

	$scope.toggle = function(scope) {
		scope.toggle();
	};

	$scope.addCacheGroups = function(node, scope) {

		if (node.type === 'EDGE_LOC') {
			// todo: better
			alert('no');
			return;
		}

		let flat = objectUtils.flatten(angular.copy($scope.massaged));
			stripped = objectUtils.removeKeysWithout(flat, 'cachegroup');

		let modalInstance = $uibModal.open({
			templateUrl: 'common/modules/table/cacheGroups/table.selectCacheGroups.tpl.html',
			controller: 'TableSelectCacheGroupsController',
			size: 'lg',
			resolve: {
				topology: function() {
					return topology;
				},
				cacheGroups: function(cacheGroupService) {
					return cacheGroupService.getCacheGroups();
				},
				usedCacheGroupNames: function() {
					return _.values(stripped);
				}
			}
		});
		modalInstance.result.then(function(selectedCacheGroups) {
			let nodeData = scope.$modelValue,
				cacheGroupNodes = _.map(selectedCacheGroups, function(cg) {
					return {
						secParent: "",
						cachegroup: cg.name,
						type: cg.typeName,
						children: []
					}
				});
			cacheGroupNodes.forEach(function(node) {
				nodeData.children.push(node);
			});

		}, function () {
			// do nothing
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
		addCacheGroupTypeToTopology();
		massage(angular.copy($scope.topology));
	};
	init();


};

FormTopologyController.$inject = ['topology', 'cacheGroups', '$scope', '$location', '$uibModal', 'formUtils', 'locationUtils', 'objectUtils'];
module.exports = FormTopologyController;
