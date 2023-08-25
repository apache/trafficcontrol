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
import { type HarnessLoader, parallel } from "@angular/cdk/testing";
import { TestbedHarnessEnvironment } from "@angular/cdk/testing/testbed";
import { HttpClientModule } from "@angular/common/http";
import { type ComponentFixture, TestBed } from "@angular/core/testing";
import { FormsModule, ReactiveFormsModule } from "@angular/forms";
import { MatButtonModule } from "@angular/material/button";
import { MatRadioModule } from "@angular/material/radio";
import { MatStepperModule } from "@angular/material/stepper";
import { MatStepperHarness } from "@angular/material/stepper/testing";
import { NoopAnimationsModule } from "@angular/platform-browser/animations";
import { RouterTestingModule } from "@angular/router/testing";
import { ReplaySubject } from "rxjs";
import { Protocol } from "trafficops-types";

import { APITestingModule } from "src/app/api/testing";
import { CurrentUserService } from "src/app/shared/current-user/current-user.service";
import { NavigationService } from "src/app/shared/navigation/navigation.service";
import { TpHeaderComponent } from "src/app/shared/navigation/tp-header/tp-header.component";

import { NewDeliveryServiceComponent } from "./new-delivery-service.component";

describe("NewDeliveryServiceComponent", () => {
	let component: NewDeliveryServiceComponent;
	let fixture: ComponentFixture<NewDeliveryServiceComponent>;
	let loader: HarnessLoader;

	beforeEach(async () => {
		// mock the API
		const mockCurrentUserService = jasmine.createSpyObj(["updateCurrentUser", "hasPermission", "login", "logout"], {currentUser: {
			tenant: "root",
			tenantId: 1
		}});
		mockCurrentUserService.updateCurrentUser.and.returnValue(new Promise(r => r(true)));
		const headerSvc = jasmine.createSpyObj([],{headerHidden: new ReplaySubject<boolean>(), headerTitle: new ReplaySubject<string>()});

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
				{ provide: CurrentUserService, useValue: mockCurrentUserService },
				{ provide: NavigationService, useValue: headerSvc}
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
		component.activeImmediately.setValue(true);
		component.activeImmediately.markAsDirty();
		component.setOriginURL();
		expect(component.deliveryService.active).toBeTrue();
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
		component.activeImmediately.setValue(false);
		component.setOriginURL();
		expect(component.deliveryService.active).toBeFalse();
		expect(component.deliveryService.protocol).toEqual(Protocol.HTTP_TO_HTTPS, "HTTP_TO_HTTPS");
	});

	it("should set meta info properly", async () => {
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
	});

	it("should set infrastructure info properly for HTTP Delivery Services", () => {
		component.cdnObject.setValue({dnssecEnabled: false, domainName: "quest", id: 2, lastUpdated: new Date(), name: "test"});
		component.dsType.setValue({description: "", id: 10, lastUpdated: new Date(), name: "HTTP", useInTable: ""});
		component.dsType.markAsDirty();
		component.protocol.setValue(Protocol.HTTPS);
		component.protocol.markAsDirty();
		component.disableIPv6.setValue(true);
		component.disableIPv6.markAsDirty();
		const bypass = "https://some-other.ds.mycdn.test";
		component.bypassLoc.setValue(bypass);
		component.bypassLoc.markAsDirty();
		component.setInfrastructureInformation();
		expect(component.deliveryService.typeId).toBe(10);
		expect(component.deliveryService.cdnId).toBe(2);
		expect(component.deliveryService.protocol).toBe(Protocol.HTTPS);
		expect(component.deliveryService.ipv6RoutingEnabled).toBeFalse();
		expect(component.deliveryService.httpBypassFqdn).toBe(bypass);
		expect(component.deliveryService.dnsBypassCname).toBeUndefined();
		expect(component.deliveryService.dnsBypassIp6).toBeUndefined();
		expect(component.deliveryService.dnsBypassIp).toBeUndefined();
	});

	it("should set infrastructure info properly for DNS Delivery Services", () => {
		component.cdnObject.setValue({dnssecEnabled: false, domainName: "quest", id: 2, lastUpdated: new Date(), name: "test"});
		component.dsType.setValue({description: "", id: 7, lastUpdated: new Date(), name: "DNS", useInTable: ""});
		component.dsType.markAsDirty();
		component.protocol.setValue(Protocol.HTTP);
		component.protocol.markAsDirty();
		component.disableIPv6.setValue(true);
		component.disableIPv6.markAsDirty();
		let bypass = "some-other.ds.mycdn.test";
		component.bypassLoc.setValue(bypass);
		component.bypassLoc.markAsDirty();
		component.setInfrastructureInformation();
		expect(component.deliveryService.typeId).toBe(7);
		expect(component.deliveryService.cdnId).toBe(2);
		expect(component.deliveryService.protocol).toBe(Protocol.HTTP);
		expect(component.deliveryService.ipv6RoutingEnabled).toBeFalse();
		expect(component.deliveryService.httpBypassFqdn).toBeNull();
		expect(component.deliveryService.dnsBypassCname).toBe(bypass);
		expect(component.deliveryService.dnsBypassIp6).toBeUndefined();
		expect(component.deliveryService.dnsBypassIp).toBeUndefined();

		bypass = "2001::abc";
		component.bypassLoc.setValue(bypass);
		component.setInfrastructureInformation();
		expect(component.deliveryService.dnsBypassIp6).toBe(bypass);

		bypass = "192.0.2.1";
		component.bypassLoc.setValue(bypass);
		component.setInfrastructureInformation();
		expect(component.deliveryService.dnsBypassIp).toBe(bypass);
	});

	it("should match hostnames", async () => {
		const invalidHostnames: Array<string> = ["h.", "h-", "h-.o", "-h.o"];
		for (const invalidHostname of invalidHostnames) {
			expect(() => component.setDNSBypass(invalidHostname)).toThrow();
		}
		const validHostnames: Array<string> = ["h", "h.o.s.T.n.a.m.e", "h-O-------s.tNaMe"];
		expect(() => validHostnames.forEach(
			hostname => {
				component.setDNSBypass(hostname);
				expect(component.deliveryService.dnsBypassCname).toBe(hostname);
			}
		)).not.toThrow();
	});

	it("goes to the previous step", async () => {
		const stepper = await loader.getHarness(MatStepperHarness);
		const steps = await stepper.getSteps();
		expect(await steps[0].isSelected()).toBeTrue();
		await steps[1].select();
		expect(await steps[0].isSelected()).toBeFalse();
		expect(await steps[1].isSelected()).toBeTrue();
		component.previous();
		expect(await steps[0].isSelected()).toBeTrue();
		expect(await steps[1].isSelected()).toBeFalse();
	});
});
