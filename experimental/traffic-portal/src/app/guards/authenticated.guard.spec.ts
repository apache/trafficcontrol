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
import { TestBed } from "@angular/core/testing";

import { CurrentUserService } from "src/app/shared/current-user/current-user.service";

import { AuthenticatedGuard } from "./authenticated-guard.service";

describe("AuthenticatedGuard", () => {
	let guard: AuthenticatedGuard;
	let mockCurrentUserService: jasmine.SpyObj<CurrentUserService>;
	let authPasses: boolean;

	beforeEach(() => {
		mockCurrentUserService = jasmine.createSpyObj(["fetchCurrentUser"]);
		mockCurrentUserService.fetchCurrentUser.and.callFake(async ()=>authPasses);
		TestBed.configureTestingModule({
			providers: [
				{ provide: CurrentUserService, useValue: mockCurrentUserService },
				AuthenticatedGuard
			]
		});
		authPasses = true;
		guard = TestBed.inject(AuthenticatedGuard);
	});

	it("should be created", () => {
		expect(guard).toBeTruthy();
	});

	it("checks activation criteria", async () => {
		expect(await guard.canActivate()).toBeTrue();
		authPasses = false;
		expect(await guard.canActivate()).toBeFalse();
	});

	it("checks load criteria", async () => {
		expect(await guard.canLoad()).toBeTrue();
		authPasses = false;
		expect(await guard.canLoad()).toBeFalse();
	});
});
