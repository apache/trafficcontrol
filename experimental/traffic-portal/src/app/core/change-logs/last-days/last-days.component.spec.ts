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
import { MAT_DIALOG_DATA, MatDialogRef } from "@angular/material/dialog";

import { LastDaysComponent } from "./last-days.component";

describe("LastDaysComponent", () => {
	let component: LastDaysComponent;
	let fixture: ComponentFixture<LastDaysComponent>;
	let mockMatDialog: jasmine.SpyObj<MatDialogRef<number>>;

	beforeEach(async () => {
		mockMatDialog = jasmine.createSpyObj("MatDialogRef", ["close", "afterClosed"]);
		await TestBed.configureTestingModule({
			declarations: [LastDaysComponent],
			providers: [{provide: MatDialogRef, useValue: mockMatDialog},
				{provide: MAT_DIALOG_DATA, useValue: (): number => 3}]
		})
			.compileComponents();
	});

	beforeEach(() => {
		fixture = TestBed.createComponent(LastDaysComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});
});
