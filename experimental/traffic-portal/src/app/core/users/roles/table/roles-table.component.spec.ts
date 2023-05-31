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
import { MatDialog, MatDialogModule, type MatDialogRef } from "@angular/material/dialog";
import { RouterTestingModule } from "@angular/router/testing";
import { of } from "rxjs";
import { ResponseRole } from "trafficops-types";

import { UserService } from "src/app/api";
import { APITestingModule } from "src/app/api/testing";
import { RolesTableComponent } from "src/app/core/users/roles/table/roles-table.component";
import { isAction } from "src/app/shared/generic-table/generic-table.component";

describe("RolesTableComponent", () => {
	let component: RolesTableComponent;
	let fixture: ComponentFixture<RolesTableComponent>;

	const role: ResponseRole = {
		description: "Test Role",
		lastUpdated: new Date(),
		name: "test"
	};

	beforeEach(async () => {
		await TestBed.configureTestingModule({
			declarations: [ RolesTableComponent ],
			imports: [ APITestingModule, RouterTestingModule, MatDialogModule ]
		}).compileComponents();

		fixture = TestBed.createComponent(RolesTableComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it("should create", () => {
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
		component.fuzzControl.setValue(text);
		component.updateURL();
		tick();
		expect(spy).toHaveBeenCalledTimes(2);
	}));

	it("handles contextmenu events", () => {
		expect(async () => component.handleContextMenu({
			action: component.contextMenuItems[0].name,
			data: {description: "Can only read", lastUpdated: new Date(), name: "test"}
		})).not.toThrow();
	});

	it("builds an 'Open in New Tab' link", () => {
		const item = component.contextMenuItems.find(i => i.name === "Open in New Tab");
		if (!item) {
			return fail("missing 'Open in New Tab' context menu item");
		}

		if (isAction(item)) {
			return fail("incorrect type for 'Open in New Tab' menu item. Expected an action, not a link");
		}

		expect(item.newTab).toBe(true);

		if (typeof(item.href) !== "function") {
			return fail("link should be built from data, not static");
		}

		expect(item.href(role)).toBe(role.name);
	});

	it("deletes Roles", fakeAsync(async () => {
		const item = component.contextMenuItems.find(i => i.name === "Delete");
		if (!item) {
			return fail("missing 'Delete' context menu item");
		}
		if (!isAction(item)) {
			return fail("incorrect type for 'Delete' menu item. Expected an action, not a link");
		}
		expect(item.multiRow).toBeFalsy();
		expect(item.disabled).toBeUndefined();

		const api = TestBed.inject(UserService);
		const spy = spyOn(api, "deleteRole").and.callThrough();
		expect(spy).not.toHaveBeenCalled();

		const dialogService = TestBed.inject(MatDialog);
		const openSpy = spyOn(dialogService, "open").and.returnValue({
			afterClosed: () => of(true)
		} as MatDialogRef<unknown>);

		const testRole = await api.createRole(role);
		expect(openSpy).not.toHaveBeenCalled();
		const asyncExpectation = expectAsync(component.handleContextMenu({action: "delete", data: testRole})).toBeResolvedTo(undefined);
		tick();

		expect(openSpy).toHaveBeenCalled();
		tick();

		expect(spy).toHaveBeenCalled();

		await asyncExpectation;
	}));

	it("generate 'View Users' context menu item href", () => {
		const item = component.contextMenuItems.find(i => i.name === "View Users");
		if (!item) {
			return fail("missing 'View Users' context menu item");
		}
		if (isAction(item)) {
			return fail("incorrect type for 'View Users' menu item. Expected an action, not a link");
		}
		if (!item.href) {
			return fail("missing 'href' property");
		}
		if (typeof(item.href) !== "string") {
			return fail("'View Users' context menu item should use a static string to determine href, instead uses a function");
		}
		expect(item.href).toBe("/core/users");
		if (typeof(item.queryParams) !== "function") {
			return fail(
				`'View Users' context menu item should use a function to determine query params, instead uses: ${item.queryParams}`
			);
		}
		expect(item.queryParams(role)).toEqual({role: role.name});
		expect(item.fragment).toBeUndefined();
		expect(item.newTab).toBeFalsy();
	});
});
