/*
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
import { ComponentFixture, fakeAsync, TestBed, tick } from "@angular/core/testing";
import { MatDialog, MatDialogModule, MatDialogRef } from "@angular/material/dialog";
import { ActivatedRoute } from "@angular/router";
import { RouterTestingModule } from "@angular/router/testing";
import type { ValueFormatterParams } from "ag-grid-community";
import { of } from "rxjs";

import { PhysicalLocationService } from "src/app/api";
import { APITestingModule } from "src/app/api/testing";
import { PhysLocTableComponent } from "src/app/core/servers/phys-loc/table/phys-loc-table.component";
import { isAction } from "src/app/shared/generic-table/generic-table.component";

const testPhysLoc = {
	address: "address",
	city: "city",
	comments: null,
	email: null,
	id: 9001,
	lastUpdated: new Date(),
	name: "test",
	phone: null,
	poc: null,
	region: "region",
	regionId: 1,
	shortName: "test",
	state: "state",
	zip: "zip"
};

describe("PhysLocTableComponent", () => {
	let component: PhysLocTableComponent;
	let fixture: ComponentFixture<PhysLocTableComponent>;

	beforeEach(async () => {
		await TestBed.configureTestingModule({
			declarations: [ PhysLocTableComponent ],
			imports: [ APITestingModule, RouterTestingModule, MatDialogModule ]
		})
			.compileComponents();

		fixture = TestBed.createComponent(PhysLocTableComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});

	it("sets the fuzzy search subject based on the search query param", fakeAsync(() => {
		const router = TestBed.inject(ActivatedRoute);
		const searchString = "testquest";
		spyOnProperty(router, "queryParamMap").and.returnValue(of(new Map([["search", searchString]])));

		let searchValue = "not the right string";
		component.fuzzySubject.subscribe(
			s => searchValue = s
		);

		component.ngOnInit();
		tick();

		expect(searchValue).toBe(searchString);
	}));

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

	it("handles unrecognized contextmenu events", async (): Promise<void> => {
		expect(async () => component.handleContextMenu({
			action: component.contextMenuItems[0].name,
			data: (await component.physLocations)[0]
		})).not.toThrow();
	});

	it("handles the 'delete' context menu item", fakeAsync(async () => {
		const item = component.contextMenuItems.find(c => c.name === "Delete");
		if (!item) {
			return fail("missing 'Delete' context menu item");
		}
		if (!isAction(item)) {
			return fail("expected an action, not a link");
		}
		expect(item.multiRow).toBeFalsy();
		expect(item.disabled).toBeUndefined();

		const api = TestBed.inject(PhysicalLocationService);
		const spy = spyOn(api, "deletePhysicalLocation").and.callThrough();
		expect(spy).not.toHaveBeenCalled();

		const dialogService = TestBed.inject(MatDialog);
		const openSpy = spyOn(dialogService, "open").and.returnValue({
			afterClosed: () => of(true)
		} as MatDialogRef<unknown>);

		const physLoc = await api.createPhysicalLocation(testPhysLoc);
		expect(openSpy).not.toHaveBeenCalled();
		const asyncExpectation = expectAsync(component.handleContextMenu({action: "delete", data: physLoc})).toBeResolvedTo(undefined);
		tick();

		expect(openSpy).toHaveBeenCalled();
		tick();

		expect(spy).toHaveBeenCalled();

		await asyncExpectation;
	}));

	it("generates 'Edit' context menu item href", () => {
		const item = component.contextMenuItems.find(i => i.name === "Edit");
		if (!item) {
			return fail("missing 'Edit' context menu item");
		}
		if (isAction(item)) {
			return fail("expected a link, not an action");
		}
		if (typeof(item.href) !== "function") {
			return fail(`'Edit' context menu item should use a function to determine href, instead uses: ${item.href}`);
		}
		expect(item.href(testPhysLoc)).toBe(String(testPhysLoc.id));
		expect(item.queryParams).toBeUndefined();
		expect(item.fragment).toBeUndefined();
		expect(item.newTab).toBeFalsy();
	});

	it("generates 'Open in New Tab' context menu item href", () => {
		const item = component.contextMenuItems.find(i => i.name === "Open in New Tab");
		if (!item) {
			return fail("missing 'Open in New Tab' context menu item");
		}
		if (isAction(item)) {
			return fail("expected a link, not an action");
		}
		if (typeof(item.href) !== "function") {
			return fail(`'Open in New Tab' context menu item should use a function to determine href, instead uses: ${item.href}`);
		}
		expect(item.href(testPhysLoc)).toBe(String(testPhysLoc.id));
		expect(item.queryParams).toBeUndefined();
		expect(item.fragment).toBeUndefined();
		expect(item.newTab).toBeTrue();
	});

	it("generates 'View Servers' context menu item href", () => {
		const item = component.contextMenuItems.find(i => i.name === "View Servers");
		if (!item) {
			return fail("missing 'View Servers' context menu item");
		}
		if (isAction(item)) {
			return fail("expected a link, not an action");
		}
		if (!item.href) {
			return fail("missing 'href' property");
		}
		if (typeof(item.href) !== "string") {
			return fail("'View Servers' context menu item should use a static string to determine href, instead uses a function");
		}
		expect(item.href).toBe("/core/servers");
		if (typeof(item.queryParams) !== "function") {
			return fail(
				`'View Servers' context menu item should use a function to determine query params, instead uses: ${item.queryParams}`
			);
		}
		expect(item.queryParams(testPhysLoc)).toEqual({physLocation: testPhysLoc.name});
		expect(item.fragment).toBeUndefined();
		expect(item.newTab).toBeFalsy();
	});

	it("generates 'View Region' context menu item href", () => {
		const item = component.contextMenuItems.find(i => i.name === "View Region");
		if (!item) {
			return fail("missing 'View Region' context menu item");
		}
		if (isAction(item)) {
			return fail("expected a link, not an action");
		}
		if (!item.href) {
			return fail("missing 'href' property");
		}
		if (typeof(item.href) !== "function") {
			return fail(`'View Region' context menu item should use a function to determine href, instead uses: ${item.href}`);
		}
		expect(item.href(testPhysLoc)).toBe(`/core/regions/${testPhysLoc.regionId}`);
		expect(item.queryParams).toBeUndefined();
		expect(item.fragment).toBeUndefined();
		expect(item.newTab).toBeFalsy();
	});

	it("formats region cell data", () => {
		const item = component.columnDefs.find(i => i.field === "region");
		if (!item) {
			return fail("missing 'region' column");
		}
		if (!item.valueFormatter) {
			return fail("'region' column should have a valueFormatter");
		}
		expect(item.valueFormatter({data: testPhysLoc} as ValueFormatterParams)).toBe(`${testPhysLoc.region} (#${testPhysLoc.regionId})`);
	});
});
