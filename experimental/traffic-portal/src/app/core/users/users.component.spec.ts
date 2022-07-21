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
import { HttpClientModule } from "@angular/common/http";
import { waitForAsync, ComponentFixture, TestBed, fakeAsync, tick } from "@angular/core/testing";
import { FormsModule, ReactiveFormsModule } from "@angular/forms";
import { MatDialogModule } from "@angular/material/dialog";
import { RouterTestingModule } from "@angular/router/testing";
import type { ValueGetterParams } from "ag-grid-community";

import { APITestingModule } from "src/app/api/testing";
import { CurrentUserService } from "src/app/shared/currentUser/current-user.service";
import { isAction } from "src/app/shared/generic-table/generic-table.component";
import { LoadingComponent } from "src/app/shared/loading/loading.component";
import { TpHeaderComponent } from "src/app/shared/tp-header/tp-header.component";

import { UsersComponent } from "./users.component";

describe("UsersComponent", () => {
	let component: UsersComponent;
	let fixture: ComponentFixture<UsersComponent>;
	const testUser = {
		addressLine1: null,
		addressLine2: null,
		city: null,
		company: null,
		country: null,
		email: "a@b.c" as const,
		fullName: "admin",
		gid: null,
		id: 1,
		lastUpdated: new Date(0),
		newUser: false,
		phoneNumber: null,
		postalCode: null,
		publicSshKey: null,
		role: 1,
		rolename: "admin",
		stateOrProvince: null,
		tenant: "root",
		tenantId: 1,
		uid: null,
		username: "admin"
	};

	beforeEach(waitForAsync(() => {
		// mock the API
		const mockCurrentUserService = jasmine.createSpyObj(["updateCurrentUser", "hasPermission", "login", "logout"]);
		mockCurrentUserService.updateCurrentUser.and.returnValue(new Promise(r => r(false)));
		mockCurrentUserService.hasPermission.and.returnValue(true);

		TestBed.configureTestingModule({
			declarations: [
				UsersComponent,
				LoadingComponent,
				TpHeaderComponent,
			],
			imports: [
				APITestingModule,
				FormsModule,
				HttpClientModule,
				ReactiveFormsModule,
				RouterTestingModule,
				MatDialogModule
			],
			providers: [
				{ provide: CurrentUserService, useValue: mockCurrentUserService }
			]
		});
		TestBed.compileComponents();
	}));

	beforeEach(() => {
		fixture = TestBed.createComponent(UsersComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it("should exist", () => {
		expect(component).toBeTruthy();
	});

	it("updates the fuzzy search output", fakeAsync(() => {
		let called = false;
		const text = "testquest";
		const spy = jasmine.createSpy("subscriber", (txt: string): void =>{
			if (!called) {
				expect(txt).toBe("");
				called = true;
			} else {
				expect(txt).toBe(text);
			}
		});
		component.fuzzySubject.subscribe(spy);
		tick();
		expect(spy).toHaveBeenCalled();
		component.searchText = text;
		component.updateURL();
		tick();
		expect(spy).toHaveBeenCalledTimes(2);
	}));

	it("gets display strings for Roles", () => {
		component.roles = new Map([[1, "admin"]]);
		const roleColDef = component.columnDefs.find(d=>d.field === "role");
		if (!roleColDef) {
			return fail("table missing column definition for the 'role' property");
		}
		if (!roleColDef.valueGetter) {
			return fail("column definition for 'role' property missing 'valueGetter' property");
		}
		expect(roleColDef.valueGetter({data: testUser} as ValueGetterParams)).toBe(`${testUser.username} (#${testUser.id})`);
		expect(()=>component.roleDisplayString(2)).toThrow();
	});

	it("gets display strings for Tenants", () => {
		const tenantColDef = component.columnDefs.find(d=>d.field === "tenant");
		if (!tenantColDef) {
			return fail("table missing column definition for the 'tenant' property");
		}
		if (!tenantColDef.valueGetter) {
			return fail("column definition for 'tenant' property missing 'valueGetter' property");
		}
		expect(tenantColDef.valueGetter({data: testUser} as ValueGetterParams)).toBe(`${testUser.tenant} (#${testUser.tenantId})`);
	});

	it("has a proper 'View User Details' context menu item", () => {
		const item = component.contextMenuItems[0];
		if (!item) {
			return fail("table is missing 'contextMenuItems' property");
		}
		if (item.name !== "View User Details") {
			return fail(`The first context menu item should've been 'View User Details', but it was '${item.name}'`);
		}
		if (isAction(item)) {
			return fail("the first context menu item should've been a link but it was an action");
		}
		if (!item.disabled) {
			return fail("missing 'disabled' property");
		}
		expect(item.disabled([])).toBe(true, "should be disabled for multiple selected users");
		expect(item.disabled(testUser)).toBe(false, "should be enabled for single selected user");
		if (!item.href) {
			return fail("missing 'href' property");
		}
		if (typeof(item.href) === "string") {
			return fail(`should use a function to generate an href, but uses static string: '${item.href}'`);
		}
		expect(item.href(testUser)).toBe(`/core/users/${testUser.id}`, "generated incorrect href");
	});

	it("has a proper 'Open in New Tab' context menu item", () => {
		const item = component.contextMenuItems[1];
		if (!item) {
			return fail("table is missing a populated 'contextMenuItems' property");
		}
		if (item.name !== "Open in New Tab") {
			return fail(`The second context menu item should've been 'Open in New Tab', but it was '${item.name}'`);
		}
		if (isAction(item)) {
			return fail("the second context menu item should've been a link but it was an action");
		}
		expect(item.newTab).toBeTrue();
		if (!item.disabled) {
			return fail("missing 'disabled' property");
		}
		expect(item.disabled([])).toBe(true, "should be disabled for multiple selected users");
		expect(item.disabled(testUser)).toBe(false, "should be enabled for single selected user");
		if (!item.href) {
			return fail("missing 'href' property");
		}
		if (typeof(item.href) === "string") {
			return fail(`should use a function to generate an href, but uses static string: '${item.href}'`);
		}
		expect(item.href(testUser)).toBe(`/core/users/${testUser.id}`, "generated incorrect href");
	});

	it("has a proper 'View User Changelogs' context menu item", () => {
		const item = component.contextMenuItems[2];
		if (!item) {
			return fail("table is missing 'contextMenuItems' property");
		}
		if (item.name !== "View User Changelogs") {
			return fail(`The third context menu item should've been 'View User Changelogs', but it was '${item.name}'`);
		}
		if (!isAction(item)) {
			return fail("the third context menu item should've been an action but it was a link");
		}
		expect(item.action).toBe("viewChangelogs", "incorrect 'action' property");
		if (item.multiRow) {
			return fail("should not act on multiple rows");
		}
		if (!item.disabled) {
			return fail("missing 'disabled' property");
		}
		expect(item.disabled(testUser)).toBeTrue();
	});
});
