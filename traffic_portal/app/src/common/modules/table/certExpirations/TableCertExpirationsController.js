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
 * @param {*} certExpirations
 * @param {*} $scope
 * @param {import("../../../service/utils/LocationUtils")} locationUtils
 * @param {*} deliveryservices
 */
var TableCertExpirationsController = function(certExpirations, $scope, locationUtils, deliveryservices) {

	/** All of the expiration fields converted to actual Dates */
	$scope.certExpirations = certExpirations.map(
		function(x) {
			// need to convert this to a date object for ag-grid filter to work properly
			x.expiration = new Date(x.expiration);
			return x;
		});

	$scope.dsXmlToIdMap = new Map(deliveryservices.map(ds=>[ds.xmlId, ds.id]));

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
			if(!selection) {
				locationUtils.navigateToPath('/delivery-services/' + $scope.dsXmlToIdMap.get(params.data.deliveryservice) + '/ssl-keys');
				// Event is outside the digest cycle, so we need to trigger one.
				$scope.$apply();
			}
		},
		rowClassRules: {
			'expired-cert': function(params) {
				const now = new Date();
				return params.data.expiration < now;
			},
			'seven-days-until-expired': function(params) {
				const sevenDays = new Date();
				sevenDays.setDate(sevenDays.getDate()+7);
				return params.data.expiration >= new Date() && params.data.expiration <= sevenDays;
			},
			'thirty-days-until-expired': function(params) {
				const thirtyDays = new Date();
				thirtyDays.setDate(thirtyDays.getDate()+30);
				return params.data.expiration >= new Date() && params.data.expiration <= thirtyDays;
			}
		}
	};

};

TableCertExpirationsController.$inject = ['certExpirations', '$scope', 'locationUtils', 'deliveryservices'];
module.exports = TableCertExpirationsController;
