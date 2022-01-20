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

import { APITestingModule } from "src/app/api/testing";
import { LoginComponent } from "src/app/login/login.component";
import { newCurrentUser, User } from "src/app/models";

import { CurrentUserService } from "./current-user.service";

describe("CurrentUserService", () => {
	let service: CurrentUserService;

	beforeEach(() => {
		const mockAPIService = jasmine.createSpyObj(["updateCurrentUser", "getCurrentUser", "saveCurrentUser"]);
		mockAPIService.getCurrentUser.and.returnValue(new Promise<User>(resolve => resolve(
			{id: 1, newUser: false, role: 1, username: "name"}
		)));
		TestBed.configureTestingModule({
			imports: [
				APITestingModule,
				RouterTestingModule.withRoutes([{component: LoginComponent, path: "login"}])
			],
			providers: [
				CurrentUserService,
			]
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
		const upd = new Date();
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
				lastUpdated: upd,
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
				addressLine1: null,
				addressLine2: null,
				city: null,
				company: null,
				country: null,
				email: "different email",
				fullName: "different full name",
				gid: null,
				id: 9001,
				lastUpdated: new Date(upd.getTime()+1000),
				localUser: false,
				newUser: true,
				phoneNumber: null,
				postalCode: null,
				publicSshKey: null,
				role: 2,
				roleName: "different role name",
				stateOrProvince: null,
				tenant: "different tenant",
				tenantId: 1,
				uid: null,
				username: "test"
			},
			new Set()
		);

		const u = service.currentUser;
		expect(u).not.toBeNull();
		if (u !== null) {
			expect(u.addressLine1).toBeNull();
			expect(u.addressLine2).toBeNull();
			expect(u.city).toBeNull();
			expect(u.company).toBeNull();
			expect(u.country).toBeNull();
			expect(u.email).toBe("different email");
			expect(u.fullName).toBe("different full name");
			expect(u.gid).toBeNull();
			expect(u.id).toEqual(9001);
			expect(u.lastUpdated).toEqual(new Date(upd.getTime()+1000));
			expect(u.localUser).toBeFalse();
			expect(u.newUser).toBeTrue();
			expect(u.phoneNumber).toBeNull();
			expect(u.postalCode).toBeNull();
			expect(u.publicSshKey).toBeNull();
			expect(u.role).toEqual(2);
			expect(u.roleName).toBe("different role name");
			expect(u.stateOrProvince).toBeNull();
			expect(u.tenant).toBe("different tenant");
			expect(u.tenantId).toEqual(1);
			expect(u.uid).toBeNull();
			expect(u.username).toBe("test");
		}
	});

	it("should update user permissions properly", () => {
		service.setUser(newCurrentUser(), new Set(["a different permission"]));
		service.logout();
		service.setUser(newCurrentUser(), new Set(["a permission"]));
		expect(service.capabilities.has("a permission")).toBeTrue();
		expect(service.hasPermission("a permission")).toBeTrue();
		expect(service.capabilities.has("a different permission")).toBeFalse();
		expect(service.hasPermission("a different permission")).toBeFalse();
	});

});
