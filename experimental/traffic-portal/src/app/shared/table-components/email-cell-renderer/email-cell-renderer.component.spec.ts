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
import { ComponentFixture, TestBed } from "@angular/core/testing";
import type { ICellRendererParams } from "ag-grid-community";

import { EmailCellRendererComponent } from "./email-cell-renderer.component";

describe("EmailCellRendererComponent", () => {
	let component: EmailCellRendererComponent;
	let fixture: ComponentFixture<EmailCellRendererComponent>;

	beforeEach(async () => {
		await TestBed.configureTestingModule({
			declarations: [ EmailCellRendererComponent ]
		}).compileComponents();
		fixture = TestBed.createComponent(EmailCellRendererComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});

	it("initializes", () => {
		expect(()=>component.agInit({value: "testquest"} as ICellRendererParams)).not.toThrow();
		expect(component.value).toBe("testquest");
	});

	it("refreshes", () => {
		expect(component.refresh({value: "testquest"} as ICellRendererParams)).toBeTrue();
	});

});
