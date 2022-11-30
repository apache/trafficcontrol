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
import { of } from "rxjs";

import { RouterDiffComponent } from "./router-diff.component";

describe("RouterDiffComponent", () => {
	let component: RouterDiffComponent;
	let fixture: ComponentFixture<RouterDiffComponent>;

	beforeEach(async () => {
		await TestBed.configureTestingModule({
			declarations: [ RouterDiffComponent ]
		})
			.compileComponents();

		fixture = TestBed.createComponent(RouterDiffComponent);
		component = fixture.componentInstance;
		component.snapshots = of({current: {}, pending: {}});
		fixture.detectChanges();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});
});
