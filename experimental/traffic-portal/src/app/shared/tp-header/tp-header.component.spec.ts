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
import {MatMenuModule} from "@angular/material/menu";
import { RouterTestingModule } from "@angular/router/testing";
import {of, ReplaySubject} from "rxjs";

import { UserService } from "src/app/api";
import { APITestingModule } from "src/app/api/testing";
import { CurrentUserService } from "src/app/shared/currentUser/current-user.service";
import {HeaderNavigation, TpHeaderService} from "src/app/shared/tp-header/tp-header.service";

import { TpHeaderComponent } from "./tp-header.component";

describe("TpHeaderComponent", () => {
	let component: TpHeaderComponent;
	let fixture: ComponentFixture<TpHeaderComponent>;
	let logOutSpy: jasmine.Spy;

	beforeEach(waitForAsync(() => {
		const mockCurrentUserService = jasmine.createSpyObj(
			["updateCurrentUser", "hasPermission", "login", "logout"], {userChanged: of(null)});
		logOutSpy = mockCurrentUserService.logout;
		const headerSvc = jasmine.createSpyObj(["addHorizontalNav", "addVerticalNav"],
			{headerHidden: new ReplaySubject<boolean>(), headerTitle: new ReplaySubject<string>(),
				horizontalNavsUpdated: new ReplaySubject<Array<HeaderNavigation>>(),
				verticalNavsUpdated: new ReplaySubject<Array<HeaderNavigation>>()});
		TestBed.configureTestingModule({
			declarations: [ TpHeaderComponent ],
			imports: [ APITestingModule, HttpClientModule, RouterTestingModule, MatMenuModule ],
			providers: [
				{ provide: CurrentUserService, useValue: mockCurrentUserService },
				{ provide: TpHeaderService, useValue: headerSvc}
			]
		}).compileComponents();
	}));

	beforeEach(() => {
		fixture = TestBed.createComponent(TpHeaderComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it("should exist", () => {
		expect(component).toBeTruthy();
	});

	it("logs the user out", async () => {
		expect(logOutSpy).not.toHaveBeenCalled();
		await component.logout();
		expect(logOutSpy).toHaveBeenCalled();
	});

	it("clears front-end user data even if server-side logout fails", async () => {
		const userService = TestBed.inject(UserService);
		const userSpy = spyOn(userService, "logout");
		userSpy.and.returnValue(new Promise(r=>r(null)));
		expect(userSpy).not.toHaveBeenCalled();
		await component.logout();
		expect(userSpy).toHaveBeenCalled();
		expect(logOutSpy).toHaveBeenCalled();
	});

	afterAll(() => {
		try{
			TestBed.resetTestingModule();
		} catch (e) {
			console.error("error in TpHeaderComponent afterAll:", e);
		}
	});
});
