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

var TableServersController = function(tableName, servers, filter, $scope, $state, $uibModal, $window, dateUtils, locationUtils, serverUtils, cdnService, serverService, statusService, propertiesModel, messageModel, userModel, $document) {
	/**** Table cell formatters/renderers ****/

	// browserify can't handle classes...
	function SSHCellRenderer() {}
	SSHCellRenderer.prototype.init = function(params) {
		this.eGui = document.createElement("A");
		this.eGui.href = "ssh://" + userModel.user.username + "@" + params.value;
		this.eGui.setAttribute("target", "_blank");
		this.eGui.textContent = params.value;
	};
	SSHCellRenderer.prototype.getGui = function() {return this.eGui;};

	function UpdateCellRenderer() {}
	UpdateCellRenderer.prototype.init = function(params) {
		this.eGui = document.createElement("I");
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

	/**
	 * Gets text with which to file a status tooltip.
	 * @returns {string | undefined} The offline reason if the server is offline, otherwise nothing.
	 */
	function offlineReasonTooltip(params) {
		if (!params.value || !serverUtils.isOffline(params.value)) {
			return;
		}
		return params.value + ': ' + params.data.offlineReason;
	}

	/**
	 * Gets value to display a default tooltip.
	 */
	function defaultTooltip(params) {
		return params.value;
	}

	/**
	 * Formats the contents of a 'lastUpdated' column cell as "relative to now".
	 */
	function dateCellFormatter(params) {
		return params.value ? dateUtils.getRelativeTime(params.value) : params.value;
	}


	/**** Constants, scope data, etc. ****/

	/** The columns of the ag-grid table */
	const columns = [
		{
			headerName: "Cache Group",
			field: "cachegroup",
			hide: false
		},
		{
			headerName: "CDN",
			field: "cdnName",
			hide: false
		},
		{
			headerName: "Domain",
			field: "domainName",
			hide: false
		},
		{
			headerName: "Host",
			field: "hostName",
			hide: false
		},
		{
			headerName: "HTTPS Port",
			field: "httpsPort",
			hide: true,
			filter: "agNumberColumnFilter"
		},
		{
			headerName: "Hash ID",
			field: "xmppId",
			hide: true
		},
		{
			headerName: "ID",
			field: "id",
			hide: true,
			filter: "agNumberColumnFilter"
		},
		{
			headerName: "ILO IP Address",
			field: "iloIpAddress",
			hide: true,
			cellRenderer: "sshCellRenderer",
			onCellClicked: null
		},
		{
			headerName: "ILO IP Gateway",
			field: "iloIpGateway",
			hide: true,
			cellRenderer: "sshCellRenderer",
			onCellClicked: null
		},
		{
			headerName: "ILO IP Netmask",
			field: "iloIpNetmask",
			hide: true
		},
		{
			headerName: "ILO Username",
			field: "iloUsername",
			hide: true
		},
		{
			headerName: "Interface Name",
			field: "interfaceName",
			hide: true
		},
		{
			headerName: "IPv6 Address",
			field: "ip6Address",
			hide: false
		},
		{
			headerName: "IPv6 Gateway",
			field: "ip6Gateway",
			hide: true
		},
		{
			headerName: "Last Updated",
			field: "lastUpdated",
			hide: true,
			filter: "agDateColumnFilter",
			valueFormatter: dateCellFormatter
		},
		{
			headerName: "Mgmt IP Address",
			field: "mgmtIpAddress",
			hide: true
		},
		{
			headerName: "Mgmt IP Gateway",
			field: "mgmtIpGateway",
			hide: true,
			filter: true,
			cellRenderer: "sshCellRenderer",
			onCellClicked: null
		},
		{
			headerName: "Mgmt IP Netmask",
			field: "mgmtIpNetmask",
			hide: true,
			filter: true,
			cellRenderer: "sshCellRenderer",
			onCellClicked: null
		},
		{
			headerName: "Network Gateway",
			field: "ipGateway",
			hide: true,
			filter: true,
			cellRenderer: "sshCellRenderer",
			onCellClicked: null
		},
		{
			headerName: "Network IP",
			field: "ipAddress",
			hide: false,
			filter: true,
			cellRenderer: "sshCellRenderer",
			onCellClicked: null
		},
		{
			headerName: "Network MTU",
			field: "interfaceMtu",
			hide: true,
			filter: "agNumberColumnFilter"
		},
		{
			headerName: "Network Subnet",
			field: "ipNetmask",
			hide: true
		},
		{
			headerName: "Offline Reason",
			field: "offlineReason",
			hide: true
		},
		{
			headerName: "Phys Location",
			field: "physLocation",
			hide: true
		},
		{
			headerName: "Profile",
			field: "profile",
			hide: false
		},
		{
			headerName: "Rack",
			field: "rack",
			hide: true
		},
		{
			headerName: "Reval Pending",
			field: "revalPending",
			hide: true,
			filter: true,
			cellRenderer: "updateCellRenderer"
		},
		{
			headerName: "Router Hostname",
			field: "routerHostName",
			hide: true
		},
		{
			headerName: "Router Port Name",
			field: "routerPortName",
			hide: true
		},
		{
			headerName: "Status",
			field: "status",
			hide: false,
			tooltip: offlineReasonTooltip
		},
		{
			headerName: "TCP Port",
			field: "tcpPort",
			hide: true
		},
		{
			headerName: "Type",
			field: "type",
			hide: false
		},
		{
			headerName: "Update Pending",
			field: "updPending",
			hide: false,
			filter: true,
			cellRenderer: "updateCellRenderer"
		},
		{
			headerName: "Status Last Updated",
			field: "statusLastUpdated",
			hide: true,
			filter: "agDateColumnFilter",
			valueFormatter: dateCellFormatter
		}
	];

	/** All of the statuses (populated on init). */
	let statuses = [];

	/** All of the servers - lastUpdated fields converted to actual Dates, ip fields populated from interfaces */
	$scope.servers = servers.map(
		function(x) {
			x.lastUpdated = x.lastUpdated ? new Date(x.lastUpdated.replace("+00", "Z")) : x.lastUpdated;
			x.statusLastUpdated = x.statusLastUpdated ? new Date(x.statusLastUpdated): x.statusLastUpdated;
			Object.assign(x, serverUtils.toLegacyIPInfo(x.interfaces));
	});

	/** The base URL to use for constructing links to server charts. */
	$scope.chartsBase = propertiesModel.properties.servers.charts.baseUrl;

	/** The currently selected server - at the moment only used by the context menu */
	$scope.server = {
		hostName: "",
		domainName: "",
		id: -1
	};

	$scope.quickSearch = '';

	$scope.pageSize = 100;

	$scope.mouseDownSelectionText = "";

	/** Options, configuration, data and callbacks for the ag-grid table. */
	$scope.gridOptions = {
		components: {
			sshCellRenderer: SSHCellRenderer,
			updateCellRenderer: UpdateCellRenderer
		},
		columnDefs: columns,
		enableCellTextSelection:true,
		suppressMenuHide: true,
		multiSortKey: 'ctrl',
		alwaysShowVerticalScroll: true,
		defaultColDef: {
			filter: true,
			sortable: true,
			resizable: true,
			tooltip: defaultTooltip
		},
		rowData: servers,
		pagination: true,
		paginationPageSize: $scope.pageSize,
		rowBuffer: 0,
		tooltipShowDelay: 500,
		allowContextMenuWithControlKey: true,
		preventDefaultOnContextMenu: true,
		colResizeDefault: "shift",
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
			$scope.server = params.data;
			$scope.$apply();
		},
		onCellMouseDown: function() {
			$scope.mouseDownSelectionText = window.getSelection().toString();
		},
		onRowClicked: function(params) {
			const selection = window.getSelection().toString();
			if(selection === "" || selection === $scope.mouseDownSelectionText) {
				locationUtils.navigateToPath('/servers/' + params.data.id);
				// Event is outside the digest cycle, so we need to trigger one.
				$scope.$apply();
			}
			$scope.mouseDownSelectionText = "";
		},
		onColumnResized: function(params) {
			localStorage.setItem(tableName + "_table_columns", JSON.stringify($scope.gridOptions.columnApi.getColumnState()));
		},

	};

	/** These three functions are used by the context menu to determine what functionality to provide for a server. */
	$scope.isCache = serverUtils.isCache;
	$scope.isEdge = serverUtils.isEdge;
	$scope.isOrigin = serverUtils.isOrigin;

	/** Used by the context menu to determine whether or not to include links to server charts. */
	$scope.showCharts = propertiesModel.properties.servers.charts.show;

	/** This is used to position the context menu under the cursor. */
	$scope.menuStyle = {
		left: 0,
		top: 0,
	};

	/** Controls whether or not the context menu is visible. */
	$scope.showMenu = false;


	/**** Miscellaneous scope functions ****/

	/** Reloads all 'resolve'd data for the view. */
	$scope.refresh = function() {
		$state.reload();
	};

	/** Toggles the visibility of a column that has the ID provided as 'col'. */
	$scope.toggleVisibility = function(col) {
		const visible = $scope.gridOptions.columnApi.getColumn(col).isVisible();
		$scope.gridOptions.columnApi.setColumnVisible(col, !visible);
	};

	/** Downloads the table as a CSV */
	$scope.exportCSV = function() {
		// TODO: figure out how to reconcile clicking on a server taking you to it
		// with row selection exports.
		const params = {
			allColumns: true,
			fileName: "servers.csv",
		};
		$scope.gridOptions.api.exportDataAsCsv(params);
	};

	/**** Context menu functions ****/

	$scope.queueServerUpdates = function(server, event) {
		event.stopPropagation();
		serverService.queueServerUpdates(server.id).then($scope.refresh);
	};

	$scope.clearServerUpdates = function(server, event) {
		event.stopPropagation();
		serverService.clearServerUpdates(server.id).then($scope.refresh);
	};

	$scope.confirmDelete = function(server, event) {
		event.stopPropagation();

		const params = {
			title: 'Delete Server: ' + server.hostName,
			key: server.hostName
		};
		const modalInstance = $uibModal.open({
			templateUrl: 'common/modules/dialog/delete/dialog.delete.tpl.html',
			controller: 'DialogDeleteController',
			size: 'md',
			resolve: {
				params: function () {
					return params;
				}
			}
		});
		modalInstance.result.then(
			function() {
				serverService.deleteServer(server.id).then(
					function(result) {
						messageModel.setMessages(result.alerts, false);
						$scope.refresh();
					},
					function(err) {
						// TODO: use template strings once the build can handle them.
						console.error("Error deleting server", server.hostName + "." + server.domainName, "(#" + String(server.id) + "):", err);
					}
				);
			},
			function() {
				// This is just a cancel event from closing the dialog, do nothing.
			}
		);
	};

	/**
	 * updateStatus sets the status of the given server to the given status value.
	 *
	 * @param {{id: number, offlineReason?: string}} status The numeric ID of the status to set along with a reason why it was set offline, if applicable.
	 * @param {{id: number}} server The server (or at least its numeric ID) which will have its status set.
	 */
	function updateStatus(status, server) {
		const params = {
			status: status.id,
			offlineReason: status.offlineReason
		};

		serverService.updateStatus(server.id, params).then(
			function(result) {
				messageModel.setMessages(result.data.alerts, false);
				$scope.refresh();
			},
			function(fault) {
				messageModel.setMessages(fault.data.alerts, false);
			}
		);
	};

	$scope.confirmStatusUpdate = function(server, event) {
		event.stopPropagation();

		const modalInstance = $uibModal.open({
			templateUrl: 'common/modules/dialog/select/status/dialog.select.status.tpl.html',
			controller: 'DialogSelectStatusController',
			size: 'md',
			resolve: {
				server: function() {
					return server;
				},
				statuses: function() {
					return statuses;
				}
			}
		});
		modalInstance.result.then(
			function(status) {
				updateStatus(status, server);
			},
			function () {
				// this is just a cancel event from closing the dialog, do nothing
			}
		);
	};

	$scope.confirmCDNQueueServerUpdates = function() {
		const params = {
			title: 'Queue Server Updates',
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
			cdnService.queueServerUpdates(cdn.id).then($scope.refresh);
		}, function () {
			// do nothing
		});
	};

	$scope.confirmCDNClearServerUpdates = function() {
		const params = {
			title: 'Clear Server Updates',
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
			cdnService.clearServerUpdates(cdn.id).then($scope.refresh);
		}, function () {
			// do nothing
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
			// apply any filter provided to the controller
			Object.assign(filterState, filter);
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

		statusService.getStatuses().then(
			function(result) {
				statuses = result;
			}
		);
	});
};

TableServersController.$inject = ['tableName', 'servers', 'filter', '$scope', '$state', '$uibModal', '$window', 'dateUtils', 'locationUtils', 'serverUtils', 'cdnService', 'serverService', 'statusService', 'propertiesModel', 'messageModel', "userModel", "$document"];
module.exports = TableServersController;
