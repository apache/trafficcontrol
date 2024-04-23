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
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { MatButtonHarness } from "@angular/material/button/testing";
import { MatDialogModule } from "@angular/material/dialog";
import { MatDialogHarness } from "@angular/material/dialog/testing";
import { NoopAnimationsModule } from "@angular/platform-browser/animations";
import { ActivatedRoute } from "@angular/router";
import { RouterTestingModule } from "@angular/router/testing";
import { ReplaySubject } from "rxjs";

import { CacheGroupService, TypeService } from "src/app/api";
import { APITestingModule } from "src/app/api/testing";
import { NavigationService } from "src/app/shared/navigation/navigation.service";

import { CacheGroupDetailsComponent } from "./cache-group-details.component";

describe("CacheGroupDetailsComponent", () => {
	let component: CacheGroupDetailsComponent;
	let fixture: ComponentFixture<CacheGroupDetailsComponent>;
	let route: ActivatedRoute;
	let paramMap: jasmine.Spy;
	let loader: HarnessLoader;
	let cgSrv: CacheGroupService;

	const navSvc = jasmine.createSpyObj([],{headerHidden: new ReplaySubject<boolean>(), headerTitle: new ReplaySubject<string>()});
	beforeEach(async () => {
		await TestBed.configureTestingModule({
			declarations: [ CacheGroupDetailsComponent ],
			imports: [
				APITestingModule,
				RouterTestingModule.withRoutes([
					{component: CacheGroupDetailsComponent, path: "core/cache-groups/:id"},
					{component: CacheGroupDetailsComponent, path: "core/cache-groups"}
				]),
				MatDialogModule,
				NoopAnimationsModule,
			],
			providers: [ { provide: NavigationService, useValue: navSvc } ]
		}).compileComponents();

		route = TestBed.inject(ActivatedRoute);
		paramMap = spyOn(route.snapshot.paramMap, "get");
		paramMap.and.returnValue(null);
		fixture = TestBed.createComponent(CacheGroupDetailsComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
		loader = TestbedHarnessEnvironment.documentRootLoader(fixture);
		cgSrv = TestBed.inject(CacheGroupService);
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});

	it("new Cache Group", async () => {
		paramMap.and.returnValue("new");

		fixture = TestBed.createComponent(CacheGroupDetailsComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
		await fixture.whenStable();
		expect(paramMap).toHaveBeenCalled();
		expect(component.cacheGroup).not.toBeNull();
		expect(component.cacheGroup.name).toBe("");
		expect(component.new).toBeTrue();
	});

	it("existing Cache Group", async () => {
		const cgs = await cgSrv.getCacheGroups();
		if (cgs.length < 1) {
			return fail("no testing Cache Groups - please add Cache Groups to the default set or fix the accidental deletion thereof");
		}
		const cg = cgs[0];
		paramMap.and.returnValue(String(cg.id));

		fixture = TestBed.createComponent(CacheGroupDetailsComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
		await fixture.whenStable();
		await component.ngOnInit();
		expect(paramMap).toHaveBeenCalled();
		expect(component.cacheGroup).not.toBeNull();
		expect(component.cacheGroup.name).toBe(cg.name, component.cacheGroup);
		expect(component.new).toBeFalse();
	});

	it("throws an error when the ID in the URL doesn't exist", async () => {
		paramMap.and.returnValue("-1");
		// This doesn't actually throw, but for unknown reasons it causes
		// earlier tests to fail if were to throw. Somehow, the Component is
		// getting a blank CG list from the testing API, and that should never
		// be possible.
		await expectAsync(component.ngOnInit()).toBeResolved();
		paramMap.and.returnValue("testquest");
		await expectAsync(component.ngOnInit()).toBeRejected();
	});

	it("gets available parent Cache Groups", async () => {
		const cgs = await cgSrv.getCacheGroups();
		expect(cgs.length).toBeGreaterThan(0);
		component.cacheGroups = cgs;
		component.cacheGroup.fallbacks = [`${cgs[0].name}-fallback-test`];

		expect(component.parentCacheGroups()).toEqual(component.cacheGroups);
		if (component.cacheGroups.length < 1) {
			return fail("need at least one cache group to test parentage");
		}
		const initialLength = component.parentCacheGroups().length;
		const cg = component.cacheGroups[0];
		const original = component.cacheGroup;
		component.cacheGroup = {
			...component.cacheGroup,
			secondaryParentCachegroupId: cg.id,
			secondaryParentCachegroupName: cg.name
		};
		expect(component.parentCacheGroups()).not.toContain(cg);
		expect(component.parentCacheGroups().length).toBe(initialLength-1);
		component.cacheGroup = original;
	});
	it("gets available secondary parent Cache Groups", async () => {
		const cgs = await cgSrv.getCacheGroups();
		expect(cgs.length).toBeGreaterThan(0);
		component.cacheGroups = cgs;
		component.cacheGroup.fallbacks = [`${cgs[0].name}-fallback-test`];

		expect(component.secondaryParentCacheGroups()).toEqual(component.cacheGroups);
		if (component.cacheGroups.length < 1) {
			return fail("need at least one cache group to test parentage");
		}
		const initialLength = component.secondaryParentCacheGroups().length;
		const cg = component.cacheGroups[0];
		const original = component.cacheGroup;
		component.cacheGroup = {
			...component.cacheGroup,
			parentCachegroupId: cg.id,
			parentCachegroupName: cg.name
		};
		expect(component.secondaryParentCacheGroups()).not.toContain(cg);
		expect(component.secondaryParentCacheGroups().length).toBe(initialLength-1);
		component.cacheGroup = original;
	});
	it("gets available fallback Cache Groups", async () => {
		const cgs = await cgSrv.getCacheGroups();
		expect(cgs.length).toBeGreaterThan(2);
		component.cacheGroups = cgs;
		component.cacheGroup.fallbacks = [`${cgs[0].name}-fallback-test`];

		expect(component.fallbacks()).toEqual(component.cacheGroups);
		if (component.cacheGroups.length < 1) {
			return fail("need at least one cache group to test parentage");
		}
		const initialLength = component.fallbacks().length;
		const [cg1, cg2] = component.cacheGroups;
		const original = component.cacheGroup;
		component.cacheGroup = {
			...component.cacheGroup,
			parentCachegroupId: cg1.id,
			parentCachegroupName: cg1.name,
			secondaryParentCachegroupId: cg2.id,
			secondaryParentCachegroupName: cg2.name
		};
		const fallbacks = component.fallbacks();
		expect(fallbacks).not.toContain(cg1);
		expect(fallbacks).not.toContain(cg2);
		expect(fallbacks.length).toBe(initialLength-2);
		component.cacheGroup = original;
	});

	it("refuses to delete new Cache Groups", async () => {
		component.new = true;
		const spy = spyOn(cgSrv, "deleteCacheGroup");

		const asyncExpectation = expectAsync(component.delete()).toBeResolvedTo(undefined);
		await component.delete();
		const dialogs = await loader.getAllHarnesses(MatDialogHarness);
		expect(dialogs.length).toBe(0);
		expect(spy).not.toHaveBeenCalled();

		await asyncExpectation;
	});
	it("deletes existing Cache Groups", async () => {
		const spy = spyOn(cgSrv, "deleteCacheGroup").and.callThrough();
		let cgs = await cgSrv.getCacheGroups();
		const initialLength = cgs.length;
		if (initialLength < 1) {
			return fail("need at least one Cache Group");
		}
		const cg = cgs[0];
		component.cacheGroup = cg;
		component.new = false;

		const asyncExpectation = expectAsync(component.delete()).toBeResolvedTo(undefined);
		const dialogs = await loader.getAllHarnesses(MatDialogHarness);
		if (dialogs.length !== 1) {
			return fail(`failed to open dialog; ${dialogs.length} dialogs found`);
		}
		const dialog = dialogs[0];
		const buttons = await dialog.getAllHarnesses(MatButtonHarness.with({text: /^[cC][oO][nN][fF][iI][rR][mM]$/}));
		if (buttons.length !== 1) {
			return fail(`'confirm' button not found; ${buttons.length} buttons found`);
		}
		await buttons[0].click();

		expect(spy).toHaveBeenCalledOnceWith(cg);

		cgs = await cgSrv.getCacheGroups();
		expect(cgs).not.toContain(cg);
		expect(cgs.length).toBe(initialLength - 1);

		await asyncExpectation;
	});

	it("creates new Cache Groups", async () => {
		const createSpy = spyOn(cgSrv, "createCacheGroup").and.callThrough();
		const updateSpy = spyOn(cgSrv, "updateCacheGroup").and.callThrough();

		component.new = true;
		const cg = component.cacheGroup;
		const typeSrv = TestBed.inject(TypeService);
		const types = await typeSrv.getTypesInTable("cachegroup");
		if (types.length < 1) {
			return fail("no cg Types");
		}
		component.typeCtrl.setValue(types[0].id);
		await expectAsync(component.submit(new Event("click"))).toBeResolvedTo(undefined);
		expect(createSpy).toHaveBeenCalledOnceWith(cg);
		expect(updateSpy).not.toHaveBeenCalled();
		expect(component.new).toBeFalse();
	});

	it("updates existing Cache Groups", async () => {
		const createSpy = spyOn(cgSrv, "createCacheGroup").and.callThrough();
		component.new = false;
		const cg = component.cacheGroup;
		const typeSrv = TestBed.inject(TypeService);
		const types = await typeSrv.getTypesInTable("cachegroup");
		if (types.length < 1) {
			return fail("no cg Types");
		}
		const updateSpy = spyOn(cgSrv, "updateCacheGroup").and.returnValue(new Promise(r => r(component.cacheGroup)));
		component.typeCtrl.setValue(types[0].id);
		await expectAsync(component.submit(new Event("click"))).toBeResolvedTo(undefined);
		expect(updateSpy).toHaveBeenCalledOnceWith(cg);
		expect(createSpy).not.toHaveBeenCalled();
		expect(component.new).toBeFalse();
	});

	it("doesn't submit a request when the form is invalid", async () => {
		const createSpy = spyOn(cgSrv, "createCacheGroup").and.callThrough();
		const updateSpy = spyOn(cgSrv, "updateCacheGroup").and.callThrough();

		component.typeCtrl.setErrors({something: true});
		await expectAsync(component.submit(new Event("click"))).toBeResolvedTo(undefined);
		expect(updateSpy).not.toHaveBeenCalled();
		expect(createSpy).not.toHaveBeenCalled();

		component.typeCtrl.setErrors(null);
		component.typeCtrl.setValue(null);
		await expectAsync(component.submit(new Event("click"))).toBeResolvedTo(undefined);
		expect(updateSpy).not.toHaveBeenCalled();
		expect(createSpy).not.toHaveBeenCalled();
	});
});
