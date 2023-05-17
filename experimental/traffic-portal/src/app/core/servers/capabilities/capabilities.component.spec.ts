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
import { ComponentFixture, TestBed, fakeAsync, tick } from "@angular/core/testing";
import { MatDialog, MatDialogModule, type MatDialogRef } from "@angular/material/dialog";
import { ActivatedRoute } from "@angular/router";
import { RouterTestingModule } from "@angular/router/testing";
import { of } from "rxjs";

import { ServerService } from "src/app/api";
import { APITestingModule } from "src/app/api/testing";
import { isAction } from "src/app/shared/generic-table/generic-table.component";

import { CapabilitiesComponent } from "./capabilities.component";

describe("CapabilitiesComponent", () => {
	let component: CapabilitiesComponent;
	let fixture: ComponentFixture<CapabilitiesComponent>;

	const capability = {
		lastUpdated: new Date(),
		name: "test"
	};

	beforeEach(async () => {
		await TestBed.configureTestingModule({
			declarations: [ CapabilitiesComponent ],
			imports: [ APITestingModule, RouterTestingModule, MatDialogModule ],
		}).compileComponents();

		fixture = TestBed.createComponent(CapabilitiesComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});

	it("builds a 'View Details' link", () => {
		const item = component.contextMenuItems.find(i => i.name === "View Details");
		if (!item) {
			return fail("missing 'View Details' context menu item");
		}

		if (isAction(item)) {
			return fail("incorrect type for 'View Details' menu item");
		}

		if (typeof(item.href) !== "function") {
			return fail("link should be built from data, not static");
		}

		expect(item.href(capability)).toBe(capability.name);
	});

	it("builds an 'Open in New Tab' link", () => {
		const item = component.contextMenuItems.find(i => i.name === "Open in New Tab");
		if (!item) {
			return fail("missing 'Open in New Tab' context menu item");
		}

		if (isAction(item)) {
			return fail("incorrect type for 'Open in New Tab' menu item");
		}

		expect(item.newTab).toBe(true);

		if (typeof(item.href) !== "function") {
			return fail("link should be built from data, not static");
		}

		expect(item.href(capability)).toBe(capability.name);
	});

	it("has context menu items that aren't implemented yet", () => {
		let item = component.contextMenuItems.find(i => i.name === "View Servers");
		if (!item) {
			return fail("missing 'View Servers' context menu item");
		}
		if (!isAction(item)) {
			return fail("incorrect type for 'View Servers' menu item");
		}
		if (!item.multiRow) {
			return fail("'View Servers' should be a multi-row action");
		}
		if (!item.disabled || !item.disabled([capability])) {
			return fail("'View Servers' should be disabled");
		}

		item = component.contextMenuItems.find(i => i.name === "Add to Server(s)");
		if (!item) {
			return fail("missing 'Add to Server(s)' context menu item");
		}
		if (!isAction(item)) {
			return fail("incorrect type for 'Add to Server(s)' menu item");
		}
		if (!item.multiRow) {
			return fail("'Add to Server(s)' should be a multi-row action");
		}
		if (!item.disabled || !item.disabled([capability])) {
			return fail("'Add to Server(s)' should be disabled");
		}

		item = component.contextMenuItems.find(i => i.name === "View Delivery Services");
		if (!item) {
			return fail("missing 'View Delivery Services' context menu item");
		}
		if (!isAction(item)) {
			return fail("incorrect type for 'View Delivery Services' menu item");
		}
		if (!item.multiRow) {
			return fail("'View Delivery Services' should be a multi-row action");
		}
		if (!item.disabled || !item.disabled([capability])) {
			return fail("'View Delivery Services' should be disabled");
		}

		item = component.contextMenuItems.find(i => i.name === "Add to Delivery Service(s)");
		if (!item) {
			return fail("missing 'Add to Delivery Service(s)' context menu item");
		}
		if (!isAction(item)) {
			return fail("incorrect type for 'Add to Delivery Service(s)' menu item");
		}
		if (!item.multiRow) {
			return fail("'Add to Delivery Service(s)' should be a multi-row action");
		}
		if (!item.disabled || !item.disabled([capability])) {
			return fail("'Add to Delivery Service(s)' should be disabled");
		}
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
			data: (await component.capabilities)[0]
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

		const api = TestBed.inject(ServerService);
		const spy = spyOn(api, "deleteCapability").and.callThrough();
		expect(spy).not.toHaveBeenCalled();

		const dialogService = TestBed.inject(MatDialog);
		const openSpy = spyOn(dialogService, "open").and.returnValue({
			afterClosed: () => of(true)
		} as MatDialogRef<unknown>);

		const cap = await api.createCapability(capability);
		expect(openSpy).not.toHaveBeenCalled();
		const asyncExpectation = expectAsync(component.handleContextMenu({action: "delete", data: cap})).toBeResolvedTo(undefined);
		tick();

		expect(openSpy).toHaveBeenCalled();
		tick();

		expect(spy).toHaveBeenCalled();

		await asyncExpectation;
	}));

	it("throws an error if improperly asked to delete more than one Capability", async () => {
		await expectAsync(component.handleContextMenu({action: "delete", data: []})).toBeRejected();
	});
});
