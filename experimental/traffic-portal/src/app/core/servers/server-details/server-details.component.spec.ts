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
import { type ComponentFixture, fakeAsync, TestBed, tick } from "@angular/core/testing";
import { FormsModule, ReactiveFormsModule } from "@angular/forms";
import { MatDialog, MatDialogRef } from "@angular/material/dialog";
import { MatFormFieldModule } from "@angular/material/form-field";
import { MatInputModule } from "@angular/material/input";
import { MatSelectModule } from "@angular/material/select";
import { BrowserAnimationsModule } from "@angular/platform-browser/animations";
import { RouterTestingModule } from "@angular/router/testing";
import { of } from "rxjs";

import { ServerService } from "src/app/api";
import { APITestingModule } from "src/app/api/testing";
import { CurrentUserService } from "src/app/shared/current-user/current-user.service";
import { SharedModule } from "src/app/shared/shared.module";

import { ServerDetailsComponent } from "./server-details.component";

describe("ServerDetailsComponent", () => {
	let component: ServerDetailsComponent;
	let fixture: ComponentFixture<ServerDetailsComponent>;

	beforeEach(async () => {
		const mockCurrentUserService = jasmine.createSpyObj(
			["updateCurrentUser", "hasPermission", "login", "logout"], {userChanged: of(null)});
		await TestBed.configureTestingModule({
			declarations: [ServerDetailsComponent],
			imports: [
				HttpClientModule,
				RouterTestingModule.withRoutes([
					{component: ServerDetailsComponent, path: "server/:id"},
					{component: ServerDetailsComponent, path: "server/new"},
				]),
				FormsModule,
				ReactiveFormsModule,
				MatSelectModule,
				MatFormFieldModule,
				MatInputModule,
				BrowserAnimationsModule,
				APITestingModule,
				SharedModule
			],
			providers: [
				{provide: CurrentUserService, useValue: mockCurrentUserService},
			]
		}).compileComponents();
		fixture = TestBed.createComponent(ServerDetailsComponent);
		const service = TestBed.inject(ServerService);
		component = fixture.componentInstance;
		component.server = await service.createServer({
			cachegroupId: 1,
			cdnId: 1,
			domainName: "",
			hostName: "",
			httpsPort: null,
			iloIpAddress: null,
			iloIpGateway: null,
			iloIpNetmask: null,
			iloPassword: null,
			iloUsername: null,
			interfaces: [],
			mgmtIpAddress: null,
			mgmtIpGateway: null,
			mgmtIpNetmask: null,
			offlineReason: null,
			physLocationId: 1,
			profileNames: ["GLOBAL"],
			statusId: 1,
			typeId: 1
		});
		fixture.detectChanges();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});

	it("gets the right status icon", () => {
		component.server.status = "ONLINE";
		expect(component.statusChangeIcon()).toBe("toggle_on");
		component.server.status = "OFFLINE";
		expect(component.statusChangeIcon()).toBe("toggle_off");
		component.server.status = "REPORTED";
		expect(component.statusChangeIcon()).toBe("toggle_on");
		component.server.status = "Anything else";
		expect(component.statusChangeIcon()).toBe("toggle_off");
	});

	it("adds and removes interfaces", () => {
		expect(component.server.interfaces.length).toBe(0);
		component.addInterface(new MouseEvent("click"));
		expect(component.server.interfaces.length).toBe(1);
		component.addInterface(new MouseEvent("click"));
		expect(component.server.interfaces.length).toBe(2);
		component.deleteInterface(new MouseEvent("click"), 1);
		expect(component.server.interfaces.length).toBe(1);
		component.deleteInterface(new MouseEvent("click"), 0);
		expect(component.server.interfaces.length).toBe(0);
	});

	it("adds and removes IP addresses to/from an interface", () => {
		component.addInterface(new MouseEvent("click"));
		expect(component.server.interfaces[0].ipAddresses.length).toBe(0);
		component.addIP(new MouseEvent("click"), component.server.interfaces[0]);
		expect(component.server.interfaces[0].ipAddresses.length).toBe(1);
		component.addIP(new MouseEvent("click"), component.server.interfaces[0]);
		expect(component.server.interfaces[0].ipAddresses.length).toBe(2);
		component.deleteIP(new MouseEvent("click"), component.server.interfaces[0], 1);
		expect(component.server.interfaces[0].ipAddresses.length).toBe(1);
		component.deleteIP(new MouseEvent("click"), component.server.interfaces[0], 0);
		expect(component.server.interfaces[0].ipAddresses.length).toBe(0);
	});

	it("knows if it's a cache", () => {
		const s = component.server;
		s.type = "";
		expect(component.isCache()).toBeFalse();
		s.type = "EDGE";
		expect(component.isCache()).toBeTrue();
		s.type = "EDGE_anything";
		expect(component.isCache()).toBeTrue();
		s.type = "MID";
		expect(component.isCache()).toBeTrue();
		s.type = "MID_anything";
		expect(component.isCache()).toBeTrue();
		s.type = "a string that merely CONTAINS 'EDGE' instead of starting with it";
		expect(component.isCache()).toBeFalse();
		s.type = "RASCAL";
		expect(component.isCache()).toBeFalse();
	});

	it("submits a server creation request", fakeAsync(() => {
		const service = TestBed.inject(ServerService);
		const spy = spyOn(service, "createServer");
		spy.and.callThrough();
		expect(spy).not.toHaveBeenCalled();
		component.isNew = true;

		component.submit(new Event("submit"));
		tick();
		expect(component.isNew).toBeFalse();
		expect(component.server.id).toBeDefined();
	}));

	it("opens the 'change status' dialog", () => {
		const mockMatDialog = TestBed.inject(MatDialog);
		const openSpy = spyOn(mockMatDialog, "open").and.returnValue({
			afterClosed: () => of(true)
		} as MatDialogRef<unknown>);
		component.isNew = true;
		expect(() => component.changeStatus(new MouseEvent("click"))).toThrow();
		expect(openSpy).not.toHaveBeenCalled();
		component.isNew = false;
		component.changeStatus(new MouseEvent("click"));
		expect(openSpy).toHaveBeenCalled();
	});
});
