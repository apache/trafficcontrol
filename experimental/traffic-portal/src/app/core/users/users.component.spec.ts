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
import { waitForAsync, ComponentFixture, TestBed } from "@angular/core/testing";
import { FormsModule, ReactiveFormsModule } from "@angular/forms";
import { RouterTestingModule } from "@angular/router/testing";

import { CurrentUserService } from "src/app/shared/currentUser/current-user.service";
import type { User } from "src/app/models";
import { TpHeaderComponent } from "src/app/shared/tp-header/tp-header.component";
import { LoadingComponent } from "src/app/shared/loading/loading.component";
import { UserService } from "src/app/shared/api";
import { UsersComponent } from "./users.component";

describe("UsersComponent", () => {
	let component: UsersComponent;
	let fixture: ComponentFixture<UsersComponent>;

	beforeEach(waitForAsync(() => {
		// mock the API
		const mockAPIService = jasmine.createSpyObj(["getUsers", "getRoles", "getCurrentUser"]);
		mockAPIService.getUsers.and.returnValue(new Promise(resolve => resolve([])));
		mockAPIService.getRoles.and.returnValue(new Promise(resolve => resolve([])));
		mockAPIService.getCurrentUser.and.returnValue(new Promise(resolve => resolve({
			id: 0,
			newUser: false,
			username: "test"
		} as User)));
		const mockCurrentUserService = jasmine.createSpyObj(["updateCurrentUser", "login", "logout"]);
		mockCurrentUserService.updateCurrentUser.and.returnValue(new Promise(r => r(false)));

		TestBed.configureTestingModule({
			declarations: [
				UsersComponent,
				LoadingComponent,
				TpHeaderComponent,
			],
			imports: [
				FormsModule,
				HttpClientModule,
				ReactiveFormsModule,
				RouterTestingModule
			],
			providers: [
				{ provide: UserService, useValue: mockAPIService },
				{ provide: CurrentUserService, useValue: mockCurrentUserService }
			]
		});
		TestBed.compileComponents();
	}));

	beforeEach(() => {
		fixture = TestBed.createComponent(UsersComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it("should exist", () => {
		expect(component).toBeTruthy();
	});

	it("can tell if a user has a location", () => {
		const u: User = {
			city: "Townsville",
			country: "Countryland",
			id: -1,
			newUser: false,
			postalCode: "00000",
			stateOrProvince: "Provincia",
			username: "test"
		};
		expect(component.userHasLocation(u)).toBeTrue();
		u.city = null;
		expect(component.userHasLocation(u)).toBeTrue();
		u.stateOrProvince = null;
		expect(component.userHasLocation(u)).toBeTrue();
		u.country = null;
		expect(component.userHasLocation(u)).toBeTrue();
		u.postalCode = null;
		expect(component.userHasLocation(u)).toBeFalse();
		u.country = "Countryland";
		expect(component.userHasLocation(u)).toBeTrue();
		u.country = null;
		u.stateOrProvince = "Provincia";
		expect(component.userHasLocation(u)).toBeTrue();
		u.stateOrProvince = null;
		u.city = "Townsville";
		expect(component.userHasLocation(u)).toBeTrue();
	});

	it("builds user location strings", () => {
		const u: User = {
			city: "Townsville",
			country: "Countryland",
			id: -1,
			newUser: false,
			postalCode: "00000",
			stateOrProvince: "Provincia",
			username: "test"
		};
		expect(component.userLocationString(u)).toBe("Townsville, Provincia, Countryland, 00000");
		u.city = null;
		expect(component.userLocationString(u)).toBe("Provincia, Countryland, 00000");
		u.stateOrProvince = null;
		expect(component.userLocationString(u)).toBe("Countryland, 00000");
		u.country = null;
		expect(component.userLocationString(u)).toBe("00000");
		u.postalCode = null;
		expect(component.userLocationString(u)).toBeNull();
		u.country = "Countryland";
		expect(component.userLocationString(u)).toBe("Countryland");
		u.country = null;
		u.stateOrProvince = "Provincia";
		expect(component.userLocationString(u)).toBe("Provincia");
		u.stateOrProvince = null;
		u.city = "Townsville";
		expect(component.userLocationString(u)).toBe("Townsville");
	});

	it("searches fuzz-ily", ()=>{
		const u = {
			id: -1,
			newUser: false,
			username: "test"
		};
		expect(component.fuzzControl.value).toBe("");
		expect(component.fuzzy(u)).toBeTrue();
		component.fuzzControl.setValue(`${u.username}z`);
		expect(component.fuzzy(u)).toBeFalse();
		component.fuzzControl.setValue(u.username);
		expect(component.fuzzy(u)).toBeTrue();
		component.fuzzControl.setValue(`${u.username[0]}${u.username.slice(-1)[0]}`);
		expect(component.fuzzy(u)).toBeTrue();
	});

	afterAll(() => {
		try{
			TestBed.resetTestingModule();
		} catch (e) {
			console.error("error in UsersComponent afterAll:", e);
		}
	});
});
