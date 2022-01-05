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

import { Component, ElementRef, EventEmitter, HostListener, Input, OnDestroy, OnInit, Output, ViewChild } from "@angular/core";
import { Router } from "@angular/router";
import { BehaviorSubject, Subscription } from "rxjs";

import type {
	CellContextMenuEvent,
	ColDef,
	ColGroupDef,
	Column,
	ColumnApi,
	CsvExportParams,
	GridApi,
	GridOptions,
	GridReadyEvent,
	ITooltipParams,
	RowNode
} from "ag-grid-community";
import {faCaretDown, faColumns, faDownload} from "@fortawesome/free-solid-svg-icons";

import { fuzzyScore } from "src/app/utils";
import { BooleanFilterComponent } from "../table-components/boolean-filter/boolean-filter.component";
import { SSHCellRendererComponent } from "../table-components/ssh-cell-renderer/ssh-cell-renderer.component";
import { UpdateCellRendererComponent } from "../table-components/update-cell-renderer/update-cell-renderer.component";

/** A context menu action that acts on a single row only. */
interface ContextMenuSingleAction<T> {
	/**
	 * The name of the action; this will be emitted along with the selected data back to the host so that it knows which option was clicked.
	 */
	action: string;
	/**
	 * If present, this method will be called to determine if the action should be disabled.
	 *
	 * @param data The selected data which can be used to make the determination.
	 * @param api A reference to the Grid's API - which must be checked for initialization, unfortunately.
	 */
	disabled?: (data: T, api?: GridApi) => boolean;
	/**
	 * If given and true, causes the action to act on all selected data instead of a single row.
	 *
	 * Actions that do not act on a single row are disabled when multiple rows are selected.
	 */
	multiRow?: false;
	/** A human-readable name for the action which is displayed to the user. */
	name: string;
}

/** A context menu action that acts on any number of selected rows. */
interface ContextMenuMultiAction<T> {
	/**
	 * The name of the action; this will be emitted along with the selected data back to the host so that it knows which option was clicked.
	 */
	action: string;
	/**
	 * If present, this method will be called to determine if the action should be disabled.
	 *
	 * @param data The selected data which can be used to make the determination.
	 * @param api A reference to the Grid's API - which must be checked for initialization, unfortunately.
	 */
	disabled?: (data: Array<T>, api?: GridApi) => boolean;
	/**
	 * If given and true, causes the action to act on all selected data instead of a single row.
	 *
	 * Actions that do not act on a single row are disabled when multiple rows are selected.
	 */
	multiRow: true;
	/** A human-readable name for the action which is displayed to the user. */
	name: string;
}

/** ContextMenuActions represent an action that can be taken in a context menu. */
type ContextMenuAction<T> = ContextMenuSingleAction<T> | ContextMenuMultiAction<T>;

/** ContextMenuLinks represent a link within a context menu. They aren't templated, so currently have limited uses. */
interface ContextMenuLink {
	/**
	 * href is inserted literally as the 'href' property of an anchor. Which means that if it's not relative it will be mangled for security
	 * reasons.
	 */
	href: string;
	/** A human-readable name for the link which is displayed to the user. */
	name: string;
	/** If given and true, sets the link to open in a new browsing context (or "tab"). */
	newTab?: boolean;
}

/** ContextMenuItems represent items in a context menu. They can be links or arbitrary actions. */
export type ContextMenuItem<T> = ContextMenuAction<T> | ContextMenuLink;

/** ContextMenuActionEvent is emitted by the GenericTableComponent when an action in its context menu was clicked. */
export interface ContextMenuActionEvent<T> {
	/** action is the 'action' property of the clicked action. */
	action: string;
	/** data is the selected data on which the action will act. */
	data: T | Array<T>;
}

/**
 * GenericTableComponent is the controller for generic tables.
 */
@Component({
	selector: "tp-generic-table[data][cols]",
	styleUrls: ["./generic-table.component.scss"],
	templateUrl: "./generic-table.component.html",
})
export class GenericTableComponent<T> implements OnInit, OnDestroy {

	/** Rows for the table */
	@Input() public data: Array<T> = [];
	/** Column and column group definitions. */
	@Input() public cols: Array<ColDef | ColGroupDef> = [];
	/** Optionally provide fuzzy search text. */
	@Input() public fuzzySearch: BehaviorSubject<string> | undefined;
	/** Optionally a context to load from localstorage. Providing a unique value for this allows for persistent filter, sort, etc. */
	@Input() public context: string | undefined;
	/** Optionally a set of context menu items. If not given, the context menu is disabled. */
	@Input() public contextMenuItems: Array<ContextMenuItem<T>> = [];
	/** Emits when context menu actions are clicked. Type safety is the host's responsibility! */
	@Output() public contextMenuAction = new EventEmitter<ContextMenuActionEvent<T>>();
	/**
	 * Checks if a context menu item is an action.
	 *
	 * @param i The menu item to check.
	 * @returns 'true' if 'i' is an action, 'false' if it's a link.
	 */
	public isAction(i: ContextMenuItem<T>): i is ContextMenuAction<T> {
		return Object.prototype.hasOwnProperty.call(i, "action");
	}

