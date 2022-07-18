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
import type { IFilterParams, RowNode } from "ag-grid-community";

import { BooleanFilterComponent } from "./boolean-filter.component";

describe("BooleanFilterComponent", () => {
	let component: BooleanFilterComponent;
	let fixture: ComponentFixture<BooleanFilterComponent>;
	let filterChangedCallback: jasmine.Spy;

	beforeEach(async () => {
		await TestBed.configureTestingModule({
			declarations: [ BooleanFilterComponent ]
		})
			.compileComponents();
	});

	beforeEach(() => {
		filterChangedCallback = jasmine.createSpy("AG-Grid filter changed callback");
		fixture = TestBed.createComponent(BooleanFilterComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
		component.agInit({
			colDef: {
				field: "test"
			},
			filterChangedCallback,
			valueGetter: (n: RowNode): boolean => n.data
		} as unknown as IFilterParams);
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});

	it("gets and sets the filter model", ()=>{
		const model = {should: false, value: false};
		expect(component.getModel()).toEqual(model);
		expect(component.isFilterActive()).toBeFalse();

		model.should = true;
		component.setModel(model);
		expect(component.getModel()).toEqual(model);
		expect(component.isFilterActive()).toBeTrue();

		model.value = true;
		component.setModel(model);
		expect(component.getModel()).toEqual(model);
		expect(component.isFilterActive()).toBeTrue();

		model.should = false;
		component.setModel(model);
		expect(component.getModel()).toEqual(model);
		expect(component.isFilterActive()).toBeFalse();
	});

	it("handles changes", () => {
		component.onChange(true, "value");
		expect(filterChangedCallback).toHaveBeenCalled();
		expect(component.getModel().value).toBeTrue();
		expect(component.isFilterActive()).toBeFalse();

		component.onChange(true, "should");
		expect(filterChangedCallback).toHaveBeenCalledTimes(2);
		expect(component.getModel().should).toBeTrue();
		expect(component.isFilterActive()).toBeTrue();
	});

	it("knows if a filter passes", () => {
		const node = {data: false} as RowNode;
		const data = {test: false};
		component.onChange(true, "should");
		expect(component.isFilterActive()).toBeTrue();
		expect(component.doesFilterPass({data, node})).toBeTrue();
		node.data = true;
		data.test = true;
		expect(component.doesFilterPass({data, node})).toBeFalse();
		component.onChange(true, "value");
		expect(component.doesFilterPass({data, node})).toBeTrue();
		node.data = false;
		data.test = false;
		expect(component.doesFilterPass({data, node})).toBeFalse();
	});
});
