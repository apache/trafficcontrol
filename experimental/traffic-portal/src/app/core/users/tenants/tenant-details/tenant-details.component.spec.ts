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
import { MatDialogModule } from "@angular/material/dialog";
import { ActivatedRoute } from "@angular/router";
import { RouterTestingModule } from "@angular/router/testing";

import { APITestingModule } from "src/app/api/testing";

import { TenantDetailsComponent } from "./tenant-details.component";

describe("TenantDetailsComponent", () => {
	let component: TenantDetailsComponent;
	let fixture: ComponentFixture<TenantDetailsComponent>;
	let route: ActivatedRoute;
	let paramMap: jasmine.Spy;

	beforeEach(async () => {
		await TestBed.configureTestingModule({
			declarations: [ TenantDetailsComponent ],
			imports: [ APITestingModule, MatDialogModule, RouterTestingModule ],
		})
			.compileComponents();

		route = TestBed.inject(ActivatedRoute);
		paramMap = spyOn(route.snapshot.paramMap, "get");
		paramMap.and.returnValue(null);
		fixture = TestBed.createComponent(TenantDetailsComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it("should create", async () => {
		expect(component).toBeTruthy();
		expect(paramMap).toHaveBeenCalled();
		expect(component.tenants.length).toBe(0);
	});

	it("new tenant", async () => {
		paramMap.and.returnValue("new");

		fixture = TestBed.createComponent(TenantDetailsComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
		await fixture.whenStable();
		expect(paramMap).toHaveBeenCalled();
		expect(component.tenants.length).toBe(2);
		expect(component.tenant.name).toBe("");
		expect(component.new).toBeTrue();
		expect(component.displayTenant.name).toBe("root");
		expect(component.disabled).toBeFalse();
	});

	it("existing root tenant", async () => {
		paramMap.and.returnValue("1");

		fixture = TestBed.createComponent(TenantDetailsComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
		await fixture.whenStable();
		expect(paramMap).toHaveBeenCalled();
		expect(component.tenants.length).toBe(2);
		expect(component.tenant.name).toBe("root");
		expect(component.new).toBeFalse();
		expect(component.displayTenant.name).toBe("root");
		expect(component.disabled).toBeTrue();
	});

	it("existing non-root tenant", async () => {
		paramMap.and.returnValue("2");

		fixture = TestBed.createComponent(TenantDetailsComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
		await fixture.whenStable();
		expect(paramMap).toHaveBeenCalled();
		expect(component.tenants.length).toBe(2);
		expect(component.tenant.name).toBe("test");
		expect(component.new).toBeFalse();
		expect(component.displayTenant.name).toBe("root");
		expect(component.disabled).toBeFalse();
	});
});
