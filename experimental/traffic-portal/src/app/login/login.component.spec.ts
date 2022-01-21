/* eslint-disable @typescript-eslint/unbound-method */
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
import { HttpClientModule } from "@angular/common/http";
import { type ComponentFixture, TestBed, fakeAsync, tick } from "@angular/core/testing";
import { FormsModule, ReactiveFormsModule } from "@angular/forms";
import { MatDialog, MatDialogModule } from "@angular/material/dialog";
import { Router } from "@angular/router";
import { RouterTestingModule } from "@angular/router/testing";

import { CurrentUserService } from "src/app/shared/currentUser/current-user.service";

import { LoginComponent } from "./login.component";

describe("LoginComponent", () => {
	let component: LoginComponent;
	let fixture: ComponentFixture<LoginComponent>;
	let mockCurrentUserService: jasmine.SpyObj<CurrentUserService>;
	let dialogSpy: jasmine.Spy;
	let router: Router;

	beforeEach(() => {
		mockCurrentUserService = jasmine.createSpyObj(["updateCurrentUser", "login", "logout"]);
		mockCurrentUserService.login.and.callFake(async (u, p)=>u === "test-admin" && p === "twelve12!");
		mockCurrentUserService.login.withArgs("tok").and.returnValue(new Promise(r=>r(true)));
		mockCurrentUserService.login.withArgs("badToken").and.callFake(
			async () => {
				throw new Error("bad token");
			}
		);
		mockCurrentUserService.login.withArgs("server error", "twelve12!").and.callFake(
			async () => {
				throw new Error("some kind of server error occurred");
			}
		);

		const dialog = jasmine.createSpyObj(["open"]);
		dialogSpy = dialog.open;

		TestBed.configureTestingModule({
			declarations: [ LoginComponent ],
			imports: [
				FormsModule,
				MatDialogModule,
				HttpClientModule,
				ReactiveFormsModule,
				RouterTestingModule.withRoutes([
					{component: LoginComponent, path: "login"},
					// This obviously isn't how this actually works, but we
					// don't care about testing anything on that page, so this
					// will do fine.
					{component: LoginComponent, path: "core/me"}
				]),
			],
			providers: [
				{ provide: CurrentUserService, useValue: mockCurrentUserService },
				{ provide: MatDialog, useValue: dialog}
			]
		}).compileComponents();
		fixture = TestBed.createComponent(LoginComponent);
		component = fixture.componentInstance;
		router = TestBed.inject(Router);
		router.initialNavigation();
		fixture.detectChanges();
	});

	it("should exist", () => {
		try{
			expect(component).toBeTruthy();
		} catch (e) {
			console.error("error in 'should exist' for LoginComponent:", e);
		}
	});

	it("submits a login request", () => {
		expect(mockCurrentUserService.login).not.toHaveBeenCalled();
		component.submitLogin();
		expect(mockCurrentUserService.login).toHaveBeenCalled();
		component.u.setValue("test-admin");
		component.p.setValue("twelve12!");
		component.submitLogin();
		expect(mockCurrentUserService.login).toHaveBeenCalledTimes(2);
		component.u.setValue("server error");
		component.submitLogin();
		expect(mockCurrentUserService.login).toHaveBeenCalledTimes(3);
	});

	it("opens the password reset dialog", () => {
		expect(dialogSpy).not.toHaveBeenCalled();
		component.resetPassword();
		expect(dialogSpy).toHaveBeenCalled();
	});

	it("redirects to the user edit page on token login", fakeAsync(() => {
		router.navigate(["/login"], {queryParams: {token: "tok"}});
		tick();
		expect(router.url).toBe("/login?token=tok");
		expect(mockCurrentUserService.login).not.toHaveBeenCalled();

		// need to re-run this to pick up the token; simulates component
		// initialization.
		component.ngOnInit();
		tick();
		expect(mockCurrentUserService.login).toHaveBeenCalled();
		expect(router.navigated).toBeTrue();
		const [path, query] = router.url.split("?");
		expect(path).toBe("/core/me");
		const kvps = query.split("&");
		expect(kvps.length).toBe(2);
		expect(kvps).toContain("edit=true");
		expect(kvps).toContain("updatePassword=true");

		router.navigate(["/login"], {queryParams: {token: "badToken"}});
		tick();
		expect(router.url).toBe("/login?token=badToken");
		component.ngOnInit();
		tick();
		expect(mockCurrentUserService.login).toHaveBeenCalledTimes(2);
		expect(router.url).toBe("/login?token=badToken");
	}));
});
