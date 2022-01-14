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
import { RouterTestingModule } from "@angular/router/testing";
import { AgGridModule } from "ag-grid-angular";
import type { RowNode } from "ag-grid-community";

import { GenericTableComponent } from "./generic-table.component";

describe("GenericTableComponent", () => {
	let component: GenericTableComponent<unknown>;
	let fixture: ComponentFixture<GenericTableComponent<unknown>>;

	beforeEach(async () => {
		await TestBed.configureTestingModule({
			declarations: [
				GenericTableComponent,

			],
			imports: [
				AgGridModule.withComponents([]),
				RouterTestingModule
			]
		}).compileComponents();

		fixture = TestBed.createComponent<GenericTableComponent<unknown>>(GenericTableComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});

	it("can tell if a context menu item is an action", () => {
		expect(component.isAction({href: "/core/dashboard", name: "Dashboard"})).toBeFalse();
		expect(component.isAction({href: "/core/dashboard", name: "Dashboard", newTab: true})).toBeFalse();
		expect(component.isAction({action: "do something", name: "Something"})).toBeTrue();
		expect(component.isAction({action: "do something", multiRow: true, name: "Something"})).toBeTrue();
	});

	it("makes all data pass the filter when there is no search box", () => {
		expect(component.filter({} as RowNode)).toBeTrue();
	});

	it("throws an error trying to check if a context menu item is disabled with no selection", () => {
		expect(()=>component.isDisabled({action: "anything", name: "whatever"})).toThrow();
		expect(()=>component.isDisabled({action: "anything", disabled: ()=>false, name: "who cares"})).toThrow();
	});

	it("throws an error trying to emit a context menu action with no selection", ()=>{
		expect(()=>component.emitContextMenuAction("anything", false, new MouseEvent("click"))).toThrow();
		expect(()=>component.emitContextMenuAction("anything", true, new MouseEvent("click"))).toThrow();
		expect(()=>component.emitContextMenuAction("anything", undefined, new MouseEvent("click"))).toThrow();
	});
});
