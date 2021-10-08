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

import { ServersTableComponent } from "./servers-table.component";
import {TpHeaderComponent} from "../../../shared/tp-header/tp-header.component";
import {ServerService, UserService} from "../../../shared/api";
import {AuthenticationService} from "../../../shared/authentication/authentication.service";


describe("ServersTableComponent", () => {
	let component: ServersTableComponent;
	let fixture: ComponentFixture<ServersTableComponent>;

	beforeEach(waitForAsync(() => {
		const mockAPIService = jasmine.createSpyObj(["getServers", "getUsers"]);
		mockAPIService.getServers.and.returnValue(new Promise(r => r([])));
		const mockAuthenticationService = jasmine.createSpyObj(["updateCurrentUser", "login", "logout"]);
		TestBed.configureTestingModule({
			declarations: [ ServersTableComponent, TpHeaderComponent ],
			imports: [HttpClientModule, RouterTestingModule],
			providers: [
				{ provide: ServerService, useValue: mockAPIService },
				{ provide: UserService, useValue: mockAPIService },
				{ provide: AuthenticationService, useValue: mockAuthenticationService }
			]
		})
			.compileComponents();
	}));

	beforeEach(() => {
		fixture = TestBed.createComponent(ServersTableComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});
});
