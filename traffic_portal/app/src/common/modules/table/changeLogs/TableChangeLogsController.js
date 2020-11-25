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

var TableChangeLogsController = function(tableName, changeLogs, $scope, $state, $uibModal, dateUtils, propertiesModel, messageModel) {

	/**
	 * Gets value to display a default tooltip.
	 */
	function defaultTooltip(params) {
		return params.value;
	}

	/**
	 * Formats the contents of a 'lastUpdated' column cell as "relative to now".
	 */
	function dateCellFormatterRelative(params) {
		return params.value ? dateUtils.getRelativeTime(params.value) : params.value;
	}

	function dateCellFormatter(params) {
		return params.value.toUTCString();
	}

	let columns = [
		{
			headerName: "Occurred",
			field: "lastUpdated",
			hide: false,
			filter: "agDateColumnFilter",
			tooltip: dateCellFormatterRelative,
			valueFormatter: dateCellFormatterRelative
		},
		{
			headerName: "Created (UTC)",
			field: "lastUpdated",
			hide: false,
			filter: "agDateColumnFilter",
			tooltip: dateCellFormatter,
			valueFormatter: dateCellFormatter
		},
		{
			headerName: "User",
			field: "user",
			hide: false
		},
		{
			headerName: "Level",
			field: "level",
			hide: true
		},
		{
			headerName: "Message",
			field: "message",
			hide: false
		}
	];

	$scope.days = (propertiesModel.properties.changeLogs) ? propertiesModel.properties.changeLogs.days : 7;

	/** All of the change logs - lastUpdated fields converted to actual Date */
	$scope.changeLogs = changeLogs.map(
		function(x) {
			x.lastUpdated = x.lastUpdated ? new Date(x.lastUpdated.replace("+00", "Z")) : x.lastUpdated;
		});

	$scope.quickSearch = '';

	$scope.pageSize = 100;

	/** Options, configuration, data and callbacks for the ag-grid table. */
	$scope.gridOptions = {
		columnDefs: columns,
		enableCellTextSelection: true,
		suppressMenuHide: true,
		multiSortKey: 'ctrl',
		alwaysShowVerticalScroll: true,
		defaultColDef: {
			filter: true,
			sortable: true,
			resizable: true,
			tooltip: defaultTooltip
		},
		rowData: changeLogs,
		pagination: true,
		paginationPageSize: $scope.pageSize,
		rowBuffer: 0,
		onColumnResized: function(params) {
			localStorage.setItem(tableName + "_table_columns", JSON.stringify($scope.gridOptions.columnApi.getColumnState()));
		},
		tooltipShowDelay: 500,
		allowContextMenuWithControlKey: true,
		preventDefaultOnContextMenu: true,
		onColumnVisible: function(params) {
			if (params.visible){
				return;
			}
			for (let column of params.columns) {
				if (column.filterActive) {
					const filterModel = $scope.gridOptions.api.getFilterModel();
					if (column.colId in filterModel) {
						delete filterModel[column.colId];
						$scope.gridOptions.api.setFilterModel(filterModel);
					}
				}
			}
		},
		colResizeDefault: "shift"
	};

	/** Allows the user to change the number of days queried for change logs. */
	$scope.changeDays = function() {
		const params = {
			title: 'Change Number of Days',
			message: 'Enter the number of days of change logs you need access to (between 1 and 365).',
			type: 'number'
		};
		const modalInstance = $uibModal.open({
			templateUrl: 'common/modules/dialog/input/dialog.input.tpl.html',
			controller: 'DialogInputController',
			size: 'md',
			resolve: {
				params: function () {
					return params;
				}
			}
		});
		modalInstance.result.then(function(days) {
			let numOfDays = parseInt(days, 10);
			if (numOfDays >= 1 && numOfDays <= 365) {
				propertiesModel.properties.changeLogs.days = numOfDays;
				$state.reload();
			} else {
				messageModel.setMessages([{level: 'error', text: 'Number of days must be between 1 and 365' }], false);
			}
		}, function () {
			console.log('Cancelled');
		});
	};

	/** Toggles the visibility of a column that has the ID provided as 'col'. */
	$scope.toggleVisibility = function(col) {
		const visible = $scope.gridOptions.columnApi.getColumn(col).isVisible();
		$scope.gridOptions.columnApi.setColumnVisible(col, !visible);
	};

	/** Downloads the table as a CSV */
	$scope.exportCSV = function() {
		const params = {
			allColumns: true,
			fileName: "change_logs.csv",
		};
		$scope.gridOptions.api.exportDataAsCsv(params);
	}

	$scope.onQuickSearchChanged = function() {
		$scope.gridOptions.api.setQuickFilter($scope.quickSearch);
		localStorage.setItem(tableName + "_quick_search", $scope.quickSearch);
	};

	$scope.onPageSizeChanged = function() {
		const value = Number($scope.pageSize);
		$scope.gridOptions.api.paginationSetPageSize(value);
		localStorage.setItem(tableName + "_page_size", value);
	};

	$scope.clearTableFilters = function() {
		// clear the quick search
		$scope.quickSearch = '';
		$scope.onQuickSearchChanged();
		// clear any column filters
		$scope.gridOptions.api.setFilterModel(null);
	};

	/**** Initialization code, including loading user columns from localstorage ****/
	angular.element(document).ready(function () {
		try {
			// need to create the show/hide column checkboxes and bind to the current visibility
			const colstates = JSON.parse(localStorage.getItem(tableName + "_table_columns"));
			if (colstates) {
				if (!$scope.gridOptions.columnApi.setColumnState(colstates)) {
					console.error("Failed to load stored column state: one or more columns not found");
				}
			} else {
				$scope.gridOptions.api.sizeColumnsToFit();
			}
		} catch (e) {
			console.error("Failure to retrieve required column info from localStorage (key=" + tableName + "_table_columns):", e);
		}

		try {
			const filterState = JSON.parse(localStorage.getItem(tableName + "_table_filters")) || {};
			$scope.gridOptions.api.setFilterModel(filterState);
		} catch (e) {
			console.error("Failure to load stored filter state:", e);
		}

		$scope.gridOptions.api.addEventListener("filterChanged", function() {
			localStorage.setItem(tableName + "_table_filters", JSON.stringify($scope.gridOptions.api.getFilterModel()));
		});

		try {
			const sortState = JSON.parse(localStorage.getItem(tableName + "_table_sort"));
			$scope.gridOptions.api.setSortModel(sortState);
		} catch (e) {
			console.error("Failure to load stored sort state:", e);
		}

		try {
			$scope.quickSearch = localStorage.getItem(tableName + "_quick_search");
			$scope.gridOptions.api.setQuickFilter($scope.quickSearch);
		} catch (e) {
			console.error("Failure to load stored quick search:", e);
		}

		try {
			const ps = localStorage.getItem(tableName + "_page_size");
			if (ps && ps > 0) {
				$scope.pageSize = Number(ps);
				$scope.gridOptions.api.paginationSetPageSize($scope.pageSize);
			}
		} catch (e) {
			console.error("Failure to load stored page size:", e);
		}

		$scope.gridOptions.api.addEventListener("sortChanged", function() {
			localStorage.setItem(tableName + "_table_sort", JSON.stringify($scope.gridOptions.api.getSortModel()));
		});

		$scope.gridOptions.api.addEventListener("columnMoved", function() {
			localStorage.setItem(tableName + "_table_columns", JSON.stringify($scope.gridOptions.columnApi.getColumnState()));
		});

		$scope.gridOptions.api.addEventListener("columnVisible", function() {
			$scope.gridOptions.api.sizeColumnsToFit();
			try {
				colStates = $scope.gridOptions.columnApi.getColumnState();
				localStorage.setItem(tableName + "_table_columns", JSON.stringify(colStates));
			} catch (e) {
				console.error("Failed to store column defs to local storage:", e);
			}
		});

	});

};

TableChangeLogsController.$inject = ['tableName', 'changeLogs', '$scope', '$state', '$uibModal', 'dateUtils', 'propertiesModel', 'messageModel'];
module.exports = TableChangeLogsController;
