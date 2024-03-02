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

import {
	Component,
	type ElementRef,
	EventEmitter,
	HostListener,
	Input,
	type OnDestroy,
	type OnInit,
	Output,
	ViewChild
} from "@angular/core";
import { MatDialog } from "@angular/material/dialog";
import { ActivatedRoute, type ParamMap, type Params, Router } from "@angular/router";
import type {
	CellContextMenuEvent,
	ColDef,
	ColGroupDef,
	Column,
	ColumnApi,
	DateFilterModel,
	FilterChangedEvent,
	GridApi,
	GridOptions,
	GridReadyEvent,
	ITooltipParams,
	NumberFilterModel,
	RowNode,
	TextFilterModel
} from "ag-grid-community";
import type { BehaviorSubject, Subscription } from "rxjs";

import { DownloadOptionsDialogComponent } from "src/app/shared/generic-table/download-options/download-options-dialog.component";
import { fuzzyScore } from "src/app/utils";

import { LoggingService } from "../logging.service";
import { BooleanFilterComponent } from "../table-components/boolean-filter/boolean-filter.component";
import { EmailCellRendererComponent } from "../table-components/email-cell-renderer/email-cell-renderer.component";
import { SSHCellRendererComponent } from "../table-components/ssh-cell-renderer/ssh-cell-renderer.component";
import { TelephoneCellRendererComponent } from "../table-components/telephone-cell-renderer/telephone-cell-renderer.component";
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
export type ContextMenuAction<T> = ContextMenuSingleAction<T> | ContextMenuMultiAction<T>;

/** ContextMenuLinks represent a link within a context menu. They aren't templated, so currently have limited uses. */
interface ContextMenuLink<T> {
	/**
	 * If present, this method will be called to determine if the link should be
	 * disabled.
	 *
	 * @param data The selected data which can be used to make the
	 * determination. This will be a single item if a single item is selected,
	 * or an array if multiple are selected.
	 * @param api A reference to the Grid's API - which must be checked for
	 * initialization, unfortunately.
	 */
	disabled?: (selection: T | Array<T>) => boolean;
	/** If present, determines the URL fragment used during navigation. */
	fragment?: string | ((selectedRow: T) => (string | null));
	/**
	 * href is inserted literally as the 'href' property of an anchor. Which means that if it's not relative it will be mangled for security
	 * reasons.
	 */
	href: string | ((selectedRow: T) => (string | Array<string>));
	/** A human-readable name for the link which is displayed to the user. */
	name: string;
	/** If given and true, sets the link to open in a new browsing context (or "tab"). */
	newTab?: boolean;
	/** If present, query string parameters to pass during navigation. */
	queryParams?: Params | ParamMap | ((selectedRow: T) => (Params | ParamMap | null));
}

/** ContextMenuItems represent items in a context menu. They can be links or arbitrary actions. */
export type ContextMenuItem<T> = ContextMenuAction<T> | ContextMenuLink<T>;

/**
 * Specifies what happens when a row in the grid is double-clicked.
 */
export interface DoubleClickLink<T> {
	/**
	 * If present, this method will be called to determine if the double click should be
	 * ignored.
	 *
	 * @param data The selected data which can be used to make the
	 * determination. This will be a single item if a single item is selected,
	 * or an array if multiple are selected.
	 * @param api A reference to the Grid's API - which must be checked for
	 * initialization, unfortunately.
	 */
	disabled?: (selection: T | Array<T>) => boolean;
	/**
	 * href is inserted literally as the 'href' property of an anchor. Which means that if it's not relative it will be mangled for security
	 * reasons.
	 */
	href: string | ((selectedRow: T) => string);
}

/** ContextMenuActionEvent is emitted by the GenericTableComponent when an action in its context menu was clicked. */
export interface ContextMenuActionEvent<T> {
	/** action is the 'action' property of the clicked action. */
	action: string;
	/** data is the selected data on which the action will act. */
	data: T | Array<T>;
}

