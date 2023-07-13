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
		}
	];

	/**
	 * Gets one or all Topologies from Traffic Ops
	 *
	 * @param name? The integral, unique identifier of a single Topology to be
	 * returned
	 * @returns Either a Map of Topology names to full Topology objects, or a single Topology, depending on whether `id` was
	 * 	passed.
	 * (In the event that `id` is passed but does not match any Topology, `null` will be emitted)
	 */
	public async getTopologies(name?: string): Promise<Array<ResponseTopology> | ResponseTopology> {
		if (name !== undefined) {
			const topology = this.topologies.find(t => t.name === name);
			if (!topology) {
				throw new Error(`no such Topology #${name}`);
			}
			return topology;
		}
		return this.topologies;
	}

	/**
	 * Deletes a Topology.
	 *
	 * @param topology The Topology to be deleted, or just its ID.
	 */
	public async deleteTopology(topology: ResponseTopology | string): Promise<void> {
		const name = typeof topology === "string" ? topology : topology.name;
		const idx = this.topologies.findIndex(t => t.name === name);
		if (idx < 0) {
			throw new Error(`no such Topology: #${name}`);
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
			if (!(node.parents instanceof Array)) {
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
	 * @param topology The full new definition of the Topology being
	 * updated, or just its ID.
	 */
	public async updateTopology(topology: ResponseTopology): Promise<ResponseTopology> {
		const idx = this.topologies.findIndex(t => t.name === topology.name);
		topology = {
			...topology,
			lastUpdated: new Date()
		};

		if (idx < 0) {
			throw new Error(`no such Topology: #${topology}`);
		}

		this.topologies[idx] = topology;
		return topology;
	}
}
