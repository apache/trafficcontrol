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
import { HarnessLoader } from "@angular/cdk/testing";
import { TestbedHarnessEnvironment } from "@angular/cdk/testing/testbed";
import { HttpClientModule } from "@angular/common/http";
import { type ComponentFixture, TestBed, fakeAsync, tick } from "@angular/core/testing";
import { ReactiveFormsModule } from "@angular/forms";
import { MatButtonHarness } from "@angular/material/button/testing";
import { MatDialogModule } from "@angular/material/dialog";
import { MatDialogHarness } from "@angular/material/dialog/testing";
import { MatSelectModule } from "@angular/material/select";
import { MatSelectHarness } from "@angular/material/select/testing";
import { NoopAnimationsModule } from "@angular/platform-browser/animations";
import { Router } from "@angular/router";
import { RouterTestingModule } from "@angular/router/testing";
import type { ValueFormatterParams, ValueGetterParams } from "ag-grid-community";
import { ReplaySubject } from "rxjs";
import { AlertLevel, LocalizationMethod, localizationMethodToString, type ResponseCacheGroup } from "trafficops-types";

import { CacheGroupService, CDNService } from "src/app/api";
import { APITestingModule } from "src/app/api/testing";
import { AlertService } from "src/app/shared/alert/alert.service";
import { isAction } from "src/app/shared/generic-table/generic-table.component";
import { NavigationService } from "src/app/shared/navigation/navigation.service";

import { CacheGroupTableComponent } from "./cache-group-table.component";

const sampleCG: ResponseCacheGroup = {
	fallbackToClosest: true,
	fallbacks: [],
	id: 1,
	lastUpdated: new Date(),
	latitude: 0,
	localizationMethods: [],
	longitude: 0,
	name: "sample",
	parentCachegroupId: null,
	parentCachegroupName: null,
	secondaryParentCachegroupId: null,
	secondaryParentCachegroupName: null,
	shortName: "sample",
	typeId: 1,
	typeName: "some type"
};

