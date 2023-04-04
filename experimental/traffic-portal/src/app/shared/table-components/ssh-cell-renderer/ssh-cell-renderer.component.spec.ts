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
import { waitForAsync, type ComponentFixture, TestBed } from "@angular/core/testing";
import { RouterTestingModule } from "@angular/router/testing";
import type { ICellRendererParams } from "ag-grid-community";

import { CurrentUserService } from "src/app/shared/current-user/current-user.service";

import { SSHCellRendererComponent } from "./ssh-cell-renderer.component";

describe("SshCellRendererComponent", () => {
	let component: SSHCellRendererComponent;
	let fixture: ComponentFixture<SSHCellRendererComponent>;

	beforeEach(waitForAsync(() => {
		const mockCurrentUserService = jasmine.createSpyObj(
			["updateCurrentUser", "login", "logout"],
			{currentUser: {username: "test-admin"}}
		);
		mockCurrentUserService.updateCurrentUser.and.returnValue(new Promise(r => r(true)));

		TestBed.configureTestingModule({
			declarations: [ SSHCellRendererComponent ],
			imports: [HttpClientModule, RouterTestingModule],
			providers: [ { provide: CurrentUserService, useValue: mockCurrentUserService} ]
		}).compileComponents();
	}));

	beforeEach(() => {
		fixture = TestBed.createComponent(SSHCellRendererComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});

	it("initializes", () => {
		component.agInit({value: "192.0.2.1"} as ICellRendererParams);
		expect(component.value).toBe("192.0.2.1");

		component.agInit({value: "192.0.2.2"} as ICellRendererParams);
		expect(component.value).toBe("192.0.2.2");
	});

	it("refreshes", () => {
		let ret = component.refresh({value: "192.0.2.1"} as ICellRendererParams);
		expect(ret).toBeTrue();
		expect(component.value).toBe("192.0.2.1");

		ret = component.refresh({value: "192.0.2.2"} as ICellRendererParams);
		expect(ret).toBeTrue();
		expect(component.value).toBe("192.0.2.2");
	});
});
