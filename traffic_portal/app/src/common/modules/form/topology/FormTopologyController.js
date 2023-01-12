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

/**
 * @param {*} topology
 * @param {*} cacheGroups
 * @param {import("angular").IAnchorScrollService} $anchorScroll
 * @param {*} $scope
 * @param {import("angular").ILocationService} $location
 * @param {import("../../../service/utils/angular.ui.bootstrap").IModalService} $uibModal
 * @param {import("../../../service/utils/FormUtils")} formUtils
 * @param {import("../../../service/utils/LocationUtils")} locationUtils
 * @param {import("../../../service/utils/TopologyUtils")} topologyUtils
 * @param {import("../../../models/MessageModel")} messageModel
 */
var FormTopologyController = function(topology, cacheGroups, $anchorScroll, $scope, $location, $uibModal, formUtils, locationUtils, topologyUtils, messageModel) {

	let cacheGroupNamesInTopology = [];

	let hydrateTopology = function() {
		// add some needed fields to each cache group (aka node) of a topology
		topology.nodes.forEach(function(node) {
			let cacheGroup = cacheGroups.find( function(cg) { return cg.name === node.cachegroup} );
			Object.assign(node, { id: cacheGroup.id, type: cacheGroup.typeName });
		});
	};

	let removeSecParentReferences = function(topologyTree, secParentName) {
		// when a cache group is removed, any references to the cache group as a secParent need to be removed
		topologyTree.forEach(function(node) {
			if (node.secParent && node.secParent.name === secParentName) {
				node.secParent = { name: '', type: '' };
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

			// EDGE_LOC can only have EDGE_LOC children
			if (parent.type === 'EDGE_LOC' && node.type !== 'EDGE_LOC') {
				$anchorScroll(); // scrolls window to top
				messageModel.setMessages([ { level: 'error', text: 'EDGE_LOC cache groups can only have EDGE_LOC children.' } ], false);
				return false;
			}

			// update the parent and secParent fields of the node on successful drop
			if (parent.cachegroup) {
				// change the node parent based on where the node is dropped
				node.parent = { name: parent.cachegroup, type: parent.type };
				if (node.parent.name === node.secParent.name) {
					// node parent and secParent cannot be the same
					node.secParent = { name: '', type: '' };
				}
			} else {
				// the node was dropped at the root of the topology. no parents.
				node.parent = { name: '', type: '' };
				node.secParent = { name: '', type: '' };
			}
			// marks the form as dirty thus enabling the save btn
			$scope.topologyForm.dirty.$setDirty();
			return true;
		}
	};

	$scope.navigateToPath = (path, unsavedChanges) => locationUtils.navigateToPath(path, unsavedChanges);

	$scope.hasError = formUtils.hasError;

	$scope.hasPropertyError = formUtils.hasPropertyError;

	$scope.clone = function(topology) {
		locationUtils.navigateToPath('/topologies/clone?name=' + topology.name);
	};

	$scope.viewCacheGroups = function() {
		$location.path('/topologies/cache-groups');
	};

	$scope.viewDeliveryServices = function() {
		$location.path('/topologies/delivery-services');
	};

	$scope.viewServers = function() {
		$location.path('/topologies/servers');
	};

	$scope.nodeLabel = function(node) {
		return node.cachegroup || 'TOPOLOGY';
	};

	$scope.editSecParent = function(node) {

		if (!node.parent) return; // if a node has no parent, it can't have a second parent

		buildCacheGroupNamesInTopology($scope.topologyTree, true);

		/*  Cache groups that can act as a second parent include:
			1. cache groups that are not the current cache group (you can't parent/sec parent yourself)
			2. cache groups that are not currently acting as the primary parent (primary parent != sec parent)
			3. cache groups that exist currently in the topology only
			4a. any cache group types (ORG_LOC, MID_LOC, EDGE_LOC) if child cache group is EDGE_LOC
			4b. only MID_LOC or ORG_LOC cache group types if child cache group is not EDGE_LOC
		 */
		let eligibleSecParentCandidates = cacheGroups.filter(function(cg) {
			return (node.cachegroup && node.cachegroup !== cg.name) &&
				(node.parent && node.parent.name !== cg.name) &&
				cacheGroupNamesInTopology.includes(cg.name) &&
				((node.type === 'EDGE_LOC') || (cg.typeName === 'MID_LOC' || cg.typeName === 'ORG_LOC'));
		}).sort(function(a,b) { return [a.name, b.name].sort().indexOf(b.name) === 0 ? 1 : -1; });

		let params = {
			title: 'Select a secondary parent',
			message: 'Please select a secondary parent that is part of the ' + topology.name + ' topology',
			key: 'name',
			required: false,
			selectedItemKeyValue: node.secParent.name,
			labelFunction: function(item) { return item['name'] + ' (' + item['typeName'] + ')' }
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
				node.secParent = { name: selectedSecParent.name, type: selectedSecParent.typeName };
			} else {
				node.secParent = { name: '', type: '' };
			}
			// marks the form as dirty thus enabling the save btn
			$scope.topologyForm.dirty.$setDirty();
		});
	};

	$scope.deleteCacheGroup = function(node, scope) {
		if (node.cachegroup) {
			removeSecParentReferences($scope.topologyTree, node.cachegroup);
			scope.remove();
			// marks the form as dirty thus enabling the save btn
			$scope.topologyForm.dirty.$setDirty();
		}
	};

	$scope.toggle = function(scope) {
		scope.toggle();
	};

	$scope.nodeWarning = function(node) {
		// EDGE_LOCs with parent/secondary parent EDGE_LOCs require special configuration
		let msg = 'Special Configuration Required';
		if (node.parent && node.parent.type === 'EDGE_LOC') {
			return msg + ' [EDGE_LOC Parent]';
		} else if (node.secParent && node.secParent.type === 'EDGE_LOC') {
			return msg + ' [EDGE_LOC 2nd Parent]';
		}
		return '';
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

		// cache groups already in the topology cannot be selected again for addition
		buildCacheGroupNamesInTopology($scope.topologyTree, true);

		// the types of child cache groups you can add depends on the parent cache group's type
		let eligibleChildrenTypes = [ { name: 'EDGE_LOC' } ];
		if (parent.type === 'ROOT') {
			eligibleChildrenTypes.push({ name: 'MID_LOC' });
			eligibleChildrenTypes.push({ name: 'ORG_LOC' });
		} else if (parent.type === 'MID_LOC' || parent.type === 'ORG_LOC') {
			eligibleChildrenTypes.push({ name: 'MID_LOC' });
		}

		let parentName = (parent.cachegroup) ? parent.cachegroup : 'ROOT';
		let params = {
			title: 'Select a child cache group type',
			message: 'Please select the type of child cache group(s) you would like to add to ' + parentName,
			key: 'name',
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
					return eligibleChildrenTypes;
				}
			}
		});
		modalInstance.result.then(function(type) {
			modalInstance = $uibModal.open({
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
					selectedType: function() {
						return type.name;
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
					cacheGroupNodes = result.selectedCacheGroups.map(function(cg) {
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
				// marks the form as dirty thus enabling the save btn
				$scope.topologyForm.dirty.$setDirty();
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
