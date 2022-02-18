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

import { APITestingModule } from "src/app/api/testing";
import { CurrentUserService } from "src/app/shared/currentUser/current-user.service";
import { LoadingComponent } from "src/app/shared/loading/loading.component";
import { TpHeaderComponent } from "src/app/shared/tp-header/tp-header.component";

import { UsersComponent } from "./users.component";

describe("UsersComponent", () => {
	let component: UsersComponent;
	let fixture: ComponentFixture<UsersComponent>;

	beforeEach(waitForAsync(() => {
		// mock the API
		const mockCurrentUserService = jasmine.createSpyObj(["updateCurrentUser", "login", "logout"]);
		mockCurrentUserService.updateCurrentUser.and.returnValue(new Promise(r => r(false)));

		TestBed.configureTestingModule({
			declarations: [
				UsersComponent,
				LoadingComponent,
				TpHeaderComponent,
			],
			imports: [
				APITestingModule,
				FormsModule,
				HttpClientModule,
				ReactiveFormsModule,
				RouterTestingModule
			],
			providers: [
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

	it("handles its context menu actions", () => {
		expect(()=>component.handleContextMenu({action: "viewDetails", data: []})).not.toThrow();
		expect(()=>component.handleContextMenu({action: "unknown action", data: []})).toThrow();
	});

	it("gets display strings for Roles", () => {
		component.roles = new Map([[1, "admin"]]);
		expect(component.roleDisplayString(1)).toBe("admin (#1)");
		expect(()=>component.roleDisplayString(2)).toThrow();
	});
});
