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

import { OriginService } from "src/app/api";
import { APITestingModule } from "src/app/api/testing";
import { NavigationService } from "src/app/shared/navigation/navigation.service";

import { OriginDetailComponent } from "./origin-detail.component";

describe("OriginDetailComponent", () => {
	let component: OriginDetailComponent;
	let fixture: ComponentFixture<OriginDetailComponent>;
	let route: ActivatedRoute;
	let paramMap: jasmine.Spy;
	let loader: HarnessLoader;
	let service: OriginService;

	const navSvc = jasmine.createSpyObj([], {
		headerHidden: new ReplaySubject<boolean>(),
		headerTitle: new ReplaySubject<string>(),
	});
	beforeEach(async () => {
		await TestBed.configureTestingModule({
			declarations: [OriginDetailComponent],
			imports: [APITestingModule, RouterTestingModule.withRoutes( [
				{component: OriginDetailComponent, path: "core/origins/:id"},
				{component: OriginDetailComponent, path: "core/origins"},
			]), MatDialogModule, NoopAnimationsModule],
			providers: [{provide: NavigationService, useValue: navSvc}],
		}).compileComponents();

		route = TestBed.inject(ActivatedRoute);
		paramMap = spyOn(route.snapshot.paramMap, "get");
		paramMap.and.returnValue(null);
		fixture = TestBed.createComponent(OriginDetailComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
		service = TestBed.inject(OriginService);
		loader = TestbedHarnessEnvironment.documentRootLoader(fixture);
		await fixture.whenStable();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});

	it("new origin", async () => {
		paramMap.and.returnValue("new");

		fixture = TestBed.createComponent(OriginDetailComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
		await fixture.whenStable();
		expect(paramMap).toHaveBeenCalled();
		expect(component.origin).not.toBeNull();
		expect(component.origin.name).toBe("");
		expect(component.new).toBeTrue();
	});

	it("existing origin", async () => {
		const id = 1;
		paramMap.and.returnValue(id);
		const origin = await service.getOrigins(id);
		fixture = TestBed.createComponent(OriginDetailComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
		await fixture.whenStable();
		expect(paramMap).toHaveBeenCalled();
		expect(component.origin).not.toBeNull();
		expect(component.origin.name).toBe(origin.name);
		expect(component.new).toBeFalse();
	});

	it("deletes existing Origins", async () => {
		const spy = spyOn(service, "deleteOrigin").and.callThrough();
		let orgs = await service.getOrigins();
		const initialLength = orgs.length;
		if (initialLength < 1) {
			return fail("need at least one Origin");
		}
		const org = orgs[0];
		component.origin = org;
		component.new = false;

		const asyncExpectation = expectAsync(component.deleteOrigin()).toBeResolvedTo(undefined);
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

		expect(spy).toHaveBeenCalledOnceWith(org);

		orgs = await service.getOrigins();
		expect(orgs).not.toContain(org);
		expect(orgs.length).toBe(initialLength - 1);

		await asyncExpectation;
	});

});
