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

import { ContentRouterComponent } from "./content-router.component";

describe("ContentRouterComponent", () => {
	let component: ContentRouterComponent;
	let fixture: ComponentFixture<ContentRouterComponent>;

	beforeEach(async () => {
		await TestBed.configureTestingModule({
			declarations: [ ContentRouterComponent ]
		})
			.compileComponents();

		fixture = TestBed.createComponent(ContentRouterComponent);
		component = fixture.componentInstance;
		component.router = {
			// eslint-disable-next-line @typescript-eslint/naming-convention
			"api.port": "",
			fqdn: "",
			httpsPort: -1,
			ip: "",
			ip6: "",
			location: "",
			port: -1,
			profile: "",
			// eslint-disable-next-line @typescript-eslint/naming-convention
			"secure.api.port": "",
			status: ""
		};
		fixture.detectChanges();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});
});
