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
import { MatDialog, MatDialogModule, MatDialogRef } from "@angular/material/dialog";
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
			imports: [
				APITestingModule,
				RouterTestingModule.withRoutes([
					{component: RoleDetailComponent, path: "core/roles/:name"},
					// This route is never actually loaded, but the tests
					// complain that it can't be routed to, so it doesn't matter
					// that it's loading the wrong component, only that it has a
					// route definition.
					{component: RoleDetailComponent, path: "core/roles"}
				]),
				MatDialogModule
			],
			providers: [ { provide: NavigationService, useValue: headerSvc } ]
		})
			.compileComponents();

		route = TestBed.inject(ActivatedRoute);
		paramMap = spyOn(route.snapshot.paramMap, "get");
		fixture = TestBed.createComponent(RoleDetailComponent);
		component = fixture.componentInstance;
		component.role = {...await TestBed.inject(UserService).createRole(component.role)};
		fixture.detectChanges();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});

	it("new role", async () => {
		paramMap.and.returnValue(null);

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
		paramMap.and.returnValue("admin");

		fixture = TestBed.createComponent(RoleDetailComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
		await fixture.whenStable();
		expect(paramMap).toHaveBeenCalled();
		expect(component.role).not.toBeNull();
		expect(component.role.name).toBe("admin");
		expect(component.new).toBeFalse();
	});

	it("opens a dialog for role deletion", async () => {
		const api = TestBed.inject(UserService);
		const spy = spyOn(api, "deleteRole").and.callThrough();
		await expect(spy).not.toHaveBeenCalled();

		const dialogService = TestBed.inject(MatDialog);
		const openSpy = spyOn(dialogService, "open").and.returnValue({
			afterClosed: () => of(true)
		} as MatDialogRef<unknown>);

		await expect(openSpy).not.toHaveBeenCalled();

		const asyncExpectation = expectAsync(component.deleteRole()).toBeResolvedTo(undefined);

		await expect(openSpy).toHaveBeenCalled();
		await expect(spy).toHaveBeenCalled();

		await asyncExpectation;
	});

	it("submits requests to create new role", async () => {
		const api = TestBed.inject(UserService);
		const spy = spyOn(api, "createRole").and.callThrough();
		paramMap.and.returnValue(null);

		fixture = TestBed.createComponent(RoleDetailComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
		await fixture.whenStable();

		expect(spy).not.toHaveBeenCalled();
		await expectAsync(component.submit(new Event("submit"))).toBeResolvedTo(undefined);
		expect(spy).toHaveBeenCalled();
		expect(component.new).toBeFalse();
	});

	it("submits requests to update role", async () => {
		const api = TestBed.inject(UserService);
		const spy = spyOn(api, "updateRole").and.callThrough();
		expect(spy).not.toHaveBeenCalled();

		component.role = {
			...component.role,
			name: `${component.role.name}quest`
		};

		await expectAsync(component.submit(new Event("submit"))).toBeResolvedTo(undefined);
		expect(spy).toHaveBeenCalled();
		expect(component.new).toBeFalse();
	});
});
