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
import { of, ReplaySubject } from "rxjs";

import { UserService } from "src/app/api";
import { APITestingModule } from "src/app/api/testing";
import { RoleDetailComponent } from "src/app/core/users/roles/detail/role-detail.component";
import { NavigationService } from "src/app/shared/navigation/navigation.service";

describe("RoleDetailComponent", () => {
	let component: RoleDetailComponent;
	let fixture: ComponentFixture<RoleDetailComponent>;
	let route: ActivatedRoute;
	let paramMap: jasmine.Spy;

	const headerSvc = jasmine.createSpyObj([],{headerHidden: new ReplaySubject<boolean>(), headerTitle: new ReplaySubject<string>()});
	beforeEach(async () => {
		await TestBed.configureTestingModule({
			declarations: [ RoleDetailComponent ],
			imports: [ APITestingModule, RouterTestingModule, MatDialogModule ],
			providers: [ { provide: NavigationService, useValue: headerSvc } ]
		})
			.compileComponents();

		route = TestBed.inject(ActivatedRoute);
		paramMap = spyOn(route.snapshot.paramMap, "get");
		paramMap.and.returnValue(null);
		fixture = TestBed.createComponent(RoleDetailComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
		expect(paramMap).toHaveBeenCalled();
	});

	it("new role", async () => {
		paramMap.and.returnValue("new");

		fixture = TestBed.createComponent(RoleDetailComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
		await fixture.whenStable();
		expect(paramMap).toHaveBeenCalled();
		expect(component.role).not.toBeNull();
		expect(component.role.name).toBe("");
		expect(component.new).toBeTrue();
	});

	it("existing role", async () => {
		expect(paramMap).toHaveBeenCalled();
		expect(component.role).not.toBeNull();
		expect(component.role.name).toBe("");
		expect(component.new).toBeFalse();
	});
});
