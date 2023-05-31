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
import { type ComponentFixture, TestBed, fakeAsync, tick } from "@angular/core/testing";
import { MatDialog } from "@angular/material/dialog";
import { Router } from "@angular/router";
import { RouterTestingModule } from "@angular/router/testing";
import { Observable, of, ReplaySubject } from "rxjs";
import type { ResponseServer } from "trafficops-types";

import { ServerService } from "src/app/api";
import { APITestingModule } from "src/app/api/testing";
import { CurrentUserService } from "src/app/shared/current-user/current-user.service";
import { NavigationService } from "src/app/shared/navigation/navigation.service";
import { TpHeaderComponent } from "src/app/shared/navigation/tp-header/tp-header.component";

import { augment, type AugmentedServer, serverIsCache, ServersTableComponent } from "./servers-table.component";

/**
 * Define the MockDialog
 */
class MockDialog {

	/**
	 * Fake opens the dialog
	 *
	 * @returns unknown
	 */
	public open(): unknown {
		return {
			afterClosed: (): Observable<boolean> => of(true)
		};
	}
}

const defaultServer: ResponseServer = {
	cachegroup: "",
	cachegroupId: 1,
	cdnId: 1,
	cdnName: "",
	domainName: "",
	guid: null,
	hostName: "",
	httpsPort: null,
	id: -1,
	iloIpAddress: null,
	iloIpGateway: null,
	iloIpNetmask: null,
	iloPassword: null,
	iloUsername: null,
	interfaces: [
		{
			ipAddresses: [
				{
					address: "1.2.3.4",
					gateway: null,
					serviceAddress: true
				}
			],
			maxBandwidth: null,
			monitor: false,
			mtu: null,
			name: "eth0"
		}
	],
	lastUpdated: new Date(),
	mgmtIpAddress: null,
	mgmtIpGateway: null,
	mgmtIpNetmask: null,
	offlineReason: null,
	physLocation: "",
	physLocationId: 1,
	profileNames: ["GLOBAL"],
	rack: null,
	revalPending: false,
	routerHostName: null,
	routerPortName: null,
	status: "",
	statusId: 1,
	statusLastUpdated: null,
	tcpPort: null,
	type: "",
	typeId: 1,
	updPending: false,
	xmppId: "",
};

