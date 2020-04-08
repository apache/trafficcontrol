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

var FormTopologyController = function(topology, cacheGroups, $anchorScroll, $scope, $location, $uibModal, formUtils, locationUtils, topologyUtils, messageModel) {

	let cacheGroupNamesInTopology = [];

	$scope.topology = topology;

	$scope.topologyTree = [];

	$scope.topologyTreeOptions = {
		beforeDrop: function(evt) {

			let node = evt.source.nodeScope.$modelValue,
				parent = evt.dest.nodesScope.$parent.$modelValue;

			if (!parent) {
				return false; // no dropping outside the toplogy tree
			}

			if (node.type === 'ORG_LOC' && parent.cachegroup !== undefined) {
				$anchorScroll(); // scrolls window to top
				messageModel.setMessages([ { level: 'error', text: 'Cache groups of ORG_LOC type must be at the top of the topology tree.' } ], false);
				return false;
			}

			if (parent.type === 'EDGE_LOC') {
				$anchorScroll(); // scrolls window to top
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

	let hydrateTopology = function() {
		topology.nodes.forEach(function(node) {
			var cg = _.findWhere(cacheGroups, { name: node.cachegroup} );
			_.extend(node, { id: cg.id, type: cg.typeName });
		});
	};

	let scrubCacheGroupByName = function(arr, name) {
		arr.forEach(function(node) {
			if (node.secParent && node.secParent === name) {
				node.secParent = '';
			}
			if (node.children && node.children.length > 0) {
				scrubCacheGroupByName(node.children, name);
			}
		});
	};

	let buildCacheGroupNamesInTopology = function(topologyTree, fromScratch) {
		if (fromScratch) cacheGroupNamesInTopology = [];
		topologyTree.forEach(function(node) {
			if (node.cachegroup) {
				cacheGroupNamesInTopology.push(node.cachegroup)
			}
			if (node.children && node.children.length > 0) {
				buildCacheGroupNamesInTopology(node.children, false);
			}
		});
	};

	$scope.navigateToPath = locationUtils.navigateToPath;

	$scope.hasError = formUtils.hasError;

	$scope.hasPropertyError = formUtils.hasPropertyError;

	$scope.nodeLabel = function(node) {
		if (node.cachegroup === undefined) return 'TOPOLOGY ROOT';
		return node.cachegroup;
	};

	$scope.editSecParent = function(node) {

		if (!node.parent) return; // if a node has no parent, it can't have a second parent

		buildCacheGroupNamesInTopology($scope.topologyTree, true);

		let eligibleSecParentCandidates = _.filter(cacheGroups, function(cg) {
			return cg.typeName !== 'EDGE_LOC' && // not an edge_loc cache group
				(node.parent && node.parent !== cg.name) && // not the primary parent cache group
				cacheGroupNamesInTopology.includes(cg.name); // a cache group that exists in the topology
		});

		let params = {
			title: 'Select a secondary parent',
			message: 'Please select a secondary parent that is part of the ' + topology.name + ' topology',
			key: 'name',
			required: false,
			selectedItemKeyValue: node.secParent
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
		modalInstance.result.then(function(selectedSecParent) {
			if (selectedSecParent) {
				node.secParent = selectedSecParent.name;
			} else {
				node.secParent = '';
			}
		});
	};

	$scope.deleteCacheGroup = function(node, scope) {
		if (node.cachegroup) {
			scrubCacheGroupByName($scope.topologyTree, node.cachegroup);
			scope.remove();
		}
	};

	$scope.toggle = function(scope) {
		scope.toggle();
	};

	$scope.hasNodeError = function(node) {
		if (node.type !== 'EDGE_LOC' && node.children.length === 0) {
			return true;
		}
		return false;
	};

	$scope.isOrigin = function(node) {
		return node.type === undefined || node.type === 'ORG_LOC';
	};

	$scope.hasChildren = function(node) {
		return node.children.length > 0;
	};

	$scope.addCacheGroups = function(parent, scope) {

		if (parent.type === 'EDGE_LOC') {
			// can't add children to EDGE_LOC
			return;
		}

		buildCacheGroupNamesInTopology($scope.topologyTree, true);

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
					return cacheGroupNamesInTopology;
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
		});
	};

	$scope.viewCacheGroupServers = function(node) {
		$uibModal.open({
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
	};

	let init = function() {
		hydrateTopology();
		$scope.topologyTree = topologyUtils.getTopologyTree($scope.topology);
	};
	init();
};

FormTopologyController.$inject = ['topology', 'cacheGroups', '$anchorScroll', '$scope', '$location', '$uibModal', 'formUtils', 'locationUtils', 'topologyUtils', 'messageModel'];
module.exports = FormTopologyController;
