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

import { CurrentUserService } from "src/app/shared/currentUser/current-user.service";
import { AuthenticatedGuard } from "./authenticated-guard.service";

describe("AuthenticationGuard", () => {
	let guard: AuthenticatedGuard;

	beforeEach(() => {
		const mockCurrentUserService = jasmine.createSpyObj(["updateCurrentUser", "login", "logout"]);
		TestBed.configureTestingModule({
			providers: [
				{ provide: CurrentUserService, useValue: mockCurrentUserService },
				AuthenticatedGuard
			]
		});
	});

	beforeEach(() => {
		guard = TestBed.inject(AuthenticatedGuard);
	});

	it("should be created", () => {
		expect(guard).toBeTruthy();
	});
});
