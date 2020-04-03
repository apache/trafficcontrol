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

var FormTopologyController = function(topology, cacheGroups, $scope, $location, $uibModal, formUtils, locationUtils, objectUtils, messageModel) {

	$scope.topology = topology;

	$scope.massaged = [];

	$scope.treeOptions = {
		beforeDrop: function(evt) {

			let node = evt.source.nodeScope.$modelValue,
				parent = evt.dest.nodesScope.$parent.$modelValue;

			if (!parent) {
				messageModel.setMessages([ { level: 'error', text: 'Please keep cache groups inside the topology tree.' } ], false);
				return false; // no dropping outside the toplogy tree
			}

			console.log(parent);

			if (node.type === 'ORG_LOC' && parent.type !== 'ROOT') {
				messageModel.setMessages([ { level: 'error', text: 'Cache groups of ORG_LOC type must be at the top of the topology tree.' } ], false);
				return false;
			}

			if (parent.type === 'EDGE_LOC') {
				messageModel.setMessages([ { level: 'error', text: 'Cache groups of EDGE_LOC type must not have children.' } ], false);
				return false;
			}

			return confirm("Move " + node.cachegroup + " under " + parent.cachegroup + "?");
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

		$scope.massaged = [
			{
				cachegroup: "TOPOLOGY ROOT (ORIGIN LAYER)",
				type: "ROOT",
				children: roots
			}
		];

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
			_.extend(node, { cachegroupId: cg.id, type: cg.typeName });
		})
	};

	$scope.navigateToPath = locationUtils.navigateToPath;

	$scope.hasError = formUtils.hasError;

	$scope.hasPropertyError = formUtils.hasPropertyError;

	$scope.nodeLabel = function(node) {
		if (node.type === 'ROOT') return 'TOPOLOGY ROOT (ORIGIN LAYER)';
		return node.cachegroup + ' [' + node.type + ']'
	};

	$scope.second = function() {
		alert('add 2nd parent');
	};

	$scope.deleteCacheGroup = function(node, scope) {

		if (node.type === 'ROOT') return;

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
			templateUrl: 'common/modules/table/topologyCacheGroups/table.selectTopologyCacheGroups.tpl.html',
			controller: 'TableSelectTopologyCacheGroupsController',
			size: 'lg',
			resolve: {
				node: function() {
					return node;
				},
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
						cachegroupId: cg.id,
						cachegroup: cg.name,
						type: cg.typeName,
						children: []
					}
				});
			cacheGroupNodes.forEach(function(node) {
				nodeData.children.unshift(node);
			});

		}, function () {
			// do nothing
		});
	};

	$scope.viewCacheGroupServers = function(node) {
		let modalInstance = $uibModal.open({
			templateUrl: 'common/modules/table/topologyCacheGroupServers/table.topologyCacheGroupServers.tpl.html',
			controller: 'TableTopologyCacheGroupServersController',
			size: 'lg',
			resolve: {
				cacheGroupName: function() {
					return node.cachegroup;
				},
				cacheGroupServers: function(serverService) {
					return serverService.getServers({ cachegroup: node.cachegroupId });
				}
			}
		});
		modalInstance.result.then(function() {

		}, function () {
			// do nothing
		});

	};

	let init = function() {
		addCacheGroupTypeToTopology();
		massage(angular.copy($scope.topology));
	};
	init();


};

FormTopologyController.$inject = ['topology', 'cacheGroups', '$scope', '$location', '$uibModal', 'formUtils', 'locationUtils', 'objectUtils', 'messageModel'];
module.exports = FormTopologyController;
