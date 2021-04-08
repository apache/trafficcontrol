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

var TableNotificationsController = function(tableName, notifications, filter, $scope, $state, $uibModal, $document, dateUtils, cdnService) {

	/**
	 * Gets value to display a default tooltip.
	 */
	function defaultTooltip(params) {
		return params.value;
	}

	function dateCellFormatter(params) {
		return params.value.toUTCString();
	}

	let columns = [
		{
			headerName: "Created (UTC)",
			field: "lastUpdated",
			hide: false,
			filter: "agDateColumnFilter",
			tooltipValueGetter: dateCellFormatter,
			valueFormatter: dateCellFormatter
		},
		{
			headerName: "User",
			field: "user",
			hide: false
		},
		{
			headerName: "CDN",
			field: "cdn",
			hide: false
		},
		{
			headerName: "Notification",
			field: "notification",
			hide: false
		}
	];

	/** All of the notifications - lastUpdated fields converted to actual Date */
	$scope.notifications = notifications.map(
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
			tooltipValueGetter: defaultTooltip
		},
		rowData: notifications,
		pagination: true,
		paginationPageSize: $scope.pageSize,
		rowBuffer: 0,
		onColumnResized: function(params) {
			localStorage.setItem(tableName + "_table_columns", JSON.stringify($scope.gridOptions.columnApi.getColumnState()));
		},
		tooltipShowDelay: 500,
		allowContextMenuWithControlKey: true,
		preventDefaultOnContextMenu: true,
		onCellContextMenu: function(params) {
			$scope.showMenu = true;
			$scope.menuStyle.left = String(params.event.clientX) + "px";
			$scope.menuStyle.top = String(params.event.clientY) + "px";
			$scope.menuStyle.bottom = "unset";
			$scope.menuStyle.right = "unset";
			$scope.$apply();
			const boundingRect = document.getElementById("context-menu").getBoundingClientRect();

			if (boundingRect.bottom > window.innerHeight){
				$scope.menuStyle.bottom = String(window.innerHeight - params.event.clientY) + "px";
				$scope.menuStyle.top = "unset";
			}
			if (boundingRect.right > window.innerWidth) {
				$scope.menuStyle.right = String(window.innerWidth - params.event.clientX) + "px";
				$scope.menuStyle.left = "unset";
			}
			$scope.notification = params.data;
			$scope.$apply();
		},
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
		onFirstDataRendered: function() {
			try {
				const filterState = JSON.parse(localStorage.getItem(tableName + "_table_filters")) || {};
				$scope.gridOptions.api.setFilterModel(filterState);
			} catch (e) {
				console.error("Failure to load stored filter state:", e);
			}

			$scope.gridOptions.api.addEventListener("filterChanged", function() {
				localStorage.setItem(tableName + "_table_filters", JSON.stringify($scope.gridOptions.api.getFilterModel()));
			});
		},
		onGridReady: function() {
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
		},
		colResizeDefault: "shift"
	};

	/** This is used to position the context menu under the cursor. */
	$scope.menuStyle = {
		left: 0,
		top: 0,
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
			fileName: "invalidation_requests.csv",
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

	$scope.selectCDNandCreateNotification = function() {
		const params = {
			title: 'Create Notification',
			message: "Please select a CDN"
		};
		const modalInstance = $uibModal.open({
			templateUrl: 'common/modules/dialog/select/dialog.select.tpl.html',
			controller: 'DialogSelectController',
			size: 'md',
			resolve: {
				params: function () {
					return params;
				},
				collection: function(cdnService) {
					return cdnService.getCDNs();
				}
			}
		});
		modalInstance.result.then(function(cdn) {
			$scope.createNotification(cdn);
		}, function () {
			// do nothing
		});
	};

	$scope.createNotification = function(cdn) {
		const params = {
			title: 'Create ' + cdn.name + ' Notification',
			message: 'What is the content of your notification for the ' + cdn.name + ' CDN?'
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
		modalInstance.result.then(function(notification) {
			cdnService.createNotification(cdn, notification).
			then(
				function() {
					$state.reload();
				}
			);
		}, function () {
			// do nothing
		});
	};

	$scope.confirmDeleteNotification = function(notification, event) {
		event.stopPropagation();
		const params = {
			title: 'Delete Notification',
			message: 'Are you sure you want to delete the notification for the ' + notification.cdn + ' CDN? This will remove the notification from the view of all users.'
		};
		const modalInstance = $uibModal.open({
			templateUrl: 'common/modules/dialog/confirm/dialog.confirm.tpl.html',
			controller: 'DialogConfirmController',
			size: 'md',
			resolve: {
				params: function () {
					return params;
				}
			}
		});
		modalInstance.result.then(function() {
			cdnService.deleteNotification({ id: notification.id }).
				then(
					function() {
						$state.reload();
					}
				);
		}, function () {
			// do nothing
		});
	};

	/**** Initialization code, including loading user columns from localstorage ****/
	angular.element(document).ready(function () {
		// clicks outside the context menu will hide it
		$document.bind("click", function(e) {
			$scope.showMenu = false;
			e.stopPropagation();
			$scope.$apply();
		});
	});

};

TableNotificationsController.$inject = ['tableName', 'notifications', 'filter', '$scope', '$state', '$uibModal', '$document', 'dateUtils', 'cdnService'];
module.exports = TableNotificationsController;