	/** Holds a reference to the context menu which is used to calculate its size. */
	@ViewChild("contextmenu") public contextmenu: ElementRef | null = null;

	/**
	 * This event handler listens to click events anywhere, since if you're clicking outside
	 * the context menu it should close (and it should also close when an action is taken).
	 *
	 * @param e The click event.
	 */
	@HostListener("document:click", ["$event"])
	public clickOutside(e: MouseEvent): void {
		e.stopPropagation();
		this.showContextMenu = false;
		this.menuClicked = false;
	}

	/** This holds a reference to the table's selected data, which is emitted on context menu action clicks. */
	public selected: T | null = null;

	/** Holds a subscription for the fuzzySearch input if one was provided (otherwise 'null') */
	private fuzzySubscription: Subscription|null = null;

	/** Options to pass into the AG-Grid object. */
	public gridOptions: GridOptions;
	/** Holds a reference to the AG-Grid API (once it has been initialized) */
	private gridAPI: GridApi | undefined;
	/** Holds a reference to the AG-Grid Column API (once it has been initialized)  */
	public columnAPI: ColumnApi | undefined;

	/** Icon used for the 'columns' dropdown item. */
	public readonly columnsIcon = faColumns;
	/** Icon used for the caret/chevron indicating menu direction on button press. */
	public readonly caretIcon = faCaretDown;
	/**
	 * Icon for the "export to CSZ" button.
	 */
	public readonly downloadIcon = faDownload;

	/** Used to handle the case that Angular loads faster than AG-Grid (as it usually does) */
	private initialize = false;

	/** Tracks whether the menu button has been clicked. */
	private menuClicked = false;

	/** Tells whether or not to show the cell context menu. */
	public showContextMenu = false;

	/** Passed as components to the ag-grid API */
	public components = {
		sshCellRenderer: SSHCellRendererComponent,
		tpBooleanFilter: BooleanFilterComponent,
		updateCellRenderer: UpdateCellRendererComponent,
	};

	/**
	 * The number of currently selected rows (-1 if the grid is not initialized).
	 */
	public get selectionCount(): number {
		if (!this.gridAPI) {
			return -1;
		}
		return this.gridAPI.getSelectedRows().length;
	}

	/**
	 * All currently selected rows.
	 */
	public get fullSelection(): Array<T> {
		if (!this.gridAPI) {
			return [];
		}
		return this.gridAPI.getSelectedRows();
	}

	/**
	 * All column definitions (regardless of whether or not they're visible).
	 */
	public get columns(): Array<Column> {
		if (!this.columnAPI) {
			return [];
		}
		return this.columnAPI.getAllColumns() ?? [];
	}

