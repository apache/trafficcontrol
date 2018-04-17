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

var TableProfileDeliveryServicesController = function(profile, deliveryServices, $scope, $state, dateUtils, deliveryServiceUtils, locationUtils, propertiesModel) {

	var protocols = deliveryServiceUtils.protocols;

	var qstrings = deliveryServiceUtils.qstrings;

	$scope.profile = profile;

	$scope.deliveryServices = deliveryServices;

	$scope.showChartsButton = propertiesModel.properties.deliveryServices.charts.show;

	$scope.openCharts = deliveryServiceUtils.openCharts;

	$scope.protocol = function(ds) {
		return protocols[ds.protocol];
	};

	$scope.qstring = function(ds) {
		return qstrings[ds.qstringIgnore];
	};

	$scope.getRelativeTime = dateUtils.getRelativeTime;

	$scope.editDeliveryService = function(ds) {
		locationUtils.navigateToPath('/delivery-services/' + ds.id + '?type=' + ds.type);
	};

	$scope.refresh = function() {
		$state.reload(); // reloads all the resolves for the view
	};

	$scope.navigateToPath = locationUtils.navigateToPath;

	angular.element(document).ready(function () {
		$('#deliveryServicesTable').dataTable({
			"aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
			"iDisplayLength": 25,
			"columnDefs": [
				{ 'orderable': false, 'targets': 12 }
			],
			"aaSorting": []
		});
	});

};

TableProfileDeliveryServicesController.$inject = ['profile', 'deliveryServices', '$scope', '$state', 'dateUtils', 'deliveryServiceUtils', 'locationUtils', 'propertiesModel'];
module.exports = TableProfileDeliveryServicesController;
