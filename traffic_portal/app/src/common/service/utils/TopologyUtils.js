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
 * @typedef NormalizedTopologyNode
 * @property {string} cachegroup
 * @property {{name: string; type: string}} parent
 * @property {{name: string; type: string}} secParent
 * @property {number[]} parents
 */

/**
 * @typedef NormalizedTopology
 * @property {string} description
 * @property {string} name
 * @property {NormalizedTopologyNode[]} nodes
 */

/**
 * @typedef TopologyTree
 * @property {string} cachegroup
 * @property {string} [type]
 * @property {{name: string; type: string}} parent
 * @property {{name: string; type: string}} secParent
 * @property {TopologyTree[]} [children]
 */

/**
 * @typedef TopologyNode
 * @property {string} cachegroup
 * @property {{name: string; type: string}} parent
 * @property {{name: string; type: string}} secParent
 * @property {number[]} parents
 * @property {string} type
 * @property {TopologyNode[]} [children]
 */

/**
 * TopologyUtils provides utilities for dealing with Topology hierarchies.
 */
class TopologyUtils {

	/**
	 * @private
	 * @param {TopologyTree[]} topologyTree
	 * @param {NormalizedTopology} topology
	 * @param {boolean} [fromScratch]
	 */
	flattenTopology(topologyTree, topology, fromScratch) {
		if (fromScratch)
			topology.nodes = [];
		for (const node of topologyTree) {
			if (node.cachegroup) {
				topology.nodes.push({
					cachegroup: node.cachegroup,
					parent: node.parent,
					secParent: node.secParent,
					parents: []
				});
			}
			if (node.children && node.children.length > 0) {
				this.flattenTopology(node.children, topology, false);
			}
		}
	}

	/**
	 * @private
	 * @param {NormalizedTopology} topology
	 */
	addNodeIndexes(topology) {
		for (const currentNode of topology.nodes) {
			const parentNodeIndex = topology.nodes.findIndex(node => currentNode.parent.name === node.cachegroup);
			const secParentNodeIndex = topology.nodes.findIndex(node => currentNode.secParent.name === node.cachegroup);
			if (parentNodeIndex > -1) {
				currentNode.parents.push(parentNodeIndex);
				if (secParentNodeIndex > -1) {
					currentNode.parents.push(secParentNodeIndex);
				}
			}
		}
	};

	/**
	 * "Normalizes" the given Topology tree.
	 *
	 * @param {string} name
	 * @param {string} description
	 * @param {TopologyTree[]} topologyTree
	 * @returns {NormalizedTopology}
	 */
	getNormalizedTopology(name, description, topologyTree) {
		// build a normalized (flat) topology with parent indexes required for topology create/update
		const normalizedTopology = {
			name,
			description,
			nodes: []
		};
		this.flattenTopology(topologyTree, normalizedTopology, true);
		this.addNodeIndexes(normalizedTopology);
		return normalizedTopology;
	};

	/**
	 * Converts a set of Topology nodes into a Topology tree.
	 *
	 * @param {{nodes: TopologyNode[]}} topology
	 * @returns {[{type: "ROOT"; children: TopologyTree[]}]}
	 */
	getTopologyTree(topology) {
		/** @type {TopologyNode[]} */
		const nodes = angular.copy(topology.nodes);
		/** @type {TopologyTree[]} */
		const roots = []; // topology items without parents (primary or secondary)
		const all = Object.fromEntries(nodes.map((n, i) => [i, n]))

		// create children based on parent definitions
		for (const item of Object.values(all)) {
			if (!("children" in item)) {
				item.children = [];
			}
			if (item.parents.length === 0) {
				item.parent = { name: "", type: "" };
				item.secParent = { name: "", type: "" };
				roots.push(item);
			} else if (item.parents[0] in all) {
				const p = all[item.parents[0]];
				if (!p.children) {
					p.children = [];
				}
				p.children.push(item);
				// add parent to each node
				item.parent = { name: all[item.parents[0]].cachegroup, type: all[item.parents[0]].type };
				// add secParent to each node
				if (item.parents.length === 2 && item.parents[1] in all) {
					item.secParent = { name: all[item.parents[1]].cachegroup, type: all[item.parents[1]].type };
				} else {
					item.secParent = { name: "", type: "" };
				}
			}
		}

		return [
			{
				type: "ROOT",
				children: roots
			}
		];
	};
}

TopologyUtils.$inject = [];
module.exports = TopologyUtils;
