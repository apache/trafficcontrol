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
import { HarnessLoader } from "@angular/cdk/testing";
import { TestbedHarnessEnvironment } from "@angular/cdk/testing/testbed";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { MatButtonHarness } from "@angular/material/button/testing";
import { MatDialogHarness } from "@angular/material/dialog/testing";
import { MatSelectHarness } from "@angular/material/select/testing";
import { NoopAnimationsModule } from "@angular/platform-browser/animations";
import type { ResponseServer } from "trafficops-types";

import {
	CDNService,
	CacheGroupService,
	MiscAPIsService,
	PhysicalLocationService,
	ProfileService,
	ServerService,
	TypeService
} from "src/app/api";
import { APITestingModule } from "src/app/api/testing";
import { FileUtilsService } from "src/app/shared/file-utils.service";
import { SharedModule } from "src/app/shared/shared.module";

import { ISOGenerationFormComponent } from "./isogeneration-form.component";

describe("ISOGenerationFormComponent", () => {
	let component: ISOGenerationFormComponent;
	let fixture: ComponentFixture<ISOGenerationFormComponent>;
	let form: typeof component.form.controls;
	let loader: HarnessLoader;
	let spy: jasmine.Spy;

	let server: ResponseServer;

	const ipv4 = {
		address: "192.0.0.1/16",
		gateway: "192.0.0.2",
		serviceAddress: true,
	};

	const ipv6 = {
		address: "::dead:beef",
		gateway: "::f1d0:f00d",
		serviceAddress: true,
	};

	const iface = {
		ipAddresses: [ipv4, ipv6],
		maxBandwidth: null,
		monitor: false,
		mtu: 2000,
		name: "eth0",
	};

	beforeEach(async () => {
		await TestBed.configureTestingModule({
			declarations: [ ISOGenerationFormComponent ],
			imports: [
				APITestingModule,
				SharedModule,
				NoopAnimationsModule
			],
			providers: [
				{
					provide: "Window",
					useValue: {
						open: (): void => {
							// do nothing
						}
					}
				},
				{
					provide: FileUtilsService,
					useValue: {
						download: (): void => {
							// do nothing
						}
					}
				}
			]
		}).compileComponents();

		fixture = TestBed.createComponent(ISOGenerationFormComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
		form = component.form.controls;
		const srv = TestBed.inject(MiscAPIsService);
		spy = spyOn(srv, "generateISO").and.callThrough();
		loader = TestbedHarnessEnvironment.documentRootLoader(fixture);

		const serverSrv = TestBed.inject(ServerService);

		const serverType = (await TestBed.inject(TypeService).getServerTypes())[0];
		if (!serverType) {
			return fail("no server types available");
		}
		const cg = (await TestBed.inject(CacheGroupService).getCacheGroups())[0];
		if (!cg) {
			return fail("no cache groups available");
		}
		const cdn = (await TestBed.inject(CDNService).getCDNs())[0];
		if (!cdn) {
			return fail("no cdns available");
		}
		const physLoc = (await TestBed.inject(PhysicalLocationService).getPhysicalLocations())[0];
		if (!physLoc) {
			return fail("no physical locations available");
		}
		const profile = (await TestBed.inject(ProfileService).getProfiles())[0];
		if (!profile) {
			return fail("no profiles available");
		}
		const status = (await serverSrv.getStatuses())[0];
		if (!status) {
			return fail("no statuses available");
		}

		server = await serverSrv.createServer({
			cachegroupId: cg.id,
			cdnId: cdn.id,
			domainName: "test",
			hostName: "quest",
			interfaces: [iface],
			physLocationId: physLoc.id,
			profileNames: [profile.name],
			statusId: status.id,
			typeId: serverType.id
		});
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
			mgmtInterface: "mgmt0",
			mgmtIpAddress: "1.3.3.7",
			mgmtIpGateway: "9.0.0.1",
			mgmtIpNetmask: "0.4.2.0",
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
		form.mgmtInterface.setValue(request.mgmtInterface);
		form.mgmtIpAddress.setValue(request.mgmtIpAddress);
		form.mgmtIpGateway.setValue(request.mgmtIpGateway);
		form.mgmtIpNetmask.setValue(request.mgmtIpNetmask);
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
			mgmtInterface: "mgmt0",
			mgmtIpAddress: "1.3.3.7",
			mgmtIpGateway: "9.0.0.1",
			mgmtIpNetmask: "0.4.2.0",
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
		form.mgmtInterface.setValue(request.mgmtInterface);
		form.mgmtIpAddress.setValue(request.mgmtIpAddress);
		form.mgmtIpGateway.setValue(request.mgmtIpGateway);
		form.mgmtIpNetmask.setValue(request.mgmtIpNetmask);
		form.osVersion.setValue(request.osVersionDir);
		form.rootPass.setValue(request.rootPass);
		form.rootPassConfirm.setValue(request.rootPass);

		await component.submit(new Event("submit"));
		expect(spy).toHaveBeenCalledOnceWith(request);
	});

	it("opens the copy dialog", async () => {
		const asyncExpectation = expectAsync(component.openCopyDialog()).toBeResolvedTo(undefined);
		const dialogs = await loader.getAllHarnesses(MatDialogHarness);
		expect(dialogs.length).toBe(1);
		dialogs[0].close();
		await asyncExpectation;

		expect(spy).not.toHaveBeenCalled();
	});

	it("copies server attributes", async () => {
		const srv = TestBed.inject(ServerService);
		expect(await srv.getServers()).toHaveSize(1);

		const asyncExpectation = expectAsync(component.openCopyDialog()).toBeResolvedTo(undefined);

		const dialogs = await loader.getAllHarnesses(MatDialogHarness);
		if (dialogs.length !== 1) {
			return fail(`exactly one dialog should exist; got: ${dialogs.length}`);
		}
		const dialog = dialogs[0];
		const selects = await dialog.getAllHarnesses(MatSelectHarness);
		if (selects.length !== 1) {
			return fail(`dialog should have contained one select input, got: ${selects.length}`);
		}
		const select = selects[0];
		await select.clickOptions();
		const buttons = await dialog.getAllHarnesses(MatButtonHarness.with({text: /^[cC][oO][nN][fF][iI][rR][mM]$/}));
		if (buttons.length !== 1) {
			return fail(`'Confirm' button not found; expected one, got: ${buttons.length}`);
		}
		await buttons[0].click();

		await asyncExpectation;

		expect(form.useDHCP.value).toBeFalse();
		expect(form.fqdn.value).toBe(`${server.hostName}.${server.domainName}`);
		expect(form.interfaceName.value).toBe(iface.name);
		expect(form.ipv4Address.value).toBe(ipv4.address.split("/")[0]);
		expect(form.ipv4Gateway.value).toBe(ipv4.gateway);
		// Ideally this wouldn't be hard-coded, but calculating it is a pain.
		expect(form.ipv4Netmask.value).toBe("255.255.0.0");
		expect(form.ipv6Address.value).toBe(ipv6.address);
		expect(form.ipv6Gateway.value).toBe(ipv6.gateway);
		expect(form.mtu.value).toBe(iface.mtu);
	});
});
