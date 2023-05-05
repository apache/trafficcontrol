/**
 * @license Apache-2.0
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
import { ComponentFixture, TestBed } from "@angular/core/testing";

import { MiscAPIsService } from "src/app/api";
import { APITestingModule } from "src/app/api/testing";
import { SharedModule } from "src/app/shared/shared.module";

import { ISOGenerationFormComponent } from "./isogeneration-form.component";

describe("ISOGenerationFormComponent", () => {
	let component: ISOGenerationFormComponent;
	let fixture: ComponentFixture<ISOGenerationFormComponent>;
	let form: typeof component.form.controls;
	let spy: jasmine.Spy;

	beforeEach(async () => {
		await TestBed.configureTestingModule({
			declarations: [ ISOGenerationFormComponent ],
			imports: [
				APITestingModule,
				SharedModule
			],
			providers: [{
				provide: "Window",
				useValue: {
					open: (): void => {
						// do nothing
					}
				}
			}]
		}).compileComponents();

		fixture = TestBed.createComponent(ISOGenerationFormComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
		form = component.form.controls;
		const srv = TestBed.inject(MiscAPIsService);
		spy = spyOn(srv, "generateISO").and.callThrough();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});

	it("hides the MTU warning when appropriate", () => {
		// should be hidden by default
		expect(component.hideMTUWarning()).toBeTrue();
		form.mtu.setValue(1501);
		expect(component.hideMTUWarning()).toBeFalse();
		form.mtu.setValue(1500);
		expect(component.hideMTUWarning()).toBeTrue();
		form.mtu.setValue(9000);
		expect(component.hideMTUWarning()).toBeTrue();
	});

	it("validates that the root password matches the confirm field", () => {
		form.rootPass.setValue("testquest");
		form.rootPassConfirm.setValue("testquest");
		expect(form.rootPass.valid).toBeTrue();
		expect(form.rootPassConfirm.valid).toBeTrue();

		form.rootPassConfirm.setValue(`${form.rootPassConfirm.value} some more stuff`);
		expect(form.rootPass.valid).toBeTrue();
		expect(form.rootPassConfirm.invalid).toBeTrue();
		if (!form.rootPassConfirm.errors) {
			return fail("rootPassConfirm had null errors when it should be invalid");
		}
		expect(form.rootPassConfirm.errors.mismatch).toBeTrue();
	});

	it("doesn't submit requests when the form is invalid", async () => {
		form.rootPass.setValue("something");
		form.rootPassConfirm.setValue("something else");
		await component.submit(new Event("submit"));
		expect(spy).not.toHaveBeenCalled();
	});

	it("submits requests for DHCP ISOs", async () => {
		const request = {
			dhcp: "yes",
			disk: "sda",
			domainName: "domain-name",
			hostName: "host-name",
			interfaceMtu: 1501,
			interfaceName: "eth0",
			ip6Address: "f1d0::f00d",
			ip6Gateway: "dead::beef",
			osVersionDir: "os version dir",
			rootPass: "root password",
		};
		form.useDHCP.setValue(request.dhcp === "yes");
		form.disk.setValue(request.disk);
		form.fqdn.setValue(`${request.hostName}.${request.domainName}`);
		form.interfaceName.setValue(request.interfaceName);
		form.mtu.setValue(request.interfaceMtu);
		form.ipv4Address.setValue("");
		form.ipv4Gateway.setValue("");
		form.ipv4Netmask.setValue("");
		form.ipv6Address.setValue(request.ip6Address);
		form.ipv6Gateway.setValue(request.ip6Gateway);
		form.osVersion.setValue(request.osVersionDir);
		form.rootPass.setValue(request.rootPass);
		form.rootPassConfirm.setValue(request.rootPass);

		await component.submit(new Event("submit"));
		expect(spy).toHaveBeenCalledOnceWith(request);
	});

	it("submits requests for non-DHCP ISOs", async () => {
		const request = {
			dhcp: "no",
			disk: "sda",
			domainName: "domain-name",
			hostName: "host-name",
			interfaceMtu: 1501,
			interfaceName: "eth0",
			ip6Address: "f1d0::f00d",
			ip6Gateway: "dead::beef",
			ipAddress: "1.2.3.4",
			ipGateway: "4.3.2.1",
			ipNetmask: "0.1.10.100",
			osVersionDir: "os version dir",
			rootPass: "root password",
		};
		form.useDHCP.setValue(request.dhcp === "yes");
		form.disk.setValue(request.disk);
		form.fqdn.setValue(`${request.hostName}.${request.domainName}`);
		form.interfaceName.setValue(request.interfaceName);
		form.mtu.setValue(request.interfaceMtu);
		form.ipv4Address.setValue(request.ipAddress);
		form.ipv4Gateway.setValue(request.ipGateway);
		form.ipv4Netmask.setValue(request.ipNetmask);
		form.ipv6Address.setValue(request.ip6Address);
		form.ipv6Gateway.setValue(request.ip6Gateway);
		form.osVersion.setValue(request.osVersionDir);
		form.rootPass.setValue(request.rootPass);
		form.rootPassConfirm.setValue(request.rootPass);

		await component.submit(new Event("submit"));
		expect(spy).toHaveBeenCalledOnceWith(request);
	});
});
