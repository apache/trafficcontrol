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

import { UpdateCellRendererComponent } from "./update-cell-renderer.component";

describe("UpdateCellRendererComponent", () => {
	let component: UpdateCellRendererComponent;
	let fixture: ComponentFixture<UpdateCellRendererComponent>;

	beforeEach(async () => {
		await TestBed.configureTestingModule({
			declarations: [ UpdateCellRendererComponent ]
		})
			.compileComponents();
	});

	beforeEach(() => {
		fixture = TestBed.createComponent(UpdateCellRendererComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});

	it("initializes", () => {
		component.agInit({value: true} as ICellRendererParams);
		expect(component.value).toBeTrue();

		component.agInit({value: false} as ICellRendererParams);
		expect(component.value).toBeFalse();
	});

	it("refreshes", () => {
		let ret = component.refresh({value: true} as ICellRendererParams);
		expect(ret).toBeTrue();
		expect(component.value).toBeTrue();

		ret = component.refresh({value: false} as ICellRendererParams);
		expect(ret).toBeTrue();
		expect(component.value).toBeFalse();
	});
});
