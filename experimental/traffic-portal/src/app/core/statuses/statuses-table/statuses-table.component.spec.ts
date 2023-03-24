/*
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import { ComponentFixture, TestBed } from "@angular/core/testing";
import { StatusesTableComponent } from "./statuses-table.component";
import { HttpClientModule } from "@angular/common/http";
import { RouterTestingModule } from "@angular/router/testing";
import { APITestingModule } from "src/app/api/testing";

describe("StatusesTableComponent", () => {
	let component: StatusesTableComponent;
	let fixture: ComponentFixture<StatusesTableComponent>;

	beforeEach(async () => {
		await TestBed.configureTestingModule({
			declarations: [ StatusesTableComponent ],
			imports:[
				HttpClientModule,
				RouterTestingModule.withRoutes([
					{component: StatusesTableComponent, path: ""},
				]),
				APITestingModule
			]
		})
			.compileComponents();

		fixture = TestBed.createComponent(StatusesTableComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});
});
