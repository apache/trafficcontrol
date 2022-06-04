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
import { type ComponentFixture, TestBed } from "@angular/core/testing";
import {MAT_DIALOG_DATA, MatDialogRef} from "@angular/material/dialog";
import {of} from "rxjs";

import { ServerService } from "src/app/api";
import { APITestingModule } from "src/app/api/testing";
import {defaultServer, Server} from "src/app/models";

import { UpdateStatusComponent } from "./update-status.component";

describe("UpdateStatusComponent", () => {
	let component: UpdateStatusComponent;
	let fixture: ComponentFixture<UpdateStatusComponent>;
	let result: boolean;
	let mockMatDialog: jasmine.SpyObj<MatDialogRef<boolean>>;

	beforeEach(() => {
		mockMatDialog = jasmine.createSpyObj("MatDialogRef", ["close", "afterClosed"]);
		TestBed.configureTestingModule({
			declarations: [ UpdateStatusComponent ],
			imports: [ HttpClientModule, APITestingModule ],
			providers: [ {provide: MatDialogRef, useValue: mockMatDialog },
				{provide: MAT_DIALOG_DATA, useValue: (): Array<Server> => []}]
		}).compileComponents();
		fixture = TestBed.createComponent(UpdateStatusComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});

	it("gets server names", () => {
		expect(component.serverName).toBe("0 servers");
		component.servers = [{
			...defaultServer,
			hostName: "host",
		}];
		expect(component.serverName).toBe("host");
		component.servers.push({...defaultServer, hostName: "a different host"});
		expect(component.serverName).toBe("2 servers");
	});

	it("sets the 'current' status ID based on selected servers", () => {
		expect(component.currentStatus).toBeNull();

		fixture = TestBed.createComponent(UpdateStatusComponent);
		component = fixture.componentInstance;
		component.servers = [{...defaultServer, statusId: 9001}];
		fixture.detectChanges();
		expect(component.currentStatus).toBe(9001);

		fixture = TestBed.createComponent(UpdateStatusComponent);
		component = fixture.componentInstance;
		component.servers = [{...defaultServer, statusId: 9001}, {...defaultServer, statusId: 9001}];
		fixture.detectChanges();
		expect(component.currentStatus).toBe(9001);

		fixture = TestBed.createComponent(UpdateStatusComponent);
		component = fixture.componentInstance;
		component.servers = [{...defaultServer, statusId: 9001}, {...defaultServer, statusId: 9}];
		fixture.detectChanges();
		expect(component.currentStatus).toBeNull();
	});

	it("cancels", () => {
		result = false;
		mockMatDialog.afterClosed.and.returnValue(of(result));
		mockMatDialog.afterClosed().subscribe(value => {
			expect(value).toBe(result);
		});
		component.cancel();
		expect(mockMatDialog.close.calls.count()).toBe(1);
	});

	it("knows if the user-selected status is an 'offline' status", () => {
		expect(component.status).toBeNull();
		expect(component.isOffline).toBeFalse();

		component.status = null;
		expect(component.isOffline).toBeFalse();

		component.status = {description: "", name: "OFFLINE"};
		expect(component.isOffline).toBeTrue();

		component.status = {description: "", name: "some weird custom status"};
		expect(component.isOffline).toBeTrue();

		component.status = {description: "", name: "ONLINE"};
		expect(component.isOffline).toBeFalse();

		component.status = {description: "", name: "REPORTED"};
		expect(component.isOffline).toBeFalse();
	});

	it("submits a request to update each server", async () => {
		result = true;
		mockMatDialog.afterClosed.and.returnValue(of(result));
		mockMatDialog.afterClosed().subscribe(value => {
			expect(value).toBe(result);
		});

		const service = TestBed.inject(ServerService);
		component.status = (await service.getStatuses()).find(s=>s.name==="ONLINE") ?? null;

		const srv = await service.createServer({...defaultServer});
		component.servers = [srv];
		await component.submit(new Event("click"));
		expect(mockMatDialog.close.calls.count()).toBe(1);

		result = true;
		mockMatDialog.afterClosed.and.returnValue(of(result));
		component.status = (await service.getStatuses()).find(s=>s.name==="OFFLINE") ?? null;
		await component.submit(new Event("click"));
		expect(mockMatDialog.close.calls.count()).toBe(2);

		result = false;
		mockMatDialog.afterClosed.and.returnValue(of(result));
		component.status = {description: "", name: "no such status"};
		await component.submit(new Event("click"));
		expect(mockMatDialog.close.calls.count()).toBe(3);
	});
});
