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
import { HttpClientTestingModule, HttpTestingController } from "@angular/common/http/testing";
import { TestBed } from "@angular/core/testing";

import { TopologyService } from "./topology.service";

describe("TopologyService", () => {
	let service: TopologyService;
	let httpTestingController: HttpTestingController;
	const topology = {
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
	};

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

	it("gets multiple Topologies", async () => {
		const responseP = service.getTopologies();
		const req = httpTestingController.expectOne(`/api/${service.apiVersion}/topologies`);
		expect(req.request.method).toBe("GET");
		expect(req.request.params.keys().length).toBe(0);
		req.flush({response: [topology]});
		await expectAsync(responseP).toBeResolvedTo([topology]);
	});

	afterEach(() => {
		httpTestingController.verify();
	});
});
