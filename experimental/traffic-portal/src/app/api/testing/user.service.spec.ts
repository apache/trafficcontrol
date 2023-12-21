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
import type { ResponseCurrentUser, ResponseRole, ResponseTenant, ResponseUser } from "trafficops-types";

import { UserService } from "./user.service";

import { APITestingModule } from ".";

describe("UserService", () => {
	let service: UserService;

	beforeEach(() => {
		TestBed.configureTestingModule({
			imports: [APITestingModule],
			providers: [
				UserService,
			]
		});
		service = TestBed.inject(UserService);
	});

	it("should be created", () => {
		expect(service).toBeTruthy();
	});

	describe("authentication methods", () => {
		let username: string;
		let password: string;
		const email = "em@i.l" as const;
		const role = "admin";

		beforeEach(() => {
			username = service.testAdminUsername;
			password = service.testAdminPassword;
		});

		it("pretends to login using username/password combinations", async () => {
			await expectAsync(service.login(username, password)).toBeResolved();
		});
		it("fails authentication that doesn't match the testing credential set", async () => {
			await expectAsync(service.login(`${username} - wrong`, password)).toBeResolvedTo(null);
			await expectAsync(service.login(username, `${password} - wrong`)).toBeResolvedTo(null);
		});

		it("pretends to login using tokens", async () => {
			service.tokens.set(password, service.users[0].email);
			await expectAsync(service.login(password)).toBeResolved();
		});
		it("fails to authenticate with a token when a token hasn't been set", async () => {
			await expectAsync(service.login(password)).toBeResolvedTo(null);
		});
		it("fails to authenticate when a login token refers to a non-existent user", async () => {
			service.tokens.set(password, email);
			await expectAsync(service.login(password)).toBeResolvedTo(null);
		});

		it("pretends to logout", async () => {
			await expectAsync(service.logout()).toBeResolved();
		});

		it("pretends to send user registration requests", async () => {
			await expectAsync(service.registerUser(email, role, 1)).toBeResolved();
		});
		it("throws an error when registration is called with an invalid call signature", async () => {
			await expectAsync((service as unknown as {registerUser: (email: string) => Promise<void>}).registerUser(email)).toBeRejected();
		});

		it("sets tokens when password reset is requested", async () => {
			const initialSize = service.tokens.size;
			const realEmail = service.users[0].email;
			await expectAsync(service.resetPassword(realEmail)).toBeResolved();
			expect(service.tokens).toHaveSize(initialSize+1);
			expect(Array.from(service.tokens.values())).toContain(realEmail);
		});
		it("does nothing when password reset is requested for a non-existent user", async () => {
			const initialSize = service.tokens.size;
			await expectAsync(service.resetPassword("")).toBeResolved();
			expect(service.tokens).toHaveSize(initialSize);
		});
	});

	describe("user methods", () => {
		let user: ResponseUser;
		let currentUser: ResponseCurrentUser;

		beforeEach(async () => {
			const initialSize = service.users.length;
			user = await service.createUser({
				confirmLocalPasswd: "test",
				email: "user.em@i.l",
				fullName: "Test McQuestington",
				localPasswd: "test",
				role: "admin",
				tenantId: 1,
				username: "testMcQuestington"
			});
			expect(service.users).toHaveSize(initialSize+1);
			expect(service.users).toContain(user);
			expect(service.users.length).toBeGreaterThanOrEqual(2);
			currentUser = await service.getCurrentUser();
			expect(currentUser).toBeTruthy();
		});

		it("gets multiple users", async () => {
			await expectAsync(service.getUsers()).toBeResolvedTo(service.users);
		});
		it("gets a single user by ID", async () => {
			await expectAsync(service.getUsers(user.id)).toBeResolvedTo(user);
		});
		it("gets a single user by username", async () => {
			await expectAsync(service.getUsers(user.username)).toBeResolvedTo(user);
		});
		it("throws an error when asked to retrieve a non-existent user", async () => {
			await expectAsync(service.getUsers(-1)).toBeRejected();
			await expectAsync(service.getUsers("")).toBeRejected(service.users.map(u=>u.username));
		});

		it("throws an error when creating a new user with an invalid Role", async () => {
			await expectAsync(service.createUser({
				confirmLocalPasswd: "test",
				email: `a-different-${user.email}`,
				fullName: "Test McQuestington",
				localPasswd: "test",
				role: "",
				tenantId: user.tenantId,
				username: `a-different-${user.username}`
			})).toBeRejected();
		});
		it("throws an error when creating a new user with an invalid Tenant", async () => {
			await expectAsync(service.createUser({
				confirmLocalPasswd: "test",
				email: `a-different-${user.email}`,
				fullName: "Test McQuestington",
				localPasswd: "test",
				role: user.role,
				tenantId: -1,
				username: `a-different-${user.username}`
			})).toBeRejected();
		});

		it("updates an existing user", async () => {
			const current = {...user};
			current.addressLine1 = `${user.addressLine1 ?? ""} But not actually}`;
			const initialSize = service.users.length;
			const updated = await service.updateUser(current);
			expect(service.users).toContain(updated);
			expect(service.users).toHaveSize(initialSize);
		});
		it("updates an existing user by ID", async () => {
			const current = {...user};
			current.fullName = `${user.fullName} But not actually}`;
			const initialSize = service.users.length;
			const updated = await service.updateUser(current.id, current);
			expect(service.users).toContain(updated);
			expect(service.users).toHaveSize(initialSize);
		});
		it("throws an error when asked to update a non-existent user", async () => {
			await expectAsync(service.updateUser(-1, user)).toBeRejected();
			await expectAsync(service.updateUser({...user, id: -1})).toBeRejected();
		});
		it("throws an error for invalid call signatures to updateUser", async () => {
			const responseP = (service as unknown as {updateUser: (id: number) => Promise<ResponseUser>}).updateUser(
				user.id
			);
			await expectAsync(responseP).toBeRejected();
		});

		it("throws an error when the current user doesn't exist", async () => {
			service.users.splice(0, service.users.length);
			await expectAsync(service.getCurrentUser()).toBeRejected();
		});
		it("returns a fake current user if the actual current user cannot be determined", async () => {
			service.testAdminUsername = "";
			await expectAsync(service.getCurrentUser()).toBeResolved();
		});

		it("updates the current user", async () => {
			const current = {...currentUser};
			current.addressLine1 = `${currentUser.addressLine1 ?? ""} But not actually`;
			await expectAsync(service.updateCurrentUser(current)).toBeResolvedTo(true);
		});
		it("returns `false` from a request to update the current user when it fails", async () => {
			const responseP = service.updateCurrentUser({...currentUser, id: -1});
			await expectAsync(responseP).toBeResolvedTo(false);
		});
	});

	describe("role methods", () => {
		let role: ResponseRole;

		beforeEach(() => {
			expect(service.roles.length).toBeGreaterThanOrEqual(2);
			role = service.roles[1];
		});

		it("gets multiple Roles", async () => {
			await expectAsync(service.getRoles()).toBeResolvedTo(service.roles);
		});
		it("gets a single Role by name", async () => {
			await expectAsync(service.getRoles(role.name)).toBeResolvedTo(role);
		});
		it("throws an error when the requested Role doesn't exist", async () => {
			await expectAsync(service.getRoles("")).toBeRejected();
		});

		it("deletes an existing Role", async () => {
			const initialSize = service.roles.length;
			await expectAsync(service.deleteRole(role)).toBeResolved();
			expect(service.roles).toHaveSize(initialSize - 1);
			expect(service.roles).not.toContain(role);
		});
		it("throws an error when asked to delete a non-existent Role", async () => {
			await expectAsync(service.deleteRole("")).toBeRejected();
		});

		it("creates a new Role", async () => {
			const initialSize = service.roles.length;
			const created = await service.createRole({
				description: "",
				name: "creaiton testing role",
			});
			expect(service.roles).toHaveSize(initialSize+1);
			expect(service.roles).toContain(created);
		});

		it("updates an existing Role", async () => {
			const initialSize = service.roles.length;
			const current = {...role};
			current.description += " and some more stuff";
			const updated = await service.updateRole(role.name, current);
			expect(service.roles).toHaveSize(initialSize);
			expect(service.roles).toContain(updated);
		});
		it("throws an error when asked to update a non-existent Role", async () => {
			await expectAsync(service.updateRole("", role)).toBeRejected();
		});
	});

	describe("tenant-related methods", () => {
		let root: ResponseTenant & {active: true; name: "root"; parentId: null};
		let tenant: ResponseTenant;

		beforeEach(() => {
			const tenants = [...service.tenants];
			if (tenants.length < 2) {
				return fail("need at least 2 testing Tenants to run the Tenant testing service tests");
			}
			const idx = tenants.findIndex(
				t => t.active && t.name === "root" && t.parentId === null
			);
			if (idx < 0) {
				return fail("no (valid) root tenant exists in the testing service");
			}
			// This is safe because we explicitly checked for these conditions
			// when finding the index.
			root = tenants.splice(idx, 1)[0] as ResponseTenant & {active: true; name: "root"; parentId: null};
			tenant = tenants[0];
		});

		it("gets multiple Tenants", async () => {
			await expectAsync(service.getTenants()).toBeResolvedTo(service.tenants);
		});
		it("gets a single Tenant by ID", async () => {
			await expectAsync(service.getTenants(tenant.id)).toBeResolvedTo(tenant);
		});
		it("gets a single Tenant by name", async () => {
			await expectAsync(service.getTenants(root.name)).toBeResolvedTo(root);
		});
		it("throws an error if asked to get a non-existent Tenant", async () => {
			await expectAsync(service.getTenants(-1)).toBeRejected();
		});

		it("creates a new Tenant", async () => {
			const initialSize = service.tenants.length;
			const created = await service.createTenant({
				active: true,
				name: "creation testing tenant",
				parentId: tenant.id,
			});
			expect(service.tenants).toHaveSize(initialSize+1);
			expect(service.tenants).toContain(created);
		});

		it("updates existing Tenants", async () => {
			const current = {...tenant};
			current.name += " - update test";
			const initialSize = service.tenants.length;
			const updated = await service.updateTenant(current);
			expect(service.tenants).toHaveSize(initialSize);
			expect(service.tenants).toContain(updated);
			expect(service.tenants).not.toContain(tenant);
		});
		it("throws an error when asked to update a non-existent Tenant", async () => {
			await expectAsync(service.updateTenant({...tenant, id: -1})).toBeRejected();
		});

		it("deletes Tenants", async () => {
			const initialSize = service.tenants.length;
			await expectAsync(service.deleteTenant(tenant)).toBeResolvedTo(tenant);
			expect(service.tenants).toHaveSize(initialSize - 1);
			expect(service.tenants).not.toContain(tenant);
		});
		it("throws an error when asked to delete a non-existent Tenant", async () => {
			await expectAsync(service.deleteTenant(-1)).toBeRejected();
		});
	});
});
