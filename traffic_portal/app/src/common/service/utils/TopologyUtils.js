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
 * @typedef TopologyNode
 * @property {string} cachegroup
 * @property {{name: string}} parent
 * @property {number[]} parents
 * @property {{name: string}} secParent
 */

/**
 * @typedef Topology
 * @property {string} description
 * @property {string} name
 * @property {TopologyNode[]} nodes
 */

/**
 * @typedef TopologyTree
 * @property {string|null|undefined} cachegroup
 * @property {TopologyTree[]|null|undefined} children
 * @property {{name: string}} parent
 * @property {{name: string}} secParent
 */

var TopologyUtils = function() {

	/** @type Topology|undefined */
	let normalizedTopology;

	/**
	 * Flattens a TopologyTree into a normalized Topology.
	 *
	 * @param {TopologyTree[]} topologyTree
	 * @param {boolean} fromScratch If true, the existing normalized Topology
	 * has its nodes wiped and rebuilt (this is used in recursion, callers
	 * should always pass `true`).
	 */
	function flattenTopology(topologyTree, fromScratch) {
		if (fromScratch) {
			normalizedTopology.nodes = [];
		}

		for (const node of topologyTree) {
			if (node.cachegroup) {
				normalizedTopology.nodes.push({
					cachegroup: node.cachegroup,
					parent: node.parent,
					secParent: node.secParent,
					parents: []
				});
			}
			if (node.children && node.children.length > 0) {
				flattenTopology(node.children, false);
			}
		}
	};

	function addNodeIndexes() {
		for (const currentNode of normalizedTopology.nodes) {
			const parentNodeIndex = normalizedTopology.nodes.findIndex(node => currentNode.parent.name === node.cachegroup);
			const secParentNodeIndex = normalizedTopology.nodes.findIndex(node => currentNode.secParent.name === node.cachegroup);
			if (parentNodeIndex > -1) {
				currentNode.parents.push(parentNodeIndex);
				if (secParentNodeIndex > -1) {
					currentNode.parents.push(secParentNodeIndex);
				}
			}
		}
	};

	this.getNormalizedTopology = function(name, description, topologyTree) {
		// build a normalized (flat) topology with parent indexes required for topology create/update
		normalizedTopology = {
			name: name,
			description: description,
			nodes: []
		};
		flattenTopology(topologyTree, true);
		addNodeIndexes();
		return normalizedTopology;
	};

	this.getTopologyTree = function(topology) {
		let nodes = angular.copy(topology.nodes);
		let roots = [], // topology items without parents (primary or secondary)
			all = {};

		nodes.forEach(function(node, index) {
			all[index] = node;
		});

		// create children based on parent definitions
		Object.keys(all).forEach(function (guid) {
			let item = all[guid];
			if (!('children' in item)) {
				item.children = [];
			}
			if (item.parents.length === 0) {
				item.parent = { name: '', type: '' };
				item.secParent = { name: '', type: '' };
				roots.push(item);
			} else if (item.parents[0] in all) {
				let p = all[item.parents[0]];
				if (!('children' in p)) {
					p.children = [];
				}
				p.children.push(item);
				// add parent to each node
				item.parent = { name: all[item.parents[0]].cachegroup, type: all[item.parents[0]].type };
				// add secParent to each node
				if (item.parents.length === 2 && item.parents[1] in all) {
					item.secParent = { name: all[item.parents[1]].cachegroup, type: all[item.parents[1]].type };
				} else {
					item.secParent = { name: '', type: '' };
				}
			}
		});

		return [
			{
				type: 'ROOT',
				children: roots
			}
		];
	};

};

TopologyUtils.$inject = [];
module.exports = TopologyUtils;