	/**
	 * Contructs the component with its required injections.
	 *
	 * @param router Used to update the 'search' query parameter on fuzzy filter changes.
	 */
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
			preventDefaultOnContextMenu: true,
			rowSelection: "multiple",
			suppressContextMenu: true,
			tooltipShowDelay: 500
		};
	}

	/**
	 * Cleans up async resources on component destruction.
	 */
	public ngOnDestroy(): void {
		if (this.fuzzySubscription) {
			this.fuzzySubscription.unsubscribe();
		}
	}

	/**
	 * Sets up async resources after component initialization is complete.
	 */
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

	/**
	 * Sets up the API once the Grid is ready, and loads the table context if one was provided.
	 *
	 * @param params The GridReadyEvent.
	 */
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

	/** When sorting changes, stores the sorting state if a context was provided. */
	public storeSort(): void {
		if (this.context && this.gridAPI) {
			localStorage.setItem(`${this.context}_table_sort`, JSON.stringify(this.gridAPI.getSortModel()));
		}
	}

	/**
	 * When column order/visibility change, stores the column state if a context was provided.
	 *
	 * @param fit If given and true, sizes columns to fit the view area. This is typically done when a column's visibility is toggled.
	 */
	public storeColumns(fit?: boolean): void {
		if (fit && this.gridAPI) {
			this.gridAPI.sizeColumnsToFit();
		}
		if (this.context && this.columnAPI) {
			localStorage.setItem(`${this.context}_table_columns`, JSON.stringify(this.columnAPI.getColumnState()));
		}
	}

	/**
	 * Checks if external filtering should be performed.
	 *
	 * @returns whether or not a fuzzy search filter was provided.
	 */
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

	/**
	 * Toggles the visibility of a column.
	 *
	 * @param col The ID of a column to toggle.
	 */
	public toggleVisibility(col: string): void {
		if (this.columnAPI) {
			const column = this.columnAPI.getColumn(col);
			if (!column) {
				console.error(`Failed to set visibility for column '${col}': no such column`);
				return;
			}
			const visible = column.isVisible();
			this.columnAPI.setColumnVisible(col, !visible);
		}
	}

	/**
	 * Toggles the column visibility menu open state.
	 *
	 * @param e The mouse event which has its propagation stopped so it doesn't interfere with the context menu.
	 */
	public toggleMenu(e: Event): void {
		e.stopPropagation();
		this.showContextMenu = false;
		this.menuClicked = !this.menuClicked;
	}

	/** This tracks whether the column visibility menu is/should be open. */
	public get showMenu(): boolean {
		return this.menuClicked && (this.columnAPI ? true : false);
	}

	/** This is the styling of the table's context menu. */
	public menuStyle = {
		bottom: "unset",
		left: "0",
		right: "unset",
		top: "0",
	};

	/**
	 * This dumb hack is necessary because for at least the past three years
	 * AG-Grid has not respected its `suppressContextMenu` and
	 * `preventDefaultOnContextMenu` settings.
	 *
	 * @param e The raw contextmenu event, so we can manually prevent its default behavior.
	 */
	public preventDefault(e: Event): void {
		e.preventDefault();
	}

	/**
	 * Handles opening the context menu when a table cell is right-clicked.
	 *
	 * @param params The AG-Grid-emitted event.
	 */
	public onCellContextMenu(params: CellContextMenuEvent): void {
		if (!params.event || !(params.event instanceof MouseEvent)) {
			console.warn("cellContextMenu fired with no underlying event");
			return;
		}

		this.menuClicked = false;

		if (!this.contextmenu) {
			console.warn("element reference to 'contextmenu' still null after view init");
			return;
		}

		this.showContextMenu = true;
		this.menuStyle.left = `${params.event.clientX}px`;
		this.menuStyle.top = `${params.event.clientY}px`;
		this.menuStyle.bottom = "unset";
		this.menuStyle.right = "unset";
		const boundingRect = this.contextmenu.nativeElement.getBoundingClientRect();
		// const boundingRect = (document.getElementById("context-menu") as HTMLMenuElement).getBoundingClientRect();

		if (boundingRect.bottom > window.innerHeight){
			this.menuStyle.bottom = `${window.innerHeight - params.event.clientY}px`;
			this.menuStyle.top = "unset";
		}
		if (boundingRect.right > window.innerWidth) {
			this.menuStyle.right = `${window.innerWidth - params.event.clientX}px`;
			this.menuStyle.left = "unset";
		}
		this.selected = params.data;
	}

	/**
	 * Checks if a context menu action is disabled.
	 *
	 * @param a The action to check.
	 * @returns Whether or not `a` should be disabled.
	 */
	public isDisabled(a: ContextMenuAction<T>): boolean {
		if (!this.selected) {
			throw new Error("cannot check if a context menu is disabled for a selection when there is no selection");
		}
		if (!a.multiRow && this.selectionCount > 1) {
			return true;
		}
		if (a.disabled) {
			if (a.multiRow) {
				return a.disabled(this.fullSelection, this.gridAPI);
			}
			return a.disabled(this.selected, this.gridAPI);
		}
		return false;
	}

	/**
	 * Handles when the user clicks on a context menu action item by emitting the proper data.
	 *
	 * @param action The action that was clicked.
	 * @param multi If 'true' the emitted data will be all rows currently selected, otherwise only the item on which the user right-clicked.
	 * @param e The mouse event that triggered this handler.
	 */
	public emitContextMenuAction(action: string, multi: boolean | undefined, e: MouseEvent): void {
		if (!this.selected) {
			throw new Error("nothing selected, cannot emit context menu action");
		}
		e.stopPropagation();
		if (multi) {
			this.contextMenuAction.emit({
				action,
				data: this.selectionCount > 0 ? this.fullSelection : [this.selected]
			});
		} else {
			this.contextMenuAction.emit({
				action,
				data: this.selected
			});
		}
		this.showContextMenu = false;
	}

	/**
	 * Downloads the table data as a CSV file.
	 */
	public download(): void {
		if (!this.gridAPI) {
			console.error("Cannot download: no grid API handle");
			return;
		}

		const params: CsvExportParams = {
			onlySelected: this.gridAPI.getSelectedNodes().length > 0,
		};

		if (this.context) {
			params.fileName = `${this.context}.csv`;
		}

		this.gridAPI.exportDataAsCsv(params);
	}


	/**
	 * Select or de-select all rows.
	 *
	 * @param de If given and true, de-selects all rows. Otherwise all rows will be selected.
	 */
	public selectAll(de?: boolean): void {
		if (!this.gridAPI) {
			console.error("Cannot de-select: no grid API handle");
			return;
		}

		if (de) {
			this.gridAPI.deselectAll();
		} else {
			this.gridAPI.selectAllFiltered();
		}
	}
}
