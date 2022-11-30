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

import { ComponentFixture, TestBed } from "@angular/core/testing";

import { ContentServerComponent } from "./content-server.component";

describe("ContentServerComponent", () => {
	let component: ContentServerComponent;
	let fixture: ComponentFixture<ContentServerComponent>;

	beforeEach(async () => {
		await TestBed.configureTestingModule({
			declarations: [ ContentServerComponent ]
		})
			.compileComponents();

		fixture = TestBed.createComponent(ContentServerComponent);
		component = fixture.componentInstance;
		component.server = {
			cacheGroup: "",
			capabilities: [],
			fqdn: "",
			hashCount: -1,
			hashId: "",
			httpsPort: -1,
			interfaceName: "",
			ip: "",
			ip6: "",
			locationId: "",
			port: -1,
			profile: "",
			routingDisabled: 0,
			status: "",
			type: "MID"
		};
		fixture.detectChanges();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});
});
