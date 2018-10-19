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

var FormEditEndpointController = function(endpoint, $scope, $controller, $uibModal, $anchorScroll, locationUtils, endpointService) {

	// extends the FormEndpointController to inherit common methods
	angular.extend(this, $controller('FormEndpointController', { endpoint: endpoint, $scope: $scope }));

	var deleteEndpoint = function(endpoint) {
		endpointService.deleteEndpoint(endpoint.id)
			.then(function() {
				locationUtils.navigateToPath('/endpoints');
			});
	};

	$scope.endpointName = angular.copy(endpoint.httpMethod) + ' /api/*/' + angular.copy(endpoint.httpRoute);

	$scope.settings = {
		isNew: false,
		saveLabel: 'Update'
	};

	$scope.save = function(endpoint) {
		endpointService.updateEndpoint(endpoint).
		then(function() {
			$scope.endpointName = angular.copy(endpoint.httpMethod) + ' ' + angular.copy(endpoint.httpRoute);
			$anchorScroll(); // scrolls window to top
		});
	};

	$scope.confirmDelete = function(endpoint) {
		var params = {
			title: 'Delete Endpoint: ' + endpoint.httpMethod + ' ' + endpoint.httpRoute,
			key: endpoint.httpMethod + ' ' + endpoint.httpRoute
		};
		var modalInstance = $uibModal.open({
			templateUrl: 'common/modules/dialog/delete/dialog.delete.tpl.html',
			controller: 'DialogDeleteController',
			size: 'md',
			resolve: {
				params: function () {
					return params;
				}
			}
		});
		modalInstance.result.then(function() {
			deleteEndpoint(endpoint);
		}, function () {
			// do nothing
		});
	};

};

FormEditEndpointController.$inject = ['endpoint', '$scope', '$controller', '$uibModal', '$anchorScroll', 'locationUtils', 'endpointService'];
module.exports = FormEditEndpointController;
