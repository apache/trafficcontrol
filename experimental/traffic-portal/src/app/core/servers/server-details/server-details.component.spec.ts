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
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { FormsModule, ReactiveFormsModule } from "@angular/forms";
import { RouterTestingModule } from "@angular/router/testing";

import {CacheGroupService, CDNService, ProfileService, ServerService, TypeService} from "../../../shared/api";
import {PhysicalLocationService} from "../../../shared/api/PhysicalLocationService";
import { ServerDetailsComponent } from "./server-details.component";

describe("ServerDetailsComponent", () => {
	let component: ServerDetailsComponent;
	let fixture: ComponentFixture<ServerDetailsComponent>;

	beforeEach(async () => {
		const mockAPIService = jasmine.createSpyObj(["getServers", "getCacheGroups", "getCDNs",
			"getProfiles", "getTypes", "getPhysicalLocations", "getStatuses", "getServerTypes"]);
		mockAPIService.getCacheGroups.and.returnValue(new Promise(r => r([])));
		mockAPIService.getCDNs.and.returnValue(new Promise(r => r([])));
		mockAPIService.getStatuses.and.returnValue(new Promise(r => r([])));
		mockAPIService.getProfiles.and.returnValue(new Promise(r => r([])));
		mockAPIService.getPhysicalLocations.and.returnValues(new Promise(r => r([])));
		mockAPIService.getServers.and.returnValues(new Promise(r => r([])));
		mockAPIService.getServerTypes.and.returnValues(new Promise(r => r([])));

		await TestBed.configureTestingModule({
			declarations: [ ServerDetailsComponent ],
			imports: [ HttpClientModule, RouterTestingModule, FormsModule, ReactiveFormsModule ],
			 providers: [
				 { provide: ServerService, useValue: mockAPIService },
				 { provide: CacheGroupService, useValue: mockAPIService },
				 { provide: CDNService, useValue: mockAPIService },
				 { provide: TypeService, useValue: mockAPIService },
				 { provide: PhysicalLocationService, useValue: mockAPIService },
				 { provide: ProfileService, useValue: mockAPIService }
			 ]
		}).compileComponents();
	});

	beforeEach(() => {
		fixture = TestBed.createComponent(ServerDetailsComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});
});
