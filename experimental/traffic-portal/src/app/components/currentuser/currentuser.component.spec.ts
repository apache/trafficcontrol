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

import { CurrentuserComponent } from "./currentuser.component";

import { User } from "../../models";
import { TpHeaderComponent } from "../tp-header/tp-header.component";

describe("CurrentuserComponent", () => {
	let component: CurrentuserComponent;
	let fixture: ComponentFixture<CurrentuserComponent>;

	beforeEach(waitForAsync(() => {
		TestBed.configureTestingModule({
			declarations: [
				CurrentuserComponent,
				TpHeaderComponent
			],
			imports: [
				HttpClientModule
			]
		})
		.compileComponents();
	}));

	beforeEach(() => {
		fixture = TestBed.createComponent(CurrentuserComponent);
		component = fixture.componentInstance;
		component.currentUser = {
			id: 1,
			newUser: false,
			username: "test"
		} as User;
		fixture.detectChanges();
	});

	it("should create", () => {
		component.currentUser = {
			id: 1,
			newUser: false,
			username: "test"
		} as User;
		expect(component).toBeTruthy();
	});

	afterAll(() => {
		TestBed.resetTestingModule();
	});
});
