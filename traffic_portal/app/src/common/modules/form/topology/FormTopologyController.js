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
				return false; // no dropping outside the toplogy tree
			}

			if (node.type === 'ORG_LOC' && parent.cachegroup !== undefined) {
				messageModel.setMessages([ { level: 'error', text: 'Cache groups of ORG_LOC type must be at the top of the topology tree.' } ], false);
				return false;
			}

			if (parent.type === 'EDGE_LOC') {
				messageModel.setMessages([ { level: 'error', text: 'Cache groups of EDGE_LOC type must not have children.' } ], false);
				return false;
			}

			// change the node parent based on where the node is dragged
			if (parent.cachegroup) {
				node.parent = parent.cachegroup;
				if (node.parent === node.secParent) {
					node.secParent = "";
				}
			} else {
				node.parent = "";
				node.secParent = "";
			}

			return true;
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
				item.parent = "";
				item.secParent = "";
				roots.push(item)
			} else if (item.parents[0] in all) {
				var p = all[item.parents[0]]
				if (!('children' in p)) {
					p.children = []
				}
				p.children.push(item);
				// add parent to each node
				item.parent = all[item.parents[0]].cachegroup;
				// add secParent to each node
				if (item.parents.length === 2 && item.parents[1] in all) {
					item.secParent = all[item.parents[1]].cachegroup;
				}
			}
		});

		$scope.massaged = [
			{
				children: roots
			}
		];

		// console.log(_.flatten(_.map(roots, _.values)) );


		// traverse(roots[0]);
	};

	var scrubSecParent = function(arr, name) {
		arr.forEach(function(node) {
			if (node.secParent && node.secParent === name) {
				node.secParent = '';
			}
			if (node.children && node.children.length > 0) {
				scrubSecParent(node.children, name);
			}
		});
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

	var addCacheGroupPropertiesToTopology = function() {
		topology.nodes.forEach(function(node) {
			var cg = _.findWhere(cacheGroups, { name: node.cachegroup} );
			_.extend(node, { id: cg.id, type: cg.typeName });
		});
	};

	$scope.navigateToPath = locationUtils.navigateToPath;

	$scope.hasError = formUtils.hasError;

	$scope.hasPropertyError = formUtils.hasPropertyError;

	$scope.nodeLabel = function(node) {
		if (node.type === undefined) return 'TOPOLOGY ROOT (ORIGIN LAYER)';
		return node.cachegroup + ' [' + node.type + ']'
	};

	$scope.editSecParent = function(node) {

		if (!node.parent) return;

		console.log(node);

		let flat = objectUtils.flatten(angular.copy($scope.massaged));
		let usedNames = _.values(objectUtils.removeKeysWithout(flat, 'cachegroup'));

		console.log(usedNames);

		let eligibleSecParentCandidates = _.filter(cacheGroups, function(cg) {

			return cg.typeName !== 'EDGE_LOC' && // not an edge_loc cache group
				(node.parent && node.parent !== cg.name) && // not the primary parent cache group
				usedNames.includes(cg.name); // a cache group that exists in the topology
		});

		console.log(eligibleSecParentCandidates);

		let params = {
			title: 'Select a secondary parent',
			message: 'Please select a secondary parent that is part of the ' + topology.name + ' topology',
			key: 'name'
		};
		let modalInstance = $uibModal.open({
			templateUrl: 'common/modules/dialog/select/dialog.select.tpl.html',
			controller: 'DialogSelectController',
			size: 'md',
			resolve: {
				params: function () {
					return params;
				},
				collection: function() {
					return eligibleSecParentCandidates;
				}
			}
		});
		modalInstance.result.then(function(cg) {
			// user selected a secondary parent
			node.secParent = cg.name;
		});
	};

	$scope.deleteCacheGroup = function(node, scope) {
		if (node.cachegroup) {
			scrubSecParent($scope.massaged, node.cachegroup);
			scope.remove();
		}
	};

	$scope.toggle = function(scope) {
		scope.toggle();
	};

	$scope.addCacheGroups = function(parent, scope) {

		if (parent.type === 'EDGE_LOC') {
			// can't add children to EDGE_LOC
			return;
		}

		let flat = objectUtils.flatten(angular.copy($scope.massaged));
		let stripped = objectUtils.removeKeysWithout(flat, 'cachegroup');

		let modalInstance = $uibModal.open({
			templateUrl: 'common/modules/table/topologyCacheGroups/table.selectTopologyCacheGroups.tpl.html',
			controller: 'TableSelectTopologyCacheGroupsController',
			size: 'lg',
			resolve: {
				parent: function() {
					return parent;
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
		modalInstance.result.then(function(result) {
			let nodeData = scope.$modelValue,
				cacheGroupNodes = _.map(result.selectedCacheGroups, function(cg) {
					return {
						id: cg.id,
						cachegroup: cg.name,
						type: cg.typeName,
						parent: result.parent,
						secParent: result.secParent,
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
					return serverService.getServers({ cachegroup: node.id });
				}
			}
		});
		modalInstance.result.then(function() {

		}, function () {
			// do nothing
		});

	};

	let init = function() {
		addCacheGroupPropertiesToTopology();
		massage(angular.copy($scope.topology));
	};
	init();


};

FormTopologyController.$inject = ['topology', 'cacheGroups', '$scope', '$location', '$uibModal', 'formUtils', 'locationUtils', 'objectUtils', 'messageModel'];
module.exports = FormTopologyController;
