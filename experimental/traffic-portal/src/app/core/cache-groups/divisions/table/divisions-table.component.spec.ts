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
import { MatDialog, MatDialogModule, MatDialogRef } from "@angular/material/dialog";
import { ActivatedRoute } from "@angular/router";
import { RouterTestingModule } from "@angular/router/testing";
import { of } from "rxjs";

import { CacheGroupService } from "src/app/api";
import { APITestingModule } from "src/app/api/testing";
import { DivisionsTableComponent } from "src/app/core/cache-groups/divisions/table/divisions-table.component";
import { isAction } from "src/app/shared/generic-table/generic-table.component";

const testDivision = {
	id: 1,
	lastUpdated: new Date(),
	name: "TestQuest",
};

describe("DivisionsTableComponent", () => {
	let component: DivisionsTableComponent;
	let fixture: ComponentFixture<DivisionsTableComponent>;

	beforeEach(async () => {
		await TestBed.configureTestingModule({
			declarations: [ DivisionsTableComponent ],
			imports: [
				APITestingModule,
				RouterTestingModule,
				MatDialogModule
			]
		}).compileComponents();

		fixture = TestBed.createComponent(DivisionsTableComponent);
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

	it("handles unrecognized contextmenu events", () => {
		expect(async () => component.handleContextMenu({
			action: component.contextMenuItems[0].name,
			data: {id: 1, lastUpdated: new Date(), name: "Div"}
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

		const api = TestBed.inject(CacheGroupService);
		const spy = spyOn(api, "deleteDivision").and.callThrough();
		expect(spy).not.toHaveBeenCalled();

		const dialogService = TestBed.inject(MatDialog);
		const openSpy = spyOn(dialogService, "open").and.returnValue({
			afterClosed: () => of(true)
		} as MatDialogRef<unknown>);

		const div = await api.createDivision({name: "test"});
		expect(openSpy).not.toHaveBeenCalled();
		const asyncExpectation = expectAsync(component.handleContextMenu({action: "delete", data: div})).toBeResolvedTo(undefined);
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
		expect(item.href(testDivision)).toBe(String(testDivision.id));
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
		expect(item.href(testDivision)).toBe(String(testDivision.id));
		expect(item.queryParams).toBeUndefined();
		expect(item.fragment).toBeUndefined();
		expect(item.newTab).toBeTrue();
	});

	it("generates 'View Regions' context menu item href", () => {
		const item = component.contextMenuItems.find(i => i.name === "View Regions");
		if (!item) {
			return fail("missing 'View Regions' context menu item");
		}
		if (isAction(item)) {
			return fail("expected a link, not an action");
		}
		if (!item.href) {
			return fail("missing 'href' property");
		}
		if (typeof(item.href) !== "string") {
			return fail("'View Regions' context menu item should use a static string to determine href, instead uses a function");
		}
		expect(item.href).toBe("/core/regions");
		if (typeof(item.queryParams) !== "function") {
			return fail(
				`'View Regions' context menu item should use a function to determine query params, instead uses: ${item.queryParams}`
			);
		}
		expect(item.queryParams(testDivision)).toEqual({divisionName: testDivision.name});
		expect(item.fragment).toBeUndefined();
		expect(item.newTab).toBeFalsy();
	});
});
