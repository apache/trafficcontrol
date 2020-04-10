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

	let hydrateTopology = function() {
		// add some needed fields to each cache group (aka node) of a topology
		topology.nodes.forEach(function(node) {
			let cg = _.findWhere(cacheGroups, { name: node.cachegroup} );
			_.extend(node, { id: cg.id, type: cg.typeName });
		});
	};

	let removeSecParentReferences = function(topologyTree, secParentName) {
		// when a cache group is removed, any references to the cache group as a secParent need to be removed
		topologyTree.forEach(function(node) {
			if (node.secParent && node.secParent === secParentName) {
				node.secParent = '';
			}
			if (node.children && node.children.length > 0) {
				removeSecParentReferences(node.children, secParentName);
			}
		});
	};

	// build a list of cache group names currently in the topology
	let buildCacheGroupNamesInTopology = function(topologyTree, fromScratch) {
		if (fromScratch) cacheGroupNamesInTopology = [];
		topologyTree.forEach(function(node) {
			if (node.cachegroup) {
				cacheGroupNamesInTopology.push(node.cachegroup);
			}
			if (node.children && node.children.length > 0) {
				buildCacheGroupNamesInTopology(node.children, false);
			}
		});
	};

	$scope.topology = topology;

	$scope.topologyTree = [];

	$scope.topologyTreeOptions = {
		beforeDrop: function(evt) {
			let node = evt.source.nodeScope.$modelValue,
				parent = evt.dest.nodesScope.$parent.$modelValue;

			if (!parent || !node) {
				return false; // no dropping outside the toplogy tree and you need a node to drop
			}

			// ORG_LOC cannot have a parent
			if (node.type === 'ORG_LOC' && parent.cachegroup) {
				$anchorScroll(); // scrolls window to top
				messageModel.setMessages([ { level: 'error', text: 'Cache groups of ORG_LOC type must not have a parent.' } ], false);
				return false;
			}

			// EDGE_LOC cannot have children
			if (parent.type === 'EDGE_LOC') {
				$anchorScroll(); // scrolls window to top
				messageModel.setMessages([ { level: 'error', text: 'Cache groups of EDGE_LOC type must not have children.' } ], false);
				return false;
			}

			// update the parent and secParent fields of the node on successful drop
			if (parent.cachegroup) {
				node.parent = parent.cachegroup; // change the node parent based on where the node is dropped
				if (node.parent === node.secParent) {
					// node parent and secParent cannot be the same
					node.secParent = "";
				}
			} else {
				// the node was dropped at the root of the topology. no parents.
				node.parent = "";
				node.secParent = "";
			}
			return true;
		}
	};

	$scope.navigateToPath = locationUtils.navigateToPath;

	$scope.hasError = formUtils.hasError;

	$scope.hasPropertyError = formUtils.hasPropertyError;

	$scope.nodeLabel = function(node) {
		if (!node.cachegroup) return 'TOPOLOGY';
		return node.cachegroup;
	};

	$scope.editSecParent = function(node) {

		if (!node.parent) return; // if a node has no parent, it can't have a second parent

		buildCacheGroupNamesInTopology($scope.topologyTree, true);

		/*  Cache groups that can act as a second parent include:
			1. cache groups of type ORG_LOC or MID_LOC (not EDGE_LOC)
			2. cache groups that are not currently acting as the primary parent
			3. cache groups that exist currently in the topology
		 */
		let eligibleSecParentCandidates = _.filter(cacheGroups, function(cg) {
			return cg.typeName !== 'EDGE_LOC' &&
				(node.parent && node.parent !== cg.name) &&
				cacheGroupNamesInTopology.includes(cg.name);
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
			removeSecParentReferences($scope.topologyTree, node.cachegroup);
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
		return node.type === 'ROOT' || node.type === 'ORG_LOC';
	};

	$scope.isMid = function(node) {
		return node.type === 'MID_LOC';
	};

	$scope.hasChildren = function(node) {
		return node.children.length > 0;
	};

	$scope.addCacheGroups = function(parent, scope) {

		if (parent.type === 'EDGE_LOC') {
			// can't add children to EDGE_LOC. button should be hidden anyhow.
			return;
		}

		// cache groups already in the topology cannot be selected again for addition
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
						parent: (result.parent) ? result.parent : '',
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
