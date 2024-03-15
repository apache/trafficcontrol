/*
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import { HarnessLoader } from "@angular/cdk/testing";
import { TestbedHarnessEnvironment } from "@angular/cdk/testing/testbed";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { MatCheckboxHarness } from "@angular/material/checkbox/testing";
import { MAT_DIALOG_DATA, MatDialogRef } from "@angular/material/dialog";
import { NoopAnimationsModule } from "@angular/platform-browser/animations";

import { AppUIModule } from "src/app/app.ui.module";
import {
	DownloadOptionsDialogComponent,
	DownloadOptionsDialogData
} from "src/app/shared/generic-table/download-options/download-options-dialog.component";

let loader: HarnessLoader;
describe("DownloadOptionsComponent", () => {
	let component: DownloadOptionsDialogComponent;
	let fixture: ComponentFixture<DownloadOptionsDialogComponent>;
	const data: DownloadOptionsDialogData = {
		allRows: 5,
		columns: [{
			hide: true
		}, {
			hide: false
		}],
		name: "test",
		selectedRows: undefined,
		visibleRows: 5
	};
	const spyRef = jasmine.createSpyObj("MatDialogRef", ["close"]);

	beforeEach(async () => {
		await TestBed.configureTestingModule({
			declarations: [ DownloadOptionsDialogComponent ],
			imports: [
				AppUIModule,
				NoopAnimationsModule
			],
			providers: [
				{provide: MatDialogRef, useValue: spyRef},
				{provide: MAT_DIALOG_DATA, useValue: data}
			]
		}).compileComponents();

		fixture = TestBed.createComponent(DownloadOptionsDialogComponent);
		component = fixture.componentInstance;
		loader = TestbedHarnessEnvironment.loader(fixture);
		fixture.detectChanges();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});

	it("defaults set", async () => {
		expect(fixture.componentInstance.allRows).toEqual(data.allRows);
		expect(fixture.componentInstance.columns).toEqual(data.columns);
		expect(fixture.componentInstance.fileName).toEqual(data.name);
		expect(fixture.componentInstance.selectedRows).toEqual(data.selectedRows);
		expect(fixture.componentInstance.visibleRows).toEqual(data.visibleRows);

		expect(fixture.componentInstance.visibleColumns).toEqual(data.columns.filter(c => !c.hide));
		expect(fixture.componentInstance.columns).toEqual(data.columns);
	});

	it("default submission", async () => {
		const cbs = await loader.getAllHarnesses(MatCheckboxHarness);
		expect(cbs.length).toBe(2);

		fixture.componentInstance.onSubmit();
		expect(spyRef.close.calls.count()).toBe(1);
	});
});
