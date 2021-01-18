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

import { Component, Input, OnDestroy, OnInit } from "@angular/core";
import { Router } from "@angular/router";
import { ColDef, ColGroupDef, ColumnApi, GridApi, GridOptions, GridReadyEvent, RowNode } from "ag-grid-community";
import { BehaviorSubject, Subscription } from "rxjs";
import { fuzzyScore } from "src/app/utils";
import { SSHCellRendererComponent } from "../table-components/ssh-cell-renderer/ssh-cell-renderer.component";

/** Tables can display any of this kind of data. */
type TableData = Record<string, string | number | bigint | Date | boolean | RegExp | null>;

@Component({
	selector: "tp-generic-table[data][cols]",
	styleUrls: ["./generic-table.component.scss"],
	templateUrl: "./generic-table.component.html",
})
export class GenericTableComponent implements OnInit, OnDestroy {

	/** Rows for the table */
	@Input() public data: Array<TableData> = [];
	/** Column and column group definitions. */
	@Input() public cols: Array<ColDef | ColGroupDef> = [];
	/** Optionally provide fuzzy search text. */
	@Input() public fuzzySearch: BehaviorSubject<string> | undefined;

	/** Holds a subscription for the fuzzySearch input if one was provided (otherwise 'null') */
	private fuzzySubscription: Subscription|null = null;

	/** Options to pass into the AG-Grid object. */
	public gridOptions: GridOptions;
	/** Holds a reference to the AG-Grid API (once it has been initialized) */
	private gridAPI: GridApi | undefined;
	/** Holds a reference to the AG-Grid Column API (once it has been initialized)  */
	private columnAPI: ColumnApi | undefined;

	/** Used to handle the case that Angular loads faster than AG-Grid (as it usually does) */
	private initialize = false;

	/** Passed as components to the ag-grid API */
	public components = {
		sshCellRenderer: SSHCellRendererComponent,
	};

	constructor(private readonly router: Router) {
		this.gridOptions = {
			doesExternalFilterPass: this.filter.bind(this),
			isExternalFilterPresent: this.shouldFilter.bind(this),
		};
	}

	public ngOnDestroy(): void {
		if (this.fuzzySubscription) {
			this.fuzzySubscription.unsubscribe();
		}
	}

	public ngOnInit(): void {
		if (this.fuzzySearch) {
			this.fuzzySubscription = this.fuzzySearch.subscribe(
				query => {
					if (this.gridAPI) {
						this.gridAPI.onFilterChanged();
					}
					this.router.navigate([], {queryParams: {search: query}, replaceUrl: true});
				}
			);
		}
	}

	public setAPI(params: GridReadyEvent): void {
		this.gridAPI = params.api;
		this.columnAPI = params.columnApi;
		if (this.initialize) {
			this.initialize = false;
			this.gridAPI.onFilterChanged();
		}
	}

	public shouldFilter(): boolean {
		return this.fuzzySearch !== undefined;
	}

	/**
	 * Checks if a node passes the user's fuzzy search filter.
	 *
	 * @param node The row to check.
	 * @returns whether or not the row passes filtering
	 */
	public filter(node: RowNode): boolean {
		// This can happen when Angular is ready before AG-Grid, which one
		// would hope is normally the case.
		if (!this.columnAPI){
			this.initialize = true;
			return true;
		}
		// ... on the other hand, maybe we just aren't being asked to filter.
		if (!this.fuzzySearch) {
			return true;
		}
		console.log("Filter query:", this.fuzzySearch.value);
		const visibleCols = new Set(this.columnAPI.getAllDisplayedColumns().map(x=>x.getColId()));
		for (const k in node.data) {
			if (!Object.prototype.hasOwnProperty.call(node.data, k) || !visibleCols.has(k)) {
				continue;
			}
			const value = node.data[k];
			let stringVal: string;
			switch (typeof value) {
				case "string":
					stringVal = value;
					break;
				case "boolean":
				case "bigint":
				case "number":
					stringVal = String(value);
					break;
				case "object":
					if (value instanceof Date) {
						stringVal = value.toLocaleString();
					} else if (value instanceof URL) {
						stringVal = value.href;
					} else if (value instanceof RegExp) {
						stringVal = value.source;
					} else {
						continue;
					}
					break;
				default:
					continue;
			}
			const score = fuzzyScore(stringVal.toLocaleLowerCase(), this.fuzzySearch.value.toLocaleLowerCase());
			if ( score < Infinity) {
				return true;
			}
		}

		return false;
	}
}
