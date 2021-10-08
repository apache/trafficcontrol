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
import { RouterTestingModule } from "@angular/router/testing";
import { LoginComponent } from "../../login/login.component";

import { CurrentUserService } from "./current-user.service";

describe("CurrentUserService", () => {
	let service: CurrentUserService;

	beforeEach(() => {
		TestBed.configureTestingModule({
			imports: [RouterTestingModule.withRoutes([{component: LoginComponent, path: "login"}])],
			providers: [CurrentUserService]
		});
		service = TestBed.inject(CurrentUserService);
	});

	it("should be created", () => {
		expect(service).toBeTruthy();
	});

	it("should clear user data on logout", () => {
		service.logout();
		expect(service.loggedIn).toBeFalse();
		expect(service.currentUser).toBeNull();
		expect(service.capabilities.size).toBe(0);
	});

	it("should update user data properly", () => {
		service.setUser(
			{
				addressLine1: "address line 1",
				addressLine2: "address line 2",
				city: "city",
				company: "company",
				country: "country",
				email: "email",
				fullName: "full name",
				gid: 0,
				id: 9000,
				lastUpdated: new Date(),
				localUser: true,
				newUser: false,
				phoneNumber: "7",
				postalCode: "also 7",
				publicSshKey: "ssh key",
				role: 1,
				roleName: "role name",
				stateOrProvince: "state or province",
				tenant: "tenant",
				tenantId: 2,
				uid: 3,
				username: "quest"
			},
			new Set(["a permission"])
		);

		service.logout();
		service.setUser(
			{
				id: 9001,
				newUser: true,
				username: "test"
			},
			new Set()
		);

		const u = service.currentUser;
		expect(u).not.toBeNull();
		if (u !== null) {
			expect(u.addressLine1).toBeUndefined();
			expect(u.addressLine2).toBeUndefined();
			expect(u.city).toBeUndefined();
			expect(u.company).toBeUndefined();
			expect(u.country).toBeUndefined();
			expect(u.email).toBeUndefined();
			expect(u.fullName).toBeUndefined();
			expect(u.gid).toBeUndefined();
			expect(u.id).toEqual(9001);
			expect(u.lastUpdated).toBeUndefined();
			expect(u.localUser).toBeUndefined();
			expect(u.newUser).toBeTrue();
			expect(u.phoneNumber).toBeUndefined();
			expect(u.postalCode).toBeUndefined();
			expect(u.publicSshKey).toBeUndefined();
			expect(u.role).toBeUndefined();
			expect(u.roleName).toBeUndefined();
			expect(u.stateOrProvince).toBeUndefined();
			expect(u.tenant).toBeUndefined();
			expect(u.tenantId).toBeUndefined();
			expect(u.uid).toBeUndefined();
			expect(u.username).toBe("test");
		}
	});

	it("should update user permissions properly", () => {
		service.setUser(
			{
				id: 9001,
				newUser: true,
				username: "test"
			},
			new Set(["a different permission"])
		);
		service.logout();
		service.setUser(
			{
				id: 9001,
				newUser: true,
				username: "test"
			},
			new Set(["a permission"])
		);
		expect(service.capabilities.has("a permission")).toBeTrue();
		expect(service.hasPermission("a permission")).toBeTrue();
		expect(service.capabilities.has("a different permission")).toBeFalse();
		expect(service.hasPermission("a different permission")).toBeFalse();
	});

});
