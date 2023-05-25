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
import { ResponseTenant, ResponseUser } from "trafficops-types";

import { UserService } from "./user.service";

describe("UserService", () => {
	let service: UserService;
	let httpTestingController: HttpTestingController;

	beforeEach(() => {
		TestBed.configureTestingModule({
			imports: [HttpClientTestingModule],
			providers: [
				UserService,
			]
		});
		service = TestBed.inject(UserService);
		httpTestingController = TestBed.inject(HttpTestingController);
	});

	it("should be created", () => {
		expect(service).toBeTruthy();
	});

	describe("authentication methods", () => {
		const username = "test";
		const password = "quest";
		const email = "em@i.l" as const;
		const role = "admin";

		const registrationRequest = {
			email,
			role,
			tenantId: 1
		};

		it("sends login requests using username/password combinations", async () => {
			const responseP = service.login(username, password);
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/user/login`);
			expect(req.request.method).toBe("POST");
			expect(req.request.body).toEqual({p: password, u: username});
			req.flush({alerts: []});
			await expectAsync(responseP).toBeResolved();
		});
		it("sends login requests using tokens", async () => {
			const responseP = service.login(password);
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/user/login/token`);
			expect(req.request.method).toBe("POST");
			expect(req.request.body).toEqual({t: password});
			req.flush({alerts: []});
			await expectAsync(responseP).toBeResolved();
		});
		it("sends logout requests", async () => {
			const responseP = service.logout();
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/user/logout`);
			expect(req.request.method).toBe("POST");
			expect(req.request.body).toBeNull();
			req.flush({alerts: []});
			await expectAsync(responseP).toBeResolved();
		});
		it("sends user registration requests", async () => {
			const responseP = service.registerUser(registrationRequest);
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/users/register`);
			expect(req.request.method).toBe("POST");
			expect(req.request.body).toEqual(registrationRequest);
			req.flush({alerts: []});
			await expectAsync(responseP).toBeResolved();
		});
		it("sends user registration requests in parts", async () => {
			const responseP = service.registerUser(email, role, registrationRequest.tenantId);
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/users/register`);
			expect(req.request.method).toBe("POST");
			expect(req.request.body).toEqual(registrationRequest);
			req.flush({alerts: []});
			await expectAsync(responseP).toBeResolved();
		});
		it("sends user registration requests in parts using full objects", async () => {
			const responseP = service.registerUser(
				email,
				{
					description: "description",
					lastUpdated: new Date(),
					name: role
				},
				{
					active: true,
					id: registrationRequest.tenantId,
					lastUpdated: new Date(),
					name: "root" as const,
					parentId: null,
				}
			);
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/users/register`);
			expect(req.request.method).toBe("POST");
			expect(req.request.body).toEqual(registrationRequest);
			req.flush({alerts: []});
			await expectAsync(responseP).toBeResolved();
		});
		it("throws an error for invalid call signatures to registerUser", async () => {
			const responseP = (service as unknown as {registerUser: (email: string) => Promise<void>}).registerUser(
				registrationRequest.email
			);
			httpTestingController.expectNone({method: "POST"});
			await expectAsync(responseP).toBeRejected();
		});
		it("throws an error when given an invalid email", async () => {
			const responseP = service.registerUser("not a valid email", "", 1);
			httpTestingController.expectNone({method: "POST"});
			await expectAsync(responseP).toBeRejected();
		});
		it("sends password reset requests", async () => {
			const responseP = service.resetPassword(email);
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/user/reset_password`);
			expect(req.request.method).toBe("POST");
			expect(req.request.body).toEqual({email});
			req.flush({alerts: []});
			await expectAsync(responseP).toBeResolved();
		});
	});

	describe("user methods", () => {
		const user = {
			addressLine1: null,
			addressLine2: null,
			changeLogCount: 1,
			city: null,
			company: null,
			get confirmLocalPasswd(): string {
				return this.localPasswd;
			},
			country: null,
			email: "em@i.l" as const,
			fullName: "full name",
			gid: null,
			id: 1,
			lastAuthenticated: null,
			lastUpdated: new Date(),
			localPasswd: "localPasswd",
			localUser: false,
			newUser: false,
			phoneNumber: null,
			postalCode: null,
			publicSshKey: null,
			registrationSent: null,
			role: "admin",
			stateOrProvince: null,
			tenant: "root",
			tenantId: 1,
			ucdn: "",
			uid: null,
			username: "username"
		};

		it("sends a request to get multiple users", async () => {
			const responseP = service.getUsers();
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/users`);
			expect(req.request.method).toBe("GET");
			req.flush({response: [user]});
			await expectAsync(responseP).toBeResolvedTo([user]);
		});

		it("sends a request to get a single user by ID", async () => {
			const responseP = service.getUsers(user.id);
			const req = httpTestingController.expectOne(r => r.url === `/api/${service.apiVersion}/users`);
			expect(req.request.method).toBe("GET");
			expect(req.request.params.keys().length).toBe(1);
			expect(req.request.params.get("id")).toBe(String(user.id));
			req.flush({response: [user]});
			await expectAsync(responseP).toBeResolvedTo(user);
		});

		it("sends a request to get a single user by username", async () => {
			const responseP = service.getUsers(user.username);
			const req = httpTestingController.expectOne(r => r.url === `/api/${service.apiVersion}/users`);
			expect(req.request.method).toBe("GET");
			expect(req.request.params.keys().length).toBe(1);
			expect(req.request.params.get("username")).toBe(user.username);
			req.flush({response: [user]});
			await expectAsync(responseP).toBeResolvedTo(user);
		});

		it("sends a request to create a new user", async () => {
			const responseP = service.createUser(user);
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/users`);
			expect(req.request.method).toBe("POST");
			expect(req.request.body).toEqual(user);
			req.flush({response: user});
			await expectAsync(responseP).toBeResolvedTo(user);
		});

		it("sends a request to update an existing user", async () => {
			const responseP = service.updateUser(user);
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/users/${user.id}`);
			expect(req.request.method).toBe("PUT");
			expect(req.request.body).toEqual(user);
			req.flush({response: user});
			await expectAsync(responseP).toBeResolvedTo(user);
		});

		it("sends a request to update an existing user by ID", async () => {
			const responseP = service.updateUser(user.id, user);
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/users/${user.id}`);
			expect(req.request.method).toBe("PUT");
			expect(req.request.body).toEqual(user);
			req.flush({response: user});
			await expectAsync(responseP).toBeResolvedTo(user);
		});
		it("throws an error for invalid call signatures to updateUser", async () => {
			const responseP = (service as unknown as {updateUser: (id: number) => Promise<ResponseUser>}).updateUser(
				user.id
			);
			httpTestingController.expectNone({method: "PUT"});
			await expectAsync(responseP).toBeRejected();
		});
		it("sends a request to get the current user", async () => {
			const responseP = service.getCurrentUser();
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/user/current`);
			expect(req.request.method).toBe("GET");
			req.flush({response: user});
			await expectAsync(responseP).toBeResolvedTo(user);
		});
		it("sends a request to update the current user", async () => {
			const responseP = service.updateCurrentUser(user);
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/user/current`);
			expect(req.request.method).toBe("PUT");
			expect(req.request.body).toEqual(user);
			req.flush({response: user});
			await expectAsync(responseP).toBeResolvedTo(true);
		});
		it("returns `false` from a request to update the current user when it fails", async () => {
			const responseP = service.updateCurrentUser(user);
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/user/current`);
			expect(req.request.method).toBe("PUT");
			expect(req.request.body).toEqual(user);
			req.flush({}, {status: 500, statusText: "Internal Server Error"});
			await expectAsync(responseP).toBeResolvedTo(false);
		});
	});

	describe("role methods", () => {
		const role = {
			description: "description",
			lastUpdated: new Date(),
			name: "test",
			parameters: ["ALL"]
		};

		it("sends a request for multiple Roles", async () => {
			const responseP = service.getRoles();
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/roles`);
			expect(req.request.method).toBe("GET");
			req.flush({response: [role]});
			await expectAsync(responseP).toBeResolvedTo([role]);
		});

		it("sends a request to get a single Role by name", async () => {
			const responseP = service.getRoles(role.name);
			const req = httpTestingController.expectOne(r => r.url === `/api/${service.apiVersion}/roles`);
			expect(req.request.method).toBe("GET");
			expect(req.request.params.keys().length).toBe(1);
			expect(req.request.params.get("name")).toBe(role.name);
			req.flush({response: [role]});
			await expectAsync(responseP).toBeResolvedTo(role);
		});

		it("throws an error when the requested Role doesn't exist", async () => {
			const responseP = service.getRoles(role.name);
			const req = httpTestingController.expectOne(r => r.url === `/api/${service.apiVersion}/roles`);
			expect(req.request.method).toBe("GET");
			expect(req.request.params.keys().length).toBe(1);
			expect(req.request.params.get("name")).toBe(role.name);
			req.flush({response: []});
			await expectAsync(responseP).toBeRejected();
		});

		it("deletes an existing Role", async () => {
			const responseP = service.deleteRole(role);
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/roles?name=${role.name}`);
			expect(req.request.method).toBe("DELETE");
			expect(req.request.params.keys().length).toBe(1);
			expect(req.request.body).toBeNull();
			req.flush({response: role});
			await expectAsync(responseP).toBeResolved();
		});

		it("deletes a Role by name", async () => {
			const responseP = service.deleteRole(role.name);
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/roles?name=${role.name}`);
			expect(req.request.method).toBe("DELETE");
			expect(req.request.body).toBeNull();
			expect(req.request.params.keys().length).toBe(1);
			req.flush({alerts: []});
			await expectAsync(responseP).toBeResolved();
		});

		it("creates a new Role", async () => {
			const responseP = service.createRole(role);
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/roles`);
			expect(req.request.method).toBe("POST");
			expect(req.request.body).toEqual(role);
			req.flush({response: role});
			await expectAsync(responseP).toBeResolved();
		});

		it("updates an existing Role", async () => {
			const responseP = service.updateRole(role.name, role);
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/roles?name=${role.name}`);
			expect(req.request.method).toBe("PUT");
			expect(req.request.body).toEqual(role);
			req.flush({response: role});
			await expectAsync(responseP).toBeResolved();
		});
	});

	describe("tenant-related methods", () => {
		const root: ResponseTenant = {
			active: true,
			id: 1,
			lastUpdated: new Date(),
			name: "root",
			parentId: null,
		};
		const tenant = {
			active: true,
			id: root.id+1,
			lastUpdated: new Date(),
			name: "testquest",
			parentId: root.id,
		};

		it("sends a request for multiple tenants", async () => {
			const responseP = service.getTenants();
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/tenants`);
			expect(req.request.method).toBe("GET");
			req.flush({response: [root, tenant]});
			await expectAsync(responseP).toBeResolvedTo([root, tenant]);
		});

		it("sends a request for a single tenant by ID", async () => {
			const responseP = service.getTenants(tenant.id);
			const req = httpTestingController.expectOne(r => r.url === `/api/${service.apiVersion}/tenants`);
			expect(req.request.method).toBe("GET");
			expect(req.request.params.keys().length).toBe(1);
			expect(req.request.params.get("id")).toBe(String(tenant.id));
			req.flush({response: [tenant]});
			await expectAsync(responseP).toBeResolvedTo(tenant);
		});

		it("sends a request for a single tenant by name", async () => {
			const responseP = service.getTenants(tenant.name);
			const req = httpTestingController.expectOne(r => r.url === `/api/${service.apiVersion}/tenants`);
			expect(req.request.method).toBe("GET");
			expect(req.request.params.keys().length).toBe(1);
			expect(req.request.params.get("name")).toBe(tenant.name);
			req.flush({response: [tenant]});
			await expectAsync(responseP).toBeResolvedTo(tenant);
		});
		it("sends a request to create a new tenant", async () => {
			const responseP = service.createTenant(tenant);
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/tenants`);
			expect(req.request.method).toBe("POST");
			expect(req.request.body).toEqual(tenant);
			req.flush({response: tenant});
			await expectAsync(responseP).toBeResolvedTo(tenant);
		});
		it("sends a request to update an existing tenant", async () => {
			const responseP = service.updateTenant(tenant);
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/tenants/${tenant.id}`);
			expect(req.request.method).toBe("PUT");
			expect(req.request.body).toEqual(tenant);
			req.flush({response: tenant});
			await expectAsync(responseP).toBeResolvedTo(tenant);
		});
		it("sends a request to delete a tenant", async () => {
			const responseP = service.deleteTenant(tenant);
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/tenants/${tenant.id}`);
			expect(req.request.method).toBe("DELETE");
			req.flush	({response: tenant});
			await expectAsync(responseP).toBeResolvedTo(tenant);
		});
		it("sends a request to delete a tenant by ID", async () => {
			const responseP = service.deleteTenant(tenant.id);
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/tenants/${tenant.id}`);
			expect(req.request.method).toBe("DELETE");
			req.flush	({response: tenant});
			await expectAsync(responseP).toBeResolvedTo(tenant);
		});
	});

	afterEach(() => {
		httpTestingController.verify();
	});
});
