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

var FormCloneDeliveryServiceController = function(deliveryService, origin, topologies, type, types, $scope, $controller) {

	// extends the FormNewDeliveryServiceController to inherit common methods
	angular.extend(this, $controller('FormNewDeliveryServiceController', { deliveryService: deliveryService, origin: origin, type: type, topologies: topologies, types: types, $scope: $scope }));

	$scope.deliveryServiceName = deliveryService.xmlId + ' clone';

	$scope.advancedShowing = true;

    $scope.restrictTLS = deliveryService.tlsVersions instanceof Array && deliveryService.tlsVersions.length > 0;

	$scope.settings = {
		isNew: true,
		saveLabel: 'Clone'
	};

	var init = function() {
		// we're going to let them select an xmlId and a type for the clone
		$scope.deliveryService.xmlId = null;
		$scope.deliveryService.typeId = null;
		$scope.loadGeoLimitCountriesRaw(deliveryService);
	};
	init();

};

FormCloneDeliveryServiceController.$inject = ['deliveryService', 'origin', 'topologies', 'type', 'types', '$scope', '$controller'];
module.exports = FormCloneDeliveryServiceController;
