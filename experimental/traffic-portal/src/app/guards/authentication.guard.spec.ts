import { TestBed } from "@angular/core/testing";

import { AuthenticatedGuard } from "./authenticated-guard.service";
import {AuthenticationService} from "../shared/authentication/authentication.service";

describe("AuthenticationGuard", () => {
	let guard: AuthenticatedGuard;

	beforeEach(() => {
		const mockAuthenticationService = jasmine.createSpyObj(["updateCurrentUser", "login", "logout"]);
		TestBed.configureTestingModule({
			providers: [
				{ provide: AuthenticationService, useValue: mockAuthenticationService },
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
