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
import { RouterTestingModule } from "@angular/router/testing";

import { CurrentUserService } from "src/app/shared/currentUser/current-user.service";
import { TpHeaderComponent } from "src/app/shared/tp-header/tp-header.component";
import { UserService } from "src/app/shared/api";
import { defaultServer, Server } from "src/app/models";
import { APITestingModule } from "src/app/api/testing";
import { augment, AugmentedServer, serverIsCache, ServersTableComponent } from "./servers-table.component";


describe("ServersTableComponent", () => {
	let component: ServersTableComponent;
	let fixture: ComponentFixture<ServersTableComponent>;

	beforeEach(waitForAsync(() => {
		const mockAPIService = jasmine.createSpyObj(["getServers", "getUsers"]);
		mockAPIService.getServers.and.returnValue(new Promise(r => r([])));
		const mockCurrentUserService = jasmine.createSpyObj(["updateCurrentUser", "login", "logout"]);
		TestBed.configureTestingModule({
			declarations: [ ServersTableComponent, TpHeaderComponent ],
			imports: [HttpClientModule, RouterTestingModule, APITestingModule],
			providers: [
				{ provide: UserService, useValue: mockAPIService },
				{ provide: CurrentUserService, useValue: mockCurrentUserService }
			]
		})
			.compileComponents();
	}));

	beforeEach(() => {
		fixture = TestBed.createComponent(ServersTableComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});

	it("knows if a server is a cache", () => {
		const s: AugmentedServer = {...defaultServer, ipv4Address: "", ipv6Address: "", type: undefined};
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
		const s: Server = {...defaultServer, interfaces: []};
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
});
