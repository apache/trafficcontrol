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
import { ComponentFixture, fakeAsync, TestBed, tick } from "@angular/core/testing";
import { BehaviorSubject } from "rxjs";

import { APITestingModule } from "src/app/api/testing";
import { CurrentUserService } from "src/app/shared/current-user/current-user.service";

import { TenantsComponent } from "./tenants.component";

describe("TenantsComponent", () => {
	let component: TenantsComponent;
	let fixture: ComponentFixture<TenantsComponent>;

	beforeEach(async () => {
		await TestBed.configureTestingModule({
			declarations: [ TenantsComponent ],
			imports: [ APITestingModule ],
			providers: [
				{
					provide: CurrentUserService,
					useValue: {
						currentUser: {
							tenantId: 1
						},
						hasPermission: (): true => true,
						userChanged: new BehaviorSubject({})
					}
				}
			]
		})
			.compileComponents();

		fixture = TestBed.createComponent(TenantsComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});

	it("updates the fuzzy search output", fakeAsync(() => {
		let called = false;
		const text = "testquest";
		const spy = jasmine.createSpy("subscriber", (txt: string): void =>{
			if (!called) {
				expect(txt).toBe("");
				called = true;
			} else {
				expect(txt).toBe(text);
			}
		});
		component.searchSubject.subscribe(spy);
		tick();
		expect(spy).toHaveBeenCalled();
		component.searchText = text;
		component.updateURL();
		tick();
		expect(spy).toHaveBeenCalledTimes(2);
	}));

	it("renders parent Tenants", () => {
		expect(component.getParentString({active: true, id: 1, lastUpdated: new Date(), name: "root", parentId: null})).toBe("");
	});

	it("handles contextmenu events", () => {
		expect(()=>component.handleContextMenu({
			action: component.contextMenuItems[0].name,
			data: {
				active: true,
				id: 1,
				lastUpdated: new Date(),
				name: "root",
				parentId: null
			}
		})).not.toThrow();
	});
});
