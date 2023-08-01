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
import {
	HttpClientTestingModule,
	HttpTestingController
} from "@angular/common/http/testing";
import { TestBed } from "@angular/core/testing";

import { TopologyService, TopTreeNode } from "./topology.service";

describe("TopologyService", () => {
	let service: TopologyService;
	let httpTestingController: HttpTestingController;
	const topology = {
		description: "",
		lastUpdated: new Date(),
		name: "my-topology",
		nodes: [
			{
				cachegroup: "Origin",
				parents: [],
			},
			{
				cachegroup: "Mid",
				parents: [0],
			},
			{
				cachegroup: "Edge",
				parents: [1],
			},
		],
	};

	const tree: Array<TopTreeNode> = [{
		cachegroup: "Origin",
		children: [{
			cachegroup: "Mid",
			children: [{
				cachegroup: "Edge",
				children: [],
				name: "Edge",
				parents: []
			}],
			name: "Mid",
			parents: []
		}],
		name: "Origin",
		parents: []
	}];

	beforeEach(() => {
		TestBed.configureTestingModule({
			imports: [HttpClientTestingModule],
			providers: [
				TopologyService,
			]
		});
		service = TestBed.inject(TopologyService);
		httpTestingController = TestBed.inject(HttpTestingController);
	});

	it("should be created", () => {
		expect(service).toBeTruthy();
	});

	it("gets a topology by name", async () => {
		const name = "my-topology";
		const responseP = service.getTopologies(name);
		const req = httpTestingController.expectOne(`/api/${service.apiVersion}/topologies?name=${name}`);
		expect(req.request.method).toBe("GET");
		expect(req.request.params.keys().length).toBe(1);
		req.flush({response: [topology]});
		await expectAsync(responseP).toBeResolvedTo([topology]);
	});

	it("throws an error when no Topology has the given name", async () => {
		const name = "nonexistent";
		const responseP = service.getTopologies(name);
		const req = httpTestingController.expectOne(`/api/${service.apiVersion}/topologies?name=${name}`);
		expect(req.request.method).toBe("GET");
		expect(req.request.params.keys().length).toBe(1);
		req.flush({response: [topology]});
		await expectAsync(responseP).toBeRejected();
	});

	it("gets multiple Topologies", async () => {
		const responseP = service.getTopologies();
		const req = httpTestingController.expectOne(`/api/${service.apiVersion}/topologies`);
		expect(req.request.method).toBe("GET");
		expect(req.request.params.keys().length).toBe(0);
		req.flush({response: [topology]});
		await expectAsync(responseP).toBeResolvedTo([topology]);
	});

	it("creates a new Topology", async () => {
		const responseP = service.createTopology(topology);
		const req = httpTestingController.expectOne(`/api/${service.apiVersion}/topologies`);
		expect(req.request.method).toBe("POST");
		expect(req.request.params.keys().length).toBe(0);
		expect(req.request.body).toBe(topology);
		req.flush({response: topology});
		await expectAsync(responseP).toBeResolved();
	});

	it("updates an existing Topology", async () => {
		const responseP = service.updateTopology(topology);
		const req = httpTestingController.expectOne(`/api/${service.apiVersion}/topologies?name=${topology.name}`);
		expect(req.request.method).toBe("PUT");
		expect(req.request.params.keys().length).toBe(1);
		expect(req.request.body).toBe(topology);
		req.flush({response: topology});
		await expectAsync(responseP).toBeResolved();
	});

	it("updates an existing Topology by name", async () => {
		const responseP = service.updateTopology(topology, topology.name);
		const req = httpTestingController.expectOne(`/api/${service.apiVersion}/topologies?name=${topology.name}`);
		expect(req.request.method).toBe("PUT");
		expect(req.request.params.keys().length).toBe(1);
		expect(req.request.body).toBe(topology);
		req.flush({response: topology});
		await expectAsync(responseP).toBeResolved();
	});

	it("deletes an existing Topology by name", async () => {
		const name = "my-topology";
		const responseP = service.deleteTopology(name);
		const req = httpTestingController.expectOne(`/api/${service.apiVersion}/topologies?name=${name}`);
		expect(req.request.method).toBe("DELETE");
		expect(req.request.params.keys().length).toBe(1);
		expect(req.request.body).toBeNull();
		req.flush({response: topology});
		await expectAsync(responseP).toBeResolved();
	});

	it("deletes an existing Topology", async () => {
		const responseP = service.deleteTopology(topology);
		const req = httpTestingController.expectOne(`/api/${service.apiVersion}/topologies?name=${topology.name}`);
		expect(req.request.method).toBe("DELETE");
		expect(req.request.params.keys().length).toBe(1);
		expect(req.request.body).toBeNull();
		req.flush({response: topology});
		await expectAsync(responseP).toBeResolved();
	});

	it("converts from a material tree to a topology", () => {
		const result = TopologyService.treeToTopology(topology.name, topology.description, tree);
		topology.lastUpdated = result.lastUpdated;
		expect(result).toEqual(topology);
	});

	it("converts from a topology to a material tree", () => {
		const result = TopologyService.topologyToTree(topology);
		result[0].children[0].parents = [];
		result[0].children[0].children[0].parents = [];
		expect(result).toEqual(tree);
	});

	afterEach(() => {
		httpTestingController.verify();
	});
});
