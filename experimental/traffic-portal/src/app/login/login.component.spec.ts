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

import { LoginComponent } from "./login.component";
import {AuthenticationService} from "../shared/authentication/authentication.service";

describe("LoginComponent", () => {
	let component: LoginComponent;
	let fixture: ComponentFixture<LoginComponent>;

	beforeEach(waitForAsync(() => {
		const mockAuthenticationService = jasmine.createSpyObj(["updateCurrentUser", "login", "logout"]);
		TestBed.configureTestingModule({
			declarations: [ LoginComponent ],
			imports: [
				FormsModule,
				HttpClientModule,
				ReactiveFormsModule,
				RouterTestingModule
			],
			providers: [ { provide: AuthenticationService, useValue: mockAuthenticationService }]
		})
			.compileComponents();
	}));

	beforeEach(async () => {
		fixture = TestBed.createComponent(LoginComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it("should exist", async () => {
		try{
			expect(component).toBeTruthy();
		} catch (e) {
			console.error("error in 'should exist' for LoginComponent:", e);
		}
	});

	afterAll(async () => {
		try{
			TestBed.resetTestingModule();
		} catch (e) {
			console.error("error in LoginComponent afterAll:", e);
		}
	});
});
