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

var TableCertExpirationsController = function(tableName, certExpirations, filter, $document, $scope, $state, $filter, locationUtils, certExpirationsService) {

    let table;

	$scope.certExpirations = certExpirations;
	$scope.days;

	$scope.editCertExpirations = function(dsId) {
		locationUtils.navigateToPath('/delivery-services/' + dsId + '/ssl-keys');
	}

	/**
	 * Formats the contents of a 'expiration' column cell as just the date.
	 */
	function dateCellFormatter(params) {
		return params.value ? $filter('date')(params.value, 'MM/dd/yyyy') : params.value;
	}

	/**
	 * Formats the contents of a 'expiration' column cell as just the date.
	 */
	function federatedCellFormatter(params) {
		if (!params.value) {
			return '';
		} else {
			return params.value;
		}
	}

    $scope.refresh = function() {
        $state.reload(); // reloads all the resolves for the view
    };

	/** Toggles the visibility of a column that has the ID provided as 'col'. */
	$scope.toggleVisibility = function(col) {
		const visible = $scope.gridOptions.columnApi.getColumn(col).isVisible();
		$scope.gridOptions.columnApi.setColumnVisible(col, !visible);
	};

	/** The columns of the ag-grid table */
	const columns = [
		{
			headerName: "Delivery Service",
			field: "deliveryservice",
			hide: false
		},
		{
			headerName: "CDN",
			field: "cdn",
			hide: false
		},
		{
			headerName: "Provider",
			field: "provider",
			hide: false
		},
		{
			headerName: "Expiration",
			field: "expiration",
			hide: false,
			valueFormatter: dateCellFormatter
		},
		{
			headerName: "Federated",
			field: "federated",
			hide: false,
			valueFormatter: federatedCellFormatter
		},
	];

	$scope.pageSize = 100;

	/** Options, configuration, data and callbacks for the ag-grid table. */
	$scope.gridOptions = {
		columnDefs: columns,
		enableCellTextSelection:true,
		suppressMenuHide: true,
		multiSortKey: 'ctrl',
		alwaysShowVerticalScroll: true,
		defaultColDef: {
			filter: true,
			sortable: true,
			resizable: true,
		},
		rowData: $scope.certExpirations,
		pagination: true,
		paginationPageSize: $scope.pageSize,
		rowBuffer: 0,
		allowContextMenuWithControlKey: true,
		preventDefaultOnContextMenu: true,
		colResizeDefault: "shift",
		onColumnVisible: function(params) {
			if (params.visible){
				return;
			}
			const filterModel = $scope.gridOptions.api.getFilterModel();
			for (let column of params.columns) {
				if (column.filterActive) {
					if (column.colId in filterModel) {
						delete filterModel[column.colId];
						$scope.gridOptions.api.setFilterModel(filterModel);
					}
				}
			}
		},
		onRowClicked: function(params) {
			const selection = window.getSelection().toString();
			if(selection === "" || selection === $scope.mouseDownSelectionText) {
				locationUtils.navigateToPath('/delivery-services/' + params.data.deliveryservice_id + '/ssl-keys');
				// Event is outside the digest cycle, so we need to trigger one.
				$scope.$apply();
			}
			$scope.mouseDownSelectionText = "";
		},
		onColumnResized: function(params) {
			localStorage.setItem(tableName + "_table_columns", JSON.stringify($scope.gridOptions.columnApi.getColumnState()));
		},
		onFirstDataRendered: function(event) {
			try {
				const filterState = JSON.parse(localStorage.getItem(tableName + "_table_filters")) || {};
				// apply any filter provided to the controller
				Object.assign(filterState, filter);
				$scope.gridOptions.api.setFilterModel(filterState);
			} catch (e) {
				console.error("Failure to load stored filter state:", e);
			}

			$scope.gridOptions.api.addEventListener("filterChanged", function() {
				localStorage.setItem(tableName + "_table_filters", JSON.stringify($scope.gridOptions.api.getFilterModel()));
			});
		},
		onGridReady: function() {
			try { // need to create the show/hide column checkboxes and bind to the current visibility
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

			try {
				const days = localStorage.getItem(tableName + "_days");
				if (days && days > 0) {
					$scope.days = Number(days);
					certExpirationsService.getCertExpirationsDaysLimit(days)
						.then(function (response) {
							newData = response;
							$scope.gridOptions.api.setRowData(newData);
						});
				} else {
					$scope.days = undefined;
					certExpirationsService.getCertExpirations()
						.then(function (response) {
							newData = response;
							$scope.gridOptions.api.setRowData(newData);
						});
				}
			} catch (e) {
				console.error("Failure to load stored days:", e);
			}

			try {
				const page = parseInt(localStorage.getItem(tableName + "_table_page"));
				const totalPages = $scope.gridOptions.api.paginationGetTotalPages();
				if (page !== undefined && page > 0 && page <= totalPages-1) {
					$scope.gridOptions.api.paginationGoToPage(page);
				}
			} catch (e) {
				console.error("Failed to load stored page number:", e);
			}

			$scope.gridOptions.api.addEventListener("paginationChanged", function() {
				localStorage.setItem(tableName + "_table_page", $scope.gridOptions.api.paginationGetCurrentPage());
			});

			$scope.gridOptions.api.addEventListener("sortChanged", function() {
				localStorage.setItem(tableName + "_table_sort", JSON.stringify($scope.gridOptions.api.getSortModel()));
			});

			$scope.gridOptions.api.addEventListener("columnMoved", function() {
				localStorage.setItem(tableName + "_table_columns", JSON.stringify($scope.gridOptions.columnApi.getColumnState()));
			});

			$scope.gridOptions.api.addEventListener("columnVisible", function() {
				$scope.gridOptions.api.sizeColumnsToFit();
				try {
					const colStates = $scope.gridOptions.columnApi.getColumnState();
					localStorage.setItem(tableName + "_table_columns", JSON.stringify(colStates));
				} catch (e) {
					console.error("Failed to store column defs to local storage:", e);
				}
			});
		}
	};

	/** This is used to position the context menu under the cursor. */
	$scope.menuStyle = {
		left: 0,
		top: 0,
	};

	/** Controls whether or not the context menu is visible. */
	$scope.showMenu = false;

	/** Downloads the table as a CSV */
	$scope.exportCSV = function() {
		const params = {
			allColumns: true,
			fileName: "certificate-expirations.csv",
		};
		$scope.gridOptions.api.exportDataAsCsv(params);
	};

	$scope.onQuickSearchChanged = function() {
		$scope.gridOptions.api.setQuickFilter($scope.quickSearch);
		localStorage.setItem(tableName + "_quick_search", $scope.quickSearch);
	};

	$scope.onPageSizeChanged = function() {
		const value = Number($scope.pageSize);
		$scope.gridOptions.api.paginationSetPageSize(value);
		localStorage.setItem(tableName + "_page_size", value);
	};

	$scope.onDaysChanged = function(days) {
		if (days && days > 0) {
			const value = Number(days);
			$scope.days = value;
			certExpirationsService.getCertExpirationsDaysLimit(value)
				.then(function (response) {
					newData = response;
					$scope.gridOptions.api.setRowData(newData);
					localStorage.setItem(tableName + "_days", value);
				});
		} else {
			$scope.days = undefined;
			certExpirationsService.getCertExpirations()
				.then(function (response) {
					newData = response;
					$scope.gridOptions.api.setRowData(newData);
					localStorage.setItem(tableName + "_days", undefined);
				});
		}
	}

	$scope.clearTableFilters = function() {
		// clear the quick search
		$scope.quickSearch = '';
		$scope.onQuickSearchChanged();
		$scope.days = undefined;
		$scope.gridOptions.api.setRowData(certExpirations);

		// clear any column filters
		$scope.gridOptions.api.setFilterModel(null);
	};

	angular.element(document).ready(function () {
		// clicks outside the context menu will hide it
		$document.bind("click", function(e) {
			$scope.showMenu = false;
			e.stopPropagation();
			$scope.$apply();
		});
	});

};

TableCertExpirationsController.$inject = ['tableName', 'certExpirations', 'filter', '$document', '$scope', '$state', '$filter', 'locationUtils', 'certExpirationsService'];
module.exports = TableCertExpirationsController;
