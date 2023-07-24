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

import { HarnessLoader } from "@angular/cdk/testing";
import { TestbedHarnessEnvironment } from "@angular/cdk/testing/testbed";
import { HttpClientModule } from "@angular/common/http";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { ReactiveFormsModule } from "@angular/forms";
import { MatButtonHarness } from "@angular/material/button/testing";
import { MatDialogModule } from "@angular/material/dialog";
import { MatDialogHarness } from "@angular/material/dialog/testing";
import { MatSelectModule } from "@angular/material/select";
import { NoopAnimationsModule } from "@angular/platform-browser/animations";
import { RouterTestingModule } from "@angular/router/testing";
import { ReplaySubject } from "rxjs";
import { ResponseCDN } from "trafficops-types";

import { CDNService } from "src/app/api";
import { APITestingModule } from "src/app/api/testing";
import { NavigationService } from "src/app/shared/navigation/navigation.service";

import { CDNTableComponent } from "./cdn-table.component";

const sampleCDN: ResponseCDN = {
	dnssecEnabled: false,
	domainName: "*",
	id: 2,
	lastUpdated: new Date(),
	name: "*",
};

describe("CDNTableComponent", () => {
	let component: CDNTableComponent;
	let fixture: ComponentFixture<CDNTableComponent>;
	let loader: HarnessLoader;

	const navService = jasmine.createSpyObj([],{headerHidden: new ReplaySubject<boolean>(), headerTitle: new ReplaySubject<string>()});

	beforeEach(async () => {
		await TestBed.configureTestingModule({
			declarations: [CDNTableComponent],
			imports: [
				APITestingModule,
				HttpClientModule,
				ReactiveFormsModule,
				RouterTestingModule,
				MatDialogModule,
				NoopAnimationsModule,
				MatSelectModule,
			],
			providers: [
				{provide: NavigationService, useValue: navService},
			],
		})
			.compileComponents();

		fixture = TestBed.createComponent(CDNTableComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
		loader = TestbedHarnessEnvironment.documentRootLoader(fixture);
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});

	it("queues CDN updates", async () => {
		component = fixture.componentInstance;
		const service = TestBed.inject(CDNService);
		const queueSpy = spyOn(service, "queueServerUpdates");
		expect(queueSpy).not.toHaveBeenCalled();

		let dialogs = await loader.getAllHarnesses(MatDialogHarness);
		expect(dialogs.length).toBe(0);

		component.handleContextMenu({action: "queue", data: sampleCDN});
		dialogs = await loader.getAllHarnesses(MatDialogHarness);
		expect(dialogs.length).toBe(1);
		const dialog = dialogs[0];
		const buttons = await dialog.getAllHarnesses(MatButtonHarness);
		expect(buttons.length).toBe(2);
		const button = buttons[0];
		await button.click();

		expect(queueSpy).toHaveBeenCalledTimes(1);
	});

});
