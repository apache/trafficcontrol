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
import { HttpParams } from "@angular/common/http";
import { TestBed } from "@angular/core/testing";

import { ChangeLogsService } from "./change-logs.service";

import { APITestingModule } from ".";

describe("ChangeLogsService", () => {
	let service: ChangeLogsService;

	beforeEach(() => {
		TestBed.configureTestingModule({
			imports: [APITestingModule],
			providers: [
				ChangeLogsService,
			]
		});
		service = TestBed.inject(ChangeLogsService);
		expect(service.changeLogs.length).toBeGreaterThanOrEqual(1);
	});

	it("should be created", () => {
		expect(service).toBeTruthy();
	});

	it("gets Change Logs", async () => {
		await expectAsync(service.getChangeLogs()).toBeResolvedTo(service.changeLogs);
	});

	it("filters Change Logs by user", async () => {
		const {user} = service.changeLogs[0];
		let cls = await service.getChangeLogs({user});
		expect(cls.every(cl => cl.user === user));
		cls = await service.getChangeLogs(new HttpParams({fromObject: {user}}));
		expect(cls.every(cl => cl.user === user));
	});
});
