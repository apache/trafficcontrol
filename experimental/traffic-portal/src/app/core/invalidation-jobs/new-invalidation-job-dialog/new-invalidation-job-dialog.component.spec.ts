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
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { MatDialogModule, MatDialogRef, MAT_DIALOG_DATA } from "@angular/material/dialog";

import { APITestingModule } from "src/app/api/testing";

import { NewInvalidationJobDialogComponent, sanitizedRegExpString, timeStringFromDate } from "./new-invalidation-job-dialog.component";

describe("NewInvalidationJobDialogComponent", () => {
	let component: NewInvalidationJobDialogComponent;
	let fixture: ComponentFixture<NewInvalidationJobDialogComponent>;

	beforeEach(async () => {
		await TestBed.configureTestingModule({
			declarations: [ NewInvalidationJobDialogComponent ],
			imports: [
				MatDialogModule,
				HttpClientModule,
				APITestingModule
			],
			providers: [
				{provide: MatDialogRef, useValue: {close: (): void => {
					console.log("dialog closed");
				}}},
				{provide: MAT_DIALOG_DATA, useValue: {dsID: -1}},
			]
		}).compileComponents();
	});

	beforeEach(() => {
		fixture = TestBed.createComponent(NewInvalidationJobDialogComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});
});

describe("NewInvalidationJobDialogComponent utility functions", () => {
	it("gets a time string from a Date", ()=>{
		const d = new Date();
		d.setHours(0);
		d.setMinutes(0);
		expect(timeStringFromDate(d)).toBe("00:00");
		d.setHours(12);
		d.setMinutes(34);
		expect(timeStringFromDate(d)).toBe("12:34");
	});
	it("sanitizes regular expressions", ()=>{
		expect(sanitizedRegExpString(/\/.+\/my\/path\.jpg/)).toBe("/.+/my/path\\.jpg");
		expect(sanitizedRegExpString(new RegExp("\\/path\\/to\\/content\\/.+\\.m3u8"))).toBe("/path/to/content/.+\\.m3u8");
	});
});
