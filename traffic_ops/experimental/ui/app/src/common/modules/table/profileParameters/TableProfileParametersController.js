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

var TableProfileParametersController = function(profile, profileParameters, $scope, $state, $uibModal, locationUtils, profileParameterService) {

	$scope.profile = profile;

	$scope.profileParameters = profileParameters;

	$scope.removeParameter = function(paramId) {
		profileParameterService.unlinkProfileParameter(profile.id, paramId)
			.then(
				function() {
					$scope.refresh();
				}
			);
	};

	$scope.refresh = function() {
		$state.reload(); // reloads all the resolves for the view
	};

	$scope.selectParams = function() {
		var modalInstance = $uibModal.open({
			templateUrl: 'common/modules/table/profileParameters/table.profileParamsUnassigned.tpl.html',
			controller: 'TableProfileParamsUnassignedController',
			size: 'lg',
			resolve: {
				profile: function(parameterService) {
					return profile;
				},
				parameters: function(parameterService) {
					return parameterService.getProfileUnassignedParams(profile.id);
				}
			}
		});
		modalInstance.result.then(function(selectedParams) {
			var massagedArray = [];
			for (i = 0; i < selectedParams.length; i++) {
				massagedArray.push( { profileId: profile.id, parameterId: selectedParams[i] } );
			}
			profileParameterService.linkProfileParameters(massagedArray)
				.then(
					function() {
						$scope.refresh();
					}
				);
		}, function () {
			// do nothing
		});
	};

	$scope.navigateToPath = locationUtils.navigateToPath;

	angular.element(document).ready(function () {
		$('#profileParametersTable').dataTable({
			"aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
			"iDisplayLength": 25,
			"aaSorting": []
		});
	});

};

TableProfileParametersController.$inject = ['profile', 'profileParameters', '$scope', '$state', '$uibModal', 'locationUtils', 'profileParameterService'];
module.exports = TableProfileParametersController;