/**
 * Checks if a context menu item is an action.
 *
 * @param i The menu item to check.
 * @returns 'true' if 'i' is an action, 'false' if it's a link.
 */
export function isAction<T=unknown>(i: ContextMenuItem<T>): i is ContextMenuAction<T> {
	return Object.prototype.hasOwnProperty.call(i, "action");
}

/**
 * TableTitleButton represents a button added to the heading of the table.
 */
export interface TableTitleButton {
	action: string;
	text: string;
}

/**
 * Gets a basic type from a column definition.
 *
 * @param col The definition of the column
 * @returns The basic type of the column - or `null` if it couldn't be
 * determined.
 */
export function getColType(col: ColDef): "string" | "number" | "date" | null {
	if (!Object.prototype.hasOwnProperty.call(col, "filter") || col.filter === true) {
		return "string";
	}
	if (typeof(col.filter) !== "string") {
		return null;
	}
	switch(col.filter) {
		case "textFilter":
		case "agTextColumnFilter":
			return "string";
		case "agNumberColumnFilter":
			return "number";
		case "agDateColumnFilter":
			return "date";
	}
	return null;
}

/**
 * Given some query parameters, the columns of a table, and a hook into the
 * AG-Grid API of said table, sets up filtering based on matches between the
 * names of query parameters and the raw data fields of the columns.
 *
 * @param params The query string parameters.
 * @param columns The column definitions.
 * @param api An API handle to the grid.
 */
