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
import { ReplaySubject } from "rxjs";

import { APITestingModule } from "src/app/api/testing";
import { DivisionDetailComponent } from "src/app/core/cache-groups/divisions/detail/division-detail.component";
import { TpHeaderService } from "src/app/shared/tp-header/tp-header.service";

describe("DetailComponent", () => {
	let component: DivisionDetailComponent;
	let fixture: ComponentFixture<DivisionDetailComponent>;
	let route: ActivatedRoute;
	let paramMap: jasmine.Spy;

	const headerSvc = jasmine.createSpyObj([],{headerHidden: new ReplaySubject<boolean>(), headerTitle: new ReplaySubject<string>()});
	beforeEach(async () => {
		await TestBed.configureTestingModule({
			declarations: [ DivisionDetailComponent ],
			imports: [ APITestingModule, RouterTestingModule, MatDialogModule ],
			providers: [ { provide: TpHeaderService, useValue: headerSvc } ]
		})
			.compileComponents();

		route = TestBed.inject(ActivatedRoute);
		paramMap = spyOn(route.snapshot.paramMap, "get");
		paramMap.and.returnValue(null);
		fixture = TestBed.createComponent(DivisionDetailComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
		expect(paramMap).toHaveBeenCalled();
	});

	it("new division", async () => {
		paramMap.and.returnValue("new");

		fixture = TestBed.createComponent(DivisionDetailComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
		await fixture.whenStable();
		expect(paramMap).toHaveBeenCalled();
		expect(component.division).not.toBeNull();
		expect(component.division.name).toBe("");
		expect(component.new).toBeTrue();
	});

	it("existing division", async () => {
		paramMap.and.returnValue("1");

		fixture = TestBed.createComponent(DivisionDetailComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
		await fixture.whenStable();
		expect(paramMap).toHaveBeenCalled();
		expect(component.division).not.toBeNull();
		expect(component.division.name).toBe("Div1");
		expect(component.new).toBeFalse();
	});
});
