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
 * @param {*} cdniRequests
 * @param {*} $scope
 * @param {import("../../../service/utils/LocationUtils")} locationUtils
 */
var TableCdniController = function(cdniRequests, $scope, locationUtils) {

	$scope.cdniRequests = cdniRequests.map(
		function(x) {
			// need to convert this to a date object for ag-grid filter to work properly
			x.data = JSON.stringify(x.data);
			return x;
		});

	/** The columns of the ag-grid table */
	$scope.columns = [
		{
			headerName: "Upstream CDN",
			field: "ucdn",
			hide: false
		},
		{
			headerName: "Host",
			field: "host",
			hide: false
		},
		{
			headerName: "Request Type",
			field: "request_type",
			hide: false
		},
		{
			headerName: "New Data",
			field: "data",
			hide: false
		}
	];

	/** Options, configuration, data and callbacks for the ag-grid table. */
	$scope.gridOptions = {
		onRowClick: function(params) {
			const selection = window.getSelection().toString();
			if(!selection) {
				locationUtils.navigateToPath('/cdni-config-requests/' + params.data.id);
				// Event is outside the digest cycle, so we need to trigger one.
				$scope.$apply();
			}
		}
	};

};

TableCdniController.$inject = ['cdniRequests', '$scope', 'locationUtils'];
module.exports = TableCdniController;
