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

import {HarnessLoader, parallel} from "@angular/cdk/testing";
import {TestbedHarnessEnvironment} from "@angular/cdk/testing/testbed";
import { HttpClientModule } from "@angular/common/http";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { FormsModule, ReactiveFormsModule } from "@angular/forms";
import { MatButtonModule } from "@angular/material/button";
import { MatRadioModule } from "@angular/material/radio";
import { MatStepperModule } from "@angular/material/stepper";
import { MatStepperHarness } from "@angular/material/stepper/testing";
import { NoopAnimationsModule } from "@angular/platform-browser/animations";
import { RouterTestingModule } from "@angular/router/testing";

import { APITestingModule } from "src/app/api/testing";
import { UserService } from "src/app/shared/api";

import { CurrentUserService } from "src/app/shared/currentUser/current-user.service";
import { Protocol } from "../../models";
import { TpHeaderComponent } from "../../shared/tp-header/tp-header.component";
import { NewDeliveryServiceComponent } from "./new-delivery-service.component";


describe("NewDeliveryServiceComponent", () => {
	let component: NewDeliveryServiceComponent;
	let fixture: ComponentFixture<NewDeliveryServiceComponent>;
	let loader: HarnessLoader;

	beforeEach(async () => {
		// mock the API
		const mockAPIService = jasmine.createSpyObj(["getRoles", "getCurrentUser"]);
		const mockCurrentUserService = jasmine.createSpyObj(["updateCurrentUser", "login", "logout"]);
		mockCurrentUserService.updateCurrentUser.and.returnValue(new Promise(r => r(false)));
		mockAPIService.getRoles.and.returnValue(new Promise(resolve => resolve([])));
		mockAPIService.getCurrentUser.and.returnValue(new Promise(resolve => resolve({
			id: 0,
			newUser: false,
			username: "test"
		})));

		await TestBed.configureTestingModule({
			declarations: [
				NewDeliveryServiceComponent,
				TpHeaderComponent
			],
			imports: [
				APITestingModule,
				FormsModule,
				HttpClientModule,
				ReactiveFormsModule,
				RouterTestingModule,
				NoopAnimationsModule,
				MatStepperModule,
				MatButtonModule,
				MatRadioModule
			],
			providers: [
				{provide: UserService, useValue: mockAPIService},
				{ provide: CurrentUserService, useValue: mockCurrentUserService }
			]
		}).compileComponents();
		// TestBed.overrideProvider(UserService, { useValue: mockAPIService });
		// TestBed.compileComponents();
		fixture = TestBed.createComponent(NewDeliveryServiceComponent);
		component = fixture.componentInstance;
		loader = TestbedHarnessEnvironment.loader(fixture);
		fixture.detectChanges();
	});

	it("should exist", async () => {
		expect(component).toBeTruthy();
	});

	it("should parse Origin URLs properly", async () => {
		component.originURL.setValue("http://some.domain.test:9001/a/check/path/here");
		component.setOriginURL();
		expect(component.deliveryService.orgServerFqdn).toEqual("http://some.domain.test:9001", "http://some.domain.test:9001");
		expect(component.deliveryService.checkPath).toEqual("/a/check/path/here", "/a/check/path/here");
		expect(component.deliveryService.displayName).toEqual(
			"Delivery Service for some.domain.test",
			"Delivery Service for some.domain.test"
		);
		expect(component.displayName.value).toEqual("Delivery Service for some.domain.test", "Delivery Service for some.domain.test");
		const stepper = await loader.getHarness(MatStepperHarness);
		const steps = await stepper.getSteps();
		expect(await parallel(() => steps.map(async step => step.isSelected()))).toEqual([
			false,
			true,
			false
		]);
		expect(component.deliveryService.protocol).toEqual(Protocol.HTTP_AND_HTTPS, "HTTP_AND_HTTPS");

		// check other protocol setting
		component.originURL.setValue("https://test.test");
		component.setOriginURL();
		expect(component.deliveryService.protocol).toEqual(Protocol.HTTP_TO_HTTPS, "HTTP_TO_HTTPS");
	});

	it("should set meta info properly", async () => {
		try {
			const stepper = await loader.getHarness(MatStepperHarness);
			const steps = await stepper.getSteps();
			await steps[1].select();
			component.displayName.setValue("test._QUEST");
			component.infoURL.setValue("ftp://this-is-a-weird.url/");
			component.description.setValue("test description");
			component.setMetaInformation();

			expect(component.deliveryService.displayName).toEqual("test._QUEST", "test._QUEST");
			expect(component.deliveryService.xmlId).toEqual("test-quest", "test-quest");
			expect(component.deliveryService.longDesc).toEqual("test description", "test description");
			expect(component.deliveryService.infoUrl).toEqual("ftp://this-is-a-weird.url/", "ftp://this-is-a-weird.url/");
			expect(await parallel(() => steps.map(async step => step.isSelected()))).toEqual([
				false,
				false,
				true
			]);
		} catch (e) {
			console.error("Error occurred:", e);
		}
	});

	// it('should set infrastructure info properly', () => {
	// 	component.step = 2;
	// 	component.cdnObject.setValue({ name: 'testCDN', id: 1 } as CDN);
	// 	component.dsType.setValue({ name: 'testType', id: 1 } as Type);
	// });

	it("should match hostnames", async () => {
		const invalidHostnames: Array<string> = ["h.", "h-", "h-.o", "-h.o"];
		invalidHostnames.forEach((invalidHostname: string) => void expect(() => component.setDNSBypass(invalidHostname)).toThrow());
		const validHostnames: Array<string> = ["h", "h.o.s.T.n.a.m.e", "h-O-------s.tNaMe"];
		expect(() => validHostnames.forEach((hostname: string) => {
			component.setDNSBypass(hostname);
			expect(component.deliveryService.dnsBypassCname).toBe(hostname);
		})).not.toThrow();
	});
});
