import { ComponentFixture, fakeAsync, TestBed, tick } from "@angular/core/testing";
import type { MatSelect } from "@angular/material/select";
import { RouterTestingModule } from "@angular/router/testing";

import { UserService } from "src/app/api";
import { APITestingModule } from "src/app/api/testing";

import { UserDetailsComponent } from "./user-details.component";

describe("UserDetailsComponent", () => {
	let component: UserDetailsComponent;
	let fixture: ComponentFixture<UserDetailsComponent>;
	let service: UserService;

	beforeEach(async () => {
		await TestBed.configureTestingModule({
			declarations: [ UserDetailsComponent ],
			imports: [ APITestingModule, RouterTestingModule ]
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
		const before = component.user.lastUpdated;
		if (!before) {
			return fail("'lastUpdated' property was null or undefined or missing");
		}
		component.submit(new Event("submit"));
		tick();
		expect(component.user.lastUpdated).toBeGreaterThan(before.getTime());
	}));

	it("throws errors for non-existent roles and tenants", () => {
		component.user.role = -1;
		expect(()=>component.role()).toThrow();
		component.user.tenantId = -1;
		expect(()=>component.tenant()).toThrow();
	});

	it("updates the user's Role-related properties", () => {
		const selectedRole = {
			capabilities: [],
			id: 2,
			name: "test",
			privLevel: 100
		};
		component.updateRole({source: {} as MatSelect, value: selectedRole});
		expect(component.user.role).toBe(selectedRole.id);
		expect(component.user.rolename).toBe(selectedRole.name);
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
		expect(component.user.tenant).toBe(selectedTenant.name);
	});
});
