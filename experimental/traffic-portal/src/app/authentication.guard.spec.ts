import { TestBed } from "@angular/core/testing";

import { AuthenticationGuard } from "./authentication.guard";

describe("AuthenticationGuard", () => {
	let guard: AuthenticationGuard;

	beforeEach(() => {
		TestBed.configureTestingModule({});
		guard = TestBed.inject(AuthenticationGuard);
	});

	it("should be created", () => {
		expect(guard).toBeTruthy();
	});
});
