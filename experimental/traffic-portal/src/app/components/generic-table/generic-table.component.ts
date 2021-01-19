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
import { ColDef, ColGroupDef, ColumnApi, GridApi, GridOptions, GridReadyEvent, ITooltipParams, RowNode } from "ag-grid-community";
import { BehaviorSubject, Subscription } from "rxjs";

import {faCaretDown, faColumns} from "@fortawesome/free-solid-svg-icons";

import { fuzzyScore } from "src/app/utils";
import { BooleanFilterComponent } from "../table-components/boolean-filter/boolean-filter.component";
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
	/** Optionally a context to load from localstorage. Providing a unique value for this allows for persistent filter, sort, etc. */
	@Input() public context: string | undefined;

	/** Holds a subscription for the fuzzySearch input if one was provided (otherwise 'null') */
	private fuzzySubscription: Subscription|null = null;

	/** Options to pass into the AG-Grid object. */
	public gridOptions: GridOptions;
	/** Holds a reference to the AG-Grid API (once it has been initialized) */
	private gridAPI: GridApi | undefined;
	/** Holds a reference to the AG-Grid Column API (once it has been initialized)  */
	public columnAPI: ColumnApi | undefined;

	/** Icon used for the 'columns' dropdown item. */
	public columnsIcon = faColumns;
	/** Icon used for the caret/chevron indicating menu direction on button press. */
	public caretIcon = faCaretDown;

	/** Used to handle the case that Angular loads faster than AG-Grid (as it usually does) */
	private initialize = false;

	/** Tracks whether the menu button has been clicked. */
	private menuClicked = false;

	/** Passed as components to the ag-grid API */
	public components = {
		sshCellRenderer: SSHCellRendererComponent,
		tpBooleanFilter: BooleanFilterComponent,
	};

	constructor(private readonly router: Router) {
		this.gridOptions = {
			defaultColDef: {
				filter: true,
				resizable: true,
				sortable: true,
				tooltipValueGetter: (params: ITooltipParams): string => params.value
			},
			doesExternalFilterPass: this.filter.bind(this),
			isExternalFilterPresent: this.shouldFilter.bind(this),
			tooltipShowDelay: 500
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

		// Context loading beyond this point
		if (!this.context) {
			this.gridAPI.sizeColumnsToFit();
			return;
		}

		try {
			const colstates = localStorage.getItem(`${this.context}_table_columns`);
			if (colstates) {
				if (!this.columnAPI.setColumnState(JSON.parse(colstates))) {
					console.error("Failed to load stored column state: one or more columns not found");
				}
			} else {
				this.gridAPI.sizeColumnsToFit();
			}
		} catch (e) {
			console.error(`Failure to retrieve required column info from localStorage (key=${this.context}_table_columns):`, e);
		}

		try {
			const storedSort = localStorage.getItem(`${this.context}_table_sort`);
			if (storedSort) {
				this.gridAPI.setSortModel(JSON.parse(storedSort));
			}
		} catch (e) {
			console.error("Failure to load stored sort state:", e);
		}

		// try {
		// 	$scope.quickSearch = localStorage.getItem(tableName + "_quick_search");
		// 	$scope.gridOptions.api.setQuickFilter($scope.quickSearch);
		// } catch (e) {
		// 	console.error("Failure to load stored quick search:", e);
		// }

		// try {
		// 	const ps = localStorage.getItem(tableName + "_page_size");
		// 	if (ps && ps > 0) {
		// 		$scope.pageSize = Number(ps);
		// 		$scope.gridOptions.api.paginationSetPageSize($scope.pageSize);
		// 	}
		// } catch (e) {
		// 	console.error("Failure to load stored page size:", e);
		// }
	}

	public storeSort(): void {
		if (this.context && this.gridAPI) {
			localStorage.setItem(`${this.context}_table_sort`, JSON.stringify(this.gridAPI.getSortModel()));
		}
	}

	public storeColumns(fit?: boolean): void {
		if (fit && this.gridAPI) {
			this.gridAPI.sizeColumnsToFit();
		}
		if (this.context && this.columnAPI) {
			localStorage.setItem(`${this.context}_table_columns`, JSON.stringify(this.columnAPI.getColumnState()));
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

	public toggleVisibility(col: string): void {
		if (this.columnAPI) {
			const visible = this.columnAPI.getColumn(col).isVisible();
			console.log(`setting column '${col}' to visible:`, !visible);
			this.columnAPI.setColumnVisible(col, !visible);
		}
	}

	public toggleMenu(e: Event): void {
		e.stopPropagation();
		this.menuClicked = !this.menuClicked;
	}

	public get showMenu(): boolean {
		return this.menuClicked && (this.columnAPI ? true : false);
	}
}
