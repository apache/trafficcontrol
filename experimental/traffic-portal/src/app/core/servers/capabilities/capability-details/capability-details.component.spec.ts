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
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { MatDialog, MatDialogModule, type MatDialogRef } from "@angular/material/dialog";
import { ActivatedRoute } from "@angular/router";
import { RouterTestingModule } from "@angular/router/testing";
import { of } from "rxjs";

import { ServerService } from "src/app/api";
import { APITestingModule } from "src/app/api/testing";

import { CapabilityDetailsComponent } from "./capability-details.component";

describe("CapabilityDetailsComponent", () => {
	let component: CapabilityDetailsComponent;
	let fixture: ComponentFixture<CapabilityDetailsComponent>;
	let paramMap: jasmine.Spy;
	let route: ActivatedRoute;
	const name = "test";

	beforeEach(async () => {
		await TestBed.configureTestingModule({
			declarations: [ CapabilityDetailsComponent ],
			imports: [
				APITestingModule,
				RouterTestingModule.withRoutes([
					{component: CapabilityDetailsComponent, path: "core/capabilities/:name"},
					// This route is never actually loaded, but the tests
					// complain that it can't be routed to, so it doesn't matter
					// that it's loading the wrong component, only that it has a
					// route definition.
					{component: CapabilityDetailsComponent, path: "core/capabilities"}
				]),
				MatDialogModule
			],
		}).compileComponents();

		route = TestBed.inject(ActivatedRoute);
		paramMap = spyOn(route.snapshot.paramMap, "get");
		paramMap.and.returnValue(name);
		fixture = TestBed.createComponent(CapabilityDetailsComponent);
		component = fixture.componentInstance;
		component.capability = {...await TestBed.inject(ServerService).createCapability({name})};
		fixture.detectChanges();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});

	it("sets up the form for a new Capability", async () => {
		paramMap.and.returnValue(null);

		fixture = TestBed.createComponent(CapabilityDetailsComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
		await fixture.whenStable();
		expect(paramMap).toHaveBeenCalled();
		expect(component.capability).not.toBeNull();
		expect(component.capability.name).toBe("");
		expect(component.new).toBeTrue();
	});

	it("existing Capability", async () => {
		expect(paramMap).toHaveBeenCalled();
		expect(component.capability).not.toBeNull();
		expect(component.capability.name).toBe(name);
		expect(component.new).toBeFalse();
	});

	it("opens a dialog for Capability deletion", async () => {
		const api = TestBed.inject(ServerService);
		const spy = spyOn(api, "deleteCapability").and.callThrough();
		await expect(spy).not.toHaveBeenCalled();

		const dialogService = TestBed.inject(MatDialog);
		const openSpy = spyOn(dialogService, "open").and.returnValue({
			afterClosed: () => of(true)
		} as MatDialogRef<unknown>);

		await expect(openSpy).not.toHaveBeenCalled();

		const asyncExpectation = expectAsync(component.deleteCapability()).toBeResolvedTo(undefined);

		await expect(openSpy).toHaveBeenCalled();
		await expect(spy).toHaveBeenCalled();

		await asyncExpectation;
	});

	it("submits requests to create new Capabilities", async () => {
		const api = TestBed.inject(ServerService);
		const spy = spyOn(api, "createCapability").and.callThrough();
		paramMap.and.returnValue(null);

		fixture = TestBed.createComponent(CapabilityDetailsComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
		await fixture.whenStable();

		expect(spy).not.toHaveBeenCalled();
		await expectAsync(component.submit(new Event("submit"))).toBeResolvedTo(undefined);
		expect(spy).toHaveBeenCalled();
		expect(component.new).toBeFalse();
	});

	it("submits requests to update Capabilities", async () => {
		const api = TestBed.inject(ServerService);
		const spy = spyOn(api, "updateCapability").and.callThrough();
		expect(spy).not.toHaveBeenCalled();

		component.capability = {
			...component.capability,
			name: `${component.capability.name}quest`
		};

		await expectAsync(component.submit(new Event("submit"))).toBeResolvedTo(undefined);
		expect(spy).toHaveBeenCalled();
		expect(component.new).toBeFalse();
	});
});
