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
import { HttpClientTestingModule } from "@angular/common/http/testing";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { FormsModule } from "@angular/forms";
import { MatCardModule } from "@angular/material/card";
import { RouterTestingModule } from "@angular/router/testing";

import { ServerService } from "src/app/api";
import { SharedModule } from "src/app/shared/shared.module";

import { StatusesTableComponent } from "./statuses-table.component";

const statuses = [{description: "test", id: 1,lastUpdated: new Date("02/02/2023"), name: "test"}];
describe("StatusesTableComponent", () => {
	let component: StatusesTableComponent;
	let fixture: ComponentFixture<StatusesTableComponent>;

	beforeEach(async () => {
		await TestBed.configureTestingModule({
			declarations: [ StatusesTableComponent ],
			imports:[
				RouterTestingModule,
				HttpClientTestingModule,
				FormsModule,
				MatCardModule,
				SharedModule
			],
			providers:[ServerService]
		})
			.compileComponents();

		fixture = TestBed.createComponent(StatusesTableComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});

	it("should get all statuses from getStatuses",(()=>{
		const service = fixture.debugElement.injector.get(ServerService);
		spyOn(service, "getStatuses").and.returnValue(Promise.resolve(statuses));

		service.getStatuses().then((result)=>{
			expect(result).toEqual(statuses);
		});
	}));
});
