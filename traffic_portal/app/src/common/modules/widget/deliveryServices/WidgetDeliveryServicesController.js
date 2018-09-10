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

var WidgetDeliveryServicesController = function($scope, $interval, deliveryServiceService, locationUtils, propertiesModel) {

	var interval,
		autoRefresh = propertiesModel.properties.dashboard.autoRefresh;

	var getDeliveryServices = function() {
		deliveryServiceService.getDeliveryServices()
			.then(function(result) {
				$scope.deliveryServices = result;
			});
	};

	// pagination
	$scope.currentDeliveryServicesPage = 1;
	$scope.deliveryServicesPerPage = 10;

	$scope.navigateToPath = locationUtils.navigateToPath;


	$scope.$on("$destroy", function() {

	});

	var init = function() {
		getDeliveryServices();
	};
	init();

};

WidgetDeliveryServicesController.$inject = ['$scope', '$interval', 'deliveryServiceService', 'locationUtils', 'propertiesModel'];
module.exports = WidgetDeliveryServicesController;
