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
import { fakeAsync, TestBed, tick } from "@angular/core/testing";
import { MatSnackBarModule } from "@angular/material/snack-bar";
import { Router } from "@angular/router";
import { RouterTestingModule } from "@angular/router/testing";
import { of } from "rxjs";

import { CurrentUserService } from "src/app/shared/current-user/current-user.service";

import { AppComponent } from "./app.component";

describe("AppComponent", () => {
	let component: AppComponent;
	let mockCurrentUserService: jasmine.SpyObj<CurrentUserService>;

	beforeEach(() => {
		mockCurrentUserService = jasmine.createSpyObj(["updateCurrentUser", "login", "logout"], {userChanged: of(null)});
		TestBed.configureTestingModule({
			declarations: [
				AppComponent
			],
			imports: [
				HttpClientModule,
				RouterTestingModule.withRoutes([{component: AppComponent, path: "login"}]),
				MatSnackBarModule
			],
			providers: [ { provide: CurrentUserService, useValue: mockCurrentUserService }]
		}).compileComponents();
		const fixture = TestBed.createComponent(AppComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it("should create the app", () => {
		expect(component).toBeTruthy();
	});

	it("logs out", fakeAsync(() => {
		// eslint-disable-next-line @typescript-eslint/unbound-method
		expect(mockCurrentUserService.logout).not.toHaveBeenCalled();
		component.logout();
		tick();
		// eslint-disable-next-line @typescript-eslint/unbound-method
		expect(mockCurrentUserService.logout).toHaveBeenCalled();
		const router = TestBed.inject(Router);
		expect(router.url).toBe("/login");
	}));
});
