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
import { RouterTestingModule } from "@angular/router/testing";

import { CurrentUserService } from "src/app/shared/currentUser/current-user.service";
import {UserService} from "../api";
import { TpHeaderComponent } from "./tp-header.component";

describe("TpHeaderComponent", () => {
	let component: TpHeaderComponent;
	let fixture: ComponentFixture<TpHeaderComponent>;

	beforeEach(waitForAsync(() => {
		const mockAPIService = jasmine.createSpyObj(["getUsers"]);
		const mockCurrentUserService = jasmine.createSpyObj(["updateCurrentUser", "login", "logout"]);
		TestBed.configureTestingModule({
			declarations: [ TpHeaderComponent ],
			imports: [ HttpClientModule, RouterTestingModule ],
			providers: [
				{ provide: CurrentUserService, useValue: mockCurrentUserService },
				{ provide: UserService, useValue: mockAPIService}
			]
		})
			.compileComponents();
	}));

	beforeEach(() => {
		fixture = TestBed.createComponent(TpHeaderComponent);
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
			console.error("error in TpHeaderComponent afterAll:", e);
		}
	});
});
