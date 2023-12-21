/**
 * @license Apache-2.0
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

import { TestBed } from "@angular/core/testing";
import type { ResponseTopology } from "trafficops-types";

import { TopologyService as ConcreteTopologyService } from "..";

import { TopologyService } from "./topology.service";

import { APITestingModule } from ".";

describe("TestingTopologyService", () => {
	let service: TopologyService;
	let topology: ResponseTopology;

	beforeEach(async () => {
		TestBed.configureTestingModule({
			imports: [APITestingModule],
			providers: [
				TopologyService,
			]
		});
		service = TestBed.inject(TopologyService);
		const topologies = await service.getTopologies();
		expect(topologies.length).toBeGreaterThanOrEqual(1);
		topology = topologies[0];
	});

	it("should be created", () => {
		expect(service).toBeTruthy();
	});

	it("gets a single Topology by name", async () => {
		await expectAsync(service.getTopologies(topology.name)).toBeResolvedTo(topology);
	});

	it("throws an error when asked to get a Topology that doesn't exist", async () => {
		await expectAsync(service.getTopologies("")).toBeRejected();
	});

	it("deletes topologies", async () => {
		const initialLength = (await service.getTopologies()).length;
		await service.deleteTopology(topology);

		const response = await service.getTopologies();
		expect(response).toHaveSize(initialLength - 1);
		expect(response).not.toContain(topology);
	});

	it("throws an error if asked to delete a non-existent Topology", async () => {
		await expectAsync(service.deleteTopology("")).toBeRejected();
	});

	it("creates a new Topology", async () => {
		const initialLength = (await service.getTopologies()).length;

		const newTop = {
			...topology,
			name: `${topology.name}-copy`,
			nodes: topology.nodes.map(n => ({cachegroup: n.cachegroup}))
		};
		const created = await service.createTopology(newTop);
		expect(created.description).toEqual(newTop.description);
		expect(created.name).toEqual(newTop.name);
		for (const node of created.nodes) {
			expect(node.parents).not.toBeNull();
			expect(node.parents).toBeDefined();
		}

		const response = await service.getTopologies();
		expect(response).toHaveSize(initialLength+1);
		expect(response).toContain(created);
	});

	it("updates existing Topologies", async () => {
		const initialLength = (await service.getTopologies()).length;

		const description = `${topology.description} - changed`;
		// need to decouple the references
		topology = {...topology};
		const updated = await service.updateTopology({...topology, description});
		expect(updated.description).toEqual(description);

		const response = await service.getTopologies();
		expect(response).toHaveSize(initialLength);
		expect(response).not.toContain(topology);
		expect(response).toContain(updated);
	});

	it("throws an error when asked to update a non-existent Topology", async () => {
		await expectAsync(service.updateTopology(topology, "")).toBeRejected();
	});

	it("converts Topologies to trees just like the concrete service", () => {
		expect(ConcreteTopologyService.topologyToTree(topology)).toEqual(TopologyService.topologyToTree(topology));
	});

	it("converts trees to Topologies just like the concrete service", () => {
		const tree = ConcreteTopologyService.topologyToTree(topology);
		expect(
			ConcreteTopologyService.treeToTopology(topology.name, topology.description, tree)
		).toEqual(TopologyService.treeToTopology(topology.name, topology.description, tree)
		);
	});
});
