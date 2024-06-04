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

import { ComponentFixture, fakeAsync, TestBed, tick } from "@angular/core/testing";
import type { MatSelect } from "@angular/material/select";
import { RouterTestingModule } from "@angular/router/testing";

import { UserService } from "src/app/api";
import { APITestingModule } from "src/app/api/testing";
import { CurrentUserService } from "src/app/shared/current-user/current-user.service";
import { CurrentUserTestingService } from "src/app/shared/current-user/current-user.testing-service.spec";

import { UserDetailsComponent } from "./user-details.component";

describe("UserDetailsComponent", () => {
	let component: UserDetailsComponent;
	let fixture: ComponentFixture<UserDetailsComponent>;
	let service: UserService;

	beforeEach(async () => {
		await TestBed.configureTestingModule({
			declarations: [ UserDetailsComponent ],
			imports: [ APITestingModule, RouterTestingModule.withRoutes( [
				{component: UserDetailsComponent, path: "core/users"},
				{component: UserDetailsComponent, path: "core/users/:id"},
			])],
			providers: [{provide: CurrentUserService, useClass: CurrentUserTestingService}]
		}).compileComponents();
		fixture = TestBed.createComponent(UserDetailsComponent);
		component = fixture.componentInstance;
		service = TestBed.inject(UserService);
		component.roles = await service.getRoles();
		component.tenants = await service.getTenants();
		component.user = (await service.getUsers())[0];
		fixture.detectChanges();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});

	it("submits a user update request", fakeAsync(() => {
		const spy = spyOn(service, "updateUser");
		spy.and.callThrough();
		expect(spy).not.toHaveBeenCalled();
		if (component.isNew(component.user)) {
			return fail("user should not be new");
		}
		const before = component.user.lastUpdated;
		if (!before) {
			return fail("'lastUpdated' property was null or undefined or missing");
		}
		component.submit(new Event("submit"));
		tick();
		expect(component.user.lastUpdated).toBeGreaterThan(before.getTime());
	}));

	it("throws errors for non-existent roles and tenants", () => {
		component.user.role = "";
		expect(()=>component.role()).toThrow();
		component.user.tenantId = -1;
		expect(()=>component.tenant()).toThrow();
	});

	it("updates the user's Role-related properties", () => {
		const selectedRole = {
			capabilities: [],
			description: "",
			id: 2,
			name: "test",
			privLevel: 100
		};
		component.updateRole({source: {} as MatSelect, value: selectedRole});
		expect(component.user.role).toBe(selectedRole.name);
	});

	it("updates the user's Tenant-related properties", () => {
		const selectedTenant = {
			active: true,
			id: 2,
			lastUpdated: new Date(),
			name: "test",
			parentId: 1
		};
		component.updateTenant({source: {} as MatSelect, value: selectedTenant});
		expect(component.user.tenantId).toBe(selectedTenant.id);
	});

	it("has 'null' Role and Tenant for new users", () => {
		const oldValue = component.new;
		component.new = true;
		expect(component.role()).toBeNull();
		expect(component.tenant()).toBeNull();
		component.new = oldValue;
	});

	it("submits a user creation request", fakeAsync(() => {
		const spy = spyOn(service, "createUser");
		spy.and.callThrough();
		expect(spy).not.toHaveBeenCalled();
		const oldValue = component.new;
		const oldUser = component.user;
		component.new = true;
		component.user.tenantId = 1;
		component.submit(new Event("submit"));
		tick();
		expect(spy).toHaveBeenCalled();
		component.new = oldValue;
		component.user = oldUser;
	}));
});