export function setUpQueryParamFilter<T>(params: ParamMap, columns: ColDef<T>[], api: GridApi): void {
	for (const col of columns) {
		if (typeof(col.field) !== "string") {
			continue;
		}

		// According to the AG-Grid docs, you can pass
		const filter = api.getFilterInstance(col.field);
		if (!filter || !col.field) {
			continue;
		}
		const values = params.getAll(col.field);
		if (values.length < 1) {
			continue;
		}

		const colType = getColType(col);
		if (!colType) {
			return;
		}

		let filterModel;
		switch(colType) {
			case "string":
				if (values.length === 1) {
					filterModel = {
						filter: values[0],
						type: "equals"
					};
				} else {
					filterModel = {
						condition1: {
							filter: values[0],
							type: "equals"
						},
						condition2: {
							filter: values[1],
							type: "equals"
						},
						operator: "OR",
					};
				}
				break;
			case "number":
				if (values.length === 1) {
					filterModel = {
						filter: parseInt(values[0], 10),
						type: "equals"
					};
					if (isNaN(filterModel.filter)) {
						continue;
					}
				} else {
					filterModel = {
						condition1: {
							filter: parseInt(values[0], 10),
							type: "equals"
						},
						condition2: {
							filter: parseInt(values[1], 10),
							type: "equals"
						},
						operator: "OR",
					};
					if (isNaN(filterModel.condition1.filter) || isNaN(filterModel.condition2.filter)) {
						continue;
					}
				}
				break;
			case "date":
				const date = new Date(values[0]);
				if (Number.isNaN(date.getTime())) {
					continue;
				}
				const pad = (num: number): string => String(num).padStart(2,"0");
				filterModel = {
					dateFrom: [
						`${date.getUTCFullYear()}-${pad(date.getUTCMonth()+1)}-${pad(date.getUTCDate())}`,
						`${pad(date.getUTCHours())}:${pad(date.getUTCMinutes())}:${pad(date.getUTCSeconds())}`,
					].join(" "),
					type: "equals"
				};
				break;
		}
		filter.setModel(filterModel);
	}
	api.onFilterChanged();
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
	@Input() public contextMenuItems: readonly ContextMenuItem<Readonly<T>>[] = [];
	/** Optionally a set of additional table title buttons. */
	@Input() public tableTitleButtons: Array<TableTitleButton> = [];
	/** Optionally a set of additional more menu buttons. */
	@Input() public moreMenuButtons: Array<TableTitleButton> = [];
	/** Optionally a link that determines the action when double-clicking a grid row */
	@Input() public doubleClickLink: DoubleClickLink<T> | undefined;
	/** Emits when context menu actions are clicked. Type safety is the host's responsibility! */
	@Output() public contextMenuAction = new EventEmitter<ContextMenuActionEvent<T>>();
	/** Emits when title button actions are clicked. Type safety is the host's responsibility! */
	@Output() public tableTitleButtonAction = new EventEmitter<string>();
	/** Emits when more menu title button actions are clicked. Type safety is the host's responsibility! */
	@Output() public moreMenuButtonAction = new EventEmitter<string>();

	public isAction = isAction;

	/** Holds a reference to the context menu which is used to calculate its size. */
	@ViewChild("contextmenu") public contextmenu!: ElementRef;

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
	}

	/** This holds a reference to the table's selected data, which is emitted on context menu action clicks. */
	public selected: T | null = null;

	/** Holds a subscription for the fuzzySearch input if one was provided (otherwise 'null') */
	private fuzzySubscription: Subscription|null = null;

	/** Options to pass into the AG-Grid object. */
	public gridOptions: GridOptions;
	/** Holds a reference to the AG-Grid API (once it has been initialized) */
	public gridAPI!: GridApi;
	/** Holds a reference to the AG-Grid Column API (once it has been initialized)  */
	public columnAPI: ColumnApi | undefined;

	/** Used to handle the case that Angular loads faster than AG-Grid (as it usually does) */
	private initialize = true;

	/** Tells whether or not to show the cell context menu. */
	public showContextMenu = false;

	/** Passed as components to the ag-grid API */
	public components = {
		emailCellRenderer: EmailCellRendererComponent,
		phoneNumberCellRenderer: TelephoneCellRendererComponent,
		sshCellRenderer: SSHCellRendererComponent,
		tpBooleanFilter: BooleanFilterComponent,
		updateCellRenderer: UpdateCellRendererComponent,
	};

	/**
	 * The number of currently selected rows (-1 if the grid is not initialized).
	 */
	public get selectionCount(): number {
		return this.gridAPI.getSelectedRows().length;
	}

	/**
	 * All currently selected rows.
	 */
	public get fullSelection(): Array<T> {
		return this.gridAPI.getSelectedRows();
	}

	/**
	 * All column definitions (regardless of whether or not they're visible).
	 */
	public get columns(): Array<Column> {
		if (!this.columnAPI) {
			return [];
		}
		return (this.columnAPI.getColumns() ?? []).reverse();
	}

	constructor(private readonly router: Router,
		private readonly route: ActivatedRoute,
		private readonly dialog: MatDialog,
		private readonly log: LoggingService) {
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
		this.gridOptions.onRowDoubleClicked = (e): void => {
			if (this.doubleClickLink !== undefined) {
				if (!this.doubleClickLink?.disabled) {
					let href = "";
					if (typeof (this.doubleClickLink.href) === "string") {
						href = this.doubleClickLink.href;
					} else {
						href = this.doubleClickLink.href(e.data);
					}
					this.router.navigate([href]);
				}
			}
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
					const queryParams = {search: query ? query : null};
					this.router.navigate([], {queryParams, queryParamsHandling: "merge", relativeTo: this.route, replaceUrl: true});
				}
			);
		}
		this.cols.sort((a, b) => a.headerName === b.headerName ? 0 : ((a.headerName ?? "") > (b.headerName ?? "" ) ? -1 : 1));
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
			try {
				const filterState = localStorage.getItem(`${this.context}_table_filter`);
				if (filterState) {
					this.gridAPI.setFilterModel(JSON.parse(filterState));
				}
			} catch (e) {
				this.log.error(`Failed to retrieve stored column sort info from localStorage (key=${this.context}_table_filter:`, e);
			}
			setUpQueryParamFilter(this.route.snapshot.queryParamMap, this.cols, this.gridAPI);
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
				if (!this.columnAPI.applyColumnState(JSON.parse(colstates))) {
					this.log.error("Failed to load stored column state: one or more columns not found");
				}
			} else {
				this.gridAPI.sizeColumnsToFit();
			}
		} catch (e) {
			this.log.error(`Failure to retrieve required column info from localStorage (key=${this.context}_table_columns):`, e);
		}

	}

	/**
	 * Triggered by a table button, clears all the filters on the table.
	 */
	public clearFilters(): void {
		this.gridAPI.setFilterModel(null);
		const queryParams = Object.fromEntries(this.cols.filter((c: ColDef) => c.field).map((c: ColDef) => [c.field, null]));
		queryParams.search = null;
		this.router.navigate([], {queryParams, queryParamsHandling: "merge", relativeTo: this.route, replaceUrl: true});
	}

	/**
	 * When filter changes, stores the filter state if a context was provided,
	 * and updates query string parameters.
	 *
	 * @param e The filter change event fired by AG-Grid.
	 */
	public storeFilter(e: FilterChangedEvent<T>): void {
		if (this.context) {
			localStorage.setItem(`${this.context}_table_filter`, JSON.stringify(this.gridAPI.getFilterModel()));
		}
		// the user can only set one filter at a time, so that's all we gotta
		// handle.
		if (e.columns.length !== 1) {
			return;
		}
		const col = e.columns[0].getColDef();
		if (!col.field) {
			return;
		}
		const filter = this.gridAPI.getFilterInstance(e.columns[0]);
		if (!filter) {
			return;
		}

		let queryParams: Params = {};
		const model = filter.getModel();
		if (!model) {
			queryParams = {...this.route.snapshot.queryParams};
			queryParams[col.field] = null;
			this.router.navigate([], {queryParams, queryParamsHandling: "merge", relativeTo: this.route, replaceUrl: true});
			return;
		}
		let value = null;
		// Default filter (indicated by 'true') is the text filter
		switch(getColType(col)) {
			case "string":
				if ((model as TextFilterModel).type !== "equals") {
					return;
				}
				value = (model as TextFilterModel).filter;
				break;
			case "number":
				if ((model as NumberFilterModel).type !== "equals") {
					return;
				}
				value = (model as NumberFilterModel).filter;
				break;
			case "date":
				if ((model as DateFilterModel).type !== "equals") {
					return;
				}
				value = (model as DateFilterModel).dateFrom;
				if (value) {
					value = `${value.replace(" ", "T")}Z`;
				}
				break;
		}
		if (value === null || value === undefined) {
			return;
		}
		queryParams[col.field] = value;
		this.router.navigate([], {queryParams, queryParamsHandling: "merge", relativeTo: this.route, replaceUrl: true});
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
	 * Generates an HTML "name" attribute value for a context menu item.
	 *
	 * @param item The item for which a name will be generated.
	 * @returns A name for the given `item`. This is **not** guaranteed to be
	 * unique - in particular if any item names differ only by non-"word"
	 * characters, the output of this will collide.
	 */
	public itemName(item: ContextMenuItem<T>): string {
		return item.name.replace(/\W+/g, "-");
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
	 * @param $event The triggering dom event.
	 * @param col The ID of a column to toggle.
	 */
	public toggleVisibility($event: Event, col: string): void {
		$event.stopPropagation();
		if (this.columnAPI) {
			const column = this.columnAPI.getColumn(col);
			if (!column) {
				this.log.error(`Failed to set visibility for column '${col}': no such column`);
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
			this.log.warn("cellContextMenu fired with no underlying event");
			return;
		}

		if (!this.contextmenu) {
			this.log.warn("element reference to 'contextmenu' still null after view init");
			return;
		}

		this.showContextMenu = true;
		this.menuStyle.left = `${params.event.clientX}px`;
		this.menuStyle.top = `${params.event.clientY}px`;
		this.menuStyle.bottom = "unset";
		this.menuStyle.right = "unset";
		const boundingRect = this.contextmenu.nativeElement.getBoundingClientRect();

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
	public isDisabled(a: ContextMenuItem<T>): boolean {
		if (!this.selected) {
			return true;
		}
		if (!isAction(a)) {
			if (a.disabled) {
				return a.disabled(this.selectionCount > 1 ? this.fullSelection : this.selected);
			}
			return false;
		}
		if (!a.multiRow && this.selectionCount > 1) {
			return true;
		}
		if (a.disabled) {
			if (a.multiRow) {
				return a.disabled(this.selectionCount > 1 ? this.fullSelection : [this.selected], this.gridAPI);
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
	 * Handles when the user clicks on a title button action item by emitting the proper data.
	 *
	 * @param action The action that was clicked.
	 */
	public emitTitleButtonAction(action: string): void {
		this.tableTitleButtonAction.emit(action);
	}

	/**
	 * Handles when the user clicks on a more menu title button action item by emitting the proper data.
	 *
	 * @param action The action that was clicked.
	 */
	public emitMoreButtonAction(action: string): void {
		this.moreMenuButtonAction.emit(action);
	}

	/**
	 * Downloads the table data as a CSV file.
	 */
	public download(): void {
		const nodes = this.gridAPI.getSelectedNodes();
		const model = this.gridAPI.getModel();
		let visible = 0;
		let all = 0;
		model.forEachNode(rowNode => {
			if(rowNode.displayed) {
				visible++;
			}
			all++;
		});
		this.dialog.open(DownloadOptionsDialogComponent, {
			data: {
				allRows: all,
				columns: this.gridAPI.getColumnDefs() ?? [],
				name: this.context,
				selectedRows: nodes.length > 0 ? nodes.length : undefined,
				visibleRows: visible
			}
		}).afterClosed().subscribe(value => {
			this.gridAPI.exportDataAsCsv(value);
		});
	}

	/**
	 * Select or de-select all rows.
	 *
	 * @param de If given and true, de-selects all rows. Otherwise all rows will be selected.
	 */
	public selectAll(de?: boolean): void {
		if (de) {
			this.gridAPI.deselectAll();
			return;
		}
		this.gridAPI.selectAllFiltered();
	}

	/**
	 * Builds a link for a link context menu item.
	 *
	 * @param item The item being constructed into a link.
	 * @returns A URL or router path as determined by the settings of `item`.
	 */
	public href(item: ContextMenuLink<T>): string | Array<string> {
		if (typeof(item.href) === "string") {
			return item.href;
		}
		if (!this.selected) {
			// This happens when the context menu is hidden.
			return "";
		}
		return item.href(this.selected);
	}

	/**
	 * Gets query string parameters for a link context menu item.
	 *
	 * @param item The item being constructed into a link.
	 * @returns A set of query string parameters to pass, or `null` if the link
	 * doesn't specify any.
	 */
	public queryParameters(item: ContextMenuLink<T>): Params | ParamMap | null {
		if (!item.queryParams) {
			return null;
		}
		if (typeof(item.queryParams) !== "function") {
			return item.queryParams;
		}
		if (!this.selected) {
			// This happens when the context menu is hidden.
			return null;
		}
		return item.queryParams(this.selected);
	}

	/**
	 * Gets a URL document fragment for a link context menu item.
	 *
	 * @param item The item being constructed into a link.
	 * @returns A document fragment to pass to the routerLink, or `null` if the
	 * link doesn't specify one.
	 */
	public fragment(item: ContextMenuLink<T>): string | null {
		if (!item.fragment) {
			return null;
		}
		if (typeof(item.fragment) !== "function") {
			return item.fragment;
		}
		if (!this.selected) {
			// This happens when the context menu is hidden.
			return null;
		}
		return item.fragment(this.selected);
	}
}
