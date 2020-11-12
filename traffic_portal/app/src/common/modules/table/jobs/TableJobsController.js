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

var TableJobsController = function(tableName, jobs, $document, $scope, $state, $uibModal, locationUtils, jobService, messageModel, dateUtils) {

	/**
	 * Gets value to display a default tooltip.
	 */
	function defaultTooltip(params) {
		return params.value;
	}

	function dateCellFormatter(params) {
		return params.value.toUTCString();
	}

	columns = [
		{
			headerName: "Delivery Service",
			field: "deliveryService",
			hide: false
		},
		{
			headerName: "Asset URL",
			field: "assetUrl",
			hide: false
		},
		{
			headerName: "Parameters",
			field: "parameters",
			hide: false
		},
		{
			headerName: "Start (UTC)",
			field: "startTime",
			hide: false,
			filter: "agDateColumnFilter",
			tooltip: dateCellFormatter,
			valueFormatter: dateCellFormatter
		},
		{
			headerName: "Expires (UTC)",
			field: "expires",
			hide: false,
			filter: "agDateColumnFilter",
			tooltip: dateCellFormatter,
			valueFormatter: dateCellFormatter
		},
		{
			headerName: "Created By",
			field: "createdBy",
			hide: false
		}
	];

	/** All of the jobs - startTime fields converted to actual Dates and derived expires field from TTL */
	$scope.jobs = jobs.map(
		function(x) {
			// need to convert this to a date object for ag-grid filter to work properly
			x.startTime = new Date(x.startTime.replace("+00", "Z"));

			// going to derive the expires date from start + TTL (hours). Format: TTL:24h
			let ttl = parseInt(x.parameters.slice('TTL:'.length, x.parameters.length-1), 10);
			x.expires = new Date(x.startTime.getTime() + ttl*3600*1000);
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
		rowClassRules: {
			'active-job': function(params) {
				return params.data.expires > new Date();
			},
			'expired-job': function(params) {
				return params.data.expires <= new Date();
			}
		},
		rowData: jobs,
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
			$scope.job = params.data;
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

	$scope.createJob = function() {
		locationUtils.navigateToPath('/jobs/new');
	};

	$scope.refresh = function() {
		$state.reload(); // reloads all the resolves for the view
	};

	$scope.confirmRemoveJob = function(job, $event) {
		$event.stopPropagation();
		const params = {
			title: 'Remove Invalidation Request?',
			message: 'Are you sure you want to remove the ' + job.assetUrl + ' invalidation request?<br><br>' +
				'NOTE: The invalidation request may have already been applied.'
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
			jobService.deleteJob(job.id)
				.then(
					function(result) {
						messageModel.setMessages(result.data.alerts, false);
						$scope.refresh(); // refresh the jobs table
					}
				);
		});
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

		// clicks outside the context menu will hide it
		$document.bind("click", function(e) {
			$scope.showMenu = false;
			e.stopPropagation();
			$scope.$apply();
		});
	});

};

TableJobsController.$inject = ['tableName', 'jobs', '$document', '$scope', '$state', '$uibModal', 'locationUtils', 'jobService', 'messageModel', 'dateUtils'];
module.exports = TableJobsController;
