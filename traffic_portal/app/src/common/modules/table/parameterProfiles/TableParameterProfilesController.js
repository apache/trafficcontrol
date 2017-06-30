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

var TableParameterProfilesController = function(parameter, parameterProfiles, $scope, $state, $uibModal, locationUtils, profileParameterService) {

	$scope.parameter = parameter;

	$scope.parameterProfiles = parameterProfiles;

	$scope.addProfile = function() {
		alert('not hooked up yet: add profile to parameter');
	};

	$scope.removeProfile = function(profileId) {
		profileParameterService.unlinkProfileParameter(profileId, parameter.id)
			.then(
				function() {
					$scope.refresh();
				}
			);
	};

	$scope.refresh = function() {
		$state.reload(); // reloads all the resolves for the view
	};

	$scope.selectProfiles = function() {
		var modalInstance = $uibModal.open({
			templateUrl: 'common/modules/table/parameterProfiles/table.paramProfilesUnassigned.tpl.html',
			controller: 'TableParamProfilesUnassignedController',
			size: 'lg',
			resolve: {
				parameter: function() {
					return parameter;
				},
				profiles: function(profileService) {
					return profileService.getParamUnassignedProfiles(parameter.id);
				}
			}
		});
		modalInstance.result.then(function(selectedProfiles) {
			var massagedArray = [];
			for (i = 0; i < selectedProfiles.length; i++) {
				massagedArray.push( { parameterId: parameter.id, profileId: selectedProfiles[i] } );
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
		$('#parameterProfilesTable').dataTable({
			"aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
			"iDisplayLength": 25,
			"aaSorting": []
		});
	});

};

TableParameterProfilesController.$inject = ['parameter', 'parameterProfiles', '$scope', '$state', '$uibModal', 'locationUtils', 'profileParameterService'];
module.exports = TableParameterProfilesController;
