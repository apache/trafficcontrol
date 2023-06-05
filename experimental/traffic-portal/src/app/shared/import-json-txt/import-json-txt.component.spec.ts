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

import { DatePipe } from "@angular/common";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { MatDialogModule, MatDialogRef } from "@angular/material/dialog";

import { ImportJsonTxtComponent } from "./import-json-txt.component";

describe("ImportJsonTxtComponent", () => {
	let component: ImportJsonTxtComponent;
	let fixture: ComponentFixture<ImportJsonTxtComponent>;

	beforeEach(async () => {
		await TestBed.configureTestingModule({
			declarations: [ ImportJsonTxtComponent ],
			imports: [
				MatDialogModule
			],
			providers: [
				DatePipe,
				{provide: MatDialogRef, useValue: {}}
			]
		}).compileComponents();

		fixture = TestBed.createComponent(ImportJsonTxtComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});
});
