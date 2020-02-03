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

var TableProfileParametersController = function(profile, parameters, $controller, $scope, $state, $uibModal, locationUtils, deliveryServiceService, profileParameterService, serverService, messageModel) {

	// extends the TableParametersController to inherit common methods
	angular.extend(this, $controller('TableParametersController', { parameters: parameters, $scope: $scope }));

	$scope.profile = profile;

	// adds some items to the base parameters context menu
	$scope.contextMenuItems.splice(2, 0,
		{
			text: 'Unlink Parameter from Profile',
			click: function ($itemScope) {
				$scope.confirmRemoveParam($itemScope.p);
			}
		},
		null // Divider
	);

	var removeParameter = function(paramId) {
		profileParameterService.unlinkProfileParameter(profile.id, paramId)
			.then(
				function() {
					$scope.refresh(); // refresh the profile parameters table
				}
			);
	};

	var linkProfileParameters = function(paramIds) {
		profileParameterService.linkProfileParameters(profile, paramIds)
			.then(
				function(result) {
					messageModel.setMessages(result.data.alerts, false);
					$scope.refresh(); // refresh the profile parameters table
				}
			);
	};

	$scope.confirmRemoveParam = function(parameter, $event) {
		if ($event) {
			$event.stopPropagation(); // this kills the click event so it doesn't trigger anything else
		}
		if (profile.type == 'DS_PROFILE') { // if this is a ds profile, then it is used by delivery service(s) so we'll fetch the ds count...
			deliveryServiceService.getDeliveryServices({ profile: profile.id }).
				then(function(result) {
					var params = {
						title: 'Remove Parameter from Profile?',
						message: 'The ' + profile.name + ' profile is used by ' + result.length + ' delivery service(s). Are you sure you want to remove the ' + parameter.name + ' parameter from this profile?'
					};
					var modalInstance = $uibModal.open({
						templateUrl: 'common/modules/dialog/confirm/dialog.confirm.tpl.html',
						controller: 'DialogConfirmController',
						size: 'md',
						resolve: {
							params: function () {
								return params;
							}
						}
					});
					modalInstance.result.then(function() {
						removeParameter(parameter.id);
					}, function () {
						// do nothing
					});
				});
		} else { // otherwise the profile is used by servers so we'll fetch the server count...
			serverService.getServers({ profileId: profile.id }).
				then(function(result) {
					var params = {
						title: 'Remove Parameter from Profile?',
						message: 'The ' + profile.name + ' profile is used by ' + result.length + ' server(s). Are you sure you want to remove the ' + parameter.name + ' parameter from this profile?'
					};
					var modalInstance = $uibModal.open({
						templateUrl: 'common/modules/dialog/confirm/dialog.confirm.tpl.html',
						controller: 'DialogConfirmController',
						size: 'md',
						resolve: {
							params: function () {
								return params;
							}
						}
					});
					modalInstance.result.then(function() {
						removeParameter(parameter.id);
					}, function () {
						// do nothing
					});
				});
		}
	};

	$scope.selectParams = function() {
		var modalInstance = $uibModal.open({
			templateUrl: 'common/modules/table/profileParameters/table.profileParamsUnassigned.tpl.html',
			controller: 'TableProfileParamsUnassignedController',
			size: 'lg',
			resolve: {
				profile: function() {
					return profile;
				},
				allParams: function(parameterService) {
					return parameterService.getParameters();
				},
				assignedParams: function() {
					return parameters;
				}
			}
		});
		modalInstance.result.then(function(selectedParamIds) {
			if (profile.type == 'DS_PROFILE') { // if this is a ds profile, then it is used by delivery service(s) so we'll fetch the ds count...
				deliveryServiceService.getDeliveryServices({ profile: profile.id }).
					then(function(result) {
						var params = {
							title: 'Modify ' + profile.name + ' parameters',
							message: 'The ' + profile.name + ' profile is used by ' + result.length + ' delivery service(s). Are you sure you want to modify the parameters?'
						};
						var modalInstance = $uibModal.open({
							templateUrl: 'common/modules/dialog/confirm/dialog.confirm.tpl.html',
							controller: 'DialogConfirmController',
							size: 'md',
							resolve: {
								params: function () {
									return params;
								}
							}
						});
						modalInstance.result.then(function() {
							linkProfileParameters(selectedParamIds);
						}, function () {
							// do nothing
						});
					});
			} else { // otherwise the profile is used by servers so we'll fetch the server count...
				serverService.getServers({ profileId: profile.id }).
					then(function(result) {
						var params = {
							title: 'Modify ' + profile.name + ' parameters',
							message: 'The ' + profile.name + ' profile is used by ' + result.length + ' server(s). Are you sure you want to modify the parameters?'
						};
						var modalInstance = $uibModal.open({
							templateUrl: 'common/modules/dialog/confirm/dialog.confirm.tpl.html',
							controller: 'DialogConfirmController',
							size: 'md',
							resolve: {
								params: function () {
									return params;
								}
							}
						});
						modalInstance.result.then(function() {
							linkProfileParameters(selectedParamIds);
						}, function () {
							// do nothing
						});
					});
			}
		}, function () {
			// do nothing
		});
	};

	$scope.navigateToPath = locationUtils.navigateToPath;

	angular.element(document).ready(function () {
		$('#profileParametersTable').dataTable({
			"aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
			"iDisplayLength": 25,
			"columnDefs": [
				{ "width": "50%", "targets": 2 },
				{ 'orderable': false, 'targets': 4 }
			],
			"aaSorting": []
		});
	});

};

TableProfileParametersController.$inject = ['profile', 'parameters', '$controller', '$scope', '$state', '$uibModal', 'locationUtils', 'deliveryServiceService', 'profileParameterService', 'serverService', 'messageModel'];
module.exports = TableProfileParametersController;
