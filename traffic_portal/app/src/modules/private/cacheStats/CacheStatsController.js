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

/**
 * @typedef CacheStat
 * @property {string} cachegroup
 * @property {number} connections
 * @property {boolean} healthy
 * @property {string} hostname
 * @property {string | null} ip
 * @property {number} kbps
 * @property {string} profile
 * @property {string} status
 */


class CacheStatsHealthyCellRenderer {
	/** @type HTMLSpanElement */
	eGui;

	/**
	 * Called by AG-Grid as a pseudo constructor, used when a renderer is
	 * instantiated.
	 *
	 * @param {{
	 * 	api: any;
	 * 	colDef: any;
	 * 	column: any;
	 * 	columnApi: any;
	 * 	context: any;
	 * 	data: CacheStat;
	 * 	eGridCell: HTMLElement;
	 * 	eParentOfValue: HTMLElement;
	 * 	formatValue: function;
	 * 	fullWidth: boolean;
	 * 	getValue: function;
	 * 	node: any;
	 * 	pinned: string | null;
	 * 	refreshCell: function;
	 * 	registerRowDragger: function;
	 * 	rowIndex: number;
	 * 	setValue: function;
	 * 	value: boolean;
	 * 	valueFormatted: null;
	 * 	$scope: null;
	 * }} params
	 */
	init(params) {
		this.eGui = document.createElement("span");
		this.eGui.setAttribute("class", params.value ? "green" : "red");
		this.eGui.textContent = String(params.value);
	}

	/**
	 * Gets a rendered cell. Parameters are available, but not currently used.
	 *
	 * @returns {HTMLElement | Text} A rendered cell element.
	 */
	getGui() {
		return this.eGui;
	}
}

/**
 *
 * @param {CacheStat[]} cacheStats
 * @param {{user: {username: string}}} userModel
 */
var CacheStatsController = function(cacheStats, $scope, userModel) {

	class CacheStatsSSHCellRenderer {
		/** @type HTMLAnchorElement | Text */
		eGui;

		/**
		 * Called by AG-Grid as a pseudo constructor, used when a renderer is
		 * instantiated.
		 *
		 * @param {{
		 * 	api: any;
		 * 	colDef: any;
		 * 	column: any;
		 * 	columnApi: any;
		 * 	context: any;
		 * 	data: CacheStat;
		 * 	eGridCell: HTMLElement;
		 * 	eParentOfValue: HTMLElement;
		 * 	formatValue: function;
		 * 	fullWidth: boolean;
		 * 	getValue: function;
		 * 	node: any;
		 * 	pinned: string | null;
		 * 	refreshCell: function;
		 * 	registerRowDragger: function;
		 * 	rowIndex: number;
		 * 	setValue: function;
		 * 	value: string;
		 * 	valueFormatted: null;
		 * 	$scope: null;
		 * }} params
		 */
		init(params) {
			if (params.data.ip === null) {
				this.eGui = document.createTextNode(params.value);
				return;
			}
			this.eGui = document.createElement("a");
			this.eGui.href = `ssh://${userModel.user.username}@${params.data.ip}`;
			this.eGui.setAttribute("class", "link");
			this.eGui.textContent = params.value;
		}

		/**
		 * Gets a rendered cell. Parameters are available, but not currently used.
		 *
		 * @returns {HTMLElement | Text} A rendered cell element.
		 */
		getGui() {
			return this.eGui;
		}
	}

	/**
	 * The columns of the ag-grid table.
	 * @type CGC.ColumnDefinition[]
	 */
	$scope.columns = [
		{
			headerName: "Cache Group",
			field: "cachegroup",
			hide: false
		},
		{
			headerName: "Connections",
			field: "connections",
			hide: false,
			filter: "agNumberColumnFilter"
		},
		{
			cellRenderer: CacheStatsHealthyCellRenderer,
			headerName: "Healthy",
			field: "healthy",
			hide: false
		},
		{
			cellRenderer: CacheStatsSSHCellRenderer,
			headerName: "Host",
			field: "hostname",
			hide: false
		},
		{
			cellRenderer: CacheStatsSSHCellRenderer,
			headerName: "IP",
			field: "ip",
			hide: true,
		},
		{
			headerName: "Data out rate",
			field: "kbps",
			filter: "agNumberColumnFilter",
			hide: false,
			/**
			 * Formats the kbps metric.
			 * @param {{value: number}} value Some AG Grid params containing the raw value to be formatted.
			 * @returns {string} The formatted value.
			 */
			valueFormatter: ({value}) => {
				if (!value || value <= 0 || !isFinite(value)) {
					return "0bps";
				}
				if (value >= 1e9) {
					return `${(value/1e9).toFixed(3)}Tb/s`;
				}
				if (value >= 1e6) {
					return `${(value/1e6).toFixed(3)}Gb/s`;
				}
				if (value >= 1000) {
					return `${(value/1000).toFixed(3)}Mb/s`;
				}
				if (value >= 1) {
					return `${value.toFixed(3)}kb/s`;
				}
				return `${(value*1000).toFixed(0)}b/s`;
			}
		},
		{
			headerName: "Profile",
			field: "profile",
			hide: false,
		},
		{
			headerName: "Status",
			field: "status",
			hide: false,
		}
	];

	$scope.data = cacheStats;

	/** Options, configuration, data and callbacks for the ag-grid table. */
	/** @type CGC.GridSettings */
	$scope.gridOptions = {
		refreshable: true,
	};

	$scope.defaultData = {
		cachegroup: "ALL",
		connections: -1,
		healthy: false,
		hostname: "ALL",
		ip: null,
		kbps: -1,
		profile: "ALL",
		status: "ALL"
	};

};

CacheStatsController.$inject = ["cacheStats", "$scope", "userModel"];
module.exports = CacheStatsController;
