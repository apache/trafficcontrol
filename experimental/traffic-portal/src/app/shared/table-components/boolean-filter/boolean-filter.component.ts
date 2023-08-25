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
import { Component } from "@angular/core";
import { AgFilterComponent } from "ag-grid-angular";
import { IDoesFilterPassParams, IFilterParams } from "ag-grid-community";

import { LoggingService } from "src/app/shared/logging.service";

/** A model that fully describes the state of a Boolean Filter. */
interface BooleanFilterModel {
	/** Whether or not filtering *should* be done. */
	should: boolean;
	/** The value to be filtered with. */
	value: boolean;
}

/**
 * BooleanFilterComponent is an AG-Grid component that provides a filter for boolean columns.
 */
@Component({
	selector: "tp-boolean-filter",
	styleUrls: ["./boolean-filter.component.scss"],
	templateUrl: "./boolean-filter.component.html",
})
export class BooleanFilterComponent implements AgFilterComponent {

	/** Describes whether or not filtering should be performed. */
	public shouldFilter = false;
	/** Describes which boolean value to match (if filtering is performed). */
	public value = false;

	/** Stores the column name */
	private field = "";

	/** Initialization parameters. */
	private params!: IFilterParams;

	constructor(private readonly log: LoggingService) {}

	/**
	 * Called by AG-Grid to check if the filter is in effect.
	 *
	 * @returns Whether or not the filter is active.
	 */
	public isFilterActive(): boolean {
		return this.shouldFilter;
	}

	/**
	 * Called by AG-Grid to filter a row.
	 *
	 * @param params Parameters given by AG-Grid which contain (among other things) the table row being filtered.
	 * @returns 'true' if the row matches the filter state, 'false' if it should be filtered out.
	 */
	public doesFilterPass(params: IDoesFilterPassParams): boolean {
		if (!params.node) {
			return false;
		}

		let colValue = params.data[this.field];
		if (colValue === undefined) {
			colValue = false;
		}
		return colValue === this.value;
	}

	/**
	 * Handles when a control for the filter is changed.
	 *
	 * @param event The new value of the control.
	 * @param input A string that describes which input has changed.
	 */
	public onChange(event: boolean, input: "should" | "value"): void {
		switch (input) {
			case "should":
				if (event !== this.shouldFilter) {
					this.shouldFilter = event;
					this.params.filterChangedCallback();
				}
				break;
			case "value":
				if (event !== this.value) {
					this.value = event;
					this.params.filterChangedCallback();
				}
		}
	}

	/**
	 * Gets the filter state.
	 *
	 * @returns the full representation of the filter state.
	 */
	public getModel(): BooleanFilterModel {
		return {
			should: this.shouldFilter,
			value: this.value
		};
	}

	/**
	 * Sets the filter state.
	 *
	 * @param model A full representation of a filter state.
	 */
	public setModel(model: BooleanFilterModel | null): void {
		if (!model) {
			model = {
				should: false,
				value: false
			};
		}
		this.shouldFilter = model.should;
		this.value = model.value;
	}

	/**
	 * Called by AG-Grid after the grid is ready.
	 *
	 * @param params All of the information required to set up a filter.
	 */
	public agInit(params: IFilterParams): void {
		this.params = params;
		if (!params.colDef.field) {
			this.log.error("No column name found on boolean-filter parameters");
			return;
		}
		this.field = params.colDef.field;
	}

}
