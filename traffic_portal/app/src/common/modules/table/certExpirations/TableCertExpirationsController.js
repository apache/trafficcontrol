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

var TableCertExpirationsController = function(tableName, certExpirations, dsXmlToIdMap, filter, $document, $scope, $state, $filter, locationUtils) {

	/** All of the expiration fields converted to actual Dates */
	$scope.certExpirations = certExpirations.map(
		function(x) {
			// need to convert this to a date object for ag-grid filter to work properly
			x.expiration = new Date(x.expiration);
			return x;
		});

	$scope.dsXmlToIdMap = dsXmlToIdMap;

	$scope.editCertExpirations = function(dsId) {
		locationUtils.navigateToPath('/delivery-services/' + dsId + '/ssl-keys');
	}

	/**
	 * Formats the contents of a 'federation' column cell as "True" or blank for visibility.
	 */
	function federatedCellFormatter(params) {
		if (!params.value) {
			return '';
		} else {
			return params.value;
		}
	}

	/** The columns of the ag-grid table */
	$scope.columns = [
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
			filter: "agDateColumnFilter"
		},
		{
			headerName: "Federated",
			field: "federated",
			hide: false,
			valueFormatter: federatedCellFormatter
		},
	];

	/** Options, configuration, data and callbacks for the ag-grid table. */
	$scope.gridOptions = {
		onRowClick: function(params) {
			const selection = window.getSelection().toString();
			if(selection === "" || selection === $scope.mouseDownSelectionText) {
				locationUtils.navigateToPath('/delivery-services/' + $scope.dsXmlToIdMap[params.data.deliveryservice] + '/ssl-keys');
				// Event is outside the digest cycle, so we need to trigger one.
				$scope.$apply();
			}
			$scope.mouseDownSelectionText = "";
		}
	};

};

TableCertExpirationsController.$inject = ['tableName', 'certExpirations', 'dsXmlToIdMap', 'filter', '$document', '$scope', '$state', '$filter', 'locationUtils'];
module.exports = TableCertExpirationsController;
