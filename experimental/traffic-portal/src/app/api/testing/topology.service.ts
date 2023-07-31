/*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*/

import { Injectable } from "@angular/core";
import {
	RequestTopology,
	ResponseTopology,
	ResponseTopologyNode
} from "trafficops-types";

import { TopologyService as ConcreteService, TopTreeNode } from "src/app/api";

/**
 * TopologyService expose API functionality relating to Topologies.
 */
@Injectable()
export class TopologyService {
	private readonly topologies: ResponseTopology[] = [
		{
			description: "",
			lastUpdated: new Date(),
			name: "test",
			nodes: [
				{
					cachegroup: "Edge",
					parents: [1],
				},
				{
					cachegroup: "Mid",
					parents: [2],
				},
				{
					cachegroup: "Origin",
					parents: [],
				},
			],
		},
	];

	/**
	 * Gets one or all Topologies from Traffic Ops
	 *
	 * @param name The unique name of a single Topology to be returned
	 * @returns An Array of Topologies
	 */
	public async getTopologies(name?: string): Promise<Array<ResponseTopology>> {
		if (name !== undefined) {
			const topology = this.topologies.find(t => t.name === name);
			if (!topology) {
				throw new Error(`no such Topology ${name}`);
			}
			return [topology];
		}
		return this.topologies;
	}

	/**
	 * Deletes a Topology.
	 *
	 * @param topology The Topology to be deleted, or just its name.
	 */
	public async deleteTopology(topology: ResponseTopology | string): Promise<void> {
		const name = typeof topology === "string" ? topology : topology.name;
		const idx = this.topologies.findIndex(t => t.name === name);
		if (idx < 0) {
			throw new Error(`no such Topology: ${name}`);
		}
		this.topologies.splice(idx, 1);
	}

	/**
	 * Creates a new Topology.
	 *
	 * @param topology The Topology to create.
	 */
	public async createTopology(topology: RequestTopology): Promise<ResponseTopology> {
		const nodes: ResponseTopologyNode[] = topology.nodes.map(node => {
			if (!Array.isArray(node.parents)) {
				node.parents = [];
			}
			const responseNode: ResponseTopologyNode = {
				cachegroup: node.cachegroup,
				parents: node.parents || [],
			};
			return responseNode;
		});
		const t: ResponseTopology = {
			description: topology.description || "",
			lastUpdated: new Date(),
			name: topology.name,
			nodes,
		};
		this.topologies.push(t);
		return t;
	}

	/**
	 * Replaces an existing Topology with the provided new definition of a
	 * Topology.
	 *
	 * @param topology The full new definition of the Topology being updated
	 * @param name What the topology was named before it was updated
	 */
	public async updateTopology(topology: ResponseTopology, name?: string): Promise<ResponseTopology> {
		if (typeof name === "undefined") {
			name = topology.name;
		}
		const idx = this.topologies.findIndex(t => t.name === name);
		topology = {
			...topology,
			lastUpdated: new Date()
		};

		if (idx < 0) {
			throw new Error(`no such Topology: ${topology}`);
		}

		this.topologies[idx] = topology;
		return topology;
	}

	/**
	 * Generates a material tree from a topology.
	 *
	 * @param topology The topology to generate a material tree from.
	 * @returns a material tree.
	 */
	public static topologyToTree(topology: ResponseTopology): Array<TopTreeNode> {
		return ConcreteService.topologyToTree(topology);
	}

	/**
	 * Generates a topology from a material tree.
	 *
	 * @param name The topology name
	 * @param description The topology description
	 * @param treeNodes The data for a material tree
	 * @returns a topology.
	 */
	public static treeToTopology(name: string, description: string, treeNodes: Array<TopTreeNode>): ResponseTopology {
		return ConcreteService.treeToTopology(name, description, treeNodes);
	}
}
