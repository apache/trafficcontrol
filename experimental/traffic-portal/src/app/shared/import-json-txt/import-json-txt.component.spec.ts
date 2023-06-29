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
import { MAT_DIALOG_DATA, MatDialogModule, MatDialogRef } from "@angular/material/dialog";

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
				{provide: MatDialogRef, useValue: {}},
				{provide: MAT_DIALOG_DATA, useValue: { title: ""}}
			]
		}).compileComponents();

		fixture = TestBed.createComponent(ImportJsonTxtComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});

	it("should set dragOn to true when dragover event occurs", () => {
		const event = new DragEvent("dragover");
		const preventDefaultSpy = spyOn(event, "preventDefault");
		const stopPropagationSpy = spyOn(event, "stopPropagation");

		fixture.nativeElement.dispatchEvent(event);
		fixture.detectChanges();

		expect(preventDefaultSpy).toHaveBeenCalled();
		expect(stopPropagationSpy).toHaveBeenCalled();
		expect(component.dragOn).toBeTrue();
	});

	it("should set dragOn to true when dragover event occurs", () => {

		const event = new DragEvent("dragleave");
		const preventDefaultSpy = spyOn(event, "preventDefault");
		const stopPropagationSpy = spyOn(event, "stopPropagation");

		fixture.nativeElement.dispatchEvent(event);
		fixture.detectChanges();

		expect(preventDefaultSpy).toHaveBeenCalled();
		expect(stopPropagationSpy).toHaveBeenCalled();
		expect(component.dragOn).toBeFalse();
	});
});
