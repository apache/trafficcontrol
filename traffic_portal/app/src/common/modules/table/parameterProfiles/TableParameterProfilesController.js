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

var TableParameterProfilesController = function(parameter, profiles, $controller, $scope, $state, $uibModal, locationUtils, deliveryServiceService, profileParameterService, serverService) {

	// extends the TableProfilesController to inherit common methods
	angular.extend(this, $controller('TableProfilesController', { profiles: profiles, $scope: $scope }));

	let parameterProfilesTable;

	var removeProfile = function(profileId) {
		profileParameterService.unlinkProfileParameter(profileId, parameter.id)
			.then(
				function() {
					$scope.refresh();
				}
			);
	};

	$scope.parameter = parameter;

	// adds some items to the base profiles context menu
	$scope.contextMenuItems.splice(2, 0,
		{
			text: 'Unlink Profile from Parameter',
			hasBottomDivider: function() {
				return true;
			},
			click: function ($itemScope) {
				$scope.confirmRemoveProfile($itemScope.p);
			}
		}
	);

	$scope.confirmRemoveProfile = function(profile, $event) {
		if ($event) {
			$event.stopPropagation();
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
						removeProfile(profile.id);
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
						removeProfile(profile.id);
					}, function () {
						// do nothing
					});
				});
		}
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
				allProfiles: function(profileService) {
					return profileService.getProfiles({ orderby: 'name' });
				},
				assignedProfiles: function() {
					return profiles;
				}
			}
		});
		modalInstance.result.then(function(selectedProfileIds) {
			var params = {
				title: 'Assign profiles to ' + parameter.name,
				message: 'Are you sure you want to modify the profiles assigned to ' + parameter.name + '?'
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
				profileParameterService.linkParamProfiles(parameter.id, selectedProfileIds)
					.then(
						function() {
							$scope.refresh(); // refresh the parameter profiles table
						}
					);
			}, function () {
				// do nothing
			});
		}, function () {
			// do nothing
		});
	};

	$scope.toggleVisibility = function(colName) {
		const col = parameterProfilesTable.column(colName + ':name');
		col.visible(!col.visible());
		parameterProfilesTable.rows().invalidate().draw();
	};

	$scope.columnFilterFn = function(column) {
		if (column.name === 'Action') {
			return false;
		}
		return true;
	};

	angular.element(document).ready(function () {
		parameterProfilesTable = $('#parameterProfilesTable').DataTable({
			"aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
			"iDisplayLength": 25,
			"columnDefs": [
				{ 'orderable': false, 'targets': 5 }
			],
			"aaSorting": [],
			"columns": [
				{ "name": "Name", "visible": true, "searchable": true },
				{ "name": "Type", "visible": true, "searchable": true },
				{ "name": "Routing Disabled", "visible": true, "searchable": true },
				{ "name": "Description", "visible": true, "searchable": true },
				{ "name": "CDN", "visible": true, "searchable": true },
				{ "name": "Action", "visible": true, "searchable": false }
			],
			"initComplete": function(settings, json) {
				try {
					// need to create the show/hide column checkboxes and bind to the current visibility
					$scope.columns = JSON.parse(localStorage.getItem('DataTables_parameterProfilesTable_/')).columns;
				} catch (e) {
					console.error("Failure to retrieve required column info from localStorage (key=DataTables_parameterProfilesTable_/):", e);
				}
			}
		});
	});

};

TableParameterProfilesController.$inject = ['parameter', 'profiles', '$controller', '$scope', '$state', '$uibModal', 'locationUtils', 'deliveryServiceService', 'profileParameterService', 'serverService'];
module.exports = TableParameterProfilesController;
