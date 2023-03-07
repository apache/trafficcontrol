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
 * @param {*} cacheChecks
 * @param {*} $scope
 * @param {import("../../../common/service/utils/LocationUtils")} locationUtils
 * @param {import("../../../common/models/PropertiesModel")} propertiesModel
 */
var CacheChecksController = function(cacheChecks, $scope, locationUtils, propertiesModel) {
	$scope.cacheChecks = cacheChecks;

	$scope.config = propertiesModel.properties.cacheChecks;

	/** @type {import("../../../common/modules/table/agGrid/CommonGridController").CGC.ColumnDefinition} */
	$scope.columns = [
		{
			headerName: "Hostname",
			field: "hostName",
		},
		{
			headerName: "Profile",
			field: "profile",
		},
		{
			headerName: "Status",
			field: "adminState"
		},
		{
			headerName: $scope.config.updatePending.key,
			field: "updPending",
			filter: true,
			cellRenderer: "updateCellRenderer",
			headerTooltip: $scope.config.updatePending.desc
		},
		{
			headerName: $scope.config.revalPending.key,
			field: "revalPending",
			filter: true,
			cellRenderer: "updateCellRenderer",
			headerTooltip: $scope.config.revalPending.desc
		},
	];

	/** @type {import("../../../common/modules/table/agGrid/CommonGridController").CGC.GridSettings} */
	$scope.gridOptions = {
		refreshable: true,
		onRowClick(row) {
			locationUtils.navigateToPath('/servers/' + row.data.id);
		}
	};

	this.unroll = function(row) {
	    $scope.config.extensions.forEach(ext => {
	        let key = "_config" + ext.key;
	        if(row.checks !== undefined) {
				row[key] = row.checks[ext.key];
			}
		});
	};

	this.$onInit = function() {
		let self = this;
		$scope.cacheChecks.forEach(cacheCheck => {
			self.unroll(cacheCheck);
		});
		$scope.config.extensions.forEach(ext => {
			let colDef = {
				headerName: ext.key,
                field: "_config" + ext.key,
				headerTooltip: ext.desc
			};
		    if(ext.type === "bool") {
		    	colDef.cellRenderer = "checkCellRenderer";
		    	colDef.filter = true;
			}
		    $scope.columns.push(colDef);
		});
	};
};

CacheChecksController.$inject = ['cacheChecks', '$scope', 'locationUtils', 'propertiesModel'];
module.exports = CacheChecksController;
