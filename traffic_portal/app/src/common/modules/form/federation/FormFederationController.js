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
 * @param {*} cdn
 * @param {*} federation
 * @param {*} deliveryServices
 * @param {*} $scope
 * @param {import("angular").ILocationService} $location
 * @param {import("../../../service/utils/FormUtils")} formUtils
 * @param {import("../../../service/utils/LocationUtils")} locationUtils
 */
var FormFederationController = function(cdn, federation, deliveryServices, $scope, $location, formUtils, locationUtils) {

	$scope.cdn = cdn;

	$scope.federation = federation;

	$scope.deliveryServices = deliveryServices;

	$scope.viewDeliveryServices = function() {
		$location.path($location.path() + '/delivery-services');
	};

	$scope.viewUsers = function() {
		$location.path($location.path() + '/users');
	};

	$scope.navigateToPath = (path, unsavedChanges) => locationUtils.navigateToPath(path, unsavedChanges);

	$scope.hasError = formUtils.hasError;

	$scope.hasPropertyError = formUtils.hasPropertyError;

};

FormFederationController.$inject = ['cdn', 'federation', 'deliveryServices', '$scope', '$location', 'formUtils', 'locationUtils'];
module.exports = FormFederationController;
