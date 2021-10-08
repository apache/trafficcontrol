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

import {InvalidationJobService} from "../../../shared/api";
import { NewInvalidationJobDialogComponent } from "./new-invalidation-job-dialog.component";

describe("NewInvalidationJobDialogComponent", () => {
	let component: NewInvalidationJobDialogComponent;
	let fixture: ComponentFixture<NewInvalidationJobDialogComponent>;

	beforeEach(async () => {
		const mockAPIService = jasmine.createSpyObj(["getInvalidationJobs"]);
		await TestBed.configureTestingModule({
			declarations: [ NewInvalidationJobDialogComponent ],
			imports: [
				MatDialogModule,
				HttpClientModule
			],
			providers: [
				{provide: MatDialogRef, useValue: {close: (): void => {
					console.log("dialog closed");
				}}},
				{provide: MAT_DIALOG_DATA, useValue: -1},
				{ provide: InvalidationJobService, useValue: mockAPIService}
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