describe("ServersTableComponent", () => {
	let component: ServersTableComponent;
	let fixture: ComponentFixture<ServersTableComponent>;
	let router: Router;

	beforeEach(() => {
		const mockCurrentUserService = jasmine.createSpyObj(["updateCurrentUser", "hasPermission", "login", "logout"]);

		const navSvc = jasmine.createSpyObj([],{headerHidden: new ReplaySubject<boolean>(), headerTitle: new ReplaySubject<string>()});
		TestBed.configureTestingModule({
			declarations: [ ServersTableComponent, TpHeaderComponent ],
			imports: [
				HttpClientModule,
				RouterTestingModule.withRoutes([
					{component: ServersTableComponent, path: ""},
					{component: ServersTableComponent, path: "core/server/:id"}
				]),
				APITestingModule
			],
			providers: [
				{ provide: CurrentUserService, useValue: mockCurrentUserService },
				{ provide: MatDialog, useClass: MockDialog },
				{ provide: NavigationService, useValue: navSvc}
			]
		}).compileComponents();
		fixture = TestBed.createComponent(ServersTableComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
		router = TestBed.inject(Router);
		router.initialNavigation();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});

	it("knows if a server is a cache", () => {
		const s: AugmentedServer = {...defaultServer, ipv4Address: "", ipv6Address: "", type: ""};
		expect(serverIsCache(s)).toBeFalse();
		s.type = "EDGE";
		expect(serverIsCache(s)).toBeTrue();
		s.type = "EDGE_anything";
		expect(serverIsCache(s)).toBeTrue();
		s.type = "MID";
		expect(serverIsCache(s)).toBeTrue();
		s.type = "MID_anything";
		expect(serverIsCache(s)).toBeTrue();
		s.type = "a string that merely CONTAINS 'EDGE' instead of starting with it";
		expect(serverIsCache(s)).toBeFalse();
		s.type = "RASCAL";
		expect(serverIsCache(s)).toBeFalse();
	});

	it("augments servers", () => {
		const s: ResponseServer = {...defaultServer, interfaces: []};
		let a = augment(s);
		expect(a.ipv4Address).toBe("");
		expect(a.ipv6Address).toBe("");

		s.interfaces.push({
			ipAddresses: [],
			maxBandwidth: null,
			monitor: false,
			mtu: null,
			name: "test"
		});
		a = augment(s);
		expect(a.ipv4Address).toBe("");
		expect(a.ipv6Address).toBe("");

		s.interfaces[0].ipAddresses.push({
			address: "192.0.2.0",
			gateway: null,
			serviceAddress: false
		});
		a = augment(s);
		expect(a.ipv4Address).toBe("");
		expect(a.ipv6Address).toBe("");

		s.interfaces[0].ipAddresses.push({
			address: "2001::1",
			gateway: null,
			serviceAddress: false
		});
		a = augment(s);
		expect(a.ipv4Address).toBe("");
		expect(a.ipv6Address).toBe("");

		s.interfaces.push({
			ipAddresses: [
				{
					address: "192.0.2.1",
					gateway: null,
					serviceAddress: false
				},
				{
					address: "2001::2",
					gateway: null,
					serviceAddress: false
				}
			],
			maxBandwidth: null,
			monitor: true,
			mtu: null,
			name: "quest"
		});
		a = augment(s);
		expect(a.ipv4Address).toBe("");
		expect(a.ipv6Address).toBe("");

		s.interfaces[1].ipAddresses.push({
			address: "192.0.2.2",
			gateway: null,
			serviceAddress: true
		});
		a = augment(s);
		expect(a.ipv4Address).toBe("192.0.2.2");
		expect(a.ipv6Address).toBe("");

		s.interfaces[1].ipAddresses.push({
			address: "2001::3",
			gateway: null,
			serviceAddress: true
		});
		a = augment(s);
		expect(a.ipv4Address).toBe("192.0.2.2");
		expect(a.ipv6Address).toBe("2001::3");
	});

	it("loads the 'search' query string parameter as the text for the fuzzy search box", fakeAsync(() => {
		expect(component.fuzzControl.value).toBe("");
		router.navigate(["/"], {queryParams: {search: "testquest"}});
		component.ngOnInit();
		tick();
		expect(component.fuzzControl.value).toBe("testquest");
	}));

	it("propagates changes to the search box to the subscription input of the generic table", () => {
		const spy = jasmine.createSpy("fuzzySubscription", (v) => {
			expect(v).toBe("testquest");
		});
		component.fuzzySubject.subscribe(spy);
		component.fuzzControl.setValue("testquest");
		component.updateURL();
		expect(spy).toHaveBeenCalled();
	});

	it("reloads servers when one or more servers' statuses are updated", async () => {
		const service = TestBed.inject(ServerService);
		await service.createServer({...defaultServer, interfaces: []});
		component.ngOnInit();
		const servers = await component.servers;
		if (!servers) {
			return fail("servers table has no servers even though I just created one");
		}

		component.servers = new Promise(r=>r([]));
		expect((await component.servers).length).toBe(0);

		await component.reloadServers();
		expect((await component.servers).length).toBeGreaterThan(0);
	});

	it("handles its context menu actions", fakeAsync(async () => {
		const augmentFields = {ipv4Address: "192.0.2.0", ipv6Address: "2001::1"};
		const server = {...defaultServer, id: 9001, type: "EDGE", ...augmentFields};

		await expectAsync(component.handleContextMenu({action: "viewDetails", data: [server]})).toBeRejected();
		await expectAsync(component.handleContextMenu({action: "viewDetails", data: server})).toBeRejected();

		component.handleContextMenu({action: "updateStatus", data: server});

		component.handleContextMenu({action: "updateStatus", data: [server]});

		const service = TestBed.inject(ServerService);
		const queueSpy = spyOn(service, "queueUpdates");
		const clearSpy = spyOn(service, "clearUpdates");
		expect(queueSpy).not.toHaveBeenCalled();
		expect(clearSpy).not.toHaveBeenCalled();

		component.handleContextMenu({action: "queue", data: server});
		expect(queueSpy).toHaveBeenCalledTimes(1);
		expect(clearSpy).not.toHaveBeenCalled();

		component.handleContextMenu({action: "queue", data: [server, server]});
		expect(queueSpy).toHaveBeenCalledTimes(3);

		component.handleContextMenu({action: "dequeue", data: server});
		expect(queueSpy).toHaveBeenCalledTimes(3);
		expect(clearSpy).toHaveBeenCalledTimes(1);

		component.handleContextMenu({action: "dequeue", data: [server, server]});
		expect(clearSpy).toHaveBeenCalledTimes(3);

		expectAsync(component.handleContextMenu({action: "not a real action", data: []})).toBeRejected();
	}));
});
