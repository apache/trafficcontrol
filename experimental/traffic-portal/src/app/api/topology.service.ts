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
import { HttpClient } from "@angular/common/http";
import { Injectable } from "@angular/core";
import type {
	RequestTopology,
	ResponseTopology,
	ResponseTopologyNode,
} from "trafficops-types";

import { APIService } from "./base-api.service";

/**
 * TopTreeNode is used to represent a topology in a format usable as a material
 * nested tree data source.
 */
export interface TopTreeNode {
	name: string;
	cachegroup: string;
	children: Array<TopTreeNode>;
	parents: Array<this>;
}

/**
 * TopologyService exposes API functionality relating to Topologies.
 */
@Injectable()
export class TopologyService extends APIService {

	constructor(http: HttpClient) {
		super(http);
	}

	/**
	 * Gets one or all Topologies from Traffic Ops
	 *
	 * @param name The name of a single Topology to be returned.
	 * @returns An Array of Topologies
	 * whether `name` was passed.
	 */
	public async getTopologies(name?: string): Promise<Array<ResponseTopology>> {
		const path = "topologies";
		if (name) {
			const topology = await this.get<[ResponseTopology]>(path, undefined, {name}).toPromise();
			if (topology.length !== 1 || topology[0].name !== name) {
				throw new Error(`${topology.length} Topologies found by name ${name}`);
			}
			return topology;
		}
		return this.get<Array<ResponseTopology>>(path).toPromise();
	}

	/**
	 * Deletes a Topology.
	 *
	 * @param topology The Topology to be deleted, or just its name.
	 */
	public async deleteTopology(topology: ResponseTopology | string): Promise<void> {
		const name = typeof topology === "string" ? topology : topology.name;
		return this.delete("topologies", undefined, {name}).toPromise();
	}

	/**
	 * Creates a new Topology.
	 *
	 * @param topology The Topology to create.
	 */
	public async createTopology(topology: RequestTopology): Promise<ResponseTopology> {
		return this.post<ResponseTopology>("topologies", topology).toPromise();
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
		return this.put<ResponseTopology>("topologies", topology, {name}).toPromise();
	}

	/**
	 * Generates a material tree from a topology.
	 *
	 * @param topology The topology to generate a material tree from.
	 * @returns a material tree.
	 */
	public static topologyToTree(topology: ResponseTopology): Array<TopTreeNode> {
		const treeNodes: Array<TopTreeNode> = [];
		const topLevel: Array<TopTreeNode> = [];
		for (const node of topology.nodes) {
			treeNodes.push({
				cachegroup: node.cachegroup,
				children: [],
				name: node.cachegroup,
				parents: [],
			});
		}
		for (let index = 0; index < topology.nodes.length; index++) {
			const node = topology.nodes[index];
			const treeNode = treeNodes[index];
			if (!Array.isArray(node.parents) || node.parents.length < 1) {
				topLevel.push(treeNode);
				continue;
			}
			for (const parent of node.parents) {
				treeNodes[parent].children.push(treeNode);
				treeNode.parents.push(treeNodes[parent]);
			}
		}
		return topLevel;
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
		const topologyNodeIndicesByCacheGroup: Map<string, number> = new Map();
		const nodes: Array<ResponseTopologyNode> = new Array<ResponseTopologyNode>();
		this.treeToTopologyInner(topologyNodeIndicesByCacheGroup, nodes, undefined, treeNodes);
		return {
			description,
			lastUpdated: new Date(),
			name,
			nodes,
		};
	}

	/**
	 * Inner recursive function for generating a Topology from a material tree.
	 *
	 * @param topologyNodeIndicesByCacheGroup A map of Topology node indices
	 * using cache group names as the key
	 * @param topologyNodes The mutable array of Topology nodes
	 * @param parent The parent, if it exists
	 * @param treeNodes The data for a material tree
	 */
	protected static treeToTopologyInner(
		topologyNodeIndicesByCacheGroup: Map<string, number>,
		topologyNodes: Array<ResponseTopologyNode>,
		parent: ResponseTopologyNode | undefined,
		treeNodes: Array<TopTreeNode>): void {

		for (const treeNode of treeNodes) {
			const cachegroup = treeNode.cachegroup;
			const parents: number[] = [];
			if (typeof parent !== "undefined") {
				const index = topologyNodeIndicesByCacheGroup.get(parent.cachegroup) as number;
				parents.push(index);
			}
			const topologyNode: ResponseTopologyNode = {
				cachegroup,
				parents,
			};
			topologyNodes.push(topologyNode);
			topologyNodeIndicesByCacheGroup.set(cachegroup, topologyNodes.length - 1);
			if (treeNode.children.length > 0) {
				this.treeToTopologyInner(topologyNodeIndicesByCacheGroup, topologyNodes, topologyNode, treeNode.children);
			}
		}
	}
}
