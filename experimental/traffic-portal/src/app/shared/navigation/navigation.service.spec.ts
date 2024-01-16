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
import { TestBed } from "@angular/core/testing";
import { MatButtonModule } from "@angular/material/button";
import { MatMenuModule } from "@angular/material/menu";
import { Title } from "@angular/platform-browser";
import { RouterTestingModule } from "@angular/router/testing";
import { of } from "rxjs";

import { UserService } from "src/app/api";
import { APITestingModule } from "src/app/api/testing";
import { CurrentUserService } from "src/app/shared/current-user/current-user.service";
import { NavigationService } from "src/app/shared/navigation/navigation.service";
import { type TpHeaderComponent } from "src/app/shared/navigation/tp-header/tp-header.component";

describe("NavigationService", () => {
	let service: NavigationService;
	let mockHeaderComp: jasmine.SpyObj<TpHeaderComponent>;
	let logOutSpy: jasmine.Spy;

	beforeEach(() => {
		const mockCurrentUserService = jasmine.createSpyObj(
			["updateCurrentUser", "hasPermission", "login", "logout"], {userChanged: of(null)});
		logOutSpy = mockCurrentUserService.logout;
		mockHeaderComp = jasmine.createSpyObj<TpHeaderComponent>([], {hidden: false, title: ""});
		TestBed.configureTestingModule({
			imports: [APITestingModule, HttpClientModule, RouterTestingModule, MatMenuModule, MatButtonModule],
			providers: [
				NavigationService,
				{provide: CurrentUserService, useValue: mockCurrentUserService},
			],
		});
		service = TestBed.inject(NavigationService);
	});

	it("should be created", () => {
		expect(service).toBeTruthy();
	});

	it("clears front-end user data even if server-side logout fails", async () => {
		const userService = TestBed.inject(UserService);
		const userSpy = spyOn(userService, "logout");
		userSpy.and.returnValue(new Promise(r => r(null)));
		expect(userSpy).not.toHaveBeenCalled();
		await service.logout();
		expect(userSpy).toHaveBeenCalled();
		expect(logOutSpy).toHaveBeenCalled();
	});

	it("logs the user out", async () => {
		expect(logOutSpy).not.toHaveBeenCalled();
		await service.logout();
		expect(logOutSpy).toHaveBeenCalled();
	});

	it("sets the page title and header", () => {
		expect(mockHeaderComp).toBeTruthy();
		expect(mockHeaderComp.hidden).toBeFalse();
		expect(mockHeaderComp.title).toBe("");

		const title = "something else";
		const titleService = TestBed.inject(Title);

		service.headerTitle.next(title);
		expect(titleService.getTitle()).toBe(title);
	});
});
