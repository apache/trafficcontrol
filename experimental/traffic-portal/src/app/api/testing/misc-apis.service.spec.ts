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

import type { MiscAPIsService } from "../misc-apis.service";

import { MiscAPIsService as TestingMiscAPIsService } from "./misc-apis.service";

import { APITestingModule } from ".";

describe("TestingMiscAPIsService", () => {
	let service: TestingMiscAPIsService;

	beforeEach(() => {
		TestBed.configureTestingModule({
			imports: [APITestingModule]
		});
		service = TestBed.inject(TestingMiscAPIsService);
	});

	it("should be created", () => {
		expect(service).toBeTruthy();
	});

	it("returns a static set of mock os versions", async () => {
		expect(await service.getISOOSVersions()).toEqual(service.osVersions);
	});

	it("gives back an empty blob when requesting an ISO be generated, no matter what you give it", async () => {
		let blob = await service.generateISO();
		expect(blob.size).toBe(0);
		blob = await (service as unknown as MiscAPIsService).generateISO({
			dhcp: "yes",
			disk: "sda",
			domainName: "domain-name",
			hostName: "host-name",
			interfaceMtu: 1500,
			osVersionDir: "a version that doesn't even exist",
			rootPass: ""
		});
		expect(blob.size).toBe(0);
	});
});
