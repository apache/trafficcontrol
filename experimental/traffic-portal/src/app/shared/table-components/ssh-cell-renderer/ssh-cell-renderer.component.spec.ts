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

import { SSHCellRendererComponent } from "./ssh-cell-renderer.component";
import {AuthenticationService} from "../../authentication/authentication.service";

describe("SshCellRendererComponent", () => {
	let component: SSHCellRendererComponent;
	let fixture: ComponentFixture<SSHCellRendererComponent>;

	beforeEach(waitForAsync(() => {
		const mockAuthenticationService = jasmine.createSpyObj(["updateCurrentUser", "login", "logout"]);
		mockAuthenticationService.updateCurrentUser.and.returnValue(new Promise(r => r(false)));
		TestBed.configureTestingModule({
			declarations: [ SSHCellRendererComponent ],
			imports: [HttpClientModule, RouterTestingModule],
			providers: [ { provide: AuthenticationService, useValue: mockAuthenticationService} ]
		})
			.compileComponents();
	}));

	beforeEach(() => {
		fixture = TestBed.createComponent(SSHCellRendererComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});
});
