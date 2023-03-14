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
import { Location } from "@angular/common";
import { fakeAsync, TestBed, tick } from "@angular/core/testing";
import { Router } from "@angular/router";
import { RouterTestingModule } from "@angular/router/testing";
import type { ResponseCurrentUser } from "trafficops-types";

import { UserService } from "src/app/api";
import { APITestingModule } from "src/app/api/testing";
import { LoginComponent } from "src/app/login/login.component";

import { CurrentUserService } from "./current-user.service";

/**
 * Creates a new user for testing purposes.
 *
 * @returns A new current user.
 */
function newCurrentUser(): ResponseCurrentUser {
	return {
		addressLine1: "addressLine1",
		addressLine2: "addressLine2",
		changeLogCount: 2,
		city: "city",
		company: "company",
		country: "country",
		email: "em@i.l",
		fullName: "fullName",
		gid: null,
		id: 1,
		lastAuthenticated: null,
		lastUpdated: new Date(),
		localUser: true,
		newUser: false,
		phoneNumber: "phoneNumber",
		postalCode: "postalCode",
		publicSshKey: "publicSshKey",
		registrationSent: null,
		role: "roleName",
		stateOrProvince: "stateOrProvince",
		tenant: "tenant",
		tenantId: 1,
		ucdn: "",
		uid: null,
		username: "username"
	};
}

describe("CurrentUserService", () => {
	let service: CurrentUserService;
	let router: Router;
	let location: Location;

	beforeEach(() => {
		const mockAPIService = jasmine.createSpyObj(["updateCurrentUser", "getCurrentUser", "saveCurrentUser"]);
		mockAPIService.getCurrentUser.and.returnValue(new Promise<ResponseCurrentUser>(resolve => resolve(
			{id: 1, newUser: false, role: "admin", username: "name"} as ResponseCurrentUser
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
		router = TestBed.inject(Router);
		location = TestBed.inject(Location);
		router.initialNavigation();
	});

	it("should be created", () => {
		expect(service).toBeTruthy();
	});

	it("should clear user data on logout", () => {
		service.logout();
		expect(service.loggedIn).toBeFalse();
		expect(service.currentUser).toBeNull();
		expect(service.permissions.getValue().size).toBe(0);
	});

	it("should update user data properly", () => {
		const upd = new Date();
		service.setUser(
			{
				addressLine1: "address line 1",
				addressLine2: "address line 2",
				changeLogCount: 2,
				city: "city",
				company: "company",
				country: "country",
				email: "em@i.l",
				fullName: "full name",
				gid: 0,
				id: 9000,
				lastAuthenticated: null,
				lastUpdated: upd,
				localUser: true,
				newUser: false,
				phoneNumber: "7",
				postalCode: "also 7",
				publicSshKey: "ssh key",
				registrationSent: null,
				role: "role name",
				stateOrProvince: "state or province",
				tenant: "tenant",
				tenantId: 2,
				ucdn: "",
				uid: 3,
				username: "quest",
			},
			new Set(["a permission"])
		);

		service.logout();
		service.setUser(
			{
				addressLine1: null,
				addressLine2: null,
				changeLogCount: 2,
				city: null,
				company: null,
				country: null,
				email: "different em@i.l",
				fullName: "different full name",
				gid: null,
				id: 9001,
				lastAuthenticated: null,
				lastUpdated: new Date(upd.getTime() + 1000),
				localUser: false,
				newUser: true,
				phoneNumber: null,
				postalCode: null,
				publicSshKey: null,
				registrationSent: null,
				role: "different role name",
				stateOrProvince: null,
				tenant: "different tenant",
				tenantId: 1,
				ucdn: "",
				uid: null,
				username: "test",
			},
			new Set()
		);

		const u = service.currentUser;
		if (u === null) {
			return fail("user is null after being set");
		}
		expect(u).toEqual({
			addressLine1: null,
			addressLine2: null,
			changeLogCount: 2,
			city: null,
			company: null,
			country: null,
			email: "different em@i.l",
			fullName: "different full name",
			gid: null,
			id: 9001,
			lastAuthenticated: null,
			lastUpdated: new Date(upd.getTime()+1000),
			localUser: false,
			newUser: true,
			phoneNumber: null,
			postalCode: null,
			publicSshKey: null,
			registrationSent: null,
			role: "different role name",
			stateOrProvince: null,
			tenant: "different tenant",
			tenantId: 1,
			ucdn: "",
			uid: null,
			username: "test",
		});
	});

	it("should update user permissions properly", () => {
		service.setUser(newCurrentUser(), new Set(["a different permission"]));
		expect(service.hasPermission("a different permission")).toBeTrue();

		service.logout();

		expect(service.hasPermission("a different permission")).toBeFalse();

		service.setUser(newCurrentUser(), new Set(["a permission"]));
		expect(service.hasPermission("a permission")).toBeTrue();
		expect(service.hasPermission("a different permission")).toBeFalse();

		service.setUser(newCurrentUser(), [{description: "", lastUpdated: new Date(), name: "a permission"}]);
		expect(service.hasPermission("a permission")).toBeTrue();
		expect(service.hasPermission("a different permission")).toBeFalse();
	});

	it("lets 'admin' users do things even when they don't have permission", () => {
		service.setUser({...newCurrentUser(), role: CurrentUserService.ADMIN_ROLE}, new Set());
		expect(service.permissions.getValue().has("a permission")).toBeFalse();
		expect(service.hasPermission("a permission")).toBeTrue();
	});

	it("fetches the current user when appropriate", async () => {
		const spy = spyOn(service, "updateCurrentUser").and.returnValue(new Promise(r=>r(true)));
		service.setUser(newCurrentUser(), []);
		expect(await service.fetchCurrentUser()).toBeTrue();
		expect(spy).not.toHaveBeenCalled();

		service.logout();
		expect(await service.fetchCurrentUser()).toBeTrue();
		expect(spy).toHaveBeenCalled();
	});

	it("logs users in", async () => {
		expect(await service.login("test-admin", "twelve12!")).toBeTrue();
		expect(await service.login("test-admin", "a misspelled password")).toBeFalse();
		expect(await service.login("there's no token that includes apostrophes")).toBeFalse();
	});

	it("logs users out", fakeAsync(() => {
		service.logout();
		tick();
		expect(location.path()).toBe("/login");
		expect(service.currentUser).toBeNull();

		service.setUser(newCurrentUser(), new Set(["perm"]));
		expect(service.currentUser).toBeTruthy();
		expect(service.hasPermission("perm")).toBeTrue();

		service.logout(true);
		tick();
		expect(location.path()).toBe(`/login?returnUrl=${encodeURIComponent("/core")}`);
		expect(service.currentUser).toBeNull();
		expect(service.hasPermission("perm")).toBeFalse();
	}));

	it("submits a request to update the current user", () => {
		const spy = spyOn(TestBed.inject(UserService), "updateCurrentUser");
		expect(spy).not.toHaveBeenCalled();
		service.saveCurrentUser(newCurrentUser());
		expect(spy).toHaveBeenCalled();
	});
});