describe("CacheGroupTableComponent", () => {
	let component: CacheGroupTableComponent;
	let fixture: ComponentFixture<CacheGroupTableComponent>;
	let loader: HarnessLoader;

	const navSvc = jasmine.createSpyObj([],{headerHidden: new ReplaySubject<boolean>(), headerTitle: new ReplaySubject<string>()});
	beforeEach(async () => {
		await TestBed.configureTestingModule({
			declarations: [ CacheGroupTableComponent ],
			imports: [
				APITestingModule,
				HttpClientModule,
				ReactiveFormsModule,
				RouterTestingModule,
				MatDialogModule,
				NoopAnimationsModule,
				MatSelectModule
			],
			providers: [
				{ provide: NavigationService, useValue: navSvc}
			]
		}).compileComponents();
		fixture = TestBed.createComponent(CacheGroupTableComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
		loader = TestbedHarnessEnvironment.documentRootLoader(fixture);
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});

	it("emits the search box value", fakeAsync(() => {
		component.fuzzControl.setValue("query");
		component.updateURL();
		expectAsync(component.fuzzySubject.toPromise()).toBeResolvedTo("query");
	}));

	it("doesn't throw errors when handling context menu events", () => {
		expect(()=>component.handleContextMenu({action: "something", data: []})).not.toThrow();
	});

	it("renders parent cache group cells", () => {
		const col = component.columnDefs.find(d => d.field === "parentCachegroupName");
		if (!col) {
			return fail("parentCachegroupName column not found");
		}
		const {valueFormatter} = col;
		if (typeof(valueFormatter) !== "function") {
			return fail(`invalid valueFormatter found on parentCachegroupName column definition: ${valueFormatter}`);
		}
		let value = valueFormatter({data: sampleCG} as ValueFormatterParams);
		expect(value).toBe("");
		value = valueFormatter({data: {...sampleCG, parentCachegroupId: 1, parentCachegroupName: "sample"}} as ValueFormatterParams);
		expect(value).toBe("sample (#1)");
	});
	it("renders secondary parent cache group cells", () => {
		const col = component.columnDefs.find(d => d.field === "secondaryParentCachegroupName");
		if (!col) {
			return fail("secondaryParentCachegroupName column not found");
		}
		const {valueFormatter} = col;
		if (typeof(valueFormatter) !== "function") {
			return fail(`invalid valueFormatter found on secondaryParentCachegroupName column definition: ${valueFormatter}`);
		}
		let value = valueFormatter({data: sampleCG} as ValueFormatterParams);
		expect(value).toBe("");
		value = valueFormatter({
			data: {
				...sampleCG,
				secondaryParentCachegroupId: 1,
				secondaryParentCachegroupName: "sample"
			}
		} as ValueFormatterParams);
		expect(value).toBe("sample (#1)");
	});
	it("renders type cells", () => {
		const col = component.columnDefs.find(d => d.field === "typeName");
		if (!col) {
			return fail("type column not found");
		}
		const {valueFormatter} = col;
		if (typeof(valueFormatter) !== "function") {
			return fail(`invalid valueFormatter found on type column definition: ${valueFormatter}`);
		}
		const value = valueFormatter({data: {...sampleCG, typeId: 1, typeName: "sample"}} as ValueFormatterParams);
		expect(value).toBe("sample (#1)");
	});
	it("renders localization methods cells", () => {
		const col = component.columnDefs.find(d => d.field === "localizationMethods");
		if (!col) {
			return fail("localizationMethods column not found");
		}
		const {valueGetter} = col;
		if (typeof(valueGetter) !== "function") {
			return fail(`invalid valueGetter found on localizationMethods column definition: ${valueGetter}`);
		}
		let value: string = valueGetter({data: {...sampleCG, localizationMethods: []}} as ValueGetterParams);
		let valueArr = value.split(", ");
		expect(valueArr.length).toBe(3);
		expect(valueArr).toContain(localizationMethodToString(LocalizationMethod.CZ));
		expect(valueArr).toContain(localizationMethodToString(LocalizationMethod.DEEP_CZ));
		expect(valueArr).toContain(localizationMethodToString(LocalizationMethod.GEO));
		value = valueGetter({data: {...sampleCG, localizationMethods: [LocalizationMethod.CZ]}} as ValueGetterParams);
		valueArr = value.split(", ");
		expect(valueArr.length).toBe(1);
		expect(valueArr).toContain(localizationMethodToString(LocalizationMethod.CZ));
	});

	it("has context menu links to individual Cache Groups", () => {
		let menuItem = component.contextMenuItems.find(i => i.name === "Open in New Tab");
		if (!menuItem) {
			return fail("'Open in New Tab' context menu item not found");
		}
		if (isAction(menuItem) || typeof(menuItem.href) !== "function") {
			return fail(`invalid 'Open in New Tab' context menu item; either not a link or has a static href: ${menuItem}`);
		}
		expect(menuItem.newTab).toBeTrue();
		expect(menuItem.href({...sampleCG, id: 5})).toBe("5");

		menuItem = component.contextMenuItems.find(i => i.name === "Edit");
		if (!menuItem) {
			return fail("'Edit' context menu item not found");
		}
		if (isAction(menuItem) || typeof(menuItem.href) !== "function") {
			return fail(`invalid 'Edit' context menu item; either not a link or has a static href: ${menuItem}`);
		}
		expect(menuItem.newTab).toBeFalsy();
		expect(menuItem.href({...sampleCG, id: 5})).toBe("5");
	});

	it("generates 'View ASNs' context menu item href", () => {
		const item = component.contextMenuItems.find(i => i.name === "View ASNs");
		if (!item) {
			return fail("missing 'View ASNs' context menu item");
		}
		if (isAction(item)) {
			return fail("expected an action, not a link");
		}
		if (!item.href) {
			return fail("missing 'href' property");
		}
		if (typeof(item.href) !== "string") {
			return fail("'View ASNs' context menu item should use a static string to determine href, instead uses a function");
		}
		expect(item.href).toBe("/core/asns");
		if (typeof(item.queryParams) !== "function") {
			return fail(
				`'Manage ASNs' context menu item should use a function to determine query params, instead uses: ${item.queryParams}`
			);
		}
		expect(item.queryParams(sampleCG)).toEqual({cachegroup: sampleCG.name});
		expect(item.fragment).toBeUndefined();
		expect(item.newTab).toBeFalsy();
	});

	it("builds links to the servers in a Cache Group", () => {
		const menuItem = component.contextMenuItems.find(i => i.name === "View Servers");
		if (!menuItem) {
			return fail("'View Servers' context menu item not found");
		}
		if (isAction(menuItem)) {
			return fail("Invalid 'View Servers' context menu item; not a link");
		}

		expect(menuItem.href).toBe("/core/servers");
		expect(menuItem.newTab).toBeFalsy();
		expect(menuItem.fragment).not.toBeDefined();
		expect(menuItem.queryParams).toBeDefined();
		if (typeof(menuItem.queryParams) !== "function") {
			return fail("invalid 'View Servers' context menu item; query params not a function");
		}
		expect(menuItem.queryParams(sampleCG)).toEqual({cachegroup: sampleCG.name});
	});

	it("initializes from query string parameters", fakeAsync(() => {
		const router = TestBed.inject(Router);
		router.navigate([], {queryParams: {search: "testquest"}});
		component.ngOnInit();
		tick();
		expectAsync(component.fuzzySubject.toPromise()).toBeResolvedTo("testquest");
	}));

	it("deletes Cache Groups", async () => {
		expect(() => component.handleContextMenu({action: "delete", data: []})).not.toThrow();
		let dialogs = await loader.getAllHarnesses(MatDialogHarness);
		expect(dialogs.length).toBe(0);

		const cgSrv = TestBed.inject(CacheGroupService);
		const spy = spyOn(cgSrv, "deleteCacheGroup");

		component.handleContextMenu({action: "delete", data: sampleCG});
		dialogs = await loader.getAllHarnesses(MatDialogHarness);
		if (dialogs.length !== 1) {
			return fail(`dialog should have opened for deleting, actual number of dialogs: ${dialogs.length}`);
		}
		let dialog = dialogs[0];
		await dialog.close();
		expect(spy).not.toHaveBeenCalled();

		component.handleContextMenu({action: "delete", data: sampleCG});
		dialogs = await loader.getAllHarnesses(MatDialogHarness);
		if (dialogs.length !== 1) {
			return fail(`dialog should have opened for deleting, actual number of dialogs: ${dialogs.length}`);
		}
		dialog = dialogs[0];
		const buttons = await dialog.getAllHarnesses(MatButtonHarness.with({text: /^[cC][oO][nN][fF][iI][rR][mM]$/}));
		if (buttons.length !== 1) {
			return fail(`'Confirm' button not found; expected one, got: ${buttons.length}`);
		}
		const button = buttons[0];
		await button.click();
		expect(spy).toHaveBeenCalled();
	});
	it("queues Cache Group updates", async () => {
		const cgSrv = TestBed.inject(CacheGroupService);
		const spy = spyOn(cgSrv, "queueCacheGroupUpdates").and.returnValue(
			new Promise(r => r({
				action: "queue",
				cachegroupID: 1,
				cachegroupName: "testquest",
				cdn: "doesn't matter",
				serverNames: ["testquest"],
			}))
		);
		const alertSrv = TestBed.inject(AlertService);
		const alertSpy = spyOn(alertSrv, "newAlert");

		expect(() => component.handleContextMenu({action: "queue", data: [sampleCG]})).not.toThrow();
		let dialogs = await loader.getAllHarnesses(MatDialogHarness);
		expect(dialogs.length).toBe(1);
		let dialog = dialogs[0];
		await dialog.close();

		expect(spy).not.toHaveBeenCalled();
		expect(alertSpy).not.toHaveBeenCalled();

		component.handleContextMenu({action: "queue", data: sampleCG});
		dialogs = await loader.getAllHarnesses(MatDialogHarness);
		if (dialogs.length !== 1) {
			return fail(`dialog should have opened for queuing, actual number of dialogs: ${dialogs.length}`);
		}
		dialog = dialogs[0];
		const selects = await dialog.getAllHarnesses(MatSelectHarness);
		if (selects.length !== 1) {
			return fail(`dialog should have contained one select input, got: ${selects.length}`);
		}
		const select = selects[0];
		const cdnSrv = TestBed.inject(CDNService);
		expect(await cdnSrv.getCDNs()).toHaveSize(2);
		await select.clickOptions();
		const buttons = await dialog.getAllHarnesses(MatButtonHarness.with({text: /^[cC][oO][nN][fF][iI][rR][mM]$/}));
		if (buttons.length !== 1) {
			return fail(`'Confirm' button not found; expected one, got: ${buttons.length}`);
		}
		const button = buttons[0];
		await button.click();
		expect(spy).toHaveBeenCalled();
		// Jasmine has trouble with the overload signatures; it thinks you can
		// only call `newAlerts` with a single `Alert` argument. Otherwise, the
		// below lines could simply be:
		// expect(alertSpy).toHaveBeenCalledOnceWith(AlertLevel.SUCCESS, "Queued Updates on 1 server");
		expect(alertSpy).toHaveBeenCalledTimes(1);
		const [level, text] = alertSpy.calls.all()[0].args as unknown as [string, string];
		expect(level).toBe(AlertLevel.SUCCESS);
		expect(text).toBe("Queued Updates on 1 server");
	});
	it("clears Cache Group updates", async () => {
		const cgSrv = TestBed.inject(CacheGroupService);
		const spy = spyOn(cgSrv, "queueCacheGroupUpdates").and.returnValue(
			new Promise(r => r({
				action: "dequeue",
				cachegroupID: 1,
				cachegroupName: "testquest",
				cdn: "doesn't matter",
				serverNames: ["testquest"],
			}))
		);
		const alertSrv = TestBed.inject(AlertService);
		const alertSpy = spyOn(alertSrv, "newAlert");

		expect(() => component.handleContextMenu({action: "dequeue", data: sampleCG})).not.toThrow();
		let dialogs = await loader.getAllHarnesses(MatDialogHarness);
		expect(dialogs.length).toBe(1);
		let dialog = dialogs[0];
		await dialog.close();

		expect(spy).not.toHaveBeenCalled();
		expect(alertSpy).not.toHaveBeenCalled();

		component.handleContextMenu({action: "dequeue", data: [sampleCG, sampleCG]});
		dialogs = await loader.getAllHarnesses(MatDialogHarness);
		if (dialogs.length !== 1) {
			return fail(`dialog should have opened for dequeuing, actual number of dialogs: ${dialogs.length}`);
		}
		dialog = dialogs[0];
		const selects = await dialog.getAllHarnesses(MatSelectHarness);
		if (selects.length !== 1) {
			return fail(`dialog should have contained one select input, got: ${selects.length}`);
		}
		const select = selects[0];
		const cdnSrv = TestBed.inject(CDNService);
		expect(await cdnSrv.getCDNs()).toHaveSize(2);
		await select.clickOptions();
		const buttons = await dialog.getAllHarnesses(MatButtonHarness.with({text: /^[cC][oO][nN][fF][iI][rR][mM]$/}));
		if (buttons.length !== 1) {
			return fail(`'Confirm' button not found; expected one, got: ${buttons.length}`);
		}
		const button = buttons[0];
		await button.click();
		expect(spy).toHaveBeenCalled();
		// Jasmine has trouble with the overload signatures; it thinks you can
		// only call `newAlerts` with a single `Alert` argument. Otherwise, the
		// below lines could simply be:
		// expect(alertSpy).toHaveBeenCalledOnceWith(AlertLevel.SUCCESS, "Queued Updates on 1 server");
		expect(alertSpy).toHaveBeenCalledTimes(1);
		const [level, text] = alertSpy.calls.all()[0].args as unknown as [string, string];
		expect(level).toBe(AlertLevel.SUCCESS);
		expect(text).toBe("Cleared Updates on 2 servers");
	});
});
