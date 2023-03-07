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

/** @typedef {import("jquery")} */

/**
 * This is the controller for the table that lists the Parameters used by a
 * Profile.
 *
 * @param {{id: number; name: string; type: string}} profile
 * @param {unknown[]} parameters
 * @param {import("angular").IControllerService} $controller
 * @param {*} $scope
 * @param {import("../../../service/utils/angular.ui.bootstrap").IModalService} $uibModal
 * @param {import("../../../service/utils/LocationUtils")} locationUtils
 * @param {import("../../../api/DeliveryServiceService")} deliveryServiceService
 * @param {import("../../../api/ProfileParameterService")} profileParameterService
 * @param {import("../../../api/ServerService")} serverService
 * @param {import("../../../models/MessageModel")} messageModel
 */
var TableProfileParametersController = function(profile, parameters, $controller, $scope, $uibModal, locationUtils, deliveryServiceService, profileParameterService, serverService, messageModel) {

	// extends the TableParametersController to inherit common methods
	angular.extend(this, $controller('TableParametersController', { parameters: parameters, $scope: $scope }));

	let profileParametersTable;

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

	$scope.confirmRemoveParam = async function(parameter, $event) {
		if ($event) {
			$event.stopPropagation(); // this kills the click event so it doesn't trigger anything else
		}

		const params = {
			message: `The ${profile.name} profile is used by `,
			title: "Remove Parameter from Profile?"
		};
		if (profile.type == 'DS_PROFILE') { // if this is a ds profile, then it is used by delivery service(s) so we'll fetch the ds count...
			const result = await deliveryServiceService.getDeliveryServices({ profile: profile.id });
			params.message += `${result.length} delivery service(s). Are you sure you want to remove the ${parameter.name} parameter from this profile?`
		} else { // otherwise the profile is used by servers so we'll fetch the server count...
			const result = await serverService.getServers({ profileName: profile.name });
			params.message += `${result.length} server(s). Are you sure you want to remove the ${parameter.name} parameter from this profile?`
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
			await modalInstance.result;
			await removeParameter(parameter.id);
		} catch {
			// modalInstances will throw if the user cancels the action.
			// it's not an actual error, so we don't need to actually handle it.
		}
	};

	$scope.selectParams = async function() {
		const modalInstance = $uibModal.open({
			templateUrl: 'common/modules/table/profileParameters/table.profileParamsUnassigned.tpl.html',
			controller: 'TableProfileParamsUnassignedController',
			size: 'lg',
			resolve: {
				allParams: parameterService => parameterService.getParameters(),
				assignedParams: () => parameters,
				profile: () => profile
			},
		});
		const selectedParamIds = await modalInstance.result;
		const params = {
			message: `The ${profile.name} profile is used by `,
			title: `Modify ${profile.name} parameters`
		}
		if (profile.type == 'DS_PROFILE') { // if this is a ds profile, then it is used by delivery service(s) so we'll fetch the ds count...
			const result = await deliveryServiceService.getDeliveryServices({ profile: profile.id });
			params.message += `${result.length} delivery service(s). Are you sure you want to modify the parameters?`;
		} else { // otherwise the profile is used by servers so we'll fetch the server count...
			const result = await serverService.getServers({ profileName: profile.name });
			params.message += `${result.length} server(s). Are you sure you want to modify the parameters?`
		}
		const confirmModal = $uibModal.open({
			templateUrl: 'common/modules/dialog/confirm/dialog.confirm.tpl.html',
			controller: 'DialogConfirmController',
			size: 'md',
			resolve: {
				params: () => params
			}
		});
		await confirmModal.result;
		await linkProfileParameters(selectedParamIds);
	};

	$scope.toggleVisibility = function(colName) {
		const col = profileParametersTable.column(colName + ':name');
		col.visible(!col.visible());
		profileParametersTable.rows().invalidate().draw();
	};

	$scope.navigateToPath = (path, unsavedChanges) => locationUtils.navigateToPath(path, unsavedChanges);

	angular.element(document).ready(function () {
		// Datatables should be replaced, so their typings aren't included.
		// @ts-ignore
		profileParametersTable = $('#profileParametersTable').DataTable({
			"aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
			"iDisplayLength": 25,
			"aaSorting": [],
			"columnDefs": [
				{ "width": "50%", "targets": 2 },
				{ "orderable": false, "targets": 4 }
			],
			"columns": [
				{ "name": "Name", "visible": true, "searchable": true },
				{ "name": "Config File", "visible": true, "searchable": true },
				{ "name": "Value", "visible": true, "searchable": true },
				{ "name": "Secure", "visible": true, "searchable": true },
				{ "name": "Action", "visible": true, "searchable": false }
			],
			"initComplete": function(settings, json) {
				try {
					// need to create the show/hide column checkboxes and bind to the current visibility
					$scope.columns = JSON.parse(localStorage.getItem("DataTables_profileParametersTable_/") ?? "null").columns;
				} catch (e) {
					console.error("Failure to retrieve required column info from localStorage (key=DataTables_profileParametersTable_/):", e);
				}
			}
		});
	});

};

TableProfileParametersController.$inject = ['profile', 'parameters', '$controller', '$scope', '$uibModal', 'locationUtils', 'deliveryServiceService', 'profileParameterService', 'serverService', 'messageModel'];
module.exports = TableProfileParametersController;
