/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

/** @typedef { import('./CommonGridController').CGC } CGC */

/**
 * Given some query parameters, the columns of a table, and a hook into the
 * AG-Grid API of said table, sets up filtering based on matches between the
 * names of query parameters and the raw data fields of the columns.
 *
 * @param {URLSearchParams} params
 * @param {{field?: string; filter?: string}[]} columns
 * @param {GridApi} api
 */
function setUpQueryParamFilter(params, columns, api) {
	for (const col of columns) {
		if (!Object.prototype.hasOwnProperty.call(col, "field")) {
			continue;
		}
		const filter = api.getFilterInstance(col.field);
		if (!filter) {
			continue;
		}
		const values = params.getAll(col.field);
		if (values.length < 1) {
			continue;
		}

		/** @type {"string" | "number" | "date"} */
		let colType;
		if (!Object.prototype.hasOwnProperty.call(col, "filter")) {
			colType = "string";
		} else if (typeof(col.filter) !== "string") {
			continue;
		} else {
			let bail = false;
			switch(col.filter) {
				case "agTextColumnFilter":
					colType = "string";
					break;
				case "agNumberColumnFilter":
					colType = "number";
					break;
				case "agDateColumnFilter":
					colType = "date";
					break;
				default:
					bail = true;
					break;
			}
			if (bail) {
				continue;
			}
		}

		let filterModel;
		switch(colType) {
			case "string":
				if (values.length === 1) {
					filterModel = {
						filter: values[0],
						type: "equals"
					}
				} else {
					filterModel = {
						operator: "OR",
						condition1: {
							filter: values[0],
							type: "equals"
						},
						condition2: {
							filter: values[1],
							type: "equals"
						}
					}
				}
				break;
			case "number":
				if (values.length === 1) {
					filterModel = {
						filter: parseInt(values[0], 10),
						type: "equals"
					}
					if (isNaN(filterModel.filter)) {
						continue;
					}
				} else {
					filterModel = {
						operator: "OR",
						condition1: {
							filter: parseInt(values[0], 10),
							type: "equals"
						},
						condition2: {
							filter: parseInt(values[1], 10),
							type: "equals"
						}
					}
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
				const pad = num => String(num).padStart(2,"0");
				filterModel = {
					dateFrom: `${date.getUTCFullYear()}-${pad(date.getUTCMonth()+1)}-${pad(date.getUTCDate())} ${pad(date.getUTCHours())}:${pad(date.getUTCMinutes())}:${pad(date.getUTCSeconds())}`,
					type: "equals"
				}
				break;
		}
		filter.setModel(filterModel);
		filter.applyModel();
	}
}

/**
 * @param {*} $scope
 * @param {import("angular").IDocumentService} $document
 * @param {*} $state
 * @param {import("../../../models/UserModel")} userModel
 * @param {import("../../../service/utils/DateUtils")} dateUtils
 */
let CommonGridController = function ($scope, $document, $state, userModel, dateUtils) {
	this.entry = null;
	this.quickSearch = "";
	this.pageSize = 100;
	this.showMenu = false;
	/**
	 * @type {{
	 * 	bottom?: string | 0;
	 * 	left: string | 0;
	 * 	right?: string | 0;
	 * 	top: string | 0;
	 * }}
	 */
	this.menuStyle = {
		left: 0,
		top: 0
	};
	this.mouseDownSelectionText = "";

	// Bound Variables
	/** @type string */
	this.tableTitle = "";
	/** @type string */
	this.tableName = "";
	/** @type CGC.GridSettings */
	this.options = {};
	/** @type any */
	this.gridOptions = {};
	/** @type any[] */
	this.columns = [];
	/** @type string[] */
	this.sensitiveColumns = [];
	/** @type any[] */
	this.data = [];
	/** @type any[] */
	this.selectedData = [];
	/** @type any */
	this.defaultData = {};
	/** @type CGC.DropDownOption[] */
	this.dropDownOptions = [];
	/** @type CGC.ContextMenuOption[] */
	this.contextMenuOptions = [];
	/** @type CGC.TitleButton */
	this.titleButton = {};
	/** @type CGC.TitleBreadCrumbs */
	this.breadCrumbs = [];

	function HTTPSCellRenderer() {}
	HTTPSCellRenderer.prototype.init = function(params) {
		this.eGui = document.createElement("a");
		this.eGui.href = "https://" + params.value;
		this.eGui.setAttribute("class", "link");
		this.eGui.setAttribute("target", "_blank");
		this.eGui.textContent = params.value;
	};
	HTTPSCellRenderer.prototype.getGui = function() {return this.eGui;};

	// browserify can't handle classes...
	function SSHCellRenderer() {}
	SSHCellRenderer.prototype.init = function(params) {
		this.eGui = document.createElement("a");
		this.eGui.href = "ssh://" + userModel.user.username + "@" + params.value;
		this.eGui.setAttribute("class", "link");
		this.eGui.textContent = params.value;
	};
	SSHCellRenderer.prototype.getGui = function() {return this.eGui;};

	function CheckCellRenderer() {}
	CheckCellRenderer.prototype.init = function(params) {
		this.eGui = document.createElement("i");
		if (params.value === null || params.value === undefined) {
			return;
		}

		this.eGui.setAttribute("aria-hidden", "true");
		this.eGui.setAttribute("title", String(params.value));
		this.eGui.classList.add("fa", "fa-lg");
		if (params.value) {
			this.eGui.classList.add("fa-check");
		} else {
			this.eGui.classList.add("fa-times");
		}
	};
	CheckCellRenderer.prototype.getGui = function() {return this.eGui;};

	function UpdateCellRenderer() {}
	UpdateCellRenderer.prototype.init = function(params) {
		this.eGui = document.createElement("i");

		this.eGui.setAttribute("aria-hidden", "true");
		this.eGui.setAttribute("title", String(params.value));
		this.eGui.classList.add("fa", "fa-lg");
		if (params.value) {
			this.eGui.classList.add("fa-clock-o");
		} else {
			this.eGui.classList.add("fa-check");
		}
	};
	UpdateCellRenderer.prototype.getGui = function() {return this.eGui;};

	function defaultTooltip(params) {
		return params.value;
	}

	function dateCellFormatterRelative(params) {
		return params.value ? dateUtils.getRelativeTime(params.value) : params.value;
	}

	function dateCellFormatterUTC(params) {
		return params.value ? params.value.toUTCString() : params.value;
	}

	this.hasContextItems = function() {
		return this.contextMenuOptions.length > 0;
	};

	this.hasSensitiveColumns = function() {
		return this.sensitiveColumns.length > 0;
	}

	/**
	 * @param {string} colID
	 */
	this.isSensitive = function(colID) {
		return this.sensitiveColumns.includes(colID);
	}

	this.sensitiveColumnsShown = false;

	this.toggleSensitiveFields = function() {
		if (this.sensitiveColumnsShown) {
			return;
		}
		for (const col of this.gridOptions.columnApi.getAllColumns()) {
			const id = col.getColId();
			if (this.isSensitive(id)) {
				this.gridOptions.columnApi.setColumnVisible(id, false);
			}
		}
	};

	this.getColumns = () => {
		/** @type {{colId: string}[]} */
		const cols = this.gridOptions.columnApi.getAllColumns();
		if (!this.hasSensitiveColumns || this.sensitiveColumnsShown) {
			return cols;
		}
		return cols.filter(c => !this.isSensitive(c.colId));
	}

	this.$onInit = () => {
		const tableName = this.tableName;

		if (this.defaultData !== undefined) {
			this.entry = this.defaultData;
		}

		for(let i = 0; i < this.columns.length; ++i) {
			if (this.columns[i].filter === "agDateColumnFilter") {
				if (this.columns[i].relative) {
					this.columns[i].tooltipValueGetter = dateCellFormatterRelative;
					this.columns[i].valueFormatter = dateCellFormatterRelative;
				}
				else {
					this.columns[i].tooltipValueGetter = dateCellFormatterUTC;
					this.columns[i].valueFormatter = dateCellFormatterUTC;
				}
			} else if (this.columns[i].filter === 'arrayTextColumnFilter') {
				this.columns[i].filter = 'agTextColumnFilter'
				this.columns[i].filterParams = {
					textCustomComparator: (filter, value, filterText) => {
						const filterTextLowerCase = filterText.toLowerCase();
						const valueLowerCase = value.toString().toLowerCase();
						const profileNameValue = valueLowerCase.split(",");
						switch (filter) {
							case 'contains':
								return valueLowerCase.indexOf(filterTextLowerCase) >= 0;
							case 'notContains':
								return valueLowerCase.indexOf(filterTextLowerCase) === -1;
							case 'equals':
								return profileNameValue.includes(filterTextLowerCase);
							case 'notEqual':
								return !profileNameValue.includes(filterTextLowerCase);
							case 'startsWith':
								return valueLowerCase.indexOf(filterTextLowerCase) === 0;
							case 'endsWith':
								let index = valueLowerCase.lastIndexOf(filterTextLowerCase);
								return index >= 0 && index === (valueLowerCase.length - filterTextLowerCase.length);
							default:
								// should never happen
								console.warn('invalid filter type ' + filter);
								return false;
						}
					}
				}
			}
		}

		// clicks outside the context menu will hide it
		$document.bind("click", e => {
			this.showMenu = false;
			e.stopPropagation();
			$scope.$apply();
		});

		this.gridOptions = {
			components: {
				httpsCellRenderer: HTTPSCellRenderer,
				sshCellRenderer: SSHCellRenderer,
				updateCellRenderer: UpdateCellRenderer,
				checkCellRenderer: CheckCellRenderer,
			},
			columnDefs: this.columns,
			enableCellTextSelection: true,
			suppressMenuHide: true,
			multiSortKey: 'ctrl',
			alwaysShowVerticalScroll: true,
			defaultColDef: {
				filter: true,
				sortable: true,
				resizable: true,
				tooltipValueGetter: defaultTooltip
			},
			rowClassRules: this.options.rowClassRules,
			rowData: this.data,
			pagination: true,
			paginationPageSize: this.pageSize,
			rowBuffer: 0,
			onColumnResized: () => {
				/** @type {{colId: string; hide?: boolean | null}[]} */
				const states = this.gridOptions.columnApi.getColumnState();
				for (const state of states) {
					state.hide = state.hide || this.isSensitive(state.colId);
				}
				localStorage.setItem(tableName + "_table_columns", JSON.stringify(states));
			},
			colResizeDefault: "shift",
			tooltipShowDelay: 500,
			allowContextMenuWithControlKey: true,
			preventDefaultOnContextMenu: this.hasContextItems(),
			onCellMouseDown: () => {
				const selection = window.getSelection();
				if (!selection) {
					this.mouseDownSelectionText = "";
				} else {
					this.mouseDownSelectionText = selection.toString();
				}
			},
			onCellContextMenu: params => {
				if (!this.hasContextItems()){
					return;
				}
				this.showMenu = true;
				this.menuStyle.left = String(params.event.clientX) + "px";
				this.menuStyle.top = String(params.event.clientY) + "px";
				this.menuStyle.bottom = "unset";
				this.menuStyle.right = "unset";
				$scope.$apply();
				const boundingRect = document.getElementById("context-menu")?.getBoundingClientRect();
				if (!boundingRect) {
					throw new Error("no bounding rectangle for context-menu; element possibly missing");
				}

				if (boundingRect.bottom > window.innerHeight){
					this.menuStyle.bottom = String(window.innerHeight - params.event.clientY) + "px";
					this.menuStyle.top = "unset";
				}
				if (boundingRect.right > window.innerWidth) {
					this.menuStyle.right = String(window.innerWidth - params.event.clientX) + "px";
					this.menuStyle.left = "unset";
				}
				this.entry = params.data;
				$scope.$apply();
			},
			onColumnVisible: params => {
				if (params.visible){
					return;
				}
				for (let column of params.columns) {
					if (column.filterActive) {
						const filterModel = this.gridOptions.api.getFilterModel();
						if (column.colId in filterModel) {
							delete filterModel[column.colId];
							this.gridOptions.api.setFilterModel(filterModel);
						}
					}
				}
			},
			onRowSelected: () => {
				this.selectedData = this.gridOptions.api.getSelectedRows();
				$scope.$apply();
			},
			onSelectionChanged: () => {
				this.selectedData = this.gridOptions.api.getSelectedRows();
				$scope.$apply();
			},
			onRowClicked: params => {
				if (params.event.target instanceof HTMLAnchorElement) {
					return;
				}
				const selection = window.getSelection();
				if (this.options.onRowClick !== undefined) {
					if (!selection || selection.toString() === "" || selection === $scope.mouseDownSelectionText) {
						this.options.onRowClick(params);
						$scope.$apply();
					}
				}
				$scope.mouseDownSelectionText = "";
			},
			onFirstDataRendered: () => {
				if(this.options.selectRows) {
					this.gridOptions.rowSelection = this.options.selectRows ? "multiple" : "";
					this.gridOptions.rowMultiSelectWithClick = this.options.selectRows;
					this.gridOptions.api.forEachNode(node => {
						if (node.data[this.options.selectionProperty] === true) {
							node.setSelected(true, false);
						}
					});
				}
				try {
					const filterState = JSON.parse(localStorage.getItem(tableName + "_table_filters") ?? "{}") || {};
					this.gridOptions.api.setFilterModel(filterState);
				} catch (e) {
					console.error("Failure to load stored filter state:", e);
				}
				// Set up filters from query string paramters.
				const params = new URLSearchParams(globalThis.location.hash.split("?").slice(1).join("?"));
				setUpQueryParamFilter(params, this.columns, this.gridOptions.api);
				this.gridOptions.api.onFilterChanged();

				this.gridOptions.api.addEventListener("filterChanged", () => {
					localStorage.setItem(tableName + "_table_filters", JSON.stringify(this.gridOptions.api.getFilterModel()));
				});
			},
			onGridReady: () => {
				try {
					// need to create the show/hide column checkboxes and bind to the current visibility
					const colstates = JSON.parse(localStorage.getItem(tableName + "_table_columns") ?? "null");
					if (colstates) {
						if (!this.gridOptions.columnApi.setColumnState(colstates)) {
							console.error("Failed to load stored column state: one or more columns not found");
						}
					} else {
						this.gridOptions.api.sizeColumnsToFit();
					}
				} catch (e) {
					console.error("Failure to retrieve required column info from localStorage (key=" + tableName + "_table_columns):", e);
				}

				try {
					const sortState = JSON.parse(localStorage.getItem(tableName + "_table_sort") ?? "{}");
					this.gridOptions.api.setSortModel(sortState);
				} catch (e) {
					console.error("Failure to load stored sort state:", e);
				}

				try {
					this.quickSearch = localStorage.getItem(tableName + "_quick_search") ?? "";
					this.gridOptions.api.setQuickFilter(this.quickSearch);
				} catch (e) {
					console.error("Failure to load stored quick search:", e);
				}

				try {
					const ps = Number(localStorage.getItem(tableName + "_page_size"));
					if (ps > 0) {
						this.pageSize = Number(ps);
						this.gridOptions.api.paginationSetPageSize(this.pageSize);
					}
				} catch (e) {
					console.error("Failure to load stored page size:", e);
				}

				try {
					const page = parseInt(localStorage.getItem(tableName + "_table_page") ?? "0", 10);
					if (page > 0 && page <= $scope.gridOptions.api.paginationGetTotalPages()-1) {
						$scope.gridOptions.api.paginationGoToPage(page);
					}
				} catch (e) {
					console.error("Failed to load stored page number:", e);
				}

				this.gridOptions.api.addEventListener("sortChanged", () => {
					localStorage.setItem(tableName + "_table_sort", JSON.stringify(this.gridOptions.api.getSortModel()));
				});

				this.gridOptions.api.addEventListener("columnMoved", () => {
					/** @type {{colId: string; hide?: boolean | null}[]} */
					const states = this.gridOptions.columnApi.getColumnState();
					for (const state of states) {
						state.hide = state.hide || this.isSensitive(state.colId);
					}

					localStorage.setItem(tableName + "_table_columns", JSON.stringify(this.gridOptions.columnApi.getColumnState()));
				});

				this.gridOptions.api.addEventListener("columnVisible", () => {
					this.gridOptions.api.sizeColumnsToFit();
					try {
						const colStates = this.gridOptions.columnApi.getColumnState();
						localStorage.setItem(tableName + "_table_columns", JSON.stringify(colStates));
					} catch (e) {
						console.error("Failed to store column defs to local storage:", e);
					}
				});
			}
		};

	};

	this.exportCSV = function() {
		const params = {
			allColumns: true,
			fileName: this.tableName + ".csv",
		};
		this.gridOptions.api.exportDataAsCsv(params);
	};

	this.toggleVisibility = function(col) {
		const visible = this.gridOptions.columnApi.getColumn(col).isVisible();
		this.gridOptions.columnApi.setColumnVisible(col, !visible);
	};

	this.onQuickSearchChanged = function() {
		this.gridOptions.api.setQuickFilter(this.quickSearch);
		localStorage.setItem(this.tableName + "_quick_search", this.quickSearch);
	};

	this.onPageSizeChanged = function() {
		const value = Number(this.pageSize);
		this.gridOptions.api.paginationSetPageSize(value);
		localStorage.setItem(this.tableName + "_page_size", value.toString());
	};

	this.clearTableFilters = () => {
		// clear the quick search
		this.quickSearch = '';
		this.onQuickSearchChanged();
		// clear any column filters
		this.gridOptions.api.setFilterModel(null);
	};

	this.contextMenuClick = function(menu, $event) {
		$event.stopPropagation();
		menu.onClick(this.entry);
	};

	this.getHref = function(menu) {
		if (menu.href !== undefined){
			return menu.href;
		}
		return menu.getHref(this.entry);
	};

	this.contextIsDisabled = function(menu) {
		if (menu.isDisabled !== undefined) {
			return menu.isDisabled(this.entry);
		}
		return false;
	};

	this.bcGetText = function (bc) {
		if(bc.text !== undefined){
			return bc.text;
		}
		return bc.getText();
	};

	this.bcHasHref = function(bc) {
		return bc.href !== undefined || bc.getHref !== undefined;
	};

	this.bcGetHref = function(bc) {
		if(bc.href !== undefined) {
			return bc.href;
		}
		return bc.getHref();
	};

	this.getText = function (menu) {
		if (menu.text !== undefined){
			return menu.text;
		}
		return menu.getText(this.entry);
	};

	this.isShown = function (menu) {
		if (menu.shown === undefined){
			return true;
		}
		return menu.shown(this.entry);
	};

	$scope.refresh = function() {
		$state.reload(); // reloads all the resolves for the view
	};
};

angular.module("trafficPortal.table").component("commonGridController", {
	templateUrl: "common/modules/table/agGrid/grid.tpl.html",
	controller: CommonGridController,
	bindings: {
		tableTitle: "@",
		tableName: "@",
		options: "<",
		columns: "<",
		data: "<",
		selectedData: "=?",
		dropDownOptions: "<?",
		contextMenuOptions: "<?",
		defaultData: "<?",
		titleButton: "<?",
		breadCrumbs: "<?",
		sensitiveColumns: "<?"
	}
});

CommonGridController.$inject = ["$scope", "$document", "$state", "userModel", "dateUtils"];
module.exports = CommonGridController;
