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

import { User } from "../../models";
import { APIService } from "../../services/api/apiservice";
import { LoadingComponent } from "../loading/loading.component";
import { TpHeaderComponent } from "../tp-header/tp-header.component";
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
		});
		TestBed.overrideProvider(APIService, { useValue: mockAPIService });
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

	afterAll(() => {
		try{
			TestBed.resetTestingModule();
		} catch (e) {
			console.error("error in UsersComponent afterAll:", e);
		}
	});
});
