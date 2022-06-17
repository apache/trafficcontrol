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
import { MatFormFieldModule } from "@angular/material/form-field";
import { MatInputModule } from "@angular/material/input";
import { MatSelectModule } from "@angular/material/select";
import { BrowserAnimationsModule } from "@angular/platform-browser/animations";
import { RouterTestingModule } from "@angular/router/testing";
import { faToggleOff, faToggleOn } from "@fortawesome/free-solid-svg-icons";

import { ServerService } from "src/app/api";
import { APITestingModule } from "src/app/api/testing";
import { defaultServer } from "src/app/models";

import { ServerDetailsComponent } from "./server-details.component";

describe("ServerDetailsComponent", () => {
	let component: ServerDetailsComponent;
	let fixture: ComponentFixture<ServerDetailsComponent>;

	beforeEach(async () => {
		await TestBed.configureTestingModule({
			declarations: [ ServerDetailsComponent ],
			imports: [
				HttpClientModule,
				RouterTestingModule.withRoutes([
					{component: ServerDetailsComponent, path: "server/:id"},
					{component: ServerDetailsComponent, path: "server/new"}
				]),
				FormsModule,
				ReactiveFormsModule,
				MatSelectModule,
				MatFormFieldModule,
				MatInputModule,
				BrowserAnimationsModule,
				APITestingModule
			],
		}).compileComponents();
		fixture = TestBed.createComponent(ServerDetailsComponent);
		const service = TestBed.inject(ServerService);
		component = fixture.componentInstance;
		component.server = await service.createServer({...defaultServer, interfaces: []});
		fixture.detectChanges();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});

	it("gets the right status icon", () => {
		component.server.status = "ONLINE";
		expect(component.statusChangeIcon).toBe(faToggleOn);
		component.server.status = "OFFLINE";
		expect(component.statusChangeIcon).toBe(faToggleOff);
		component.server.status = "REPORTED";
		expect(component.statusChangeIcon).toBe(faToggleOn);
		component.server.status = "Anything else";
		expect(component.statusChangeIcon).toBe(faToggleOff);
	});

	it("adds and removes interfaces", () => {
		expect(component.server.interfaces.length).toBe(0);
		component.addInterface(new MouseEvent("click"));
		expect(component.server.interfaces.length).toBe(1);
		component.addInterface(new MouseEvent("click"));
		expect(component.server.interfaces.length).toBe(2);
		component.deleteInterface(1);
		expect(component.server.interfaces.length).toBe(1);
		component.deleteInterface(0);
		expect(component.server.interfaces.length).toBe(0);
	});

	it("adds and removes IP addresses to/from an interface", () => {
		component.addInterface(new MouseEvent("click"));
		expect(component.server.interfaces[0].ipAddresses.length).toBe(0);
		component.addIP(component.server.interfaces[0]);
		expect(component.server.interfaces[0].ipAddresses.length).toBe(1);
		component.addIP(component.server.interfaces[0]);
		expect(component.server.interfaces[0].ipAddresses.length).toBe(2);
		component.deleteIP(component.server.interfaces[0], 1);
		expect(component.server.interfaces[0].ipAddresses.length).toBe(1);
		component.deleteIP(component.server.interfaces[0], 0);
		expect(component.server.interfaces[0].ipAddresses.length).toBe(0);
	});

	it("knows if it's a cache", () => {
		const s = component.server;
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
		expect(component.changeStatusDialogOpen).toBeFalse();
		component.changeStatus(new MouseEvent("click"));
		expect(component.changeStatusDialogOpen).toBeTrue();
		component.isNew = true;
		expect(()=>component.changeStatus(new MouseEvent("click"))).toThrow();
	});

	it("closes the 'change status' dialog when done", () => {
		component.changeStatusDialogOpen = true;
		component.doneUpdatingStatus(true);
		expect(component.changeStatusDialogOpen).toBeFalse();
	});
});
