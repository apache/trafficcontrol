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
	 * Gets a specific Topology from Traffic Ops
	 *
	 * @param name The name of the Topology to be returned.
	 * @returns The Topology with the given name.
	 */
	public async getTopologies(name: string): Promise<ResponseTopology>;
	/**
	 * Gets all Topologies from Traffic Ops
	 *
	 * @returns An Array of all Topologies configured in Traffic Ops.
	 */
	public async getTopologies(): Promise<Array<ResponseTopology>>;
	/**
	 * Gets one or all Topologies from Traffic Ops
	 *
	 * @param name The name of a single Topology to be returned.
	 * @returns Either an Array of Topology objects, or a single Topology, depending on
	 * whether `name` was	passed.
	 */
	public async getTopologies(name?: string): Promise<Array<ResponseTopology> | ResponseTopology> {
		const path = "topologies";
		if (name) {
			const topology = await this.get<[ResponseTopology]>(path, undefined, {name}).toPromise();
			if (topology.length !== 1) {
				throw new Error(`${topology.length} Topologies found by name ${name}`);
			}
			return topology[0];
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
		return this.delete(`topologies?name=${name}`).toPromise();
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
	 */
	public async updateTopology(topology: ResponseTopology): Promise<ResponseTopology> {
		return this.put<ResponseTopology>(`topologies?name=${topology.name}`, topology).toPromise();
	}

	/**
	 * Generates a material tree from a topology.
	 *
	 * @param topology The topology to generate a material tree from.
	 * @returns a material tree.
	 */
	public topologyToTree(topology: ResponseTopology): Array<TopTreeNode> {
		const treeNodes: Array<TopTreeNode> = [];
		const topLevel: Array<TopTreeNode> = [];
		for (const node of topology.nodes) {
			const name = node.cachegroup;
			const cachegroup = node.cachegroup;
			const children: Array<TopTreeNode> = [];
			const parents: Array<TopTreeNode> = [];
			treeNodes.push({
				cachegroup,
				children,
				name,
				parents,
			});
		}
		for (let index = 0; index < topology.nodes.length; index++) {
			const node = topology.nodes[index];
			const treeNode = treeNodes[index];
			if (!(node.parents instanceof Array) || node.parents.length < 1) {
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
	 * @returns a material tree.
	 */
	public treeToTopology(name: string, description: string, treeNodes: Array<TopTreeNode>): ResponseTopology {
		const topologyNodeIndicesByCacheGroup: Map<string, number> = new Map();
		const nodes: Array<ResponseTopologyNode> = new Array<ResponseTopologyNode>();
		this.treeToTopologyInner(topologyNodeIndicesByCacheGroup, nodes, undefined, treeNodes);
		const topology: ResponseTopology = {
			description,
			lastUpdated: new Date(),
			name,
			nodes,
		};
		return topology;
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
	protected treeToTopologyInner(topologyNodeIndicesByCacheGroup: Map<string, number>,
	                              topologyNodes: Array<ResponseTopologyNode>, parent: ResponseTopologyNode | undefined, treeNodes: Array<TopTreeNode>): void {

		for (const treeNode of treeNodes) {
			const cachegroup = treeNode.cachegroup;
			const parents: number[] = [];
			if (parent instanceof Object) {
				const index = topologyNodeIndicesByCacheGroup.get(parent.cachegroup);
				if (!(typeof index === "number")) {
					throw new Error(`index of cachegroup ${parent?.cachegroup} not found in topologyNodeIndicesByCacheGroup`);
				}
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
