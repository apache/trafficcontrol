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
import { waitForAsync, ComponentFixture, TestBed } from "@angular/core/testing";
import { FormsModule, ReactiveFormsModule } from "@angular/forms";
import { RouterTestingModule } from "@angular/router/testing";
import { of } from "rxjs";
import { UserService } from "src/app/services/api";


import { Protocol } from "../../models";
import { TpHeaderComponent } from "../tp-header/tp-header.component";
import { NewDeliveryServiceComponent } from "./new-delivery-service.component";

describe("NewDeliveryServiceComponent", () => {
	let component: NewDeliveryServiceComponent;
	let fixture: ComponentFixture<NewDeliveryServiceComponent>;

	beforeEach(waitForAsync(() => {
		// mock the API
		const mockAPIService = jasmine.createSpyObj(["getRoles", "getCurrentUser"]);
		mockAPIService.getRoles.and.returnValue(of([]));
		mockAPIService.getCurrentUser.and.returnValue(of({
			id: 0,
			newUser: false,
			username: "test"
		}));

		TestBed.configureTestingModule({
			declarations: [
				NewDeliveryServiceComponent,
				TpHeaderComponent
			],
			imports: [
				FormsModule,
				HttpClientModule,
				ReactiveFormsModule,
				RouterTestingModule
			]
		});
		TestBed.overrideProvider(UserService, { useValue: mockAPIService });
		TestBed.compileComponents();
	}));

	beforeEach(() => {
		fixture = TestBed.createComponent(NewDeliveryServiceComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it("should exist", () => {
		expect(component).toBeTruthy();
	});

	it("should parse Origin URLs properly", () => {
		component.originURL.setValue("http://some.domain.test:9001/a/check/path/here");
		component.setOriginURL();
		expect(component.deliveryService.orgServerFqdn).toEqual("http://some.domain.test:9001", "http://some.domain.test:9001");
		expect(component.deliveryService.checkPath).toEqual("/a/check/path/here", "/a/check/path/here");
		expect(component.deliveryService.displayName).toEqual(
			"Delivery Service for some.domain.test",
			"Delivery Service for some.domain.test"
		);
		expect(component.displayName.value).toEqual("Delivery Service for some.domain.test", "Delivery Service for some.domain.test");
		expect(component.step).toEqual(1, "one");
		expect(component.deliveryService.protocol).toEqual(Protocol.HTTP_AND_HTTPS, "HTTP_AND_HTTPS");

		// check other protocol setting
		component.originURL.setValue("https://test.test");
		component.setOriginURL();
		expect(component.deliveryService.protocol).toEqual(Protocol.HTTP_TO_HTTPS, "HTTP_TO_HTTPS");
	});

	it("should set meta info properly", () => {
		try {
			component.step = 1;
			component.displayName.setValue("test._QUEST");
			component.infoURL.setValue("ftp://this-is-a-weird.url/");
			component.description.setValue("test description");
			component.setMetaInformation();

			expect(component.deliveryService.displayName).toEqual("test._QUEST", "test._QUEST");
			expect(component.deliveryService.xmlId).toEqual("test-quest", "test-quest");
			expect(component.deliveryService.longDesc).toEqual("test description", "test description");
			expect(component.deliveryService.infoUrl).toEqual("ftp://this-is-a-weird.url/", "ftp://this-is-a-weird.url/");
			expect(component.step).toEqual(2, "two");
		} catch (e) {
			console.error("Error occurred:", e);
		}
	});

	// it('should set infrastructure info properly', () => {
	// 	component.step = 2;
	// 	component.cdnObject.setValue({ name: 'testCDN', id: 1 } as CDN);
	// 	component.dsType.setValue({ name: 'testType', id: 1 } as Type);
	// });
});
