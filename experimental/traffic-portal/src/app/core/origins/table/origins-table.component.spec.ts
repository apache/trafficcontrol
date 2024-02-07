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

import {
	ComponentFixture,
	fakeAsync,
	TestBed,
	tick,
} from "@angular/core/testing";
import {
	MatDialog,
	MatDialogModule,
	MatDialogRef,
} from "@angular/material/dialog";
import { ActivatedRoute } from "@angular/router";
import { RouterTestingModule } from "@angular/router/testing";
import { of } from "rxjs";
import type { ResponseOrigin } from "trafficops-types";

import { OriginService } from "src/app/api";
import { APITestingModule } from "src/app/api/testing";
import { isAction } from "src/app/shared/generic-table/generic-table.component";

import { OriginsTableComponent } from "./origins-table.component";

const testOrigin: ResponseOrigin = {
	cachegroup: "",
	cachegroupId: 1,
	coordinate: "",
	coordinateId: 1,
	deliveryService: "",
	deliveryServiceId: 1,
	fqdn: "0",
	id: 1,
	ip6Address: "",
	ipAddress: "",
	isPrimary: false,
	lastUpdated: new Date(),
	name: "TestOrigin",
	port: 80,
	profile: "",
	profileId: 1,
	protocol: "https",
	tenant: "*",
	tenantId: 0,
};

describe("OriginsTableComponent", () => {
	let component: OriginsTableComponent;
	let fixture: ComponentFixture<OriginsTableComponent>;

	beforeEach(async () => {
		await TestBed.configureTestingModule({
			declarations: [OriginsTableComponent],
			imports: [APITestingModule, RouterTestingModule, MatDialogModule],
		}).compileComponents();

		fixture = TestBed.createComponent(OriginsTableComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});

	it("sets the fuzzy search subject based on the search query param", fakeAsync(() => {
		const router = TestBed.inject(ActivatedRoute);
		const searchString = "testorigin";
		spyOnProperty(router, "queryParamMap").and.returnValue(
			of(new Map([["search", searchString]]))
		);

		let searchValue = "not the right string";
		component.fuzzySubject.subscribe((s) => (searchValue = s));

		component.ngOnInit();
		tick();

		expect(searchValue).toBe(searchString);
	}));

	it("updates the fuzzy search output", fakeAsync(() => {
		let called = false;
		const text = "testorigin";
		const spy = jasmine.createSpy("subscriber", (txt: string): void => {
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

	it("handle the 'delete' context menu item", fakeAsync(async () => {
		const item = component.contextMenuItems.find(
			(c) => c.name === "Delete"
		);
		if (!item) {
			return fail("missing 'Delete' context menu item");
		}
		if (!isAction(item)) {
			return fail("expected an action, not a link");
		}
		expect(item.multiRow).toBeFalsy();
		expect(item.disabled).toBeUndefined();

		const api = TestBed.inject(OriginService);
		const spy = spyOn(api, "deleteOrigin").and.callThrough();
		expect(spy).not.toHaveBeenCalled();

		const dialogService = TestBed.inject(MatDialog);
		const openSpy = spyOn(dialogService, "open").and.returnValue({
			afterClosed: () => of(true),
		} as MatDialogRef<unknown>);

		const origin = await api.createOrigin({
			deliveryServiceId: 1,
			fqdn: "0",
			name: "*",
			protocol: "https",
			tenantID: 1,
		});
		expect(openSpy).not.toHaveBeenCalled();
		const asyncExpectation = expectAsync(
			component.handleContextMenu({
				action: "delete",
				data: origin,
			})
		).toBeResolvedTo(undefined);
		tick();

		expect(openSpy).toHaveBeenCalled();
		tick();

		expect(spy).toHaveBeenCalled();

		await asyncExpectation;
	}));

	it("generates 'Edit' context menu item href", () => {
		const item = component.contextMenuItems.find((i) => i.name === "Edit");
		if (!item) {
			return fail("missing 'Edit' context menu item");
		}
		if (isAction(item)) {
			return fail("expected a link, not an action");
		}
		if (typeof item.href !== "function") {
			return fail(
				`'Edit' context menu item should use a function to determine href, instead uses: ${item.href}`
			);
		}
		expect(item.href(testOrigin)).toBe(String(testOrigin.id));
		expect(item.queryParams).toBeUndefined();
		expect(item.fragment).toBeUndefined();
		expect(item.newTab).toBeFalsy();
	});

	it("generates 'Open in New Tab' context menu item href", () => {
		const item = component.contextMenuItems.find(
			(i) => i.name === "Open in New Tab"
		);
		if (!item) {
			return fail("missing 'Open in New Tab' context menu item");
		}
		if (isAction(item)) {
			return fail("expected a link, not an action");
		}
		if (typeof item.href !== "function") {
			return fail(
				`'Open in New Tab' context menu item should use a function to determine href, instead uses: ${item.href}`
			);
		}
		expect(item.href(testOrigin)).toBe(String(testOrigin.id));
		expect(item.queryParams).toBeUndefined();
		expect(item.fragment).toBeUndefined();
		expect(item.newTab).toBeTrue();
	});
});
