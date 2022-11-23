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

/** @typedef {import("jquery")}*/

/**
 * @param {{id: number; name: string}} parameter
 * @param {unknown[]} profiles
 * @param {*} $scope
 * @param {*} $state
 * @param {import("../../../service/utils/angular.ui.bootstrap").IModalService} $uibModal
 * @param {import("angular").IWindowService} $window
 * @param {import("angular").ILocationService} $location
 * @param {import("../../../service/utils/LocationUtils")} locationUtils
 * @param {import("../../../api/DeliveryServiceService")} deliveryServiceService
 * @param {import("../../../api/ProfileParameterService")} profileParameterService
 * @param {import("../../../api/ServerService")} serverService
 * @param {import("../../../api/ProfileService")} profileService
 * @param {import("../../../models/MessageModel")} messageModel
 * @param {import("../../../service/utils/FileUtils")} fileUtils
 */
const TableParameterProfilesController = function (parameter, profiles, $scope, $state, $uibModal, $window, $location, locationUtils, deliveryServiceService, profileParameterService, serverService, profileService, messageModel, fileUtils) {
	const deleteProfile = function (profile) {
		profileService.deleteProfile(profile.id)
			.then(function (result) {
				messageModel.setMessages(result.alerts, false);
				$scope.refresh();
			});
	};
	let parameterProfilesTable;

	const removeProfile = function (profileId) {
		profileParameterService.unlinkProfileParameter(profileId, parameter.id)
			.then(
				function () {
					$scope.refresh();
				}
			);
	};

	$scope.profiles = profiles;
	$scope.parameter = parameter;

	$scope.editProfile = function (id) {
		locationUtils.navigateToPath('/profiles/' + id);
	};

	$scope.createProfile = function () {
		locationUtils.navigateToPath('/profiles/new');
	};

	$scope.importProfile = function () {
		const params = {
			title: 'Import Profile',
			message: "Drop Profile Here"
		};
		const modalInstance = $uibModal.open({
			templateUrl: 'common/modules/dialog/import/dialog.import.tpl.html',
			controller: 'DialogImportController',
			size: 'lg',
			resolve: {
				params: function () {
					return params;
				}
			}
		});
		modalInstance.result.then(function (importJSON) {
			profileService.importProfile(importJSON);
		}, function () {
			// do nothing
		});
	};

	$scope.compareProfiles = function () {
		const params = {
			title: 'Compare Profiles',
			message: 'Please select 2 profiles to compare',
			labelFunction: function (item) {
				return item['name'] + ' (' + item['type'] + ')'
			}
		};
		const modalInstance = $uibModal.open({
			templateUrl: 'common/modules/dialog/compare/dialog.compare.tpl.html',
			controller: 'DialogCompareController',
			size: 'md',
			resolve: {
				params: function () {
					return params;
				},
				collection: function (profileService) {
					return profileService.getProfiles({orderby: 'name'});
				}
			}
		});
		modalInstance.result.then(function (profiles) {
			$location.path($location.path() + '/' + profiles[0].id + '/' + profiles[1].id + '/compare/diff');
		}, function () {
			// do nothing
		});
	};

	$scope.refresh = function () {
		$state.reload(); // reloads all the resolves for the view
	};

	$scope.navigateToPath = (path, unsavedChanges) => locationUtils.navigateToPath(path, unsavedChanges);

	const confirmDelete = function (profile) {
		const params = {
			title: 'Delete Profile: ' + profile.name,
			key: profile.name
		};
		const modalInstance = $uibModal.open({
			templateUrl: 'common/modules/dialog/delete/dialog.delete.tpl.html',
			controller: 'DialogDeleteController',
			size: 'md',
			resolve: {
				params: function () {
					return params;
				}
			}
		});
		modalInstance.result.then(function () {
			deleteProfile(profile);
		}, function () {
			// do nothing
		});
	};


	const cloneProfile = function (profile) {
		const params = {
			title: 'Clone Profile',
			message: "You're about to clone the " + profile.name + " profile. Your clone will have the same attributes and parameter assignments as the " + profile.name + " profile.<br><br>Please enter a name for your cloned profile."
		};
		const modalInstance = $uibModal.open({
			templateUrl: 'common/modules/dialog/input/dialog.input.tpl.html',
			controller: 'DialogInputController',
			size: 'md',
			resolve: {
				params: function () {
					return params;
				}
			}
		});
		modalInstance.result.then(function (clonedProfileName) {
			profileService.cloneProfile(profile.name, clonedProfileName);
		}, function () {
			// do nothing
		});
	};

	const exportProfile = function (profile) {
		profileService.exportProfile(profile.id).then(
			function (result) {
				fileUtils.exportJSON(result, profile.name, 'traffic_ops');
			}
		);

	};

	// adds some items to the base profiles context menu
	$scope.contextMenuItems = [
		{
			text: 'Open in New Tab',
			click: function ($itemScope) {
				$window.open('/#!/profiles/' + $itemScope.p.id, '_blank');
			}
		},
		null, // Divider
		{
			text: 'Unlink Profile from Parameter',
			hasBottomDivider: function () {
				return true;
			},
			click: function ($itemScope) {
				$scope.confirmRemoveProfile($itemScope.p);
			}
		},
		{
			text: 'Edit',
			click: function ($itemScope) {
				$scope.editProfile($itemScope.p.id);
			}
		},
		{
			text: 'Delete',
			click: function ($itemScope) {
				confirmDelete($itemScope.p);
			}
		},
		null, // Divider
		{
			text: 'Clone Profile',
			click: function ($itemScope) {
				cloneProfile($itemScope.p);
			}
		},
		{
			text: 'Export Profile',
			click: function ($itemScope) {
				exportProfile($itemScope.p);
			}
		},
		null, // Divider
		{
			text: 'Manage Parameters',
			click: function ($itemScope) {
				locationUtils.navigateToPath('/profiles/' + $itemScope.p.id + '/parameters');
			}
		},
		{
			text: 'Manage Servers',
			click: function ($itemScope) {
				locationUtils.navigateToPath('/profiles/' + $itemScope.p.id + '/servers');
			}
		},
	];

	$scope.confirmRemoveProfile = async function (profile, $event) {
		if ($event) {
			$event.stopPropagation();
		}
		const params = {
			message: `The ${profile.name} profile is used by `,
			title: "Remove Parameter from Profile?",
		};
		if (profile.type === 'DS_PROFILE') { // if this is a ds profile, then it is used by delivery service(s) so we'll fetch the ds count...
			const result = await deliveryServiceService.getDeliveryServices({profile: profile.id});
			params.message += `${result.length} delivery service(s).`;
		} else { // otherwise the profile is used by servers so we'll fetch the server count...
			const result = await serverService.getServers({profileName: profile.name});
			params.message += `${result.length} server(s).`;
		}
		const modalInstance = $uibModal.open({
			templateUrl: 'common/modules/dialog/confirm/dialog.confirm.tpl.html',
			controller: 'DialogConfirmController',
			size: 'md',
			resolve: {
				params: () => params
			}
		});
		try {
			await modalInstance.result
			await removeProfile(profile.id);
		} catch {
			// modalInstances will throw if the user cancels the action.
			// it's not an actual error, so we don't need to actually handle it.
		}
	};

	$scope.selectProfiles = function () {
		const modalInstance = $uibModal.open({
			templateUrl: 'common/modules/table/parameterProfiles/table.paramProfilesUnassigned.tpl.html',
			controller: 'TableParamProfilesUnassignedController',
			size: 'lg',
			resolve: {
				parameter: function () {
					return parameter;
				},
				allProfiles: function (profileService) {
					return profileService.getProfiles({orderby: 'name'});
				},
				assignedProfiles: function () {
					return profiles;
				}
			}
		});
		modalInstance.result.then(function (selectedProfileIds) {
			const params = {
				title: 'Assign profiles to ' + parameter.name,
				message: 'Are you sure you want to modify the profiles assigned to ' + parameter.name + '?'
			};
			const modalInstance = $uibModal.open({
				templateUrl: 'common/modules/dialog/confirm/dialog.confirm.tpl.html',
				controller: 'DialogConfirmController',
				size: 'md',
				resolve: {
					params: function () {
						return params;
					}
				}
			});
			modalInstance.result.then(function () {
				profileParameterService.linkParamProfiles(parameter.id, selectedProfileIds)
					.then(
						function () {
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

	$scope.toggleVisibility = function (colName) {
		const col = parameterProfilesTable.column(colName + ':name');
		col.visible(!col.visible());
		parameterProfilesTable.rows().invalidate().draw();
	};

	$scope.columnFilterFn = function (column) {
		return column.name !== 'Action';

	};

	angular.element(document).ready(function () {
		// JQuery Datatables typings aren't in the project, because they
		// probably don't exist, but even if they do we should just get rid of
		// them instead of making them better.
		// @ts-ignore
		parameterProfilesTable = $("#parameterProfilesTable").DataTable({
			"aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
			"iDisplayLength": 25,
			"columnDefs": [
				{'orderable': false, 'targets': 5}
			],
			"aaSorting": [],
			"columns": [
				{"name": "Name", "visible": true, "searchable": true},
				{"name": "Type", "visible": true, "searchable": true},
				{"name": "Routing Disabled", "visible": true, "searchable": true},
				{"name": "Description", "visible": true, "searchable": true},
				{"name": "CDN", "visible": true, "searchable": true},
				{"name": "Action", "visible": true, "searchable": false}
			],
			"initComplete": function (settings, json) {
				try {
					// need to create the show/hide column checkboxes and bind to the current visibility
					$scope.columns = JSON.parse(localStorage.getItem("DataTables_parameterProfilesTable_/") ?? "null").columns;
				} catch (e) {
					console.error("Failure to retrieve required column info from localStorage (key=DataTables_parameterProfilesTable_/):", e);
				}
			}
		});
	});
};

TableParameterProfilesController.$inject = ["parameter", "profiles", "$scope", "$state", "$uibModal", "$window", "$location", "locationUtils", "deliveryServiceService", "profileParameterService", "serverService", "profileService", "messageModel", "fileUtils"];
module.exports = TableParameterProfilesController;
