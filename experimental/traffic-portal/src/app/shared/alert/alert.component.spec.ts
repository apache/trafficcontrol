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
import {TestBed, ComponentFixture} from "@angular/core/testing";
import {TestbedHarnessEnvironment} from "@angular/cdk/testing/testbed";
import {MatSnackBarHarness} from "@angular/material/snack-bar/testing";
import {HarnessLoader} from "@angular/cdk/testing";
import {MatSnackBarModule} from "@angular/material/snack-bar";
import {NoopAnimationsModule} from "@angular/platform-browser/animations";

import { AlertService } from "./alert.service";
import { AlertLevel } from "../../models/alert.model";
import { AlertComponent } from "./alert.component";

describe("AlertComponent", () => {
	let component: AlertComponent;
	let fixture: ComponentFixture<AlertComponent>;
	let loader: HarnessLoader;
	let service: AlertService;

	beforeEach(async () => {
		await TestBed.configureTestingModule({
			declarations: [ AlertComponent ],
			imports: [MatSnackBarModule, NoopAnimationsModule],
			providers: [ AlertService ]
		}).compileComponents();
		service = TestBed.inject(AlertService);
		fixture = TestBed.createComponent(AlertComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
		fixture.componentInstance.duration = undefined;
		loader = TestbedHarnessEnvironment.documentRootLoader(fixture);
	});

	it("should exist", () => {
		expect(component).toBeTruthy();
		expect(service).toBeTruthy();
	});

	it("should load simple alerts", async () => {
		const levels: Array<AlertLevel> = ["error", "warning", "info", "success"];
		for (const errLevel of levels) {
			const msg = `An alert at the '${errLevel}' level`;
			service.newAlert(errLevel, msg);

			let snackBars = await loader.getAllHarnesses(MatSnackBarHarness);
			expect(snackBars.length).toBe(1);

			const snackBar = await loader.getHarness(MatSnackBarHarness);
			expect(await snackBar.getMessage()).toBe(msg);

			fixture.componentInstance.clear();
			snackBars = await loader.getAllHarnesses(MatSnackBarHarness);
			expect(snackBars.length).toBe(0);
		}
	});
});
